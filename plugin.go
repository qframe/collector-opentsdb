package qcollector_opentsdb

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
	"encoding/json"
	"github.com/zpatrick/go-config"
	"github.com/albrow/negroni-json-recovery"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/qframe/collector-opentsdb/models"
	"compress/gzip"
	"github.com/pkg/errors"
	"github.com/qframe/types/plugin"
	"github.com/qframe/functions"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/metrics"
)

const (
	version = "0.1.1"
	pluginTyp = "collector"
	pluginPkg = "opentsdb"
)

type Plugin struct {
	*qtypes_plugin.Plugin
}

func Start(qChan qtypes_qchannel.QChan, cfg *config.Config, name string) (err error) {
	p, err := New(qChan, cfg, name)
	if err != nil {
		return errors.Wrap(err, "Failed to create new plugin")
	}
	p.Run()
	return errors.New("Serving of http endpoint finished")

}

func New(qChan qtypes_qchannel.QChan, cfg *config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes_plugin.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
	}
	return p, err
}


func (p *Plugin) Run() {
	qfunctions.Log(p, "notice", fmt.Sprintf("Start collector v%s", p.Version))
	host := p.CfgStringOr("bind-host", "0.0.0.0")
	port := p.CfgStringOr("bind-port", "8070")
	router := mux.NewRouter()
	router.HandleFunc("/api/put", p.putHandler).Methods("POST")

	logger := negroni.NewLogger()
	logger.SetDateFormat(time.RFC3339Nano)
	n := negroni.New(logger)
	n.Use(recovery.JSONRecovery(true))
	n.UseHandler(router)
	n.Run(host+":"+port)
}


func (p *Plugin) putHandler(w http.ResponseWriter, r *http.Request) {
	result := models.NewHttpResponse()
	if r.Header.Get("Content-Encoding") == "gzip" {
		r.Header.Del("Content-Length")
		zr, err := gzip.NewReader(r.Body)
		if err != nil {
			errResponse := fmt.Sprintf(`{"failed": 1,"success": 0, "error": %s}`, err.Error())
			http.Error(w, errResponse, http.StatusBadRequest)
			qfunctions.Log(p, "error", errResponse)
			return
		}
		r.Body = models.GZIPreadCloser{zr, r.Body}
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	var m qtypes_metrics.OpenTSDBMetric
	if err := json.Unmarshal(body, &m); err != nil {
		var ms []qtypes_metrics.OpenTSDBMetric
		if err := json.Unmarshal(body, &ms); err != nil {
			result.Failed += 1
			errResponse := fmt.Sprintf(`{"failed": 1,"success": 0, "error": %s >> %s}`, err.Error(), string(body))
			http.Error(w, errResponse, http.StatusBadRequest)
			qfunctions.Log(p,"error", errResponse)
			return
		} else {
			for _, m := range ms {
				p.QChan.SendData(m)
				qfunctions.Log(p,"trace", fmt.Sprintf("Received: %s", m.String()))
				result.Success += 1
			}
			err := writeResult(w, result)
			if err != nil {
				qfunctions.Log(p,"error", err.Error())

			}
			return
		}
	} else {
		result.Success += 1
		p.QChan.SendData(m)
		qfunctions.Log(p,"trace", fmt.Sprintf("Received: %s", m.String()))
		err := writeResult(w, result)
		if err != nil {
			qfunctions.Log(p,"error", err.Error())

		}
		return
	}
}

func writeResult(w http.ResponseWriter, result models.HttpResponse) (err error) {
	w.Header().Set("Content-Type", "application/json")
	outgoingJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write(outgoingJSON)
	return
}
