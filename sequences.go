package convertkit

import (
	"net/http"
	"time"
)

// Sequence is a series of emails that a user might receive. It was previously
// called a course, which is why some the JSON is a little wonky.
type Sequence struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// SequencesResponse is returned after making a Sequence call.
//
// Sequences were previously called courses, hence the mismatched json naming.
type SequencesResponse struct {
	Sequences []Sequence `json:"courses"`
}

// Sequences lists the sequences from your account.
func (c *Client) Sequences() (*SequencesResponse, error) {
	var ret SequencesResponse
	err := c.Do(http.MethodGet, "sequences", nil, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
