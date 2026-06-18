//此文件静态导入，不需要编译

const shelfText = (key, fallback) => {
  const value = window.i18next?.isInitialized ? window.i18next.t(key) : "";
  return value && value !== key ? value : fallback;
};

// shelfRescanStoreSuccessMessage 统一格式化书架页重扫结果。
const shelfRescanStoreSuccessMessage = (data) =>
  shelfText("rescan_store_success", "Store scan completed, {0} new books added, {1} books removed")
    .replace("{0}", data.newBooksCount ?? 0)
    .replace("{1}", data.removedBooksCount ?? 0);

window.ComiGoShelf = window.ComiGoShelf || {};
window.ComiGoShelf.rescanAllStores = async () => {
  window.showToast?.(shelfText("rescan_store_in_progress", "Scanning store, please wait..."), "info");

  try {
    const response = await fetch(window.ComiGoPath("/api/rescan-all-stores"), {
      method: "POST",
    });

    if (!response.ok) {
      window.showToast?.(shelfText("err_rescan_store_failed", "Rescan failed"), "error");
      return;
    }

    const data = await response.json();
    window.showToast?.(shelfRescanStoreSuccessMessage(data), "success");
  } catch (error) {
    console.error("重新扫描全部书库失败:", error);
    window.showToast?.(shelfText("err_network_error", "Failed"), "error");
  }
};

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
  // 删除源文件只走 Wails 绑定；普通网页没有 window.go，不能调用这个能力。
  const deleteBookFile = () => window.go?.main?.App?.DeleteBookFile;
  const hideMenu = () => menu.classList.add("hidden");
  const toggleDeleteAction = () => {
    const canDelete = typeof deleteBookFile() === "function";
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
  const placeMenu = (event) => {
    menu.classList.remove("hidden");
    const gap = 8;
    const left = Math.min(event.clientX, window.innerWidth - menu.offsetWidth - gap);
    const top = Math.min(event.clientY, window.innerHeight - menu.offsetHeight - gap);
    menu.style.left = Math.max(gap, left) + "px";
    menu.style.top = Math.max(gap, top) + "px";
  };

  // 只接管 Wails 书籍卡片的右键，保留普通浏览器环境不变。
  document.addEventListener("contextmenu", (event) => {
    const card = event.target.closest("[data-wails-book-card]");
    if (!card) return;
    event.preventDefault();
    currentURL = readURL(card);
    currentBookID = card.dataset.wailsBookId || "";
    currentCard = card;
    toggleDeleteAction();
    placeMenu(event);
  });
  document.addEventListener("click", hideMenu);
  window.addEventListener("resize", hideMenu);
  window.addEventListener("scroll", hideMenu, true);

  menu.addEventListener("click", async (event) => {
    const action = event.target.closest("button")?.dataset.action;
    if (!action) return;
    hideMenu();
    if (action === "delete") {
      const deleteFn = deleteBookFile();
      if (!currentBookID || typeof deleteFn !== "function") return;
      try {
        const deleted = await deleteFn(currentBookID);
        if (!deleted) return;
        currentCard?.remove();
        window.showToast?.(shelfText("wails_delete_file_success", "Moved to system trash"), "success");
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
})();
