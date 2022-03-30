package config

import (
	"encoding/csv"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/ywengineer/mr.smart/log"
	"go.uber.org/zap"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func RegisterCsvAs(file, name string, cfg interface{}) error {
	cf, err := os.Open(file)
	if err != nil {
		return err
	}
	defer cf.Close()
	reader := csv.NewReader(cf)
	reader.Comment = '#' // 可以设置读入文件中的注释符
	reader.Comma = ','   // 默认是逗号，也可以自己设置
	// 还可以设置以下信息
	// FieldsPerRecord  int  // Number of expected fields per record
	// LazyQuotes       bool // Allow lazy quotes
	// TrailingComma    bool // Allow trailing comma
	// TrimLeadingSpace bool // Trim leading space
	// line             int
	// column           int
	k := 0 // 第一行是字段名，不需要
	var headers []*header
	ch := &cache{
		name: name,
		data: make(map[int32]interface{}, 1000),
	}
	cfgType := reflect.ValueOf(cfg).Type()
	if cfgType.Kind() == reflect.Ptr {
		cfgType = cfgType.Elem()
	}
	for {
		r, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		// header
		if k == 0 {
			headers, err = buildHeader(r)
			if err != nil {
				return fmt.Errorf("[%s] %s", name, err.Error())
			}
		} else { // data
			// store row data
			tMap := make(map[string]interface{}, 10)
			var id int32
			// all exposed columns
			for _, headerDefine := range headers {
				if headerDefine.ignore {
					continue
				}
				if tMap[headerDefine.name], err = headerDefine.convert(r[headerDefine.index]); err != nil {
					// id col
					if headerDefine.index == 0 {
						id = tMap[headerDefine.name].(int32)
					}
					// data col
					log.GetLogger().Error("parse config data error", zap.String("name", name), zap.String("type", headerDefine.dataType),
						zap.String("file", file), zap.Int("row", k+1), zap.Int("col", headerDefine.index),
					)
				}
			}
			// marshal row data to struct and add to cache
			rowValue := reflect.New(cfgType).Interface()
			jBytes, _ := sonic.Marshal(tMap)
			_ = sonic.Unmarshal(jBytes, rowValue)
			ch.data[id] = rowValue
		}
		k = k + 1
	}
	gCache[name] = ch
	return nil
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

func buildHeader(record []string) ([]*header, error) {
	hs := make([]*header, len(record))
	for i, rcd := range record {
		hDef := headerRegexp.FindStringSubmatch(rcd)
		h := &header{index: i, ignore: len(hDef) < 3}
		if h.ignore == false { // [all, name, type, ....]
			h.name, h.dataType = hDef[1], hDef[2]
		} else if strings.ContainsRune(rcd, '{') {
			log.GetLogger().Error("wrong export column definition", zap.String("def", rcd))
		}
		if h.index == 0 && !strings.EqualFold(h.name, "id") {
			return nil, fmt.Errorf("the first col definition must be {id,int32}")
		}
		hs[i] = h
	}
	return hs, nil
}

var headerRegexp = regexp.MustCompile(".*\\{([a-zA-Z]+)\\,(int|bool|int32|json|string)\\}.*")
