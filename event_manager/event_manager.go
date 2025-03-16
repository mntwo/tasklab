package event_manager

import (
	"sync"

	"github.com/mntwo/tasklab/gen_event"
)

var manager *Manager

// Manager is a struct that manages multiple EventManagers identified by aliases.
type Manager struct {
	eventManagers map[string]*gen_event.EventManager // A map of aliases to EventManagers.
	mu            sync.RWMutex                       // A read-write mutex to protect the eventManagers map.
}

func init() {
	manager = &Manager{
		eventManagers: make(map[string]*gen_event.EventManager),
	}
}

func Stop() {
	for _, em := range manager.eventManagers {
		em.Close()
	}
}

// addEventManager adds a new EventManager with a specified alias.
func (m *Manager) addEventManager(alias string, em *gen_event.EventManager) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.eventManagers[alias] = em
}

// getEventManager retrieves an EventManager by its alias.
func (m *Manager) getEventManager(alias string) (*gen_event.EventManager, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	em, exists := m.eventManagers[alias]
	return em, exists
}

// removeEventManager removes an EventManager by its alias.
func (m *Manager) removeEventManager(alias string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if em, exists := m.eventManagers[alias]; exists {
		em.Close()
		delete(m.eventManagers, alias)
	}
}

func AddEventManager(alias string, em *gen_event.EventManager) {
	manager.addEventManager(alias, em)
}

func GetEventManager(alias string) (*gen_event.EventManager, bool) {
	return manager.getEventManager(alias)
}

func RemoveEventManager(alias string) {
	manager.removeEventManager(alias)
}
