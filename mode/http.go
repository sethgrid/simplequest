package mode

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/sethgrid/simplequest/quester"
	"github.com/sethgrid/simplequest/quests/rickroll"
	"github.com/sethgrid/simplequest/twilio"
	"github.com/sethgrid/simplequest/utils"
)

type server struct {
	mu         sync.Mutex
	games      map[string]*quester.Quest
	startTime  time.Time
	totalGames int64
}

var srv *server

func init() {
	srv = &server{
		games:     make(map[string]*quester.Quest),
		startTime: time.Now(),
	}

	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			// if we wanted to the be real, we would not stop the world to poll
			srv.mu.Lock()
			for phoneNumber, quest := range srv.games {
				if quest.IsExpired() {
					log.Println("expiring session for ", phoneNumber)
					quest.Stop()
					delete(srv.games, phoneNumber)
				}
			}
			srv.mu.Unlock()
		}
	}()
}

// ValidGameMode validates the game mode as http, sms, or cmd
func ValidGameMode(mode string) bool {
	if mode == "http" || mode == "cmd" {
		return true
	}
	return false
}

func RunCMDServer() {
	player := quester.NewPlayer("local")
	game := rickroll.NewRickRoll(player)

	go game.Start()
	fmt.Println("press enter to start text quest")
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _ := reader.ReadString('\n')
		response := game.TakeCommand(line)
		if response == "you have existed text quest" {
			break
		}
		fmt.Printf(response + "> ")
	}
	fmt.Println("thanks for playing!")
}

// RunHTTPServer is a blocking call that starts a new http server
func RunHTTPServer(port int) {
	mux := http.NewServeMux()

	mux.HandleFunc("/sms", smsHandler)
	mux.HandleFunc("/plain", nonSMSHandler)
	mux.HandleFunc("/metrics", metricsHandler)

	log.Printf("running on :%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatal(err)
	}
}

func nonSMSHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: DRY it up with sms handler
	sms := twilio.SMSParseForm(r)
	utils.Debugf("new http command received: from %s - body %s", sms.From, sms.Body)

	if sms.From == "" || sms.Body == "" {
		utils.Debugf("'Form' and 'Body' required. Try curl -s -X POST $URL:$PORT/plain -d From $NAME -d Body $STRING_COMMAND")
		return
	}

	srv.mu.Lock()
	game, ok := srv.games[sms.From]
	if !ok {
		utils.Debugf("creating new player and quest")
		player := quester.NewPlayer(sms.From)
		game = rickroll.NewRickRoll(player)
		srv.games[sms.From] = game
		srv.totalGames++
		go game.Start()
	}
	srv.mu.Unlock()

	utils.Debugf("issuing command...")
	prompt := game.TakeCommand(sms.Body)

	if prompt == "you have existed text quest" {
		srv.mu.Lock()
		delete(srv.games, sms.From)
		srv.mu.Unlock()
	}

	w.Write([]byte(prompt))
}

// TODO: DRY it up with nonSMS handler
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
		srv.totalGames++
		go game.Start()
	}
	srv.mu.Unlock()

	utils.Debugf("issuing command...")
	prompt := game.TakeCommand(sms.Body)

	if prompt == "you have existed text quest" {
		srv.mu.Lock()
		delete(srv.games, sms.From)
		srv.mu.Unlock()
	}

	w.Write([]byte(twilio.SimpleTwiML(prompt)))
}

type metrics struct {
	ActiveGames int
	TotalGames  int64
	Uptime      string
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	activeGames := len(srv.games)
	uptime := time.Now().Sub(srv.startTime).String()
	totalGames := srv.totalGames

	m := metrics{
		ActiveGames: activeGames,
		TotalGames:  totalGames,
		Uptime:      uptime,
	}

	b, err := json.Marshal(m)
	if err != nil {
		log.Println("unable to marshal metrics ", err.Error())
	}
	w.Write(b)
}
