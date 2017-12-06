package models

type TimestampVersion struct {
	Version string `json:"version"`
}

type Metadata map[string]string

type InRequest struct {
	Source  Source  `json:"source"`
	Version TimestampVersion `json:"version"`
}

type InResponse struct {
	Version  TimestampVersion  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version  TimestampVersion  `json:"version"`
}

type CheckResponse []TimestampVersion

type Source struct {}

