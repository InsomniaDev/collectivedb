test:
	make clean
	go test ./... 

clean:
	find ./ -name "*.db" -exec rm -rf {} \;
