package main

import (
	"bytes"
	dc "discordhttpclient"
	"encoding/json"
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
	log.Printf("Updating embed")
	updateEmbed(formWaitingEmbed())
	ip := getIP()
	for {
		resp, err := http.Get(fmt.Sprintf("%v:7032/status", ip))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				status := StatusResponse{}
				json.NewDecoder(resp.Body).Decode(&status)
				if status.Running {
					log.Printf("Online")
					advanceWorkflow(&AdvanceWorkflowInput{
						ExecutionName: executionName,
					})
					updateEmbed(formCompleteEmbed())
					break
				}
			}
		}
		log.Printf("Offline. Updating embed")
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
