package logging

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

// for github.com/yangchenxing/go-map2struct
type WriterFactory struct {
	MapUnmarshaler func(interface{}, map[string]interface{}) error
}

func (factory *WriterFactory) GetInstanceType() reflect.Type {
	return reflect.TypeOf((*Writer)(nil)).Elem()
}

func (factory *WriterFactory) Create(data map[string]interface{}) (interface{}, error) {
	if typeName, ok := data["type"].(string); ok {
		switch typeName {
		case "stderr":
			return os.Stderr, nil
		case "timerotate":
			writer := new(TimeRotateWriter)
			if err := factory.MapUnmarshaler(writer, data); err != nil {
				return nil, err
			} else if err := writer.Initialize(); err != nil {
				return nil, err
			}
			return writer, nil
		case "email":
			writer := new(EMailWriter)
			if err := factory.MapUnmarshaler(writer, data); err != nil {
				return nil, err
			}
			writer.Initialize()
			return writer, nil
		default:
			return nil, fmt.Errorf("unknown write type: %q", typeName)
		}
	}
	return nil, errors.New("missing writer type")
}
