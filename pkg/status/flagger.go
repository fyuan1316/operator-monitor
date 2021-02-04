package status

import (
	operatorv1alpha1 "gitlab-ce.alauda.cn/micro-service/operator-monitor/api/v1alpha1"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/pkg/util"
)

type FlaggerStatusGetter struct {
	OperatorStatusGet
}

var _ OperatorStatusGetter = FlaggerStatusGetter{}

func NewFlaggerStatusGetter(g OperatorStatusGet) FlaggerStatusGetter {
	getter := FlaggerStatusGetter{
		OperatorStatusGet: g,
	}
	return getter
}
func (g FlaggerStatusGetter) Status() (operatorv1alpha1.OperatorStatusStatus, error) {
	allComponentsQueryFn := ResourceQueryHelper.Queryers(g.OperatorStatus.Spec.InstalledNamespace)
	matches := map[string]string{
		FlaggerOwningResourceKey:          g.OperatorStatus.Spec.CR.Name,
		FlaggerOwningResourceNamespaceKey: g.OperatorStatus.Spec.CR.Namespace,
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
		statusLog.Error(errs.ErrorOrNil(), "get flagger components status")
	}
	oss := operatorv1alpha1.OperatorStatusStatus{}
	oss.ComponentStatus = mergedMap
	return oss, errs.ErrorOrNil()
}
