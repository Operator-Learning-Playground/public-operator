package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// Generics
type Generic struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GenericSpec `json:"spec,omitempty"`
}

type GenericSpec struct {
	Gvr      string       `json:"gvr,omitempty"`
	Template SpecTemplate `json:"template,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// GenericList
type GenericList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Generic `json:"items"`
}

type SpecTemplate map[string]interface{}

func (in *SpecTemplate) DeepCopyInto(out *SpecTemplate) {
	if in == nil {
		return
	}
	b, err := yaml.Marshal(in)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(b, &out)
	if err != nil {
		return
	}

}
