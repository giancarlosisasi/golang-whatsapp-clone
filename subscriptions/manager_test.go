package subscriptions

import (
	"fmt"
	"golang-whatsapp-clone/graph/model"
	"golang-whatsapp-clone/repository"
	"testing"
	"time"
)

func TestBasicMessageSubscription(t *testing.T) {
	fmt.Println("=== Testing message added subscription")

	sm := NewSubscriptionManager()
	conversationID := "conv-123"

	fmt.Println("1. Created subscription manager")

	fmt.Printf("2. User A subscribes to conversation %s\n", conversationID)
	userAChan := sm.SubscribeToMessages(conversationID)

	count := sm.GetSubscriberCount(conversationID)
	if count != 1 {
		t.Errorf("Expected 1 subscriber, got %d", count)
	}
	fmt.Printf("   ✓ Subscriber count: %d\n", count)

	fmt.Printf("3. User A starts listening for messages...\n")
	messageReceived := make(chan bool) // to sync our test

	go func() {
		fmt.Printf("    User A's client is now listening")

		// wait for a message on the channel
		// it will block until a message is received
		msg := <-userAChan

		fmt.Printf("   📨 User A received message: '%s' from %s\n",
			msg.Content, msg.SenderUserID)

		messageReceived <- true
	}()

	// wait for a moment for the goroutine to start listening
	time.Sleep(100 * time.Millisecond)

	fmt.Println("4. User B sends message to the conversation")
	testMessage := &model.MessageAddedEvent{
		ID:           "msg-001",
		Content:      "HI! lets play!",
		SenderUserID: "user-b",
		MessageType:  model.MessageTypeEnum(repository.MESSAGE_TYPE_TEXT),
	}

	sm.BroadcastMessage(conversationID, testMessage)

	// wait for the message to be received (with timeout)
	select {
	case <-messageReceived:
		fmt.Println("   ✓  Message successfully delivered to User A!")
	case <-time.After(2 * time.Second):
		t.Error("Timeout: Message was not received within 2 seconds")
	}

	fmt.Println(">>>> ✅ Test completed successfully!")
}

func TestOneONONeChatScenario(t *testing.T) {
	sm := NewSubscriptionManager()
	conversationID := "conv-456"

	userAChan := sm.SubscribeToMessages(conversationID)
	userBChan := sm.SubscribeToMessages(conversationID)

	count := sm.GetSubscriberCount(conversationID)
	if count != 2 {
		t.Errorf("Expected 2 subscribers, got %d", count)
	}
	fmt.Printf("   ✅ Both users subscribed. Total subscribers: %d\n", count)

	userAReceived := 0
	userBReceived := 0
	testComplete := make(chan bool)

	go func() {
		for msg := range userAChan {
			userAReceived++
			fmt.Printf("    > User A received: '%s'\n", msg.Content)

			if userAReceived == 1 {
				testComplete <- true
			}
		}
	}()
	go func() {
		for msg := range userBChan {
			userBReceived++
			fmt.Printf("    > User B received: '%s'\n", msg.Content)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	messageFromA := &model.MessageAddedEvent{
		ID:           "msg-002",
		Content:      "Hey! welcome!",
		SenderUserID: "user-a",
		MessageType:  model.MessageTypeEnum(repository.MESSAGE_TYPE_TEXT),
	}

	// this will send the message to the A and B user channels
	sm.BroadcastMessage(conversationID, messageFromA)

	// wait for message to be received
	select {
	case <-testComplete:
		fmt.Println("    ✅ Message delivered to both users!")
		fmt.Printf("     > User A received %d messages, User B received %d messages\n",
			userAReceived, userBReceived)
	case <-time.After(2 * time.Second):
		t.Errorf("Timeout: Messages were not received")
	}

	fmt.Println(">>>> ✅ 1-on-1 Chat test completed! ===")
}
