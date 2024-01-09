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
      //debug条件下，显示右上角的debug标签。我不需要，所以设置为false。
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        // 这是应用程序的主题。
        //
        // 尝试一下：运行应用程序，使用 "flutter run" 命令。会看到应用程序有一个蓝色的工具栏。
        // 然后，在不退出应用程序的情况下，尝试将 colorScheme 中的 seedColor 更改为 Colors.green，
        // 然后触发 "hot reload"（保存更改或在支持Flutter的IDE中按下 "hot reload" 按钮，或者如果你使用命令行启动应用程序，则按下 "r"）。
        //
        // 注意，计数器没有重置为零；应用程序状态在重新加载期间不会丢失。要重置状态，请使用热重启。
        //
        // 这对于代码也适用，不仅仅是值：大多数代码更改只需要进行热重载即可测试。
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.blueAccent),
        useMaterial3: true,
      ),
      home: const ComigoHomePage(title: 'Comigo Home Page'),
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
  int _counter = 0;

  Future<List<Book>>? booksFuture;
  @override
  void initState() {
    super.initState();
    booksFuture = fetchBooks(); // 调用函数并初始化参数
    booksFuture?.then((value) => {
      for (var book in value) {
        print(book.title)
      }
    });
  }

  void _incrementCounter() {
    setState(() {
      // 这个调用setState告诉Flutter框架，这个状态发生了变化，
      // 导致重新运行下面的build方法，以便显示可以反映更新后的值。
      // 如果我们在不调用setState的情况下更改_counter，
      // 那么build方法将不会被再次调用，因此不会有任何变化显示出来。
      _counter = _counter + 2;
    });
  }

  void _minusCounter() {
    setState(() {
      _counter = _counter - 2;
    });
  }

  @override
  Widget build(BuildContext context) {
    // 每次调用setState时，此方法都会重新运行，例如上面的_incrementCounter方法。
    //
    // Flutter框架已经进行了优化，使重新运行build方法变得快速，因此您只需重新构建需要更新的内容，
    // 而不必逐个更改小部件的实例。
    return Scaffold(
      backgroundColor: Colors.yellowAccent,
      appBar: AppBar(
        // 试试这个：尝试在这里将颜色更改为特定的颜色（例如Colors.amber），然后触发热重载，看看AppBar的颜色是否改变，而其他颜色保持不变。
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        // 在这里，我们从由App.build方法创建的MyHomePage对象中获取值，并将其用于设置我们的应用栏标题。
        title: Text(widget.title),
      ),
      body: Center(
        // Center是一个布局小部件。它接受一个子部件并将其放置在父级的中间位置。
        child: Row(
          // Column也是一个布局小部件。它接受一个子部件列表并垂直排列它们。默认情况下，它会水平调整自身大小以适应其子部件，并尽量与其父级一样高。
          // Column有各种属性来控制其自身的大小和子部件的位置。在这里，我们使用mainAxisAlignment来垂直居中子部件；主轴在这里是垂直轴，因为Column是垂直的（交叉轴将是水平的）。
          // 试试这个：调用“调试绘制”（在IDE中选择“切换调试绘制”操作，或在控制台中按“p”），以查看每个小部件的线框。
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            const Text(
              '你已经点击按钮：',
            ),
            Text(
              '$_counter',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            FloatingActionButton.extended(
              onPressed: _minusCounter,
              label: const Text('Minus'),
              icon: const Icon(Icons.article),
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _incrementCounter,
        tooltip: 'Increment',
        child: const Icon(Icons.add),
      ), // 这个尾随逗号使构建方法的自动格式化更加美观。
    );
  }
}
