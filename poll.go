package lineapi

import "github.com/mopeneko/linethrift"

type PollingController struct {
	client, poll *linethrift.TalkServiceClient
}

func NewPollingController(talk *linethrift.TalkServiceClient) (*PollingController, error) {
	headers := map[string]string{
		"User-Agent": USER_AGENT,
		"X-Line-Application": LINE_APP,
		"X-Line-Access": authToken,
		"x-lpqs": "/P4",
		"X-LPV": "1",
		"X-LHM": "POST",
	}
	client, err := NewThriftClient(HOST + POLLING_ENDPOINT, headers)
	poll := linethrift.NewTalkServiceClient(client)
	if err != nil {
		return nil, err
	}
	pollcon := &PollingController(talk, poll)
	return pollcon, nil
}
