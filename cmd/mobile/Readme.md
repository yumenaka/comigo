## Android
```
gomobile bind -ldflags="-w -s" -o ../../app/android/app/libs/server.aar -target=android -androidapi 21 -javapkg="net.yumenaka.comigo" github.com/yumenaka/comigo/cmd/mobile
```

## IOS
```
gomobile bind -ldflags="-w -s" -o ../app/ios/Frameworks/server.xcframework -target=ios github.com/yumenaka/comigo/cmd/mobile
```

## MacOS
```
go build -ldflags="-w -s" -buildmode=c-shared -o _temp/output/libserver.dylib github.com/yumenaka/comigo/cmd/desktop
cp _temp/output/libserver.h ../app/include/
cp _temp/output/libserver.dylib ../app/macos/Frameworks/
```

## Linux
```
go build -ldflags="-w -s" -buildmode=c-shared -o libserver.so github.com/yumenaka/comigo/cmd/desktop
```