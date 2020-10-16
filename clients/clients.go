package clients

import (
testResource "github.com/nitin.github.io/pkg/client/clientset/versioned"
"k8s.io/client-go/kubernetes"
"k8s.io/client-go/rest"
)

type Interface interface {
	GetTestResourceClient() testResource.Interface
	GetKubernetesClient() kubernetes.Interface
}

type client struct {
	testResourceClientSet    *testResource.Clientset
	kubernetesClientSet     *kubernetes.Clientset
}

func NewClientManager(config *rest.Config) (Interface, error) {

	c := new(client)

	tc, err := testResource.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	c.testResourceClientSet = tc

	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	c.kubernetesClientSet = k

	return c, nil
}

func (c *client) GetTestResourceClient() testResource.Interface {
	return c.testResourceClientSet
}

func (c *client) GetKubernetesClient() kubernetes.Interface {
	return c.kubernetesClientSet
}
