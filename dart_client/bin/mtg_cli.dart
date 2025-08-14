#!/usr/bin/env dart

import 'dart:io';
import 'package:args/args.dart';
import '../lib/mtgtracker_client.dart';

const String version = '1.0.0';

ArgParser buildParser() {
  return ArgParser()
    ..addFlag(
      'help',
      abbr: 'h',
      negatable: false,
      help: 'Print this usage information.',
    )
    ..addFlag(
      'version',
      abbr: 'v',
      negatable: false,
      help: 'Print the tool version.',
    )
    ..addOption(
      'server',
      abbr: 's',
      help: 'MTG Tracker server URL',
      defaultsTo: 'http://localhost:8080',
    )
    ..addOption(
      'token',
      abbr: 't',
      help: 'Authentication token',
    );
}

void printUsage(ArgParser argParser) {
  print('MTG Tracker CLI - Manage players and games\n');
  print('Usage: mtg_cli [options] <command> [arguments]\n');
  print('Global options:');
  print(argParser.usage);
  print('\nCommands:');
  print('  player                   Manage players');
  print('    signup <name> <email>  Sign up a new player');
  print('    list [search]          List all players (with optional search)');
  print('    get <id>               Get player by ID');
  print('    me                     Get current player info');
  print('');
  print('  game                     Manage games');
  print('    create <comments>      Create a new game');
  print('    mock-game <players>    Create a mock Commander game (2-4 players)');
  print('    list                   List all games');
  print('    get <id>               Get game by ID');
  print('    update <id> <field>    Update game field (finish|reopen)');
  print('    delete <id>            Delete a game');
  print('    event <id> <type> <damage> <life>  Add event to game');
  print('');
  print('  ranking                  Manage rankings');
  print('    pending                List pending rankings');
  print('    accept <id>            Accept a ranking');
  print('    decline <id>           Decline a ranking');
  print('\nExamples:');
  print('  mtg_cli player signup "John Doe" "john@example.com"');
  print('  mtg_cli game create "Commander game night"');
  print('  mtg_cli game event 1 damage -5 35');
  print('  mtg_cli --server https://api.example.com game list');
  print('  mtg_cli --token abc123 ranking pending');
}

Future<void> main(List<String> arguments) async {
  final ArgParser argParser = buildParser();

  try {
    final ArgResults results = argParser.parse(arguments);

    if (results['help'] == true) {
      printUsage(argParser);
      return;
    }

    if (results['version'] == true) {
      print('MTG Tracker CLI version: $version');
      return;
    }

    if (results.rest.isEmpty) {
      print('Error: No command specified.\n');
      printUsage(argParser);
      exit(1);
    }

    final String serverUrl = results['server'];
    final String? token = results['token'];
    final client = MTGTrackerClient(
      baseUrl: serverUrl,
      authToken: token,
    );

    final String command = results.rest[0];
    final List<String> args = results.rest.skip(1).toList();

    try {
      switch (command) {
        case 'player':
          await handlePlayerCommand(client, args);
          break;
        case 'game':
          await handleGameCommand(client, args);
          break;
        case 'ranking':
          await handleRankingCommand(client, args);
          break;
        default:
          print('Error: Unknown command "$command".\n');
          printUsage(argParser);
          exit(1);
      }
    } catch (e) {
      print('Error: $e');
      exit(1);
    } finally {
      client.dispose();
    }
  } on FormatException catch (e) {
    print('Error: ${e.message}\n');
    printUsage(argParser);
    exit(1);
  }
}

Future<void> handlePlayerCommand(
    MTGTrackerClient client, List<String> args) async {
  if (args.isEmpty) {
    print('Error: No player subcommand specified.');
    print('Available subcommands: signup, list, get, me');
    exit(1);
  }

  switch (args[0]) {
    case 'signup':
      if (args.length < 2) {
        print('Error: Player name required for signup.');
        print('Usage: player signup <name>');
        exit(1);
      }
      final player = await client.signupPlayer(SignupPlayerRequest(
        name: args[1],
      ));
      print('Player created: ${player.name} (ID: ${player.id})');
      break;

    case 'list':
      final search = args.length > 1 ? args[1] : null;
      final players = await client.getPlayers(search: search);
      if (players.isEmpty) {
        print('No players found.');
      } else {
        print('Players:');
        for (final player in players) {
          print('  ${player.id}: ${player.name}');
        }
      }
      break;

    case 'get':
      if (args.length < 2) {
        print('Error: Player ID required.');
        exit(1);
      }
      final playerId = int.tryParse(args[1]);
      if (playerId == null) {
        print('Error: Invalid player ID "${args[1]}".');
        exit(1);
      }
      final player = await client.getPlayer(playerId);
      print('Player: ${player.name} (ID: ${player.id})');
      break;

    case 'me':
      final player = await client.getMyPlayer();
      print('Current player: ${player.name} (ID: ${player.id})');
      break;

    default:
      print('Error: Unknown player subcommand "${args[0]}".');
      print('Available subcommands: signup, list, get, me');
      exit(1);
  }
}

