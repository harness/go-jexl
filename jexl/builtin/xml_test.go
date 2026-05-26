// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import (
	"encoding/xml"
	"testing"
)

type xmlItem struct {
	XMLName xml.Name `xml:"item"`
	Value   string   `xml:"value"`
}

// Ensure XMLMarshal encodes a struct to XML.
func TestXMLMarshal_basic(t *testing.T) {
	out, err := XMLMarshal(xmlItem{Value: "hi"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "<item><value>hi</value></item>" {
		t.Fatalf("unexpected output: %s", out)
	}
}

// Ensure XMLSelect returns inner text of all matching nodes.
func TestXMLSelect_match(t *testing.T) {
	const doc = `<root><value>hello</value><value>world</value></root>`
	got, err := XMLSelect("//value", doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "hello" || got[1] != "world" {
		t.Fatalf("unexpected result: %v", got)
	}
}

// Ensure XMLSelect returns an empty slice when no nodes match.
func TestXMLSelect_noMatch(t *testing.T) {
	const doc = `<root><value>hello</value></root>`
	got, err := XMLSelect("//missing", doc)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

// Ensure XMLSelect returns an error for malformed XML.
func TestXMLSelect_invalidXML(t *testing.T) {
	_, err := XMLSelect("//value", "not xml")
	if err == nil {
		t.Fatal("expected error")
	}
}

// Ensure XMLSelect returns an error for an invalid XPath expression.
func TestXMLSelect_invalidXPath(t *testing.T) {
	_, err := XMLSelect("///[[[", "<root/>")
	if err == nil {
		t.Fatal("expected error")
	}
}

// Ensure XMLMarshalIndent produces indented XML output.
func TestXMLMarshalIndent_basic(t *testing.T) {
	out, err := XMLMarshalIndent(xmlItem{Value: "hi"}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	want := "<item>\n  <value>hi</value>\n</item>"
	if out != want {
		t.Fatalf("expected %q, got %q", want, out)
	}
}
