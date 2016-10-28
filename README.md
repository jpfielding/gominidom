gominidom
======

Python's minidom in Go

[![Build Status](https://travis-ci.org/jpfielding/gominidom.svg?branch=master)](https://travis-ci.org/jpfielding/gominidom.


    ```
    	in := ioutil.NopCloser(...)
    	parser := xml.NewDecoder(in)
    	listings := syndication.Listings{}

    	// minidom is crazy useful for massive streams
    	md := MiniDom{
    			StartFunc: func(start xml.StartElement) {
    				switch start.Name.Local {
    				case "Listings":
                        attrs := map[string]string{}
        				for _, v := range start.Attr {
        					attrs[v.Name.Local] = v.Value
        				}
        				listings.ListingsKey = attrs["listingsKey"]
        				listings.Version = attrs["version"]
        				listings.VersionTimestamp = attrs["versionTimestamp"]
        				listings.Language = attrs["lang"]
                    case "Disclaimer":
        				parser.DecodeElement(listings.Disclaimer, &start)
    				}
    			},
    			// quit on the the xml tag
    			EndFunc: QuitAt("Listings"),
    		}
    	}
    	err := md.Walk(parser, ByName("Listing"), syndication.ToListing(func(l Listing, err error) error {
    		// .... process the listing here
    		return err
    	}))

    ```
