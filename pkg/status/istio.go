package status

import (
	operatorv1alpha1 "gitlab-ce.alauda.cn/micro-service/operator-monitor/api/v1alpha1"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/pkg/util"
)

var _ OperatorStatusGetter = IstioStatusGetter{}

type IstioStatusGetter struct {
	OperatorStatusGet
}

func NewIstioStatusGetter(g OperatorStatusGet) IstioStatusGetter {
	getter := IstioStatusGetter{
		OperatorStatusGet: g,
	}
	return getter
}
func (g IstioStatusGetter) Status() (operatorv1alpha1.OperatorStatusStatus, error) {
	allComponentsQueryFn := ResourceQueryHelper.Queryers(g.OperatorStatus.Spec.InstalledNamespace)
	matches := map[string]string{
		IstioOwningResourceKey:          g.OperatorStatus.Spec.CR.Name,
		IstioOwningResourceNamespaceKey: g.OperatorStatus.Spec.CR.Namespace,
	}
	var mergedMap = make(map[string]operatorv1alpha1.ComponentStatus)
	var errs util.Errors
	for _, fn := range allComponentsQueryFn {
		m, err := fn(g.OperatorStatusGet, matches)
		for k, v := range m {
			mergedMap[k] = v
		}
		errs.Append(err)
	}
	if !errs.Empty() {
		statusLog.Error(errs.ErrorOrNil(), "get istio components status")
	}
	oss := operatorv1alpha1.OperatorStatusStatus{}
	oss.ComponentStatus = mergedMap
	return oss, errs.ErrorOrNil()

}
