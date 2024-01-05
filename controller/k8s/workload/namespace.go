package workload

import (
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"wuchenyanghaoshuai/trident/controller/k8s/public"
)

// 查看和删除namespace
type Namespace struct {
	NameSpace string
	Status    string
	Age       string
}

func ListNamespace(c *gin.Context) {

	clientset, err := public.SetKubernetesConfig()
	if err != nil {
		panic(err)
	}
	var nsList []Namespace
	namespaceList, err := clientset.CoreV1().Namespaces().List(c, metav1.ListOptions{})

	for _, item := range namespaceList.Items {
		var namespace Namespace
		namespace.NameSpace = item.Name
		namespace.Status = string(item.Status.Phase)
		namespace.Age, _ = public.CalculateDays(item.CreationTimestamp.Format("2006-01-02 15:04:05"))
		nsList = append(nsList, namespace)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": nsList,
	})
}

func DeleteNamespace(c *gin.Context) {
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

	err = clientset.CoreV1().Namespaces().Delete(c, namespace, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
	})
}
