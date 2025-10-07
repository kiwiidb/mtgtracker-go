// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'deck.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Deck _$DeckFromJson(Map<String, dynamic> json) => Deck(
      id: (json['id'] as num?)?.toInt(),
      commander: json['commander'] as String,
      crop: json['crop'] as String,
      secondaryImg: json['secondary_image'] as String,
      image: json['image'] as String,
      colors:
          (json['colors'] as List<dynamic>?)?.map((e) => e as String).toList(),
      moxfieldUrl: json['moxfield_url'] as String?,
      bracket: (json['bracket'] as num?)?.toInt(),
    );

Map<String, dynamic> _$DeckToJson(Deck instance) => <String, dynamic>{
      'id': instance.id,
      'commander': instance.commander,
      'crop': instance.crop,
      'secondary_image': instance.secondaryImg,
      'image': instance.image,
      'colors': instance.colors,
      'moxfield_url': instance.moxfieldUrl,
      'bracket': instance.bracket,
    };

DeckWithCount _$DeckWithCountFromJson(Map<String, dynamic> json) =>
    DeckWithCount(
      deck: Deck.fromJson(json['deck'] as Map<String, dynamic>),
      count: (json['count'] as num).toInt(),
      wins: (json['wins'] as num).toInt(),
    );

Map<String, dynamic> _$DeckWithCountToJson(DeckWithCount instance) =>
    <String, dynamic>{
      'deck': instance.deck,
      'count': instance.count,
      'wins': instance.wins,
    };
