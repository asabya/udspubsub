package udspubsub

import (
	"os"
	"testing"
)

var ps *PubSubHandler
var ClientOne, clientTwo, clientThree *PubSubClient
var socketPath = "/tmp/udspubsub.sock"

func TestMain(m *testing.M) {
	ps = Listener(socketPath)
	code := m.Run()
	os.Remove(socketPath)
	os.Exit(code)
}

func TestPSHandler(t *testing.T) {
	if ps == nil {
		t.Fatalf("PubSubHandler is not initialized")
	}
}

func TestPSClientOneCreation(t *testing.T) {
	clientOne, err := NewPubSubClient(socketPath)
	ClientOne = clientOne
	if err != nil {
		t.Fatalf("Client Creation failed : %s", err.Error())
	}
	clientTwo, _ = NewPubSubClient(socketPath)
	clientThree, _ = NewPubSubClient(socketPath)
}
