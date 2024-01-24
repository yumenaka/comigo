import 'package:flutter/material.dart';
import 'models/book.dart';

void main() {
  runApp(const ComigoApp());
}

class ComigoApp extends StatelessWidget {
  const ComigoApp({super.key});

  // 这个小部件是应用的根部件。
  @override
  Widget build(BuildContext context) {
    // 谷歌推荐的Material（ Android 默认的视觉风格）的组件库
    return MaterialApp(
      title: 'Comigo Demo v1.0',
      //隐藏右上角的debug标签
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        // 这是应用程序的主题。
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.blueAccent),
        useMaterial3: true,
      ),
      home: const ComigoHomePage(title: 'Comigo v1.0'),
    );
  }
}

class ComigoHomePage extends StatefulWidget {
  const ComigoHomePage({super.key, required this.title});
  // 这个小部件是应用的主页。它是有状态的，意味着它有一个包含影响其外观的字段的 State 对象（在下面定义）。
  // 这个类是状态的配置。它保存了父级（在这种情况下是 App 小部件）提供的值（在这种情况下是标题），并由状态的 build 方法使用。
  // Widget 子类中的字段始终标记为 "final"。
  final String title;

  @override
  State<ComigoHomePage> createState() => _ComigoHomePageState();
}

class _ComigoHomePageState extends State<ComigoHomePage> {


  Future<List<Book>>? Books;
  @override
  void initState() {
    super.initState();
    Books = fetchBooks(); // 调用函数并初始化参数
    Books?.then((value) => {
      for (var book in value) {
        debugPrint(book.title)
      }
    });
  }

  Future<List<Book>> initBooks() async {
    Future<List<Book>>? books = fetchBooks(); // 调用函数并初始化参数
    return books.then((value) => value);
  }

  @override
  Widget build(BuildContext context) {
    //  画面の高さを取得する
    final mediaQueryData = MediaQuery.of(context);
    // トップナビゲーションバーの高さ
    final headerHeight = mediaQueryData.size.height * 0.07;
    // ボトムナビゲーションバーの高さ
    // final bottomHeight = mediaQueryData.size.height * 0.1;
    // 异步UI更新（FutureBuilder、StreamBuilder）
    // https://book.flutterchina.club/chapter7/futurebuilder_and_streambuilder.html
    Widget booksWidget = FutureBuilder<List<Book>>(
      future: fetchBooks(),
      initialData: [Book(title: 'Loading...', id: '0', type: '0')],
      // snapshot会包含当前异步任务的状态信息及结果信息
      builder: (context, snapshot) {
        if (snapshot.hasData) {
          debugPrint(snapshot.data!.length.toString());
          return ListView.builder(
            itemCount: snapshot.data!.length,
            itemBuilder: (context, index) {
              return ListTile(
                title: Text(snapshot.data![index].title),
                subtitle: Text(snapshot.data![index].id),
              );
            },
          );
        } else if (snapshot.hasError) {
          return Text('${snapshot.error}');
        }
        return const CircularProgressIndicator();
      },
    );

    // 每次调用setState时，此方法都会重新运行，例如上面的_incrementCounter方法。
    // Flutter框架已经进行了优化，使重新运行build方法变得快速，因此您只需重新构建需要更新的内容，
    // 而不必逐个更改小部件的实例。
    return Scaffold(
      backgroundColor: Colors.yellow[200],
      appBar: AppBar(
        title: Text(widget.title,
            style: Theme.of(context).textTheme.headlineMedium),
        backgroundColor: Colors.lightBlue, ////appbarの背景色を設定する
        toolbarHeight: headerHeight, //appbarの高さを設定する
        centerTitle: true, //タイトルを中央に配置
      ),
      body: booksWidget,
    );
  }
}
