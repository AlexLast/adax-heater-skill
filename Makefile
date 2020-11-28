darwin:
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -a -o build/adax-skill github.com/alexlast/adax-heater-skill/cmd/skill
lambda:
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -a -o build/adax-skill-linux github.com/alexlast/adax-heater-skill/cmd/skill
	zip -j build/adax-lambda.zip build/adax-skill-linux
