/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"

	"net/http"

	"github.com/spf13/cobra"
	"log"
)

var serverPort *string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		http.HandleFunc("/clean",CleanServer)
		err := http.ListenAndServe(fmt.Sprintf(":"+ *serverPort),nil)
		if err != nil {
			fmt.Println(err)
		}
	},
}


func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverPort = serveCmd.Flags().StringP("port","p","8080","define the port on wich the server will listen")
}

func CleanServer(w http.ResponseWriter, req *http.Request){
	
	//req.ParseForm()
	//manifest := req.Form
	//fmt.Println(manifest)
	manifest, err := ioutil.ReadAll(req.Body)
	if err!=nil {log.Fatal("request",err)}
	
	CleanManifestByteSlice := CleanManifest(manifest)
	fmt.Println(string(CleanManifestByteSlice))
}