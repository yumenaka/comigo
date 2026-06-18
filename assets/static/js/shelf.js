//此文件静态导入，不需要编译

(() => {
  if (!(window.ComiGoIsWails?.() || window.ComiGoWails)) return;

  const text = (key, fallback) => {
    const value = window.i18next?.isInitialized ? window.i18next.t(key) : "";
    return value && value !== key ? value : fallback;
  };
  const menu = document.createElement("div");
  menu.id = "wails-book-card-menu";
  menu.className =
    "fixed z-50 hidden min-w-44 overflow-hidden rounded border border-gray-500 bg-base-100 text-base-content shadow-lg";
  menu.innerHTML = `
    <button type="button" data-action="open" class="block w-full px-3 py-2 text-left text-sm hover:bg-base-200"></button>
    <button type="button" data-action="copy" class="block w-full px-3 py-2 text-left text-sm hover:bg-base-200"></button>
  `;
  const setMenuText = () => {
    menu.querySelector('[data-action="open"]').textContent = text(
      "open_external_browser",
      "Open in external browser",
    );
    menu.querySelector('[data-action="copy"]').textContent = text(
      "systray_copy_url",
      "Copy Reading URL",
    );
  };
  setMenuText();
  window.i18next?.on?.("initialized", setMenuText);
  window.i18next?.on?.("languageChanged", setMenuText);
  document.body.appendChild(menu);

  let currentURL = "";
  const hideMenu = () => menu.classList.add("hidden");
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
    placeMenu(event);
  });
  document.addEventListener("click", hideMenu);
  window.addEventListener("resize", hideMenu);
  window.addEventListener("scroll", hideMenu, true);

  menu.addEventListener("click", async (event) => {
    const action = event.target.closest("button")?.dataset.action;
    if (!action || !currentURL) return;
    hideMenu();
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
      window.showToast?.(text("comigo_xyz_cli_install_copied", "Copied"), "success");
    } catch (error) {
      console.error("复制阅读地址失败:", error);
      window.showToast?.(text("err_network_error", "Failed"), "error");
    }
  });
})();
