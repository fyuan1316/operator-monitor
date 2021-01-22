package controllers

import (
	"context"
	"fyuan1316/operator-monitor/api/v1alpha1"
	"fyuan1316/operator-monitor/pkg/operator"
	"k8s.io/apimachinery/pkg/api/equality"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/operator-framework/api/pkg/operators/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func isOperatorSuccess(rf v1.RichReference) bool {
	if !(rf.Kind == csvV1alpha1.Kind && rf.APIVersion == csvV1alpha1.APIVersion) {
		return false
	}
	for _, cond := range rf.Conditions {
		if cond.Reason == csvV1alpha1.Conditions[0].Reason &&
			cond.Type == csvV1alpha1.Conditions[0].Type &&
			cond.Status == csvV1alpha1.Conditions[0].Status {
			return true
		}
	}
	return false
}

var csvV1alpha1 = v1.RichReference{
	ObjectReference: &corev1.ObjectReference{
		Kind:       "ClusterServiceVersion",
		APIVersion: "operators.coreos.com/v1alpha1",
	},
	Conditions: []v1.Condition{
		{
			Reason: "InstallSucceeded",
			Type:   "Succeeded",
			Status: "True",
		},
	},
}

func (r *OperatorWatcherReconciler) createOperatorStatus(opr *v1.Operator) error {
	var err error
	crs, err := operator.CrsForOperator(opr, r.Client)
	if err != nil {
		//if apierrors.IsNotFound(err) {
		//	return nil
		//}
		return err
	}
	for _, cr := range crs {
		key := client.ObjectKey{Name: cr.GetName(), Namespace: cr.GetNamespace()}
		current := v1alpha1.OperatorStatus{}
		err = r.Client.Get(context.Background(), key, &current)
		if err != nil {
			if !apierrors.IsNotFound(err) {
				return err
			}
			if err = r.Client.Create(context.Background(), &cr); err != nil {
				return err
			}
			continue
		}
		if !equality.Semantic.DeepDerivative(cr.Spec, current.Spec) {
			current.Spec = cr.Spec
			if err = r.Client.Update(context.Background(), &current); err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

func (r *OperatorWatcherReconciler) deleteOperatorStatus(opr *v1.Operator) error {
	var err error
	crs, err := operator.CrsForOperator(opr, r.Client)
	if err != nil {
		//if apierrors.IsNotFound(err) {
		//	return nil
		//}
		return err
	}
	for _, cr := range crs {
		key := client.ObjectKey{Name: cr.GetName(), Namespace: cr.GetNamespace()}
		current := v1alpha1.OperatorStatus{}
		err = r.Client.Get(context.Background(), key, &current)
		if err != nil {
			if apierrors.IsNotFound(err) {
				continue
			}
			return err
		}
		if err = r.Client.Delete(context.Background(), &current); err != nil {
			return err
		}
	}
	return nil
}
