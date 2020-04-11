package udspubsub

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
)

const MAX_BUFFER = 1024

type PubSubClient struct {
	ID string
	c net.Conn
}

func autoId() string {
	return uuid.Must(uuid.NewV4(), nil).String()
}

func Listener(socketPath string) *PubSubHandler {
	ps := &PubSubHandler{}
	l, err := net.Listen("unix",socketPath)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			client := ps.handleConnection(conn)
			go func() {
				for {
					b := make([]byte, MAX_BUFFER)
					n, err := conn.Read(b)
					if err != nil {
						return
					}
					ps.handleReceiveMessage(client, b[:n])
					b = nil
				}
			}()
		}
	}()
	return ps
}

func NewPubSubClient(socketPath string) (*PubSubClient, error){
	c, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, err
	}
	b := make([]byte, 1024)
	n, err := c.Read(b)
	if err != nil {
		return nil, err
	}
	p := &PubSubClient{
		ID: string(b[:n]),
		c:  c,
	}
	b = nil
	return p, nil
}

func (p *PubSubClient) Subscribe(topic string) (chan string, error) {
	onMessage := make(chan string)
	payload, err := json.Marshal(Message{
		Action:  SUBSCRIBE,
		Topic:   topic,
	})
	if err != nil {
		return onMessage, err
	}
	_, err = p.c.Write(payload)
	if err != nil {
		return onMessage, err
	}
	go func() {
		defer p.c.Close()
		for {
			b := make([]byte, MAX_BUFFER)
			n, err := p.c.Read(b)
			if err != nil {
				return
			}
			onMessage <- string(b[:n])
			b = nil
		}
	}()
	return onMessage, nil
}

func (p *PubSubClient) Publish(topic string, message []byte) error {
	payload, err := json.Marshal(Message{
		Action:  PUBLISH,
		Topic:   topic,
		Message: message,
	})
	if err != nil {
		return err
	}
	_, err = p.c.Write(payload)
	if err != nil {
		return err
	}
	return nil
}

func (p *PubSubClient) Unsubscribe(topic string) error {
	payload, err := json.Marshal(Message{
		Action:  UNSUBSCRIBE,
		Topic:   topic,
	})
	if err != nil {
		return err
	}
	_, err = p.c.Write(payload)
	if err != nil {
		return err
	}
	return nil
}