import 'dart:convert';
import 'package:http/http.dart' as http;
import 'models/models.dart';

class MTGTrackerClient {
  final String baseUrl;
  final String? authToken;
  final http.Client _httpClient;

  MTGTrackerClient({
    required this.baseUrl,
    this.authToken,
    http.Client? httpClient,
  }) : _httpClient = httpClient ?? http.Client();

  Map<String, String> get _headers {
    final headers = <String, String>{
      'Content-Type': 'application/json',
    };

    if (authToken != null) {
      headers['Authorization'] = 'Bearer $authToken';
    }

    return headers;
  }

  Future<T> _handleResponse<T>(
    http.Response response,
    T Function(Map<String, dynamic>) fromJson,
  ) async {
    if (response.statusCode >= 200 && response.statusCode < 300) {
      final Map<String, dynamic> json = jsonDecode(response.body);
      return fromJson(json);
    } else {
      throw MTGTrackerException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }

  Future<PaginatedResult<T>> _handlePaginatedResponse<T>(
    http.Response response,
    T Function(Map<String, dynamic>) fromJson,
  ) async {
    if (response.statusCode >= 200 && response.statusCode < 300) {
      final Map<String, dynamic> json = jsonDecode(response.body);
      return PaginatedResult.fromJson(json, (itemJson) {
        return fromJson(itemJson as Map<String, dynamic>);
      });
    } else {
      throw MTGTrackerException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }

  // Player endpoints
  Future<Player> signupPlayer(SignupPlayerRequest request) async {
    final response = await _httpClient.post(
      Uri.parse('$baseUrl/player/v1/signup'),
      headers: _headers,
      body: jsonEncode(request.toJson()),
    );
    return _handleResponse(response, Player.fromJson);
  }

  Future<PaginatedResult<Player>> getPlayers({
    String? search,
    int? page,
    int? perPage,
  }) async {
    final queryParams = <String, String>{};
    if (search != null) queryParams['search'] = search;
    if (page != null) queryParams['page'] = page.toString();
    if (perPage != null) queryParams['per_page'] = perPage.toString();

    final uri = Uri.parse('$baseUrl/player/v1/players')
        .replace(queryParameters: queryParams.isEmpty ? null : queryParams);

    final response = await _httpClient.get(uri, headers: _headers);
    return _handlePaginatedResponse(response, Player.fromJson);
  }

  Future<Player> getPlayer(String playerId) async {
    final response = await _httpClient.get(
      Uri.parse('$baseUrl/player/v1/players/$playerId'),
      headers: _headers,
    );
    return _handleResponse(response, Player.fromJson);
  }

  Future<Player> getMyPlayer() async {
    final response = await _httpClient.get(
      Uri.parse('$baseUrl/player/v1/me'),
      headers: _headers,
    );
    return _handleResponse(response, Player.fromJson);
  }

  Future<PaginatedResult<Deck>> getPlayerDecks(
    String playerId, {
    int? page,
    int? perPage,
  }) async {
    final queryParams = <String, String>{};
    if (page != null) queryParams['page'] = page.toString();
    if (perPage != null) queryParams['per_page'] = perPage.toString();

    final uri = Uri.parse('$baseUrl/player/v1/players/$playerId/decks')
        .replace(queryParameters: queryParams.isEmpty ? null : queryParams);

    final response = await _httpClient.get(uri, headers: _headers);
    return _handlePaginatedResponse(response, Deck.fromJson);
  }

  Future<PaginatedResult<Game>> getPlayerGames(
    String playerId, {
    int? page,
    int? perPage,
  }) async {
    final queryParams = <String, String>{};
    if (page != null) queryParams['page'] = page.toString();
    if (perPage != null) queryParams['per_page'] = perPage.toString();

    final uri = Uri.parse('$baseUrl/player/v1/players/$playerId/games')
        .replace(queryParameters: queryParams.isEmpty ? null : queryParams);

    final response = await _httpClient.get(uri, headers: _headers);
    return _handlePaginatedResponse(response, Game.fromJson);
  }

  Future<Deck> createDeck(CreateDeckRequest request) async {
    final response = await _httpClient.post(
      Uri.parse('$baseUrl/deck/v1/decks'),
      headers: _headers,
      body: jsonEncode(request.toJson()),
    );
    return _handleResponse(response, Deck.fromJson);
  }

  Future<List<String>> getThemes() async {
    final response = await _httpClient.get(
      Uri.parse('$baseUrl/moxfield/v1/themes'),
      headers: _headers,
    );

    if (response.statusCode >= 200 && response.statusCode < 300) {
      final List<dynamic> jsonList = jsonDecode(response.body);
      return jsonList.cast<String>();
    } else {
      throw MTGTrackerException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }

  /// Get a signed URL for uploading a profile image to S3
  Future<ProfileImageUploadUrlResponse> getProfileImageUploadUrl(
      ProfileImageUploadUrlRequest request) async {
    final response = await _httpClient.post(
      Uri.parse('$baseUrl/player/v1/profile-image/upload-url'),
      headers: _headers,
      body: jsonEncode(request.toJson()),
    );
    return _handleResponse(response, ProfileImageUploadUrlResponse.fromJson);
  } // Game endpoints

  Future<Game> createGame(CreateGameRequest request) async {
    final response = await _httpClient.post(
      Uri.parse('$baseUrl/game/v1/games'),
      headers: _headers,
      body: jsonEncode(request.toJson()),
    );
    return _handleResponse(response, Game.fromJson);
  }

  Future<PaginatedResult<Game>> getGames({
    int? page,
    int? perPage,
  }) async {
    final queryParams = <String, String>{};
    if (page != null) queryParams['page'] = page.toString();
    if (perPage != null) queryParams['per_page'] = perPage.toString();

    final uri = Uri.parse('$baseUrl/game/v1/games')
        .replace(queryParameters: queryParams.isEmpty ? null : queryParams);

    final response = await _httpClient.get(uri, headers: _headers);
    return _handlePaginatedResponse(response, Game.fromJson);
  }

  Future<PaginatedResult<Game>> searchGames(
    SearchGamesRequest request, {
    int? page,
    int? perPage,
  }) async {
    final queryParams = <String, String>{};
    if (page != null) queryParams['page'] = page.toString();
    if (perPage != null) queryParams['per_page'] = perPage.toString();

    final uri = Uri.parse('$baseUrl/game/v1/games/search')
        .replace(queryParameters: queryParams.isEmpty ? null : queryParams);

    final response = await _httpClient.post(
      uri,
      headers: _headers,
      body: jsonEncode(request.toJson()),
    );
    return _handlePaginatedResponse(response, Game.fromJson);
  }

  Future<Game?> getActiveGame() async {
    final response = await _httpClient.get(
      Uri.parse('$baseUrl/game/v1/games/active'),
      headers: _headers,
    );

    // 204 No Content means no active game
    if (response.statusCode == 204) {
      return null;
    }

    return _handleResponse(response, Game.fromJson);
  }

  Future<Game> getGame(int gameId) async {
    final response = await _httpClient.get(
      Uri.parse('$baseUrl/game/v1/games/$gameId'),
      headers: _headers,
    );
    return _handleResponse(response, Game.fromJson);
  }

  Future<Game> updateGame(int gameId, UpdateGameRequest request) async {
    final response = await _httpClient.put(
      Uri.parse('$baseUrl/game/v1/games/$gameId'),
      headers: _headers,
      body: jsonEncode(request.toJson()),
    );
    return _handleResponse(response, Game.fromJson);
  }

  Future<void> deleteGame(int gameId) async {
    final response = await _httpClient.delete(
      Uri.parse('$baseUrl/game/v1/games/$gameId'),
      headers: _headers,
    );

    if (response.statusCode != 204) {
      throw MTGTrackerException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }

  Future<GameEvent> addGameEvent(int gameId, GameEventRequest request) async {
    final response = await _httpClient.post(
      Uri.parse('$baseUrl/game/v1/games/$gameId/events'),
      headers: _headers,
      body: jsonEncode(request.toJson()),
    );
    return _handleResponse(response, GameEvent.fromJson);
  }

  // Notification endpoints
  Future<PaginatedResult<MtgNotification>> getNotifications({
    bool? read,
    int? page,
    int? perPage,
  }) async {
    final queryParams = <String, String>{};
    if (read != null) queryParams['read'] = read.toString();
    if (page != null) queryParams['page'] = page.toString();
    if (perPage != null) queryParams['per_page'] = perPage.toString();

    final uri = Uri.parse('$baseUrl/notification/v1/notifications')
        .replace(queryParameters: queryParams.isEmpty ? null : queryParams);

    final response = await _httpClient.get(uri, headers: _headers);
    return _handlePaginatedResponse(response, MtgNotification.fromJson);
  }

  Future<void> markNotificationAsRead(int notificationId) async {
    final response = await _httpClient.put(
      Uri.parse('$baseUrl/notification/v1/notifications/$notificationId/read'),
      headers: _headers,
    );

    if (response.statusCode != 204) {
      throw MTGTrackerException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }

  // Ranking endpoints
  Future<Ranking> updateRanking(
      int rankingId, UpdateRankingRequest request) async {
    final response = await _httpClient.put(
      Uri.parse('$baseUrl/ranking/v1/rankings/$rankingId'),
      headers: _headers,
      body: jsonEncode(request.toJson()),
    );

    return _handleResponse(response, Ranking.fromJson);
  }

  Future<void> deleteRanking(int rankingId) async {
    final response = await _httpClient.delete(
      Uri.parse('$baseUrl/ranking/v1/rankings/$rankingId'),
      headers: _headers,
    );

    if (response.statusCode != 204) {
      throw MTGTrackerException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }

  // Opponent endpoints
  Future<Opponent> createOpponent(String playerId) async {
    final response = await _httpClient.post(
      Uri.parse('$baseUrl/opponent/v1/opponents/$playerId'),
      headers: _headers,
    );
    return _handleResponse(response, Opponent.fromJson);
  }

  Future<void> deleteOpponent(String playerId) async {
    final response = await _httpClient.delete(
      Uri.parse('$baseUrl/opponent/v1/opponents/$playerId'),
      headers: _headers,
    );

    if (response.statusCode != 204) {
      throw MTGTrackerException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }

  Future<PaginatedResult<PlayerOpponentWithCount>> getMyOpponents({
    int? page,
    int? perPage,
  }) async {
    final queryParams = <String, String>{};
    if (page != null) queryParams['page'] = page.toString();
    if (perPage != null) queryParams['per_page'] = perPage.toString();

    final uri = Uri.parse('$baseUrl/opponent/v1/opponents')
        .replace(queryParameters: queryParams.isEmpty ? null : queryParams);

    final response = await _httpClient.get(uri, headers: _headers);
    return _handlePaginatedResponse(response, PlayerOpponentWithCount.fromJson);
  }

  Future<PaginatedResult<PlayerOpponentWithCount>> getPlayerOpponents(
    String playerId, {
    int? page,
    int? perPage,
  }) async {
    final queryParams = <String, String>{};
    if (page != null) queryParams['page'] = page.toString();
    if (perPage != null) queryParams['per_page'] = perPage.toString();

    final uri = Uri.parse('$baseUrl/opponent/v1/players/$playerId/opponents')
        .replace(queryParameters: queryParams.isEmpty ? null : queryParams);

    final response = await _httpClient.get(uri, headers: _headers);
    return _handlePaginatedResponse(response, PlayerOpponentWithCount.fromJson);
  }

  // Feed endpoints
  Future<PaginatedResult<Game>> getFeedItems({
    int? page,
    int? perPage,
  }) async {
    final queryParams = <String, String>{};
    if (page != null) queryParams['page'] = page.toString();
    if (perPage != null) queryParams['per_page'] = perPage.toString();

    final uri = Uri.parse('$baseUrl/feed/v1/items')
        .replace(queryParameters: queryParams.isEmpty ? null : queryParams);

    final response = await _httpClient.get(uri, headers: _headers);
    return _handlePaginatedResponse(response, Game.fromJson);
  }

  // Push notification endpoints
  Future<List<DeviceToken>> getMyPushTokens() async {
    final response = await _httpClient.get(
      Uri.parse('$baseUrl/push/v1/tokens'),
      headers: _headers,
    );

    if (response.statusCode >= 200 && response.statusCode < 300) {
      final List<dynamic> jsonList = jsonDecode(response.body);
      return jsonList
          .map((json) => DeviceToken.fromJson(json as Map<String, dynamic>))
          .toList();
    } else {
      throw MTGTrackerException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }

  Future<void> registerPushToken(String token, String platform) async {
    final response = await _httpClient.post(
      Uri.parse('$baseUrl/push/v1/tokens'),
      headers: _headers,
      body: jsonEncode({
        'token': token,
        'platform': platform,
      }),
    );

    if (response.statusCode != 204) {
      throw MTGTrackerException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }

  Future<void> unregisterPushToken(String token) async {
    final response = await _httpClient.delete(
      Uri.parse('$baseUrl/push/v1/tokens'),
      headers: _headers,
      body: jsonEncode({
        'token': token,
      }),
    );

    if (response.statusCode != 204) {
      throw MTGTrackerException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }

  Future<void> sendTestNotification({
    required String title,
    required String body,
    String? imageUrl,
  }) async {
    final requestBody = <String, dynamic>{
      'title': title,
      'body': body,
    };
    if (imageUrl != null) {
      requestBody['image_url'] = imageUrl;
    }

    final response = await _httpClient.post(
      Uri.parse('$baseUrl/push/v1/debug/send'),
      headers: _headers,
      body: jsonEncode(requestBody),
    );

    if (response.statusCode != 204) {
      throw MTGTrackerException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }

  void dispose() {
    _httpClient.close();
  }
}

class MTGTrackerException implements Exception {
  final int statusCode;
  final String message;

  const MTGTrackerException({
    required this.statusCode,
    required this.message,
  });

  @override
  String toString() {
    return 'MTGTrackerException: $statusCode - $message';
  }
}
