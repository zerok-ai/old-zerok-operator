package v1alpha1

type StringMatch struct {
	Exact string `json:"exact"`
}

type HeaderMatch struct {
	Name         string      `json:"name"`
	String_match StringMatch `json:"string_match"`
	Invert_match bool        `json:"invert_match,omitempty"`
}

type HttpResponseHeadersMatch struct {
	Headers []HeaderMatch `json:"headers"`
}
