// Copyright 2023 NJWS Inc.

package agent

import (
	"strings"

	"github.com/koron/go-ssdp"
)

func (a *Agent) monitor() (err error) {
	log.Info("ssdp monitor")

	m := &ssdp.Monitor{
		Alive: a.onAlive,
		Bye:   a.onBye,
	}

	return m.Start()
}

func (a *Agent) onAlive(m *ssdp.AliveMessage) {
	if !strings.Contains(m.Type, "redfish-rest") {
		return
	}

	if err := a.createOrUpdateAliveMessage(m); err != nil {
		log.Error(err)
		return
	}
}

func (a *Agent) onBye(m *ssdp.ByeMessage) {
	if !strings.Contains(m.Type, "redfish-rest") {
		return
	}

	log.Infof("Bye: From=%s Type=%s USN=%s", m.From.String(), m.Type, m.USN)
}
