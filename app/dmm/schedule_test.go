package dmm

import (
	"reflect"
	"strconv"
	"testing"
)

func TestDiffSlots(t *testing.T) {
	t.Helper()
	tests := []map[string][]Slot{
		{
			"news": {
				{TeacherId: "a", DateTime: "11-14_22-30"},
				{TeacherId: "a", DateTime: "11-24_11-30"},
				{TeacherId: "a", DateTime: "03-05_08-30"},
				{TeacherId: "a", DateTime: "08-30_00-00"},
				{TeacherId: "a", DateTime: "06-21_11-30"},
			},
			"exists": {
				{TeacherId: "a", DateTime: "11-24_11-30"},
				{TeacherId: "a", DateTime: "03-30_01-00"},
				{TeacherId: "a", DateTime: "11-14_22-30"},
				{TeacherId: "a", DateTime: "06-21_11-30"},
				{TeacherId: "a", DateTime: "03-05_08-30"},
			},
			"adds": {{TeacherId: "a", DateTime: "08-30_00-00"}},
			"dels": {{TeacherId: "a", DateTime: "03-30_01-00"}},
		},
		{
			"news": {
				{TeacherId: "a", DateTime: "11-14_22-30"},
			},
			"exists": {},
			"adds":   {{TeacherId: "a", DateTime: "11-14_22-30"}},
			"dels":   {},
		},
		{
			"news":   {},
			"exists": {{TeacherId: "a", DateTime: "11-14_22-30"}},
			"adds":   {},
			"dels":   {{TeacherId: "a", DateTime: "11-14_22-30"}},
		},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			adds, dels, _ := DiffSlots(tt["news"], tt["exists"])
			if !reflect.DeepEqual(adds, tt["adds"]) {
				t.Errorf("got %+v, want %+v", adds, tt["adds"])
			}
			if !reflect.DeepEqual(dels, tt["dels"]) {
				t.Errorf("got %+v, want %+v", dels, tt["dels"])
			}
		})
	}
}

func TestAssertTeacherId(t *testing.T) {
	t.Helper()
	tests := []struct {
		input string
		want  bool
	}{
		{
			input: "53216",
			want:  true,
		},
		{
			input: "532160",
			want:  false,
		},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			got := AssertTeacherId(tt.input)
			if got != tt.want {
				t.Errorf("input %+v, got %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}
