package udspubsub

import (
	"os"
	"testing"
	"time"
)

var ps *PubSubHandler
var clientOne, clientTwo, clientThree *PubSubClient
var clientOneChannel, clientTwoChannel, clientThreeChannel chan string
var socketPath = "/tmp/udspubsub.sock"

func TestMain(m *testing.M) {
	ps = Listener(socketPath)
	code := m.Run()
	os.Remove(socketPath)
	<-time.After(time.Second * 1)
	os.Exit(code)
}

func TestPSHandler(t *testing.T) {
	if ps == nil {
		t.Fatal("PubSubHandler is not initialized")
	}
}

func TestPSClientOneCreation(t *testing.T) {
	c, err := NewPubSubClient(socketPath)
	clientOne = c
	if err != nil {
		t.Fatalf("Client Creation failed : %s", err.Error())
	}
	clientTwo, _ = NewPubSubClient(socketPath)
	clientThree, _ = NewPubSubClient(socketPath)

	if len(ps.Clients) !=3 {
		t.Fatal("Client count must be 3")
	}
}

func TestSubscription(t *testing.T) {
	chanOne, err := clientOne.Subscribe("TopicOne")
	if err != nil {
		t.Fatalf("Client Subscription failed : %s", err.Error())
	}
	clientOneChannel = chanOne
	<-time.After(time.Second * 1)
	clientTwoChannel, err = clientTwo.Subscribe("TopicTwo")
	if err != nil {
		t.Fatalf("Client Subscription failed : %s", err.Error())
	}
	<-time.After(time.Second * 1)
	if len(ps.Subscriptions) != 2 {
		t.Fatalf("Subscription count should be two, its %d", len(ps.Subscriptions))
	}
}

func TestClientThreePublishClientOne(t *testing.T) {
	m := "This is a message to Topic one"

	err := clientThree.Publish("TopicOne", []byte(m))
	if err != nil {
		t.Fatalf("Publish failed : %s", err.Error())
	}
	message := <-clientOneChannel
	if m != message {
		t.Fatal("Got wrong message")
	}
	<-time.After(time.Second * 1)
}

func TestClientThreePublishClientTwo(t *testing.T) {
	m := "This is a message to Topic two"
	err := clientThree.Publish("TopicTwo", []byte(m))
	if err != nil {
		t.Fatalf("Publish failed : %s", err.Error())
	}
	message := <-clientTwoChannel
	if m != message {
		t.Fatal("Got wrong message")
	}
	<-time.After(time.Second * 1)
}

func TestClientThreeSubscribe(t *testing.T) {
	var err error
	clientThreeChannel, err = clientThree.Subscribe("TopicThree")
	if err != nil {
		t.Fatalf("Client Subscription failed : %s", err.Error())
	}
	m := "This is a message to Topic three"
	go func() {
		message := <-clientThreeChannel
		if m != message {
			t.Fatal("Got wrong message")
		}
	}()
	err = clientTwo.Publish("TopicThree", []byte(m))
	if err != nil {
		t.Fatalf("Publish failed : %s", err.Error())
	}
	<-time.After(time.Second * 1)
}

func TestClientUnSubscribe(t *testing.T) {
	var err error
	err = clientThree.Unsubscribe("TopicThree")
	if err != nil {
		t.Fatalf("Client Unsubscribe failed : %s", err.Error())
	}

	<-time.After(time.Second * 1)
	if len(ps.Subscriptions) != 2 {
		t.Fatalf("Subscription count should be two, its %d", len(ps.Subscriptions))
	}

	err = clientTwo.Unsubscribe("TopicTwo")
	if err != nil {
		t.Fatalf("Client Unsubscribe failed : %s", err.Error())
	}

	<-time.After(time.Second * 1)
	if len(ps.Subscriptions) != 1 {
		t.Fatalf("Subscription count should be one, its %d", len(ps.Subscriptions))
	}

	err = clientOne.Unsubscribe("TopicOne")
	if err != nil {
		t.Fatalf("Client Unsubscribe failed : %s", err.Error())
	}

	<-time.After(time.Second * 1)
	if len(ps.Subscriptions) != 0 {
		t.Fatalf("Subscription count should be zero, its %d", len(ps.Subscriptions))
	}
}