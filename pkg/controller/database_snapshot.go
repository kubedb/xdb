package controller

import (
	tapi "github.com/k8sdb/apimachinery/api"
	kbatch "k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/runtime"
)

func (c *Controller) ValidateSnapshot(dbSnapshot *tapi.DatabaseSnapshot) error {
	return nil
}

func (c *Controller) GetDatabase(snapshot *tapi.DatabaseSnapshot) (runtime.Object, error) {
	return nil, nil
}

func (c *Controller) GetSnapshotter(snapshot *tapi.DatabaseSnapshot) (*kbatch.Job, error) {
	return nil, nil
}

func (c *Controller) DestroySnapshot(dbSnapshot *tapi.DatabaseSnapshot) error {
	return nil
}
