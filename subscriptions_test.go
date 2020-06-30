package convertkit_test

import (
	"testing"

	"github.com/joncalhoun/convertkit"
)

func TestClient_SubscribeToForm(t *testing.T) {
	c := client(t, "fake-secret-key")
	t.Run("basic check", func(t *testing.T) {
		email := "jonsnow@example.com"
		resp, err := c.SubscribeToForm(convertkit.SubscribeToFormRequest{
			FormID: 213,
			Email:  email,
		})
		if err != nil {
			t.Fatalf("SubscribeToForm() err = %v; want %v", err, nil)
		}
		if resp.Subscription.Subscriber.ID != 1 {
			t.Errorf("Subscriber.ID = %v; want 1", resp.Subscription.Subscriber.ID)
		}
	})
}

func TestClient_FormSubscriptions(t *testing.T) {
	c := client(t, "fake-secret-key")
	t.Run("basic check", func(t *testing.T) {
		resp, err := c.FormSubscriptions(convertkit.FormSubscriptionsRequest{
			FormID: 213,
		})
		if err != nil {
			t.Fatalf("FormSubscriptions() err = %v; want %v", err, nil)
		}
		if resp.Subscriptions[1].Subscriber.ID != 2 {
			t.Errorf("resp.Subscriptions[1].Subscriber.ID = %v; want 2", resp.Subscriptions[1].Subscriber.ID)
		}
	})
}

func TestClient_SubscribeToSequence(t *testing.T) {
	c := client(t, "fake-secret-key")
	t.Run("basic check", func(t *testing.T) {
		email := "jonsnow@example.com"
		resp, err := c.SubscribeToSequence(convertkit.SubscribeToSequenceRequest{
			SequenceID: 55,
			Email:      email,
		})
		if err != nil {
			t.Fatalf("SubscribeToSequence() err = %v; want %v", err, nil)
		}
		if resp.Subscription.Subscriber.ID != 1 {
			t.Errorf("Subscriber.ID = %v; want 1", resp.Subscription.Subscriber.ID)
		}
	})
}

func TestClient_SequenceSubscriptions(t *testing.T) {
	c := client(t, "fake-secret-key")
	t.Run("basic check", func(t *testing.T) {
		resp, err := c.SequenceSubscriptions(convertkit.SequenceSubscriptionsRequest{
			SequenceID: 55,
		})
		if err != nil {
			t.Fatalf("SequenceSubscriptions() err = %v; want %v", err, nil)
		}
		if resp.Subscriptions[1].Subscriber.ID != 2 {
			t.Errorf("resp.Subscriptions[1].Subscriber.ID = %v; want 2", resp.Subscriptions[1].Subscriber.ID)
		}
	})
}

func TestClient_TagSubscriberSequence(t *testing.T) {
	c := client(t, "fake-secret-key")
	t.Run("basic check", func(t *testing.T) {
		email := "jonsnow@example.com"
		resp, err := c.TagSubscriber(convertkit.TagSubscriberRequest{
			TagID: 14,
			Email: email,
		})
		if err != nil {
			t.Fatalf("TagSubscriber() err = %v; want %v", err, nil)
		}
		if resp.Subscription.Subscriber.ID != 1 {
			t.Errorf("Subscriber.ID = %v; want 1", resp.Subscription.Subscriber.ID)
		}
	})
}

func TestClient_UntagSubscriberSequence(t *testing.T) {
	c := client(t, "fake-secret-key")
	t.Run("basic check", func(t *testing.T) {
		_, err := c.UntagSubscriber(convertkit.UntagSubscriberRequest{
			SubscriberID: 88,
			TagID:        71,
		})
		if err != nil {
			t.Fatalf("UntagSubscriber() err = %v; want %v", err, nil)
		}
	})
}
