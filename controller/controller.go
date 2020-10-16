package controller

import (
	testresourcescheme "github.com/nitin.github.io/pkg/client/clientset/versioned/scheme"
	testresourceinformer "github.com/nitin.github.io/pkg/client/informers/externalversions/testResource/v1beta1"
	testresourcelister "github.com/nitin.github.io/pkg/client/listers/testResource/v1beta1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/nitin.github.io/clients"
	"log"
)

type Controller struct {
	client                  clients.Interface
	testResourceLister      testresourcelister.TestResourceLister
	testResourceSynced      cache.InformerSynced
	queue                   workqueue.RateLimitingInterface
	testResourceInformer    testresourceinformer.TestResourceInformer
	processItem             ProcessItem
}

func NewController(cl clients.Interface, trInf testresourceinformer.TestResourceInformer, q workqueue.RateLimitingInterface,
	pItem ProcessItem) *Controller {

	utilruntime.Must(testresourcescheme.AddToScheme(scheme.Scheme))

	c := &Controller{
		client:                  cl,
		testResourceLister:      trInf.Lister(),
		testResourceSynced:      trInf.Informer().HasSynced,
		queue:                   q,
		testResourceInformer:    trInf,
		processItem:             pItem,
	}

	trInf.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Println("Event add")
			if key := onAdd(obj); key != "" {
				c.queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Println("Event update")
			if key := onUpdate(oldObj, newObj); key != "" {
				c.queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			log.Println("Event delete")
			if key := onDelete(obj); key != "" {
				c.queue.Forget(key)
			}
		},
	})

	return c
}
