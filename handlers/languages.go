package handlers

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

func GetAllUnits() []*Unit {
	// allUnits := Units{}
	alllangs := []*Unit{}
	var i = 1
	for {

		if i > 100 {
			break
		}

		b, err := os.Open(fmt.Sprintf("./handlers/pimfiles/pashto/%d.json", i))
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

func SortLanguagesToUnits(units []*Unit) map[string][]Unit {

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

type Section struct {
	Name     string `json:"name"`
	Feedid   string `json:"feedid"`
	FeedName string `json:"feedName"`
}
type LanguageSect struct {
	Name     string    `json:"name"`
	Sections []Section `json:"sections"`
}

func SortLikeTalkFS(all map[string][]Unit) []LanguageSect {

	sects := []LanguageSect{}

	for k, v := range all {

		var sects2 = []Section{}

		for _, u := range v {
			sects2 = append(sects2, Section{
				Name:     u.Name,
				Feedid:   u.Name,
				FeedName: u.Name,
			})
		}

		sects = append(sects, LanguageSect{
			Name:     k,
			Sections: sects2,
		})

	}
	return sects
}

func DedublicateUnits(units []*Unit) []*Unit {

	var newUnits = []*Unit{}

	UnitsSeen := make(map[string]bool)

	for _, x := range units {
		if UnitsSeen[x.Name] == false {
			newUnits = append(newUnits, x)
			UnitsSeen[x.Name] = true
		}
	}

	return newUnits

}

func AllLanguages(w http.ResponseWriter, r *http.Request) {

	units := GetAllUnits()

	units = DedublicateUnits(units)

	var units2 = SortLanguagesToUnits(units)

	b, _ := json.Marshal(units2)
	w.Write(b)

}
