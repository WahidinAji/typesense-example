package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

func main() {
	searchUrl, ok := os.LookupEnv("TYPESENSE_HOST")
	if !ok {
		searchUrl = "http://localhost:8108"
	}

	searchApiKey, ok := os.LookupEnv("TYPESENSE_API_KEY")
	if !ok {
		searchApiKey = ""
	}
	fmt.Print(searchApiKey)
	client := typesense.NewClient(
		typesense.WithServer(searchUrl),
		typesense.WithAPIKey(searchApiKey),
	)
	// client := typesense.NewClient(
	// 	typesense.WithServer("http://localhost:8008"),
	// 	typesense.WithAPIKey(os.Getenv("TYPESENSE_API_KEY")),
	// )
	// fmt.Print(os.Getenv("TYPESENSE_API_KEY"))

	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!")
	})
	app.Post("/api/documents", func(c *fiber.Ctx) error {
		// d, err := client.Collection("companies").Retrieve()
		// if err != nil {
		// 	if err != nil {
		// 		return c.Status(404).JSON(fiber.Map{
		// 			"error": err,
		// 		})
		// 	}
		// }

		_, err := client.Collection("companies").Retrieve()
		if err != nil {
			schema := &api.CollectionSchema{
				Name: "companies",
				Fields: []api.Field{
					{
						Name: "company_name",
						Type: "string",
					},
					{
						Name: "num_employees",
						Type: "int32",
					},
					{
						Name:  "country",
						Type:  "string",
						Facet: pointer.True(),
					},
					{
						Name: "created_at",
						Type: "string",
					},
					{
						Name: "updated_at",
						Type: "string",
					},
					{
						Name: "posts",
						Type: "geopoint[]",
					},
				},

				//DefaultSortingField: pointer.String("num_employees"),
			}
			bc, errSchema := client.Collections().Create(schema)
			if errSchema != nil {
				return c.Status(500).JSON(fiber.Map{
					"error":   500,
					"message": "failed to create schema collection: " + errSchema.Error(),
				})
			}
			log.Printf("Creating collection: %v", bc)
		}

		// document := struct {
		// 	ID           string `json:"id"`
		// 	CompanyName  string `json:"company_name"`
		// 	NumEmployees int32  `json:"num_employees"`
		// 	Country      string `json:"country"`
		// }{
		// 	ID:           "123",
		// 	CompanyName:  "Stark Industries",
		// 	NumEmployees: 5215,
		// 	Country:      "USA",
		// }
		document := createNewDocument()
		result, err := client.Collection("companies").Documents().Create(document)
		fmt.Println(result)
		if err != nil {
			log.Printf("Error collection document: %v", err)
			return c.Status(404).JSON(fiber.Map{
				"error":   404,
				"message": err,
			})
		}
		fmt.Printf("data: %v\n", document)
		fmt.Printf("data: %v\n", result)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"data":   result,
			"status": fiber.StatusCreated,
		})
	})
	app.Get("/api/documents", func(c *fiber.Ctx) error {
		d, err := client.Collection("companies").Document("123").Retrieve()
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err,
			})
		}
		//fmt.Println(d)
		//document := struct {
		//	ID           string `json:"id"`
		//	CompanyName  string `json:"company_name"`
		//	NumEmployees int    `json:"num_employees"`
		//	Country      string `json:"country"`
		//	CreatedAt    string `json:"created_at"`
		//	UpdatedAt    string `json:"updated_at"`
		//}{}

		//client.Collection("companies").Document("123").Update(d)
		//client.Collections("companies").documents().
		return c.Status(200).JSON(d)
	})
	app.Get("/api/collections", func(c *fiber.Ctx) error {
		d, err := client.Collection("companies").Retrieve()
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err,
			})
		}
		fmt.Print(d.CollectionSchema.Name)
		a := client.Collection("companies").Synonyms()
		client.Operations()
		fmt.Print(a, d.Name)
		return c.Status(200).JSON(d)
	})
	app.Delete("/api/collections", func(c *fiber.Ctx) error {
		d, err := client.Collection("companies").Delete()
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": err,
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": string("companies collection successfully deleted!") + d.Name,
		})
	})

	app.Get("/home", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello world!",
		})
	})

	app.Post("/api/24", func(c *fiber.Ctx) (err error) {

		_, err = client.Collection("docs").Retrieve()
		if err != nil {
			schema := &api.CollectionSchema{
				Name: "docs",
				Fields: []api.Field{
					{Name: "person", Type: "string"},
					{Name: "details", Type: "object[]"},
				},
			}
			res, errSchema := client.Collections().Create(schema)
			if errSchema != nil {
				return c.Status(500).JSON(fiber.Map{
					"status":  false,
					"message": errSchema.Error(),
				})
			}
			log.Printf("Creating collection: %v", res)
		}
		doc := createNewHuman()
		resMap, err := client.Collection("docs").Documents().Create(doc)
		if err != nil {
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"status":  false,
					"message": err.Error(),
				})
			}
		}
		return c.Status(201).JSON(fiber.Map{
			"status": true,
			"data":   resMap,
		})
	})
	app.Get("/api/24", func(c *fiber.Ctx) (err error) {
		res, err := client.Collection("docs").Retrieve()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"status": true,
			"data":   res.Fields,
		})
	})
	app.Delete("/api/24", func(c *fiber.Ctx) (err error) {
		res, err := client.Collection("docs").Delete()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"status": true,
			"data":   res.Name,
		})
	})
	log.Fatal(app.Listen(":3000"))
}

type Human struct {
	Id      string   `json:"id"`
	Person  string   `json:"person"`
	Details []Detail `json:"details"`
}
type Detail struct {
	Id      string `json:"id"`
	Address string `json:"address"`
}

func createNewHuman() interface{} {
	var details []Detail
	for i := 1; i <= 5; i++ {
		detail := Detail{
			Id:      fmt.Sprint(i),
			Address: fmt.Sprintf("address of %d", i),
		}
		details = append(details, detail)
	}
	doc := Human{
		Id:      "1",
		Person:  "person one",
		Details: details,
	}
	return &doc
	//Thanks, mate.
	//	now, I can fix this issue with the object[] type
}
func createNewDocument(docIDs ...string) interface{} {
	docID := "123"
	if len(docIDs) > 0 {
		docID = docIDs[0]
	}
	type Post struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	var posts []Post
	for i := 1; i <= 5; i++ {
		h := Post{
			ID:   fmt.Sprint(i),
			Name: fmt.Sprintf("name %d", i),
		}
		posts = append(posts, h)
	}
	//fmt.Printf("%T", posts)
	fmt.Println(posts)
	document := struct {
		ID           string `json:"id"`
		CompanyName  string `json:"company_name"`
		NumEmployees int    `json:"num_employees"`
		Country      string `json:"country"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		Posts        []Post `json:"posts"`
	}{
		ID:           docID,
		CompanyName:  "Stark Industries",
		NumEmployees: 5215,
		Country:      "USA",
		CreatedAt:    time.Now().String(),
		UpdatedAt:    time.Now().String(),
		Posts:        posts,
	}
	return &document
}
