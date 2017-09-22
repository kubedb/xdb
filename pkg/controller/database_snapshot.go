package controller

import (
	tapi "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	kbatch "k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/runtime"
)

func (c *Controller) ValidateSnapshot(dbSnapshot *tapi.Snapshot) error {
	return nil
}

func (c *Controller) GetDatabase(snapshot *tapi.Snapshot) (runtime.Object, error) {
	return nil, nil
}

func (c *Controller) GetSnapshotter(snapshot *tapi.Snapshot) (*kbatch.Job, error) {
	return nil, nil
}

func (c *Controller) DestroySnapshot(dbSnapshot *tapi.Snapshot) error {
	return nil
}
