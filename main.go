package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"gingob/config"
	"gingob/model"
	v "gingob/pkg/version"
	"gingob/router"
	"gingob/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg     = pflag.StringP("config", "c", "", "gingob config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	pflag.Parse()

	// Show version info
	if *version {
		v := v.Get()
		marshalled, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshalled))
		return
	}

	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// init db
	model.DB.Init()
	defer model.DB.Close()

	// Set gin mode.
	gin.SetMode(viper.GetString("runmode"))

	// create the gin engine
	g := gin.New()

	// routes
	router.Load(g,
		// Middiewares
		middleware.RequestId(),
	)

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()

	// If open https, start listening https request
	openHttps := viper.GetBool("tls.https_open")
	if openHttps == true {
		cert := viper.GetString("tls.cert")
		key := viper.GetString("tls.key")
		if cert != "" && key != "" {
			go func() {
				log.Infof("Start to listening the incoming requests on https address: %s", viper.GetString("tls.addr"))
				log.Info(http.ListenAndServeTLS(viper.GetString("tls.addr"), cert, key, g).Error())
			}()
		}
	}

	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}
