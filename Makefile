GET_HANDLER_FILE_PATH=lambda/get-redirect-handler
GET_HANDLER_FILE_NAME=cmd/handler/main.go
POST_HANDLER_FILE_PATH=lambda/post-register-handler
POST_HANDLER_FILE_NAME=cmd/handler/main.go
AWS_PROFILE=brad-personal
AWS_DEFAULT_REGION=ap-northeast-1

build_go_get_linux:
	@echo Building Golang File For Linux...
	cd $(GET_HANDLER_FILE_PATH) && GOOS=linux go build $(GET_HANDLER_FILE_NAME)

compress_go_get:
	@echo Compressing Golang File...
	cd $(GET_HANDLER_FILE_PATH) && zip -r main.zip main

build_go_post_linux:
	@echo Building Golang File For Linux...
	cd $(POST_HANDLER_FILE_PATH) && GOOS=linux go build $(POST_HANDLER_FILE_NAME)

compress_go_post:
	@echo Compressing Golang File...
	cd $(POST_HANDLER_FILE_PATH) && zip -r main.zip main

generate_lambda_files: build_go_get_linux compress_go_get build_go_post_linux compress_go_post
	@echo Complete Generate Golang File For Lambda

bootstrap:
	@echo node_modules/aws-cdk/bin/cdk execute bootstrap
	AWS_PROFILE=$(AWS_PROFILE) \
	AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) \
		node_modules/aws-cdk/bin/cdk bootstrap \
			--toolkit-stack-name=CDKToolkit

deploy: generate_lambda_files
	@echo Deploying...
	AWS_PROFILE=$(AWS_PROFILE) \
	AWS_DEFAULT_REGION=$(AWS_PROFILE) \
		node_modules/aws-cdk/bin/cdk deploy GolangShortenUrlServerlessStack \
			--app="npx ts-node ./bin/golang-shorten-url-serverless.ts" \
			--toolkit-stack-name=CDKToolkit

diff:
	@echo Diff The Stack
	AWS_PROFILE=$(AWS_PROFILE) \
	AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) \
	node_modules/aws-cdk/bin/cdk diff GolangShortenUrlServerlessStack \
			--app="npx ts-node ./bin/golang-shorten-url-serverless.ts" \
			--toolkit-stack-name=CDKToolkit

synth:
	@echo Diff The Stack
	AWS_PROFILE=$(AWS_PROFILE) \
	AWS_DEFAULT_REGION=$(AWS_PROFILE) \
		node_modules/aws-cdk/bin/cdk synth GolangShortenUrlServerlessStack \
			--app="npx ts-node ./bin/golang-shorten-url-serverless.ts" \
			--toolkit-stack-name=CDKToolkit

destroy:
	@echo Destroying...
	AWS_PROFILE=$(AWS_PROFILE) \
	AWS_DEFAULT_REGION=$(AWS_PROFILE) \
		node_modules/aws-cdk/bin/cdk destroy GolangShortenUrlServerlessStack \
			--app="npx ts-node ./bin/golang-shorten-url-serverless.ts" \
			--toolkit-stack-name=CDKToolkit
