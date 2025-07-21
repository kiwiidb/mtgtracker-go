# MTG Tracker CLI

A command-line interface for the MTG Tracker API, built using the dart_client library.

## Installation

1. Make sure you have Dart SDK installed
2. Navigate to the dart_client directory:
   ```bash
   cd dart_client
   ```
3. Get dependencies:
   ```bash
   dart pub get
   ```
4. Install the CLI globally (optional):
   ```bash
   dart pub global activate .
   ```

## Usage

### Basic Commands

Run the CLI using:
```bash
dart bin/mtg_cli.dart [options] <command> [arguments]
```

Or if installed globally:
```bash
mtg_cli [options] <command> [arguments]
```

### Global Options

- `-h, --help`: Show help information
- `-v, --version`: Show version
- `-s, --server <url>`: Set MTG Tracker server URL (default: http://localhost:8080)
- `-t, --token <token>`: Set authentication token

### Player Management

```bash
# Sign up a new player
mtg_cli player signup "John Doe" "john@example.com"

# List all players
mtg_cli player list

# Search for players
mtg_cli player list "John"

# Get player by ID
mtg_cli player get 1

# Get current authenticated player info
mtg_cli player me
```

### Game Management

```bash
# Create a new game
mtg_cli game create "Commander night at the shop"

# List all games
mtg_cli game list

# Get game details
mtg_cli game get 1

# Update game status
mtg_cli game update 1 finish
mtg_cli game update 1 reopen

# Delete a game
mtg_cli game delete 1

# Add game event
mtg_cli game event 1 damage -5 35
```

### Ranking Management

```bash
# List pending rankings
mtg_cli ranking pending

# Accept a ranking
mtg_cli ranking accept 1

# Decline a ranking
mtg_cli ranking decline 1
```

### Examples with Server and Authentication

```bash
# Use custom server
mtg_cli --server https://mtgtracker.example.com player list

# Use authentication token
mtg_cli --token abc123 player me

# Combined options
mtg_cli --server https://mtgtracker.example.com --token abc123 game list
```

## Error Handling

The CLI provides clear error messages for:
- Invalid commands or arguments
- Network connection issues
- API authentication failures
- Invalid data formats

All errors are displayed with helpful suggestions for fixing the issue.

## Development

To run the CLI during development:
```bash
dart bin/mtg_cli.dart --help
```

To test specific functionality:
```bash
# Test with local server (make sure server is running)
dart bin/mtg_cli.dart player list
dart bin/mtg_cli.dart game create "Test game"
```