package service

import (
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"wuchenyanghaoshuai/trident/controller/k8s/public"
)

type Ingress struct {
	Name    string
	Class   any
	Host    string
	Address string
	Ports   any
	Age     string
}

func ListIngress(c *gin.Context) {

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

	ingress, err := clientset.NetworkingV1().Ingresses(namespace).List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	var ingressList []Ingress
	for _, item := range ingress.Items {
		var ingressInfo Ingress
		ingressInfo.Name = item.Name
		ingressInfo.Class = item.Spec.IngressClassName
		ingressInfo.Host = item.Spec.Rules[0].Host
		//		ingressInfo.Address = item.Status.LoadBalancer.Ingress[0].IP
		ingressInfo.Ports = item.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Number
		// 计算svc的创建时间
		ingressInfo.Age, _ = public.CalculateDays(item.CreationTimestamp.Format("2006-01-02 15:04:05"))
		ingressList = append(ingressList, ingressInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": ingressList,
	})
}

func DelIngress(c *gin.Context) {
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
	err = clientset.NetworkingV1().Ingresses(namespace).Delete(c, parms.IngressName, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Ingress: " + parms.IngressName + " delete success",
	})
}
