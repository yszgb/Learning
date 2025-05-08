import 'package:flutter/material.dart';
import 'package:lazy_load_indexed_stack/lazy_load_indexed_stack.dart';
import 'package:learning/constants/app_strings.dart';
import 'package:learning/pages/main_tab/account/account_page.dart';
import 'package:learning/pages/main_tab/account/account_page_logic.dart';
import 'package:learning/pages/main_tab/main_tab_logic.dart';
import 'package:provider/provider.dart';

class MainTab extends StatefulWidget {
  const MainTab({super.key});

  @override
  State<MainTab> createState() => _MainTabState();
}

class _MainTabState extends State<MainTab> {
  late List<Widget> _pages;

  @override
  void initState() {
    // 重启才会调用，重载不会
    super.initState();

    _pages = [
      // 4 个页面
      Container(color: Colors.red),
      Container(color: Colors.green),
      Container(color: Colors.blue),
      ChangeNotifierProvider(
        create: (context) => AccountPageLogic(),
        child: const AccountPage(),
      ),
    ];

    context.read<MainTabLogic>().init();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // 基础页面结构（脚手架）
      body: _buildPageView(context), // 不是 child
      bottomNavigationBar: _buildBottomNavigationBar(context),
    );
  }

  Widget _buildPageView(BuildContext context) {
    return Selector<MainTabLogic, int>(
      selector: (context, logic) => logic.mainTabIndex,
      builder: (context, index, child) {
        return LazyLoadIndexedStack(index: index, children: _pages);
      },
    );
  }

  Widget _buildBottomNavigationBar(BuildContext context) {
    final bottomTabs = <BottomNavigationBarItem>[
      const BottomNavigationBarItem(
        icon: Icon(Icons.home),
        label: AppStrings.home,
      ),
      const BottomNavigationBarItem(
        icon: Icon(Icons.grid_view_rounded),
        label: AppStrings.category,
      ),
      const BottomNavigationBarItem(
        icon: Icon(Icons.menu_book_sharp),
        label: AppStrings.learning,
      ),
      const BottomNavigationBarItem(
        icon: Icon(Icons.person),
        label: AppStrings.profile,
      ),
    ];

    return Selector<MainTabLogic, int>(
      selector: (_, logic) => logic.mainTabIndex,
      builder: (context, index, child) {
        return BottomNavigationBar(
          // 底部导航栏
          // 选中时
          selectedFontSize: 10,
          // selectedItemColor: Theme.of(context).colorScheme.primary,
          selectedItemColor: Colors.black,
          // 未选中时
          unselectedFontSize: 10,
          // unselectedItemColor: Theme.of(context).colorScheme.secondary,
          unselectedItemColor: Colors.grey,
          type: BottomNavigationBarType.fixed, // 默认为 shifting。fixed 固定在底部
          items: bottomTabs,
          currentIndex: index,
          onTap: (index) {
            final logic = context.read<MainTabLogic>();
            logic.setMainTabIndex(index);
            // setState(() {
            //   _currentIndex = index;
            // });
          },
        );
      },
    );
  }
}
