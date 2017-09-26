package controller

import (
	"fmt"

	tapi "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/k8sdb/apimachinery/pkg/monitor"
)

func (c *Controller) newMonitorController(xdb *tapi.Xdb) (monitor.Monitor, error) {
	monitorSpec := xdb.Spec.Monitor

	if monitorSpec == nil {
		return nil, fmt.Errorf("MonitorSpec not found in %v", xdb.Spec)
	}

	if monitorSpec.Prometheus != nil {
		return monitor.NewPrometheusController(c.Client, c.ApiExtKubeClient, c.promClient, c.opt.OperatorNamespace), nil
	}

	return nil, fmt.Errorf("Monitoring controller not found for %v", monitorSpec)
}

func (c *Controller) addMonitor(xdb *tapi.Xdb) error {
	ctrl, err := c.newMonitorController(xdb)
	if err != nil {
		return err
	}
	return ctrl.AddMonitor(xdb.ObjectMeta, xdb.Spec.Monitor)
}

func (c *Controller) deleteMonitor(xdb *tapi.Xdb) error {
	ctrl, err := c.newMonitorController(xdb)
	if err != nil {
		return err
	}
	return ctrl.DeleteMonitor(xdb.ObjectMeta, xdb.Spec.Monitor)
}

func (c *Controller) updateMonitor(oldXdb, updatedXdb *tapi.Xdb) error {
	var err error
	var ctrl monitor.Monitor
	if updatedXdb.Spec.Monitor == nil {
		ctrl, err = c.newMonitorController(oldXdb)
	} else {
		ctrl, err = c.newMonitorController(updatedXdb)
	}
	if err != nil {
		return err
	}
	return ctrl.UpdateMonitor(updatedXdb.ObjectMeta, oldXdb.Spec.Monitor, updatedXdb.Spec.Monitor)
}
