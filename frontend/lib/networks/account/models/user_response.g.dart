// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

UserBrief _$UserBriefFromJson(Map<String, dynamic> json) => UserBrief(
  id: (json['id'] as num).toInt(),
  type: json['type'] as String,
  name: json['name'] as String,
  gender: json['gender'] as String?,
  avatar: json['avatar'] as String?,
  createdAt: DateTime.parse(json['created_at'] as String),
  purchasingCourses:
      (json['purchasing_courses'] as List<dynamic>?)
          ?.map((e) => (e as num).toInt())
          .toList() ??
      const [],
  purchasedCourses:
      (json['purchased_courses'] as List<dynamic>?)
          ?.map((e) => (e as num).toInt())
          .toList() ??
      const [],
);

Map<String, dynamic> _$UserBriefToJson(UserBrief instance) => <String, dynamic>{
  'id': instance.id,
  'type': instance.type,
  'name': instance.name,
  'gender': instance.gender,
  'avatar': instance.avatar,
  'created_at': instance.createdAt.toIso8601String(),
  'purchasing_courses': instance.purchasingCourses,
  'purchased_courses': instance.purchasedCourses,
};

UserBriefResponse _$UserBriefResponseFromJson(Map<String, dynamic> json) =>
    UserBriefResponse(
      code: (json['code'] as num).toInt(),
      token: json['token'] as String?,
      brief:
          json['brief'] == null
              ? null
              : UserBrief.fromJson(json['brief'] as Map<String, dynamic>),
      message: json['message'] as String?,
    );

Map<String, dynamic> _$UserBriefResponseToJson(UserBriefResponse instance) =>
    <String, dynamic>{
      'code': instance.code,
      'message': instance.message,
      'token': instance.token,
      'brief': instance.brief,
    };
