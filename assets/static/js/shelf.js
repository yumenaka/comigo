//此文件静态导入，不需要编译

const shelfText = (key, fallback) => {
  const value = window.i18next?.isInitialized ? window.i18next.t(key) : "";
  return value && value !== key ? value : fallback;
};

// shelfRescanStoreSuccessMessage 统一格式化书架页重扫结果。
const shelfRescanStoreSuccessMessage = (data) => {
  const added = data.newBooksCount ?? 0;
  const removed = data.removedBooksCount ?? 0;
  if (added > 0 && removed > 0) {
    return shelfText("rescan_store_added_removed", "{0} books added, {1} books removed")
      .replace("{0}", added)
      .replace("{1}", removed);
  }
  if (added > 0) {
    return shelfText("rescan_store_added", "{0} books added").replace("{0}", added);
  }
  if (removed > 0) {
    return shelfText("rescan_store_removed", "{0} books removed").replace("{0}", removed);
  }
  return shelfText("rescan_store_no_change", "No book count changes");
};

window.ComiGoShelf = window.ComiGoShelf || {};

// getShelfContentRoot 找到当前页可替换的书架主体，兼容主书架、子书架和空书架。
const getShelfContentRoot = (root = document) =>
  root.getElementById("ShelfMainArea") ||
  root.getElementById("book-shelf") ||
  root.getElementById("tab-contents");

// initShelfFragment 让新插入的 Alpine 片段重新绑定交互。
const initShelfFragment = (element) => {
  if (!element || !window.Alpine?.initTree) return;
  window.Alpine.initTree(element);
};

// refreshShelfHTML 重新获取当前页 HTML，只替换书架主体和标题，避免整页 reload。
const refreshShelfHTML = async () => {
  const response = await fetch(window.location.href, { cache: "no-store" });
  if (!response.ok) {
    throw new Error("refresh shelf html failed");
  }

  const nextDoc = new DOMParser().parseFromString(await response.text(), "text/html");
  const currentShelf = getShelfContentRoot();
  const nextShelf = getShelfContentRoot(nextDoc);
  if (!currentShelf || !nextShelf) {
    throw new Error("shelf content not found");
  }

  currentShelf.replaceWith(nextShelf);
  initShelfFragment(nextShelf);

  const currentTitle = document.getElementById("headerTitle");
  const nextTitle = nextDoc.getElementById("headerTitle");
  if (currentTitle && nextTitle) {
    currentTitle.replaceWith(nextTitle);
    initShelfFragment(nextTitle);
  }
};

// rescanShelfStore 调用书架重扫接口；单书库和全量重扫只差 URL 与请求体。
const rescanShelfStore = async (path, body, logMessage) => {
  window.showToast?.(shelfText("rescan_store_in_progress", "Scanning store, please wait..."), "info");

  try {
    const options = {
      method: "POST",
    };
    if (body) {
      options.headers = { "Content-Type": "application/json" };
      options.body = JSON.stringify(body);
    }
    const response = await fetch(window.ComiGoPath(path), options);

    if (!response.ok) {
      window.showToast?.(shelfText("err_rescan_store_failed", "Rescan failed"), "error");
      return;
    }

    const data = await response.json();
    window.showToast?.(shelfRescanStoreSuccessMessage(data), "success");
    await refreshShelfHTML();
  } catch (error) {
    console.error(logMessage, error);
    window.showToast?.(shelfText("err_network_error", "Failed"), "error");
  }
};

window.ComiGoShelf.rescanStore = (storeUrlB64) =>
  rescanShelfStore("/api/rescan-store", { storeUrl: storeUrlB64 }, "重新扫描书库失败:");

window.ComiGoShelf.rescanAllStores = () =>
  rescanShelfStore("/api/rescan-all-stores", null, "重新扫描全部书库失败:");

