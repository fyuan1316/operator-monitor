package operator

import (
	"fmt"
	lib "github.com/fyuan1316/operatorlib/api"
	v1 "github.com/operator-framework/api/pkg/operators/v1"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/api/v1alpha1"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/pkg/gvk"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type FlaggerImpl struct {
	Observer
}

func NewFlaggerImpl(opr *v1.Operator, client client.Client) Observable {
	s := opr.GroupVersionKind().GroupVersion().String()
	return &FlaggerImpl{
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

/**
找到当前处于成功状态的operator所管控的cr实例
*/
func (o FlaggerImpl) CRs() ([]v1alpha1.OperatorStatus, error) {
	uList := GetUsList(
		gvk.FlaggerOperatorKind,
		gvk.FlaggerOperatorVersionv1alpha1,
		gvk.FlaggerOperatorGroup)
	// asm中flagger用法会watch所有的namespaces，相当于cluster范围
	//TODO fy 对watch namespace 逻辑进行处理
	// 从deployment中找出具体watch的namespace，
	// 如果namespace为空，则watch 所有ns，任意cr都被监控
	// 如果namespace不为空，则需要查看cr资源创建到的ns是否与watch的ns匹配
	//  - cr的namespace为指定ns或缺省为istio-system
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
	return GenStatus(o.Operator, uList, getFlaggerInstalledNs)
}
func getFlaggerInstalledNs(us unstructured.Unstructured) ([]string, error) {
	crVal, err := getFlaggerComponentNs(us)
	if err != nil {
		return []string{}, nil
	}
	if len(crVal) != 0 {
		return crVal, nil
	}
	return []string{DefaultNSInstallFlaggerTo}, nil
}

const (
	DefaultNSInstallFlaggerTo                     = "istio-system"
	FlaggerDeployNamespaceKeyInFlaggerChartValues = "istioNamespace"
)

func getFlaggerComponentNs(us unstructured.Unstructured) ([]string, error) {
	oprSpec := lib.OperatorSpec{}
	if err := util.JsonConvert(us.Object, &oprSpec); err != nil {
		return nil, err
	}
	if oprSpec.Namespace != "" {
		return []string{oprSpec.Namespace}, nil
	}

	var m map[string]interface{}
	if err := util.StringConvert(oprSpec.Parameters, &m); err != nil {
		return nil, err
	}
	if v, exist := m[FlaggerDeployNamespaceKeyInFlaggerChartValues]; exist {
		if installed, ok := v.(string); ok {
			return []string{installed}, nil
		}
	}

	return []string{}, nil
}

//
//// flagger 在启动参数中通过namespace来指定监控范围 https://github.com/fluxcd/flagger/blob/main/cmd/flagger/main.go#L104
//// 如果namespace值为空代表当前flagger实例的监控范围是整个集群范围，所以任意Canary均可有效调度。
//// 反之，则仅针对指定namespace的Canary起作用。
////
//// TODO fy 在ASM 内部operator（asm/flagger）实现中，增加watchnamespace 这种可见性限制
////   还需要增加控制平台真正要安装到的 namespace，在spec中增加一个值，另外需要能与values中的值进行merge，
////   spec中的值优先
//
//func (o FlaggerImpl) getWatchedNamespaces() (string, error) {
//	var err error
//	deploy := &appsv1.Deployment{}
//
//	name, ns, err := ValidOperatorName(o.Operator.Name)
//	if err != nil {
//		return "", err
//	}
//	err = o.Client.Get(context.Background(), types.NamespacedName{
//		Namespace: ns,
//		Name:      name,
//	}, deploy)
//	if err != nil {
//		return "", err
//	}
//
//	var watchedNamespace string
//	for _, c := range deploy.Spec.Template.Spec.Containers {
//		for _, arg := range c.Args {
//			if strings.Contains(arg, "namespace") {
//				kvPair := strings.Split(arg, "=")
//				if len(kvPair) == 2 {
//					watchedNamespace = kvPair[1]
//				}
//				break
//			}
//		}
//	}
//	return watchedNamespace, err
//}
