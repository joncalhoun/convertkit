package convertkit

import (
	"fmt"
	"net/http"
	"time"
)

// Subscriber is a user subscribed to your mailing list. This type is returned
// from several API endpoints.
type Subscriber struct {
	ID        int               `json:"id"`
	FirstName string            `json:"first_name"`
	Email     string            `json:"email_address"`
	State     string            `json:"state"`
	CreatedAt time.Time         `json:"created_at"`
	Fields    map[string]string `json:"fields"`
}

// SubscribersRequest is used to narrow down the list of subscribers being
// returned. Please note that all time.Time fields are converted to "yyyy-mm-dd"
// when interacting with the API, so any finer time increments will be lost.
type SubscribersRequest struct {
	// Optional
	Page        int       `json:"page,omitempty"`
	From        *Date     `json:"from,omitempty"`
	To          *Date     `json:"to,omitempty"`
	UpdatedFrom *Date     `json:"updated_from,omitempty"`
	UpdatedTo   *Date     `json:"updated_to,omitempty"`
	SortOrder   SortOrder `json:"sort_order,omitempty"`
	SortField   string    `json:"sort_field,omitempty"`
	Email       string    `json:"email_address,omitempty"`
}

// SubscribersResponse is the data returned from a Subscribers call.
type SubscribersResponse struct {
	TotalSubscribers int          `json:"total_subscribers"`
	Page             int          `json:"page"`
	TotalPages       int          `json:"total_pages"`
	Subscribers      []Subscriber `json:"subscribers"`
}

// Subscribers lists subscribers for an account.
func (c *Client) Subscribers(req SubscribersRequest) (*SubscribersResponse, error) {
	var ret SubscribersResponse
	err := c.Do(http.MethodGet, "subscribers", req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// UpdateSubscriberRequest is used to update a subscriber.
type UpdateSubscriberRequest struct {
	// Required
	SubscriberID int `json:"-"`
	// Optional
	FirstName string            `json:"first_name,omitempty"`
	Email     string            `json:"email_address,omitempty"`
	Fields    map[string]string `json:"fields,omitempty"`
}

// UpdateSubscriberResponse is the data returned from a UpdateSubscriber call.
type UpdateSubscriberResponse struct {
	Subscriber `json:"subscriber"`
}

// UpdateSubscriber will update a subscriber's information.
func (c *Client) UpdateSubscriber(req UpdateSubscriberRequest) (*UpdateSubscriberResponse, error) {
	var ret UpdateSubscriberResponse
	err := c.Do(http.MethodPut, fmt.Sprintf("subscribers/%v", req.SubscriberID), req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// UnsubscribeSubscriberResponse is the data returned from an UnsubscribeSubscriber call.
type UnsubscribeSubscriberResponse struct {
	Subscriber `json:"subscriber"`
}

// UnsubscribeSubscriber will update a subscriber's information.
func (c *Client) UnsubscribeSubscriber(email string) (*UnsubscribeSubscriberResponse, error) {
	var req struct {
		Email string `json:"email"`
	}
	req.Email = email
	var ret UnsubscribeSubscriberResponse
	err := c.Do(http.MethodPut, "unsubscribe", req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
