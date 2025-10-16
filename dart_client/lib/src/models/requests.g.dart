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

CreateRankingRequest _$CreateRankingRequestFromJson(
        Map<String, dynamic> json) =>
    CreateRankingRequest(
      playerId: json['player_id'] as String?,
      deckId: (json['deck_id'] as num?)?.toInt(),
      deck: json['deck'] == null
          ? null
          : Deck.fromJson(json['deck'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$CreateRankingRequestToJson(
        CreateRankingRequest instance) =>
    <String, dynamic>{
      'player_id': instance.playerId,
      'deck_id': instance.deckId,
      'deck': instance.deck,
    };

CreateGameRequest _$CreateGameRequestFromJson(Map<String, dynamic> json) =>
    CreateGameRequest(
      duration: (json['duration'] as num?)?.toInt(),
      date:
          json['date'] == null ? null : DateTime.parse(json['date'] as String),
      comments: json['comments'] as String,
      finished: json['finished'] as bool,
      rankings: (json['rankings'] as List<dynamic>)
          .map((e) => CreateRankingRequest.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$CreateGameRequestToJson(CreateGameRequest instance) =>
    <String, dynamic>{
      'duration': instance.duration,
      'date': instance.date?.toIso8601String(),
      'comments': instance.comments,
      'finished': instance.finished,
      'rankings': instance.rankings,
    };

UpdateRanking _$UpdateRankingFromJson(Map<String, dynamic> json) =>
    UpdateRanking(
      rankingId: (json['ranking_id'] as num).toInt(),
      position: (json['position'] as num).toInt(),
      description: json['description'] == null
          ? null
          : GameDescription.fromJson(
              json['description'] as Map<String, dynamic>),
      startingPlayer: json['starting_player'] as bool?,
      playerId: json['player_id'] as String?,
    );

Map<String, dynamic> _$UpdateRankingToJson(UpdateRanking instance) =>
    <String, dynamic>{
      'ranking_id': instance.rankingId,
      'position': instance.position,
      'description': instance.description,
      'starting_player': instance.startingPlayer,
      'player_id': instance.playerId,
    };

UpdateGameRequest _$UpdateGameRequestFromJson(Map<String, dynamic> json) =>
    UpdateGameRequest(
      gameId: (json['game_id'] as num).toInt(),
      finished: json['finished'] as bool?,
      rankings: (json['rankings'] as List<dynamic>)
          .map((e) => UpdateRanking.fromJson(e as Map<String, dynamic>))
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

ProfileImageUploadUrlRequest _$ProfileImageUploadUrlRequestFromJson(
        Map<String, dynamic> json) =>
    ProfileImageUploadUrlRequest(
      fileName: json['file_name'] as String,
    );

Map<String, dynamic> _$ProfileImageUploadUrlRequestToJson(
        ProfileImageUploadUrlRequest instance) =>
    <String, dynamic>{
      'file_name': instance.fileName,
    };

ProfileImageUploadUrlResponse _$ProfileImageUploadUrlResponseFromJson(
        Map<String, dynamic> json) =>
    ProfileImageUploadUrlResponse(
      uploadUrl: json['upload_url'] as String,
      imageUrl: json['image_url'] as String,
    );

Map<String, dynamic> _$ProfileImageUploadUrlResponseToJson(
        ProfileImageUploadUrlResponse instance) =>
    <String, dynamic>{
      'upload_url': instance.uploadUrl,
      'image_url': instance.imageUrl,
    };

UpdateProfileImageRequest _$UpdateProfileImageRequestFromJson(
        Map<String, dynamic> json) =>
    UpdateProfileImageRequest(
      imageUrl: json['image_url'] as String,
    );

Map<String, dynamic> _$UpdateProfileImageRequestToJson(
        UpdateProfileImageRequest instance) =>
    <String, dynamic>{
      'image_url': instance.imageUrl,
    };

CreateDeckRequest _$CreateDeckRequestFromJson(Map<String, dynamic> json) =>
    CreateDeckRequest(
      moxfieldUrl: json['moxfield_url'] as String?,
      themes:
          (json['themes'] as List<dynamic>).map((e) => e as String).toList(),
      bracket: (json['bracket'] as num?)?.toInt(),
      commander: json['commander'] as String,
      colors:
          (json['colors'] as List<dynamic>).map((e) => e as String).toList(),
      image: json['image'] as String,
      secondaryImage: json['secondary_image'] as String,
      crop: json['crop'] as String,
    );

Map<String, dynamic> _$CreateDeckRequestToJson(CreateDeckRequest instance) =>
    <String, dynamic>{
      'moxfield_url': instance.moxfieldUrl,
      'themes': instance.themes,
      'bracket': instance.bracket,
      'commander': instance.commander,
      'colors': instance.colors,
      'image': instance.image,
      'secondary_image': instance.secondaryImage,
      'crop': instance.crop,
    };

SearchGamesRequest _$SearchGamesRequestFromJson(Map<String, dynamic> json) =>
    SearchGamesRequest(
      playerIds: (json['player_ids'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
      commanders: (json['commanders'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
      allPlayers: (json['all_players'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
      allCommanders: (json['all_commanders'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
    );

Map<String, dynamic> _$SearchGamesRequestToJson(SearchGamesRequest instance) =>
    <String, dynamic>{
      'player_ids': instance.playerIds,
      'commanders': instance.commanders,
      'all_players': instance.allPlayers,
      'all_commanders': instance.allCommanders,
    };

UpdateRankingRequest _$UpdateRankingRequestFromJson(
        Map<String, dynamic> json) =>
    UpdateRankingRequest(
      description: json['description'] == null
          ? null
          : GameDescription.fromJson(
              json['description'] as Map<String, dynamic>),
      startingPlayer: json['starting_player'] as bool?,
      couldHaveWon: json['could_have_won'] as bool?,
      earlySolRing: json['early_sol_ring'] as bool?,
    );

Map<String, dynamic> _$UpdateRankingRequestToJson(
        UpdateRankingRequest instance) =>
    <String, dynamic>{
      'description': instance.description,
      'starting_player': instance.startingPlayer,
      'could_have_won': instance.couldHaveWon,
      'early_sol_ring': instance.earlySolRing,
    };
