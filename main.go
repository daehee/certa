package main

import (
	"errors"
	"flag"
	"os"
	"regexp"
	"strings"

	"github.com/daehee/certstream-ws"
	"go.uber.org/zap"
)

var (
	matchPattern = flag.String("d", "", "regex pattern for domain matching")
	debug        = flag.Bool("v", false, "verbose debug mode")
)

var (
	logger, _ = zap.NewDevelopment()
	sugar     = logger.Sugar()
)

func main() {
	defer logger.Sync()
	flag.Parse()

	drx, err := checkDomainRegex()
	if err != nil {
		sugar.Fatal(err)
	}

	slackWebhook := os.Getenv("SLACK_WEBHOOK_URL")
	// if slackWebhook == "" {
	//     sugar.Fatal("slack webhook url not set")
	// }

	config := &Config{
		Match: drx,
		Slack: slackWebhook,
		Debug: *debug,
	}

	stream, errStream := certstream.CertStreamEventStream(config.Debug)

	for {
		select {
		case v := <-stream:
			dvs := v.GetArray("data", "leaf_cert", "all_domains")
			if len(dvs) == 0 {
				break
			}
			for _, dv := range dvs {
				domain := string(dv.GetStringBytes())
				domain = cleanDomain(domain)
				if !drx.MatchString(domain) {
					continue
				}

				sugar.Infof("matched domain: %s", domain)
				go AddDomain(domain)
				// go sendSlack(config.Slack, domain)
			}
		case err := <-errStream:
			sugar.Error(err)
		}
	}
}

func checkDomainRegex() (*regexp.Regexp, error) {
	if *matchPattern == "" {
		return nil, errors.New("domain regex pattern flag not set")
	}
	drx, err := regexp.Compile(*matchPattern)
	if err != nil {
		return nil, errors.New("invalid domain regex pattern")
	}
	return drx, nil
}

// cleanDomain strips DNS entry characters *, %, \. from beginning of domain
// credit: tomnomnom/assetfinder
// https://github.com/tomnomnom/assetfinder/blob/4e95d8701aae8cff1c27af2626eb22ba110ad583/main.go#L116
func cleanDomain(d string) string {
	d = strings.ToLower(d)

	if len(d) < 2 {
		return d
	}
	if d[0] == '*' || d[0] == '%' {
		d = d[1:]
	}
	if d[0] == '.' {
		d = d[1:]
	}

	return d
}
