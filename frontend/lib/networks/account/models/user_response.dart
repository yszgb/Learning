import 'package:json_annotation/json_annotation.dart';
import 'package:learning/networks/common/base_response.dart';

part 'user_response.g.dart';

@JsonSerializable()
class UserBrief {
  @JsonKey(name: 'id')
  int id;
  @JsonKey(name: 'type')
  String type;
  @JsonKey(name: 'name')
  String name;
  @JsonKey(name: 'gender')
  String? gender;
  @JsonKey(name: 'avatar')
  String? avatar;
  @JsonKey(name: 'created_at')
  DateTime createdAt;
  @JsonKey(name: 'purchasing_courses')
  List<int> purchasingCourses;
  @JsonKey(name: 'purchased_courses')
  List<int> purchasedCourses;

  UserBrief({
    required this.id,
    required this.type,
    required this.name,
    this.gender,
    this.avatar,
    required this.createdAt,
    this.purchasingCourses = const [],
    this.purchasedCourses = const [],
  });

  factory UserBrief.fromJson(Map<String, dynamic> json) =>
      _$UserBriefFromJson(json);
  Map<String, dynamic> toJson() => _$UserBriefToJson(this);
}

@JsonSerializable()
class UserBriefResponse extends BaseResponse {
  @JsonKey(name: 'token')
  String? token;
  @JsonKey(name: 'brief')
  UserBrief? brief;

  UserBriefResponse({
    required super.code,
    this.token,
    this.brief,
    super.message,
  });

  factory UserBriefResponse.fromJson(Map<String, dynamic> json) =>
      _$UserBriefResponseFromJson(json);

  Map<String, dynamic> toJson() => _$UserBriefResponseToJson(this);
}
