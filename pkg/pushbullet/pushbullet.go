package pushbullet

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var StreamAddress = url.URL{
	Scheme: "wss",
	Host:   "stream.pushbullet.com",
	Path:   "/websocket",
}

var RESTAddress = url.URL{
	Scheme: "https",
	Host:   "api.pushbullet.com",
	Path:   "v2",
}

type MessageHandler func(Pushbullet, string, Push)

type Pushbullet struct {
	apiToken        string
	websocketClient *websocket.Conn
	httpClient      http.Client
	MessageHandler  MessageHandler
	modified        time.Time
}

func (pb *Pushbullet) Close() error {
	log.Print("Closing PushBullet connection")
	return pb.websocketClient.Close()
}

func (pb *Pushbullet) SetHandler(handler MessageHandler) *Pushbullet {
	pb.MessageHandler = handler

	return pb
}

func (pb *Pushbullet) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", RESTAddress.String(), url), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Access-Token", pb.apiToken)

	return pb.httpClient.Do(req)
}

func (pb *Pushbullet) GetChannelById(id string) (*Channel, error) {
	res, err := pb.get("/subscriptions")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var subscriptions Subscriptions
	err = json.Unmarshal(body, &subscriptions)
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}

	var channel *Channel
	for _, s := range subscriptions.Subscriptions {
		if s.Active && s.Channel.Iden == id {
			channel = &s.Channel
			break
		}
	}

	if channel == nil {
		err = fmt.Errorf("unknown channel ID %s", id)
		log.Printf("error: %s", err)
		return nil, err
	}

	return channel, nil
}

// This is unlikely to receive more than one push, but not impossible which is why we return a slice
func (pb *Pushbullet) Pushes() ([]PushAndTag, error) {
	res, err := pb.get(fmt.Sprintf("/pushes?modified_after=%d", pb.modified.Unix()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var pushes Pushes
	err = json.Unmarshal(body, &pushes)
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}

	var pushAndTags []PushAndTag
	for _, p := range pushes.Pushes {
		err := p.GeneratePayload()
		if err != nil {
			return nil, err
		}

		if channelId := p.ChannelID; channelId != "" {
			channel, err := pb.GetChannelById(p.ChannelID)
			if err != nil {
				return nil, err
			}

			pushAndTags = append(pushAndTags, PushAndTag{
				Push: p,
				Tag:  channel.Tag,
			})
		}
	}

	return pushAndTags, nil
}

func New(apiToken string) (*Pushbullet, error) {
	log.Printf("Connecting to websocket: %s", StreamAddress.String())

	// Append the API token to the path
	StreamAddress.Path += fmt.Sprintf("/%s", apiToken)

	c, _, err := websocket.DefaultDialer.Dial(StreamAddress.String(), nil)
	if err != nil {
		return nil, err
	}

	pb := &Pushbullet{
		apiToken:        apiToken,
		websocketClient: c,
		httpClient: http.Client{
			Timeout: 5 * time.Second,
		},
	}

	pb.modified = time.Now()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := pb.websocketClient.ReadMessage()
			if err != nil {
				log.Printf("error: %s", err)
				return
			}

			var streamEvent StreamEvent
			err = json.Unmarshal(message, &streamEvent)
			if err != nil {
				log.Printf("error: %s", err)
				return
			}

			switch streamEvent.Type {
			case StreamEventTypeNOP:
				log.Print("Health ping received")
			case StreamEventTypePush:
				log.Printf("Ephemeral push received: %s", streamEvent.SubType)
			case StreamEventTypeTickle:
				log.Printf("Tickle received: %s", streamEvent.SubType)

				if streamEvent.SubType == StreamEventSubtypePush {
					pushes, err := pb.Pushes()
					if err != nil {
						log.Printf("push err: %s", err)
					} else {
						for _, p := range pushes {
							log.Printf("new push to %s: %v", p.Tag, p.Payload())
							pb.MessageHandler(*pb, p.Tag, p.Push)
						}
					}
				}
			}

			log.Print("Updating last modified time")
			pb.modified = time.Now()
		}
	}()

	return pb, nil
}
