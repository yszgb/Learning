import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:learning/constants/app_strings.dart';
import 'package:learning/pages/main_tab/main_tab.dart';
import 'package:learning/pages/main_tab/main_tab_logic.dart';
import 'package:learning/pages/signin/welcome/welcome_page.dart';
import 'package:learning/utils/account_info.dart';
import 'package:learning/utils/app_preferences.dart';
import 'package:learning/utils/app_theme.dart';
import 'package:learning/utils/data_fetcher.dart';
import 'package:provider/provider.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  await AppPreferences.shared.init();
  await AppTheme.shared.init();
  await DataFetcher.shared.init();
  await AccountInfo.shared.init();

  runApp(const MyApp());
}

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  late final MaterialApp _mainApp;

  bool _showMainWidget = false;

  @override
  void initState() {
    super.initState();
    
    _mainApp = MaterialApp(
      title: AppStrings.appName,
      theme: AppTheme.shared.currentTheme,
      // debugShowCheckedModeBanner: false, // 隐藏调试横幅
      home: ChangeNotifierProvider(
        create: (context) => MainTabLogic(mainTabIndex: 3), // 初始选中第 4 个页面
        child: const MainTab(),
      ),
      builder:EasyLoading.init(),
    );
  }

  // root of application
  @override
  Widget build(BuildContext context) {
    if (AppPreferences.shared.isFirstRun) {
      // 第一次运行
      if (_showMainWidget) {
        // 欢迎页显示完成后，显示主页面
        return _mainApp;
      } else {
        // 显示欢迎页
        return MaterialApp(
          title: AppStrings.appName,
          theme: AppTheme.shared.currentTheme,
          // debugShowCheckedModeBanner: false,
          home: WelcomePage(
            skipAction: () {
              AppPreferences.shared.setFirstOpened();
              setState(() {
                _showMainWidget = true;
              });
            },
          ),
        );
      }
    } else {
      return _mainApp;
    }
  }
}
