package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"wuchenyanghaoshuai/trident/controller/k8s/public"
)

type Service struct {
	Name      string
	Type      string
	ClusterIP string
	Ports     string
	Age       string
}

func ListSvc(c *gin.Context) {

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

	svc, err := clientset.CoreV1().Services(namespace).List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	var svcList []Service
	for _, item := range svc.Items {
		var svcInfo Service
		svcInfo.Name = item.Name
		svcInfo.Type = string(item.Spec.Type)
		svcInfo.ClusterIP = item.Spec.ClusterIP
		svcInfo.Ports = fmt.Sprintf("%d:%d/%s", item.Spec.Ports[0].Port, item.Spec.Ports[0].NodePort, item.Spec.Ports[0].Protocol)
		// 计算svc的创建时间
		svcInfo.Age, _ = public.CalculateDays(item.CreationTimestamp.Format("2006-01-02 15:04:05"))
		svcList = append(svcList, svcInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": svcList,
	})
}

func DelSvc(c *gin.Context) {
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

	err = clientset.CoreV1().Services(namespace).Delete(c, parms.SvcName, metav1.DeleteOptions{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "Service: " + parms.SvcName + " not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Service: " + parms.SvcName + " delete success",
	})
}
