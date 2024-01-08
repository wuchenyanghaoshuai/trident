package storage

import (
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"wuchenyanghaoshuai/trident/controller/k8s/public"
)

type Sc struct {
	Name        string `json:"name"`
	Provisioner string `json:"provisioner"`
	Age         string `json:"age"`
}

func ListSc(c *gin.Context) {
	clientset, err := public.SetKubernetesConfig()
	if err != nil {
		panic(err)
	}
	var parms public.Params
	if err := c.ShouldBindJSON(&parms); err != nil {
		return
	}
	namespace := parms.NameSpace

	res := public.IsNamespaceExists(namespace)
	if !res {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "namespace not exists",
		})
		return
	}
	sc, err := clientset.StorageV1().StorageClasses().List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	var scList []Sc
	for _, item := range sc.Items {
		var scInfo Sc
		scInfo.Name = item.Name
		scInfo.Provisioner = item.Provisioner
		scInfo.Age, _ = public.CalculateDays(item.CreationTimestamp.Format("2006-01-02 15:04:05"))
		scList = append(scList, scInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": scList,
	})
}
