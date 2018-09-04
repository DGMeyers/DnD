package character

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const dbURL = "mongodb://localhost:12345"

// Character is a representation for a DnD character
type Character struct {
	Name  string
	Class string
	Race  string
	Age   int
}

func DisplayCharacters(w http.ResponseWriter, r *http.Request) {
	var characters []Character
	client := connectToDatabase(dbURL)
	collection := client.Database("DnD").Collection("characters")
	cur, err := collection.Find(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		elem := bson.NewDocument()
		err = cur.Decode(elem)
		if err != nil {
			log.Fatal(err)
		}
		var currentCharacter Character
		bsonBytes, _ := bson.Marshal(elem)
		bson.Unmarshal(bsonBytes, &currentCharacter)
		fmt.Printf("connect to currentCharacter: %+v", currentCharacter)
		characters = append(characters, currentCharacter)
	}
	if err = cur.Err(); err != nil {
		log.Fatal(err)
	}
	// Initialize a new template
	t := template.New("characters")
	// Parse the template
	parsedT, err := t.ParseFiles("assets/characters.tmpl")
	if err != nil {
		fmt.Printf("Error Parsing file: %+v", err)
	}
	// Execute the template to the indicated io writer
	fmt.Printf("Here are my characters: %+v", characters)
	err = parsedT.Execute(w, characters)
	if err != nil {
		fmt.Printf("\n\nError executing template: %+v\n\n", err)
	}
}

// CreateCharacter is the handler for the GET/POST requests on the /createCharacter endpoint
func CreateCharacter(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Initialize a new template
		t := template.New("createCharacter")
		// Parse the template
		parsedT, err := t.ParseFiles("assets/createCharacter.tmpl")
		if err != nil {
			fmt.Printf("Error Parsing file: %+v", err)
		}
		// Execute the template to the indicated io writer
		err = parsedT.Execute(w, nil)
		if err != nil {
			fmt.Printf("\n\nError executing template: %+v\n\n", err)
		}
	} else if r.Method == "POST" {
		r.ParseForm()
		name := r.Form.Get("name")
		class := r.Form.Get("class")
		race := r.Form.Get("race")
		age, _ := strconv.Atoi(r.Form.Get("age"))
		newCharacter := Character{
			Name:  name,
			Age:   age,
			Class: class,
			Race:  race,
		}
		SaveCharacter(newCharacter)
	}
}

// SaveCharacter will save a character struct to a mongodb instance
func SaveCharacter(characterToSave Character) {
	client := connectToDatabase(dbURL)
	collection := client.Database("DnD").Collection("characters")
	res, err := collection.InsertOne(context.Background(), characterToSave)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("responseobject: %+v", res)
}

func connectToDatabase(url string) *mongo.Client {
	client, err := mongo.NewClient(url)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return client
}
