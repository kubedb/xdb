package controller

import (
	"fmt"

	"github.com/appscode/kutil/tools/monitoring/agents"
	mona "github.com/appscode/kutil/tools/monitoring/api"
	api "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
)

func (c *Controller) newMonitorController(xdb *api.Xdb) (mona.Agent, error) {
	monitorSpec := xdb.Spec.Monitor

	if monitorSpec == nil {
		return nil, fmt.Errorf("MonitorSpec not found in %v", xdb.Spec)
	}

	if monitorSpec.Prometheus != nil {
		return agents.New(monitorSpec.Agent, c.Client, c.ApiExtKubeClient, c.promClient), nil
	}

	return nil, fmt.Errorf("monitoring controller not found for %v", monitorSpec)
}

func (c *Controller) addMonitor(xdb *api.Xdb) error {
	agent, err := c.newMonitorController(xdb)
	if err != nil {
		return err
	}
	return agent.Add(xdb.StatsAccessor(), xdb.Spec.Monitor)
}

func (c *Controller) deleteMonitor(xdb *api.Xdb) error {
	agent, err := c.newMonitorController(xdb)
	if err != nil {
		return err
	}
	return agent.Delete(xdb.StatsAccessor(), xdb.Spec.Monitor)
}

func (c *Controller) updateMonitor(oldXdb, updatedXdb *api.Xdb) error {
	var err error
	var agent mona.Agent
	if updatedXdb.Spec.Monitor == nil {
		agent, err = c.newMonitorController(oldXdb)
	} else {
		agent, err = c.newMonitorController(updatedXdb)
	}
	if err != nil {
		return err
	}
	return agent.Update(updatedXdb.StatsAccessor(), oldXdb.Spec.Monitor, updatedXdb.Spec.Monitor)
}
