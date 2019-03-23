package api

type QueueManager interface {
	GetTopic(topic string) (int64, error)
	ProcessTopic(topic int64, f func(queue int64) error) error
	GetQueue(queue string) (int64, error)

	MarkPublished(queue int64)
	WaitPublication(queue int64) bool // true - new publication, false - timeout
}

var (
	queueManagerProvider = make(map[string]QueueManager)
)

func RegisterQueueManager(name string, queueManger QueueManager) {
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

type QueueChangeEvent struct {
}

type QueueChangeListener interface {
	OnQueueCreated(queue int64, event *QueueChangeEvent)
	OnQueueDeleted(queue int64)
}

var (
	queueChangeListener = make(map[string]QueueChangeListener)
)

func RegisterQueueListener(name string, listener QueueChangeListener) {
	queueChangeListener[name] = listener
}

func GetQueueListener(name string) QueueChangeListener {
	q, ok := queueChangeListener[name]
	if !ok {
		return nil
	} else {
		return q
	}
}

func ForEachQueueListener(f func(listener QueueChangeListener) error) error {
	for _, v := range queueChangeListener {
		err := f(v)
		if err != nil {
			return err
		}
	}
	return nil
}
