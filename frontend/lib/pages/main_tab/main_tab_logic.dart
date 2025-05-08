import 'package:flutter/material.dart';

class MainTabLogic extends ChangeNotifier {
  late int _mainTabIndex;

  int get mainTabIndex => _mainTabIndex;

  MainTabLogic({int mainTabIndex = 0}) {
    _mainTabIndex = mainTabIndex;
  }

  void init(){}

  void setMainTabIndex(int index) {
    _mainTabIndex = index;
    notifyListeners(); // 通知监听器
  }
}
