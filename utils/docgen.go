package utils

import (
	"log"
	"os"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/docgen"
)

func DocGen(r chi.Router) {
	// define doc config
	doc := docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
		ProjectPath: "github.com/ccau1/test-blockchain-client",
		Intro:       "generated docs.",
	})
	// generate file
	file, err := os.Create("docs/docgen.md")
	if err != nil {
		log.Printf("err: %v", err)
	}
	defer file.Close()
	// write doc to file
	file.Write(([]byte)(doc))
	log.Println("docgen.md updated")
}
