package utility

import (
	"github.com/bytedance/sonic"
	"gopkg.in/yaml.v2"
	"net/url"
	"testing"
)

func TestRdbConfigProperties(t *testing.T) {
	rp := &RdbProperties{}
	t.Log(sonic.MarshalString(rp))
	rpYaml, _ := yaml.Marshal(rp)
	t.Log(string(rpYaml))
}

func TestParseUrl(t *testing.T) {
	_case := []string{
		"mem://?size=120000",
		"redis://127.0.0.1:6379/?db=0",
		"redis://user:password@127.0.0.1:6379/1",
		"redis-cluster://user:password@127.0.0.1:6379/?node=127.0.0.1:6381&node=127.0.0.1:6380",
	}
	for _, v := range _case {
		r, _ := url.Parse(v)
		t.Logf("%+v", r.Query())
	}
}
