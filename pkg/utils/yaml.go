package utils

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

func DecodeYaml(reader io.Reader, out any) error {
	decoder := yaml.NewDecoder(reader)
	return decoder.Decode(out)
}
