// 此文件只实现模板实际使用的 drawer、modal 和 dropdown 数据属性交互。
// 调用点集中在 templ/common/header.templ、drawer.templ 与 qrcode.templ。
const backdropClassNames = ['bg-dark-backdrop/70', 'fixed', 'inset-0']

// data-* 布尔属性只接受字符串 true；未配置时沿用调用方给出的默认值。
function getBoolAttr(element, name, defaultValue) {
    const value = element.getAttribute(name)
    if (value === null) return defaultValue
    return value === 'true'
}

// 将抽屉方向转换为固定边和显隐位移类，和模板初始的 translate-x-full 配合。
function getPlacementClasses(placement) {
    switch (placement) {
        case 'right':
            return {
                active: ['transform-none'],
                inactive: ['translate-x-full'],
                base: ['right-0', 'top-0'],
            }
        case 'top':
            return {
                active: ['transform-none'],
                inactive: ['-translate-y-full'],
                base: ['top-0', 'left-0', 'right-0'],
            }
        case 'bottom':
            return {
                active: ['transform-none'],
                inactive: ['translate-y-full'],
                base: ['bottom-0', 'left-0', 'right-0'],
            }
        case 'left':
        default:
            return {
                active: ['transform-none'],
                inactive: ['-translate-x-full'],
                base: ['left-0', 'top-0'],
            }
    }
}

function addClasses(element, classes) {
    element.classList.add(...classes)
}

function removeClasses(element, classes) {
    element.classList.remove(...classes)
}

// 控制单个抽屉的位移、遮罩、页面滚动和无障碍属性。
function createDrawerController(drawer, options) {
    let visible = false
    const placementClasses = getPlacementClasses(options.placement)

    addClasses(drawer, placementClasses.base)
    drawer.classList.add('transition-transform')
    drawer.setAttribute('aria-hidden', 'true')

    function removeBackdrop() {
        document.querySelector('[drawer-backdrop]')?.remove()
    }

    function createBackdrop() {
        if (!options.backdrop || visible) return
        const backdrop = document.createElement('div')
        backdrop.setAttribute('drawer-backdrop', '')
        addClasses(backdrop, [...backdropClassNames, 'z-30'])
        backdrop.addEventListener('click', hide)
        document.body.append(backdrop)
    }

    // 打开/关闭抽屉时只切换既有 class 与 aria 属性，避免影响抽屉内部 Alpine 状态。
    function show() {
        removeClasses(drawer, placementClasses.inactive)
        addClasses(drawer, placementClasses.active)
        drawer.setAttribute('aria-modal', 'true')
        drawer.setAttribute('role', 'dialog')
        drawer.removeAttribute('aria-hidden')
        if (!options.bodyScrolling) {
            document.body.classList.add('overflow-hidden')
        }
        createBackdrop()
        visible = true
        // 通知依赖可见状态的内容；二维码据此按需检查 URL，不监听每次翻页。
        drawer.dispatchEvent(
            new CustomEvent('comigo-drawer-shown', { bubbles: true }),
        )
    }

    function hide() {
        removeClasses(drawer, placementClasses.active)
        addClasses(drawer, placementClasses.inactive)
        drawer.setAttribute('aria-hidden', 'true')
        drawer.removeAttribute('aria-modal')
        drawer.removeAttribute('role')
        if (!options.bodyScrolling) {
            document.body.classList.remove('overflow-hidden')
        }
        if (options.backdrop) {
            removeBackdrop()
        }
        visible = false
    }

    function toggle() {
        if (visible) {
            hide()
        } else {
            show()
        }
    }

    return {
        hide,
        show,
        toggle,
        isVisible: () => visible,
    }
}

// 先由 data-drawer-target 建立控制器，再为 show/toggle/hide 触发器绑定点击行为。
function initDrawers() {
    const drawers = new Map()

    document.querySelectorAll('[data-drawer-target]').forEach((trigger) => {
        const drawerId = trigger.getAttribute('data-drawer-target')
        const drawer = document.getElementById(drawerId)
        if (!drawer || drawers.has(drawerId)) return

        drawers.set(
            drawerId,
            createDrawerController(drawer, {
                backdrop: getBoolAttr(trigger, 'data-drawer-backdrop', true),
                bodyScrolling: getBoolAttr(
                    trigger,
                    'data-drawer-body-scrolling',
                    false,
                ),
                placement:
                    trigger.getAttribute('data-drawer-placement') || 'left',
            }),
        )
    })

    document.querySelectorAll('[data-drawer-show]').forEach((trigger) => {
        trigger.addEventListener('click', () => {
            drawers.get(trigger.getAttribute('data-drawer-show'))?.show()
        })
    })

    document.querySelectorAll('[data-drawer-toggle]').forEach((trigger) => {
        trigger.addEventListener('click', () => {
            drawers.get(trigger.getAttribute('data-drawer-toggle'))?.toggle()
        })
    })

    document
        .querySelectorAll('[data-drawer-hide], [data-drawer-dismiss]')
        .forEach((trigger) => {
            trigger.addEventListener('click', () => {
                const drawerId =
                    trigger.getAttribute('data-drawer-hide') ||
                    trigger.getAttribute('data-drawer-dismiss')
                drawers.get(drawerId)?.hide()
            })
        })

    document.addEventListener('keydown', (event) => {
        if (event.key !== 'Escape') return
        drawers.forEach((drawer) => {
            if (drawer.isVisible()) drawer.hide()
        })
    })
}

