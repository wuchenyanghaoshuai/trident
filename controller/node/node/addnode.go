package node

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"wuchenyanghaoshuai/trident/controller/dao/mysql"
)

// 思路 在获取到前端传递过来的ip和用户名密码以后，模拟使用ssh去登录，如果登录成功，就说明这个ip是可以使用的,然后拷贝公钥到目标机器上
// 如果登录成功，就把这个ip和用户名存入数据库，如果登录失败，直接return
// 登录成功以后就去获取linux的版本号，以及cpu核数，内存大小，如果是centos就返回centos，如果是ubuntu就返回ubuntu
// host节点信息
/* postman调接口新增node，参数如下，其中nodtes字段可以为空相当于备注
{
    "hostname":"trident111",
    "ip":"192.168.3.102",
    "username":"root",
    "password":"centos",
    "port" :"22",
    "label":"devops",
    "notes":"这个是我创建的一个备注关于node机器"
}

*/
type HostParams struct {
	Id         int               `json:"id" gorm:"primaryKey"`
	Hostname   string            `json:"hostname"gorm:"unique"`
	Username   string            `json:"username"`
	Password   string            `json:"password" gorm:"-"`
	Port       string            `json:"port"`
	Ip         string            `json:"ip"`
	Status     bool              `json:"status"`
	Osinfo     map[string]string `json:"osinfo" gorm:"serializer:json"`
	Label      string            `json:"label"`
	Notes      string            `json:"notes"`
	PrivateKey string            `json:"private_key" gorm:"-"` //这个字段不会存入数据库
}
type RespHostParams struct {
	Hostname string            `json:"hostname"`
	Ip       string            `json:"ip"`
	Status   bool              `json:"status"`
	Osinfo   map[string]string `json:"osinfo" gorm:"serializer:json"`
}

func AddHost(c *gin.Context) {
	var hosts HostParams
	err := c.ShouldBindJSON(&hosts)
	if err != nil {
		panic(err)
	}
	if hosts.Port == "" || hosts.Label == "" || hosts.Username == "" || hosts.Password == "" || hosts.Ip == "" || hosts.Hostname == "" {
		c.JSON(200, gin.H{
			"message": "参数不能为空",
		})
		return
	}
	//给机器加一个linux的标签
	if hosts.Osinfo == nil {
		hosts.Osinfo = map[string]string{"os": "linux"}
	}

	_, osinfo := GetHostIsLinuxOrUbuntu(hosts.Ip, hosts.Port, hosts.Username, hosts.Password)
	if osinfo == "" {
		c.JSON(200, gin.H{
			"message": "获取目标机器信息失败",
		})
		return
	}
	//新增一个函数就是判断这个是否能登录到这个机器，如果可以的话就执行一下ssh-copyid 的这个操作把本机的公钥复制到目标机器上
	info := HostParams{
		Username:   hosts.Username,
		Password:   hosts.Password,
		Ip:         hosts.Ip,
		Port:       hosts.Port,
		PrivateKey: "controller/config/id_rsa",
	}
	if err := CopyID(info); err != nil {
		c.JSON(200, gin.H{
			"message": "复制公钥失败",
		})
		return
	}
	//
	hosts.Osinfo["osinfo"] = osinfo
	//这块密码设置为空是因为在存入数据库的时候，密码不需要存入数据库，密码仅仅作为验证是否能登录上目标机器的一个标准
	hosts.Password = ""
	//status字段主要是就为了判断用户名密码是否正确，如果正确的话这个位置就是true
	//在数据库里false=0 true=1 所以看到0不必惊讶
	hosts.Status = true
	mysql.DB.AutoMigrate(&hosts)
	res := mysql.DB.WithContext(c).Table("host_params").Create(&hosts)
	if res.Error != nil {
		// 如果返回的错误是因为唯一性约束违反，可以返回一个特定的错误信息
		if strings.Contains(res.Error.Error(), "Duplicate entry") {
			c.JSON(http.StatusConflict, gin.H{"error": "主机名称重复,请尝试其他名称"})
		} else {
			// 如果是其他类型的数据库错误
			c.JSON(http.StatusInternalServerError, gin.H{"error2": res.Error.Error()})
		}
		return
	}

	resphostparams := RespHostParams{
		Hostname: hosts.Hostname,
		Ip:       hosts.Ip,
		Osinfo:   hosts.Osinfo,
		Status:   hosts.Status,
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    resphostparams,
	})
}

