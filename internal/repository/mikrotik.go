package repository

import (
	"fmt"

	"gopkg.in/routeros.v2"
)

//TODO:Implement custom errors

type Resources struct {
	Cpu    string
	Uptime string
}

type Traffic struct {
	Rx string
	Tx string
}

type MikrotikRepository interface {
	GetTraffic(i string) (Traffic, error)
	GetResources() (Resources, error)
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

	//defer r.client.Close()

	var name string

	for _, x := range identity.Re {
		name = x.Map["name"]
	}

	return name, nil
}

func (r *mikrotikRepository) GetTraffic(i string) (Traffic, error) {

	traffic, err := r.client.Run("/interface/monitor-traffic", "=interface="+i, "=once")

	if err != nil {
		return Traffic{"", ""}, err
	}

	//defer r.client.Close()

	var rx, tx string

	for _, x := range traffic.Re {
		rx = x.Map["rx-bits-per-second"]
		tx = x.Map["tx-bits-per-second"]
	}

	return Traffic{rx, tx}, nil

}

func (r *mikrotikRepository) GetResources() (Resources, error) {
	resources, err := r.client.Run("/system/resource/print")

	if err != nil {
		return Resources{"", ""}, err
	}

	var cpu, uptime string

	for _, x := range resources.Re {
		cpu = x.Map["cpu-load"]
		uptime = x.Map["uptime"]
	}

	return Resources{cpu, uptime}, nil
}
