import 'package:json_annotation/json_annotation.dart';
import 'player.dart';

part 'opponent.g.dart';

@JsonSerializable()
class Opponent {
  final int id;
  final Player player1;
  final Player player2;

  const Opponent({
    required this.id,
    required this.player1,
    required this.player2,
  });

  factory Opponent.fromJson(Map<String, dynamic> json) => _$OpponentFromJson(json);
  Map<String, dynamic> toJson() => _$OpponentToJson(this);
}
