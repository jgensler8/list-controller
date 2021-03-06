package controller

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	common "github.com/jgensler8/list-controller/common/v1"
)

// Watcher is an example of watching on resource create/update/delete events
type ManagedListController struct {
	ManagedListClient *rest.RESTClient
	ManagedListScheme *runtime.Scheme
}

// Run starts an Example resource controller
func (c *ManagedListController) Run(ctx context.Context) error {
	fmt.Print("Watch Example objects\n")

	// Watch Example objects
	_, err := c.watchExamples(ctx)
	if err != nil {
		fmt.Printf("Failed to register watch for Example resource: %v\n", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (c *ManagedListController) watchExamples(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.ManagedListClient,
		common.ManagedListPlural,
		apiv1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,

		// The object type.
		//&crv1.Example{},
		&common.ManagedList{},

		// resyncPeriod
		// Every resyncPeriod, all resources in the cache will retrigger events.
		// Set to 0 to disable the resync.
		0,

		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		})

	go controller.Run(ctx.Done())
	return controller, nil
}

func (c *ManagedListController) onAdd(obj interface{}) {
	example := obj.(*common.ManagedList)
	fmt.Printf("[CONTROLLER] OnAdd %s\n", example.ObjectMeta.SelfLink)

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use exampleScheme.Copy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	copyObj, err := c.ManagedListScheme.Copy(example)
	if err != nil {
		fmt.Printf("ERROR creating a deep copy of example object: %v\n", err)
		return
	}

	exampleCopy := copyObj.(*common.ManagedList)
	exampleCopy.Status = common.ManagedListStatus{
		State:   common.ManagedListStateProcessed,
		Message: "Successfully processed by controller",
	}

	//for _, sublist := range exampleCopy.Spec.List {
	//	for _, item := range sublist.Items {
	//		fmt.Printf("[CONTROLLER] OnDelete %s", item.Object.GetObjectKind().GroupVersionKind().String())
	//	}
	//}

	err = c.ManagedListClient.Put().
		Name(example.ObjectMeta.Name).
		Namespace(example.ObjectMeta.Namespace).
		Resource(common.ManagedListPlural).
		Body(exampleCopy).
		Do().
		Error()

	if err != nil {
		fmt.Printf("ERROR updating status: %v\n", err)
	} else {
		fmt.Printf("UPDATED status: %#v\n", exampleCopy)
	}
}

func (c *ManagedListController) onUpdate(oldObj, newObj interface{}) {
	oldExample := oldObj.(*common.ManagedList)
	newExample := newObj.(*common.ManagedList)
	fmt.Printf("[CONTROLLER] OnUpdate oldObj: %s\n", oldExample.ObjectMeta.SelfLink)
	fmt.Printf("[CONTROLLER] OnUpdate newObj: %s\n", newExample.ObjectMeta.SelfLink)
}

func (c *ManagedListController) onDelete(obj interface{}) {
	example := obj.(*common.ManagedList)
	fmt.Printf("[CONTROLLER] OnDelete %s\n", example.ObjectMeta.SelfLink)
}