(() => {
  const isAndroidWails = window.ComiGoWails && /Android/i.test(navigator.userAgent);
  if (!isAndroidWails) return;
  const path = window.ComiGoRelativePath?.(window.location.pathname) || window.location.pathname;
  if (path !== "/" && path !== "/index.html") return;

  let startY = 0;
  let pulling = false;
  let refreshing = false;
  const pullThreshold = 72;
  const scrollTop = () => getShelfContentRoot()?.scrollTop || window.scrollY || 0;

  document.addEventListener(
    "touchstart",
    (event) => {
      if (event.touches.length !== 1 || scrollTop() > 0 || refreshing) return;
      startY = event.touches[0].clientY;
      pulling = true;
    },
    { passive: true },
  );

  document.addEventListener(
    "touchmove",
    (event) => {
      if (!pulling || event.touches.length !== 1) return;
      if (event.touches[0].clientY - startY < pullThreshold) return;
      pulling = false;
      refreshing = true;
      window.showToast?.(shelfText("loading", "Loading..."), "info");
      // Android 首页下拉刷新只需要重新取当前书架 HTML，不重启整个 WebView。
      refreshShelfHTML()
        .catch((error) => {
          console.error("下拉刷新书架失败:", error);
          window.showToast?.(shelfText("err_network_error", "Failed"), "error");
        })
        .finally(() => {
          refreshing = false;
        });
    },
    { passive: true },
  );

  ["touchend", "touchcancel"].forEach((eventName) => {
    document.addEventListener(eventName, () => {
      pulling = false;
    }, { passive: true });
  });
})();

