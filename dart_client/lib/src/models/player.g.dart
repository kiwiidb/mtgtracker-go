// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'player.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Player _$PlayerFromJson(Map<String, dynamic> json) => Player(
      id: (json['ID'] as num).toInt(),
      name: json['name'] as String,
      winrateAllTime: (json['winrate_all_time'] as num).toDouble(),
      numberOfGamesAllTime: (json['number_of_games_all_time'] as num).toInt(),
      decksAllTime: (json['decks_all_time'] as List<dynamic>)
          .map((e) => DeckWithCount.fromJson(e as Map<String, dynamic>))
          .toList(),
      coPlayersAllTime: (json['co_players_all_time'] as List<dynamic>)
          .map((e) => PlayerWithCount.fromJson(e as Map<String, dynamic>))
          .toList(),
      games: (json['games'] as List<dynamic>)
          .map((e) => Game.fromJson(e as Map<String, dynamic>))
          .toList(),
      currentGame: json['current_game'] == null
          ? null
          : Game.fromJson(json['current_game'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$PlayerToJson(Player instance) => <String, dynamic>{
      'ID': instance.id,
      'name': instance.name,
      'winrate_all_time': instance.winrateAllTime,
      'number_of_games_all_time': instance.numberOfGamesAllTime,
      'decks_all_time': instance.decksAllTime,
      'co_players_all_time': instance.coPlayersAllTime,
      'games': instance.games,
      'current_game': instance.currentGame,
    };

PlayerWithCount _$PlayerWithCountFromJson(Map<String, dynamic> json) =>
    PlayerWithCount(
      player: Player.fromJson(json['player'] as Map<String, dynamic>),
      count: (json['count'] as num).toInt(),
    );

Map<String, dynamic> _$PlayerWithCountToJson(PlayerWithCount instance) =>
    <String, dynamic>{
      'player': instance.player,
      'count': instance.count,
    };
