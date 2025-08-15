// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'ranking.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Ranking _$RankingFromJson(Map<String, dynamic> json) => Ranking(
      id: (json['id'] as num).toInt(),
      playerId: (json['player_id'] as num).toInt(),
      position: (json['position'] as num).toInt(),
      lifeTotal: (json['life_total'] as num?)?.toInt(),
      deck: Deck.fromJson(json['deck'] as Map<String, dynamic>),
      player: json['player'] == null
          ? null
          : Player.fromJson(json['player'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$RankingToJson(Ranking instance) => <String, dynamic>{
      'id': instance.id,
      'player_id': instance.playerId,
      'position': instance.position,
      'life_total': instance.lifeTotal,
      'deck': instance.deck,
      'player': instance.player,
    };
