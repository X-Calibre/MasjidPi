package masjidboardlive

import (
	"encoding/json"
	"testing"
)

func TestApplyDisplayTimesAddsIstiwaAndZawaal(t *testing.T) {
	astronomical, err := applyDisplayTimes(json.RawMessage(`["11:55","12:08","theme","English","12:01"]`), nil)
	if err != nil {
		t.Fatalf("applyDisplayTimes() error = %v", err)
	}
	if astronomical == nil {
		t.Fatal("astronomical is nil")
	}
	if astronomical.IstiwaCaution == nil || astronomical.IstiwaCaution.Hour != 11 || astronomical.IstiwaCaution.Minute != 55 {
		t.Fatalf("IstiwaCaution = %+v", astronomical.IstiwaCaution)
	}
	if astronomical.Istiwa == nil || astronomical.Istiwa.Hour != 12 || astronomical.Istiwa.Minute != 1 {
		t.Fatalf("Istiwa = %+v", astronomical.Istiwa)
	}
	if astronomical.ZawaalEnd == nil || astronomical.ZawaalEnd.Hour != 12 || astronomical.ZawaalEnd.Minute != 8 {
		t.Fatalf("ZawaalEnd = %+v", astronomical.ZawaalEnd)
	}
}
