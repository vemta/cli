unix_64:
	SET GOOS=linux
	SET GOARCH=amd64
	go build -o bin/vemta-amd64-unix app/main.go  
unix_86:
	SET GOOS=linux
	SET GOARCH=386
	go build -o bin/vemta-x86-unix app/main.go  
darwin_64:
	SET GOOS=darwin
	SET GOARCH=amd64
	go build -o bin/vemta-amd64-darwin app/main.go  
darwin_86:
	SET GOOS=darwin
	SET GOARCH=386
	go build -o bin/vemta-x86-darwin app/main.go 
windows_64:
	SET GOOS=windows
	SET GOARCH=amd64
	go build -o bin/vemta-amd64.exe app/main.go 
windows_86:
	SET GOOS=windows
	SET GOARCH=386
	go build -o bin/vemta-x86.exe app/main.go 

