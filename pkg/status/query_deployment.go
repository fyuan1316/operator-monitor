package status

import (
	"context"
	operatorv1alpha1 "gitlab-ce.alauda.cn/micro-service/operator-monitor/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
)

type DeploymentQuery struct{}

func (d DeploymentQuery) GenQueryFn(namespaces []string) []QueryComponentStatusFn {
	var fnList []QueryComponentStatusFn
	for _, ns := range namespaces {
		if ns == "" {
			continue
		}
		fnList = append(fnList, getDeploymentsStatus(ns))
	}
	return fnList
}

var _ ComponentStatusQuery = DeploymentQuery{}

func getDeploymentsStatus(ns string) QueryComponentStatusFn {
	return func(get OperatorStatusGet, labels map[string]string) (map[string]operatorv1alpha1.ComponentStatus, error) {
		opts := Selector(labels)
		opts.Namespace = ns
		deployList := appsv1.DeploymentList{}
		if err := get.Client.List(context.Background(), &deployList, opts); err != nil {
			return nil, err
		}
		return statusForDeploymentList(&deployList), nil
	}
}

func statusForDeploymentList(deployList *appsv1.DeploymentList) map[string]operatorv1alpha1.ComponentStatus {
	ComponentsStatus := make(map[string]operatorv1alpha1.ComponentStatus)
	for _, deploy := range deployList.Items {
		ready := operatorv1alpha1.ResourceStates.NotReady
		if deploy.Status.ReadyReplicas > 0 {
			ready = operatorv1alpha1.ResourceStates.Ready
		}
		c := operatorv1alpha1.ComponentStatus{}
		c.Namespace = deploy.Namespace
		c.Name = deploy.Name
		c.Kind = deploy.Kind
		s := deploy.GroupVersionKind().GroupVersion().String()
		c.APIGroup = &s
		c.Status = ready
		ComponentsStatus[c.Name] = c
	}
	return ComponentsStatus
}
