package counters

import "testing"

func TestNoCounters(t *testing.T) {
	manager := GetManager()
	domains, err := manager.GetDomains()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if len(domains) != 0 {
		t.Error("Expected 0 counters, got", len(domains))
	}
}

func TestCounter(t *testing.T) {
	manager := GetManager()

	err := manager.CreateDomain("marvel", "immutable")
	if err != nil {
		t.Error("Expected no errors while creating domain, got", err)
	}

	domains, err := manager.GetDomains()
	if err != nil {
		t.Error("Expected no errors while getting domains, got", err)
	}
	if len(domains) != 1 {
		t.Error("Expected 1 counters, got", len(domains))
	}

	err = manager.AddToDomain("marvel", []string{"hulk", "thor"})
	if err != nil {
		t.Error("Expected no errors while adding to domain, got", err)
	}

	count, err := manager.GetCountForDomain("marvel")
	if len(domains) != 1 {
		t.Error("Expected 1 counters, got", len(domains))
	}

	if count != 2 {
		t.Error("Expected count == 2, got", count)
	}
}
