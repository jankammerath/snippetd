package snippetd

type CodeExecution struct {
	Id         string `json:"id"`
	Container  string `json:"container"`
	Language   string `json:"language"`
	SourceCode string `json:"source_code"`
	TempPath   string `json:"-"`
}
