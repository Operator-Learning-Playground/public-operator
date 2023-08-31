package init_controller

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

// K8sRestConfigInPod
func K8sRestConfigInPod() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	return config
}

// K8sRestConfig 初始化rest.Config
func K8sRestConfig() *rest.Config {
	if os.Getenv("release") == "1" { //自定义环境
		log.Println("run in cluster")
		return K8sRestConfigInPod()
	}
	log.Println("run outside cluster")
	config, err := clientcmd.BuildConfigFromFlags("", "./config1")
	if err != nil {
		log.Fatal(err)
	}
	config.Insecure = true
	return config
}
