package xml

import (
	"encoding/xml"
	"fmt"
	"github.com/the-prophet1/datamapper/mapping/json"
	"reflect"

	"github.com/clbanning/mxj/v2"
)

func Unmarshal(data []byte, v interface{}) error {
	if err := xml.Unmarshal(data, v); err == nil {
		return nil
	}

	// 此时默认的xml转换器已经无法满足需求
	// 判断v的值是否为*map[string]interface{}
	m, ok := v.(*map[string]interface{})
	if !ok {
		return fmt.Errorf("can't unmarshal xml to %s", reflect.TypeOf(v).String())
	}

	mv, err := mxj.NewMapXml(data) // unmarshal
	if err != nil {
		return err
	}

	*m = mv
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	if data, err := xml.Marshal(v); err == nil {
		return data, err
	}

	// 此时默认的xml转换器已经无法满足需求
	// 判断v的值是否为map[string]interface{}
	switch v := v.(type) {
	case map[string]interface{}:
		res := mxj.Map(v)
		return res.Xml()
	case map[string]*interface{}:
		data, _ := json.Marshal(v)
		var tmpMap map[string]interface{}
		_ = json.Unmarshal(data, &tmpMap)
		res := mxj.Map(tmpMap)
		return res.Xml()
	default:
		return nil, fmt.Errorf("can't marshal xml to %s", reflect.TypeOf(v).String())
	}

}
