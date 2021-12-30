LDFLAGS=-w -extldflags

scorecard-webapp: ## Runs go build on repo
	# Run go build and generate scorecard-webapp executable
	CGO_ENABLED=0 go build -trimpath -a -tags netgo -ldflags '$(LDFLAGS)' -o scorecard-webapp

scorecard-webapp-docker:
	DOCKER_BUILDKIT=1 docker build . --file Dockerfile --tag scorecard-webapp
