/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	operatorv1alpha1 "fyuan1316/operator-monitor/api/v1alpha1"
	"fyuan1316/operator-monitor/pkg/gvk"
	"fyuan1316/operator-monitor/pkg/operator"
	error2 "fyuan1316/operator-monitor/pkg/operator/error"
	"fyuan1316/operator-monitor/pkg/util"
	"github.com/fyuan1316/operatorlib/event"
	"github.com/go-logr/logr"
	v1 "github.com/operator-framework/api/pkg/operators/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	finalizer           = "operatorwatchers.operator.alauda.io"
	finalizerMaxRetries = 1
)

// OperatorWatcherReconciler reconciles a OperatorWatcher object
type OperatorWatcherReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=operator.alauda.io,resources=operatorwatchers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operator.alauda.io,resources=operatorwatchers/status,verbs=get;update;patch

func (r *OperatorWatcherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error
	log := r.Log.WithValues("operatorwatcher", req.NamespacedName)
	log.Info(fmt.Sprintf("Starting reconcile loop for %v", req.NamespacedName))
	defer log.Info(fmt.Sprintf("Finish reconcile loop for %v", req.NamespacedName))

	u := gvk.NewOperatorUnstructured().
		WithName(client.ObjectKey{Namespace: req.Namespace, Name: req.Name}).
		GetUnstructured()
	err = r.Get(context.Background(), req.NamespacedName, &u)
	if err != nil {
		if errors.IsNotFound(err) {
			// CR not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	//check
	_, _, err = operator.ValidOperatorName(u.GetName())
	if err != nil {
		r.Recorder.Event(&u, event.WarningEvent, error2.OperatorNameUnsupported, err.Error())
		return reconcile.Result{}, err
	}

	opr := &v1.Operator{}
	if err = util.JsonConvert(u.Object, opr); err != nil {
		r.Recorder.Event(&u, event.WarningEvent, error2.OperatorCastError, err.Error())
		return reconcile.Result{}, err
	}
	deleted := !u.GetDeletionTimestamp().IsZero()
	finalizers := sets.NewString(opr.GetFinalizers()...)
	if deleted {
		if !finalizers.Has(finalizer) {
			return reconcile.Result{}, nil
		}
		if err := r.deleteOperatorStatus(opr); err != nil {
			return reconcile.Result{}, err
		}
		finalizers.Delete(finalizer)
		opr.SetFinalizers(finalizers.List())
		finalizerError := r.Client.Update(context.TODO(), opr)
		for retryCount := 0; errors.IsConflict(finalizerError) && retryCount < finalizerMaxRetries; retryCount++ {
			log.Info("conflict during finalizer removal, retrying")
			_ = r.Client.Get(context.TODO(), req.NamespacedName, opr)
			finalizers = sets.NewString(opr.GetFinalizers()...)
			finalizers.Delete(finalizer)
			opr.SetFinalizers(finalizers.List())
			finalizerError = r.Client.Update(context.TODO(), opr)
		}
		if finalizerError != nil {
			if errors.IsNotFound(finalizerError) {
				log.Info("Could not remove finalizer from %v: the object was deleted", req)
				return reconcile.Result{}, nil
			} else if errors.IsConflict(finalizerError) {
				log.Info("Could not remove finalizer from %v due to conflict. Operation will be retried in next reconcile attempt", req)
				return reconcile.Result{}, nil
			}
			log.Error(finalizerError, "error removing finalizer")
			return reconcile.Result{}, finalizerError
		}
		return reconcile.Result{}, nil
	}

	// if operator install successful, create operatorStatus for cr
	if opr.Status.Components == nil {
		return ctrl.Result{}, nil
	}
	success := false
	for _, c := range opr.Status.Components.Refs {
		if isOperatorSuccess(c) {
			success = true
			break
		}
	}
	if success {
		// 正常运行的operator 并且没有finalizer的需要设置
		if !finalizers.Has(finalizer) {
			log.Info("Adding finalizer %v to %v", finalizer, req)
			finalizers.Insert(finalizer)
			opr.SetFinalizers(finalizers.List())
			err := r.Client.Update(context.TODO(), opr)
			if err != nil {
				if errors.IsNotFound(err) {
					log.Info("Could not add finalizer to %v: the object was deleted", req)
					return reconcile.Result{}, nil
				} else if errors.IsConflict(err) {
					log.Info("Could not add finalizer to %v due to conflict. Operation will be retried in next reconcile attempt", req)
					return reconcile.Result{}, nil
				}
				log.Error(err, fmt.Sprintf("Failed to update Operator %s with finalizer", opr.Name))
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil
		}
		// find all crs and gen operatorStatus
		if err = r.createOperatorStatus(opr); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *OperatorWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	u := gvk.NewOperatorUnstructured().GetUnstructured()
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1alpha1.OperatorWatcher{}).
		Watches(&source.Kind{Type: &u}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
