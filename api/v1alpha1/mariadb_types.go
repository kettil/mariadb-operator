/*
Copyright 2022.

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
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BootstrapFromBackup struct {
	// +kubebuilder:validation:Required
	BackupRef corev1.LocalObjectReference `json:"backupRef"`
}

type Exporter struct {
	// +kubebuilder:validation:Required
	Image     Image                        `json:"image"`
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

type Metrics struct {
	Exporter Exporter `json:"exporter"`
	// +kubebuilder:default=false
	CreateServiceMonitor bool `json:"createServiceMonitor"`
}

// MariaDBSpec defines the desired state of MariaDB
type MariaDBSpec struct {
	// +kubebuilder:validation:Required
	RootPasswordSecretKeyRef corev1.SecretKeySelector `json:"rootPasswordSecretKeyRef"`

	Database             *string                   `json:"database,omitempty"`
	Username             *string                   `json:"username,omitempty"`
	PasswordSecretKeyRef *corev1.SecretKeySelector `json:"passwordSecretKeyRef,omitempty"`

	// +kubebuilder:validation:Required
	Image            Image                         `json:"image"`
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// +kubebuilder:default=3306
	Port int32 `json:"port,omitempty"`

	// +kubebuilder:validation:Required
	Storage Storage `json:"storage"`

	BootstrapFromBackup *BootstrapFromBackup `json:"bootstrapFromBackup,omitempty"`

	Metrics *Metrics `json:"metrics,omitempty"`

	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	Env     []v1.EnvVar        `json:"env,omitempty"`
	EnvFrom []v1.EnvFromSource `json:"envFrom,omitempty"`
}

// MariaDBStatus defines the observed state of MariaDB
type MariaDBStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func (s *MariaDBStatus) SetCondition(condition metav1.Condition) {
	if s.Conditions == nil {
		s.Conditions = make([]metav1.Condition, 0)
	}
	meta.SetStatusCondition(&s.Conditions, condition)
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=mdb
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].message"
// +kubebuilder:printcolumn:name="Storage Class",type="string",JSONPath=".spec.storage.className"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// MariaDB is the Schema for the mariadbs API
type MariaDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MariaDBSpec   `json:"spec"`
	Status MariaDBStatus `json:"status,omitempty"`
}

func (m *MariaDB) IsReady() bool {
	return meta.IsStatusConditionTrue(m.Status.Conditions, ConditionTypeReady)
}

func (m *MariaDB) IsBootstrapped() bool {
	return meta.IsStatusConditionTrue(m.Status.Conditions, ConditionTypeBootstrapped)
}

func (m *MariaDB) IsBeingDeleted() bool {
	return !m.DeletionTimestamp.IsZero()
}

// +kubebuilder:object:root=true

// MariaDBList contains a list of MariaDB
type MariaDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MariaDB `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MariaDB{}, &MariaDBList{})
}
