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
	Readings     Reading  `json:"readings"`
}

type Units []Unit

type Lesson struct {
	Name      string `json:"name"`
	AduioLink string `json:"audioLink"`
	S3Audio   string
	Image     LessonImage `json:"image"`
}

type LessonImage struct {
	FullImage  string `json:"fullImageAddress"`
	ThumbImage string `json:"thumbImageAddress"`
}

type Reading struct {
	Pdf     string  `json:"pdf"`
	PdfName string  `json:"pdfName"`
	Audios  []Audio `json:"audios"`
}

type Audio struct {
	Title     string `json:"title"`
	AudioLink string `json:"audioLink"`
	StartPage int    `json:"startPage"`
	PageCount int    `json:"pageCount"`
}

func PerformLessonDownload(dst *os.File, link string, guard chan struct{}) {
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

func PerformReadingDownload(guard chan struct{}, audio string, out *os.File) {

	defer out.Close()
	r, err := http.Get(audio)
	if err != nil {
		fmt.Println(err)
	}

	defer r.Body.Close()

	io.Copy(out, r.Body)

	<-guard
}

func DownloadLanguage(units []*Unit, concurrentLevel int) {

	guard := make(chan struct{}, concurrentLevel)

	for _, u := range units {
		err := os.Mkdir(u.Name, 0777)
		if err != nil {
		}

		os.Mkdir(u.Name+"/readings", 0777)

		for _, l := range u.Readings.Audios {

			guard <- struct{}{}

			readingOs, _ := os.Create(u.Name + "/readings/" + l.Title + ".mp3")

			fmt.Println(":", u.Name)

			go PerformReadingDownload(guard, l.AudioLink, readingOs)

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

			go PerformLessonDownload(dst, l.AduioLink, guard)

		}

	}

	close(guard)
	os.Exit(1)

}

func getAllUnits() []*Unit {
	// allUnits := Units{}
	alllangs := []*Unit{}
	var i = 1
	for {

		if i > 100 {
			break
		}

		b, err := os.Open(fmt.Sprintf("../../handlers/pimfiles/pashto/%d.json", i))
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
					AduioLink: ll.AduioLink,
					S3Audio:   units[0].Name + "/" + ll.Name + ".mp3",
					Image:     ll.Image,
				})
			}

			alllangs = append(alllangs, &Unit{
				Name:         units[0].Name,
				Language:     units[0].Language,
				Lessons:      les,
				TotalLessons: len(units[0].Lessons),
				Readings:     units[0].Readings,
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

	// var units2 = sortLanguagesToUnits(units)

	DownloadLanguage(units, 1)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(units)
		w.Write(b)
	})

	err := http.ListenAndServe(":5004", nil)
	if err != nil {
		panic(err)
	}
}
