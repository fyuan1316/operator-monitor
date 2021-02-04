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
	"github.com/go-logr/logr"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/pkg/status"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	"k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	operatorv1alpha1 "gitlab-ce.alauda.cn/micro-service/operator-monitor/api/v1alpha1"
)

// OperatorStatusReconciler reconciles a OperatorStatus object
type OperatorStatusReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=operator.alauda.io,resources=operatorstatuses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operator.alauda.io,resources=operatorstatuses/status,verbs=get;update;patch

func (r *OperatorStatusReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error
	log := r.Log.WithValues("operatorstatus", req.NamespacedName)
	log.Info(fmt.Sprintf("Starting reconcile loop for %v", req.NamespacedName))
	defer log.Info(fmt.Sprintf("Finish reconcile loop for %v", req.NamespacedName))

	ins := &operatorv1alpha1.OperatorStatus{}
	err = r.Get(context.Background(), req.NamespacedName, ins)
	if err != nil {
		if errors.IsNotFound(err) {
			// CR not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if err := r.updateStatus(ins); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
func (r *OperatorStatusReconciler) updateStatus(ins *operatorv1alpha1.OperatorStatus) error {
	wantedStatus, err := status.CheckComponentsStatus(r.Client, ins)
	if err != nil {
		return err
	}
	currentOprStatus := operatorv1alpha1.OperatorStatus{}
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		err := r.Client.Get(context.Background(), client.ObjectKeyFromObject(ins), &currentOprStatus)
		if err != nil {
			return err
		}
		current := currentOprStatus.DeepCopy()
		if !equality.Semantic.DeepDerivative(wantedStatus, current.Status) {
			current.Status = wantedStatus
			if errUpd := r.Client.Status().Update(context.Background(), current); errUpd != nil {
				return errUpd
			}
		}
		return nil
	})

}

func (r *OperatorStatusReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1alpha1.OperatorStatus{}).
		Complete(r)
}
