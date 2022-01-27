package mapping

import (
	"encoding/xml"
	"fmt"
	"reflect"
)

type XmlMap struct {
	m map[string]interface{}
}

func (m *XmlMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	m.m = make(map[string]interface{})

	for {

	}
	return nil
}

func XMLUnmarshal(data []byte, v interface{}) error {
	if err := xml.Unmarshal(data, v); err == nil {
		return nil
	}

	// 此时默认的xml转换器已经无法满足需求
	// 判断v的值是否为map[string]interface{} point
	xmlmap, ok := v.(map[string]interface{})
	if !ok {
		return fmt.Errorf("can't unmarshal xml to %s", reflect.TypeOf(v).String())
	}

	xm := &XmlMap{m: xmlmap}

	if err := xml.Unmarshal(data, xm); err == nil {
		return err
	}
	return nil
}
