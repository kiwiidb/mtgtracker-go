import 'package:json_annotation/json_annotation.dart';
import 'deck.dart';
import 'player.dart';

part 'ranking.g.dart';

@JsonSerializable()
class CardReference {
  final String name;
  @JsonKey(name: 'oracle_text')
  final String oracleText;
  @JsonKey(name: 'image_uri')
  final String? imageUri;
  @JsonKey(name: 'art_crop_uri')
  final String? artCropUri;
  @JsonKey(name: 'secondary_image_uri')
  final String? secondaryImageUri;
  @JsonKey(name: 'secondary_art_crop_uri')
  final String? secondaryArtCropUri;
  @JsonKey(name: 'color_identity')
  final List<String> colorIdentity;

  CardReference({
    required this.name,
    required this.oracleText,
    this.imageUri,
    this.artCropUri,
    this.secondaryImageUri,
    this.secondaryArtCropUri,
    required this.colorIdentity,
  });

  factory CardReference.fromJson(Map<String, dynamic> json) =>
      _$CardReferenceFromJson(json);
  Map<String, dynamic> toJson() => _$CardReferenceToJson(this);
}

@JsonSerializable()
class PlayerReference {
  final String id;
  final String name;

  PlayerReference({
    required this.id,
    required this.name,
  });

  factory PlayerReference.fromJson(Map<String, dynamic> json) =>
      _$PlayerReferenceFromJson(json);
  Map<String, dynamic> toJson() => _$PlayerReferenceToJson(this);
}

@JsonSerializable()
class GameDescription {
  final String text;
  @JsonKey(name: 'card_references')
  final Map<String, CardReference> cardReferences;
  @JsonKey(name: 'player_references')
  final List<PlayerReference> playerReferences;

  GameDescription({
    required this.text,
    required this.cardReferences,
    required this.playerReferences,
  });

  factory GameDescription.fromJson(Map<String, dynamic> json) =>
      _$GameDescriptionFromJson(json);
  Map<String, dynamic> toJson() => _$GameDescriptionToJson(this);
}

@JsonSerializable()
class Ranking {
  final int id;
  @JsonKey(name: 'player_id')
  final String? playerId;
  int position;
  @JsonKey(name: 'life_total')
  final int? lifeTotal;
  @JsonKey(name: 'last_life_total')
  final int? lastLifeTotal;
  @JsonKey(name: 'last_life_total_timestamp')
  final DateTime? lastLifeTotalTimestamp;
  final Deck deck;
  final Player? player;
  final String? status;
  final GameDescription? description;

  Ranking({
    required this.id,
    this.playerId,
    required this.position,
    this.lifeTotal,
    this.lastLifeTotal,
    this.lastLifeTotalTimestamp,
    required this.deck,
    this.player,
    this.status,
    this.description,
  });

  factory Ranking.fromJson(Map<String, dynamic> json) =>
      _$RankingFromJson(json);
  Map<String, dynamic> toJson() => _$RankingToJson(this);
}
