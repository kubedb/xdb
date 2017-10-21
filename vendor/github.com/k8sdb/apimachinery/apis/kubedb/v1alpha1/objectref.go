package v1alpha1

import (
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (d DormantDatabase) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            ResourceKindDormantDatabase,
		Namespace:       d.Namespace,
		Name:            d.Name,
		UID:             d.UID,
		ResourceVersion: d.ResourceVersion,
	}
}

func (p Postgres) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            ResourceKindPostgres,
		Namespace:       p.Namespace,
		Name:            p.Name,
		UID:             p.UID,
		ResourceVersion: p.ResourceVersion,
	}
}

func (e Elasticsearch) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            ResourceKindElasticsearch,
		Namespace:       e.Namespace,
		Name:            e.Name,
		UID:             e.UID,
		ResourceVersion: e.ResourceVersion,
	}
}

func (s Xdb) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            ResourceKindXdb,
		Namespace:       s.Namespace,
		Name:            s.Name,
		UID:             s.UID,
		ResourceVersion: s.ResourceVersion,
	}
}

func (s Snapshot) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            ResourceKindSnapshot,
		Namespace:       s.Namespace,
		Name:            s.Name,
		UID:             s.UID,
		ResourceVersion: s.ResourceVersion,
	}
}

func ObjectReferenceFor(obj runtime.Object) *core.ObjectReference {
	switch u := obj.(type) {
	case *DormantDatabase:
		return u.ObjectReference()
	case *Postgres:
		return u.ObjectReference()
	case *Elasticsearch:
		return u.ObjectReference()
	case *Snapshot:
		return u.ObjectReference()
	case *Xdb:
		return u.ObjectReference()
	}
	return &core.ObjectReference{}
}
