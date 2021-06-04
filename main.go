package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/uolter/cptfcvars/tfcloud"
)

var (
	do        string
	fileName  string
	workspace string
	token     string
	format    string
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
	flag.StringVar(&token, "token", LookupEnvOrString("TF_TOKEN", ""), "bearer token for authenticatio. If not defined it reads the env variable TF_TOKEN or the credeintial storage file: credentials.tfrc.json")
	flag.StringVar(&format, "format", "json", "Output format [json|tfvars]")

	flag.Parse()
}

func main() {

	if workspace == "" {
		log.Println("[INFO] workspace required")
		Usage()
		os.Exit(0)
	}

	if token == "" {

		c := tfcloud.TfConfig{}
		dirname, err := os.UserHomeDir()

		if err != nil {
			log.Fatal(err)
			Usage()
			os.Exit(0)
		}

		err = c.Read(filepath.Join(dirname, ".terraform.d", "credentials.tfrc.json"))

		if err == nil {
			token = c.Credentials.App_terraform_io.Token
		} else {
			log.Println("[INFO] token required")
			Usage()
			os.Exit(0)
		}
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

	var err error
	var j string

	if err := t.Get(workspace, token); err != nil {
		log.Printf("[ERROR]: %t", err)
		os.Exit(1)
	}

	switch format {
	case "json":
		j, err = t.Json(true)
	case "tfvars":
		j, err = t.ToTfVars(true)
	default:
		log.Println("[INFO] wrong format value.")
		Usage()
	}

	if err != nil {
		log.Printf("[ERROR]: %t", err)
		os.Exit(1)
	}
	fmt.Println(j)

}

func save() {

	t := tfcloud.TerraformVars{}

	var content string

	// Read from the standard input
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		content += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	err := t.Load(content)

	if err != nil {
		log.Println(fmt.Sprintf("[ERROR] %s", err.Error()))
		os.Exit(1)
	}

	log.Println("Load into workspace")

	err = t.Post(workspace, token)

	if err != nil {
		log.Println(fmt.Sprintf("[ERROR] %s", err.Error()))
		os.Exit(1)
	}

}
