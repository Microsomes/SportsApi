package handlers

import (
	"encoding/json"
	"net/http"
)

type HttpDoc struct {
	Path        string `json:"path"`
	Description string `json:"description"`
}

type Docs []*HttpDoc

func AllDocs() *Docs {

	docs := Docs{}

	docs = append(docs, &HttpDoc{
		Path:        "/leagues",
		Description: "returns a list of supported leagues",
	})

	docs = append(docs, &HttpDoc{
		Path:        "/languages",
		Description: "returns a list of supported languages for talkfs.com",
	})

	return &docs
}

func (*Docs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	docs := AllDocs()
	b, _ := json.Marshal(docs)
	w.Write(b)
}
