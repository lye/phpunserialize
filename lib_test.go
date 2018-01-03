package main

import (
	"testing"
)

func testUnmarshal(t *testing.T, bs string, out interface{}) {
	er := Unmarshal([]byte(bs), out)
	if er != nil {
		t.Fatalf("Unexpected error: %s", er.Error())
	}
}

func TestUnserializeInt(t *testing.T) {
	var ival int

	testUnmarshal(t, "i:1;", &ival)

	if ival != 1 {
		t.Errorf("expected 1, got %d", ival)
	}
}

func TestUnmarshalInt64(t *testing.T) {
	var ival int64

	testUnmarshal(t, "i:42;", &ival)

	if ival != 42 {
		t.Errorf("expected 42, got %d", ival)
	}
}

func TestUnmarshalString(t *testing.T) {
	var sval string

	testUnmarshal(t, `s:5:"hello";`, &sval)

	if sval != "hello" {
		t.Errorf("expected 'hello', got '%s'", sval)
	}
}

func TestUnmarshalBool(t *testing.T) {
	var bval bool

	testUnmarshal(t, "b:1;", &bval)
	if !bval {
		t.Errorf("expected true, got %#v", bval)
	}

	testUnmarshal(t, "b:0;", &bval)
	if bval {
		t.Errorf("expected false, got %#v", bval)
	}
}

func TestUnmarshalFloat64(t *testing.T) {
	var fval float64

	testUnmarshal(t, "d:6.5;", &fval)

	if fval != 6.5 {
		t.Errorf("expected 6.5, got %#v", fval)
	}
}

func TestUnmarshalFloat32(t *testing.T) {
	var fval float32

	testUnmarshal(t, "d:6.5;", &fval)

	if fval != 6.5 {
		t.Errorf("expected 6.5, got %#v", fval)
	}
}

func TestUnmarshalIntArray(t *testing.T) {
	var ary []interface{}

	testUnmarshal(t, "a:3:{i:0;i:1;i:1;i:2;i:2;i:3;}", &ary)

	if len(ary) != 3 {
		t.Fatalf("expected 3 elements, got %#v", ary)
	}

	for i, val := range ary {
		if ival, ok := val.(int64); !ok {
			t.Errorf("expected all int64's got %#v", ary)
		} else if ival != int64(i+1) {
			t.Errorf("expected %d, got %d", i+1, ival)
		}
	}
}

func TestUnmarshalMap(t *testing.T) {
	mapping := map[string]interface{}{}

	testUnmarshal(t, `a:2:{s:3:"one";i:1;s:3:"two";i:2;}`, &mapping)

	if len(mapping) != 2 {
		t.Fatalf("expected 2 elements, got %#v", mapping)
	}

	if val, ok := mapping["one"]; !ok {
		t.Errorf("no 'one' in %#v", mapping)
	} else if ival, ok := val.(int64); !ok || ival != 1 {
		t.Errorf("one is wrong %#v", mapping)
	}

	if val, ok := mapping["two"]; !ok {
		t.Errorf("no 'two' in %#v", mapping)
	} else if ival, ok := val.(int64); !ok || ival != 2 {
		t.Errorf("two is wrong %#v", mapping)
	}
}
