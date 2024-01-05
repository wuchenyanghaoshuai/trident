package workload

import (
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
	"wuchenyanghaoshuai/trident/controller/k8s/public"
)

type StatefulSet struct {
	Name      string
	Replicas  int32
	Available int32
	Image     string
	Age       string
}

func ListSts(c *gin.Context) {
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

	sts, err := clientset.AppsV1().StatefulSets(namespace).List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	var stsList []StatefulSet
	for _, item := range sts.Items {
		var stsInfo StatefulSet
		stsInfo.Name = item.Name
		stsInfo.Replicas = *item.Spec.Replicas
		stsInfo.Available = item.Status.AvailableReplicas
		stsInfo.Image = item.Spec.Template.Spec.Containers[0].Image
		stsInfo.Age, _ = public.CalculateDays(item.CreationTimestamp.Format("2006-01-02 15:04:05"))
		stsList = append(stsList, stsInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": stsList,
	})
}

// k8sapi 原生不支持重启sts，只能通过更新annotations来触发重启
func RestartSts(c *gin.Context) {

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

	ststName := parms.StsName
	sts, err := clientset.AppsV1().StatefulSets(namespace).Get(c, ststName, metav1.GetOptions{})

	if err != nil {
		panic(err.Error())
	}
	// 更新 Deployment 的 annotations 来触发重启
	if sts.Spec.Template.Annotations == nil {
		sts.Spec.Template.Annotations = make(map[string]string)
	}
	sts.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	// 应用更新
	_, err = clientset.AppsV1().StatefulSets(namespace).Update(c, sts, metav1.UpdateOptions{})
	if err != nil {
		panic(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "statefulset: " + ststName + " restart success",
	})
}

func DeleteSts(c *gin.Context) {
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
	err = clientset.AppsV1().StatefulSets(namespace).Delete(c, parms.StsName, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "statefulset: " + parms.StsName + " delete success",
	})
}
