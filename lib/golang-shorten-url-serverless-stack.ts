import * as apigatewayv2 from '@aws-cdk/aws-apigatewayv2'
import * as apiIntegration from '@aws-cdk/aws-apigatewayv2-integrations'
import * as cdk from '@aws-cdk/core';
import * as cloudfront from '@aws-cdk/aws-cloudfront';
import * as origins from '@aws-cdk/aws-cloudfront-origins';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as iam from '@aws-cdk/aws-iam';
import * as lambda from '@aws-cdk/aws-lambda';
export class GolangShortenUrlServerlessStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here

    // lambda
    const postRegisterLambda = new lambda.Function(this, 'PostRegisterLambda', {
      functionName: 'post-register-lambda',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(15),
      handler: 'main',
      code: lambda.Code.fromAsset(`lambda/post-register-handler/main.zip`),
      environment: {
        "DYNAMODB_TABLE_NAME": "ShortenUrl",
      },
      initialPolicy: [new iam.PolicyStatement({
        actions: ['dynamodb:PutItem'],
        effect: iam.Effect.ALLOW,
        resources: [`arn:aws:dynamodb:${this.region}:${this.account}:table/ShortenUrl`],
      })],
    });

    const getRedirectLambda = new lambda.Function(this, 'GetRedirectLambda', {
      functionName: 'get-redirect-lambda',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(15),
      handler: 'main',
      code: lambda.Code.fromAsset(`lambda/get-redirect-handler/main.zip`),
      environment: {
        "DYNAMODB_TABLE_NAME": "ShortenUrl",
      },
      initialPolicy: [new iam.PolicyStatement({
        actions: ['dynamodb:GetItem'],
        effect: iam.Effect.ALLOW,
        resources: [`arn:aws:dynamodb:${this.region}:${this.account}:table/ShortenUrl`],
      })],
    });

    // apigateway
    const shortenUrlApi = new apigatewayv2.HttpApi(this, 'ShortenUrlApi', {
      apiName: 'shorten-url-api',
    })

    shortenUrlApi.addRoutes({
      path: '/shortened_id/{shortened_id}',
      methods: [apigatewayv2.HttpMethod.GET],
      integration: new apiIntegration.LambdaProxyIntegration({
        handler: getRedirectLambda,
      }),
    });

    shortenUrlApi.addRoutes({
      path: '/register',
      methods: [apigatewayv2.HttpMethod.POST],
      integration: new apiIntegration.LambdaProxyIntegration({
        handler: postRegisterLambda,
      }),
    });

    // dynamoDb
    const dynamoDb = new dynamodb.Table(this, 'ShortenUrl', {
      tableName: "ShortenUrl",
      partitionKey: {
        name: "shortened_id",
        type: dynamodb.AttributeType.STRING,
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: cdk.RemovalPolicy.DESTROY, // NOT FOR PRODUCTION
    });

    // cloudfront
    const cdn = new cloudfront.Distribution(this, 'ShortenUrlCDN', {
      defaultBehavior: {
        origin: new origins.HttpOrigin(shortenUrlApi.url!.replace("https://", "").replace("/", "")),
        cachedMethods: cloudfront.CachedMethods.CACHE_GET_HEAD,
        allowedMethods: cloudfront.AllowedMethods.ALLOW_ALL,
        viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.ALLOW_ALL,
        cachePolicy: new cloudfront.CachePolicy(this, 'ShortenUrlDistributionCachePolicy', {
          defaultTtl: cdk.Duration.days(30),
          maxTtl: cdk.Duration.days(90),
          minTtl: cdk.Duration.days(0),
        }),
      },
    });

    new cdk.CfnOutput(this, `ShortenApiUrl (Apigateway)`, {
      value: shortenUrlApi.url!,
    });

    new cdk.CfnOutput(this, `ShortenApiCDNDomainName`, {
      value: cdn.distributionDomainName,
    });

  }
}
