// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'follow.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Follow _$FollowFromJson(Map<String, dynamic> json) => Follow(
      id: (json['id'] as num).toInt(),
      player1: Player.fromJson(json['player1'] as Map<String, dynamic>),
      player2: Player.fromJson(json['player2'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$FollowToJson(Follow instance) => <String, dynamic>{
      'id': instance.id,
      'player1': instance.player1,
      'player2': instance.player2,
    };
