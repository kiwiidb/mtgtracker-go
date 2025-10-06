import 'package:json_annotation/json_annotation.dart';
import 'deck.dart';
import 'game.dart';

part 'player.g.dart';

@JsonSerializable()
class Player {
  final String id;
  final String name;
  @JsonKey(name: 'profile_image_url')
  final String? profileImageUrl;
  @JsonKey(name: 'moxfield_username')
  final String? moxfieldUsername;
  final List<String>? colors; // Top 2 most played colors
  @JsonKey(name: 'winrate_all_time')
  final double winrateAllTime;
  @JsonKey(name: 'number_of_games_all_time')
  final int numberOfGamesAllTime;
  @JsonKey(name: 'decks_all_time')
  final List<DeckWithCount>? decksAllTime;
  @JsonKey(name: 'co_players_all_time')
  final List<PlayerWithCount>? coPlayersAllTime;
  final List<Game>? games;
  @JsonKey(name: 'current_game')
  final Game? currentGame;

  const Player({
    required this.id,
    required this.name,
    this.profileImageUrl,
    this.moxfieldUsername,
    this.colors,
    required this.winrateAllTime,
    required this.numberOfGamesAllTime,
    required this.decksAllTime,
    required this.coPlayersAllTime,
    required this.games,
    this.currentGame,
  });

  factory Player.fromJson(Map<String, dynamic> json) => _$PlayerFromJson(json);
  Map<String, dynamic> toJson() => _$PlayerToJson(this);
}

@JsonSerializable()
class PlayerWithCount {
  final Player player;
  final int count;

  const PlayerWithCount({
    required this.player,
    required this.count,
  });

  factory PlayerWithCount.fromJson(Map<String, dynamic> json) =>
      _$PlayerWithCountFromJson(json);
  Map<String, dynamic> toJson() => _$PlayerWithCountToJson(this);
}
