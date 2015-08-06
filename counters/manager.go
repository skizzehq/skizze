package counters

/*
Manager is responsible for manipulating the counters and syncing to disk
*/
type Manager struct {
}

/*
CreateDomain ...
*/
func (m *Manager) CreateDomain(domain string, domainType string) error {
	return nil
}

/*
DeleteDomain ...
*/
func (m *Manager) DeleteDomain(domain string) error {
	return nil
}

/*
GetDomains ...
*/
func (m *Manager) GetDomains() ([]string, error) {
	// TODO: Remove dummy data and implement proper result
	return []string{"foo", "bar"}, nil
}

/*
AddToDomain ...
*/
func (m *Manager) AddToDomain(domain string, values []string) error {
	return nil
}

/*
DeleteFromDomain ...
*/
func (m *Manager) DeleteFromDomain(domain string, values []string) error {
	return nil
}

func newManager() *Manager {
	return &Manager{}
}

// Manager handles creation counters of all types and
// adding/removing from them
var ManagerProxy = newManager()
