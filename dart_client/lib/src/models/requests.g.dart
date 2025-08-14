// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'requests.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SignupPlayerRequest _$SignupPlayerRequestFromJson(Map<String, dynamic> json) =>
    SignupPlayerRequest(
      name: json['name'] as String,
    );

Map<String, dynamic> _$SignupPlayerRequestToJson(
        SignupPlayerRequest instance) =>
    <String, dynamic>{
      'name': instance.name,
    };

CreateGameRequest _$CreateGameRequestFromJson(Map<String, dynamic> json) =>
    CreateGameRequest(
      duration: (json['duration'] as num?)?.toInt(),
      date:
          json['date'] == null ? null : DateTime.parse(json['date'] as String),
      comments: json['comments'] as String,
      image: json['image'] as String,
      finished: json['finished'] as bool,
      rankings: (json['rankings'] as List<dynamic>)
          .map((e) => Ranking.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$CreateGameRequestToJson(CreateGameRequest instance) =>
    <String, dynamic>{
      'duration': instance.duration,
      'date': instance.date?.toIso8601String(),
      'comments': instance.comments,
      'image': instance.image,
      'finished': instance.finished,
      'rankings': instance.rankings,
    };

UpdateGameRequest _$UpdateGameRequestFromJson(Map<String, dynamic> json) =>
    UpdateGameRequest(
      gameId: (json['game_id'] as num).toInt(),
      finished: json['finished'] as bool?,
      rankings: (json['rankings'] as List<dynamic>)
          .map((e) => Ranking.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$UpdateGameRequestToJson(UpdateGameRequest instance) =>
    <String, dynamic>{
      'game_id': instance.gameId,
      'finished': instance.finished,
      'rankings': instance.rankings,
    };

GameEventRequest _$GameEventRequestFromJson(Map<String, dynamic> json) =>
    GameEventRequest(
      eventType: json['event_type'] as String,
      eventImageName: json['event_image_name'] as String?,
      comment: json['comment'] as String?,
      damageDelta: (json['damage_delta'] as num).toInt(),
      targetLifeTotalAfter: (json['life_total_after'] as num).toInt(),
      sourceRankingId: (json['source_ranking_id'] as num?)?.toInt(),
      targetRankingId: (json['target_ranking_id'] as num?)?.toInt(),
    );

Map<String, dynamic> _$GameEventRequestToJson(GameEventRequest instance) =>
    <String, dynamic>{
      'event_type': instance.eventType,
      'event_image_name': instance.eventImageName,
      'comment': instance.comment,
      'damage_delta': instance.damageDelta,
      'life_total_after': instance.targetLifeTotalAfter,
      'source_ranking_id': instance.sourceRankingId,
      'target_ranking_id': instance.targetRankingId,
    };
