package config

import (
	"fmt"
	"github.com/bytedance/sonic"
	"regexp"
	"strconv"
)

type CfgType int

const (
	Csv CfgType = iota
	Excel
)

func (ct CfgType) String() string {
	switch ct {
	case Csv:
		return "Csv"
	case Excel:
		return "Excel"
	}
	return "Unknown"
}

const DataTypeInt = "int"
const DataTypeString = "string"
const DataTypeInt32 = "int32"
const DataTypeJSON = "json"
const DataTypeBool = "bool"

var gCache = make(map[string]*cache)
var _csvLoader = &baseLoader{&csvLoader{}}
var _xlsxLoader = &baseLoader{&xlsxLoader{}}

type cache struct {
	name string
	data map[int32]interface{}
}

// LoadConfig int(params[0]) is the sheet number of an excel file when CfgType is Excel
func LoadConfig(t CfgType, file, name string, cfg interface{}, params ...interface{}) error {
	var c *cache
	var err error
	switch t {
	case Csv:
		c, err = _csvLoader.load(file, name, cfg, params)
	case Excel:
		c, err = _xlsxLoader.load(file, name, cfg, params)
	default:
		err = fmt.Errorf("unsupported CfgType: %s", t)
	}
	if err != nil {
		return err
	}
	return putConfig(name, c)
}

func putConfig(name string, ch *cache) error {
	if _, ok := gCache[name]; ok {
		return fmt.Errorf("config[%s] already registerd, skip register again", name)
	}
	gCache[name] = ch
	return nil
}

func FindConfig(name string, id int32) (interface{}, bool) {
	ch, ok := gCache[name]
	if ok == false {
		return nil, false
	}
	cfg, ok := ch.data[id]
	return cfg, ok
}

// header: {name[string],dataType[int,bool,int32,json,string]}
type header struct {
	name     string
	index    int
	dataType string
	ignore   bool
}

func (h *header) convert(str string) (interface{}, error) {
	switch h.dataType {
	case DataTypeBool:
		return strconv.ParseBool(str)
	case DataTypeInt:
		return strconv.Atoi(str)
	case DataTypeInt32:
		i, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return 0, err
		} else {
			return int32(i), err
		}
	case DataTypeJSON:
		m := make(map[string]interface{})
		err := sonic.UnmarshalString(str, &m)
		return m, err
	case DataTypeString:
		return str, nil
	}
	return nil, nil
}

var headerRegexp = regexp.MustCompile(".*\\{([a-zA-Z]+).*\\,.*(int|bool|int32|json|string)\\}.*")
