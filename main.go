package main

import (
    "flag"
    "regexp"

    "github.com/CaliDog/certstream-go"
    "github.com/jmoiron/jsonq"
    "go.uber.org/zap"
)

var (
    d = flag.String("d", "","regex pattern for domain matching")
)

var (
    logger, _ = zap.NewProduction()
    sugar = logger.Sugar()
)

func main() {
    defer logger.Sync()

    flag.Parse()

    if *d == "" {
        sugar.Fatal("domain regex pattern flag not set")
    }
    drx, err := regexp.Compile(*d)
    if err != nil {
        sugar.Fatal("invalid domain regex pattern")
    }

    stream, errStream := certstream.CertStreamEventStream(false) // heartbeats prevent disconnect

    for {
        select {
        case jq := <-stream:
            go drink(drx, jq)
        case err := <-errStream:
            sugar.Error(err)
        }
    }
}

func drink(drx *regexp.Regexp, jq jsonq.JsonQuery) {
    messageType, err := jq.String("message_type")
    if err != nil{
        sugar.Fatal("Error decoding jq string", err)
    }
    if messageType != "certificate_update" {
        sugar.Info("Not certificate_update:", jq)
        return
    }

    go func() {
        domains, err := jq.ArrayOfStrings("data", "leaf_cert", "all_domains")
        if err != nil {
            sugar.Fatal("Error decoding domains")
        }
        for _, d := range domains {
            if !drx.MatchString(d) {
                continue
            }
            sugar.Info(d)
        }

    }()
}
