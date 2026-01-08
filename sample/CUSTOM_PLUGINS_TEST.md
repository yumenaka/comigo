# 自定义插件功能测试说明

## 功能概述

用户自定义插件加载功能，支持按路径分类加载 HTML/CSS/JS 插件文件。

## 实现内容

### 1. 核心代码实现

- ✅ `config/plugin.go`: 插件扫描和加载逻辑
  - `CustomPlugin` 结构体：定义插件数据结构
  - `ScanUserPlugins()`: 扫描并加载用户自定义插件
  - `GetCustomPluginsByScope()`: 按范围获取插件列表
  - `LoadBookPlugins()`: 按需加载特定书籍的插件（新功能）
  
- ✅ `config/config.go`: 添加 `CustomPlugins` 字段到 Config 结构体

- ✅ `templ/plugins/custom_plugin.templ`: 自定义插件渲染模板
  - `RenderCustomPlugin()`: 渲染单个插件
  - `RenderCustomPlugins()`: 根据页面路径渲染对应范围的插件
  - `extractBookID()`: 从 URL 路径中提取书籍 ID（新功能）

- ✅ `templ/plugins/manager.templ`: 集成自定义插件到主插件管理器

- ✅ `cmd/set_stores.go`: 添加 `LoadUserPlugins()` 函数

- ✅ `main.go`: 在程序启动时调用插件初始化

### 2. 测试插件文件

已在 `~/.config/comigo/plugins/` 目录下创建了以下测试插件：

```
~/.config/comigo/plugins/
├── global-test.js          # 全局 JS 插件
├── global-style.css        # 全局 CSS 插件  
├── test-html.html          # 全局 HTML 插件
├── shelf/
│   └── shelf-plugin.js     # Shelf 页面插件
├── flip/
│   ├── flip-plugin.js      # Flip 页面通用插件
│   └── aBcE4Fz/            # 书籍 aBcE4Fz 的专属插件（新功能）
│       └── book-plugin.js
└── scroll/
    ├── scroll-plugin.js    # Scroll 页面通用插件
    └── aBcE4Fz/            # 书籍 aBcE4Fz 的专属插件（新功能）
        └── book-plugin.js
```

## 测试步骤

### 1. 启用插件系统

在配置文件中设置：
```toml
EnablePlugin = true
```

或通过命令行参数启动：
```bash
./comigo --enable-plugin
```

### 2. 启动程序

```bash
cd /Users/bai/cvgo
go run . --enable-plugin --debug
```

### 3. 验证插件加载

#### 方法1：查看日志
启动程序时，如果开启了 Debug 模式，会在日志中看到：
```
成功加载 6 个自定义插件
  - [global] global-test.js (js)
  - [global] global-style.css (css)
  - [global] test-html.html (html)
  - [shelf] shelf-plugin.js (js)
  - [flip] flip-plugin.js (js)
  - [scroll] scroll-plugin.js (js)
```

#### 方法2：浏览器控制台检查

1. **访问 Shelf 页面** (`http://localhost:1234/` 或 `http://localhost:1234/shelf/xxx`)
   - 打开浏览器控制台，应该看到：
     - `✓ 全局自定义插件已加载`
     - `✓ Shelf 页面自定义插件已加载`
   - 页面右上角应显示红色边框的提示：`✓ 自定义 HTML 插件已加载`
   - 页面右下角应显示绿色文字：`✓ Global Plugin CSS Loaded`

2. **访问 Flip 页面** (`http://localhost:1234/flip/xxx`)
   - 控制台应显示：
     - `✓ 全局自定义插件已加载`
     - `✓ Flip 页面自定义插件已加载`
   - 全局 HTML 和 CSS 插件效果也应可见

3. **访问 Scroll 页面** (`http://localhost:1234/scroll/xxx`)
   - 控制台应显示：
     - `✓ 全局自定义插件已加载`
     - `✓ Scroll 页面自定义插件已加载`
   - 全局 HTML 和 CSS 插件效果也应可见

## 插件开发说明

### 插件文件放置规则

- **全局插件**：直接放在 `~/.config/comigo/plugins/` 目录下
  - 在所有页面生效
  
- **Shelf 插件**：放在 `~/.config/comigo/plugins/shelf/` 目录下
  - 在 `/` 和 `/shelf/*` 路径的页面生效
  
- **Flip 插件**：放在 `~/.config/comigo/plugins/flip/` 目录下
  - 在 `/flip/*` 路径的页面生效
  
- **Scroll 插件**：放在 `~/.config/comigo/plugins/scroll/` 目录下
  - 在 `/scroll/*` 路径的页面生效

### 支持的文件类型

- `.js` - JavaScript 脚本（自动包裹在 `<script>` 标签中）
- `.css` - 样式表（自动包裹在 `<style>` 标签中）
- `.html` - HTML 片段（直接插入页面）

### 加载顺序

1. 内置插件（按 EnabledPluginList 顺序）
2. 用户自定义全局插件
3. 用户自定义页面特定插件

### 注意事项

