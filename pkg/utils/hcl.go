package utils

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func hclFormatter(w io.Writer, data interface{}) error {
	hclFile := hclwrite.NewEmptyFile()
	rootBody := hclFile.Body()

	encodeIntoBody(data, rootBody)

	_, err := hclFile.WriteTo(w)

	return err
}

func encodeIntoBody(data any, rootBody *hclwrite.Body) {
	rv := reflect.ValueOf(data)
	ty := rv.Type()
	if ty.Kind() == reflect.Ptr {
		rv = rv.Elem()
		ty = rv.Type()
	}
	if ty.Kind() != reflect.Struct {
		panic(fmt.Sprintf("value is %s, not struct", ty.Kind()))
	}

	tags := getFieldTags(ty)

	populateBody(rv, ty, tags, rootBody)
}

func EncodeAsBlock(val interface{}, blockType string) *hclwrite.Block {
	rv := reflect.ValueOf(val)
	ty := rv.Type()
	if ty.Kind() == reflect.Ptr {
		rv = rv.Elem()
		ty = rv.Type()
	}
	if ty.Kind() != reflect.Struct {
		panic(fmt.Sprintf("value is %s, not struct", ty.Kind()))
	}

	tags := getFieldTags(ty)
	labels := make([]string, len(tags.Labels))
	for i, lf := range tags.Labels {
		lv := rv.Field(lf.FieldIndex)
		// We just stringify whatever we find. It should always be a string
		// but if not then we'll still do something reasonable.
		labels[i] = fmt.Sprintf("%s", lv.Interface())
	}

	block := hclwrite.NewBlock(blockType, labels)
	populateBody(rv, ty, tags, block.Body())
	return block
}

func populateBody(rv reflect.Value, ty reflect.Type, tags *fieldTags, dst *hclwrite.Body) {
	nameIdxs := make(map[string]int, len(tags.Attributes)+len(tags.Blocks))
	namesOrder := make([]string, 0, len(tags.Attributes)+len(tags.Blocks))
	for n, i := range tags.Attributes {
		nameIdxs[n] = i
		namesOrder = append(namesOrder, n)
	}
	for n, i := range tags.Blocks {
		nameIdxs[n] = i
		namesOrder = append(namesOrder, n)
	}
	sort.SliceStable(namesOrder, func(i, j int) bool {
		ni, nj := namesOrder[i], namesOrder[j]
		return nameIdxs[ni] < nameIdxs[nj]
	})

	dst.Clear()

	prevWasBlock := false
	for _, name := range namesOrder {
		fieldIdx := nameIdxs[name]
		field := ty.Field(fieldIdx)
		fieldTy := field.Type
		fieldVal := rv.Field(fieldIdx)

		if fieldTy.Kind() == reflect.Ptr {
			fieldTy = fieldTy.Elem()
			fieldVal = fieldVal.Elem()
		}

		if _, isAttr := tags.Attributes[name]; isAttr {

			if !fieldVal.IsValid() {
				continue // ignore (field value is nil pointer)
			}
			if fieldTy.Kind() == reflect.Ptr && fieldVal.IsNil() {
				continue // ignore
			}
			if prevWasBlock {
				dst.AppendNewline()
				prevWasBlock = false
			}

			valTy, err := gocty.ImpliedType(fieldVal.Interface())
			if err != nil {
				panic(fmt.Sprintf("cannot encode %T as HCL expression: %s", fieldVal.Interface(), err))
			}

			val, err := gocty.ToCtyValue(fieldVal.Interface(), valTy)
			if err != nil {
				// This should never happen, since we should always be able
				// to decode into the implied type.
				panic(fmt.Sprintf("failed to encode %T as %#v: %s", fieldVal.Interface(), valTy, err))
			}

			dst.SetAttributeValue(name, val)
		} else { // must be a block
			elemTy := fieldTy
			isSeq := false
			if elemTy.Kind() == reflect.Slice || elemTy.Kind() == reflect.Array {
				isSeq = true
				elemTy = elemTy.Elem()
			}

			prevWasBlock = false

			if isSeq {
				l := fieldVal.Len()
				for i := 0; i < l; i++ {
					elemVal := fieldVal.Index(i)
					if !elemVal.IsValid() {
						continue // ignore (elem value is nil pointer)
					}
					if elemTy.Kind() == reflect.Ptr && elemVal.IsNil() {
						continue // ignore
					}
					block := EncodeAsBlock(elemVal.Interface(), name)
					if !prevWasBlock {
						dst.AppendNewline()
						prevWasBlock = true
					}
					dst.AppendBlock(block)
				}
			} else {
				if !fieldVal.IsValid() {
					continue // ignore (field value is nil pointer)
				}
				if elemTy.Kind() == reflect.Ptr && fieldVal.IsNil() {
					continue // ignore
				}
				block := EncodeAsBlock(fieldVal.Interface(), name)
				if !prevWasBlock {
					dst.AppendNewline()
					prevWasBlock = true
				}
				dst.AppendBlock(block)
			}
		}
	}

	// Encode Body
	if tags.Body != nil {
		bodyFieldIdx := *tags.Body
		field := ty.Field(bodyFieldIdx)
		fieldTy := field.Type
		fieldVal := rv.Field(bodyFieldIdx)
		if fieldTy.Kind() == reflect.Ptr {
			fieldTy = fieldTy.Elem()
			fieldVal = fieldVal.Elem()
		}

		if fieldTy.Kind() == reflect.Map {
			formatMapBody(fieldVal.Interface().(map[string]any), dst)
		}
	}
}

