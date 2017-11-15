package v1alpha1

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

	MySQLKey             = ResourceTypeMySQL + "." + GenericKey
	MySQLDatabaseVersion = MySQLKey + "/version"

	MongoDBKey             = ResourceTypeMongoDB + "." + GenericKey
	MongoDBDatabaseVersion = MongoDBKey + "/version"

	XdbKey             = ResourceTypeXdb + "." + GenericKey
	XdbDatabaseVersion = XdbKey + "/version"

	SnapshotKey         = ResourceTypeSnapshot + "." + GenericKey
	LabelSnapshotStatus = SnapshotKey + "/status"

	PostgresInitSpec      = PostgresKey + "/init"
	ElasticsearchInitSpec = ElasticsearchKey + "/init"
	MySQLInitSpec         = MySQLKey + "/init"
	MongoDBInitSpec       = MongoDBKey + "/init"
	XdbInitSpec           = XdbKey + "/init"

	PostgresIgnore      = PostgresKey + "/ignore"
	ElasticsearchIgnore = ElasticsearchKey + "/ignore"
	MySQLIgnore         = MySQLKey + "/ignore"
	MongoDBIgnore       = MongoDBKey + "/ignore"
	XdbIgnore           = XdbKey + "/ignore"

	AgentCoreosPrometheus        = "coreos-prometheus-operator"
	PrometheusExporterPortNumber = 56790
	PrometheusExporterPortName   = "http"
)
