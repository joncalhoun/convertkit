package convertkit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Tag can be applied to subscribers to help filter and customize your mailing
// list actions. Eg you might tag a subscriber "beginner" and send them
// beginner-oriented emails, or you might tag them as interested in a paid
// course so they get information about future sales.
type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// TagsResponse is the response data from Tags.
type TagsResponse struct {
	Tags []Tag `json:"tags"`
}

// Tags lists the sequences from your account.
func (c *Client) Tags() (*TagsResponse, error) {
	var ret TagsResponse
	err := c.Do(http.MethodGet, "tags", nil, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// CreateTagsResponse is the data returned from a CreateTags call.
type CreateTagsResponse struct {
	Tags []Tag
}

// UnmarshalJSON implements json.Unmarshaler
func (ctr *CreateTagsResponse) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return fmt.Errorf("no bytes to unmarshal")
	}
	// See if we can guess based on the first character
	switch b[0] {
	case '{':
		return ctr.unmarshalSingle(b)
	case '[':
		return ctr.unmarshalMany(b)
	}
	// This shouldn't really happen as the standard library seems to strip
	// whitespace from the bytes being passed in, but just in case let's guess at
	// multiple tags and fall back to a single one if that doesn't work.
	err := ctr.unmarshalMany(b)
	if err != nil {
		return ctr.unmarshalSingle(b)
	}
	return nil
}

func (ctr *CreateTagsResponse) unmarshalSingle(b []byte) error {
	var t Tag
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	ctr.Tags = []Tag{t}
	return nil
}

func (ctr *CreateTagsResponse) unmarshalMany(b []byte) error {
	var tags []Tag
	err := json.Unmarshal(b, &tags)
	if err != nil {
		return err
	}
	ctr.Tags = tags
	return nil
}

// CreateTags will create tags using the provided values as their names.
func (c *Client) CreateTags(tags ...string) (*CreateTagsResponse, error) {
	type newTag struct {
		Name string `json:"name"`
	}
	var data struct {
		Tags []newTag `json:"tag"`
	}
	for _, tag := range tags {
		data.Tags = append(data.Tags, newTag{tag})
	}
	var ret CreateTagsResponse
	err := c.Do(http.MethodPost, fmt.Sprintf("tags"), data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// Delete isn't supported by the API despite the UI doing it via a path similar to this.
// // DeleteTag will create tags using the provided values as their names.
// func (c *Client) DeleteTag(id int, name string) (interface{}, error) {
// 	var data struct {
// 		Name string `json:"name"`
// 	}
// 	data.Name = name
// 	var ret interface{}
// 	err := c.Do(http.MethodDelete, fmt.Sprintf("tags/%v", id), data, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &ret, nil
// }
