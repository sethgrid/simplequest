package mode

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/sethgrid/simplequest/quester"
	"github.com/sethgrid/simplequest/quests/rickroll"
	"github.com/sethgrid/simplequest/utils"

	"github.com/sethgrid/simplequest/twilio"
)

type server struct {
	mu    sync.Mutex
	games map[string]*quester.Quest
}

var srv *server

func init() {
	srv = &server{
		games: make(map[string]*quester.Quest),
	}
}

// ValidGameMode validates the game mode as http, sms, or cmd
func ValidGameMode(mode string) bool {
	if mode == "http" || mode == "sms" || mode == "cmd" {
		return true
	}
	return false
}

// RunHTTPServer is a blocking call that starts a new http server
func RunHTTPServer(port int) {
	mux := http.NewServeMux()

	mux.HandleFunc("/sms", smsHandler)
	log.Printf("running on :%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatal(err)
	}
}

func smsHandler(w http.ResponseWriter, r *http.Request) {
	sms := twilio.SMSParseForm(r)
	utils.Debugf("new sms received: from %s - body %s", sms.From, sms.Body)

	if sms.From == "" || sms.Body == "" {
		utils.Debugf("no from or body in sms")
		return
	}

	srv.mu.Lock()
	game, ok := srv.games[sms.From]
	if !ok {
		utils.Debugf("creating new player and quest")
		player := quester.NewPlayer(sms.From)
		game = rickroll.NewRickRoll(player)
		srv.games[sms.From] = game
		go game.Start()
	}
	srv.mu.Unlock()

	utils.Debugf("issuing command...")
	prompt := game.TakeCommand(sms.Body)
	utils.Debugf("got prompt")
	w.Write([]byte(twilio.SimpleTwiML(prompt)))
}
