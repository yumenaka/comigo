import 'package:comigo/pages/book_shelf.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'models/setting.dart';
import 'models/remote_server.dart';

void main() {
  runApp(const ComigoApp());
}

// 这个Widget是应用的根部件。
class ComigoApp extends StatelessWidget {
  const ComigoApp({super.key});

  @override
  Widget build(BuildContext context) {
    // https://flutter.cn/docs/development/data-and-backend/state-mgmt/simple#changenotifierprovider
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => Setting()),
        ChangeNotifierProvider(create: (_) => RemoteServer()),
      ],
      child: MaterialApp(
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
          "/": (context) => const BookShelf(title: 'Comigo Mobile'),
          "ScrollMode": (context) =>
              const BookShelf(title: 'Comigo Mobile'),
          "FlipMode": (context) => const BookShelf(title: 'Comigo Mobile'),
        },
        initialRoute: "/", //设定名为"/"的route为首页
      ),
    );
  }
}
