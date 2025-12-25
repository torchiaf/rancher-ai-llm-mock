package types

type Text struct {
	Chunks []string `json:"chunks"`
}

type Tool struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type MockResponse struct {
	Text Text `json:"text"`
	Tool Tool `json:"tool"`
}
