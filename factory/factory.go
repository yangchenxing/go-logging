package factory

import (
	"errors"
	"fmt"
	"github.com/yangchenxing/go-logging"
	"github.com/yangchenxing/go-map2struct"
	"os"
	"reflect"
)

func init() {
	map2struct.RegisterFactory(new(WriterFactory))
}

// WriterFactory is used by github.com/yangchenxing/go-map2struct
type WriterFactory struct {
	MapUnmarshaler func(interface{}, map[string]interface{}) error
}

// GetInstanceType return type of Writer
func (factory *WriterFactory) GetInstanceType() reflect.Type {
	return reflect.TypeOf((*logging.Writer)(nil)).Elem()
}

// Create create Writer instance with a map[string]interface{}
func (factory *WriterFactory) Create(data map[string]interface{}) (interface{}, error) {
	if typeName, ok := data["type"].(string); ok {
		switch typeName {
		case "stderr":
			return os.Stderr, nil
		case "stdout":
			return os.Stdout, nil
		case "timerotate":
			writer := new(logging.TimeRotateWriter)
			if err := factory.MapUnmarshaler(writer, data); err != nil {
				return nil, err
			} else if err := writer.Initialize(); err != nil {
				return nil, err
			}
			return writer, nil
		case "email":
			writer := new(logging.EmailWriter)
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
