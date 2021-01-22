package v1alpha1

const (
	KindAsm           = "Asm"
	KindFlagger       = "Flagger"
	KindJaeger        = "Jaeger"
	KindIstioOperator = "IstioOperator"
)

type ResourceState string

var ResourceStates = struct {
	NotPresent string
	//Present    string
	NotReady string
	Ready    string
	Healthy  string
}{
	"NotPresent",
	//"Present",
	"NotReady",
	"Ready",
	"Healthy",
}
