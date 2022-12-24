package main

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
)

//go:embed module.nix.template
var moduleTemplate string

type Schema struct {
	Includes []string  `xml:"include,omitempty" json:"includes,omitempty"`
	Files    []*File   `xml:"kcfgfile,omitempty" json:"files,omitempty"`
	Signals  []*Signal `xml:"signal,omitempty" json:"signals,omitempty"`
	Groups   []*Group  `xml:"group,omitempty" json:"groups,omitempty"`
}

type File struct {
	Name string `xml:"name,attr,omitempty" json:"name,omitempty"`
}

type Signal struct {
	Name      string      `xml:"name,attr,omitempty" json:"name,omitempty"`
	Arguments []*Argument `xml:"argument,omitempty" json:"arguments,omitempty"`
}

type Argument struct {
	Name string `xml:",innerxml" json:"name,omitempty"`
	Type string `xml:"type,attr,omitempty" json:"type,omitempty"`
}

type Group struct {
	Name    string   `xml:"name,attr,omitempty" json:"name,omitempty"`
	Entries []*Entry `xml:"entry,omitempty" json:"entries,omitempty"`
}

type Entry struct {
	Name    string    `xml:"name,attr,omitempty" json:"name,omitempty"`
	Key     string    `xml:"key,attr,omitempty" json:"key,omitempty"`
	Type    Type      `xml:"type,attr,omitempty" json:"type,omitempty"`
	Choices []*Choice `xml:"choices>choice,omitempty" json:"choices,omitempty"`
	Label   string    `xml:"label,omitempty" json:"label,omitempty"`
	Default *Default  `xml:"default,omitempty" json:"default,omitempty"`
	Emit    *Emit     `xml:"emit,omitempty" json:"emit,omitempty"`
}

func (e *Entry) ActualKey() string {
	if e.Key != "" {
		return e.Key
	}

	return e.Name
}

func (e *Entry) DefaultValue() (interface{}, error) {
	if e.Default == nil {
		return nil, nil
	}

	if e.Default.Code {
		Printfln("Note: default values with `code` flag set are not supported (%s)", e.Name)
		return nil, nil
	}

	switch t := e.Type; t {
	case TypeBool:
		return strconv.ParseBool(e.Default.Value)
	case TypeDouble:
		return strconv.ParseFloat(e.Default.Value, 64)
	case TypeInt:
		return strconv.ParseInt(e.Default.Value, 10, 64)
	case TypeEnum:
		if len(e.Choices) != 0 {
			return e.Choices[0].Name, nil
		}

		return e.Default.Value, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", t)
	}
}

type Type string

const (
	TypeBool   = "bool"
	TypeInt    = "int"
	TypeDouble = "double"
	TypeEnum   = "enum"
)

func (t *Type) UnmarshalText(text []byte) error {
	*t = Type(strings.ToLower(string(text)))
	return nil
}

type Choice struct {
	Name string `xml:"name,attr,omitempty" json:"name,omitempty"`
}

type Default struct {
	Value string `xml:",innerxml" json:"value,omitempty"`
	Code  bool   `xml:"code,attr,omitempty" json:"code,omitempty"`
}

type Emit struct {
	Signal string `xml:"signal,attr,omitempty" json:"signal,omitempty"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "kcfg2home: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		return err
	}

	var schema Schema
	if err := xml.Unmarshal(data, &schema); err != nil {
		return err
	}

	for _, group := range schema.Groups {
		for _, entry := range group.Entries {
			defaultValue, err := entry.DefaultValue()
			if err != nil {
				return err
			}

			Printfln("%s: %s", entry.ActualKey(), entry.Label)
			Printfln("  Type: %s", entry.Type)
			if len(entry.Choices) > 0 {
				fmt.Println("  Choices:")
				for _, choice := range entry.Choices {
					Printfln("    %s", choice.Name)
				}
			}

			if entry.Default != nil {
				Printfln("  Default: %v", defaultValue)
			}
		}
	}

	tmpl, err := template.New("kcfg2home").Parse(moduleTemplate)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(os.Stdout, schema); err != nil {
		return err
	}

	return nil
}

func Printfln(format string, a ...any) (n int, err error) {
	return fmt.Printf(format+"\n", a...)
}
