package controller

import (
	tapi "github.com/k8sdb/apimachinery/api"
	kapi "k8s.io/kubernetes/pkg/api"
)

func (c *Controller) Exists(om *kapi.ObjectMeta) (bool, error) {
	return false, nil
}

func (c *Controller) DeleteDatabase(deletedDb *tapi.DeletedDatabase) error {
	return nil
}

func (c *Controller) DestroyDatabase(deletedDb *tapi.DeletedDatabase) error {
	return nil
}

func (c *Controller) RecoverDatabase(deletedDb *tapi.DeletedDatabase) error {
	return nil
}
