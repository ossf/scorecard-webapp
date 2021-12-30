LDFLAGS=-w -extldflags

build-scorecard-webapp: ## Runs go build on repo
	# Run go build and generate scorecard-webapp executable
	CGO_ENABLED=0 go build -trimpath -a -tags netgo -ldflags '$(LDFLAGS)' -o scorecard-webapp
