package config

import "testing"

func Test_CsvHeader(t *testing.T) {
	t.Logf("%v", headerRegexp.FindStringSubmatch("{name,int}"))
	t.Logf("%v", headerRegexp.FindStringSubmatch("{name,int64}"))
	t.Logf("%v", headerRegexp.FindStringSubmatch("{name,string}"))
	t.Logf("%v", headerRegexp.FindStringSubmatch("{name}"))
}
