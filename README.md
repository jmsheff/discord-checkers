# Discord checkers bot
[<img src="images/checkers.png" width="100" />](https://discordapp.com/oauth2/authorize?client_id=469537029546442752&permissions=93248&scope=bot)

## Add to your server
1. Click the logo above or [this link](https://discordapp.com/oauth2/authorize?client_id=469537029546442752&permissions=93248&scope=bot) to invite the bot to your server.
2. Type the command `!checkers ping` to test the bot
3. That's it! Everything else you need to know can be obtained by typing `!checkers help`.

## Demo
<img src="images/checkers_demo.gif" width="200" />

## Running locally
This project is written in [Go](golang.org) using [discordgo](https://github.com/bwmarrin/discordgo) which also happens to be it's only external dependency. Different from what you might expect, there is **no database**. All game states are stored in the messages sent. 
1. Clone the repository
2. Install dependencies by running `go get github.com/bwmarrin/discordgo`
3. If you haven't already, go to the [Discord developer portal](https://discordapp.com/developers/applications) and create a new application to obtain a token.
4. Set the environment variable `BOT_TOKEN` to the token of your bot(which can also be obtained in the previous step)
5. Run the bot by running `go run main.go`
