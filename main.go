package main

import (
	"fmt"
	"net/http"

	"github.com/DGMeyers/DnD/character"
)

func main() {
	fmt.Println("starting server")
	http.HandleFunc("/createCharacter", character.CreateCharacter)
	http.HandleFunc("/", character.DisplayCharacters)
	http.ListenAndServe(":8080", nil)
	fmt.Println("listening on server")
}
