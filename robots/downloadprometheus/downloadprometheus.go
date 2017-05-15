package downloadprometheus

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/trinchan/slackbot/robots"
)

type bot struct{}

// Registers the bot with the server for command /test.
func init() {
	r := &bot{}
	robots.RegisterRobot("download-prometheus", r)
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (r bot) Run(p *robots.Payload) string {
	// If you (optionally) want to do some asynchronous work (like sending API calls to slack)
	// you can put it in a go routine like this
	go r.DeferredAction(p)
	// The string returned here will be shown only to the user who executed the command
	// and will show up as a message from slackbot.
	return "Baby dingos are building a URL just for you for the latest Dingo Prometheus tile."
}

func (r bot) DeferredAction(p *robots.Payload) {
	// Let's use the IncomingWebhook struct defined in payload.go to form and send an
	// IncomingWebhook message to slack that can be seen by everyone in the room. You can
	// read the Slack API Docs (https://api.slack.com/) to know which fields are required, etc.
	// You can also see what data is available from the command structure in definitions.go
	// Alternatively, you can make a SlashCommandResponse, with the same fields, and call
	// reponse.Send(p)

	productName := "Dingo Prometheus"
	productLabel := "dingo-prometheus"
	version, err := r.lookupLatestTileVersion(productLabel)
	if err != nil {
		r.sendErrorResponse(p, err)
		return
	}
	url, err := r.generateTempURL(productLabel, version)
	if err != nil {
		r.sendErrorResponse(p, err)
		return
	}

	response := &robots.SlashCommandResponse{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "dingobot",
		Text:        fmt.Sprintf("Download %s v%s tile at %s", productName, version, url),
		IconEmoji:   ":dingo:",
		UnfurlLinks: false,
		Parse:       robots.ParseStyleFull,
	}
	if err := response.Send(p); err != nil {
		r.sendErrorResponse(p, err)
	}

	webbookResponse := &robots.IncomingWebhook{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "dingobot",
		Text:        fmt.Sprintf("Another happy %s v%s tile sent on its way to a new home! (via `/download-prometheus`)", productName, version),
		IconEmoji:   ":dingo:",
		UnfurlLinks: false,
		Parse:       robots.ParseStyleFull,
	}
	if err := webbookResponse.Send(); err != nil {
		r.sendErrorResponse(p, err)
	}

	salesResponse := &robots.IncomingWebhook{
		Domain:      p.TeamDomain,
		Channel:     "G0N9JP199", // #sales-announcements channel ID
		Username:    "dingobot",
		Text:        fmt.Sprintf("Product %s v%s was requested by @%s in channel @%s", productName, version, p.UserName, p.ChannelName),
		IconEmoji:   ":dingo:",
		UnfurlLinks: false,
		Parse:       robots.ParseStyleFull,
	}
	if err := salesResponse.Send(); err != nil {
		r.sendErrorResponse(p, err)
	}
}

func (r bot) Description() string {
	// In addition to a Run method, each Robot must implement a Description method which
	// is just a simple string describing what the Robot does. This is used in the included
	// /c command which gives users a list of commands and descriptions
	return "Fetch URL to download latest Dingo Prometheus tile."
}

func (r bot) awsBucket(productLabel string) (bucket *s3.Bucket, err error) {
	auth, err := aws.EnvAuth()
	if err != nil {
		return
	}
	client := s3.New(auth, aws.APSoutheast)
	bucket = client.Bucket(fmt.Sprintf("%s-public-pivotaltile", productLabel))
	return
}

func (r bot) lookupLatestTileVersion(productLabel string) (productVersion string, err error) {
	bucket, err := r.awsBucket(productLabel)
	if err != nil {
		return
	}
	log.Println("Getting bucket contents...")
	listResp, err := bucket.List(productLabel, "/", "", 1000)
	if err != nil {
		return
	}

	latestVersion, _ := version.NewVersion("0.0.0")

	versionRegexp, err := regexp.Compile(fmt.Sprintf("%s-(.*)\\.pivotal", productLabel))
	if err != nil {
		return
	}
	for _, key := range listResp.Contents {
		filename := key.Key
		match := versionRegexp.FindStringSubmatch(filename)
		if match != nil {
			fileVersion, err := version.NewVersion(match[1])
			if err == nil && fileVersion.Prerelease() == "" && latestVersion.LessThan(fileVersion) {
				latestVersion = fileVersion
			}
		}
	}
	if latestVersion.String() == "0.0.0" {
		return "", fmt.Errorf("No published releases yet for %s", productLabel)
	}
	return latestVersion.String(), nil
}

func (r bot) generateTempURL(productLabel, version string) (url string, err error) {
	bucket, err := r.awsBucket(productLabel)
	if err != nil {
		return
	}
	expiryTime := time.Now().Add(5 * time.Minute)
	url = bucket.SignedURL(fmt.Sprintf("%s-%s.pivotal", productLabel, version), expiryTime)
	return
}

func (r bot) sendErrorResponse(p *robots.Payload, err error) {
	log.Printf("ERROR in %s/%s: %s", p.TeamDomain, p.ChannelName, err.Error())
	response := &robots.SlashCommandResponse{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "dingobot",
		Text:        err.Error(),
		IconEmoji:   ":dingo:",
		UnfurlLinks: false,
		Parse:       robots.ParseStyleFull,
	}
	response.Send(p)
}
