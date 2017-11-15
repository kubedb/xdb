package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/appscode/kutil/tools/monitoring/api"
	core "k8s.io/api/core/v1"
)

func (p Xdb) OffshootName() string {
	return p.Name
}

func (p Xdb) OffshootLabels() map[string]string {
	return map[string]string{
		LabelDatabaseName: p.Name,
		LabelDatabaseKind: ResourceKindXdb,
	}
}

func (p Xdb) StatefulSetLabels() map[string]string {
	labels := p.OffshootLabels()
	for key, val := range p.Labels {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, XdbKey+"/") {
			labels[key] = val
		}
	}
	return labels
}

func (p Xdb) StatefulSetAnnotations() map[string]string {
	annotations := make(map[string]string)
	for key, val := range p.Annotations {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, XdbKey+"/") {
			annotations[key] = val
		}
	}
	annotations[XdbDatabaseVersion] = string(p.Spec.Version)
	return annotations
}

var _ ResourceInfo = &Xdb{}

func (p Xdb) ResourceCode() string {
	return ResourceCodeXdb
}

func (p Xdb) ResourceKind() string {
	return ResourceKindXdb
}

func (p Xdb) ResourceName() string {
	return ResourceNameXdb
}

func (p Xdb) ResourceType() string {
	return ResourceTypeXdb
}

func (p Xdb) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            p.ResourceKind(),
		Namespace:       p.Namespace,
		Name:            p.Name,
		UID:             p.UID,
		ResourceVersion: p.ResourceVersion,
	}
}

func (p Xdb) ServiceName() string {
	return p.OffshootName()
}

func (p Xdb) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", p.Namespace, p.Name)
}

func (p Xdb) Path() string {
	return fmt.Sprintf("/kubedb.com/v1alpha1/namespaces/%s/%s/%s/metrics", p.Namespace, p.ResourceType(), p.Name)
}

func (p Xdb) Scheme() string {
	return ""
}

func (p *Xdb) StatsAccessor() api.StatsAccessor {
	return p
}
