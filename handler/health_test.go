package handler

import (
	"reflect"
	"testing"

	"gofr.dev/pkg/gofr"
)

func TestHealth(t *testing.T) {
	handler := func(_ *gofr.Context) (any, error) {
		return map[string]string{"status": "UP"}, nil
	}

	got, err := handler(nil)
	if err != nil {
		t.Fatalf("Handler returned an error: %v", err)
	}

	want := map[string]string{"status": "UP"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}
