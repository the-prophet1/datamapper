package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var logger *log.Logger

func init() {
	logger = log.New()
	logger.SetReportCaller(true)
}

type Specification interface {
	Get(name string) ([]byte, error)
}

func NewSpecification() Specification {
	return &FileReader{}
}

type FileReader struct{}

func (FileReader) Get(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	res, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type DataSpec struct {
	Type     string `yaml:"type"`
	TypeRef  string `yaml:"typeRef"`
	Multiple string `yaml:"multiple"`
}

type ComplexDefine map[string]*DataSpec

type DataDefine struct {
	SourceType string                   `yaml:"sourceType"`
	TargetType string                   `yaml:"targetType"`
	Source     ComplexDefine            `yaml:"source"`
	Target     ComplexDefine            `yaml:"target"`
	Mapper     map[string]string        `yaml:"mapper"`
	Complex    map[string]ComplexDefine `yaml:"complex"`
}

func GenerateDataDefine(data []byte) (*DataDefine, error) {
	var define DataDefine
	if err := yaml.Unmarshal(data, &define); err != nil {
		return nil, err
	}
	return &define, nil
}

func (d *DataSpec) IsNumber() bool {
	if d.TypeRef == "number" {
		return true
	}
	return false
}

func (d *DataSpec) IsString() bool {
	if d.TypeRef == "string" {
		return true
	}
	return false
}

func (d *DataSpec) IsComplex() bool {
	if d.Type == "complex" {
		return true
	}
	return false
}

func (d *DataSpec) IsSimple() bool {
	if d.Type == "simple" {
		return true
	}
	return false
}

func (d *DataSpec) IsArray() bool {
	if d.Multiple == "true" {
		return true
	}
	return false
}

// To 将输入数据转换为源数据，再将源数据转换为目标数据
func (d *DataDefine) To(input []byte) ([]byte, error) {
	if d.SourceType != "json" {
		panic("not implemented")
	}

	// 获取数据数据的map
	var inputMap map[string]interface{}
	if err := json.Unmarshal(input, &inputMap); err != nil {
		return nil, err
	}

	sourceMap := d.ParseSource(d.Source, inputMap)
	targetMap := d.generateMap(d.Target)
	d.mapping(sourceMap,targetMap)

	logger.Debug(sourceMap)
	logger.Debug(targetMap)
	return nil, nil
}

func (d *DataDefine) mapping(sourceMap map[string]interface{}, targetMap map[string]*interface{}) {

	for source, target := range d.Mapper {
		sourcePaths := strings.Split(source, ".")
		sourceData := sourceMap[sourcePaths[0]]

		for i := 1; i < len(sourcePaths); i++ {
			switch sourceData.(type) {
			case map[string]interface{}:
				sourceData = sourceData.(map[string]interface{})[sourcePaths[i]]
			}
		}

		targetPaths := strings.Split(target, ".")
		targetData := targetMap[targetPaths[0]]

		for i := 1; i < len(targetPaths); i++ {
			switch (*targetData).(type) {
			case map[string]*interface{}:
				var val interface{} = (*targetData).(map[string]*interface{})
				targetData = &val
			}
		}
		*targetData = sourceData
	}
}


//
func (d *DataDefine) generateMap(complexDefine ComplexDefine) map[string]*interface{} {
	res := make(map[string]*interface{})
	for key, def := range complexDefine {
		if def.IsArray() {
			if def.IsComplex() {
				val := make([]map[string]*interface{}, 0)
				val = append(val, d.generateMap(d.Complex[def.TypeRef]))
				var vals interface{} = val
				res[key] = &vals
			}
			if def.IsSimple() {
				if def.IsNumber() {
					val := make([]float64, 0)
					val = append(val, float64(0))
					var vals interface{} = val
					res[key] = &vals
				}
				if def.IsString() {
					val := make([]string, 0)
					val = append(val, "")
					var vals interface{} = val
					res[key] = &vals
				}
			}
		} else {
			if def.IsComplex() {
				var val interface{}  = d.generateMap(d.Complex[def.TypeRef])
				res[key] = &val
			}
			if def.IsSimple() {
				if def.IsNumber() {
					var val interface{} = float64(0)
					res[key] = &val
				}
				if def.IsString() {
					var val interface{} = ""
					res[key] = &val
				}
			}
		}
	}
	return res
}

//ParseSource 根据数据定义将input转换为由ComplexDefine定义的output
func (d *DataDefine) ParseSource(complexDefine ComplexDefine, inputMap map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})

	// 从复合类型定义中获取各个字段的定义
	for key, def := range complexDefine {
		// 从复合类型的名称取出inputMap的数据
		inValue, ok := inputMap[key]
		if !ok { //不存在则跳过
			continue
		}
		//判断所需是否为数组
		if def.IsArray() {
			// 如果input也是数组
			if typ := reflect.TypeOf(inValue); typ.Kind() == reflect.Slice {
				if def.IsComplex() {
					switch inValue.(type) {
					case []map[string]interface{}:
						sliMap := inValue.([]map[string]interface{})
						slim := make([]map[string]interface{}, 0)
						for _, m := range sliMap {
							e := d.ParseSource(d.Complex[key], m)
							slim = append(slim, e)
						}
						res[key] = slim
					default:
						logger.Warn()
					}
				}
				if def.IsSimple() { // 此时数据类型为简单类型
					// 将[]interface{}转换为对应的类型切片
					is, ok := inValue.([]interface{})
					if ok && len(is) > 0 {
						elemTyp := reflect.TypeOf(is[0]).Kind()

						if elemTyp == reflect.String {
							strs := make([]string, 0)
							for _, i := range is {
								strs = append(strs, i.(string))
							}
							inValue = strs
						} else if elemTyp == reflect.Float64 {
							fs := make([]float64, 0)
							for _, i := range is {
								fs = append(fs, i.(float64))
							}
							inValue = fs
						}
					}

					switch inValue.(type) {
					case []string:
						if def.IsString() { // 类型匹配，直接赋值
							res[key] = inValue
						}
						if def.IsNumber() {
							// 类型不匹配，尝试将string转为float
							slif := make([]float64, 0)
							for _, s := range inValue.([]string) {
								f, err := strconv.ParseFloat(s, 64)
								if err != nil {
									log.Warn()
								}
								slif = append(slif, f)
							}
							res[key] = slif
						}
					case []float64:
						if def.IsNumber() {
							res[key] = inValue
						}
						if def.IsString() {
							// 类型不匹配，尝试将float转为string
							slis := make([]string, 0)
							for _, f := range inValue.([]float64) {
								s := strconv.FormatFloat(f, 'f', -1, 64)
								slis = append(slis, s)
							}
							res[key] = slis
						}
					default:
						logger.Warn("")
					}
				}
			} else {
				// 此时输入为对象或简单数据类型
				if def.IsComplex() {
					switch inValue.(type) {
					case map[string]interface{}:
						m := inValue.(map[string]interface{})
						c := d.ParseSource(d.Complex[key], m)
						sli := make([]map[string]interface{}, 0)
						res[key] = append(sli, c)
					default:
						logger.Warn()
					}
				}
				if def.IsSimple() {
					switch inValue.(type) {
					case string:
						if def.IsString() {
							res[key] = []string{inValue.(string)}
						}
						if def.IsNumber() {
							f, err := strconv.ParseFloat(inValue.(string), 64)
							if err != nil {
								log.Warn(err)
							}
							res[key] = []float64{f}
						}
					case float64:
						if def.IsNumber() {
							res[key] = []float64{inValue.(float64)}
						}
						if def.IsString() {
							s := strconv.FormatFloat(inValue.(float64), 'f', -1, 64)
							res[key] = []string{s}
						}
					default:
						logger.Warn()
					}
				}
			}
		} else {
			if typ := reflect.TypeOf(inValue); typ.Kind() == reflect.Slice {
				if def.IsComplex() {
					switch inValue.(type) {
					case []map[string]interface{}:
						slim := inValue.([]map[string]interface{})
						if len(slim) == 0 {
							res[key] = nil
						} else {
							res[key] = d.ParseSource(d.Complex[key], slim[0])
						}
					default:
						logger.Warn()
					}
				}
				if def.IsSimple() {
					switch inValue.(type) {
					case []string:
						strs := inValue.([]string)
						if len(strs) == 0 {
							res[key] = nil
						} else {
							if def.IsString() {
								res[key] = strs[0]
							}
							if def.IsNumber() {
								f, err := strconv.ParseFloat(strs[0], 64)
								if err != nil {
									log.Warn(err)
								}
								res[key] = f
							}
						}
					case []float64:
						floats := inValue.([]float64)
						if len(floats) == 0 {
							res[key] = nil
						} else {
							if def.IsNumber() {
								res[key] = floats[0]
							}
							if def.IsString() {
								s := strconv.FormatFloat(floats[0], 'f', -1, 64)
								res[key] = s
							}
						}
					default:
						logger.Warn()
					}
				}
			} else {
				if def.IsComplex() {
					switch inValue.(type) {
					case map[string]interface{}:
						m := inValue.(map[string]interface{})
						c := d.ParseSource(d.Complex[key], m)
						res[key] = c
					default:
						logger.Warn()
					}
				}
				if def.IsSimple() {
					switch inValue.(type) {
					case string:
						if def.IsString() {
							res[key] = inValue.(string)
						}
						if def.IsNumber() {
							f, err := strconv.ParseFloat(inValue.(string), 64)
							if err != nil {
								log.Warn(err)
							}
							res[key] = f
						}
					case float64:
						if def.IsNumber() {
							res[key] = inValue.(float64)
						}
						if def.IsString() {
							s := strconv.FormatFloat(inValue.(float64), 'f', -1, 64)
							res[key] = s
						}
					default:
						logger.Warn()
					}
				}
			}
		}
	}
	return res
}

func main() {

	data, _ := NewSpecification().Get("ammeter.yaml")
	dataDefine, _ := GenerateDataDefine(data)
	dataDefine.To([]byte(`{"id":"电表1","description":"描述","data":{"voltage":[220,180],"current":["10","10"],"power":2200}}`))

}
