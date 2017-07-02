package main


import (
	"context"
	"flag"
	//"fmt"

	//apiv1 "k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	//crv1 "k8s.io/apiextensions-apiserver/examples/client-go/apis/cr/v1"
	//exampleclient "k8s.io/apiextensions-apiserver/examples/client-go/client"
	//examplecontroller "k8s.io/apiextensions-apiserver/examples/client-go/controller"
	client "github.com/jgensler8/list-controller/pkg/client"
	controller "github.com/jgensler8/list-controller/pkg/controller"
)



func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := buildConfig(*kubeconfig)
	if err != nil {
		panic(err)
	}

	apiextensionsclientset, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// initialize custom resource using a CustomResourceDefinition if it does not exist
	_, err = client.CreateCustomResourceDefinition(apiextensionsclientset)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		panic(err)
	}
	//defer apiextensionsclientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crd.Name, nil)

	// make a new config for our extension's API group, using the first config as a baseline
	exampleClient, exampleScheme, err := client.NewClient(config)
	if err != nil {
		panic(err)
	}

	// start a controller on instances of our custom resource
	controller := controller.ManagedListController{
		ManagedListClient: exampleClient,
		ManagedListScheme: exampleScheme,
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	//go controller.Run(ctx)
	controller.Run(ctx)

	// Create an instance of our custom resource
	//example := &crv1.Example{
	//	ObjectMeta: metav1.ObjectMeta{
	//		Name: "example1",
	//	},
	//	Spec: crv1.ExampleSpec{
	//		Foo: "hello",
	//		Bar: true,
	//	},
	//	Status: crv1.ExampleStatus{
	//		State:   crv1.ExampleStateCreated,
	//		Message: "Created, not processed yet",
	//	},
	//}
	//var result crv1.Example
	//err = exampleClient.Post().
	//	Resource(crv1.ExampleResourcePlural).
	//	Namespace(apiv1.NamespaceDefault).
	//	Body(example).
	//	Do().Into(&result)
	//if err == nil {
	//	fmt.Printf("CREATED: %#v\n", result)
	//} else if apierrors.IsAlreadyExists(err) {
	//	fmt.Printf("ALREADY EXISTS: %#v\n", result)
	//} else {
	//	panic(err)
	//}
	//
	//// Poll until Example object is handled by controller and gets status updated to "Processed"
	//err = exampleclient.WaitForExampleInstanceProcessed(exampleClient, "example1")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Print("PROCESSED\n")
	//
	//// Fetch a list of our TPRs
	//exampleList := crv1.ExampleList{}
	//err = exampleClient.Get().Resource(crv1.ExampleResourcePlural).Do().Into(&exampleList)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("LIST: %#v\n", exampleList)
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
