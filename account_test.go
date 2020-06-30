package convertkit_test

import (
	"testing"
)

func TestClient_Account(t *testing.T) {
	c := client(t, "fake-secret-key")
	resp, err := c.Account()
	if err != nil {
		t.Fatalf("Account() err = %v; want %v", err, nil)
	}
	if resp.Name != "Acme Corp." {
		t.Errorf("Name = %v; want %v", resp.Name, "Acme Corp.")
	}
	if resp.PrimaryEmail != "you@example.com" {
		t.Errorf("PrimaryEmail = %v; want %v", resp.PrimaryEmail, "you@example.com")
	}
}
