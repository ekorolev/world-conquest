## About

This game is about territorial gains. Join the game and make %COUNTRY_NAME% GREAT AGAIN! At the moment there is available five countries: Russia, United States, The United Kingdom, China and Mongolia. Choose your favourite country and join us!

![Screenshot](https://github.com/ekorolev/websocket-shooter/blob/master/screenshot.png?raw=true)

## Installation

To install game on your computer you need to follow the instructions:
- Install NodeJS and NPM https://nodejs.org/en/download/package-manager/
- Install Golang https://golang.org/doc/install
- Clone the repository `git clone https://github.com/ekorolev/websocket-shooter`
- Go to the game client directory `cd websocket-shooter/client`
- Install client dependecies `npm install`
- Compile client files `npm run build`
- Install golang dependecies: `go get github.com/alehano/wsgame/utils` `go get github.com/gorilla/websocket`
- Run program `go run server/main.go server/game.go` and two webservices `go run microservices/statistic.go`, `go run microservices/savemap.go`
- That's all! Open http://localhost:9001/static/index.html
