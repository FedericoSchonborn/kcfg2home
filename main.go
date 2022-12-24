package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
)

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
	Name string `xml:",innerxml" xml:",innerxml"`
	Type string `xml:"type,attr,omitempty" json:"type,omitempty"`
}

type Group struct {
	Name    string   `xml:"name,attr,omitempty" json:"name,omitempty"`
	Entries []*Entry `xml:"entry,omitempty" json:"entries,omitempty"`
}

type Entry struct {
	Name    string    `xml:"name,attr,omitempty" json:"name,omitempty"`
	Key string  `xml:"key,attr,omitempty" json:"key,omitempty"`
	Type    string    `xml:"type,attr,omitempty" json:"type,omitempty"`
	Choices []*Choice `xml:"choices>choice,omitempty" json:"choices,omitempty"`
	Label   string    `xml:"label,omitempty" json:"label,omitempty"`
	Default *Default  `xml:"default,omitempty" json:"default,omitempty"`
	Emit    *Emit     `xml:"emit,omitempty" json:"emit,omitempty"`
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
		fmt.Fprintf(os.Stderr, "kcfg2nix: %v\n", err)
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

	json, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("%s", json)
	return nil
}
