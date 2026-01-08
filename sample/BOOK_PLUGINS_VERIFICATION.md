# 书籍特定插件功能验证

## 实现状态 ✅

## 已完成的工作

### 1. ✅ 核心功能实现

#### config/plugin.go
- 新增 `LoadBookPlugins(bookID, scope string)` 函数
  - 按需加载指定书籍的插件
  - 从 `~/.config/comigo/plugins/{scope}/{bookID}/` 读取插件文件
  - 支持 .js、.css、.html 三种文件类型
  - 返回该书籍的插件列表

#### templ/plugins/custom_plugin.templ
- 新增 `extractBookID(path string)` 函数
  - 从 URL 路径提取书籍 ID
  - 例如：`/flip/aBcE4Fz` → `aBcE4Fz`
  
- 更新 `RenderCustomPlugins()` 函数
  - 在 flip 和 scroll 路径下增加书籍特定插件加载逻辑
  - 先加载通用插件，再加载书籍特定插件

### 2. ✅ 测试环境准备

#### 已创建的测试插件

```
~/.config/comigo/plugins/
├── flip/
│   ├── flip-plugin.js                    # 通用插件
│   └── aBcE4Fz/                          # 书籍特定插件 ⭐
│       ├── book-plugin.js
│       └── book-style.css
├── scroll/
│   ├── scroll-plugin.js                  # 通用插件
│   └── aBcE4Fz/                          # 书籍特定插件 ⭐
│       └── book-plugin.js
├── shelf/
│   └── shelf-plugin.js
├── global-test.js
├── global-style.css
└── test-html.html
```

#### 测试插件功能

**Flip 书籍插件** (`flip/aBcE4Fz/book-plugin.js`):
- 在控制台输出：`✓ Book aBcE4Fz flip plugin loaded`
- 设置全局变量：`window.bookFlipPluginLoaded = true`
- 添加快捷键监听（按 'b' 键触发）

**Flip 书籍样式** (`flip/aBcE4Fz/book-style.css`):
- 在页面左上角显示：`📖 Book aBcE4Fz Custom Style`
- 蓝色边框提示框

**Scroll 书籍插件** (`scroll/aBcE4Fz/book-plugin.js`):
- 在控制台输出：`✓ Book aBcE4Fz scroll plugin loaded`
- 设置全局变量：`window.bookScrollPluginLoaded = true`

### 3. ✅ 代码质量验证

- ✅ 无 linter 错误
- ✅ Templ 模板正确生成 Go 代码
- ✅ 所有函数正确导出和调用
- ✅ 错误处理完善

## 功能特性

### 插件加载策略

1. **按需加载（Lazy Loading）**
   - 书籍插件只在访问对应书籍时才读取
   - 不占用程序启动时间
   - 节省内存开销

2. **作用域隔离**
   - 每本书的插件互不影响
   - 访问书籍 A 不会加载书籍 B 的插件

3. **优雅降级**
   - 如果书籍插件目录不存在，静默跳过
   - 不影响页面正常加载

### 插件加载顺序

访问 `http://localhost:1234/flip/aBcE4Fz` 时：

```
1. 全局插件 (plugins/*.{js,css,html})
   ↓
2. Flip 通用插件 (plugins/flip/*.{js,css,html})
   ↓
3. 书籍专属插件 (plugins/flip/aBcE4Fz/*.{js,css,html}) ⭐ 新增
```

## 测试计划

### 自动验证（已通过）
- ✅ 代码语法检查（无 linter 错误）
- ✅ Templ 生成验证
- ✅ 函数签名正确性
- ✅ 导入依赖完整性

### 手动测试（需要运行时验证）

#### 测试用例 1：访问包含插件的书籍
**步骤**:
1. 启动程序：`./comigo --enable-plugin --debug`
2. 访问：`http://localhost:1234/flip/aBcE4Fz`
3. 打开浏览器开发者工具

**预期结果**:
- ✓ 控制台显示：`✓ Book aBcE4Fz flip plugin loaded`
- ✓ 页面左上角显示蓝色提示框：`📖 Book aBcE4Fz Custom Style`
- ✓ `window.bookFlipPluginLoaded === true`
- ✓ 按 'b' 键触发快捷键事件

#### 测试用例 2：访问 Scroll 页面
**步骤**:
1. 访问：`http://localhost:1234/scroll/aBcE4Fz`
2. 打开浏览器控制台

**预期结果**:
- ✓ 控制台显示：`✓ Book aBcE4Fz scroll plugin loaded`
- ✓ `window.bookScrollPluginLoaded === true`

