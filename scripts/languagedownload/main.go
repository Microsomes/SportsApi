package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

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
	S3Audio   string
}

func PerformDownload(dst *os.File, link string, guard chan struct{}) {
	defer dst.Close()

	resp, err := http.Get(link)

	if err != nil {
		fmt.Println("cannot download")
	}

	defer resp.Body.Close()

	lt, err := io.Copy(dst, resp.Body)

	if err != nil {
		panic(err)
	}

	fmt.Println("downloaded:", lt)

	fmt.Println(dst.Name())

	<-guard

}

func DownloadLanguage(units []*Unit, concurrentLevel int) {

	guard := make(chan struct{}, concurrentLevel)

	for _, u := range units {
		err := os.Mkdir(u.Name, 0777)
		if err != nil {
		}

		for _, l := range u.Lessons {

			fname := u.Name + "/" + l.Name
			dst, err := os.Create(fname + ".mp3")
			if err != nil {
				panic(err)
			}
			fmt.Println(u.Name)
			fmt.Println("performing download:", l.Name)

			guard <- struct{}{}
			fmt.Println("perform download")

			go PerformDownload(dst, l.AduioLink, guard)

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

			var les = []Lesson{}

			for _, ll := range units[0].Lessons {
				les = append(les, Lesson{
					Name:      ll.Name,
					AduioLink: "--",
					S3Audio:   units[0].Name + "/" + ll.Name + ".mp3",
				})
			}

			alllangs = append(alllangs, &Unit{
				Name:         units[0].Name,
				Language:     units[0].Language,
				Lessons:      les,
				TotalLessons: len(units[0].Lessons),
			})
		}
		i++

	}

	return alllangs

}

// type Language struct {
// 	Name  string
// 	Units []Unit
// }

func sortLanguagesToUnits(units []*Unit) map[string][]Unit {

	var langmap = make(map[string][]Unit)

	for _, ll := range units {

		langmap[ll.Language] = append(langmap[ll.Language], Unit{
			Name:         ll.Name,
			Language:     ll.Language,
			TotalLessons: ll.TotalLessons,
			Lessons:      ll.Lessons,
		})

	}

	return langmap
}

func main() {
	units := getAllUnits()

	var units2 = sortLanguagesToUnits(units)

	DownloadLanguage(units, 10)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(units2)
		w.Write(b)
	})

	err := http.ListenAndServe(":5004", nil)
	if err != nil {
		panic(err)
	}
}
