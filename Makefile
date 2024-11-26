
all:

codegen:
	@go generate ./...
	@openapi-generator-cli generate -i docs/interpret-api.yaml -g python-flask -o interpret --global-property models,apis,apiTests=false,modelTests=false,supportingFiles
	@openapi-generator-cli generate -i docs/formalize-api.yaml -g python-flask -o formalize_convert --global-property models,apis,apiTests=false,modelTests=false,supportingFiles
	@openapi-generator-cli generate -i docs/llm-api.yaml -g python -o formalize_convert/llm_client/ --global-property apis,models,supportingFiles --package-name llm_client

.PHONY: all codegen
