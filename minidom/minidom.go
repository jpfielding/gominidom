package minidom

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
)

// EachDOM provides a walking function for a substream
type EachDOM func(io.ReadCloser, error) error

// QuitAt is a default implementation for EndFunc
func QuitAt(stop string) func(xml.EndElement) bool {
	return func(end xml.EndElement) bool {
		return end.Name.Local == stop
	}
}

// MiniDom provides the lifecycle management for walking streams of doms
type MiniDom struct {
	// StartFunc listens to start elems outside of Prefix
	StartFunc func(xml.StartElement)
	// EndFunc listens to the end elems outside of Prefix, bool returns whether an exit is requested
	EndFunc func(xml.EndElement) bool
}

// Matcher provides a swappable matcher for finding elems
type Matcher func(xml.StartElement) bool

// ByName is the simple default for matching
func ByName(match string) Matcher {
	return func(t xml.StartElement) bool {
		return t.Name.Local == match
	}
}

// Walk finds the next <prefix> and produces an io.ReadCloser of the <prefix>...</prefix> sub elem
func (md MiniDom) Walk(parser *xml.Decoder, match Matcher, each EachDOM) error {
	for {
		token, err := parser.Token()
		if err != nil {
			return err
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch {
			case match(t):
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
	err := collect.EncodeToken(start)
	if err != nil {
		return nil
	}
	for {
		token, err := parser.Token()
		if err != nil {
			return err
		}
		// namespaces arent handled here
		switch t := token.(type) {
		case xml.StartElement:
			// recurse
			err = mini(collect, t, parser)
			if err != nil {
				return nil
			}
		case xml.EndElement:
			// write end elem, and return
			return collect.EncodeToken(t)
		default:
			// write other element
			err = collect.EncodeToken(t)
			if err != nil {
				return nil
			}
		}
	}

}
