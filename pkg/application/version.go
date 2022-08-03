package application

const Name = ""
const Version = ""
const Build = ""
const Commit = ""

type InfoFields struct {
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	Build   string `json:"build,omitempty" yaml:"build,omitempty"`
	Commit  string `json:"commit,omitempty" yaml:"commit,omitempty"`
}

func Info() InfoFields {
	return InfoFields{
		Name:    Name,
		Version: Version,
		Build:   Build,
		Commit:  Commit,
	}
}
