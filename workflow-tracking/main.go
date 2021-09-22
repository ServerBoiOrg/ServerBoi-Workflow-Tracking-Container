package main

import (
	"bytes"
	dc "discordhttpclient"
	"encoding/json"
	"fmt"
	gu "generalutils"
	"io"
	"log"
	"net/http"
	ru "responseutils"
	"strings"
	"time"

	dt "github.com/awlsring/discordtypes"
)

var (
	ip               = getIP()
	applicationID    = gu.GetEnvVar("APPLICATION_ID")
	interactionToken = gu.GetEnvVar("INTERACTION_TOKEN")
	executionName    = gu.GetEnvVar("EXECUTION_NAME")
	client           = dc.CreateClient(&dc.CreateClientInput{
		ApiVersion: "v9",
	})
)

type StatusResponse struct {
	Running bool `json:"Running"`
}

type UpdateEmbedInput struct {
	Embed  *dt.Embed
	Status string
	Stage  string
	Color  int
}

type AdvanceWorkflowInput struct {
	ExecutionName string `json:"execution_name"`
}

func main() {
	log.Printf("Updating workflow embed")
	updateEmbed(formWaitingEmbed())
	for {
		log.Print("Checking it local server is up.")
		resp, err := http.Get(fmt.Sprintf("%v:7032/status", ip))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				status := StatusResponse{}
				json.NewDecoder(resp.Body).Decode(&status)
				if status.Running {
					log.Printf("Local server is online")
					advanceWorkflow(&AdvanceWorkflowInput{
						ExecutionName: executionName,
					})
					log.Printf("Updating workflow embed")
					updateEmbed(formCompleteEmbed())
					break
				}
			} else {
				log.Printf("Status code not 200, retrying")
			}
		}
		// This may be too frequent
		log.Printf("Server is offline. Updating embed")
		updateEmbed(formWaitingEmbed())
		time.Sleep(time.Duration(30) * time.Second)
	}

}

func advanceWorkflow(input *AdvanceWorkflowInput) {
	body, _ := json.Marshal(input)
	http.Post("https://api.serverboi.io/bootstrap", "application/json", bytes.NewBuffer(body))
}

func formWaitingEmbed() *dt.Embed {
	return ru.CreateWorkflowEmbed(&ru.CreateWorkflowEmbedInput{
		Status: "🟢 Running",
		Stage:  "Waiting for Application Container",
		Color:  ru.DiscordGreen,
	})
}

func formCompleteEmbed() *dt.Embed {
	return ru.CreateWorkflowEmbed(&ru.CreateWorkflowEmbedInput{
		Status: "🟢 Running",
		Stage:  "Application Container Complete",
		Color:  ru.DiscordGreen,
	})
}

func updateEmbed(embed *dt.Embed) {
	for {
		_, headers, err := client.EditInteractionResponse(&dc.InteractionFollowupInput{
			ApplicationID:    applicationID,
			InteractionToken: interactionToken,
			Data: &dt.InteractionCallbackData{
				Embeds: []*dt.Embed{embed},
			},
		})
		if err != nil {
			log.Fatalf("Error editing embed message: %v", err)
		}
		done := dc.StatusCodeHandler(*headers)
		if done {
			break
		}
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
