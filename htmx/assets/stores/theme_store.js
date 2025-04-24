// 自定义主题
Alpine.store('theme', {
    theme: Alpine.$persist('light').as('theme'),
    interfaceColor: '#F5F5E4',
    backgroundColor: '#E0D9CD',
    textColor: '#000000',
    toggleTheme() {
        this.theme = this.theme === 'light' ? 'dark' : 'light'
    },
}) 