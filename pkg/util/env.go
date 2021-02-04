package util

import "os"

const (
	DefaultAsmOperatorName     = "asm-operator"
	DefaultFlaggerOperatorName = "flagger-operator"
	DefaultIstioOperatorName   = "istio-operator"
	DefaultJaegerOperatorName  = "jaeger-operator"

	AsmOperatorNameKey     = "ASM_OPERATOR_NAME"
	FlaggerOperatorNameKey = "FLAGGEROPERATOR_NAME"
	IstioOperatorNameKey   = "ISTIO_OPERATOR_NAME"
	JaegerOperatorNameKey  = "JAEGER_OPERATOR_NAME"
)

func GetAsmOperatorName() string {
	return findNameInEnv(AsmOperatorNameKey, DefaultAsmOperatorName)
}

func GetFlaggerOperatorName() string {
	return findNameInEnv(FlaggerOperatorNameKey, DefaultFlaggerOperatorName)
}
func GetIstioOperatorName() string {
	return findNameInEnv(IstioOperatorNameKey, DefaultIstioOperatorName)
}
func GetJaegerOperatorName() string {
	return findNameInEnv(JaegerOperatorNameKey, DefaultJaegerOperatorName)
}
func findNameInEnv(key string, fallback string) string {
	if v, found := os.LookupEnv(key); found && v != "" {
		return v
	}
	return fallback
}
