package controller

import (
	//"fmt"
	//"log"
	//"time"

	"github.com/nitin.github.io/clients"
	testresourcelister "github.com/nitin.github.io/pkg/client/listers/testResource/v1beta1"
	//utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	//"k8s.io/client-go/tools/cache"

)

type ProcessItem func(string, clients.Interface, testresourcelister.TestResourceLister) error

// ProcItem contains logic to reconcile resource, needs to be idempotent
// returns nil for success or an error string, returning nil will stop re-processing
func ProcItem(key string, clt clients.Interface, tl testresourcelister.TestResourceLister) error {
    //controller  business logic goes here

	return nil
}
