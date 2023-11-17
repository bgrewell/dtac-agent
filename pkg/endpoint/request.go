package endpoint

// EndpointRequest represents the data structure for a request made to an endpoint in the framework.
// It encapsulates all the necessary information that an endpoint might need to process a request.
type EndpointRequest struct {
	// Metadata holds additional data that might be relevant to the processing of the request.
	// This can include things like authentication tokens, trace IDs, etc.
	Metadata map[string]string `json:"metadata,omitempty"`

	// Headers represent the HTTP-style headers that might accompany the request.
	// This is useful for passing extra context or information to the endpoint.
	Headers map[string][]string `json:"headers,omitempty"`

	// Parameters are the query or path parameters that may be used by the endpoint to process the request.
	// These are typically key-value pairs.
	Parameters map[string][]string `json:"parameters,omitempty"`

	// Body holds the raw data of the request. This could be in any format (binary, JSON, XML, etc.),
	// and is intended to be interpreted by the endpoint as per its requirements.
	Body []byte `json:"body,omitempty"`
}
