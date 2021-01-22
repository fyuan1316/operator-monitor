package util

import (
	"encoding/json"
	pkgerrors "github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

func JsonConvert(from interface{}, to interface{}) error {
	var data []byte
	var err error
	if data, err = json.Marshal(from); err != nil {
		return pkgerrors.WithStack(err)
	}
	return pkgerrors.WithStack(json.Unmarshal(data, to))
}
func StringConvert(from string, to interface{}) error {
	return pkgerrors.WithStack(yaml.Unmarshal([]byte(from), to))
}
