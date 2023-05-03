package dmm

import (
	"reflect"
	"testing"
)

func TestDiffSlots(t *testing.T) {
	t.Helper()
	news := []Slot{
		{TeacherId: "a", DateTime: "11-14_22-30"},
		{TeacherId: "a", DateTime: "11-24_11-30"},
		{TeacherId: "a", DateTime: "03-05_08-30"},
		{TeacherId: "a", DateTime: "08-30_00-00"},
		{TeacherId: "a", DateTime: "06-21_11-30"},
	}
	exists := []Slot{
		{TeacherId: "a", DateTime: "11-24_11-30"},
		{TeacherId: "a", DateTime: "03-30_01-00"},
		{TeacherId: "a", DateTime: "11-14_22-30"},
		{TeacherId: "a", DateTime: "06-21_11-30"},
		{TeacherId: "a", DateTime: "03-05_08-30"},
	}
	want := map[string][]Slot{
		"adds": {{TeacherId: "a", DateTime: "08-30_00-00"}},
		"dels": {{TeacherId: "a", DateTime: "03-30_01-00"}},
	}
	adds, dels, _ := DiffSlots(news, exists)
	if !reflect.DeepEqual(adds, want["adds"]) {
		t.Errorf("got %+v, want %+v", adds, want["adds"])
	}
	if !reflect.DeepEqual(dels, want["dels"]) {
		t.Errorf("got %+v, want %+v", adds, want["dels"])
	}
}
