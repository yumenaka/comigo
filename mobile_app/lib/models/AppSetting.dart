
import 'package:flutter/cupertino.dart';
import 'package:shared_preferences/shared_preferences.dart';

class AppSetting extends ChangeNotifier {
  String comigoHost = "";
  AppSetting() {
    initHost();
  }

  Future<void> initHost() async {
    final prefs = await SharedPreferences.getInstance();
    comigoHost = prefs.getString('comigo_host') ?? "http://127.0.0.1:1234";
    //模型发生改变并且需要更新 UI 的时候调用该方法
    notifyListeners();
  }

  Future<void> setHost(String host) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('comigo_host', host);
    comigoHost = host;
    //模型发生改变并且需要更新 UI 的时候调用该方法
    notifyListeners();
  }
}