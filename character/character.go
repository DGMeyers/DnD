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

// compile all templates and cache them
var templates = template.Must(template.ParseGlob("assets/*"))

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
		characters = append(characters, currentCharacter)
	}
	if err = cur.Err(); err != nil {
		log.Fatal(err)
	}
	// Initialize a new template
	//	t := template.New("characters")
	// Parse the template
	//	parsedT, err := t.ParseFiles("assets/characters.tmpl")
	//	if err != nil {
	//		fmt.Printf("Error Parsing file: %+v", err)
	//	}
	// Execute the template to the indicated io writer
	//	fmt.Printf("Here are my characters: %+v", characters)
	//	err = parsedT.Execute(w, characters)
	//	if err != nil {
	//		fmt.Printf("\n\nError executing template: %+v\n\n", err)
	//	}
	// you access the cached templates with the defined name, not the filename
	err = templates.ExecuteTemplate(w, "characters", characters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateCharacter is the handler for the GET/POST requests on the /createCharacter endpoint
func CreateCharacter(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("hewwo fwiend")
		err := templates.ExecuteTemplate(w, "createCharacter", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// SaveCharacter will save a character struct to a mongodb instance
func SaveCharacter(characterToSave Character) {
	client := connectToDatabase(dbURL)
	collection := client.Database("DnD").Collection("characters")
	_, err := collection.InsertOne(context.Background(), characterToSave)
	if err != nil {
		log.Fatal(err)
	}
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
