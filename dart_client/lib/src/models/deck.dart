import 'package:json_annotation/json_annotation.dart';

part 'deck.g.dart';

@JsonSerializable()
class Deck {
  final int? id;
  final String commander;
  final String crop;
  @JsonKey(name: 'secondary_image')
  final String secondaryImg;
  final String image;
  final List<String>? colors; // Scryfall color codes: W, U, B, R, G, C

  const Deck({
    this.id,
    required this.commander,
    required this.crop,
    required this.secondaryImg,
    required this.image,
    this.colors,
  });

  factory Deck.fromJson(Map<String, dynamic> json) => _$DeckFromJson(json);
  Map<String, dynamic> toJson() => _$DeckToJson(this);
}

@JsonSerializable()
class DeckWithCount {
  final Deck deck;
  final int count;
  final int wins;

  const DeckWithCount({
    required this.deck,
    required this.count,
    required this.wins,
  });

  factory DeckWithCount.fromJson(Map<String, dynamic> json) =>
      _$DeckWithCountFromJson(json);
  Map<String, dynamic> toJson() => _$DeckWithCountToJson(this);
}
