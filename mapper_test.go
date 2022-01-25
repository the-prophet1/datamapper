package datamapper

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

type TestStruct struct {
	ID            string
	Specification []byte
	Input         string
	Output        string
}

func Spec(path string) []byte {
	file, _ := os.Open(path)
	data, _ := ioutil.ReadAll(file)
	return data
}

var mapperTest = []TestStruct{
	{
		"test1",
		Spec("./test/test1.yaml"),
		`{"id":"test1","description":"描述映射test1","data":{"voltage":220,"current": 10,"power":2200}}`,
		`{"id":"test1","description":"描述映射test1","va":{"A":10,"P":2200,"V":220}}`,
	},
	{
		"test2",
		Spec("./test/test2.yaml"),
		`{"id":"test2","description":"描述映射test2","data":[{"voltage":220,"current":10,"power":2200},{"voltage":36,"current":1,"power":36}]}`,
		`{"id":"test2","description":"描述映射test2","va":[{"V":220,"A":10,"P":2200},{"V":36,"A":1,"P":36}]}`,
	},
	{
		"test3",
		Spec("./test/test3.yaml"),
		`{"id":"test2","description":"描述映射test2","voltage":220,"ampere":10,"power":2200}`,
		`{"id":"test2","description":"描述映射test2","V":220,"A":10,"P":2200}`,
	},
	{
		"test4",
		Spec("./test/test4.yaml"),
		`{"msg":"成功","headers":{"qos":1,"oneofCase":5,"token":"kCBQLBlvOp+9fOsRWKN3VD6V5DSNgnpNnU2U1M6cOYg="},"code":"SUCCESS","fromMessageId":"","messageId":"f09856be6ae947a79ca21d24a33e7239","properties":[{"val":"7.00","name":"CPU使用率","time":"1642757405418","desc":{"unit":"%","plugName":"sysinfo","source":"","type":"DOUBLE","group":"","quality":0}}],"timestamp":1642757411915}`,
		``,
	},
}

func TestGenerateDataDefine(t *testing.T) {

	for _, testStruct := range mapperTest {
		dataDefine, err := GenerateDataDefine(testStruct.Specification)
		if err != nil {
			t.Fatal(err)
		}

		output, err := dataDefine.To([]byte(testStruct.Input))
		if err != nil {
			t.Fatal(err)
		}

		var m1 map[string]interface{}
		var m2 map[string]interface{}
		if err := json.Unmarshal(output, &m1); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(testStruct.Output), &m2); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(m1, m2) {
			t.Error(testStruct.ID)
		}
	}
}
