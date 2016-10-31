package minidom

import (
	"bytes"
	"encoding/xml"
	"fmt"
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
			xpath := strings.Join(path, "/")
			// check for repeatable modification, if so, replace the past path
			if i, ok := counter[xpath]; ok {
				counter[xpath] = (i + 1)
				path[len(path)-1] = fmt.Sprintf("%s[%d]", t.Name.Local, i+1)
				xpath = strings.Join(path, "/")
			}
			// attributes for this element
			for _, a := range t.Attr {
				key := fmt.Sprintf("%s/@%s", xpath, a.Name.Local)
				tmp[key] = a.Value
			}
		case xml.EndElement:
			switch t.Name.Local {
			case f.Prefix:
				return tmp, nil
			}
			// see if we collected data for this element
			value := strings.TrimSpace(buf.String())
			if value != "" {
				key := strings.Join(path, "/")
				tmp[key] = value
				buf.Reset()
			}
			// pop
			path = path[0 : len(path)-1]
		case xml.CharData:
			buf.Write(xml.CharData(t))
		}
	}
}
