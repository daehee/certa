package main

import "regexp"

type Config struct {
    Match *regexp.Regexp
    Slack string
    Debug bool
}