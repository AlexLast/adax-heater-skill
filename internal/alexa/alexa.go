package alexa

import (
	"context"

	"github.com/alexlast/adax-heater-skill/internal/adax"
)

// TODO: move this to something more dynamic
const (
	version = "1.0"
)

// Context defines any context required
// functions in this package
type Context struct {
	Adax *adax.Client
}

// EventPayload defines the stucture for payloads
// sent from Alexa to our skill, following the schema defined by Amazon
type EventPayload struct {
	Version string         `json:"version"`
	Context *LambdaContext `json:"context"`
	Session *Session       `json:"session"`
	Request *Request       `json:"request"`
}

// ResponsePayload defines the payload returned by
// the lambda handler
type ResponsePayload struct {
	Version           string                 `json:"version"`
	Response          *Response              `json:"response"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
}

// Response defines the response inside the response payload
type Response struct {
	ShouldSessionEnd bool          `json:"shouldEndSession"`
	Reprompt         *Reprompt     `json:"reprompt,omitempty"`
	OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
}

// OutputSpeech defines the voice response returned to the user
type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	SSML string `json:"ssml,omitempty"`
}

// Reprompt defines whether Alexa should ask the user for more information
type Reprompt struct {
	OutputSpeech *OutputSpeech `json:"outputSpeech,omitempty"`
}

// LambdaContext defines the context inside the Alexa payload
type LambdaContext struct{}

// Session defines the session inside the Alexa payload
type Session struct{}

// Request defines the request inside the Alexa payload
type Request struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
}

// Handler is the function that handles requests from Alexa
// it will route request based on type to the correct function
func (c *Context) Handler(ctx context.Context, event *EventPayload) (*ResponsePayload, error) {
	switch event.Request.Type {
	case "LaunchRequest":
		response := &ResponsePayload{
			Version: version,
			Response: &Response{
				ShouldSessionEnd: false,
				OutputSpeech: &OutputSpeech{
					Type: "SSML",
					SSML: launchRequest,
				},
				Reprompt: &Reprompt{
					OutputSpeech: &OutputSpeech{
						Type: "SSML",
						SSML: launchRequestReprompt,
					},
				},
			},
		}

		return response, nil
	}

	return &ResponsePayload{}, nil
}
