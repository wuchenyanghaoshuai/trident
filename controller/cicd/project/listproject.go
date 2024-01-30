package project

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"wuchenyanghaoshuai/trident/controller/dao/mysql"
)

type RespProject struct {
	//	Id          int    `json:"id"`
	ProjectName   string   `json:"project_name"`
	ProjectID     int      `json:"project_id"`
	ProjectBranch []string `json:"project_branch"`
}

// 这块返回的就是projectname,projectid,projectbranch
/*
{
    "message": "成功列出所有GitLab项目",
    "projects": [
        {
            "project_name": "wx-go",
            "project_id": 90,
            "project_branch": [
                "dev",
                "main",
                "master",
                "prod",
                "test",
                "wuchenyangtest"
            ]
        }
}
*/
func ListProject(c *gin.Context) {
	//查询所有项目
	var projects []GitlabProject
	err := mysql.DB.WithContext(c).Find(&projects)
	if err.Error != nil {
		c.JSON(500, gin.H{
			"error": err.Error,
		})
		return
	}

	//返所有的项目名称
	var respProject []RespProject

	for _, project := range projects {

		branch, _ := GetBranches(project.ProjectID)
		respProject = append(respProject, RespProject{
			//			Id:          project.Id,
			ProjectID:     project.ProjectID,
			ProjectName:   project.ProjectName,
			ProjectBranch: branch,
		})
	}

	c.JSON(200, gin.H{
		"message":  "成功列出所有GitLab项目",
		"projects": respProject,
	})
}

type Branch struct {
	Name string `json:"name"`
}

func GetBranches(projectid int) ([]string, error) {
	// Prepare the request
	gitlabURL := fmt.Sprintf("http://192.168.3.101/api/v4/projects/%d", projectid)
	client := &http.Client{}
	req, err := http.NewRequest("GET", gitlabURL+"/repository/branches", nil)
	if err != nil {
		return nil, err
	}

	// Set the private token header
	req.Header.Add("PRIVATE-TOKEN", "BWqn7UKq3BhUvqeAvJxN")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch branches: %s", resp.Status)
	}

	// Parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var branches []Branch
	err = json.Unmarshal(body, &branches)
	if err != nil {
		return nil, err
	}

	// Extract branch names
	var branchNames []string
	for _, branch := range branches {
		branchNames = append(branchNames, branch.Name)
	}

	return branchNames, nil
}
