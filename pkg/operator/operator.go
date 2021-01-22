package operator

import (
	"context"
	"fmt"
	"fyuan1316/operator-monitor/api/v1alpha1"
	"fyuan1316/operator-monitor/pkg/util"
	v1 "github.com/operator-framework/api/pkg/operators/v1"
	pkgerrors "github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

var (
	oprLog = ctrl.Log.WithName("operatorWatcher")
)

type Observable interface {
	CRs() ([]v1alpha1.OperatorStatus, error)
	//SetNamespace(string)
}

type Observer struct {
	Client   client.Client
	Operator corev1.TypedLocalObjectReference
}

//
//func (o *Observer) SetNamespace(ns string) {
//	o.Operator.Namespace = ns
//}

func CrsForOperator(opr *v1.Operator, client client.Client) ([]v1alpha1.OperatorStatus, error) {
	name, _, _ := ValidOperatorName(opr.Name)

	var impl Observable
	switch name {
	case util.GetAsmOperatorName():
		impl = NewAsmImpl(opr, client)
	case util.GetFlaggerOperatorName():
		impl = NewFlaggerImpl(opr, client)
	case util.GetIstioOperatorName():
		impl = NewIstioImpl(opr, client)
	case util.GetJaegerOperatorName():
		impl = NewJaegerImpl(opr, client)
	default:
		return nil, pkgerrors.New(fmt.Sprintf("not supported operator %s", name))
	}
	return impl.CRs()
}

func ValidOperatorName(name string) (string, string, error) {
	nameTuple := strings.Split(name, ".")
	if len(nameTuple) != 2 {
		return "", "", pkgerrors.New("not a valid operator name")
	}
	return nameTuple[0], nameTuple[1], nil
}

func GetUsList(kind, apiversion, group string) unstructured.UnstructuredList {
	uList := unstructured.UnstructuredList{}
	uList.SetKind(kind)
	uList.SetAPIVersion(apiversion)
	uList.SetGroupVersionKind(
		schema.GroupVersionKind{
			Group:   group,
			Kind:    kind,
			Version: apiversion},
	)
	return uList
}

func getStatusName(name, ns string) string {
	if ns == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", name, ns)
}
func GenStatus(opr corev1.TypedLocalObjectReference, uList unstructured.UnstructuredList, getInstalledNs func(unstructured.Unstructured) ([]string, error)) ([]v1alpha1.OperatorStatus, error) {
	var statuses []v1alpha1.OperatorStatus
	for _, us := range uList.Items {
		status := v1alpha1.OperatorStatus{}
		status.Name = getStatusName(us.GetName(), us.GetNamespace())
		status.Spec.Operator = opr
		group := us.GroupVersionKind().Group
		status.Spec.CR = v1alpha1.NameSpacedTypedObjectReference{
			TypedLocalObjectReference: corev1.TypedLocalObjectReference{
				Kind:     us.GetKind(),
				APIGroup: &group,
				Name:     us.GetName(),
			},
			Namespace: us.GetNamespace(),
		}
		installTo, err := getInstalledNs(us)
		if err != nil {
			return statuses, err
		}
		status.Spec.InstalledNamespace = installTo
		statuses = append(statuses, status)
	}
	return statuses, nil
}

type WatchNamespaceType string

var WatchNamespaces = struct {
	All   WatchNamespaceType
	One   WatchNamespaceType
	Multi WatchNamespaceType
}{"", "one", "multi"}

type WatchNamespace struct {
	sets.String
}

func NewWatchNamespace() WatchNamespace {
	return WatchNamespace{
		sets.String{},
	}
}

func (w WatchNamespace) getType() WatchNamespaceType {
	if w.Len() > 1 {
		return WatchNamespaces.Multi
	} else if w.Len() == 1 && w.List()[0] != "" {
		return WatchNamespaces.One
	} else {
		return WatchNamespaces.All
	}

}
func (o Observer) GetWatchNamespace(f func(*appsv1.Deployment) WatchNamespace) (watchNs WatchNamespace, err error) {
	deploy := appsv1.Deployment{}
	name, ns, err := ValidOperatorName(o.Operator.Name)
	key := client.ObjectKey{Name: name, Namespace: ns}
	err = o.Client.Get(context.Background(), key, &deploy)
	if err != nil {
		return
	}
	watchNs = f(&deploy)

	return
}

var CheckWatchNamespaceInArgs = func(d *appsv1.Deployment) WatchNamespace {
	var watchedNamespace string
	for _, c := range d.Spec.Template.Spec.Containers {
		for _, arg := range c.Args {
			if strings.Contains(arg, "namespace") {
				kvPair := strings.Split(arg, "=")
				if len(kvPair) == 2 {
					watchedNamespace = kvPair[1]
				}
				break
			}
		}
	}
	ns := NewWatchNamespace()
	if watchedNamespace != "" {
		ns.Insert(strings.Split(watchedNamespace, ",")...)
	}
	return ns
}
var CheckWatchNamespaceInEnv = func(d *appsv1.Deployment) WatchNamespace {
	var watchedNamespace string
	for _, c := range d.Spec.Template.Spec.Containers {
		for _, env := range c.Env {
			if env.Name == "WATCH_NAMESPACE" {
				watchedNamespace = env.Value
				break
			}
		}
	}
	ns := NewWatchNamespace()
	if watchedNamespace != "" {
		ns.Insert(strings.Split(watchedNamespace, ",")...)
	}
	return ns
}

func (o Observer) GetValidCRs(watchNs WatchNamespace, uList unstructured.UnstructuredList) (unstructured.UnstructuredList, error) {
	opts := client.ListOptions{}

	if watchNs.getType() == WatchNamespaces.One {
		opts.Namespace = watchNs.List()[0]
	}

	var err = o.Client.List(context.Background(), &uList, &opts)
	if err != nil {
		return uList, err
	}
	var uns []unstructured.Unstructured
	if watchNs.getType() == WatchNamespaces.Multi {
		for _, cr := range uList.Items {
			if watchNs.Has(cr.GetNamespace()) {
				uns = append(uns, cr)
			}
		}
		uList.Items = uns
	}
	return uList, nil
}
