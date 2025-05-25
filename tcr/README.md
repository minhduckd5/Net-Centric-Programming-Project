# Text-Based Clash Royale (TCR)

A network-based multiplayer game implementing TCP/UDP communication protocols, built in Go.

## Features

- Two game modes: Simple (turn-based) and Enhanced (real-time)
- Player authentication system
- Mana and EXP systems
- Tower and troop management
- Critical hit mechanics
- Leveling system

## Prerequisites

- Go 1.21 or higher
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd tcr
```

2. Install dependencies:
```bash
go mod download
```

## Building

```bash
# Build server
go build -o bin/server ./server

# Build client
go build -o bin/client ./client
```

## Running

1. Start the server:
```bash
./bin/server -config config/prod.json
```

2. Start the client:
```bash
./bin/client -server localhost:8080
```

## Configuration

Configuration files are located in the `config` directory:
- `dev.json`: Development settings
- `prod.json`: Production settings

## Project Structure

```
tcr/
├── client/           # Client implementation
├── server/           # Server implementation
├── config/           # Configuration files
├── models.go         # Data structures
├── network.go        # Network communication
├── game.go           # Game logic
└── specs.json        # Game specifications
```

## Game Rules

### Simple Mode
- Turn-based gameplay
- Players take turns deploying troops
- Must destroy Guard Towers before King Tower

### Enhanced Mode
- Real-time gameplay (3-minute matches)
- Mana regeneration (1 per second)
- Critical hit system
- EXP and leveling system

## API Documentation

### Client-Server Protocol

#### Authentication
- `login`: Send username and password
- `login_resp`: Server response with status

#### Game Commands
- `deploy`: Deploy a troop
- `state_update`: Game state update
- `game_end`: Match conclusion

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 