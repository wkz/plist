package plist

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	encDecl = `<?xml version="1.0" encoding="UTF-8"?>`
	dtdDecl = `<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">`
)

type indentBuffer struct {
	bytes.Buffer
	Prefix, Indent, Newline string
}

func (ib *indentBuffer) Writeln(d int, ln string) {
	ib.WriteString(
		ib.Prefix +
		strings.Repeat(ib.Indent, d) +
		ln +
		ib.Newline)
}

func Marshal(v interface{}) ([]byte, error) {
	ib := &indentBuffer{}
	err := marshalHead(v, ib)
	return ib.Bytes(), err
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	ib := &indentBuffer{Prefix: prefix, Indent: indent, Newline: "\n"}
	err := marshalHead(v, ib)
	return ib.Bytes(), err
}

func marshalHead(v interface{}, ib *indentBuffer) error {
	ib.Writeln(0, encDecl)
	ib.Writeln(0, dtdDecl)
	ib.Writeln(0, "<plist version=\"1.0\">")
	err := marshal(1, reflect.ValueOf(v), ib)
	if err != nil {
		return err
	}
	ib.Writeln(0, "</plist>")
	return nil
}

func marshal (d int, v reflect.Value, ib *indentBuffer) error {
	t := v.Type()

	// auto-dereference pointers
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.ConvertibleTo(reflect.TypeOf(time.Time{})) {
		tm := v.Interface().(time.Time)
		ib.Writeln(d, "<date>"+tm.Format(time.RFC3339)+"</date>")
		return nil
	}

	switch t.Kind() {
	case reflect.Bool:
		if v.Bool() {
			ib.Writeln(d, "<true/>")
		} else {
			ib.Writeln(d, "<false/>")
		}

	case reflect.Int:    fallthrough
	case reflect.Int8:   fallthrough
	case reflect.Int16:  fallthrough
	case reflect.Int32:  fallthrough
	case reflect.Int64:
		ib.Writeln(d, fmt.Sprintf("<integer>%d</integer>", v.Int()))

	case reflect.Uint:   fallthrough
	case reflect.Uint8:  fallthrough
	case reflect.Uint16: fallthrough
	case reflect.Uint32: fallthrough
	case reflect.Uint64:
		ib.Writeln(d, fmt.Sprintf("<integer>%d</integer>", v.Uint()))
		
	case reflect.Float32: fallthrough
	case reflect.Float64:
		ib.Writeln(d, fmt.Sprintf("<real>%f</real>", v.Float()))
	
	case reflect.Array:
		v = v.Slice(0, v.Len())
		fallthrough
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			bytes := v.Interface().([]byte)
			ib.Writeln(d, 
				"<data>" + 
				base64.StdEncoding.EncodeToString(bytes) +
				"</data>")
			return nil
		}
		
		if v.Len() == 0 {
			ib.Writeln(d, "<array/>")
			return nil
		}

		ib.Writeln(d, "<array>")
		for i := 0; i < v.Len(); i++ {
			err := marshal(d+1, v.Index(i), ib)
			if err != nil {
				return err
			}
		}
		ib.Writeln(d, "</array>")

	case reflect.String:
		ib.Writeln(d, "<string>" + v.String() + "</string>")

	case reflect.Struct:
		ib.Writeln(d, "<dict>")
		for i := 0; i < t.NumField(); i++ {
			key, ok := fieldName(t.Field(i))
			if !ok {
				continue
			}

			ib.Writeln(d+1, "<key>"+key+"</key>")
			err := marshal(d+1, v.Field(i), ib)
			if err != nil {
				return err
			}
		}
		ib.Writeln(d, "</dict>")

	case reflect.Map:
		ib.Writeln(d, "<dict>")
		for _, k := range v.MapKeys() {
			ib.Writeln(d + 1, "<key>" + k.String() + "</key>")
			err := marshal(d + 1, v.MapIndex(k), ib)
			if err != nil {
				return err
			}
		}
		ib.Writeln(d, "</dict>")

	default:
		return errors.New(fmt.Sprintf("unable to marshal \"%v\" of type \"%v\"", v, t))

	}
	return nil
}

func fieldName(f reflect.StructField) (string, bool) {
	if f.PkgPath != "" {
		return "", false
	}

	tn := f.Tag.Get("plist")
	if len(tn) > 0 {
		return tn, tn != "-"
	}

	return f.Name, true
}
