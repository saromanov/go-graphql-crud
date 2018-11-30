package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
)

// User defines ident of the user
type User struct {
	ID          string       `json:"id,omitempty"`
	Firstname   string       `json:"firstName"`
	Lastname    string       `json:"lastName"`
	Phone       string       `json:"phone"`
	Gender      string       `json:"gender"`
	Attachments []Attachment `json:"attachments"`
	CreatedDate time.Time
}

// Attachment defines additional link to attached info
type Attachment struct {
	ID   string `json:"id,omitempty"`
	Link string `json:"link"`
}

func Filter(users []User, f func(User) bool) []User {
	vsf := make([]User, 0)
	for _, v := range users {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func main() {
	var users []User = []User{
		User{
			ID:          "1",
			Firstname:   "Sergey",
			Lastname:    "Romanov",
			Phone:       "+7123456789",
			Gender:      "M",
			CreatedDate: time.Now().UTC(),
		},
	}

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Firstname",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"lastName": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return users, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{})
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
		})
		json.NewEncoder(w).Encode(result)
	})
	http.ListenAndServe(":8080", nil)
}
