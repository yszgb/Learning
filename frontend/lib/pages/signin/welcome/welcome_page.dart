import 'package:flutter/material.dart';
import 'package:learning/constants/app_images.dart';
import 'package:learning/constants/app_strings.dart';
import 'package:smooth_page_indicator/smooth_page_indicator.dart';

class WelcomePage extends StatelessWidget {
  final _pageController = PageController(initialPage: 0);

  final Function skipAction;
  WelcomePage({super.key, required this.skipAction});

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Theme.of(context).scaffoldBackgroundColor,
      child: SafeArea(
        child: Column(
          children: [
            Expanded(
              flex: 3, // 占据 3/4 屏幕高度
              child: _buildPageView(context),
            ),
            Expanded(
              flex: 1, // 占据 1/4 屏幕高度
              child: Column(
                children: [
                  const SizedBox(height: 20), // 与文字底部距离 20
                  _buildIndicator(context),
                  const Spacer(), // 占据剩余空间。指示器与跳过按钮隔一段距离
                  _buildSkipButton(context),
                  const SizedBox(height: 40),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  /// 构建多个页面排布方式
  Widget _buildPageView(BuildContext context) {
    return PageView(
      controller: _pageController,
      scrollDirection: Axis.horizontal,
      reverse: false, // left -> right
      pageSnapping: true, // 滚动时对齐
      // physics: const ClampingScrollPhysics(), // 防止弹过界
      children: [
        _buildSinglePage(
          context: context,
          image: AppImages.welcome1,
          title: AppStrings.welcomeTitle1,
          description: AppStrings.welcomeDesc1,
        ),
        _buildSinglePage(
          context: context,
          image: AppImages.welcome2,
          title: AppStrings.welcomeTitle2,
          description: AppStrings.welcomeDesc2,
        ),
        _buildSinglePage(
          context: context,
          image: AppImages.welcome3,
          title: AppStrings.welcomeTitle3,
          description: AppStrings.welcomeDesc3,
        ),
      ],
    );
  }

  /// 构建单页布局
  Widget _buildSinglePage({
    required BuildContext context,
    required String image,
    required String title,
    required String description,
  }) {
    return Column(
      children: [
        Image.asset(image, fit: BoxFit.fitHeight),
        Container(
          margin: const EdgeInsets.only(top: 15), // 与边缘距离 15
          child: Text(
            title,
            textAlign: TextAlign.center,
            style: Theme.of(context).textTheme.headlineSmall,
          ),
        ),
        Container(
          margin: const EdgeInsets.only(top: 15),
          padding: const EdgeInsets.symmetric(horizontal: 30), // 对称水平距离 30
          child: Text(
            description,
            textAlign: TextAlign.center,
            style: Theme.of(context).textTheme.bodyLarge,
          ),
        ),
      ],
    );
  }

  // 构建底部按钮
  Widget _buildIndicator(BuildContext context) {
    return SmoothPageIndicator(
      controller: _pageController,
      count: 3,
      effect: WormEffect(activeDotColor: Theme.of(context).colorScheme.primary),
    );
  }

  // 构建跳过按钮
  Widget _buildSkipButton(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.end, // 排列在末尾，即右对齐
      children: [
        MaterialButton(
          onPressed: () async {
            skipAction();
          },
          child: Text(
            AppStrings.skip,
            style: Theme.of(context).textTheme.headlineSmall,
          ),
        ),
        const SizedBox(width: 15), // 与右侧距离 15，避免太靠右
      ],
    );
  }
}
