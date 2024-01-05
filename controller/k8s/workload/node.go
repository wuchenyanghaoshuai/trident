package workload

import (
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strconv"
	"strings"
	"wuchenyanghaoshuai/trident/controller/k8s/public"
)

// 查看node
type NodeInfo struct {
	Name      string
	CPU       string
	Memory    string
	Status    string
	Version   string
	IpAddress string
	Pods      string
	Age       string
}

// 下面的三个函数主要就是把内存获取的kb转换为gb
func parseResourceQuantity(quantity string) int64 {
	// 去除资源数量字符串中的 "Ki" 后缀
	trimmed := strings.TrimSuffix(quantity, "Ki")
	// 将剩余的数字部分转换为 int64
	value, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		// 处理错误，可能返回 0 或者记录日志
		return 0
	}
	return value
}
func kiBToGB(kiB int64) float64 {
	return float64(kiB) / (1024 * 1024)
}
func formatKiBToGB(kiB int64) string {
	gb := kiBToGB(kiB)
	return fmt.Sprintf("%.2f GB", gb)
}

func GetNode(c *gin.Context) {
	clientset, err := public.SetKubernetesConfig()
	if err != nil {
		panic(err)
	}
	nodes, err := clientset.CoreV1().Nodes().List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	var nodeList []NodeInfo
	for _, item := range nodes.Items {
		var nodeinfo NodeInfo
		nodeinfo.Name = item.Name
		nodeinfo.Status = string(item.Status.Conditions[4].Type)
		nodeinfo.IpAddress = item.Status.Addresses[0].Address
		nodeinfo.Version = item.Status.NodeInfo.KubeletVersion
		nodeinfo.CPU = item.Status.Capacity.Cpu().String()
		//	nodeinfo.Memory = item.Status.Capacity.Memory().String()
		nodeinfo.Memory = formatKiBToGB(parseResourceQuantity(item.Status.Capacity.Memory().String()))
		nodeinfo.Pods = item.Status.Capacity.Pods().String()
		nodeinfo.Age, _ = public.CalculateDays(item.CreationTimestamp.Format("2006-01-02 15:04:05"))
		nodeList = append(nodeList, nodeinfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": nodeList,
	})
}
