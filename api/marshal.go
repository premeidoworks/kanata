package api

type MessageMarshal interface {
	MarshalPublishRequest(p *PublishRequest) ([]byte, error)
	UnmarshalPublishRequest(data []byte) (*PublishRequest, error)
	MarshalPublishResponse(p *PublishResponse) ([]byte, error)
	UnmarshalPublishResponse(data []byte) (*PublishResponse, error)
}

var (
	marshallingProvider = make(map[string]MessageMarshal)
)

func RegisterMarshallingProvider(name string, m MessageMarshal) {
	marshallingProvider[name] = m
}

func GetmarshallingProvider(name string) MessageMarshal {
	m, ok := marshallingProvider[name]
	if !ok {
		return nil
	} else {
		return m
	}
}

type PublishRequest struct {
	MessageList []*struct {
		MsgId    string
		MsgOutId string
		MsgBody  []byte
	}
}

type PublishResponse struct {
	SuccessIdList []*struct {
		MsgId    string
		MsgOutId string
	}
	FailIdList []*struct {
		MsgId    string
		MsgOutId string
		Code     string
	}
}
