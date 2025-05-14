package config

const BaseURL = "http://host.docker.internal"

type Statistics struct {
	ID          string
	Rank        string
	IsDuplicate bool
}

type Text struct {
	Value   string
	Country string
}
