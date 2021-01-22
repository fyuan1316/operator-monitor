package operator

import (
	"fmt"
	"fyuan1316/operator-monitor/api/v1alpha1"
	"fyuan1316/operator-monitor/pkg/gvk"
	v1 "github.com/operator-framework/api/pkg/operators/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type JaegerImpl struct {
	Observer
}

func NewJaegerImpl(opr *v1.Operator, client client.Client) Observable {
	s := opr.GroupVersionKind().GroupVersion().String()
	return &JaegerImpl{
		Observer{
			Client: client,
			Operator: corev1.TypedLocalObjectReference{
				//TypedLocalObjectReference: corev1.TypedLocalObjectReference{
				APIGroup: &s,
				Kind:     opr.Kind,
				Name:     opr.Name,
				//},
			},
		},
	}
}

func (o JaegerImpl) CRs() ([]v1alpha1.OperatorStatus, error) {
	uList := GetUsList(
		gvk.JaegerOperatorKind,
		gvk.JaegerOperatorVersionv1alpha1,
		gvk.JaegerOperatorGroup)

	watchNs, err := o.GetWatchNamespace(CheckWatchNamespaceInEnv)
	if err != nil {
		return nil, err
	}
	uList, err = o.GetValidCRs(watchNs, uList)
	if err != nil {
		return nil, err
	}
	if len(uList.Items) == 0 {
		oprLog.Info(fmt.Sprintf("not found cr for operator %s", o.Operator.Name))
		return nil, nil
	}
	return GenStatus(o.Operator, uList, getJaegerInstalledNs)
}

func getJaegerInstalledNs(us unstructured.Unstructured) ([]string, error) {
	// jaeger resources created in same namespace with jaeger operator cr
	// TODO fy 检查olm targetnamespace 是否会对此有副作用
	return []string{us.GetNamespace()}, nil
}

/*
func (o JaegerImpl) CRs() ([]v1alpha1.OperatorStatus, error) {
	uList := GetUsList(
		gvk.JaegerOperatorKind,
		gvk.JaegerOperatorVersionv1alpha1,
		gvk.JaegerOperatorGroup)
	var err error
	// jaegertracing.io/operated-by: istio-system.jaeger-operator
	name, ns, err := ValidOperatorName(o.Operator.Name)
	if err != nil {
		return nil, err
	}
	//TODO fy CR与operator的匹配逻辑
	// operator 对应的 deployment处于相同ns，
	// 从deployment envVar中取key为OPERATOR_NAME的name值，加上ns，组成ns.name。
	// 如果当前cr 的ns.name 与其不匹配则忽略
	err = o.Client.List(context.Background(), &uList,
		Selector(map[string]string{"jaegertracing.io/operated-by": fmt.Sprintf("%s.%s", ns, name)}))
	if err != nil {
		return nil, err
	}
	if len(uList.Items) == 0 {
		oprLog.Info(fmt.Sprintf("not found cr for operator %s", o.Operator.Name))
		return nil, nil
	}
	return GenStatus(o.Operator, uList), nil
}
*/
