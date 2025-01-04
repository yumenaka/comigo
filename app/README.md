## 开发笔记

flutter_launcher_icons自动生成图标 ，可以在项目文件夹下运行此命令：
```bash
dart run flutter_launcher_icons && rm -r android/app/src/main/res/mipmap-anydpi-v26
```

重命名项目（https://pub.dev/packages/rename）：
```bash
flutter pub global activate rename
flutter pub global run rename setAppName --targets ios,android,macos,windows,linux,web --value "Comigo"
```
