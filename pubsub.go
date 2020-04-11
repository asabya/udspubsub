/*
This is inspired by https://github.com/tabvn/golang-pubsub-youtube
 */

package udspubsub

import (
	"encoding/json"
	"log"
	"net"
)

const (
	PUBLISH     = "publish"
	SUBSCRIBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"
)

type PubSubHandler struct {
	Clients       []Client
	Subscriptions []Subscription
}

type Client struct {
	Id         string
	Connection net.Conn
}

type Message struct {
	Action  string          `json:"action"`
	Topic   string          `json:"topic"`
	Message []byte 			`json:"message"`
}

type Subscription struct {
	Topic  string
	Client *Client
}

var Topics = []string{}


func (ps *PubSubHandler) handleConnection(c net.Conn) Client {
	client := Client{
		Id:         autoId(),
		Connection: c,
	}
	ps.addClient(client)
	return client
}


func (ps *PubSubHandler) addClient(client Client) error {
	ps.Clients = append(ps.Clients, client)
	payload := []byte(client.Id)
	_, err := client.Connection.Write(payload)
	return err
}

func (ps *PubSubHandler) removeClient(client Client) {
	for index, sub := range ps.Subscriptions {
		if client.Id == sub.Client.Id {
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}
	for index, c := range ps.Clients {
		if c.Id == client.Id {
			ps.Clients = append(ps.Clients[:index], ps.Clients[index+1:]...)
		}
	}
}

func (ps *PubSubHandler) GetTopicSubscriptions(topic string) ([]Subscription) {
	var subscriptionList []Subscription
	for _, subscription := range ps.Subscriptions {
		if subscription.Topic == topic {
			subscriptionList = append(subscriptionList, subscription)
		}
	}
	return subscriptionList
}

func (ps *PubSubHandler) GetClientSubscriptions(client *Client) ([]Subscription) {
	var subscriptionList []Subscription
	for _, subscription := range ps.Subscriptions {
			if subscription.Client.Id == client.Id{
				subscriptionList = append(subscriptionList, subscription)
			}
	}
	return subscriptionList
}

func (ps *PubSubHandler) IsClientSubscribed(topic string, client *Client) bool {
	for _, subscription := range ps.Subscriptions {
		if subscription.Client.Id == client.Id && subscription.Topic == topic {
			return true
		}
	}

	return false
}

func (ps *PubSubHandler) subscribe(client *Client, topic string) {
	if !ps.IsClientSubscribed(topic, client) {
		newSubscription := Subscription{
			Topic:  topic,
			Client: client,
		}
		ps.Subscriptions = append(ps.Subscriptions, newSubscription)
	}
}

func (ps *PubSubHandler) publish(topic string, message []byte, excludeClient *Client) {
	subscriptions := ps.GetTopicSubscriptions(topic)
	for _, sub := range subscriptions {
		sub.Client.send(message)
	}
}

func (client *Client) send(message [] byte) (error) {
	_, err := client.Connection.Write(message)
	return err
}

func (ps *PubSubHandler) unsubscribe(client *Client, topic string) {
	for index, sub := range ps.Subscriptions {
		if sub.Client.Id == client.Id && sub.Topic == topic {
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}
}

func (ps *PubSubHandler) handleReceiveMessage(client Client, payload []byte) error {
	m := Message{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println("This is not correct message payload ", err.Error() )
		return err
	}

	switch m.Action {
	case PUBLISH:
		if !isTopicAvailable(m.Topic) {
			log.Println("No subscribers for this topic")
			break
		}
		ps.publish(m.Topic, m.Message, nil)
		break

	case SUBSCRIBE:
		if !isTopicAvailable(m.Topic) {
			Topics = append(Topics, m.Topic)
		}
		ps.subscribe(&client, m.Topic)
		break

	case UNSUBSCRIBE:
		ps.unsubscribe(&client, m.Topic)
		break

	default:
		break
	}
	return nil
}

func isTopicAvailable(topic string) bool {
	for _, t := range Topics {
		if t == topic {
			return true
		}
	}
	return false
}