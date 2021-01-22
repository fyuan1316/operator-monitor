package status

const (
	//install.operator.istio.io/owning-resource-namespace=istio-system,
	//install.operator.istio.io/owning-resource=example-istiocontrolplane
	IstioOwningResourceKey          = "install.operator.istio.io/owning-resource"
	IstioOwningResourceNamespaceKey = "install.operator.istio.io/owning-resource-namespace"

	// jaeger控制平面的k8s资源会创建到与cr相同的namespace中
	//app.kubernetes.io/instance=jaeger-prod,
	//app.kubernetes.io/managed-by=jaeger-operator
	JaegerInstanceKey = "app.kubernetes.io/instance"
	JaegerOperatorKey = "app.kubernetes.io/managed-by"

	AsmOwningResourceKey          = "asm.operator.alauda.io/owning-resource"
	AsmOwningResourceNamespaceKey = "asm.operator.alauda.io/owning-resource-namespace"

	FlaggerOwningResourceKey          = "flagger.operator.alauda.io/owning-resource"
	FlaggerOwningResourceNamespaceKey = "flagger.operator.alauda.io/owning-resource-namespace"
)
