package kanata_discovery

type QuorumManager interface {
	Start() error
	Shutdown() error

	AcquireSession() (int64, error)
	SessionHeatbeat(session int64) error

	WatchServiceChange(prefix string, fn func(eventType WatchEvent, k, v string)) error
	PutService(key string, value string, session int64) error

	GetAllServices(prefix string) ([]struct {
		Key   string
		Value string
	}, error)
}
