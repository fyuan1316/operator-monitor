package gvk

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	OperatorKind      = "Operator"
	OperatorGroup     = "operators.coreos.com"
	OperatorVersionv1 = "v1"

	ASMOperatorKind            = "Asm"
	ASMOperatorGroup           = "operator.alauda.io"
	ASMOperatorVersionv1alpha1 = "v1alpha1"

	FlaggerOperatorKind            = "Flagger"
	FlaggerOperatorGroup           = "operator.alauda.io"
	FlaggerOperatorVersionv1alpha1 = "v1alpha1"

	IstioOperatorKind            = "IstioOperator"
	IstioOperatorGroup           = "install.istio.io"
	IstioOperatorVersionv1alpha1 = "v1alpha1"

	JaegerOperatorKind            = "Jaeger"
	JaegerOperatorGroup           = "jaegertracing.io"
	JaegerOperatorVersionv1alpha1 = "v1"
)

type OperatorUnstructured struct {
	unstructured.Unstructured
}

func NewOperatorUnstructured() *OperatorUnstructured {
	ou := OperatorUnstructured{}
	ou.Unstructured.SetAPIVersion(OperatorVersionv1)
	ou.Unstructured.SetKind(OperatorKind)
	ou.Unstructured.SetGroupVersionKind(schema.GroupVersionKind{Group: OperatorGroup, Kind: OperatorKind, Version: OperatorVersionv1})
	return &ou
}

func (ou *OperatorUnstructured) WithName(key client.ObjectKey) *OperatorUnstructured {
	ou.SetName(key.Name)
	ou.SetNamespace(key.Namespace)
	return ou
}
func (ou *OperatorUnstructured) GetUnstructured() unstructured.Unstructured {
	return ou.Unstructured
}
