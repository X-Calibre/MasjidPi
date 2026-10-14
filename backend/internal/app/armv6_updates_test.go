package app

import "testing"

func TestDisabledApplianceUpdatesDoNotCheckOrInstall(t *testing.T) {
	t.Setenv("MASJIDFRAME_APPLIANCE_UPDATES", "disabled")
	checker := &fakeScheduledUpdateChecker{}
	monitorUpdateChecks(t.Context(), checker, updateCheckTestLogger())
	if checker.statusCalls+checker.checkCalls+checker.prepareCalls+checker.installCalls != 0 {
		t.Fatal("disabled updater contacted the appliance update service")
	}
}

func TestApplianceUpdatesRemainEnabledByDefault(t *testing.T) {
	t.Setenv("MASJIDFRAME_APPLIANCE_UPDATES", "")
	if applianceUpdatesDisabled() {
		t.Fatal("standard installations must retain appliance updates")
	}
}
