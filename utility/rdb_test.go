package utility

import (
	"github.com/bytedance/sonic"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestRdbConfigProperties(t *testing.T) {
	rp := &RdbProperties{}
	t.Log(sonic.MarshalString(rp))
	rpYaml, _ := yaml.Marshal(rp)
	t.Log(string(rpYaml))
}
