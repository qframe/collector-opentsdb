package main

import (
	"log"
	"fmt"
	"time"
	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
	"github.com/docker/docker/api/types"
	"github.com/qframe/collector-opentsdb/lib"
)


func Run(qChan qtypes.QChan, cfg *config.Config, name string) {
	p, _ := collector_opentsdb.New(qChan, cfg, name)
	p.Run()
}


func main() {
	qChan := qtypes.NewQChan()
	qChan.Broadcast()
	cfgMap := map[string]string{}

	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)

	p, err := collector_opentsdb.New(qChan, cfg, "opentsdb")
	if err != nil {
		log.Printf("[EE] Failed to create collector: %v", err)
		return
	}
	go p.Run()
	time.Sleep(2*time.Second)
	bg := qChan.Data.Join()
	done := false
	for {
		select {
		case val := <- bg.Read:
			switch val.(type) {
			case qtypes.QMsg:
				qm := val.(qtypes.QMsg)
				if qm.Source == "tcp" {
					switch qm.Data.(type) {
					case types.ContainerJSON:
						cnt := qm.Data.(types.ContainerJSON)
						p.Log("info", fmt.Sprintf("Got inventory response for msg: '%s'", qm.Msg))
						p.Log("info", fmt.Sprintf("        Container{Name:%s, Image: %s}", cnt.Name, cnt.Image))
						done = true

					}
				}
			}
		}
		if done {
			break
		}
	}
}
