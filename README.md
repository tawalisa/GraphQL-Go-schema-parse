# GraphQL-Go-Schema-Parser

The project provides a tool to support from a SDL to GraphQL-go struct.
Since I have not found a schema parsing for [GraphQL-GO](https://github.com/graphql-go/graphql). Then I try to implement it. So far, there are some limitations. But it has been used to my project.

## Getting started for GraphQL-GO

The guide's example is from [Getting started with Spring for GraphQL](https://www.graphql-java.com/tutorials/getting-started-with-spring-boot)

In this tutorial, you will create a GraphQL in go using [GraphQL-GO](https://github.com/graphql-go/graphql). It requires a little Golang knowledge. While we give a brief introduction to GraphQL, the focus of this tutorial is developing a GraphQL server in Golang.

## A very short introduction to GraphQL


GraphQL is a query language to retrieve data from a server. It is an alternative to REST, SOAP or gRPC.

Let's suppose we want to query the details for a specific book from an online store backend.

With GraphQL you send the following query to the server to get the details for the book with the id "book-1":

```graphql
query bookDetails {
  bookById(id: "book-1"){
    id
    name
    pageCount
    author {
      firstName
      lastName
    }
  }
}
```
This is not JSON (even though it looks deliberately similar), it is a GraphQL query. It basically says:

query a book with a specific 

1. get me the id, name, pageCount and author from that book
2. for the author, I want to know the firstName and lastName
3. The response is normal JSON:

```json
{
  "bookById": {
    "id":"book-1",
    "name":"Harry Potter and the Philosopher's Stone",
    "pageCount":223,
    "author": {
      "firstName":"Joanne",
      "lastName":"Rowling"
    }
  }
}
```

One very important property of GraphQL is that it is statically typed: the server knows exactly the shape of every object you can query and any client can actually "introspect" the server and ask for the "schema". The schema describes what queries are possible and what fields you can get back. (Note: when we refer to schema here, we always refer to a "GraphQL Schema", which is not related to other schemas like "JSON Schema" or "Database Schema")

The schema for the above query looks like this:

```graphql
type Query {
  bookById(id: ID): Book
}

type Book {
  id: ID
  name: String
  pageCount: Int
  author: Author
}

type Author {
  id: ID
  firstName: String
  lastName: String
}
```

This tutorial will focus on how to implement a GraphQL-Engine with this schema in GraphQL-GO.

The main steps of creating a GraphQL-GO engine are:

Defining a GraphQL Schema.
Deciding on how the actual data for a query is fetched.

## Our example API: getting book details

Our example app will be a simple function to get details for a specific book. This is in no way a comprehensive API, but it is enough for this tutorial.

## Prepare the environment

go 1.21.3 (This is my view, I support it should lower version)

### Schema

Here is a schema file `resource/schema/examplev1.sdl`
So far, There is a problem with dependency. This feature has not yet been implemented, so the order of the schema maintains a dependency relationship, ensuring that referenced types must be defined beforehand. This issue should be corrected in the future.

```schema
type Author {
    id: ID
    firstName: String
    lastName: String
}

type Book {
    id: ID
    name: String
    pageCount: Int
    author: Author
}

type Query {
    bookById(id: ID): Book
}
```

This schema defines one top level field (in the type `Query`): `bookById` which returns the details of a specific book.

It also defines the type `Book` which has the fields: `id`, `name`, `pageCount` and `author`. `author` is of type `Author`, which is defined after `Book`.

The Domain Specific Language (shown above) used to describe a schema is called the Schema Definition Language or SDL. More details about it can be found [here](https://graphql.org/learn/schema/).

### Source of the data

To simplify the tutorial, book and author data will come from static lists inside their respective classes. It is very important to understand that GraphQL doesn't dictate in any way where the data comes from. This is the power of GraphQL: it can come from a static in-memory list, from a database or an external service.

### Create the Book

The code is here `schemaParser/schemaParsing_test.go`

```go
type Book struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PageCount int    `json:"pageCount"`
	AuthorID  string `json:"authorId"`
}

var books = []Book{
	{ID: "book-1", Name: "Harry Potter and the Philosopher's Stone", PageCount: 223, AuthorID: "author-1"},
	{ID: "book-2", Name: "Moby Dick", PageCount: 635, AuthorID: "author-2"},
	{ID: "book-3", Name: "Interview with the Vampire", PageCount: 371, AuthorID: "author-3"},
}

func getById(id string) *Book {
	for _, book := range books {
		if book.ID == id {
			return &book
		}
	}
	return nil
}
```
# Create the Author

The code is here `schemaParser/schemaParsing_test.go`


```go
type Author struct {
	ID        string
	FirstName string
	LastName  string
}

var authors = []Author{
	{ID: "author-1", FirstName: "Joanne", LastName: "Rowling"},
	{ID: "author-2", FirstName: "Herman", LastName: "Melville"},
	{ID: "author-3", FirstName: "Anne", LastName: "Rice"},
}

func getAuthorById(id string) *Author {
	for _, author := range authors {
		if author.ID == id {
			return &author
		}
	}
	return nil
}
```

### Adding code to fetch data

Datafetcher need define a map and pass to func Parsing.

```go
func Parsing(sdlContent string, queryFucMap map[string]graphql.FieldResolveFn, schemaFucMap map[string]graphql.FieldResolveFn) (*graphql.Object, error) 
```

```go
QueryType, _ := Parsing(string(sdl),
		map[string]graphql.FieldResolveFn{
			"bookById": func(p graphql.ResolveParams) (interface{}, error) {
				return getById(p.Args["id"].(string)), nil
			}},
		map[string]graphql.FieldResolveFn{
			"Author": func(p graphql.ResolveParams) (interface{}, error) {
				if book, ok := p.Source.(*Book); ok {
					return getAuthorById(book.AuthorID), nil
				}
				return nil, nil
			},
		})
```

The 2nd parameter `queryFucMap map[string]graphql.FieldResolveFn` is a `map` which binds this function to a `query`, a field under the Query type. 
The `key` of `map` binds function name in `Query`. Here `bookById` will bind the `value` `function` of `map` 

The 3rd parameter `schemaFucMap map[string]graphql.FieldResolveFn` is a `map` which binds this function to a `Type`
The `key` of `map` binds function name in type 'Author'. Here `author` in `Book` will bind the `value` `function` of `map` 

### Run the query

Here is graphQL query. 
```graphql
query bookDetails {
  bookById(id: "book-1") {
    id
    name
    pageCount
    author {
      id
      firstName
      lastName
    }
  }
}
```

The code is here `schemaParser/schemaParsing_test.go`

```go
func TestGraphql(t *testing.T) {
	root, _ := os.Getwd()
	println(root)
	sdl, e := fileutil.ReadFile("resource/schema/examplev1.sdl")
	require.NoError(t, e)
	assert.NotEqual(t, sdl, "")
	QueryType, _ := Parsing(string(sdl),
		map[string]graphql.FieldResolveFn{
			"bookById": func(p graphql.ResolveParams) (interface{}, error) {
				return getById(p.Args["id"].(string)), nil
			}},
		map[string]graphql.FieldResolveFn{
			"Author": func(p graphql.ResolveParams) (interface{}, error) {
				if book, ok := p.Source.(*Book); ok {
					return getAuthorById(book.AuthorID), nil
				}
				return nil, nil
			},
		})
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: QueryType,
	})

	query := `
query bookDetails {
  bookById(id: "book-1") {
    id
    name
    pageCount
    author {
      id
      firstName
      lastName
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
	fmt.Printf("%s \n", rJSON) // {"data":{"hello":"world"}}
}

```

we can see the result in console.

```text
{
    "data": {
        "bookById": {
            "author": {
                "firstName": "Herman",
                "id": "author-2",
                "lastName": "Melville"
            },
            "id": "book-2",
            "name": "Moby Dick",
            "pageCount": 635
        }
    }
}
```
