package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/satori/go.uuid"
)

// User defines ident of the user
type User struct {
	ID        string `json:"id"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Phone     string `json:"phone"`
	//Gender      string       `json:"gender"`
	//Attachments []Attachment `json:"attachments"`
	//CreatedDate time.Time
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

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       context.WithValue(context.Background(), "currentUser", nil),
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {
	var users []User = []User{
		User{
			ID:        "1",
			Firstname: "Sergey",
			Lastname:  "Romanov",
			Phone:     "+7123456789",
			//Gender:      "M",
			//CreatedDate: time.Now().UTC(),
		},
	}

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"lastName": &graphql.Field{
				Type: graphql.String,
			},
			"firstName": &graphql.Field{
				Type: graphql.String,
			},
			"phone": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Args: graphql.FieldConfigArgument{
					"user": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					idQuery, ok := params.Args["id"].(string)
					if !ok {
						return User{}, nil
					}
					for _, usr := range users {
						if usr.ID == idQuery {
							return usr, nil
						}
					}

					return User{}, nil
				},
			},
		},
	})
	fmt.Println(rootQuery)

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"firstName": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"lastName": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"phone": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var user User
					user.ID = uuid.Must(uuid.NewV4()).String()
					user.Firstname = params.Args["firstName"].(string)
					user.Lastname = params.Args["lastName"].(string)
					user.Phone = params.Args["phone"].(string)
					fmt.Println(user)
					users = append(users, user)
					return user, nil
				},
			},

			"deleteUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					idQuery, ok := params.Args["id"].(string)
					if !ok {
						return User{}, nil
					}
					user := User{}
					for i, u := range users {
						if u.ID == idQuery {
							users = append(users[:i], users[:i]...)
							user = u
							break
						}
					}
					return user, nil
				},
			},
			"updateUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"firstName": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"lastName": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"phone": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					idQuery, ok := params.Args["id"].(string)
					if !ok {
						return User{}, nil
					}
					user := User{}
					user.Firstname = params.Args["firstName"].(string)
					user.Lastname = params.Args["lastName"].(string)
					user.Phone = params.Args["phone"].(string)
					for i, u := range users {
						if u.ID == idQuery {
							users[i] = user
							break
						}
					}
					return user, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})
	http.ListenAndServe(":8080", nil)
}
