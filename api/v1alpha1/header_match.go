package v1alpha1

type StringMatch struct {
	exact string `json:"items"`
}

type HeaderMatch struct {
	name         string      `json:"name"`
	string_match StringMatch `json:"string_match"`
	invert_match bool        `json:"invert_match,omitempty"`
}

type HttpResponseHeadersMatch struct {
	headers []HeaderMatch `json:"headers"`
}
