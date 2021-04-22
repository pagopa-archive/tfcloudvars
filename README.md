# Terraform Cloud Vars (TFCLOUDVARS)

This simple golang program helps to manage terraform cloud variables.
It allow:

1. Read all variables in a workspace.
2. Load variables - in the format provided at the step 1 - into an existeing workspace.


## Requirements

1. [Terraform cloud](https://app.terraform.io/) account with an organization and one or more workspace.
2. Terraform cloud access [token](https://app.terraform.io/app/settings/tokens).
3. golang 1.15.* installed

## Run from the source code

```bash
> git clone https://github.com/pagopa/tfcloudvars.git
> cd tfcloudvars
> go run main.go help

Usage of /tmp/go-build148387471/b001/exe/main:
  -do string
        Operation: [read | save|help] (default "help")
  -file string
        json file with variables to load in a workspace
  -token string
        bearer token for authenticatio. If not defined it reads the env variable TF_TOKEN
  -ws string
        Terraform cloud workspace id to read from or to save in.
```
