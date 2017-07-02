package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "cr.client-go.k8s.io"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

const (
	RebootAnnotation           = "list-agent.v1.demo.local/reboot"
	RebootNeededAnnotation     = "list-agent.v1.demo.local/reboot-needed"
	RebootInProgressAnnotation = "list-agent.v1.demo.local/reboot-in-progress"
)

const ManagedListPlural = "managedlists"

type ManagedList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ManagedListSpec   `json:"spec"`
	Status            ManagedListStatus `json:"status,omitempty"`
}

type ManagedListSpec struct {
	List []v1.List `json:"List"`
}


type ManagedListStatus struct {
	State   ManagedListState `json:"state,omitempty"`
	Message string       `json:"message,omitempty"`
}

type ManagedListState string

const (
	ManagedListStateCreated   ManagedListState = "Created"
	ManagedListStateProcessed ManagedListState = "Processed"
	ManagedListStateOutOfSync ManagedListState = "OutOfSync"
)

type ManagedListList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ManagedList `json:"items"`
}