buid:
	@go build -o ./bin/heimdallr

clean:
	@rm -rf ./bin/

run: clean buid
	./bin/heimdallr