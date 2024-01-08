package storage

import (
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"wuchenyanghaoshuai/trident/controller/k8s/public"
)

type Pv struct {
	Name          string `json:"name"`
	Capacity      string `json:"capacity"`
	AccessModes   string `json:"accessmodes"`
	ReclaimPolicy string `json:"reclaimpolicy"`
	Status        string `json:"status"`
	Claim         string `json:"claim"`
	StorageClass  string `json:"storageclass"`
	Age           string `json:"age"`
}

func ListPv(c *gin.Context) {

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
	pv, err := clientset.CoreV1().PersistentVolumes().List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	var pvList []Pv
	for _, item := range pv.Items {
		var pvInfo Pv
		pvInfo.Name = item.Name
		pvInfo.Capacity = item.Spec.Capacity.Storage().String()
		pvInfo.AccessModes = string(item.Spec.AccessModes[0])
		pvInfo.ReclaimPolicy = string(item.Spec.PersistentVolumeReclaimPolicy)
		pvInfo.Status = string(item.Status.Phase)
		pvInfo.Claim = item.Spec.ClaimRef.Name
		pvInfo.StorageClass = item.Spec.StorageClassName
		pvInfo.Age, _ = public.CalculateDays(item.CreationTimestamp.Format("2006-01-02 15:04:05"))
		pvList = append(pvList, pvInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": pvList,
	})
}
