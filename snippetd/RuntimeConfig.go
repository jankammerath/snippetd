package snippetd

type RuntimeConfig struct {
	Container string   `json:"container"`
	FileName  string   `json:"file_name"`
	MimeTypes []string `json:"mime_types"`
	RunScript string   `json:"-"`
}
