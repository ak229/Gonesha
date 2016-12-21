package models

type Result struct {

	Name string `json:"name"`
	Status bool `json:"status"`
	Messages []string `json:"messages"`

}
