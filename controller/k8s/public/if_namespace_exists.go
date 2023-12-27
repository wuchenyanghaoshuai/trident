package public

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IsNamespaceExists(namespace string) bool {
	clientset, err := SetKubernetesConfig()
	if err != nil {
		panic(err)
	}
	namespaceList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	found := false
	for _, item := range namespaceList.Items {
		if namespace == item.Name {
			found = true
			break
		}
	}
	return found
}
