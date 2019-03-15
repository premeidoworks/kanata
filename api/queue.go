package api

type QueueManager interface {
	GetTopic(topic string) (int64, error)
	ProcessTopic(topic int64, f func(queue int64) error) error
}

var (
	queueManagerProvider = make(map[string]QueueManager)
)

func RegisterQueueManger(name string, queueManger QueueManager) {
	queueManagerProvider[name] = queueManger
}

func GetQueueManager(name string) QueueManager {
	q, ok := queueManagerProvider[name]
	if !ok {
		return nil
	} else {
		return q
	}
}
