package workload

import (
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"wuchenyanghaoshuai/trident/controller/k8s/public"
)

type Pod struct {
	Name    string `json:"name"`
	Ready   bool   `json:"ready"`
	Status  string `json:"status"`
	Restart int32  `json:"restart"`
	Age     string `json:"age"`
}

func ListPods(c *gin.Context) {
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
	pods, err := clientset.CoreV1().Pods(namespace).List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	var podList []Pod
	for _, item := range pods.Items {
		var podInfo Pod
		podInfo.Name = item.Name
		podInfo.Ready = item.Status.ContainerStatuses[0].Ready
		podInfo.Status = string(item.Status.Phase)
		podInfo.Restart = item.Status.ContainerStatuses[0].RestartCount

		podInfo.Age, _ = public.CalculateDays(item.CreationTimestamp.Format("2006-01-02 15:04:05"))
		podList = append(podList, podInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": podList,
	})
}

func DeletePod(c *gin.Context) {
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
	err = clientset.CoreV1().Pods(namespace).Delete(c, parms.PodName, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Pod: " + parms.PodName + " delete success",
	})
}
