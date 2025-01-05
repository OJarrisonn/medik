package config

import "gopkg.in/yaml.v3"

type Medik struct {
	Exams     []Exam              `yaml:"exams,omitempty"`
	Protocols map[string]Protocol `yaml:"protocols,omitempty"`
}

type Protocol struct {
	Exams []Exam `yaml:"exams,omitempty"`
}

type Exam struct {
	Type     string      `yaml:"exam"`
	Level    string      `yaml:"level,omitempty"`
	Vars     []string    `yaml:"vars,omitempty"`
	Paths    []string    `yaml:"paths,omitempty"`
	Options  []string    `yaml:"options,omitempty"`
	Regex    string      `yaml:"regex,omitempty"`
	Min      interface{} `yaml:"min,omitempty"`
	Max      interface{} `yaml:"max,omitempty"`
	Protocol string      `yaml:"protocol,omitempty"`
	Exists   bool        `yaml:"exists,omitempty"`
}

// Given the contents of a Medik configuration file, parse it and return a config.Medik object
func Parse(content string) (*Medik, error) {
	var m Medik
	err := yaml.Unmarshal([]byte(content), &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
