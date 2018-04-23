package models

import (
	"testing"
)

func TestTrain_LoadFromStorage(t *testing.T) {
	LoadedMaterials = Materials{
		Material{Name: "M1", Size: 5},
		Material{Name: "M2", Size: 2},
		Material{Name: "M3", Size: 1},
	}

	t1 := Train{Name: "T1", MaxCapacity: 300}
	t2 := Train{Name: "T2", MaxCapacity: 111}

	testStations := Workstations{&Workstation{Name: "WS1",
		Requirements: []string{"M1,10", "M2,10", "M3,5"},
	}, &Workstation{Name: "WS2",
		Requirements: []string{"M1,10", "M2,20", "M3,15"},
	},
	}
	t1.LoadFromStorage(testStations)
	t2.LoadFromStorage(testStations)

	tests := map[string]struct {
		input    int
		expected int
		err      error
	}{
		"T1_M1": {
			t1.Stock.GetAmount("M1"),
			100,
			nil,
		},
		"T1_M2": {
			t1.Stock.GetAmount("M2"),
			60,
			nil,
		},
		"T1_M3": {
			t1.Stock.GetAmount("M3"),
			20,
			nil,
		},
		"T2_M1": {
			t2.Stock.GetAmount("M1"),
			85,
			nil,
		},
		"T2_M2": {
			t2.Stock.GetAmount("M2"),
			20,
			nil},
		"T2_M3": {
			t2.Stock.GetAmount("M3"),
			6,
			nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		if test.input != test.expected {
			t.Errorf("%s: Expected: %d - Got: %d", testName, test.expected, test.input)
		}

	}

}

func TestTrain_Unload(t *testing.T) {

}