- 插件内容在程序启动时加载到内存，修改插件后需要重启程序
- 插件文件必须使用 UTF-8 编码
- 如果 `EnablePlugin = false` 或插件目录不存在，不会加载任何自定义插件
- 自定义插件不需要在 `EnabledPluginList` 中配置

## 代码验证状态

- ✅ 代码实现完整
- ✅ 无 linter 错误
- ✅ Templ 文件已正确生成
- ✅ 测试插件文件已创建
- ⏳ 运行时测试（需要用户在正常环境中执行）

## 故障排查

### 插件未加载

1. 检查 `EnablePlugin` 是否为 `true`
2. 检查插件目录是否存在：`~/.config/comigo/plugins/`
3. 检查插件文件扩展名是否正确（`.js`, `.css`, `.html`）
4. 开启 Debug 模式查看详细日志

### 插件效果不生效

1. 确认浏览器控制台是否有 JavaScript 错误
2. 检查插件代码语法是否正确
3. 确认访问的页面路径是否匹配插件的作用域

### 查看已加载的插件

开启 Debug 模式启动程序，查看启动日志中的插件加载信息。

## 书籍特定插件功能（新增）

### 功能说明

除了全局和页面范围的插件外，现在还支持为特定书籍创建专属插件。这些插件只在访问对应书籍时才会加载。

### 插件目录结构

```
~/.config/comigo/plugins/
├── flip/
│   ├── common-plugin.js        # 所有 flip 页面都会加载
│   ├── aBcE4Fz/                # 仅在访问 /flip/aBcE4Fz 时加载
│   │   ├── custom-viewer.js
│   │   └── book-style.css
│   └── anotherID/              # 仅在访问 /flip/anotherID 时加载
│       └── special-feature.js
└── scroll/
    ├── common-plugin.js        # 所有 scroll 页面都会加载
    └── aBcE4Fz/                # 仅在访问 /scroll/aBcE4Fz 时加载
        └── auto-scroll.js
```

### 插件加载顺序

访问 `http://localhost:1234/flip/aBcE4Fz` 时，插件按以下顺序加载：

1. **全局插件** - `plugins/*.{js,css,html}`
2. **Flip 通用插件** - `plugins/flip/*.{js,css,html}`
3. **书籍专属插件** - `plugins/flip/aBcE4Fz/*.{js,css,html}` ⭐ 新功能

### 使用场景

- 为特定漫画添加自定义阅读模式
- 为某本书添加特殊的键盘快捷键
- 针对特定内容添加专属样式调整
- 为特定书籍添加注释或标注功能

### 性能特性

- **按需加载**：书籍插件只在访问对应书籍时才读取，不占用启动时间
- **内存优化**：不会在启动时加载所有书籍的插件，节省内存
- **快速失败**：如果书籍插件目录不存在，直接跳过，不影响页面加载速度

### 测试书籍特定插件

#### 1. 创建书籍插件目录

```bash
# 为书籍 aBcE4Fz 创建 flip 专属插件
mkdir -p ~/.config/comigo/plugins/flip/aBcE4Fz

# 为书籍 aBcE4Fz 创建 scroll 专属插件
mkdir -p ~/.config/comigo/plugins/scroll/aBcE4Fz
```

#### 2. 创建测试插件文件

```bash
# Flip 页面的书籍专属插件
cat > ~/.config/comigo/plugins/flip/aBcE4Fz/book-plugin.js << 'EOF'
console.log('✓ Book aBcE4Fz flip plugin loaded');
window.bookFlipPluginLoaded = true;
EOF

# Scroll 页面的书籍专属插件
cat > ~/.config/comigo/plugins/scroll/aBcE4Fz/book-plugin.js << 'EOF'
console.log('✓ Book aBcE4Fz scroll plugin loaded');
window.bookScrollPluginLoaded = true;
EOF
```

#### 3. 验证加载效果

1. **访问 flip 页面**: `http://localhost:1234/flip/aBcE4Fz`
   - 打开浏览器控制台
   - 应该看到: `✓ Book aBcE4Fz flip plugin loaded`
   - 检查 `window.bookFlipPluginLoaded` 为 `true`

2. **访问 scroll 页面**: `http://localhost:1234/scroll/aBcE4Fz`
   - 打开浏览器控制台
   - 应该看到: `✓ Book aBcE4Fz scroll plugin loaded`
   - 检查 `window.bookScrollPluginLoaded` 为 `true`

3. **访问其他书籍**: `http://localhost:1234/flip/otherID`
   - 不应该加载 aBcE4Fz 的专属插件
   - 只会加载全局和 flip 通用插件

#### 4. Debug 模式查看

开启 Debug 模式时，访问包含书籍插件的页面会在日志中看到：

```
加载书籍 aBcE4Fz 的 flip 插件: 1 个
```

### 兼容性说明

- ✅ 完全向后兼容现有的全局和范围插件
- ✅ 不影响已有的插件加载逻辑
- ✅ 如果某本书没有专属插件，页面正常显示
- ✅ 书籍 ID 从 URL 路径中自动提取，无需额外配置

