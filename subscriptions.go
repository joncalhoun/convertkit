package convertkit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Subscription is a shared object across a few endpoints. It represents an
// entity being subscribed to something. Typically this is a Subscriber being
// subscribed to a Sequence of Form.
type Subscription struct {
	ID               int         `json:"id"`
	State            string      `json:"state"`
	CreatedAt        time.Time   `json:"created_at"`
	Source           interface{} `json:"source"`
	Referrer         interface{} `json:"referrer"`
	SubscribableID   int         `json:"subscribable_id"`
	SubscribableType string      `json:"subscribable_type"`
	Subscriber       Subscriber  `json:"subscriber"`
}

// SubscribeToFormRequest is used when making SubscribeToForm calls.
type SubscribeToFormRequest struct {
	// Required
	FormID int    `json:"-"`
	Email  string `json:"email"`
	// Optional
	FirstName string            `json:"first_name,omitempty"`
	Fields    map[string]string `json:"fields,omitempty"`
	TagIDs    []int             `json:"tags,omitempty"`
}

// SubscribeToFormResponse is the response data from SubscribeToForm.
type SubscribeToFormResponse struct {
	Subscription Subscription `json:"subscription"`
}

// SubscribeToForm will subscribe an email address to a form.
func (c *Client) SubscribeToForm(req SubscribeToFormRequest) (*SubscribeToFormResponse, error) {
	var ret SubscribeToFormResponse
	err := c.Do(http.MethodPost, fmt.Sprintf("forms/%v/subscribe", req.FormID), req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// FormSubscriptionsRequest is used when making FormSubscriptions calls.
type FormSubscriptionsRequest struct {
	// Required
	FormID int `json:"-"`
	// Optional
	SortOrder       SortOrder       `json:"sort_order,omitempty"`
	SubscriberState SubscriberState `json:"subscriber_state,omitempty"`
}

// FormSubscriptionsResponse is the response data from FormSubscriptions.
type FormSubscriptionsResponse struct {
	TotalSubscriptions int            `json:"total_subscriptions"`
	Page               int            `json:"page"`
	TotalPages         int            `json:"total_pages"`
	Subscriptions      []Subscription `json:"subscriptions"`
}

// FormSubscriptions will subscribe an email address to a form.
func (c *Client) FormSubscriptions(req FormSubscriptionsRequest) (*FormSubscriptionsResponse, error) {
	var ret FormSubscriptionsResponse
	err := c.Do(http.MethodGet, fmt.Sprintf("forms/%v/subscriptions", req.FormID), req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// SubscribeToSequenceRequest is used when making SubscribeToSequence calls.
type SubscribeToSequenceRequest struct {
	// Required
	SequenceID int    `json:"-"`
	Email      string `json:"email"`
	// Optional
	FirstName string            `json:"first_name,omitempty"`
	Fields    map[string]string `json:"fields,omitempty"`
	TagIDs    []int             `json:"tags,omitempty"`
}

// SubscribeToSequenceResponse is the response data from SubscribeToSequence.
type SubscribeToSequenceResponse struct {
	Subscription Subscription `json:"subscription"`
}

// SubscribeToSequence will subscribe an email address to a form.
func (c *Client) SubscribeToSequence(req SubscribeToSequenceRequest) (*SubscribeToSequenceResponse, error) {
	var ret SubscribeToSequenceResponse
	err := c.Do(http.MethodPost, fmt.Sprintf("sequences/%v/subscribe", req.SequenceID), req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// SequenceSubscriptionsRequest is used when making SequenceSubscriptions calls.
type SequenceSubscriptionsRequest struct {
	// Required
	SequenceID int `json:"-"`
	// Optional
	SortOrder       SortOrder       `json:"sort_order,omitempty"`
	SubscriberState SubscriberState `json:"subscriber_state,omitempty"`
}

// SequenceSubscriptionsResponse is the response data from SequenceSubscriptions.
type SequenceSubscriptionsResponse struct {
	TotalSubscriptions int            `json:"total_subscriptions"`
	Page               int            `json:"page"`
	TotalPages         int            `json:"total_pages"`
	Subscriptions      []Subscription `json:"subscriptions"`
}

// SequenceSubscriptions will subscribe an email address to a form.
func (c *Client) SequenceSubscriptions(req SequenceSubscriptionsRequest) (*SequenceSubscriptionsResponse, error) {
	var ret SequenceSubscriptionsResponse
	err := c.Do(http.MethodGet, fmt.Sprintf("sequences/%v/subscriptions", req.SequenceID), req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// TagSubscriptionsRequest is used when making TagSubscriptions calls.
type TagSubscriptionsRequest struct {
	// Required
	TagID int `json:"-"`
	// Optional
	SortOrder       SortOrder       `json:"sort_order,omitempty"`
	SubscriberState SubscriberState `json:"subscriber_state,omitempty"`
}

// TagSubscriptionsResponse is the response data from TagSubscriptions.
type TagSubscriptionsResponse struct {
	TotalSubscriptions int            `json:"total_subscriptions"`
	Page               int            `json:"page"`
	TotalPages         int            `json:"total_pages"`
	Subscriptions      []Subscription `json:"subscriptions"`
}

// TagSubscriptions will subscribe an email address to a form.
func (c *Client) TagSubscriptions(req TagSubscriptionsRequest) (*TagSubscriptionsResponse, error) {
	var ret TagSubscriptionsResponse
	err := c.Do(http.MethodGet, fmt.Sprintf("tags/%v/subscriptions", req.TagID), req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// TagSubscriberRequest is used when making TagSubscriber calls.
type TagSubscriberRequest struct {
	// Required
	TagID int    `json:"-"`
	Email string `json:"email"`
	// Optional
	FirstName string            `json:"first_name,omitempty"`
	Fields    map[string]string `json:"fields,omitempty"`
	// Additional TagIDs you wish to apply to the user.
	TagIDs []int `json:"tags,omitempty"`
}

// TagSubscriberResponse is the response data from TagSubscriber.
type TagSubscriberResponse struct {
	Subscription Subscription `json:"subscription"`
}

// TagSubscriber will subscribe an email address to a form.
func (c *Client) TagSubscriber(req TagSubscriberRequest) (*TagSubscriberResponse, error) {
	var ret TagSubscriberResponse
	err := c.Do(http.MethodPost, fmt.Sprintf("tags/%v/subscribe", req.TagID), req, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// UntagSubscriberRequest is used when making UntagSubscriber calls.
type UntagSubscriberRequest struct {
	// Required
	SubscriberID int `json:"-"`
	TagID        int `json:"-"`
}

// UntagSubscriberResponse is the response data from UntagSubscriber.
type UntagSubscriberResponse struct {
	Tag Tag
}

func (usr *UntagSubscriberResponse) UnmarshalJSON(b []byte) error {
	var t Tag
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	usr.Tag = t
	return nil
}

// UntagSubscriber will subscribe an email address to a form.
func (c *Client) UntagSubscriber(req UntagSubscriberRequest) (*UntagSubscriberResponse, error) {
	var ret UntagSubscriberResponse
	var tag Tag
	err := c.Do(http.MethodDelete, fmt.Sprintf("subscribers/%v/tags/%v", req.SubscriberID, req.TagID), req, &ret)
	if err != nil {
		return nil, err
	}
	ret.Tag = tag
	return &ret, nil
}
