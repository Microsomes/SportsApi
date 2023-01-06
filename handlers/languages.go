package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
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
	S3Link    string `json:"s3link"`
	StartPage int    `json:"startPage"`
	PageCount int    `json:"pageCount"`
}

func getAllUnits() []*Unit {
	// allUnits := Units{}
	alllangs := []*Unit{}
	var i = 1
	for {

		if i > 100 {
			break
		}

		b, err := os.Open(fmt.Sprintf("./handlers/pimfiles/pashto/%d.json", i))
		if err != nil {
			panic(err)
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
			// var res = []Reading{}

			for _, ll := range units[0].Lessons {
				les = append(les, Lesson{
					Name:      ll.Name,
					AduioLink: ll.AduioLink,
					S3Audio:   units[0].Name + "/" + ll.Name + ".mp3",
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
			Readings:     ll.Readings,
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

func LimitSectionOne(units map[string][]Unit) map[string][]Unit {

	var units2 = make(map[string][]Unit)

	for tri, l := range units {

		units := []Unit{}

		for _, uni := range l {
			units = append(units, Unit{
				Name:         uni.Name,
				Language:     uni.Language,
				TotalLessons: uni.TotalLessons,
			})
		}

		units2[tri] = units
	}

	return units2
}

func AllLanguages(w http.ResponseWriter, r *http.Request) {

	re := regexp.MustCompile("\\?")

	lt := re.MatchString(r.URL.String())

	units := getAllUnits()

	units = DedublicateUnits(units)

	var units2 = sortLanguagesToUnits(units)

	if lt {
		units2 = LimitSectionOne(units2)
	}

	b, _ := json.Marshal(units2)
	w.Write(b)

}

func OneLanguage(w http.ResponseWriter, r *http.Request) {

	uri, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	params, err := url.ParseQuery(uri.RawQuery)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if params.Get("lang") == "" {
		http.Error(w, "Param ?lang= missing", http.StatusBadRequest)
		return
	}

	units := getAllUnits()

	units = DedublicateUnits(units)

	var units2 = sortLanguagesToUnits(units)

	var toUnits = []Unit{}

	for _, un := range units2 {
		for _, unit := range un {
			if unit.Language == params.Get("lang") {
				toUnits = append(toUnits, unit)
			}
		}
	}

	b, _ := json.Marshal(toUnits)

	w.Write(b)
}