(() => {
  if (!window.ComiGoIsWails?.()) return;

  const menu = document.createElement("div");
  menu.id = "wails-book-card-menu";
  menu.className =
    "fixed z-50 hidden min-w-44 overflow-hidden rounded border border-gray-500 bg-base-100 text-base-content shadow-lg";
  const style = document.createElement("style");
  // 菜单样式放在脚本里，避免静态 JS 中的新 Tailwind hover 类未进入编译产物。
  style.textContent = `
    #wails-book-card-menu button {
      transition: background-color 120ms ease, color 120ms ease;
    }
    #wails-book-card-menu button:hover,
    #wails-book-card-menu button:focus-visible {
      background: var(--primary, #3b82f6);
      color: var(--primary-content, #fff);
      outline: none;
    }
    #wails-book-card-menu button[data-danger]:hover,
    #wails-book-card-menu button[data-danger]:focus-visible {
      background: #dc2626;
      color: #fff;
    }
  `;
  menu.innerHTML = `
    <button type="button" data-action="open" class="block w-full px-3 py-2 text-left text-sm"></button>
    <div role="separator" class="border-t border-gray-500"></div>
    <button type="button" data-action="copy" class="block w-full px-3 py-2 text-left text-sm"></button>
    <div role="separator" data-delete-separator class="border-t border-gray-500"></div>
    <button type="button" data-action="delete" data-danger class="block w-full px-3 py-2 text-left text-sm"></button>
  `;
  const deleteButton = menu.querySelector('[data-action="delete"]');
  const deleteSeparator = menu.querySelector("[data-delete-separator]");
  const setMenuText = () => {
    menu.querySelector('[data-action="open"]').textContent = shelfText(
      "open_external_browser",
      "Open in external browser",
    );
    menu.querySelector('[data-action="copy"]').textContent = shelfText(
      "systray_copy_url",
      "Copy Reading URL",
    );
    deleteButton.textContent = shelfText("wails_delete_file", "Delete source file");
  };
  setMenuText();
  window.i18next?.on?.("initialized", setMenuText);
  window.i18next?.on?.("languageChanged", setMenuText);
  document.head.appendChild(style);
  document.body.appendChild(menu);

  let currentURL = "";
  let currentBookID = "";
  let currentCard = null;
  let longPressTimer = 0;
  let suppressNextClick = false;
  const hideMenu = () => menu.classList.add("hidden");
  const toggleDeleteAction = () => {
    const canDelete = currentCard?.dataset.wailsCanDeleteSource === "true";
    deleteButton.hidden = !canDelete;
    deleteSeparator.hidden = !canDelete;
  };
  const shareBase = () =>
    document.querySelector("[data-qrcode-base]")?.dataset.qrcodeBase ||
    window.location.origin + "/";
  const readURL = (card) =>
    window.ComiGoShareURL
      ? window.ComiGoShareURL(card.href, shareBase())
      : card.href;
  const placeMenu = (clientX, clientY) => {
    menu.classList.remove("hidden");
    const gap = 8;
    const left = Math.min(clientX, window.innerWidth - menu.offsetWidth - gap);
    const top = Math.min(clientY, window.innerHeight - menu.offsetHeight - gap);
    menu.style.left = Math.max(gap, left) + "px";
    menu.style.top = Math.max(gap, top) + "px";
  };
  const openCardMenu = (card, clientX, clientY) => {
    currentURL = readURL(card);
    currentBookID = card.dataset.wailsBookId || "";
    currentCard = card;
    toggleDeleteAction();
    placeMenu(clientX, clientY);
  };

  // 只接管 Wails 书籍卡片的右键，保留普通浏览器环境不变。
  document.addEventListener("contextmenu", (event) => {
    const card = event.target.closest("[data-wails-book-card]");
    if (!card) return;
    event.preventDefault();
    openCardMenu(card, event.clientX, event.clientY);
  });
  document.addEventListener(
    "click",
    (event) => {
      if (suppressNextClick && !event.target.closest("#wails-book-card-menu")) {
        suppressNextClick = false;
        event.preventDefault();
        event.stopPropagation();
        return;
      }
      hideMenu();
    },
    true,
  );
  window.addEventListener("resize", hideMenu);
  window.addEventListener("scroll", hideMenu, true);

  // Android WebView 长按不一定触发 contextmenu，手动给书卡补一个长按菜单。
  document.addEventListener(
    "touchstart",
    (event) => {
      const card = event.target.closest("[data-wails-book-card]");
      if (!card || event.touches.length !== 1) return;
      const touch = event.touches[0];
      clearTimeout(longPressTimer);
      longPressTimer = window.setTimeout(() => {
        suppressNextClick = true;
        openCardMenu(card, touch.clientX, touch.clientY);
      }, 650);
    },
    { passive: true },
  );
  ["touchmove", "touchend", "touchcancel"].forEach((eventName) => {
    document.addEventListener(eventName, () => clearTimeout(longPressTimer), {
      passive: true,
    });
  });

  menu.addEventListener("click", async (event) => {
    suppressNextClick = false;
    const action = event.target.closest("button")?.dataset.action;
    if (!action) return;
    hideMenu();
    if (action === "delete") {
      if (!currentBookID) return;
      try {
        const result = await deleteBookSource(currentBookID);
        if (!result.deleted) return;
        const deletedCard = currentCard;
        try {
          await refreshShelfHTML();
        } catch (refreshError) {
          console.error("刷新删除后的书架失败:", refreshError);
          deletedCard?.remove();
        }
        window.showToast?.(
          result.message || shelfText("wails_delete_file_success", "Moved to system trash"),
          "success",
        );
      } catch (error) {
        console.error("删除书籍源文件失败:", error);
        window.showToast?.(
          String(error?.message || error || shelfText("wails_delete_file_failed", "Delete failed")),
          "error",
        );
      }
      return;
    }
    if (!currentURL) return;
    if (action === "open") {
      window.ComiGoOpenExternalURL?.(currentURL);
      return;
    }
    try {
      if (navigator.clipboard?.writeText) {
        await navigator.clipboard.writeText(currentURL);
      } else {
        const input = document.createElement("textarea");
        input.value = currentURL;
        input.style.position = "fixed";
        input.style.opacity = "0";
        document.body.appendChild(input);
        input.select();
        document.execCommand("copy");
        input.remove();
      }
      window.showToast?.(shelfText("comigo_xyz_cli_install_copied", "Copied"), "success");
    } catch (error) {
      console.error("复制阅读地址失败:", error);
      window.showToast?.(shelfText("err_network_error", "Failed"), "error");
    }
  });

  // deleteBookSource 统一走 HTTP 子路由，由 Wails 后端确认并按平台删除。
  async function deleteBookSource(bookID) {
    const response = await fetch(window.ComiGoPath("/api/wails/delete-book-file"), {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ bookId: bookID }),
    });
    if (!response.ok) {
      throw new Error(await response.text());
    }
    return response.json();
  }
})();
