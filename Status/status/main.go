package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rumblefrog/go-a2s"
)

var statusFile = "server_status.json"

func main() {
	info, players := waitForClientStart()
	writeServerInfo(info, players)

	router := configureRouter()

	go continouslyUpdateStatus()
	router.Run(":7032")
}

func continouslyUpdateStatus() {
	for {
		log.Printf("Updating server status")
		serverConfig := getConfig()
		address := fmt.Sprintf("%v:%v", serverConfig.IP, serverConfig.Port)

		info, players, _ := a2sQuery(address)
		writeServerInfo(info, players)
		time.Sleep(time.Duration(30) * time.Second)
	}
}

func waitForClientStart() (a2sInfo *a2s.ServerInfo, a2sPlayers *a2s.PlayerInfo) {
	serverConfig := getConfig()
	address := fmt.Sprintf("%v:%v", serverConfig.IP, serverConfig.Port)

	for {
		info, players, err := a2sQuery(address)
		if err == nil {
			if info != nil {
				a2sInfo = info
				a2sPlayers = players
				log.Printf("Client started")
				break
			} else {
				fmt.Printf("bad")
			}
		}
	}
	return a2sInfo, a2sPlayers
}

type Player struct {
	Name     string  `json:"Name"`
	Duration float32 `json:"Duration"`
}

type ServerInfo struct {
	LastUpdate        string    `json:"LastUpdate"`
	Name              string    `json:"Name"`
	Application       string    `json:"Application"`
	ServerType        string    `json:"ServerType"`
	OS                string    `json:"OS"`
	PlayerCount       int       `json:"PlayerCount"`
	MaxPlayers        int       `json:"MaxPlayers"`
	Map               string    `json:"Map,omitempty"`
	PasswordProtected bool      `json:"PasswordProtected"`
	VAC               bool      `json:"VAC"`
	IP                string    `json:"IP"`
	Port              int       `json:"Port"`
	Players           []*Player `json:"Players"`
}

func writeServerInfo(info *a2s.ServerInfo, players *a2s.PlayerInfo) {
	config := getConfig()
	port, _ := strconv.Atoi(config.Port)

	var playerInfo []*Player
	for _, player := range players.Players {
		newPlayer := &Player{
			Name:     player.Name,
			Duration: player.Duration,
		}
		playerInfo = append(playerInfo, newPlayer)
	}

	serverInfo := ServerInfo{
		LastUpdate:        time.Now().UTC().String(),
		Name:              info.Name,
		Application:       info.Game,
		ServerType:        info.ServerType.String(),
		OS:                info.ServerOS.String(),
		PlayerCount:       int(info.Players),
		MaxPlayers:        int(info.MaxPlayers),
		Map:               info.Map,
		PasswordProtected: info.Visibility,
		VAC:               info.VAC,
		IP:                config.IP,
		Port:              port,
		Players:           playerInfo,
	}
	file, _ := json.MarshalIndent(serverInfo, "", " ")
	_ = ioutil.WriteFile(statusFile, file, 0644)
}

func a2sQuery(address string) (info *a2s.ServerInfo, players *a2s.PlayerInfo, err error) {
	client, err := a2s.NewClient(address)
	if err == nil {
		defer client.Close()
		info, _ = client.QueryInfo()
		players, _ = client.QueryPlayer()
		client.Close()
		return info, players, nil
	}
	return info, players, err
}

type Config struct {
	IP          string
	Port        string
	QueryMethod string
}

func getConfig() *Config {
	env := godotenv.Load(".env")
	if env == nil {
		log.Fatalf("Error loading .env file")
	}

	return &Config{
		"sea-1.us.uncletopia.com",
		os.Getenv("PORT"),
		os.Getenv("QUERY_METHOD"),
	}
}

func getIP() string {
	resp, err := http.Get("http://checkip.amazonaws.com")
	if err == nil {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return strings.TrimSpace(string(b))
	} else {
		return ""
	}
}
