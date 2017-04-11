package config

import (
	"fmt"
	"testing"
)

func TestPrecedence(t *testing.T) {
	low := map[string]string{
		"without_override": "false",
		"with_override":    "false",
	}

	high := map[string]string{
		"with_override": "true",
	}

	c := NewConfig(
		[]Provider{
			NewStatic(low),
			NewStatic(high),
		},
	)

	without, err := c.Bool("without_override")
	if err != nil {
		t.Error(err)
	}

	if without == true {
		t.Errorf("Setting 'without_override' was true, expected false")
	}

	with, err := c.Bool("with_override")
	if err != nil {
		t.Error(err)
	}

	if with == false {
		t.Errorf("Setting 'with_override' was 'false', expected 'true'")
	}
}

func TestTypeLookups(t *testing.T) {
	settings := map[string]string{
		"string": "some_string",
		"bool":   "true",
		"int":    "1",
		"float":  "1.5",
	}

	c := NewConfig([]Provider{NewStatic(settings)})

	s, err := c.String("string")
	if err != nil {
		t.Error(err)
	}

	if s != "some_string" {
		t.Errorf("String setting was '%s', expected 'some_string'", s)
	}

	b, err := c.Bool("bool")
	if err != nil {
		t.Error(err)
	}

	if b != true {
		t.Errorf("Bool setting was 'false', expected 'true'")
	}

	i, err := c.Int("int")
	if err != nil {
		t.Error(err)
	}

	if i != 1 {
		t.Errorf("Int setting was '%d', expected '1'", i)
	}

	f, err := c.Float("float")
	if err != nil {
		t.Error(err)
	}

	if f != 1.5 {
		t.Errorf("Float setting was '%f', expected '1.5'", f)
	}
}

func TestTypeOrLookups(t *testing.T) {
	c := NewConfig(nil)

	s, err := c.StringOr("string", "some_string")
	if err != nil {
		t.Error(err)
	}

	if s != "some_string" {
		t.Errorf("String setting was '%s', expected 'some_string'", s)
	}

	b, err := c.BoolOr("bool", true)
	if err != nil {
		t.Error(err)
	}

	if b != true {
		t.Errorf("Bool setting was 'false', expected 'true'")
	}

	i, err := c.IntOr("int", 1)
	if err != nil {
		t.Error(err)
	}

	if i != 1 {
		t.Errorf("Int setting was '%d', expected '1'", i)
	}

	f, err := c.FloatOr("float", 1.5)
	if err != nil {
		t.Error(err)
	}

	if f != 1.5 {
		t.Errorf("Float setting was '%f', expected '1.5'", f)
	}
}

func TestValidate(t *testing.T) {
	c := NewConfig(nil)
	c.Validate = func(map[string]string) error {
		return fmt.Errorf("some error")
	}

	if err := c.Load(); err == nil {
		t.Errorf("Error was nil")
	}
}
