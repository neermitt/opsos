package utils

import (
	"fmt"
	"io"
)

type Formatter func(w io.Writer, any interface{}) error

func GetFormatter(format string) Formatter {
	switch format {
	case "json":
		return jsonFormatter
	case "yaml":
		return yamlFormatter
	case "hcl":
		return hclFormatter
	default:
		return func(w io.Writer, any interface{}) error {
			return fmt.Errorf("invalid format type: %s", format)
		}
	}
}
