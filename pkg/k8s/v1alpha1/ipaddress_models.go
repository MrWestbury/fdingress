package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type IpAddressSpec struct {
	IpAddress string `json:"ipAddress"`
}

type IpAddress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec IpAddressSpec `json:"spec"`
}

type IpAddressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []IpAddress `json:"items"`
}

func (in *IpAddress) DeepCopyInto(out *IpAddress) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = IpAddressSpec{
		IpAddress: in.Spec.IpAddress,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *IpAddress) DeepCopyObject() runtime.Object {
	out := IpAddress{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *IpAddressList) DeepCopyObject() runtime.Object {
	out := IpAddressList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]IpAddress, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
