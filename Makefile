PKGS = $(shell go list ./... | grep -v /vendor/ | grep -v /test/)
VERSION = $(shell cat VERSION)
PROMU := $(GOPATH)/bin/promu
GITHUB_RELEASE := $(GOPATH)/bin/github-release

build: promu
	@$(PROMU) build --prefix $(shell pwd)

crossbuild:
	@GOOS=linux GOARCH=amd64 go build -o .build/linux-amd64/dexy ./cmd
	@GOOS=darwin GOARCH=amd64 go build -o .build/darwin-amd64/dexy ./cmd
	@GOOS=windows GOARCH=amd64 go build -o .build/windows-amd64/dexy.exe ./cmd

test:
	@go test --short $(PKGS)
	@go vet $(PKGS)

release: github-release
	@git tag -a $(VERSION) -m "Version $(VERSION)"
	@git push --tags
	@$(GITHUB_RELEASE) release \
		--user chronojam
		--repo dexy \
		--tag $(VERSION) \
		--name "dexy-$(VERSION)" \
		--description ""

	@$(GITHUB_RELEASE) upload \
		--user chronojam \
		--repo dexy \
		--tag $(VERSION) \
		--name "linux-amd64" \
		--file .build/linux-amd64/dexy

	@$(GITHUB_RELEASE) upload \
		--user chronojam \
		--repo dexy \
		--tag $(VERSION) \
		--name "darwin-amd64" \
		--file .build/darwin-amd64/dexy

	@$(GITHUB_RELEASE) upload \
		--user chronojam \
		--repo dexy \
		--tag $(VERSION) \
		--name "windows-amd64" \
		--file .build/windows-amd64/dexy.exe

github-release:
	@go get -u github.com/aktau/github-release

promu:
	@go get -u github.com/prometheus/promu
