package lineapi

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/mopeneko/linethrift"
)

// Generate LINE TalkService client
func NewLineClient(authToken string) (*linethrift.TalkServiceClient, *thrift.THttpClient, error) {
	headers := map[string]string{
		"User-Agent":         USER_AGENT,
		"X-Line-Application": LINE_APP,
		"X-Line-Access":      authToken,
	}
	client, transport, err := NewThriftClient(HOST+TALKSERVICE_ENDPOINT, headers)
	talk := linethrift.NewTalkServiceClient(client)
	talk.AuthToken = authToken
	if err != nil {
		return nil, nil, err
	}
	return talk, transport, nil
}

func NewLineClient_(authToken string) (*linethrift.TalkServiceClient, *thrift.THttpClient, error) {
	headers := map[string]string{
		"User-Agent":         USER_AGENT,
		"X-Line-Application": LINE_APP,
		"X-Line-Access":      authToken,
	}
	client, transport, err := NewThriftClientForLP(HOST+TALKSERVICE_ENDPOINT, headers)
	talk := linethrift.NewTalkServiceClient(client)
	talk.AuthToken = authToken
	if err != nil {
		return nil, nil, err
	}
	return talk, transport, nil
}
