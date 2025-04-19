# Base flags
GO = go

# Override the base flags using makerc.mk
-include .makerc.mk

# Flags that makerc.mk doesn't touch
GO_LOCAL_MODS = $(shell $(GO) list ./...)
RELEASE_FLAGS = -ldflags "-s -w -extldflags '-static'"
DEBUG_FLAGS =
TARGET = keruu
SOURCES = $(shell find . -name '*.go')

# Targets

.PHONY: debug
debug: FLAGS = $(DEBUG_FLAGS)
debug: $(TARGET)

.PHONY: release
release: FLAGS = $(RELEASE_FLAGS)
release: $(TARGET)

$(TARGET): go.* $(SOURCES)
	$(GO) version
	CGO_ENABLED=0 $(GO) build $(FLAGS) -o $@ .

.PHONY: lint
lint:
	golangci-lint run -v --skip-dirs '(^|/)\.go($|/)'

.PHONY: test
test:
	$(GO) version
	CGO_ENABLED=1 $(GO) test \
		-race \
		-coverprofile=coverage.txt \
		-covermode=atomic \
		$(GO_LOCAL_MODS)

coverage.txt: test

coverage.xml: coverage.txt
	$(GO) tool gocover-cobertura < $< > $@

.PHONY: coverage
coverage: coverage.xml

.PHONY: clean
clean:
	rm -f keruu
	rm -f coverage.txt
	rm -f coverage.xml

.PHONY: update-deps
update-deps:
	$(GO) get -u
	$(GO) mod tidy

.PHONY: download-deps
download-deps:
	$(GO) mod download

