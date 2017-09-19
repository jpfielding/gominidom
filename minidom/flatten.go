package minidom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"path"
	"strings"
)

// Flatten ...
type Flatten struct {
	Prefix     string
	Repeatable []string
	OmitPrefix bool
}

// Map ...
func (f Flatten) Map(parser *xml.Decoder) (map[string]string, error) {
	s, err := f.advanceToPrefix(parser)
	if err != nil {
		return nil, err
	}
	out, err := f.recurse(parser, s)
	if err != nil {
		return nil, err
	}
	if f.OmitPrefix {
		return f.trimPrefix(out), nil
	}
	return out, nil
}

func (f Flatten) trimPrefix(in map[string]string) map[string]string {
	tmp := map[string]string{}
	for k, v := range in {
		tmp[strings.TrimPrefix(k, fmt.Sprintf("%s/", f.Prefix))] = v
	}
	return tmp
}

func (f Flatten) advanceToPrefix(parser *xml.Decoder) (xml.StartElement, error) {
	for {
		token, err := parser.Token()
		if err != nil {
			return xml.StartElement{}, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == f.Prefix {
				return t, nil
			}
		}
	}
}

func (f Flatten) recurse(parser *xml.Decoder, start xml.StartElement) (map[string]string, error) {
	tmp := map[string]string{}
	// process this xml.StartElement's attributes
	for _, a := range start.Attr {
		tmp[fmt.Sprintf("%s/@%s", start.Name.Local, a.Name.Local)] = a.Value
	}
	// counter := map[string]int{}
	// for _, e := range f.Repeatable {
	// 	counter[e] = 0
	// }
	// prep for any char data
	var buf bytes.Buffer
	for {
		token, err := parser.Token()
		if err != nil {
			return tmp, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			// recurse into sub elements
			sub, err := f.recurse(parser, t)
			if err != nil {
				return tmp, err
			}
			// append child to the parent
			for k, v := range sub {
				tmp[fmt.Sprintf("%s/%s", start.Name.Local, k)] = v
			}
			// TODO figure out counters
		case xml.EndElement:
			// append the char data
			if value := strings.TrimSpace(buf.String()); value != "" {
				tmp[t.Name.Local] = value
			}
			switch t.Name.Local {
			case start.Name.Local:
				return tmp, nil
			}
		case xml.CharData:
			buf.Write(xml.CharData(t))
		}
	}
}

// XPath provides helpers for a path
type XPath []string

// String is the default string view
func (xp XPath) String() string {
	return path.Join(xp...)
}

// Index creates indexes a known path given count state
func (xp XPath) Index(counts map[string]int) string {
	var tmp []string
	for i, e := range xp {
		tmp = append(tmp, e)
		// create the raw sub path for a key
		sub := path.Join(xp[0 : i+1]...)
		// if we have a counter, index this elem
		if c, ok := counts[sub]; ok {
			tmp[i] = fmt.Sprintf("%s[%d]", e, c)
		}
	}
	return path.Join(tmp...)
}
