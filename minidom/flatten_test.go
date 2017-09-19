package minidom

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var data = `
    <Listing source="mls">
     <ListingID>one</ListingID>
     <Status>active</Status>
     <URL>http://example.com/property/one.html</URL>
     <Photos LastModified="2016-10-20T05:23:23Z">
				<Photo sequence="1" href="http://example.com/property/1.jpg">
					<Label name="main"/>
					<Label name="front"/>
				</Photo>
				<Photo sequence="2" href="http://example.com/property/2.jpg">
					<Label name="interior"/>
					<Label name="kitchen"/>
				</Photo>
				<Photo sequence="3" href="http://example.com/property/3.jpg">
					<Label name="backyard"/>
					<Label name="pool"/>
				</Photo>
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

func TestXPather(t *testing.T) {
	path := []string{
		"Photos",
		"Photo",
		"Label",
	}
	counts := map[string]int{
		"Photos/Photo":       2,
		"Photos/Photo/Label": 2,
	}
	xpath := XPath(path).Index(counts)
	assert.Equal(t, "Photos/Photo[2]/Label[2]", xpath)
}

func TestFlattenNoPrefix(t *testing.T) {
	flatten := Flatten{
		Prefix: "Listing",
		Repeatable: []string{
			"Photos/Photo",
			"Photos/Photo/Label",
		},
		OmitPrefix: true,
	}
	parser := xml.NewDecoder(strings.NewReader(data))
	flat, err := flatten.Map(parser)
	assert.Nil(t, err)

	assert.Equal(t, "one", flat["ListingID"])
	assert.Equal(t, "active", flat["Status"])
	assert.Equal(t, "http://example.com/property/one.html", flat["URL"])
	assert.Equal(t, "http://example.com/property/1.jpg", flat["Photos/Photo[1]/@href"])
	assert.Equal(t, "1", flat["Photos/Photo[1]/@sequence"])
	assert.Equal(t, "main", flat["Photos/Photo[1]/Label[1]/@name"])
	assert.Equal(t, "http://example.com/property/2.jpg", flat["Photos/Photo[2]/@href"])
	assert.Equal(t, "2", flat["Photos/Photo[2]/@sequence"])
	assert.Equal(t, "interior", flat["Photos/Photo[2]/Label[1]/@name"])
	assert.Equal(t, "http://example.com/property/3.jpg", flat["Photos/Photo[3]/@href"])
	assert.Equal(t, "3", flat["Photos/Photo[3]/@sequence"])
	assert.Equal(t, "2245 Don Knotts Blvd.", flat["Address/FullStreetAddress"])
	assert.Equal(t, "2", flat["Address/UnitNumber"])
	assert.Equal(t, "WV", flat["Address/StateOrProvince"])
	assert.Equal(t, "26501", flat["Address/PostalCode"])
	assert.Equal(t, "US", flat["Address/Country"])
	assert.Equal(t, "234000", flat["ListPrice"])
	assert.Equal(t, "Public", flat["ListPrice/@isgSecurityClass"])

	if t.Failed() {
		b, _ := json.MarshalIndent(flat, "", "  ")
		fmt.Print(string(b))
	}
}

func TestFlattenWithPrefix(t *testing.T) {
	flatten := Flatten{
		Prefix: "Listing",
		Repeatable: []string{
			"Listing/Photos/Photo",
			"Listing/Photos/Photo/Label",
		},
		OmitPrefix: false,
	}
	parser := xml.NewDecoder(strings.NewReader(data))
	flat, err := flatten.Map(parser)
	assert.Nil(t, err)

	assert.Equal(t, "mls", flat["Listing/@source"])
	assert.Equal(t, "one", flat["Listing/ListingID"])
	assert.Equal(t, "active", flat["Listing/Status"])
	assert.Equal(t, "http://example.com/property/one.html", flat["Listing/URL"])
	assert.Equal(t, "http://example.com/property/1.jpg", flat["Listing/Photos/Photo[1]/@href"])
	assert.Equal(t, "1", flat["Listing/Photos/Photo[1]/@sequence"])
	assert.Equal(t, "http://example.com/property/2.jpg", flat["Listing/Photos/Photo[2]/@href"])
	assert.Equal(t, "2", flat["Listing/Photos/Photo[2]/@sequence"])
	assert.Equal(t, "http://example.com/property/3.jpg", flat["Listing/Photos/Photo[3]/@href"])
	assert.Equal(t, "3", flat["Listing/Photos/Photo[3]/@sequence"])
	assert.Equal(t, "2245 Don Knotts Blvd.", flat["Listing/Address/FullStreetAddress"])
	assert.Equal(t, "2", flat["Listing/Address/UnitNumber"])
	assert.Equal(t, "WV", flat["Listing/Address/StateOrProvince"])
	assert.Equal(t, "26501", flat["Listing/Address/PostalCode"])
	assert.Equal(t, "US", flat["Listing/Address/Country"])
	assert.Equal(t, "234000", flat["Listing/ListPrice"])
	assert.Equal(t, "Public", flat["Listing/ListPrice/@isgSecurityClass"])

	if t.Failed() {
		b, _ := json.MarshalIndent(flat, "", "  ")
		fmt.Print(string(b))
	}
}

func TestFlattenNoRepeatable(t *testing.T) {
	flatten := Flatten{
		Prefix:     "Listing",
		OmitPrefix: true,
	}
	parser := xml.NewDecoder(strings.NewReader(data))
	flat, err := flatten.Map(parser)
	assert.Nil(t, err)

	assert.Equal(t, "one", flat["ListingID"])
	assert.Equal(t, "active", flat["Status"])
	assert.Equal(t, "http://example.com/property/one.html", flat["URL"])
	assert.Equal(t, "", flat["Photos/Photo[1]"])
	assert.Equal(t, "", flat["Photos/Photo[1]/@sequence"])
	assert.Equal(t, "", flat["Photos/Photo[2]"])
	assert.Equal(t, "", flat["Photos/Photo[2]/@sequence"])
	assert.Equal(t, "http://example.com/property/3.jpg", flat["Photos/Photo/@href"])
	assert.Equal(t, "3", flat["Photos/Photo/@sequence"])
	assert.Equal(t, "2245 Don Knotts Blvd.", flat["Address/FullStreetAddress"])
	assert.Equal(t, "2", flat["Address/UnitNumber"])
	assert.Equal(t, "WV", flat["Address/StateOrProvince"])
	assert.Equal(t, "26501", flat["Address/PostalCode"])
	assert.Equal(t, "US", flat["Address/Country"])
	assert.Equal(t, "234000", flat["ListPrice"])
	assert.Equal(t, "Public", flat["ListPrice/@isgSecurityClass"])

	if t.Failed() {
		b, _ := json.MarshalIndent(flat, "", "  ")
		fmt.Print(string(b))
	}
}
