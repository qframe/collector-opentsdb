package main

import (
	"log"
	"time"
	"github.com/zpatrick/go-config"
	"github.com/qframe/collector-opentsdb"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/metrics"
)

func Run(qChan qtypes_qchannel.QChan, cfg *config.Config, name string) {
	p, _ := qcollector_opentsdb.New(qChan, cfg, name)
	p.Run()
}

func main() {
	qChan := qtypes_qchannel.NewQChan()
	qChan.Broadcast()
	cfgMap := map[string]string{}

	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)

	p, err := qcollector_opentsdb.New(qChan, cfg, "opentsdb")
	if err != nil {
		log.Printf("[EE] Failed to create collector: %v", err)
		return
	}
	go p.Run()
	time.Sleep(2 * time.Second)
	bg := qChan.Data.Join()
	done := false
	for {
		select {
		case val := <-bg.Read:
			switch val.(type) {
			case qtypes_metrics.OpenTSDBMetric:
				otm := val.(qtypes_metrics.OpenTSDBMetric)
				log.Println(otm.String())
			}
		}
		if done {
			time.Sleep(time.Second)
			break
		}
	}
}
