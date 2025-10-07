import 'package:json_annotation/json_annotation.dart';
import 'ranking.dart';

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
  @JsonKey(name: 'source_ranking')
  final Ranking? sourceRanking;
  @JsonKey(name: 'target_ranking')
  final Ranking? targetRanking;
  @JsonKey(name: 'image_url')
  final String imageUrl;
  @JsonKey(name: 'upload_image_url')
  final String? uploadImageUrl;
  @JsonKey(name: 'comment')
  final String? comment;

  const GameEvent({
    required this.gameId,
    required this.eventType,
    required this.damageDelta,
    required this.createdAt,
    required this.targetLifeTotalAfter,
    this.sourceRanking,
    this.targetRanking,
    required this.imageUrl,
    this.uploadImageUrl,
    this.comment,
  });

  factory GameEvent.fromJson(Map<String, dynamic> json) =>
      _$GameEventFromJson(json);
  Map<String, dynamic> toJson() => _$GameEventToJson(this);
}
