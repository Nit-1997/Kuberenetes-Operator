package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"

	"github.com/nitin.github.io/controller"

	informers "github.com/nitin.github.io/pkg/client/informers/externalversions"

	"github.com/nitin.github.io/clients"
	//"github.com/nitin.github.io/pkg/client/clientset/versioned/scheme"
	"github.com/nitin.github.io/signals"
)

func main() {

	var (
		masterURL  string
		kubeConfig string
	)

	log.Println("cluster config to be read")

	flag.StringVar(&kubeConfig, "kubeConfig", defaultKubeConfig(), "path to a k8s config")
	flag.StringVar(&masterURL, "master", "", "address of k8s API server")
	flag.Parse()

	cfg, err := rest.InClusterConfig()
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags(masterURL, kubeConfig)
		if err != nil {
			log.Println("error building k8s config: %s", err.Error())
		}
	}

	log.Println("cluster config read : ",cfg)



	client, err := clients.NewClientManager(cfg)
	if err != nil {
		log.Fatalf("error creating clients: %v", err)
		os.Exit(1)
	}

	//	TODO: Create event broadcaster

	//eventRecorder.InitRecorder(scheme.Scheme,client.GetKubernetesClient(),"kube-operator")
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	customInformerFactory := informers.NewSharedInformerFactory(client.GetTestResourceClient(), time.Minute*1)
	//informer setup
	informer := customInformerFactory.Nitin().V1beta1().TestResources()
	//controller runtime object
	ctrlObj := controller.NewController(client, informer, queue, controller.ProcItem)

	stopCh := signals.SetupSignalHandler()
	customInformerFactory.Start(stopCh)
	customInformerFactory.WaitForCacheSync(stopCh)

	if err = ctrlObj.Run(2, stopCh); err != nil {
		log.Fatalf("error running ctrlObj: %s", err.Error())
	}

}

func defaultKubeConfig() string {
	cfgFileName := os.Getenv("KUBECONFIG")
	if cfgFileName != "" {
		return cfgFileName
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get user home directory: %s", err.Error())
		return ""
	}

	return filepath.Join(home, ".kube", "config")
}
