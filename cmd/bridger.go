package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/mastrogiovanni/bridger/src/config"
	"github.com/mastrogiovanni/bridger/src/executor"
)

func usage() {
	log.Println("usage: bridger <path to config file> env_1:service_1:port_1 ... env_n:service_n:port_n")
}

func main() {

	var stopChan = make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	configuration, err := config.LoadConfig(os.Args[1])
	if err != nil {
		usage()
		log.Fatal(err)
		return
	}

	mappings := os.Args[2:]

	// WaitGroup to ensure graceful shutdown
	var wg sync.WaitGroup

	interruptCh := make(chan struct{})

	for _, item := range mappings {

		mapping, err := config.GetService(configuration, item)
		if err != nil {
			log.Fatal(err)
			usage()
			return
		}

		log.Printf("Connecting to %s:%d in %s and exposing it on local port %s", mapping.Service, mapping.Component.Port, mapping.Enviroment, mapping.Port)

		sshCmd := ""

		if mapping.Component.Type == "docker" {

			ipCmd := "ssh %s docker inspect %s | jq -r .[0].NetworkSettings.Networks.compose_default.IPAddress"
			ipCmd = fmt.Sprintf(ipCmd, mapping.HostName, mapping.Component.Service)

			log.Printf("Fetching IP address of the component\n")
			ipOut := executor.ExecCmd(ipCmd)

			log.Printf("IP address: %s", ipOut)

			sshCmd = fmt.Sprintf("ssh -v -L 0.0.0.0:%s:%s:%d %s -N", mapping.Port, ipOut, mapping.Component.Port, mapping.HostName)

			log.Printf("Connecting...")
			log.Printf("Component %s will be available on http://127.0.0.1:%s", mapping.Service, mapping.Port)

		} else if mapping.Component.Type == "kubernetes" {

			if mapping.Component.BridgePort != "" {

				forward := fmt.Sprintf("kubectl port-forward svc/%s %s:%d", mapping.Component.Service, mapping.Component.BridgePort, mapping.Component.Port)
				sshCmd = fmt.Sprintf("ssh -v -L 0.0.0.0:%s:127.0.0.1:%s %s %s", mapping.Port, mapping.Component.BridgePort, mapping.HostName, forward)

			} else {

				forward := fmt.Sprintf("kubectl port-forward svc/%s %d:%d", mapping.Component.Service, mapping.Component.Port, mapping.Component.Port)
				sshCmd = fmt.Sprintf("ssh -v -L 0.0.0.0:%s:127.0.0.1:%d %s %s", mapping.Port, mapping.Component.Port, mapping.HostName, forward)

			}

			log.Printf("Connecting...")
			log.Printf("Component %s will be available on http://127.0.0.1:%s", mapping.Service, mapping.Port)

		} else {

			log.Println("Skip")
			continue

		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			outputCh := make(chan string)
			errorCh := make(chan string)

			go func() {
				for line := range outputCh {
					fmt.Println(line)
				}
				fmt.Println("SSH Tunnel Closed (output)")
			}()

			go func() {
				for line := range errorCh {
					if strings.Contains(line, "requested") { // || strings.Contains(line, "free:") {
						fmt.Println(line)
					}
				}
				fmt.Println("SSH Tunnel Closed (error)")
			}()

			comps := strings.Split(sshCmd, " ")
			executor.ExecuteCommandAsync(comps[0], comps[1:], outputCh, errorCh, interruptCh)
			// out := executor.ExecCmd(sshCmd)
			// log.Println(out)
		}()

	}

	<-stopChan // wait for SIGINT
	log.Println("Interrupted")
	close(interruptCh)

	wg.Wait()
	fmt.Println("Program terminated.")

}
