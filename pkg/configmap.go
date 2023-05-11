package pkg

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ReadConfigMap(secretName string, clientset *kubernetes.Clientset) map[string][]byte {
	secretNamespace := "default"
	secret, err := clientset.CoreV1().Secrets(secretNamespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	return secret.Data
}