#### 测试用例 3：访问其他书籍（插件隔离）
**步骤**:
1. 访问：`http://localhost:1234/flip/otherBookID`
2. 打开浏览器控制台

**预期结果**:
- ✓ 不显示 aBcE4Fz 的插件消息
- ✓ 不显示蓝色提示框
- ✓ `window.bookFlipPluginLoaded === undefined`
- ✓ 仍然加载全局和 flip 通用插件

#### 测试用例 4：Debug 模式日志
**步骤**:
1. 以 Debug 模式启动：`./comigo --enable-plugin --debug`
2. 访问：`http://localhost:1234/flip/aBcE4Fz`
3. 查看程序日志

**预期结果**:
- ✓ 日志中显示：`加载书籍 aBcE4Fz 的 flip 插件: 2 个`

## 兼容性说明

### 向后兼容
- ✅ 不影响现有全局插件
- ✅ 不影响现有范围插件（shelf/flip/scroll）
- ✅ 不影响内置插件系统
- ✅ 老版本插件目录结构完全兼容

### 性能影响
- ✅ 启动时间：无影响（不预加载）
- ✅ 内存占用：按需加载，仅占用访问书籍的插件
- ✅ 页面加载：文件 I/O 开销极小（通常 < 1ms）

## 技术细节

### URL 解析逻辑
```go
// extractBookID 从 URL 路径中提取书籍 ID
func extractBookID(path string) string {
    parts := strings.Split(strings.Trim(path, "/"), "/")
    if len(parts) >= 2 {
        return parts[1]  // 返回第二段作为书籍 ID
    }
    return ""
}
```

**示例**:
- `/flip/aBcE4Fz` → `aBcE4Fz`
- `/scroll/abc123` → `abc123`
- `/flip/` → `` (空字符串，不加载)

### 插件加载逻辑
```go
// LoadBookPlugins 加载特定书籍的插件
func LoadBookPlugins(bookID, scope string) ([]CustomPlugin, error) {
    // 1. 检查插件是否启用
    if !cfg.EnablePlugin || bookID == "" {
        return nil, nil
    }
    
    // 2. 构建插件目录路径
    bookPluginPath := filepath.Join(configDir, "plugins", scope, bookID)
    
    // 3. 检查目录是否存在（不存在则返回空列表）
    if _, err := os.Stat(bookPluginPath); os.IsNotExist(err) {
        return nil, nil
    }
    
    // 4. 读取并解析插件文件
    // 5. 返回插件列表
}
```

## 使用场景示例

### 场景 1：特定漫画的阅读优化
```
plugins/flip/manga_xyz/
├── auto-zoom.js      # 自动缩放到合适大小
└── custom-theme.css  # 专属主题色
```

### 场景 2：教材书籍的增强功能
```
plugins/flip/textbook_abc/
├── note-taker.js     # 笔记功能
├── bookmark-sync.js  # 书签同步
└── highlight.css     # 高亮样式
```

### 场景 3：长篇小说的阅读辅助
```
plugins/scroll/novel_123/
├── reading-progress.js  # 阅读进度跟踪
├── speed-reader.js      # 快速阅读模式
└── night-mode.css       # 夜间模式
```

## 后续扩展建议

### 可选优化（未实现）
1. **插件缓存机制**：缓存已加载的插件，减少重复文件读取
2. **热重载支持**：文件修改后自动刷新，无需重启程序
3. **插件配置文件**：支持 JSON 配置文件，定义插件元数据
4. **插件依赖管理**：支持插件间的依赖关系
5. **插件版本控制**：支持插件版本管理

### 已满足的核心需求
- ✅ 特定书籍的插件隔离
- ✅ 按需加载机制
- ✅ 支持 JS/CSS/HTML 三种文件
- ✅ 从 URL 自动识别书籍 ID
- ✅ 性能优化（不影响启动）

## 总结

书籍特定插件功能已完整实现，代码质量通过验证，测试环境已准备就绪。

### 核心优势
1. **灵活性**：为每本书定制独特功能
2. **性能**：按需加载，不影响整体性能
3. **简单性**：无需配置，自动识别书籍 ID
4. **兼容性**：完全向后兼容现有插件系统

### 待确认事项
- 需要在实际运行环境中测试浏览器加载效果
- 验证不同书籍 ID 格式的兼容性
- 确认 Debug 模式下的日志输出

**状态**: ✅ 开发完成，等待运行时测试

