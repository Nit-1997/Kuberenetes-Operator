package controller

import (
	"k8s.io/client-go/tools/cache"
	"log"

	testresource "github.com/nitin.github.io/pkg/apis/testResource/v1beta1"
)

func onAdd(obj interface{}) string {

	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Println("error while fetching the obj", err)
	}
	//logic for add wil go here
	log.Println("key for onAdd event is ",key)
	return key
}

func onUpdate(oldObj, newObj interface{}) string {

	newObject := newObj.(*testresource.TestResource)
	oldObject := oldObj.(*testresource.TestResource)

	if oldObject.ObjectMeta.ResourceVersion == newObject.ObjectMeta.ResourceVersion &&
		newObject.GetDeletionTimestamp().IsZero() {
		log.Println("re sync operation started")
		//	Uncomment below line to stop periodic reconciliation
		//return ""
	}

	key, err := cache.MetaNamespaceKeyFunc(newObj)
	if err != nil {
		log.Println("error while fetching the obj", err)
	}

	log.Println("updation started for key: %s", key)

	return key
}

func onDelete(obj interface{}) string {

	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err == nil {
		log.Println("deleting key: %s", key)
	}

	return key
}
