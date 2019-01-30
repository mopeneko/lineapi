package lineapi

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/mopeneko/linethrift"
)

// Generate LINE TalkService client
func NewLineClient(authToken string) (*linethrift.TalkServiceClient, error) {
	headers := map[string]string{
		"User-Agent": USER_AGENT,
		"X-Line-Application": LINE_APP,
		"X-Line-Access": authToken,
		"x-lpqs": "/S4",
		"X-LPV": "1",
		"X-LHM": "POST",
	}
	client, err := NewThriftClient(HOST + TALKSERVICE_ENDPOINT, headers)
	talk := linethrift.NewTalkServiceClient(client)
	if err != nil {
		return nil, err
	}
	return talk, nil
}
