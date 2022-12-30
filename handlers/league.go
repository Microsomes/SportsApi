package handlers

import (
	"encoding/json"
	"net/http"
)

type League struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type Leagues []League

func AllLeagues() *Leagues {
	leagues := Leagues{}

	leagues = append(leagues, League{Name: "Premier League", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "Women's Football", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "Championship", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "League One", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "League Two", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "National League", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "FA Cup", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "EFL Cup", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "Scottish Football", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "Scottish Premiership", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "Scottish Championship", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "Scottish League One", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "Scottish League Two", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "Scottish Cup", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "Scottish League Cup", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "Scottish Challenge Cup", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "Welsh Football", Country: "UK & Ireland"})
	leagues = append(leagues, League{Name: "Irish Football", Country: "UK & Ireland"})

	leagues = append(leagues, League{Name: "European Football", Country: "Europe"})
	leagues = append(leagues, League{Name: "Champions League", Country: "Europe"})
	leagues = append(leagues, League{Name: "Europa League", Country: "Europe"})

	leagues = append(leagues, League{Name: "FIFA Men's World Cup", Country: "International"})
	leagues = append(leagues, League{Name: "UEFA Men's Euro", Country: "International"})
	leagues = append(leagues, League{Name: "African Football", Country: "International"})
	leagues = append(leagues, League{Name: "FIFA Women's World Cup", Country: "International"})
	leagues = append(leagues, League{Name: "UEFA Women's Euro", Country: "International"})

	return &leagues
}

func (*Leagues) ServeHTTP(w http.ResponseWriter, rr *http.Request) {
	leagues := AllLeagues()
	b, _ := json.Marshal(leagues)
	w.Write(b)
}
