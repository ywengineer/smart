package server_config

import (
	"testing"
)

func TestLocalJSONLoader(t *testing.T) {
	//
	loader := &LocalLoader{
		Path: "./conf.json",
	}
	if c, err := loader.Load(); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("%v", *c)
	}
}

func TestLocalYamlLoader(t *testing.T) {
	loader := &LocalLoader{
		Path: "./conf.yaml",
	}
	if c, err := loader.Load(); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("%v", *c)
	}
}

func TestUnknownLocalLoader(t *testing.T) {
	loader := &LocalLoader{
		Path: "./con.json",
	}
	if c, err := loader.Load(); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("%v", *c)
	}
}

func TestNoSuffixLocalLoader(t *testing.T) {
	loader := &LocalLoader{
		Path: "/etc/hosts",
	}
	if c, err := loader.Load(); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Logf("%v", *c)
	}
}
