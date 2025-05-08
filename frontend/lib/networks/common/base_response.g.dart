// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'base_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

BaseResponse _$BaseResponseFromJson(Map<String, dynamic> json) => BaseResponse(
  code: (json['code'] as num).toInt(),
  message: json['message'] as String?,
);

Map<String, dynamic> _$BaseResponseToJson(BaseResponse instance) =>
    <String, dynamic>{'code': instance.code, 'message': instance.message};

SimpleResponse _$SimpleResponseFromJson(Map<String, dynamic> json) =>
    SimpleResponse(
      code: (json['code'] as num).toInt(),
      message: json['message'] as String?,
    );

Map<String, dynamic> _$SimpleResponseToJson(SimpleResponse instance) =>
    <String, dynamic>{'code': instance.code, 'message': instance.message};
