package subscriptions

import (
	"fmt"
	"golang-whatsapp-clone/graph/model"
	"sync"
)

// SubscriptionManager handles real-time message subscriptions
type SubscriptionManager struct {
	// Map of conversationID -> list of channels listening to that conversation
	//
	// EXAMPLE 1: 1-on-1 chat between User A and User B
	// "conv-123" -> [channelA, channelB]  // Both users listening to same conversation
	//
	// EXAMPLE 2: Group chat with 3 users
	// "conv-456" -> [channelA, channelB, channelC]  // All 3 users listening
	//
	// EXAMPLE 3: User with multiple browser tabs/devices
	// "conv-123" -> [channelA1, channelA2, channelB]  // User A has 2 connections, User B has 1
	//
	// IMPORTANT: Each channel = one client connection (browser tab, mobile app, etc.)
	//            NOT one channel per user!
	messageSubscribers map[string][]chan *model.MessageAddedEvent

	// Mutex to protect concurrent access to the subscribers map
	mutex sync.RWMutex
}

// NewSubscriptionManager creates a new subscription manager
func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		messageSubscribers: make(map[string][]chan *model.MessageAddedEvent),
	}
}

// SubscribeToMessages creates a subscription for new messages in a conversation
// Returns a channel that will receive Message objects for the given conversationID
//
// EXAMPLE USAGE:
// - User A opens WhatsApp web -> calls SubscribeToMessages("conv-123")
// - User B opens mobile app -> calls SubscribeToMessages("conv-123")
// - Now both will receive messages sent to "conv-123"
//
// MULTIPLE DEVICES:
// - User A opens 2 browser tabs -> 2 calls to SubscribeToMessages("conv-123")
// - User A will receive the same message on both tabs
func (sm *SubscriptionManager) SubscribeToMessages(conversationID string) <-chan *model.MessageAddedEvent {
	// Create a buffered channel to prevent blocking
	// Buffer size of 10 means we can queue up to 10 messages
	ch := make(chan *model.MessageAddedEvent, 10)

	// Lock for writing to the map
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Add this channel to the list of subscribers for this conversation
	sm.messageSubscribers[conversationID] = append(sm.messageSubscribers[conversationID], ch)

	fmt.Printf("New subscriber added to conversation %s. Total subscribers: %d\n",
		conversationID, len(sm.messageSubscribers[conversationID]))

	// Return read-only channel
	return ch
}

// BroadcastMessage sends a message to all subscribers of a conversation
//
// EXAMPLE: User A sends "Hello!" to conversation "conv-123"
// WHAT HAPPENS:
// 1. Message gets saved to database
// 2. BroadcastMessage() is called with the new message
// 3. Message is sent to ALL channels subscribed to "conv-123"
// 4. User A receives it (for message confirmation/UI update)
// 5. User B receives it (new message notification)
// 6. Any other devices/tabs also receive it
func (sm *SubscriptionManager) BroadcastMessage(conversationID string, message *model.MessageAddedEvent) {
	// Lock for reading from the map
	sm.mutex.RLock()
	subscribers := sm.messageSubscribers[conversationID]
	sm.mutex.RUnlock()

	fmt.Printf("Broadcasting message to %d subscribers of conversation %s\n",
		len(subscribers), conversationID)

	// Send message to all subscribers
	for i, ch := range subscribers {
		select {
		case ch <- message:
			fmt.Printf("Message sent to subscriber %d\n", i+1)
		default:
			// Channel is full or closed, skip this subscriber
			fmt.Printf("Failed to send message to subscriber %d (channel full or closed)\n", i+1)
		}
	}
}

// GetSubscriberCount returns the number of active subscribers for a conversation
func (sm *SubscriptionManager) GetSubscriberCount(conversationID string) int {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	return len(sm.messageSubscribers[conversationID])
}

func (sm *SubscriptionManager) Unsubscribe(conversationID string, ch <-chan *model.MessageAddedEvent) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	subscribers := sm.messageSubscribers[conversationID]

	for i, subscriber := range subscribers {
		if subscriber == ch {
			close(subscriber)

			// remove from slice by combining the parts before and after
			sm.messageSubscribers[conversationID] = append(subscribers[:i], subscribers[i+1:]...)

			fmt.Printf("Channel unsubscribed from conversation %s. Remaining:%d\n", conversationID, len(sm.messageSubscribers[conversationID]))

			if len(sm.messageSubscribers[conversationID]) == 0 {
				delete(sm.messageSubscribers, conversationID)
				fmt.Printf("No more subscribers for conversation %s. Cleaned up map entry.\n", conversationID)
			}

			return
		}
	}

	fmt.Printf("Warning: channel not found for conversation %s during unsubscribe\n", conversationID)
}

func (sm *SubscriptionManager) UnsubscribeAll(conversationID string) int {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	subscribers := sm.messageSubscribers[conversationID]
	count := len(subscribers)

	// close all channels
	for _, ch := range subscribers {
		close(ch)
	}

	delete(sm.messageSubscribers, conversationID)

	fmt.Printf("Unsubscribed all %d channels for conversation %s\n", count, conversationID)

	return count
}

func (sm *SubscriptionManager) GetAllConversationIDs() []string {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	ids := make([]string, 0, len(sm.messageSubscribers))

	for id := range sm.messageSubscribers {
		ids = append(ids, id)
	}

	return ids
}

func (sm *SubscriptionManager) GetTotalSubscribers() int {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	total := 0
	for _, subscribers := range sm.messageSubscribers {
		total += len(subscribers)
	}

	return total
}
