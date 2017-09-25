package validator

import (
	"fmt"

	tapi "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/k8sdb/apimachinery/pkg/docker"
	amv "github.com/k8sdb/apimachinery/pkg/validator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

func ValidateXdb(client clientset.Interface, xdb *tapi.Xdb) error {
	if xdb.Spec.Version == "" {
		return fmt.Errorf(`Object 'Version' is missing in '%v'`, xdb.Spec)
	}

	version := fmt.Sprintf("%v-db", xdb.Spec.Version)
	if err := docker.CheckDockerImageVersion(docker.ImageXdb, version); err != nil {
		return fmt.Errorf(`Image %v:%v not found`, docker.ImageXdb, version)
	}

	if xdb.Spec.Storage != nil {
		var err error
		if err = amv.ValidateStorage(client, xdb.Spec.Storage); err != nil {
			return err
		}
	}

	databaseSecret := xdb.Spec.DatabaseSecret
	if databaseSecret != nil {
		if _, err := client.CoreV1().Secrets(xdb.Namespace).Get(databaseSecret.SecretName, metav1.GetOptions{}); err != nil {
			return err
		}
	}

	backupScheduleSpec := xdb.Spec.BackupSchedule
	if backupScheduleSpec != nil {
		if err := amv.ValidateBackupSchedule(client, backupScheduleSpec, xdb.Namespace); err != nil {
			return err
		}
	}

	monitorSpec := xdb.Spec.Monitor
	if monitorSpec != nil {
		if err := amv.ValidateMonitorSpec(monitorSpec); err != nil {
			return err
		}

	}
	return nil
}
