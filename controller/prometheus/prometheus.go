package prometheus

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"net/http"
	"os"
)

type PromeSql struct {
	SQL string `json:"sql"`
}

func Prometheus(c *gin.Context) {

	cfg, err := ini.Load("controller/config/config.ini")
	if err != nil {
		fmt.Println("mysql配置文件读取失败", err)
		os.Exit(1)
	}
	vm_host := cfg.Section("victoriaMetrics").Key("host").String()
	vm_port := cfg.Section("victoriaMetrics").Key("port").String()

	vmAddr := fmt.Sprintf("http://%s:%s/api/v1/query", vm_host, vm_port)
	fmt.Println(vmAddr)

	var promeSql PromeSql
	if err := c.ShouldBindJSON(&promeSql); err != nil {
		fmt.Println("绑定JSON失败", err)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s?query=%s", vmAddr, promeSql.SQL))

	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return

	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		fmt.Println("解析 JSON 数据时出错:", err)
		return
	}

	// 提取并打印 value 下的内容
	if data, ok := result["data"].(map[string]interface{}); ok {
		if results, ok := data["result"].([]interface{}); ok {
			for _, r := range results {
				if value, ok := r.(map[string]interface{})["value"].([]interface{}); ok {
					fmt.Println("Value:", value[1])
				}
			}
		}
	}

	//	fmt.Println(string(body))

	c.JSON(200, gin.H{
		"message": "Prometheus",
		"vmAddr":  vmAddr,
		"res":     string(body),
	})
}
