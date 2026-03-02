package types

import "encoding/json"

type Text struct {
	Chunks []string `json:"chunks"`
}

type Args []interface{}

// UnmarshalJSON implements custom unmarshaling logic for Args to handle both array and single object formats
func (a *Args) UnmarshalJSON(data []byte) error {
	var arr []interface{}
	if err := json.Unmarshal(data, &arr); err == nil {
		*a = arr
		return nil
	}

	// If that fails, unmarshal as single object and wrap in array
	var single interface{}
	if err := json.Unmarshal(data, &single); err != nil {
		return err
	}
	*a = []interface{}{single}
	return nil
}

type Tool struct {
	Name string `json:"name"`
	Args Args   `json:"args"`
}

type MockResponse struct {
	Agent string `json:"agent,omitempty"`
	Text  Text   `json:"text"`
	Tool  Tool   `json:"tool"`
}
