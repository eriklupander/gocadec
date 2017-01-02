package main

import (
	"github.com/callistaenterprise/gocadec/compositeservice/service"
	ct "github.com/eriklupander/cloudtoolkit"
	"github.com/spf13/viper"
	"sync"
	"github.com/callistaenterprise/gocadec/compositeservice/client"
)

var appName = "compservice"

// var EnvProfile string = "dev"

var configServerDefaultUrl = "http://configserver:8888"
var amqpClient *ct.MessagingClient

func main() {
	ct.Log.Println("Starting " + appName + "...")
	// First of all, dump various hostname ips to log to see what mood the DNS resolver is in... :(
	ct.DumpDNS()
	ct.LoadSpringCloudConfig(appName, ct.ResolveProfile(), configServerDefaultUrl)
	ct.InitTracingFromConfigProperty(appName)

	amqpClient = ct.InitMessagingClientFromConfigProperty()
	defer amqpClient.GetConn().Close()

	ct.ConfigureHystrix([]string{"get_account_secured"}, amqpClient)

	// Configure the HTTP client (disable Keep-alives so Docker Swarm round-robins for us)
	client.ConfigureClient()

	// Starts HTTP service  (async)
	go service.StartWebServer(viper.GetString("server_port"))

	// Block...
	wg := sync.WaitGroup{} // Use a WaitGroup to block main() exit
	wg.Add(1)
	wg.Wait()
}

