GO_LOCAL_MODS = $(shell go list ./...)
COVERAGE_XML_TOOL = "github.com/boumenot/gocover-cobertura"
RELEASE_FLAGS = -ldflags "-s -w -extldflags '-static'"
DEBUG_FLAGS =
TARGET = keruu
SOURCES = $(shell find . -name '*.go')

.PHONY: debug
debug: FLAGS = $(DEBUG_FLAGS)
debug: $(TARGET)

.PHONY: release
release: FLAGS = $(RELEASE_FLAGS)
release: $(TARGET)

$(TARGET): go.* $(SOURCES)
	CGO_ENABLED=0 go build $(FLAGS) -o $@ .

.PHONY: lint
lint:
	golangci-lint run -v --skip-dirs '(^|/)\.go($|/)'

.PHONY: test
test:
	go test \
		-race \
		-coverprofile=coverage.txt \
		-covermode=atomic \
		$(GO_LOCAL_MODS)

coverage.txt: test

coverage.xml: coverage.txt
	go run $(COVERAGE_XML_TOOL) < $< > $@

.PHONY: init-tools
init-tools:
	go get $(COVERAGE_XML_TOOL)

.PHONY: coverage
coverage: coverage.xml

.PHONY: clean
clean:
	rm -f keruu
	rm -f coverage.txt
	rm -f coverage.xml

.PHONY: update-deps
update-deps:
	go get -u
	go mod tidy

.PHONY: download-deps
download-deps:
	go mod download

