package server_config

import (
	"testing"
)

func TestLocalJSONLoader(t *testing.T) {
	//
	loader := &LocalLoader{
		Path: "./conf.json",
	}
	var c = &Conf{}
	if err := loader.Load(c); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("%v", *c)
	}
}

func TestLocalYamlLoader(t *testing.T) {
	loader := &LocalLoader{
		Path: "./conf.yaml",
	}
	var c = &Conf{}
	if err := loader.Load(c); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("%v", *c)
	}
}

func TestUnknownLocalLoader(t *testing.T) {
	loader := &LocalLoader{
		Path: "./con.json",
	}
	var c = &Conf{}
	if err := loader.Load(c); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("%v", *c)
	}
}

func TestNoSuffixLocalLoader(t *testing.T) {
	loader := &LocalLoader{
		Path: "/etc/hosts",
	}
	var c = &Conf{}
	if err := loader.Load(c); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("%v", *c)
	}
}
