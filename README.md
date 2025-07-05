# wrong-answer-client

A terminal-based WebSocket client for the party game "Wrong Answer".

This is the official CLI client that connects to a `wrong-answer` server instance and allows users to participate in the game via the terminal.

## Features

- Connects to the server using WebSocket
- Supports answering the secret question
- Displays all player answers after submission
- Handles voting phase and displays the result

## Getting Started

### Requirements

- Go 1.20 or newer

### Installation

```bash
git clone https://github.com/yourname/wrong-answer-client.git
cd wrong-answer-client
go run .
```

## Configuration

A basic config file is stored at:

```
~/.config/wrong-answer-client/config.yaml
```

This is currently unused but reserved for future features.

## Usage

Run the client and follow the prompts:

```bash
go run .
```

You will be asked to enter a username and then wait for the game to start.

## Server

This client requires a running instance of the `wrong-answer` server.

Repository: https://github.com/yourname/wrong-answer

## License

MIT

## Author

https://github.com/jad0s
