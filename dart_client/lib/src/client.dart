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

  Future<List<T>> _handleListResponse<T>(
    http.Response response,
    T Function(Map<String, dynamic>) fromJson,
  ) async {
    if (response.statusCode >= 200 && response.statusCode < 300) {
      final List<dynamic> jsonList = jsonDecode(response.body);
      return jsonList.map((json) => fromJson(json as Map<String, dynamic>)).toList();
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

  Future<List<Player>> getPlayers({String? search}) async {
    final uri = Uri.parse('$baseUrl/player/v1/players');
    final finalUri = search != null 
        ? uri.replace(queryParameters: {'search': search})
        : uri;
    
    final response = await _httpClient.get(finalUri, headers: _headers);
    return _handleListResponse(response, Player.fromJson);
  }

  Future<Player> getPlayer(int playerId) async {
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

  // Game endpoints
  Future<Game> createGame(CreateGameRequest request) async {
    final response = await _httpClient.post(
      Uri.parse('$baseUrl/game/v1/games'),
      headers: _headers,
      body: jsonEncode(request.toJson()),
    );
    return _handleResponse(response, Game.fromJson);
  }

  Future<List<Game>> getGames() async {
    final response = await _httpClient.get(
      Uri.parse('$baseUrl/game/v1/games'),
      headers: _headers,
    );
    return _handleListResponse(response, Game.fromJson);
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

  // Ranking endpoints
  Future<List<Ranking>> getPendingRankings() async {
    final response = await _httpClient.get(
      Uri.parse('$baseUrl/ranking/v1/rankings/pending'),
      headers: _headers,
    );
    return _handleListResponse(response, Ranking.fromJson);
  }

  Future<void> acceptRanking(int rankingId) async {
    final response = await _httpClient.put(
      Uri.parse('$baseUrl/ranking/v1/rankings/$rankingId/accept'),
      headers: _headers,
    );
    
    if (response.statusCode != 204) {
      throw MTGTrackerException(
        statusCode: response.statusCode,
        message: response.body,
      );
    }
  }

  Future<void> declineRanking(int rankingId) async {
    final response = await _httpClient.put(
      Uri.parse('$baseUrl/ranking/v1/rankings/$rankingId/decline'),
      headers: _headers,
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