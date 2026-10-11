package main

import (
	"os"
	"testing"
)

func TestOfficialExample(t *testing.T) {
	input, err := os.ReadFile("testdata/example.txt")
	if err != nil {
		t.Fatal(err)
	}

	versionSum, value, err := solve(string(input))

	if err != nil || versionSum != 6 || value != 2021 {
		t.Fatalf("got (%d, %d, %v), want (6, 2021, nil)", versionSum, value, err)
	}
}

func TestOfficialVersionSums(t *testing.T) {
	tests := []struct {
		transmission string
		versionSum   uint64
	}{
		{"8A004A801A8002F478", 16},
		{"620080001611562C8802118E34", 12},
		{"C0015000016115A2E0802F182340", 23},
		{"A0016C880162017C3686B18A3D4780", 31},
	}
	for _, tc := range tests {
		t.Run(tc.transmission, func(t *testing.T) {
			versionSum, _, err := solve(tc.transmission)

			if err != nil || versionSum != tc.versionSum {
				t.Fatalf("got (%d, %v), want (%d, nil)", versionSum, err, tc.versionSum)
			}
		})
	}
}

func TestOfficialValues(t *testing.T) {
	tests := []struct {
		transmission string
		value        uint64
	}{
		{"C200B40A82", 3},
		{"04005AC33890", 54},
		{"880086C3E88112", 7},
		{"CE00C43D881120", 9},
		{"D8005AC2A8F0", 1},
		{"F600BC2D8F", 0},
		{"9C005AC2F8F0", 0},
		{"9C0141080250320F1802104A08", 1},
	}
	for _, tc := range tests {
		t.Run(tc.transmission, func(t *testing.T) {
			_, value, err := solve(tc.transmission)

			if err != nil || value != tc.value {
				t.Fatalf("got (%d, %v), want (%d, nil)", value, err, tc.value)
			}
		})
	}
}

func TestTransmissionEdgeCases(t *testing.T) {
	tests := []struct {
		name, transmission string
		value              uint64
		wantError          bool
	}{
		{name: "lowercase digits and whitespace", transmission: "  d2fe28\n", value: 2021},
		{name: "length in bits", transmission: "38006F45291200", value: 1},
		{name: "count of sub-packets", transmission: "EE00D40C823060", value: 3},
		{name: "empty", transmission: "", wantError: true},
		{name: "blank", transmission: "\n", wantError: true},
		{name: "not hexadecimal", transmission: "D2FG28", wantError: true},
		{name: "truncated literal", transmission: "D2", wantError: true},
		{name: "truncated sub-packets", transmission: "38006F452912", wantError: true},
		{name: "data after the outermost packet", transmission: "D2FE28F", wantError: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, value, err := solve(tc.transmission)

			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError = %v", err, tc.wantError)
			}
			if !tc.wantError && value != tc.value {
				t.Fatalf("value = %d, want %d", value, tc.value)
			}
		})
	}
}
