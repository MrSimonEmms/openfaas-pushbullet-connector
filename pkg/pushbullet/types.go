package pushbullet

import (
	"encoding/json"
)

type StreamEventType string

const (
	StreamEventTypeNOP    StreamEventType = "nop"
	StreamEventTypeTickle StreamEventType = "tickle"
	StreamEventTypePush   StreamEventType = "push"
)

type StreamEventSubtype string

const (
	StreamEventSubtypePush   StreamEventSubtype = "push"   // A change to the /v2/pushes resources
	StreamEventSubtypeDevice StreamEventSubtype = "device" // A change to the /v2/devices resources - not supported
)

type StreamEvent struct {
	Type    StreamEventType    `json:"type"`
	SubType StreamEventSubtype `json:"subtype,omitempty"`
}

type Subscriptions struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type Subscription struct {
	Iden     string  `json:"iden"`
	Active   bool    `json:"active"`
	Created  float32 `json:"created"`
	Modified float32 `json:"modified"`
	Channel  Channel `json:"channel,omitempty"`
}

type Channel struct {
	Iden string `json:"iden"`
	Tag  string `json:"tag"`
	Name string `json:"name"`
}

type PushType string

const (
	PushTypeNote PushType = "note"
	PushTypeLink PushType = "link"
)

type Push struct {
	Iden                    string   `json:"iden"`
	Active                  bool     `json:"active"`
	Created                 float32  `json:"created"`
	Modified                float32  `json:"modified"`
	Type                    PushType `json:"type"`
	Dismissed               bool     `json:"dismissed"`
	GUID                    string   `json:"guid"`
	Direction               string   `json:"direction"`
	SenderIden              string   `json:"sender_iden"`
	SenderEmail             string   `json:"sender_email"`
	SenderEmailNormalized   string   `json:"sender_email_normalized"`
	SenderName              string   `json:"sender_name"`
	ReceiverIden            string   `json:"receiver_iden"`
	ReceiverEmail           string   `json:"receiver_email"`
	ReceiverEmailNormalized string   `json:"receiver_email_normalized"`
	TargetDeviceIden        string   `json:"target_device_iden"`
	SourceDeviceIden        string   `json:"source_device_iden"`
	Body                    string   `json:"body,omitempty"`
	Title                   string   `json:"title,omitempty"`
	URL                     string   `json:"url,omitempty"`
	ChannelID               string   `json:"channel_iden,omitempty"`
	payload                 string
}

func (push *Push) GeneratePayload() error {
	payload, err := json.Marshal(push)
	if err != nil {
		return err
	}

	push.payload = string(payload)
	return nil
}

func (push *Push) Payload() string {
	return push.payload
}

type Pushes struct {
	Pushes []Push `json:"pushes"`
}

type PushAndTag struct {
	Push
	Tag string
}
