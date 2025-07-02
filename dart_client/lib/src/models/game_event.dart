import 'package:json_annotation/json_annotation.dart';

part 'game_event.g.dart';

@JsonSerializable()
class GameEvent {
  @JsonKey(name: 'game_id')
  final int gameId;
  @JsonKey(name: 'event_type')
  final String eventType;
  @JsonKey(name: 'damage_delta')
  final int damageDelta;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'target_life_total_after')
  final int targetLifeTotalAfter;
  @JsonKey(name: 'source_player')
  final String sourcePlayer;
  @JsonKey(name: 'target_player')
  final String targetPlayer;
  @JsonKey(name: 'source_commander')
  final String sourceCommander;
  @JsonKey(name: 'target_commander')
  final String targetCommander;
  @JsonKey(name: 'image_url')
  final String imageUrl;
  @JsonKey(name: 'upload_image_url')
  final String? uploadImageUrl;

  const GameEvent({
    required this.gameId,
    required this.eventType,
    required this.damageDelta,
    required this.createdAt,
    required this.targetLifeTotalAfter,
    required this.sourcePlayer,
    required this.targetPlayer,
    required this.sourceCommander,
    required this.targetCommander,
    required this.imageUrl,
    this.uploadImageUrl,
  });

  factory GameEvent.fromJson(Map<String, dynamic> json) => _$GameEventFromJson(json);
  Map<String, dynamic> toJson() => _$GameEventToJson(this);
}