package handler

import (
	"reflect"
	"testing"

	"gofr.dev/pkg/gofr"
)

func TestHealth(t *testing.T) {
	handler := func(_ *gofr.Context) any {
		return map[string]string{"status": "UP"}
	}

	got := handler(nil)

	want := map[string]string{"status": "UP"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}
