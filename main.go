package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/uolter/cptfcvars/tfcloud"
)

var (
	do        string
	fileName  string
	workspace string
	token     string
)

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

	flag.PrintDefaults()
}

func init() {
	flag.StringVar(&do, "do", "help", "Operation: [read|load|help]")
	flag.StringVar(&workspace, "ws", "", "Terraform cloud workspace id to read from or to save in.")
	flag.StringVar(&fileName, "file", "", "json file with variables to load in a workspace")
	flag.StringVar(&token, "token", LookupEnvOrString("TF_TOKEN", ""), "bearer token for authenticatio. If not defined it reads the env variable TF_TOKEN")

	flag.Parse()
}

func main() {

	if workspace == "" {
		log.Println("[INFO] workspace required")
		Usage()
		os.Exit(0)
	}

	if token == "" {
		log.Println("[INFO] token required")
		Usage()
		os.Exit(0)
	}

	switch do {
	case "read":
		read()
	case "load":
		save()
	default:
		Usage()
	}
}

func read() {

	t := tfcloud.TerraformVars{}

	if err := t.Get(workspace, token); err != nil {
		log.Printf("[ERROR]: %t", err)
		os.Exit(1)
	}

	j, err := t.Json(true)

	if err != nil {
		log.Printf("[ERROR]: %t", err)
		os.Exit(1)
	}

	fmt.Println(j)
}

func save() {

	if fileName == "" {
		log.Println("[INFO] file name required")
		Usage()
		os.Exit(0)
	}

	t := tfcloud.TerraformVars{}

	err := t.Load(fileName)

	if err != nil {
		log.Println(fmt.Sprintf("[ERROR] %s", err.Error()))
		os.Exit(1)
	}

	err = t.Post(workspace, token)

	if err != nil {
		log.Println(fmt.Sprintf("[ERROR] %s", err.Error()))
		os.Exit(1)
	}

}
