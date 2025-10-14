import 'package:json_annotation/json_annotation.dart';
import 'ranking.dart';
import 'game_event.dart';
import 'player.dart';

part 'game.g.dart';

@JsonSerializable()
class Game {
  final int id;
  @JsonKey(name: 'creator_id')
  final String? creatorId;
  final int? duration;
  final DateTime? date;
  @JsonKey(name: 'end_date')
  final DateTime? endDate;
  final String? comments;
  final List<Ranking> rankings;
  final bool? finished;
  @JsonKey(name: 'game_events')
  final List<GameEvent>? gameEvents;
  final Player? creator;

  const Game({
    required this.id,
    this.creatorId,
    this.duration,
    this.date,
    this.endDate,
    required this.comments,
    required this.rankings,
    required this.finished,
    required this.gameEvents,
    this.creator,
  });

  factory Game.fromJson(Map<String, dynamic> json) => _$GameFromJson(json);
  Map<String, dynamic> toJson() => _$GameToJson(this);
}
