import 'package:json_annotation/json_annotation.dart';

part 'paginated_result.g.dart';

@JsonSerializable(genericArgumentFactories: true)
class PaginatedResult<T> {
  final List<T> items;
  @JsonKey(name: 'total_count')
  final int totalCount;
  final int page;
  @JsonKey(name: 'per_page')
  final int perPage;

  const PaginatedResult({
    required this.items,
    required this.totalCount,
    required this.page,
    required this.perPage,
  });

  factory PaginatedResult.fromJson(
    Map<String, dynamic> json,
    T Function(Object? json) fromJsonT,
  ) =>
      _$PaginatedResultFromJson(json, fromJsonT);

  Map<String, dynamic> toJson(Object Function(T value) toJsonT) =>
      _$PaginatedResultToJson(this, toJsonT);
}
