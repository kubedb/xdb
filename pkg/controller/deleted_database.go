package controller

import (
	tapi "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	kapi "k8s.io/kubernetes/pkg/api"
)

func (c *Controller) Exists(om *kapi.ObjectMeta) (bool, error) {
	return false, nil
}

func (c *Controller) DeleteDatabase(deletedDb *tapi.DormantDatabase) error {
	return nil
}

func (c *Controller) DestroyDatabase(deletedDb *tapi.DormantDatabase) error {
	return nil
}

func (c *Controller) RecoverDatabase(deletedDb *tapi.DormantDatabase) error {
	return nil
}
