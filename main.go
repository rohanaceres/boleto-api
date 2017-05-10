package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"os"

	"bitbucket.org/mundipagg/boletoapi/api"
	"bitbucket.org/mundipagg/boletoapi/config"
	"bitbucket.org/mundipagg/boletoapi/log"
	"bitbucket.org/mundipagg/boletoapi/robot"
)

var (
	processID  = os.Getpid()
	totalProcs = runtime.NumCPU()
	devMode    = flag.Bool("dev", false, "-dev To run in dev mode")
	mockMode   = flag.Bool("mock", false, "-mock To run mock requests")
)

func init() {
	createPIDfile()
	// Configure signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go handleSignal(c)
}

func handleSignal(c chan os.Signal) {
	<-c
	config.Stop()
	log.Info("Quiting BoletoApi")
	log.Close()
	//db.GetDB().Close()
	fmt.Println("Done")
	os.Exit(1)
}

func createPIDfile() {
	p := strconv.Itoa(processID)
	ioutil.WriteFile("boletoapi.pid", []byte(p), 0644)
}
func configFlags() {
	if *devMode {
		fmt.Println("Execute in DEV mode")
		os.Setenv("API_PORT", "3000")
		os.Setenv("API_VERSION", "0.0.1")
		os.Setenv("ENVIROMENT", "Development")
		os.Setenv("SEQ_URL", "http://192.168.8.119:5341") // http://stglog.mundipagg.com/
		os.Setenv("SEQ_API_KEY", "4jZzTybZ9bUHtJiPdh6")   //4jZzTybZ9bUHtJiPdh6
		os.Setenv("ENABLE_REQUEST_LOG", "false")
		os.Setenv("ENABLE_PRINT_REQUEST", "true")
		os.Setenv("URL_BB_REGISTER_BOLETO", "https://cobranca.homologa.bb.com.br:7101/registrarBoleto")
		os.Setenv("URL_BB_TOKEN", "https://oauth.hm.bb.com.br:43000/oauth/token")
		os.Setenv("APP_URL", "http://localhost:3000/boleto")
		os.Setenv("ELASTIC_URL", "http://localhost:9200")
		os.Setenv("MONGODB_URL", "localhost:27017")
		os.Setenv("BOLETO_JSON_STORE", "/home/philippe/boletodb/upMongo")
		if *mockMode {
			os.Setenv("URL_BB_REGISTER_BOLETO", "http://localhost:4000/registrarBoleto")
			os.Setenv("URL_BB_TOKEN", "http://localhost:4000/oauth/token")
		}
	}
	config.Install(*mockMode, *devMode)
}
func main() {
	flag.Parse()
	configFlags()
	logo1()
	installLog()
	robot.GoRobots()
	api.InstallRestAPI()

}

func logo1() {
	l := `
$$$$$$$\            $$\            $$\                $$$$$$\            $$\ 
$$  __$$\           $$ |           $$ |              $$  __$$\           \__|
$$ |  $$ | $$$$$$\  $$ | $$$$$$\ $$$$$$\    $$$$$$\  $$ /  $$ | $$$$$$\  $$\ 
$$$$$$$\ |$$  __$$\ $$ |$$  __$$\\_$$  _|  $$  __$$\ $$$$$$$$ |$$  __$$\ $$ |
$$  __$$\ $$ /  $$ |$$ |$$$$$$$$ | $$ |    $$ /  $$ |$$  __$$ |$$ /  $$ |$$ |
$$ |  $$ |$$ |  $$ |$$ |$$   ____| $$ |$$\ $$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |
$$$$$$$  |\$$$$$$  |$$ |\$$$$$$$\  \$$$$  |\$$$$$$  |$$ |  $$ |$$$$$$$  |$$ |
\_______/  \______/ \__| \_______|  \____/  \______/ \__|  \__|$$  ____/ \__|
                                                               $$ |          
                                                               $$ |          
                                                               \__|          
	`
	fmt.Println(l)
	fmt.Println("Version: " + config.Get().Version)
}

func installLog() {
	err := log.Install()
	if err != nil {
		fmt.Println("Log SEQ Fails")
		os.Exit(-1)
	}
}
