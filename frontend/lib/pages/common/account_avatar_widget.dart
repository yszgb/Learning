import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:learning/constants/app_images.dart';
import 'package:learning/constants/config.dart';
import 'package:learning/utils/account_info.dart';
import 'package:learning/utils/app_preferences.dart';
import 'package:learning/utils/app_theme.dart';

class AccountAvatarWidget extends StatefulWidget {
  const AccountAvatarWidget({super.key});

  @override
  State<AccountAvatarWidget> createState() => _AccountAvatarWidgetState();
}

class _AccountAvatarWidgetState extends State<AccountAvatarWidget> {
  int _index = 0;

  @override
  void initState() {
    super.initState();
    // 监听账号信息变化
    AccountInfo.shared.avatarChangeIndexNotifier.addListener(_updateIndex);
  }

  @override
  void dispose() {
    super.dispose();
    // 移除监听
    AccountInfo.shared.avatarChangeIndexNotifier.removeListener(_updateIndex);
  }

  void _updateIndex() {
    setState(() {
      // setState 会导致 build 方法重新执行，所以这里的 _index 会被更新
      _index = AccountInfo.shared.avatarChangeIndexNotifier.value;
    });
  }

  @override
  Widget build(BuildContext context) {
    return _InnerAccountUserAvatar(key: ValueKey(_index));
  }
}

class _InnerAccountUserAvatar extends StatelessWidget {
  const _InnerAccountUserAvatar({super.key}); // super 是必须的，因为这里调用了父类的构造函数
  // 这里把子类的 key 传递给了父类的构造函数，这样父类的 key 就会被更新，从而导致 build 方法重新执行

  @override
  Widget build(BuildContext context) {
    final token = AppPreferences.shared.userToken;
    final avatar = AccountInfo.shared.userBrief?.avatar;
    if (token == null || avatar == null) {
      return Image.asset(AppImages.avatarPlaceholder);
    }
    return CachedNetworkImage(
      imageUrl: '${Config.baseUrl}/v1/user/avatar',
      httpHeaders: {'authorization': 'Bearer $token'},
      placeholder:
          (context, url) => ColorFiltered(
            colorFilter: ColorFilter.mode(
              AppTheme.shared.currentTheme.colorScheme.primary,
              BlendMode.modulate,
            ),
            child: Image.asset(AppImages.avatarPlaceholder),
          ),
      fit: BoxFit.cover,
      // 加载失败时显示默认头像
      errorWidget:
          (context, url, error) => Image.asset(AppImages.avatarPlaceholder),
      errorListener: (e) {},
    );
  }
}
