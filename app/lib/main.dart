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
        title: 'Comigo Demo',
        //隐藏右上角的debug标签
        debugShowCheckedModeBanner: false,
        theme: appTheme,
        // 下面使用了 const，这意味着它是一个编译时常量，无法接受运行时才传进来的参数。
        // 也就是说，这里写死的 title: 'ScrollMode', bookID: '' 跟想要的动态参数并不匹配。
        // routes: <String, WidgetBuilder>{
        //   "/": (context) => const BookShelf(title: 'Comigo Demo'),
        //   "ScrollMode": (context) => const ScrollMode(title: 'ScrollMode', bookID: '',),
        //   "FlipMode": (context) => const BookShelf(title: 'FlipMode'),
        // },
        ///路由表注册
        onGenerateRoute: (settings) {
          // 这里可以判断当前的 settings.name
          switch (settings.name) {
            case '/':
              return MaterialPageRoute(
                builder: (_) => const BookShelf(title: 'Comigo Demo'),
              );
            case 'ScrollMode':
              final args = settings.arguments as Map<String, dynamic>;
              return MaterialPageRoute(
                builder: (_) => ScrollMode(
                  title: args['title'],
                  bookID: args['bookID'],
                ),
              );
            case 'FlipMode':
              return MaterialPageRoute(
                builder: (_) => const BookShelf(title: 'FlipMode'),
              );
            default:
            // 如果路由表里没有匹配，返回一个默认处理
              return MaterialPageRoute(
                builder: (_) => const Scaffold(
                  body: Center(child: Text('Unknown Route')),
                ),
              );
          }
        },
        initialRoute: "/", //设定名为"/"的route为首页
      ),
    );
  }
}
