package workload

import (
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"time"
	"wuchenyanghaoshuai/trident/controller/k8s/public"
)

type Deployment struct {
	Name          string
	Replicas      int32
	Available     int32
	Image         string
	CpuRequest    string
	CpuLimit      string
	MemoryLimit   string
	MemoryRequest string
}

func ListDeploy(c *gin.Context) {
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

	deploy, err := clientset.AppsV1().Deployments(namespace).List(c, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	var deployList []Deployment
	for _, item := range deploy.Items {
		var deployInfo Deployment
		deployInfo.Name = item.Name
		deployInfo.Replicas = *item.Spec.Replicas
		deployInfo.Available = item.Status.AvailableReplicas
		deployInfo.Image = item.Spec.Template.Spec.Containers[0].Image
		deployInfo.CpuRequest = item.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String()
		deployInfo.CpuLimit = item.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()
		deployInfo.MemoryRequest = item.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String()
		deployInfo.MemoryLimit = item.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()
		deployList = append(deployList, deployInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": deployList,
	})
}

// k8sapi 原生不支持重启deployment，只能通过更新annotations来触发重启
func RestartDeploy(c *gin.Context) {
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

	deploymentName := parms.DeployName
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(c, deploymentName, metav1.GetOptions{})

	if err != nil {
		panic(err.Error())
	}
	// 更新 Deployment 的 annotations 来触发重启
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	// 应用更新
	_, err = clientset.AppsV1().Deployments(namespace).Update(c, deployment, metav1.UpdateOptions{})
	if err != nil {
		panic(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "deployment: " + deploymentName + " restart success",
	})
}

func DeleteDeploy(c *gin.Context) {
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
	err = clientset.AppsV1().Deployments(namespace).Delete(c, parms.DeployName, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "deployment: " + parms.DeployName + " delete success",
	})
}
