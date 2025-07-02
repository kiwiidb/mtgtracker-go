import 'package:json_annotation/json_annotation.dart';
import 'ranking.dart';
import 'game_event.dart';

part 'game.g.dart';

@JsonSerializable()
class Game {
  @JsonKey(name: 'ID')
  final int id;
  final int? duration;
  final DateTime? date;
  final String comments;
  final String image;
  final List<Ranking> rankings;
  final bool finished;
  @JsonKey(name: 'GameEvents')
  final List<GameEvent> gameEvents;

  const Game({
    required this.id,
    this.duration,
    this.date,
    required this.comments,
    required this.image,
    required this.rankings,
    required this.finished,
    required this.gameEvents,
  });

  factory Game.fromJson(Map<String, dynamic> json) => _$GameFromJson(json);
  Map<String, dynamic> toJson() => _$GameToJson(this);
}