// 对已启动的隔离 Manual 预览运行：playwright-cli run-code "$(< scripts/manual-browser-check.js)"
// 使用专用浏览器会话；会修改该会话的手册语言和视口，不启动 ComiGo 书库。
// prettier-ignore
async (page) => {
  const check = (ok, message) => {
    if (!ok) throw new Error(message);
  };
  const origin = await page.evaluate(() => location.origin);
  // 禁用测试会话 HTTP 缓存，避免热更新时拿到上一版静态脚本。
  await page.route("**/*", (route) => route.continue());
  const dialog = page.locator("[data-manual-search-dialog]");
  const input = page.locator("[data-manual-search-modal]");
  const trigger = page.locator(".manual-sidebar [data-manual-search-open]");
  const close = page.locator("dialog form button");
  const errors = [];
  page.on("pageerror", (error) => errors.push(error.message));
  page.on("request", (request) => {
    if (request.url().includes("/api/sse"))
      errors.push("Manual connected to SSE");
  });
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.evaluate(() => localStorage.setItem("i18nextLng", "zh-CN"));
  const returnTo = "/reader?sort=name#reading";
  await page.goto(`${origin}/manual/?from=${encodeURIComponent(returnTo)}`);
  await trigger.waitFor();
  await trigger.click();
  await close.click();
  // 恢复焦点是异步浏览器行为：下一帧仍必须保持关闭。
  await page.evaluate(() => new Promise(requestAnimationFrame));
  check(!(await dialog.isVisible()), "Close button reopened search");
  check(
    await trigger.evaluate((el) => el === document.activeElement),
    "Focus was not restored",
  );
  await trigger.press("Enter");
  await input.fill("ComiGo");
  await input.press("Escape");
  await page.evaluate(() => new Promise(requestAnimationFrame));
  check(!(await dialog.isVisible()), "Escape reopened search");
  await page.keyboard.press("ControlOrMeta+k");
  await input.fill("远程 ComiGo");
  check(
    (await page.locator(".manual-search-result").count()) > 0,
    "Chinese results missing",
  );
  check(
    (await page.locator(".manual-search-result mark").count()) === (await page.locator(".manual-search-result").count()),
    "Match context missing",
  );
  await input.fill("<script>");
  check(
    (await page.locator(".manual-search-results script").count()) === 0,
    "Search parsed HTML",
  );
  await input.fill("no-match-849264");
  check(
    (await page.locator("[role=status]").textContent()).includes("没有"),
    "Empty state missing",
  );
  await page.mouse.click(5, 5);
  check(!(await dialog.isVisible()), "Backdrop did not close search");

  // 每种语言都实际加载页面、搜索并关闭；日志不属于内容页索引。
  for (const [language, segment, query] of [
    ["zh-CN", "", "远程"],
    ["ja", "ja-JP/", "リモート"],
    ["en-US", "en-US/", "remote"],
  ]) {
    await page.locator("[data-manual-language]").selectOption(language);
    await page.waitForURL(`${origin}/manual/${segment}`);
    await page.locator("[data-manual-nav] a").first().waitFor();
    check(await page.locator("[data-manual-language] option").count() === 3, "Manual must have only three languages");
    check(await page.evaluate(() => i18next.language === localStorage.getItem("i18nextLng")) && await page.locator("[data-manual-language]").inputValue() === language, "Manual and app language differ");
    check((await page.locator("[data-manual-back]").getAttribute("href")) === returnTo, "Language switch lost entry page");
    check(
      (await page.locator("[data-manual-nav] a").count()) === 10,
      "Navigation incomplete",
    );
    await trigger.click();
    await input.fill(query);
    check(
      (await page.locator(".manual-search-result").count()) > 0,
      `${language} search failed`,
    );
    await close.click();
    check(!(await dialog.isVisible()), `${language} close failed`);
  }
  await page.goto(`${origin}/manual/en-US/install#install-the-cli`);
  await page.locator(".manual-copy").first().waitFor();
  const headingTop = await page
    .locator("#install-the-cli")
    .evaluate((el) => el.getBoundingClientRect().top);
  check(
    headingTop >= 56 && headingTop < 300,
    "Deep-link heading hidden or not reached",
  );
  await page.locator(".manual-copy").first().click();
  await page.waitForFunction(
    () => document.querySelector(".manual-copy").textContent !== "Copy",
  );
  check(
    (await page.locator(".manual-copy").first().textContent()) === "Copied",
    "Copy did not succeed",
  );
  // 拒绝剪贴板权限时仍给出可操作的手动复制提示。
  await page.evaluate(() => {
    navigator.clipboard.writeText = async () => { throw new Error("Denied by test"); };
  });
  await page.locator(".manual-copy").first().click();
  await page.waitForFunction(() => document.querySelector(".manual-copy").textContent.includes("manually"));
  check(await page.evaluate(() => getSelection().toString().length > 0), "Copy fallback did not select code");

  // 手机抽屉、焦点、搜索叠层、目录收起和窄屏溢出。
  await page.setViewportSize({ width: 390, height: 844 });
  const menu = page.locator("[data-manual-menu]");
  await menu.click();
  check(
    (await menu.getAttribute("aria-expanded")) === "true",
    "Drawer did not open",
  );
  await trigger.click();
  await close.click();
  check(!(await dialog.isVisible()), "Drawer search reopened");
  await page.keyboard.press("Escape");
  check(
    (await menu.getAttribute("aria-expanded")) === "false",
    "Drawer Escape failed",
  );
  check(
    await page.locator(".manual-sidebar").evaluate((el) => el.inert),
    "Hidden drawer accepts focus",
  );
  await page.locator(".manual-header [data-manual-search-open]").click();
  await input.fill("remote");
  await page.mouse.click(2, 2);
  check(!(await dialog.isVisible()), "Mobile backdrop failed");
  await page.locator(".manual-mobile-outline summary").click();
  await page.locator("[data-manual-outline-mobile] a").first().click();
  check(
    !(await page.locator(".manual-mobile-outline").evaluate((el) => el.open)),
    "Outline stayed expanded",
  );
  for (const width of [320, 390, 768, 1024, 1440]) {
    await page.setViewportSize({ width, height: 900 });
    check(
      await page.evaluate(
        () => document.documentElement.scrollWidth <= innerWidth,
      ),
      `Overflow at ${width}px`,
    );
  }
  await page.locator("[data-manual-language]").selectOption("zh-CN");
  await page.goto(`${origin}/manual/`);
  await page.locator("[data-manual-nav] a").first().waitFor();
  await page.locator('[data-system="Linux"]').click();
  const packages = page.locator(".manual-download-actions select");
  await packages.focus();
  await packages.selectOption({ label: "Debian/Ubuntu 包 · ARM64（64 位）" });
  check(await packages.evaluate((el) => el === document.activeElement), "Package selection lost focus");
  check((await page.locator(".manual-download-actions a").getAttribute("href")).endsWith("comi_latest_arm64.deb"), "Wrong package URL");
  // 逐一验证下载选项，所有安装包必须使用官网 latest 固定转发地址。
  for (const system of ["Windows", "macOS", "Linux"]) {
    await page.locator(`[data-system="${system}"]`).click();
    const options = await packages.locator("option").count();
    for (let index = 0; index < options; index++) {
      await packages.selectOption(String(index));
      const href = await page.locator(".manual-download-actions a").getAttribute("href");
      check(href.startsWith("https://comigo.xyz/yumenaka/comigo/releases/download/latest/") && href.includes("_latest_"), "Download URL is not a fixed latest URL");
    }
  }
  // 所有章节的新旧代码块共用样式；明暗主题和窄屏下复制按钮均位于右下角。
  for (const language of ["", "en-US/", "ja-JP/"]) {
    for (const chapter of ["index", "install", "quick-start", "reading", "library", "desktop", "comigo-omarchy", "deployment", "development", "faq"]) {
      await page.goto(`${origin}/manual/${language}${chapter}`);
      await page.locator("[data-manual-nav] a").first().waitFor();
      if (await page.locator("[data-manual-content] pre").count() === 0) continue;
      for (const theme of ["light", "dark"]) {
        await page.evaluate((value) => { document.body.dataset.theme = value; }, theme);
        for (const width of [390, 1440]) {
          await page.setViewportSize({ width, height: 900 });
          const problem = await page.locator("[data-manual-content]").evaluate((content) => {
            for (const pre of content.querySelectorAll("pre")) {
              const wrapper = pre.parentElement;
              const buttons = wrapper.querySelectorAll(".manual-copy");
              if (!wrapper.classList.contains("manual-code") || buttons.length !== 1) return "Missing code wrapper/button";
              const style = getComputedStyle(wrapper);
              if (style.borderTopWidth !== "1px" || style.borderTopStyle !== "solid") return "Missing code border";
              const box = wrapper.getBoundingClientRect();
              const button = buttons[0].getBoundingClientRect();
              if (Math.abs(box.right - button.right - 10) > 3 || Math.abs(box.bottom - button.bottom - 10) > 3) return "Copy button is not bottom-right";
              const preStyle = getComputedStyle(pre);
              const buttonStyle = getComputedStyle(buttons[0]);
              if (buttonStyle.position !== "absolute" || preStyle.paddingBottom !== preStyle.paddingTop) return "Copy button reserves a row";
              const alpha = Number(buttonStyle.backgroundColor.match(/rgba\(.+, ([\d.]+)\)/)?.[1]);
              if (!(alpha > 0 && alpha < 1)) return "Copy background is not translucent";
              const type = [...pre.querySelector("code").classList, ...wrapper.classList].find((name) => name.startsWith("language-"))?.slice(9);
              const labels = wrapper.querySelectorAll(".lang");
              if (!type || labels.length !== 1 || labels[0].textContent !== type) return "Missing or incorrect language label";
            }
            return "";
          });
          check(!problem, `${language}${chapter}/${theme}/${width}: ${problem}`);
        }
      }
      // 捕获复制内容，不改动系统剪贴板；保留之前对真实剪贴板及拒绝权限的验证。
      await page.evaluate(() => { navigator.clipboard.writeText = async (text) => { window.manualCopiedText = text; }; });
      for (const block of await page.locator(".manual-code").all()) {
        const expected = await block.locator("pre").textContent();
        await block.locator(".manual-copy").click();
        check(await page.evaluate((text) => window.manualCopiedText === text, expected), "Copied text changed");
      }
    }
  }
  check(errors.length === 0, errors.join("\n"));
  // URL 来源不能将返回按钮变成站外跳转或再次进入手册。
  const home = await page.evaluate(() => window.ComiGoPath("/"));
  for (const source of ["", "https://example.invalid/", "javascript:alert(1)", "/manual/install"]) {
    await page.evaluate(() => sessionStorage.removeItem("manual.returnTo"));
    await page.goto(`${origin}/manual/?from=${encodeURIComponent(source)}`);
    await page.locator("[data-manual-nav] a").first().waitFor();
    check((await page.locator("[data-manual-back]").getAttribute("href")) === home, "Invalid entry was accepted");
  }
  return "PASS: search dismissal, focus, 3 languages, snippets, clipboard, anchors, drawer, 5 widths, all code blocks in light/dark, translucent overlay copy and language labels, no SSE/errors";
}
