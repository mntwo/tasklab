package event_manager

import (
	"sync"

	"github.com/mntwo/tasklab/gen_event"
)

// Manager is a struct that manages multiple EventManagers identified by aliases.
type Manager struct {
	eventManagers map[string]*gen_event.EventManager // A map of aliases to EventManagers.
	mu            sync.RWMutex                       // A read-write mutex to protect the eventManagers map.
}

// NewManager creates a new Manager.
func NewManager() *Manager {
	return &Manager{
		eventManagers: make(map[string]*gen_event.EventManager),
	}
}

// AddEventManager adds a new EventManager with a specified alias.
func (m *Manager) AddEventManager(alias string, em *gen_event.EventManager) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.eventManagers[alias] = em
}

// GetEventManager retrieves an EventManager by its alias.
func (m *Manager) GetEventManager(alias string) (*gen_event.EventManager, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	em, exists := m.eventManagers[alias]
	return em, exists
}

// RemoveEventManager removes an EventManager by its alias.
func (m *Manager) RemoveEventManager(alias string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if em, exists := m.eventManagers[alias]; exists {
		em.Close()
		delete(m.eventManagers, alias)
	}
}
