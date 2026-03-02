package types

type Text struct {
	Chunks []string `json:"chunks"`
}

type Args interface{}

type Tool struct {
	Name string `json:"name"`
	Args Args   `json:"args"`
}

type MockResponse struct {
	Agent string `json:"agent,omitempty"`
	Text  Text   `json:"text"`
	Tool  Tool   `json:"tool"`
}
