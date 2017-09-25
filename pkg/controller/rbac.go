package controller

import (
	kutilcore "github.com/appscode/kutil/core/v1"
	kutilrbac "github.com/appscode/kutil/rbac/v1beta1"
	"github.com/k8sdb/apimachinery/apis/kubedb"
	tapi "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	rbac "k8s.io/client-go/pkg/apis/rbac/v1beta1"
)

func (c *Controller) deleteRole(xdb *tapi.Xdb) error {
	// Delete existing Roles
	if err := c.Client.RbacV1beta1().Roles(xdb.Namespace).Delete(xdb.OffshootName(), nil); err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (c *Controller) createRole(xdb *tapi.Xdb) error {
	// Create new Roles
	_, err := kutilrbac.EnsureRole(
		c.Client,
		metav1.ObjectMeta{
			Name:      xdb.OffshootName(),
			Namespace: xdb.Namespace,
		},
		func(in *rbac.Role) *rbac.Role {
			in.Rules = []rbac.PolicyRule{
				{
					APIGroups:     []string{kubedb.GroupName},
					Resources:     []string{tapi.ResourceTypeXdb},
					ResourceNames: []string{xdb.Name},
					Verbs:         []string{"get"},
				},
				{
					// Use this if secret is necessary, Otherwise remove it
					APIGroups:     []string{apiv1.GroupName},
					Resources:     []string{"secrets"},
					ResourceNames: []string{xdb.Spec.DatabaseSecret.SecretName},
					Verbs:         []string{"get"},
				},
			}
			return in
		},
	)
	return err
}

func (c *Controller) deleteServiceAccount(xdb *tapi.Xdb) error {
	// Delete existing ServiceAccount
	if err := c.Client.CoreV1().ServiceAccounts(xdb.Namespace).Delete(xdb.OffshootName(), nil); err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (c *Controller) createServiceAccount(xdb *tapi.Xdb) error {
	// Create new ServiceAccount
	_, err := kutilcore.EnsureServiceAccount(
		c.Client,
		metav1.ObjectMeta{
			Name:      xdb.OffshootName(),
			Namespace: xdb.Namespace,
		},
		func(in *apiv1.ServiceAccount) *apiv1.ServiceAccount {
			return in
		},
	)
	return err
}

func (c *Controller) deleteRoleBinding(xdb *tapi.Xdb) error {
	// Delete existing RoleBindings
	if err := c.Client.RbacV1beta1().RoleBindings(xdb.Namespace).Delete(xdb.OffshootName(), nil); err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (c *Controller) createRoleBinding(xdb *tapi.Xdb) error {
	// Ensure new RoleBindings
	_, err := kutilrbac.EnsureRoleBinding(
		c.Client,
		metav1.ObjectMeta{
			Name:      xdb.OffshootName(),
			Namespace: xdb.Namespace,
		},
		func(in *rbac.RoleBinding) *rbac.RoleBinding {
			in.RoleRef = rbac.RoleRef{
				APIGroup: rbac.GroupName,
				Kind:     "Role",
				Name:     xdb.OffshootName(),
			}
			in.Subjects = []rbac.Subject{
				{
					Kind:      rbac.ServiceAccountKind,
					Name:      xdb.OffshootName(),
					Namespace: xdb.Namespace,
				},
			}
			return in
		},
	)
	return err
}

func (c *Controller) createRBACStuff(xdb *tapi.Xdb) error {
	// Delete Existing Role
	if err := c.deleteRole(xdb); err != nil {
		return err
	}
	// Create New Role
	if err := c.createRole(xdb); err != nil {
		return err
	}

	// Create New ServiceAccount
	if err := c.createServiceAccount(xdb); err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
	}

	// Create New RoleBinding
	if err := c.createRoleBinding(xdb); err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}

func (c *Controller) deleteRBACStuff(xdb *tapi.Xdb) error {
	// Delete Existing Role
	if err := c.deleteRole(xdb); err != nil {
		return err
	}

	// Delete ServiceAccount
	if err := c.deleteServiceAccount(xdb); err != nil {
		return err
	}

	// Delete New RoleBinding
	if err := c.deleteRoleBinding(xdb); err != nil {
		return err
	}

	return nil
}
