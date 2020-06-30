package convertkit

import (
	"net/http"
	"time"
)

// Form is an entry point for a user joining a mailing list. Typically this is
// an HTML form, but you can subscribe someone to a form via API as well. It is
// returned from several API endpoints.
type Form struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	CreatedAt        time.Time `json:"created_at"`
	Type             string    `json:"type"`
	URL              string    `json:"url"`
	EmbedJs          string    `json:"embed_js"`
	EmbedURL         string    `json:"embed_url"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	SignUpButtonText string    `json:"sign_up_button_text"`
	SuccessMessage   string    `json:"success_message"`
}

// FormsResponse defines the data returned from a Forms API call.
type FormsResponse struct {
	Forms []Form `json:"forms"`
}

// Forms lists the forms from your account.
func (c *Client) Forms() (*FormsResponse, error) {
	var ret FormsResponse
	err := c.Do(http.MethodGet, "forms", nil, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
