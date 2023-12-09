package delete

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"stressTest/defs"
	"stressTest/util"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
)

func Delete() {
	auth := "Bearer " + defs.Token
	req, err := http.NewRequest("DELETE", "https://192.168.12.127:6443/api/v1/namespaces/myx-test/configmaps", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", auth)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(string(body))
}
func ClearAll(ns, labelSelector, auth string) {
	for _, v := range defs.Reslist {
		ClearPost(v, ns, labelSelector, auth)
	}
}
func ClearPost(res, ns, labelSelector, auth string) {
	_, _, request := util.GetBasic(res, ns)
	// req, err := http.NewRequest("DELETE", request+"?labelSelector="+labelSelector, nil)
	req, err := http.NewRequest("DELETE", request, nil)
	log.Println("DELETE", request, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", auth)
	client := &http.Client{
		Timeout: time.Minute * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	_, err = client.Do(req)
	if err != nil {
		log.Fatal("req err : ", err)
	}
	// defer resp.Body.Close()
	// log.Println("delete post : ", resp.StatusCode, resp.Status)
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("delete : ", string(body))
}
func DeleteNameSpace(res, ns, labelSelector, auth string) {
	_, _, request := util.GetBasic(res, ns)
	req, err := http.NewRequest("GET", request+"?labelSelector="+labelSelector, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", auth)

	client := &http.Client{
		Timeout: time.Minute * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var namespaceList v1.NamespaceList
	err = json.Unmarshal(respBody, &namespaceList)
	if err != nil {
		panic(err)
	}
	for _, namespace := range namespaceList.Items {
		req, err := http.NewRequest("DELETE", request+"/"+namespace.Name, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Authorization", auth)
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Deleted namespace %s, response: %s\n", namespace.Name, respBody)
	}
}
func DeleteTerminatedNamespace(res, ns, labelSelector, auth string) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	_, _, request := util.GetBasic(res, ns)
	req, err := http.NewRequest("GET", request+"?labelSelector="+labelSelector, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", auth)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var namespaceList v1.NamespaceList
	err = json.Unmarshal(respBody, &namespaceList)
	if err != nil {
		panic(err)
	}
	for _, namespace := range namespaceList.Items {
		if namespace.Status.Phase == "Terminating" {
			payload := strings.NewReader(`{"metadata":{"finalizers":null}}`)
			req, err := http.NewRequest("PATCH", request+"/"+namespace.Name+"/"+"finalize", payload)
			if err != nil {
				panic(err)
			}
			req.Header.Set("Authorization", auth)
			req.Header.Set("Content-Type", "application/strategic-merge-patch+json")
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Finalized namespace %s, response: %s\n", namespace.Name, respBody)
		}
	}
}
