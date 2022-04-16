package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"github.com/ghodss/yaml"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strings"
	"github.com/spf13/cobra"
)

func main () {

	var cleanManifest string

	yamlFileInput, err := ioutil.ReadFile("pod-output.yaml")
	if err != nil {
		log.Fatal(err)
	}

	
	jsonBytesOutput , _ := yaml.YAMLToJSON(yamlFileInput)
	jsonStringOutput := string(jsonBytesOutput)


	cleanManifest, _ = CleanStatus(jsonStringOutput)
	if err != nil{
		log.Panic("Error while cleaning manifest: %v", err)
	}
	
	cleanManifest,_ = CleanDefaultSA(cleanManifest)
	cleanManifest, _= CleanMedata(cleanManifest)
	cleanManifest, _ = CleanTolerations(cleanManifest)
	jsonByteOutput := []byte(cleanManifest)

	yamlOutput, err := yaml.JSONToYAML(jsonByteOutput)
	if err != nil{
		log.Panic("Cannot convert JSON to YAML: %v",err)
	}

	ioutil.WriteFile("podClaned.yaml",yamlOutput, 0644)

}

//CleanStatus delete the status field of the given manifest
func CleanStatus(manifest string) (string,error) {
	return sjson.Delete(manifest,"status")
}

//CleanDefaultSA
func CleanDefaultSA(manifest string) (string,error) {
	
	var err error 
	//Iterate over the volumes
	for i, volume := range gjson.Get(manifest, "spec.volumes").Array() {
		volumeName := volume.Get("name").String()
		if strings.HasPrefix(volumeName, "kube-api-") {
			manifest, err = sjson.Delete(manifest,fmt.Sprintf("spec.volumes.%d",i))	
			if err != nil {
				continue
			}
		}
	}
	//Iterate over container and VolumeMounts
	for ci, containers := range gjson.Get(manifest, "spec.containers").Array() {
		for vmi, volumeMount := range containers.Get("volumeMounts").Array() {
			volumeMountName := volumeMount.Get("name").String()
			if strings.HasPrefix(volumeMountName, "kube-api-") {
				manifest, err = sjson.Delete(manifest, fmt.Sprintf("spec.containers.%d.volumeMounts.%d", ci, vmi))
				if err != nil {
					continue
				}
			}
		}
	}
	return manifest,nil
}
 

//CleanMedata Delete non needed metadata in the manifests
func CleanMedata(manifest string) (string, error) {

	var err error
	nonNeededMetadata := []string{"creationTimestamp","resourceVersion","uid"}

	for _ , metadata := range nonNeededMetadata {
		manifest, err = sjson.Delete(manifest,fmt.Sprintf("metadata.%s", metadata))
		if err != nil {
			continue
		}
	}

	//also clean last applied configuration 
	manifest, _ = sjson.Delete(manifest,  `metadata.annotations.kubectl\.kubernetes\.io/last-applied-configuration`)
	return manifest, err
}

//CleanTolerations will clean tolerations
func CleanTolerations(manifest string) (string, error) {
	var err error 
	manifest, err = sjson.Delete(manifest, fmt.Sprintf("spec.tolerations"))
	return manifest, err
}