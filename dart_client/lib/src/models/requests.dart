import 'package:json_annotation/json_annotation.dart';
import 'ranking.dart';

part 'requests.g.dart';

@JsonSerializable()
class SignupPlayerRequest {
  final String name;
  final String email;

  const SignupPlayerRequest({
    required this.name,
    required this.email,
  });

  factory SignupPlayerRequest.fromJson(Map<String, dynamic> json) => _$SignupPlayerRequestFromJson(json);
  Map<String, dynamic> toJson() => _$SignupPlayerRequestToJson(this);
}

@JsonSerializable()
class CreateGameRequest {
  final int? duration;
  final DateTime? date;
  final String comments;
  final String image;
  final bool finished;
  final List<Ranking> rankings;

  const CreateGameRequest({
    this.duration,
    this.date,
    required this.comments,
    required this.image,
    required this.finished,
    required this.rankings,
  });

  factory CreateGameRequest.fromJson(Map<String, dynamic> json) => _$CreateGameRequestFromJson(json);
  Map<String, dynamic> toJson() => _$CreateGameRequestToJson(this);
}

@JsonSerializable()
class UpdateGameRequest {
  @JsonKey(name: 'game_id')
  final int gameId;
  final bool? finished;
  final List<Ranking> rankings;

  const UpdateGameRequest({
    required this.gameId,
    this.finished,
    required this.rankings,
  });

  factory UpdateGameRequest.fromJson(Map<String, dynamic> json) => _$UpdateGameRequestFromJson(json);
  Map<String, dynamic> toJson() => _$UpdateGameRequestToJson(this);
}

@JsonSerializable()
class GameEventRequest {
  @JsonKey(name: 'event_type')
  final String eventType;
  @JsonKey(name: 'event_image_name')
  final String? eventImageName;
  final String? comment;
  @JsonKey(name: 'damage_delta')
  final int damageDelta;
  @JsonKey(name: 'life_total_after')
  final int targetLifeTotalAfter;
  @JsonKey(name: 'source_ranking_id')
  final int? sourceRankingId;
  @JsonKey(name: 'target_ranking_id')
  final int? targetRankingId;

  const GameEventRequest({
    required this.eventType,
    this.eventImageName,
    this.comment,
    required this.damageDelta,
    required this.targetLifeTotalAfter,
    this.sourceRankingId,
    this.targetRankingId,
  });

  factory GameEventRequest.fromJson(Map<String, dynamic> json) => _$GameEventRequestFromJson(json);
  Map<String, dynamic> toJson() => _$GameEventRequestToJson(this);
}