package cicd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

// 如果要buildimage的话首先第一步先拉代码
// 然后再 切到正确的分之
// 然后根据Dockerfile来buildimage
/* postman请求方式
{
    "job_name":"wx-go",
    "change_type":"test",
    "job_url":"git@192.168.3.101:wx/wx-go.git",
    "job_branch":"wuchenyangtest"
}
*/
type Build struct {
	JobName    string `json:"job_name"`
	ChangeType string `json:"change_type"`
	JobUrl     string `json:"job_url"`
	JobBranch  string `json:"job_branch"`
}

func CICD(c *gin.Context) {
	var build Build
	if err := c.ShouldBindJSON(&build); err != nil {
		fmt.Println("err:", err)
		return
	}

	imageTag, err := BuildImage(build.JobName, build.ChangeType, build.JobUrl, build.JobBranch)
	if err != nil {
		fmt.Println("err:", "buildimage的时候出现了错误", err)
		return
	}
	c.JSON(200, gin.H{
		"message":  "buildimage成功",
		"imageTag": imageTag,
		"status":   http.StatusOK,
	})

}

func BuildImage(jobname string, changetype string, joburl string, jobbranch string) (string, error) {
	// 1. clone代码
	// 2. 切到正确的分之
	// 3. 根据Dockerfile来buildimage
	//4. demo   192.168.3.103:8888/lieyun/ailieyun-ms:test2-1.0.0-120820
	harborurl := "192.168.3.103:8888/lieyun/"
	script := `
		#!/bin/bash
        cd ./gitcodedic
	    git clone ` + joburl + `
		cd ` + jobname + `
		git checkout ` + jobbranch + `
        commitid=$(git rev-parse HEAD | cut -c 1-6)
        docker build -t ` + harborurl + jobname + `:` + changetype + `-` + "${commitid}" + ` .
		docker push  ` + harborurl + jobname + `:` + changetype + `-` + "${commitid}" + `  
		cd ..
        rm -rf ` + jobname + `
		echo ` + harborurl + jobname + `:` + changetype + `-` + "${commitid}" + ` > result.txt
	`
	cmd := exec.Command("bash", "-c", script)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("err111:", err)
		return "", err
	}
	fmt.Println(string(output))
	imageTagOutput, err := ioutil.ReadFile("/Users/wuchenyang/code/NewCodeForTrident/trident/gitcodedic/result.txt")
	if err != nil {
		fmt.Println("err222:", err)
		return "", err
	}
	imageTag := strings.TrimSpace(string(imageTagOutput))
	deletefilescript := `
		#!/bin/bash
		pwd
		rm -rf ./gitcodedic/result.txt
	`
	cmd = exec.Command("bash", "-c", deletefilescript)
	outputdel, err := cmd.Output()
	if err != nil {
		fmt.Println("err333:", err)
		return "", err
	}
	fmt.Println(string(outputdel))
	return imageTag, err
}

//df4d59f37d6b6cb75e876efc9da67849
//imagePullSecrets:
//      - name: prepullimage
//imagePullSecrets:
////- name: prepullimage
//registry.cn-zhangjiakou.aliyuncs.com/tianjinlieyun
//registry.cn-zhangjiakou.aliyuncs.com/tianjinlieyun/ailieyun-h5:pre-1.0.0-0388e0
