package cicd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os/exec"
	"strings"
)

// 如果要buildimage的话首先第一步先拉代码
// 然后再 切到正确的分之
// 然后根据Dockerfile来buildimage

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

	data, err := BuildImage(build.JobName, build.ChangeType, build.JobUrl, build.JobBranch)
	if err != nil {
		fmt.Println("err:", "buildimage的时候出现了错误", err)
		return
	}
	c.JSON(200, gin.H{
		"message": "cicd",
		"data":    data,
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
		pwd
	    git clone ` + joburl + `
		cd ` + jobname + `
		pwd
		git checkout ` + jobbranch + `
        commitid=$(git rev-parse HEAD | cut -c 1-6)
        docker build -t ` + harborurl + jobname + `:` + changetype + `-` + "${commitid}" + ` .
		docker push  ` + harborurl + jobname + `:` + changetype + `-` + "${commitid}" + `  
		cd .. 
		pwd
        rm -rf ` + jobname + `
		echo ` + harborurl + jobname + `:` + changetype + `-` + "${commitid}" + `
	`
	cmd := exec.Command("bash", "-c", script)

	output, err := cmd.Output()

	if err != nil {
		fmt.Println("err111:", err)
		return "", err
	}
	data := string(output)
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.Contains(line, "192.168.3.103:8888/lieyunzjk") {
			data = line
			break
		}
	}
	return data, err
}

// 获取项目到底多少个分支
