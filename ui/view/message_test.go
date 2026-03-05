package view

import (
	"mushin/internal/audio"
	"testing"
	"time"

	"github.com/CoyAce/wi"
)

func TestMessagePersistence(t *testing.T) {
	wi.DefaultClient = &wi.Client{Identity: wi.Identity{UUID: "#00001"}}
	wi.Mkdir(GetDir(wi.DefaultClient.ID()))
	wi.RemoveFile(GetDataPath("message.log"))
	mk := MessageKeeper{MessageChannel: make(chan *Message, 1)}
	go mk.Loop()
	mk.MessageChannel <- &Message{TextControl: NewTextControl("hello world"), Contacts: Contacts{Sender: "test#00001", UUID: "#00001"}}
	mk.MessageChannel <- &Message{TextControl: NewTextControl("hello beautiful world"), Contacts: Contacts{Sender: "test#00001", UUID: "#00001"}}
	mk.Flush()
	messages := mk.Messages(audio.StreamConfig{})
	if len(messages) != 2 {
		t.Errorf("Messages length should be 2, but %d", len(messages))
	}
	if messages[0].UUID != "#00001" {
		t.Errorf("Messages[0].UUID should be #00001, but %s", messages[0].UUID)
	}
	if messages[1].Text != "hello beautiful world" {
		t.Errorf("Messages[1].Text != hello beautiful world, but %s", messages[1].Text)
	}
}

func TestMessageOrder(t *testing.T) {
	list := MessageList{Messages: make([]*Message, 0)}
	m2 := &Message{TextControl: NewTextControl("hello world"), Contacts: Contacts{Sender: "test#00001", UUID: "#00001"}, CreatedAt: time.Now()}
	m1 := &Message{TextControl: NewTextControl("hello beautiful world"), Contacts: Contacts{Sender: "test#00001", UUID: "#00001"}, CreatedAt: time.Now()}
	m1.AddTo(&list)
	m2.AddTo(&list)
	if list.Messages[0].Text != m2.Text || list.Messages[1].Text != m1.Text {
		t.Errorf("messages should be ordered by CreatedAt time")
	}
}
