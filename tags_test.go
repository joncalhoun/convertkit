package convertkit_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/joncalhoun/convertkit"
)

func TestClient_Tags(t *testing.T) {
	c := client(t, "fake-secret-key")
	t.Run("multi", func(t *testing.T) {
		resp, err := c.Tags()
		if err != nil {
			t.Fatalf("Tags() err = %v; want %v", err, nil)
		}
		if len(resp.Tags) != 2 {
			t.Errorf("len(Tags) = %v; want 2", len(resp.Tags))
		}
	})
}

func TestClient_CreateTags(t *testing.T) {
	c := client(t, "fake-secret-key")
	t.Run("multi", func(t *testing.T) {
		resp, err := c.CreateTags("Example Tag", "Example Tag 2")
		if err != nil {
			t.Fatalf("CreateTags() err = %v; want %v", err, nil)
		}
		if resp.Tags[0].Name != "Example Tag" {
			t.Errorf("Tag[0].Name = %v; want %v", resp.Tags[0].Name, "Example Tag")
		}
	})

	t.Run("singular", func(t *testing.T) {
		c := clientWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
			testdataHandler(t, "POST_tags_singular")(w, r)
		})
		want := "Single Tag"
		resp, err := c.CreateTags(want)
		if err != nil {
			t.Fatalf("CreateTags() err = %v; want nil", err)
		}
		if len(resp.Tags) != 1 {
			t.Errorf("len(Tags) = %v; want 1", len(resp.Tags))
		}
		if resp.Tags[0].Name != want {
			t.Errorf("Tag.Name = %v; want %v", resp.Tags[0].Name, want)
		}
	})
}

func TestCreateTagsResponse_UnmarshalJSON(t *testing.T) {
	parseTime := func(str string) time.Time {
		ti, err := time.Parse(time.RFC3339, str)
		if err != nil {
			t.Fatalf("time.Parse() err = %v; want nil", err)
		}
		return ti
	}
	for name, tc := range map[string]struct {
		bytes []byte
		want  convertkit.CreateTagsResponse
	}{
		"single": {
			bytes: []byte(`{
        "id": 1,
        "name": "House Stark",
        "created_at": "2016-02-28T08:07:00Z"
      }`),
			want: convertkit.CreateTagsResponse{
				Tags: []convertkit.Tag{
					{
						ID:        1,
						Name:      "House Stark",
						CreatedAt: parseTime("2016-02-28T08:07:00Z"),
					},
				},
			},
		},
		"single leading spaces": {
			bytes: []byte(` {
        "id": 1,
        "name": "House Stark",
        "created_at": "2016-02-28T08:07:00Z"
      }`),
			want: convertkit.CreateTagsResponse{
				Tags: []convertkit.Tag{
					{
						ID:        1,
						Name:      "House Stark",
						CreatedAt: parseTime("2016-02-28T08:07:00Z"),
					},
				},
			},
		},
		"multi": {
			bytes: []byte(`[{
        "id": 1,
        "name": "House Stark",
        "created_at": "2016-02-28T08:07:00Z"
      },{
        "id": 2,
        "name": "House Lannister",
        "created_at": "2016-02-28T08:10:00Z"
      }]`),
			want: convertkit.CreateTagsResponse{
				Tags: []convertkit.Tag{
					{
						ID:        1,
						Name:      "House Stark",
						CreatedAt: parseTime("2016-02-28T08:07:00Z"),
					}, {
						ID:        2,
						Name:      "House Lannister",
						CreatedAt: parseTime("2016-02-28T08:10:00Z"),
					},
				},
			},
		},
		"multi leading spaces": {
			bytes: []byte(`  [  {
        "id": 1,
        "name": "House Stark",
        "created_at": "2016-02-28T08:07:00Z"
      },{
        "id": 2,
        "name": "House Lannister",
        "created_at": "2016-02-28T08:10:00Z"
      }]`),
			want: convertkit.CreateTagsResponse{
				Tags: []convertkit.Tag{
					{
						ID:        1,
						Name:      "House Stark",
						CreatedAt: parseTime("2016-02-28T08:07:00Z"),
					}, {
						ID:        2,
						Name:      "House Lannister",
						CreatedAt: parseTime("2016-02-28T08:10:00Z"),
					},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			var got convertkit.CreateTagsResponse
			err := json.Unmarshal(tc.bytes, &got)
			if err != nil {
				t.Fatalf("Unmarshal() err = %v; want nil", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got = %+v; want %+v", got, tc.want)
			}
		})
	}
}
