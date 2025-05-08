import 'package:flutter/material.dart';
import 'package:learning/networks/account/models/user_response.dart';
import 'package:learning/utils/app_preferences.dart';
import 'package:learning/utils/data_fetcher.dart';

import 'dart:developer';

class AccountInfo {
  final userBriefNotifier = ValueNotifier<UserBrief?>(null);
  final avatarChangeIndexNotifier = ValueNotifier<int>(0);

  bool get isLoggedIn => userBriefNotifier.value != null;
  UserBrief? get userBrief => userBriefNotifier.value;

  AccountInfo._();
  static final shared = AccountInfo._();

  Future<void> init() async {
    _loadUserBrief();
  }

  Future<void> login(UserBrief brief, String token) async {
    AppPreferences.shared.userToken = token; // 前端持久化保存 token
    DataFetcher.shared.updateAuthorizationToken(token); // 更新 dio 的 token

    _setUserBrief(brief);
  }

  void logout() {
    AppPreferences.shared.userToken = null;
    DataFetcher.shared.updateAuthorizationToken(null);
    userBriefNotifier.value = null;
    avatarChangeIndexNotifier.value = 0;
  }

  void _setUserBrief(UserBrief? brief) {
    userBriefNotifier.value = brief;
    avatarChangeIndexNotifier.value++;
  }

  // 测试函数
  void _loadUserBrief() async {
    if (AppPreferences.shared.userToken == null) {
      final response = await DataFetcher.shared.userApi.login(
        'admin@qq.com',
        '123456789',
        null,
      );
      if (response.code == 0) {
        login(response.brief!, response.token!);
      }
    } else {
      try {
        final response = await DataFetcher.shared.userApi.getUserBrief();
        if (response.code == 0) {
          _setUserBrief(response.brief);
        }
      } catch (e) {
        log('Error when facing user breif: ${e.toString()}');
      }
    }
  }
}
