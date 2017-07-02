# List Controller

## Use Case

Your organization has a recommended way of deploying applications on Kubernetes.
This might include creating an initial set of resource for a particular container.
Some of these resources could include Ingress or Service objects based on variables supplied by a build process.
In the end, this process requirements lifecycle management.
The resources that should be managed are subject to change so your organization will need a controller that can manage an arbitrary set of resources.

Ideally, another Controller should be written to manage the ListResources.
That Controller should manage the health of the List.

## How It Works

* Deploy the List Controller to a particular namespace.
  ```
  kubectl apply -f examples/list-controller.yaml
  ```
* Upload a CustomResourceDefinition containing a List of Resource
  ```
  kubectl apply -f examples/example-list.yaml
  ```

# Thank You

* Initially, I started with Aaron Levy's [kube-controller-demo](https://github.com/aaronlevy/kube-controller-demo). This demo uses ThirdPartyResources so I decided to switch to CustomResourceDefinitions.
* Check out the CustomResourceDefinition example located [here](https://github.com/kubernetes/apiextensions-apiserver). 99% of the initial commit of this project is from that example.
* All of the other Contributors to Kubernetes and Go