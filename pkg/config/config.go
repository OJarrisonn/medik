package config

import "gopkg.in/yaml.v3"

type Medik struct {
	Vitals []Exam			`yaml:"vitals,omitempty"`
	Checks []Exam			`yaml:"checks,omitempty"`
	Protocols []Protocol	`yaml:"protocols,omitempty"`
}

type Protocol struct {
	Name string			`yaml:"name"`
	Vitals []Exam		`yaml:"vitals,omitempty"`
	Checks []Exam 		`yaml:"checks,omitempty"`
}

type Exam struct {
	Type string			`yaml:"type"`
	Vars []string		`yaml:"vars,omitempty"`
	Options []string	`yaml:"options,omitempty"`
	Regex string		`yaml:"regex,omitempty"`
	Min interface{}		`yaml:"min,omitempty"`
	Max interface{}		`yaml:"max,omitempty"`
	Protocol string		`yaml:"protocol,omitempty"`
}

// Given the contents of a Medik configuration file, parse it and return a config.Medik object
func Parse(file string) (*Medik, error) {
	var m Medik
	err := yaml.Unmarshal([]byte(file), &m)
	
	if err != nil {
		return nil, err
	}
	
	return &m, nil
}