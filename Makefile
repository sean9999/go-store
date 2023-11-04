REPO=github.com/sean9999/go-store
SEMVER := $$(git tag --sort=-version:refname | head -n 1)

.PHONY: test

info:
	echo REPO is ${REPO} and SEMVER is ${SEMVER}

run:
	go run cmd/main.go 

tidy:
	go mod tidy

vendor:
	go mod vendor

publish:
	GOPROXY=proxy.golang.org go list -m ${REPO}@${SEMVER}