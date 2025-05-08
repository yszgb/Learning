import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:json_theme/json_theme.dart';

class AppTheme {
  static final shared = AppTheme._();
  AppTheme._(); // 私有构造函数。外部只能通过 shared 访问

  late final ThemeData currentTheme; // 主题
  bool get isDarkMode => currentTheme.brightness == Brightness.dark;

  Future<void> init() async {
    final themeString = await rootBundle.loadString(
      'assets/jsons/advanced_theme.json',
    );
    final themeJson = jsonDecode(themeString);
    currentTheme =
        ThemeDecoder.decodeThemeData(themeJson)!; // "!" 表示 themeJson 不为 null
  }
}
