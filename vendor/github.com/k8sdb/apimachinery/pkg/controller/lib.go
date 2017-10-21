package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/appscode/go/log"
	"github.com/appscode/go/types"
	kutilapps "github.com/appscode/kutil/apps/v1beta1"
	"github.com/graymeta/stow"
	_ "github.com/graymeta/stow/azure"
	_ "github.com/graymeta/stow/google"
	_ "github.com/graymeta/stow/s3"
	tapi "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/k8sdb/apimachinery/pkg/eventer"
	"github.com/k8sdb/apimachinery/pkg/storage"
	apps "k8s.io/api/apps/v1beta1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
)

func (c *Controller) CheckStatefulSetPodStatus(statefulSet *apps.StatefulSet, checkDuration time.Duration) error {
	podName := fmt.Sprintf("%v-%v", statefulSet.Name, 0)

	podReady := false
	then := time.Now()
	now := time.Now()
	for now.Sub(then) < checkDuration {
		pod, err := c.Client.CoreV1().Pods(statefulSet.Namespace).Get(podName, metav1.GetOptions{})
		if err != nil {
			if kerr.IsNotFound(err) {
				_, err := c.Client.AppsV1beta1().StatefulSets(statefulSet.Namespace).Get(statefulSet.Name, metav1.GetOptions{})
				if kerr.IsNotFound(err) {
					break
				}

				time.Sleep(sleepDuration)
				now = time.Now()
				continue
			} else {
				return err
			}
		}
		log.Debugf("Pod Phase: %v", pod.Status.Phase)

		// If job is success
		if pod.Status.Phase == core.PodRunning {
			podReady = true
			break
		}

		time.Sleep(sleepDuration)
		now = time.Now()
	}
	if !podReady {
		return errors.New("Database fails to be Ready")
	}
	return nil
}

