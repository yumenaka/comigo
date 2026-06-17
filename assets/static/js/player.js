// 播放器 Alpine 组件。
// 此文件由播放器模板提前加载，保证 window.playerData 在 Alpine.start() 前存在。
'use strict';

(function () {
  function readJSONScript(id, fallback) {
    const el = document.getElementById(id);
    if (!el) return fallback;
    try {
      return JSON.parse(el.textContent);
    } catch (error) {
      console.error('[player] Failed to parse JSONScript:', id, error);
      return fallback;
    }
  }

  function toast(message, type) {
    if (typeof window.showToast === 'function') {
      window.showToast(message, type);
      return;
    }
    console.log('[toast]', type || 'info', message);
  }

  function comigoPath(path) {
    return window.ComiGoPath ? window.ComiGoPath(path) : path;
  }

  // Header 也会异步更新 document.title；播放器切歌前再探测一次，避免丢掉版本后缀。
  function readComigoTitleSuffix() {
    const title = document.title || '';
    const index = title.indexOf(' - Comigo');
    return index >= 0 ? title.slice(index) : '';
  }

  // 根据文件扩展名猜测 MIME，后端未提供时作为兜底。
  function guessMimeType(filename, fallbackType) {
    if (!filename) {
      if (fallbackType === 'video') return 'video/mp4';
      if (fallbackType === 'audio') return 'audio/mpeg';
      return '';
    }
    const lower = String(filename).toLowerCase();
    const dot = lower.lastIndexOf('.');
    const ext = dot >= 0 ? lower.slice(dot) : '';
    switch (ext) {
      case '.mp3':
        return 'audio/mpeg';
      case '.m4a':
        return 'audio/mp4';
      case '.aac':
        return 'audio/aac';
      case '.wav':
        return 'audio/wav';
      case '.ogg':
        return 'audio/ogg';
      case '.wma':
        return 'audio/x-ms-wma';
      case '.mp4':
        return 'video/mp4';
      case '.m4v':
        return 'video/x-m4v';
      case '.mov':
        return 'video/quicktime';
      case '.webm':
        return 'video/webm';
      case '.avi':
        return 'video/x-msvideo';
      case '.flv':
        return 'video/x-flv';
      default:
        return '';
    }
  }

  // normalizeMediaItem 只使用后端公开的播放器字段，不读取本地 BookPath/StoreUrl。
  function normalizeMediaItem(item) {
    if (!item || !item.id || !item.title) return null;
    const rawUrl =
      item.rawUrl || comigoPath(`/api/raw/${item.id}/${encodeURIComponent(item.title)}`);
    return {
      id: item.id,
      title: item.title,
      type: item.type,
      rawUrl,
      coverUrl: item.coverUrl || comigoPath(`/api/get-cover?id=${encodeURIComponent(item.id)}&resize_height=352`),
      mimeType: item.mimeType || guessMimeType(item.title, item.type),
    };
  }

  window.playerData = function () {
    const playerPayload = readJSONScript('PlayerData', {});
    const currentMedia = normalizeMediaItem(playerPayload.current);
    const playlistRaw = Array.isArray(playerPayload.playlist) ? playerPayload.playlist : [];
    const playlistData = playlistRaw.map(normalizeMediaItem).filter(Boolean);

    if (playlistData.length === 0 && currentMedia) {
      playlistData.push(currentMedia);
    }

    let currentIndex = 0;
    if (currentMedia) {
      const idx = playlistData.findIndex((x) => x.id === currentMedia.id);
      if (idx >= 0) currentIndex = idx;
    }

    return {
      playlistData,
      currentIndex,
      currentMedia: playlistData[currentIndex] || currentMedia || null,
      showFallbackCover: false,
      titleSuffix: '',
      isPlaying: false,
      currentTime: 0,
      duration: 0,
      progress: 0,
      player: null,

      // 初始化播放器状态，保证动态 source 与持久化音量在首帧同步。
      init() {
        this.player = document.getElementById('mediaPlayer');
        if (!this.titleSuffix) {
          this.titleSuffix = readComigoTitleSuffix();
        }
        this.applyNowPlayingTitle();
        this.$nextTick(() => {
          this.applyVolumeFromStore();
          if (this.player) this.player.load();
          this.scrollCurrentPlaylistItem('auto');
        });
      },

      // applyVolumeFromStore 将全局存储的音量/静音状态应用到媒体元素。
      applyVolumeFromStore() {
        if (!this.player) return;
        try {
          if (!(window.Alpine && Alpine.store && Alpine.store('global'))) return;
          const globalStore = Alpine.store('global');
          let volume = Number(globalStore.playerVolume);
          if (Number.isNaN(volume)) volume = 100;
          volume = Math.max(0, Math.min(100, volume));
          this.player.volume = volume / 100;
          this.player.muted = !!globalStore.playerMuted;
        } catch (error) {
          console.error('[player] applyVolumeFromStore error:', error);
        }
      },

      // toggleMute 切换静音状态并持久化到全局 store。
      toggleMute() {
        if (!(window.Alpine && Alpine.store && Alpine.store('global'))) return;
        const globalStore = Alpine.store('global');
        globalStore.playerMuted = !globalStore.playerMuted;
        if (this.player) this.player.muted = !!globalStore.playerMuted;
      },

      // applyNowPlayingTitle 同步 header 标题与浏览器标签页标题。
      applyNowPlayingTitle() {
        const name = this.currentMedia?.title || '';
        if (!name) return;
        if (!this.titleSuffix) {
          this.titleSuffix = readComigoTitleSuffix();
        }
        document.title = name + (this.titleSuffix || '');
        const el = document.getElementById('headerTitle');
        if (!el) return;
        try {
          el.setAttribute('title', name);
          const titleText = el.querySelector('[data-header-title-text]') || el.querySelector('a');
          if (titleText) {
            titleText.textContent = name;
          } else {
            el.textContent = name;
          }
        } catch (error) {
          console.error('[player] applyNowPlayingTitle error:', error);
        }
      },

      onLoadedMetadata(event) {
        this.duration = event?.target?.duration || 0;
      },

      onTimeUpdate(event) {
        this.currentTime = event?.target?.currentTime || 0;
        this.progress = this.duration > 0 ? (this.currentTime / this.duration) * 100 : 0;
      },

      onEnded() {
        try {
          if (window.Alpine && Alpine.store && Alpine.store('global')) {
            const globalStore = Alpine.store('global');
            if (globalStore.autoPlayNext === false) return;
            if (this.currentIndex >= this.playlistData.length - 1) {
              if (globalStore.loopPlaylist === true) {
                this.playMedia(0, { forcePlay: true });
              } else {
                if (this.player) this.player.pause();
                this.isPlaying = false;
              }
              return;
            }
          }
        } catch (error) {
          console.error('[player] onEnded flags error:', error);
        }
        this.playNext({ forcePlay: true });
      },

      togglePlay() {
        if (!this.player) return;
        if (this.player.paused) {
          this.player.play().catch((error) => {
            console.error('[player] play failed:', error);
            toast(i18next.t('play_failed') || '播放失败', 'error');
          });
        } else {
          this.player.pause();
        }
      },

      // playMedia 切换曲目；ended 场景使用 forcePlay 强制自动播放下一项。
      playMedia(index, opts = {}) {
        if (index < 0 || index >= this.playlistData.length) return;

        const forcePlay = !!opts.forcePlay;
        this.currentIndex = index;
        this.currentMedia = this.playlistData[index] || null;
        this.showFallbackCover = false;
        this.applyNowPlayingTitle();

        if (!this.player) return;

        const wasPlaying = !this.player.paused;
        const shouldAutoPlay = forcePlay ? true : wasPlaying;

        this.player.pause();
        this.isPlaying = false;

        this.$nextTick(() => {
          this.applyVolumeFromStore();
          this.player.load();
          this.applyVolumeFromStore();
          this.scrollCurrentPlaylistItem('smooth');
          if (shouldAutoPlay) {
            this.player
              .play()
              .catch((error) => console.error('[player] autoplay failed:', error));
          }
        });
      },

      // scrollCurrentPlaylistItem 刷新或切歌后，将播放列表滚到当前媒体。
      scrollCurrentPlaylistItem(behavior = 'smooth') {
        const root = document.getElementById('PlayerMainArea');
        const activeItem = root?.querySelector('[data-player-active="true"]');
        if (!activeItem) return;
        activeItem.scrollIntoView({ block: 'center', inline: 'nearest', behavior });
      },

      playPrevious(opts = {}) {
        if (this.currentIndex <= 0) {
          toast(i18next.t('first_media') || '已经是第一个了', 'info');
          return;
        }
        this.playMedia(this.currentIndex - 1, opts);
      },

      playNext(opts = {}) {
        if (this.currentIndex >= this.playlistData.length - 1) {
          toast(i18next.t('last_media') || '已经是最后一个了', 'info');
          if (this.player) this.player.pause();
          this.isPlaying = false;
          return;
        }
        this.playMedia(this.currentIndex + 1, opts);
      },

      seekTo(value) {
        if (!this.player || !this.duration) return;
        this.player.currentTime = (Number(value) / 100) * this.duration;
      },

      formatTime(seconds) {
        if (!seconds || Number.isNaN(seconds)) return '00:00';
        const minutes = Math.floor(seconds / 60);
        const secs = Math.floor(seconds % 60);
        return `${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
      },
    };
  };
})();
