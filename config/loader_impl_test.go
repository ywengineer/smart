package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CsvHeader(t *testing.T) {
	t.Logf("%v", headerRegexp.FindStringSubmatch("{name,int}"))
	t.Logf("%v", headerRegexp.FindStringSubmatch("{name,int64}"))
	t.Logf("%v", headerRegexp.FindStringSubmatch("{name,string}"))
	t.Logf("%v", headerRegexp.FindStringSubmatch("{name, string}"))
	t.Logf("%v", headerRegexp.FindStringSubmatch("{name}"))
}

type Address struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
}
type ConfigTest struct {
	Id      int32    `json:"id"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Male    bool     `json:"male"`
	Address *Address `json:"address"`
}

func Test_LoadCsv(t *testing.T) {
	err := LoadConfig(Csv, "./csv_config_test.csv", "test", &ConfigTest{})
	assert.Nil(t, err)
	t.Log(FindConfig("test", 10001))
	t.Log(FindConfig("test", 10002))
}

func Test_LoadExcel(t *testing.T) {
	err := LoadConfig(Excel, "./xlsx_config_test.xlsx", "test", &ConfigTest{}, 0)
	assert.Nil(t, err)
	t.Log(FindConfig("test", 10003))
	t.Log(FindConfig("test", 10004))
}
