package status

import (
	operatorv1alpha1 "gitlab-ce.alauda.cn/micro-service/operator-monitor/api/v1alpha1"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/pkg/operator"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/pkg/util"
)

var _ OperatorStatusGetter = JaegerStatusGetter{}

type JaegerStatusGetter struct {
	OperatorStatusGet
}

func NewJaegerStatusGetter(g OperatorStatusGet) JaegerStatusGetter {
	getter := JaegerStatusGetter{
		OperatorStatusGet: g,
	}
	return getter
}
func (g JaegerStatusGetter) Status() (operatorv1alpha1.OperatorStatusStatus, error) {
	allComponentsQueryFn := ResourceQueryHelper.Queryers(g.OperatorStatus.Spec.InstalledNamespace)
	// just use the name without namespace
	name, _, _ := operator.ValidOperatorName(g.OperatorStatus.Spec.Operator.Name)
	matches := map[string]string{
		JaegerInstanceKey: g.OperatorStatus.Spec.CR.Name,
		JaegerOperatorKey: name,
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
		statusLog.Error(errs.ErrorOrNil(), "get jaeger components status")
	}
	oss := operatorv1alpha1.OperatorStatusStatus{}
	oss.ComponentStatus = mergedMap
	return oss, errs.ErrorOrNil()

}
