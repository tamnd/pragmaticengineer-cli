package pragmaticengineer

import (
	"testing"
)

// These tests are offline: they exercise the domain's pure string functions
// and info, which need no network. The client's HTTP behaviour is covered
// in pragmaticengineer_test.go.

func TestDomainInfo(t *testing.T) {
	info := Domain{}.Info()
	if info.Scheme != "pragmaticengineer" {
		t.Errorf("Scheme = %q, want pragmaticengineer", info.Scheme)
	}
	if len(info.Hosts) == 0 || info.Hosts[0] != Host {
		t.Errorf("Hosts = %v, want [%s]", info.Hosts, Host)
	}
	if info.Identity.Binary != "pragmaticengineer" {
		t.Errorf("Identity.Binary = %q, want pragmaticengineer", info.Identity.Binary)
	}
}
