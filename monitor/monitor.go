package monitor

import (
	"sync"

	"github.com/pkg/errors"
)

// CadMonitor defines the interface for all monitors used to watch CAD systems.
type CadMonitor interface {
	// ConfigureFromValues populates fields specific to an implementation of
	// CadMonitor from a map[string]string.
	ConfigureFromValues(map[string]string) error
	// Login authenticates to a CAD system using the provided username and password
	Login(string, string) error
	// GetActiveCalls returns a list of active call URLs or identifiers
	GetActiveCalls() ([]string, error)
	// GetActiveAndUnassignedCalls returns a list of CallStatus objects for unassigned and
	// current calls
	GetActiveAndUnassignedCalls() (map[string]CallStatus, error)
	// GetStatus given an identifier retrieves a CallStatus entry that describes a call
	GetStatus([]byte, string) (CallStatus, error)
	// GetStatusFromURL given an identifier retrieves a CallStatus entry that describes a call
	GetStatusFromURL(string) (CallStatus, error)
	// GetClearedCalls retrieves a map of ids for cleared calls for a specific date
	GetClearedCalls(string) (map[string]string, error)
	// SetDebug determines whether debug is enabled or not
	SetDebug(bool)
	// KeepAlive represents some manner of maintaining a persistent connection
	KeepAlive() error
	TerminateMonitor() bool
	SetTerminateMonitor(bool)
	// Monitor actively runs a monitoring function with a callback
	Monitor(func(CallStatus) error, int) error
}

var (
	// ErrCadMonitorLoggedOut represents a status where the application needs to reauthenticate
	ErrCadMonitorLoggedOut = errors.New("logged out")

	cadMonitorRegistry     = map[string]func() CadMonitor{}
	cadMonitorRegistryLock = new(sync.Mutex)
)

// RegisterCadMonitor adds a new CadMonitor instance to the registry
func RegisterCadMonitor(name string, m func() CadMonitor) {
	cadMonitorRegistryLock.Lock()
	defer cadMonitorRegistryLock.Unlock()
	cadMonitorRegistry[name] = m
}

// InstantiateCadMonitor instantiates a CadMonitor by name
func InstantiateCadMonitor(name string) (m CadMonitor, err error) {
	var f func() CadMonitor
	var found bool
	if f, found = cadMonitorRegistry[name]; !found {
		err = errors.New("unable to locate cad monitor " + name)
		return
	}
	m = f()
	err = nil
	return
}
