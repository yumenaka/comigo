// Book lBjY4Qz 的 Flip 页面专属插件
console.log('✓ Book lBjY4Qz flip plugin loaded');
window.bookFlipPluginLoaded = true;

// 示例：为这本书添加特殊快捷键
document.addEventListener('keydown', function(e) {
    if (e.key === 'b') {
        console.log('Book-specific shortcut "b" pressed for lBjY4Qz');
    }
});
