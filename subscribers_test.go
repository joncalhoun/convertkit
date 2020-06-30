package convertkit_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/joncalhoun/convertkit"
)

// This test is incredibly exhaustive because the requests here use a few custom
// types. Other endpoints likely don't need to have such exhaustive tests.
func TestClient_Subscribers(t *testing.T) {
	c := client(t, "fake-secret-key")
	t.Run("no options", func(t *testing.T) {
		resp, err := c.Subscribers(convertkit.SubscribersRequest{})
		if err != nil {
			t.Fatalf("Subscribers() err = %v; want %v", err, nil)
		}
		if len(resp.Subscribers) != 2 {
			t.Errorf("len(.Subscribers) = %d; want 2", len(resp.Subscribers))
		}
		if resp.TotalSubscribers != 2 {
			t.Errorf("TotalSubscribers = %d; want 2", resp.TotalSubscribers)
		}
		if resp.Page != 1 {
			t.Errorf("Page = %d; want 1", resp.Page)
		}
		for _, f := range resp.Subscribers {
			if f.ID <= 0 {
				t.Errorf("ID = %d; want > 0", f.ID)
			}
			if f.FirstName == "" {
				t.Errorf("FirstName is empty")
			}
		}
	})

	t.Run("page option", func(t *testing.T) {
		c := clientWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			got := r.FormValue("page")
			if got != "2" {
				t.Errorf("Server recieved page %v; want 2", got)
			}
			testdataHandler(t, "GET_subscribers_page_2")(w, r)
		})
		resp, err := c.Subscribers(convertkit.SubscribersRequest{
			Page: 2,
		})
		// Most of these checks don't matter much. The point of this test is
		// verifying that the options made it to the server.
		if err != nil {
			t.Fatalf("Subscribers() err = %v; want %v", err, nil)
		}
		if len(resp.Subscribers) != 1 {
			t.Errorf("len(.Subscribers) = %d; want 1", len(resp.Subscribers))
		}
		if resp.Page != 2 {
			t.Errorf("Page = %d; want 2", resp.Page)
		}
	})

	t.Run("options are sent to server", func(t *testing.T) {
		ckDate := func(year, month, day int) *convertkit.Date {
			d := convertkit.NewDate(year, month, day)
			return &d
		}
		for name, tc := range map[string]struct {
			want map[string]string
			req  convertkit.SubscribersRequest
		}{
			"From": {
				want: map[string]string{
					"from": "1999-10-21",
				},
				req: convertkit.SubscribersRequest{
					From: ckDate(1999, 10, 21),
				},
			},
			"To": {
				want: map[string]string{
					"to": "2001-09-18",
				},
				req: convertkit.SubscribersRequest{
					To: ckDate(2001, 9, 18),
				},
			},
			"UpdatedFrom": {
				want: map[string]string{
					"updated_from": "2020-02-01",
				},
				req: convertkit.SubscribersRequest{
					UpdatedFrom: ckDate(2020, 2, 1),
				},
			},
			"UpdatedTo": {
				want: map[string]string{
					"updated_to": "2019-05-05",
				},
				req: convertkit.SubscribersRequest{
					UpdatedTo: ckDate(2019, 5, 5),
				},
			},
			"SortOrder": {
				want: map[string]string{
					"sort_order": "asc",
				},
				req: convertkit.SubscribersRequest{
					SortOrder: convertkit.SortOldToNew,
				},
			},
			"SortField": {
				want: map[string]string{
					"sort_field": "cancelled_at",
				},
				req: convertkit.SubscribersRequest{
					SortField: "cancelled_at",
				},
			},
			"Email": {
				want: map[string]string{
					"email_address": "test@user.com",
				},
				req: convertkit.SubscribersRequest{
					Email: "test@user.com",
				},
			},
		} {
			t.Run(name, func(t *testing.T) {
				c := clientWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
					r.ParseForm()
					for k, want := range tc.want {
						got := r.FormValue(k)
						if got != want {
							t.Errorf("%v = %v; want %v", k, got, want)
						}
					}
					testdataHandler(t, "GET_subscribers")(w, r)
				})
				_, err := c.Subscribers(tc.req)
				if err != nil {
					t.Fatalf("Subscribers() err = %v; want nil", err)
				}
			})
		}
	})
}

func TestClient_UpdateSubscriber(t *testing.T) {
	c := client(t, "fake-secret-key")
	t.Run("response data", func(t *testing.T) {
		resp, err := c.UpdateSubscriber(convertkit.UpdateSubscriberRequest{
			SubscriberID: 123,
			FirstName:    "New First Name",
		})
		if err != nil {
			t.Fatalf("UpdateSubscriber() err = %v; want %v", err, nil)
		}
		if resp.Subscriber.Email != "jonsnow@example.com" {
			t.Errorf("Email = %v; want %v", resp.Subscriber.Email, "jonsnow@example.com")
		}
	})

	t.Run("options are sent to server", func(t *testing.T) {
		for name, tc := range map[string]struct {
			want map[string]interface{}
			req  convertkit.UpdateSubscriberRequest
		}{
			"FirstName": {
				want: map[string]interface{}{
					"first_name": "Bob",
				},
				req: convertkit.UpdateSubscriberRequest{
					FirstName: "Bob",
				},
			},
			"Email": {
				want: map[string]interface{}{
					"email_address": "new@email.com",
				},
				req: convertkit.UpdateSubscriberRequest{
					Email: "new@email.com",
				},
			},
			"Fields": {
				want: map[string]interface{}{
					"fields": map[string]interface{}{
						"last_name": "Parker",
					},
				},
				req: convertkit.UpdateSubscriberRequest{
					Fields: map[string]string{
						"last_name": "Parker",
					},
				},
			},
		} {
			t.Run(name, func(t *testing.T) {
				c := clientWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
					var gotBody map[string]interface{}
					err := json.NewDecoder(r.Body).Decode(&gotBody)
					if err != nil {
						t.Fatalf("decode: %v", err)
					}
					for k, want := range tc.want {
						got := gotBody[k]
						if !reflect.DeepEqual(got, want) {
							t.Errorf("%v = %v; want %v", k, got, want)
						}
					}
					testdataHandler(t, "PUT_subscribers_123")(w, r)
				})
				tc.req.SubscriberID = 123
				_, err := c.UpdateSubscriber(tc.req)
				if err != nil {
					t.Fatalf("UpdateSubscriber() err = %v; want nil", err)
				}
			})
		}
	})
}

func TestClient_UnsubscribeSubscriber(t *testing.T) {
	c := client(t, "fake-secret-key")
	t.Run("response data", func(t *testing.T) {
		email := "jonsnow@example.com"
		resp, err := c.UnsubscribeSubscriber(email)
		if err != nil {
			t.Fatalf("UnsubscribeSubscriber() err = %v; want %v", err, nil)
		}
		if resp.Subscriber.Email != email {
			t.Errorf("Email = %v; want %v", resp.Subscriber.Email, email)
		}
	})

	t.Run("options are sent to server", func(t *testing.T) {
		want := "jonsnow@example.com"
		c := clientWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
			var gotBody map[string]string
			err := json.NewDecoder(r.Body).Decode(&gotBody)
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if gotBody["email"] != want {
				t.Logf("body: %+v", gotBody)
				t.Errorf("email = %v; want %v", gotBody["email"], want)
			}
			testdataHandler(t, "PUT_unsubscribe")(w, r)
		})
		_, err := c.UnsubscribeSubscriber(want)
		if err != nil {
			t.Fatalf("UnsubscribeSubscriber() err = %v; want nil", err)
		}
	})
}
