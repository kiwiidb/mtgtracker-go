import 'package:mtgtracker_client/mtgtracker_client.dart';

void main() async {
  // Initialize the client with your API host and auth token
  final client = MTGTrackerClient(
    baseUrl: 'https://api-staging.plowshare.social/api',
    authToken: 'your-firebase-auth-token',
  );

  try {
    // Example: Sign up a new player
    print('Signing up a new player...');
    final signupRequest = SignupPlayerRequest(
      name: 'Alice Smith',
    );
    final newPlayer = await client.signupPlayer(signupRequest);
    print('New player created: ${newPlayer.name} (ID: ${newPlayer.id})');

    // Example: Get all players
    print('\nFetching all players...');
    final players = await client.getPlayers();
    print('Found ${players.length} players');

    // Example: Search for players
    print('\nSearching for players named "Alice"...');
    final searchResults = await client.getPlayers(search: 'Alice');
    print('Found ${searchResults.length} players matching "Alice"');

    // Example: Get current user's profile
    print('\nFetching my player profile...');
    final myPlayer = await client.getMyPlayer();
    print('My profile: ${myPlayer.name}');
    print('Games played: ${myPlayer.numberOfGamesAllTime}');
    print('Win rate: ${(myPlayer.winrateAllTime * 100).toStringAsFixed(1)}%');

    // Uncomment when you have proper ranking data:
    // final newGame = await client.createGame(createGameRequest);
    // print('New game created: ${newGame.id}');

    // Example: Get all games
    print('\nFetching all games...');
    final games = await client.getGames();
    print('Found ${games.length} games');

    // Example: Add a game event (if you have a game)
    if (games.isNotEmpty) {
      final gameId = games.first.id;
      print('\nAdding event to game $gameId...');

      final eventRequest = GameEventRequest(
        eventType: 'damage',
        damageDelta: -5,
        targetLifeTotalAfter: 35,
        sourceRankingId: 1,
        targetRankingId: 2,
        comment: 'Lightning Bolt to the face!',
      );

      final event = await client.addGameEvent(gameId, eventRequest);
      print('Event added: ${event.eventType} for ${event.damageDelta} damage');
    }
  } on MTGTrackerException catch (e) {
    print('API Error: ${e.statusCode} - ${e.message}');
  } catch (e) {
    print('Unexpected error: $e');
  } finally {
    // Always dispose of the client when done
    client.dispose();
    print('\nClient disposed');
  }
}
