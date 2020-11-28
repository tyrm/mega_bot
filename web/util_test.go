package web

import "testing"

func TestIsValidUUID(t *testing.T) {
	tables := []struct {
		x string
		n bool
	}{
		{"c3ac5fe1-561f-4d69-8af1-df4cd386ffbb", true},
		{"53042eea-c077-489a-a623-96355562c9f9", true},
		{"dcf14a94-749d-4be1-b355-7cbda1eecb07", true},
		{"56ca3ece-ce5d-43b1-87f6-fe8cae2d50cd", true},
		{"81495fe8-8ead-4950-a2ab-5338d79c6b4e", true},
		{"not a uuid", false},
		{"1", false},
		{"53042eea-c077-989a-a623-96355562c9f9", false}, // invalid version
		{"d1bcc518-3184-11eb-adc1-0242ac120002", false}, // gen 1 uuid
	}

	for _, table := range tables {
		total := isValidUUID4(table.x)
		if total != table.n {
			t.Errorf("RoundUp of (%s) was incorrect, got: %v, want: %v.", table.x, total, table.n)
		}
	}
}

func TestRoundUp(t *testing.T) {
	tables := []struct {
		x float64
		n uint64
	}{
		{0, 0},
		{1, 1},
		{10, 10},
		{5783, 5783},
		{15687, 15687},
		{44, 44},
		{1.5, 2},
		{33.8, 34},
		{48.55555555555555555, 49},
		{82214.221387435, 82215},
	}

	for _, table := range tables {
		total := roundUp(table.x)
		if total != table.n {
			t.Errorf("RoundUp of (%f) was incorrect, got: %d, want: %d.", table.x, total, table.n)
		}
	}
}