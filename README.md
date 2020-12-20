# What is this?

A Serverless shorten URL solution with CDK, application is written in Golang.

## Prerequisites
Go 1.14 [installed](https://golang.org/doc/install)
npm installed
cdk [installed](https://docs.aws.amazon.com/cdk/latest/guide/work-with-cdk-typescript.html)
aws access with S3, Dynamodb, Lambda, Apigateway, Cloudfront, CloudFormation

## How to use?
 * Update AWS_PROFILE and AWS_DEFAULT_REGION in Makefile
 * run `cdk bootstrap` if this is the first time to run CDK in the region
 * run `make -f Makefile deploy` for shorten url service deploy, it might take about 5 minutes. Endpoint of ApiGateway and CloudFront will be displayed on the console right after deploy.

## POST Register API:
Register API respond shortened ID in body directly
```
curl -X POST \
  https://${DOMAIN}/register \
  -H 'content-type: application/json' \
  -d '{"register_url":"SAMPLE_URL"}'
```
```
Respond in JSON format:
{
    "shortened_id": "${shortened_id}"
}
```
## Get Redirect API:
```
curl -X GET \
  https://${DOMAIN}/shortened_id/${shortened_id} \
  -H 'content-type: application/json' \
```
```
- Redirect directly if shortened_id exists
- Return 403 if shortened_id doesn't exists
```

## Command:
 * `make -f Makefile deploy`      deploy this stack to your default AWS account/region
 * `make -f Makefile diff`        compare deployed stack with current state
 * `make -f Makefile synth`       emits the synthesized CloudFormation template
 * `make -f Makefile destroy`     destroy this stack along with resource defined on this stack
