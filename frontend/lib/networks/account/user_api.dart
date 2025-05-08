import 'package:dio/dio.dart';
import 'package:learning/constants/config.dart';
import 'package:learning/networks/account/models/user_response.dart';
import 'package:learning/networks/common/base_response.dart';
import 'package:retrofit/retrofit.dart';

part 'user_api.g.dart';

abstract class UserApi {
  Future<UserBriefResponse> login(
    String account,
    String code,
    CancelToken? cancelToken,
  );
  Future<SimpleResponse> logout();

  // 后端接口
  Future<UserBriefResponse> getUserBrief();
}

@RestApi(baseUrl: Config.baseUrl)
abstract class RemoteUserApi extends UserApi {
  factory RemoteUserApi(Dio dio, {String baseUrl}) = _RemoteUserApi;

  @override
  @POST('/v1/user/login')
  Future<UserBriefResponse> login(
    @Field('account') String account,
    @Field('code') String code,
    @CancelRequest() CancelToken? cancelToken,
  );

  @override
  @POST('/v1/user/logout')
  Future<SimpleResponse> logout();

  @override
  @GET('/v1/user/brief')
  Future<UserBriefResponse> getUserBrief();
}
