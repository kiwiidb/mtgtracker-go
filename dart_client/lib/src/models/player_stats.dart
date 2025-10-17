import 'package:json_annotation/json_annotation.dart';

part 'player_stats.g.dart';

@JsonSerializable()
class PlayerStats {
  final int id;
  @JsonKey(name: 'player_id')
  final String playerId;
  final DateTime timestamp;
  @JsonKey(name: 'total_wins')
  final int totalWins;
  final double winrate;
  @JsonKey(name: 'rolling_winrate')
  final double rollingWinrate;
  @JsonKey(name: 'game_count')
  final int gameCount;
  @JsonKey(name: 'game_duration')
  final int gameDuration;
  final int streak;
  final int elo;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;

  PlayerStats({
    required this.id,
    required this.playerId,
    required this.timestamp,
    required this.totalWins,
    required this.winrate,
    required this.rollingWinrate,
    required this.gameCount,
    required this.gameDuration,
    required this.streak,
    required this.elo,
    required this.createdAt,
    required this.updatedAt,
  });

  factory PlayerStats.fromJson(Map<String, dynamic> json) =>
      _$PlayerStatsFromJson(json);

  Map<String, dynamic> toJson() => _$PlayerStatsToJson(this);
}
