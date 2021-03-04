package mock

import (
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/api/v1alpha1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Case struct {
	Deployment     *v1.Deployment
	OperatorStatus *v1alpha1.OperatorStatus
	ExpectStatus   *v1alpha1.OperatorStatusStatus
}

func NewAsmCase() Case {
	c := Case{}
	c.Deployment = getDeployment(
		"asm-dep",
		"default",
		map[string]string{
			"asm.operator.alauda.io/owning-resource":           "asm",
			"asm.operator.alauda.io/owning-resource-namespace": "",
		},
	)
	c.OperatorStatus = getAsmStatus()
	c.ExpectStatus = c.getWantStatus()
	return c
}

func NewFlaggerCase() Case {
	c := Case{}
	c.Deployment = getDeployment(
		"flagger-dep",
		"default",
		map[string]string{
			"flagger.operator.alauda.io/owning-resource":           "flagger",
			"flagger.operator.alauda.io/owning-resource-namespace": "istio-system",
		},
	)
	c.OperatorStatus = getFlaggerStatus()
	c.ExpectStatus = c.getWantStatus()
	return c
}
func NewJaegerCase() Case {
	c := Case{}
	c.Deployment = getDeployment(
		"jaeger-dep",
		"default",

		map[string]string{
			"app.kubernetes.io/instance":   "jaeger-prod",
			"app.kubernetes.io/managed-by": "jaeger-operator",
		},
	)
	c.OperatorStatus = getJaegerStatus()
	c.ExpectStatus = c.getWantStatus()
	return c
}
func NewIstioCase() Case {
	c := Case{}
	c.Deployment = getDeployment(
		"istio-dep",
		"default",
		map[string]string{
			"install.operator.istio.io/owning-resource":           "istio",
			"install.operator.istio.io/owning-resource-namespace": "istio-system",
		},
	)
	c.OperatorStatus = getIstioStatus()
	c.ExpectStatus = c.getWantStatus()
	return c
}

func (c Case) getWantStatus() *v1alpha1.OperatorStatusStatus {
	create := v1alpha1.OperatorStatusStatus{}
	depApiGroup := "apps/v1"
	create.ComponentStatus = map[string]v1alpha1.ComponentStatus{
		c.Deployment.Name: {
			NameSpacedTypedObjectReference: v1alpha1.NameSpacedTypedObjectReference{
				TypedLocalObjectReference: corev1.TypedLocalObjectReference{
					APIGroup: &depApiGroup,
					Kind:     "Deployment",
					Name:     c.Deployment.Name,
				},
				Namespace: c.Deployment.Namespace,
			},
			Status: "NotReady",
		},
	}
	return &create
}
func getAsmStatus() *v1alpha1.OperatorStatus {
	create := v1alpha1.OperatorStatus{}
	create.Name = "asm"
	create.Spec.Operator.Name = "asm-operator.istio-system"
	create.Spec.InstalledNamespace = []string{"default"}
	apiGroup := "operators.coreos.com/v1"
	create.Spec.Operator.APIGroup = &apiGroup
	create.Spec.Operator.Kind = "Operator"
	create.Spec.CR.Kind = "Asm"
	create.Spec.CR.Name = "asm"
	crApiGroup := "operator.alauda.io"
	create.Spec.CR.APIGroup = &crApiGroup
	return &create
}
func getFlaggerStatus() *v1alpha1.OperatorStatus {
	create := v1alpha1.OperatorStatus{}
	create.Name = "flagger"
	create.Spec.Operator.Name = "flagger-operator.istio-system"
	create.Spec.InstalledNamespace = []string{"default"}
	apiGroup := "operators.coreos.com/v1"
	create.Spec.Operator.APIGroup = &apiGroup
	create.Spec.Operator.Kind = "Operator"
	create.Spec.CR.Kind = "Flagger"
	create.Spec.CR.Name = "flagger"
	create.Spec.CR.Namespace = "istio-system"
	crApiGroup := "operator.alauda.io"
	create.Spec.CR.APIGroup = &crApiGroup
	return &create
}
func getJaegerStatus() *v1alpha1.OperatorStatus {
	create := v1alpha1.OperatorStatus{}
	create.Name = "jaeger"
	create.Spec.Operator.Name = "jaeger-operator.istio-system"
	create.Spec.InstalledNamespace = []string{"default"}
	apiGroup := "operators.coreos.com/v1"
	create.Spec.Operator.APIGroup = &apiGroup
	create.Spec.Operator.Kind = "Operator"
	create.Spec.CR.Kind = "Jaeger"
	create.Spec.CR.Name = "jaeger-prod"
	create.Spec.CR.Namespace = "istio-system"
	crApiGroup := "jaegertracing.io/v1"
	create.Spec.CR.APIGroup = &crApiGroup
	return &create
}
func getIstioStatus() *v1alpha1.OperatorStatus {
	create := v1alpha1.OperatorStatus{}
	create.Name = "istio"
	create.Spec.Operator.Name = "istio-operator.istio-system"
	create.Spec.InstalledNamespace = []string{"default"}
	apiGroup := "operators.coreos.com/v1"
	create.Spec.Operator.APIGroup = &apiGroup
	create.Spec.Operator.Kind = "Operator"
	create.Spec.CR.Kind = "IstioOperator"
	create.Spec.CR.Name = "istio"
	create.Spec.CR.Namespace = "istio-system"
	crApiGroup := "install.istio.io/v1alpha1"
	create.Spec.CR.APIGroup = &crApiGroup
	return &create
}

func getDeployment(name, namespace string, labels map[string]string) *v1.Deployment {
	d := v1.Deployment{}
	d.Name = name
	d.Namespace = namespace
	d.Labels = map[string]string{
		"app": name,
	}
	for k, v := range labels {
		d.Labels[k] = v
	}
	d.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app": name,
		},
	}
	d.Spec.Template.Labels = map[string]string{
		"app": name,
	}
	d.Spec.Template.Spec.Containers = []corev1.Container{
		{Name: name, Image: "sleep"},
	}
	return &d
}

func GetObjects() []client.Object {
	asm := NewAsmCase()
	flagger := NewFlaggerCase()
	jaeger := NewJaegerCase()
	istio := NewIstioCase()

	return []client.Object{
		asm.Deployment, asm.OperatorStatus,
		flagger.Deployment, flagger.OperatorStatus,
		jaeger.Deployment, jaeger.OperatorStatus,
		istio.Deployment, istio.OperatorStatus,
	}
}