Future<void> handleGameCommand(
    MTGTrackerClient client, List<String> args) async {
  if (args.isEmpty) {
    print('Error: No game subcommand specified.');
    print(
        'Available subcommands: create, mock-game, list, get, update, delete, event');
    exit(1);
  }

  switch (args[0]) {
    case 'create':
      if (args.length < 1) {
        print('Error: At least one argument required for game creation.');
        exit(1);
      }
      final game = await client.createGame(CreateGameRequest(
        comments: args.isNotEmpty ? args[0] : '',
        image: '',
        finished: false,
        rankings: [],
      ));
      print('Game created with ID: ${game.id}');
      break;

    case 'mock-game':
      if (args.length < 2) {
        print('Error: Number of players required for mock game.');
        print('Usage: game mock-game <players>');
        print('Players must be between 2 and 4.');
        exit(1);
      }
      
      final playerCount = int.tryParse(args[1]);
      if (playerCount == null || playerCount < 2 || playerCount > 4) {
        print('Error: Invalid player count "${args[1]}".');
        print('Players must be between 2 and 4.');
        exit(1);
      }

      // Create a mock Commander game with decks
      final allMockRankings = [
        Ranking(
          id: 0,
          playerId: 1,
          position: 1,
          lifeTotal: 40,
          deck: Deck(
            id: 0,
            commander: 'Teysa Karlov',
            crop: 'https://cards.scryfall.io/art_crop/front/c/d/cd14f1ce-7fcd-485c-b7ca-01c5b45fdc01.jpg?1689999296',
            secondaryImg: '',
            image: 'https://cards.scryfall.io/normal/front/c/d/cd14f1ce-7fcd-485c-b7ca-01c5b45fdc01.jpg?1689999296',
          ),
        ),
        Ranking(
          id: 0,
          playerId: 2,
          position: 2,
          lifeTotal: 40,
          deck: Deck(
            id: 0,
            commander: 'Ojer Axonil, Deepest Might',
            crop: 'https://cards.scryfall.io/art_crop/front/5/0/50f8e2b6-98c7-4f28-bb39-e1fbe841f1ee.jpg?1699044315',
            secondaryImg: 'https://cards.scryfall.io/art_crop/back/5/0/50f8e2b6-98c7-4f28-bb39-e1fbe841f1ee.jpg?1699044315',
            image: 'https://cards.scryfall.io/normal/front/5/0/50f8e2b6-98c7-4f28-bb39-e1fbe841f1ee.jpg?1699044315',
          ),
        ),
        Ranking(
          id: 0,
          playerId: 3,
          position: 3,
          lifeTotal: 40,
          deck: Deck(
            id: 0,
            commander: 'Queen Marchesa',
            crop: 'https://cards.scryfall.io/art_crop/front/0/f/0fdae05f-7bdc-45fb-b9b9-e5ec3766f965.jpg?1712354769',
            secondaryImg: '',
            image: 'https://cards.scryfall.io/normal/front/0/f/0fdae05f-7bdc-45fb-b9b9-e5ec3766f965.jpg?1712354769',
          ),
        ),
        Ranking(
          id: 0,
          playerId: 4,
          position: 4,
          lifeTotal: 40,
          deck: Deck(
            id: 0,
            commander: 'Lord Windgrace',
            crop: 'https://cards.scryfall.io/art_crop/front/2/1/213d6fb8-5624-4804-b263-51f339482754.jpg?1592710275',
            secondaryImg: '',
            image: 'https://cards.scryfall.io/normal/front/2/1/213d6fb8-5624-4804-b263-51f339482754.jpg?1592710275',
          ),
        ),
      ];

      final selectedRankings = allMockRankings.take(playerCount).toList();

      final mockGame = await client.createGame(CreateGameRequest(
        comments: 'Mock Commander Game - $playerCount players',
        image: '',
        finished: false,
        rankings: selectedRankings,
      ));
      print('Mock game created with ID: ${mockGame.id}');
      print('Players:');
      
      final commanders = ['Teysa Karlov', 'Ojer Axonil, Deepest Might', 'Queen Marchesa', 'Lord Windgrace'];
      for (int i = 0; i < playerCount; i++) {
        print('  ${i + 1}. ${commanders[i]} (Player ${i + 1})');
      }
      break;

    case 'list':
      final games = await client.getGames();
      if (games.isEmpty) {
        print('No games found.');
      } else {
        print('Games:');
        for (final game in games) {
          final status = game.finished ? 'finished' : 'active';
          print('  ${game.id}: ${game.comments} - $status');
        }
      }
      break;

    case 'get':
      if (args.length < 2) {
        print('Error: Game ID required.');
        exit(1);
      }
      final gameId = int.tryParse(args[1]);
      if (gameId == null) {
        print('Error: Invalid game ID "${args[1]}".');
        exit(1);
      }
      final game = await client.getGame(gameId);
      final status = game.finished ? 'finished' : 'active';
      print('Game: ${game.comments} (ID: ${game.id}) - Status: $status');
      break;

    case 'update':
      if (args.length < 3) {
        print('Error: Game ID and field required for update.');
        print('Usage: game update <id> <field>');
        print('Available fields: finish, reopen');
        exit(1);
      }
      final gameId = int.tryParse(args[1]);
      if (gameId == null) {
        print('Error: Invalid game ID "${args[1]}".');
        exit(1);
      }

      final field = args[2];

      UpdateGameRequest request;
      switch (field) {
        case 'finish':
          request = UpdateGameRequest(
            gameId: gameId,
            finished: true,
            rankings: [],
          );
          break;
        case 'reopen':
          request = UpdateGameRequest(
            gameId: gameId,
            finished: false,
            rankings: [],
          );
          break;
        default:
          print(
              'Error: Unknown field "$field". Available fields: finish, reopen');
          exit(1);
      }

      final game = await client.updateGame(gameId, request);
      final status = game.finished ? 'finished' : 'active';
      print(
          'Game updated: ${game.comments} (ID: ${game.id}) - Status: $status');
      break;

    case 'delete':
      if (args.length < 2) {
        print('Error: Game ID required for deletion.');
        exit(1);
      }
      final gameId = int.tryParse(args[1]);
      if (gameId == null) {
        print('Error: Invalid game ID "${args[1]}".');
        exit(1);
      }
      await client.deleteGame(gameId);
      print('Game $gameId deleted successfully.');
      break;

    case 'event':
      if (args.length < 5) {
        print(
            'Error: Game ID, event type, damage delta, and target life total required.');
        print(
            'Usage: game event <game_id> <event_type> <damage_delta> <target_life_total>');
        exit(1);
      }
      final gameId = int.tryParse(args[1]);
      if (gameId == null) {
        print('Error: Invalid game ID "${args[1]}".');
        exit(1);
      }

      final eventType = args[2];
      final damageDelta = int.tryParse(args[3]);
      final targetLifeTotal = int.tryParse(args[4]);

      if (damageDelta == null || targetLifeTotal == null) {
        print('Error: Damage delta and target life total must be integers.');
        exit(1);
      }

      final event = await client.addGameEvent(
          gameId,
          GameEventRequest(
            eventType: eventType,
            damageDelta: damageDelta,
            targetLifeTotalAfter: targetLifeTotal,
          ));
      print('Event added to game $gameId: ${event.eventType}');
      break;

    default:
      print('Error: Unknown game subcommand "${args[0]}".');
      print(
          'Available subcommands: create, mock-game, list, get, update, delete, event');
      exit(1);
  }
}

