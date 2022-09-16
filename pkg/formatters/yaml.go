package formatters

import (
	"bufio"
	"io"

	"gopkg.in/yaml.v2"
)

// ConvertToYAML converts the provided value to a YAML string
func ConvertToYAML(data any) (string, error) {

	y, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(y), nil
}

func yamlFormatter(w io.Writer, data interface{}) error {

	bw := bufio.NewWriter(w)

	encoder := yaml.NewEncoder(bw)

	defer encoder.Close()

	if err := encoder.Encode(data); err != nil {
		return err
	}

	if err := bw.Flush(); err != nil {
		return err
	}

	return nil
}
