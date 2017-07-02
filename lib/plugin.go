package collector_opentsdb

import (
	"fmt"
	"net/http"
	"io/ioutil"

	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
	"github.com/albrow/negroni-json-recovery"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const (
	version = "0.0.0"
	pluginTyp = "collector"
	pluginPkg = "opentsdb"
)

type Plugin struct {
	qtypes.Plugin
}

func New(qChan qtypes.QChan, cfg *config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
	}
	return p, err
}


func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start collector v%s", p.Version))
	host := p.CfgStringOr("bind-host", "0.0.0.0")
	port := p.CfgStringOr("bind-port", "8070")
	//
	router := mux.NewRouter()
	router.HandleFunc("/api/put", p.putHandler).Methods("POST")

	n := negroni.New(negroni.NewLogger())
	n.Use(recovery.JSONRecovery(true))
	n.UseHandler(router)

	n.Run(host+":"+port)

	/* Listen for incoming connections.
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		p.Log("error", fmt.Sprintln("Error listening:", err.Error()))
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	p.Log("info", fmt.Sprintln("Listening on " + host + ":" + port))
	go p.handleRequests(l)
	for {
		select {
		case msg := <- p.buffer:
			switch msg.(type) {
			case IncommingMsg:
				im := msg.(IncommingMsg)
				base := qtypes.NewTimedBase("tcp", time.Now())
				qm := qtypes.NewMessage(base, p.Name, qtypes.MsgTCP, im.Msg)
				qm.KV["host"] = im.Host
				go p.HandleInventoryRequest(qm)
			default:
				p.Log("warn", fmt.Sprintf("Unkown data type: %s", reflect.TypeOf(msg)))
			}
		}
	}
	*/
}

func (p *Plugin) putHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		p.Log("info", fmt.Sprintf("Received: %v", string(body)))
		fmt.Fprint(w, "OK")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
