// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'game.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Game _$GameFromJson(Map<String, dynamic> json) => Game(
      id: (json['id'] as num).toInt(),
      duration: (json['duration'] as num?)?.toInt(),
      date:
          json['date'] == null ? null : DateTime.parse(json['date'] as String),
      comments: json['comments'] as String?,
      image: json['image'] as String?,
      rankings: (json['rankings'] as List<dynamic>)
          .map((e) => Ranking.fromJson(e as Map<String, dynamic>))
          .toList(),
      finished: json['finished'] as bool?,
      gameEvents: (json['GameEvents'] as List<dynamic>?)
          ?.map((e) => GameEvent.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$GameToJson(Game instance) => <String, dynamic>{
      'id': instance.id,
      'duration': instance.duration,
      'date': instance.date?.toIso8601String(),
      'comments': instance.comments,
      'image': instance.image,
      'rankings': instance.rankings,
      'finished': instance.finished,
      'GameEvents': instance.gameEvents,
    };