func (c *Controller) DeletePersistentVolumeClaims(namespace string, selector labels.Selector) error {
	pvcList, err := c.Client.CoreV1().PersistentVolumeClaims(namespace).List(
		metav1.ListOptions{
			LabelSelector: selector.String(),
		},
	)
	if err != nil {
		return err
	}

	for _, pvc := range pvcList.Items {
		if err := c.Client.CoreV1().PersistentVolumeClaims(pvc.Namespace).Delete(pvc.Name, nil); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) DeleteSnapshotData(snapshot *tapi.Snapshot) error {
	cfg, err := storage.NewOSMContext(c.Client, snapshot.Spec.SnapshotStorageSpec, snapshot.Namespace)
	if err != nil {
		return err
	}

	loc, err := stow.Dial(cfg.Provider, cfg.Config)
	if err != nil {
		return err
	}
	bucket, err := snapshot.Spec.SnapshotStorageSpec.Container()
	if err != nil {
		return err
	}
	container, err := loc.Container(bucket)
	if err != nil {
		return err
	}

	prefix, _ := snapshot.Location() // error checked by .Container()
	cursor := stow.CursorStart
	for {
		items, next, err := container.Items(prefix, cursor, 50)
		if err != nil {
			return err
		}
		for _, item := range items {
			if err := container.RemoveItem(item.ID()); err != nil {
				return err
			}
		}
		cursor = next
		if stow.IsCursorEnd(cursor) {
			break
		}
	}

	return nil
}

func (c *Controller) DeleteSnapshots(namespace string, selector labels.Selector) error {
	snapshotList, err := c.ExtClient.Snapshots(namespace).List(
		metav1.ListOptions{
			LabelSelector: selector.String(),
		},
	)
	if err != nil {
		return err
	}

	for _, snapshot := range snapshotList.Items {
		if err := c.ExtClient.Snapshots(snapshot.Namespace).Delete(snapshot.Name, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) CheckDatabaseRestoreJob(
	job *batch.Job,
	runtimeObj runtime.Object,
	recorder record.EventRecorder,
	checkDuration time.Duration,
) bool {
	var jobSuccess bool = false
	var err error

	then := time.Now()
	now := time.Now()
	for now.Sub(then) < checkDuration {
		log.Debugln("Checking for Job ", job.Name)
		job, err = c.Client.BatchV1().Jobs(job.Namespace).Get(job.Name, metav1.GetOptions{})
		if err != nil {
			if kerr.IsNotFound(err) {
				time.Sleep(sleepDuration)
				now = time.Now()
				continue
			}
			recorder.Eventf(
				tapi.ObjectReferenceFor(runtimeObj),
				core.EventTypeWarning,
				eventer.EventReasonFailedToList,
				"Failed to get Job. Reason: %v",
				err,
			)
			log.Errorln(err)
			return jobSuccess
		}
		log.Debugf("Pods Statuses:	%d Running / %d Succeeded / %d Failed",
			job.Status.Active, job.Status.Succeeded, job.Status.Failed)
		// If job is success
		if job.Status.Succeeded > 0 {
			jobSuccess = true
			break
		} else if job.Status.Failed > 0 {
			break
		}

		time.Sleep(sleepDuration)
		now = time.Now()
	}

	if err != nil {
		return false
	}

	deleteJobResources(c.Client, recorder, runtimeObj, job)

	err = c.Client.CoreV1().Secrets(job.Namespace).Delete(job.Name, &metav1.DeleteOptions{})
	if err != nil && !kerr.IsNotFound(err) {
		return false
	}

	return jobSuccess
}

func (c *Controller) checkGoverningService(name, namespace string) (bool, error) {
	_, err := c.Client.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func (c *Controller) CreateGoverningService(name, namespace string) error {
	// Check if service name exists
	found, err := c.checkGoverningService(name, namespace)
	if err != nil {
		return err
	}
	if found {
		return nil
	}

	service := &core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: core.ServiceSpec{
			Type:      core.ServiceTypeClusterIP,
			ClusterIP: core.ClusterIPNone,
		},
	}
	_, err = c.Client.CoreV1().Services(namespace).Create(service)
	return err
}

func (c *Controller) DeleteService(name, namespace string) error {
	service, err := c.Client.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	if service.Spec.Selector[tapi.LabelDatabaseName] != name {
		return nil
	}

	return c.Client.CoreV1().Services(namespace).Delete(name, nil)
}

func (c *Controller) DeleteStatefulSet(name, namespace string) error {
	statefulSet, err := c.Client.AppsV1beta1().StatefulSets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	// Update StatefulSet
	_, err = kutilapps.TryPatchStatefulSet(c.Client, statefulSet.ObjectMeta, func(in *apps.StatefulSet) *apps.StatefulSet {
		in.Spec.Replicas = types.Int32P(0)
		return in
	})
	if err != nil {
		return err
	}

	var checkSuccess bool = false
	then := time.Now()
	now := time.Now()
	for now.Sub(then) < time.Minute*10 {
		podList, err := c.Client.CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{
			LabelSelector: labels.Set(statefulSet.Spec.Selector.MatchLabels).AsSelector().String(),
		})
		if err != nil {
			return err
		}
		if len(podList.Items) == 0 {
			checkSuccess = true
			break
		}

		time.Sleep(sleepDuration)
		now = time.Now()
	}

	if !checkSuccess {
		return errors.New("Fail to delete StatefulSet Pods")
	}

	// Delete StatefulSet
	return c.Client.AppsV1beta1().StatefulSets(statefulSet.Namespace).Delete(statefulSet.Name, nil)
}

func (c *Controller) DeleteSecret(name, namespace string) error {
	if _, err := c.Client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{}); err != nil {
		if kerr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	return c.Client.CoreV1().Secrets(namespace).Delete(name, nil)
}

func deleteJobResources(
	client clientset.Interface,
	recorder record.EventRecorder,
	runtimeObj runtime.Object,
	job *batch.Job,
) {
	if err := client.BatchV1().Jobs(job.Namespace).Delete(job.Name, nil); err != nil && !kerr.IsNotFound(err) {
		recorder.Eventf(
			tapi.ObjectReferenceFor(runtimeObj),
			core.EventTypeWarning,
			eventer.EventReasonFailedToDelete,
			"Failed to delete Job. Reason: %v",
			err,
		)
		log.Errorln(err)
	}

	r, err := metav1.LabelSelectorAsSelector(job.Spec.Selector)
	if err != nil {
		log.Errorln(err)
	} else {
		err = client.CoreV1().Pods(job.Namespace).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{
			LabelSelector: r.String(),
		})
		if err != nil {
			recorder.Eventf(
				tapi.ObjectReferenceFor(runtimeObj),
				core.EventTypeWarning,
				eventer.EventReasonFailedToDelete,
				"Failed to delete Pods. Reason: %v",
				err,
			)
			log.Errorln(err)
		}
	}

	for _, volume := range job.Spec.Template.Spec.Volumes {
		claim := volume.PersistentVolumeClaim
		if claim != nil {
			err := client.CoreV1().PersistentVolumeClaims(job.Namespace).Delete(claim.ClaimName, nil)
			if err != nil && !kerr.IsNotFound(err) {
				recorder.Eventf(
					tapi.ObjectReferenceFor(runtimeObj),
					core.EventTypeWarning,
					eventer.EventReasonFailedToDelete,
					"Failed to delete PersistentVolumeClaim. Reason: %v",
					err,
				)
				log.Errorln(err)
			}
		}
	}

	return
}
