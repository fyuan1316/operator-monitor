package operator

import (
	"fmt"
	lib "github.com/fyuan1316/operatorlib/api"
	v1 "github.com/operator-framework/api/pkg/operators/v1"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/api/v1alpha1"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/pkg/gvk"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/pkg/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AsmImpl struct {
	Observer
}

func NewAsmImpl(opr *v1.Operator, client client.Client) Observable {
	s := opr.GroupVersionKind().GroupVersion().String()
	impl := AsmImpl{}
	impl.Observer.Client = client
	impl.Observer.Operator.Kind = opr.Kind
	impl.Observer.Operator.Name = opr.Name
	impl.Observer.Operator.APIGroup = &s
	return &impl
}

func (o AsmImpl) CRs() ([]v1alpha1.OperatorStatus, error) {
	uList := GetUsList(
		gvk.ASMOperatorKind,
		gvk.ASMOperatorVersionv1alpha1,
		gvk.ASMOperatorGroup)
	// Asm 为集群级别的资源，将所有cr均视为合法存在
	//TODO fy 创建cr时增加cr与operator的对应关系
	// 当前Asm 为集群级别的资源，将所有cr均视为合法存在
	// 对watch namespace 处理的逻辑
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
	return GenStatus(o.Operator, uList, getAsmInstalledNs)
}

func getAsmInstalledNs(us unstructured.Unstructured) ([]string, error) {
	crVal, err := getAsmComponentNs(us)
	if err != nil {
		return []string{}, nil
	}
	if len(crVal) != 0 {
		return crVal, nil
	}
	return []string{DefaultNSInstallAsmTo}, nil
}

const (
	DefaultNSInstallAsmTo = "cpaas-system"
)

func getAsmComponentNs(us unstructured.Unstructured) ([]string, error) {
	oprSpec := lib.OperatorSpec{}
	if err := util.JsonConvert(us.Object, &oprSpec); err != nil {
		return nil, err
	}
	if oprSpec.Namespace != "" {
		return []string{oprSpec.Namespace}, nil
	}
	/*  asm 资源安装namespace 由$release.namespace 决定，此处不需要考虑chart values
	var m map[string]interface{}
	if err := util.StringConvert(oprSpec.Parameters, &m); err != nil {
		return nil, err
	}
	if v, exist := m[FlaggerDeployNamespaceKeyInFlaggerChartValues]; exist {
		if installed, ok := v.(string); ok {
			return []string{installed}, nil
		}
	}
	*/
	return []string{}, nil
}
