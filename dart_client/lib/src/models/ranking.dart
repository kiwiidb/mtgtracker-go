import 'package:json_annotation/json_annotation.dart';
import 'deck.dart';
import 'player.dart';

part 'ranking.g.dart';

@JsonSerializable()
class Ranking {
  final int id;
  @JsonKey(name: 'player_id')
  final String? playerId;
  int position;
  @JsonKey(name: 'life_total')
  final int? lifeTotal;
  @JsonKey(name: 'last_life_total')
  final int? lastLifeTotal;
  @JsonKey(name: 'last_life_total_timestamp')
  final DateTime? lastLifeTotalTimestamp;
  final Deck deck;
  final Player? player;
  final String? status;

  Ranking({
    required this.id,
    this.playerId,
    required this.position,
    this.lifeTotal,
    this.lastLifeTotal,
    this.lastLifeTotalTimestamp,
    required this.deck,
    this.player,
    this.status,
  });

  factory Ranking.fromJson(Map<String, dynamic> json) =>
      _$RankingFromJson(json);
  Map<String, dynamic> toJson() => _$RankingToJson(this);
}
