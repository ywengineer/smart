package config

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/ywengineer/smart-kit/pkg/logk"
	"go.uber.org/zap"
	"reflect"
	"strings"
)

type loader interface {
	load(file string, processor func(row int, record []string) (interface{}, error), params []interface{}) error
}

type baseLoader struct {
	loader
}

func (bl *baseLoader) buildHeader(record []string) ([]*header, error) {
	hs := make([]*header, len(record))
	for i, rcd := range record {
		hDef := headerRegexp.FindStringSubmatch(rcd)
		h := &header{index: i, ignore: len(hDef) < 3}
		if h.ignore == false { // [all, name, type, ....]
			h.name, h.dataType = hDef[1], hDef[2]
		} else if strings.ContainsRune(rcd, '{') {
			logk.Error("wrong export column definition", zap.String("def", rcd))
		}
		if h.index == 0 && !strings.EqualFold(h.name, "id") {
			return nil, fmt.Errorf("the first col definition must be {id,int32}")
		}
		hs[i] = h
	}
	return hs, nil
}

func (bl *baseLoader) load(file, name string, cfg interface{}, params []interface{}) (*cache, error) {
	var headers []*header
	ch := &cache{
		name: name,
		data: make(map[int32]interface{}, 1000),
	}
	rowDataType := reflect.ValueOf(cfg).Type()
	if rowDataType.Kind() == reflect.Ptr {
		rowDataType = rowDataType.Elem()
	}
	err := bl.loader.load(file, bl.rowProcessor(&headers, ch, rowDataType), params)
	return ch, err
}

func (bl *baseLoader) rowProcessor(headers *[]*header, ch *cache, rowType reflect.Type) func(row int, record []string) (interface{}, error) {
	return func(row int, record []string) (interface{}, error) {
		// header
		if row == 0 {
			hdr, err := bl.buildHeader(record)
			*headers = hdr
			return hdr, err
		} else { // data
			// marshal row data to struct and add to cache
			rowValue := reflect.New(rowType).Interface()
			id, err := bl._load(*headers, record, rowValue)
			if err != nil {
				return nil, err
			}
			ch.data[id] = rowValue
			return rowValue, nil
		}
	}
}

func (bl *baseLoader) _load(headers []*header, r []string, dest interface{}) (int32, error) {
	// data
	// store row data
	tMap := make(map[string]interface{}, 10)
	var id int32
	// all exposed columns
	for _, headerDefine := range headers {
		if headerDefine.ignore {
			continue
		}
		hv, err := headerDefine.convert(r[headerDefine.index])
		if err != nil {
			return 0, fmt.Errorf("parse config data error. type[%s], col[%d]", headerDefine.dataType, headerDefine.index)
		}
		tMap[headerDefine.name] = hv
		// id col
		if headerDefine.index == 0 {
			id = hv.(int32)
		}
	}
	// marshal row data to struct
	jBytes, _ := sonic.Marshal(tMap)
	err := sonic.Unmarshal(jBytes, dest)
	return id, err
}
