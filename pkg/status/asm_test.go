package status

import (
	"context"
	"fmt"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/api/v1alpha1"
	operatorv1alpha1 "gitlab-ce.alauda.cn/micro-service/operator-monitor/api/v1alpha1"
	"gitlab-ce.alauda.cn/micro-service/operator-monitor/pkg/mock"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

var cb = fake.ClientBuilder{}
var k8sClient = getClient()

func getClient() client.Client {
	cl := cb.WithScheme(func() *runtime.Scheme {
		var sch = runtime.NewScheme()
		_ = appsv1.AddToScheme(sch)
		_ = operatorv1alpha1.AddToScheme(sch)
		return sch
	}()).WithObjects(mock.GetObjects()...).Build()
	return cl
}

func TestStatusGetter_Status(t *testing.T) {
	type fields struct {
		OperatorStatusGet OperatorStatusGetter
	}
	tests := []struct {
		name    string
		fields  fields
		want    v1alpha1.OperatorStatusStatus
		wantErr bool
	}{
		{
			name: "get asm status",
			fields: fields{
				OperatorStatusGet: NewAsmStatusGetter(OperatorStatusGet{
					Client: k8sClient,
					OperatorStatus: func() *operatorv1alpha1.OperatorStatus {
						c := mock.NewAsmCase()
						return c.OperatorStatus
					}(),
				}),
			},
			want: func() v1alpha1.OperatorStatusStatus {
				c := mock.NewAsmCase()
				return *c.ExpectStatus
			}(),
			wantErr: false,
		},
		{
			name: "get flagger status",
			fields: fields{
				OperatorStatusGet: NewFlaggerStatusGetter(OperatorStatusGet{
					Client: k8sClient,
					OperatorStatus: func() *operatorv1alpha1.OperatorStatus {
						c := mock.NewFlaggerCase()
						return c.OperatorStatus
					}(),
				}),
			},
			want: func() v1alpha1.OperatorStatusStatus {
				c := mock.NewFlaggerCase()
				return *c.ExpectStatus
			}(),
			wantErr: false,
		},
		{
			name: "get jaeger status",
			fields: fields{
				OperatorStatusGet: NewJaegerStatusGetter(OperatorStatusGet{
					Client: k8sClient,
					OperatorStatus: func() *operatorv1alpha1.OperatorStatus {
						c := mock.NewJaegerCase()
						return c.OperatorStatus
					}(),
				}),
			},
			want: func() v1alpha1.OperatorStatusStatus {
				c := mock.NewJaegerCase()
				return *c.ExpectStatus
			}(),
			wantErr: false,
		},
		{
			name: "get istio status",
			fields: fields{
				OperatorStatusGet: NewIstioStatusGetter(OperatorStatusGet{
					Client: k8sClient,
					OperatorStatus: func() *operatorv1alpha1.OperatorStatus {
						c := mock.NewIstioCase()
						return c.OperatorStatus
					}(),
				}),
			},
			want: func() v1alpha1.OperatorStatusStatus {
				c := mock.NewIstioCase()
				return *c.ExpectStatus
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep := &appsv1.DeploymentList{}
			k8sClient.List(context.TODO(), dep)
			fmt.Println("")
			got, err := tt.fields.OperatorStatusGet.Status()
			if (err != nil) != tt.wantErr {
				t.Errorf("Status() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !ignoreSchemeEqual(got, tt.want) {
				t.Errorf(diff.ObjectDiff(tt.want, got))
			}
		})
	}
}
func ignoreSchemeEqual(got, want operatorv1alpha1.OperatorStatusStatus) bool {
	for k, v := range want.ComponentStatus {
		empty := ""
		v.APIGroup = &empty
		v.Kind = empty
		want.ComponentStatus[k] = v
	}

	return reflect.DeepEqual(got, want)
}
