package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	static := http.FileServer(http.Dir("./webapp/dist"))
	router := httprouter.New()
	router.Handler("GET", "/webapp/*path", http.StripPrefix("/webapp/", static))
	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		http.Redirect(w, req, "/webapp", http.StatusFound)
	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err.Error())
	}
}

// func main() {
// 	rand.Seed(time.Now().Unix())
// 	lib := card.LoadLibrary()
// 	gs := state.NewGame(deck.BuildDeck("assets/decks/test.yaml", lib), deck.BuildDeck("assets/decks/test2.yaml", lib))

// 	g := &machine.GameMachine{
// 		Current: gs,
// 	}
// 	g.RegisterHandler(constant.Player1, &HumanConsole{
// 		player:  constant.Player1,
// 		color:   "\033[31m",
// 		machine: g,
// 	})
// 	g.RegisterHandler(constant.Player2, &HumanConsole{
// 		player:  constant.Player2,
// 		color:   "\033[34m",
// 		machine: g,
// 	})
// 	g.Start(gs)
// }
