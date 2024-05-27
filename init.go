package bamboo

import (
	"flag"
	"net/http"

	"github.com/Rachelgill00/simulation-for-calvin/config"
	"github.com/Rachelgill00/simulation-for-calvin/log"
)

func Init() {
	flag.Parse()
	log.Setup()
	config.Configuration.Load()
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1000
}
