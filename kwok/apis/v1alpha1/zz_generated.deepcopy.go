//go:build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/awslabs/operatorpkg/status"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KWOKNodeClass) DeepCopyInto(out *KWOKNodeClass) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KWOKNodeClass.
func (in *KWOKNodeClass) DeepCopy() *KWOKNodeClass {
	if in == nil {
		return nil
	}
	out := new(KWOKNodeClass)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KWOKNodeClass) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KWOKNodeClassList) DeepCopyInto(out *KWOKNodeClassList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KWOKNodeClass, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KWOKNodeClassList.
func (in *KWOKNodeClassList) DeepCopy() *KWOKNodeClassList {
	if in == nil {
		return nil
	}
	out := new(KWOKNodeClassList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KWOKNodeClassList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KWOKNodeClassSpec) DeepCopyInto(out *KWOKNodeClassSpec) {
	*out = *in
	out.NodeRegistrationDelay = in.NodeRegistrationDelay
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KWOKNodeClassSpec.
func (in *KWOKNodeClassSpec) DeepCopy() *KWOKNodeClassSpec {
	if in == nil {
		return nil
	}
	out := new(KWOKNodeClassSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KWOKNodeClassStatus) DeepCopyInto(out *KWOKNodeClassStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]status.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KWOKNodeClassStatus.
func (in *KWOKNodeClassStatus) DeepCopy() *KWOKNodeClassStatus {
	if in == nil {
		return nil
	}
	out := new(KWOKNodeClassStatus)
	in.DeepCopyInto(out)
	return out
}
