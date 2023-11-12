#!/bin/bash -xe
# title           :deploy.sh
# description     :Deploy Go binary to AWS Lambda via AWS CLI
# date            :2023-10-30
# version         :0.1.0
# usage           :bash depoy.sh || ./depoy.sh
# assumption      :The AWS function already exists. Zip is downloaded, Bash is available
#                 :Sorry Windows. AWS CLI credentials are available, e.g., ~/.aws/credentials
# variables
#                 :REGION - The AWS Region the lambda function resides
#                 :FUNC_NAME  - Name of the AWS lambda function to update
#==============================================================================

# Delete an previous build if exists
rm -f -- bootstrap function.zip 

REGION=us-east-1
FUNC_NAME=goodidea
RELEASE_PATH=$PWD # Where to start looking for files
ARCHIVE=${LAMBDA_ARCHIVE:-function.zip} # the name of the zip file to deploy
BINARY_PATH=${RELEASE_PATH}/lambda

deploy () {
    npm run build
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
        -o bootstrap \
        -ldflags="-s -w" \
        -tags lambda.norpc \
        $BINARY_PATH

    zip -r function.zip bootstrap static templates

    aws lambda update-function-code \
        --region $REGION \
        --function-name $FUNC_NAME \
        --zip-file fileb://$RELEASE_PATH/$ARCHIVE
}

deploy
