import 'package:comigo/pages/book_shelf.dart';
import 'package:comigo/pages/scroll_mode.dart';
import 'package:flutter/material.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:provider/provider.dart';
import 'common/theme.dart';
import 'models/setting.dart';
import 'models/remote_server.dart';

Future main() async {
  //https://www.dhiwise.com/post/flutter-dotenv-comprehensive-guide-on-environment-management
  await dotenv.load(fileName: '.env');
  runApp(const ComigoApp());
}

// 这个Widget是应用的根部件。
class ComigoApp extends StatelessWidget {
  const ComigoApp({super.key});

  @override
  Widget build(BuildContext context) {
    // https://github.com/rrousselGit/provider/blob/master/resources/translations/zh-CN/README.md#multiprovider
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => Setting()),
        ChangeNotifierProvider(
            create: (_) =>
                RemoteServer(defaultHost:  dotenv.env['DEFAULT_HOST']!)),
      ],
      child: MaterialApp(
        title: 'Comigo Demo Mobile',
        //隐藏右上角的debug标签
        debugShowCheckedModeBanner: false,
        theme: appTheme,
        //路由表注册
        routes: <String, WidgetBuilder>{
          "/": (context) => const ScrollMode(title: 'Comigo Mobile'),
          "ScrollMode": (context) => const ScrollMode(title: 'ScrollMode'),
          "FlipMode": (context) => const BookShelf(title: 'FlipMode'),
        },
        initialRoute: "/", //设定名为"/"的route为首页
      ),
    );
  }
}
