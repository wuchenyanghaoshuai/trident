package public

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

// 公共部分，引入k8s.io/client-go/kubernetes包，用于连接k8s集群
func SetKubernetesConfig() (*kubernetes.Clientset, error) {
	kubeconfig := "controller/config/kube_config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// 这块主要是用于删除的时候使用的
type Params struct {
	NameSpace   string `json:"namespace"`
	PodName     string `json:"podname"`
	DeployName  string `json:"deployname"`
	StsName     string `json:"stsname"`
	IngressName string `json:"ingressname"`
	SvcName     string `json:"svcname"`
}

// 判断传入的时间戳到现在是多久
func CalculateDays(timestamp string) (string, error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	const layout = "2006-01-02 15:04:05"
	// 解析时间字符串为 time.Time 对象
	parsedTime, err := time.ParseInLocation(layout, timestamp, loc)
	if err != nil {
		return "", err
	}

	// 获取当前时间，考虑时区
	now := time.Now().In(loc)

	// 确保给定的时间戳不是未来时间
	if parsedTime.After(now) {
		return "", fmt.Errorf("the given timestamp is in the future")
	}

	// 计算时间差
	duration := now.Sub(parsedTime)

	// 根据时间差的大小决定输出格式
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	var timeString string
	if days > 0 {
		timeString = fmt.Sprintf("%dd", days)
	} else if hours > 0 {
		timeString = fmt.Sprintf("%dh%dm", hours, minutes)
	} else if minutes > 0 {
		timeString = fmt.Sprintf("%dm", minutes)
	} else {
		timeString = fmt.Sprintf("%ds", seconds)
	}

	return timeString, nil

}
