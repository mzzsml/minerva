buid:
	@go build -o ./bin/minerva

clean:
	@rm -rf ./bin/

run: clean buid
	./bin/minerva