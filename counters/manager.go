package counters

// This class handles creation counters of all types and
// adding/removing from them
type manager struct {
}

func newManager() *manager {
	return &manager{}
}

func (m *manager) createDomain(domain string, domainType string) error {
	return nil
}

func (m *manager) deleteDomain(domain string) error {
	return nil
}

func (m *manager) getDomains() ([]string, error) {
	return nil, nil
}

func (m *manager) addToDomain(domain string, values []string) error {
	return nil
}

func (m *manager) deleteFromDomain(domain string, values []string) error {
	return nil
}
