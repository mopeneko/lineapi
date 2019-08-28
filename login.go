package lineapi

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/mopeneko/linethrift"
)

// Generate LINE TalkService client
func NewLineClient(authToken string, options thrift.THttpClientOptions) (*linethrift.TalkServiceClient, error) {
	headers := map[string]string{
		"User-Agent":         USER_AGENT,
		"X-Line-Application": LINE_APP,
		"X-Line-Access":      authToken,
	}
	client, err := NewThriftClient(HOST+TALKSERVICE_ENDPOINT, headers, options)
	talk := linethrift.NewTalkServiceClient(client)
	talk.AuthToken = authToken
	if err != nil {
		return nil, err
	}
	return talk, nil
}
