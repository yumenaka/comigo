## [v1.1.4 阅读历史、远程连接(Tailscale)](https://github.com/yumenaka/comigo/releases/tag/v1.1.4)

**新功能：**
1. 自动保存与恢复阅读进度，让阅读体验更加连贯。
2. 内置 Tailscale 远程连接功能，轻松实现跨设备访问。
3. 支持多个书库分批展示，加载更高效。

**优化：**
1. 同一浏览器的不同标签页之间也可同步翻页操作。
2. 自动忽略以「.」开头的隐藏文件。
3. 命令行模式下，打印帮助或版本信息后直接退出。
4. 修复网络较慢时书架与设置页面出现闪烁的问题。

**New Features:**
1. Automatically save and restore reading progress for a seamless experience.
2. Built-in Tailscale remote connection for easy cross-device access.
3. Support for displaying multiple libraries in batches for improved loading efficiency.

**Improvements:**
1. Page flipping is now synchronized across different tabs in the same browser.
2. Hidden files (starting with “.”) are now automatically ignored.
3. In command-line mode, the program now exits immediately after printing help or version info.
4. Fixed flickering issues on the bookshelf and settings pages when the network is slow.

## [v1.0.3 全面优化与重构](https://github.com/yumenaka/comigo/releases/tag/v1.0.3)

1.  **翻页阅读模式** 已重写，新增 **滑动翻页** 与 **自动对齐** 功能。
2. **设置面板** 全面优化，修改选项后立即生效。
3. **卷轴模式** 新增键盘快捷键；可选择 一次性全部加载 或 分页加载（32 页）。
4. **实验功能**：静态模式——在卷轴一次性加载时，可将页面另存为 单一 HTML 文件。
5. **技术栈** 由 Vue 迁移至 HTMX + Alpine.js，并修复大量（主要是书籍扫描相关）错误。                         

   

1. **Page-flip mode** has been rewritten with auto-alignment and **Swipe page-turning**.
2. **The Settings panel** has been overhauled; changes now take effect instantly.
3. **Scroll mode** gains keyboard shortcuts and offers two loading strategies: load-all-at-once or paged loading (32 pages per batch).
4. **Experimental:** Static Mode – when using load-all-at-once scroll mode, the entire page can be saved as a single HTML file.
5. **Tech stack** migration from Vue to HTMX + Alpine.js, fixing numerous scan-related bugs.



## [v0.9.12 监控文件变更](https://github.com/yumenaka/comigo/releases/tag/v0.9.12)

1、删除文件的时候同步删除对应书籍。
2、重做上传功能，上传到第一个可用书库。
3、修复网页切边功能，页码显示等BUG。

1、Synchronize the deletion of corresponding books when deleting files.
2、Redesign the upload function: files are directly uploaded to the first library, fixing related bugs.
3、Fix the webpage cropping feature and the page number display bug.




## [v0.9.7 快速切换下一话，保存服务器设置到文件](https://github.com/yumenaka/comigo/releases/tag/v0.9.7)

1、网页端：【快速切换】功能。直接在阅读界面顶部或底部，切换到下一回。
2、网页端：完善【服务器设置】功能。配置文件可以保存到服务器。
3、网页端：优化界面，补充图标。顶部按钮的鼠标悬停提示，补充英语，日语翻译。
4、修复BUG、 画了一个新图标。                       

## [v0.9.0 优化Epub，优化卷轴模式](https://github.com/yumenaka/comigo/releases/tag/v0.9.0)

1、默认读取Epub的opf设置，对图片重新排序。
2、优化卷轴模式，支持自定义页边距。

1、Read Epub's opf settings by default and reorder the images.
2、Scroll mode supports custom margins.



## [v0.8.9 重要更新](https://github.com/yumenaka/comigo/releases/tag/v0.8.9)

1、跨设备同步，用手机当遥控器。
2、文件上传
3、Windows 系统右键菜单
4、卷轴模式优化
5、查看服务器状态。
6、快捷图标“显示二维码”与“切换全屏”。
7、简化书名显示
8、修复一些Linux平台的BUG

## [v0.8.2](https://github.com/yumenaka/comigo/releases/tag/v0.8.2)

web: 设置界面加切换全屏按钮。
web: A  full screen button has been added to the settings screen.

## [v0.8.1](https://github.com/yumenaka/comigo/releases/tag/v0.8.1)

1、重新排序书籍或书架功能，可以按照文件名、修改时间、文件大小（正向与反向）排序。
2、sqlite数据库保存扫描数据，避免重复扫描。需要通过配置文件启用。示例配置文件，可以在web端下载。
3、支持pdf与MP4等。虽然只是简单地用浏览器打开。

1、Reorder function, can be sorted by file name, modification time, file size (forward and reverse).
2、Sqlite database to save scanned data. Need to be enabled through the configuration file. Example configuration file can be downloaded from the web.
3、Support pdf and MP4. although simply open with a browser.

## [v0.7.3](https://github.com/yumenaka/comigo/releases/tag/v0.7.3)

1、webui：添加qrcode、修改设置区域为正方形。
2、支持GIF。修复开启web缓存时，切边与黑白化出错的BUG。                 

1、webui: add qrcode. set area to square.
2、Support GIF. fix the bug when web cache is on.

## [v0.7.1](https://github.com/yumenaka/comigo/releases/tag/v0.7.1)

1、优化逻辑，省略解压过程。
2、新增书架，在线阅读多本书。
3、新功能：自动切边、黑白模式、设定背景色等。
4、增加编译版本：Windows Arm64，修复部分BUG。

1、Optimize the reading logic, omit the decompression process.
2、New bookshelf, you can read multiple books online.
3、New features: auto cut edge, black and white mode, set background color, etc.
4、Add compiled version: Windows Arm64, fix some bugs.

## [v0.6.0](https://github.com/yumenaka/comigo/releases/tag/v0.6.0)

Vue3.0！



## [v0.4.5](https://github.com/yumenaka/comigo/releases/tag/v0.4.5)
一、新女仆图标。
二、支持 ,cbr .cbz 文件。
三、可指定zip文件编码、临时文件夹路径。

1. New maid icon.
2. cbr .cbz files are supported,
3. zip files can be specified with file encoding，Temporary decompression folder for images.


## [v0.4.1](https://github.com/yumenaka/comigo/releases/tag/v0.4.1)

一、国际化、支持汉语、日语与英语。
二、web模板：卷轴模式、双页模式、单页模式、速写模式。
三、限制线程占用数，并修复部分BUG。   

1. Internationalization, support Chinese, Japanese and English.
2. Web template: scroll mode, double page mode, single page mode, sketch mode.
3. Limit the number of threads, and fix some bugs.



## [v0.3.0](https://github.com/yumenaka/comigo/releases/tag/v0.3.0)

一、优化cmd显示，修复bug。
二、初步实现国际化。
三、go升级到1.16beta。   

1. Optimize cmd display and fix bugs.
2. Initially implement multi-language.
3. Update golang to 1.16beta.

