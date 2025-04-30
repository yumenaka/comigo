// Alpine 使用 Persist 插件，用会话 cookie 作为存储
// https://alpinejs.dev/plugins/persist#custom-storage
// 定义自定义存储对象，公开 getItem 函数和 setItem 函数
window.cookieStorage = {
    getItem(key) {
        let cookies = document.cookie.split(";");
        for (let i = 0; i < cookies.length; i++) {
            let cookie = cookies[i].split("=");
            if (key === cookie[0].trim()) {
                return decodeURIComponent(cookie[1]);
            }
        }
        return null;
    },
    setItem(key, value) {
        document.cookie = `${key}=${encodeURIComponent(value)}; SameSite=Lax`;//SameSite设置默认值（Lax），防止控制台报错。加载图像或框架（frame）的请求将不会包含用户的 Cookie。
    }
}

// // 然后就可以这样使用使用 cookieStorage 作为 Persist 插件的存储了
// Alpine.store('cookie', {
//     someCookieKey: Alpine.$persist(false).using(cookieStorage).as('someCookieKey'),
// })