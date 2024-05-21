import 'package:flutter/cupertino.dart';
import 'package:shared_preferences/shared_preferences.dart';

enum ReadMode { scrollMode, flipMode }

class Setting extends ChangeNotifier {
  // 服务器地址,当前为远程地址
  String serverHost = "";
  // 阅读模式，默认卷轴模式
  ReadMode readMode = ReadMode.scrollMode;
  Setting() {
    init();
  }

  // 初始化设置，利用shared_preferences插件，从本地读取设置
  Future<void> init() async {
    final prefs = await SharedPreferences.getInstance();
    serverHost = prefs.getString('comigo_host') ?? "http://127.0.0.1:1234";
    final mode = prefs.getString('read_mode');
    if (mode == null) {
      readMode = ReadMode.scrollMode;
    } else {
      readMode = ReadMode.values.firstWhere((e) => e.toString() == mode);
    }
    //模型发生改变并且需要更新 UI 的时候调用该方法
    notifyListeners();
  }

  // 设置服务器地址
  Future<void> setHost(String host) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('comigo_host', host);
    serverHost = host;
    //模型发生改变并且需要更新 UI 的时候调用该方法
    notifyListeners();
  }
  // 设置阅读模式
  Future<void> setReadMode(ReadMode mode) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('read_mode', mode.toString());
    readMode = mode;
    //模型发生改变并且需要更新 UI 的时候调用该方法
    notifyListeners();
  }

}