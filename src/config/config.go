package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Component struct {
	Type       string `yaml:"type"`
	Name       string `yaml:"name"`
	Service    string `yaml:"service"`
	Port       int    `yaml:"port"`
	BridgePort string `yaml:"bridge-port"`
}

type Host struct {
	Name       string      `yaml:"name"`
	Hostname   string      `yaml:"hostname"`
	Components []Component `yaml:"components"`
}

type Mapping struct {
	Component  *Component
	HostName   string
	Enviroment string
	Service    string
	Port       string
}

func LoadConfig(configFileName string) ([]Host, error) {
	jsonFile, err := os.Open(configFileName)
	if err != nil {
		return nil, errors.Join(err)
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var configuration []Host
	err = yaml.Unmarshal(byteValue, &configuration)
	if err != nil {
		return nil, errors.Join(err)
	}
	return configuration, nil
}

func GetService(configuration []Host, cmd string) (*Mapping, error) {

	components := strings.Split(cmd, ":")
	if len(components) != 3 {
		return nil, fmt.Errorf("malformed command: '%s'", cmd)
	}

	envPtr := strings.TrimSpace(components[0])
	componentPtr := strings.TrimSpace(components[1])
	portPtr := strings.TrimSpace(components[2])

	if envPtr == "" {
		return nil, errors.New("env missing")
	}

	if componentPtr == "" {
		return nil, errors.New("component missing")
	}

	if portPtr == "" {
		return nil, errors.New("port missing")
	}

	var hostFound *Host = nil
	for _, host := range configuration {
		if host.Name == envPtr {
			hostFound = &host
			break
		}
	}
	if hostFound == nil {
		return nil, fmt.Errorf("env not found: %s", envPtr)
	}

	hostname := hostFound.Hostname

	if hostname == "" {
		return nil, errors.New("hostname missing")
	}

	var compPort *Component = nil
	for _, item := range hostFound.Components {
		if item.Service == componentPtr || item.Name == componentPtr {
			compPort = &item
			break
		}
	}
	if compPort == nil {
		return nil, fmt.Errorf("component not found: %s", componentPtr)
	}

	return &Mapping{
		Component:  compPort,
		HostName:   hostname,
		Enviroment: envPtr,
		Service:    componentPtr,
		Port:       portPtr,
	}, nil

}
