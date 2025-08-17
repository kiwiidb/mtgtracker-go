import 'dart:math';
import 'models/models.dart';

class MockGameGenerator {
  static final Random _random = Random();
  
  static final List<MockCommander> _commanders = [
    MockCommander(
      name: 'Teysa Karlov',
      crop: 'https://cards.scryfall.io/art_crop/front/c/d/cd14f1ce-7fcd-485c-b7ca-01c5b45fdc01.jpg?1689999296',
      image: 'https://cards.scryfall.io/normal/front/c/d/cd14f1ce-7fcd-485c-b7ca-01c5b45fdc01.jpg?1689999296',
    ),
    MockCommander(
      name: 'Ojer Axonil, Deepest Might',
      crop: 'https://cards.scryfall.io/art_crop/front/5/0/50f8e2b6-98c7-4f28-bb39-e1fbe841f1ee.jpg?1699044315',
      secondaryImg: 'https://cards.scryfall.io/art_crop/back/5/0/50f8e2b6-98c7-4f28-bb39-e1fbe841f1ee.jpg?1699044315',
      image: 'https://cards.scryfall.io/normal/front/5/0/50f8e2b6-98c7-4f28-bb39-e1fbe841f1ee.jpg?1699044315',
    ),
    MockCommander(
      name: 'Queen Marchesa',
      crop: 'https://cards.scryfall.io/art_crop/front/0/f/0fdae05f-7bdc-45fb-b9b9-e5ec3766f965.jpg?1712354769',
      image: 'https://cards.scryfall.io/normal/front/0/f/0fdae05f-7bdc-45fb-b9b9-e5ec3766f965.jpg?1712354769',
    ),
    MockCommander(
      name: 'Lord Windgrace',
      crop: 'https://cards.scryfall.io/art_crop/front/2/1/213d6fb8-5624-4804-b263-51f339482754.jpg?1592710275',
      image: 'https://cards.scryfall.io/normal/front/2/1/213d6fb8-5624-4804-b263-51f339482754.jpg?1592710275',
    ),
    MockCommander(
      name: 'Tasigur, the Golden Fang',
      crop: 'https://cards.scryfall.io/art_crop/front/1/7/175ad810-3cdd-43c7-99a9-8a2e8ad6dbae.jpg?1743206587',
      image: 'https://cards.scryfall.io/normal/front/1/7/175ad810-3cdd-43c7-99a9-8a2e8ad6dbae.jpg?1743206587',
    ),
    MockCommander(
      name: 'Borborygmos Enraged',
      crop: 'https://cards.scryfall.io/art_crop/front/0/2/02df18c5-07b9-4e85-9c4b-58b63fa59437.jpg?1702429604',
      image: 'https://cards.scryfall.io/normal/front/0/2/02df18c5-07b9-4e85-9c4b-58b63fa59437.jpg?1702429604',
    ),
    MockCommander(
      name: 'Anzrag, the Quake-Mole',
      crop: 'https://cards.scryfall.io/art_crop/front/7/0/70e9d8b8-4b32-4414-b32f-1f47523239c5.jpg?1706242114',
      image: 'https://cards.scryfall.io/normal/front/7/0/70e9d8b8-4b32-4414-b32f-1f47523239c5.jpg?1706242114',
    ),
    MockCommander(
      name: 'The Mimeoplasm',
      crop: 'https://cards.scryfall.io/art_crop/front/9/9/998b86f3-e53f-4ebb-b111-6f14577fded1.jpg?1712354748',
      image: 'https://cards.scryfall.io/normal/front/9/9/998b86f3-e53f-4ebb-b111-6f14577fded1.jpg?1712354748',
    ),
  ];

  /// Generates a list of mock rankings with randomized commanders
  /// [playerCount] must be between 2 and 4
  static List<Ranking> generateMockRankings(int playerCount) {
    if (playerCount < 2 || playerCount > 4) {
      throw ArgumentError('Player count must be between 2 and 4');
    }

    // Shuffle commanders and take the required number
    final shuffledCommanders = List<MockCommander>.from(_commanders)..shuffle(_random);
    final selectedCommanders = shuffledCommanders.take(playerCount).toList();

    final rankings = <Ranking>[];
    for (int i = 0; i < playerCount; i++) {
      final commander = selectedCommanders[i];
      rankings.add(Ranking(
        id: 0,
        playerId: i + 1,
        position: i + 1,
        lifeTotal: 40,
        deck: Deck(
          commander: commander.name,
          crop: commander.crop,
          secondaryImg: commander.secondaryImg,
          image: commander.image,
        ),
        player: null,
      ));
    }

    return rankings;
  }

  /// Gets a random commander for the game image
  static String getRandomGameImage() {
    final randomCommander = _commanders[_random.nextInt(_commanders.length)];
    return randomCommander.crop;
  }
}

class MockCommander {
  final String name;
  final String crop;
  final String image;
  final String secondaryImg;

  const MockCommander({
    required this.name,
    required this.crop,
    required this.image,
    this.secondaryImg = '',
  });
}