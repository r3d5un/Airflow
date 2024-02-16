# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## test: run all tests found
.PHONY: test
test:
	@echo 'Running tests'
	go test ./...

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	golines . -w
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/docker/brreg/units: Run the BRREG Units ETL Process
.PHONY: run/docker/brreg/units
run/docker/brreg/units:
	@echo 'Running BRREG Units ETL Process...'
	docker run --network="host" brreg:latest -dataset units

## run/docker/brreg/subunits: Run the BRREG Subunits ETL Process
.PHONY: run/docker/brreg/subunits
run/docker/brreg/subunits:
	@echo 'Running BRREG Subunits ETL Process...'
	docker run --network="host" brreg:latest -dataset subunits
