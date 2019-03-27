package subscription

import (
	"encoding/json"
	"fmt"

	"github.com/ably/ably-go/ably"
)

// An AblySubscription represents a subscription to an Ably REST channel.
type AblySubscription struct {
	// Channel represents the Ably REST channel associated to the
	// subscription's topic.
	channel *ably.RestChannel
}

// NewAblySubscription creates a new subscription to the given channel.
func NewAblySubscription(channel *ably.RestChannel) (Subscription, error) {
	if channel == nil {
		return nil, fmt.Errorf("subscription.AblySubscription: channel cannot be nil")
	}

	return &AblySubscription{channel}, nil
}

// Publish sends a message on the subscription's topic.
func (s *AblySubscription) Publish(msg *Message) error {
	if msg == nil {
		return fmt.Errorf("subscription.AblySubscription [topic=%s]: message cannot be nil", s.Topic())
	}

	payload, err := json.Marshal(msg.Data)
	if err != nil {
		return fmt.Errorf("subscription.AblySubscription [topic=%s]: failed to marshal message (%s)", s.Topic(), err)
	}

	err = s.channel.Publish(msg.Type, string(payload))
	if err != nil {
		return fmt.Errorf("subscription.AblySubscription [topic=%s]: failed to publish message (%s)", s.Topic(), err)
	}
	return nil
}

// Topic returns the subscriptions's topic.
func (s *AblySubscription) Topic() string {
	return s.channel.Name
}
