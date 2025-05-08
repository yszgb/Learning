import 'package:flutter/material.dart';
import 'package:learning/constants/app_icons.dart';
import 'package:learning/constants/app_strings.dart';
import 'package:learning/pages/common/account_avatar_widget.dart';
import 'package:learning/pages/main_tab/account/account_page_logic.dart';
import 'package:learning/utils/app_theme.dart';
import 'package:provider/provider.dart';

class AccountPage extends StatefulWidget {
  const AccountPage({super.key});

  @override
  State<AccountPage> createState() => _AccountPageState();
}

class _AccountPageState extends State<AccountPage> {
  @override
  void initState() {
    super.initState();

    context.read<AccountPageLogic>().init();
  }

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Padding(
        // 左右边距 20
        padding: const EdgeInsets.symmetric(horizontal: 20),
        // 占位。不会遮挡状态栏
        child: ScrollConfiguration(
          // 不显示滚动条
          behavior: ScrollConfiguration.of(context).copyWith(scrollbars: false),
          child: ListView(
            children: [
              const SizedBox(height: 10),
              // const _NoneLogginHeader(),
              _buildPageHeader(context),
              const SizedBox(height: 30),
              _buildFirstSection(),
              const SizedBox(height: 20),
              _buildSecondSection(),
              _buildExitRow(context),
              const SizedBox(height: 10),
              const _Footer(),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildFirstSection() {
    const divider = Divider(
      height: 0.5,
      indent: 50,
      endIndent: 15,
      color: Colors.grey,
    );

    return Card(
      child: Column(
        children: [
          _ListViewCell(
            icon: Icons.shopping_cart_outlined,
            title: AppStrings.shoppingCart,
            onTap: () {
              print('About Learning');
            },
          ),
          divider,
          _ListViewCell(
            icon: Icons.assignment,
            title: AppStrings.myOrder,
            onTap: () {
              print('About Learning');
            },
          ),
          divider,
          _ListViewCell(
            icon: Icons.favorite,
            title: AppStrings.favoriteCourse,
            onTap: () {
              print('About Learning');
            },
          ),
        ],
      ),
    );
  }

  Widget _buildSecondSection() {
    const divider = Divider(
      height: 0.5,
      indent: 50,
      endIndent: 15,
      color: Colors.grey,
    );

    return Card(
      child: Column(
        children: [
          _ListViewCell(
            icon: Icons.info,
            title: AppStrings.aboutLearning,
            onTap: () {
              print('About Learning');
            },
          ),
          divider,
          _ListViewCell(
            icon: Icons.support_outlined,
            title: AppStrings.welcomePage,
            onTap: () {
              print('About Learning');
            },
          ),
          divider,
          _ListViewCell(
            icon: Icons.star_rate,
            title: AppStrings.ratingMe,
            onTap: () {
              print('About Learning');
            },
          ),
          divider,
          _ListViewCell(
            icon: Icons.edit,
            title: AppStrings.editProfile,
            onTap: () {
              print('About Learning');
            },
          ),
        ],
      ),
    );
  }

  Widget _buildExitRow(BuildContext context) {
    return Selector<AccountPageLogic, bool>(
      selector: (context, logic) => logic.isLoggedIn,
      builder: (context, v, child) {
        if (v) {
          return Column(
            children: [
              const SizedBox(height: 20),
              Card(
                child: _ListViewCell(
                  icon: Icons.exit_to_app,
                  title: AppStrings.logout,
                  onTap: () {
                    context.read<AccountPageLogic>().logout();
                  },
                ),
              ),
            ],
          );
        } else {
          // 如果没有登录，就不显示退出按钮
          return const SizedBox.shrink();
        }
      },
    );
  }

  /// 构建页面头部组件
  Widget _buildPageHeader(BuildContext context) {
    return Selector<AccountPageLogic, bool>(
      selector: (context, logic) => logic.isLoggedIn,
      builder: (context, v, child) {
        if (v) {
          return const _BriefProfileHeader();
        } else {
          return const _NoneLogginHeader();
        }
      },
    );
  }
}

// 列表组件
class _ListViewCell extends StatelessWidget {
  final IconData icon;
  final String title;
  final GestureTapCallback onTap;

  const _ListViewCell({
    required this.icon,
    required this.title,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final style = Theme.of(
      context,
    ).textTheme.titleLarge?.copyWith(fontWeight: FontWeight.w600, fontSize: 16);
    return InkWell(
      onTap: onTap,
      child: ListTile(
        leading: Icon(
          icon,
          color: Theme.of(context).textTheme.titleLarge?.color,
          size: 20,
        ),
        trailing: const Icon(
          Icons.arrow_forward_ios,
          color: Colors.grey,
          size: 16,
        ),
        title: Text(title, style: style),
        onTap: null,
      ),
    );
  }
}

class _NoneLogginHeader extends StatelessWidget {
  const _NoneLogginHeader();

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () {},
      child: Container(
        color: Theme.of(context).scaffoldBackgroundColor,
        child: Row(
          children: [
            SizedBox(
              height: 60,
              width: 60,
              child: CircleAvatar(
                radius: 30,
                backgroundColor: Theme.of(context).colorScheme.surface,
                child: ColorFiltered(
                  colorFilter: ColorFilter.mode(
                    Colors.grey[600]!,
                    BlendMode.srcIn,
                  ),
                  child: Image.asset(AppIcons.greyAvatar),
                ),
              ),
            ),
            const SizedBox(width: 16),
            SizedBox(
              height: 60,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    AppStrings.clickLogin,
                    style: Theme.of(context).textTheme.headlineSmall,
                  ),
                  const SizedBox(height: 5),
                  Text(
                    AppStrings.welcomeLogin,
                    style: Theme.of(
                      context,
                    ).textTheme.bodyMedium?.copyWith(color: Colors.grey),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

// 底部组件
class _Footer extends StatelessWidget {
  const _Footer();

  @override
  Widget build(BuildContext context) {
    final style = Theme.of(
      context,
    ).textTheme.bodySmall?.copyWith(color: Colors.grey);

    return Column(
      children: [
        Text('www.learning.com', style: style),
        const SizedBox(height: 2),
        Text(AppStrings.profileFooter, style: style),
        const SizedBox(height: 10),
      ],
    );
  }
}

class _BriefProfileHeader extends StatelessWidget {
  const _BriefProfileHeader();

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Theme.of(context).scaffoldBackgroundColor,
      child: Row(
        children: [
          _buildUserAvatar(context),
          const SizedBox(width: 16),
          _buildUserTextInfo(context),
        ],
      ),
      // Spacer 是一个占位组件，它会占据尽可能多的空间
    );
  }

  /// 构建用户简介组件
  Widget _buildUserTextInfo(BuildContext context) {
    return SizedBox(
      height: 60,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Selector<AccountPageLogic, String?>(
                selector: (context, logic) => logic.nickName,
                builder: (context, v, child) {
                  return Text(
                    v ?? '',
                    style: Theme.of(context).textTheme.headlineSmall,
                  );
                },
              ),
              const SizedBox(width: 5),
              _buildGenderImage(context),
            ],
          ),
          const SizedBox(height: 5),
          Row(
            children: [
              _buildLearnedDaysText(context),
              const SizedBox(width: 10),
              _buildLearnedMinutesText(context),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildGenderImage(BuildContext context) {
    return Selector<AccountPageLogic, String?>(
      selector: (context, logic) => logic.gender,
      builder: (context, v, child) {
        if (v != null) {
          final icon = v == 'm' ? Icons.male : Icons.female;
          return ClipOval(
            child: Container(
              width: 20,
              height: 20,
              color: Theme.of(
                context,
              ).colorScheme.surface.withValues(alpha: 0.75), // 设置透明度
              child: Center(
                child: Icon(
                  icon,
                  size: 15.5,
                  color: Theme.of(context).colorScheme.primary,
                ),
              ),
            ),
          );
        } else {
          return const SizedBox.shrink();
        }
      },
    );
  }

  /// 显示已学习天数的组件
  Widget _buildLearnedDaysText(BuildContext context) {
    final isDarkMode = AppTheme.shared.isDarkMode;
    final nomalColor = isDarkMode ? Colors.grey[300] : Colors.grey[700];
    final normalStyle = Theme.of(
      context,
    ).textTheme.bodyMedium?.copyWith(color: nomalColor);
    final boldStyle = Theme.of(
      context,
    ).textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600);
    final span1 = TextSpan(text: AppStrings.learnedDays, style: normalStyle);
    final span3 = TextSpan(text: AppStrings.day, style: normalStyle);

    return Selector<AccountPageLogic, int?>(
      selector: (context, logic) => logic.learnedDays,
      builder: (context, v, child) {
        final span2 = TextSpan(text: '${v ?? 0}', style: boldStyle);
        return RichText(text: TextSpan(children: [span1, span2, span3]));
      },
    );
  }

  Widget _buildLearnedMinutesText(BuildContext context) {
    final isDarkMode = AppTheme.shared.isDarkMode;
    final nomalColor = isDarkMode ? Colors.grey[300] : Colors.grey[700];
    final normalStyle = Theme.of(
      context,
    ).textTheme.bodyMedium?.copyWith(color: nomalColor);
    final boldStyle = Theme.of(
      context,
    ).textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600);
    final span1 = TextSpan(text: AppStrings.learnedMinutes, style: normalStyle);

    return Selector<AccountPageLogic, int>(
      selector: (context, logic) => logic.learnedMinutes,
      builder: (context, v, child) {
        TextSpan span2, span3;
        if (v / 60 > 1) {
          span2 = TextSpan(text: '${v ~/ 60}', style: boldStyle);
          span3 = TextSpan(text: AppStrings.hour, style: normalStyle);
        } else {
          span2 = TextSpan(text: '$v', style: boldStyle);
          span3 = TextSpan(text: AppStrings.minute, style: normalStyle);
        }
        return RichText(text: TextSpan(children: [span1, span2, span3]));
      },
    );
  }

  Widget _buildUserAvatar(BuildContext context) {
    return const SizedBox(
      height: 60,
      width: 60,
      child: ClipOval(child: AccountAvatarWidget()), // 裁剪为圆形头像
    );
  }
}
