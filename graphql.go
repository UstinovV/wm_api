package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"log"
)

type Offer struct {
	ID 		string
	Title   string
	Content string
	Author  Author
}

type Author struct {
	Name      string
	Offers 	  []string
	Company   Company
}

type Company struct {
	Title string
}

var companyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Company",
		Fields: graphql.Fields{
			"title": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var authorType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Author",
		Fields: graphql.Fields{
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Offers": &graphql.Field{
				// we'll use NewList to deal with an array
				// of int values
				Type: graphql.NewList(graphql.String),
			},
			"Company": &graphql.Field{
				Type: companyType,
			},
		},
	},
)

var offerType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "OFfer",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: authorType,
			},

		},
	},
)

func populate() []Offer {
	author := &Author{Name: "Elliot Forbes", Offers: []string{"123qwe"}}
	tutorial := Offer{
		ID:     "123qwe",
		Title:  "Go GraphQL Tutorial",
		Author: *author,
	}

	var offers []Offer
	offers = append(offers, tutorial)

	return offers
}

func main() {
	offers := populate()

	fields := graphql.Fields{
		"offer": &graphql.Field{
			Type:        offerType,
			// it's good form to add a description
			// to each field.
			Description: "Get Offer By ID",
			// We can define arguments that allow us to
			// pick specific tutorials. In this case
			// we want to be able to specify the ID of the
			// tutorial we want to retrieve
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// take in the ID argument
				id, ok := p.Args["id"].(string)
				if ok {
					// Parse our tutorial array for the matching id
					for _, offer := range offers {
						if string(offer.ID) == id {
							// return our tutorial
							return offer, nil
						}
					}
				}
				return nil, nil
			},
		},
		// this is our `list` endpoint which will return all
		// tutorials available
		"list": &graphql.Field{
			Type:        graphql.NewList(offerType),
			Description: "Get Offers List",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return offers, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := `
    {
        list {
            id
            title
            author {
                Name
                Offers
				Company {
					title
				}
            }
        }
    }
`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {“data”:{“hello”:”world”}}
}
