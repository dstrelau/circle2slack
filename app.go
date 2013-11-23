package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var slackUrl string = os.ExpandEnv("https://$SLACK_ORGANIZATION.slack.com/services/hooks/incoming-webhook?token=$SLACK_TOKEN")

type SlackPayload struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

func WebhookText(circlePost []byte) string {
	var jsonData interface{}
	json.Unmarshal(circlePost, &jsonData)
	payloadI := jsonData.(map[string]interface{})["payload"]
	payload := payloadI.(map[string]interface{})

	vcsUrl := payload["vcs_url"].(string)
	splits := strings.Split(vcsUrl, "/")
	repo := splits[len(splits)-2]

	buildUrl    := payload["build_url"].(string)
	branch      := payload["branch"].(string)
	vcsRevision := payload["vcs_revision"].(string)
	status      := payload["status"].(string)
	subject     := payload["subject"].(string)
	authorName  := payload["author_name"].(string)

	return fmt.Sprintf("[%s/%s] (<%s|%s>) <%s/commit/%s|%s>: %s - %s",
		repo, branch,
		buildUrl, strings.ToUpper(status),
		vcsUrl, vcsRevision, vcsRevision[0:12],
		subject, authorName)
}

func HandleBuild(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	payload := SlackPayload{
		os.Getenv("SLACK_CHANNEL"),
		os.Getenv("SLACK_BOTNAME"),
		WebhookText(body),
	}

	payloadEnc, err := json.Marshal(payload)
	payloadReader := bytes.NewReader(payloadEnc)

	slackResp, err := http.Post(slackUrl, "application/json", payloadReader)
	defer slackResp.Body.Close()
	if err != nil {
		panic(err)
	}

	slackBody, err := ioutil.ReadAll(slackResp.Body)
	fmt.Println(slackResp)
	fmt.Println(string(slackBody))
	fmt.Fprintln(w, "Thanks!")
}

func main() {
	if os.Getenv("SLACK_ORGANIZATION") == "" ||
	os.Getenv("SLACK_CHANNEL") == "" ||
	os.Getenv("SLACK_BOTNAME") == "" ||
	os.Getenv("SLACK_TOKEN") == "" {
		panic("SLACK_* variables not found")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	http.HandleFunc("/build", HandleBuild)

	fmt.Println(slackUrl)
	log.Fatal(http.ListenAndServe(os.ExpandEnv(":$PORT"), nil))
}
