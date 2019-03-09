package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/premeidoworks/kanata/api"
	"github.com/premeidoworks/kanata/core"
)

func Publish(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	// GET method only support simplest publish
	case "GET":
		{
			query := r.URL.Query()

			topic, err := requiredString(query.Get("topic"))
			if err != nil {
				log.Println("[ERROR] topic is empty when Publish.", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			messageBody := []byte(query.Get("message"))
			if len(messageBody) == 0 {
				log.Println("[ERROR] message body is empty when Publish.", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			topicId, err := core.GetTopic(topic)
			if err != nil {
				log.Println("[ERROR] no such topic when Publish.", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			idVal, err := UUID_Generator.Generate()
			if err != nil {
				log.Println("[ERROR] generate messageId error when Publish.", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			queueList := core.GetTopicQueues(topicId)
			for _, queueId := range queueList {
				var _ = topic
				msg := &api.Message{
					MessageId: idVal,
					Body:      messageBody,
					Queue:     queueId,
					Topic:     topicId,
					Status:    0,
				}
				err = StoreProvider.SaveMessage(msg)
				if err != nil {
					log.Println("[ERROR] insert message error when Publish.", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			w.WriteHeader(http.StatusOK)
			result := map[string][]string{
				"successList": []string{idVal},
			}
			data, _ := json.Marshal(result)
			_, err = w.Write(data)
			if err != nil {
				log.Println("[ERROR] write response to client when Publish.", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			return
		}
	case "POST":
		{
			err := prepareParams(w, r)
			if err != nil {
				log.Println("[ERROR] prepareParams error when Publish.", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			form := r.Form
			//TODO support GET request
			topic, err := requiredString(form.Get("topic"))
			if err != nil {
				log.Println("[ERROR] topic is empty when Publish.", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			//TODO need to follow protocol instruction
			var _ = topic
			messageBody := []byte(form.Get("message_body"))
			msg := &api.Message{
				Body: messageBody,
			}
			err = StoreProvider.SaveMessage(msg)
			if err != nil {
				log.Println("[ERROR] insert message error when Publish.", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			//TODO need to response correctly
			w.WriteHeader(http.StatusOK)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func PrePublish(w http.ResponseWriter, r *http.Request) {

}

func CommitPublish(w http.ResponseWriter, r *http.Request) {

}
