import 'package:json_annotation/json_annotation.dart';
import 'player.dart';

part 'follow.g.dart';

@JsonSerializable()
class Follow {
  final int id;
  final Player player1;
  final Player player2;

  const Follow({
    required this.id,
    required this.player1,
    required this.player2,
  });

  factory Follow.fromJson(Map<String, dynamic> json) => _$FollowFromJson(json);
  Map<String, dynamic> toJson() => _$FollowToJson(this);
}
