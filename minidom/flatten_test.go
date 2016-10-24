package minidom

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/jpfielding/gotest/testutils"
)

var data = `
    <Listing>
     <ListingID>one</ListingID>
     <Status>active</Status>
     <URL>http://example.com/property/one.html</URL>
     <Photos LastModified="2016-10-20T05:23:23Z">
        <Photo sequence="1">http://example.com/property/1.jpg</Photo>
        <Photo sequence="2">http://example.com/property/2.jpg</Photo>
        <Photo sequence="3">http://example.com/property/3.jpg</Photo>
     </Photos>
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
    `

func TestFlattenNoPrefix(t *testing.T) {
	flatten := Flatten{
		Prefix:     "Listing",
		Repeatable: []string{"Photos/Photo"},
		OmitPrefix: true,
	}
	parser := xml.NewDecoder(strings.NewReader(data))
	flat, err := flatten.Map(parser)
	testutils.Ok(t, err)

	testutils.Equals(t, "one", flat["ListingID"])
	testutils.Equals(t, "active", flat["Status"])
	testutils.Equals(t, "http://example.com/property/one.html", flat["URL"])
	testutils.Equals(t, "http://example.com/property/1.jpg", flat["Photos/Photo[1]"])
	testutils.Equals(t, "1", flat["Photos/Photo[1]/@sequence"])
	testutils.Equals(t, "http://example.com/property/2.jpg", flat["Photos/Photo[2]"])
	testutils.Equals(t, "2", flat["Photos/Photo[2]/@sequence"])
	testutils.Equals(t, "http://example.com/property/3.jpg", flat["Photos/Photo[3]"])
	testutils.Equals(t, "3", flat["Photos/Photo[3]/@sequence"])
	testutils.Equals(t, "2245 Don Knotts Blvd.", flat["Address/FullStreetAddress"])
	testutils.Equals(t, "2", flat["Address/UnitNumber"])
	testutils.Equals(t, "WV", flat["Address/StateOrProvince"])
	testutils.Equals(t, "26501", flat["Address/PostalCode"])
	testutils.Equals(t, "US", flat["Address/Country"])
	testutils.Equals(t, "234000", flat["ListPrice"])
	testutils.Equals(t, "Public", flat["ListPrice/@isgSecurityClass"])
}

func TestFlattenWithPrefix(t *testing.T) {
	flatten := Flatten{
		Prefix:     "Listing",
		Repeatable: []string{"Listing/Photos/Photo"},
		OmitPrefix: false,
	}
	parser := xml.NewDecoder(strings.NewReader(data))
	flat, err := flatten.Map(parser)
	testutils.Ok(t, err)

	testutils.Equals(t, "one", flat["Listing/ListingID"])
	testutils.Equals(t, "active", flat["Listing/Status"])
	testutils.Equals(t, "http://example.com/property/one.html", flat["Listing/URL"])
	testutils.Equals(t, "http://example.com/property/1.jpg", flat["Listing/Photos/Photo[1]"])
	testutils.Equals(t, "1", flat["Listing/Photos/Photo[1]/@sequence"])
	testutils.Equals(t, "http://example.com/property/2.jpg", flat["Listing/Photos/Photo[2]"])
	testutils.Equals(t, "2", flat["Listing/Photos/Photo[2]/@sequence"])
	testutils.Equals(t, "http://example.com/property/3.jpg", flat["Listing/Photos/Photo[3]"])
	testutils.Equals(t, "3", flat["Listing/Photos/Photo[3]/@sequence"])
	testutils.Equals(t, "2245 Don Knotts Blvd.", flat["Listing/Address/FullStreetAddress"])
	testutils.Equals(t, "2", flat["Listing/Address/UnitNumber"])
	testutils.Equals(t, "WV", flat["Listing/Address/StateOrProvince"])
	testutils.Equals(t, "26501", flat["Listing/Address/PostalCode"])
	testutils.Equals(t, "US", flat["Listing/Address/Country"])
	testutils.Equals(t, "234000", flat["Listing/ListPrice"])
	testutils.Equals(t, "Public", flat["Listing/ListPrice/@isgSecurityClass"])
}

func TestFlattenNoRepeatable(t *testing.T) {
	flatten := Flatten{
		Prefix:     "Listing",
		OmitPrefix: true,
	}
	parser := xml.NewDecoder(strings.NewReader(data))
	flat, err := flatten.Map(parser)
	testutils.Ok(t, err)

	testutils.Equals(t, "one", flat["ListingID"])
	testutils.Equals(t, "active", flat["Status"])
	testutils.Equals(t, "http://example.com/property/one.html", flat["URL"])
	testutils.Equals(t, "", flat["Photos/Photo[1]"])
	testutils.Equals(t, "", flat["Photos/Photo[1]/@sequence"])
	testutils.Equals(t, "", flat["Photos/Photo[2]"])
	testutils.Equals(t, "", flat["Photos/Photo[2]/@sequence"])
	testutils.Equals(t, "http://example.com/property/3.jpg", flat["Photos/Photo"])
	testutils.Equals(t, "3", flat["Photos/Photo/@sequence"])
	testutils.Equals(t, "2245 Don Knotts Blvd.", flat["Address/FullStreetAddress"])
	testutils.Equals(t, "2", flat["Address/UnitNumber"])
	testutils.Equals(t, "WV", flat["Address/StateOrProvince"])
	testutils.Equals(t, "26501", flat["Address/PostalCode"])
	testutils.Equals(t, "US", flat["Address/Country"])
	testutils.Equals(t, "234000", flat["ListPrice"])
	testutils.Equals(t, "Public", flat["ListPrice/@isgSecurityClass"])
}
