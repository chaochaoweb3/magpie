package cli

import (
	"strings"
	"testing"

	"github.com/yetone/magpie/internal/settings"
)

// magpie quota alert --back turns the allowance-back alert on and off, and
// leaves the usage and balance alerts as they were (#6821).
func TestQuotaAlertBack(t *testing.T) {
	settingsHome(t)
	out, err := runOut(t, "quota", "alert", "--back", "on")
	if err != nil {
		t.Fatal(err)
	}
	if !settings.Load().AllowanceBackAlert || settings.Load().UsageAlert != 0 || settings.Load().BalanceAlert != 0 {
		t.Fatalf("saved %+v", settings.Load())
	}
	if !strings.Contains(out, "once when a window that was used up renews") || !strings.Contains(out, "no usage alert") {
		t.Fatalf("said %q", out)
	}
	if _, err := runOut(t, "quota", "alert", "80", "--back", "off"); err != nil {
		t.Fatal(err)
	}
	if s := settings.Load(); s.AllowanceBackAlert || s.UsageAlert != 80 {
		t.Fatalf("after off: %+v", s)
	}
	if _, err := runOut(t, "quota", "alert", "--back", "later"); err == nil || !strings.Contains(err.Error(), "on or off") {
		t.Fatalf("bad value: %v", err)
	}
	if settings.Load().UsageAlert != 80 {
		t.Fatal("a refused --back wrote the settings")
	}
	if _, err := runOut(t, "quota", "alert", "--back"); err == nil {
		t.Fatal("--back with no value was accepted")
	}
}
