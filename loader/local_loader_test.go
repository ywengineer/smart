package loader

import (
	"github.com/bytedance/sonic"
	"testing"
)

func TestLocalJSONLoader(t *testing.T) {
	//
	loader := &localLoader{
		path: "./conf.json",
	}
	var c = &Conf{}
	if err := loader.Load(c); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Log(sonic.MarshalString(c))
	}
}

func TestLocalYamlLoader(t *testing.T) {
	loader := &localLoader{
		path: "./conf.yaml",
	}
	var c = &Conf{}
	if err := loader.Load(c); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("%v", *c)
	}
}

func TestUnknownLocalLoader(t *testing.T) {
	loader := &localLoader{
		path: "./con.json",
	}
	var c = &Conf{}
	if err := loader.Load(c); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("%v", *c)
	}
}

func TestNoSuffixLocalLoader(t *testing.T) {
	loader := &localLoader{
		path: "/etc/hosts",
	}
	var c = &Conf{}
	if err := loader.Load(c); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("%v", *c)
	}
}
