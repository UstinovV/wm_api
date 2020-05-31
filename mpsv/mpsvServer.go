package mpsv

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type MpsvServer struct {

}

func (s *MpsvServer) ParseMpsvUrl (url *MpsvUrl, outStream MpsvParser_ParseMpsvUrlServer) error {
	fmt.Printf("Received urls %s", url.Url)
	xmlFile, err := os.Open("example.xml")
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)

	for {
		offer := MpsvOffer{}
		token, tokenErr := decoder.Token()
		if tokenErr != nil && tokenErr != io.EOF {
			fmt.Println("error happend", tokenErr)
			break

		} else if tokenErr == io.EOF {
			break
		}
		if token == nil {
			fmt.Println("t is nil break")
		}

		switch tok := token.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
			case "VOLNEMISTO":
				fmt.Println("found ")
				for _, attr := range tok.Attr {
					if attr.Name.Local == "uid" {
						offer.MpsvId = attr.Value
					}
					if attr.Name.Local == "zmena" {
						//offer.CreatedAt = attr.Value
					}
				}
			case "PROFESE":
				for _, attr := range tok.Attr {
					if attr.Name.Local == "nazev" {
						//fmt.Println(attr.Value)
						offer.Title = attr.Value
					}
				}
			case "POZNAMKA":
				decoder.DecodeElement(&offer.Content, &tok)

			}

		case xml.EndElement:
			if tok.Name.Local == "VOLNEMISTO" {
				 if err = outStream.Send(&MpsvOffer{
					 MpsvId:  "123",
					 Title:   "some title",
					 Content: "asdasd",
				 }); err != nil {
				 	return err
				}
			}
		}
	}
	return nil
}

