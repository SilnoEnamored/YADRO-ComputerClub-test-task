build:
	go build -o bin/task.exe cmd/main.go

run:
	@echo "Usage: make run ARGS=<filename>"
	./bin/task.exe $(TEST)
