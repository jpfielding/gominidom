package minidom

import (
	"encoding/xml"
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
		mini = append(mini, tmp)
		return err
	})
	testutils.Ok(t, err)
	testutils.Equals(t, 3, len(mini))
	testutils.Equals(t, 1, mini[0].ID)
	testutils.Equals(t, 2, mini[1].ID)
	testutils.Equals(t, 3, mini[2].ID)
}

func TestComplex(t *testing.T) {
	type Address struct {
		FullStreetAddres string
		UnitNumber       int
		City             string
		State            string
		PostalCode       string
		Country          string
	}
	type ListPrice struct {
		SecurityClass string `xml:"isgSecurityClass,attr"`
		Value         int    `xml:",chardata"`
	}
	type Listing struct {
		ID        string `xml:"ListingID"`
		Address   Address
		ListPrice ListPrice
	}

	var data = `<xml>
    <Listings>
        <Listing>
         <ListingID>one</ListingID>
         <Address>
             <FullStreetAddress>2245 Don Knotts Blvd.</FullStreetAddress>
             <UnitNumber>2</UnitNumber>
             <City>Morgantown</City>
             <StateOrProvince>WV</StateOrProvince>
             <PostalCode>26501</PostalCode>
             <Country>US</Country>
         </Address>
         <ListPrice isgSecurityClass="Public">234000</ListPrice>
         </Listing>
         <Listing>
         <ListingID>two</ListingID>
          <Address>
              <FullStreetAddress>453 Suncrest Towncentre.</FullStreetAddress>
              <UnitNumber>200</UnitNumber>
              <City>Morgantown</City>
              <StateOrProvince>WV</StateOrProvince>
              <PostalCode>26505</PostalCode>
              <Country>US</Country>
          </Address>
          <ListPrice isgSecurityClass="Public">5000000</ListPrice>
          </Listing>
    </Listings>`

	doms := ioutil.NopCloser(strings.NewReader(data))
	parser := xml.NewDecoder(doms)
	var listings []Listing
	md := MiniDom{
		// quit on the the xml tag
		EndFunc: func(end xml.EndElement) bool {
			return end.Name.Local == "Listings"
		},
	}
	err := md.Walk(parser, "Listing", func(segment io.ReadCloser, err error) error {
		tmp := Listing{}
		xml.NewDecoder(segment).Decode(&tmp)
		listings = append(listings, tmp)
		return err
	})
	testutils.Ok(t, err)
	testutils.Equals(t, 2, len(listings))
	testutils.Equals(t, "one", listings[0].ID)
	testutils.Equals(t, "two", listings[1].ID)
}

func TestComplexStartEnd(t *testing.T) {
	type Response struct {
		Count   int `xml:"Records,attr"`
		Maxrows bool
	}
	type Listing struct {
		ID string `xml:"ListingID"`
	}
	var data = `
    <?xml version="1.0" encoding="utf-8"?>
    <RETS ReplyCode="0" ReplyText="Operation successful.">
      <COUNT Records="74" />
      <REData>
        <REProperties>
          <Residential>
            <PropertyListing><ListingID>one</ListingID></PropertyListing>
            <PropertyListing><ListingID>two</ListingID></PropertyListing>
          </Residential>
        </REProperties>
      </REData>
      <MAXROWS/>
    </RETS>
    `

	response := Response{}
	doms := ioutil.NopCloser(strings.NewReader(data))
	parser := xml.NewDecoder(doms)
	var listings []Listing
	md := MiniDom{
		StartFunc: func(start xml.StartElement) {
			switch start.Name.Local {
			case "COUNT":
				parser.DecodeElement(&response, &start)
			case "MAXROWS":
				response.Maxrows = true
			}
		},
		// quit on the the xml tag
		EndFunc: func(end xml.EndElement) bool {
			switch end.Name.Local {
			case "RETS", "RETS-STATUS":
				return true
			}
			return false
		},
	}
	err := md.Walk(parser, "PropertyListing", func(segment io.ReadCloser, err error) error {
		tmp := Listing{}
		xml.NewDecoder(segment).Decode(&tmp)
		listings = append(listings, tmp)
		return err
	})
	testutils.Ok(t, err)
	testutils.Equals(t, 2, len(listings))
	testutils.Equals(t, "one", listings[0].ID)
	testutils.Equals(t, "two", listings[1].ID)
	testutils.Equals(t, 74, response.Count)
	testutils.Equals(t, true, response.Maxrows)
}
