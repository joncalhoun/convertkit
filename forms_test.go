package convertkit_test

import (
	"testing"
)

func TestClient_Forms(t *testing.T) {
	c := client(t, "fake-secret-key")
	resp, err := c.Forms()
	if err != nil {
		t.Fatalf("Forms() err = %v; want %v", err, nil)
	}
	if len(resp.Forms) != 2 {
		t.Errorf("len(.Forms) = %d; want 2", len(resp.Forms))
	}
	for _, f := range resp.Forms {
		if f.ID <= 0 {
			t.Errorf("ID = %d; want > 0", f.ID)
		}
		if f.Name == "" {
			t.Errorf("Name is empty")
		}
	}
}
