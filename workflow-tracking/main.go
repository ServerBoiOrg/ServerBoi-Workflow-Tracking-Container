package main

import (
	"bytes"
	dc "discordhttpclient"
	"encoding/json"
	"fmt"
	gu "generalutils"
	"log"
	"net/http"
	ru "responseutils"
	"time"

	dt "github.com/awlsring/discordtypes"
)

var (
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
	advanceWorkflow(&AdvanceWorkflowInput{
		ExecutionName: executionName,
	})
	address := "http://service-monitor:7032/status"
	log.Printf("Updating workflow embed")
	updateEmbed(formWaitingEmbed())
	for {
		log.Printf("Checking if local server is up on %v", address)
		resp, err := http.Get(address)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				status := StatusResponse{}
				json.NewDecoder(resp.Body).Decode(&status)
				log.Printf("Status response: %v", status.Running)
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
		} else {
			log.Printf("Error talking to local server %v", err)
		}
		log.Printf("Server is offline. Updating embed")
		updateEmbed(formWaitingEmbed())
		time.Sleep(time.Duration(1) * time.Minute)
	}

}

func advanceWorkflow(input *AdvanceWorkflowInput) {
	body, _ := json.Marshal(input)
	http.Post("https://api.serverboi.io/bootstrap", "application/json", bytes.NewBuffer(body))
}

func formWaitingEmbed() *dt.Embed {
	return embedTemplate("Waiting for application to start")
}

func formCompleteEmbed() *dt.Embed {
	return embedTemplate("Application started")
}

func embedTemplate(stage string) *dt.Embed {
	return ru.CreateWorkflowEmbed(&ru.CreateWorkflowEmbedInput{
		Name:        "Provision-Server",
		Description: fmt.Sprintf("WorkflowID: %s", executionName),
		Status:      "🟢 Running",
		Stage:       stage,
		Color:       ru.DiscordGreen,
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
