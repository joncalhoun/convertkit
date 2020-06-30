package convertkit_test

import "testing"

func TestClient_Sequences(t *testing.T) {
	c := client(t, "fake-secret-key")
	resp, err := c.Sequences()
	if err != nil {
		t.Fatalf("Sequences() err = %v; want %v", err, nil)
	}
	if len(resp.Sequences) != 2 {
		t.Errorf("len(.Sequences) = %d; want 2", len(resp.Sequences))
	}
	for _, f := range resp.Sequences {
		if f.ID <= 0 {
			t.Errorf("ID = %d; want > 0", f.ID)
		}
		if f.Name == "" {
			t.Errorf("Name is empty")
		}
	}
}
