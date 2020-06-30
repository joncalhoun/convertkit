package convertkit

import (
	"net/http"
)

// AccountResponse defines the data returned from an Account API call.
type AccountResponse struct {
	Name         string `json:"name"`
	PrimaryEmail string `json:"primary_email_address"`
}

// Account shows the account information for the provided secret.
func (c *Client) Account() (*AccountResponse, error) {
	var ret AccountResponse
	err := c.Do(http.MethodGet, "account", nil, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
