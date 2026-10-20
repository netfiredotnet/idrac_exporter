package collector

import (
	"encoding/json"
	"testing"
)

func TestVolumeSizesAcceptStrings(t *testing.T) {
	// A Supermicro sends StripSizeBytes as a string, which used to fail the
	// whole volume.
	data := `{"Id": "0", "CapacityBytes": 1919850381312, "StripSizeBytes": "65536", "BlockSizeBytes": 512, "OptimumIOSizeBytes": null}`
	var v StorageVolume
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v.StripSizeBytes != 65536 || v.BlockSizeBytes != 512 || v.OptimumIOSizeBytes != 0 {
		t.Fatalf("got strip=%d block=%d optimum=%d", v.StripSizeBytes, v.BlockSizeBytes, v.OptimumIOSizeBytes)
	}
	if v.CapacityBytes != 1919850381312 {
		t.Fatalf("capacity = %v", v.CapacityBytes)
	}
}

func TestXintIgnoresNonNumbers(t *testing.T) {
	for _, in := range []string{`"n/a"`, `null`, `true`, `{}`} {
		var x xint = 7
		if err := json.Unmarshal([]byte(in), &x); err != nil {
			t.Fatalf("%s: %v", in, err)
		}
		if x != 0 {
			t.Fatalf("%s: got %d, want 0", in, x)
		}
	}
}

func TestIntrusion2Value(t *testing.T) {
	cases := map[string]int{
		"":                  -1,
		"Normal":            0,
		"HardwareIntrusion": 1,
		"TamperingDetected": 1,
	}
	for in, want := range cases {
		if got := intrusion2value(in); got != want {
			t.Errorf("intrusion2value(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestChassisPhysicalSecurity(t *testing.T) {
	data := `{"Status": {"Health": "Critical"}, "PhysicalSecurity": {"IntrusionSensor": "HardwareIntrusion", "IntrusionSensorReArm": "Manual"}}`
	var c ChassisResponse
	if err := json.Unmarshal([]byte(data), &c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if c.Status.Health != "Critical" || c.PhysicalSecurity == nil || c.PhysicalSecurity.IntrusionSensor != "HardwareIntrusion" {
		t.Fatalf("got %+v", c)
	}
}
