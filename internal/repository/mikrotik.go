package repository

import (
	"fmt"

	"gopkg.in/routeros.v2"
)

//TODO:Implement custom errors

type Interface struct {
	Name string
}

type Resources struct {
	cpu    string
	uptime string
	err    error
}

type Traffic struct {
	rx  string
	tx  string
	err error
}

type MikrotikRepository interface {
	GetTraffic(interfaceName Interface) Traffic
	GetResources() Resources
	GetIndentity() (string, error)
}

type mikrotikRepository struct {
	client *routeros.Client
} //Receiver

func New(address, username, password, port string) (MikrotikRepository, error) {

	formatAddress := fmt.Sprintf("%s:%s", address, port)

	client, err := routeros.Dial(formatAddress, username, password)

	if err != nil {
		return nil, err
	}

	return &mikrotikRepository{client}, nil
}

func (r *mikrotikRepository) GetIndentity() (string, error) {

	identity, err := r.client.Run("/system/identity/print")

	if err != nil {
		return "", err
	}

	defer r.client.Close()

	var name string

	for _, x := range identity.Re {
		name = x.Map["name"]
	}

	return name, nil
}

func (r *mikrotikRepository) GetTraffic(interfaceName Interface) Traffic {

	traffic, err := r.client.Run("/interface/monitor-traffic", "=interface="+interfaceName.Name, "=once")

	if err != nil {
		return Traffic{"", "", err}
	}

	var rx, tx string

	for _, x := range traffic.Re {
		rx = x.Map["tx-bits-per-second"]
		tx = x.Map["tx-bits-per-second"]
	}

	return Traffic{rx, tx, nil}

}

func (r *mikrotikRepository) GetResources() Resources {
	resources, err := r.client.Run("/system/resources/print")

	if err != nil {
		return Resources{"", "", nil}
	}

	var cpu, uptime string

	for _, x := range resources.Re {
		cpu = x.Map["cpu-load"] + "%"
		uptime = x.Map["uptime"]
	}

	return Resources{cpu, uptime, nil}
}
