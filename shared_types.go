package convertkit

import (
	"encoding/json"
	"fmt"
	"time"
)

// SortOrder is used to define sorting order of some API requests.
type SortOrder string

// SortOrders supported by the API
const (
	SortOldToNew SortOrder = "asc"
	SortNewToOld           = "desc"
)

// SubscriberState is used to filter subscribers in some API calls.
type SubscriberState string

// SubscriberStates supported by the API
const (
	SubscriberStateActive    SubscriberState = "active"
	SubscriberStateCancelled                 = "cancelled"
)

// NewDate is a helper for constructing dates.
func NewDate(year, month, day int) Date {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	d := Date(t)
	return d
}

// Date is used to represent dates that you might query the Convert Kit API
// with. It is basically a time.Time to make it easier to pull in times from
// other code, but please note that all Date objects are converted to
// "yyyy-mm-dd" when interacting with the API, so any finer time increments will
// be lost.
type Date time.Time

// MarshalJSON converts a Date into yyyy-mm-dd format
func (d *Date) MarshalJSON() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	t := time.Time(*d)
	year, month, day := t.Date()
	return json.Marshal(fmt.Sprintf("%4d-%02d-%02d", year, month, day))
}

// TODO: Test this, assuming we even need it. Not so sure we do at all.
// // UnmarhsalJSON parses yyyy-mm-dd format into a Date
// func (d *Date) UnmarhsalJSON(b []byte) error {
// 	if len(b) == 0 {
// 		return nil
// 	}
// 	var year, month, day int
// 	fmt.Fscanf(bytes.NewReader(b), "\"%4d-%2d-%2d\"", &year, &month, &day)
// 	*d = Date(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC))
// 	return nil
// }
