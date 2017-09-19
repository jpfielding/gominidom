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
	tmp := map[string]string{}
	var buf bytes.Buffer
	var path []string
	counter := map[string]int{}
	for _, e := range f.Repeatable {
		counter[e] = 0
	}
	for {
		token, err := parser.Token()
		if err != nil {
			return tmp, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			if f.OmitPrefix && t.Name.Local == f.Prefix && len(path) == 0 {
				continue
			}
			// push the path
			path = append(path, t.Name.Local)
			xpath := XPath(path).String()
			// check for repeatable modification, if so, replace the past path
			if i, ok := counter[xpath]; ok {
				counter[xpath] = (i + 1)
			}
			// apply the counts for the path
			indexed := XPath(path).Index(counter)
			// attributes for this element
			for _, a := range t.Attr {
				tmp[fmt.Sprintf("%s/@%s", indexed, a.Name.Local)] = a.Value
			}
		case xml.EndElement:
			switch t.Name.Local {
			case f.Prefix:
				return tmp, nil
			}
			// see if we collected data for this element
			value := strings.TrimSpace(buf.String())
			if value != "" {
				tmp[XPath(path).Index(counter)] = value
				buf.Reset()
			}
			// clear any counters for sub children
			xpath := XPath(path).String()
			for k := range counter {
				if strings.HasPrefix(k, xpath+"/") {
					counter[k] = 0
				}
			}
			// pop
			path = path[0 : len(path)-1]
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
