package minidom

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	testutils "github.com/jpfielding/goTest/testutils"
)

func TestSimple(t *testing.T) {
	type Mini struct {
		ID int `xml:"id,attr"`
	}
	var data = `<xml>
    <mini id="1"></mini>
    <mini id="2"></mini>
    <mini id="3"></mini>
    </xml>`

	doms := ioutil.NopCloser(strings.NewReader(data))
	parser := xml.NewDecoder(doms)
	var mini []Mini
	md := MiniDom{
		// quit on the the xml tag
		EndFunc: func(end xml.EndElement) bool {
			return end.Name.Local == "xml"
		},
	}
	err := md.Walk(parser, "mini", func(segment io.ReadCloser, err error) error {
		tmp := Mini{}
		xml.NewDecoder(segment).Decode(&tmp)
		fmt.Printf("found mini %v\n", tmp)
		mini = append(mini, tmp)
		return err
	})
	testutils.Ok(t, err)
	testutils.Equals(t, 3, len(mini))
	testutils.Equals(t, 1, mini[0].ID)
}
