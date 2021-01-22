package status

import (
	operatorv1alpha1 "fyuan1316/operator-monitor/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	statusLog = ctrl.Log.WithName("operatorStatus")
)

func CheckComponentsStatus(client client.Client, ins *operatorv1alpha1.OperatorStatus) (operatorv1alpha1.OperatorStatusStatus, error) {
	getter := NewOperatorStatusGet(client, ins)
	return getter.Status()
}

type OperatorStatusGetter interface {
	Status() (operatorv1alpha1.OperatorStatusStatus, error)
}
type OperatorStatusGet struct {
	Client         client.Client
	OperatorStatus *operatorv1alpha1.OperatorStatus
}

func NewOperatorStatusGet(client client.Client, ins *operatorv1alpha1.OperatorStatus) OperatorStatusGet {
	get := OperatorStatusGet{
		Client:         client,
		OperatorStatus: ins,
	}
	return get
}
func (g OperatorStatusGet) OperatorCRKind() string {
	return g.OperatorStatus.Spec.CR.Kind
}
func (g OperatorStatusGet) Status() (operatorv1alpha1.OperatorStatusStatus, error) {
	var getter OperatorStatusGetter
	switch g.OperatorCRKind() {
	case operatorv1alpha1.KindAsm:
		getter = NewAsmStatusGetter(g)
	case operatorv1alpha1.KindFlagger:
		getter = NewFlaggerStatusGetter(g)
	case operatorv1alpha1.KindJaeger:
		getter = NewJaegerStatusGetter(g)
	case operatorv1alpha1.KindIstioOperator:
		getter = NewIstioStatusGetter(g)
	}
	if getter == nil {
		return operatorv1alpha1.OperatorStatusStatus{}, nil
	}
	return getter.Status()
}
