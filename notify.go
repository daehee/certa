package main

import (
    "fmt"

    "github.com/ashwanthkumar/slack-go-webhook"
)

func sendSlack(webhook string, domain string) {
    p := slack.Payload {
        Text: fmt.Sprintf("[certa] new domain: %s", domain),
    }
    for _, err := range slack.Send(webhook, "", p) {
        sugar.Error(err)
    }
}
