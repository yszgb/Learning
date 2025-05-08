import 'dart:developer';

import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:learning/constants/app_strings.dart';
import 'package:learning/utils/account_info.dart';
import 'package:learning/utils/data_fetcher.dart';

class AccountPageLogic extends ChangeNotifier {
  bool get isLoggedIn => AccountInfo.shared.isLoggedIn;

  void init() {
    AccountInfo.shared.userBriefNotifier.addListener(_userBriefListener);
  }

  @override
  void dispose() {
    AccountInfo.shared.userBriefNotifier.removeListener(_userBriefListener);
    super.dispose();
  }

  void _userBriefListener() {
    notifyListeners(); // 通知监听器，更新状态
  }

  int? get learnedDays {
    final userBrief = AccountInfo.shared.userBrief;
    if (userBrief == null) {
      return null;
    }
    // 计算学习天数，从注册时间开始算
    final start = userBrief.createdAt;
    return DateTime.now().difference(start).inDays;
  }

  int get learnedMinutes {
    return 0;
  }

  String? get nickName {
    return AccountInfo.shared.userBrief?.name;
  }

  String? get gender {
    return AccountInfo.shared.userBrief?.gender;
  }

  void logout() async {
    try {
      EasyLoading.show(
        indicator: const CircularProgressIndicator(),
        maskType: EasyLoadingMaskType.clear,
        dismissOnTap: false,
      );

      final response = await DataFetcher.shared.userApi.logout();
      if (response.code == 0) {
        EasyLoading.showSuccess(AppStrings.logoutSuccess);
      } else {
        EasyLoading.showError(AppStrings.logoutFailed);
        if (response.message != null) {
          log('logout failed with remote message: ${response.message}');
        }
      }
    } on DioException catch (e) {
      log('Error when logout: ${e.toString()}');
    } finally {
      EasyLoading.dismiss();
      AccountInfo.shared.logout();
    }
  }
}
