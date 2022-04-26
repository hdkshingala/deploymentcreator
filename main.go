package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	dcClient "github.com/hdkshingala/deploymentcreator/pkg/client/clientset/versioned"
	dcFactory "github.com/hdkshingala/deploymentcreator/pkg/client/informers/externalversions"
	dcController "github.com/hdkshingala/deploymentcreator/pkg/controller"
)

func main() {
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Printf("Building config from flags, %s", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("Building config from InClusterConfig, %s", err.Error())
			return
		}
	}

	dcClientSet, err := dcClient.NewForConfig(config)
	if err != nil {
		log.Printf("Getting klient set, %s", err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Getting client set, %s", err.Error())
	}

	dcInformerFactory := dcFactory.NewSharedInformerFactory(dcClientSet, 10*time.Minute)

	k := dcController.NewController(clientSet, dcClientSet, dcInformerFactory.Hardik().V1alpha1().DeploymentCreators())

	ch := make(chan struct{})

	dcInformerFactory.Start(ch)

	if err = k.Run(ch); err != nil {
		log.Printf("Error running controller, %s", err.Error())
	}

}
