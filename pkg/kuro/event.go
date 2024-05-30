package kuro

import (
	"encoding/json"
)

type EventCallback func(socket *Client, payload []byte)

type Event struct {
	Event   string `json:"event"`
	Payload []byte `json:"payload"`
}

func (e *Event) toBytes() []byte {
	var tempStruct struct {
		Event   string `json:"event"`
		Payload string `json:"payload"`
	}

	tempStruct.Event = e.Event
	tempStruct.Payload = string(e.Payload)

	output, err := json.Marshal(tempStruct)
	if err != nil {
		return []byte{}
	}
	return output
}

func bytesToEvent(data []byte) (*Event, error) {
	var output Event

	var tempStruct struct {
		Event   string `json:"event"`
		Payload any    `json:"payload"`
	}

	err := json.Unmarshal(data, &tempStruct)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(tempStruct.Payload)
	if err != nil {
		return nil, err
	}

	output.Event = tempStruct.Event
	output.Payload = bytes

	return &output, nil
}
