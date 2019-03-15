package handler

import (
	"log"
	"net/http"

	"github.com/premeidoworks/kanata/api"
)

func OnPublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := prepareParams(w, r)
	if err != nil {
		log.Println("[ERROR] prepareParams error when Publish.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req, err := extractReq(r)
	if err != nil {
		log.Println("[ERROR] parse param error when Publish.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pubreq, err := MarshalProvider.UnmarshalPublishRequest(req)
	if err != nil {
		log.Println("[ERROR] PublishRequest cannot be unmarshal when Publish.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	topicId, err := QueueManager.GetTopic(pubreq.Topic)
	if err != nil {
		log.Println("[ERROR] no such topic when Publish.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := &api.PublishResponse{}
	for _, m := range pubreq.MessageList {
		idVal, err := UUID_Generator.Generate()
		if err != nil {
			log.Println("[ERROR] generate messageId error when Publish.", err)
			response.FailIdList = append(response.FailIdList, &struct {
				MsgId    string
				MsgOutId string
				Code     string
			}{
				MsgOutId: m.MsgOutId,
				Code:     "generate messageId error",
			})
			continue
		}
		err = QueueManager.ProcessTopic(topicId, func(queueId int64) error {
			msg := &api.Message{
				MessageId: idVal,
				Body:      m.MsgBody,
				Queue:     queueId,
				Topic:     topicId,
				Status:    0,
				OutId:     m.MsgOutId,
			}
			return StoreProvider.SaveMessage(msg)
		})
		if err != nil {
			log.Println("[ERROR] insert message error when Publish.", err)
			response.FailIdList = append(response.FailIdList, &struct {
				MsgId    string
				MsgOutId string
				Code     string
			}{
				MsgOutId: m.MsgOutId,
				MsgId:    idVal,
				Code:     "save message error",
			})
		} else {
			response.SuccessIdList = append(response.SuccessIdList, &struct {
				MsgId    string
				MsgOutId string
			}{
				MsgId:    idVal,
				MsgOutId: m.MsgOutId,
			})
		}
	}

	responseData, err := MarshalProvider.MarshalPublishResponse(response)
	if err != nil {
		log.Println("[ERROR] marshal publish response error when Publish.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseData)
	if err != nil {
		log.Println("[ERROR] write response error when Publish.", err)
		return
	}
}

func OnRollbackPublish(w http.ResponseWriter, r *http.Request) {

}

func OnCommitPublish(w http.ResponseWriter, r *http.Request) {

}
