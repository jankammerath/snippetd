package snippetd

type CodeExecution struct {
	Uuid           string `json:"uuid"`
	Language       string `json:"language"`
	StandardOutput string `json:"standard_output,omitempty"`
	StandardError  string `json:"standard_error,omitempty"`
	Completed      bool   `json:"completed,omitempty"`
	ExitCode       int    `json:"exit_code,omitempty"`
}
