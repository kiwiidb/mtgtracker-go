import 'package:json_annotation/json_annotation.dart';
import 'deck.dart';
import 'player.dart';

part 'ranking.g.dart';

@JsonSerializable()
class Ranking {
  @JsonKey(name: 'ID')
  final int id;
  @JsonKey(name: 'player_id')
  final int playerId;
  final int position;
  @JsonKey(name: 'life_total')
  final int? lifeTotal;
  final Deck deck;
  final Player? player;

  const Ranking({
    required this.id,
    required this.playerId,
    required this.position,
    this.lifeTotal,
    required this.deck,
    this.player,
  });

  factory Ranking.fromJson(Map<String, dynamic> json) => _$RankingFromJson(json);
  Map<String, dynamic> toJson() => _$RankingToJson(this);
}