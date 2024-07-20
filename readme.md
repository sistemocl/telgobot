# Chat Bot using Telegram and Golang

This project implements a chat bot using Telegram and Golang, telebot, chromedp y .env. The bot is designed to respond to specific commands and perform certain actions.

## Requirements
- Golang
- telebot.v2
- chromedp
- .env

To check if you have golang installed, run the following command:
```bash
go version
```
If you dont have Golang installed on your system visit https://golang.org/dl/

## Intalation
* 1.Clone this repository on your local machine:
```bash
git clone https://github.com/N-Ignacio-Bouffanais/telgobot

```

* 2.Install the dependencies
```go
go mod tidy
```

* 3.Run the command
```bash
go run ./src/main.go
```

## Another way to install 
* 1.Create a folder for your project and named, this name will be the same for the next command.
* 2.Run the command
```bash
go mod init your_project_name
```
* 3.Create a folder inside named src, and inside src create a file named main.go
* 4.Create a .env file with credentials.
* 5.Run the command:
```bash
go run ./src/main.go
```
* 6 If all was correct well, you must see a message in your terminal like:"Bot is running. Press CTRL+C to exit."

## Configuration

Before run the script, you need to set the environment variables on .env.example, I suggest you to change the name of .env.example to .env
```bash
BOT_TOKEN="fwejnviewnvwenvwenvja"
PASSWORD="password"
USER="USER"
ADMIN_USER="admin"
ADMIN_PASS=""
```

## Telegram Token
To get an access token for telegram bot, you need to search on Telegram one bot call "BotFather", this bot will let you build a bot on telegram api.
* 1.send the /start command to Botfather.
* 2.send the /newbot command to Botfather.
* 3.select a name for your bot like bot_example or something like that.
* 4.save the toket that Botfather give to you and that token use it on your project.

## License

This project is licensed under the MIT License.