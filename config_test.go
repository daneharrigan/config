package config_test

import (
	"github.com/daneharrigan/config"
	"testing"
)

func TestSimple(t *testing.T) {
	t.SkipNow()
	config, err := config.New("./examples/simple")
	if err != nil {
		t.Errorf("fn=New error=%q", err)
	}

	v, err := config.Get("section", "key")
	if err != nil {
		t.Errorf("fn=Get error=%q", err)
	}

	if v != "value" {
		t.Errorf("fn=Get error=%q", v+" was the wrong value")
	}

	v, err = config.Get("empty", "key")
	if err == nil {
		t.Errorf("fn=Get error=%q", err)
	}

	v, err = config.Get("multi", "foo")
	if err != nil {
		t.Errorf("fn=Get error=%q", err)
	}

	if v != "multi word value" {
		t.Errorf("fn=Get error=%q", v+" was the wrong value")
	}

	v, err = config.Get("multi", "bar")
	if err != nil {
		t.Errorf("fn=Get error=%q", err)
	}

	if v != "another multi word value" {
		t.Errorf("fn=Get error=%q", v+" was the wrong value")
	}
}

func TestSpaced(t *testing.T) {
	config, err := config.New("./examples/spaced")
	if err != nil {
		t.Errorf("fn=New error=%q", err)
	}

	v, err := config.Get("section", "foo")
	if err != nil {
		t.Errorf("fn=Get error=%q", err)
	}

	if v != "two spaces" {
		t.Errorf("fn=Get error=%q", v+" was the wrong value")
	}

	v, err = config.Get("section", "bar")
	if err != nil {
		t.Errorf("fn=Get error=%q", err)
	}

	if v != "tabs" {
		t.Errorf("fn=Get error=%q", v+" was the wrong value")
	}

	v, err = config.Get("tab", "baz")
	if err != nil {
		t.Errorf("fn=Get error=%q", err)
	}

	if v != "1" {
		t.Errorf("fn=Get error=%q", v+" was the wrong value")
	}

	v, err = config.Get("tab", "cux")
	if err != nil {
		t.Errorf("fn=Get error=%q", err)
	}

	if v != "2" {
		t.Errorf("fn=Get error=%q", v+" was the wrong value")
	}
}

func TestMalformed(t *testing.T) {
	_, err := config.New("./examples/malformed")
	if err == nil {
		t.Errorf("fn=New error=%q", "malformed formed file did not error")
	}
}
