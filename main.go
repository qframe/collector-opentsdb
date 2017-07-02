package main

import (
	"log"
	"time"
	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
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
			case qtypes.OpenTSDBMetric:
				otm := val.(qtypes.OpenTSDBMetric)
				log.Println(otm.String())
			}
		}
		if done {
			time.Sleep(time.Second)
			break
		}
	}
}
