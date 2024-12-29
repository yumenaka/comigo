import 'package:flutter/cupertino.dart';
import 'package:shared_preferences/shared_preferences.dart';

class RemoteServer extends ChangeNotifier {
  late String remoteHost;

  RemoteServer({required String defaultHost}) {
    remoteHost= defaultHost;
    loadHost(defaultHost: defaultHost);
  }

  Future<void> loadHost({required String defaultHost}) async {
    final prefs = await SharedPreferences.getInstance();
    var h = prefs.getString('comigo_host');
    if (h != null) {
      remoteHost = h;
      //模型发生改变并且需要更新 UI 的时候调用该方法
      notifyListeners();
    }
  }

  Future<void> saveHost(String host) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('comigo_host', host);
    remoteHost = host;
    //模型发生改变并且需要更新 UI 的时候调用该方法
    notifyListeners();
  }
}
