package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Language struct {
	Name         string `json:"name"`
	UnitAmount   int    `json:"units"`
	TotalLessons int    `json:"total_lessons"`
}

type Unit struct {
	Name         string   `json:"courseName"`
	Language     string   `json:"languageName"`
	TotalLessons int      `json:"total_lessons"`
	Lessons      []Lesson `json:"lessons"`
}

type Units []Unit

type Lesson struct {
	Name      string `json:"name"`
	AduioLink string `json:"audioLink"`
}

type Languages []*Language

func PerformDownload(dst *os.File, link string) {
	defer dst.Close()
	dst.Write([]byte("audio file"))
}

func DownloadLanguage(units []*Unit) {

	for _, u := range units {
		err := os.Mkdir(u.Name, 0777)
		if err != nil {
			fmt.Println("folder does not exist")
		}

		for _, l := range u.Lessons {

			fname := u.Name + "/" + l.Name
			fmt.Println(fname)
			dst, err := os.Create(fname)
			if err != nil {
				panic(err)
			}

			go PerformDownload(dst, l.AduioLink)

		}

	}

}

func getAllUnits() []*Unit {
	// allUnits := Units{}
	alllangs := []*Unit{}
	var i = 1
	for {

		if i > 100 {
			break
		}

		b, err := os.Open(fmt.Sprintf("../../pimfiles/pashto/%d.json", i))
		if err != nil {
			break
		}
		defer b.Close()

		//get size

		info, _ := b.Stat()

		buf := make([]byte, info.Size())

		b.Read(buf)

		units := Units{}

		json.Unmarshal(buf, &units)

		if len(units) >= 1 {
			alllangs = append(alllangs, &Unit{
				Name:         units[0].Name,
				Language:     units[0].Language,
				Lessons:      units[0].Lessons,
				TotalLessons: len(units[0].Lessons),
			})
		}
		i++

	}

	return alllangs

}

func main() {
	units := getAllUnits()

	DownloadLanguage(units)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(units)
		w.Write(b)
	})

	http.ListenAndServe(":5003", nil)
}
