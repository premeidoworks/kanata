package api

import "time"

type Message struct {
	Id           int64
	MessageId    string
	Queue        int64
	Topic        int64
	Body         []byte
	ScheduleTime int64
	Status       int64
	Type         int16 // 0 - at most once(pull mode)/at least once(push mode), 1 - exactly once
	OutId        string
	CreateTime   *time.Time
}

var (
	storeProviders = make(map[string]Store)
)

type StoreInitConfig struct {
	ConnectionString string

	Details map[string]string
}

func RegisterStoreProvider(name string, store Store) {
	storeProviders[name] = store
}

func GetStoreProvider(name string) Store {
	s, ok := storeProviders[name]
	if !ok {
		return nil
	} else {
		return s
	}
}

type Store interface {
	Init(config *StoreInitConfig) error

	SaveMessage(message *Message) error

	ObtainOnceMessage(queue int64, maxCount int) ([]*Message, error)
}
