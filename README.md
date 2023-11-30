# Good Idea
Good Idea is a slightly fancier TODO app where anonymous users may
submit requests and vote on their importance.
Users do **not** need to login or create an account to make a request or cast votes.

### Requirements
 - Postgresql Database
 - environment variable `DATABASE_URL` to point the server at the database

#### Dev Requirements
 - go 1.21+
 - node.js 18+

## Build
 1. Build and minify CSS and JS with: `npm install && npm run build`
 2. Build the server binary with `go build ./app/main.go`

 To make it easy, use [make](https://www.gnu.org/software/make/) to build

 ```bash
 make build
 ```

 which will populate the [static](./static) directory and build a binary called `server`

### Configure
By default, the server will attempt to user local storage for images under `static/img`.
This can be changed by setting `LOCAL_DIR` within the environment. 

If you prefer to use object storage instead, you can set the environment variables to provided
the credentials needed to authenticate: 

| variable | description |
|:---|:---|
|AWS\_ACCESS\_KEY\_ID | The access key provided by your object storage provider for the S3 compatible sdk |
|AWS\_SECRET\_ACCESS\_KEY | The secret key for authenticating to object storage |
|AWS\_SESSION\_TOKEN | The token for the current session which can be used |
|AWS\_DEFAULT\_REGION | The region the bucket is in|

[aws sdk env variables](https://docs.aws.amazon.com/sdkref/latest/guide/environment-variables.html)

In addition, the following variables may be used as well. 

| variable | description |
|:---|:---|
|AWS\_BUCKET | The name of the bucket to use within the object store |
|AWS\_IMAGE\_CDN | The URL (including the protocol, e.g. https) of the CDN which serves the images out of the bucket |
|AWS\_ENDPOINT | The URI for the object storage to try and connect to if not AWS S3 |

The system checks for the environment variables `AWS_BUCKET` to decide if it should be using object storage or not
and thus, the rest are optional.

# Status
> **warning**
> This project is in _beta_

## Deployment
For deployment information, please see the project [wiki](https://github.com/lwileczek/goodidea/wiki)

