package mapping

import (
	"encoding/xml"
	"fmt"
	"reflect"

	"github.com/clbanning/mxj/v2"
)



func XMLUnmarshal(data []byte, v interface{}) error {
	if err := xml.Unmarshal(data, v); err == nil {
		return nil
	}

	// 此时默认的xml转换器已经无法满足需求
	// 判断v的值是否为map[string]interface{} point
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
