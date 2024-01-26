import 'package:comigo/pages/home.dart';
import 'package:flutter/material.dart';

void main() {
  runApp(const ComigoApp());
}

// 这个Widget是应用的根部件。
class ComigoApp extends StatelessWidget {
  const ComigoApp({super.key});

  @override
  Widget build(BuildContext context) {
    // 谷歌推荐的Material（ Android 默认的视觉风格）的组件库
    return MaterialApp(
      title: 'Comigo Demo Mobile',
      //隐藏右上角的debug标签
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        // 应用程序的主题。
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.blueAccent),
        useMaterial3: true,
      ),
      //路由表注册
      routes: <String, WidgetBuilder>{
        "/": (context) => const ComigoHomePage(title: 'Comigo Mobile'), //注册首页路由
        "ScrollMode": (context) => const ComigoHomePage(title: 'Comigo Mobile'),
        "FlipMode": (context) => const ComigoHomePage(title: 'Comigo Mobile'),
      },
      initialRoute: "/", //设定名为"/"的route为首页
    );
  }
}

