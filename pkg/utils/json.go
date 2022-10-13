package utils

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

var (
	jsonator jsoniter.API
)

func init() {
	jsonator = jsoniter.Config{
		EscapeHTML:                    true,
		ObjectFieldMustBeSimpleString: false,
		SortMapKeys:                   true,
		ValidateJsonRawMessage:        true,
		IndentionStep:                 3,
	}.Froze()

	jsonator = jsoniter.ConfigCompatibleWithStandardLibrary
}

func jsonFormatter(w io.Writer, data interface{}) error {
	stream := jsonator.BorrowStream(w)
	defer jsonator.ReturnStream(stream)

	stream.WriteVal(data)
	if stream.Error != nil {
		return stream.Error
	}

	if err := stream.Flush(); err != nil {
		return err
	}

	return nil
}

// ConvertToJSONFast converts the provided value to a JSON-encoded string using 'ConfigFastest' config and json.Marshal without indents
func ConvertToJSONFast(data any) (string, error) {
	var json = jsoniter.Config{
		EscapeHTML:                    false,
		MarshalFloatWith6Digits:       true,
		ObjectFieldMustBeSimpleString: true,
		SortMapKeys:                   true,
		ValidateJsonRawMessage:        true,
	}

	j, err := json.Froze().MarshalToString(data)
	if err != nil {
		return "", err
	}
	return j, nil
}
