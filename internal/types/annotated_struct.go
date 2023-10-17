package types

// AnnotatedStruct is a struct helper for adding descriptions to objects that are returned by the API
type AnnotatedStruct struct {
	Description string      `json:"description" yaml:"description"`
	Value       interface{} `json:"value" yaml:"value"`
}
