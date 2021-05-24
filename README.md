# Terraform Cloud Vars (TFCLOUDVARS)

This simple golang program helps to manage terraform cloud variables.
It allow:

1. Read all variables from a workspace.
2. Load variables - in the format provided at the step 1 - into an existing workspace.

## Requirements

1. [Terraform cloud](https://app.terraform.io/) account with an organization and one or more workspaces.
2. Terraform cloud access [token](https://app.terraform.io/app/settings/tokens).
3. golang 1.15.* installed

## Run from the source code

```bash
> git clone https://github.com/pagopa/tfcloudvars.git
> cd tfcloudvars
> # run tests
> go test ./...
# help
> go run main.go help

Usage of /tmp/go-build1852085667/b001/exe/main:
  -do string
        Operation: [read|load|help] (default "help")
  -file string
        json file with variables to load in a workspace
  -format string
        Output format [json|tfvars] (default "json")
  -token string
        bearer token for authenticatio. If not defined it reads the env variable TF_TOKEN
  -ws string
        Terraform cloud workspace id to read from or to save in.

# set terraform cloud token.
export TF_TOKEN=5i*****......................................*****2Ls

> go run main.go -do read -ws ws-<my ws> > ./vars.json

# edit the file. eg: change values
# load the variables in another workspace.
> go run main.go -do load -ws ws-<new ws> -file ./vars.json
```

## Build

```bash
# build binary file:
> make
```

## Run unite tests

```bash
> make test
```
