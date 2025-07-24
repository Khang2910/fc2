GOOS=windows GOARCH=amd64 go build -o serving/victim.txt victim.go
sudo nohup python3 -m http.server 8082