// modal 本身是 flex 容器，placement 决定内容在视口中的对齐位置。
function getModalPlacementClasses(placement) {
    switch (placement) {
        case 'top-left':
            return ['justify-start', 'items-start']
        case 'top-center':
            return ['justify-center', 'items-start']
        case 'top-right':
            return ['justify-end', 'items-start']
        case 'center-left':
            return ['justify-start', 'items-center']
        case 'center-right':
            return ['justify-end', 'items-center']
        case 'bottom-left':
            return ['justify-start', 'items-end']
        case 'bottom-center':
            return ['justify-center', 'items-end']
        case 'bottom-right':
            return ['justify-end', 'items-end']
        case 'center':
        default:
            return ['justify-center', 'items-center']
    }
}

// 控制单个模态框；dynamic 遮罩允许点击空白处关闭，所有模态框都支持 ESC。
function createModalController(modal, options) {
    let visible = false
    let backdrop = null

    addClasses(modal, getModalPlacementClasses(options.placement))

    function removeBackdrop() {
        backdrop?.remove()
        backdrop = null
    }

    function createBackdrop() {
        if (backdrop) return
        backdrop = document.createElement('div')
        addClasses(backdrop, [...backdropClassNames, 'z-40'])
        document.body.append(backdrop)
    }

    function handleOutsideClick(event) {
        if (
            options.backdrop === 'dynamic' &&
            (event.target === modal || event.target === backdrop)
        ) {
            hide()
        }
    }

    function handleKeydown(event) {
        if (event.key === 'Escape') {
            hide()
        }
    }

    // Modal 保留原先的 dynamic backdrop 与 ESC 关闭能力，便携 reader 页也会复用该路径。
    function show() {
        if (visible) return
        modal.classList.remove('hidden')
        modal.classList.add('flex')
        modal.setAttribute('aria-modal', 'true')
        modal.setAttribute('role', 'dialog')
        modal.removeAttribute('aria-hidden')
        createBackdrop()
        document.body.classList.add('overflow-hidden')
        modal.addEventListener('click', handleOutsideClick, true)
        document.body.addEventListener('keydown', handleKeydown, true)
        visible = true
        // 模态框真正显示后再刷新其按需内容，关闭动作不会触发后台请求。
        modal.dispatchEvent(
            new CustomEvent('comigo-modal-shown', { bubbles: true }),
        )
    }

    function hide() {
        if (!visible && modal.classList.contains('hidden')) return
        modal.classList.add('hidden')
        modal.classList.remove('flex')
        modal.setAttribute('aria-hidden', 'true')
        modal.removeAttribute('aria-modal')
        modal.removeAttribute('role')
        removeBackdrop()
        document.body.classList.remove('overflow-hidden')
        modal.removeEventListener('click', handleOutsideClick, true)
        document.body.removeEventListener('keydown', handleKeydown, true)
        visible = false
    }

    function toggle() {
        if (visible) {
            hide()
        } else {
            show()
        }
    }

    return {
        hide,
        show,
        toggle,
        isVisible: () => visible,
    }
}

// data-modal-target 负责注册模态框，toggle/show/hide 只调用已注册的控制器。
function initModals() {
    const modals = new Map()

    document.querySelectorAll('[data-modal-target]').forEach((trigger) => {
        const modalId = trigger.getAttribute('data-modal-target')
        const modal = document.getElementById(modalId)
        if (!modal || modals.has(modalId)) return

        modals.set(
            modalId,
            createModalController(modal, {
                backdrop:
                    modal.getAttribute('data-modal-backdrop') || 'dynamic',
                placement:
                    modal.getAttribute('data-modal-placement') || 'center',
            }),
        )
    })

    document.querySelectorAll('[data-modal-toggle]').forEach((trigger) => {
        trigger.addEventListener('click', () => {
            modals.get(trigger.getAttribute('data-modal-toggle'))?.toggle()
        })
    })

    document.querySelectorAll('[data-modal-show]').forEach((trigger) => {
        trigger.addEventListener('click', () => {
            modals.get(trigger.getAttribute('data-modal-show'))?.show()
        })
    })

    document.querySelectorAll('[data-modal-hide]').forEach((trigger) => {
        trigger.addEventListener('click', () => {
            modals.get(trigger.getAttribute('data-modal-hide'))?.hide()
        })
    })
}

