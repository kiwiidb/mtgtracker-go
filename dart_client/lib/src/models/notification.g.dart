// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'notification.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

MtgNotification _$MtgNotificationFromJson(Map<String, dynamic> json) =>
    MtgNotification(
      id: (json['id'] as num).toInt(),
      title: json['title'] as String,
      body: json['body'] as String,
      type: json['type'] as String,
      actions: (json['actions'] as List<dynamic>)
          .map((e) => $enumDecode(_$NotificationActionEnumMap, e))
          .toList(),
      read: json['read'] as bool,
      createdAt: DateTime.parse(json['created_at'] as String),
      gameId: (json['game_id'] as num?)?.toInt(),
      referredPlayerId: json['referred_player_id'] as String?,
      game: json['game'] == null
          ? null
          : Game.fromJson(json['game'] as Map<String, dynamic>),
      referredPlayer: json['referred_player'] == null
          ? null
          : Player.fromJson(json['referred_player'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$MtgNotificationToJson(MtgNotification instance) =>
    <String, dynamic>{
      'id': instance.id,
      'title': instance.title,
      'body': instance.body,
      'type': instance.type,
      'actions':
          instance.actions.map((e) => _$NotificationActionEnumMap[e]!).toList(),
      'read': instance.read,
      'created_at': instance.createdAt.toIso8601String(),
      'game_id': instance.gameId,
      'referred_player_id': instance.referredPlayerId,
      'game': instance.game,
      'referred_player': instance.referredPlayer,
    };

const _$NotificationActionEnumMap = {
  NotificationAction.deleteRanking: 'delete_ranking',
  NotificationAction.viewGame: 'view_game',
};
