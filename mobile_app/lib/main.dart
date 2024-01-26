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
      //ルーティング表を登録する
      routes: <String, WidgetBuilder>{
        "webview": (context) => const ComigoHomePage(title: 'Comigo Mobile'),
        "/": (context) => const ComigoHomePage(title: 'Comigo Mobile'), //ホームページのルートを登録する
      },
      initialRoute: "/", //"/"という名前のルートをアプリのホーム（スタートページ）とする
      home: const ComigoHomePage(title: 'Comigo Mobile'),
    );
  }
}

