PKGS := k2m m2k pkgs

install:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

workspace:
	@go work init
	@go work use pkgs
	@go work use k2m
	@go work use m2k

download:
	@echo "Downloading external packages..."
	@for dir in ${PKGS}; do \
		echo "Downloading packages for: $$dir..."; \
		cd $$dir; \
		go mod download; \
		cd ..; \
	done
	@echo "External packages downloaded successfully!"

update:
	@echo "Updating external packages..."
	@for dir in ${PKGS}; do \
		echo "Updating packages for: $$dir..."; \
		cd $$dir; \
		go get -u all; \
		go mod tidy; \
		cd ..; \
	done
	@echo "External packages updated successfully!"

tests:
	@echo "Running unit tests..."
	@for dir in ${PKGS}; do \
		echo "Testing package: $$dir..."; \
		cd $$dir; \
		go test ./... -v -covermode atomic -coverprofile=coverage.out; \
		cd ..; \
	done
	@echo "All unit test runned successfully!"

lint:
	@echo "Running golangci-lint..."
	@for dir in ${PKGS}; do \
		echo "Testing package: $$dir..."; \
		cd $$dir; \
		golangci-lint run --print-issued-lines=false --print-linter-name=false --issues-exit-code=0 --enable=revive -- ./...; \
		cd ..; \
	done

test-cov:
	@go test ./... -v -covermode atomic -coverprofile=coverage.out
