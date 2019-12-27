package tradfri

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dustin/go-coap"
	"github.com/pion/dtls/v2"
)

type Tradfri struct {
	client *dtls.Conn

	msgID uint32
	mux   sync.Mutex
}

func New(address, id string, psk []byte) (*Tradfri, error) {
	if id == "" {
		return nil, errors.New("invalid id")
	}

	if len(psk) == 0 {
		return nil, errors.New("invalid psk")
	}

	if !strings.Contains(address, ":") {
		address = address + ":5684"
	}

	addr, err := net.ResolveUDPAddr("udp", address)

	if err != nil {
		return nil, err
	}

	config := &dtls.Config{
		PSKIdentityHint: []byte(id),
		PSK: func(hint []byte) ([]byte, error) {
			return psk, nil
		},
		ConnectTimeout: dtls.ConnectTimeoutOption(30 * time.Second),
		CipherSuites:   []dtls.CipherSuiteID{dtls.TLS_PSK_WITH_AES_128_CCM_8},
	}

	client, err := dtls.Dial("udp", addr, config)

	if err != nil {
		return nil, err
	}

	return &Tradfri{
		client: client,
	}, nil
}

func (t *Tradfri) Close() {
	if t == nil {
		return
	}

	if t.client != nil {
		t.client.Close()
	}
}

func (t *Tradfri) Info() (*GatewayInfo, error) {
	var gateway GatewayInfo

	if err := t.GetAsJson("/15011/15012", &gateway); err != nil {
		return nil, err
	}

	return &gateway, nil
}

func (t *Tradfri) Devices() ([]string, error) {
	return t.GetAsIDs("/15001")
}

func (t *Tradfri) Device(deviceID string) (*DeviceInfo, error) {
	var device DeviceInfo

	if err := t.GetAsJson("/15001/"+deviceID, &device); err != nil {
		return nil, err
	}

	return &device, nil
}

func (t *Tradfri) UpdateDevice(deviceID string, settings DeviceSettings) error {
	return t.PutJsonChange("/15001/"+deviceID, settings)
}

func (t *Tradfri) Groups() ([]string, error) {
	return t.GetAsIDs("/15004")
}

func (t *Tradfri) Group(groupID string) (*GroupInfo, error) {
	var group GroupInfo

	if err := t.GetAsJson("/15004/"+groupID, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

func (t *Tradfri) Scenes(groupID string) ([]string, error) {
	return t.GetAsIDs("/15005/" + groupID)
}

func (t *Tradfri) Scene(groupID, sceneID string) (*SceneInfo, error) {
	var scene SceneInfo

	if err := t.GetAsJson("/15005/"+groupID+"/"+sceneID, &scene); err != nil {
		return nil, err
	}

	return &scene, nil
}

func (t *Tradfri) MessageID() uint16 {
	return uint16(atomic.AddUint32(&t.msgID, 1) % 0xffff)
}

func (t *Tradfri) RoundTrip(request coap.Message) (*coap.Message, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	payload, err := request.MarshalBinary()

	if err != nil {
		return nil, err
	}

	if _, err = t.client.Write(payload); err != nil {
		return nil, err
	}

	data := make([]byte, 65*1024)

	count, err := t.client.Read(data)

	if err != nil {
		return nil, err
	}

	data = append([]byte(nil), data[:count]...)

	message, err := coap.ParseMessage(data)

	if err != nil {
		return nil, err
	}

	//println(string(message.Payload))

	return &message, nil
}

func (t *Tradfri) GetAsIDs(path string) ([]string, error) {
	request := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: t.MessageID(),
	}

	request.SetPathString(path)

	message, err := t.RoundTrip(request)

	if err != nil {
		return nil, err
	}

	if message.Code != coap.Content {
		return nil, errors.New("response is not of type content")
	}

	var data []int

	if err := json.Unmarshal(message.Payload, &data); err != nil {
		return nil, err
	}

	ids := make([]string, 0)

	for _, id := range data {
		ids = append(ids, fmt.Sprintf("%d", id))
	}

	return ids, nil
}

func (t *Tradfri) GetAsJson(path string, out interface{}) error {
	request := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: t.MessageID(),
	}

	request.SetPathString(path)

	message, err := t.RoundTrip(request)

	if err != nil {
		return err
	}

	if message.Code != coap.Content {
		return errors.New("response is not of type content")
	}

	if err := json.Unmarshal(message.Payload, &out); err != nil {
		return err
	}

	return nil
}

func (t *Tradfri) PutJsonChange(path string, data interface{}) error {
	payload, err := json.Marshal(data)

	if err != nil {
		return err
	}

	request := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.PUT,
		MessageID: t.MessageID(),

		Payload: payload,
	}

	request.SetPathString(path)

	message, err := t.RoundTrip(request)

	if err != nil {
		return err
	}

	if message.Code != coap.Changed {
		return errors.New("response is not of type changed")
	}

	return nil
}