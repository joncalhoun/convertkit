# ConvertKit API Library

`convertkit` is an API library for the [Convert Kit](https://convertkit.com/) API.

## Usage

To start using the library, you first need to setup a client:

```go
client := convertkit.Client{
  Secret: "you-convert-kit-secret",
}
```

Accessing the API via public API key is not supported. The secret key can access all of these endpoints plus additional ones, and given that this won't be public like a JavaScript client, it didn't make sense in the first version.

To make API calls, you simply need to call methods on the client. The general format used is that any GET request will be handled by a method named after the resources you are getting. Eg `client.Forms` or `client.Subscribers`.

For PUT/POST/DELETE requests, the function will start with a verb describing what is happening. A couple examples of this are `client.SubscribeToForm` and `client.UnsubscribeSubscriber`.

Below is a non-exhaustive list of examples:

```go
// List all forms
GET /v3/forms => client.Forms(...)
// List subscribers
GET /v3/subscribers => client.Subscribers(...)
// Get single subscriber
GET /v3/subscribers/:id => client.Subscriber(...)

// Subscribe an email address to a form
POST /v3/forms/:id/subscribe => client.SubscribeToForm(...)
// Subscribe an email address to a sequence
POST /v3/sequences/:id/subscribe => client.SubscribeToSequence(...)
// Update a subscriber
PUT /v3/subscribers/:id => client.UpdateSubscriber(...)
// Unsubscribe a subscriber
PUT /v3/unsubscribe => client.UnsubscribeSubscriber(...)
// Delete a webhook
DELETE /v3/automations/hooks/:id => client.DeleteWebhook(...)
```

`convertkit` does not support every API endpoint offered by Convert Kit, nor does it support official integration parameters & keys, because I didn't need them. Adding new endpoints is VERY easy, so if you need to add a specific endpoint feel free to check out the existing code and submit a PR, or create an issue detailing which endpoint you need.

_If you do submit a PR, please try to follow the same style used in the rest of the library. If this isn't followed, I will not accept the PR without some modifications. This is done in an attempt to keep this library more maintainable._

## What API endpoints are supported?

Below is a table of all the API endpoints along with which are and are not supported.


| CLIENT METHOD           | SUPPORT? | HTTP METHOD | PATH                                 |
|-------------------------|----------|-------------|--------------------------------------|
| `Account`               | Y        | GET         | /v3/account                          |
| `Forms`                 | Y        | GET         | /v3/forms                            |
| `SubscribeToForm`       | Y        | POST        | /v3/forms/:id/subscribe              |
| `FormSubscriptions`     | Y        | GET         | /v3/forms/:id/subscriptions          |
| `Sequences`             | Y        | GET         | /v3/sequences                        |
| `SubscribeToSequence`   | Y        | POST        | /v3/sequences/:id/subscribe          |
| `SequenceSubscriptions` | Y        | GET         | /v3/sequences/:id/subscriptions      |
| `Tags`                  | Y        | GET         | /v3/tags                             |
| `CreateTags`             | Y        | POST        | /v3/tags                             |
| `TagSubscriber`         | Y        | POST        | /v3/tags/:id/subscribe               |
| `UntagSubscriber`       | Y        | POST        | /v3/tags/:id/unsubscribe             |
| `UntagSubscriberByID`   | N        | DELETE      | /v3/subscribers/:sub_id/tags/:tag_id |
| `TagSubscriptions`      | N        | GET         | /v3/tags/:id/subscriptions           |
| `Subscribers`           | Y        | GET         | /v3/subscribers                      |
| `Subscriber`            | N        | GET         | /v3/subscribers/:id                  |
| `UpdateSubscriber`      | Y        | PUT         | /v3/subscribers/:id                  |
| `UnsubscribeSubscriber` | Y        | PUT         | /v3/unsubscribe                      |
| `SubscriberTags`        | N        | GET         | /v3/subscribers/:id/tags             |
| `Broadcasts`            | N        | GET         | /v3/broadcasts                       |
| `BroadcastStats`        | N        | GET         | /v3/broadcasts/:id/stats             |
| `CreateWebhook`         | N        | POST        | /v3/automations/hooks                |
| `DeleteWebhook`         | N        | DELETE      | /v3/automations/hooks/:rule_id       |
| `Fields`                | N        | GET         | /v3/custom_fields                    |
| `CreateField`           | N        | POST        | /v3/custom_fields                    |
| `UpdateField`           | N        | PUT         | /v3/custom_fields/:id                |
| `DeleteField`           | N        | DELETE      | /v3/custom_fields/:id                |
| `Purchases`             | N        | GET         | /v3/purchases                        |
| `Purchase`              | N        | GET         | /v3/purchases/:id                    |
| `CreatePurchase`        | N        | POST        | /v3/purchases                        |


**NOTE: client methods are listed for every endpoint, even if they aren't coded yet. This is done to help assist in planning. Please refer to the supported column before trying to use a method.**

An API endpoint is not considered to be implemented until the client method has been coded AND the method has at least one test.

## Adding new endpoints with `client.Do`

Almost all of the logic for the API library is handled in the `client.Do` method. This handles:

1. Adding the API Secret wherever it is needed.
2. Encoding the params.
3. Performing the HTTP request for the API call with data from (1) and (2).
4. Decoding the response body to the provided response variable.
5. Handling errors from 400+ status codes and parsing the body into an `ErrorResponse` error.

(1) even occurs on API endpoints where the API Secret is meant to be part of the HTTP Request Body. This is done by encoding the entire body into JSON, decoding it into a `map[string]interface{}`, adding the `"api_secret"` key and value, and then finally encoding the map back into JSON as either an io.Reader for PUT and POST requests, or as URL query params for GET and DELETE requests.

Because `client.Do` handles all of this logic, adding a new API endpoint (or manually calling an unsupported endpoint) is almost always as simple as:

1. Adding a new request type (if needed).
2. Adding a new response type.
3. Passing instances of (1) and (2) along with an HTTP method and path into `client.Do`

You should also write a test, but again this is pretty simple because the current testing tools use `testdata` and JSON files to do most of the heavy lifting.

TODO: Add more detailed instructions on adding tests.

## Roadmap

At some point I want to improve the fields that are not the best type. Eg some fields are of type `interface{}` which isn't very telling, but I need to use the API a bit more to verify what they should be and this first pass is getting me to a point where I can do that.
