package minidom

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
)

// EachDOM provides a walking function for a substream
type EachDOM func(io.ReadCloser, error) error

// MiniDom ...
type MiniDom struct {
	// StartFunc listens to start elems outside of Prefix
	StartFunc func(xml.StartElement)
	// EndFunc listens to the end elems outside of Prefix, bool returns whether an exit is requested
	EndFunc func(xml.EndElement) bool
}

// Walk finds the next <prefix> and produces an io.ReadCloser of the <prefix>...</prefix> sub elem
func (md MiniDom) Walk(parser *xml.Decoder, prefix string, each EachDOM) error {
	for {
		token, err := parser.Token()
		if err != nil {
			return err
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case prefix:
				var buf bytes.Buffer
				enc := xml.NewEncoder(&buf)
				err = mini(enc, t, parser)
				enc.Flush()
				if err = each(ioutil.NopCloser(&buf), err); err != nil {
					return err
				}
			default:
				if md.StartFunc != nil {
					md.StartFunc(t)
				}
			}
		case xml.EndElement:
			if md.EndFunc != nil {
				exit := md.EndFunc(t)
				if exit {
					return nil
				}
			}
		}
	}
}

// recurse into the <prefix> and pipe them into the buffer
func mini(collect *xml.Encoder, start xml.StartElement, parser *xml.Decoder) error {
	// write start elem
	collect.EncodeToken(start)
	for {
		token, err := parser.Token()
		if err != nil {
			return err
		}
		// TODO we need to think about namespaces and how we transplant them
		switch t := token.(type) {
		case xml.StartElement:
			// recurse
			return mini(collect, t, parser)
		case xml.EndElement:
			// write end elem
			collect.EncodeToken(t)
			// return on end elem
			return nil
		default:
			// write other element
			collect.EncodeToken(t)
		}
	}

}
