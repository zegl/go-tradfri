package tradfri

import (
	"encoding/json"
	"errors"

	"github.com/dustin/go-coap"
)

func PSK(address, id, code string) ([]byte, error) {
	if code == "" {
		return nil, errors.New("invalid security code")
	}

	t, err := New(address, "Client_identity", []byte(code))

	if err != nil {
		return nil, err
	}

	defer t.Close()

	type requestType struct {
		ID string `json:"9090"`
	}

	type responseTyype struct {
		PSK string `json:"9091"`
	}

	payload, err := json.Marshal(requestType{
		ID: id,
	})

	request := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.POST,
		MessageID: t.MessageID(),
		Payload:   payload,
	}

	request.SetPathString("/15011/9063")

	response, err := t.RoundTrip(request)

	if err != nil {
		return nil, err
	}

	if response.Code != coap.Created {
		return nil, errors.New("response is not of type created")
	}

	var data responseTyype

	if err := json.Unmarshal(response.Payload, &data); err != nil {
		return nil, err
	}

	if len(data.PSK) == 0 {
		return nil, errors.New("unexpected response")
	}

	return []byte(data.PSK), nil
}
