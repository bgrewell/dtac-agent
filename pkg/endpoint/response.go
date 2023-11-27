package endpoint

// Response represents the data structure for the response returned from an endpoint.
// It contains all the data that an endpoint would return in response to a request.
type Response struct {
	// Metadata holds additional data about the response, similar to the request metadata.
	// This can include things like execution time, server information, etc.
	Metadata map[string]string `json:"metadata,omitempty"`

	// Headers represent the HTTP-style headers that might accompany the response.
	// This can be used to pass extra information back to the client.
	Headers map[string][]string `json:"headers,omitempty"`

	// Parameters may include any additional data that the endpoint wants to return,
	// which isn't part of the main response body, similar to headers.
	Parameters map[string][]string `json:"parameters,omitempty"`

	// Value is the primary data of the response. Like the request body, this could be of any format
	// (binary, JSON, XML, etc.) and is intended to be interpreted by the client as per its requirements.
	Value []byte `json:"value,omitempty"`
}
