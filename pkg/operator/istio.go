package operator

import (
	"fmt"
	v1 "github.com/operator-framework/api/pkg/operators/v1"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/api/v1alpha1"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/pkg/gvk"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/pkg/util"
	iopv1alpha1 "istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IstioImpl struct {
	Observer
}

func NewIstioImpl(opr *v1.Operator, client client.Client) Observable {
	s := opr.GroupVersionKind().GroupVersion().String()
	return &IstioImpl{
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

func (o IstioImpl) CRs() ([]v1alpha1.OperatorStatus, error) {
	var err error
	uList := GetUsList(
		gvk.IstioOperatorKind,
		gvk.IstioOperatorVersionv1alpha1,
		gvk.IstioOperatorGroup)

	// IstioOperator 要求必须传入WATCH_NAMESPACE envVar https://github.com/istio/istio/blob/master/operator/cmd/operator/server.go#L109-L113
	// 如果WATCH_NAMESPACE值为空代表当前operator的监控范围是整个集群范围，所以任意CR均可有效调度。
	// 反之，则仅针对指定WATCH_NAMESPACE的CR起作用。
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
	return GenStatus(o.Operator, uList, getIstioInstalledNs)
}

func getIstioInstalledNs(us unstructured.Unstructured) ([]string, error) {
	iops := iopv1alpha1.IstioOperator{}
	if err := util.JsonConvert(us.Object, &iops); err != nil {
		return nil, err
	}
	ns := iopv1alpha1.Namespace(iops.Spec)
	if ns == "" {
		// Namespace to install control plane resources into. If unset, Istio will be installed into the same namespace
		// as the IstioOperator CR.
		ns = us.GetNamespace()
	}
	return []string{ns}, nil
}

//const DefaultNSInstallIstioTo = "istio-system"

//
//// IstioOperator 要求必须传入WATCH_NAMESPACE envVar https://github.com/istio/istio/blob/master/operator/cmd/operator/server.go#L109-L113
//// 如果WATCH_NAMESPACE值为空代表当前operator的监控范围是整个集群范围，所以任意CR均可有效调度。
//// 反之，则仅针对指定WATCH_NAMESPACE的CR起作用。
//
//// 当前operator是否启用watchedNamespaces，如果有，cr只能在其中才有效
//// 如果没有启用，默认watch istio-system，则cr在其中才有效
//func (o IstioImpl) getWatchedNamespaces() (string, error) {
//	var err error
//	deployList := &appsv1.DeploymentList{}
//	// operators.coreos.com/istio-operator.operators: ""
//	opts := Selector(map[string]string{fmt.Sprintf("operators.coreos.com/%s", o.Operator.Name): ""})
//	opts.Namespace = o.Operator.Namespace
//	err = o.Client.List(context.Background(), deployList, opts)
//	if err != nil {
//		return "", err
//	}
//	if len(deployList.Items) == 0 {
//		return "", pkgerrors.New("not found")
//	}
//	deploy := deployList.Items[0]
//	var watchedNamespace string
//	for _, c := range deploy.Spec.Template.Spec.Containers {
//		for _, env := range c.Env {
//			if env.Name == "WATCH_NAMESPACE" {
//				watchedNamespace = env.Value
//				break
//			}
//		}
//	}
//	if watchedNamespace == "" {
//		watchedNamespace = "istio-system"
//	}
//	return watchedNamespace, err
//}
