// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'opponent.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Opponent _$OpponentFromJson(Map<String, dynamic> json) => Opponent(
      id: (json['id'] as num).toInt(),
      player1: Player.fromJson(json['player1'] as Map<String, dynamic>),
      player2: Player.fromJson(json['player2'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$OpponentToJson(Opponent instance) => <String, dynamic>{
      'id': instance.id,
      'player1': instance.player1,
      'player2': instance.player2,
    };
