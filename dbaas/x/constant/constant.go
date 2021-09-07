package constant

import "fmt"

const (
	y = `apiVersion: %v
kind: %v
metadata:
  namespace: %v
  name: %v`
	MysqlV1 = "mysql.presslabs.org/v1alpha1"
	V1      = "v1"

	KindMysqlBackup  = "MysqlBackup"
	KindMysqlCluster = "MysqlCluster"
	KindSecret       = "Secret"
)

func Yaml(apiVersion, kind, namespace, name string) string {
	return fmt.Sprintf(y, apiVersion, kind, namespace, name)
}

func MysqlBackupYaml(namespace, name string) string {
	return Yaml(MysqlV1, KindMysqlBackup, namespace, name)
}

func MysqlClusterYaml(namespace, name string) string {
	return Yaml(MysqlV1, KindMysqlCluster, namespace, name)
}

func SecretYaml(namespace, name string) string {
	return Yaml(V1, KindSecret, namespace, name)
}
