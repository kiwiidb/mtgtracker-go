import 'package:json_annotation/json_annotation.dart';
import 'game.dart';
import 'player.dart';

part 'notification.g.dart';

enum NotificationAction {
  @JsonValue('delete_ranking')
  deleteRanking,
  @JsonValue('view_game')
  viewGame,
}

@JsonSerializable()
class Notification {
  final int id;
  final String title;
  final String body;
  final String type;
  final List<NotificationAction> actions;
  final bool read;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'game_id')
  final int? gameId;
  @JsonKey(name: 'referred_player_id')
  final String? referredPlayerId;
  final Game? game;
  @JsonKey(name: 'referred_player')
  final Player? referredPlayer;

  const Notification({
    required this.id,
    required this.title,
    required this.body,
    required this.type,
    required this.actions,
    required this.read,
    required this.createdAt,
    this.gameId,
    this.referredPlayerId,
    this.game,
    this.referredPlayer,
  });

  factory Notification.fromJson(Map<String, dynamic> json) =>
      _$NotificationFromJson(json);
  Map<String, dynamic> toJson() => _$NotificationToJson(this);
}