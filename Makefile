BINARY=pmguard

dev:
	go build -o $(BINARY) .

fmt:
	gofmt -w .

tidy:
	go mod tidy

release:
	@[ "$(VERSION)" ] || ( echo "❌ Usage: make release VERSION=v0.1.0"; exit 1 )
	git tag $(VERSION)
	git push origin $(VERSION)