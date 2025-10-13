import 'package:json_annotation/json_annotation.dart';

part 'device_token.g.dart';

@JsonSerializable()
class DeviceToken {
  final int id;
  final String token;
  final String platform;
  @JsonKey(name: 'created_at')
  final DateTime createdAt;
  @JsonKey(name: 'updated_at')
  final DateTime updatedAt;

  const DeviceToken({
    required this.id,
    required this.token,
    required this.platform,
    required this.createdAt,
    required this.updatedAt,
  });

  factory DeviceToken.fromJson(Map<String, dynamic> json) =>
      _$DeviceTokenFromJson(json);
  Map<String, dynamic> toJson() => _$DeviceTokenToJson(this);
}