// 当前排序菜单固定显示在触发器下方并水平居中，同时限制在视口宽度内。
function positionDropdown(trigger, menu, offsetDistance) {
    const triggerRect = trigger.getBoundingClientRect()
    const menuRect = menu.getBoundingClientRect()
    const viewportWidth =
        window.innerWidth || document.documentElement.clientWidth || 0
    const scrollX = window.scrollX || window.pageXOffset || 0
    const scrollY = window.scrollY || window.pageYOffset || 0
    const centeredLeft =
        triggerRect.left + scrollX + triggerRect.width / 2 - menuRect.width / 2
    const maxLeft = Math.max(0, viewportWidth - menuRect.width)
    const left = Math.max(0, Math.min(centeredLeft, maxLeft))
    const top = triggerRect.bottom + scrollY + offsetDistance

    menu.style.position = 'absolute'
    menu.style.inset = '0px auto auto 0px'
    menu.style.margin = '0px'
    menu.style.transform = `translate(${Math.round(left)}px, ${Math.round(top)}px)`
    menu.setAttribute('data-popper-placement', 'bottom')
}

// 排序菜单支持模板声明的 hover 或 click 触发，并在点击菜单外部时关闭。
function createDropdownController(trigger, menu, options) {
    let visible = false
    let hideTimer = null

    function clearHideTimer() {
        if (!hideTimer) return
        window.clearTimeout(hideTimer)
        hideTimer = null
    }

    function show() {
        clearHideTimer()
        menu.classList.remove('hidden')
        menu.classList.add('block')
        menu.removeAttribute('aria-hidden')
        positionDropdown(trigger, menu, options.offsetDistance)
        visible = true
        document.body.addEventListener('click', handleOutsideClick, true)
    }

    function hide() {
        clearHideTimer()
        menu.classList.remove('block')
        menu.classList.add('hidden')
        menu.setAttribute('aria-hidden', 'true')
        visible = false
        document.body.removeEventListener('click', handleOutsideClick, true)
    }

    function toggle() {
        if (visible) {
            hide()
        } else {
            show()
        }
    }

    function handleOutsideClick(event) {
        if (trigger.contains(event.target) || menu.contains(event.target))
            return
        hide()
    }

    function scheduleHide() {
        clearHideTimer()
        hideTimer = window.setTimeout(() => {
            if (!menu.matches(':hover') && !trigger.matches(':hover')) {
                hide()
            }
        }, options.delay)
    }

    if (options.triggerType === 'hover') {
        trigger.addEventListener('mouseenter', () => {
            window.setTimeout(show, options.delay)
        })
        trigger.addEventListener('click', toggle)
        trigger.addEventListener('mouseleave', scheduleHide)
        menu.addEventListener('mouseenter', show)
        menu.addEventListener('mouseleave', scheduleHide)
    } else if (options.triggerType === 'click') {
        trigger.addEventListener('click', toggle)
    }
}

// 读取 header.templ 中的 data-dropdown-* 参数并初始化排序菜单。
function initDropdowns() {
    document.querySelectorAll('[data-dropdown-toggle]').forEach((trigger) => {
        const dropdownId = trigger.getAttribute('data-dropdown-toggle')
        const menu = document.getElementById(dropdownId)
        if (!menu) return

        createDropdownController(trigger, menu, {
            delay: Number.parseInt(
                trigger.getAttribute('data-dropdown-delay') || '300',
                10,
            ),
            offsetDistance: Number.parseInt(
                trigger.getAttribute('data-dropdown-offset-distance') || '10',
                10,
            ),
            triggerType:
                trigger.getAttribute('data-dropdown-trigger') || 'click',
        })
    })
}

// 主入口只初始化当前模板真实使用的三类控件。
export function initComigoUIControls() {
    initDrawers()
    initModals()
    initDropdowns()
}

// 主包脚本可能先于页面节点执行，等待 DOM 就绪后再扫描 data-* 属性。
document.addEventListener('DOMContentLoaded', initComigoUIControls)
