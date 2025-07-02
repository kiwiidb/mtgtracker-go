import 'package:json_annotation/json_annotation.dart';

part 'deck.g.dart';

@JsonSerializable()
class Deck {
  @JsonKey(name: 'ID')
  final int id;
  final String commander;
  final String crop;
  @JsonKey(name: 'secondary_image')
  final String secondaryImg;
  final String image;

  const Deck({
    required this.id,
    required this.commander,
    required this.crop,
    required this.secondaryImg,
    required this.image,
  });

  factory Deck.fromJson(Map<String, dynamic> json) => _$DeckFromJson(json);
  Map<String, dynamic> toJson() => _$DeckToJson(this);
}

@JsonSerializable()
class DeckWithCount {
  final Deck deck;
  final int count;

  const DeckWithCount({
    required this.deck,
    required this.count,
  });

  factory DeckWithCount.fromJson(Map<String, dynamic> json) => _$DeckWithCountFromJson(json);
  Map<String, dynamic> toJson() => _$DeckWithCountToJson(this);
}