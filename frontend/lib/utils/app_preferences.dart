import 'package:shared_preferences/shared_preferences.dart';

class AppPreferences {
  AppPreferences._();
  static final shared = AppPreferences._();

  late SharedPreferences _preferences;

  static const _kFirstOpened = 'first_opened';
  static const _kUserToken = 'user_token';

  // bool get isFirstRun => _preferences.getBool(_kFirstOpened) != true;
  bool get isFirstRun => !(_preferences.getBool(_kFirstOpened) ?? false);
  void setFirstOpened() => _preferences.setBool(_kFirstOpened, true);

  String? get userToken => _preferences.getString(_kUserToken);
  // 设置用户 token
  set userToken(String? token) {
    if (token == null) {
      _preferences.remove(_kUserToken);
    } else {
      // 持久化保存 token
      _preferences.setString(_kUserToken, token);
    }
  }

  Future<void> init() async {
    _preferences = await SharedPreferences.getInstance();
  }
}
