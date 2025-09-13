import 'package:json_annotation/json_annotation.dart';
import 'deck.dart';

part 'requests.g.dart';

@JsonSerializable()
class SignupPlayerRequest {
  final String name;

  const SignupPlayerRequest({
    required this.name,
  });

  factory SignupPlayerRequest.fromJson(Map<String, dynamic> json) =>
      _$SignupPlayerRequestFromJson(json);
  Map<String, dynamic> toJson() => _$SignupPlayerRequestToJson(this);
}

@JsonSerializable()
class CreateRankingRequest {
  @JsonKey(name: 'player_id')
  final String? playerId;
  final Deck deck;

  const CreateRankingRequest({
    this.playerId,
    required this.deck,
  });

  factory CreateRankingRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateRankingRequestFromJson(json);
  Map<String, dynamic> toJson() => _$CreateRankingRequestToJson(this);
}

@JsonSerializable()
class CreateGameRequest {
  final int? duration;
  final DateTime? date;
  final String comments;
  final bool finished;
  final List<CreateRankingRequest> rankings;

  const CreateGameRequest({
    this.duration,
    this.date,
    required this.comments,
    required this.finished,
    required this.rankings,
  });

  factory CreateGameRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateGameRequestFromJson(json);
  Map<String, dynamic> toJson() => _$CreateGameRequestToJson(this);
}

@JsonSerializable()
class UpdateRanking {
  @JsonKey(name: 'player_id')
  final String? playerId;
  final int position;

  const UpdateRanking({
    this.playerId,
    required this.position,
  });

  factory UpdateRanking.fromJson(Map<String, dynamic> json) =>
      _$UpdateRankingFromJson(json);
  Map<String, dynamic> toJson() => _$UpdateRankingToJson(this);
}

@JsonSerializable()
class UpdateGameRequest {
  @JsonKey(name: 'game_id')
  final int gameId;
  final bool? finished;
  final List<UpdateRanking> rankings;

  const UpdateGameRequest({
    required this.gameId,
    this.finished,
    required this.rankings,
  });

  factory UpdateGameRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateGameRequestFromJson(json);
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

  factory GameEventRequest.fromJson(Map<String, dynamic> json) =>
      _$GameEventRequestFromJson(json);
  Map<String, dynamic> toJson() => _$GameEventRequestToJson(this);
}

@JsonSerializable()
class ProfileImageUploadUrlRequest {
  @JsonKey(name: 'file_name')
  final String fileName;

  const ProfileImageUploadUrlRequest({
    required this.fileName,
  });

  factory ProfileImageUploadUrlRequest.fromJson(Map<String, dynamic> json) =>
      _$ProfileImageUploadUrlRequestFromJson(json);
  Map<String, dynamic> toJson() => _$ProfileImageUploadUrlRequestToJson(this);
}

@JsonSerializable()
class ProfileImageUploadUrlResponse {
  @JsonKey(name: 'upload_url')
  final String uploadUrl;
  @JsonKey(name: 'image_url')
  final String imageUrl;

  const ProfileImageUploadUrlResponse({
    required this.uploadUrl,
    required this.imageUrl,
  });

  factory ProfileImageUploadUrlResponse.fromJson(Map<String, dynamic> json) =>
      _$ProfileImageUploadUrlResponseFromJson(json);
  Map<String, dynamic> toJson() => _$ProfileImageUploadUrlResponseToJson(this);
}

@JsonSerializable()
class UpdateProfileImageRequest {
  @JsonKey(name: 'image_url')
  final String imageUrl;

  const UpdateProfileImageRequest({
    required this.imageUrl,
  });

  factory UpdateProfileImageRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateProfileImageRequestFromJson(json);
  Map<String, dynamic> toJson() => _$UpdateProfileImageRequestToJson(this);
}
