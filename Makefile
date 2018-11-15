include .env

# found yolo at: https://azer.bike/journal/a-good-makefile-for-go/

AWS_STACK_NAME ?= $(PROJECT_NAME)-$(ENV)
AWS_BUCKET_NAME ?= $(PROJECT_NAME)

deploy: build awspackage awsdeploy

clean:
	@rm -rf dist
	@mkdir -p dist

build: clean
	@for dir in `ls handler`; do \
		GOOS=linux go build -o dist/$$dir github.com/pulpfree/gdps-fs-sum-dwnld/handler/$$dir; \
	done
	@cp ./config/defaults.yml dist/
	@echo "build successful"

# "create-bucket-configuration" line only if not in us-east-1, or apparently some other regions as well...
configure:
	@ aws s3api create-bucket \
		--bucket $(AWS_BUCKET_NAME) \
		--region $(AWS_REGION) \
		--create-bucket-configuration LocationConstraint=$(AWS_REGION)

# watch: Run given command when code changes. e.g; make watch run="echo 'hey'"
# @yolo -i . -e vendor -e bin -e dist -c $(run)

watch:
	@yolo -i . -e vendor -e dist -c "make build"

run: build
	sam local start-api -n env.json

validate:
	sam validate

awspackage:
	@aws cloudformation package \
   --template-file ${FILE_TEMPLATE} \
   --output-template-file ${FILE_PACKAGE} \
   --s3-bucket $(AWS_BUCKET_NAME) \
   --s3-prefix $(AWS_BUCKET_PREFIX) \
   --profile $(AWS_PROFILE) \
   --region $(AWS_REGION)

awsdeploy:
	@aws cloudformation deploy \
   --template-file ${FILE_PACKAGE} \
   --region $(AWS_REGION) \
   --stack-name $(AWS_STACK_NAME) \
   --capabilities CAPABILITY_IAM \
   --profile $(AWS_PROFILE) \
   --force-upload
#   --parameter-overrides \
#	 	  ParamProjectName=$(PROJECT_NAME) \
#			ParamKeyExpiration=$(EXPIRATION) \
#			ParamENV=$(ENV)

describe:
	@aws cloudformation describe-stacks \
		--region $(AWS_REGION) \
		--stack-name $(AWS_STACK_NAME)

outputs:
	@ make describe \
		| jq -r '.Stacks[0].Outputs'
