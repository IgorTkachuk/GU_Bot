package github

type Response struct {
	Type     string `json:"type,omitempty"`
	Name     string `json:"name,omitempty"`
	Encoding string `json:"encoding,omitempty"`
	Content  string `json:"content,omitempty"`
}
