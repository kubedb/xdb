package api

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

const (
	ResourceKindXdb = "Xdb"
	ResourceNameXdb = "xdb"
	ResourceTypeXdb = "xdbes"
)

// Xdb defines a Xdb database.
type Xdb struct {
	unversioned.TypeMeta `json:",inline,omitempty"`
	api.ObjectMeta       `json:"metadata,omitempty"`
	Spec                 XdbSpec   `json:"spec,omitempty"`
	Status               XdbStatus `json:"status,omitempty"`
}

type XdbSpec struct {
	// Version of Xdb to be deployed.
	Version string `json:"version,omitempty"`
	// Number of instances to deploy for a Xdb database.
	Replicas int32 `json:"replicas,omitempty"`
	// Storage spec to specify how storage shall be used.
	Storage *StorageSpec `json:"storage,omitempty"`
	// ServiceAccountName is the name of the ServiceAccount to use to run the
	// Prometheus Pods.
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// Database authentication secret
	DatabaseSecret *api.SecretVolumeSource `json:"databaseSecret,omitempty"`
	// NodeSelector is a selector which must be true for the pod to fit on a node
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Run initial script when starting Xdb master
	InitialScript *InitialScriptSpec `json:"initialScript,omitempty"`
	// BackupSchedule spec to specify how database backup will be taken
	// +optional
	BackupSchedule *BackupScheduleSpec `json:"backupSchedule,omitempty"`
	// If DoNotDelete is true, controller will prevent to delete this Xdb object.
	// Controller will create same Xdb object and ignore other process.
	// +optional
	DoNotDelete bool `json:"doNotDelete,omitempty"`
}

type XdbStatus struct {
	Created        *unversioned.Time `json:"created,omitempty"`
	DatabaseStatus `json:",inline,omitempty"`
	Reason         string `json:"reason,omitempty"`
}

type XdbList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Xdb TPR objects
	Items []*Xdb `json:"items,omitempty"`
}
