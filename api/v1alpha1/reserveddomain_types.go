/*
MIT License

Copyright (c) 2022 ngrok, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ReservedDomainSpec defines the desired state of ReservedDomain
type ReservedDomainSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Domain is the domain name to reserve
	// +kubebuilder:validation:Required
	Domain string `json:"domain"`

	// Region is the region in which to reserve the domain
	// +kubebuilder:validation:Required
	Region string `json:"region,omitempty"`

	// Description is a human-readable description of the reserved domain
	// +kubebuilder:default:=`Created by ngrok-ingress-controller`
	Description string `json:"description,omitempty"`

	// Metadata is a string of arbitrary data associated with the reserved domain
	// +kubebuilder:default:=`{"owned-by":"ngrok-ingress-controller"}`
	Metadata string `json:"metadata,omitempty"`
}

// ReservedDomainStatus defines the observed state of ReservedDomain
type ReservedDomainStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ID is the unique identifier of the reserved domain
	ID string `json:"id,omitempty"`

	// Domain is the domain that was reserved
	Domain string `json:"domain,omitempty"`

	// Region is the region in which the reserved domain was created
	Region string `json:"region,omitempty"`

	// URI of the reserved domain API resource
	URI string `json:"uri,omitempty"`

	// CNAMETarget is the CNAME target for the reserved domain
	CNAMETarget *string `json:"cname_target,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="ID",type=string,JSONPath=`.status.id`,description="Reserved Domain ID"
//+kubebuilder:printcolumn:name="Region",type=string,JSONPath=`.status.region`,description="Region"
//+kubebuilder:printcolumn:name="Domain",type=string,JSONPath=`.status.domain`,description="Domain"
//+kubebuilder:printcolumn:name="CNAME Target",type=string,JSONPath=`.status.cname_target`,description="CNAME Target"
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="Age"

// ReservedDomain is the Schema for the reserveddomains API
type ReservedDomain struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReservedDomainSpec   `json:"spec,omitempty"`
	Status ReservedDomainStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ReservedDomainList contains a list of ReservedDomain
type ReservedDomainList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReservedDomain `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ReservedDomain{}, &ReservedDomainList{})
}
