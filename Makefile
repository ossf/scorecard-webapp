include Makefile.swagger

LDFLAGS=-w -extldflags

scorecard-webapp: main.go app/server/*.go ## Runs go build on repo
	# Run go build and generate scorecard-webapp executable
	CGO_ENABLED=0 go build -trimpath -a -tags netgo -ldflags '$(LDFLAGS)' -o scorecard-webapp

scorecard-webapp-docker: scorecard-webapp Dockerfile
	DOCKER_BUILDKIT=1 docker build . --file Dockerfile --tag scorecard-webapp

linter:
	golangci-lint run -c .golangci.yml

swagger: openapi.yaml
	swagger validate openapi.yaml
	swagger generate server -f openapi.yaml -q -r COPYRIGHT.txt -t app/generated  -A scorecard --flag-strategy=pflag --keep-spec-order --exclude-main --default-produces application/json --additional-initialism=SSF
	swagger generate client -f openapi.yaml -q -r COPYRIGHT.txt -t app/generated --keep-spec-order --additional-initialism=SSF
	@echo "# This file is generated after swagger runs as part of the build; do not edit!" > Makefile.swagger
	@echo "SWAGGER_GEN=`find app/generated/client app/generated/models/ app/generated/restapi/ -iname '*.go' | grep -v 'configure_scorecard' | sort -d | tr '\n' ' ' | sed 's/ $$//'`" >> Makefile.swagger

clean:
	rm -rf $(SWAGGER_GEN)

