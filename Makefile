test:
	go test ./... 

clean:
	find ./ -name "database.db" -exec rm -rf {} \;