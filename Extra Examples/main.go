package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
)

type country struct {
	Name    string `json:"name"`
	Capital string `json:"capital"`
	ID      int64  `json:"id"`
}

//Define country gql type
var countryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Country",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"capital": &graphql.Field{
				Type: graphql.String,
			},
			"id": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

//Define the query
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		/*
		   curl -g 'http://localhost:8080/graphql?query={countries{name,capital,id}}'
		*/
		"countries": &graphql.Field{
			Type:        graphql.NewList(countryType),
			Description: "List of countries",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				_, countries := getCountryDataFrmHTTP()
				fmt.Println(countries)
				return countries, nil
			},
		},
	},
})

// define schema, with our rootQuery and rootMutation
var Countryschema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

func getCountryDataFrmHTTP() (bool, []country) {
	var result map[string][]country
	resp, err := http.Get("https://api.jsonbin.io/b/5f169718c58dc34bf5d7aaec")
	if err != nil {
		fmt.Print("Error 1:", err)
		return false, nil
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&result)
	countries := result["countries"]

	return true, countries
}

//Helper function to import json from file to map
func getCountriesData() (bool, []country) {
	var result map[string][]country
	content := `{
					"countries": [
					  {
						"name": "Andorra",
						"capital": "Andorra la Vella",
                        "id": 1
					  },
					  {
						"name": "United Arab Emirates",
						"capital": "Abu Dhabi",
  						"id": 2
					  }
					]
				}`

	err := json.Unmarshal([]byte(content), &result)
	if err != nil {
		fmt.Print("Error 1:", err)
		return false, nil
	}
	countries := result["countries"]
	return true, countries
}

func executeCountryQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {

	http.HandleFunc("/graphql/country-info", func(w http.ResponseWriter, r *http.Request) {
		result := executeCountryQuery(r.URL.Query().Get("query"), Countryschema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Load country list: curl -g 'http://localhost:8080/graphql/country-info?query={countries{name,capital}}'")
	http.ListenAndServe(":8080", nil)
}
