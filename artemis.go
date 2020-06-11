package main

import (
	"sync"

	"github.com/linushung/artemis/cmd/server"
	"github.com/linushung/artemis/cmd/server/rest"
	"github.com/linushung/artemis/internal/app/authorization"
	"github.com/linushung/artemis/internal/pkg/configs"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	log "github.com/sirupsen/logrus"
	//"go.elastic.co/apm/module/apmlogrus"
)

func initLogrus() {
	logstashHost := configs.GetConfigStr("connection.logstash.host")

	if logstashHost == "" {
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{
			// DisableColors: true,
			TimestampFormat: "2006-01-02 15:04:05.000",
			FullTimestamp:   true,
		})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		hook, err := logrustash.NewHookWithFields("udp", logstashHost, "artemis", log.Fields{})
		if err != nil {
			log.Fatal(err)
		}

		log.AddHook(hook)
		//log.AddHook(&apmlogrus.Hook{})
	}
}

func initService() {
	var wg sync.WaitGroup
	authorization.InitJWTService()
	server.InitCircuitBreakerMgr()
	baseServer := server.NewBaseServer()

	wg.Add(1)
	go rest.InitRestServer(*baseServer)

	wg.Wait()
}

func main() {
	log.Infof("***** [INIT:ARTEMIS] ***** Start to launch Artemis 🤓 ...")
	configs.InitConfig()
	initLogrus()
	initService()
}
