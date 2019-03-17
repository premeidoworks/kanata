package handler

import (
	"log"
	"net/http"

	"github.com/premeidoworks/kanata/api"
)

func OnAcquire(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := prepareParams(w, r)
	if err != nil {
		log.Println("[ERROR] prepareParams error when Acquire.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req, err := extractReq(r)
	if err != nil {
		log.Println("[ERROR] parse param error when Acquire.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acReq, err := MarshalProvider.UnmarshalAcquireRequest(req)
	if err != nil {
		log.Println("[ERROR] Acquire cannot be unmarshal when Acquire.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	queueId, err := QueueManager.GetQueue(acReq.Queue)
	if err != nil {
		log.Println("[ERROR] Acquire cannot get queue when Acquire.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var msgs []*api.Message
	for {
		msgs, err = StoreProvider.ObtainOnceMessage(queueId, 16) //TODO should use configuration instead of magic number
		if err != nil {
			log.Println("[ERROR] Acquire cannot obtain once message when Acquire.", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(msgs) == 0 {
			ok := QueueManager.WaitPublication(queueId)
			if ok {
				continue
			} else {
				break
			}
		} else {
			QueueManager.MarkPublished(queueId)
			break
		}
	}

	response := &api.AcquireResponse{
		MessageList: make([]*struct {
			MsgId   string
			MsgBody []byte
		}, len(msgs)),
	}
	for i := 0; i < len(msgs); i++ {
		response.MessageList[i] = &struct {
			MsgId   string
			MsgBody []byte
		}{
			MsgId:   msgs[i].MessageId,
			MsgBody: msgs[i].Body,
		}
	}
	responseData, err := MarshalProvider.MarshalAcquireResponse(response)
	if err != nil {
		log.Println("[ERROR] marshal acquire response error when Acquire.", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseData)
	if err != nil {
		log.Println("[ERROR] write response error when Acquire.", err)
		return
	}
}

func OnCommit(w http.ResponseWriter, r *http.Request) {
}
