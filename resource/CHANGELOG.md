# changelog

2024.02.24 v0.9.8:
1. 完善老代码，重新支持 PDF 格式。
2. 卷轴模式：现有的“无限滚动”以外，新增“分页加载”功能。


## 2023.12.29 v0.9.7：
1. 网页端：【章节跳转】。直接在阅读界面顶部或底部，切换下一回。
2. 网页端：完善【Comigo服务器设置】功能，配置文件保存到服务器。
3. 网页端：优化界面，补充图标。顶部按钮的鼠标悬停提示，补充英语，日语翻译。
4. 修复BUG：某些情况下的ePub排序错误。返回上一页错误。
5. 修复BUG：命令行comi --config不生效问题。
6. 修复BUG：设置页面添加路径，扫描时崩溃问题。
7. 画了一个新图标。

## 2023.09.16 v0.9.5：
1、可以在网页端添加书库，修改服务器配置。  
2、新书架显示模式（卡片模式，列表模式，文字模式）。

## 2022.12.11 v0.9.0：
卷轴模式自定义页间距。  
epub文件、按opf文件排序。

## 2022.06.16 v0.8.2：
1、网页设置界面，加了一个切换全屏按钮。  
2、重新支持32位 Windows。

## 2022.06.15 v0.8.1：
1、重新排序书籍或书架功能，可以按照文件名、修改时间、文件大小（正向与反向）排序。  
2、sqlite数据库保存扫描数据，避免重复扫描。需要通过配置文件启用。示例配置文件，可以在web端下载。  
3、支持pdf与MP4等。虽然只是简单地用浏览器打开。


## 2022.07.24 v0.8.9：
1、跨设备同步翻页，手机当遥控器。
不同浏览器与设备，看同一本书的时候，可以通过服务器中转、实现同步翻页。 目前卷轴模式只能当遥控端，不能当被控端。

2、文件上传  
网页上传漫画压缩包，扫描并观看。

3、生成Windows reg文件，注册系统右键菜单。    
需要从书架的设置页面下载reg文件，然后双击导入。
comi.exe文件路径改变，请再次重新生成、下载并导入。

4、卷轴模式优化  
加上了缺少的保存、载入本地阅读记录功能。  
减少单次载入量，下拉到20页以后，才加载接下来的20页。

5、右上角新增两个快捷图标，分别是“显示当前页二维码”与“切换全屏”。

6、简化书名，隐藏类似 [作者名] 【出版社】这样的字段。因为容易误伤，本功能可关闭。
















