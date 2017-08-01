PKGS = $(shell go list ./... | grep -v /vendor/ | grep -v /test/)
VERSION = $(shell cat VERSION)
PROMU := $(GOPATH)/bin/promu
GITHUB_RELEASE := $(GOPATH)/bin/github-release

build: promu
	@$(PROMU) --verbose crossbuild
	# @$(PROMU) build --prefix $(shell pwd)

test:
	@go test --short $(PKGS)
	@go vet $(PKGS)

release: github-release
	@git tag -a $(VERSION) -m "Version $(VERSION)"
	@$(GITHUB_RELEASE) upload \
		--user chronojam \
		--repo dexy \
		--tag $(VERSION) \
		--name "linux-amd64" \
		--file dexy

github-release:
	@go get -u github.com/aktau/github-release

promu:
	@go get -u github.com/prometheus/promu
