import 'package:json_annotation/json_annotation.dart';

part 'base_response.g.dart';

@JsonSerializable()
class BaseResponse {
  @JsonKey(name: 'code')
  int code;
  @JsonKey(name: 'message')
  String? message;

  BaseResponse({required this.code, this.message});
}

@JsonSerializable()
class SimpleResponse extends BaseResponse {
  SimpleResponse({required super.code, super.message});

  factory SimpleResponse.fromJson(Map<String, dynamic> json) =>
      _$SimpleResponseFromJson(json);
  Map<String, dynamic> toJson() => _$SimpleResponseToJson(this);
}
