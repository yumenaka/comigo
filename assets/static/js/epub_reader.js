// EPUB 独立阅读器：只接收浏览器文件或本站原文件接口，不复用 ZIP 漫画的状态与 WASM。
(() => {
  const root = document.getElementById('EpubReader')
  if (!root) return
  const element = (name) => document.getElementById(`Epub${name}`)
  const t = (key, values) => i18next.t(`epub_${key}`, values)
  let book, rendition, progressKey, opening = false
  const palettes = { light: ['#fffdf8', '#252525'], sepia: ['#f3e7cf', '#453a2b'], dark: ['#202329', '#e5e7eb'] }

  // 隐私模式可能禁用存储，不能因此阻止阅读。
  function recall(key) {
    try { return localStorage.getItem(key) } catch { return null }
  }

  function remember(key, value) {
    try { localStorage.setItem(key, value) } catch { element('Status').textContent = t('storage_unavailable') }
  }

  // 在章节进入 iframe 之前限制资源来源；书内脚本、表单和外部跟踪资源都不执行。
  function isolateChapter(doc) {
    doc.querySelectorAll('script, iframe, frame, object, embed, form, meta[http-equiv], link:not([rel="stylesheet"])').forEach(node => node.remove())
    const policy = doc.createElementNS('http://www.w3.org/1999/xhtml', 'meta')
    policy.setAttribute('http-equiv', 'Content-Security-Policy')
    policy.setAttribute('content', "default-src 'none'; img-src blob: data:; style-src 'unsafe-inline' blob: data:; font-src blob: data:; media-src blob: data:; form-action 'none'; frame-src 'none'")
    doc.querySelector('head').prepend(policy)
    doc.querySelectorAll('a[href]').forEach(link => {
      // 仅允许包内章节和锚点跳转，书中链接不能把本地内容带到外部网站。
      if (/^(?:[a-z][a-z\d+.-]*:|\/\/)/i.test(link.getAttribute('href').trim())) link.removeAttribute('href')
    })
  }

  // 偏好独立持久化，调整字号后仍停留在当前文本位置。
  function applyAppearance() {
    const theme = element('Theme').value
    const font = element('Font').value
    const [paper, ink] = palettes[theme]
    root.dataset.theme = theme
    if (rendition) {
      rendition.themes.default({
        body: { 'background-color': `${paper} !important`, color: `${ink} !important`, 'line-height': '1.7 !important' },
        'img, svg': { 'max-width': '100% !important', 'object-fit': 'contain' },
        a: { color: theme === 'dark' ? '#91c8ff !important' : '#245f9b !important' }
      })
      rendition.themes.fontSize(`${font}%`)
    }
    remember('epub.theme', theme)
    remember('epub.fontSize', font)
  }

  // 目录用文本节点构建，书籍标题和章节名不作为 HTML 插入。
  function fillContents(items, depth = 0) {
    for (const item of items) {
      const option = new Option(`${'　'.repeat(depth)}${item.label.trim()}`, item.href)
      element('TOC').add(option)
      if (item.subitems?.length) fillContents(item.subitems, depth + 1)
    }
  }

  // 以 CFI 保存文本位置，翻页尺寸或字号变化不会把记录变成错误的图片页码。
  function relocated(location) {
    const { start, atStart, atEnd } = location
    remember(progressKey, start.cfi)
    element('Previous').disabled = atStart
    element('Next').disabled = atEnd
    const options = Array.from(element('TOC').options)
    const match = options.find(option => option.value.split('#')[0] === start.href)
    if (match) element('TOC').value = match.value
    else element('TOC').selectedIndex = -1
    element('Position').textContent = t('position', {
      chapter: start.index + 1, chapters: book.spine.length,
      page: start.displayed.page, pages: start.displayed.total
    })
  }

  // 章节内和工具栏外的方向键都按书本方向翻页，不拦截文字选择与输入控件。
  function onKey(event) {
    if (!rendition || event.altKey || event.ctrlKey || event.metaKey || /INPUT|SELECT|TEXTAREA/.test(event.target.tagName)) return
    const rtl = book.package.metadata.direction === 'rtl'
    if (event.key === 'ArrowRight' || event.key === 'ArrowLeft') {
      event.preventDefault()
      turn((event.key === 'ArrowRight') !== rtl)
    }
  }

  async function turn(next) {
    if (!rendition || opening) return
    try { await (next ? rendition.next() : rendition.prev()) }
    catch { element('Status').textContent = t('chapter_error') }
  }

  // 每次打开都销毁上一册的 iframe 和资源 URL，选择新文件不会沿用旧章节。
  async function openBook(input, name, key) {
    if (opening) return
    opening = true
    element('File').disabled = true
    element('Status').textContent = t('loading')
    element('Welcome').hidden = true
    for (const name of ['Controls', 'Navigation', 'Viewport']) element(name).hidden = true
    book?.destroy()
    rendition = null
    element('TOC').replaceChildren()
    element('Position').textContent = ''
    element('Previous').disabled = element('Next').disabled = true
    element('Viewport').replaceChildren()
    element('Title').textContent = name
    progressKey = `epub.progress.${key}`
    try {
      const data = await input()
      book = ePub({ replacements: 'blobUrl' })
      book.spine.hooks.content.register(isolateChapter)
      await book.open(data, 'binary')
      await book.ready
      if (!book.spine.length) throw new Error('Empty EPUB spine')
      element('Title').textContent = book.package.metadata.title || name
      document.title = `${element('Title').textContent} - Comigo`
      fillContents(book.navigation.toc)
      // 部分 EPUB 没有目录，仍可按 spine 阅读全部章节。
      if (!element('TOC').options.length) fillContents(book.spine.spineItems.map((item, i) => ({ href: item.href, label: t('chapter', { number: i + 1 }) })))
      for (const name of ['Controls', 'Navigation', 'Viewport']) element(name).hidden = false
      rendition = book.renderTo(element('Viewport'), {
        width: '100%', height: '100%', spread: 'none', flow: 'paginated',
        allowScriptedContent: false, allowPopups: false
      })
      rendition.on('relocated', relocated)
      rendition.on('keydown', onKey)
      rendition.hooks.content.register(contents => {
        contents.document.documentElement.lang = book.package.metadata.language || 'en'
        contents.document.defaultView.frameElement.title = element('Title').textContent
        // 手机横向滑动翻页；纵向手势和选中文字保持浏览器默认行为。
        let touch
        contents.document.addEventListener('touchstart', event => { touch = event.touches.length === 1 ? event.touches[0] : null }, { passive: true })
        contents.document.addEventListener('touchend', event => {
          if (!touch || contents.window.getSelection().toString()) return
          const dx = event.changedTouches[0].clientX - touch.clientX
          const dy = event.changedTouches[0].clientY - touch.clientY
          if (Math.abs(dx) > 60 && Math.abs(dx) > Math.abs(dy) * 2) turn((dx < 0) !== (book.package.metadata.direction === 'rtl'))
          touch = null
        }, { passive: true })
      })
      applyAppearance()
      const saved = recall(progressKey)
      try { await rendition.display(saved || undefined) }
      catch (error) {
        if (!saved) throw error
        await rendition.display()
      }
      element('Status').textContent = ''
    } catch (error) {
      console.error('EPUB:', error)
      book?.destroy()
      book = rendition = null
      for (const name of ['Controls', 'Navigation', 'Viewport']) element(name).hidden = true
      element('Welcome').hidden = false
      element('Status').textContent = t('open_error')
    } finally {
      opening = false
      element('File').disabled = false
    }
  }

  function openFile(file) {
    if (!file) return
    if (!/\.epub$/i.test(file.name)) { element('Status').textContent = t('select_epub'); return }
    return openBook(() => file.arrayBuffer(), file.name, `${file.name}:${file.size}:${file.lastModified}`)
  }

  element('File').addEventListener('change', event => {
    openFile(event.target.files[0])
    event.target.value = ''
  })
  root.addEventListener('dragover', event => event.preventDefault())
  root.addEventListener('drop', event => { event.preventDefault(); openFile(event.dataTransfer.files[0]) })
  element('Previous').addEventListener('click', () => turn(false))
  element('Next').addEventListener('click', () => turn(true))
  element('TOC').addEventListener('change', async event => {
    try { await rendition.display(event.target.value) }
    catch { element('Status').textContent = t('chapter_error') }
  })
  element('Theme').value = palettes[recall('epub.theme')] ? recall('epub.theme') : 'light'
  const font = recall('epub.fontSize')
  if (Array.from(element('Font').options).some(option => option.value === font)) element('Font').value = font
  else element('Font').value = '100'
  for (const name of ['Font', 'Theme']) element(name).addEventListener('change', applyAppearance)
  document.addEventListener('keydown', onKey)
  applyAppearance()
  if (root.dataset.source) {
    openBook(async () => {
      const response = await fetch(root.dataset.source, { credentials: 'same-origin' })
      if (!response.ok) throw new Error(`HTTP ${response.status}`)
      return response.arrayBuffer()
    }, t('reader'), root.dataset.source)
  }
})()
