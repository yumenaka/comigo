// 此文件保留占位：播放器逻辑已提前注入到模板中（见 templ/pages/player/player_main_area.templ）
// 原因：本项目会先执行 /script/main.js（其中 Alpine.start() 会触发 alpine:init），
// 然后才加载页面级脚本（如 /script/player.js），因此在这里监听 alpine:init 会错过事件。
//
// 注意：不要在此处覆盖 window.playerData。
'use strict';
