package status

import (
	operatorv1alpha1 "gitlab-ce.alauda.cn/micro-service/operator-monitor/api/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type QueryComponentStatusFn func(OperatorStatusGet, map[string]string) (map[string]operatorv1alpha1.ComponentStatus, error)

type ComponentStatusQuery interface {
	GenQueryFn(namespaces []string) []QueryComponentStatusFn
}

func Selector(m labels.Set) *client.ListOptions {
	var listOpts = client.ListOptions{}
	listOpts.LabelSelector = labels.SelectorFromSet(m)
	return &listOpts
}

type ResourceQuery struct {
	Queryer []ComponentStatusQuery
}

func (rq *ResourceQuery) Add(q ComponentStatusQuery) {
	if rq != nil {
		rq.Queryer = append(rq.Queryer, q)
	}
}

func (rq *ResourceQuery) Queryers(namespaces []string) []QueryComponentStatusFn {
	var allComponentsQueryFn []QueryComponentStatusFn
	for _, query := range rq.Queryer {
		allComponentsQueryFn = append(allComponentsQueryFn, query.GenQueryFn(namespaces)...)
	}
	return allComponentsQueryFn
}

var ResourceQueryHelper ResourceQuery

func init() {
	ResourceQueryHelper = ResourceQuery{}
	ResourceQueryHelper.Add(DeploymentQuery{})
}
