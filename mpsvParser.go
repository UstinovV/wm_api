package main

import (
	"encoding/xml"
	"fmt"
	"github.com/UstinovV/wm_api/database"
	"io"
	"os"
)

func main() {
	xmlFile, err := os.Open("example.xml")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)

	for {
		offer := database.Offer{}
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
							fmt.Println(attr.Value)
						}
						if attr.Name.Local == "zmena" {
							//offer.CreatedAt = attr.Value
							fmt.Println(attr.Value)
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
					var str string
					decoder.DecodeElement(&str, &tok)
					offer.Content = str
					//fmt.Println(str)

			}

		case xml.EndElement:
			if tok.Name.Local == "VOLNEMISTO" {
				fmt.Println(offer.Content)
				fmt.Println("save to db")
			}
		}
	}
}
