package demo

// AnythingResponse is the response type for Anything.
type AnythingResponse struct {
	Args    map[string]string      `json:"args"`
	Data    string                 `json:"data"`
	Files   map[string]interface{} `json:"files"`
	Form    map[string]interface{} `json:"form"`
	Headers map[string]string      `json:"headers"`
	JSON    map[string]interface{} `json:"json"`
	Method  string                 `json:"method"`
	Origin  string                 `json:"origin"`
	URL     string                 `json:"url"`
}
