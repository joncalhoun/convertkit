package convertkit_test

import (
	"flag"
	"math/rand"
	"testing"
	"time"

	"github.com/joncalhoun/convertkit"
)

var (
	APISecret string
)

func init() {
	flag.StringVar(&APISecret, "secret", "", "Your API secret key. Only provide this if you have a ConvertKit developer account. If present, integration tests will be run using this key.")
}

// I DO NOT recommend using this test. It is mostly a sanity check I used for
// the API library.
func TestClient_Integration(t *testing.T) {
	if APISecret == "" {
		t.Skip("skipping integration tests - API secret flag missing")
	}
	c := convertkit.Client{
		Secret: APISecret,
	}
	seed := time.Now().UnixNano()
	t.Logf("seed is: %v", seed)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomString := func() string {
		const alphabet = "abcdefghijklmnopqrstuvwxyz"
		var ret []byte
		for i := 0; i < 10; i++ {
			ret = append(ret, alphabet[random.Intn(len(alphabet))])
		}
		return string(ret)
	}

	account, err := c.Account()
	if err != nil {
		t.Fatalf("Account() err = %v; want nil", err)
	}
	if account.PrimaryEmail == "" {
		t.Errorf("Account() PrimaryEmail is empty")
	}

	// A form must exist and you should probably turn off all incentive emails and confirming for this to work
	t.Run("subscribe via form and unsub", func(t *testing.T) {
		formsResponse, err := c.Forms()
		if err != nil {
			t.Fatalf("Forms() err = %v; want nil", err)
		}
		if len(formsResponse.Forms) == 0 {
			t.Skipf("form required for this test")
		}
		formID := formsResponse.Forms[0].ID
		email := randomString() + "@calhoun.io"
		subResp, err := c.SubscribeToForm(convertkit.SubscribeToFormRequest{
			FormID: formID,
			Email:  email,
		})
		if err != nil {
			t.Errorf("SubscribeToForm() err = %v; want nil", err)
		}
		t.Logf("Subscribe Response: %+v", subResp)
		formSubsResp, err := c.FormSubscriptions(convertkit.FormSubscriptionsRequest{
			FormID: formID,
		})
		if err != nil {
			t.Errorf("FormSubscriptions() err = %v; want nil", err)
		}
		found := false
		for _, sub := range formSubsResp.Subscriptions {
			if sub.Subscriber.Email == email {
				found = true
				break
			}
		}
		if !found {
			t.Logf("%+v", formSubsResp)
			t.Errorf("FormSubscriptions() didn't return subscriber with Email = %v", email)
		}
		_, err = c.UnsubscribeSubscriber(email)
		if err != nil {
			t.Errorf("UnsubscribeSubscriber() err = %v; want nil", err)
		}
	})

	t.Run("creating one tag", func(t *testing.T) {
		want := randomString()
		resp, err := c.CreateTags(want)
		if err != nil {
			t.Fatalf("CreateTags() err = %v; want nil", err)
		}
		if len(resp.Tags) != 1 {
			t.Fatalf("len(resp.Tags) = %v; want 1", len(resp.Tags))
		}
		if resp.Tags[0].Name != want {
			t.Fatalf("Tag.Name = %v; want %v", resp.Tags[0].Name, want)
		}
	})

	t.Run("creating multiple tags", func(t *testing.T) {
		tags := []string{randomString(), randomString(), randomString()}
		resp, err := c.CreateTags(tags...)
		if err != nil {
			t.Fatalf("CreateTags() err = %v; want nil", err)
		}
		if len(resp.Tags) != len(tags) {
			t.Fatalf("len(resp.Tags) = %v; want %v", len(resp.Tags), len(tags))
		}
		got := make(map[string]struct{})
		for _, tag := range resp.Tags {
			got[tag.Name] = struct{}{}
		}
		for _, want := range tags {
			if _, ok := got[want]; !ok {
				t.Errorf("Tags missing %v; got %v", want, got)
			}
		}
	})

}
