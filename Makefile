test:
	go test ./... 

clean:
	find ./ -name "*.db" -exec rm -rf {} \;
