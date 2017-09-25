package v1alpha1

import (
	"errors"
	"path/filepath"
	"strings"
)

const (
	DatabaseNamePrefix = "kubedb"

	GenericKey = "kubedb.com"

	LabelDatabaseKind = GenericKey + "/kind"
	LabelDatabaseName = GenericKey + "/name"
	LabelJobType      = GenericKey + "/job-type"

	PostgresKey             = ResourceTypePostgres + "." + GenericKey
	PostgresDatabaseVersion = PostgresKey + "/version"

	ElasticsearchKey             = ResourceTypeElasticsearch + "." + GenericKey
	ElasticsearchDatabaseVersion = ElasticsearchKey + "/version"

	XdbKey             = ResourceTypeXdb + "." + GenericKey
	XdbDatabaseVersion = XdbKey + "/version"

	SnapshotKey         = ResourceTypeSnapshot + "." + GenericKey
	LabelSnapshotStatus = SnapshotKey + "/status"

	PostgresInitSpec      = PostgresKey + "/init"
	ElasticsearchInitSpec = ElasticsearchKey + "/init"
	XdbInitSpec           = XdbKey + "/init"

	PostgresIgnore      = PostgresKey + "/ignore"
	ElasticsearchIgnore = ElasticsearchKey + "/ignore"
	XdbIgnore           = XdbKey + "/ignore"
)

type RuntimeObject interface {
	ResourceCode() string
	ResourceKind() string
	ResourceName() string
	ResourceType() string
}

func (p Postgres) OffshootName() string {
	return p.Name
}

func (p Postgres) OffshootLabels() map[string]string {
	return map[string]string{
		LabelDatabaseName: p.Name,
		LabelDatabaseKind: ResourceKindPostgres,
	}
}

func (p Postgres) StatefulSetLabels() map[string]string {
	labels := p.OffshootLabels()
	for key, val := range p.Labels {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, PostgresKey+"/") {
			labels[key] = val
		}
	}
	return labels
}

func (p Postgres) StatefulSetAnnotations() map[string]string {
	annotations := make(map[string]string)
	for key, val := range p.Annotations {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, PostgresKey+"/") {
			annotations[key] = val
		}
	}
	annotations[PostgresDatabaseVersion] = string(p.Spec.Version)
	return annotations
}

func (p Postgres) ResourceCode() string {
	return ResourceCodePostgres
}

func (p Postgres) ResourceKind() string {
	return ResourceKindPostgres
}

func (p Postgres) ResourceName() string {
	return ResourceNamePostgres
}

func (p Postgres) ResourceType() string {
	return ResourceTypePostgres
}

func (e Elasticsearch) OffshootName() string {
	return e.Name
}

func (e Elasticsearch) OffshootLabels() map[string]string {
	return map[string]string{
		LabelDatabaseKind: ResourceKindElasticsearch,
		LabelDatabaseName: e.Name,
	}
}

func (e Elasticsearch) StatefulSetLabels() map[string]string {
	labels := e.OffshootLabels()
	for key, val := range e.Labels {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, ElasticsearchKey+"/") {
			labels[key] = val
		}
	}
	return labels
}

func (e Elasticsearch) StatefulSetAnnotations() map[string]string {
	annotations := make(map[string]string)
	for key, val := range e.Annotations {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, ElasticsearchKey+"/") {
			annotations[key] = val
		}
	}
	annotations[ElasticsearchDatabaseVersion] = string(e.Spec.Version)
	return annotations
}

func (e Elasticsearch) ResourceCode() string {
	return ResourceCodeElasticsearch
}

func (e Elasticsearch) ResourceKind() string {
	return ResourceKindElasticsearch
}

func (e Elasticsearch) ResourceName() string {
	return ResourceNameElasticsearch
}

func (e Elasticsearch) ResourceType() string {
	return ResourceTypeElasticsearch
}

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

func (d DormantDatabase) OffshootName() string {
	return d.Name
}

func (d DormantDatabase) ResourceCode() string {
	return ResourceCodeDormantDatabase
}

func (d DormantDatabase) ResourceKind() string {
	return ResourceKindDormantDatabase
}

func (d DormantDatabase) ResourceName() string {
	return ResourceNameDormantDatabase
}

func (d DormantDatabase) ResourceType() string {
	return ResourceTypeDormantDatabase
}

func (s Snapshot) OffshootName() string {
	return s.Name
}

func (s Snapshot) Location() (string, error) {
	spec := s.Spec.SnapshotStorageSpec
	if spec.S3 != nil {
		return filepath.Join(spec.S3.Prefix, DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	} else if spec.GCS != nil {
		return filepath.Join(spec.GCS.Prefix, DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	} else if spec.Azure != nil {
		return filepath.Join(spec.Azure.Prefix, DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	} else if spec.Local != nil {
		return filepath.Join(DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	} else if spec.Swift != nil {
		return filepath.Join(spec.Swift.Prefix, DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	}
	return "", errors.New("No storage provider is configured.")
}

func (s Snapshot) ResourceCode() string {
	return ResourceCodeSnapshot
}

func (s Snapshot) ResourceKind() string {
	return ResourceKindSnapshot
}

func (s Snapshot) ResourceName() string {
	return ResourceNameSnapshot
}

func (s Snapshot) ResourceType() string {
	return ResourceTypeSnapshot
}

func (s SnapshotStorageSpec) Container() (string, error) {
	if s.S3 != nil {
		return s.S3.Bucket, nil
	} else if s.GCS != nil {
		return s.GCS.Bucket, nil
	} else if s.Azure != nil {
		return s.Azure.Container, nil
	} else if s.Local != nil {
		return s.Local.Path, nil
	} else if s.Swift != nil {
		return s.Swift.Container, nil
	}
	return "", errors.New("No storage provider is configured.")
}

func (s SnapshotStorageSpec) Location() (string, error) {
	if s.S3 != nil {
		return "s3:" + s.S3.Bucket, nil
	} else if s.GCS != nil {
		return "gs:" + s.GCS.Bucket, nil
	} else if s.Azure != nil {
		return "azure:" + s.Azure.Container, nil
	} else if s.Local != nil {
		return "local:" + s.Local.Path, nil
	} else if s.Swift != nil {
		return "swift:" + s.Swift.Container, nil
	}
	return "", errors.New("No storage provider is configured.")
}
