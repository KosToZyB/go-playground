package consul

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	consul "github.com/hashicorp/consul/api"
)

//Client provides an interface for getting data out of Consul
type Client interface {
	// Get a Service from consul
	Service(string, string) ([]*consul.ServiceEntry, *consul.QueryMeta, error)
	// Register a service with local agent
	Register(string, string, string, int) error
	// Deregister a service with local agent
	DeRegister(string) error
}

type client struct {
	consul *consul.Client
}

//NewConsul returns a Client interface for given consul address
func NewConsulClient(addr string) (Client, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	c, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &client{consul: c}, nil
}

// Register a service with consul local agent
func (c *client) Register(service, hostname, protocol string, port int) error {
	rand.Seed(time.Now().UnixNano())
	sid := rand.Intn(65534)

	serviceID := service + "-" + strconv.Itoa(sid)

	reg := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    service,
		Port:    port,
		Address: hostname,
		Check: &consul.AgentServiceCheck{
			Script:   "curl --connect-timeout=5 " + protocol + "://" + hostname + ":" + strconv.Itoa(port) + "/health",
			Interval: "10s",
			Timeout:  "8s",
			TTL:      "",
			HTTP:     protocol + "://" + hostname + ":" + strconv.Itoa(port) + "/health",
			Status:   "passing",
		},
		Checks: consul.AgentServiceChecks{},
	}

	return c.consul.Agent().ServiceRegister(reg)
}

// DeRegister a service with consul local agent
func (c *client) DeRegister(id string) error {
	return c.consul.Agent().ServiceDeregister(id)
}

// Service return a service
func (c *client) Service(service, tag string) ([]*consul.ServiceEntry, *consul.QueryMeta, error) {
	passingOnly := true
	addrs, meta, err := c.consul.Health().Service(service, tag, passingOnly, nil)
	if len(addrs) == 0 && err == nil {
		return nil, nil, fmt.Errorf("service ( %s ) was not found", service)
	}
	if err != nil {
		return nil, nil, err
	}
	return addrs, meta, nil
}
