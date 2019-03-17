package core

import (
	"github.com/premeidoworks/kanata/api"
	"time"
)

func Init() {
	qm := new(CoreQueueManager)
	//TODO hardcoded init
	qm.queueMap = map[int64]*struct {
		Notify      bool
		Queue       chan *api.Message
		NotifyQueue chan struct{}
	}{
		1: {
			NotifyQueue: make(chan struct{}, 1),
		},
	}
	api.RegisterQueueManger("default", qm)
}

type CoreQueueManager struct {
	queueMap map[int64]*struct {
		Notify      bool
		Queue       chan *api.Message
		NotifyQueue chan struct{}
	}
}

func (this *CoreQueueManager) MarkPublished(queue int64) {
	q, ok := this.queueMap[queue]
	if ok {
		select {
		case q.NotifyQueue <- struct{}{}:
		default:
		}
	}
}

func (this *CoreQueueManager) WaitPublication(queue int64) bool {
	q, ok := this.queueMap[queue]
	if ok {
		timer := time.NewTimer(5 * time.Second) //TODO should use configuration instead of magic number
		select {
		case <-q.NotifyQueue:
			timer.Stop()
			return true
		case <-timer.C:
			timer.Stop()
			return false
		}
	} else {
		return false
	}
}

func (this *CoreQueueManager) GetTopic(topic string) (int64, error) {
	return getTopic(topic)
}

func (this *CoreQueueManager) ProcessTopic(topic int64, f func(queue int64) error) error {
	l := getTopicQueues(topic)
	for _, q := range l {
		err := f(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *CoreQueueManager) GetQueue(queue string) (int64, error) {
	return getQueue(queue)
}

func getTopic(topic string) (int64, error) {
	//TODO
	return 1, nil
}

func getQueue(queue string) (int64, error) {
	//TODO
	return 1, nil
}

func getTopicQueues(topic int64) []int64 {
	//TODO
	return []int64{1}
}
