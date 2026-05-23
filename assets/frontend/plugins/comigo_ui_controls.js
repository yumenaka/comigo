// 此文件提供当前项目实际用到的轻量 UI 交互，避免为了少量组件引入整套前端依赖。
const backdropClassNames = ['bg-dark-backdrop/70', 'fixed', 'inset-0']

function getBoolAttr(element, name, defaultValue) {
    const value = element.getAttribute(name)
    if (value === null) return defaultValue
    return value === 'true'
}

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
                placement: trigger.getAttribute('data-drawer-placement') || 'left',
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

function initModals() {
    const modals = new Map()

    document.querySelectorAll('[data-modal-target]').forEach((trigger) => {
        const modalId = trigger.getAttribute('data-modal-target')
        const modal = document.getElementById(modalId)
        if (!modal || modals.has(modalId)) return

        modals.set(
            modalId,
            createModalController(modal, {
                backdrop: modal.getAttribute('data-modal-backdrop') || 'dynamic',
                placement: modal.getAttribute('data-modal-placement') || 'center',
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
        if (trigger.contains(event.target) || menu.contains(event.target)) return
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
            triggerType: trigger.getAttribute('data-dropdown-trigger') || 'click',
        })
    })
}

export function initComigoUIControls() {
    initDrawers()
    initModals()
    initDropdowns()
}

document.addEventListener('DOMContentLoaded', initComigoUIControls)
