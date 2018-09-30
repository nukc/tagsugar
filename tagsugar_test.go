package tagsugar

import (
	"log"
	"testing"
)

func TestModel(t *testing.T) {
	Http = "https://cdn.github.com/"

	json := "{\"id\": 1, \"post\": 2}"
	array := "[{\"id\": 1, \"post\": 3},{\"id\": 2, \"post\": 66}]"
	model := Model{Id: 1, Name: "test", Sex: 1, Image: "test.png", Json: json, PostJson: json, ArrayJson: array}
	log.Print(model)
	Lick(&model)
	log.Print("---------- Lick after -----------")
	log.Print(model)
}

func TestUnmarshal(t *testing.T) {
	json := "{\"id\": 1, \"post\": 2}"
	model := Model{Id: 2, Json: json}
	log.Print(model.Object)
	Lick(&model)
	log.Print(model.Object)
}

func TestPostUnmarshal(t *testing.T) {
	json := "{\"id\": 1, \"post\": 2}"
	model := Model{Id: 3, PostJson: json}
	log.Print(model.Post)
	Lick(&model)
	log.Print(model.Post)
}

func TestArrayJsonUnmarshal(t *testing.T) {
	json := "[{\"id\": 1, \"post\": 3},{\"id\": 2, \"post\": 66}]"
	model := Model{Id: 4, ArrayJson: json}
	log.Print(model.Array)
	Lick(&model)
	log.Print(model.Array)
}
