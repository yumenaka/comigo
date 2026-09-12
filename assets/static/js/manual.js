(() => {
  const app = document.getElementById("ManualApp");
  if (!app) return;

  const localeFiles = { zh: "zh_CN.json", ja: "ja_JP.json", en: "en_US.json" };
  const localePaths = { zh: "", ja: "ja-JP/", en: "en-US/" };
  const routeLanguages = { zh: "zh", "ja-JP": "ja", "en-US": "en" };
  const page = app.dataset.page || "index";
  const routeLanguage = routeLanguages[app.dataset.language] || "zh";

  // 同一标签页保留手册入口；先记录再处理语言跳转，章节切换不会覆盖来源。
  const returnLink = document.querySelector("[data-manual-back]");
  for (const source of [
    new URLSearchParams(location.search).get("from"),
    document.referrer,
    sessionStorage.getItem("manual.returnTo"),
  ]) {
    if (!source) continue;
    try {
      const url = new URL(source, location.origin);
      // 来源参数来自 URL，只接受同源且不属于手册的地址。
      if (
        url.origin !== location.origin ||
        /^\/manual(?:\/|$)/.test(url.pathname)
      )
        continue;
      const target = url.pathname + url.search + url.hash;
      sessionStorage.setItem("manual.returnTo", target);
      returnLink.href = target;
      break;
    } catch {
      // 无效来源不影响手册加载，返回按钮保留首页链接。
    }
  }

  // 与 drawer 共用 i18next 当前语言及其持久化值，正文 JSON 仍独立加载。
  const appLanguage = (i18next.resolvedLanguage || i18next.language).split(
    "-",
  )[0];
  const language = localeFiles[appLanguage] ? appLanguage : "en";

  // 生成当前语言的手册章节地址。
  function pageURL(slug, language) {
    return `/manual/${localePaths[language]}${slug === "index" ? "" : slug}`;
  }

  if (routeLanguage !== language) {
    location.replace(pageURL(page, language) + location.hash);
    return;
  }

  const contentURL = window.ComiGoPath(
    `/assets/static/manual-content/${localeFiles[language]}`,
  );

  // 仅加载当前语言的一份 JSON；内容、导航和搜索索引共用这份数据。
  fetch(contentURL)
    .then((response) => {
      if (!response.ok) throw new Error(`Manual content: ${response.status}`);
      return response.json();
    })
    .then((manual) => renderManual(manual))
    .catch((error) => {
      console.error(error);
      document.querySelector("[data-manual-content]").textContent =
        i18next.t("err_network_error");
    });

  // 填充当前语言的标题、导航及正文。
  function renderManual(manual) {
    const ui = manual.ui;
    const current = manual.pages[page];
    if (!current) return;

    document.documentElement.lang = manual.lang;
    document.title = current.title === ui.manualTitle ? ui.manualTitle : `${current.title} - ${ui.manualTitle}`;
    document.querySelector('meta[name="description"]').content = manual.home.tagline;
    document.querySelectorAll("[data-manual-text]").forEach((element) => {
      element.textContent =
        ui[element.dataset.manualText] || element.textContent;
    });
    document.querySelectorAll("[data-manual-label]").forEach((element) => {
      element.setAttribute("aria-label", ui[element.dataset.manualLabel]);
    });
    document.querySelectorAll("[data-manual-search-modal]").forEach((input) => {
      input.placeholder = ui.search;
      input.setAttribute("aria-label", ui.search);
    });

    const languageSelect = document.querySelector("[data-manual-language]");
    languageSelect.value = { zh: "zh-CN", ja: "ja", en: "en-US" }[language];
    languageSelect.addEventListener("change", async () => {
      await i18next.changeLanguage(languageSelect.value);
      location.assign(pageURL(page, languageSelect.value.split("-")[0]));
    });

    renderNavigation(manual, ui);
    renderContent(manual, current, ui);
    renderSearch(manual, ui);
    bindMenu();
  }

  // 按入门与使用指南生成章节目录。
  function renderNavigation(manual, ui) {
    const groups = [
      [ui.start, manual.order.slice(0, 3)],
      [ui.guide, manual.order.slice(3)],
    ];
    document.querySelector("[data-manual-nav]").innerHTML = groups
      .map(
        ([label, slugs]) =>
          `<section class="manual-nav-group"><strong>${label}</strong>${slugs
            .map(
              (slug) =>
                `<a href="${pageURL(slug, routeLanguage)}"${slug === page ? ' aria-current="page"' : ""}>${manual.pages[slug].title}</a>`,
            )
            .join("")}</section>`,
      )
      .join("");
  }

  // 显示正文，并接入下载、复制和章节导航。
  function renderContent(manual, current, ui) {
    const content = document.querySelector("[data-manual-content]");
    content.innerHTML =
      page === "index" ? homeHTML(manual) + current.html : current.html;
    // 手册固定在根路径，正文中的 ComiGo 链接仍需跟随部署的 BasePath。
    content.querySelectorAll('a[href="/reader"]').forEach((link) => {
      link.href = window.ComiGoPath("/reader");
    });
    const download = document.querySelector("[data-manual-download]");
    if (download) renderDownload(download, ui);
    addCopyButtons(content, ui);
    renderOutline(content, ui);
    renderPageNavigation(manual, ui);
    // 正文异步加载，浏览器首次导航时锚点尚不存在，需要在渲染后定位。
    if (location.hash) {
      try {
        document
          .getElementById(decodeURIComponent(location.hash.slice(1)))
          ?.scrollIntoView();
      } catch {
        // 外部链接的无效百分号编码不应阻断菜单与搜索初始化。
      }
    }
  }

  // 生成手册首页的介绍与快捷入口。
  function homeHTML(manual) {
    const home = manual.home;
    return `<section class="manual-hero">
      <h1>${manual.ui.manualTitle}</h1><h2>${home.text}</h2><p>${home.tagline}</p>
      <div class="manual-hero-actions">${home.actions.map(([label, slug], index) => `<a class="manual-button${index ? " secondary" : ""}" href="${pageURL(slug, routeLanguage)}">${label}</a>`).join("")}</div>
      <div class="manual-features">${home.features.map(([title, description]) => `<section class="manual-feature"><strong>${title}</strong><p>${description}</p></section>`).join("")}</div>
    </section>`;
  }

  // 统一新旧代码块容器；复制失败时选中代码，允许用户手动复制。
  function addCopyButtons(content, ui) {
    content.querySelectorAll("pre").forEach((pre) => {
      let wrapper = pre.parentElement;
      if (!wrapper.matches('div[class*="language-"], .manual-code')) {
        wrapper = document.createElement("div");
        pre.before(wrapper);
        wrapper.append(pre);
      }
      wrapper.classList.add("manual-code");
      // 新代码块的类型在 code 上，旧格式在容器上；统一显示已有类型标签。
      if (!wrapper.querySelector(".lang")) {
        const language = [...(pre.querySelector("code")?.classList || []), ...wrapper.classList]
          .find((name) => name.startsWith("language-"))?.slice(9) || "text";
        const label = document.createElement("span");
        label.className = "lang";
        label.textContent = language;
        wrapper.prepend(label);
      }
      if (wrapper.querySelector(".manual-copy")) return;
      const button = document.createElement("button");
      button.type = "button";
      button.className = "manual-copy";
      button.textContent = ui.copy;
      button.setAttribute("aria-live", "polite");
      button.addEventListener("click", async () => {
        try {
          await navigator.clipboard.writeText(pre.textContent);
          button.textContent = ui.copied;
        } catch {
          const range = document.createRange();
          range.selectNodeContents(pre);
          const selection = window.getSelection();
          selection.removeAllRanges();
          selection.addRange(range);
          button.textContent = ui.copyFailed;
        }
        setTimeout(() => {
          button.textContent = ui.copy;
        }, 1500);
      });
      wrapper.append(button);
    });
  }

  // 根据正文标题生成桌面与手机目录。
  function renderOutline(content, ui) {
    const headings = [...content.querySelectorAll("h2[id], h3[id]")];
    const links = headings
      .map(
        (heading) =>
          `<a class="${heading.tagName.toLowerCase()}" href="#${heading.id}">${heading.childNodes[0]?.textContent.trim() || heading.textContent.trim()}</a>`,
      )
      .join("");
    document.querySelector("[data-manual-outline]").innerHTML = links;
    document.querySelector("[data-manual-outline-mobile]").innerHTML =
      links || `<span>${ui.onThisPage}</span>`;
    document
      .querySelector("[data-manual-outline-mobile]")
      .addEventListener("click", (event) => {
        if (event.target.closest("a"))
          event.currentTarget.closest("details").open = false;
      });
  }

  // 连接相邻章节。
  function renderPageNavigation(manual, ui) {
    const index = manual.order.indexOf(page);
    const previous = manual.order[index - 1];
    const next = manual.order[index + 1];
    document.querySelector("[data-manual-page-nav]").innerHTML =
      (previous
        ? `<a href="${pageURL(previous, routeLanguage)}">← ${ui.previous}<br><strong>${manual.pages[previous].title}</strong></a>`
        : "<span></span>") +
      (next
        ? `<a href="${pageURL(next, routeLanguage)}">${ui.next} →<br><strong>${manual.pages[next].title}</strong></a>`
        : "<span></span>");
  }

  // 共用原生 dialog 管理焦点与 Esc；入口只响应点击，避免关闭后恢复焦点又重开。
  function renderSearch(manual, ui) {
    const dialog = document.querySelector("[data-manual-search-dialog]");
    const modalInput = document.querySelector("[data-manual-search-modal]");
    const results = document.querySelector("[data-manual-search-results]");
    const status = document.querySelector("[data-manual-search-status]");
    document.querySelector("[data-manual-shortcut]").textContent =
      /Mac|iPhone|iPad/.test(navigator.userAgent) ? "⌘ K" : "Ctrl K";
    const index = manual.order.map((slug) => {
      const holder = document.createElement("div");
      holder.innerHTML =
        (slug === "index" ? homeHTML(manual) : "") + manual.pages[slug].html;
      holder
        .querySelectorAll(".lang, .header-anchor")
        .forEach((label) => label.remove());
      holder
        .querySelectorAll("p, h1, h2, h3, pre, li, td, th")
        .forEach((block) => block.append(" "));
      return {
        slug,
        title: manual.pages[slug].title,
        text: holder.textContent.replace(/\s+/g, " ").trim(),
      };
    });

    // 按关键词筛选章节，并显示命中内容。
    function search(query) {
      const value = query.trim().toLocaleLowerCase(manual.lang);
      results.replaceChildren();
      if (!value) {
        status.textContent = ui.searchHint;
        return;
      }
      const matches = index.filter((item) =>
        `${item.title} ${item.text}`
          .toLocaleLowerCase(manual.lang)
          .includes(value),
      );
      status.textContent = matches.length
        ? ui.resultCount.replace("{count}", matches.length)
        : ui.noResults;
      // 截取命中位置附近的正文；用文本节点展示代码，避免搜索摘要被解析为 HTML。
      matches.forEach((item) => {
        const link = document.createElement("a");
        link.className = "manual-search-result";
        link.href = pageURL(item.slug, routeLanguage);
        const title = document.createElement("strong");
        title.textContent = item.title;
        const excerpt = document.createElement("small");
        const at = item.text.toLocaleLowerCase(manual.lang).indexOf(value);
        const start = Math.max(0, at - 45);
        const snippet = item.text.slice(
          start,
          start + Math.max(160, value.length + 90),
        );
        excerpt.append(start ? "…" : "");
        if (at >= 0) {
          const offset = at - start;
          const mark = document.createElement("mark");
          mark.textContent = snippet.slice(offset, offset + value.length);
          excerpt.append(
            snippet.slice(0, offset),
            mark,
            snippet.slice(offset + value.length),
          );
        } else {
          excerpt.append(snippet);
        }
        if (start + snippet.length < item.text.length) excerpt.append("…");
        link.append(title, excerpt);
        results.append(link);
      });
    }

    // 打开搜索弹窗，将焦点移到输入框。
    function openSearch() {
      if (!dialog.open) dialog.showModal();
      search(modalInput.value);
      modalInput.focus();
    }

    document
      .querySelectorAll("[data-manual-search-open]")
      .forEach((button) => button.addEventListener("click", openSearch));
    // 仅在一次点击始末均落在遮罩上时关闭，拖选输入文字不会误关弹窗。
    let backdropPressed = false;
    dialog.addEventListener("pointerdown", (event) => {
      backdropPressed = event.target === dialog;
    });
    dialog.addEventListener("click", (event) => {
      if (backdropPressed && event.target === dialog) dialog.close();
      backdropPressed = false;
    });
    dialog.addEventListener("keydown", (event) => {
      // search 输入框会先吞掉 Esc 清空内容，这里统一为一次按键关闭。
      if (event.key === "Escape") {
        event.preventDefault();
        event.stopPropagation();
        dialog.close();
      }
    });
    modalInput.addEventListener("input", (event) => search(event.target.value));
    document.addEventListener("keydown", (event) => {
      if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") {
        event.preventDefault();
        openSearch();
      }
    });
  }

  // 处理手机目录的展开、关闭和键盘导航。
  function bindMenu() {
    const sidebar = document.querySelector(".manual-sidebar");
    const backdrop = document.querySelector(".manual-backdrop");
    const button = document.querySelector("[data-manual-menu]");
    const mobile = matchMedia("(max-width: 959px)");
    // 同步目录的显示状态与焦点可达性。
    function setOpen(open) {
      sidebar.classList.toggle("open", open);
      backdrop.classList.toggle("open", open);
      button.setAttribute("aria-expanded", String(open));
      sidebar.inert = mobile.matches && !open;
      if (open) sidebar.querySelector("button").focus();
    }
    setOpen(false);
    mobile.addEventListener("change", () => setOpen(false));
    button.addEventListener("click", () =>
      setOpen(!sidebar.classList.contains("open")),
    );
    // 关闭目录后将焦点还给菜单按钮。
    function closeMenu() {
      setOpen(false);
      button.focus();
    }
    backdrop.addEventListener("click", closeMenu);
    // 抽屉打开时保持键盘焦点在菜单内；搜索弹窗由原生 dialog 接管。
    document.addEventListener("keydown", (event) => {
      if (
        !sidebar.classList.contains("open") ||
        document.querySelector("dialog[open]")
      )
        return;
      if (event.key === "Escape") closeMenu();
      if (event.key === "Tab") {
        const items = [button, ...sidebar.querySelectorAll("button, a")];
        const index = items.indexOf(document.activeElement);
        event.preventDefault();
        items[
          (index + (event.shiftKey ? -1 : 1) + items.length) % items.length
        ].focus();
      }
    });
  }

  const systems = {
    Windows: [
      ["tray", "x86_64", "comigo-tray_latest_Windows_x86_64.zip"],
      ["tray", "arm64", "comigo-tray_latest_Windows_arm64.zip"],
      ["desktop", "x86_64", "comigo-desktop_latest_Windows_x86_64.zip"],
      ["desktop", "arm64", "comigo-desktop_latest_Windows_arm64.zip"],
      ["cli", "x86_64", "comi_latest_Windows_x86_64.zip"],
      ["cli", "arm64", "comi_latest_Windows_arm64.zip"],
      ["cli", "i386", "comi_latest_Windows_i386.zip"],
    ],
    macOS: [
      ["tray", "universal", "comigo-tray_latest_MacOS_universal.dmg"],
      ["desktop", "universal", "comigo-desktop_latest_MacOS_universal.dmg"],
      ["cli", "arm64", "comi_latest_MacOS_arm64.tar.gz"],
      ["cli", "x86_64", "comi_latest_MacOS_x86_64.tar.gz"],
    ],
    Linux: [
      ["tray", "x86_64", "comigo-tray_latest_Linux_x86_64.tar.gz"],
      ["tray", "arm64", "comigo-tray_latest_Linux_arm64.tar.gz"],
      ["desktop", "x86_64", "comigo-desktop_latest_Linux_x86_64.tar.gz"],
      ["desktop", "arm64", "comigo-desktop_latest_Linux_arm64.tar.gz"],
      ["cli", "x86_64", "comi_latest_Linux_x86_64.tar.gz"],
      ["cli", "arm64", "comi_latest_Linux_arm64.tar.gz"],
      ["cli", "armv7", "comi_latest_Linux_armv7.tar.gz"],
      ["cli", "i386", "comi_latest_Linux_i386.tar.gz"],
      ["deb", "amd64", "comi_latest_amd64.deb"],
      ["deb", "arm64", "comi_latest_arm64.deb"],
    ],
  };

  // 浏览器不能可靠判断 CPU 架构，因此只自动选系统，安装包架构仍由用户选择。
  function renderDownload(root, ui) {
    const detected = /Android|iPhone|iPad/i.test(navigator.userAgent)
      ? null
      : /Windows NT/i.test(navigator.userAgent)
        ? "Windows"
        : /Mac/i.test(navigator.userAgent)
          ? "macOS"
          : /Linux/i.test(navigator.userAgent)
            ? "Linux"
            : null;
    let selectedSystem = detected || "Windows";
    const selectedOptions = { Windows: 0, macOS: 0, Linux: 0 };
    root.className = "manual-download-panel";

    // 更新系统选项与安装包下载地址。
    function draw() {
      const option = systems[selectedSystem][selectedOptions[selectedSystem]];
      root.innerHTML = `<p>${detected ? `${ui.detected}：<strong>${detected}</strong>` : ui.choosePlatform}</p>
        <div class="manual-download-platforms" role="group" aria-label="${ui.switchPlatform}">${Object.keys(
          systems,
        )
          .map(
            (system) =>
              `<button type="button" data-system="${system}" class="${system === selectedSystem ? "active" : ""}" aria-pressed="${system === selectedSystem}"><span>${system}</span>${system === detected ? `<small>${ui.detectedBadge}</small>` : ""}</button>`,
          )
          .join("")}</div>
        <div class="manual-download-actions"><label><span>${ui.package}</span><select>${systems[selectedSystem].map(([edition, arch], index) => `<option value="${index}"${index === selectedOptions[selectedSystem] ? " selected" : ""}>${ui.editions[edition]} · ${ui.architectures[arch]}</option>`).join("")}</select>
        </label><a href="https://comigo.xyz/yumenaka/comigo/releases/download/latest/${option[2]}" target="_blank" rel="noopener">${ui.download} ${selectedSystem}</a></div>
        <code class="manual-download-file">${option[2]}</code>`;
      root.querySelectorAll("[data-system]").forEach((button) =>
        button.addEventListener("click", () => {
          selectedSystem = button.dataset.system;
          draw();
          root.querySelector(`[data-system="${selectedSystem}"]`).focus();
        }),
      );
      root.querySelector("select").addEventListener("change", (event) => {
        selectedOptions[selectedSystem] = Number(event.target.value);
        // 不重建 select，键盘切换安装包后仍保留焦点。
        const file =
          systems[selectedSystem][selectedOptions[selectedSystem]][2];
        root.querySelector(".manual-download-actions a").href =
          `https://comigo.xyz/yumenaka/comigo/releases/download/latest/${file}`;
        root.querySelector(".manual-download-file").textContent = file;
      });
    }
    draw();
  }
})();
