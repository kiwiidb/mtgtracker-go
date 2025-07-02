# MTG Tracker Dart Client

A Dart HTTP client for the MTG Tracker API that allows for injecting host and auth token.

## Installation

Add this to your package's `pubspec.yaml` file:

```yaml
dependencies:
  mtgtracker_client:
    path: path/to/dart_client
```

Then run:

```bash
dart pub get
```

## Setup

Before using the client, you need to generate the JSON serialization code:

```bash
cd dart_client
dart pub get
dart pub run build_runner build
```

## Usage

### Basic Setup

```dart
import 'package:mtgtracker_client/mtgtracker_client.dart';

// Create client with host and optional auth token
final client = MTGTrackerClient(
  baseUrl: 'https://your-api-host.com',
  authToken: 'your-auth-token', // Optional
);
```

### Player Operations

```dart
// Sign up a new player
final signupRequest = SignupPlayerRequest(
  name: 'John Doe',
  email: 'john@example.com',
);
final player = await client.signupPlayer(signupRequest);

// Get all players
final players = await client.getPlayers();

// Search players
final searchResults = await client.getPlayers(search: 'John');

// Get specific player
final player = await client.getPlayer(1);

// Get current user's player profile
final myPlayer = await client.getMyPlayer();
```

### Game Operations

```dart
// Create a new game
final createRequest = CreateGameRequest(
  duration: 3600, // seconds
  date: DateTime.now(),
  comments: 'Great game!',
  image: 'image-url',
  finished: true,
  rankings: [
    // Add your rankings here
  ],
);
final game = await client.createGame(createRequest);

// Get all games
final games = await client.getGames();

// Get specific game
final game = await client.getGame(1);

// Update a game
final updateRequest = UpdateGameRequest(
  gameId: 1,
  finished: true,
  rankings: [
    // Updated rankings
  ],
);
final updatedGame = await client.updateGame(1, updateRequest);

// Delete a game
await client.deleteGame(1);

// Add game event
final eventRequest = GameEventRequest(
  eventType: 'damage',
  damageDelta: -5,
  targetLifeTotalAfter: 35,
  sourceRankingId: 1,
  targetRankingId: 2,
);
final event = await client.addGameEvent(1, eventRequest);
```

### Error Handling

The client throws `MTGTrackerException` for HTTP errors:

```dart
try {
  final player = await client.getPlayer(999);
} on MTGTrackerException catch (e) {
  print('Error ${e.statusCode}: ${e.message}');
}
```

### Clean Up

Don't forget to dispose of the client when done:

```dart
client.dispose();
```

## API Endpoints

The client supports all MTG Tracker API endpoints:

### Player Endpoints
- `POST /player/v1/signup` - Sign up a new player
- `GET /player/v1/players` - Get all players (with optional search)
- `GET /player/v1/players/{playerId}` - Get specific player
- `GET /player/v1/me` - Get current user's player profile

### Game Endpoints
- `POST /game/v1/games` - Create a new game
- `GET /game/v1/games` - Get all games
- `GET /game/v1/games/{gameId}` - Get specific game
- `PUT /game/v1/games/{gameId}` - Update a game
- `DELETE /game/v1/games/{gameId}` - Delete a game
- `POST /game/v1/games/{gameId}/events` - Add game event

## Models

The client includes generated Dart models for all API request/response types:

- `Player`, `PlayerWithCount`
- `Game`, `GameEvent`
- `Deck`, `DeckWithCount`
- `Ranking`
- `SignupPlayerRequest`
- `CreateGameRequest`, `UpdateGameRequest`
- `GameEventRequest`