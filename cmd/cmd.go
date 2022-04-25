/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"fmt"
	"strings"
	"io/ioutil"
	"log"
	"github.com/ghodss/yaml"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)


var inputFilePath *string
var outputFormat *string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubeclean",
	Short: "Tool that clean Kubernetes Manifest",
	Long: `kubeclean is a tool that clean kuberntes manifest, removing the non-needed data so that the manifest can be re-used again in another context.
Also, kubeclean can be launched as a server, that can receive and clean data in JSON format.
Data that will be cleaned: Default SA, Kubernetes defined MetaData, Status and Tolerations`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { 

		var input []byte 
		var err error 

		if *inputFilePath == "-"{
			stdin := cmd.InOrStdin() 
			input, _ = ioutil.ReadAll(stdin)
		} else {
			input, err = ioutil.ReadFile(*inputFilePath) 
			if err != nil {
				log.Fatal(err)
			}	
		}

		cleanManifestByteSlice := CleanManifest(input)
		fmt.Print(string(cleanManifestByteSlice))
		//ioutil.WriteFile("podCleaned.yaml",cleanManifestByteSlice, 0644)
	},
}


// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubeclean.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	inputFilePath = rootCmd.Flags().StringP("file","f", "-","path to the file that need to be cleand, default to stdin")
	outputFormat = rootCmd.Flags().StringP("output","o","yaml","output format, is either json or yaml, default to yaml")

}


func CleanManifest(manifest []byte) ([]byte) {
	var cleanManifestString string

	jsonBytesOutput , _ := yaml.YAMLToJSON(manifest)
	jsonStringOutput := string(jsonBytesOutput)

	cleanManifestString, _ = CleanStatus(jsonStringOutput)
	
	cleanManifestString,_ = CleanDefaultSA(cleanManifestString)
	cleanManifestString, _= CleanMedata(cleanManifestString)
	cleanManifestString, _ = CleanTolerations(cleanManifestString)
	jsonByteOutput := []byte(cleanManifestString)


	if *outputFormat == "json" {

		return jsonByteOutput
		//ioutil.WriteFile("podCleaned.json",jsonByteOutput, 0644)
	}else{

		yamlByteOutput, err := yaml.JSONToYAML(jsonByteOutput)
		if err != nil{
			log.Panic("Cannot convert JSON to YAML: %v",err)
		}
	
		return yamlByteOutput

		//ioutil.WriteFile("podCleaned.yaml",yamlOutput, 0644)
	}
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
