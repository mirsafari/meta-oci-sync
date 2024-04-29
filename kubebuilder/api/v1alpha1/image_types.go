/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ImageSpec defines the desired state of Image
type ImageSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+kubebuilder:validation:MinLength=5
	// SourceRegistry defines location FROM which the image should be mirrored
	SourceRegistry string `json:"sourceRegistry"`
	//+kubebuilder:validation:MinLength=0
	// SourceImage defines the image that should be mirrored
	SourceImage string `json:"sourceImage"`
	//+kkubebuilder:validation:MaxItems=10
	// SourceImageTags specifies all the tags that should be mirrored for a given image
	SourceImageTags []string `json:"sourceImageTags"`
	//+kubebuilder:validation:MinLength=0
	// DestinationRegistry specifies location TO which the image should be mirrored
	DestinationRegistry string `json:"destinationRegistry"`

	//+kubebuilder:validation:MinLength=0
	// CheckInterval specifies how often will the DestinationRegistry be queried
	// to detect drift (missing image).
	// +optional
	CheckInterval string `json:"checkInterval,omitempty"`

	//+kubebuilder:validation:MinLength=0
	// RetryInterval specifies how long to wait between syncs that have errors.
	// +optional
	RetryInterval string `json:"RetryInterval,omitempty"`
}

// ImageStatus defines the observed state of Image
type ImageStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	MirroredState MirroredState `json:"mirroredState,omitempty"`

	// Information when was the last time the mirror was successfully executed.
	// +optional
	LastMirrorTime *metav1.Time `json:"lastMirrorTime,omitempty"`
}

// MirroredState describes the state of the syncing process.
// Only one of the following states can be be specified.
// If none of the following policies is specified, the default one
// is Pending.
// +kubebuilder:validation:Enum=Allow;Forbid;Replace
type MirroredState string

const (
	// PendingMirroredState is the default state and shows that an image is scheduled for sync.
	PendingMirroredState MirroredState = "Pending"

	// SyncedMirroredState shows that all tags for a given image are synced to mirrored repository.
	SyncedMirroredState MirroredState = "Synced"

	// PartiallySyncedMirroredState shows that only some tags for a given image are synced to
	// mirrored repository.
	PartiallySyncedMirroredState MirroredState = "PartiallySynced"

	// ErroredMirroredState shows that the sync failed for all tags of a given image
	ErroredMirroredState MirroredState = "Errored"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Image is the Schema for the images API
type Image struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSpec   `json:"spec,omitempty"`
	Status ImageStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ImageList contains a list of Image
type ImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Image `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Image{}, &ImageList{})
}
