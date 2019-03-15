package core

import (
	"github.com/premeidoworks/kanata/api"
)

func Init() {
	api.RegisterQueueManger("default", new(CoreQueueManager))
}

type CoreQueueManager struct {
	queueMap map[int64]*struct {
		Notify bool
		Queue  chan *api.Message
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
