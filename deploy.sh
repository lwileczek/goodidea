rm -f -- bootstrap function.zip 

# BINARY_PATH - path to the specific endpoint to build
# FUNC_NAME - The name of the lambda function in AWS we are deploying to
RELEASE_PATH=${PROJECT_PATH:-$PWD} # Where to start looking for files
ARCHIVE=${LAMBDA_ARCHIVE:-function.zip} # the name of the zip file to deploy
FUNC_NAME=goodidea
BINARY_PATH=${RELEASE_PATH}/lambda
REGION=us-west-2

deploy () {
    npm run build
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
        -o bootstrap \
        -ldflags="-s -w" \
        -tags lambda.norpc \
        $BINARY_PATH

    zip -r function.zip bootstrap kodata

    aws lambda update-function-code \
        --region $REGION \
        --function-name $FUNC_NAME \
        --zip-file fileb://$RELEASE_PATH/$ARCHIVE
}

deploy
