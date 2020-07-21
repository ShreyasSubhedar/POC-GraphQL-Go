package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"
)

type categories struct {
	ID           string `json:"id"`
	CategoryName string `json:"categoryName"`
	TimeStamp    string `json:"timeStamp"`
}

var data []categories

func restAPICall() {
	url := "https://blogbid.000webhostapp.com/api/categories/read.php"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
}

/*
   Create User object type with fields "name", "categories" and "timeStamp" by using GraphQLObjectTypeConfig:
       - Name: name of object type
       - Fields: a map of fields by using GraphQLFields
   Setup type of field use GraphQLFieldConfig
*/
var categoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "categories",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"categoryName": &graphql.Field{
				Type: graphql.String,
			},
			"timeStamp": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

/*
   Create Query object type with fields "data"  by using GraphQLObjectTypeConfig:
*/
var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"data": &graphql.Field{
				Type: graphql.NewList(categoryType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {

					restAPICall()
					return data, nil

				},
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func main() {
	query := `
	{
		data
		{
			id
			categoryName
			timeStamp
		}
	}
	`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	fmt.Println(result.Data)
}
