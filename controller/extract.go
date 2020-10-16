package controller

import (
	"fmt"
	"time"
	"log"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) Run(threads int, stopCh <-chan struct{}) error {

	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	if !cache.WaitForCacheSync(stopCh, c.testResourceSynced) {
		log.Println("timed out waiting for caches to sync")
		return fmt.Errorf("failed to wait for caches to sync")
	}

	log.Println("Starting workers")

	for i := 0; i < threads; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	log.Println("Started workers")
	<-stopCh
	log.Println("Shutting down workers")
	return nil
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

// processNextWorkItem deals with one key off the queue.
// It returns false when it's time to quit.
func (c *Controller) processNextItem() bool {
	obj, quit := c.queue.Get()
	if quit {
		return false
	}

	defer c.queue.Done(obj)

	var key string
	var ok bool

	if key, ok = obj.(string); !ok {
		c.queue.Forget(obj)
		utilruntime.HandleError(fmt.Errorf("expected string in queue but got %#v", obj))
		return true
	}

	// processItem must return nil for stop processing, else re-processing is triggered
	err := c.processItem(key, c.client, c.testResourceLister)

	if err == nil {
		log.Println("successfully processed: ", key)
		c.queue.Forget(key)
		return true
		//	Here we check how many times we can re-process same key(job)
	} else if c.queue.NumRequeues(key) < 5 {
		c.queue.AddRateLimited(key)
		//log.Warn("re-processing: ", key)
	} else {
		c.queue.Forget(key)
		//log.Error("max retries reached while processing: ", key)
	}

	utilruntime.HandleError(fmt.Errorf("processItem failed with: %v", err))
	return true
}
