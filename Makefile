run:
	go build . && ./bof

clean:
	rm ./bof

db:
	sqlite3 ./bof.db
