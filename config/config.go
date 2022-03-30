package config

const DataTypeInt = "int"
const DataTypeString = "string"
const DataTypeInt32 = "int32"
const DataTypeJSON = "json"
const DataTypeBool = "bool"

var gCache = make(map[string]*cache)

type cache struct {
	name string
	data map[int32]interface{}
}

func GetConfig(name string, id int32) (interface{}, bool) {
	ch, ok := gCache[name]
	if ok == false {
		return nil, false
	}
	cfg, ok := ch.data[id]
	return cfg, ok
}