func formatMapBody(v map[string]any, dst *hclwrite.Body) {
	backendConfigSortedKeys := StringKeysFromMap(v)
	for _, name := range backendConfigSortedKeys {
		v := v[name]
		formatField(dst, name, v)

	}
}

func formatField(blockBody *hclwrite.Body, name string, v any) {
	if v == nil {
		blockBody.SetAttributeValue(name, cty.NilVal)
	} else if i, ok := v.(string); ok {
		blockBody.SetAttributeValue(name, cty.StringVal(i))
	} else if i, ok := v.(bool); ok {
		blockBody.SetAttributeValue(name, cty.BoolVal(i))
	} else if i, ok := v.(int64); ok {
		blockBody.SetAttributeValue(name, cty.NumberIntVal(i))
	} else if i, ok := v.(uint64); ok {
		blockBody.SetAttributeValue(name, cty.NumberUIntVal(i))
	} else if i, ok := v.(float64); ok {
		blockBody.SetAttributeValue(name, cty.NumberFloatVal(i))
	}
}

type fieldTags struct {
	Attributes map[string]int
	Blocks     map[string]int
	Labels     []labelField
	Remain     *int
	Body       *int
	Optional   map[string]bool
}
type labelField struct {
	FieldIndex int
	Name       string
}

func getFieldTags(ty reflect.Type) *fieldTags {
	ret := &fieldTags{
		Attributes: map[string]int{},
		Blocks:     map[string]int{},
		Optional:   map[string]bool{},
	}

	ct := ty.NumField()
	for i := 0; i < ct; i++ {
		field := ty.Field(i)
		tag := field.Tag.Get("hcle")
		if tag == "" {
			continue
		}

		comma := strings.Index(tag, ",")
		var name, kind string
		if comma != -1 {
			name = tag[:comma]
			kind = tag[comma+1:]
		} else {
			name = tag
			kind = "attr"
		}

		switch kind {
		case "attr":
			ret.Attributes[name] = i
		case "block":
			ret.Blocks[name] = i
		case "label":
			ret.Labels = append(ret.Labels, labelField{
				FieldIndex: i,
				Name:       name,
			})
		case "body":
			if ret.Body != nil {
				panic("only one 'body' tag is permitted")
			}
			idx := i // copy, because this loop will continue assigning to i
			ret.Body = &idx
		default:
			panic(fmt.Sprintf("invalid hcl field tag kind %q on %s %q", kind, field.Type.String(), field.Name))
		}
	}

	return ret
}
