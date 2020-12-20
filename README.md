# What is this?

A serverless shorten URL solution with CDK, and the application is written in Golang.

## Design
In my point of view, it’s a read-heavy service and the content is quite static since it won’t be changed after created.
So I use a CDN with long cache time to lower the request time from user to server and optimize the user experience.
And because it is a new service so we don’t have information about peak hours and count. It might be a waste if we build and maintain a server with too many resource allocated. And it might be potential threat if we allocate resource that is not enough to handle the requests. So I choose API Gateway, Lambda and Dynamo DB for the shorten URL solution. We only have to pay based on the request instead of the whole instance. So we can focus on the application and business logic first, and consider about server migration if needed. 

![Screenshot](get-redirect-api.png)
![Screenshot](post-register-api.png)

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
  -H 'content-type: application/json'
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
