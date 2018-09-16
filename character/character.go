package character

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const dbURL = "mongodb://localhost:12345"

// compile all templates and cache them
var templates = template.Must(template.ParseGlob("assets/*"))

// Fight will take two forms for characters and display them back
func Fight(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates.ExecuteTemplate(w, "fight", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// Parse form for two characters
	// Grab the two characters from the database
	// Pass back the that booty and the two characters through the template execute
	// The template will be called "fight"
}

// Character is a representation for a DnD character
type Character struct {
	Name             string
	Career           string
	Race             string
	CurrentXP        int
	TotalXP          int
	PersonalDetails  personalDetails
	CharacterProfile characterProfile
	Weapons          []weapons
	Armor            []armour
}

// DisplayCharacters is being displayed for user
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
	err = templates.ExecuteTemplate(w, "characters", characters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateCharacter is the handler for the GET/POST requests on the /createCharacter endpoint
func CreateCharacter(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates.ExecuteTemplate(w, "createCharacter", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		r.ParseForm()
		name := r.Form.Get("name")
		career := r.Form.Get("career")
		race := r.Form.Get("race")
		newCharacter := Character{
			Name:   name,
			Career: career,
			Race:   race,
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
