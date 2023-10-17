# Trapa API

TODO: put lots of stuff in here. Figure out how to write a proper API in go. 

## Setup
1. Install Go 1.21.3: https://go.dev/dl/
1. Install the [gcloud cli](https://cloud.google.com/sdk/docs/install)
2. Login with your account `gcloud login`. This allows you to execute actions on the command line
3. Login with your account to create "application default credentials". These credentials can be automatically detected by the Go google client libraries to authenticate when running this project locally: `gcloud auth application-default login`

## Running 

### VSCode
1. Install [Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)
2. Run `Trapa API (Dev)` from the VSCode Run and Debug menu

### Docker
1. Find the platform-specific location of the `application_default_credentials.json` file createdin the setup steps
2. Copy and rename `compose.override.example.yaml` to `compose.override.yaml`
3. Edit `compose.override.yaml` so that the volume source path points to `application_default_credentials.json` file
4. Run `docker compose up --build`