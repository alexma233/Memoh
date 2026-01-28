package settings

const (
	DefaultMaxContextLoadTime = 24 * 60
	DefaultLanguage           = "Same as user input"
)

type Settings struct {
	MaxContextLoadTime int    `json:"max_context_load_time"`
	Language           string `json:"language"`
}

type UpsertRequest struct {
	MaxContextLoadTime *int   `json:"max_context_load_time,omitempty"`
	Language           string `json:"language,omitempty"`
}

