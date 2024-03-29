# basicKubeOperator
Basic kubernetes operator written in golang using client go library and code-generator

# References 
   - https://stackoverflow.com/questions/60094368/kubernetes-client-code-generator-can-the-code-exist-only-locally-and-not-on-a-r
   - https://insujang.github.io/2020-02-13/programming-kubernetes-crd/#get-generated-go-code-with-code-generator
   - https://trstringer.com/extending-k8s-custom-controllers/
   - https://itnext.io/how-to-create-a-kubernetes-custom-controller-using-client-go-f36a7a7536cc
# Setup 

**How to setup a basic controller using code-generator in 2020**
   
   - Pull the code generator repo into your $GOPATH
   - Make sure your directory structure resembles this
          
                 
               somedir (src in this case)
               ├── code-generator
               └── myorg.com (github.com in my case)
                   └── controller-repo (nitin.github.io in this case) 
                       ├── go.mod
                       ├── go.sum
                       └── pkg
                           └── apis
                               └── customResourceName (testResource in this case)
                                   ├── register.go
                                   └── api-group-version(v1beta1 in my case)
                                       ├── doc.go
                                       ├── register.go
                                       └── types.go
                               
  
   - Structure of doc.go for code generator
       
        
          // +k8s:deepcopy-gen=package
          // +groupName=nitin.github.io
          
          package v1beta1
        **Make sure to give the same group name as in your CRD defination
    
   - Structure of register.go for code generator 
         
           
             package v1beta1
             
             import (
             	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
             	"k8s.io/apimachinery/pkg/runtime"
             	"k8s.io/apimachinery/pkg/runtime/schema"
             )
             
             var (
             	// SchemeBuilder initializes a scheme builder
             	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
             	// AddToScheme is a global function that registers this API group & version to a scheme
             	AddToScheme = SchemeBuilder.AddToScheme
             )
             
             // SchemeGroupVersion is group version used to register these objects.
             var SchemeGroupVersion = schema.GroupVersion{
             	Group:   GroupName,
             	Version: Version,
             }
             
             func Resource(resource string) schema.GroupResource {
             	return SchemeGroupVersion.WithResource(resource).GroupResource()
             }
             
             func addKnownTypes(scheme *runtime.Scheme) error {
             	scheme.AddKnownTypes(SchemeGroupVersion,
             		&TestResource{},
             		&TestResourceList{},
             	)
             	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
             	return nil
             } 
     
   **change the Resource name as per need , this is used to add your custom resource to scheme
      
   - Structure of types.go for code generator
          
          
          
             package v1beta1
             
             import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
             
             // These const variables are used in our custom controller.
             const (
             	GroupName string = "nitin.github.io"
             	Kind      string = "TestResource"
             	Version   string = "v1beta1"
             	Plural    string = "testresources"
             	Singluar  string = "testresource"
             	ShortName string = "ts"
             	Name      string = Plural + "." + GroupName
             )
             
             // TestResourceSpec specifies the 'spec' of TestResource CRD.
             type TestResourceSpec struct {
             	Command        string `json:"command"`
             	CustomProperty string `json:"customProperty"`
             }
             
             // +genclient
             // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
             
             // TestResource describes a TestResource custom resource.
             type TestResource struct {
             	metav1.TypeMeta   `json:",inline"`
             	metav1.ObjectMeta `json:"metadata,omitempty"`
             
             	Spec   TestResourceSpec `json:"spec"`
             	Status string           `json:"status"`
             }
             
             // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
             
             // TestResourceList is a list of TestResource resources.
             type TestResourceList struct {
             	metav1.TypeMeta `json:",inline"`
             	metav1.ListMeta `json:"metadata"`
             
             	Items []TestResource `json:"items"`
             }
     
   **Make sure to give the same group name as in your CRD defination and define the structure of the schema as per your crd
   
   - Now Your code is ready to be used by code generator
      
          
          
           RUN THIS COMMAND IN TERMINAL TO GENERATE CLIENTSET/INFORMERS/LISTERS
           
           $GOPATH/src/code-generator/generate-groups.sh all github.com/nitin.github.io/pkg/client github.com/nitin.github.io/pkg/apis "testResource:v1beta1" \
           -h $GOPATH/src/code-generator/hack/boilerplate.go.txt \
           -o $GOPATH/src
   
   
   - Now you should see this structure       
     
              
              
              nitin.github.io
              ├── go.mod
              ├── go.sum
              └── pkg
                  ├── apis
                  │   └── testResource
                  │       ├── register.go
                  │       └── v1beta1
                  │           ├── doc.go
                  │           ├── register.go
                  │           ├── types.go
                  │           └── zz_generated.deepcopy.go
                  └── client
                      ├── clientset
                      │   └── versioned
                      │       ├── clientset.go
                      │       ├── doc.go
                      │       ├── fake
                      │       ├── scheme
                      │       └── typed
                      ├── informers
                      │   └── externalversions
                      │       ├── factory.go
                      │       ├── generic.go
                      │       ├── internalinterfaces
                      │       └── testResource
                      └── listers
                          └── testResource
                              └── v1beta1
      
   - Now time to setup our controller directory structure
     
     
     
                  nitin.github.io        
                  ├── go.mod
                  ├── go.sum
                  ├── pkg
                  ├── signals               #used to handle client
                  ├── clients               #contains logic to access both custom as well as kubernetes client 
                  └── controller
                      ├── controller.go     #contains controller setup
                      ├── events.go         #contains event handling
                      ├── extract.go        #queue actions are performed here
                      └── process_item.go   #contains the core logic for operations on each key 
                  
   - Final Step before we can start working on the controller is to remove the context ctx parameter from typed clientset and from informer generated by code-generator

#Wehoo!! we are good to go to write our custom controller now.                   