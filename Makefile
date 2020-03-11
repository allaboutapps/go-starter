download:
	@echo Downloading go.mod dependencies...
	@go mod download

init:
	@echo Installing tools from tools.go...
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

clean:
	@rm -rf bin