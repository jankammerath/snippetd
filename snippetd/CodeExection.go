package snippetd

type CodeExecution struct {
	Uuid           string `json:"uuid"`
	StandardOutput string `json:"standard_output,omitempty"`
	StandardError  string `json:"standard_error,omitempty"`
	ExitCode       int    `json:"exit_code,omitempty"`
}
