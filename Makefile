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

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/docker/brreg: run the app using docker compose
.PHONY: build/docker/brreg
build/docker/brreg:
	@echo 'Building containerized BRREG CLI app...'
	docker build -f ./dockerfiles/brreg.Dockerfile -t brreg:latest .

# ==================================================================================== #
# PUBLISH
# ==================================================================================== #

## publish/docker/brreg: publish the brreg app to ghcr.io
.PHONY: publish/docker/brreg
publish/docker/brreg:
	docker build -f ./dockerfiles/brreg.Dockerfile -t brreg:latest .
	@echo 'Retag BRREG CLI app...'
	docker tag brreg:latest ghcr.io/r3d5un/brreg:latest
	@echo 'Publishing BRREG CLI app to ghcr.io...'
	docker push ghcr.io/r3d5un/brreg:latest
