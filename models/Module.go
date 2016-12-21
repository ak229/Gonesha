package models

type Module struct {

	Name string `json:"name"`
	FileInfo []FileMeta `json:"files"`
	TestsToRun []string `json:"tests_to_run"`
	EvaluationFile string `json:"evaluation_file"`
}
