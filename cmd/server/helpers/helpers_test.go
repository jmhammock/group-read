package helpers

import (
	"reflect"
	"testing"
)

func TestStrToUint16(t *testing.T) {
	v, err := StrToUint32("16")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if v != uint32(16) {
		t.Logf("expected 16 got %d\n", v)
		t.Fail()
	}

	tp := reflect.ValueOf(v).Kind()
	if tp != reflect.Uint32 {
		t.Logf("expected type uint16 got %s\n", tp.String())
		t.Fail()
	}

	v, err = StrToUint32("owen")
	if err == nil {
		t.Log("expected error")
		t.Fail()
	}
}
