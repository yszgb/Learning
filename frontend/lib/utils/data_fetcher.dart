import 'dart:io';

import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:learning/constants/config.dart';
import 'package:learning/networks/account/user_api.dart';
import 'package:learning/utils/app_preferences.dart';
import 'package:pretty_dio_logger/pretty_dio_logger.dart';

class DataFetcher {
  late final Dio dio;
  late final UserApi userApi = RemoteUserApi(dio);

  static final shared = DataFetcher._();
  DataFetcher._();

  Future<void> init() async {
    HttpOverrides.global = _MyHttpOverrides();
    dio = _buildDio();
  }

  void updateAuthorizationToken(String? token) {
    if (token != null) {
      final headers = dio.options.headers;
      headers['authorization'] = 'Bearer $token';
      final options = dio.options.copyWith(headers: headers);
      dio.options = options;
    } else {
      final headers = dio.options.headers;
      headers.remove('authorization');
      final options = dio.options.copyWith(headers: headers);
      dio.options = options;
    }
  }

  Dio _buildDio() {
    final dio = Dio();
    const timeOut = Duration(seconds: 60);
    final headers = <String, String>{
      'Content-Type': 'application/json',
      'accept': 'application/json',
    };

    final token = AppPreferences.shared.userToken;
    if (token != null) {
      headers['authorization'] = 'Bearer $token';
    }

    dio.options = BaseOptions(
      baseUrl: Config.baseUrl,
      connectTimeout: timeOut,
      receiveTimeout: timeOut,
      headers: headers,
      validateStatus: (status) {
        return true;
      },
    );

    if (kDebugMode) {
      dio.interceptors.add(
        PrettyDioLogger(
          requestHeader: true,
          requestBody: true,
          responseHeader: true,
        ),
      );
    }

    return dio;
  }
}

class _MyHttpOverrides extends HttpOverrides {
  @override
  HttpClient createHttpClient(SecurityContext? context) {
    return super.createHttpClient(context)
      ..badCertificateCallback = (X509Certificate cert, String host, int port) {
        return true; // 忽略证书验证
      };
  }
}
