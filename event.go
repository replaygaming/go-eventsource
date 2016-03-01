package main

import "encoding/json"
import "log"

type payload struct {
	Event    string
	Data     *json.RawMessage
	Channels []string
}

// event name is packed inside the data struct to allow the client to subscribe
// to all events on .addEventListener("message") without knowning all the possible
// event names
type event struct {
	Event string           `json:"event"`
	Data  *json.RawMessage `json:"data"`
}

func parseMessage(message []byte) ([]byte, []string, error) {
	p := payload{}
	err := json.Unmarshal(message, &p)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("[INFO] parseMessage")
	log.Printf("[INFO]     Incomming - event=%s, data=%s", p.Event, p.Data)

	e := event{Event: p.Event, Data: p.Data}
	data, err := json.Marshal(e)
	if err != nil {
		return nil, nil, err
	}

	// defaults messages without channels to '*' channel
	if len(p.Channels) == 0 {
		p.Channels = []string{"*"}
	}

	log.Printf("[INFO]     Outgoing - %s", data)

	return data, p.Channels, nil
}
