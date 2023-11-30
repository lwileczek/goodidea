#!/bin/bash
#Required DB Env
export DATABASE_URL="postgres://fairy:goodidea@localhost:5555/tasks"

# Required for building contianer
export KO_DATA_PATH="./kodata/"
export KO_DOCKER_REPO="alomio/goodidea"

## AWS SDK environment variables which may be used to explicitly authenticate to 
## the object storage
#export AWS_ACCESS_KEY_ID=AKIDZXCVBNM
#export AWS_SECRET_ACCESS_KEY=ASDDFFGHJKLZXCVBNMPOIUYTRE
#export AWS_SESSION_TOKEN=
#export AWS_DEFAULT_REGION=us-west-2


##The name of the bucket to store image in. *Required to use object storage
#export AWS_BUCKET=
## If images will be served by a CDN after being uploaded to object storage
## provide the URL here
#export AWS_IMAGE_CDN=https://mycdn.example.com
## Endpoint to use if not using AWS S3 for object storage
#export AWS_ENDPOINT=


# Local Storage Directory override
export LOCAL_DIR=static/img
