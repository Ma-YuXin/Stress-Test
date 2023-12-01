package util

import (
	"log"
	"strconv"
	"stressTest/defs"
	"strings"
)

func GetBasic(res, namespace string) (string, string, string) {
	kind := Res2kind(res)
	if kind == "" {
		log.Println("invalid resource type")
	}
	uri := Res2uri(res)
	request := ""
	namespaced, api_version := Kind2ns_and_version(kind)
	if api_version == "v1" {
		if namespaced {
			request = defs.Endpoint + "api/" + api_version + "/namespaces/" + namespace + "/" + uri
		} else {
			request = defs.Endpoint + "api/" + api_version + "/" + uri
		}
	} else {
		request = defs.Endpoint + "apis/" + api_version + "/namespaces/" + namespace + "/" + uri
	}
	return kind, api_version, request
}
func GetPostDataAndUrl(res, namespace string, antNum, num, id int) (data []byte, request string) {
	kind, api_version, request := GetBasic(res, namespace)
	// fmt.Println(kind, api_version, request)
	body := ""
	annotation := GetAnnotations(antNum)
	if kind == "Namespace" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "labels": {"env":"test"}` + annotation + `}}`
	} else if kind == "Node" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"unschedulable": true}}`
	} else if kind == "PersistentVolume" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"accessModes": ["ReadWriteOnce"], "capacity": {"storage": "100Ki"}, "hostPath": {"path": "/root/data"}}}`
	} else if kind == "ConfigMap" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "data": {"test-data":""}}`
	} else if kind == "Endpoints" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}}`
	} else if kind == "LimitRange" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"limits": [{"type":"Container"}]}}`
	} else if kind == "PersistentVolumeClaim" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"accessModes": ["ReadWriteOnce"], "resources": {"requests": {"storage": "200Ki"}}}}`
	} else if kind == "Pod" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"containers": [{"name": "test", "image": "nginx:1.17"}]}}`
	} else if kind == "PodTemplate" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "template": {"metadata": {"name": "pod-template"}, "spec": {"containers": [{"name": "test", "image": "nginx:1.17"}]}}}`
	} else if kind == "ReplicationController" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"replicas": 0, "selector": {"app": "test"}, "template": {"metadata": {"name": "test", "labels": {"app": "test"}}, "spec": {"containers": [{"name": "test", "image": "nginx:1.17"}]}}}}`
	} else if kind == "ResourceQuota" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}}`
	} else if kind == "Secret" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "data": {"test-secret":""}, "type": "Opaque"}`
	} else if kind == "ServiceAccount" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}}`
	} else if kind == "Service" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `},"spec": {"ports": [{"port": 80}]}}`
	} else if kind == "ControllerRevision" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "revision": 0, "data": ""}`
	} else if kind == "DaemonSet" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"selector": {"matchLabels": {"app": "test"}}, "template": {"metadata": {"name": "test", "labels": {"app": "test"}}, "spec": {"containers": [{"name": "test", "image": "nginx:1.17"}]}}}}`
	} else if kind == "Deployment" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"replicas": 0, "selector": {"matchLabels": {"app": "test"}}, "template":{"metadata": {"name": "test", "labels": {"app": "test"}}, "spec": {"containers": [{"name": "test","image": "nginx:1.17"}]}}}}`
	} else if kind == "ReplicaSet" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"replicas": 0, "selector": {"matchLabels": {"app": "test"}}, "template":{"metadata": {"name": "test", "labels": {"app": "test"}}, "spec": {"containers": [{"name": "test","image": "nginx:1.17"}]}}}}`
	} else if kind == "StatefulSet" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"replicas": 0, "selector": {"matchLabels": {"app": "test"}}, "serviceName": "", "template": {"metadata": {"name": "test", "labels": {"app": "test"}}, "spec": {"containers": [{"name": "test", "image": "nginx:1.17"}]}}}}`
	} else if kind == "CronJob" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"schedule": "0 */1 * * *", "suspend": true, "jobTemplate": {"spec": {"template": {"spec": {"restartPolicy": "Never", "containers": [{"name": "test", "image": "busybox:1.30", "command": ["bin/sh", "-c", "sleep 10"]}]}}}}}}`
	} else if kind == "Job" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "test-` + strings.ToLower(kind) + `-` + strconv.Itoa(num) + `-` + strconv.Itoa(id) + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"template": {"spec": {"restartPolicy": "Never", "containers": [{"name": "test", "image": "busybox:1.30", "command": ["bin/sh", "-c", "sleep 10"]}]}}}}`
	}
	return []byte(body), request
}
func GetAnnotations(num int) string {
	res := `, "annotations": {`
	for i := 0; i < num; i++ {
		if i < 10 {
			res += `"key-00` + strconv.Itoa(i) + `": "value-00` + strconv.Itoa(i) + `"`
		} else if i < 100 {
			res += `"key-0` + strconv.Itoa(i) + `": "value-0` + strconv.Itoa(i) + `"`
		} else {
			res += `"key-` + strconv.Itoa(i) + `": "value-` + strconv.Itoa(i) + `"`
		}
		if i != num-1 {
			res += ", "
		}
	}

	res += "}"
	return res
}
func GetPatchAnnotations(num int) string {
	res := `{"metadata": {"annotations": {`
	for i := 0; i < num; i++ {
		if i < 10 {
			res += `"key-00` + strconv.Itoa(i) + `": "value-00` + strconv.Itoa(i) + `"`
		} else if i < 100 {
			res += `"key-0` + strconv.Itoa(i) + `": "value-0` + strconv.Itoa(i) + `"`
		} else {
			res += `"key-` + strconv.Itoa(i) + `": "value-` + strconv.Itoa(i) + `"`
		}
		if i != num-1 {
			res += ", "
		}
	}
	res += "}}}"
	return res
}

func GetPutDataAndUrl(res, namespace, resName string, antNum int) (data []byte, request string) {
	kind, api_version, request := GetBasic(res, namespace)
	// fmt.Println(kind, api_version, request)
	body := ""
	annotation := GetAnnotations(antNum)
	if kind == "Namespace" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "labels": {"env":"test"}` + annotation + `}}`
	} else if kind == "Node" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"unschedulable": true}}`
	} else if kind == "PersistentVolume" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"accessModes": ["ReadWriteOnce"], "capacity": {"storage": "100Ki"}, "hostPath": {"path": "/root/data"}}}`
	} else if kind == "ConfigMap" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "data": {"test-data":""}}`
	} else if kind == "Endpoints" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}}`
	} else if kind == "LimitRange" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"limits": [{"type":"Container"}]}}`
	} else if kind == "PersistentVolumeClaim" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"accessModes": ["ReadWriteOnce"], "resources": {"requests": {"storage": "200Ki"}}}}`
	} else if kind == "Pod" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"containers": [{"name": "test", "image": "nginx:1.17"}]}}`
	} else if kind == "PodTemplate" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "template": {"metadata": {"name": "pod-template"}, "spec": {"containers": [{"name": "test", "image": "nginx:1.17"}]}}}`
	} else if kind == "ReplicationController" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"replicas": 0, "selector": {"app": "test"}, "template": {"metadata": {"name": "test", "labels": {"app": "test"}}, "spec": {"containers": [{"name": "test", "image": "nginx:1.17"}]}}}}`
	} else if kind == "ResourceQuota" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}}`
	} else if kind == "Secret" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "data": {"test-secret":""}, "type": "Opaque"}`
	} else if kind == "ServiceAccount" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}}`
	} else if kind == "Service" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `},"spec": {"ports": [{"port": 80}]}}`
	} else if kind == "ControllerRevision" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "revision": 0, "data": ""}`
	} else if kind == "DaemonSet" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"selector": {"matchLabels": {"app": "test"}}, "template": {"metadata": {"name": "test", "labels": {"app": "test"}}, "spec": {"containers": [{"name": "test", "image": "nginx:1.17"}]}}}}`
	} else if kind == "Deployment" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"replicas": 0, "selector": {"matchLabels": {"app": "test"}}, "template":{"metadata": {"name": "test", "labels": {"app": "test"}}, "spec": {"containers": [{"name": "test","image": "nginx:1.17"}]}}}}`
	} else if kind == "ReplicaSet" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"replicas": 0, "selector": {"matchLabels": {"app": "test"}}, "template":{"metadata": {"name": "test", "labels": {"app": "test"}}, "spec": {"containers": [{"name": "test","image": "nginx:1.17"}]}}}}`
	} else if kind == "StatefulSet" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"replicas": 0, "selector": {"matchLabels": {"app": "test"}}, "serviceName": "", "template": {"metadata": {"name": "test", "labels": {"app": "test"}}, "spec": {"containers": [{"name": "test", "image": "nginx:1.17"}]}}}}`
	} else if kind == "CronJob" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"schedule": "0 */1 * * *", "suspend": true, "jobTemplate": {"spec": {"template": {"spec": {"restartPolicy": "Never", "containers": [{"name": "test", "image": "busybox:1.30", "command": ["bin/sh", "-c", "sleep 10"]}]}}}}}}`
	} else if kind == "Job" {
		body = `{"apiVersion": "` + api_version + `", "kind": "` + kind + `", "metadata": {"name": "` + resName + `", "namespace": "` + namespace + `", "labels": {"env":"test"}` + annotation + `}, "spec": {"template": {"spec": {"restartPolicy": "Never", "containers": [{"name": "test", "image": "busybox:1.30", "command": ["bin/sh", "-c", "sleep 10"]}]}}}}`
	}
	return []byte(body), request
}