Future<void> handleRankingCommand(
    MTGTrackerClient client, List<String> args) async {
  if (args.isEmpty) {
    print('Error: No ranking subcommand specified.');
    print('Available subcommands: pending, accept, decline');
    exit(1);
  }

  switch (args[0]) {
    case 'pending':
      final rankings = await client.getPendingRankings();
      if (rankings.isEmpty) {
        print('No pending rankings found.');
      } else {
        print('Pending rankings:');
        for (final ranking in rankings) {
          print(
              '  ${ranking.id}: Player ${ranking.playerId} - Position ${ranking.position}');
        }
      }
      break;

    case 'accept':
      if (args.length < 2) {
        print('Error: Ranking ID required.');
        exit(1);
      }
      final rankingId = int.tryParse(args[1]);
      if (rankingId == null) {
        print('Error: Invalid ranking ID "${args[1]}".');
        exit(1);
      }
      await client.acceptRanking(rankingId);
      print('Ranking $rankingId accepted.');
      break;

    case 'decline':
      if (args.length < 2) {
        print('Error: Ranking ID required.');
        exit(1);
      }
      final rankingId = int.tryParse(args[1]);
      if (rankingId == null) {
        print('Error: Invalid ranking ID "${args[1]}".');
        exit(1);
      }
      await client.declineRanking(rankingId);
      print('Ranking $rankingId declined.');
      break;

    default:
      print('Error: Unknown ranking subcommand "${args[0]}".');
      print('Available subcommands: pending, accept, decline');
      exit(1);
  }
}
