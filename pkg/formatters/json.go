package formatters

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
