package resources

type UpdateEnvoyConfRequest struct {
	RateLimit string `json:"rate_limit"`
	Weight    string `json:"weight"`
}

type UpdateEnvoyConfResponse struct {
}
