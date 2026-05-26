// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import (
	"encoding/xml"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/harness/go-jexl/jexl/coerce"
)

// XMLMarshal returns the XML encoding of v.
func XMLMarshal(v any) (string, error) {
	b, err := xml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// XMLSelect returns the text content of all nodes matching the XPath expression.
func XMLSelect(expr any, doc any) ([]string, error) {
	root, err := xmlquery.Parse(
		strings.NewReader(coerce.ToString(doc)),
	)
	if err != nil {
		return nil, err
	}
	nodes, err := xmlquery.QueryAll(root, coerce.ToString(expr))
	if err != nil {
		return nil, err
	}
	out := make([]string, len(nodes))
	for i, n := range nodes {
		out[i] = n.InnerText()
	}
	return out, nil
}

// XMLMarshalIndent returns the indented XML encoding of v.
func XMLMarshalIndent(v any, prefix, indent any) (string, error) {
	b, err := xml.MarshalIndent(
		v,
		coerce.ToString(prefix),
		coerce.ToString(indent),
	)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