// 获取目标机器是luinx的centos还是ubuntu,以及版本号以及几核几G
func GetHostIsLinuxOrUbuntu(ip string, port string, username string, password string) (string, string) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		fmt.Println("连接失败", err)
		return "连接失败", ""
	}
	defer conn.Close()

	// 创建一个新会话并运行命令
	runCommand := func(cmd string) (string, error) {
		session, err := conn.NewSession()
		if err != nil {
			return "", err
		}
		defer session.Close()

		output, err := session.Output(cmd)
		if err != nil {
			return "", err
		}
		return string(output), nil
	}

	//获取系统是centos
	oscmd := " cat /etc/redhat-release | awk '{print $1, $4, $5, $6}'"
	outputos, err := runCommand(oscmd)
	if err != nil {

		return "执行命令失败1", ""
	}
	outputos = strings.TrimSpace(outputos)
	//获取系统CPU核数
	cpucmd := "cat /proc/cpuinfo | grep 'cpu cores' | wc -l"
	outputcpu, err := runCommand(cpucmd)
	if err != nil {

		return "执行命令失败2", ""
	}
	outputcpu = strings.TrimSpace(outputcpu)
	//获取系统内存大小
	memcmd := "cat /proc/meminfo | grep MemTotal | awk '{print $2}'"
	outputmem, err := runCommand(memcmd)
	if err != nil {
		return "执行命令失败3", ""
	}
	mem, err := MemKBtoGBStringToInt(outputmem)
	if err != nil {
		fmt.Printf("Error converting KB to GB: %s\n", err)
		return "error", ""
	}

	return "err", fmt.Sprintf("%s  %sC  %dG ", string(outputos), string(outputcpu), mem)
}

// MemKBtoGBStringToInt 将内存从KB（字符串）转换为GB（整数），结果四舍五入
func MemKBtoGBStringToInt(kbStr string) (int, error) {
	kbStr = strings.TrimSpace(kbStr)
	kb, err := strconv.ParseFloat(kbStr, 64)
	if err != nil {
		return 0, err
	}
	gb := kb / (1024 * 1024)        // 将KB转换为GB
	return int(math.Round(gb)), nil // 使用math.Round四舍五入到最近的整数，并转换为int
}

func CopyID(info HostParams) error {

	key, err := ioutil.ReadFile(info.PrivateKey)
	if err != nil {
		return fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: info.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
			ssh.Password(info.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 注意：这是不安全的，实际使用时应该使用更安全的方法
	}

	// 连接到远程服务器
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", info.Ip, info.Port), config)
	if err != nil {
		return fmt.Errorf("unable to connect: %v", err)
	}
	defer conn.Close()

	// 获取公钥的内容
	pubKeyPath := info.PrivateKey + ".pub"
	pubKeyData, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return fmt.Errorf("unable to read public key data: %v", err)
	}

	// 创建远程.ssh目录（如果它不存在的话）
	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	mkdirCmd := "mkdir -p ~/.ssh"
	if err := session.Run(mkdirCmd); err != nil {
		return fmt.Errorf("failed to run: %s, error: %v", mkdirCmd, err)
	}

	// 将公钥复制到远程服务器的authorized_keys文件
	pubKeyCmd := fmt.Sprintf("echo '%s' >> ~/.ssh/authorized_keys", pubKeyData)
	session, err = conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	if err := session.Run(pubKeyCmd); err != nil {
		return fmt.Errorf("failed to run: %s, error: %v", pubKeyCmd, err)
	}

	fmt.Println("Public key copied successfully.")
	return nil
}
