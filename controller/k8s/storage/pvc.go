package storage

import (
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"wuchenyanghaoshuai/trident/controller/k8s/public"
)

type Pvc struct {
	Name         string  `json:"name"`
	Status       string  `json:"status"`
	Volume       string  `json:"volume"`
	Capacity     string  `json:"capacity"`
	AccessModes  string  `json:"accessmodes"`
	StorageClass *string `json:"storageclass"`
	Age          string  `json:"age"`
}

func ListPvc(c *gin.Context) {
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
	pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	var pvcList []Pvc
	for _, item := range pvc.Items {
		var pvcInfo Pvc
		pvcInfo.Name = item.Name
		pvcInfo.Status = string(item.Status.Phase)
		pvcInfo.Volume = item.Spec.VolumeName
		pvcInfo.Capacity = item.Spec.Resources.Requests.Storage().String()
		pvcInfo.AccessModes = string(item.Spec.AccessModes[0])
		pvcInfo.StorageClass = item.Spec.StorageClassName
		pvcInfo.Age, _ = public.CalculateDays(item.CreationTimestamp.Format("2006-01-02 15:04:05"))
		pvcList = append(pvcList, pvcInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": pvcList,
	})
}
