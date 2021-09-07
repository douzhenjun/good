.PHONY: all
all:
	@$(MAKE) --no-print-directory deps
	@$(MAKE) --no-print-directory dbaasbuild

.PHONY: deps
deps:
	go mod download

.PHONY: dbaasbuild
dbaasbuild:
	go build


