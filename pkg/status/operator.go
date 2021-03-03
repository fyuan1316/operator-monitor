package status

import (
	pkgerrors "github.com/pkg/errors"
	"strings"
)

func ValidOperatorName(name string) (string, string, error) {
	nameTuple := strings.Split(name, ".")
	if len(nameTuple) != 2 {
		return "", "", pkgerrors.New("not a valid operator name")
	}
	return nameTuple[0], nameTuple[1], nil
}
