package types

type AnnotatedStruct struct {
	Description string      `json:"description" yaml:"description"`
	Value       interface{} `json:"value" yaml:"value"`
}
