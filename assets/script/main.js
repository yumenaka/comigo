(function () {

function $parcel$interopDefault(a) {
  return a && a.__esModule ? a.default : a;
}
//此文件需要编译，编译指令请参考 package.json
// import 'htmx.org'
var $4c6a2f83122c31a3$var$Events = /** @class */ function() {
    function Events(eventType, eventFunctions) {
        if (eventFunctions === void 0) eventFunctions = [];
        this._eventType = eventType;
        this._eventFunctions = eventFunctions;
    }
    Events.prototype.init = function() {
        var _this = this;
        this._eventFunctions.forEach(function(eventFunction) {
            if (typeof window !== 'undefined') window.addEventListener(_this._eventType, eventFunction);
        });
    };
    return Events;
}();
var $4c6a2f83122c31a3$export$2e2bcd8739ae039 = $4c6a2f83122c31a3$var$Events;


var $43290e7e11b6c4c9$var$Instances = /** @class */ function() {
    function Instances() {
        this._instances = {
            Accordion: {},
            Carousel: {},
            Collapse: {},
            Dial: {},
            Dismiss: {},
            Drawer: {},
            Dropdown: {},
            Modal: {},
            Popover: {},
            Tabs: {},
            Tooltip: {},
            InputCounter: {},
            CopyClipboard: {},
            Datepicker: {}
        };
    }
    Instances.prototype.addInstance = function(component, instance, id, override) {
        if (override === void 0) override = false;
        if (!this._instances[component]) {
            console.warn("Flowbite: Component ".concat(component, " does not exist."));
            return false;
        }
        if (this._instances[component][id] && !override) {
            console.warn("Flowbite: Instance with ID ".concat(id, " already exists."));
            return;
        }
        if (override && this._instances[component][id]) this._instances[component][id].destroyAndRemoveInstance();
        this._instances[component][id ? id : this._generateRandomId()] = instance;
    };
    Instances.prototype.getAllInstances = function() {
        return this._instances;
    };
    Instances.prototype.getInstances = function(component) {
        if (!this._instances[component]) {
            console.warn("Flowbite: Component ".concat(component, " does not exist."));
            return false;
        }
        return this._instances[component];
    };
    Instances.prototype.getInstance = function(component, id) {
        if (!this._componentAndInstanceCheck(component, id)) return;
        if (!this._instances[component][id]) {
            console.warn("Flowbite: Instance with ID ".concat(id, " does not exist."));
            return;
        }
        return this._instances[component][id];
    };
    Instances.prototype.destroyAndRemoveInstance = function(component, id) {
        if (!this._componentAndInstanceCheck(component, id)) return;
        this.destroyInstanceObject(component, id);
        this.removeInstance(component, id);
    };
    Instances.prototype.removeInstance = function(component, id) {
        if (!this._componentAndInstanceCheck(component, id)) return;
        delete this._instances[component][id];
    };
    Instances.prototype.destroyInstanceObject = function(component, id) {
        if (!this._componentAndInstanceCheck(component, id)) return;
        this._instances[component][id].destroy();
    };
    Instances.prototype.instanceExists = function(component, id) {
        if (!this._instances[component]) return false;
        if (!this._instances[component][id]) return false;
        return true;
    };
    Instances.prototype._generateRandomId = function() {
        return Math.random().toString(36).substr(2, 9);
    };
    Instances.prototype._componentAndInstanceCheck = function(component, id) {
        if (!this._instances[component]) {
            console.warn("Flowbite: Component ".concat(component, " does not exist."));
            return false;
        }
        if (!this._instances[component][id]) {
            console.warn("Flowbite: Instance with ID ".concat(id, " does not exist."));
            return false;
        }
        return true;
    };
    return Instances;
}();
var $43290e7e11b6c4c9$var$instances = new $43290e7e11b6c4c9$var$Instances();
var $43290e7e11b6c4c9$export$2e2bcd8739ae039 = $43290e7e11b6c4c9$var$instances;
if (typeof window !== 'undefined') window.FlowbiteInstances = $43290e7e11b6c4c9$var$instances;


var $946dcfd8fe0754f4$var$__assign = undefined && undefined.__assign || function() {
    $946dcfd8fe0754f4$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $946dcfd8fe0754f4$var$__assign.apply(this, arguments);
};
var $946dcfd8fe0754f4$var$Default = {
    alwaysOpen: false,
    activeClasses: 'bg-neutral-secondary-medium text-heading',
    inactiveClasses: 'bg-neutral-primary text-body',
    onOpen: function() {},
    onClose: function() {},
    onToggle: function() {}
};
var $946dcfd8fe0754f4$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $946dcfd8fe0754f4$var$Accordion = /** @class */ function() {
    function Accordion(accordionEl, items, options, instanceOptions) {
        if (accordionEl === void 0) accordionEl = null;
        if (items === void 0) items = [];
        if (options === void 0) options = $946dcfd8fe0754f4$var$Default;
        if (instanceOptions === void 0) instanceOptions = $946dcfd8fe0754f4$var$DefaultInstanceOptions;
        this._instanceId = instanceOptions.id ? instanceOptions.id : accordionEl.id;
        this._accordionEl = accordionEl;
        this._items = items;
        this._options = $946dcfd8fe0754f4$var$__assign($946dcfd8fe0754f4$var$__assign({}, $946dcfd8fe0754f4$var$Default), options);
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('Accordion', this, this._instanceId, instanceOptions.override);
    }
    Accordion.prototype.init = function() {
        var _this = this;
        if (this._items.length && !this._initialized) {
            // show accordion item based on click
            this._items.forEach(function(item) {
                if (item.active) _this.open(item.id);
                var clickHandler = function() {
                    _this.toggle(item.id);
                };
                item.triggerEl.addEventListener('click', clickHandler);
                // Store the clickHandler in a property of the item for removal later
                item.clickHandler = clickHandler;
            });
            this._initialized = true;
        }
    };
    Accordion.prototype.destroy = function() {
        if (this._items.length && this._initialized) {
            this._items.forEach(function(item) {
                item.triggerEl.removeEventListener('click', item.clickHandler);
                // Clean up by deleting the clickHandler property from the item
                delete item.clickHandler;
            });
            this._initialized = false;
        }
    };
    Accordion.prototype.removeInstance = function() {
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('Accordion', this._instanceId);
    };
    Accordion.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    Accordion.prototype.getItem = function(id) {
        return this._items.filter(function(item) {
            return item.id === id;
        })[0];
    };
    Accordion.prototype.open = function(id) {
        var _a, _b;
        var _this = this;
        var item = this.getItem(id);
        // don't hide other accordions if always open
        if (!this._options.alwaysOpen) this._items.map(function(i) {
            var _a, _b;
            if (i !== item) {
                (_a = i.triggerEl.classList).remove.apply(_a, _this._options.activeClasses.split(' '));
                (_b = i.triggerEl.classList).add.apply(_b, _this._options.inactiveClasses.split(' '));
                i.targetEl.classList.add('hidden');
                i.triggerEl.setAttribute('aria-expanded', 'false');
                i.active = false;
                // rotate icon if set
                if (i.iconEl) i.iconEl.classList.add('rotate-180');
            }
        });
        // show active item
        (_a = item.triggerEl.classList).add.apply(_a, this._options.activeClasses.split(' '));
        (_b = item.triggerEl.classList).remove.apply(_b, this._options.inactiveClasses.split(' '));
        item.triggerEl.setAttribute('aria-expanded', 'true');
        item.targetEl.classList.remove('hidden');
        item.active = true;
        // rotate icon if set
        if (item.iconEl) item.iconEl.classList.remove('rotate-180');
        // callback function
        this._options.onOpen(this, item);
    };
    Accordion.prototype.toggle = function(id) {
        var item = this.getItem(id);
        if (item.active) this.close(id);
        else this.open(id);
        // callback function
        this._options.onToggle(this, item);
    };
    Accordion.prototype.close = function(id) {
        var _a, _b;
        var item = this.getItem(id);
        (_a = item.triggerEl.classList).remove.apply(_a, this._options.activeClasses.split(' '));
        (_b = item.triggerEl.classList).add.apply(_b, this._options.inactiveClasses.split(' '));
        item.targetEl.classList.add('hidden');
        item.triggerEl.setAttribute('aria-expanded', 'false');
        item.active = false;
        // rotate icon if set
        if (item.iconEl) item.iconEl.classList.add('rotate-180');
        // callback function
        this._options.onClose(this, item);
    };
    Accordion.prototype.updateOnOpen = function(callback) {
        this._options.onOpen = callback;
    };
    Accordion.prototype.updateOnClose = function(callback) {
        this._options.onClose = callback;
    };
    Accordion.prototype.updateOnToggle = function(callback) {
        this._options.onToggle = callback;
    };
    return Accordion;
}();
function $946dcfd8fe0754f4$export$226c1dc98323ee4d() {
    document.querySelectorAll('[data-accordion]').forEach(function($accordionEl) {
        var alwaysOpen = $accordionEl.getAttribute('data-accordion');
        var activeClasses = $accordionEl.getAttribute('data-active-classes');
        var inactiveClasses = $accordionEl.getAttribute('data-inactive-classes');
        var items = [];
        $accordionEl.querySelectorAll('[data-accordion-target]').forEach(function($triggerEl) {
            // Consider only items that directly belong to $accordionEl
            // (to make nested accordions work).
            if ($triggerEl.closest('[data-accordion]') === $accordionEl) {
                var item = {
                    id: $triggerEl.getAttribute('data-accordion-target'),
                    triggerEl: $triggerEl,
                    targetEl: document.querySelector($triggerEl.getAttribute('data-accordion-target')),
                    iconEl: $triggerEl.querySelector('[data-accordion-icon]'),
                    active: $triggerEl.getAttribute('aria-expanded') === 'true' ? true : false
                };
                items.push(item);
            }
        });
        new $946dcfd8fe0754f4$var$Accordion($accordionEl, items, {
            alwaysOpen: alwaysOpen === 'open' ? true : false,
            activeClasses: activeClasses ? activeClasses : $946dcfd8fe0754f4$var$Default.activeClasses,
            inactiveClasses: inactiveClasses ? inactiveClasses : $946dcfd8fe0754f4$var$Default.inactiveClasses
        });
    });
}
if (typeof window !== 'undefined') {
    window.Accordion = $946dcfd8fe0754f4$var$Accordion;
    window.initAccordions = $946dcfd8fe0754f4$export$226c1dc98323ee4d;
}
var $946dcfd8fe0754f4$export$2e2bcd8739ae039 = $946dcfd8fe0754f4$var$Accordion;



var $16ba7ad25d51d4f8$var$__assign = undefined && undefined.__assign || function() {
    $16ba7ad25d51d4f8$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $16ba7ad25d51d4f8$var$__assign.apply(this, arguments);
};
var $16ba7ad25d51d4f8$var$Default = {
    onCollapse: function() {},
    onExpand: function() {},
    onToggle: function() {}
};
var $16ba7ad25d51d4f8$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $16ba7ad25d51d4f8$var$Collapse = /** @class */ function() {
    function Collapse(targetEl, triggerEl, options, instanceOptions) {
        if (targetEl === void 0) targetEl = null;
        if (triggerEl === void 0) triggerEl = null;
        if (options === void 0) options = $16ba7ad25d51d4f8$var$Default;
        if (instanceOptions === void 0) instanceOptions = $16ba7ad25d51d4f8$var$DefaultInstanceOptions;
        this._instanceId = instanceOptions.id ? instanceOptions.id : targetEl.id;
        this._targetEl = targetEl;
        this._triggerEl = triggerEl;
        this._options = $16ba7ad25d51d4f8$var$__assign($16ba7ad25d51d4f8$var$__assign({}, $16ba7ad25d51d4f8$var$Default), options);
        this._visible = false;
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('Collapse', this, this._instanceId, instanceOptions.override);
    }
    Collapse.prototype.init = function() {
        var _this = this;
        if (this._triggerEl && this._targetEl && !this._initialized) {
            if (this._triggerEl.hasAttribute('aria-expanded')) this._visible = this._triggerEl.getAttribute('aria-expanded') === 'true';
            else // fix until v2 not to break previous single collapses which became dismiss
            this._visible = !this._targetEl.classList.contains('hidden');
            this._clickHandler = function() {
                _this.toggle();
            };
            this._triggerEl.addEventListener('click', this._clickHandler);
            this._initialized = true;
        }
    };
    Collapse.prototype.destroy = function() {
        if (this._triggerEl && this._initialized) {
            this._triggerEl.removeEventListener('click', this._clickHandler);
            this._initialized = false;
        }
    };
    Collapse.prototype.removeInstance = function() {
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('Collapse', this._instanceId);
    };
    Collapse.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    Collapse.prototype.collapse = function() {
        this._targetEl.classList.add('hidden');
        if (this._triggerEl) this._triggerEl.setAttribute('aria-expanded', 'false');
        this._visible = false;
        // callback function
        this._options.onCollapse(this);
    };
    Collapse.prototype.expand = function() {
        this._targetEl.classList.remove('hidden');
        if (this._triggerEl) this._triggerEl.setAttribute('aria-expanded', 'true');
        this._visible = true;
        // callback function
        this._options.onExpand(this);
    };
    Collapse.prototype.toggle = function() {
        if (this._visible) this.collapse();
        else this.expand();
        // callback function
        this._options.onToggle(this);
    };
    Collapse.prototype.updateOnCollapse = function(callback) {
        this._options.onCollapse = callback;
    };
    Collapse.prototype.updateOnExpand = function(callback) {
        this._options.onExpand = callback;
    };
    Collapse.prototype.updateOnToggle = function(callback) {
        this._options.onToggle = callback;
    };
    return Collapse;
}();
function $16ba7ad25d51d4f8$export$355ba5a528b4009a() {
    document.querySelectorAll('[data-collapse-toggle]').forEach(function($triggerEl) {
        var targetId = $triggerEl.getAttribute('data-collapse-toggle');
        var $targetEl = document.getElementById(targetId);
        // check if the target element exists
        if ($targetEl) {
            if (!(0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).instanceExists('Collapse', $targetEl.getAttribute('id'))) new $16ba7ad25d51d4f8$var$Collapse($targetEl, $triggerEl);
            else // if instance exists already for the same target element then create a new one with a different trigger element
            new $16ba7ad25d51d4f8$var$Collapse($targetEl, $triggerEl, {}, {
                id: $targetEl.getAttribute('id') + '_' + (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039)._generateRandomId()
            });
        } else console.error("The target element with id \"".concat(targetId, "\" does not exist. Please check the data-collapse-toggle attribute."));
    });
}
if (typeof window !== 'undefined') {
    window.Collapse = $16ba7ad25d51d4f8$var$Collapse;
    window.initCollapses = $16ba7ad25d51d4f8$export$355ba5a528b4009a;
}
var $16ba7ad25d51d4f8$export$2e2bcd8739ae039 = $16ba7ad25d51d4f8$var$Collapse;



var $6a8fcd5cc87cb99e$var$__assign = undefined && undefined.__assign || function() {
    $6a8fcd5cc87cb99e$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $6a8fcd5cc87cb99e$var$__assign.apply(this, arguments);
};
var $6a8fcd5cc87cb99e$var$Default = {
    defaultPosition: 0,
    indicators: {
        items: [],
        activeClasses: 'bg-white dark:bg-gray-800',
        inactiveClasses: 'bg-white/50 dark:bg-gray-800/50 hover:bg-white dark:hover:bg-gray-800'
    },
    interval: 3000,
    onNext: function() {},
    onPrev: function() {},
    onChange: function() {}
};
var $6a8fcd5cc87cb99e$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $6a8fcd5cc87cb99e$var$Carousel = /** @class */ function() {
    function Carousel(carouselEl, items, options, instanceOptions) {
        if (carouselEl === void 0) carouselEl = null;
        if (items === void 0) items = [];
        if (options === void 0) options = $6a8fcd5cc87cb99e$var$Default;
        if (instanceOptions === void 0) instanceOptions = $6a8fcd5cc87cb99e$var$DefaultInstanceOptions;
        this._instanceId = instanceOptions.id ? instanceOptions.id : carouselEl.id;
        this._carouselEl = carouselEl;
        this._items = items;
        this._options = $6a8fcd5cc87cb99e$var$__assign($6a8fcd5cc87cb99e$var$__assign($6a8fcd5cc87cb99e$var$__assign({}, $6a8fcd5cc87cb99e$var$Default), options), {
            indicators: $6a8fcd5cc87cb99e$var$__assign($6a8fcd5cc87cb99e$var$__assign({}, $6a8fcd5cc87cb99e$var$Default.indicators), options.indicators)
        });
        this._activeItem = this.getItem(this._options.defaultPosition);
        this._indicators = this._options.indicators.items;
        this._intervalDuration = this._options.interval;
        this._intervalInstance = null;
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('Carousel', this, this._instanceId, instanceOptions.override);
    }
    /**
     * initialize carousel and items based on active one
     */ Carousel.prototype.init = function() {
        var _this = this;
        if (this._items.length && !this._initialized) {
            this._items.map(function(item) {
                item.el.classList.add('absolute', 'inset-0', 'transition-transform', 'transform');
            });
            // if no active item is set then first position is default
            if (this.getActiveItem()) this.slideTo(this.getActiveItem().position);
            else this.slideTo(0);
            this._indicators.map(function(indicator, position) {
                indicator.el.addEventListener('click', function() {
                    _this.slideTo(position);
                });
            });
            this._initialized = true;
        }
    };
    Carousel.prototype.destroy = function() {
        if (this._initialized) this._initialized = false;
    };
    Carousel.prototype.removeInstance = function() {
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('Carousel', this._instanceId);
    };
    Carousel.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    Carousel.prototype.getItem = function(position) {
        return this._items[position];
    };
    /**
     * Slide to the element based on id
     * @param {*} position
     */ Carousel.prototype.slideTo = function(position) {
        var nextItem = this._items[position];
        var rotationItems = {
            left: nextItem.position === 0 ? this._items[this._items.length - 1] : this._items[nextItem.position - 1],
            middle: nextItem,
            right: nextItem.position === this._items.length - 1 ? this._items[0] : this._items[nextItem.position + 1]
        };
        this._rotate(rotationItems);
        this._setActiveItem(nextItem);
        if (this._intervalInstance) {
            this.pause();
            this.cycle();
        }
        this._options.onChange(this);
    };
    /**
     * Based on the currently active item it will go to the next position
     */ Carousel.prototype.next = function() {
        var activeItem = this.getActiveItem();
        var nextItem = null;
        // check if last item
        if (activeItem.position === this._items.length - 1) nextItem = this._items[0];
        else nextItem = this._items[activeItem.position + 1];
        this.slideTo(nextItem.position);
        // callback function
        this._options.onNext(this);
    };
    /**
     * Based on the currently active item it will go to the previous position
     */ Carousel.prototype.prev = function() {
        var activeItem = this.getActiveItem();
        var prevItem = null;
        // check if first item
        if (activeItem.position === 0) prevItem = this._items[this._items.length - 1];
        else prevItem = this._items[activeItem.position - 1];
        this.slideTo(prevItem.position);
        // callback function
        this._options.onPrev(this);
    };
    /**
     * This method applies the transform classes based on the left, middle, and right rotation carousel items
     * @param {*} rotationItems
     */ Carousel.prototype._rotate = function(rotationItems) {
        // reset
        this._items.map(function(item) {
            item.el.classList.add('hidden');
        });
        // Handling the case when there is only one item
        if (this._items.length === 1) {
            rotationItems.middle.el.classList.remove('-translate-x-full', 'translate-x-full', 'translate-x-0', 'hidden', 'z-10');
            rotationItems.middle.el.classList.add('translate-x-0', 'z-20');
            return;
        }
        // left item (previously active)
        rotationItems.left.el.classList.remove('-translate-x-full', 'translate-x-full', 'translate-x-0', 'hidden', 'z-20');
        rotationItems.left.el.classList.add('-translate-x-full', 'z-10');
        // currently active item
        rotationItems.middle.el.classList.remove('-translate-x-full', 'translate-x-full', 'translate-x-0', 'hidden', 'z-10');
        rotationItems.middle.el.classList.add('translate-x-0', 'z-30');
        // right item (upcoming active)
        rotationItems.right.el.classList.remove('-translate-x-full', 'translate-x-full', 'translate-x-0', 'hidden', 'z-30');
        rotationItems.right.el.classList.add('translate-x-full', 'z-20');
    };
    /**
     * Set an interval to cycle through the carousel items
     */ Carousel.prototype.cycle = function() {
        var _this = this;
        if (typeof window !== 'undefined') this._intervalInstance = window.setInterval(function() {
            _this.next();
        }, this._intervalDuration);
    };
    /**
     * Clears the cycling interval
     */ Carousel.prototype.pause = function() {
        clearInterval(this._intervalInstance);
    };
    /**
     * Get the currently active item
     */ Carousel.prototype.getActiveItem = function() {
        return this._activeItem;
    };
    /**
     * Set the currently active item and data attribute
     * @param {*} position
     */ Carousel.prototype._setActiveItem = function(item) {
        var _a, _b;
        var _this = this;
        this._activeItem = item;
        var position = item.position;
        // update the indicators if available
        if (this._indicators.length) {
            this._indicators.map(function(indicator) {
                var _a, _b;
                indicator.el.setAttribute('aria-current', 'false');
                (_a = indicator.el.classList).remove.apply(_a, _this._options.indicators.activeClasses.split(' '));
                (_b = indicator.el.classList).add.apply(_b, _this._options.indicators.inactiveClasses.split(' '));
            });
            (_a = this._indicators[position].el.classList).add.apply(_a, this._options.indicators.activeClasses.split(' '));
            (_b = this._indicators[position].el.classList).remove.apply(_b, this._options.indicators.inactiveClasses.split(' '));
            this._indicators[position].el.setAttribute('aria-current', 'true');
        }
    };
    Carousel.prototype.updateOnNext = function(callback) {
        this._options.onNext = callback;
    };
    Carousel.prototype.updateOnPrev = function(callback) {
        this._options.onPrev = callback;
    };
    Carousel.prototype.updateOnChange = function(callback) {
        this._options.onChange = callback;
    };
    return Carousel;
}();
function $6a8fcd5cc87cb99e$export$3ab77386b16b9e58() {
    document.querySelectorAll('[data-carousel]').forEach(function($carouselEl) {
        var interval = $carouselEl.getAttribute('data-carousel-interval');
        var slide = $carouselEl.getAttribute('data-carousel') === 'slide' ? true : false;
        var items = [];
        var defaultPosition = 0;
        if ($carouselEl.querySelectorAll('[data-carousel-item]').length) Array.from($carouselEl.querySelectorAll('[data-carousel-item]')).map(function($carouselItemEl, position) {
            items.push({
                position: position,
                el: $carouselItemEl
            });
            if ($carouselItemEl.getAttribute('data-carousel-item') === 'active') defaultPosition = position;
        });
        var indicators = [];
        if ($carouselEl.querySelectorAll('[data-carousel-slide-to]').length) Array.from($carouselEl.querySelectorAll('[data-carousel-slide-to]')).map(function($indicatorEl) {
            indicators.push({
                position: parseInt($indicatorEl.getAttribute('data-carousel-slide-to')),
                el: $indicatorEl
            });
        });
        var carousel = new $6a8fcd5cc87cb99e$var$Carousel($carouselEl, items, {
            defaultPosition: defaultPosition,
            indicators: {
                items: indicators
            },
            interval: interval ? interval : $6a8fcd5cc87cb99e$var$Default.interval
        });
        if (slide) carousel.cycle();
        // check for controls
        var carouselNextEl = $carouselEl.querySelector('[data-carousel-next]');
        var carouselPrevEl = $carouselEl.querySelector('[data-carousel-prev]');
        if (carouselNextEl) carouselNextEl.addEventListener('click', function() {
            carousel.next();
        });
        if (carouselPrevEl) carouselPrevEl.addEventListener('click', function() {
            carousel.prev();
        });
    });
}
if (typeof window !== 'undefined') {
    window.Carousel = $6a8fcd5cc87cb99e$var$Carousel;
    window.initCarousels = $6a8fcd5cc87cb99e$export$3ab77386b16b9e58;
}
var $6a8fcd5cc87cb99e$export$2e2bcd8739ae039 = $6a8fcd5cc87cb99e$var$Carousel;



var $fa24e680ac3fba2d$var$__assign = undefined && undefined.__assign || function() {
    $fa24e680ac3fba2d$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $fa24e680ac3fba2d$var$__assign.apply(this, arguments);
};
var $fa24e680ac3fba2d$var$Default = {
    transition: 'transition-opacity',
    duration: 300,
    timing: 'ease-out',
    onHide: function() {}
};
var $fa24e680ac3fba2d$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $fa24e680ac3fba2d$var$Dismiss = /** @class */ function() {
    function Dismiss(targetEl, triggerEl, options, instanceOptions) {
        if (targetEl === void 0) targetEl = null;
        if (triggerEl === void 0) triggerEl = null;
        if (options === void 0) options = $fa24e680ac3fba2d$var$Default;
        if (instanceOptions === void 0) instanceOptions = $fa24e680ac3fba2d$var$DefaultInstanceOptions;
        this._instanceId = instanceOptions.id ? instanceOptions.id : targetEl.id;
        this._targetEl = targetEl;
        this._triggerEl = triggerEl;
        this._options = $fa24e680ac3fba2d$var$__assign($fa24e680ac3fba2d$var$__assign({}, $fa24e680ac3fba2d$var$Default), options);
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('Dismiss', this, this._instanceId, instanceOptions.override);
    }
    Dismiss.prototype.init = function() {
        var _this = this;
        if (this._triggerEl && this._targetEl && !this._initialized) {
            this._clickHandler = function() {
                _this.hide();
            };
            this._triggerEl.addEventListener('click', this._clickHandler);
            this._initialized = true;
        }
    };
    Dismiss.prototype.destroy = function() {
        if (this._triggerEl && this._initialized) {
            this._triggerEl.removeEventListener('click', this._clickHandler);
            this._initialized = false;
        }
    };
    Dismiss.prototype.removeInstance = function() {
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('Dismiss', this._instanceId);
    };
    Dismiss.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    Dismiss.prototype.hide = function() {
        var _this = this;
        this._targetEl.classList.add(this._options.transition, "duration-".concat(this._options.duration), this._options.timing, 'opacity-0');
        setTimeout(function() {
            _this._targetEl.classList.add('hidden');
        }, this._options.duration);
        // callback function
        this._options.onHide(this, this._targetEl);
    };
    Dismiss.prototype.updateOnHide = function(callback) {
        this._options.onHide = callback;
    };
    return Dismiss;
}();
function $fa24e680ac3fba2d$export$69c84787a503daef() {
    document.querySelectorAll('[data-dismiss-target]').forEach(function($triggerEl) {
        var targetId = $triggerEl.getAttribute('data-dismiss-target');
        var $dismissEl = document.querySelector(targetId);
        if ($dismissEl) new $fa24e680ac3fba2d$var$Dismiss($dismissEl, $triggerEl);
        else console.error("The dismiss element with id \"".concat(targetId, "\" does not exist. Please check the data-dismiss-target attribute."));
    });
}
if (typeof window !== 'undefined') {
    window.Dismiss = $fa24e680ac3fba2d$var$Dismiss;
    window.initDismisses = $fa24e680ac3fba2d$export$69c84787a503daef;
}
var $fa24e680ac3fba2d$export$2e2bcd8739ae039 = $fa24e680ac3fba2d$var$Dismiss;


function $e992165d27cf1e50$export$2e2bcd8739ae039(node) {
    if (node == null) return window;
    if (node.toString() !== '[object Window]') {
        var ownerDocument = node.ownerDocument;
        return ownerDocument ? ownerDocument.defaultView || window : window;
    }
    return node;
}


function $b83e73df1a45d10b$export$45a5e7f76e0caa8d(node) {
    var OwnElement = (0, $e992165d27cf1e50$export$2e2bcd8739ae039)(node).Element;
    return node instanceof OwnElement || node instanceof Element;
}
function $b83e73df1a45d10b$export$1b3bfaa9684536aa(node) {
    var OwnElement = (0, $e992165d27cf1e50$export$2e2bcd8739ae039)(node).HTMLElement;
    return node instanceof OwnElement || node instanceof HTMLElement;
}
function $b83e73df1a45d10b$export$af51f0f06c0f328a(node) {
    // IE 11 has no ShadowRoot
    if (typeof ShadowRoot === 'undefined') return false;
    var OwnElement = (0, $e992165d27cf1e50$export$2e2bcd8739ae039)(node).ShadowRoot;
    return node instanceof OwnElement || node instanceof ShadowRoot;
}


var $9c234cc6f4ed3b92$export$8960430cfd85939f = Math.max;
var $9c234cc6f4ed3b92$export$96ec731ed4dcb222 = Math.min;
var $9c234cc6f4ed3b92$export$2077e0241d6afd3c = Math.round;



function $5e056ec93b493b76$export$2e2bcd8739ae039() {
    var uaData = navigator.userAgentData;
    if (uaData != null && uaData.brands && Array.isArray(uaData.brands)) return uaData.brands.map(function(item) {
        return item.brand + "/" + item.version;
    }).join(' ');
    return navigator.userAgent;
}


function $ab0cb3910b57b210$export$2e2bcd8739ae039() {
    return !/^((?!chrome|android).)*safari/i.test((0, $5e056ec93b493b76$export$2e2bcd8739ae039)());
}


function $d5b6692e6e485772$export$2e2bcd8739ae039(element, includeScale, isFixedStrategy) {
    if (includeScale === void 0) includeScale = false;
    if (isFixedStrategy === void 0) isFixedStrategy = false;
    var clientRect = element.getBoundingClientRect();
    var scaleX = 1;
    var scaleY = 1;
    if (includeScale && (0, $b83e73df1a45d10b$export$1b3bfaa9684536aa)(element)) {
        scaleX = element.offsetWidth > 0 ? (0, $9c234cc6f4ed3b92$export$2077e0241d6afd3c)(clientRect.width) / element.offsetWidth || 1 : 1;
        scaleY = element.offsetHeight > 0 ? (0, $9c234cc6f4ed3b92$export$2077e0241d6afd3c)(clientRect.height) / element.offsetHeight || 1 : 1;
    }
    var _ref = (0, $b83e73df1a45d10b$export$45a5e7f76e0caa8d)(element) ? (0, $e992165d27cf1e50$export$2e2bcd8739ae039)(element) : window, visualViewport = _ref.visualViewport;
    var addVisualOffsets = !(0, $ab0cb3910b57b210$export$2e2bcd8739ae039)() && isFixedStrategy;
    var x = (clientRect.left + (addVisualOffsets && visualViewport ? visualViewport.offsetLeft : 0)) / scaleX;
    var y = (clientRect.top + (addVisualOffsets && visualViewport ? visualViewport.offsetTop : 0)) / scaleY;
    var width = clientRect.width / scaleX;
    var height = clientRect.height / scaleY;
    return {
        width: width,
        height: height,
        top: y,
        right: x + width,
        bottom: y + height,
        left: x,
        x: x,
        y: y
    };
}



function $2270e2f70c3d3867$export$2e2bcd8739ae039(node) {
    var win = (0, $e992165d27cf1e50$export$2e2bcd8739ae039)(node);
    var scrollLeft = win.pageXOffset;
    var scrollTop = win.pageYOffset;
    return {
        scrollLeft: scrollLeft,
        scrollTop: scrollTop
    };
}




function $7bd577cc58a332b0$export$2e2bcd8739ae039(element) {
    return {
        scrollLeft: element.scrollLeft,
        scrollTop: element.scrollTop
    };
}


function $1416989062dadf75$export$2e2bcd8739ae039(node) {
    if (node === (0, $e992165d27cf1e50$export$2e2bcd8739ae039)(node) || !(0, $b83e73df1a45d10b$export$1b3bfaa9684536aa)(node)) return (0, $2270e2f70c3d3867$export$2e2bcd8739ae039)(node);
    else return (0, $7bd577cc58a332b0$export$2e2bcd8739ae039)(node);
}


function $560fb2f7f26d74c6$export$2e2bcd8739ae039(element) {
    return element ? (element.nodeName || '').toLowerCase() : null;
}





function $d56bbea5dd2fe00d$export$2e2bcd8739ae039(element) {
    // $FlowFixMe[incompatible-return]: assume body is always available
    return (((0, $b83e73df1a45d10b$export$45a5e7f76e0caa8d)(element) ? element.ownerDocument : element.document) || window.document).documentElement;
}



function $c4822cfd8628a2db$export$2e2bcd8739ae039(element) {
    // If <html> has a CSS width greater than the viewport, then this will be
    // incorrect for RTL.
    // Popper 1 is broken in this case and never had a bug report so let's assume
    // it's not an issue. I don't think anyone ever specifies width on <html>
    // anyway.
    // Browsers where the left scrollbar doesn't cause an issue report `0` for
    // this (e.g. Edge 2019, IE11, Safari)
    return (0, $d5b6692e6e485772$export$2e2bcd8739ae039)((0, $d56bbea5dd2fe00d$export$2e2bcd8739ae039)(element)).left + (0, $2270e2f70c3d3867$export$2e2bcd8739ae039)(element).scrollLeft;
}




function $7d74e8e12f28eb07$export$2e2bcd8739ae039(element) {
    return (0, $e992165d27cf1e50$export$2e2bcd8739ae039)(element).getComputedStyle(element);
}


function $3e03f71b3bbb98eb$export$2e2bcd8739ae039(element) {
    // Firefox wants us to check `-x` and `-y` variations as well
    var _getComputedStyle = (0, $7d74e8e12f28eb07$export$2e2bcd8739ae039)(element), overflow = _getComputedStyle.overflow, overflowX = _getComputedStyle.overflowX, overflowY = _getComputedStyle.overflowY;
    return /auto|scroll|overlay|hidden/.test(overflow + overflowY + overflowX);
}



function $cd2b59ec72f3c621$var$isElementScaled(element) {
    var rect = element.getBoundingClientRect();
    var scaleX = (0, $9c234cc6f4ed3b92$export$2077e0241d6afd3c)(rect.width) / element.offsetWidth || 1;
    var scaleY = (0, $9c234cc6f4ed3b92$export$2077e0241d6afd3c)(rect.height) / element.offsetHeight || 1;
    return scaleX !== 1 || scaleY !== 1;
} // Returns the composite rect of an element relative to its offsetParent.
function $cd2b59ec72f3c621$export$2e2bcd8739ae039(elementOrVirtualElement, offsetParent, isFixed) {
    if (isFixed === void 0) isFixed = false;
    var isOffsetParentAnElement = (0, $b83e73df1a45d10b$export$1b3bfaa9684536aa)(offsetParent);
    var offsetParentIsScaled = (0, $b83e73df1a45d10b$export$1b3bfaa9684536aa)(offsetParent) && $cd2b59ec72f3c621$var$isElementScaled(offsetParent);
    var documentElement = (0, $d56bbea5dd2fe00d$export$2e2bcd8739ae039)(offsetParent);
    var rect = (0, $d5b6692e6e485772$export$2e2bcd8739ae039)(elementOrVirtualElement, offsetParentIsScaled, isFixed);
    var scroll = {
        scrollLeft: 0,
        scrollTop: 0
    };
    var offsets = {
        x: 0,
        y: 0
    };
    if (isOffsetParentAnElement || !isOffsetParentAnElement && !isFixed) {
        if ((0, $560fb2f7f26d74c6$export$2e2bcd8739ae039)(offsetParent) !== 'body' || // https://github.com/popperjs/popper-core/issues/1078
        (0, $3e03f71b3bbb98eb$export$2e2bcd8739ae039)(documentElement)) scroll = (0, $1416989062dadf75$export$2e2bcd8739ae039)(offsetParent);
        if ((0, $b83e73df1a45d10b$export$1b3bfaa9684536aa)(offsetParent)) {
            offsets = (0, $d5b6692e6e485772$export$2e2bcd8739ae039)(offsetParent, true);
            offsets.x += offsetParent.clientLeft;
            offsets.y += offsetParent.clientTop;
        } else if (documentElement) offsets.x = (0, $c4822cfd8628a2db$export$2e2bcd8739ae039)(documentElement);
    }
    return {
        x: rect.left + scroll.scrollLeft - offsets.x,
        y: rect.top + scroll.scrollTop - offsets.y,
        width: rect.width,
        height: rect.height
    };
}



function $f4ac90f91f825ce1$export$2e2bcd8739ae039(element) {
    var clientRect = (0, $d5b6692e6e485772$export$2e2bcd8739ae039)(element); // Use the clientRect sizes if it's not been transformed.
    // Fixes https://github.com/popperjs/popper-core/issues/1223
    var width = element.offsetWidth;
    var height = element.offsetHeight;
    if (Math.abs(clientRect.width - width) <= 1) width = clientRect.width;
    if (Math.abs(clientRect.height - height) <= 1) height = clientRect.height;
    return {
        x: element.offsetLeft,
        y: element.offsetTop,
        width: width,
        height: height
    };
}





function $f019b932829cd8d0$export$2e2bcd8739ae039(element) {
    if ((0, $560fb2f7f26d74c6$export$2e2bcd8739ae039)(element) === 'html') return element;
    return(// $FlowFixMe[incompatible-return]
    // $FlowFixMe[prop-missing]
    element.assignedSlot || // step into the shadow DOM of the parent of a slotted node
    element.parentNode || ((0, $b83e73df1a45d10b$export$af51f0f06c0f328a)(element) ? element.host : null) || // ShadowRoot detected
    // $FlowFixMe[incompatible-call]: HTMLElement is a Node
    (0, $d56bbea5dd2fe00d$export$2e2bcd8739ae039)(element) // fallback
    );
}





function $7b745ff2b79fb468$export$2e2bcd8739ae039(node) {
    if ([
        'html',
        'body',
        '#document'
    ].indexOf((0, $560fb2f7f26d74c6$export$2e2bcd8739ae039)(node)) >= 0) // $FlowFixMe[incompatible-return]: assume body is always available
    return node.ownerDocument.body;
    if ((0, $b83e73df1a45d10b$export$1b3bfaa9684536aa)(node) && (0, $3e03f71b3bbb98eb$export$2e2bcd8739ae039)(node)) return node;
    return $7b745ff2b79fb468$export$2e2bcd8739ae039((0, $f019b932829cd8d0$export$2e2bcd8739ae039)(node));
}





function $ebe9d77e8ccdd6fa$export$2e2bcd8739ae039(element, list) {
    var _element$ownerDocumen;
    if (list === void 0) list = [];
    var scrollParent = (0, $7b745ff2b79fb468$export$2e2bcd8739ae039)(element);
    var isBody = scrollParent === ((_element$ownerDocumen = element.ownerDocument) == null ? void 0 : _element$ownerDocumen.body);
    var win = (0, $e992165d27cf1e50$export$2e2bcd8739ae039)(scrollParent);
    var target = isBody ? [
        win
    ].concat(win.visualViewport || [], (0, $3e03f71b3bbb98eb$export$2e2bcd8739ae039)(scrollParent) ? scrollParent : []) : scrollParent;
    var updatedList = list.concat(target);
    return isBody ? updatedList : updatedList.concat($ebe9d77e8ccdd6fa$export$2e2bcd8739ae039((0, $f019b932829cd8d0$export$2e2bcd8739ae039)(target)));
}







function $224b24caa0d1b578$export$2e2bcd8739ae039(element) {
    return [
        'table',
        'td',
        'th'
    ].indexOf((0, $560fb2f7f26d74c6$export$2e2bcd8739ae039)(element)) >= 0;
}




function $01dd27adea182f9a$var$getTrueOffsetParent(element) {
    if (!(0, $b83e73df1a45d10b$export$1b3bfaa9684536aa)(element) || // https://github.com/popperjs/popper-core/issues/837
    (0, $7d74e8e12f28eb07$export$2e2bcd8739ae039)(element).position === 'fixed') return null;
    return element.offsetParent;
} // `.offsetParent` reports `null` for fixed elements, while absolute elements
// return the containing block
function $01dd27adea182f9a$var$getContainingBlock(element) {
    var isFirefox = /firefox/i.test((0, $5e056ec93b493b76$export$2e2bcd8739ae039)());
    var isIE = /Trident/i.test((0, $5e056ec93b493b76$export$2e2bcd8739ae039)());
    if (isIE && (0, $b83e73df1a45d10b$export$1b3bfaa9684536aa)(element)) {
        // In IE 9, 10 and 11 fixed elements containing block is always established by the viewport
        var elementCss = (0, $7d74e8e12f28eb07$export$2e2bcd8739ae039)(element);
        if (elementCss.position === 'fixed') return null;
    }
    var currentNode = (0, $f019b932829cd8d0$export$2e2bcd8739ae039)(element);
    if ((0, $b83e73df1a45d10b$export$af51f0f06c0f328a)(currentNode)) currentNode = currentNode.host;
    while((0, $b83e73df1a45d10b$export$1b3bfaa9684536aa)(currentNode) && [
        'html',
        'body'
    ].indexOf((0, $560fb2f7f26d74c6$export$2e2bcd8739ae039)(currentNode)) < 0){
        var css = (0, $7d74e8e12f28eb07$export$2e2bcd8739ae039)(currentNode); // This is non-exhaustive but covers the most common CSS properties that
        // create a containing block.
        // https://developer.mozilla.org/en-US/docs/Web/CSS/Containing_block#identifying_the_containing_block
        if (css.transform !== 'none' || css.perspective !== 'none' || css.contain === 'paint' || [
            'transform',
            'perspective'
        ].indexOf(css.willChange) !== -1 || isFirefox && css.willChange === 'filter' || isFirefox && css.filter && css.filter !== 'none') return currentNode;
        else currentNode = currentNode.parentNode;
    }
    return null;
} // Gets the closest ancestor positioned element. Handles some edge cases,
function $01dd27adea182f9a$export$2e2bcd8739ae039(element) {
    var window = (0, $e992165d27cf1e50$export$2e2bcd8739ae039)(element);
    var offsetParent = $01dd27adea182f9a$var$getTrueOffsetParent(element);
    while(offsetParent && (0, $224b24caa0d1b578$export$2e2bcd8739ae039)(offsetParent) && (0, $7d74e8e12f28eb07$export$2e2bcd8739ae039)(offsetParent).position === 'static')offsetParent = $01dd27adea182f9a$var$getTrueOffsetParent(offsetParent);
    if (offsetParent && ((0, $560fb2f7f26d74c6$export$2e2bcd8739ae039)(offsetParent) === 'html' || (0, $560fb2f7f26d74c6$export$2e2bcd8739ae039)(offsetParent) === 'body' && (0, $7d74e8e12f28eb07$export$2e2bcd8739ae039)(offsetParent).position === 'static')) return window;
    return offsetParent || $01dd27adea182f9a$var$getContainingBlock(element) || window;
}


var $8478d2765df37982$export$1e95b668f3b82d = 'top';
var $8478d2765df37982$export$40e543e69a8b3fbb = 'bottom';
var $8478d2765df37982$export$79ffe56a765070d2 = 'right';
var $8478d2765df37982$export$eabcd2c8791e7bf4 = 'left';
var $8478d2765df37982$export$dfb5619354ba860 = 'auto';
var $8478d2765df37982$export$aec2ce47c367b8c3 = [
    $8478d2765df37982$export$1e95b668f3b82d,
    $8478d2765df37982$export$40e543e69a8b3fbb,
    $8478d2765df37982$export$79ffe56a765070d2,
    $8478d2765df37982$export$eabcd2c8791e7bf4
];
var $8478d2765df37982$export$b3571188c770cc5a = 'start';
var $8478d2765df37982$export$bd5df0f255a350f8 = 'end';
var $8478d2765df37982$export$390fd549c5303b4d = 'clippingParents';
var $8478d2765df37982$export$d7b7311ec04a3e8f = 'viewport';
var $8478d2765df37982$export$ae5ab1c730825774 = 'popper';
var $8478d2765df37982$export$ca50aac9f3ba507f = 'reference';
var $8478d2765df37982$export$368f9a87e87fa4e1 = /*#__PURE__*/ $8478d2765df37982$export$aec2ce47c367b8c3.reduce(function(acc, placement) {
    return acc.concat([
        placement + "-" + $8478d2765df37982$export$b3571188c770cc5a,
        placement + "-" + $8478d2765df37982$export$bd5df0f255a350f8
    ]);
}, []);
var $8478d2765df37982$export$803cd8101b6c182b = /*#__PURE__*/ [].concat($8478d2765df37982$export$aec2ce47c367b8c3, [
    $8478d2765df37982$export$dfb5619354ba860
]).reduce(function(acc, placement) {
    return acc.concat([
        placement,
        placement + "-" + $8478d2765df37982$export$b3571188c770cc5a,
        placement + "-" + $8478d2765df37982$export$bd5df0f255a350f8
    ]);
}, []); // modifiers that need to read the DOM
var $8478d2765df37982$export$421679a7c3d56e = 'beforeRead';
var $8478d2765df37982$export$aafa59e2e03f2942 = 'read';
var $8478d2765df37982$export$6964f6c886723980 = 'afterRead'; // pure-logic modifiers
var $8478d2765df37982$export$c65e99957a05207c = 'beforeMain';
var $8478d2765df37982$export$f22da7240b7add18 = 'main';
var $8478d2765df37982$export$bab79516f2d662fe = 'afterMain'; // modifier with the purpose to write to the DOM (or write into a framework state)
var $8478d2765df37982$export$8d4d2d70e7d46032 = 'beforeWrite';
var $8478d2765df37982$export$68d8715fc104d294 = 'write';
var $8478d2765df37982$export$70a6e5159acce2e6 = 'afterWrite';
var $8478d2765df37982$export$d087d3878fdf71d5 = [
    $8478d2765df37982$export$421679a7c3d56e,
    $8478d2765df37982$export$aafa59e2e03f2942,
    $8478d2765df37982$export$6964f6c886723980,
    $8478d2765df37982$export$c65e99957a05207c,
    $8478d2765df37982$export$f22da7240b7add18,
    $8478d2765df37982$export$bab79516f2d662fe,
    $8478d2765df37982$export$8d4d2d70e7d46032,
    $8478d2765df37982$export$68d8715fc104d294,
    $8478d2765df37982$export$70a6e5159acce2e6
];


function $18833a8410471e6e$var$order(modifiers) {
    var map = new Map();
    var visited = new Set();
    var result = [];
    modifiers.forEach(function(modifier) {
        map.set(modifier.name, modifier);
    }); // On visiting object, check for its dependencies and visit them recursively
    function sort(modifier) {
        visited.add(modifier.name);
        var requires = [].concat(modifier.requires || [], modifier.requiresIfExists || []);
        requires.forEach(function(dep) {
            if (!visited.has(dep)) {
                var depModifier = map.get(dep);
                if (depModifier) sort(depModifier);
            }
        });
        result.push(modifier);
    }
    modifiers.forEach(function(modifier) {
        if (!visited.has(modifier.name)) // check for visited object
        sort(modifier);
    });
    return result;
}
function $18833a8410471e6e$export$2e2bcd8739ae039(modifiers) {
    // order based on dependencies
    var orderedModifiers = $18833a8410471e6e$var$order(modifiers); // order based on phase
    return (0, $8478d2765df37982$export$d087d3878fdf71d5).reduce(function(acc, phase) {
        return acc.concat(orderedModifiers.filter(function(modifier) {
            return modifier.phase === phase;
        }));
    }, []);
}


function $e8ac179c28da1322$export$2e2bcd8739ae039(fn) {
    var pending;
    return function() {
        if (!pending) pending = new Promise(function(resolve) {
            Promise.resolve().then(function() {
                pending = undefined;
                resolve(fn());
            });
        });
        return pending;
    };
}


function $12b179ff5efeeee3$export$2e2bcd8739ae039(modifiers) {
    var merged = modifiers.reduce(function(merged, current) {
        var existing = merged[current.name];
        merged[current.name] = existing ? Object.assign({}, existing, current, {
            options: Object.assign({}, existing.options, current.options),
            data: Object.assign({}, existing.data, current.data)
        }) : current;
        return merged;
    }, {}); // IE11 does not support Object.values
    return Object.keys(merged).map(function(key) {
        return merged[key];
    });
}




var $f3aacc53d895d057$var$DEFAULT_OPTIONS = {
    placement: 'bottom',
    modifiers: [],
    strategy: 'absolute'
};
function $f3aacc53d895d057$var$areValidElements() {
    for(var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++)args[_key] = arguments[_key];
    return !args.some(function(element) {
        return !(element && typeof element.getBoundingClientRect === 'function');
    });
}
function $f3aacc53d895d057$export$ed5e13716264f202(generatorOptions) {
    if (generatorOptions === void 0) generatorOptions = {};
    var _generatorOptions = generatorOptions, _generatorOptions$def = _generatorOptions.defaultModifiers, defaultModifiers = _generatorOptions$def === void 0 ? [] : _generatorOptions$def, _generatorOptions$def2 = _generatorOptions.defaultOptions, defaultOptions = _generatorOptions$def2 === void 0 ? $f3aacc53d895d057$var$DEFAULT_OPTIONS : _generatorOptions$def2;
    return function createPopper(reference, popper, options) {
        if (options === void 0) options = defaultOptions;
        var state = {
            placement: 'bottom',
            orderedModifiers: [],
            options: Object.assign({}, $f3aacc53d895d057$var$DEFAULT_OPTIONS, defaultOptions),
            modifiersData: {},
            elements: {
                reference: reference,
                popper: popper
            },
            attributes: {},
            styles: {}
        };
        var effectCleanupFns = [];
        var isDestroyed = false;
        var instance = {
            state: state,
            setOptions: function setOptions(setOptionsAction) {
                var options = typeof setOptionsAction === 'function' ? setOptionsAction(state.options) : setOptionsAction;
                cleanupModifierEffects();
                state.options = Object.assign({}, defaultOptions, state.options, options);
                state.scrollParents = {
                    reference: (0, $b83e73df1a45d10b$export$45a5e7f76e0caa8d)(reference) ? (0, $ebe9d77e8ccdd6fa$export$2e2bcd8739ae039)(reference) : reference.contextElement ? (0, $ebe9d77e8ccdd6fa$export$2e2bcd8739ae039)(reference.contextElement) : [],
                    popper: (0, $ebe9d77e8ccdd6fa$export$2e2bcd8739ae039)(popper)
                }; // Orders the modifiers based on their dependencies and `phase`
                // properties
                var orderedModifiers = (0, $18833a8410471e6e$export$2e2bcd8739ae039)((0, $12b179ff5efeeee3$export$2e2bcd8739ae039)([].concat(defaultModifiers, state.options.modifiers))); // Strip out disabled modifiers
                state.orderedModifiers = orderedModifiers.filter(function(m) {
                    return m.enabled;
                });
                runModifierEffects();
                return instance.update();
            },
            // Sync update – it will always be executed, even if not necessary. This
            // is useful for low frequency updates where sync behavior simplifies the
            // logic.
            // For high frequency updates (e.g. `resize` and `scroll` events), always
            // prefer the async Popper#update method
            forceUpdate: function forceUpdate() {
                if (isDestroyed) return;
                var _state$elements = state.elements, reference = _state$elements.reference, popper = _state$elements.popper; // Don't proceed if `reference` or `popper` are not valid elements
                // anymore
                if (!$f3aacc53d895d057$var$areValidElements(reference, popper)) return;
                 // Store the reference and popper rects to be read by modifiers
                state.rects = {
                    reference: (0, $cd2b59ec72f3c621$export$2e2bcd8739ae039)(reference, (0, $01dd27adea182f9a$export$2e2bcd8739ae039)(popper), state.options.strategy === 'fixed'),
                    popper: (0, $f4ac90f91f825ce1$export$2e2bcd8739ae039)(popper)
                }; // Modifiers have the ability to reset the current update cycle. The
                // most common use case for this is the `flip` modifier changing the
                // placement, which then needs to re-run all the modifiers, because the
                // logic was previously ran for the previous placement and is therefore
                // stale/incorrect
                state.reset = false;
                state.placement = state.options.placement; // On each update cycle, the `modifiersData` property for each modifier
                // is filled with the initial data specified by the modifier. This means
                // it doesn't persist and is fresh on each update.
                // To ensure persistent data, use `${name}#persistent`
                state.orderedModifiers.forEach(function(modifier) {
                    return state.modifiersData[modifier.name] = Object.assign({}, modifier.data);
                });
                for(var index = 0; index < state.orderedModifiers.length; index++){
                    if (state.reset === true) {
                        state.reset = false;
                        index = -1;
                        continue;
                    }
                    var _state$orderedModifie = state.orderedModifiers[index], fn = _state$orderedModifie.fn, _state$orderedModifie2 = _state$orderedModifie.options, _options = _state$orderedModifie2 === void 0 ? {} : _state$orderedModifie2, name = _state$orderedModifie.name;
                    if (typeof fn === 'function') state = fn({
                        state: state,
                        options: _options,
                        name: name,
                        instance: instance
                    }) || state;
                }
            },
            // Async and optimistically optimized update – it will not be executed if
            // not necessary (debounced to run at most once-per-tick)
            update: (0, $e8ac179c28da1322$export$2e2bcd8739ae039)(function() {
                return new Promise(function(resolve) {
                    instance.forceUpdate();
                    resolve(state);
                });
            }),
            destroy: function destroy() {
                cleanupModifierEffects();
                isDestroyed = true;
            }
        };
        if (!$f3aacc53d895d057$var$areValidElements(reference, popper)) return instance;
        instance.setOptions(options).then(function(state) {
            if (!isDestroyed && options.onFirstUpdate) options.onFirstUpdate(state);
        }); // Modifiers have the ability to execute arbitrary code before the first
        // update cycle runs. They will be executed in the same order as the update
        // cycle. This is useful when a modifier adds some persistent data that
        // other modifiers need to use, but the modifier is run after the dependent
        // one.
        function runModifierEffects() {
            state.orderedModifiers.forEach(function(_ref) {
                var name = _ref.name, _ref$options = _ref.options, options = _ref$options === void 0 ? {} : _ref$options, effect = _ref.effect;
                if (typeof effect === 'function') {
                    var cleanupFn = effect({
                        state: state,
                        name: name,
                        instance: instance,
                        options: options
                    });
                    var noopFn = function noopFn() {};
                    effectCleanupFns.push(cleanupFn || noopFn);
                }
            });
        }
        function cleanupModifierEffects() {
            effectCleanupFns.forEach(function(fn) {
                return fn();
            });
            effectCleanupFns = [];
        }
        return instance;
    };
}
var $f3aacc53d895d057$export$8f7491d57c8f97a9 = /*#__PURE__*/ $f3aacc53d895d057$export$ed5e13716264f202(); // eslint-disable-next-line import/no-unused-modules



var $48430453eb0d053b$var$passive = {
    passive: true
};
function $48430453eb0d053b$var$effect(_ref) {
    var state = _ref.state, instance = _ref.instance, options = _ref.options;
    var _options$scroll = options.scroll, scroll = _options$scroll === void 0 ? true : _options$scroll, _options$resize = options.resize, resize = _options$resize === void 0 ? true : _options$resize;
    var window = (0, $e992165d27cf1e50$export$2e2bcd8739ae039)(state.elements.popper);
    var scrollParents = [].concat(state.scrollParents.reference, state.scrollParents.popper);
    if (scroll) scrollParents.forEach(function(scrollParent) {
        scrollParent.addEventListener('scroll', instance.update, $48430453eb0d053b$var$passive);
    });
    if (resize) window.addEventListener('resize', instance.update, $48430453eb0d053b$var$passive);
    return function() {
        if (scroll) scrollParents.forEach(function(scrollParent) {
            scrollParent.removeEventListener('scroll', instance.update, $48430453eb0d053b$var$passive);
        });
        if (resize) window.removeEventListener('resize', instance.update, $48430453eb0d053b$var$passive);
    };
} // eslint-disable-next-line import/no-unused-modules
var $48430453eb0d053b$export$2e2bcd8739ae039 = {
    name: 'eventListeners',
    enabled: true,
    phase: 'write',
    fn: function fn() {},
    effect: $48430453eb0d053b$var$effect,
    data: {}
};



function $f33c93a011fffe7b$export$2e2bcd8739ae039(placement) {
    return placement.split('-')[0];
}


function $c5e5c390efc1bad4$export$2e2bcd8739ae039(placement) {
    return placement.split('-')[1];
}


function $247260022a0ad6ee$export$2e2bcd8739ae039(placement) {
    return [
        'top',
        'bottom'
    ].indexOf(placement) >= 0 ? 'x' : 'y';
}



function $87655a20e47aab91$export$2e2bcd8739ae039(_ref) {
    var reference = _ref.reference, element = _ref.element, placement = _ref.placement;
    var basePlacement = placement ? (0, $f33c93a011fffe7b$export$2e2bcd8739ae039)(placement) : null;
    var variation = placement ? (0, $c5e5c390efc1bad4$export$2e2bcd8739ae039)(placement) : null;
    var commonX = reference.x + reference.width / 2 - element.width / 2;
    var commonY = reference.y + reference.height / 2 - element.height / 2;
    var offsets;
    switch(basePlacement){
        case 0, $8478d2765df37982$export$1e95b668f3b82d:
            offsets = {
                x: commonX,
                y: reference.y - element.height
            };
            break;
        case 0, $8478d2765df37982$export$40e543e69a8b3fbb:
            offsets = {
                x: commonX,
                y: reference.y + reference.height
            };
            break;
        case 0, $8478d2765df37982$export$79ffe56a765070d2:
            offsets = {
                x: reference.x + reference.width,
                y: commonY
            };
            break;
        case 0, $8478d2765df37982$export$eabcd2c8791e7bf4:
            offsets = {
                x: reference.x - element.width,
                y: commonY
            };
            break;
        default:
            offsets = {
                x: reference.x,
                y: reference.y
            };
    }
    var mainAxis = basePlacement ? (0, $247260022a0ad6ee$export$2e2bcd8739ae039)(basePlacement) : null;
    if (mainAxis != null) {
        var len = mainAxis === 'y' ? 'height' : 'width';
        switch(variation){
            case 0, $8478d2765df37982$export$b3571188c770cc5a:
                offsets[mainAxis] = offsets[mainAxis] - (reference[len] / 2 - element[len] / 2);
                break;
            case 0, $8478d2765df37982$export$bd5df0f255a350f8:
                offsets[mainAxis] = offsets[mainAxis] + (reference[len] / 2 - element[len] / 2);
                break;
            default:
        }
    }
    return offsets;
}


function $b19a98045ea3e977$var$popperOffsets(_ref) {
    var state = _ref.state, name = _ref.name;
    // Offsets are the actual position the popper needs to have to be
    // properly positioned near its reference element
    // This is the most basic placement, and will be adjusted by
    // the modifiers in the next step
    state.modifiersData[name] = (0, $87655a20e47aab91$export$2e2bcd8739ae039)({
        reference: state.rects.reference,
        element: state.rects.popper,
        strategy: 'absolute',
        placement: state.placement
    });
} // eslint-disable-next-line import/no-unused-modules
var $b19a98045ea3e977$export$2e2bcd8739ae039 = {
    name: 'popperOffsets',
    enabled: true,
    phase: 'read',
    fn: $b19a98045ea3e977$var$popperOffsets,
    data: {}
};










var $c5e18f9acfdbb8a2$var$unsetSides = {
    top: 'auto',
    right: 'auto',
    bottom: 'auto',
    left: 'auto'
}; // Round the offsets to the nearest suitable subpixel based on the DPR.
// Zooming can change the DPR, but it seems to report a value that will
// cleanly divide the values into the appropriate subpixels.
function $c5e18f9acfdbb8a2$var$roundOffsetsByDPR(_ref, win) {
    var x = _ref.x, y = _ref.y;
    var dpr = win.devicePixelRatio || 1;
    return {
        x: (0, $9c234cc6f4ed3b92$export$2077e0241d6afd3c)(x * dpr) / dpr || 0,
        y: (0, $9c234cc6f4ed3b92$export$2077e0241d6afd3c)(y * dpr) / dpr || 0
    };
}
function $c5e18f9acfdbb8a2$export$378fa78a8fea596f(_ref2) {
    var _Object$assign2;
    var popper = _ref2.popper, popperRect = _ref2.popperRect, placement = _ref2.placement, variation = _ref2.variation, offsets = _ref2.offsets, position = _ref2.position, gpuAcceleration = _ref2.gpuAcceleration, adaptive = _ref2.adaptive, roundOffsets = _ref2.roundOffsets, isFixed = _ref2.isFixed;
    var _offsets$x = offsets.x, x = _offsets$x === void 0 ? 0 : _offsets$x, _offsets$y = offsets.y, y = _offsets$y === void 0 ? 0 : _offsets$y;
    var _ref3 = typeof roundOffsets === 'function' ? roundOffsets({
        x: x,
        y: y
    }) : {
        x: x,
        y: y
    };
    x = _ref3.x;
    y = _ref3.y;
    var hasX = offsets.hasOwnProperty('x');
    var hasY = offsets.hasOwnProperty('y');
    var sideX = (0, $8478d2765df37982$export$eabcd2c8791e7bf4);
    var sideY = (0, $8478d2765df37982$export$1e95b668f3b82d);
    var win = window;
    if (adaptive) {
        var offsetParent = (0, $01dd27adea182f9a$export$2e2bcd8739ae039)(popper);
        var heightProp = 'clientHeight';
        var widthProp = 'clientWidth';
        if (offsetParent === (0, $e992165d27cf1e50$export$2e2bcd8739ae039)(popper)) {
            offsetParent = (0, $d56bbea5dd2fe00d$export$2e2bcd8739ae039)(popper);
            if ((0, $7d74e8e12f28eb07$export$2e2bcd8739ae039)(offsetParent).position !== 'static' && position === 'absolute') {
                heightProp = 'scrollHeight';
                widthProp = 'scrollWidth';
            }
        } // $FlowFixMe[incompatible-cast]: force type refinement, we compare offsetParent with window above, but Flow doesn't detect it
        offsetParent;
        if (placement === (0, $8478d2765df37982$export$1e95b668f3b82d) || (placement === (0, $8478d2765df37982$export$eabcd2c8791e7bf4) || placement === (0, $8478d2765df37982$export$79ffe56a765070d2)) && variation === (0, $8478d2765df37982$export$bd5df0f255a350f8)) {
            sideY = (0, $8478d2765df37982$export$40e543e69a8b3fbb);
            var offsetY = isFixed && offsetParent === win && win.visualViewport ? win.visualViewport.height : offsetParent[heightProp];
            y -= offsetY - popperRect.height;
            y *= gpuAcceleration ? 1 : -1;
        }
        if (placement === (0, $8478d2765df37982$export$eabcd2c8791e7bf4) || (placement === (0, $8478d2765df37982$export$1e95b668f3b82d) || placement === (0, $8478d2765df37982$export$40e543e69a8b3fbb)) && variation === (0, $8478d2765df37982$export$bd5df0f255a350f8)) {
            sideX = (0, $8478d2765df37982$export$79ffe56a765070d2);
            var offsetX = isFixed && offsetParent === win && win.visualViewport ? win.visualViewport.width : offsetParent[widthProp];
            x -= offsetX - popperRect.width;
            x *= gpuAcceleration ? 1 : -1;
        }
    }
    var commonStyles = Object.assign({
        position: position
    }, adaptive && $c5e18f9acfdbb8a2$var$unsetSides);
    var _ref4 = roundOffsets === true ? $c5e18f9acfdbb8a2$var$roundOffsetsByDPR({
        x: x,
        y: y
    }, (0, $e992165d27cf1e50$export$2e2bcd8739ae039)(popper)) : {
        x: x,
        y: y
    };
    x = _ref4.x;
    y = _ref4.y;
    if (gpuAcceleration) {
        var _Object$assign;
        return Object.assign({}, commonStyles, (_Object$assign = {}, _Object$assign[sideY] = hasY ? '0' : '', _Object$assign[sideX] = hasX ? '0' : '', _Object$assign.transform = (win.devicePixelRatio || 1) <= 1 ? "translate(" + x + "px, " + y + "px)" : "translate3d(" + x + "px, " + y + "px, 0)", _Object$assign));
    }
    return Object.assign({}, commonStyles, (_Object$assign2 = {}, _Object$assign2[sideY] = hasY ? y + "px" : '', _Object$assign2[sideX] = hasX ? x + "px" : '', _Object$assign2.transform = '', _Object$assign2));
}
function $c5e18f9acfdbb8a2$var$computeStyles(_ref5) {
    var state = _ref5.state, options = _ref5.options;
    var _options$gpuAccelerat = options.gpuAcceleration, gpuAcceleration = _options$gpuAccelerat === void 0 ? true : _options$gpuAccelerat, _options$adaptive = options.adaptive, adaptive = _options$adaptive === void 0 ? true : _options$adaptive, _options$roundOffsets = options.roundOffsets, roundOffsets = _options$roundOffsets === void 0 ? true : _options$roundOffsets;
    var commonStyles = {
        placement: (0, $f33c93a011fffe7b$export$2e2bcd8739ae039)(state.placement),
        variation: (0, $c5e5c390efc1bad4$export$2e2bcd8739ae039)(state.placement),
        popper: state.elements.popper,
        popperRect: state.rects.popper,
        gpuAcceleration: gpuAcceleration,
        isFixed: state.options.strategy === 'fixed'
    };
    if (state.modifiersData.popperOffsets != null) state.styles.popper = Object.assign({}, state.styles.popper, $c5e18f9acfdbb8a2$export$378fa78a8fea596f(Object.assign({}, commonStyles, {
        offsets: state.modifiersData.popperOffsets,
        position: state.options.strategy,
        adaptive: adaptive,
        roundOffsets: roundOffsets
    })));
    if (state.modifiersData.arrow != null) state.styles.arrow = Object.assign({}, state.styles.arrow, $c5e18f9acfdbb8a2$export$378fa78a8fea596f(Object.assign({}, commonStyles, {
        offsets: state.modifiersData.arrow,
        position: 'absolute',
        adaptive: false,
        roundOffsets: roundOffsets
    })));
    state.attributes.popper = Object.assign({}, state.attributes.popper, {
        'data-popper-placement': state.placement
    });
} // eslint-disable-next-line import/no-unused-modules
var $c5e18f9acfdbb8a2$export$2e2bcd8739ae039 = {
    name: 'computeStyles',
    enabled: true,
    phase: 'beforeWrite',
    fn: $c5e18f9acfdbb8a2$var$computeStyles,
    data: {}
};




// and applies them to the HTMLElements such as popper and arrow
function $1c7a027d393b751c$var$applyStyles(_ref) {
    var state = _ref.state;
    Object.keys(state.elements).forEach(function(name) {
        var style = state.styles[name] || {};
        var attributes = state.attributes[name] || {};
        var element = state.elements[name]; // arrow is optional + virtual elements
        if (!(0, $b83e73df1a45d10b$export$1b3bfaa9684536aa)(element) || !(0, $560fb2f7f26d74c6$export$2e2bcd8739ae039)(element)) return;
         // Flow doesn't support to extend this property, but it's the most
        // effective way to apply styles to an HTMLElement
        // $FlowFixMe[cannot-write]
        Object.assign(element.style, style);
        Object.keys(attributes).forEach(function(name) {
            var value = attributes[name];
            if (value === false) element.removeAttribute(name);
            else element.setAttribute(name, value === true ? '' : value);
        });
    });
}
function $1c7a027d393b751c$var$effect(_ref2) {
    var state = _ref2.state;
    var initialStyles = {
        popper: {
            position: state.options.strategy,
            left: '0',
            top: '0',
            margin: '0'
        },
        arrow: {
            position: 'absolute'
        },
        reference: {}
    };
    Object.assign(state.elements.popper.style, initialStyles.popper);
    state.styles = initialStyles;
    if (state.elements.arrow) Object.assign(state.elements.arrow.style, initialStyles.arrow);
    return function() {
        Object.keys(state.elements).forEach(function(name) {
            var element = state.elements[name];
            var attributes = state.attributes[name] || {};
            var styleProperties = Object.keys(state.styles.hasOwnProperty(name) ? state.styles[name] : initialStyles[name]); // Set all values to an empty string to unset them
            var style = styleProperties.reduce(function(style, property) {
                style[property] = '';
                return style;
            }, {}); // arrow is optional + virtual elements
            if (!(0, $b83e73df1a45d10b$export$1b3bfaa9684536aa)(element) || !(0, $560fb2f7f26d74c6$export$2e2bcd8739ae039)(element)) return;
            Object.assign(element.style, style);
            Object.keys(attributes).forEach(function(attribute) {
                element.removeAttribute(attribute);
            });
        });
    };
} // eslint-disable-next-line import/no-unused-modules
var $1c7a027d393b751c$export$2e2bcd8739ae039 = {
    name: 'applyStyles',
    enabled: true,
    phase: 'write',
    fn: $1c7a027d393b751c$var$applyStyles,
    effect: $1c7a027d393b751c$var$effect,
    requires: [
        'computeStyles'
    ]
};




function $f0fe73ce831ce17b$export$7fa02d8595b015ed(placement, rects, offset) {
    var basePlacement = (0, $f33c93a011fffe7b$export$2e2bcd8739ae039)(placement);
    var invertDistance = [
        (0, $8478d2765df37982$export$eabcd2c8791e7bf4),
        (0, $8478d2765df37982$export$1e95b668f3b82d)
    ].indexOf(basePlacement) >= 0 ? -1 : 1;
    var _ref = typeof offset === 'function' ? offset(Object.assign({}, rects, {
        placement: placement
    })) : offset, skidding = _ref[0], distance = _ref[1];
    skidding = skidding || 0;
    distance = (distance || 0) * invertDistance;
    return [
        (0, $8478d2765df37982$export$eabcd2c8791e7bf4),
        (0, $8478d2765df37982$export$79ffe56a765070d2)
    ].indexOf(basePlacement) >= 0 ? {
        x: distance,
        y: skidding
    } : {
        x: skidding,
        y: distance
    };
}
function $f0fe73ce831ce17b$var$offset(_ref2) {
    var state = _ref2.state, options = _ref2.options, name = _ref2.name;
    var _options$offset = options.offset, offset = _options$offset === void 0 ? [
        0,
        0
    ] : _options$offset;
    var data = (0, $8478d2765df37982$export$803cd8101b6c182b).reduce(function(acc, placement) {
        acc[placement] = $f0fe73ce831ce17b$export$7fa02d8595b015ed(placement, state.rects, offset);
        return acc;
    }, {});
    var _data$state$placement = data[state.placement], x = _data$state$placement.x, y = _data$state$placement.y;
    if (state.modifiersData.popperOffsets != null) {
        state.modifiersData.popperOffsets.x += x;
        state.modifiersData.popperOffsets.y += y;
    }
    state.modifiersData[name] = data;
} // eslint-disable-next-line import/no-unused-modules
var $f0fe73ce831ce17b$export$2e2bcd8739ae039 = {
    name: 'offset',
    enabled: true,
    phase: 'main',
    requires: [
        'popperOffsets'
    ],
    fn: $f0fe73ce831ce17b$var$offset
};


var $c0bd8d5905df73b7$var$hash = {
    left: 'right',
    right: 'left',
    bottom: 'top',
    top: 'bottom'
};
function $c0bd8d5905df73b7$export$2e2bcd8739ae039(placement) {
    return placement.replace(/left|right|bottom|top/g, function(matched) {
        return $c0bd8d5905df73b7$var$hash[matched];
    });
}



var $a0a642892f75993d$var$hash = {
    start: 'end',
    end: 'start'
};
function $a0a642892f75993d$export$2e2bcd8739ae039(placement) {
    return placement.replace(/start|end/g, function(matched) {
        return $a0a642892f75993d$var$hash[matched];
    });
}







function $b4831f59b241eaa3$export$2e2bcd8739ae039(element, strategy) {
    var win = (0, $e992165d27cf1e50$export$2e2bcd8739ae039)(element);
    var html = (0, $d56bbea5dd2fe00d$export$2e2bcd8739ae039)(element);
    var visualViewport = win.visualViewport;
    var width = html.clientWidth;
    var height = html.clientHeight;
    var x = 0;
    var y = 0;
    if (visualViewport) {
        width = visualViewport.width;
        height = visualViewport.height;
        var layoutViewport = (0, $ab0cb3910b57b210$export$2e2bcd8739ae039)();
        if (layoutViewport || !layoutViewport && strategy === 'fixed') {
            x = visualViewport.offsetLeft;
            y = visualViewport.offsetTop;
        }
    }
    return {
        width: width,
        height: height,
        x: x + (0, $c4822cfd8628a2db$export$2e2bcd8739ae039)(element),
        y: y
    };
}







function $f086f128b1379452$export$2e2bcd8739ae039(element) {
    var _element$ownerDocumen;
    var html = (0, $d56bbea5dd2fe00d$export$2e2bcd8739ae039)(element);
    var winScroll = (0, $2270e2f70c3d3867$export$2e2bcd8739ae039)(element);
    var body = (_element$ownerDocumen = element.ownerDocument) == null ? void 0 : _element$ownerDocumen.body;
    var width = (0, $9c234cc6f4ed3b92$export$8960430cfd85939f)(html.scrollWidth, html.clientWidth, body ? body.scrollWidth : 0, body ? body.clientWidth : 0);
    var height = (0, $9c234cc6f4ed3b92$export$8960430cfd85939f)(html.scrollHeight, html.clientHeight, body ? body.scrollHeight : 0, body ? body.clientHeight : 0);
    var x = -winScroll.scrollLeft + (0, $c4822cfd8628a2db$export$2e2bcd8739ae039)(element);
    var y = -winScroll.scrollTop;
    if ((0, $7d74e8e12f28eb07$export$2e2bcd8739ae039)(body || html).direction === 'rtl') x += (0, $9c234cc6f4ed3b92$export$8960430cfd85939f)(html.clientWidth, body ? body.clientWidth : 0) - width;
    return {
        width: width,
        height: height,
        x: x,
        y: y
    };
}










function $6bd41d59857881c0$export$2e2bcd8739ae039(parent, child) {
    var rootNode = child.getRootNode && child.getRootNode(); // First, attempt with faster native method
    if (parent.contains(child)) return true;
    else if (rootNode && (0, $b83e73df1a45d10b$export$af51f0f06c0f328a)(rootNode)) {
        var next = child;
        do {
            if (next && parent.isSameNode(next)) return true;
             // $FlowFixMe[prop-missing]: need a better way to handle this...
            next = next.parentNode || next.host;
        }while (next);
    } // Give up, the result is false
    return false;
}



function $fe38ab5b3ef9f518$export$2e2bcd8739ae039(rect) {
    return Object.assign({}, rect, {
        left: rect.x,
        top: rect.y,
        right: rect.x + rect.width,
        bottom: rect.y + rect.height
    });
}



function $202ca78a1a2af2c0$var$getInnerBoundingClientRect(element, strategy) {
    var rect = (0, $d5b6692e6e485772$export$2e2bcd8739ae039)(element, false, strategy === 'fixed');
    rect.top = rect.top + element.clientTop;
    rect.left = rect.left + element.clientLeft;
    rect.bottom = rect.top + element.clientHeight;
    rect.right = rect.left + element.clientWidth;
    rect.width = element.clientWidth;
    rect.height = element.clientHeight;
    rect.x = rect.left;
    rect.y = rect.top;
    return rect;
}
function $202ca78a1a2af2c0$var$getClientRectFromMixedType(element, clippingParent, strategy) {
    return clippingParent === (0, $8478d2765df37982$export$d7b7311ec04a3e8f) ? (0, $fe38ab5b3ef9f518$export$2e2bcd8739ae039)((0, $b4831f59b241eaa3$export$2e2bcd8739ae039)(element, strategy)) : (0, $b83e73df1a45d10b$export$45a5e7f76e0caa8d)(clippingParent) ? $202ca78a1a2af2c0$var$getInnerBoundingClientRect(clippingParent, strategy) : (0, $fe38ab5b3ef9f518$export$2e2bcd8739ae039)((0, $f086f128b1379452$export$2e2bcd8739ae039)((0, $d56bbea5dd2fe00d$export$2e2bcd8739ae039)(element)));
} // A "clipping parent" is an overflowable container with the characteristic of
// clipping (or hiding) overflowing elements with a position different from
// `initial`
function $202ca78a1a2af2c0$var$getClippingParents(element) {
    var clippingParents = (0, $ebe9d77e8ccdd6fa$export$2e2bcd8739ae039)((0, $f019b932829cd8d0$export$2e2bcd8739ae039)(element));
    var canEscapeClipping = [
        'absolute',
        'fixed'
    ].indexOf((0, $7d74e8e12f28eb07$export$2e2bcd8739ae039)(element).position) >= 0;
    var clipperElement = canEscapeClipping && (0, $b83e73df1a45d10b$export$1b3bfaa9684536aa)(element) ? (0, $01dd27adea182f9a$export$2e2bcd8739ae039)(element) : element;
    if (!(0, $b83e73df1a45d10b$export$45a5e7f76e0caa8d)(clipperElement)) return [];
     // $FlowFixMe[incompatible-return]: https://github.com/facebook/flow/issues/1414
    return clippingParents.filter(function(clippingParent) {
        return (0, $b83e73df1a45d10b$export$45a5e7f76e0caa8d)(clippingParent) && (0, $6bd41d59857881c0$export$2e2bcd8739ae039)(clippingParent, clipperElement) && (0, $560fb2f7f26d74c6$export$2e2bcd8739ae039)(clippingParent) !== 'body';
    });
} // Gets the maximum area that the element is visible in due to any number of
function $202ca78a1a2af2c0$export$2e2bcd8739ae039(element, boundary, rootBoundary, strategy) {
    var mainClippingParents = boundary === 'clippingParents' ? $202ca78a1a2af2c0$var$getClippingParents(element) : [].concat(boundary);
    var clippingParents = [].concat(mainClippingParents, [
        rootBoundary
    ]);
    var firstClippingParent = clippingParents[0];
    var clippingRect = clippingParents.reduce(function(accRect, clippingParent) {
        var rect = $202ca78a1a2af2c0$var$getClientRectFromMixedType(element, clippingParent, strategy);
        accRect.top = (0, $9c234cc6f4ed3b92$export$8960430cfd85939f)(rect.top, accRect.top);
        accRect.right = (0, $9c234cc6f4ed3b92$export$96ec731ed4dcb222)(rect.right, accRect.right);
        accRect.bottom = (0, $9c234cc6f4ed3b92$export$96ec731ed4dcb222)(rect.bottom, accRect.bottom);
        accRect.left = (0, $9c234cc6f4ed3b92$export$8960430cfd85939f)(rect.left, accRect.left);
        return accRect;
    }, $202ca78a1a2af2c0$var$getClientRectFromMixedType(element, firstClippingParent, strategy));
    clippingRect.width = clippingRect.right - clippingRect.left;
    clippingRect.height = clippingRect.bottom - clippingRect.top;
    clippingRect.x = clippingRect.left;
    clippingRect.y = clippingRect.top;
    return clippingRect;
}








function $8e2ffbe2ceace94c$export$2e2bcd8739ae039() {
    return {
        top: 0,
        right: 0,
        bottom: 0,
        left: 0
    };
}


function $2cee3ab2c79c79df$export$2e2bcd8739ae039(paddingObject) {
    return Object.assign({}, (0, $8e2ffbe2ceace94c$export$2e2bcd8739ae039)(), paddingObject);
}


function $4aa2b0a0a983c60c$export$2e2bcd8739ae039(value, keys) {
    return keys.reduce(function(hashMap, key) {
        hashMap[key] = value;
        return hashMap;
    }, {});
}


function $8380aa231acdf82c$export$2e2bcd8739ae039(state, options) {
    if (options === void 0) options = {};
    var _options = options, _options$placement = _options.placement, placement = _options$placement === void 0 ? state.placement : _options$placement, _options$strategy = _options.strategy, strategy = _options$strategy === void 0 ? state.strategy : _options$strategy, _options$boundary = _options.boundary, boundary = _options$boundary === void 0 ? (0, $8478d2765df37982$export$390fd549c5303b4d) : _options$boundary, _options$rootBoundary = _options.rootBoundary, rootBoundary = _options$rootBoundary === void 0 ? (0, $8478d2765df37982$export$d7b7311ec04a3e8f) : _options$rootBoundary, _options$elementConte = _options.elementContext, elementContext = _options$elementConte === void 0 ? (0, $8478d2765df37982$export$ae5ab1c730825774) : _options$elementConte, _options$altBoundary = _options.altBoundary, altBoundary = _options$altBoundary === void 0 ? false : _options$altBoundary, _options$padding = _options.padding, padding = _options$padding === void 0 ? 0 : _options$padding;
    var paddingObject = (0, $2cee3ab2c79c79df$export$2e2bcd8739ae039)(typeof padding !== 'number' ? padding : (0, $4aa2b0a0a983c60c$export$2e2bcd8739ae039)(padding, (0, $8478d2765df37982$export$aec2ce47c367b8c3)));
    var altContext = elementContext === (0, $8478d2765df37982$export$ae5ab1c730825774) ? (0, $8478d2765df37982$export$ca50aac9f3ba507f) : (0, $8478d2765df37982$export$ae5ab1c730825774);
    var popperRect = state.rects.popper;
    var element = state.elements[altBoundary ? altContext : elementContext];
    var clippingClientRect = (0, $202ca78a1a2af2c0$export$2e2bcd8739ae039)((0, $b83e73df1a45d10b$export$45a5e7f76e0caa8d)(element) ? element : element.contextElement || (0, $d56bbea5dd2fe00d$export$2e2bcd8739ae039)(state.elements.popper), boundary, rootBoundary, strategy);
    var referenceClientRect = (0, $d5b6692e6e485772$export$2e2bcd8739ae039)(state.elements.reference);
    var popperOffsets = (0, $87655a20e47aab91$export$2e2bcd8739ae039)({
        reference: referenceClientRect,
        element: popperRect,
        strategy: 'absolute',
        placement: placement
    });
    var popperClientRect = (0, $fe38ab5b3ef9f518$export$2e2bcd8739ae039)(Object.assign({}, popperRect, popperOffsets));
    var elementClientRect = elementContext === (0, $8478d2765df37982$export$ae5ab1c730825774) ? popperClientRect : referenceClientRect; // positive = overflowing the clipping rect
    // 0 or negative = within the clipping rect
    var overflowOffsets = {
        top: clippingClientRect.top - elementClientRect.top + paddingObject.top,
        bottom: elementClientRect.bottom - clippingClientRect.bottom + paddingObject.bottom,
        left: clippingClientRect.left - elementClientRect.left + paddingObject.left,
        right: elementClientRect.right - clippingClientRect.right + paddingObject.right
    };
    var offsetData = state.modifiersData.offset; // Offsets can be applied only to the popper element
    if (elementContext === (0, $8478d2765df37982$export$ae5ab1c730825774) && offsetData) {
        var offset = offsetData[placement];
        Object.keys(overflowOffsets).forEach(function(key) {
            var multiply = [
                (0, $8478d2765df37982$export$79ffe56a765070d2),
                (0, $8478d2765df37982$export$40e543e69a8b3fbb)
            ].indexOf(key) >= 0 ? 1 : -1;
            var axis = [
                (0, $8478d2765df37982$export$1e95b668f3b82d),
                (0, $8478d2765df37982$export$40e543e69a8b3fbb)
            ].indexOf(key) >= 0 ? 'y' : 'x';
            overflowOffsets[key] += offset[axis] * multiply;
        });
    }
    return overflowOffsets;
}






function $bfa1c3bbb83e4d1a$export$2e2bcd8739ae039(state, options) {
    if (options === void 0) options = {};
    var _options = options, placement = _options.placement, boundary = _options.boundary, rootBoundary = _options.rootBoundary, padding = _options.padding, flipVariations = _options.flipVariations, _options$allowedAutoP = _options.allowedAutoPlacements, allowedAutoPlacements = _options$allowedAutoP === void 0 ? (0, $8478d2765df37982$export$803cd8101b6c182b) : _options$allowedAutoP;
    var variation = (0, $c5e5c390efc1bad4$export$2e2bcd8739ae039)(placement);
    var placements = variation ? flipVariations ? (0, $8478d2765df37982$export$368f9a87e87fa4e1) : (0, $8478d2765df37982$export$368f9a87e87fa4e1).filter(function(placement) {
        return (0, $c5e5c390efc1bad4$export$2e2bcd8739ae039)(placement) === variation;
    }) : (0, $8478d2765df37982$export$aec2ce47c367b8c3);
    var allowedPlacements = placements.filter(function(placement) {
        return allowedAutoPlacements.indexOf(placement) >= 0;
    });
    if (allowedPlacements.length === 0) allowedPlacements = placements;
     // $FlowFixMe[incompatible-type]: Flow seems to have problems with two array unions...
    var overflows = allowedPlacements.reduce(function(acc, placement) {
        acc[placement] = (0, $8380aa231acdf82c$export$2e2bcd8739ae039)(state, {
            placement: placement,
            boundary: boundary,
            rootBoundary: rootBoundary,
            padding: padding
        })[(0, $f33c93a011fffe7b$export$2e2bcd8739ae039)(placement)];
        return acc;
    }, {});
    return Object.keys(overflows).sort(function(a, b) {
        return overflows[a] - overflows[b];
    });
}




function $2c04a4fb35bcf3ba$var$getExpandedFallbackPlacements(placement) {
    if ((0, $f33c93a011fffe7b$export$2e2bcd8739ae039)(placement) === (0, $8478d2765df37982$export$dfb5619354ba860)) return [];
    var oppositePlacement = (0, $c0bd8d5905df73b7$export$2e2bcd8739ae039)(placement);
    return [
        (0, $a0a642892f75993d$export$2e2bcd8739ae039)(placement),
        oppositePlacement,
        (0, $a0a642892f75993d$export$2e2bcd8739ae039)(oppositePlacement)
    ];
}
function $2c04a4fb35bcf3ba$var$flip(_ref) {
    var state = _ref.state, options = _ref.options, name = _ref.name;
    if (state.modifiersData[name]._skip) return;
    var _options$mainAxis = options.mainAxis, checkMainAxis = _options$mainAxis === void 0 ? true : _options$mainAxis, _options$altAxis = options.altAxis, checkAltAxis = _options$altAxis === void 0 ? true : _options$altAxis, specifiedFallbackPlacements = options.fallbackPlacements, padding = options.padding, boundary = options.boundary, rootBoundary = options.rootBoundary, altBoundary = options.altBoundary, _options$flipVariatio = options.flipVariations, flipVariations = _options$flipVariatio === void 0 ? true : _options$flipVariatio, allowedAutoPlacements = options.allowedAutoPlacements;
    var preferredPlacement = state.options.placement;
    var basePlacement = (0, $f33c93a011fffe7b$export$2e2bcd8739ae039)(preferredPlacement);
    var isBasePlacement = basePlacement === preferredPlacement;
    var fallbackPlacements = specifiedFallbackPlacements || (isBasePlacement || !flipVariations ? [
        (0, $c0bd8d5905df73b7$export$2e2bcd8739ae039)(preferredPlacement)
    ] : $2c04a4fb35bcf3ba$var$getExpandedFallbackPlacements(preferredPlacement));
    var placements = [
        preferredPlacement
    ].concat(fallbackPlacements).reduce(function(acc, placement) {
        return acc.concat((0, $f33c93a011fffe7b$export$2e2bcd8739ae039)(placement) === (0, $8478d2765df37982$export$dfb5619354ba860) ? (0, $bfa1c3bbb83e4d1a$export$2e2bcd8739ae039)(state, {
            placement: placement,
            boundary: boundary,
            rootBoundary: rootBoundary,
            padding: padding,
            flipVariations: flipVariations,
            allowedAutoPlacements: allowedAutoPlacements
        }) : placement);
    }, []);
    var referenceRect = state.rects.reference;
    var popperRect = state.rects.popper;
    var checksMap = new Map();
    var makeFallbackChecks = true;
    var firstFittingPlacement = placements[0];
    for(var i = 0; i < placements.length; i++){
        var placement = placements[i];
        var _basePlacement = (0, $f33c93a011fffe7b$export$2e2bcd8739ae039)(placement);
        var isStartVariation = (0, $c5e5c390efc1bad4$export$2e2bcd8739ae039)(placement) === (0, $8478d2765df37982$export$b3571188c770cc5a);
        var isVertical = [
            (0, $8478d2765df37982$export$1e95b668f3b82d),
            (0, $8478d2765df37982$export$40e543e69a8b3fbb)
        ].indexOf(_basePlacement) >= 0;
        var len = isVertical ? 'width' : 'height';
        var overflow = (0, $8380aa231acdf82c$export$2e2bcd8739ae039)(state, {
            placement: placement,
            boundary: boundary,
            rootBoundary: rootBoundary,
            altBoundary: altBoundary,
            padding: padding
        });
        var mainVariationSide = isVertical ? isStartVariation ? (0, $8478d2765df37982$export$79ffe56a765070d2) : (0, $8478d2765df37982$export$eabcd2c8791e7bf4) : isStartVariation ? (0, $8478d2765df37982$export$40e543e69a8b3fbb) : (0, $8478d2765df37982$export$1e95b668f3b82d);
        if (referenceRect[len] > popperRect[len]) mainVariationSide = (0, $c0bd8d5905df73b7$export$2e2bcd8739ae039)(mainVariationSide);
        var altVariationSide = (0, $c0bd8d5905df73b7$export$2e2bcd8739ae039)(mainVariationSide);
        var checks = [];
        if (checkMainAxis) checks.push(overflow[_basePlacement] <= 0);
        if (checkAltAxis) checks.push(overflow[mainVariationSide] <= 0, overflow[altVariationSide] <= 0);
        if (checks.every(function(check) {
            return check;
        })) {
            firstFittingPlacement = placement;
            makeFallbackChecks = false;
            break;
        }
        checksMap.set(placement, checks);
    }
    if (makeFallbackChecks) {
        // `2` may be desired in some cases – research later
        var numberOfChecks = flipVariations ? 3 : 1;
        var _loop = function _loop(_i) {
            var fittingPlacement = placements.find(function(placement) {
                var checks = checksMap.get(placement);
                if (checks) return checks.slice(0, _i).every(function(check) {
                    return check;
                });
            });
            if (fittingPlacement) {
                firstFittingPlacement = fittingPlacement;
                return "break";
            }
        };
        for(var _i = numberOfChecks; _i > 0; _i--){
            var _ret = _loop(_i);
            if (_ret === "break") break;
        }
    }
    if (state.placement !== firstFittingPlacement) {
        state.modifiersData[name]._skip = true;
        state.placement = firstFittingPlacement;
        state.reset = true;
    }
} // eslint-disable-next-line import/no-unused-modules
var $2c04a4fb35bcf3ba$export$2e2bcd8739ae039 = {
    name: 'flip',
    enabled: true,
    phase: 'main',
    fn: $2c04a4fb35bcf3ba$var$flip,
    requiresIfExists: [
        'offset'
    ],
    data: {
        _skip: false
    }
};





function $64f3a18c323e0caa$export$2e2bcd8739ae039(axis) {
    return axis === 'x' ? 'y' : 'x';
}



function $77ead1178fd8688c$export$f28d906d67a997f3(min, value, max) {
    return (0, $9c234cc6f4ed3b92$export$8960430cfd85939f)(min, (0, $9c234cc6f4ed3b92$export$96ec731ed4dcb222)(value, max));
}
function $77ead1178fd8688c$export$86c8af6d3ef0b4a(min, value, max) {
    var v = $77ead1178fd8688c$export$f28d906d67a997f3(min, value, max);
    return v > max ? max : v;
}








function $135f2e2733b618c2$var$preventOverflow(_ref) {
    var state = _ref.state, options = _ref.options, name = _ref.name;
    var _options$mainAxis = options.mainAxis, checkMainAxis = _options$mainAxis === void 0 ? true : _options$mainAxis, _options$altAxis = options.altAxis, checkAltAxis = _options$altAxis === void 0 ? false : _options$altAxis, boundary = options.boundary, rootBoundary = options.rootBoundary, altBoundary = options.altBoundary, padding = options.padding, _options$tether = options.tether, tether = _options$tether === void 0 ? true : _options$tether, _options$tetherOffset = options.tetherOffset, tetherOffset = _options$tetherOffset === void 0 ? 0 : _options$tetherOffset;
    var overflow = (0, $8380aa231acdf82c$export$2e2bcd8739ae039)(state, {
        boundary: boundary,
        rootBoundary: rootBoundary,
        padding: padding,
        altBoundary: altBoundary
    });
    var basePlacement = (0, $f33c93a011fffe7b$export$2e2bcd8739ae039)(state.placement);
    var variation = (0, $c5e5c390efc1bad4$export$2e2bcd8739ae039)(state.placement);
    var isBasePlacement = !variation;
    var mainAxis = (0, $247260022a0ad6ee$export$2e2bcd8739ae039)(basePlacement);
    var altAxis = (0, $64f3a18c323e0caa$export$2e2bcd8739ae039)(mainAxis);
    var popperOffsets = state.modifiersData.popperOffsets;
    var referenceRect = state.rects.reference;
    var popperRect = state.rects.popper;
    var tetherOffsetValue = typeof tetherOffset === 'function' ? tetherOffset(Object.assign({}, state.rects, {
        placement: state.placement
    })) : tetherOffset;
    var normalizedTetherOffsetValue = typeof tetherOffsetValue === 'number' ? {
        mainAxis: tetherOffsetValue,
        altAxis: tetherOffsetValue
    } : Object.assign({
        mainAxis: 0,
        altAxis: 0
    }, tetherOffsetValue);
    var offsetModifierState = state.modifiersData.offset ? state.modifiersData.offset[state.placement] : null;
    var data = {
        x: 0,
        y: 0
    };
    if (!popperOffsets) return;
    if (checkMainAxis) {
        var _offsetModifierState$;
        var mainSide = mainAxis === 'y' ? (0, $8478d2765df37982$export$1e95b668f3b82d) : (0, $8478d2765df37982$export$eabcd2c8791e7bf4);
        var altSide = mainAxis === 'y' ? (0, $8478d2765df37982$export$40e543e69a8b3fbb) : (0, $8478d2765df37982$export$79ffe56a765070d2);
        var len = mainAxis === 'y' ? 'height' : 'width';
        var offset = popperOffsets[mainAxis];
        var min = offset + overflow[mainSide];
        var max = offset - overflow[altSide];
        var additive = tether ? -popperRect[len] / 2 : 0;
        var minLen = variation === (0, $8478d2765df37982$export$b3571188c770cc5a) ? referenceRect[len] : popperRect[len];
        var maxLen = variation === (0, $8478d2765df37982$export$b3571188c770cc5a) ? -popperRect[len] : -referenceRect[len]; // We need to include the arrow in the calculation so the arrow doesn't go
        // outside the reference bounds
        var arrowElement = state.elements.arrow;
        var arrowRect = tether && arrowElement ? (0, $f4ac90f91f825ce1$export$2e2bcd8739ae039)(arrowElement) : {
            width: 0,
            height: 0
        };
        var arrowPaddingObject = state.modifiersData['arrow#persistent'] ? state.modifiersData['arrow#persistent'].padding : (0, $8e2ffbe2ceace94c$export$2e2bcd8739ae039)();
        var arrowPaddingMin = arrowPaddingObject[mainSide];
        var arrowPaddingMax = arrowPaddingObject[altSide]; // If the reference length is smaller than the arrow length, we don't want
        // to include its full size in the calculation. If the reference is small
        // and near the edge of a boundary, the popper can overflow even if the
        // reference is not overflowing as well (e.g. virtual elements with no
        // width or height)
        var arrowLen = (0, $77ead1178fd8688c$export$f28d906d67a997f3)(0, referenceRect[len], arrowRect[len]);
        var minOffset = isBasePlacement ? referenceRect[len] / 2 - additive - arrowLen - arrowPaddingMin - normalizedTetherOffsetValue.mainAxis : minLen - arrowLen - arrowPaddingMin - normalizedTetherOffsetValue.mainAxis;
        var maxOffset = isBasePlacement ? -referenceRect[len] / 2 + additive + arrowLen + arrowPaddingMax + normalizedTetherOffsetValue.mainAxis : maxLen + arrowLen + arrowPaddingMax + normalizedTetherOffsetValue.mainAxis;
        var arrowOffsetParent = state.elements.arrow && (0, $01dd27adea182f9a$export$2e2bcd8739ae039)(state.elements.arrow);
        var clientOffset = arrowOffsetParent ? mainAxis === 'y' ? arrowOffsetParent.clientTop || 0 : arrowOffsetParent.clientLeft || 0 : 0;
        var offsetModifierValue = (_offsetModifierState$ = offsetModifierState == null ? void 0 : offsetModifierState[mainAxis]) != null ? _offsetModifierState$ : 0;
        var tetherMin = offset + minOffset - offsetModifierValue - clientOffset;
        var tetherMax = offset + maxOffset - offsetModifierValue;
        var preventedOffset = (0, $77ead1178fd8688c$export$f28d906d67a997f3)(tether ? (0, $9c234cc6f4ed3b92$export$96ec731ed4dcb222)(min, tetherMin) : min, offset, tether ? (0, $9c234cc6f4ed3b92$export$8960430cfd85939f)(max, tetherMax) : max);
        popperOffsets[mainAxis] = preventedOffset;
        data[mainAxis] = preventedOffset - offset;
    }
    if (checkAltAxis) {
        var _offsetModifierState$2;
        var _mainSide = mainAxis === 'x' ? (0, $8478d2765df37982$export$1e95b668f3b82d) : (0, $8478d2765df37982$export$eabcd2c8791e7bf4);
        var _altSide = mainAxis === 'x' ? (0, $8478d2765df37982$export$40e543e69a8b3fbb) : (0, $8478d2765df37982$export$79ffe56a765070d2);
        var _offset = popperOffsets[altAxis];
        var _len = altAxis === 'y' ? 'height' : 'width';
        var _min = _offset + overflow[_mainSide];
        var _max = _offset - overflow[_altSide];
        var isOriginSide = [
            (0, $8478d2765df37982$export$1e95b668f3b82d),
            (0, $8478d2765df37982$export$eabcd2c8791e7bf4)
        ].indexOf(basePlacement) !== -1;
        var _offsetModifierValue = (_offsetModifierState$2 = offsetModifierState == null ? void 0 : offsetModifierState[altAxis]) != null ? _offsetModifierState$2 : 0;
        var _tetherMin = isOriginSide ? _min : _offset - referenceRect[_len] - popperRect[_len] - _offsetModifierValue + normalizedTetherOffsetValue.altAxis;
        var _tetherMax = isOriginSide ? _offset + referenceRect[_len] + popperRect[_len] - _offsetModifierValue - normalizedTetherOffsetValue.altAxis : _max;
        var _preventedOffset = tether && isOriginSide ? (0, $77ead1178fd8688c$export$86c8af6d3ef0b4a)(_tetherMin, _offset, _tetherMax) : (0, $77ead1178fd8688c$export$f28d906d67a997f3)(tether ? _tetherMin : _min, _offset, tether ? _tetherMax : _max);
        popperOffsets[altAxis] = _preventedOffset;
        data[altAxis] = _preventedOffset - _offset;
    }
    state.modifiersData[name] = data;
} // eslint-disable-next-line import/no-unused-modules
var $135f2e2733b618c2$export$2e2bcd8739ae039 = {
    name: 'preventOverflow',
    enabled: true,
    phase: 'main',
    fn: $135f2e2733b618c2$var$preventOverflow,
    requiresIfExists: [
        'offset'
    ]
};











var $19cc820fe7335f18$var$toPaddingObject = function toPaddingObject(padding, state) {
    padding = typeof padding === 'function' ? padding(Object.assign({}, state.rects, {
        placement: state.placement
    })) : padding;
    return (0, $2cee3ab2c79c79df$export$2e2bcd8739ae039)(typeof padding !== 'number' ? padding : (0, $4aa2b0a0a983c60c$export$2e2bcd8739ae039)(padding, (0, $8478d2765df37982$export$aec2ce47c367b8c3)));
};
function $19cc820fe7335f18$var$arrow(_ref) {
    var _state$modifiersData$;
    var state = _ref.state, name = _ref.name, options = _ref.options;
    var arrowElement = state.elements.arrow;
    var popperOffsets = state.modifiersData.popperOffsets;
    var basePlacement = (0, $f33c93a011fffe7b$export$2e2bcd8739ae039)(state.placement);
    var axis = (0, $247260022a0ad6ee$export$2e2bcd8739ae039)(basePlacement);
    var isVertical = [
        (0, $8478d2765df37982$export$eabcd2c8791e7bf4),
        (0, $8478d2765df37982$export$79ffe56a765070d2)
    ].indexOf(basePlacement) >= 0;
    var len = isVertical ? 'height' : 'width';
    if (!arrowElement || !popperOffsets) return;
    var paddingObject = $19cc820fe7335f18$var$toPaddingObject(options.padding, state);
    var arrowRect = (0, $f4ac90f91f825ce1$export$2e2bcd8739ae039)(arrowElement);
    var minProp = axis === 'y' ? (0, $8478d2765df37982$export$1e95b668f3b82d) : (0, $8478d2765df37982$export$eabcd2c8791e7bf4);
    var maxProp = axis === 'y' ? (0, $8478d2765df37982$export$40e543e69a8b3fbb) : (0, $8478d2765df37982$export$79ffe56a765070d2);
    var endDiff = state.rects.reference[len] + state.rects.reference[axis] - popperOffsets[axis] - state.rects.popper[len];
    var startDiff = popperOffsets[axis] - state.rects.reference[axis];
    var arrowOffsetParent = (0, $01dd27adea182f9a$export$2e2bcd8739ae039)(arrowElement);
    var clientSize = arrowOffsetParent ? axis === 'y' ? arrowOffsetParent.clientHeight || 0 : arrowOffsetParent.clientWidth || 0 : 0;
    var centerToReference = endDiff / 2 - startDiff / 2; // Make sure the arrow doesn't overflow the popper if the center point is
    // outside of the popper bounds
    var min = paddingObject[minProp];
    var max = clientSize - arrowRect[len] - paddingObject[maxProp];
    var center = clientSize / 2 - arrowRect[len] / 2 + centerToReference;
    var offset = (0, $77ead1178fd8688c$export$f28d906d67a997f3)(min, center, max); // Prevents breaking syntax highlighting...
    var axisProp = axis;
    state.modifiersData[name] = (_state$modifiersData$ = {}, _state$modifiersData$[axisProp] = offset, _state$modifiersData$.centerOffset = offset - center, _state$modifiersData$);
}
function $19cc820fe7335f18$var$effect(_ref2) {
    var state = _ref2.state, options = _ref2.options;
    var _options$element = options.element, arrowElement = _options$element === void 0 ? '[data-popper-arrow]' : _options$element;
    if (arrowElement == null) return;
     // CSS selector
    if (typeof arrowElement === 'string') {
        arrowElement = state.elements.popper.querySelector(arrowElement);
        if (!arrowElement) return;
    }
    if (!(0, $6bd41d59857881c0$export$2e2bcd8739ae039)(state.elements.popper, arrowElement)) return;
    state.elements.arrow = arrowElement;
} // eslint-disable-next-line import/no-unused-modules
var $19cc820fe7335f18$export$2e2bcd8739ae039 = {
    name: 'arrow',
    enabled: true,
    phase: 'main',
    fn: $19cc820fe7335f18$var$arrow,
    effect: $19cc820fe7335f18$var$effect,
    requires: [
        'popperOffsets'
    ],
    requiresIfExists: [
        'preventOverflow'
    ]
};




function $1f5fc639c8bd03d5$var$getSideOffsets(overflow, rect, preventedOffsets) {
    if (preventedOffsets === void 0) preventedOffsets = {
        x: 0,
        y: 0
    };
    return {
        top: overflow.top - rect.height - preventedOffsets.y,
        right: overflow.right - rect.width + preventedOffsets.x,
        bottom: overflow.bottom - rect.height + preventedOffsets.y,
        left: overflow.left - rect.width - preventedOffsets.x
    };
}
function $1f5fc639c8bd03d5$var$isAnySideFullyClipped(overflow) {
    return [
        (0, $8478d2765df37982$export$1e95b668f3b82d),
        (0, $8478d2765df37982$export$79ffe56a765070d2),
        (0, $8478d2765df37982$export$40e543e69a8b3fbb),
        (0, $8478d2765df37982$export$eabcd2c8791e7bf4)
    ].some(function(side) {
        return overflow[side] >= 0;
    });
}
function $1f5fc639c8bd03d5$var$hide(_ref) {
    var state = _ref.state, name = _ref.name;
    var referenceRect = state.rects.reference;
    var popperRect = state.rects.popper;
    var preventedOffsets = state.modifiersData.preventOverflow;
    var referenceOverflow = (0, $8380aa231acdf82c$export$2e2bcd8739ae039)(state, {
        elementContext: 'reference'
    });
    var popperAltOverflow = (0, $8380aa231acdf82c$export$2e2bcd8739ae039)(state, {
        altBoundary: true
    });
    var referenceClippingOffsets = $1f5fc639c8bd03d5$var$getSideOffsets(referenceOverflow, referenceRect);
    var popperEscapeOffsets = $1f5fc639c8bd03d5$var$getSideOffsets(popperAltOverflow, popperRect, preventedOffsets);
    var isReferenceHidden = $1f5fc639c8bd03d5$var$isAnySideFullyClipped(referenceClippingOffsets);
    var hasPopperEscaped = $1f5fc639c8bd03d5$var$isAnySideFullyClipped(popperEscapeOffsets);
    state.modifiersData[name] = {
        referenceClippingOffsets: referenceClippingOffsets,
        popperEscapeOffsets: popperEscapeOffsets,
        isReferenceHidden: isReferenceHidden,
        hasPopperEscaped: hasPopperEscaped
    };
    state.attributes.popper = Object.assign({}, state.attributes.popper, {
        'data-popper-reference-hidden': isReferenceHidden,
        'data-popper-escaped': hasPopperEscaped
    });
} // eslint-disable-next-line import/no-unused-modules
var $1f5fc639c8bd03d5$export$2e2bcd8739ae039 = {
    name: 'hide',
    enabled: true,
    phase: 'main',
    requiresIfExists: [
        'preventOverflow'
    ],
    fn: $1f5fc639c8bd03d5$var$hide
};




var $725309ecee3adaa0$export$d34966752335dd47 = [
    (0, $48430453eb0d053b$export$2e2bcd8739ae039),
    (0, $b19a98045ea3e977$export$2e2bcd8739ae039),
    (0, $c5e18f9acfdbb8a2$export$2e2bcd8739ae039),
    (0, $1c7a027d393b751c$export$2e2bcd8739ae039),
    (0, $f0fe73ce831ce17b$export$2e2bcd8739ae039),
    (0, $2c04a4fb35bcf3ba$export$2e2bcd8739ae039),
    (0, $135f2e2733b618c2$export$2e2bcd8739ae039),
    (0, $19cc820fe7335f18$export$2e2bcd8739ae039),
    (0, $1f5fc639c8bd03d5$export$2e2bcd8739ae039)
];
var $725309ecee3adaa0$export$8f7491d57c8f97a9 = /*#__PURE__*/ (0, $f3aacc53d895d057$export$ed5e13716264f202)({
    defaultModifiers: $725309ecee3adaa0$export$d34966752335dd47
}); // eslint-disable-next-line import/no-unused-modules



var $ece16b3047178a9c$var$__assign = undefined && undefined.__assign || function() {
    $ece16b3047178a9c$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $ece16b3047178a9c$var$__assign.apply(this, arguments);
};
var $ece16b3047178a9c$var$__spreadArray = undefined && undefined.__spreadArray || function(to, from, pack) {
    if (pack || arguments.length === 2) {
        for(var i = 0, l = from.length, ar; i < l; i++)if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
};
var $ece16b3047178a9c$var$Default = {
    placement: 'bottom',
    triggerType: 'click',
    offsetSkidding: 0,
    offsetDistance: 10,
    delay: 300,
    ignoreClickOutsideClass: false,
    onShow: function() {},
    onHide: function() {},
    onToggle: function() {}
};
var $ece16b3047178a9c$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $ece16b3047178a9c$var$Dropdown = /** @class */ function() {
    function Dropdown(targetElement, triggerElement, options, instanceOptions) {
        if (targetElement === void 0) targetElement = null;
        if (triggerElement === void 0) triggerElement = null;
        if (options === void 0) options = $ece16b3047178a9c$var$Default;
        if (instanceOptions === void 0) instanceOptions = $ece16b3047178a9c$var$DefaultInstanceOptions;
        this._instanceId = instanceOptions.id ? instanceOptions.id : targetElement.id;
        this._targetEl = targetElement;
        this._triggerEl = triggerElement;
        this._options = $ece16b3047178a9c$var$__assign($ece16b3047178a9c$var$__assign({}, $ece16b3047178a9c$var$Default), options);
        this._popperInstance = null;
        this._visible = false;
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('Dropdown', this, this._instanceId, instanceOptions.override);
    }
    Dropdown.prototype.init = function() {
        if (this._triggerEl && this._targetEl && !this._initialized) {
            this._popperInstance = this._createPopperInstance();
            this._setupEventListeners();
            this._initialized = true;
        }
    };
    Dropdown.prototype.destroy = function() {
        var _this = this;
        var triggerEvents = this._getTriggerEvents();
        // Remove click event listeners for trigger element
        if (this._options.triggerType === 'click') triggerEvents.showEvents.forEach(function(ev) {
            _this._triggerEl.removeEventListener(ev, _this._clickHandler);
        });
        // Remove hover event listeners for trigger and target elements
        if (this._options.triggerType === 'hover') {
            triggerEvents.showEvents.forEach(function(ev) {
                _this._triggerEl.removeEventListener(ev, _this._hoverShowTriggerElHandler);
                _this._targetEl.removeEventListener(ev, _this._hoverShowTargetElHandler);
            });
            triggerEvents.hideEvents.forEach(function(ev) {
                _this._triggerEl.removeEventListener(ev, _this._hoverHideHandler);
                _this._targetEl.removeEventListener(ev, _this._hoverHideHandler);
            });
        }
        this._popperInstance.destroy();
        this._initialized = false;
    };
    Dropdown.prototype.removeInstance = function() {
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('Dropdown', this._instanceId);
    };
    Dropdown.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    Dropdown.prototype._setupEventListeners = function() {
        var _this = this;
        var triggerEvents = this._getTriggerEvents();
        this._clickHandler = function() {
            _this.toggle();
        };
        // click event handling for trigger element
        if (this._options.triggerType === 'click') triggerEvents.showEvents.forEach(function(ev) {
            _this._triggerEl.addEventListener(ev, _this._clickHandler);
        });
        this._hoverShowTriggerElHandler = function(ev) {
            if (ev.type === 'click') _this.toggle();
            else setTimeout(function() {
                _this.show();
            }, _this._options.delay);
        };
        this._hoverShowTargetElHandler = function() {
            _this.show();
        };
        this._hoverHideHandler = function() {
            setTimeout(function() {
                if (!_this._targetEl.matches(':hover')) _this.hide();
            }, _this._options.delay);
        };
        // hover event handling for trigger element
        if (this._options.triggerType === 'hover') {
            triggerEvents.showEvents.forEach(function(ev) {
                _this._triggerEl.addEventListener(ev, _this._hoverShowTriggerElHandler);
                _this._targetEl.addEventListener(ev, _this._hoverShowTargetElHandler);
            });
            triggerEvents.hideEvents.forEach(function(ev) {
                _this._triggerEl.addEventListener(ev, _this._hoverHideHandler);
                _this._targetEl.addEventListener(ev, _this._hoverHideHandler);
            });
        }
    };
    Dropdown.prototype._createPopperInstance = function() {
        return (0, $725309ecee3adaa0$export$8f7491d57c8f97a9)(this._triggerEl, this._targetEl, {
            placement: this._options.placement,
            modifiers: [
                {
                    name: 'offset',
                    options: {
                        offset: [
                            this._options.offsetSkidding,
                            this._options.offsetDistance
                        ]
                    }
                }
            ]
        });
    };
    Dropdown.prototype._setupClickOutsideListener = function() {
        var _this = this;
        this._clickOutsideEventListener = function(ev) {
            _this._handleClickOutside(ev, _this._targetEl);
        };
        document.body.addEventListener('click', this._clickOutsideEventListener, true);
    };
    Dropdown.prototype._removeClickOutsideListener = function() {
        document.body.removeEventListener('click', this._clickOutsideEventListener, true);
    };
    Dropdown.prototype._handleClickOutside = function(ev, targetEl) {
        var clickedEl = ev.target;
        // Ignore clicks on the trigger element (ie. a datepicker input)
        var ignoreClickOutsideClass = this._options.ignoreClickOutsideClass;
        var isIgnored = false;
        if (ignoreClickOutsideClass) {
            var ignoredClickOutsideEls = document.querySelectorAll(".".concat(ignoreClickOutsideClass));
            ignoredClickOutsideEls.forEach(function(el) {
                if (el.contains(clickedEl)) {
                    isIgnored = true;
                    return;
                }
            });
        }
        // Ignore clicks on the target element (ie. dropdown itself)
        if (clickedEl !== targetEl && !targetEl.contains(clickedEl) && !this._triggerEl.contains(clickedEl) && !isIgnored && this.isVisible()) this.hide();
    };
    Dropdown.prototype._getTriggerEvents = function() {
        switch(this._options.triggerType){
            case 'hover':
                return {
                    showEvents: [
                        'mouseenter',
                        'click'
                    ],
                    hideEvents: [
                        'mouseleave'
                    ]
                };
            case 'click':
                return {
                    showEvents: [
                        'click'
                    ],
                    hideEvents: []
                };
            case 'none':
                return {
                    showEvents: [],
                    hideEvents: []
                };
            default:
                return {
                    showEvents: [
                        'click'
                    ],
                    hideEvents: []
                };
        }
    };
    Dropdown.prototype.toggle = function() {
        if (this.isVisible()) this.hide();
        else this.show();
        this._options.onToggle(this);
    };
    Dropdown.prototype.isVisible = function() {
        return this._visible;
    };
    Dropdown.prototype.show = function() {
        this._targetEl.classList.remove('hidden');
        this._targetEl.classList.add('block');
        this._targetEl.removeAttribute('aria-hidden');
        // Enable the event listeners
        this._popperInstance.setOptions(function(options) {
            return $ece16b3047178a9c$var$__assign($ece16b3047178a9c$var$__assign({}, options), {
                modifiers: $ece16b3047178a9c$var$__spreadArray($ece16b3047178a9c$var$__spreadArray([], options.modifiers, true), [
                    {
                        name: 'eventListeners',
                        enabled: true
                    }
                ], false)
            });
        });
        this._setupClickOutsideListener();
        // Update its position
        this._popperInstance.update();
        this._visible = true;
        // callback function
        this._options.onShow(this);
    };
    Dropdown.prototype.hide = function() {
        this._targetEl.classList.remove('block');
        this._targetEl.classList.add('hidden');
        this._targetEl.setAttribute('aria-hidden', 'true');
        // Disable the event listeners
        this._popperInstance.setOptions(function(options) {
            return $ece16b3047178a9c$var$__assign($ece16b3047178a9c$var$__assign({}, options), {
                modifiers: $ece16b3047178a9c$var$__spreadArray($ece16b3047178a9c$var$__spreadArray([], options.modifiers, true), [
                    {
                        name: 'eventListeners',
                        enabled: false
                    }
                ], false)
            });
        });
        this._visible = false;
        this._removeClickOutsideListener();
        // callback function
        this._options.onHide(this);
    };
    Dropdown.prototype.updateOnShow = function(callback) {
        this._options.onShow = callback;
    };
    Dropdown.prototype.updateOnHide = function(callback) {
        this._options.onHide = callback;
    };
    Dropdown.prototype.updateOnToggle = function(callback) {
        this._options.onToggle = callback;
    };
    return Dropdown;
}();
function $ece16b3047178a9c$export$8cb65a02593bf108() {
    document.querySelectorAll('[data-dropdown-toggle]').forEach(function($triggerEl) {
        var dropdownId = $triggerEl.getAttribute('data-dropdown-toggle');
        var $dropdownEl = document.getElementById(dropdownId);
        if ($dropdownEl) {
            var placement = $triggerEl.getAttribute('data-dropdown-placement');
            var offsetSkidding = $triggerEl.getAttribute('data-dropdown-offset-skidding');
            var offsetDistance = $triggerEl.getAttribute('data-dropdown-offset-distance');
            var triggerType = $triggerEl.getAttribute('data-dropdown-trigger');
            var delay = $triggerEl.getAttribute('data-dropdown-delay');
            var ignoreClickOutsideClass = $triggerEl.getAttribute('data-dropdown-ignore-click-outside-class');
            new $ece16b3047178a9c$var$Dropdown($dropdownEl, $triggerEl, {
                placement: placement ? placement : $ece16b3047178a9c$var$Default.placement,
                triggerType: triggerType ? triggerType : $ece16b3047178a9c$var$Default.triggerType,
                offsetSkidding: offsetSkidding ? parseInt(offsetSkidding) : $ece16b3047178a9c$var$Default.offsetSkidding,
                offsetDistance: offsetDistance ? parseInt(offsetDistance) : $ece16b3047178a9c$var$Default.offsetDistance,
                delay: delay ? parseInt(delay) : $ece16b3047178a9c$var$Default.delay,
                ignoreClickOutsideClass: ignoreClickOutsideClass ? ignoreClickOutsideClass : $ece16b3047178a9c$var$Default.ignoreClickOutsideClass
            });
        } else console.error("The dropdown element with id \"".concat(dropdownId, "\" does not exist. Please check the data-dropdown-toggle attribute."));
    });
}
if (typeof window !== 'undefined') {
    window.Dropdown = $ece16b3047178a9c$var$Dropdown;
    window.initDropdowns = $ece16b3047178a9c$export$8cb65a02593bf108;
}
var $ece16b3047178a9c$export$2e2bcd8739ae039 = $ece16b3047178a9c$var$Dropdown;



var $8d24653c8cdd4795$var$__assign = undefined && undefined.__assign || function() {
    $8d24653c8cdd4795$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $8d24653c8cdd4795$var$__assign.apply(this, arguments);
};
var $8d24653c8cdd4795$var$Default = {
    placement: 'center',
    backdropClasses: 'bg-dark-backdrop/70 fixed inset-0 z-40',
    backdrop: 'dynamic',
    closable: true,
    onHide: function() {},
    onShow: function() {},
    onToggle: function() {}
};
var $8d24653c8cdd4795$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $8d24653c8cdd4795$var$Modal = /** @class */ function() {
    function Modal(targetEl, options, instanceOptions) {
        if (targetEl === void 0) targetEl = null;
        if (options === void 0) options = $8d24653c8cdd4795$var$Default;
        if (instanceOptions === void 0) instanceOptions = $8d24653c8cdd4795$var$DefaultInstanceOptions;
        this._eventListenerInstances = [];
        this._instanceId = instanceOptions.id ? instanceOptions.id : targetEl.id;
        this._targetEl = targetEl;
        this._options = $8d24653c8cdd4795$var$__assign($8d24653c8cdd4795$var$__assign({}, $8d24653c8cdd4795$var$Default), options);
        this._isHidden = true;
        this._backdropEl = null;
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('Modal', this, this._instanceId, instanceOptions.override);
    }
    Modal.prototype.init = function() {
        var _this = this;
        if (this._targetEl && !this._initialized) {
            this._getPlacementClasses().map(function(c) {
                _this._targetEl.classList.add(c);
            });
            this._initialized = true;
        }
    };
    Modal.prototype.destroy = function() {
        if (this._initialized) {
            this.removeAllEventListenerInstances();
            this._destroyBackdropEl();
            this._initialized = false;
        }
    };
    Modal.prototype.removeInstance = function() {
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('Modal', this._instanceId);
    };
    Modal.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    Modal.prototype._createBackdrop = function() {
        var _a;
        if (this._isHidden) {
            var backdropEl = document.createElement('div');
            (_a = backdropEl.classList).add.apply(_a, this._options.backdropClasses.split(' '));
            document.querySelector('body').append(backdropEl);
            this._backdropEl = backdropEl;
        }
    };
    Modal.prototype._destroyBackdropEl = function() {
        if (!this._isHidden && this._backdropEl) {
            this._backdropEl.remove();
            this._backdropEl = null;
        }
    };
    Modal.prototype._setupModalCloseEventListeners = function() {
        var _this = this;
        if (this._options.backdrop === 'dynamic') {
            this._clickOutsideEventListener = function(ev) {
                _this._handleOutsideClick(ev.target);
            };
            this._targetEl.addEventListener('click', this._clickOutsideEventListener, true);
        }
        this._keydownEventListener = function(ev) {
            if (ev.key === 'Escape') _this.hide();
        };
        document.body.addEventListener('keydown', this._keydownEventListener, true);
    };
    Modal.prototype._removeModalCloseEventListeners = function() {
        if (this._options.backdrop === 'dynamic') this._targetEl.removeEventListener('click', this._clickOutsideEventListener, true);
        document.body.removeEventListener('keydown', this._keydownEventListener, true);
    };
    Modal.prototype._handleOutsideClick = function(target) {
        if (target === this._targetEl || target === this._backdropEl && this.isVisible()) this.hide();
    };
    Modal.prototype._getPlacementClasses = function() {
        switch(this._options.placement){
            // top
            case 'top-left':
                return [
                    'justify-start',
                    'items-start'
                ];
            case 'top-center':
                return [
                    'justify-center',
                    'items-start'
                ];
            case 'top-right':
                return [
                    'justify-end',
                    'items-start'
                ];
            // center
            case 'center-left':
                return [
                    'justify-start',
                    'items-center'
                ];
            case 'center':
                return [
                    'justify-center',
                    'items-center'
                ];
            case 'center-right':
                return [
                    'justify-end',
                    'items-center'
                ];
            // bottom
            case 'bottom-left':
                return [
                    'justify-start',
                    'items-end'
                ];
            case 'bottom-center':
                return [
                    'justify-center',
                    'items-end'
                ];
            case 'bottom-right':
                return [
                    'justify-end',
                    'items-end'
                ];
            default:
                return [
                    'justify-center',
                    'items-center'
                ];
        }
    };
    Modal.prototype.toggle = function() {
        if (this._isHidden) this.show();
        else this.hide();
        // callback function
        this._options.onToggle(this);
    };
    Modal.prototype.show = function() {
        if (this.isHidden) {
            this._targetEl.classList.add('flex');
            this._targetEl.classList.remove('hidden');
            this._targetEl.setAttribute('aria-modal', 'true');
            this._targetEl.setAttribute('role', 'dialog');
            this._targetEl.removeAttribute('aria-hidden');
            this._createBackdrop();
            this._isHidden = false;
            // Add keyboard event listener to the document
            if (this._options.closable) this._setupModalCloseEventListeners();
            // prevent body scroll
            document.body.classList.add('overflow-hidden');
            // callback function
            this._options.onShow(this);
        }
    };
    Modal.prototype.hide = function() {
        if (this.isVisible) {
            this._targetEl.classList.add('hidden');
            this._targetEl.classList.remove('flex');
            this._targetEl.setAttribute('aria-hidden', 'true');
            this._targetEl.removeAttribute('aria-modal');
            this._targetEl.removeAttribute('role');
            this._destroyBackdropEl();
            this._isHidden = true;
            // re-apply body scroll
            document.body.classList.remove('overflow-hidden');
            if (this._options.closable) this._removeModalCloseEventListeners();
            // callback function
            this._options.onHide(this);
        }
    };
    Modal.prototype.isVisible = function() {
        return !this._isHidden;
    };
    Modal.prototype.isHidden = function() {
        return this._isHidden;
    };
    Modal.prototype.addEventListenerInstance = function(element, type, handler) {
        this._eventListenerInstances.push({
            element: element,
            type: type,
            handler: handler
        });
    };
    Modal.prototype.removeAllEventListenerInstances = function() {
        this._eventListenerInstances.map(function(eventListenerInstance) {
            eventListenerInstance.element.removeEventListener(eventListenerInstance.type, eventListenerInstance.handler);
        });
        this._eventListenerInstances = [];
    };
    Modal.prototype.getAllEventListenerInstances = function() {
        return this._eventListenerInstances;
    };
    Modal.prototype.updateOnShow = function(callback) {
        this._options.onShow = callback;
    };
    Modal.prototype.updateOnHide = function(callback) {
        this._options.onHide = callback;
    };
    Modal.prototype.updateOnToggle = function(callback) {
        this._options.onToggle = callback;
    };
    return Modal;
}();
function $8d24653c8cdd4795$export$e14d3bc873a8b9b1() {
    // initiate modal based on data-modal-target
    document.querySelectorAll('[data-modal-target]').forEach(function($triggerEl) {
        var modalId = $triggerEl.getAttribute('data-modal-target');
        var $modalEl = document.getElementById(modalId);
        if ($modalEl) {
            var placement = $modalEl.getAttribute('data-modal-placement');
            var backdrop = $modalEl.getAttribute('data-modal-backdrop');
            new $8d24653c8cdd4795$var$Modal($modalEl, {
                placement: placement ? placement : $8d24653c8cdd4795$var$Default.placement,
                backdrop: backdrop ? backdrop : $8d24653c8cdd4795$var$Default.backdrop
            });
        } else console.error("Modal with id ".concat(modalId, " does not exist. Are you sure that the data-modal-target attribute points to the correct modal id?."));
    });
    // toggle modal visibility
    document.querySelectorAll('[data-modal-toggle]').forEach(function($triggerEl) {
        var modalId = $triggerEl.getAttribute('data-modal-toggle');
        var $modalEl = document.getElementById(modalId);
        if ($modalEl) {
            var modal_1 = (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).getInstance('Modal', modalId);
            if (modal_1) {
                var toggleModal = function() {
                    modal_1.toggle();
                };
                $triggerEl.addEventListener('click', toggleModal);
                modal_1.addEventListenerInstance($triggerEl, 'click', toggleModal);
            } else console.error("Modal with id ".concat(modalId, " has not been initialized. Please initialize it using the data-modal-target attribute."));
        } else console.error("Modal with id ".concat(modalId, " does not exist. Are you sure that the data-modal-toggle attribute points to the correct modal id?"));
    });
    // show modal on click if exists based on id
    document.querySelectorAll('[data-modal-show]').forEach(function($triggerEl) {
        var modalId = $triggerEl.getAttribute('data-modal-show');
        var $modalEl = document.getElementById(modalId);
        if ($modalEl) {
            var modal_2 = (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).getInstance('Modal', modalId);
            if (modal_2) {
                var showModal = function() {
                    modal_2.show();
                };
                $triggerEl.addEventListener('click', showModal);
                modal_2.addEventListenerInstance($triggerEl, 'click', showModal);
            } else console.error("Modal with id ".concat(modalId, " has not been initialized. Please initialize it using the data-modal-target attribute."));
        } else console.error("Modal with id ".concat(modalId, " does not exist. Are you sure that the data-modal-show attribute points to the correct modal id?"));
    });
    // hide modal on click if exists based on id
    document.querySelectorAll('[data-modal-hide]').forEach(function($triggerEl) {
        var modalId = $triggerEl.getAttribute('data-modal-hide');
        var $modalEl = document.getElementById(modalId);
        if ($modalEl) {
            var modal_3 = (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).getInstance('Modal', modalId);
            if (modal_3) {
                var hideModal = function() {
                    modal_3.hide();
                };
                $triggerEl.addEventListener('click', hideModal);
                modal_3.addEventListenerInstance($triggerEl, 'click', hideModal);
            } else console.error("Modal with id ".concat(modalId, " has not been initialized. Please initialize it using the data-modal-target attribute."));
        } else console.error("Modal with id ".concat(modalId, " does not exist. Are you sure that the data-modal-hide attribute points to the correct modal id?"));
    });
}
if (typeof window !== 'undefined') {
    window.Modal = $8d24653c8cdd4795$var$Modal;
    window.initModals = $8d24653c8cdd4795$export$e14d3bc873a8b9b1;
}
var $8d24653c8cdd4795$export$2e2bcd8739ae039 = $8d24653c8cdd4795$var$Modal;



var $c647e009fa86f7f0$var$__assign = undefined && undefined.__assign || function() {
    $c647e009fa86f7f0$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $c647e009fa86f7f0$var$__assign.apply(this, arguments);
};
var $c647e009fa86f7f0$var$Default = {
    placement: 'left',
    bodyScrolling: false,
    backdrop: true,
    edge: false,
    edgeOffset: 'bottom-[60px]',
    backdropClasses: 'bg-dark-backdrop/70 fixed inset-0 z-30',
    onShow: function() {},
    onHide: function() {},
    onToggle: function() {}
};
var $c647e009fa86f7f0$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $c647e009fa86f7f0$var$Drawer = /** @class */ function() {
    function Drawer(targetEl, options, instanceOptions) {
        if (targetEl === void 0) targetEl = null;
        if (options === void 0) options = $c647e009fa86f7f0$var$Default;
        if (instanceOptions === void 0) instanceOptions = $c647e009fa86f7f0$var$DefaultInstanceOptions;
        this._eventListenerInstances = [];
        this._instanceId = instanceOptions.id ? instanceOptions.id : targetEl.id;
        this._targetEl = targetEl;
        this._options = $c647e009fa86f7f0$var$__assign($c647e009fa86f7f0$var$__assign({}, $c647e009fa86f7f0$var$Default), options);
        this._visible = false;
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('Drawer', this, this._instanceId, instanceOptions.override);
    }
    Drawer.prototype.init = function() {
        var _this = this;
        // set initial accessibility attributes
        if (this._targetEl && !this._initialized) {
            this._targetEl.setAttribute('aria-hidden', 'true');
            this._targetEl.classList.add('transition-transform');
            // set base placement classes
            this._getPlacementClasses(this._options.placement).base.map(function(c) {
                _this._targetEl.classList.add(c);
            });
            this._handleEscapeKey = function(event) {
                if (event.key === 'Escape') // if 'Escape' key is pressed
                {
                    if (_this.isVisible()) // if the Drawer is visible
                    _this.hide(); // hide the Drawer
                }
            };
            // add keyboard event listener to document
            document.addEventListener('keydown', this._handleEscapeKey);
            this._initialized = true;
        }
    };
    Drawer.prototype.destroy = function() {
        if (this._initialized) {
            this.removeAllEventListenerInstances();
            this._destroyBackdropEl();
            // Remove the keyboard event listener
            document.removeEventListener('keydown', this._handleEscapeKey);
            this._initialized = false;
        }
    };
    Drawer.prototype.removeInstance = function() {
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('Drawer', this._instanceId);
    };
    Drawer.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    Drawer.prototype.hide = function() {
        var _this = this;
        // based on the edge option show placement classes
        if (this._options.edge) {
            this._getPlacementClasses(this._options.placement + '-edge').active.map(function(c) {
                _this._targetEl.classList.remove(c);
            });
            this._getPlacementClasses(this._options.placement + '-edge').inactive.map(function(c) {
                _this._targetEl.classList.add(c);
            });
        } else {
            this._getPlacementClasses(this._options.placement).active.map(function(c) {
                _this._targetEl.classList.remove(c);
            });
            this._getPlacementClasses(this._options.placement).inactive.map(function(c) {
                _this._targetEl.classList.add(c);
            });
        }
        // set accessibility attributes
        this._targetEl.setAttribute('aria-hidden', 'true');
        this._targetEl.removeAttribute('aria-modal');
        this._targetEl.removeAttribute('role');
        // enable body scroll
        if (!this._options.bodyScrolling) document.body.classList.remove('overflow-hidden');
        // destroy backdrop
        if (this._options.backdrop) this._destroyBackdropEl();
        this._visible = false;
        // callback function
        this._options.onHide(this);
    };
    Drawer.prototype.show = function() {
        var _this = this;
        if (this._options.edge) {
            this._getPlacementClasses(this._options.placement + '-edge').active.map(function(c) {
                _this._targetEl.classList.add(c);
            });
            this._getPlacementClasses(this._options.placement + '-edge').inactive.map(function(c) {
                _this._targetEl.classList.remove(c);
            });
        } else {
            this._getPlacementClasses(this._options.placement).active.map(function(c) {
                _this._targetEl.classList.add(c);
            });
            this._getPlacementClasses(this._options.placement).inactive.map(function(c) {
                _this._targetEl.classList.remove(c);
            });
        }
        // set accessibility attributes
        this._targetEl.setAttribute('aria-modal', 'true');
        this._targetEl.setAttribute('role', 'dialog');
        this._targetEl.removeAttribute('aria-hidden');
        // disable body scroll
        if (!this._options.bodyScrolling) document.body.classList.add('overflow-hidden');
        // show backdrop
        if (this._options.backdrop) this._createBackdrop();
        this._visible = true;
        // callback function
        this._options.onShow(this);
    };
    Drawer.prototype.toggle = function() {
        if (this.isVisible()) this.hide();
        else this.show();
    };
    Drawer.prototype._createBackdrop = function() {
        var _a;
        var _this = this;
        if (!this._visible) {
            var backdropEl = document.createElement('div');
            backdropEl.setAttribute('drawer-backdrop', '');
            (_a = backdropEl.classList).add.apply(_a, this._options.backdropClasses.split(' '));
            document.querySelector('body').append(backdropEl);
            backdropEl.addEventListener('click', function() {
                _this.hide();
            });
        }
    };
    Drawer.prototype._destroyBackdropEl = function() {
        if (this._visible && document.querySelector('[drawer-backdrop]') !== null) document.querySelector('[drawer-backdrop]').remove();
    };
    Drawer.prototype._getPlacementClasses = function(placement) {
        switch(placement){
            case 'top':
                return {
                    base: [
                        'top-0',
                        'left-0',
                        'right-0'
                    ],
                    active: [
                        'transform-none'
                    ],
                    inactive: [
                        '-translate-y-full'
                    ]
                };
            case 'right':
                return {
                    base: [
                        'right-0',
                        'top-0'
                    ],
                    active: [
                        'transform-none'
                    ],
                    inactive: [
                        'translate-x-full'
                    ]
                };
            case 'bottom':
                return {
                    base: [
                        'bottom-0',
                        'left-0',
                        'right-0'
                    ],
                    active: [
                        'transform-none'
                    ],
                    inactive: [
                        'translate-y-full'
                    ]
                };
            case 'left':
                return {
                    base: [
                        'left-0',
                        'top-0'
                    ],
                    active: [
                        'transform-none'
                    ],
                    inactive: [
                        '-translate-x-full'
                    ]
                };
            case 'bottom-edge':
                return {
                    base: [
                        'left-0',
                        'top-0'
                    ],
                    active: [
                        'transform-none'
                    ],
                    inactive: [
                        'translate-y-full',
                        this._options.edgeOffset
                    ]
                };
            default:
                return {
                    base: [
                        'left-0',
                        'top-0'
                    ],
                    active: [
                        'transform-none'
                    ],
                    inactive: [
                        '-translate-x-full'
                    ]
                };
        }
    };
    Drawer.prototype.isHidden = function() {
        return !this._visible;
    };
    Drawer.prototype.isVisible = function() {
        return this._visible;
    };
    Drawer.prototype.addEventListenerInstance = function(element, type, handler) {
        this._eventListenerInstances.push({
            element: element,
            type: type,
            handler: handler
        });
    };
    Drawer.prototype.removeAllEventListenerInstances = function() {
        this._eventListenerInstances.map(function(eventListenerInstance) {
            eventListenerInstance.element.removeEventListener(eventListenerInstance.type, eventListenerInstance.handler);
        });
        this._eventListenerInstances = [];
    };
    Drawer.prototype.getAllEventListenerInstances = function() {
        return this._eventListenerInstances;
    };
    Drawer.prototype.updateOnShow = function(callback) {
        this._options.onShow = callback;
    };
    Drawer.prototype.updateOnHide = function(callback) {
        this._options.onHide = callback;
    };
    Drawer.prototype.updateOnToggle = function(callback) {
        this._options.onToggle = callback;
    };
    return Drawer;
}();
function $c647e009fa86f7f0$export$319df2b10e87c8ec() {
    document.querySelectorAll('[data-drawer-target]').forEach(function($triggerEl) {
        // mandatory
        var drawerId = $triggerEl.getAttribute('data-drawer-target');
        var $drawerEl = document.getElementById(drawerId);
        if ($drawerEl) {
            var placement = $triggerEl.getAttribute('data-drawer-placement');
            var bodyScrolling = $triggerEl.getAttribute('data-drawer-body-scrolling');
            var backdrop = $triggerEl.getAttribute('data-drawer-backdrop');
            var edge = $triggerEl.getAttribute('data-drawer-edge');
            var edgeOffset = $triggerEl.getAttribute('data-drawer-edge-offset');
            new $c647e009fa86f7f0$var$Drawer($drawerEl, {
                placement: placement ? placement : $c647e009fa86f7f0$var$Default.placement,
                bodyScrolling: bodyScrolling ? bodyScrolling === 'true' ? true : false : $c647e009fa86f7f0$var$Default.bodyScrolling,
                backdrop: backdrop ? backdrop === 'true' ? true : false : $c647e009fa86f7f0$var$Default.backdrop,
                edge: edge ? edge === 'true' ? true : false : $c647e009fa86f7f0$var$Default.edge,
                edgeOffset: edgeOffset ? edgeOffset : $c647e009fa86f7f0$var$Default.edgeOffset
            });
        } else console.error("Drawer with id ".concat(drawerId, " not found. Are you sure that the data-drawer-target attribute points to the correct drawer id?"));
    });
    document.querySelectorAll('[data-drawer-toggle]').forEach(function($triggerEl) {
        var drawerId = $triggerEl.getAttribute('data-drawer-toggle');
        var $drawerEl = document.getElementById(drawerId);
        if ($drawerEl) {
            var drawer_1 = (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).getInstance('Drawer', drawerId);
            if (drawer_1) {
                var toggleDrawer = function() {
                    drawer_1.toggle();
                };
                $triggerEl.addEventListener('click', toggleDrawer);
                drawer_1.addEventListenerInstance($triggerEl, 'click', toggleDrawer);
            } else console.error("Drawer with id ".concat(drawerId, " has not been initialized. Please initialize it using the data-drawer-target attribute."));
        } else console.error("Drawer with id ".concat(drawerId, " not found. Are you sure that the data-drawer-target attribute points to the correct drawer id?"));
    });
    document.querySelectorAll('[data-drawer-dismiss], [data-drawer-hide]').forEach(function($triggerEl) {
        var drawerId = $triggerEl.getAttribute('data-drawer-dismiss') ? $triggerEl.getAttribute('data-drawer-dismiss') : $triggerEl.getAttribute('data-drawer-hide');
        var $drawerEl = document.getElementById(drawerId);
        if ($drawerEl) {
            var drawer_2 = (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).getInstance('Drawer', drawerId);
            if (drawer_2) {
                var hideDrawer = function() {
                    drawer_2.hide();
                };
                $triggerEl.addEventListener('click', hideDrawer);
                drawer_2.addEventListenerInstance($triggerEl, 'click', hideDrawer);
            } else console.error("Drawer with id ".concat(drawerId, " has not been initialized. Please initialize it using the data-drawer-target attribute."));
        } else console.error("Drawer with id ".concat(drawerId, " not found. Are you sure that the data-drawer-target attribute points to the correct drawer id"));
    });
    document.querySelectorAll('[data-drawer-show]').forEach(function($triggerEl) {
        var drawerId = $triggerEl.getAttribute('data-drawer-show');
        var $drawerEl = document.getElementById(drawerId);
        if ($drawerEl) {
            var drawer_3 = (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).getInstance('Drawer', drawerId);
            if (drawer_3) {
                var showDrawer = function() {
                    drawer_3.show();
                };
                $triggerEl.addEventListener('click', showDrawer);
                drawer_3.addEventListenerInstance($triggerEl, 'click', showDrawer);
            } else console.error("Drawer with id ".concat(drawerId, " has not been initialized. Please initialize it using the data-drawer-target attribute."));
        } else console.error("Drawer with id ".concat(drawerId, " not found. Are you sure that the data-drawer-target attribute points to the correct drawer id?"));
    });
}
if (typeof window !== 'undefined') {
    window.Drawer = $c647e009fa86f7f0$var$Drawer;
    window.initDrawers = $c647e009fa86f7f0$export$319df2b10e87c8ec;
}
var $c647e009fa86f7f0$export$2e2bcd8739ae039 = $c647e009fa86f7f0$var$Drawer;



var $ff078cf47b3c0903$var$__assign = undefined && undefined.__assign || function() {
    $ff078cf47b3c0903$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $ff078cf47b3c0903$var$__assign.apply(this, arguments);
};
var $ff078cf47b3c0903$var$Default = {
    defaultTabId: null,
    activeClasses: 'text-fg-brand hover:text-fg-brand border-brand',
    inactiveClasses: 'border-transparent text-body hover:text-heading border-soft hover:border-default',
    onShow: function() {}
};
var $ff078cf47b3c0903$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $ff078cf47b3c0903$var$Tabs = /** @class */ function() {
    function Tabs(tabsEl, items, options, instanceOptions) {
        if (tabsEl === void 0) tabsEl = null;
        if (items === void 0) items = [];
        if (options === void 0) options = $ff078cf47b3c0903$var$Default;
        if (instanceOptions === void 0) instanceOptions = $ff078cf47b3c0903$var$DefaultInstanceOptions;
        this._instanceId = instanceOptions.id ? instanceOptions.id : tabsEl.id;
        this._tabsEl = tabsEl;
        this._items = items;
        this._activeTab = options ? this.getTab(options.defaultTabId) : null;
        this._options = $ff078cf47b3c0903$var$__assign($ff078cf47b3c0903$var$__assign({}, $ff078cf47b3c0903$var$Default), options);
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('Tabs', this, this._instanceId, instanceOptions.override);
    }
    Tabs.prototype.init = function() {
        var _this = this;
        if (this._items.length && !this._initialized) {
            // set the first tab as active if not set by explicitly
            if (!this._activeTab) this.setActiveTab(this._items[0]);
            // force show the first default tab
            this.show(this._activeTab.id, true);
            // show tab content based on click
            this._items.map(function(tab) {
                tab.triggerEl.addEventListener('click', function(event) {
                    event.preventDefault();
                    _this.show(tab.id);
                });
            });
        }
    };
    Tabs.prototype.destroy = function() {
        if (this._initialized) this._initialized = false;
    };
    Tabs.prototype.removeInstance = function() {
        this.destroy();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('Tabs', this._instanceId);
    };
    Tabs.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    Tabs.prototype.getActiveTab = function() {
        return this._activeTab;
    };
    Tabs.prototype.setActiveTab = function(tab) {
        this._activeTab = tab;
    };
    Tabs.prototype.getTab = function(id) {
        return this._items.filter(function(t) {
            return t.id === id;
        })[0];
    };
    Tabs.prototype.show = function(id, forceShow) {
        var _a, _b;
        var _this = this;
        if (forceShow === void 0) forceShow = false;
        var tab = this.getTab(id);
        // don't do anything if already active
        if (tab === this._activeTab && !forceShow) return;
        // hide other tabs
        this._items.map(function(t) {
            var _a, _b;
            if (t !== tab) {
                (_a = t.triggerEl.classList).remove.apply(_a, _this._options.activeClasses.split(' '));
                (_b = t.triggerEl.classList).add.apply(_b, _this._options.inactiveClasses.split(' '));
                t.targetEl.classList.add('hidden');
                t.triggerEl.setAttribute('aria-selected', 'false');
            }
        });
        // show active tab
        (_a = tab.triggerEl.classList).add.apply(_a, this._options.activeClasses.split(' '));
        (_b = tab.triggerEl.classList).remove.apply(_b, this._options.inactiveClasses.split(' '));
        tab.triggerEl.setAttribute('aria-selected', 'true');
        tab.targetEl.classList.remove('hidden');
        this.setActiveTab(tab);
        // callback function
        this._options.onShow(this, tab);
    };
    Tabs.prototype.updateOnShow = function(callback) {
        this._options.onShow = callback;
    };
    return Tabs;
}();
function $ff078cf47b3c0903$export$c92e8f569ad2976f() {
    document.querySelectorAll('[data-tabs-toggle]').forEach(function($parentEl) {
        var tabItems = [];
        var activeClasses = $parentEl.getAttribute('data-tabs-active-classes');
        var inactiveClasses = $parentEl.getAttribute('data-tabs-inactive-classes');
        var defaultTabId = null;
        $parentEl.querySelectorAll('[role="tab"]').forEach(function($triggerEl) {
            var isActive = $triggerEl.getAttribute('aria-selected') === 'true';
            var tab = {
                id: $triggerEl.getAttribute('data-tabs-target'),
                triggerEl: $triggerEl,
                targetEl: document.querySelector($triggerEl.getAttribute('data-tabs-target'))
            };
            tabItems.push(tab);
            if (isActive) defaultTabId = tab.id;
        });
        new $ff078cf47b3c0903$var$Tabs($parentEl, tabItems, {
            defaultTabId: defaultTabId,
            activeClasses: activeClasses ? activeClasses : $ff078cf47b3c0903$var$Default.activeClasses,
            inactiveClasses: inactiveClasses ? inactiveClasses : $ff078cf47b3c0903$var$Default.inactiveClasses
        });
    });
}
if (typeof window !== 'undefined') {
    window.Tabs = $ff078cf47b3c0903$var$Tabs;
    window.initTabs = $ff078cf47b3c0903$export$c92e8f569ad2976f;
}
var $ff078cf47b3c0903$export$2e2bcd8739ae039 = $ff078cf47b3c0903$var$Tabs;




var $7ab0a41fd48c3b4d$var$__assign = undefined && undefined.__assign || function() {
    $7ab0a41fd48c3b4d$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $7ab0a41fd48c3b4d$var$__assign.apply(this, arguments);
};
var $7ab0a41fd48c3b4d$var$__spreadArray = undefined && undefined.__spreadArray || function(to, from, pack) {
    if (pack || arguments.length === 2) {
        for(var i = 0, l = from.length, ar; i < l; i++)if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
};
var $7ab0a41fd48c3b4d$var$Default = {
    placement: 'top',
    triggerType: 'hover',
    onShow: function() {},
    onHide: function() {},
    onToggle: function() {}
};
var $7ab0a41fd48c3b4d$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $7ab0a41fd48c3b4d$var$Tooltip = /** @class */ function() {
    function Tooltip(targetEl, triggerEl, options, instanceOptions) {
        if (targetEl === void 0) targetEl = null;
        if (triggerEl === void 0) triggerEl = null;
        if (options === void 0) options = $7ab0a41fd48c3b4d$var$Default;
        if (instanceOptions === void 0) instanceOptions = $7ab0a41fd48c3b4d$var$DefaultInstanceOptions;
        this._instanceId = instanceOptions.id ? instanceOptions.id : targetEl.id;
        this._targetEl = targetEl;
        this._triggerEl = triggerEl;
        this._options = $7ab0a41fd48c3b4d$var$__assign($7ab0a41fd48c3b4d$var$__assign({}, $7ab0a41fd48c3b4d$var$Default), options);
        this._popperInstance = null;
        this._visible = false;
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('Tooltip', this, this._instanceId, instanceOptions.override);
    }
    Tooltip.prototype.init = function() {
        if (this._triggerEl && this._targetEl && !this._initialized) {
            this._setupEventListeners();
            this._popperInstance = this._createPopperInstance();
            this._initialized = true;
        }
    };
    Tooltip.prototype.destroy = function() {
        var _this = this;
        if (this._initialized) {
            // remove event listeners associated with the trigger element
            var triggerEvents = this._getTriggerEvents();
            triggerEvents.showEvents.forEach(function(ev) {
                _this._triggerEl.removeEventListener(ev, _this._showHandler);
            });
            triggerEvents.hideEvents.forEach(function(ev) {
                _this._triggerEl.removeEventListener(ev, _this._hideHandler);
            });
            // remove event listeners for keydown
            this._removeKeydownListener();
            // remove event listeners for click outside
            this._removeClickOutsideListener();
            // destroy the Popper instance if you have one (assuming this._popperInstance is the Popper instance)
            if (this._popperInstance) this._popperInstance.destroy();
            this._initialized = false;
        }
    };
    Tooltip.prototype.removeInstance = function() {
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('Tooltip', this._instanceId);
    };
    Tooltip.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    Tooltip.prototype._setupEventListeners = function() {
        var _this = this;
        var triggerEvents = this._getTriggerEvents();
        this._showHandler = function() {
            _this.show();
        };
        this._hideHandler = function() {
            _this.hide();
        };
        triggerEvents.showEvents.forEach(function(ev) {
            _this._triggerEl.addEventListener(ev, _this._showHandler);
        });
        triggerEvents.hideEvents.forEach(function(ev) {
            _this._triggerEl.addEventListener(ev, _this._hideHandler);
        });
    };
    Tooltip.prototype._createPopperInstance = function() {
        return (0, $725309ecee3adaa0$export$8f7491d57c8f97a9)(this._triggerEl, this._targetEl, {
            placement: this._options.placement,
            modifiers: [
                {
                    name: 'offset',
                    options: {
                        offset: [
                            0,
                            8
                        ]
                    }
                }
            ]
        });
    };
    Tooltip.prototype._getTriggerEvents = function() {
        switch(this._options.triggerType){
            case 'hover':
                return {
                    showEvents: [
                        'mouseenter',
                        'focus'
                    ],
                    hideEvents: [
                        'mouseleave',
                        'blur'
                    ]
                };
            case 'click':
                return {
                    showEvents: [
                        'click',
                        'focus'
                    ],
                    hideEvents: [
                        'focusout',
                        'blur'
                    ]
                };
            case 'none':
                return {
                    showEvents: [],
                    hideEvents: []
                };
            default:
                return {
                    showEvents: [
                        'mouseenter',
                        'focus'
                    ],
                    hideEvents: [
                        'mouseleave',
                        'blur'
                    ]
                };
        }
    };
    Tooltip.prototype._setupKeydownListener = function() {
        var _this = this;
        this._keydownEventListener = function(ev) {
            if (ev.key === 'Escape') _this.hide();
        };
        document.body.addEventListener('keydown', this._keydownEventListener, true);
    };
    Tooltip.prototype._removeKeydownListener = function() {
        document.body.removeEventListener('keydown', this._keydownEventListener, true);
    };
    Tooltip.prototype._setupClickOutsideListener = function() {
        var _this = this;
        this._clickOutsideEventListener = function(ev) {
            _this._handleClickOutside(ev, _this._targetEl);
        };
        document.body.addEventListener('click', this._clickOutsideEventListener, true);
    };
    Tooltip.prototype._removeClickOutsideListener = function() {
        document.body.removeEventListener('click', this._clickOutsideEventListener, true);
    };
    Tooltip.prototype._handleClickOutside = function(ev, targetEl) {
        var clickedEl = ev.target;
        if (clickedEl !== targetEl && !targetEl.contains(clickedEl) && !this._triggerEl.contains(clickedEl) && this.isVisible()) this.hide();
    };
    Tooltip.prototype.isVisible = function() {
        return this._visible;
    };
    Tooltip.prototype.toggle = function() {
        if (this.isVisible()) this.hide();
        else this.show();
    };
    Tooltip.prototype.show = function() {
        this._targetEl.classList.remove('opacity-0', 'invisible');
        this._targetEl.classList.add('opacity-100', 'visible');
        // Enable the event listeners
        this._popperInstance.setOptions(function(options) {
            return $7ab0a41fd48c3b4d$var$__assign($7ab0a41fd48c3b4d$var$__assign({}, options), {
                modifiers: $7ab0a41fd48c3b4d$var$__spreadArray($7ab0a41fd48c3b4d$var$__spreadArray([], options.modifiers, true), [
                    {
                        name: 'eventListeners',
                        enabled: true
                    }
                ], false)
            });
        });
        // handle click outside
        this._setupClickOutsideListener();
        // handle esc keydown
        this._setupKeydownListener();
        // Update its position
        this._popperInstance.update();
        // set visibility
        this._visible = true;
        // callback function
        this._options.onShow(this);
    };
    Tooltip.prototype.hide = function() {
        this._targetEl.classList.remove('opacity-100', 'visible');
        this._targetEl.classList.add('opacity-0', 'invisible');
        // Disable the event listeners
        this._popperInstance.setOptions(function(options) {
            return $7ab0a41fd48c3b4d$var$__assign($7ab0a41fd48c3b4d$var$__assign({}, options), {
                modifiers: $7ab0a41fd48c3b4d$var$__spreadArray($7ab0a41fd48c3b4d$var$__spreadArray([], options.modifiers, true), [
                    {
                        name: 'eventListeners',
                        enabled: false
                    }
                ], false)
            });
        });
        // handle click outside
        this._removeClickOutsideListener();
        // handle esc keydown
        this._removeKeydownListener();
        // set visibility
        this._visible = false;
        // callback function
        this._options.onHide(this);
    };
    Tooltip.prototype.updateOnShow = function(callback) {
        this._options.onShow = callback;
    };
    Tooltip.prototype.updateOnHide = function(callback) {
        this._options.onHide = callback;
    };
    Tooltip.prototype.updateOnToggle = function(callback) {
        this._options.onToggle = callback;
    };
    return Tooltip;
}();
function $7ab0a41fd48c3b4d$export$8f2a38c1f9d9dc3e() {
    document.querySelectorAll('[data-tooltip-target]').forEach(function($triggerEl) {
        var tooltipId = $triggerEl.getAttribute('data-tooltip-target');
        var $tooltipEl = document.getElementById(tooltipId);
        if ($tooltipEl) {
            var triggerType = $triggerEl.getAttribute('data-tooltip-trigger');
            var placement = $triggerEl.getAttribute('data-tooltip-placement');
            new $7ab0a41fd48c3b4d$var$Tooltip($tooltipEl, $triggerEl, {
                placement: placement ? placement : $7ab0a41fd48c3b4d$var$Default.placement,
                triggerType: triggerType ? triggerType : $7ab0a41fd48c3b4d$var$Default.triggerType
            });
        } else console.error("The tooltip element with id \"".concat(tooltipId, "\" does not exist. Please check the data-tooltip-target attribute."));
    });
}
if (typeof window !== 'undefined') {
    window.Tooltip = $7ab0a41fd48c3b4d$var$Tooltip;
    window.initTooltips = $7ab0a41fd48c3b4d$export$8f2a38c1f9d9dc3e;
}
var $7ab0a41fd48c3b4d$export$2e2bcd8739ae039 = $7ab0a41fd48c3b4d$var$Tooltip;




var $5370c2686863c309$var$__assign = undefined && undefined.__assign || function() {
    $5370c2686863c309$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $5370c2686863c309$var$__assign.apply(this, arguments);
};
var $5370c2686863c309$var$__spreadArray = undefined && undefined.__spreadArray || function(to, from, pack) {
    if (pack || arguments.length === 2) {
        for(var i = 0, l = from.length, ar; i < l; i++)if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
};
var $5370c2686863c309$var$Default = {
    placement: 'top',
    offset: 10,
    triggerType: 'hover',
    onShow: function() {},
    onHide: function() {},
    onToggle: function() {}
};
var $5370c2686863c309$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $5370c2686863c309$var$Popover = /** @class */ function() {
    function Popover(targetEl, triggerEl, options, instanceOptions) {
        if (targetEl === void 0) targetEl = null;
        if (triggerEl === void 0) triggerEl = null;
        if (options === void 0) options = $5370c2686863c309$var$Default;
        if (instanceOptions === void 0) instanceOptions = $5370c2686863c309$var$DefaultInstanceOptions;
        this._instanceId = instanceOptions.id ? instanceOptions.id : targetEl.id;
        this._targetEl = targetEl;
        this._triggerEl = triggerEl;
        this._options = $5370c2686863c309$var$__assign($5370c2686863c309$var$__assign({}, $5370c2686863c309$var$Default), options);
        this._popperInstance = null;
        this._visible = false;
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('Popover', this, instanceOptions.id ? instanceOptions.id : this._targetEl.id, instanceOptions.override);
    }
    Popover.prototype.init = function() {
        if (this._triggerEl && this._targetEl && !this._initialized) {
            this._setupEventListeners();
            this._popperInstance = this._createPopperInstance();
            this._initialized = true;
        }
    };
    Popover.prototype.destroy = function() {
        var _this = this;
        if (this._initialized) {
            // remove event listeners associated with the trigger element and target element
            var triggerEvents = this._getTriggerEvents();
            triggerEvents.showEvents.forEach(function(ev) {
                _this._triggerEl.removeEventListener(ev, _this._showHandler);
                _this._targetEl.removeEventListener(ev, _this._showHandler);
            });
            triggerEvents.hideEvents.forEach(function(ev) {
                _this._triggerEl.removeEventListener(ev, _this._hideHandler);
                _this._targetEl.removeEventListener(ev, _this._hideHandler);
            });
            // remove event listeners for keydown
            this._removeKeydownListener();
            // remove event listeners for click outside
            this._removeClickOutsideListener();
            // destroy the Popper instance if you have one (assuming this._popperInstance is the Popper instance)
            if (this._popperInstance) this._popperInstance.destroy();
            this._initialized = false;
        }
    };
    Popover.prototype.removeInstance = function() {
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('Popover', this._instanceId);
    };
    Popover.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    Popover.prototype._setupEventListeners = function() {
        var _this = this;
        var triggerEvents = this._getTriggerEvents();
        this._showHandler = function() {
            _this.show();
        };
        this._hideHandler = function() {
            setTimeout(function() {
                if (!_this._targetEl.matches(':hover')) _this.hide();
            }, 100);
        };
        triggerEvents.showEvents.forEach(function(ev) {
            _this._triggerEl.addEventListener(ev, _this._showHandler);
            _this._targetEl.addEventListener(ev, _this._showHandler);
        });
        triggerEvents.hideEvents.forEach(function(ev) {
            _this._triggerEl.addEventListener(ev, _this._hideHandler);
            _this._targetEl.addEventListener(ev, _this._hideHandler);
        });
    };
    Popover.prototype._createPopperInstance = function() {
        return (0, $725309ecee3adaa0$export$8f7491d57c8f97a9)(this._triggerEl, this._targetEl, {
            placement: this._options.placement,
            modifiers: [
                {
                    name: 'offset',
                    options: {
                        offset: [
                            0,
                            this._options.offset
                        ]
                    }
                }
            ]
        });
    };
    Popover.prototype._getTriggerEvents = function() {
        switch(this._options.triggerType){
            case 'hover':
                return {
                    showEvents: [
                        'mouseenter',
                        'focus'
                    ],
                    hideEvents: [
                        'mouseleave',
                        'blur'
                    ]
                };
            case 'click':
                return {
                    showEvents: [
                        'click',
                        'focus'
                    ],
                    hideEvents: [
                        'focusout',
                        'blur'
                    ]
                };
            case 'none':
                return {
                    showEvents: [],
                    hideEvents: []
                };
            default:
                return {
                    showEvents: [
                        'mouseenter',
                        'focus'
                    ],
                    hideEvents: [
                        'mouseleave',
                        'blur'
                    ]
                };
        }
    };
    Popover.prototype._setupKeydownListener = function() {
        var _this = this;
        this._keydownEventListener = function(ev) {
            if (ev.key === 'Escape') _this.hide();
        };
        document.body.addEventListener('keydown', this._keydownEventListener, true);
    };
    Popover.prototype._removeKeydownListener = function() {
        document.body.removeEventListener('keydown', this._keydownEventListener, true);
    };
    Popover.prototype._setupClickOutsideListener = function() {
        var _this = this;
        this._clickOutsideEventListener = function(ev) {
            _this._handleClickOutside(ev, _this._targetEl);
        };
        document.body.addEventListener('click', this._clickOutsideEventListener, true);
    };
    Popover.prototype._removeClickOutsideListener = function() {
        document.body.removeEventListener('click', this._clickOutsideEventListener, true);
    };
    Popover.prototype._handleClickOutside = function(ev, targetEl) {
        var clickedEl = ev.target;
        if (clickedEl !== targetEl && !targetEl.contains(clickedEl) && !this._triggerEl.contains(clickedEl) && this.isVisible()) this.hide();
    };
    Popover.prototype.isVisible = function() {
        return this._visible;
    };
    Popover.prototype.toggle = function() {
        if (this.isVisible()) this.hide();
        else this.show();
        this._options.onToggle(this);
    };
    Popover.prototype.show = function() {
        this._targetEl.classList.remove('opacity-0', 'invisible');
        this._targetEl.classList.add('opacity-100', 'visible');
        // Enable the event listeners
        this._popperInstance.setOptions(function(options) {
            return $5370c2686863c309$var$__assign($5370c2686863c309$var$__assign({}, options), {
                modifiers: $5370c2686863c309$var$__spreadArray($5370c2686863c309$var$__spreadArray([], options.modifiers, true), [
                    {
                        name: 'eventListeners',
                        enabled: true
                    }
                ], false)
            });
        });
        // handle click outside
        this._setupClickOutsideListener();
        // handle esc keydown
        this._setupKeydownListener();
        // Update its position
        this._popperInstance.update();
        // set visibility to true
        this._visible = true;
        // callback function
        this._options.onShow(this);
    };
    Popover.prototype.hide = function() {
        this._targetEl.classList.remove('opacity-100', 'visible');
        this._targetEl.classList.add('opacity-0', 'invisible');
        // Disable the event listeners
        this._popperInstance.setOptions(function(options) {
            return $5370c2686863c309$var$__assign($5370c2686863c309$var$__assign({}, options), {
                modifiers: $5370c2686863c309$var$__spreadArray($5370c2686863c309$var$__spreadArray([], options.modifiers, true), [
                    {
                        name: 'eventListeners',
                        enabled: false
                    }
                ], false)
            });
        });
        // handle click outside
        this._removeClickOutsideListener();
        // handle esc keydown
        this._removeKeydownListener();
        // set visibility to false
        this._visible = false;
        // callback function
        this._options.onHide(this);
    };
    Popover.prototype.updateOnShow = function(callback) {
        this._options.onShow = callback;
    };
    Popover.prototype.updateOnHide = function(callback) {
        this._options.onHide = callback;
    };
    Popover.prototype.updateOnToggle = function(callback) {
        this._options.onToggle = callback;
    };
    return Popover;
}();
function $5370c2686863c309$export$200409e83b2a0dd4() {
    document.querySelectorAll('[data-popover-target]').forEach(function($triggerEl) {
        var popoverID = $triggerEl.getAttribute('data-popover-target');
        var $popoverEl = document.getElementById(popoverID);
        if ($popoverEl) {
            var triggerType = $triggerEl.getAttribute('data-popover-trigger');
            var placement = $triggerEl.getAttribute('data-popover-placement');
            var offset = $triggerEl.getAttribute('data-popover-offset');
            new $5370c2686863c309$var$Popover($popoverEl, $triggerEl, {
                placement: placement ? placement : $5370c2686863c309$var$Default.placement,
                offset: offset ? parseInt(offset) : $5370c2686863c309$var$Default.offset,
                triggerType: triggerType ? triggerType : $5370c2686863c309$var$Default.triggerType
            });
        } else console.error("The popover element with id \"".concat(popoverID, "\" does not exist. Please check the data-popover-target attribute."));
    });
}
if (typeof window !== 'undefined') {
    window.Popover = $5370c2686863c309$var$Popover;
    window.initPopovers = $5370c2686863c309$export$200409e83b2a0dd4;
}
var $5370c2686863c309$export$2e2bcd8739ae039 = $5370c2686863c309$var$Popover;



var $4fe34b241a1ad5f0$var$__assign = undefined && undefined.__assign || function() {
    $4fe34b241a1ad5f0$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $4fe34b241a1ad5f0$var$__assign.apply(this, arguments);
};
var $4fe34b241a1ad5f0$var$Default = {
    triggerType: 'hover',
    onShow: function() {},
    onHide: function() {},
    onToggle: function() {}
};
var $4fe34b241a1ad5f0$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $4fe34b241a1ad5f0$var$Dial = /** @class */ function() {
    function Dial(parentEl, triggerEl, targetEl, options, instanceOptions) {
        if (parentEl === void 0) parentEl = null;
        if (triggerEl === void 0) triggerEl = null;
        if (targetEl === void 0) targetEl = null;
        if (options === void 0) options = $4fe34b241a1ad5f0$var$Default;
        if (instanceOptions === void 0) instanceOptions = $4fe34b241a1ad5f0$var$DefaultInstanceOptions;
        this._instanceId = instanceOptions.id ? instanceOptions.id : targetEl.id;
        this._parentEl = parentEl;
        this._triggerEl = triggerEl;
        this._targetEl = targetEl;
        this._options = $4fe34b241a1ad5f0$var$__assign($4fe34b241a1ad5f0$var$__assign({}, $4fe34b241a1ad5f0$var$Default), options);
        this._visible = false;
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('Dial', this, this._instanceId, instanceOptions.override);
    }
    Dial.prototype.init = function() {
        var _this = this;
        if (this._triggerEl && this._targetEl && !this._initialized) {
            var triggerEventTypes = this._getTriggerEventTypes(this._options.triggerType);
            this._showEventHandler = function() {
                _this.show();
            };
            triggerEventTypes.showEvents.forEach(function(ev) {
                _this._triggerEl.addEventListener(ev, _this._showEventHandler);
                _this._targetEl.addEventListener(ev, _this._showEventHandler);
            });
            this._hideEventHandler = function() {
                if (!_this._parentEl.matches(':hover')) _this.hide();
            };
            triggerEventTypes.hideEvents.forEach(function(ev) {
                _this._parentEl.addEventListener(ev, _this._hideEventHandler);
            });
            this._initialized = true;
        }
    };
    Dial.prototype.destroy = function() {
        var _this = this;
        if (this._initialized) {
            var triggerEventTypes = this._getTriggerEventTypes(this._options.triggerType);
            triggerEventTypes.showEvents.forEach(function(ev) {
                _this._triggerEl.removeEventListener(ev, _this._showEventHandler);
                _this._targetEl.removeEventListener(ev, _this._showEventHandler);
            });
            triggerEventTypes.hideEvents.forEach(function(ev) {
                _this._parentEl.removeEventListener(ev, _this._hideEventHandler);
            });
            this._initialized = false;
        }
    };
    Dial.prototype.removeInstance = function() {
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('Dial', this._instanceId);
    };
    Dial.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    Dial.prototype.hide = function() {
        this._targetEl.classList.add('hidden');
        if (this._triggerEl) this._triggerEl.setAttribute('aria-expanded', 'false');
        this._visible = false;
        // callback function
        this._options.onHide(this);
    };
    Dial.prototype.show = function() {
        this._targetEl.classList.remove('hidden');
        if (this._triggerEl) this._triggerEl.setAttribute('aria-expanded', 'true');
        this._visible = true;
        // callback function
        this._options.onShow(this);
    };
    Dial.prototype.toggle = function() {
        if (this._visible) this.hide();
        else this.show();
    };
    Dial.prototype.isHidden = function() {
        return !this._visible;
    };
    Dial.prototype.isVisible = function() {
        return this._visible;
    };
    Dial.prototype._getTriggerEventTypes = function(triggerType) {
        switch(triggerType){
            case 'hover':
                return {
                    showEvents: [
                        'mouseenter',
                        'focus'
                    ],
                    hideEvents: [
                        'mouseleave',
                        'blur'
                    ]
                };
            case 'click':
                return {
                    showEvents: [
                        'click',
                        'focus'
                    ],
                    hideEvents: [
                        'focusout',
                        'blur'
                    ]
                };
            case 'none':
                return {
                    showEvents: [],
                    hideEvents: []
                };
            default:
                return {
                    showEvents: [
                        'mouseenter',
                        'focus'
                    ],
                    hideEvents: [
                        'mouseleave',
                        'blur'
                    ]
                };
        }
    };
    Dial.prototype.updateOnShow = function(callback) {
        this._options.onShow = callback;
    };
    Dial.prototype.updateOnHide = function(callback) {
        this._options.onHide = callback;
    };
    Dial.prototype.updateOnToggle = function(callback) {
        this._options.onToggle = callback;
    };
    return Dial;
}();
function $4fe34b241a1ad5f0$export$33aa68e0bfaa065d() {
    document.querySelectorAll('[data-dial-init]').forEach(function($parentEl) {
        var $triggerEl = $parentEl.querySelector('[data-dial-toggle]');
        if ($triggerEl) {
            var dialId = $triggerEl.getAttribute('data-dial-toggle');
            var $dialEl = document.getElementById(dialId);
            if ($dialEl) {
                var triggerType = $triggerEl.getAttribute('data-dial-trigger');
                new $4fe34b241a1ad5f0$var$Dial($parentEl, $triggerEl, $dialEl, {
                    triggerType: triggerType ? triggerType : $4fe34b241a1ad5f0$var$Default.triggerType
                });
            } else console.error("Dial with id ".concat(dialId, " does not exist. Are you sure that the data-dial-toggle attribute points to the correct modal id?"));
        } else console.error("Dial with id ".concat($parentEl.id, " does not have a trigger element. Are you sure that the data-dial-toggle attribute exists?"));
    });
}
if (typeof window !== 'undefined') {
    window.Dial = $4fe34b241a1ad5f0$var$Dial;
    window.initDials = $4fe34b241a1ad5f0$export$33aa68e0bfaa065d;
}
var $4fe34b241a1ad5f0$export$2e2bcd8739ae039 = $4fe34b241a1ad5f0$var$Dial;



var $df0fc5a35c362eab$var$__assign = undefined && undefined.__assign || function() {
    $df0fc5a35c362eab$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $df0fc5a35c362eab$var$__assign.apply(this, arguments);
};
var $df0fc5a35c362eab$var$Default = {
    minValue: null,
    maxValue: null,
    onIncrement: function() {},
    onDecrement: function() {}
};
var $df0fc5a35c362eab$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $df0fc5a35c362eab$var$InputCounter = /** @class */ function() {
    function InputCounter(targetEl, incrementEl, decrementEl, options, instanceOptions) {
        if (targetEl === void 0) targetEl = null;
        if (incrementEl === void 0) incrementEl = null;
        if (decrementEl === void 0) decrementEl = null;
        if (options === void 0) options = $df0fc5a35c362eab$var$Default;
        if (instanceOptions === void 0) instanceOptions = $df0fc5a35c362eab$var$DefaultInstanceOptions;
        this._instanceId = instanceOptions.id ? instanceOptions.id : targetEl.id;
        this._targetEl = targetEl;
        this._incrementEl = incrementEl;
        this._decrementEl = decrementEl;
        this._options = $df0fc5a35c362eab$var$__assign($df0fc5a35c362eab$var$__assign({}, $df0fc5a35c362eab$var$Default), options);
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('InputCounter', this, this._instanceId, instanceOptions.override);
    }
    InputCounter.prototype.init = function() {
        var _this = this;
        if (this._targetEl && !this._initialized) {
            this._inputHandler = function(event) {
                var target = event.target;
                // check if the value is numeric
                if (!/^\d*$/.test(target.value)) // Regex to check if the value is numeric
                target.value = target.value.replace(/[^\d]/g, ''); // Remove non-numeric characters
                // check for max value
                if (_this._options.maxValue !== null && parseInt(target.value) > _this._options.maxValue) target.value = _this._options.maxValue.toString();
                // check for min value
                if (_this._options.minValue !== null && parseInt(target.value) < _this._options.minValue) target.value = _this._options.minValue.toString();
            };
            this._incrementClickHandler = function() {
                _this.increment();
            };
            this._decrementClickHandler = function() {
                _this.decrement();
            };
            // Add event listener to restrict input to numeric values only
            this._targetEl.addEventListener('input', this._inputHandler);
            if (this._incrementEl) this._incrementEl.addEventListener('click', this._incrementClickHandler);
            if (this._decrementEl) this._decrementEl.addEventListener('click', this._decrementClickHandler);
            this._initialized = true;
        }
    };
    InputCounter.prototype.destroy = function() {
        if (this._targetEl && this._initialized) {
            this._targetEl.removeEventListener('input', this._inputHandler);
            if (this._incrementEl) this._incrementEl.removeEventListener('click', this._incrementClickHandler);
            if (this._decrementEl) this._decrementEl.removeEventListener('click', this._decrementClickHandler);
            this._initialized = false;
        }
    };
    InputCounter.prototype.removeInstance = function() {
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('InputCounter', this._instanceId);
    };
    InputCounter.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    InputCounter.prototype.getCurrentValue = function() {
        return parseInt(this._targetEl.value) || 0;
    };
    InputCounter.prototype.increment = function() {
        // don't increment if the value is already at the maximum value
        if (this._options.maxValue !== null && this.getCurrentValue() >= this._options.maxValue) return;
        this._targetEl.value = (this.getCurrentValue() + 1).toString();
        this._options.onIncrement(this);
    };
    InputCounter.prototype.decrement = function() {
        // don't decrement if the value is already at the minimum value
        if (this._options.minValue !== null && this.getCurrentValue() <= this._options.minValue) return;
        this._targetEl.value = (this.getCurrentValue() - 1).toString();
        this._options.onDecrement(this);
    };
    InputCounter.prototype.updateOnIncrement = function(callback) {
        this._options.onIncrement = callback;
    };
    InputCounter.prototype.updateOnDecrement = function(callback) {
        this._options.onDecrement = callback;
    };
    return InputCounter;
}();
function $df0fc5a35c362eab$export$4352dd3b7a1676ab() {
    document.querySelectorAll('[data-input-counter]').forEach(function($targetEl) {
        var targetId = $targetEl.id;
        var $incrementEl = document.querySelector('[data-input-counter-increment="' + targetId + '"]');
        var $decrementEl = document.querySelector('[data-input-counter-decrement="' + targetId + '"]');
        var minValue = $targetEl.getAttribute('data-input-counter-min');
        var maxValue = $targetEl.getAttribute('data-input-counter-max');
        // check if the target element exists
        if ($targetEl) {
            if (!(0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).instanceExists('InputCounter', $targetEl.getAttribute('id'))) new $df0fc5a35c362eab$var$InputCounter($targetEl, $incrementEl ? $incrementEl : null, $decrementEl ? $decrementEl : null, {
                minValue: minValue ? parseInt(minValue) : null,
                maxValue: maxValue ? parseInt(maxValue) : null
            });
        } else console.error("The target element with id \"".concat(targetId, "\" does not exist. Please check the data-input-counter attribute."));
    });
}
if (typeof window !== 'undefined') {
    window.InputCounter = $df0fc5a35c362eab$var$InputCounter;
    window.initInputCounters = $df0fc5a35c362eab$export$4352dd3b7a1676ab;
}
var $df0fc5a35c362eab$export$2e2bcd8739ae039 = $df0fc5a35c362eab$var$InputCounter;



var $f2c2b31007186eb1$var$__assign = undefined && undefined.__assign || function() {
    $f2c2b31007186eb1$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $f2c2b31007186eb1$var$__assign.apply(this, arguments);
};
var $f2c2b31007186eb1$var$Default = {
    htmlEntities: false,
    contentType: 'input',
    onCopy: function() {}
};
var $f2c2b31007186eb1$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $f2c2b31007186eb1$var$CopyClipboard = /** @class */ function() {
    function CopyClipboard(triggerEl, targetEl, options, instanceOptions) {
        if (triggerEl === void 0) triggerEl = null;
        if (targetEl === void 0) targetEl = null;
        if (options === void 0) options = $f2c2b31007186eb1$var$Default;
        if (instanceOptions === void 0) instanceOptions = $f2c2b31007186eb1$var$DefaultInstanceOptions;
        this._instanceId = instanceOptions.id ? instanceOptions.id : targetEl.id;
        this._triggerEl = triggerEl;
        this._targetEl = targetEl;
        this._options = $f2c2b31007186eb1$var$__assign($f2c2b31007186eb1$var$__assign({}, $f2c2b31007186eb1$var$Default), options);
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('CopyClipboard', this, this._instanceId, instanceOptions.override);
    }
    CopyClipboard.prototype.init = function() {
        var _this = this;
        if (this._targetEl && this._triggerEl && !this._initialized) {
            this._triggerElClickHandler = function() {
                _this.copy();
            };
            // clicking on the trigger element should copy the value of the target element
            if (this._triggerEl) this._triggerEl.addEventListener('click', this._triggerElClickHandler);
            this._initialized = true;
        }
    };
    CopyClipboard.prototype.destroy = function() {
        if (this._triggerEl && this._targetEl && this._initialized) {
            if (this._triggerEl) this._triggerEl.removeEventListener('click', this._triggerElClickHandler);
            this._initialized = false;
        }
    };
    CopyClipboard.prototype.removeInstance = function() {
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('CopyClipboard', this._instanceId);
    };
    CopyClipboard.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    CopyClipboard.prototype.getTargetValue = function() {
        if (this._options.contentType === 'input') return this._targetEl.value;
        if (this._options.contentType === 'innerHTML') return this._targetEl.innerHTML;
        if (this._options.contentType === 'textContent') return this._targetEl.textContent.replace(/\s+/g, ' ').trim();
    };
    CopyClipboard.prototype.copy = function() {
        var textToCopy = this.getTargetValue();
        // Check if HTMLEntities option is enabled
        if (this._options.htmlEntities) // Encode the text using HTML entities
        textToCopy = this.decodeHTML(textToCopy);
        // Create a temporary textarea element
        var tempTextArea = document.createElement('textarea');
        tempTextArea.value = textToCopy;
        document.body.appendChild(tempTextArea);
        // Select the text inside the textarea and copy it to the clipboard
        tempTextArea.select();
        document.execCommand('copy');
        // Remove the temporary textarea
        document.body.removeChild(tempTextArea);
        // Callback function
        this._options.onCopy(this);
        return textToCopy;
    };
    // Function to encode text into HTML entities
    CopyClipboard.prototype.decodeHTML = function(html) {
        var textarea = document.createElement('textarea');
        textarea.innerHTML = html;
        return textarea.textContent;
    };
    CopyClipboard.prototype.updateOnCopyCallback = function(callback) {
        this._options.onCopy = callback;
    };
    return CopyClipboard;
}();
function $f2c2b31007186eb1$export$ad4117886fd6cf74() {
    document.querySelectorAll('[data-copy-to-clipboard-target]').forEach(function($triggerEl) {
        var targetId = $triggerEl.getAttribute('data-copy-to-clipboard-target');
        var $targetEl = document.getElementById(targetId);
        var contentType = $triggerEl.getAttribute('data-copy-to-clipboard-content-type');
        var htmlEntities = $triggerEl.getAttribute('data-copy-to-clipboard-html-entities');
        // check if the target element exists
        if ($targetEl) {
            if (!(0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).instanceExists('CopyClipboard', $targetEl.getAttribute('id'))) new $f2c2b31007186eb1$var$CopyClipboard($triggerEl, $targetEl, {
                htmlEntities: htmlEntities && htmlEntities === 'true' ? true : $f2c2b31007186eb1$var$Default.htmlEntities,
                contentType: contentType ? contentType : $f2c2b31007186eb1$var$Default.contentType
            });
        } else console.error("The target element with id \"".concat(targetId, "\" does not exist. Please check the data-copy-to-clipboard-target attribute."));
    });
}
if (typeof window !== 'undefined') {
    window.CopyClipboard = $f2c2b31007186eb1$var$CopyClipboard;
    window.initClipboards = $f2c2b31007186eb1$export$ad4117886fd6cf74;
}
var $f2c2b31007186eb1$export$2e2bcd8739ae039 = $f2c2b31007186eb1$var$CopyClipboard;



function $a18c32cb29f9cbaf$var$_arrayLikeToArray(r, a) {
    (null == a || a > r.length) && (a = r.length);
    for(var e = 0, n = Array(a); e < a; e++)n[e] = r[e];
    return n;
}
function $a18c32cb29f9cbaf$var$_arrayWithHoles(r) {
    if (Array.isArray(r)) return r;
}
function $a18c32cb29f9cbaf$var$_arrayWithoutHoles(r) {
    if (Array.isArray(r)) return $a18c32cb29f9cbaf$var$_arrayLikeToArray(r);
}
function $a18c32cb29f9cbaf$var$_assertThisInitialized(e) {
    if (void 0 === e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
    return e;
}
function $a18c32cb29f9cbaf$var$_callSuper(t, o, e) {
    return o = $a18c32cb29f9cbaf$var$_getPrototypeOf(o), $a18c32cb29f9cbaf$var$_possibleConstructorReturn(t, $a18c32cb29f9cbaf$var$_isNativeReflectConstruct() ? Reflect.construct(o, e || [], $a18c32cb29f9cbaf$var$_getPrototypeOf(t).constructor) : o.apply(t, e));
}
function $a18c32cb29f9cbaf$var$_classCallCheck(a, n) {
    if (!(a instanceof n)) throw new TypeError("Cannot call a class as a function");
}
function $a18c32cb29f9cbaf$var$_defineProperties(e, r) {
    for(var t = 0; t < r.length; t++){
        var o = r[t];
        o.enumerable = o.enumerable || !1, o.configurable = !0, "value" in o && (o.writable = !0), Object.defineProperty(e, $a18c32cb29f9cbaf$var$_toPropertyKey(o.key), o);
    }
}
function $a18c32cb29f9cbaf$var$_createClass(e, r, t) {
    return r && $a18c32cb29f9cbaf$var$_defineProperties(e.prototype, r), t && $a18c32cb29f9cbaf$var$_defineProperties(e, t), Object.defineProperty(e, "prototype", {
        writable: !1
    }), e;
}
function $a18c32cb29f9cbaf$var$_get() {
    return $a18c32cb29f9cbaf$var$_get = "undefined" != typeof Reflect && Reflect.get ? Reflect.get.bind() : function(e, t, r) {
        var p = $a18c32cb29f9cbaf$var$_superPropBase(e, t);
        if (p) {
            var n = Object.getOwnPropertyDescriptor(p, t);
            return n.get ? n.get.call(arguments.length < 3 ? e : r) : n.value;
        }
    }, $a18c32cb29f9cbaf$var$_get.apply(null, arguments);
}
function $a18c32cb29f9cbaf$var$_getPrototypeOf(t) {
    return $a18c32cb29f9cbaf$var$_getPrototypeOf = Object.setPrototypeOf ? Object.getPrototypeOf.bind() : function(t) {
        return t.__proto__ || Object.getPrototypeOf(t);
    }, $a18c32cb29f9cbaf$var$_getPrototypeOf(t);
}
function $a18c32cb29f9cbaf$var$_inherits(t, e) {
    if ("function" != typeof e && null !== e) throw new TypeError("Super expression must either be null or a function");
    t.prototype = Object.create(e && e.prototype, {
        constructor: {
            value: t,
            writable: !0,
            configurable: !0
        }
    }), Object.defineProperty(t, "prototype", {
        writable: !1
    }), e && $a18c32cb29f9cbaf$var$_setPrototypeOf(t, e);
}
function $a18c32cb29f9cbaf$var$_isNativeReflectConstruct() {
    try {
        var t = !Boolean.prototype.valueOf.call(Reflect.construct(Boolean, [], function() {}));
    } catch (t) {}
    return ($a18c32cb29f9cbaf$var$_isNativeReflectConstruct = function() {
        return !!t;
    })();
}
function $a18c32cb29f9cbaf$var$_iterableToArray(r) {
    if ("undefined" != typeof Symbol && null != r[Symbol.iterator] || null != r["@@iterator"]) return Array.from(r);
}
function $a18c32cb29f9cbaf$var$_iterableToArrayLimit(r, l) {
    var t = null == r ? null : "undefined" != typeof Symbol && r[Symbol.iterator] || r["@@iterator"];
    if (null != t) {
        var e, n, i, u, a = [], f = !0, o = !1;
        try {
            if (i = (t = t.call(r)).next, 0 === l) {
                if (Object(t) !== t) return;
                f = !1;
            } else for(; !(f = (e = i.call(t)).done) && (a.push(e.value), a.length !== l); f = !0);
        } catch (r) {
            o = !0, n = r;
        } finally{
            try {
                if (!f && null != t.return && (u = t.return(), Object(u) !== u)) return;
            } finally{
                if (o) throw n;
            }
        }
        return a;
    }
}
function $a18c32cb29f9cbaf$var$_nonIterableRest() {
    throw new TypeError("Invalid attempt to destructure non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method.");
}
function $a18c32cb29f9cbaf$var$_nonIterableSpread() {
    throw new TypeError("Invalid attempt to spread non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method.");
}
function $a18c32cb29f9cbaf$var$_possibleConstructorReturn(t, e) {
    if (e && ("object" == typeof e || "function" == typeof e)) return e;
    if (void 0 !== e) throw new TypeError("Derived constructors may only return object or undefined");
    return $a18c32cb29f9cbaf$var$_assertThisInitialized(t);
}
function $a18c32cb29f9cbaf$var$_setPrototypeOf(t, e) {
    return $a18c32cb29f9cbaf$var$_setPrototypeOf = Object.setPrototypeOf ? Object.setPrototypeOf.bind() : function(t, e) {
        return t.__proto__ = e, t;
    }, $a18c32cb29f9cbaf$var$_setPrototypeOf(t, e);
}
function $a18c32cb29f9cbaf$var$_slicedToArray(r, e) {
    return $a18c32cb29f9cbaf$var$_arrayWithHoles(r) || $a18c32cb29f9cbaf$var$_iterableToArrayLimit(r, e) || $a18c32cb29f9cbaf$var$_unsupportedIterableToArray(r, e) || $a18c32cb29f9cbaf$var$_nonIterableRest();
}
function $a18c32cb29f9cbaf$var$_superPropBase(t, o) {
    for(; !({}).hasOwnProperty.call(t, o) && null !== (t = $a18c32cb29f9cbaf$var$_getPrototypeOf(t)););
    return t;
}
function $a18c32cb29f9cbaf$var$_toConsumableArray(r) {
    return $a18c32cb29f9cbaf$var$_arrayWithoutHoles(r) || $a18c32cb29f9cbaf$var$_iterableToArray(r) || $a18c32cb29f9cbaf$var$_unsupportedIterableToArray(r) || $a18c32cb29f9cbaf$var$_nonIterableSpread();
}
function $a18c32cb29f9cbaf$var$_toPrimitive(t, r) {
    if ("object" != typeof t || !t) return t;
    var e = t[Symbol.toPrimitive];
    if (void 0 !== e) {
        var i = e.call(t, r || "default");
        if ("object" != typeof i) return i;
        throw new TypeError("@@toPrimitive must return a primitive value.");
    }
    return ("string" === r ? String : Number)(t);
}
function $a18c32cb29f9cbaf$var$_toPropertyKey(t) {
    var i = $a18c32cb29f9cbaf$var$_toPrimitive(t, "string");
    return "symbol" == typeof i ? i : i + "";
}
function $a18c32cb29f9cbaf$var$_typeof(o) {
    "@babel/helpers - typeof";
    return $a18c32cb29f9cbaf$var$_typeof = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function(o) {
        return typeof o;
    } : function(o) {
        return o && "function" == typeof Symbol && o.constructor === Symbol && o !== Symbol.prototype ? "symbol" : typeof o;
    }, $a18c32cb29f9cbaf$var$_typeof(o);
}
function $a18c32cb29f9cbaf$var$_unsupportedIterableToArray(r, a) {
    if (r) {
        if ("string" == typeof r) return $a18c32cb29f9cbaf$var$_arrayLikeToArray(r, a);
        var t = ({}).toString.call(r).slice(8, -1);
        return "Object" === t && r.constructor && (t = r.constructor.name), "Map" === t || "Set" === t ? Array.from(r) : "Arguments" === t || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(t) ? $a18c32cb29f9cbaf$var$_arrayLikeToArray(r, a) : void 0;
    }
}
function $a18c32cb29f9cbaf$var$hasProperty(obj, prop) {
    return Object.prototype.hasOwnProperty.call(obj, prop);
}
function $a18c32cb29f9cbaf$var$lastItemOf(arr) {
    return arr[arr.length - 1];
}
// push only the items not included in the array
function $a18c32cb29f9cbaf$var$pushUnique(arr) {
    for(var _len = arguments.length, items = new Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++)items[_key - 1] = arguments[_key];
    items.forEach(function(item) {
        if (arr.includes(item)) return;
        arr.push(item);
    });
    return arr;
}
function $a18c32cb29f9cbaf$var$stringToArray(str, separator) {
    // convert empty string to an empty array
    return str ? str.split(separator) : [];
}
function $a18c32cb29f9cbaf$var$isInRange(testVal, min, max) {
    var minOK = min === undefined || testVal >= min;
    var maxOK = max === undefined || testVal <= max;
    return minOK && maxOK;
}
function $a18c32cb29f9cbaf$var$limitToRange(val, min, max) {
    if (val < min) return min;
    if (val > max) return max;
    return val;
}
function $a18c32cb29f9cbaf$var$createTagRepeat(tagName, repeat) {
    var attributes = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : {};
    var index = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : 0;
    var html = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : '';
    var openTagSrc = Object.keys(attributes).reduce(function(src, attr) {
        var val = attributes[attr];
        if (typeof val === 'function') val = val(index);
        return "".concat(src, " ").concat(attr, "=\"").concat(val, "\"");
    }, tagName);
    html += "<".concat(openTagSrc, "></").concat(tagName, ">");
    var next = index + 1;
    return next < repeat ? $a18c32cb29f9cbaf$var$createTagRepeat(tagName, repeat, attributes, next, html) : html;
}
// Remove the spacing surrounding tags for HTML parser not to create text nodes
// before/after elements
function $a18c32cb29f9cbaf$var$optimizeTemplateHTML(html) {
    return html.replace(/>\s+/g, '>').replace(/\s+</, '<');
}
function $a18c32cb29f9cbaf$var$stripTime(timeValue) {
    return new Date(timeValue).setHours(0, 0, 0, 0);
}
function $a18c32cb29f9cbaf$var$today() {
    return new Date().setHours(0, 0, 0, 0);
}
// Get the time value of the start of given date or year, month and day
function $a18c32cb29f9cbaf$var$dateValue() {
    switch(arguments.length){
        case 0:
            return $a18c32cb29f9cbaf$var$today();
        case 1:
            return $a18c32cb29f9cbaf$var$stripTime(arguments.length <= 0 ? undefined : arguments[0]);
    }
    // use setFullYear() to keep 2-digit year from being mapped to 1900-1999
    var newDate = new Date(0);
    newDate.setFullYear.apply(newDate, arguments);
    return newDate.setHours(0, 0, 0, 0);
}
function $a18c32cb29f9cbaf$var$addDays(date, amount) {
    var newDate = new Date(date);
    return newDate.setDate(newDate.getDate() + amount);
}
function $a18c32cb29f9cbaf$var$addWeeks(date, amount) {
    return $a18c32cb29f9cbaf$var$addDays(date, amount * 7);
}
function $a18c32cb29f9cbaf$var$addMonths(date, amount) {
    // If the day of the date is not in the new month, the last day of the new
    // month will be returned. e.g. Jan 31 + 1 month → Feb 28 (not Mar 03)
    var newDate = new Date(date);
    var monthsToSet = newDate.getMonth() + amount;
    var expectedMonth = monthsToSet % 12;
    if (expectedMonth < 0) expectedMonth += 12;
    var time = newDate.setMonth(monthsToSet);
    return newDate.getMonth() !== expectedMonth ? newDate.setDate(0) : time;
}
function $a18c32cb29f9cbaf$var$addYears(date, amount) {
    // If the date is Feb 29 and the new year is not a leap year, Feb 28 of the
    // new year will be returned.
    var newDate = new Date(date);
    var expectedMonth = newDate.getMonth();
    var time = newDate.setFullYear(newDate.getFullYear() + amount);
    return expectedMonth === 1 && newDate.getMonth() === 2 ? newDate.setDate(0) : time;
}
// Calculate the distance bettwen 2 days of the week
function $a18c32cb29f9cbaf$var$dayDiff(day, from) {
    return (day - from + 7) % 7;
}
// Get the date of the specified day of the week of given base date
function $a18c32cb29f9cbaf$var$dayOfTheWeekOf(baseDate, dayOfWeek) {
    var weekStart = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : 0;
    var baseDay = new Date(baseDate).getDay();
    return $a18c32cb29f9cbaf$var$addDays(baseDate, $a18c32cb29f9cbaf$var$dayDiff(dayOfWeek, weekStart) - $a18c32cb29f9cbaf$var$dayDiff(baseDay, weekStart));
}
// Get the ISO week of a date
function $a18c32cb29f9cbaf$var$getWeek(date) {
    // start of ISO week is Monday
    var thuOfTheWeek = $a18c32cb29f9cbaf$var$dayOfTheWeekOf(date, 4, 1);
    // 1st week == the week where the 4th of January is in
    var firstThu = $a18c32cb29f9cbaf$var$dayOfTheWeekOf(new Date(thuOfTheWeek).setMonth(0, 4), 4, 1);
    return Math.round((thuOfTheWeek - firstThu) / 604800000) + 1;
}
// Get the start year of the period of years that includes given date
// years: length of the year period
function $a18c32cb29f9cbaf$var$startOfYearPeriod(date, years) {
    /* @see https://en.wikipedia.org/wiki/Year_zero#ISO_8601 */ var year = new Date(date).getFullYear();
    return Math.floor(year / years) * years;
}
// pattern for format parts
var $a18c32cb29f9cbaf$var$reFormatTokens = /dd?|DD?|mm?|MM?|yy?(?:yy)?/;
// pattern for non date parts
var $a18c32cb29f9cbaf$var$reNonDateParts = /[\s!-/:-@[-`{-~年月日]+/;
// cache for persed formats
var $a18c32cb29f9cbaf$var$knownFormats = {};
// parse funtions for date parts
var $a18c32cb29f9cbaf$var$parseFns = {
    y: function y(date, year) {
        return new Date(date).setFullYear(parseInt(year, 10));
    },
    m: function m(date, month, locale) {
        var newDate = new Date(date);
        var monthIndex = parseInt(month, 10) - 1;
        if (isNaN(monthIndex)) {
            if (!month) return NaN;
            var monthName = month.toLowerCase();
            var compareNames = function compareNames(name) {
                return name.toLowerCase().startsWith(monthName);
            };
            // compare with both short and full names because some locales have periods
            // in the short names (not equal to the first X letters of the full names)
            monthIndex = locale.monthsShort.findIndex(compareNames);
            if (monthIndex < 0) monthIndex = locale.months.findIndex(compareNames);
            if (monthIndex < 0) return NaN;
        }
        newDate.setMonth(monthIndex);
        return newDate.getMonth() !== $a18c32cb29f9cbaf$var$normalizeMonth(monthIndex) ? newDate.setDate(0) : newDate.getTime();
    },
    d: function d(date, day) {
        return new Date(date).setDate(parseInt(day, 10));
    }
};
// format functions for date parts
var $a18c32cb29f9cbaf$var$formatFns = {
    d: function d(date) {
        return date.getDate();
    },
    dd: function dd(date) {
        return $a18c32cb29f9cbaf$var$padZero(date.getDate(), 2);
    },
    D: function D(date, locale) {
        return locale.daysShort[date.getDay()];
    },
    DD: function DD(date, locale) {
        return locale.days[date.getDay()];
    },
    m: function m(date) {
        return date.getMonth() + 1;
    },
    mm: function mm(date) {
        return $a18c32cb29f9cbaf$var$padZero(date.getMonth() + 1, 2);
    },
    M: function M(date, locale) {
        return locale.monthsShort[date.getMonth()];
    },
    MM: function MM(date, locale) {
        return locale.months[date.getMonth()];
    },
    y: function y(date) {
        return date.getFullYear();
    },
    yy: function yy(date) {
        return $a18c32cb29f9cbaf$var$padZero(date.getFullYear(), 2).slice(-2);
    },
    yyyy: function yyyy(date) {
        return $a18c32cb29f9cbaf$var$padZero(date.getFullYear(), 4);
    }
};
// get month index in normal range (0 - 11) from any number
function $a18c32cb29f9cbaf$var$normalizeMonth(monthIndex) {
    return monthIndex > -1 ? monthIndex % 12 : $a18c32cb29f9cbaf$var$normalizeMonth(monthIndex + 12);
}
function $a18c32cb29f9cbaf$var$padZero(num, length) {
    return num.toString().padStart(length, '0');
}
function $a18c32cb29f9cbaf$var$parseFormatString(format) {
    if (typeof format !== 'string') throw new Error("Invalid date format.");
    if (format in $a18c32cb29f9cbaf$var$knownFormats) return $a18c32cb29f9cbaf$var$knownFormats[format];
    // sprit the format string into parts and seprators
    var separators = format.split($a18c32cb29f9cbaf$var$reFormatTokens);
    var parts = format.match(new RegExp($a18c32cb29f9cbaf$var$reFormatTokens, 'g'));
    if (separators.length === 0 || !parts) throw new Error("Invalid date format.");
    // collect format functions used in the format
    var partFormatters = parts.map(function(token) {
        return $a18c32cb29f9cbaf$var$formatFns[token];
    });
    // collect parse function keys used in the format
    // iterate over parseFns' keys in order to keep the order of the keys.
    var partParserKeys = Object.keys($a18c32cb29f9cbaf$var$parseFns).reduce(function(keys, key) {
        var token = parts.find(function(part) {
            return part[0] !== 'D' && part[0].toLowerCase() === key;
        });
        if (token) keys.push(key);
        return keys;
    }, []);
    return $a18c32cb29f9cbaf$var$knownFormats[format] = {
        parser: function parser(dateStr, locale) {
            var dateParts = dateStr.split($a18c32cb29f9cbaf$var$reNonDateParts).reduce(function(dtParts, part, index) {
                if (part.length > 0 && parts[index]) {
                    var token = parts[index][0];
                    if (token === 'M') dtParts.m = part;
                    else if (token !== 'D') dtParts[token] = part;
                }
                return dtParts;
            }, {});
            // iterate over partParserkeys so that the parsing is made in the oder
            // of year, month and day to prevent the day parser from correcting last
            // day of month wrongly
            return partParserKeys.reduce(function(origDate, key) {
                var newDate = $a18c32cb29f9cbaf$var$parseFns[key](origDate, dateParts[key], locale);
                // ingnore the part failed to parse
                return isNaN(newDate) ? origDate : newDate;
            }, $a18c32cb29f9cbaf$var$today());
        },
        formatter: function formatter(date, locale) {
            var dateStr = partFormatters.reduce(function(str, fn, index) {
                return str += "".concat(separators[index]).concat(fn(date, locale));
            }, '');
            // separators' length is always parts' length + 1,
            return dateStr += $a18c32cb29f9cbaf$var$lastItemOf(separators);
        }
    };
}
function $a18c32cb29f9cbaf$var$parseDate(dateStr, format, locale) {
    if (dateStr instanceof Date || typeof dateStr === 'number') {
        var date = $a18c32cb29f9cbaf$var$stripTime(dateStr);
        return isNaN(date) ? undefined : date;
    }
    if (!dateStr) return undefined;
    if (dateStr === 'today') return $a18c32cb29f9cbaf$var$today();
    if (format && format.toValue) {
        var _date = format.toValue(dateStr, format, locale);
        return isNaN(_date) ? undefined : $a18c32cb29f9cbaf$var$stripTime(_date);
    }
    return $a18c32cb29f9cbaf$var$parseFormatString(format).parser(dateStr, locale);
}
function $a18c32cb29f9cbaf$var$formatDate(date, format, locale) {
    if (isNaN(date) || !date && date !== 0) return '';
    var dateObj = typeof date === 'number' ? new Date(date) : date;
    if (format.toDisplay) return format.toDisplay(dateObj, format, locale);
    return $a18c32cb29f9cbaf$var$parseFormatString(format).formatter(dateObj, locale);
}
var $a18c32cb29f9cbaf$var$listenerRegistry = new WeakMap();
var $a18c32cb29f9cbaf$var$_EventTarget$prototyp = EventTarget.prototype, $a18c32cb29f9cbaf$var$addEventListener = $a18c32cb29f9cbaf$var$_EventTarget$prototyp.addEventListener, $a18c32cb29f9cbaf$var$removeEventListener = $a18c32cb29f9cbaf$var$_EventTarget$prototyp.removeEventListener;
// Register event listeners to a key object
// listeners: array of listener definitions;
//   - each definition must be a flat array of event target and the arguments
//     used to call addEventListener() on the target
function $a18c32cb29f9cbaf$var$registerListeners(keyObj, listeners) {
    var registered = $a18c32cb29f9cbaf$var$listenerRegistry.get(keyObj);
    if (!registered) {
        registered = [];
        $a18c32cb29f9cbaf$var$listenerRegistry.set(keyObj, registered);
    }
    listeners.forEach(function(listener) {
        $a18c32cb29f9cbaf$var$addEventListener.call.apply($a18c32cb29f9cbaf$var$addEventListener, $a18c32cb29f9cbaf$var$_toConsumableArray(listener));
        registered.push(listener);
    });
}
function $a18c32cb29f9cbaf$var$unregisterListeners(keyObj) {
    var listeners = $a18c32cb29f9cbaf$var$listenerRegistry.get(keyObj);
    if (!listeners) return;
    listeners.forEach(function(listener) {
        $a18c32cb29f9cbaf$var$removeEventListener.call.apply($a18c32cb29f9cbaf$var$removeEventListener, $a18c32cb29f9cbaf$var$_toConsumableArray(listener));
    });
    $a18c32cb29f9cbaf$var$listenerRegistry["delete"](keyObj);
}
// Event.composedPath() polyfill for Edge
// based on https://gist.github.com/kleinfreund/e9787d73776c0e3750dcfcdc89f100ec
if (!Event.prototype.composedPath) {
    var $a18c32cb29f9cbaf$var$getComposedPath = function getComposedPath(node) {
        var path = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : [];
        path.push(node);
        var parent;
        if (node.parentNode) parent = node.parentNode;
        else if (node.host) // ShadowRoot
        parent = node.host;
        else if (node.defaultView) // Document
        parent = node.defaultView;
        return parent ? getComposedPath(parent, path) : path;
    };
    Event.prototype.composedPath = function() {
        return $a18c32cb29f9cbaf$var$getComposedPath(this.target);
    };
}
function $a18c32cb29f9cbaf$var$findFromPath(path, criteria, currentTarget) {
    var index = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : 0;
    var el = path[index];
    if (criteria(el)) return el;
    else if (el === currentTarget || !el.parentElement) // stop when reaching currentTarget or <html>
    return;
    return $a18c32cb29f9cbaf$var$findFromPath(path, criteria, currentTarget, index + 1);
}
// Search for the actual target of a delegated event
function $a18c32cb29f9cbaf$var$findElementInEventPath(ev, selector) {
    var criteria = typeof selector === 'function' ? selector : function(el) {
        return el.matches(selector);
    };
    return $a18c32cb29f9cbaf$var$findFromPath(ev.composedPath(), criteria, ev.currentTarget);
}
// default locales
var $a18c32cb29f9cbaf$var$locales = {
    en: {
        days: [
            "Sunday",
            "Monday",
            "Tuesday",
            "Wednesday",
            "Thursday",
            "Friday",
            "Saturday"
        ],
        daysShort: [
            "Sun",
            "Mon",
            "Tue",
            "Wed",
            "Thu",
            "Fri",
            "Sat"
        ],
        daysMin: [
            "Su",
            "Mo",
            "Tu",
            "We",
            "Th",
            "Fr",
            "Sa"
        ],
        months: [
            "January",
            "February",
            "March",
            "April",
            "May",
            "June",
            "July",
            "August",
            "September",
            "October",
            "November",
            "December"
        ],
        monthsShort: [
            "Jan",
            "Feb",
            "Mar",
            "Apr",
            "May",
            "Jun",
            "Jul",
            "Aug",
            "Sep",
            "Oct",
            "Nov",
            "Dec"
        ],
        today: "Today",
        clear: "Clear",
        titleFormat: "MM y"
    }
};
// config options updatable by setOptions() and their default values
var $a18c32cb29f9cbaf$var$defaultOptions = {
    autohide: false,
    beforeShowDay: null,
    beforeShowDecade: null,
    beforeShowMonth: null,
    beforeShowYear: null,
    calendarWeeks: false,
    clearBtn: false,
    dateDelimiter: ',',
    datesDisabled: [],
    daysOfWeekDisabled: [],
    daysOfWeekHighlighted: [],
    defaultViewDate: undefined,
    // placeholder, defaults to today() by the program
    disableTouchKeyboard: false,
    format: 'mm/dd/yyyy',
    language: 'en',
    maxDate: null,
    maxNumberOfDates: 1,
    maxView: 3,
    minDate: null,
    nextArrow: '<svg class="w-4 h-4 rtl:rotate-180" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 10"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M1 5h12m0 0L9 1m4 4L9 9"/></svg>',
    orientation: 'auto',
    pickLevel: 0,
    prevArrow: '<svg class="w-4 h-4 rtl:rotate-180" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 10"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5H1m0 0 4 4M1 5l4-4"/></svg>',
    showDaysOfWeek: true,
    showOnClick: true,
    showOnFocus: true,
    startView: 0,
    title: '',
    todayBtn: false,
    todayBtnMode: 0,
    todayHighlight: false,
    updateOnBlur: true,
    weekStart: 0
};
var $a18c32cb29f9cbaf$var$range = null;
function $a18c32cb29f9cbaf$var$parseHTML(html) {
    if ($a18c32cb29f9cbaf$var$range == null) $a18c32cb29f9cbaf$var$range = document.createRange();
    return $a18c32cb29f9cbaf$var$range.createContextualFragment(html);
}
function $a18c32cb29f9cbaf$var$hideElement(el) {
    if (el.style.display === 'none') return;
    // back up the existing display setting in data-style-display
    if (el.style.display) el.dataset.styleDisplay = el.style.display;
    el.style.display = 'none';
}
function $a18c32cb29f9cbaf$var$showElement(el) {
    if (el.style.display !== 'none') return;
    if (el.dataset.styleDisplay) {
        // restore backed-up dispay property
        el.style.display = el.dataset.styleDisplay;
        delete el.dataset.styleDisplay;
    } else el.style.display = '';
}
function $a18c32cb29f9cbaf$var$emptyChildNodes(el) {
    if (el.firstChild) {
        el.removeChild(el.firstChild);
        $a18c32cb29f9cbaf$var$emptyChildNodes(el);
    }
}
function $a18c32cb29f9cbaf$var$replaceChildNodes(el, newChildNodes) {
    $a18c32cb29f9cbaf$var$emptyChildNodes(el);
    if (newChildNodes instanceof DocumentFragment) el.appendChild(newChildNodes);
    else if (typeof newChildNodes === 'string') el.appendChild($a18c32cb29f9cbaf$var$parseHTML(newChildNodes));
    else if (typeof newChildNodes.forEach === 'function') newChildNodes.forEach(function(node) {
        el.appendChild(node);
    });
}
var $a18c32cb29f9cbaf$var$defaultLang = $a18c32cb29f9cbaf$var$defaultOptions.language, $a18c32cb29f9cbaf$var$defaultFormat = $a18c32cb29f9cbaf$var$defaultOptions.format, $a18c32cb29f9cbaf$var$defaultWeekStart = $a18c32cb29f9cbaf$var$defaultOptions.weekStart;
// Reducer function to filter out invalid day-of-week from the input
function $a18c32cb29f9cbaf$var$sanitizeDOW(dow, day) {
    return dow.length < 6 && day >= 0 && day < 7 ? $a18c32cb29f9cbaf$var$pushUnique(dow, day) : dow;
}
function $a18c32cb29f9cbaf$var$calcEndOfWeek(startOfWeek) {
    return (startOfWeek + 6) % 7;
}
// validate input date. if invalid, fallback to the original value
function $a18c32cb29f9cbaf$var$validateDate(value, format, locale, origValue) {
    var date = $a18c32cb29f9cbaf$var$parseDate(value, format, locale);
    return date !== undefined ? date : origValue;
}
// Validate viewId. if invalid, fallback to the original value
function $a18c32cb29f9cbaf$var$validateViewId(value, origValue) {
    var max = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : 3;
    var viewId = parseInt(value, 10);
    return viewId >= 0 && viewId <= max ? viewId : origValue;
}
// Create Datepicker configuration to set
function $a18c32cb29f9cbaf$var$processOptions(options, datepicker) {
    var inOpts = Object.assign({}, options);
    var config = {};
    var locales = datepicker.constructor.locales;
    var _ref = datepicker.config || {}, format = _ref.format, language = _ref.language, locale = _ref.locale, maxDate = _ref.maxDate, maxView = _ref.maxView, minDate = _ref.minDate, pickLevel = _ref.pickLevel, startView = _ref.startView, weekStart = _ref.weekStart;
    if (inOpts.language) {
        var lang;
        if (inOpts.language !== language) {
            if (locales[inOpts.language]) lang = inOpts.language;
            else {
                // Check if langauge + region tag can fallback to the one without
                // region (e.g. fr-CA → fr)
                lang = inOpts.language.split('-')[0];
                if (locales[lang] === undefined) lang = false;
            }
        }
        delete inOpts.language;
        if (lang) {
            language = config.language = lang;
            // update locale as well when updating language
            var origLocale = locale || locales[$a18c32cb29f9cbaf$var$defaultLang];
            // use default language's properties for the fallback
            locale = Object.assign({
                format: $a18c32cb29f9cbaf$var$defaultFormat,
                weekStart: $a18c32cb29f9cbaf$var$defaultWeekStart
            }, locales[$a18c32cb29f9cbaf$var$defaultLang]);
            if (language !== $a18c32cb29f9cbaf$var$defaultLang) Object.assign(locale, locales[language]);
            config.locale = locale;
            // if format and/or weekStart are the same as old locale's defaults,
            // update them to new locale's defaults
            if (format === origLocale.format) format = config.format = locale.format;
            if (weekStart === origLocale.weekStart) {
                weekStart = config.weekStart = locale.weekStart;
                config.weekEnd = $a18c32cb29f9cbaf$var$calcEndOfWeek(locale.weekStart);
            }
        }
    }
    if (inOpts.format) {
        var hasToDisplay = typeof inOpts.format.toDisplay === 'function';
        var hasToValue = typeof inOpts.format.toValue === 'function';
        var validFormatString = $a18c32cb29f9cbaf$var$reFormatTokens.test(inOpts.format);
        if (hasToDisplay && hasToValue || validFormatString) format = config.format = inOpts.format;
        delete inOpts.format;
    }
    //*** dates ***//
    // while min and maxDate for "no limit" in the options are better to be null
    // (especially when updating), the ones in the config have to be undefined
    // because null is treated as 0 (= unix epoch) when comparing with time value
    var minDt = minDate;
    var maxDt = maxDate;
    if (inOpts.minDate !== undefined) {
        minDt = inOpts.minDate === null ? $a18c32cb29f9cbaf$var$dateValue(0, 0, 1) // set 0000-01-01 to prevent negative values for year
         : $a18c32cb29f9cbaf$var$validateDate(inOpts.minDate, format, locale, minDt);
        delete inOpts.minDate;
    }
    if (inOpts.maxDate !== undefined) {
        maxDt = inOpts.maxDate === null ? undefined : $a18c32cb29f9cbaf$var$validateDate(inOpts.maxDate, format, locale, maxDt);
        delete inOpts.maxDate;
    }
    if (maxDt < minDt) {
        minDate = config.minDate = maxDt;
        maxDate = config.maxDate = minDt;
    } else {
        if (minDate !== minDt) minDate = config.minDate = minDt;
        if (maxDate !== maxDt) maxDate = config.maxDate = maxDt;
    }
    if (inOpts.datesDisabled) {
        config.datesDisabled = inOpts.datesDisabled.reduce(function(dates, dt) {
            var date = $a18c32cb29f9cbaf$var$parseDate(dt, format, locale);
            return date !== undefined ? $a18c32cb29f9cbaf$var$pushUnique(dates, date) : dates;
        }, []);
        delete inOpts.datesDisabled;
    }
    if (inOpts.defaultViewDate !== undefined) {
        var viewDate = $a18c32cb29f9cbaf$var$parseDate(inOpts.defaultViewDate, format, locale);
        if (viewDate !== undefined) config.defaultViewDate = viewDate;
        delete inOpts.defaultViewDate;
    }
    //*** days of week ***//
    if (inOpts.weekStart !== undefined) {
        var wkStart = Number(inOpts.weekStart) % 7;
        if (!isNaN(wkStart)) {
            weekStart = config.weekStart = wkStart;
            config.weekEnd = $a18c32cb29f9cbaf$var$calcEndOfWeek(wkStart);
        }
        delete inOpts.weekStart;
    }
    if (inOpts.daysOfWeekDisabled) {
        config.daysOfWeekDisabled = inOpts.daysOfWeekDisabled.reduce($a18c32cb29f9cbaf$var$sanitizeDOW, []);
        delete inOpts.daysOfWeekDisabled;
    }
    if (inOpts.daysOfWeekHighlighted) {
        config.daysOfWeekHighlighted = inOpts.daysOfWeekHighlighted.reduce($a18c32cb29f9cbaf$var$sanitizeDOW, []);
        delete inOpts.daysOfWeekHighlighted;
    }
    //*** multi date ***//
    if (inOpts.maxNumberOfDates !== undefined) {
        var maxNumberOfDates = parseInt(inOpts.maxNumberOfDates, 10);
        if (maxNumberOfDates >= 0) {
            config.maxNumberOfDates = maxNumberOfDates;
            config.multidate = maxNumberOfDates !== 1;
        }
        delete inOpts.maxNumberOfDates;
    }
    if (inOpts.dateDelimiter) {
        config.dateDelimiter = String(inOpts.dateDelimiter);
        delete inOpts.dateDelimiter;
    }
    //*** pick level & view ***//
    var newPickLevel = pickLevel;
    if (inOpts.pickLevel !== undefined) {
        newPickLevel = $a18c32cb29f9cbaf$var$validateViewId(inOpts.pickLevel, 2);
        delete inOpts.pickLevel;
    }
    if (newPickLevel !== pickLevel) pickLevel = config.pickLevel = newPickLevel;
    var newMaxView = maxView;
    if (inOpts.maxView !== undefined) {
        newMaxView = $a18c32cb29f9cbaf$var$validateViewId(inOpts.maxView, maxView);
        delete inOpts.maxView;
    }
    // ensure max view >= pick level
    newMaxView = pickLevel > newMaxView ? pickLevel : newMaxView;
    if (newMaxView !== maxView) maxView = config.maxView = newMaxView;
    var newStartView = startView;
    if (inOpts.startView !== undefined) {
        newStartView = $a18c32cb29f9cbaf$var$validateViewId(inOpts.startView, newStartView);
        delete inOpts.startView;
    }
    // ensure pick level <= start view <= max view
    if (newStartView < pickLevel) newStartView = pickLevel;
    else if (newStartView > maxView) newStartView = maxView;
    if (newStartView !== startView) config.startView = newStartView;
    //*** template ***//
    if (inOpts.prevArrow) {
        var prevArrow = $a18c32cb29f9cbaf$var$parseHTML(inOpts.prevArrow);
        if (prevArrow.childNodes.length > 0) config.prevArrow = prevArrow.childNodes;
        delete inOpts.prevArrow;
    }
    if (inOpts.nextArrow) {
        var nextArrow = $a18c32cb29f9cbaf$var$parseHTML(inOpts.nextArrow);
        if (nextArrow.childNodes.length > 0) config.nextArrow = nextArrow.childNodes;
        delete inOpts.nextArrow;
    }
    //*** misc ***//
    if (inOpts.disableTouchKeyboard !== undefined) {
        config.disableTouchKeyboard = 'ontouchstart' in document && !!inOpts.disableTouchKeyboard;
        delete inOpts.disableTouchKeyboard;
    }
    if (inOpts.orientation) {
        var orientation = inOpts.orientation.toLowerCase().split(/\s+/g);
        config.orientation = {
            x: orientation.find(function(x) {
                return x === 'left' || x === 'right';
            }) || 'auto',
            y: orientation.find(function(y) {
                return y === 'top' || y === 'bottom';
            }) || 'auto'
        };
        delete inOpts.orientation;
    }
    if (inOpts.todayBtnMode !== undefined) {
        switch(inOpts.todayBtnMode){
            case 0:
            case 1:
                config.todayBtnMode = inOpts.todayBtnMode;
        }
        delete inOpts.todayBtnMode;
    }
    //*** copy the rest ***//
    Object.keys(inOpts).forEach(function(key) {
        if (inOpts[key] !== undefined && $a18c32cb29f9cbaf$var$hasProperty($a18c32cb29f9cbaf$var$defaultOptions, key)) config[key] = inOpts[key];
    });
    return config;
}
var $a18c32cb29f9cbaf$var$pickerTemplate = $a18c32cb29f9cbaf$var$optimizeTemplateHTML("<div class=\"datepicker hidden\">\n  <div class=\"datepicker-picker inline-block rounded-base bg-neutral-primary-medium border border-default-medium p-4\">\n    <div class=\"datepicker-header\">\n      <div class=\"datepicker-title bg-neutral-primary-medium text-heading px-2 py-3 text-center font-medium\"></div>\n      <div class=\"datepicker-controls flex justify-between mb-2\">\n        <button type=\"button\" class=\"bg-neutral-primary-medium rounded-base text-body hover:bg-neutral-tertiary-medium hover:text-heading text-lg p-2.5 focus:outline-none focus:ring-2 focus:ring-neutral-tertiary prev-btn\"></button>\n        <button type=\"button\" class=\"text-sm rounded-base text-heading bg-neutral-primary-medium font-medium py-2.5 px-5 hover:bg-neutral-tertiary-medium focus:outline-none focus:ring-2 focus:ring-neutral-tertiary view-switch\"></button>\n        <button type=\"button\" class=\"bg-neutral-primary-medium rounded-base text-body hover:bg-neutral-tertiary-medium hover:text-heading text-lg p-2.5 focus:outline-none focus:ring-2 focus:ring-neutral-tertiary next-btn\"></button>\n      </div>\n    </div>\n    <div class=\"datepicker-main p-1\"></div>\n    <div class=\"datepicker-footer\">\n      <div class=\"datepicker-controls flex space-x-2 rtl:space-x-reverse mt-2\">\n        <button type=\"button\" class=\"%buttonClass% today-btn text-white bg-brand hover:bg-brand-strong focus:ring-4 focus:ring-brand-medium font-medium rounded-base text-sm px-5 py-2 text-center w-1/2\"></button>\n        <button type=\"button\" class=\"%buttonClass% clear-btn text-body bg-neutral-secondary-medium border border-default-medium hover:bg-neutral-tertiary-medium focus:ring-4 focus:ring-neutral-tertiary font-medium rounded-base text-sm px-5 py-2 text-center w-1/2\"></button>\n      </div>\n    </div>\n  </div>\n</div>");
var $a18c32cb29f9cbaf$var$daysTemplate = $a18c32cb29f9cbaf$var$optimizeTemplateHTML("<div class=\"days\">\n  <div class=\"days-of-week grid grid-cols-7 mb-1\">".concat($a18c32cb29f9cbaf$var$createTagRepeat('span', 7, {
    "class": 'dow block flex-1 leading-9 border-0 rounded-base cursor-default text-center text-body font-medium text-sm'
}), "</div>\n  <div class=\"datepicker-grid w-64 grid grid-cols-7\">").concat($a18c32cb29f9cbaf$var$createTagRepeat('span', 42, {
    "class": 'block flex-1 leading-9 border-0 rounded-base cursor-default text-center text-body font-medium text-sm h-6 leading-6 text-sm font-medium text-fg-disabled'
}), "</div>\n</div>"));
var $a18c32cb29f9cbaf$var$calendarWeeksTemplate = $a18c32cb29f9cbaf$var$optimizeTemplateHTML("<div class=\"calendar-weeks\">\n  <div class=\"days-of-week flex\"><span class=\"dow h-6 leading-6 text-sm font-medium text-fg-disabled\"></span></div>\n  <div class=\"weeks\">".concat($a18c32cb29f9cbaf$var$createTagRepeat('span', 6, {
    "class": 'week block flex-1 leading-9 border-0 rounded-base cursor-default text-center text-body font-medium text-sm'
}), "</div>\n</div>"));
// Base class of the view classes
var $a18c32cb29f9cbaf$var$View = /*#__PURE__*/ function() {
    function View(picker, config) {
        $a18c32cb29f9cbaf$var$_classCallCheck(this, View);
        Object.assign(this, config, {
            picker: picker,
            element: $a18c32cb29f9cbaf$var$parseHTML("<div class=\"datepicker-view flex\"></div>").firstChild,
            selected: []
        });
        this.init(this.picker.datepicker.config);
    }
    return $a18c32cb29f9cbaf$var$_createClass(View, [
        {
            key: "init",
            value: function init(options) {
                if (options.pickLevel !== undefined) this.isMinView = this.id === options.pickLevel;
                this.setOptions(options);
                this.updateFocus();
                this.updateSelection();
            }
        },
        {
            key: "performBeforeHook",
            value: function performBeforeHook(el, current, timeValue) {
                var result = this.beforeShow(new Date(timeValue));
                switch($a18c32cb29f9cbaf$var$_typeof(result)){
                    case 'boolean':
                        result = {
                            enabled: result
                        };
                        break;
                    case 'string':
                        result = {
                            classes: result
                        };
                }
                if (result) {
                    if (result.enabled === false) {
                        el.classList.add('disabled');
                        $a18c32cb29f9cbaf$var$pushUnique(this.disabled, current);
                    }
                    if (result.classes) {
                        var _el$classList;
                        var extraClasses = result.classes.split(/\s+/);
                        (_el$classList = el.classList).add.apply(_el$classList, $a18c32cb29f9cbaf$var$_toConsumableArray(extraClasses));
                        if (extraClasses.includes('disabled')) $a18c32cb29f9cbaf$var$pushUnique(this.disabled, current);
                    }
                    if (result.content) $a18c32cb29f9cbaf$var$replaceChildNodes(el, result.content);
                }
            }
        }
    ]);
}();
var $a18c32cb29f9cbaf$var$DaysView = /*#__PURE__*/ function(_View) {
    function DaysView(picker) {
        $a18c32cb29f9cbaf$var$_classCallCheck(this, DaysView);
        return $a18c32cb29f9cbaf$var$_callSuper(this, DaysView, [
            picker,
            {
                id: 0,
                name: 'days',
                cellClass: 'day'
            }
        ]);
    }
    $a18c32cb29f9cbaf$var$_inherits(DaysView, _View);
    return $a18c32cb29f9cbaf$var$_createClass(DaysView, [
        {
            key: "init",
            value: function init(options) {
                var onConstruction = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : true;
                if (onConstruction) {
                    var inner = $a18c32cb29f9cbaf$var$parseHTML($a18c32cb29f9cbaf$var$daysTemplate).firstChild;
                    this.dow = inner.firstChild;
                    this.grid = inner.lastChild;
                    this.element.appendChild(inner);
                }
                $a18c32cb29f9cbaf$var$_get($a18c32cb29f9cbaf$var$_getPrototypeOf(DaysView.prototype), "init", this).call(this, options);
            }
        },
        {
            key: "setOptions",
            value: function setOptions(options) {
                var _this = this;
                var updateDOW;
                if ($a18c32cb29f9cbaf$var$hasProperty(options, 'minDate')) this.minDate = options.minDate;
                if ($a18c32cb29f9cbaf$var$hasProperty(options, 'maxDate')) this.maxDate = options.maxDate;
                if (options.datesDisabled) this.datesDisabled = options.datesDisabled;
                if (options.daysOfWeekDisabled) {
                    this.daysOfWeekDisabled = options.daysOfWeekDisabled;
                    updateDOW = true;
                }
                if (options.daysOfWeekHighlighted) this.daysOfWeekHighlighted = options.daysOfWeekHighlighted;
                if (options.todayHighlight !== undefined) this.todayHighlight = options.todayHighlight;
                if (options.weekStart !== undefined) {
                    this.weekStart = options.weekStart;
                    this.weekEnd = options.weekEnd;
                    updateDOW = true;
                }
                if (options.locale) {
                    var locale = this.locale = options.locale;
                    this.dayNames = locale.daysMin;
                    this.switchLabelFormat = locale.titleFormat;
                    updateDOW = true;
                }
                if (options.beforeShowDay !== undefined) this.beforeShow = typeof options.beforeShowDay === 'function' ? options.beforeShowDay : undefined;
                if (options.calendarWeeks !== undefined) {
                    if (options.calendarWeeks && !this.calendarWeeks) {
                        var weeksElem = $a18c32cb29f9cbaf$var$parseHTML($a18c32cb29f9cbaf$var$calendarWeeksTemplate).firstChild;
                        this.calendarWeeks = {
                            element: weeksElem,
                            dow: weeksElem.firstChild,
                            weeks: weeksElem.lastChild
                        };
                        this.element.insertBefore(weeksElem, this.element.firstChild);
                    } else if (this.calendarWeeks && !options.calendarWeeks) {
                        this.element.removeChild(this.calendarWeeks.element);
                        this.calendarWeeks = null;
                    }
                }
                if (options.showDaysOfWeek !== undefined) {
                    if (options.showDaysOfWeek) {
                        $a18c32cb29f9cbaf$var$showElement(this.dow);
                        if (this.calendarWeeks) $a18c32cb29f9cbaf$var$showElement(this.calendarWeeks.dow);
                    } else {
                        $a18c32cb29f9cbaf$var$hideElement(this.dow);
                        if (this.calendarWeeks) $a18c32cb29f9cbaf$var$hideElement(this.calendarWeeks.dow);
                    }
                }
                // update days-of-week when locale, daysOfweekDisabled or weekStart is changed
                if (updateDOW) Array.from(this.dow.children).forEach(function(el, index) {
                    var dow = (_this.weekStart + index) % 7;
                    el.textContent = _this.dayNames[dow];
                    el.className = _this.daysOfWeekDisabled.includes(dow) ? 'dow disabled text-center h-6 leading-6 text-sm font-medium text-fg-disabled cursor-not-allowed' : 'dow text-center h-6 leading-6 text-sm font-medium text-body';
                });
            }
        },
        {
            key: "updateFocus",
            value: function updateFocus() {
                var viewDate = new Date(this.picker.viewDate);
                var viewYear = viewDate.getFullYear();
                var viewMonth = viewDate.getMonth();
                var firstOfMonth = $a18c32cb29f9cbaf$var$dateValue(viewYear, viewMonth, 1);
                var start = $a18c32cb29f9cbaf$var$dayOfTheWeekOf(firstOfMonth, this.weekStart, this.weekStart);
                this.first = firstOfMonth;
                this.last = $a18c32cb29f9cbaf$var$dateValue(viewYear, viewMonth + 1, 0);
                this.start = start;
                this.focused = this.picker.viewDate;
            }
        },
        {
            key: "updateSelection",
            value: function updateSelection() {
                var _this$picker$datepick = this.picker.datepicker, dates = _this$picker$datepick.dates, rangepicker = _this$picker$datepick.rangepicker;
                this.selected = dates;
                if (rangepicker) this.range = rangepicker.dates;
            }
        },
        {
            key: "render",
            value: function render() {
                var _this2 = this;
                // update today marker on ever render
                this.today = this.todayHighlight ? $a18c32cb29f9cbaf$var$today() : undefined;
                // refresh disabled dates on every render in order to clear the ones added
                // by beforeShow hook at previous render
                this.disabled = $a18c32cb29f9cbaf$var$_toConsumableArray(this.datesDisabled);
                var switchLabel = $a18c32cb29f9cbaf$var$formatDate(this.focused, this.switchLabelFormat, this.locale);
                this.picker.setViewSwitchLabel(switchLabel);
                this.picker.setPrevBtnDisabled(this.first <= this.minDate);
                this.picker.setNextBtnDisabled(this.last >= this.maxDate);
                if (this.calendarWeeks) {
                    // start of the UTC week (Monday) of the 1st of the month
                    var startOfWeek = $a18c32cb29f9cbaf$var$dayOfTheWeekOf(this.first, 1, 1);
                    Array.from(this.calendarWeeks.weeks.children).forEach(function(el, index) {
                        el.textContent = $a18c32cb29f9cbaf$var$getWeek($a18c32cb29f9cbaf$var$addWeeks(startOfWeek, index));
                    });
                }
                Array.from(this.grid.children).forEach(function(el, index) {
                    var classList = el.classList;
                    var current = $a18c32cb29f9cbaf$var$addDays(_this2.start, index);
                    var date = new Date(current);
                    var day = date.getDay();
                    el.className = "datepicker-cell hover:bg-neutral-tertiary-medium block flex-1 leading-9 border-0 rounded-base cursor-pointer text-center text-body font-medium text-sm ".concat(_this2.cellClass);
                    el.dataset.date = current;
                    el.textContent = date.getDate();
                    if (current < _this2.first) classList.add('prev', 'text-fg-disabled');
                    else if (current > _this2.last) classList.add('next', 'text-fg-disabled');
                    if (_this2.today === current) classList.add('today', 'bg-gray-100', 'dark:bg-gray-600');
                    if (current < _this2.minDate || current > _this2.maxDate || _this2.disabled.includes(current)) {
                        classList.add('disabled', 'cursor-not-allowed', 'text-fg-disabled');
                        classList.remove('hover:bg-neutral-tertiary-medium', 'text-body', 'cursor-pointer');
                    }
                    if (_this2.daysOfWeekDisabled.includes(day)) {
                        classList.add('disabled', 'cursor-not-allowed', 'text-fg-disabled');
                        classList.remove('hover:bg-neutral-tertiary-medium', 'text-body', 'cursor-pointer');
                        $a18c32cb29f9cbaf$var$pushUnique(_this2.disabled, current);
                    }
                    if (_this2.daysOfWeekHighlighted.includes(day)) classList.add('highlighted');
                    if (_this2.range) {
                        var _this2$range = $a18c32cb29f9cbaf$var$_slicedToArray(_this2.range, 2), rangeStart = _this2$range[0], rangeEnd = _this2$range[1];
                        if (current > rangeStart && current < rangeEnd) {
                            classList.add('range', 'bg-neutral-tertiary-medium');
                            classList.remove('rounded-base', 'rounded-s-base', 'rounded-e-base');
                        }
                        if (current === rangeStart) {
                            classList.add('range-start', 'bg-brand', 'rounded-s-base');
                            classList.remove('rounded-base', 'rounded-e-base');
                        }
                        if (current === rangeEnd) {
                            classList.add('range-end', 'bg-neutral-tertiary-medium', 'rounded-e-base');
                            classList.remove('rounded-base', 'rounded-s-base');
                        }
                    }
                    if (_this2.selected.includes(current)) {
                        classList.add('selected', 'bg-brand', 'text-white');
                        classList.remove('text-body', 'hover:bg-neutral-tertiary-medium', 'bg-neutral-tertiary-medium');
                    }
                    if (current === _this2.focused) classList.add('focused');
                    if (_this2.beforeShow) _this2.performBeforeHook(el, current, current);
                });
            }
        },
        {
            key: "refresh",
            value: function refresh() {
                var _this3 = this;
                var _ref = this.range || [], _ref2 = $a18c32cb29f9cbaf$var$_slicedToArray(_ref, 2), rangeStart = _ref2[0], rangeEnd = _ref2[1];
                this.grid.querySelectorAll('.range, .range-start, .range-end, .selected, .focused').forEach(function(el) {
                    el.classList.remove('range', 'range-start', 'range-end', 'selected', 'bg-brand', 'text-white', 'focused');
                    el.classList.add('text-body', 'rounded-base');
                });
                Array.from(this.grid.children).forEach(function(el) {
                    var current = Number(el.dataset.date);
                    var classList = el.classList;
                    classList.remove('bg-neutral-tertiary-medium', 'rounded-s-base', 'rounded-e-base');
                    if (current > rangeStart && current < rangeEnd) {
                        classList.add('range', 'bg-neutral-tertiary-medium');
                        classList.remove('rounded-base');
                    }
                    if (current === rangeStart) {
                        classList.add('range-start', 'bg-brand', 'text-white', 'rounded-s-base');
                        classList.remove('rounded-base');
                    }
                    if (current === rangeEnd) {
                        classList.add('range-end', 'bg-neutral-tertiary-medium', 'rounded-e-base');
                        classList.remove('rounded-base');
                    }
                    if (_this3.selected.includes(current)) {
                        classList.add('selected', 'bg-brand', 'text-white');
                        classList.remove('text-body', 'hover:bg-neutral-tertiary-medium', 'bg-neutral-tertiary-medium');
                    }
                    if (current === _this3.focused) classList.add('focused');
                });
            }
        },
        {
            key: "refreshFocus",
            value: function refreshFocus() {
                var index = Math.round((this.focused - this.start) / 86400000);
                this.grid.querySelectorAll('.focused').forEach(function(el) {
                    el.classList.remove('focused');
                });
                this.grid.children[index].classList.add('focused');
            }
        }
    ]);
}($a18c32cb29f9cbaf$var$View);
function $a18c32cb29f9cbaf$var$computeMonthRange(range, thisYear) {
    if (!range || !range[0] || !range[1]) return;
    var _range = $a18c32cb29f9cbaf$var$_slicedToArray(range, 2), _range$ = $a18c32cb29f9cbaf$var$_slicedToArray(_range[0], 2), startY = _range$[0], startM = _range$[1], _range$2 = $a18c32cb29f9cbaf$var$_slicedToArray(_range[1], 2), endY = _range$2[0], endM = _range$2[1];
    if (startY > thisYear || endY < thisYear) return;
    return [
        startY === thisYear ? startM : -1,
        endY === thisYear ? endM : 12
    ];
}
var $a18c32cb29f9cbaf$var$MonthsView = /*#__PURE__*/ function(_View) {
    function MonthsView(picker) {
        $a18c32cb29f9cbaf$var$_classCallCheck(this, MonthsView);
        return $a18c32cb29f9cbaf$var$_callSuper(this, MonthsView, [
            picker,
            {
                id: 1,
                name: 'months',
                cellClass: 'month'
            }
        ]);
    }
    $a18c32cb29f9cbaf$var$_inherits(MonthsView, _View);
    return $a18c32cb29f9cbaf$var$_createClass(MonthsView, [
        {
            key: "init",
            value: function init(options) {
                var onConstruction = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : true;
                if (onConstruction) {
                    this.grid = this.element;
                    this.element.classList.add('months', 'datepicker-grid', 'w-64', 'grid', 'grid-cols-4');
                    this.grid.appendChild($a18c32cb29f9cbaf$var$parseHTML($a18c32cb29f9cbaf$var$createTagRepeat('span', 12, {
                        'data-month': function dataMonth(ix) {
                            return ix;
                        }
                    })));
                }
                $a18c32cb29f9cbaf$var$_get($a18c32cb29f9cbaf$var$_getPrototypeOf(MonthsView.prototype), "init", this).call(this, options);
            }
        },
        {
            key: "setOptions",
            value: function setOptions(options) {
                if (options.locale) this.monthNames = options.locale.monthsShort;
                if ($a18c32cb29f9cbaf$var$hasProperty(options, 'minDate')) {
                    if (options.minDate === undefined) this.minYear = this.minMonth = this.minDate = undefined;
                    else {
                        var minDateObj = new Date(options.minDate);
                        this.minYear = minDateObj.getFullYear();
                        this.minMonth = minDateObj.getMonth();
                        this.minDate = minDateObj.setDate(1);
                    }
                }
                if ($a18c32cb29f9cbaf$var$hasProperty(options, 'maxDate')) {
                    if (options.maxDate === undefined) this.maxYear = this.maxMonth = this.maxDate = undefined;
                    else {
                        var maxDateObj = new Date(options.maxDate);
                        this.maxYear = maxDateObj.getFullYear();
                        this.maxMonth = maxDateObj.getMonth();
                        this.maxDate = $a18c32cb29f9cbaf$var$dateValue(this.maxYear, this.maxMonth + 1, 0);
                    }
                }
                if (options.beforeShowMonth !== undefined) this.beforeShow = typeof options.beforeShowMonth === 'function' ? options.beforeShowMonth : undefined;
            }
        },
        {
            key: "updateFocus",
            value: function updateFocus() {
                var viewDate = new Date(this.picker.viewDate);
                this.year = viewDate.getFullYear();
                this.focused = viewDate.getMonth();
            }
        },
        {
            key: "updateSelection",
            value: function updateSelection() {
                var _this$picker$datepick = this.picker.datepicker, dates = _this$picker$datepick.dates, rangepicker = _this$picker$datepick.rangepicker;
                this.selected = dates.reduce(function(selected, timeValue) {
                    var date = new Date(timeValue);
                    var year = date.getFullYear();
                    var month = date.getMonth();
                    if (selected[year] === undefined) selected[year] = [
                        month
                    ];
                    else $a18c32cb29f9cbaf$var$pushUnique(selected[year], month);
                    return selected;
                }, {});
                if (rangepicker && rangepicker.dates) this.range = rangepicker.dates.map(function(timeValue) {
                    var date = new Date(timeValue);
                    return isNaN(date) ? undefined : [
                        date.getFullYear(),
                        date.getMonth()
                    ];
                });
            }
        },
        {
            key: "render",
            value: function render() {
                var _this = this;
                // refresh disabled months on every render in order to clear the ones added
                // by beforeShow hook at previous render
                this.disabled = [];
                this.picker.setViewSwitchLabel(this.year);
                this.picker.setPrevBtnDisabled(this.year <= this.minYear);
                this.picker.setNextBtnDisabled(this.year >= this.maxYear);
                var selected = this.selected[this.year] || [];
                var yrOutOfRange = this.year < this.minYear || this.year > this.maxYear;
                var isMinYear = this.year === this.minYear;
                var isMaxYear = this.year === this.maxYear;
                var range = $a18c32cb29f9cbaf$var$computeMonthRange(this.range, this.year);
                Array.from(this.grid.children).forEach(function(el, index) {
                    var classList = el.classList;
                    var date = $a18c32cb29f9cbaf$var$dateValue(_this.year, index, 1);
                    el.className = "datepicker-cell hover:bg-neutral-tertiary-medium block flex-1 leading-9 border-0 rounded-base cursor-pointer text-center text-body font-medium text-sm ".concat(_this.cellClass);
                    if (_this.isMinView) el.dataset.date = date;
                    // reset text on every render to clear the custom content set
                    // by beforeShow hook at previous render
                    el.textContent = _this.monthNames[index];
                    if (yrOutOfRange || isMinYear && index < _this.minMonth || isMaxYear && index > _this.maxMonth) classList.add('disabled');
                    if (range) {
                        var _range2 = $a18c32cb29f9cbaf$var$_slicedToArray(range, 2), rangeStart = _range2[0], rangeEnd = _range2[1];
                        if (index > rangeStart && index < rangeEnd) classList.add('range');
                        if (index === rangeStart) classList.add('range-start');
                        if (index === rangeEnd) classList.add('range-end');
                    }
                    if (selected.includes(index)) {
                        classList.add('selected', 'bg-brand', 'text-white', 'dark:text-white');
                        classList.remove('text-body', 'hover:bg-neutral-tertiary-medium', 'dark:text-white');
                    }
                    if (index === _this.focused) classList.add('focused');
                    if (_this.beforeShow) _this.performBeforeHook(el, index, date);
                });
            }
        },
        {
            key: "refresh",
            value: function refresh() {
                var _this2 = this;
                var selected = this.selected[this.year] || [];
                var _ref = $a18c32cb29f9cbaf$var$computeMonthRange(this.range, this.year) || [], _ref2 = $a18c32cb29f9cbaf$var$_slicedToArray(_ref, 2), rangeStart = _ref2[0], rangeEnd = _ref2[1];
                this.grid.querySelectorAll('.range, .range-start, .range-end, .selected, .focused').forEach(function(el) {
                    el.classList.remove('range', 'range-start', 'range-end', 'selected', 'bg-brand', 'dark:text-white', 'text-white', 'focused');
                    el.classList.add('text-body', 'hover:bg-neutral-tertiary-medium', 'dark:text-white');
                });
                Array.from(this.grid.children).forEach(function(el, index) {
                    var classList = el.classList;
                    if (index > rangeStart && index < rangeEnd) classList.add('range');
                    if (index === rangeStart) classList.add('range-start');
                    if (index === rangeEnd) classList.add('range-end');
                    if (selected.includes(index)) {
                        classList.add('selected', 'bg-brand', 'text-white', 'dark:text-white');
                        classList.remove('text-body', 'hover:bg-neutral-tertiary-medium', 'dark:text-white');
                    }
                    if (index === _this2.focused) classList.add('focused');
                });
            }
        },
        {
            key: "refreshFocus",
            value: function refreshFocus() {
                this.grid.querySelectorAll('.focused').forEach(function(el) {
                    el.classList.remove('focused');
                });
                this.grid.children[this.focused].classList.add('focused');
            }
        }
    ]);
}($a18c32cb29f9cbaf$var$View);
function $a18c32cb29f9cbaf$var$toTitleCase(word) {
    return $a18c32cb29f9cbaf$var$_toConsumableArray(word).reduce(function(str, ch, ix) {
        return str += ix ? ch : ch.toUpperCase();
    }, '');
}
// Class representing the years and decades view elements
var $a18c32cb29f9cbaf$var$YearsView = /*#__PURE__*/ function(_View) {
    function YearsView(picker, config) {
        $a18c32cb29f9cbaf$var$_classCallCheck(this, YearsView);
        return $a18c32cb29f9cbaf$var$_callSuper(this, YearsView, [
            picker,
            config
        ]);
    }
    $a18c32cb29f9cbaf$var$_inherits(YearsView, _View);
    return $a18c32cb29f9cbaf$var$_createClass(YearsView, [
        {
            key: "init",
            value: function init(options) {
                var onConstruction = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : true;
                if (onConstruction) {
                    this.navStep = this.step * 10;
                    this.beforeShowOption = "beforeShow".concat($a18c32cb29f9cbaf$var$toTitleCase(this.cellClass));
                    this.grid = this.element;
                    this.element.classList.add(this.name, 'datepicker-grid', 'w-64', 'grid', 'grid-cols-4');
                    this.grid.appendChild($a18c32cb29f9cbaf$var$parseHTML($a18c32cb29f9cbaf$var$createTagRepeat('span', 12)));
                }
                $a18c32cb29f9cbaf$var$_get($a18c32cb29f9cbaf$var$_getPrototypeOf(YearsView.prototype), "init", this).call(this, options);
            }
        },
        {
            key: "setOptions",
            value: function setOptions(options) {
                if ($a18c32cb29f9cbaf$var$hasProperty(options, 'minDate')) {
                    if (options.minDate === undefined) this.minYear = this.minDate = undefined;
                    else {
                        this.minYear = $a18c32cb29f9cbaf$var$startOfYearPeriod(options.minDate, this.step);
                        this.minDate = $a18c32cb29f9cbaf$var$dateValue(this.minYear, 0, 1);
                    }
                }
                if ($a18c32cb29f9cbaf$var$hasProperty(options, 'maxDate')) {
                    if (options.maxDate === undefined) this.maxYear = this.maxDate = undefined;
                    else {
                        this.maxYear = $a18c32cb29f9cbaf$var$startOfYearPeriod(options.maxDate, this.step);
                        this.maxDate = $a18c32cb29f9cbaf$var$dateValue(this.maxYear, 11, 31);
                    }
                }
                if (options[this.beforeShowOption] !== undefined) {
                    var beforeShow = options[this.beforeShowOption];
                    this.beforeShow = typeof beforeShow === 'function' ? beforeShow : undefined;
                }
            }
        },
        {
            key: "updateFocus",
            value: function updateFocus() {
                var viewDate = new Date(this.picker.viewDate);
                var first = $a18c32cb29f9cbaf$var$startOfYearPeriod(viewDate, this.navStep);
                var last = first + 9 * this.step;
                this.first = first;
                this.last = last;
                this.start = first - this.step;
                this.focused = $a18c32cb29f9cbaf$var$startOfYearPeriod(viewDate, this.step);
            }
        },
        {
            key: "updateSelection",
            value: function updateSelection() {
                var _this = this;
                var _this$picker$datepick = this.picker.datepicker, dates = _this$picker$datepick.dates, rangepicker = _this$picker$datepick.rangepicker;
                this.selected = dates.reduce(function(years, timeValue) {
                    return $a18c32cb29f9cbaf$var$pushUnique(years, $a18c32cb29f9cbaf$var$startOfYearPeriod(timeValue, _this.step));
                }, []);
                if (rangepicker && rangepicker.dates) this.range = rangepicker.dates.map(function(timeValue) {
                    if (timeValue !== undefined) return $a18c32cb29f9cbaf$var$startOfYearPeriod(timeValue, _this.step);
                });
            }
        },
        {
            key: "render",
            value: function render() {
                var _this2 = this;
                // refresh disabled years on every render in order to clear the ones added
                // by beforeShow hook at previous render
                this.disabled = [];
                this.picker.setViewSwitchLabel("".concat(this.first, "-").concat(this.last));
                this.picker.setPrevBtnDisabled(this.first <= this.minYear);
                this.picker.setNextBtnDisabled(this.last >= this.maxYear);
                Array.from(this.grid.children).forEach(function(el, index) {
                    var classList = el.classList;
                    var current = _this2.start + index * _this2.step;
                    var date = $a18c32cb29f9cbaf$var$dateValue(current, 0, 1);
                    el.className = "datepicker-cell hover:bg-neutral-tertiary-medium block flex-1 leading-9 border-0 rounded-base cursor-pointer text-center text-body font-medium text-sm ".concat(_this2.cellClass);
                    if (_this2.isMinView) el.dataset.date = date;
                    el.textContent = el.dataset.year = current;
                    if (index === 0) classList.add('prev');
                    else if (index === 11) classList.add('next');
                    if (current < _this2.minYear || current > _this2.maxYear) classList.add('disabled');
                    if (_this2.range) {
                        var _this2$range = $a18c32cb29f9cbaf$var$_slicedToArray(_this2.range, 2), rangeStart = _this2$range[0], rangeEnd = _this2$range[1];
                        if (current > rangeStart && current < rangeEnd) classList.add('range');
                        if (current === rangeStart) classList.add('range-start');
                        if (current === rangeEnd) classList.add('range-end');
                    }
                    if (_this2.selected.includes(current)) {
                        classList.add('selected', 'bg-brand', 'text-white', 'dark:text-white');
                        classList.remove('text-body', 'hover:bg-neutral-tertiary-medium', 'dark:text-white');
                    }
                    if (current === _this2.focused) classList.add('focused');
                    if (_this2.beforeShow) _this2.performBeforeHook(el, current, date);
                });
            }
        },
        {
            key: "refresh",
            value: function refresh() {
                var _this3 = this;
                var _ref = this.range || [], _ref2 = $a18c32cb29f9cbaf$var$_slicedToArray(_ref, 2), rangeStart = _ref2[0], rangeEnd = _ref2[1];
                this.grid.querySelectorAll('.range, .range-start, .range-end, .selected, .focused').forEach(function(el) {
                    el.classList.remove('range', 'range-start', 'range-end', 'selected', 'bg-brand', 'text-white', 'dark:text-white', 'focused');
                });
                Array.from(this.grid.children).forEach(function(el) {
                    var current = Number(el.textContent);
                    var classList = el.classList;
                    if (current > rangeStart && current < rangeEnd) classList.add('range');
                    if (current === rangeStart) classList.add('range-start');
                    if (current === rangeEnd) classList.add('range-end');
                    if (_this3.selected.includes(current)) {
                        classList.add('selected', 'bg-brand', 'text-white', 'hover:text-heading');
                        classList.remove('text-body', 'hover:bg-neutral-tertiary-medium', 'hover:text-heading');
                    }
                    if (current === _this3.focused) classList.add('focused');
                });
            }
        },
        {
            key: "refreshFocus",
            value: function refreshFocus() {
                var index = Math.round((this.focused - this.start) / this.step);
                this.grid.querySelectorAll('.focused').forEach(function(el) {
                    el.classList.remove('focused');
                });
                this.grid.children[index].classList.add('focused');
            }
        }
    ]);
}($a18c32cb29f9cbaf$var$View);
function $a18c32cb29f9cbaf$var$triggerDatepickerEvent(datepicker, type) {
    var detail = {
        date: datepicker.getDate(),
        viewDate: new Date(datepicker.picker.viewDate),
        viewId: datepicker.picker.currentView.id,
        datepicker: datepicker
    };
    datepicker.element.dispatchEvent(new CustomEvent(type, {
        detail: detail
    }));
}
// direction: -1 (to previous), 1 (to next)
function $a18c32cb29f9cbaf$var$goToPrevOrNext(datepicker, direction) {
    var _datepicker$config = datepicker.config, minDate = _datepicker$config.minDate, maxDate = _datepicker$config.maxDate;
    var _datepicker$picker = datepicker.picker, currentView = _datepicker$picker.currentView, viewDate = _datepicker$picker.viewDate;
    var newViewDate;
    switch(currentView.id){
        case 0:
            newViewDate = $a18c32cb29f9cbaf$var$addMonths(viewDate, direction);
            break;
        case 1:
            newViewDate = $a18c32cb29f9cbaf$var$addYears(viewDate, direction);
            break;
        default:
            newViewDate = $a18c32cb29f9cbaf$var$addYears(viewDate, direction * currentView.navStep);
    }
    newViewDate = $a18c32cb29f9cbaf$var$limitToRange(newViewDate, minDate, maxDate);
    datepicker.picker.changeFocus(newViewDate).render();
}
function $a18c32cb29f9cbaf$var$switchView(datepicker) {
    var viewId = datepicker.picker.currentView.id;
    if (viewId === datepicker.config.maxView) return;
    datepicker.picker.changeView(viewId + 1).render();
}
function $a18c32cb29f9cbaf$var$unfocus(datepicker) {
    if (datepicker.config.updateOnBlur) datepicker.update({
        autohide: true
    });
    else {
        datepicker.refresh('input');
        datepicker.hide();
    }
}
function $a18c32cb29f9cbaf$var$goToSelectedMonthOrYear(datepicker, selection) {
    var picker = datepicker.picker;
    var viewDate = new Date(picker.viewDate);
    var viewId = picker.currentView.id;
    var newDate = viewId === 1 ? $a18c32cb29f9cbaf$var$addMonths(viewDate, selection - viewDate.getMonth()) : $a18c32cb29f9cbaf$var$addYears(viewDate, selection - viewDate.getFullYear());
    picker.changeFocus(newDate).changeView(viewId - 1).render();
}
function $a18c32cb29f9cbaf$var$onClickTodayBtn(datepicker) {
    var picker = datepicker.picker;
    var currentDate = $a18c32cb29f9cbaf$var$today();
    if (datepicker.config.todayBtnMode === 1) {
        if (datepicker.config.autohide) {
            datepicker.setDate(currentDate);
            return;
        }
        datepicker.setDate(currentDate, {
            render: false
        });
        picker.update();
    }
    if (picker.viewDate !== currentDate) picker.changeFocus(currentDate);
    picker.changeView(0).render();
}
function $a18c32cb29f9cbaf$var$onClickClearBtn(datepicker) {
    datepicker.setDate({
        clear: true
    });
}
function $a18c32cb29f9cbaf$var$onClickViewSwitch(datepicker) {
    $a18c32cb29f9cbaf$var$switchView(datepicker);
}
function $a18c32cb29f9cbaf$var$onClickPrevBtn(datepicker) {
    $a18c32cb29f9cbaf$var$goToPrevOrNext(datepicker, -1);
}
function $a18c32cb29f9cbaf$var$onClickNextBtn(datepicker) {
    $a18c32cb29f9cbaf$var$goToPrevOrNext(datepicker, 1);
}
// For the picker's main block to delegete the events from `datepicker-cell`s
function $a18c32cb29f9cbaf$var$onClickView(datepicker, ev) {
    var target = $a18c32cb29f9cbaf$var$findElementInEventPath(ev, '.datepicker-cell');
    if (!target || target.classList.contains('disabled')) return;
    var _datepicker$picker$cu = datepicker.picker.currentView, id = _datepicker$picker$cu.id, isMinView = _datepicker$picker$cu.isMinView;
    if (isMinView) datepicker.setDate(Number(target.dataset.date));
    else if (id === 1) $a18c32cb29f9cbaf$var$goToSelectedMonthOrYear(datepicker, Number(target.dataset.month));
    else $a18c32cb29f9cbaf$var$goToSelectedMonthOrYear(datepicker, Number(target.dataset.year));
}
function $a18c32cb29f9cbaf$var$onClickPicker(datepicker) {
    if (!datepicker.inline && !datepicker.config.disableTouchKeyboard) datepicker.inputField.focus();
}
function $a18c32cb29f9cbaf$var$processPickerOptions(picker, options) {
    if (options.title !== undefined) {
        if (options.title) {
            picker.controls.title.textContent = options.title;
            $a18c32cb29f9cbaf$var$showElement(picker.controls.title);
        } else {
            picker.controls.title.textContent = '';
            $a18c32cb29f9cbaf$var$hideElement(picker.controls.title);
        }
    }
    if (options.prevArrow) {
        var prevBtn = picker.controls.prevBtn;
        $a18c32cb29f9cbaf$var$emptyChildNodes(prevBtn);
        options.prevArrow.forEach(function(node) {
            prevBtn.appendChild(node.cloneNode(true));
        });
    }
    if (options.nextArrow) {
        var nextBtn = picker.controls.nextBtn;
        $a18c32cb29f9cbaf$var$emptyChildNodes(nextBtn);
        options.nextArrow.forEach(function(node) {
            nextBtn.appendChild(node.cloneNode(true));
        });
    }
    if (options.locale) {
        picker.controls.todayBtn.textContent = options.locale.today;
        picker.controls.clearBtn.textContent = options.locale.clear;
    }
    if (options.todayBtn !== undefined) {
        if (options.todayBtn) $a18c32cb29f9cbaf$var$showElement(picker.controls.todayBtn);
        else $a18c32cb29f9cbaf$var$hideElement(picker.controls.todayBtn);
    }
    if ($a18c32cb29f9cbaf$var$hasProperty(options, 'minDate') || $a18c32cb29f9cbaf$var$hasProperty(options, 'maxDate')) {
        var _picker$datepicker$co = picker.datepicker.config, minDate = _picker$datepicker$co.minDate, maxDate = _picker$datepicker$co.maxDate;
        picker.controls.todayBtn.disabled = !$a18c32cb29f9cbaf$var$isInRange($a18c32cb29f9cbaf$var$today(), minDate, maxDate);
    }
    if (options.clearBtn !== undefined) {
        if (options.clearBtn) $a18c32cb29f9cbaf$var$showElement(picker.controls.clearBtn);
        else $a18c32cb29f9cbaf$var$hideElement(picker.controls.clearBtn);
    }
}
// Compute view date to reset, which will be...
// - the last item of the selected dates or defaultViewDate if no selection
// - limitted to minDate or maxDate if it exceeds the range
function $a18c32cb29f9cbaf$var$computeResetViewDate(datepicker) {
    var dates = datepicker.dates, config = datepicker.config;
    var viewDate = dates.length > 0 ? $a18c32cb29f9cbaf$var$lastItemOf(dates) : config.defaultViewDate;
    return $a18c32cb29f9cbaf$var$limitToRange(viewDate, config.minDate, config.maxDate);
}
// Change current view's view date
function $a18c32cb29f9cbaf$var$setViewDate(picker, newDate) {
    var oldViewDate = new Date(picker.viewDate);
    var newViewDate = new Date(newDate);
    var _picker$currentView = picker.currentView, id = _picker$currentView.id, year = _picker$currentView.year, first = _picker$currentView.first, last = _picker$currentView.last;
    var viewYear = newViewDate.getFullYear();
    picker.viewDate = newDate;
    if (viewYear !== oldViewDate.getFullYear()) $a18c32cb29f9cbaf$var$triggerDatepickerEvent(picker.datepicker, 'changeYear');
    if (newViewDate.getMonth() !== oldViewDate.getMonth()) $a18c32cb29f9cbaf$var$triggerDatepickerEvent(picker.datepicker, 'changeMonth');
    // return whether the new date is in different period on time from the one
    // displayed in the current view
    // when true, the view needs to be re-rendered on the next UI refresh.
    switch(id){
        case 0:
            return newDate < first || newDate > last;
        case 1:
            return viewYear !== year;
        default:
            return viewYear < first || viewYear > last;
    }
}
function $a18c32cb29f9cbaf$var$getTextDirection(el) {
    return window.getComputedStyle(el).direction;
}
// Class representing the picker UI
var $a18c32cb29f9cbaf$var$Picker = /*#__PURE__*/ function() {
    function Picker(datepicker) {
        $a18c32cb29f9cbaf$var$_classCallCheck(this, Picker);
        this.datepicker = datepicker;
        var template = $a18c32cb29f9cbaf$var$pickerTemplate.replace(/%buttonClass%/g, datepicker.config.buttonClass);
        var element = this.element = $a18c32cb29f9cbaf$var$parseHTML(template).firstChild;
        var _element$firstChild$c = $a18c32cb29f9cbaf$var$_slicedToArray(element.firstChild.children, 3), header = _element$firstChild$c[0], main = _element$firstChild$c[1], footer = _element$firstChild$c[2];
        var title = header.firstElementChild;
        var _header$lastElementCh = $a18c32cb29f9cbaf$var$_slicedToArray(header.lastElementChild.children, 3), prevBtn = _header$lastElementCh[0], viewSwitch = _header$lastElementCh[1], nextBtn = _header$lastElementCh[2];
        var _footer$firstChild$ch = $a18c32cb29f9cbaf$var$_slicedToArray(footer.firstChild.children, 2), todayBtn = _footer$firstChild$ch[0], clearBtn = _footer$firstChild$ch[1];
        var controls = {
            title: title,
            prevBtn: prevBtn,
            viewSwitch: viewSwitch,
            nextBtn: nextBtn,
            todayBtn: todayBtn,
            clearBtn: clearBtn
        };
        this.main = main;
        this.controls = controls;
        var elementClass = datepicker.inline ? 'inline' : 'dropdown';
        element.classList.add("datepicker-".concat(elementClass));
        elementClass === 'dropdown' && element.classList.add('dropdown', 'absolute', 'top-0', 'left-0', 'z-50', 'pt-2');
        $a18c32cb29f9cbaf$var$processPickerOptions(this, datepicker.config);
        this.viewDate = $a18c32cb29f9cbaf$var$computeResetViewDate(datepicker);
        // set up event listeners
        $a18c32cb29f9cbaf$var$registerListeners(datepicker, [
            [
                element,
                'click',
                $a18c32cb29f9cbaf$var$onClickPicker.bind(null, datepicker),
                {
                    capture: true
                }
            ],
            [
                main,
                'click',
                $a18c32cb29f9cbaf$var$onClickView.bind(null, datepicker)
            ],
            [
                controls.viewSwitch,
                'click',
                $a18c32cb29f9cbaf$var$onClickViewSwitch.bind(null, datepicker)
            ],
            [
                controls.prevBtn,
                'click',
                $a18c32cb29f9cbaf$var$onClickPrevBtn.bind(null, datepicker)
            ],
            [
                controls.nextBtn,
                'click',
                $a18c32cb29f9cbaf$var$onClickNextBtn.bind(null, datepicker)
            ],
            [
                controls.todayBtn,
                'click',
                $a18c32cb29f9cbaf$var$onClickTodayBtn.bind(null, datepicker)
            ],
            [
                controls.clearBtn,
                'click',
                $a18c32cb29f9cbaf$var$onClickClearBtn.bind(null, datepicker)
            ]
        ]);
        // set up views
        this.views = [
            new $a18c32cb29f9cbaf$var$DaysView(this),
            new $a18c32cb29f9cbaf$var$MonthsView(this),
            new $a18c32cb29f9cbaf$var$YearsView(this, {
                id: 2,
                name: 'years',
                cellClass: 'year',
                step: 1
            }),
            new $a18c32cb29f9cbaf$var$YearsView(this, {
                id: 3,
                name: 'decades',
                cellClass: 'decade',
                step: 10
            })
        ];
        this.currentView = this.views[datepicker.config.startView];
        this.currentView.render();
        this.main.appendChild(this.currentView.element);
        datepicker.config.container.appendChild(this.element);
    }
    return $a18c32cb29f9cbaf$var$_createClass(Picker, [
        {
            key: "setOptions",
            value: function setOptions(options) {
                $a18c32cb29f9cbaf$var$processPickerOptions(this, options);
                this.views.forEach(function(view) {
                    view.init(options, false);
                });
                this.currentView.render();
            }
        },
        {
            key: "detach",
            value: function detach() {
                this.datepicker.config.container.removeChild(this.element);
            }
        },
        {
            key: "show",
            value: function show() {
                if (this.active) return;
                this.element.classList.add('active', 'block');
                this.element.classList.remove('hidden');
                this.active = true;
                var datepicker = this.datepicker;
                if (!datepicker.inline) {
                    // ensure picker's direction matches input's
                    var inputDirection = $a18c32cb29f9cbaf$var$getTextDirection(datepicker.inputField);
                    if (inputDirection !== $a18c32cb29f9cbaf$var$getTextDirection(datepicker.config.container)) this.element.dir = inputDirection;
                    else if (this.element.dir) this.element.removeAttribute('dir');
                    this.place();
                    if (datepicker.config.disableTouchKeyboard) datepicker.inputField.blur();
                }
                $a18c32cb29f9cbaf$var$triggerDatepickerEvent(datepicker, 'show');
            }
        },
        {
            key: "hide",
            value: function hide() {
                if (!this.active) return;
                this.datepicker.exitEditMode();
                this.element.classList.remove('active', 'block');
                this.element.classList.add('active', 'block', 'hidden');
                this.active = false;
                $a18c32cb29f9cbaf$var$triggerDatepickerEvent(this.datepicker, 'hide');
            }
        },
        {
            key: "place",
            value: function place() {
                var _this$element = this.element, classList = _this$element.classList, style = _this$element.style;
                var _this$datepicker = this.datepicker, config = _this$datepicker.config, inputField = _this$datepicker.inputField;
                var container = config.container;
                var _this$element$getBoun = this.element.getBoundingClientRect(), calendarWidth = _this$element$getBoun.width, calendarHeight = _this$element$getBoun.height;
                var _container$getBoundin = container.getBoundingClientRect(), containerLeft = _container$getBoundin.left, containerTop = _container$getBoundin.top, containerWidth = _container$getBoundin.width;
                var _inputField$getBoundi = inputField.getBoundingClientRect(), inputLeft = _inputField$getBoundi.left, inputTop = _inputField$getBoundi.top, inputWidth = _inputField$getBoundi.width, inputHeight = _inputField$getBoundi.height;
                var _config$orientation = config.orientation, orientX = _config$orientation.x, orientY = _config$orientation.y;
                var scrollTop;
                var left;
                var top;
                if (container === document.body) {
                    scrollTop = window.scrollY;
                    left = inputLeft + window.scrollX;
                    top = inputTop + scrollTop;
                } else {
                    scrollTop = container.scrollTop;
                    left = inputLeft - containerLeft;
                    top = inputTop - containerTop + scrollTop;
                }
                if (orientX === 'auto') {
                    if (left < 0) {
                        // align to the left and move into visible area if input's left edge < window's
                        orientX = 'left';
                        left = 10;
                    } else if (left + calendarWidth > containerWidth) // align to the right if canlendar's right edge > container's
                    orientX = 'right';
                    else orientX = $a18c32cb29f9cbaf$var$getTextDirection(inputField) === 'rtl' ? 'right' : 'left';
                }
                if (orientX === 'right') left -= calendarWidth - inputWidth;
                if (orientY === 'auto') orientY = top - calendarHeight < scrollTop ? 'bottom' : 'top';
                if (orientY === 'top') top -= calendarHeight;
                else top += inputHeight;
                classList.remove('datepicker-orient-top', 'datepicker-orient-bottom', 'datepicker-orient-right', 'datepicker-orient-left');
                classList.add("datepicker-orient-".concat(orientY), "datepicker-orient-".concat(orientX));
                style.top = top ? "".concat(top, "px") : top;
                style.left = left ? "".concat(left, "px") : left;
            }
        },
        {
            key: "setViewSwitchLabel",
            value: function setViewSwitchLabel(labelText) {
                this.controls.viewSwitch.textContent = labelText;
            }
        },
        {
            key: "setPrevBtnDisabled",
            value: function setPrevBtnDisabled(disabled) {
                this.controls.prevBtn.disabled = disabled;
            }
        },
        {
            key: "setNextBtnDisabled",
            value: function setNextBtnDisabled(disabled) {
                this.controls.nextBtn.disabled = disabled;
            }
        },
        {
            key: "changeView",
            value: function changeView(viewId) {
                var oldView = this.currentView;
                var newView = this.views[viewId];
                if (newView.id !== oldView.id) {
                    this.currentView = newView;
                    this._renderMethod = 'render';
                    $a18c32cb29f9cbaf$var$triggerDatepickerEvent(this.datepicker, 'changeView');
                    this.main.replaceChild(newView.element, oldView.element);
                }
                return this;
            }
        },
        {
            key: "changeFocus",
            value: function changeFocus(newViewDate) {
                this._renderMethod = $a18c32cb29f9cbaf$var$setViewDate(this, newViewDate) ? 'render' : 'refreshFocus';
                this.views.forEach(function(view) {
                    view.updateFocus();
                });
                return this;
            }
        },
        {
            key: "update",
            value: function update() {
                var newViewDate = $a18c32cb29f9cbaf$var$computeResetViewDate(this.datepicker);
                this._renderMethod = $a18c32cb29f9cbaf$var$setViewDate(this, newViewDate) ? 'render' : 'refresh';
                this.views.forEach(function(view) {
                    view.updateFocus();
                    view.updateSelection();
                });
                return this;
            }
        },
        {
            key: "render",
            value: function render() {
                var quickRender = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : true;
                var renderMethod = quickRender && this._renderMethod || 'render';
                delete this._renderMethod;
                this.currentView[renderMethod]();
            }
        }
    ]);
}();
// Find the closest date that doesn't meet the condition for unavailable date
// Returns undefined if no available date is found
// addFn: function to calculate the next date
//   - args: time value, amount
// increase: amount to pass to addFn
// testFn: function to test the unavailablity of the date
//   - args: time value; retun: true if unavailable
function $a18c32cb29f9cbaf$var$findNextAvailableOne(date, addFn, increase, testFn, min, max) {
    if (!$a18c32cb29f9cbaf$var$isInRange(date, min, max)) return;
    if (testFn(date)) {
        var newDate = addFn(date, increase);
        return $a18c32cb29f9cbaf$var$findNextAvailableOne(newDate, addFn, increase, testFn, min, max);
    }
    return date;
}
// direction: -1 (left/up), 1 (right/down)
// vertical: true for up/down, false for left/right
function $a18c32cb29f9cbaf$var$moveByArrowKey(datepicker, ev, direction, vertical) {
    var picker = datepicker.picker;
    var currentView = picker.currentView;
    var step = currentView.step || 1;
    var viewDate = picker.viewDate;
    var addFn;
    var testFn;
    switch(currentView.id){
        case 0:
            if (vertical) viewDate = $a18c32cb29f9cbaf$var$addDays(viewDate, direction * 7);
            else if (ev.ctrlKey || ev.metaKey) viewDate = $a18c32cb29f9cbaf$var$addYears(viewDate, direction);
            else viewDate = $a18c32cb29f9cbaf$var$addDays(viewDate, direction);
            addFn = $a18c32cb29f9cbaf$var$addDays;
            testFn = function testFn(date) {
                return currentView.disabled.includes(date);
            };
            break;
        case 1:
            viewDate = $a18c32cb29f9cbaf$var$addMonths(viewDate, vertical ? direction * 4 : direction);
            addFn = $a18c32cb29f9cbaf$var$addMonths;
            testFn = function testFn(date) {
                var dt = new Date(date);
                var year = currentView.year, disabled = currentView.disabled;
                return dt.getFullYear() === year && disabled.includes(dt.getMonth());
            };
            break;
        default:
            viewDate = $a18c32cb29f9cbaf$var$addYears(viewDate, direction * (vertical ? 4 : 1) * step);
            addFn = $a18c32cb29f9cbaf$var$addYears;
            testFn = function testFn(date) {
                return currentView.disabled.includes($a18c32cb29f9cbaf$var$startOfYearPeriod(date, step));
            };
    }
    viewDate = $a18c32cb29f9cbaf$var$findNextAvailableOne(viewDate, addFn, direction < 0 ? -step : step, testFn, currentView.minDate, currentView.maxDate);
    if (viewDate !== undefined) picker.changeFocus(viewDate).render();
}
function $a18c32cb29f9cbaf$var$onKeydown(datepicker, ev) {
    if (ev.key === 'Tab') {
        $a18c32cb29f9cbaf$var$unfocus(datepicker);
        return;
    }
    var picker = datepicker.picker;
    var _picker$currentView = picker.currentView, id = _picker$currentView.id, isMinView = _picker$currentView.isMinView;
    if (!picker.active) switch(ev.key){
        case 'ArrowDown':
        case 'Escape':
            picker.show();
            break;
        case 'Enter':
            datepicker.update();
            break;
        default:
            return;
    }
    else if (datepicker.editMode) switch(ev.key){
        case 'Escape':
            picker.hide();
            break;
        case 'Enter':
            datepicker.exitEditMode({
                update: true,
                autohide: datepicker.config.autohide
            });
            break;
        default:
            return;
    }
    else switch(ev.key){
        case 'Escape':
            picker.hide();
            break;
        case 'ArrowLeft':
            if (ev.ctrlKey || ev.metaKey) $a18c32cb29f9cbaf$var$goToPrevOrNext(datepicker, -1);
            else if (ev.shiftKey) {
                datepicker.enterEditMode();
                return;
            } else $a18c32cb29f9cbaf$var$moveByArrowKey(datepicker, ev, -1, false);
            break;
        case 'ArrowRight':
            if (ev.ctrlKey || ev.metaKey) $a18c32cb29f9cbaf$var$goToPrevOrNext(datepicker, 1);
            else if (ev.shiftKey) {
                datepicker.enterEditMode();
                return;
            } else $a18c32cb29f9cbaf$var$moveByArrowKey(datepicker, ev, 1, false);
            break;
        case 'ArrowUp':
            if (ev.ctrlKey || ev.metaKey) $a18c32cb29f9cbaf$var$switchView(datepicker);
            else if (ev.shiftKey) {
                datepicker.enterEditMode();
                return;
            } else $a18c32cb29f9cbaf$var$moveByArrowKey(datepicker, ev, -1, true);
            break;
        case 'ArrowDown':
            if (ev.shiftKey && !ev.ctrlKey && !ev.metaKey) {
                datepicker.enterEditMode();
                return;
            }
            $a18c32cb29f9cbaf$var$moveByArrowKey(datepicker, ev, 1, true);
            break;
        case 'Enter':
            if (isMinView) datepicker.setDate(picker.viewDate);
            else picker.changeView(id - 1).render();
            break;
        case 'Backspace':
        case 'Delete':
            datepicker.enterEditMode();
            return;
        default:
            if (ev.key.length === 1 && !ev.ctrlKey && !ev.metaKey) datepicker.enterEditMode();
            return;
    }
    ev.preventDefault();
    ev.stopPropagation();
}
function $a18c32cb29f9cbaf$var$onFocus(datepicker) {
    if (datepicker.config.showOnFocus && !datepicker._showing) datepicker.show();
}
// for the prevention for entering edit mode while getting focus on click
function $a18c32cb29f9cbaf$var$onMousedown(datepicker, ev) {
    var el = ev.target;
    if (datepicker.picker.active || datepicker.config.showOnClick) {
        el._active = el === document.activeElement;
        el._clicking = setTimeout(function() {
            delete el._active;
            delete el._clicking;
        }, 2000);
    }
}
function $a18c32cb29f9cbaf$var$onClickInput(datepicker, ev) {
    var el = ev.target;
    if (!el._clicking) return;
    clearTimeout(el._clicking);
    delete el._clicking;
    if (el._active) datepicker.enterEditMode();
    delete el._active;
    if (datepicker.config.showOnClick) datepicker.show();
}
function $a18c32cb29f9cbaf$var$onPaste(datepicker, ev) {
    if (ev.clipboardData.types.includes('text/plain')) datepicker.enterEditMode();
}
// for the `document` to delegate the events from outside the picker/input field
function $a18c32cb29f9cbaf$var$onClickOutside(datepicker, ev) {
    var element = datepicker.element;
    if (element !== document.activeElement) return;
    var pickerElem = datepicker.picker.element;
    if ($a18c32cb29f9cbaf$var$findElementInEventPath(ev, function(el) {
        return el === element || el === pickerElem;
    })) return;
    $a18c32cb29f9cbaf$var$unfocus(datepicker);
}
function $a18c32cb29f9cbaf$var$stringifyDates(dates, config) {
    return dates.map(function(dt) {
        return $a18c32cb29f9cbaf$var$formatDate(dt, config.format, config.locale);
    }).join(config.dateDelimiter);
}
// parse input dates and create an array of time values for selection
// returns undefined if there are no valid dates in inputDates
// when origDates (current selection) is passed, the function works to mix
// the input dates into the current selection
function $a18c32cb29f9cbaf$var$processInputDates(datepicker, inputDates) {
    var clear = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : false;
    var config = datepicker.config, origDates = datepicker.dates, rangepicker = datepicker.rangepicker;
    if (inputDates.length === 0) // empty input is considered valid unless origiDates is passed
    return clear ? [] : undefined;
    var rangeEnd = rangepicker && datepicker === rangepicker.datepickers[1];
    var newDates = inputDates.reduce(function(dates, dt) {
        var date = $a18c32cb29f9cbaf$var$parseDate(dt, config.format, config.locale);
        if (date === undefined) return dates;
        if (config.pickLevel > 0) {
            // adjust to 1st of the month/Jan 1st of the year
            // or to the last day of the monh/Dec 31st of the year if the datepicker
            // is the range-end picker of a rangepicker
            var _dt = new Date(date);
            if (config.pickLevel === 1) date = rangeEnd ? _dt.setMonth(_dt.getMonth() + 1, 0) : _dt.setDate(1);
            else date = rangeEnd ? _dt.setFullYear(_dt.getFullYear() + 1, 0, 0) : _dt.setMonth(0, 1);
        }
        if ($a18c32cb29f9cbaf$var$isInRange(date, config.minDate, config.maxDate) && !dates.includes(date) && !config.datesDisabled.includes(date) && !config.daysOfWeekDisabled.includes(new Date(date).getDay())) dates.push(date);
        return dates;
    }, []);
    if (newDates.length === 0) return;
    if (config.multidate && !clear) // get the synmetric difference between origDates and newDates
    newDates = newDates.reduce(function(dates, date) {
        if (!origDates.includes(date)) dates.push(date);
        return dates;
    }, origDates.filter(function(date) {
        return !newDates.includes(date);
    }));
    // do length check always because user can input multiple dates regardless of the mode
    return config.maxNumberOfDates && newDates.length > config.maxNumberOfDates ? newDates.slice(config.maxNumberOfDates * -1) : newDates;
}
// refresh the UI elements
// modes: 1: input only, 2, picker only, 3 both
function $a18c32cb29f9cbaf$var$refreshUI(datepicker) {
    var mode = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : 3;
    var quickRender = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : true;
    var config = datepicker.config, picker = datepicker.picker, inputField = datepicker.inputField;
    if (mode & 2) {
        var newView = picker.active ? config.pickLevel : config.startView;
        picker.update().changeView(newView).render(quickRender);
    }
    if (mode & 1 && inputField) inputField.value = $a18c32cb29f9cbaf$var$stringifyDates(datepicker.dates, config);
}
function $a18c32cb29f9cbaf$var$_setDate(datepicker, inputDates, options) {
    var clear = options.clear, render = options.render, autohide = options.autohide;
    if (render === undefined) render = true;
    if (!render) autohide = false;
    else if (autohide === undefined) autohide = datepicker.config.autohide;
    var newDates = $a18c32cb29f9cbaf$var$processInputDates(datepicker, inputDates, clear);
    if (!newDates) return;
    if (newDates.toString() !== datepicker.dates.toString()) {
        datepicker.dates = newDates;
        $a18c32cb29f9cbaf$var$refreshUI(datepicker, render ? 3 : 1);
        $a18c32cb29f9cbaf$var$triggerDatepickerEvent(datepicker, 'changeDate');
    } else $a18c32cb29f9cbaf$var$refreshUI(datepicker, 1);
    if (autohide) datepicker.hide();
}
/**
 * Class representing a date picker
 */ var $a18c32cb29f9cbaf$export$7235422eca03ec90 = /*#__PURE__*/ function() {
    /**
   * Create a date picker
   * @param  {Element} element - element to bind a date picker
   * @param  {Object} [options] - config options
   * @param  {DateRangePicker} [rangepicker] - DateRangePicker instance the
   * date picker belongs to. Use this only when creating date picker as a part
   * of date range picker
   */ function Datepicker(element) {
        var options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
        var rangepicker = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : undefined;
        $a18c32cb29f9cbaf$var$_classCallCheck(this, Datepicker);
        element.datepicker = this;
        this.element = element;
        // set up config
        var config = this.config = Object.assign({
            buttonClass: options.buttonClass && String(options.buttonClass) || 'button',
            container: document.body,
            defaultViewDate: $a18c32cb29f9cbaf$var$today(),
            maxDate: undefined,
            minDate: undefined
        }, $a18c32cb29f9cbaf$var$processOptions($a18c32cb29f9cbaf$var$defaultOptions, this));
        this._options = options;
        Object.assign(config, $a18c32cb29f9cbaf$var$processOptions(options, this));
        // configure by type
        var inline = this.inline = element.tagName !== 'INPUT';
        var inputField;
        var initialDates;
        if (inline) {
            config.container = element;
            initialDates = $a18c32cb29f9cbaf$var$stringToArray(element.dataset.date, config.dateDelimiter);
            delete element.dataset.date;
        } else {
            var container = options.container ? document.querySelector(options.container) : null;
            if (container) config.container = container;
            inputField = this.inputField = element;
            inputField.classList.add('datepicker-input');
            initialDates = $a18c32cb29f9cbaf$var$stringToArray(inputField.value, config.dateDelimiter);
        }
        if (rangepicker) {
            // check validiry
            var index = rangepicker.inputs.indexOf(inputField);
            var datepickers = rangepicker.datepickers;
            if (index < 0 || index > 1 || !Array.isArray(datepickers)) throw Error('Invalid rangepicker object.');
            // attach itaelf to the rangepicker here so that processInputDates() can
            // determine if this is the range-end picker of the rangepicker while
            // setting inital values when pickLevel > 0
            datepickers[index] = this;
            // add getter for rangepicker
            Object.defineProperty(this, 'rangepicker', {
                get: function get() {
                    return rangepicker;
                }
            });
        }
        // set initial dates
        this.dates = [];
        // process initial value
        var inputDateValues = $a18c32cb29f9cbaf$var$processInputDates(this, initialDates);
        if (inputDateValues && inputDateValues.length > 0) this.dates = inputDateValues;
        if (inputField) inputField.value = $a18c32cb29f9cbaf$var$stringifyDates(this.dates, config);
        var picker = this.picker = new $a18c32cb29f9cbaf$var$Picker(this);
        if (inline) this.show();
        else {
            // set up event listeners in other modes
            var onMousedownDocument = $a18c32cb29f9cbaf$var$onClickOutside.bind(null, this);
            var listeners = [
                [
                    inputField,
                    'keydown',
                    $a18c32cb29f9cbaf$var$onKeydown.bind(null, this)
                ],
                [
                    inputField,
                    'focus',
                    $a18c32cb29f9cbaf$var$onFocus.bind(null, this)
                ],
                [
                    inputField,
                    'mousedown',
                    $a18c32cb29f9cbaf$var$onMousedown.bind(null, this)
                ],
                [
                    inputField,
                    'click',
                    $a18c32cb29f9cbaf$var$onClickInput.bind(null, this)
                ],
                [
                    inputField,
                    'paste',
                    $a18c32cb29f9cbaf$var$onPaste.bind(null, this)
                ],
                [
                    document,
                    'mousedown',
                    onMousedownDocument
                ],
                [
                    document,
                    'touchstart',
                    onMousedownDocument
                ],
                [
                    window,
                    'resize',
                    picker.place.bind(picker)
                ]
            ];
            $a18c32cb29f9cbaf$var$registerListeners(this, listeners);
        }
    }
    /**
   * Format Date object or time value in given format and language
   * @param  {Date|Number} date - date or time value to format
   * @param  {String|Object} format - format string or object that contains
   * toDisplay() custom formatter, whose signature is
   * - args:
   *   - date: {Date} - Date instance of the date passed to the method
   *   - format: {Object} - the format object passed to the method
   *   - locale: {Object} - locale for the language specified by `lang`
   * - return:
   *     {String} formatted date
   * @param  {String} [lang=en] - language code for the locale to use
   * @return {String} formatted date
   */ return $a18c32cb29f9cbaf$var$_createClass(Datepicker, [
        {
            key: "active",
            get: /**
     * @type {Boolean} - Whether the picker element is shown. `true` whne shown
     */ function get() {
                return !!(this.picker && this.picker.active);
            }
        },
        {
            key: "pickerElement",
            get: function get() {
                return this.picker ? this.picker.element : undefined;
            }
        },
        {
            key: "setOptions",
            value: function setOptions(options) {
                var picker = this.picker;
                var newOptions = $a18c32cb29f9cbaf$var$processOptions(options, this);
                Object.assign(this._options, options);
                Object.assign(this.config, newOptions);
                picker.setOptions(newOptions);
                $a18c32cb29f9cbaf$var$refreshUI(this, 3);
            }
        },
        {
            key: "show",
            value: function show() {
                if (this.inputField) {
                    if (this.inputField.disabled) return;
                    if (this.inputField !== document.activeElement) {
                        this._showing = true;
                        this.inputField.focus();
                        delete this._showing;
                    }
                }
                this.picker.show();
            }
        },
        {
            key: "hide",
            value: function hide() {
                if (this.inline) return;
                this.picker.hide();
                this.picker.update().changeView(this.config.startView).render();
            }
        },
        {
            key: "destroy",
            value: function destroy() {
                this.hide();
                $a18c32cb29f9cbaf$var$unregisterListeners(this);
                this.picker.detach();
                if (!this.inline) this.inputField.classList.remove('datepicker-input');
                delete this.element.datepicker;
                return this;
            }
        },
        {
            key: "getDate",
            value: function getDate() {
                var _this = this;
                var format = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : undefined;
                var callback = format ? function(date) {
                    return $a18c32cb29f9cbaf$var$formatDate(date, format, _this.config.locale);
                } : function(date) {
                    return new Date(date);
                };
                if (this.config.multidate) return this.dates.map(callback);
                if (this.dates.length > 0) return callback(this.dates[0]);
            }
        },
        {
            key: "setDate",
            value: function setDate() {
                for(var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++)args[_key] = arguments[_key];
                var dates = [].concat(args);
                var opts = {};
                var lastArg = $a18c32cb29f9cbaf$var$lastItemOf(args);
                if ($a18c32cb29f9cbaf$var$_typeof(lastArg) === 'object' && !Array.isArray(lastArg) && !(lastArg instanceof Date) && lastArg) Object.assign(opts, dates.pop());
                var inputDates = Array.isArray(dates[0]) ? dates[0] : dates;
                $a18c32cb29f9cbaf$var$_setDate(this, inputDates, opts);
            }
        },
        {
            key: "update",
            value: function update() {
                var options = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : undefined;
                if (this.inline) return;
                var opts = {
                    clear: true,
                    autohide: !!(options && options.autohide)
                };
                var inputDates = $a18c32cb29f9cbaf$var$stringToArray(this.inputField.value, this.config.dateDelimiter);
                $a18c32cb29f9cbaf$var$_setDate(this, inputDates, opts);
            }
        },
        {
            key: "refresh",
            value: function refresh() {
                var target = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : undefined;
                var forceRender = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : false;
                if (target && typeof target !== 'string') {
                    forceRender = target;
                    target = undefined;
                }
                var mode;
                if (target === 'picker') mode = 2;
                else if (target === 'input') mode = 1;
                else mode = 3;
                $a18c32cb29f9cbaf$var$refreshUI(this, mode, !forceRender);
            }
        },
        {
            key: "enterEditMode",
            value: function enterEditMode() {
                if (this.inline || !this.picker.active || this.editMode) return;
                this.editMode = true;
                this.inputField.classList.add('in-edit', 'border-brand');
            }
        },
        {
            key: "exitEditMode",
            value: function exitEditMode() {
                var options = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : undefined;
                if (this.inline || !this.editMode) return;
                var opts = Object.assign({
                    update: false
                }, options);
                delete this.editMode;
                this.inputField.classList.remove('in-edit', 'border-brand');
                if (opts.update) this.update(opts);
            }
        }
    ], [
        {
            key: "formatDate",
            value: function formatDate$1(date, format, lang) {
                return $a18c32cb29f9cbaf$var$formatDate(date, format, lang && $a18c32cb29f9cbaf$var$locales[lang] || $a18c32cb29f9cbaf$var$locales.en);
            }
        },
        {
            key: "parseDate",
            value: function parseDate$1(dateStr, format, lang) {
                return $a18c32cb29f9cbaf$var$parseDate(dateStr, format, lang && $a18c32cb29f9cbaf$var$locales[lang] || $a18c32cb29f9cbaf$var$locales.en);
            }
        },
        {
            key: "locales",
            get: function get() {
                return $a18c32cb29f9cbaf$var$locales;
            }
        }
    ]);
}();
// filter out the config options inapproprite to pass to Datepicker
function $a18c32cb29f9cbaf$var$filterOptions(options) {
    var newOpts = Object.assign({}, options);
    delete newOpts.inputs;
    delete newOpts.allowOneSidedRange;
    delete newOpts.maxNumberOfDates; // to ensure each datepicker handles a single date
    return newOpts;
}
function $a18c32cb29f9cbaf$var$setupDatepicker(rangepicker, changeDateListener, el, options) {
    $a18c32cb29f9cbaf$var$registerListeners(rangepicker, [
        [
            el,
            'changeDate',
            changeDateListener
        ]
    ]);
    new $a18c32cb29f9cbaf$export$7235422eca03ec90(el, options, rangepicker);
}
function $a18c32cb29f9cbaf$var$onChangeDate(rangepicker, ev) {
    // to prevent both datepickers trigger the other side's update each other
    if (rangepicker._updating) return;
    rangepicker._updating = true;
    var target = ev.target;
    if (target.datepicker === undefined) return;
    var datepickers = rangepicker.datepickers;
    var setDateOptions = {
        render: false
    };
    var changedSide = rangepicker.inputs.indexOf(target);
    var otherSide = changedSide === 0 ? 1 : 0;
    var changedDate = datepickers[changedSide].dates[0];
    var otherDate = datepickers[otherSide].dates[0];
    if (changedDate !== undefined && otherDate !== undefined) {
        // if the start of the range > the end, swap them
        if (changedSide === 0 && changedDate > otherDate) {
            datepickers[0].setDate(otherDate, setDateOptions);
            datepickers[1].setDate(changedDate, setDateOptions);
        } else if (changedSide === 1 && changedDate < otherDate) {
            datepickers[0].setDate(changedDate, setDateOptions);
            datepickers[1].setDate(otherDate, setDateOptions);
        }
    } else if (!rangepicker.allowOneSidedRange) // to prevent the range from becoming one-sided, copy changed side's
    // selection (no matter if it's empty) to the other side
    {
        if (changedDate !== undefined || otherDate !== undefined) {
            setDateOptions.clear = true;
            datepickers[otherSide].setDate(datepickers[changedSide].dates, setDateOptions);
        }
    }
    datepickers[0].picker.update().render();
    datepickers[1].picker.update().render();
    delete rangepicker._updating;
}
/**
 * Class representing a date range picker
 */ var $a18c32cb29f9cbaf$export$17334619f3ac2224 = /*#__PURE__*/ function() {
    /**
   * Create a date range picker
   * @param  {Element} element - element to bind a date range picker
   * @param  {Object} [options] - config options
   */ function DateRangePicker(element) {
        var options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
        $a18c32cb29f9cbaf$var$_classCallCheck(this, DateRangePicker);
        var inputs = Array.isArray(options.inputs) ? options.inputs : Array.from(element.querySelectorAll('input'));
        if (inputs.length < 2) return;
        element.rangepicker = this;
        this.element = element;
        this.inputs = inputs.slice(0, 2);
        this.allowOneSidedRange = !!options.allowOneSidedRange;
        var changeDateListener = $a18c32cb29f9cbaf$var$onChangeDate.bind(null, this);
        var cleanOptions = $a18c32cb29f9cbaf$var$filterOptions(options);
        // in order for initial date setup to work right when pcicLvel > 0,
        // let Datepicker constructor add the instance to the rangepicker
        var datepickers = [];
        Object.defineProperty(this, 'datepickers', {
            get: function get() {
                return datepickers;
            }
        });
        $a18c32cb29f9cbaf$var$setupDatepicker(this, changeDateListener, this.inputs[0], cleanOptions);
        $a18c32cb29f9cbaf$var$setupDatepicker(this, changeDateListener, this.inputs[1], cleanOptions);
        Object.freeze(datepickers);
        // normalize the range if inital dates are given
        if (datepickers[0].dates.length > 0) $a18c32cb29f9cbaf$var$onChangeDate(this, {
            target: this.inputs[0]
        });
        else if (datepickers[1].dates.length > 0) $a18c32cb29f9cbaf$var$onChangeDate(this, {
            target: this.inputs[1]
        });
    }
    /**
   * @type {Array} - selected date of the linked date pickers
   */ return $a18c32cb29f9cbaf$var$_createClass(DateRangePicker, [
        {
            key: "dates",
            get: function get() {
                return this.datepickers.length === 2 ? [
                    this.datepickers[0].dates[0],
                    this.datepickers[1].dates[0]
                ] : undefined;
            }
        },
        {
            key: "setOptions",
            value: function setOptions(options) {
                this.allowOneSidedRange = !!options.allowOneSidedRange;
                var cleanOptions = $a18c32cb29f9cbaf$var$filterOptions(options);
                this.datepickers[0].setOptions(cleanOptions);
                this.datepickers[1].setOptions(cleanOptions);
            }
        },
        {
            key: "destroy",
            value: function destroy() {
                this.datepickers[0].destroy();
                this.datepickers[1].destroy();
                $a18c32cb29f9cbaf$var$unregisterListeners(this);
                delete this.element.rangepicker;
            }
        },
        {
            key: "getDates",
            value: function getDates() {
                var _this = this;
                var format = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : undefined;
                var callback = format ? function(date) {
                    return $a18c32cb29f9cbaf$var$formatDate(date, format, _this.datepickers[0].config.locale);
                } : function(date) {
                    return new Date(date);
                };
                return this.dates.map(function(date) {
                    return date === undefined ? date : callback(date);
                });
            }
        },
        {
            key: "setDates",
            value: function setDates(rangeStart, rangeEnd) {
                var _this$datepickers = $a18c32cb29f9cbaf$var$_slicedToArray(this.datepickers, 2), datepicker0 = _this$datepickers[0], datepicker1 = _this$datepickers[1];
                var origDates = this.dates;
                // If range normalization runs on every change, we can't set a new range
                // that starts after the end of the current range correctly because the
                // normalization process swaps start↔︎end right after setting the new start
                // date. To prevent this, the normalization process needs to run once after
                // both of the new dates are set.
                this._updating = true;
                datepicker0.setDate(rangeStart);
                datepicker1.setDate(rangeEnd);
                delete this._updating;
                if (datepicker1.dates[0] !== origDates[1]) $a18c32cb29f9cbaf$var$onChangeDate(this, {
                    target: this.inputs[1]
                });
                else if (datepicker0.dates[0] !== origDates[0]) $a18c32cb29f9cbaf$var$onChangeDate(this, {
                    target: this.inputs[0]
                });
            }
        }
    ]);
}();


var $6fd16f8bbe324eab$var$__assign = undefined && undefined.__assign || function() {
    $6fd16f8bbe324eab$var$__assign = Object.assign || function(t) {
        for(var s, i = 1, n = arguments.length; i < n; i++){
            s = arguments[i];
            for(var p in s)if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return $6fd16f8bbe324eab$var$__assign.apply(this, arguments);
};
var $6fd16f8bbe324eab$var$Default = {
    defaultDatepickerId: null,
    autohide: false,
    format: 'mm/dd/yyyy',
    maxDate: null,
    minDate: null,
    orientation: 'bottom',
    buttons: false,
    autoSelectToday: 0,
    title: null,
    language: 'en',
    rangePicker: false,
    onShow: function() {},
    onHide: function() {}
};
var $6fd16f8bbe324eab$var$DefaultInstanceOptions = {
    id: null,
    override: true
};
var $6fd16f8bbe324eab$var$Datepicker = /** @class */ function() {
    function Datepicker(datepickerEl, options, instanceOptions) {
        if (datepickerEl === void 0) datepickerEl = null;
        if (options === void 0) options = $6fd16f8bbe324eab$var$Default;
        if (instanceOptions === void 0) instanceOptions = $6fd16f8bbe324eab$var$DefaultInstanceOptions;
        this._instanceId = instanceOptions.id ? instanceOptions.id : datepickerEl.id;
        this._datepickerEl = datepickerEl;
        this._datepickerInstance = null;
        this._options = $6fd16f8bbe324eab$var$__assign($6fd16f8bbe324eab$var$__assign({}, $6fd16f8bbe324eab$var$Default), options);
        this._initialized = false;
        this.init();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).addInstance('Datepicker', this, this._instanceId, instanceOptions.override);
    }
    Datepicker.prototype.init = function() {
        if (this._datepickerEl && !this._initialized) {
            if (this._options.rangePicker) this._datepickerInstance = new (0, $a18c32cb29f9cbaf$export$17334619f3ac2224)(this._datepickerEl, this._getDatepickerOptions(this._options));
            else this._datepickerInstance = new (0, $a18c32cb29f9cbaf$export$7235422eca03ec90)(this._datepickerEl, this._getDatepickerOptions(this._options));
            this._initialized = true;
        }
    };
    Datepicker.prototype.destroy = function() {
        if (this._initialized) {
            this._initialized = false;
            this._datepickerInstance.destroy();
        }
    };
    Datepicker.prototype.removeInstance = function() {
        this.destroy();
        (0, $43290e7e11b6c4c9$export$2e2bcd8739ae039).removeInstance('Datepicker', this._instanceId);
    };
    Datepicker.prototype.destroyAndRemoveInstance = function() {
        this.destroy();
        this.removeInstance();
    };
    Datepicker.prototype.getDatepickerInstance = function() {
        return this._datepickerInstance;
    };
    Datepicker.prototype.getDate = function() {
        if (this._options.rangePicker && this._datepickerInstance instanceof (0, $a18c32cb29f9cbaf$export$17334619f3ac2224)) return this._datepickerInstance.getDates();
        if (!this._options.rangePicker && this._datepickerInstance instanceof (0, $a18c32cb29f9cbaf$export$7235422eca03ec90)) return this._datepickerInstance.getDate();
    };
    Datepicker.prototype.setDate = function(date) {
        if (this._options.rangePicker && this._datepickerInstance instanceof (0, $a18c32cb29f9cbaf$export$17334619f3ac2224)) return this._datepickerInstance.setDates(date);
        if (!this._options.rangePicker && this._datepickerInstance instanceof (0, $a18c32cb29f9cbaf$export$7235422eca03ec90)) return this._datepickerInstance.setDate(date);
    };
    Datepicker.prototype.show = function() {
        this._datepickerInstance.show();
        this._options.onShow(this);
    };
    Datepicker.prototype.hide = function() {
        this._datepickerInstance.hide();
        this._options.onHide(this);
    };
    Datepicker.prototype._getDatepickerOptions = function(options) {
        var datepickerOptions = {};
        if (options.buttons) {
            datepickerOptions.todayBtn = true;
            datepickerOptions.clearBtn = true;
            if (options.autoSelectToday) datepickerOptions.todayBtnMode = 1;
        }
        if (options.autohide) datepickerOptions.autohide = true;
        if (options.format) datepickerOptions.format = options.format;
        if (options.maxDate) datepickerOptions.maxDate = options.maxDate;
        if (options.minDate) datepickerOptions.minDate = options.minDate;
        if (options.orientation) datepickerOptions.orientation = options.orientation;
        if (options.title) datepickerOptions.title = options.title;
        if (options.language) datepickerOptions.language = options.language;
        return datepickerOptions;
    };
    Datepicker.prototype.updateOnShow = function(callback) {
        this._options.onShow = callback;
    };
    Datepicker.prototype.updateOnHide = function(callback) {
        this._options.onHide = callback;
    };
    return Datepicker;
}();
function $6fd16f8bbe324eab$export$81c9f52a292162ae() {
    document.querySelectorAll('[datepicker], [inline-datepicker], [date-rangepicker]').forEach(function($datepickerEl) {
        if ($datepickerEl) {
            var buttons = $datepickerEl.hasAttribute('datepicker-buttons');
            var autoselectToday = $datepickerEl.hasAttribute('datepicker-autoselect-today');
            var autohide = $datepickerEl.hasAttribute('datepicker-autohide');
            var format = $datepickerEl.getAttribute('datepicker-format');
            var maxDate = $datepickerEl.getAttribute('datepicker-max-date');
            var minDate = $datepickerEl.getAttribute('datepicker-min-date');
            var orientation_1 = $datepickerEl.getAttribute('datepicker-orientation');
            var title = $datepickerEl.getAttribute('datepicker-title');
            var language = $datepickerEl.getAttribute('datepicker-language');
            var rangePicker = $datepickerEl.hasAttribute('date-rangepicker');
            new $6fd16f8bbe324eab$var$Datepicker($datepickerEl, {
                buttons: buttons ? buttons : $6fd16f8bbe324eab$var$Default.buttons,
                autoSelectToday: autoselectToday ? autoselectToday : $6fd16f8bbe324eab$var$Default.autoSelectToday,
                autohide: autohide ? autohide : $6fd16f8bbe324eab$var$Default.autohide,
                format: format ? format : $6fd16f8bbe324eab$var$Default.format,
                maxDate: maxDate ? maxDate : $6fd16f8bbe324eab$var$Default.maxDate,
                minDate: minDate ? minDate : $6fd16f8bbe324eab$var$Default.minDate,
                orientation: orientation_1 ? orientation_1 : $6fd16f8bbe324eab$var$Default.orientation,
                title: title ? title : $6fd16f8bbe324eab$var$Default.title,
                language: language ? language : $6fd16f8bbe324eab$var$Default.language,
                rangePicker: rangePicker ? rangePicker : $6fd16f8bbe324eab$var$Default.rangePicker
            });
        } else console.error("The datepicker element does not exist. Please check the datepicker attribute.");
    });
}
if (typeof window !== 'undefined') {
    window.Datepicker = $6fd16f8bbe324eab$var$Datepicker;
    window.initDatepickers = $6fd16f8bbe324eab$export$81c9f52a292162ae;
}
var $6fd16f8bbe324eab$export$2e2bcd8739ae039 = $6fd16f8bbe324eab$var$Datepicker;
















function $abd8bfe2e23ba646$export$4cd50b9c69a85fd5() {
    (0, $946dcfd8fe0754f4$export$226c1dc98323ee4d)();
    (0, $16ba7ad25d51d4f8$export$355ba5a528b4009a)();
    (0, $6a8fcd5cc87cb99e$export$3ab77386b16b9e58)();
    (0, $fa24e680ac3fba2d$export$69c84787a503daef)();
    (0, $ece16b3047178a9c$export$8cb65a02593bf108)();
    (0, $8d24653c8cdd4795$export$e14d3bc873a8b9b1)();
    (0, $c647e009fa86f7f0$export$319df2b10e87c8ec)();
    (0, $ff078cf47b3c0903$export$c92e8f569ad2976f)();
    (0, $7ab0a41fd48c3b4d$export$8f2a38c1f9d9dc3e)();
    (0, $5370c2686863c309$export$200409e83b2a0dd4)();
    (0, $4fe34b241a1ad5f0$export$33aa68e0bfaa065d)();
    (0, $df0fc5a35c362eab$export$4352dd3b7a1676ab)();
    (0, $f2c2b31007186eb1$export$ad4117886fd6cf74)();
    (0, $6fd16f8bbe324eab$export$81c9f52a292162ae)();
}
if (typeof window !== 'undefined') window.initFlowbite = $abd8bfe2e23ba646$export$4cd50b9c69a85fd5;






























































// setup events for data attributes
var $68c8d9028d8b1de7$var$events = new (0, $4c6a2f83122c31a3$export$2e2bcd8739ae039)('load', [
    (0, $946dcfd8fe0754f4$export$226c1dc98323ee4d),
    (0, $16ba7ad25d51d4f8$export$355ba5a528b4009a),
    (0, $6a8fcd5cc87cb99e$export$3ab77386b16b9e58),
    (0, $fa24e680ac3fba2d$export$69c84787a503daef),
    (0, $ece16b3047178a9c$export$8cb65a02593bf108),
    (0, $8d24653c8cdd4795$export$e14d3bc873a8b9b1),
    (0, $c647e009fa86f7f0$export$319df2b10e87c8ec),
    (0, $ff078cf47b3c0903$export$c92e8f569ad2976f),
    (0, $7ab0a41fd48c3b4d$export$8f2a38c1f9d9dc3e),
    (0, $5370c2686863c309$export$200409e83b2a0dd4),
    (0, $4fe34b241a1ad5f0$export$33aa68e0bfaa065d),
    (0, $df0fc5a35c362eab$export$4352dd3b7a1676ab),
    (0, $f2c2b31007186eb1$export$ad4117886fd6cf74),
    (0, $6fd16f8bbe324eab$export$81c9f52a292162ae)
]);
$68c8d9028d8b1de7$var$events.init();


const $1bac384020b50752$var$isString = (obj)=>typeof obj === 'string';
const $1bac384020b50752$var$defer = ()=>{
    let res;
    let rej;
    const promise = new Promise((resolve, reject)=>{
        res = resolve;
        rej = reject;
    });
    promise.resolve = res;
    promise.reject = rej;
    return promise;
};
const $1bac384020b50752$var$makeString = (object)=>{
    if (object == null) return '';
    return '' + object;
};
const $1bac384020b50752$var$copy = (a, s, t)=>{
    a.forEach((m)=>{
        if (s[m]) t[m] = s[m];
    });
};
const $1bac384020b50752$var$lastOfPathSeparatorRegExp = /###/g;
const $1bac384020b50752$var$cleanKey = (key)=>key && key.indexOf('###') > -1 ? key.replace($1bac384020b50752$var$lastOfPathSeparatorRegExp, '.') : key;
const $1bac384020b50752$var$canNotTraverseDeeper = (object)=>!object || $1bac384020b50752$var$isString(object);
const $1bac384020b50752$var$getLastOfPath = (object, path, Empty)=>{
    const stack = !$1bac384020b50752$var$isString(path) ? path : path.split('.');
    let stackIndex = 0;
    while(stackIndex < stack.length - 1){
        if ($1bac384020b50752$var$canNotTraverseDeeper(object)) return {};
        const key = $1bac384020b50752$var$cleanKey(stack[stackIndex]);
        if (!object[key] && Empty) object[key] = new Empty();
        if (Object.prototype.hasOwnProperty.call(object, key)) object = object[key];
        else object = {};
        ++stackIndex;
    }
    if ($1bac384020b50752$var$canNotTraverseDeeper(object)) return {};
    return {
        obj: object,
        k: $1bac384020b50752$var$cleanKey(stack[stackIndex])
    };
};
const $1bac384020b50752$var$setPath = (object, path, newValue)=>{
    const { obj: obj, k: k } = $1bac384020b50752$var$getLastOfPath(object, path, Object);
    if (obj !== undefined || path.length === 1) {
        obj[k] = newValue;
        return;
    }
    let e = path[path.length - 1];
    let p = path.slice(0, path.length - 1);
    let last = $1bac384020b50752$var$getLastOfPath(object, p, Object);
    while(last.obj === undefined && p.length){
        e = `${p[p.length - 1]}.${e}`;
        p = p.slice(0, p.length - 1);
        last = $1bac384020b50752$var$getLastOfPath(object, p, Object);
        if (last?.obj && typeof last.obj[`${last.k}.${e}`] !== 'undefined') last.obj = undefined;
    }
    last.obj[`${last.k}.${e}`] = newValue;
};
const $1bac384020b50752$var$pushPath = (object, path, newValue, concat)=>{
    const { obj: obj, k: k } = $1bac384020b50752$var$getLastOfPath(object, path, Object);
    obj[k] = obj[k] || [];
    obj[k].push(newValue);
};
const $1bac384020b50752$var$getPath = (object, path)=>{
    const { obj: obj, k: k } = $1bac384020b50752$var$getLastOfPath(object, path);
    if (!obj) return undefined;
    if (!Object.prototype.hasOwnProperty.call(obj, k)) return undefined;
    return obj[k];
};
const $1bac384020b50752$var$getPathWithDefaults = (data, defaultData, key)=>{
    const value = $1bac384020b50752$var$getPath(data, key);
    if (value !== undefined) return value;
    return $1bac384020b50752$var$getPath(defaultData, key);
};
const $1bac384020b50752$var$deepExtend = (target, source, overwrite)=>{
    for(const prop in source)if (prop !== '__proto__' && prop !== 'constructor') {
        if (prop in target) {
            if ($1bac384020b50752$var$isString(target[prop]) || target[prop] instanceof String || $1bac384020b50752$var$isString(source[prop]) || source[prop] instanceof String) {
                if (overwrite) target[prop] = source[prop];
            } else $1bac384020b50752$var$deepExtend(target[prop], source[prop], overwrite);
        } else target[prop] = source[prop];
    }
    return target;
};
const $1bac384020b50752$var$regexEscape = (str)=>str.replace(/[\-\[\]\/\{\}\(\)\*\+\?\.\\\^\$\|]/g, '\\$&');
var $1bac384020b50752$var$_entityMap = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#39;',
    '/': '&#x2F;'
};
const $1bac384020b50752$var$escape = (data)=>{
    if ($1bac384020b50752$var$isString(data)) return data.replace(/[&<>"'\/]/g, (s)=>$1bac384020b50752$var$_entityMap[s]);
    return data;
};
class $1bac384020b50752$var$RegExpCache {
    constructor(capacity){
        this.capacity = capacity;
        this.regExpMap = new Map();
        this.regExpQueue = [];
    }
    getRegExp(pattern) {
        const regExpFromCache = this.regExpMap.get(pattern);
        if (regExpFromCache !== undefined) return regExpFromCache;
        const regExpNew = new RegExp(pattern);
        if (this.regExpQueue.length === this.capacity) this.regExpMap.delete(this.regExpQueue.shift());
        this.regExpMap.set(pattern, regExpNew);
        this.regExpQueue.push(pattern);
        return regExpNew;
    }
}
const $1bac384020b50752$var$chars = [
    ' ',
    ',',
    '?',
    '!',
    ';'
];
const $1bac384020b50752$var$looksLikeObjectPathRegExpCache = new $1bac384020b50752$var$RegExpCache(20);
const $1bac384020b50752$var$looksLikeObjectPath = (key, nsSeparator, keySeparator)=>{
    nsSeparator = nsSeparator || '';
    keySeparator = keySeparator || '';
    const possibleChars = $1bac384020b50752$var$chars.filter((c)=>nsSeparator.indexOf(c) < 0 && keySeparator.indexOf(c) < 0);
    if (possibleChars.length === 0) return true;
    const r = $1bac384020b50752$var$looksLikeObjectPathRegExpCache.getRegExp(`(${possibleChars.map((c)=>c === '?' ? '\\?' : c).join('|')})`);
    let matched = !r.test(key);
    if (!matched) {
        const ki = key.indexOf(keySeparator);
        if (ki > 0 && !r.test(key.substring(0, ki))) matched = true;
    }
    return matched;
};
const $1bac384020b50752$var$deepFind = (obj, path, keySeparator = '.')=>{
    if (!obj) return undefined;
    if (obj[path]) {
        if (!Object.prototype.hasOwnProperty.call(obj, path)) return undefined;
        return obj[path];
    }
    const tokens = path.split(keySeparator);
    let current = obj;
    for(let i = 0; i < tokens.length;){
        if (!current || typeof current !== 'object') return undefined;
        let next;
        let nextPath = '';
        for(let j = i; j < tokens.length; ++j){
            if (j !== i) nextPath += keySeparator;
            nextPath += tokens[j];
            next = current[nextPath];
            if (next !== undefined) {
                if ([
                    'string',
                    'number',
                    'boolean'
                ].indexOf(typeof next) > -1 && j < tokens.length - 1) continue;
                i += j - i + 1;
                break;
            }
        }
        current = next;
    }
    return current;
};
const $1bac384020b50752$var$getCleanedCode = (code)=>code?.replace('_', '-');
const $1bac384020b50752$var$consoleLogger = {
    type: 'logger',
    log (args) {
        this.output('log', args);
    },
    warn (args) {
        this.output('warn', args);
    },
    error (args) {
        this.output('error', args);
    },
    output (type, args) {
        console?.[type]?.apply?.(console, args);
    }
};
class $1bac384020b50752$var$Logger {
    constructor(concreteLogger, options = {}){
        this.init(concreteLogger, options);
    }
    init(concreteLogger, options = {}) {
        this.prefix = options.prefix || 'i18next:';
        this.logger = concreteLogger || $1bac384020b50752$var$consoleLogger;
        this.options = options;
        this.debug = options.debug;
    }
    log(...args) {
        return this.forward(args, 'log', '', true);
    }
    warn(...args) {
        return this.forward(args, 'warn', '', true);
    }
    error(...args) {
        return this.forward(args, 'error', '');
    }
    deprecate(...args) {
        return this.forward(args, 'warn', 'WARNING DEPRECATED: ', true);
    }
    forward(args, lvl, prefix, debugOnly) {
        if (debugOnly && !this.debug) return null;
        if ($1bac384020b50752$var$isString(args[0])) args[0] = `${prefix}${this.prefix} ${args[0]}`;
        return this.logger[lvl](args);
    }
    create(moduleName) {
        return new $1bac384020b50752$var$Logger(this.logger, {
            prefix: `${this.prefix}:${moduleName}:`,
            ...this.options
        });
    }
    clone(options) {
        options = options || this.options;
        options.prefix = options.prefix || this.prefix;
        return new $1bac384020b50752$var$Logger(this.logger, options);
    }
}
var $1bac384020b50752$var$baseLogger = new $1bac384020b50752$var$Logger();
class $1bac384020b50752$var$EventEmitter {
    constructor(){
        this.observers = {};
    }
    on(events, listener) {
        events.split(' ').forEach((event)=>{
            if (!this.observers[event]) this.observers[event] = new Map();
            const numListeners = this.observers[event].get(listener) || 0;
            this.observers[event].set(listener, numListeners + 1);
        });
        return this;
    }
    off(event, listener) {
        if (!this.observers[event]) return;
        if (!listener) {
            delete this.observers[event];
            return;
        }
        this.observers[event].delete(listener);
    }
    emit(event, ...args) {
        if (this.observers[event]) {
            const cloned = Array.from(this.observers[event].entries());
            cloned.forEach(([observer, numTimesAdded])=>{
                for(let i = 0; i < numTimesAdded; i++)observer(...args);
            });
        }
        if (this.observers['*']) {
            const cloned = Array.from(this.observers['*'].entries());
            cloned.forEach(([observer, numTimesAdded])=>{
                for(let i = 0; i < numTimesAdded; i++)observer.apply(observer, [
                    event,
                    ...args
                ]);
            });
        }
    }
}
class $1bac384020b50752$var$ResourceStore extends $1bac384020b50752$var$EventEmitter {
    constructor(data, options = {
        ns: [
            'translation'
        ],
        defaultNS: 'translation'
    }){
        super();
        this.data = data || {};
        this.options = options;
        if (this.options.keySeparator === undefined) this.options.keySeparator = '.';
        if (this.options.ignoreJSONStructure === undefined) this.options.ignoreJSONStructure = true;
    }
    addNamespaces(ns) {
        if (this.options.ns.indexOf(ns) < 0) this.options.ns.push(ns);
    }
    removeNamespaces(ns) {
        const index = this.options.ns.indexOf(ns);
        if (index > -1) this.options.ns.splice(index, 1);
    }
    getResource(lng, ns, key, options = {}) {
        const keySeparator = options.keySeparator !== undefined ? options.keySeparator : this.options.keySeparator;
        const ignoreJSONStructure = options.ignoreJSONStructure !== undefined ? options.ignoreJSONStructure : this.options.ignoreJSONStructure;
        let path;
        if (lng.indexOf('.') > -1) path = lng.split('.');
        else {
            path = [
                lng,
                ns
            ];
            if (key) {
                if (Array.isArray(key)) path.push(...key);
                else if ($1bac384020b50752$var$isString(key) && keySeparator) path.push(...key.split(keySeparator));
                else path.push(key);
            }
        }
        const result = $1bac384020b50752$var$getPath(this.data, path);
        if (!result && !ns && !key && lng.indexOf('.') > -1) {
            lng = path[0];
            ns = path[1];
            key = path.slice(2).join('.');
        }
        if (result || !ignoreJSONStructure || !$1bac384020b50752$var$isString(key)) return result;
        return $1bac384020b50752$var$deepFind(this.data?.[lng]?.[ns], key, keySeparator);
    }
    addResource(lng, ns, key, value, options = {
        silent: false
    }) {
        const keySeparator = options.keySeparator !== undefined ? options.keySeparator : this.options.keySeparator;
        let path = [
            lng,
            ns
        ];
        if (key) path = path.concat(keySeparator ? key.split(keySeparator) : key);
        if (lng.indexOf('.') > -1) {
            path = lng.split('.');
            value = ns;
            ns = path[1];
        }
        this.addNamespaces(ns);
        $1bac384020b50752$var$setPath(this.data, path, value);
        if (!options.silent) this.emit('added', lng, ns, key, value);
    }
    addResources(lng, ns, resources, options = {
        silent: false
    }) {
        for(const m in resources)if ($1bac384020b50752$var$isString(resources[m]) || Array.isArray(resources[m])) this.addResource(lng, ns, m, resources[m], {
            silent: true
        });
        if (!options.silent) this.emit('added', lng, ns, resources);
    }
    addResourceBundle(lng, ns, resources, deep, overwrite, options = {
        silent: false,
        skipCopy: false
    }) {
        let path = [
            lng,
            ns
        ];
        if (lng.indexOf('.') > -1) {
            path = lng.split('.');
            deep = resources;
            resources = ns;
            ns = path[1];
        }
        this.addNamespaces(ns);
        let pack = $1bac384020b50752$var$getPath(this.data, path) || {};
        if (!options.skipCopy) resources = JSON.parse(JSON.stringify(resources));
        if (deep) $1bac384020b50752$var$deepExtend(pack, resources, overwrite);
        else pack = {
            ...pack,
            ...resources
        };
        $1bac384020b50752$var$setPath(this.data, path, pack);
        if (!options.silent) this.emit('added', lng, ns, resources);
    }
    removeResourceBundle(lng, ns) {
        if (this.hasResourceBundle(lng, ns)) delete this.data[lng][ns];
        this.removeNamespaces(ns);
        this.emit('removed', lng, ns);
    }
    hasResourceBundle(lng, ns) {
        return this.getResource(lng, ns) !== undefined;
    }
    getResourceBundle(lng, ns) {
        if (!ns) ns = this.options.defaultNS;
        return this.getResource(lng, ns);
    }
    getDataByLanguage(lng) {
        return this.data[lng];
    }
    hasLanguageSomeTranslations(lng) {
        const data = this.getDataByLanguage(lng);
        const n = data && Object.keys(data) || [];
        return !!n.find((v)=>data[v] && Object.keys(data[v]).length > 0);
    }
    toJSON() {
        return this.data;
    }
}
var $1bac384020b50752$var$postProcessor = {
    processors: {},
    addPostProcessor (module) {
        this.processors[module.name] = module;
    },
    handle (processors, value, key, options, translator) {
        processors.forEach((processor)=>{
            value = this.processors[processor]?.process(value, key, options, translator) ?? value;
        });
        return value;
    }
};
const $1bac384020b50752$var$PATH_KEY = Symbol('i18next/PATH_KEY');
function $1bac384020b50752$var$createProxy() {
    const state = [];
    const handler = Object.create(null);
    let proxy;
    handler.get = (target, key)=>{
        proxy?.revoke?.();
        if (key === $1bac384020b50752$var$PATH_KEY) return state;
        state.push(key);
        proxy = Proxy.revocable(target, handler);
        return proxy.proxy;
    };
    return Proxy.revocable(Object.create(null), handler).proxy;
}
function $1bac384020b50752$export$ec2cccc18ab4c2ae(selector, opts) {
    const { [$1bac384020b50752$var$PATH_KEY]: path } = selector($1bac384020b50752$var$createProxy());
    return path.join(opts?.keySeparator ?? '.');
}
const $1bac384020b50752$var$checkedLoadedFor = {};
const $1bac384020b50752$var$shouldHandleAsObject = (res)=>!$1bac384020b50752$var$isString(res) && typeof res !== 'boolean' && typeof res !== 'number';
class $1bac384020b50752$var$Translator extends $1bac384020b50752$var$EventEmitter {
    constructor(services, options = {}){
        super();
        $1bac384020b50752$var$copy([
            'resourceStore',
            'languageUtils',
            'pluralResolver',
            'interpolator',
            'backendConnector',
            'i18nFormat',
            'utils'
        ], services, this);
        this.options = options;
        if (this.options.keySeparator === undefined) this.options.keySeparator = '.';
        this.logger = $1bac384020b50752$var$baseLogger.create('translator');
    }
    changeLanguage(lng) {
        if (lng) this.language = lng;
    }
    exists(key, o = {
        interpolation: {}
    }) {
        const opt = {
            ...o
        };
        if (key == null) return false;
        const resolved = this.resolve(key, opt);
        if (resolved?.res === undefined) return false;
        const isObject = $1bac384020b50752$var$shouldHandleAsObject(resolved.res);
        if (opt.returnObjects === false && isObject) return false;
        return true;
    }
    extractFromKey(key, opt) {
        let nsSeparator = opt.nsSeparator !== undefined ? opt.nsSeparator : this.options.nsSeparator;
        if (nsSeparator === undefined) nsSeparator = ':';
        const keySeparator = opt.keySeparator !== undefined ? opt.keySeparator : this.options.keySeparator;
        let namespaces = opt.ns || this.options.defaultNS || [];
        const wouldCheckForNsInKey = nsSeparator && key.indexOf(nsSeparator) > -1;
        const seemsNaturalLanguage = !this.options.userDefinedKeySeparator && !opt.keySeparator && !this.options.userDefinedNsSeparator && !opt.nsSeparator && !$1bac384020b50752$var$looksLikeObjectPath(key, nsSeparator, keySeparator);
        if (wouldCheckForNsInKey && !seemsNaturalLanguage) {
            const m = key.match(this.interpolator.nestingRegexp);
            if (m && m.length > 0) return {
                key: key,
                namespaces: $1bac384020b50752$var$isString(namespaces) ? [
                    namespaces
                ] : namespaces
            };
            const parts = key.split(nsSeparator);
            if (nsSeparator !== keySeparator || nsSeparator === keySeparator && this.options.ns.indexOf(parts[0]) > -1) namespaces = parts.shift();
            key = parts.join(keySeparator);
        }
        return {
            key: key,
            namespaces: $1bac384020b50752$var$isString(namespaces) ? [
                namespaces
            ] : namespaces
        };
    }
    translate(keys, o, lastKey) {
        let opt = typeof o === 'object' ? {
            ...o
        } : o;
        if (typeof opt !== 'object' && this.options.overloadTranslationOptionHandler) opt = this.options.overloadTranslationOptionHandler(arguments);
        if (typeof opt === 'object') opt = {
            ...opt
        };
        if (!opt) opt = {};
        if (keys == null) return '';
        if (typeof keys === 'function') keys = $1bac384020b50752$export$ec2cccc18ab4c2ae(keys, {
            ...this.options,
            ...opt
        });
        if (!Array.isArray(keys)) keys = [
            String(keys)
        ];
        const returnDetails = opt.returnDetails !== undefined ? opt.returnDetails : this.options.returnDetails;
        const keySeparator = opt.keySeparator !== undefined ? opt.keySeparator : this.options.keySeparator;
        const { key: key, namespaces: namespaces } = this.extractFromKey(keys[keys.length - 1], opt);
        const namespace = namespaces[namespaces.length - 1];
        let nsSeparator = opt.nsSeparator !== undefined ? opt.nsSeparator : this.options.nsSeparator;
        if (nsSeparator === undefined) nsSeparator = ':';
        const lng = opt.lng || this.language;
        const appendNamespaceToCIMode = opt.appendNamespaceToCIMode || this.options.appendNamespaceToCIMode;
        if (lng?.toLowerCase() === 'cimode') {
            if (appendNamespaceToCIMode) {
                if (returnDetails) return {
                    res: `${namespace}${nsSeparator}${key}`,
                    usedKey: key,
                    exactUsedKey: key,
                    usedLng: lng,
                    usedNS: namespace,
                    usedParams: this.getUsedParamsDetails(opt)
                };
                return `${namespace}${nsSeparator}${key}`;
            }
            if (returnDetails) return {
                res: key,
                usedKey: key,
                exactUsedKey: key,
                usedLng: lng,
                usedNS: namespace,
                usedParams: this.getUsedParamsDetails(opt)
            };
            return key;
        }
        const resolved = this.resolve(keys, opt);
        let res = resolved?.res;
        const resUsedKey = resolved?.usedKey || key;
        const resExactUsedKey = resolved?.exactUsedKey || key;
        const noObject = [
            '[object Number]',
            '[object Function]',
            '[object RegExp]'
        ];
        const joinArrays = opt.joinArrays !== undefined ? opt.joinArrays : this.options.joinArrays;
        const handleAsObjectInI18nFormat = !this.i18nFormat || this.i18nFormat.handleAsObject;
        const needsPluralHandling = opt.count !== undefined && !$1bac384020b50752$var$isString(opt.count);
        const hasDefaultValue = $1bac384020b50752$var$Translator.hasDefaultValue(opt);
        const defaultValueSuffix = needsPluralHandling ? this.pluralResolver.getSuffix(lng, opt.count, opt) : '';
        const defaultValueSuffixOrdinalFallback = opt.ordinal && needsPluralHandling ? this.pluralResolver.getSuffix(lng, opt.count, {
            ordinal: false
        }) : '';
        const needsZeroSuffixLookup = needsPluralHandling && !opt.ordinal && opt.count === 0;
        const defaultValue = needsZeroSuffixLookup && opt[`defaultValue${this.options.pluralSeparator}zero`] || opt[`defaultValue${defaultValueSuffix}`] || opt[`defaultValue${defaultValueSuffixOrdinalFallback}`] || opt.defaultValue;
        let resForObjHndl = res;
        if (handleAsObjectInI18nFormat && !res && hasDefaultValue) resForObjHndl = defaultValue;
        const handleAsObject = $1bac384020b50752$var$shouldHandleAsObject(resForObjHndl);
        const resType = Object.prototype.toString.apply(resForObjHndl);
        if (handleAsObjectInI18nFormat && resForObjHndl && handleAsObject && noObject.indexOf(resType) < 0 && !($1bac384020b50752$var$isString(joinArrays) && Array.isArray(resForObjHndl))) {
            if (!opt.returnObjects && !this.options.returnObjects) {
                if (!this.options.returnedObjectHandler) this.logger.warn('accessing an object - but returnObjects options is not enabled!');
                const r = this.options.returnedObjectHandler ? this.options.returnedObjectHandler(resUsedKey, resForObjHndl, {
                    ...opt,
                    ns: namespaces
                }) : `key '${key} (${this.language})' returned an object instead of string.`;
                if (returnDetails) {
                    resolved.res = r;
                    resolved.usedParams = this.getUsedParamsDetails(opt);
                    return resolved;
                }
                return r;
            }
            if (keySeparator) {
                const resTypeIsArray = Array.isArray(resForObjHndl);
                const copy = resTypeIsArray ? [] : {};
                const newKeyToUse = resTypeIsArray ? resExactUsedKey : resUsedKey;
                for(const m in resForObjHndl)if (Object.prototype.hasOwnProperty.call(resForObjHndl, m)) {
                    const deepKey = `${newKeyToUse}${keySeparator}${m}`;
                    if (hasDefaultValue && !res) copy[m] = this.translate(deepKey, {
                        ...opt,
                        defaultValue: $1bac384020b50752$var$shouldHandleAsObject(defaultValue) ? defaultValue[m] : undefined,
                        joinArrays: false,
                        ns: namespaces
                    });
                    else copy[m] = this.translate(deepKey, {
                        ...opt,
                        joinArrays: false,
                        ns: namespaces
                    });
                    if (copy[m] === deepKey) copy[m] = resForObjHndl[m];
                }
                res = copy;
            }
        } else if (handleAsObjectInI18nFormat && $1bac384020b50752$var$isString(joinArrays) && Array.isArray(res)) {
            res = res.join(joinArrays);
            if (res) res = this.extendTranslation(res, keys, opt, lastKey);
        } else {
            let usedDefault = false;
            let usedKey = false;
            if (!this.isValidLookup(res) && hasDefaultValue) {
                usedDefault = true;
                res = defaultValue;
            }
            if (!this.isValidLookup(res)) {
                usedKey = true;
                res = key;
            }
            const missingKeyNoValueFallbackToKey = opt.missingKeyNoValueFallbackToKey || this.options.missingKeyNoValueFallbackToKey;
            const resForMissing = missingKeyNoValueFallbackToKey && usedKey ? undefined : res;
            const updateMissing = hasDefaultValue && defaultValue !== res && this.options.updateMissing;
            if (usedKey || usedDefault || updateMissing) {
                this.logger.log(updateMissing ? 'updateKey' : 'missingKey', lng, namespace, key, updateMissing ? defaultValue : res);
                if (keySeparator) {
                    const fk = this.resolve(key, {
                        ...opt,
                        keySeparator: false
                    });
                    if (fk && fk.res) this.logger.warn('Seems the loaded translations were in flat JSON format instead of nested. Either set keySeparator: false on init or make sure your translations are published in nested format.');
                }
                let lngs = [];
                const fallbackLngs = this.languageUtils.getFallbackCodes(this.options.fallbackLng, opt.lng || this.language);
                if (this.options.saveMissingTo === 'fallback' && fallbackLngs && fallbackLngs[0]) for(let i = 0; i < fallbackLngs.length; i++)lngs.push(fallbackLngs[i]);
                else if (this.options.saveMissingTo === 'all') lngs = this.languageUtils.toResolveHierarchy(opt.lng || this.language);
                else lngs.push(opt.lng || this.language);
                const send = (l, k, specificDefaultValue)=>{
                    const defaultForMissing = hasDefaultValue && specificDefaultValue !== res ? specificDefaultValue : resForMissing;
                    if (this.options.missingKeyHandler) this.options.missingKeyHandler(l, namespace, k, defaultForMissing, updateMissing, opt);
                    else if (this.backendConnector?.saveMissing) this.backendConnector.saveMissing(l, namespace, k, defaultForMissing, updateMissing, opt);
                    this.emit('missingKey', l, namespace, k, res);
                };
                if (this.options.saveMissing) {
                    if (this.options.saveMissingPlurals && needsPluralHandling) lngs.forEach((language)=>{
                        const suffixes = this.pluralResolver.getSuffixes(language, opt);
                        if (needsZeroSuffixLookup && opt[`defaultValue${this.options.pluralSeparator}zero`] && suffixes.indexOf(`${this.options.pluralSeparator}zero`) < 0) suffixes.push(`${this.options.pluralSeparator}zero`);
                        suffixes.forEach((suffix)=>{
                            send([
                                language
                            ], key + suffix, opt[`defaultValue${suffix}`] || defaultValue);
                        });
                    });
                    else send(lngs, key, defaultValue);
                }
            }
            res = this.extendTranslation(res, keys, opt, resolved, lastKey);
            if (usedKey && res === key && this.options.appendNamespaceToMissingKey) res = `${namespace}${nsSeparator}${key}`;
            if ((usedKey || usedDefault) && this.options.parseMissingKeyHandler) res = this.options.parseMissingKeyHandler(this.options.appendNamespaceToMissingKey ? `${namespace}${nsSeparator}${key}` : key, usedDefault ? res : undefined, opt);
        }
        if (returnDetails) {
            resolved.res = res;
            resolved.usedParams = this.getUsedParamsDetails(opt);
            return resolved;
        }
        return res;
    }
    extendTranslation(res, key, opt, resolved, lastKey) {
        if (this.i18nFormat?.parse) res = this.i18nFormat.parse(res, {
            ...this.options.interpolation.defaultVariables,
            ...opt
        }, opt.lng || this.language || resolved.usedLng, resolved.usedNS, resolved.usedKey, {
            resolved: resolved
        });
        else if (!opt.skipInterpolation) {
            if (opt.interpolation) this.interpolator.init({
                ...opt,
                interpolation: {
                    ...this.options.interpolation,
                    ...opt.interpolation
                }
            });
            const skipOnVariables = $1bac384020b50752$var$isString(res) && (opt?.interpolation?.skipOnVariables !== undefined ? opt.interpolation.skipOnVariables : this.options.interpolation.skipOnVariables);
            let nestBef;
            if (skipOnVariables) {
                const nb = res.match(this.interpolator.nestingRegexp);
                nestBef = nb && nb.length;
            }
            let data = opt.replace && !$1bac384020b50752$var$isString(opt.replace) ? opt.replace : opt;
            if (this.options.interpolation.defaultVariables) data = {
                ...this.options.interpolation.defaultVariables,
                ...data
            };
            res = this.interpolator.interpolate(res, data, opt.lng || this.language || resolved.usedLng, opt);
            if (skipOnVariables) {
                const na = res.match(this.interpolator.nestingRegexp);
                const nestAft = na && na.length;
                if (nestBef < nestAft) opt.nest = false;
            }
            if (!opt.lng && resolved && resolved.res) opt.lng = this.language || resolved.usedLng;
            if (opt.nest !== false) res = this.interpolator.nest(res, (...args)=>{
                if (lastKey?.[0] === args[0] && !opt.context) {
                    this.logger.warn(`It seems you are nesting recursively key: ${args[0]} in key: ${key[0]}`);
                    return null;
                }
                return this.translate(...args, key);
            }, opt);
            if (opt.interpolation) this.interpolator.reset();
        }
        const postProcess = opt.postProcess || this.options.postProcess;
        const postProcessorNames = $1bac384020b50752$var$isString(postProcess) ? [
            postProcess
        ] : postProcess;
        if (res != null && postProcessorNames?.length && opt.applyPostProcessor !== false) res = $1bac384020b50752$var$postProcessor.handle(postProcessorNames, res, key, this.options && this.options.postProcessPassResolved ? {
            i18nResolved: {
                ...resolved,
                usedParams: this.getUsedParamsDetails(opt)
            },
            ...opt
        } : opt, this);
        return res;
    }
    resolve(keys, opt = {}) {
        let found;
        let usedKey;
        let exactUsedKey;
        let usedLng;
        let usedNS;
        if ($1bac384020b50752$var$isString(keys)) keys = [
            keys
        ];
        keys.forEach((k)=>{
            if (this.isValidLookup(found)) return;
            const extracted = this.extractFromKey(k, opt);
            const key = extracted.key;
            usedKey = key;
            let namespaces = extracted.namespaces;
            if (this.options.fallbackNS) namespaces = namespaces.concat(this.options.fallbackNS);
            const needsPluralHandling = opt.count !== undefined && !$1bac384020b50752$var$isString(opt.count);
            const needsZeroSuffixLookup = needsPluralHandling && !opt.ordinal && opt.count === 0;
            const needsContextHandling = opt.context !== undefined && ($1bac384020b50752$var$isString(opt.context) || typeof opt.context === 'number') && opt.context !== '';
            const codes = opt.lngs ? opt.lngs : this.languageUtils.toResolveHierarchy(opt.lng || this.language, opt.fallbackLng);
            namespaces.forEach((ns)=>{
                if (this.isValidLookup(found)) return;
                usedNS = ns;
                if (!$1bac384020b50752$var$checkedLoadedFor[`${codes[0]}-${ns}`] && this.utils?.hasLoadedNamespace && !this.utils?.hasLoadedNamespace(usedNS)) {
                    $1bac384020b50752$var$checkedLoadedFor[`${codes[0]}-${ns}`] = true;
                    this.logger.warn(`key "${usedKey}" for languages "${codes.join(', ')}" won't get resolved as namespace "${usedNS}" was not yet loaded`, 'This means something IS WRONG in your setup. You access the t function before i18next.init / i18next.loadNamespace / i18next.changeLanguage was done. Wait for the callback or Promise to resolve before accessing it!!!');
                }
                codes.forEach((code)=>{
                    if (this.isValidLookup(found)) return;
                    usedLng = code;
                    const finalKeys = [
                        key
                    ];
                    if (this.i18nFormat?.addLookupKeys) this.i18nFormat.addLookupKeys(finalKeys, key, code, ns, opt);
                    else {
                        let pluralSuffix;
                        if (needsPluralHandling) pluralSuffix = this.pluralResolver.getSuffix(code, opt.count, opt);
                        const zeroSuffix = `${this.options.pluralSeparator}zero`;
                        const ordinalPrefix = `${this.options.pluralSeparator}ordinal${this.options.pluralSeparator}`;
                        if (needsPluralHandling) {
                            if (opt.ordinal && pluralSuffix.indexOf(ordinalPrefix) === 0) finalKeys.push(key + pluralSuffix.replace(ordinalPrefix, this.options.pluralSeparator));
                            finalKeys.push(key + pluralSuffix);
                            if (needsZeroSuffixLookup) finalKeys.push(key + zeroSuffix);
                        }
                        if (needsContextHandling) {
                            const contextKey = `${key}${this.options.contextSeparator || '_'}${opt.context}`;
                            finalKeys.push(contextKey);
                            if (needsPluralHandling) {
                                if (opt.ordinal && pluralSuffix.indexOf(ordinalPrefix) === 0) finalKeys.push(contextKey + pluralSuffix.replace(ordinalPrefix, this.options.pluralSeparator));
                                finalKeys.push(contextKey + pluralSuffix);
                                if (needsZeroSuffixLookup) finalKeys.push(contextKey + zeroSuffix);
                            }
                        }
                    }
                    let possibleKey;
                    while(possibleKey = finalKeys.pop())if (!this.isValidLookup(found)) {
                        exactUsedKey = possibleKey;
                        found = this.getResource(code, ns, possibleKey, opt);
                    }
                });
            });
        });
        return {
            res: found,
            usedKey: usedKey,
            exactUsedKey: exactUsedKey,
            usedLng: usedLng,
            usedNS: usedNS
        };
    }
    isValidLookup(res) {
        return res !== undefined && !(!this.options.returnNull && res === null) && !(!this.options.returnEmptyString && res === '');
    }
    getResource(code, ns, key, options = {}) {
        if (this.i18nFormat?.getResource) return this.i18nFormat.getResource(code, ns, key, options);
        return this.resourceStore.getResource(code, ns, key, options);
    }
    getUsedParamsDetails(options = {}) {
        const optionsKeys = [
            'defaultValue',
            'ordinal',
            'context',
            'replace',
            'lng',
            'lngs',
            'fallbackLng',
            'ns',
            'keySeparator',
            'nsSeparator',
            'returnObjects',
            'returnDetails',
            'joinArrays',
            'postProcess',
            'interpolation'
        ];
        const useOptionsReplaceForData = options.replace && !$1bac384020b50752$var$isString(options.replace);
        let data = useOptionsReplaceForData ? options.replace : options;
        if (useOptionsReplaceForData && typeof options.count !== 'undefined') data.count = options.count;
        if (this.options.interpolation.defaultVariables) data = {
            ...this.options.interpolation.defaultVariables,
            ...data
        };
        if (!useOptionsReplaceForData) {
            data = {
                ...data
            };
            for (const key of optionsKeys)delete data[key];
        }
        return data;
    }
    static hasDefaultValue(options) {
        const prefix = 'defaultValue';
        for(const option in options){
            if (Object.prototype.hasOwnProperty.call(options, option) && prefix === option.substring(0, prefix.length) && undefined !== options[option]) return true;
        }
        return false;
    }
}
class $1bac384020b50752$var$LanguageUtil {
    constructor(options){
        this.options = options;
        this.supportedLngs = this.options.supportedLngs || false;
        this.logger = $1bac384020b50752$var$baseLogger.create('languageUtils');
    }
    getScriptPartFromCode(code) {
        code = $1bac384020b50752$var$getCleanedCode(code);
        if (!code || code.indexOf('-') < 0) return null;
        const p = code.split('-');
        if (p.length === 2) return null;
        p.pop();
        if (p[p.length - 1].toLowerCase() === 'x') return null;
        return this.formatLanguageCode(p.join('-'));
    }
    getLanguagePartFromCode(code) {
        code = $1bac384020b50752$var$getCleanedCode(code);
        if (!code || code.indexOf('-') < 0) return code;
        const p = code.split('-');
        return this.formatLanguageCode(p[0]);
    }
    formatLanguageCode(code) {
        if ($1bac384020b50752$var$isString(code) && code.indexOf('-') > -1) {
            let formattedCode;
            try {
                formattedCode = Intl.getCanonicalLocales(code)[0];
            } catch (e) {}
            if (formattedCode && this.options.lowerCaseLng) formattedCode = formattedCode.toLowerCase();
            if (formattedCode) return formattedCode;
            if (this.options.lowerCaseLng) return code.toLowerCase();
            return code;
        }
        return this.options.cleanCode || this.options.lowerCaseLng ? code.toLowerCase() : code;
    }
    isSupportedCode(code) {
        if (this.options.load === 'languageOnly' || this.options.nonExplicitSupportedLngs) code = this.getLanguagePartFromCode(code);
        return !this.supportedLngs || !this.supportedLngs.length || this.supportedLngs.indexOf(code) > -1;
    }
    getBestMatchFromCodes(codes) {
        if (!codes) return null;
        let found;
        codes.forEach((code)=>{
            if (found) return;
            const cleanedLng = this.formatLanguageCode(code);
            if (!this.options.supportedLngs || this.isSupportedCode(cleanedLng)) found = cleanedLng;
        });
        if (!found && this.options.supportedLngs) codes.forEach((code)=>{
            if (found) return;
            const lngScOnly = this.getScriptPartFromCode(code);
            if (this.isSupportedCode(lngScOnly)) return found = lngScOnly;
            const lngOnly = this.getLanguagePartFromCode(code);
            if (this.isSupportedCode(lngOnly)) return found = lngOnly;
            found = this.options.supportedLngs.find((supportedLng)=>{
                if (supportedLng === lngOnly) return supportedLng;
                if (supportedLng.indexOf('-') < 0 && lngOnly.indexOf('-') < 0) return;
                if (supportedLng.indexOf('-') > 0 && lngOnly.indexOf('-') < 0 && supportedLng.substring(0, supportedLng.indexOf('-')) === lngOnly) return supportedLng;
                if (supportedLng.indexOf(lngOnly) === 0 && lngOnly.length > 1) return supportedLng;
            });
        });
        if (!found) found = this.getFallbackCodes(this.options.fallbackLng)[0];
        return found;
    }
    getFallbackCodes(fallbacks, code) {
        if (!fallbacks) return [];
        if (typeof fallbacks === 'function') fallbacks = fallbacks(code);
        if ($1bac384020b50752$var$isString(fallbacks)) fallbacks = [
            fallbacks
        ];
        if (Array.isArray(fallbacks)) return fallbacks;
        if (!code) return fallbacks.default || [];
        let found = fallbacks[code];
        if (!found) found = fallbacks[this.getScriptPartFromCode(code)];
        if (!found) found = fallbacks[this.formatLanguageCode(code)];
        if (!found) found = fallbacks[this.getLanguagePartFromCode(code)];
        if (!found) found = fallbacks.default;
        return found || [];
    }
    toResolveHierarchy(code, fallbackCode) {
        const fallbackCodes = this.getFallbackCodes((fallbackCode === false ? [] : fallbackCode) || this.options.fallbackLng || [], code);
        const codes = [];
        const addCode = (c)=>{
            if (!c) return;
            if (this.isSupportedCode(c)) codes.push(c);
            else this.logger.warn(`rejecting language code not found in supportedLngs: ${c}`);
        };
        if ($1bac384020b50752$var$isString(code) && (code.indexOf('-') > -1 || code.indexOf('_') > -1)) {
            if (this.options.load !== 'languageOnly') addCode(this.formatLanguageCode(code));
            if (this.options.load !== 'languageOnly' && this.options.load !== 'currentOnly') addCode(this.getScriptPartFromCode(code));
            if (this.options.load !== 'currentOnly') addCode(this.getLanguagePartFromCode(code));
        } else if ($1bac384020b50752$var$isString(code)) addCode(this.formatLanguageCode(code));
        fallbackCodes.forEach((fc)=>{
            if (codes.indexOf(fc) < 0) addCode(this.formatLanguageCode(fc));
        });
        return codes;
    }
}
const $1bac384020b50752$var$suffixesOrder = {
    zero: 0,
    one: 1,
    two: 2,
    few: 3,
    many: 4,
    other: 5
};
const $1bac384020b50752$var$dummyRule = {
    select: (count)=>count === 1 ? 'one' : 'other',
    resolvedOptions: ()=>({
            pluralCategories: [
                'one',
                'other'
            ]
        })
};
class $1bac384020b50752$var$PluralResolver {
    constructor(languageUtils, options = {}){
        this.languageUtils = languageUtils;
        this.options = options;
        this.logger = $1bac384020b50752$var$baseLogger.create('pluralResolver');
        this.pluralRulesCache = {};
    }
    clearCache() {
        this.pluralRulesCache = {};
    }
    getRule(code, options = {}) {
        const cleanedCode = $1bac384020b50752$var$getCleanedCode(code === 'dev' ? 'en' : code);
        const type = options.ordinal ? 'ordinal' : 'cardinal';
        const cacheKey = JSON.stringify({
            cleanedCode: cleanedCode,
            type: type
        });
        if (cacheKey in this.pluralRulesCache) return this.pluralRulesCache[cacheKey];
        let rule;
        try {
            rule = new Intl.PluralRules(cleanedCode, {
                type: type
            });
        } catch (err) {
            if (!Intl) {
                this.logger.error('No Intl support, please use an Intl polyfill!');
                return $1bac384020b50752$var$dummyRule;
            }
            if (!code.match(/-|_/)) return $1bac384020b50752$var$dummyRule;
            const lngPart = this.languageUtils.getLanguagePartFromCode(code);
            rule = this.getRule(lngPart, options);
        }
        this.pluralRulesCache[cacheKey] = rule;
        return rule;
    }
    needsPlural(code, options = {}) {
        let rule = this.getRule(code, options);
        if (!rule) rule = this.getRule('dev', options);
        return rule?.resolvedOptions().pluralCategories.length > 1;
    }
    getPluralFormsOfKey(code, key, options = {}) {
        return this.getSuffixes(code, options).map((suffix)=>`${key}${suffix}`);
    }
    getSuffixes(code, options = {}) {
        let rule = this.getRule(code, options);
        if (!rule) rule = this.getRule('dev', options);
        if (!rule) return [];
        return rule.resolvedOptions().pluralCategories.sort((pluralCategory1, pluralCategory2)=>$1bac384020b50752$var$suffixesOrder[pluralCategory1] - $1bac384020b50752$var$suffixesOrder[pluralCategory2]).map((pluralCategory)=>`${this.options.prepend}${options.ordinal ? `ordinal${this.options.prepend}` : ''}${pluralCategory}`);
    }
    getSuffix(code, count, options = {}) {
        const rule = this.getRule(code, options);
        if (rule) return `${this.options.prepend}${options.ordinal ? `ordinal${this.options.prepend}` : ''}${rule.select(count)}`;
        this.logger.warn(`no plural rule found for: ${code}`);
        return this.getSuffix('dev', count, options);
    }
}
const $1bac384020b50752$var$deepFindWithDefaults = (data, defaultData, key, keySeparator = '.', ignoreJSONStructure = true)=>{
    let path = $1bac384020b50752$var$getPathWithDefaults(data, defaultData, key);
    if (!path && ignoreJSONStructure && $1bac384020b50752$var$isString(key)) {
        path = $1bac384020b50752$var$deepFind(data, key, keySeparator);
        if (path === undefined) path = $1bac384020b50752$var$deepFind(defaultData, key, keySeparator);
    }
    return path;
};
const $1bac384020b50752$var$regexSafe = (val)=>val.replace(/\$/g, '$$$$');
class $1bac384020b50752$var$Interpolator {
    constructor(options = {}){
        this.logger = $1bac384020b50752$var$baseLogger.create('interpolator');
        this.options = options;
        this.format = options?.interpolation?.format || ((value)=>value);
        this.init(options);
    }
    init(options = {}) {
        if (!options.interpolation) options.interpolation = {
            escapeValue: true
        };
        const { escape: escape$1, escapeValue: escapeValue, useRawValueToEscape: useRawValueToEscape, prefix: prefix, prefixEscaped: prefixEscaped, suffix: suffix, suffixEscaped: suffixEscaped, formatSeparator: formatSeparator, unescapeSuffix: unescapeSuffix, unescapePrefix: unescapePrefix, nestingPrefix: nestingPrefix, nestingPrefixEscaped: nestingPrefixEscaped, nestingSuffix: nestingSuffix, nestingSuffixEscaped: nestingSuffixEscaped, nestingOptionsSeparator: nestingOptionsSeparator, maxReplaces: maxReplaces, alwaysFormat: alwaysFormat } = options.interpolation;
        this.escape = escape$1 !== undefined ? escape$1 : $1bac384020b50752$var$escape;
        this.escapeValue = escapeValue !== undefined ? escapeValue : true;
        this.useRawValueToEscape = useRawValueToEscape !== undefined ? useRawValueToEscape : false;
        this.prefix = prefix ? $1bac384020b50752$var$regexEscape(prefix) : prefixEscaped || '{{';
        this.suffix = suffix ? $1bac384020b50752$var$regexEscape(suffix) : suffixEscaped || '}}';
        this.formatSeparator = formatSeparator || ',';
        this.unescapePrefix = unescapeSuffix ? '' : unescapePrefix || '-';
        this.unescapeSuffix = this.unescapePrefix ? '' : unescapeSuffix || '';
        this.nestingPrefix = nestingPrefix ? $1bac384020b50752$var$regexEscape(nestingPrefix) : nestingPrefixEscaped || $1bac384020b50752$var$regexEscape('$t(');
        this.nestingSuffix = nestingSuffix ? $1bac384020b50752$var$regexEscape(nestingSuffix) : nestingSuffixEscaped || $1bac384020b50752$var$regexEscape(')');
        this.nestingOptionsSeparator = nestingOptionsSeparator || ',';
        this.maxReplaces = maxReplaces || 1000;
        this.alwaysFormat = alwaysFormat !== undefined ? alwaysFormat : false;
        this.resetRegExp();
    }
    reset() {
        if (this.options) this.init(this.options);
    }
    resetRegExp() {
        const getOrResetRegExp = (existingRegExp, pattern)=>{
            if (existingRegExp?.source === pattern) {
                existingRegExp.lastIndex = 0;
                return existingRegExp;
            }
            return new RegExp(pattern, 'g');
        };
        this.regexp = getOrResetRegExp(this.regexp, `${this.prefix}(.+?)${this.suffix}`);
        this.regexpUnescape = getOrResetRegExp(this.regexpUnescape, `${this.prefix}${this.unescapePrefix}(.+?)${this.unescapeSuffix}${this.suffix}`);
        this.nestingRegexp = getOrResetRegExp(this.nestingRegexp, `${this.nestingPrefix}((?:[^()"']+|"[^"]*"|'[^']*'|\\((?:[^()]|"[^"]*"|'[^']*')*\\))*?)${this.nestingSuffix}`);
    }
    interpolate(str, data, lng, options) {
        let match;
        let value;
        let replaces;
        const defaultData = this.options && this.options.interpolation && this.options.interpolation.defaultVariables || {};
        const handleFormat = (key)=>{
            if (key.indexOf(this.formatSeparator) < 0) {
                const path = $1bac384020b50752$var$deepFindWithDefaults(data, defaultData, key, this.options.keySeparator, this.options.ignoreJSONStructure);
                return this.alwaysFormat ? this.format(path, undefined, lng, {
                    ...options,
                    ...data,
                    interpolationkey: key
                }) : path;
            }
            const p = key.split(this.formatSeparator);
            const k = p.shift().trim();
            const f = p.join(this.formatSeparator).trim();
            return this.format($1bac384020b50752$var$deepFindWithDefaults(data, defaultData, k, this.options.keySeparator, this.options.ignoreJSONStructure), f, lng, {
                ...options,
                ...data,
                interpolationkey: k
            });
        };
        this.resetRegExp();
        const missingInterpolationHandler = options?.missingInterpolationHandler || this.options.missingInterpolationHandler;
        const skipOnVariables = options?.interpolation?.skipOnVariables !== undefined ? options.interpolation.skipOnVariables : this.options.interpolation.skipOnVariables;
        const todos = [
            {
                regex: this.regexpUnescape,
                safeValue: (val)=>$1bac384020b50752$var$regexSafe(val)
            },
            {
                regex: this.regexp,
                safeValue: (val)=>this.escapeValue ? $1bac384020b50752$var$regexSafe(this.escape(val)) : $1bac384020b50752$var$regexSafe(val)
            }
        ];
        todos.forEach((todo)=>{
            replaces = 0;
            while(match = todo.regex.exec(str)){
                const matchedVar = match[1].trim();
                value = handleFormat(matchedVar);
                if (value === undefined) {
                    if (typeof missingInterpolationHandler === 'function') {
                        const temp = missingInterpolationHandler(str, match, options);
                        value = $1bac384020b50752$var$isString(temp) ? temp : '';
                    } else if (options && Object.prototype.hasOwnProperty.call(options, matchedVar)) value = '';
                    else if (skipOnVariables) {
                        value = match[0];
                        continue;
                    } else {
                        this.logger.warn(`missed to pass in variable ${matchedVar} for interpolating ${str}`);
                        value = '';
                    }
                } else if (!$1bac384020b50752$var$isString(value) && !this.useRawValueToEscape) value = $1bac384020b50752$var$makeString(value);
                const safeValue = todo.safeValue(value);
                str = str.replace(match[0], safeValue);
                if (skipOnVariables) {
                    todo.regex.lastIndex += value.length;
                    todo.regex.lastIndex -= match[0].length;
                } else todo.regex.lastIndex = 0;
                replaces++;
                if (replaces >= this.maxReplaces) break;
            }
        });
        return str;
    }
    nest(str, fc, options = {}) {
        let match;
        let value;
        let clonedOptions;
        const handleHasOptions = (key, inheritedOptions)=>{
            const sep = this.nestingOptionsSeparator;
            if (key.indexOf(sep) < 0) return key;
            const c = key.split(new RegExp(`${sep}[ ]*{`));
            let optionsString = `{${c[1]}`;
            key = c[0];
            optionsString = this.interpolate(optionsString, clonedOptions);
            const matchedSingleQuotes = optionsString.match(/'/g);
            const matchedDoubleQuotes = optionsString.match(/"/g);
            if ((matchedSingleQuotes?.length ?? 0) % 2 === 0 && !matchedDoubleQuotes || matchedDoubleQuotes.length % 2 !== 0) optionsString = optionsString.replace(/'/g, '"');
            try {
                clonedOptions = JSON.parse(optionsString);
                if (inheritedOptions) clonedOptions = {
                    ...inheritedOptions,
                    ...clonedOptions
                };
            } catch (e) {
                this.logger.warn(`failed parsing options string in nesting for key ${key}`, e);
                return `${key}${sep}${optionsString}`;
            }
            if (clonedOptions.defaultValue && clonedOptions.defaultValue.indexOf(this.prefix) > -1) delete clonedOptions.defaultValue;
            return key;
        };
        while(match = this.nestingRegexp.exec(str)){
            let formatters = [];
            clonedOptions = {
                ...options
            };
            clonedOptions = clonedOptions.replace && !$1bac384020b50752$var$isString(clonedOptions.replace) ? clonedOptions.replace : clonedOptions;
            clonedOptions.applyPostProcessor = false;
            delete clonedOptions.defaultValue;
            const keyEndIndex = /{.*}/.test(match[1]) ? match[1].lastIndexOf('}') + 1 : match[1].indexOf(this.formatSeparator);
            if (keyEndIndex !== -1) {
                formatters = match[1].slice(keyEndIndex).split(this.formatSeparator).map((elem)=>elem.trim()).filter(Boolean);
                match[1] = match[1].slice(0, keyEndIndex);
            }
            value = fc(handleHasOptions.call(this, match[1].trim(), clonedOptions), clonedOptions);
            if (value && match[0] === str && !$1bac384020b50752$var$isString(value)) return value;
            if (!$1bac384020b50752$var$isString(value)) value = $1bac384020b50752$var$makeString(value);
            if (!value) {
                this.logger.warn(`missed to resolve ${match[1]} for nesting ${str}`);
                value = '';
            }
            if (formatters.length) value = formatters.reduce((v, f)=>this.format(v, f, options.lng, {
                    ...options,
                    interpolationkey: match[1].trim()
                }), value.trim());
            str = str.replace(match[0], value);
            this.regexp.lastIndex = 0;
        }
        return str;
    }
}
const $1bac384020b50752$var$parseFormatStr = (formatStr)=>{
    let formatName = formatStr.toLowerCase().trim();
    const formatOptions = {};
    if (formatStr.indexOf('(') > -1) {
        const p = formatStr.split('(');
        formatName = p[0].toLowerCase().trim();
        const optStr = p[1].substring(0, p[1].length - 1);
        if (formatName === 'currency' && optStr.indexOf(':') < 0) {
            if (!formatOptions.currency) formatOptions.currency = optStr.trim();
        } else if (formatName === 'relativetime' && optStr.indexOf(':') < 0) {
            if (!formatOptions.range) formatOptions.range = optStr.trim();
        } else {
            const opts = optStr.split(';');
            opts.forEach((opt)=>{
                if (opt) {
                    const [key, ...rest] = opt.split(':');
                    const val = rest.join(':').trim().replace(/^'+|'+$/g, '');
                    const trimmedKey = key.trim();
                    if (!formatOptions[trimmedKey]) formatOptions[trimmedKey] = val;
                    if (val === 'false') formatOptions[trimmedKey] = false;
                    if (val === 'true') formatOptions[trimmedKey] = true;
                    if (!isNaN(val)) formatOptions[trimmedKey] = parseInt(val, 10);
                }
            });
        }
    }
    return {
        formatName: formatName,
        formatOptions: formatOptions
    };
};
const $1bac384020b50752$var$createCachedFormatter = (fn)=>{
    const cache = {};
    return (v, l, o)=>{
        let optForCache = o;
        if (o && o.interpolationkey && o.formatParams && o.formatParams[o.interpolationkey] && o[o.interpolationkey]) optForCache = {
            ...optForCache,
            [o.interpolationkey]: undefined
        };
        const key = l + JSON.stringify(optForCache);
        let frm = cache[key];
        if (!frm) {
            frm = fn($1bac384020b50752$var$getCleanedCode(l), o);
            cache[key] = frm;
        }
        return frm(v);
    };
};
const $1bac384020b50752$var$createNonCachedFormatter = (fn)=>(v, l, o)=>fn($1bac384020b50752$var$getCleanedCode(l), o)(v);
class $1bac384020b50752$var$Formatter {
    constructor(options = {}){
        this.logger = $1bac384020b50752$var$baseLogger.create('formatter');
        this.options = options;
        this.init(options);
    }
    init(services, options = {
        interpolation: {}
    }) {
        this.formatSeparator = options.interpolation.formatSeparator || ',';
        const cf = options.cacheInBuiltFormats ? $1bac384020b50752$var$createCachedFormatter : $1bac384020b50752$var$createNonCachedFormatter;
        this.formats = {
            number: cf((lng, opt)=>{
                const formatter = new Intl.NumberFormat(lng, {
                    ...opt
                });
                return (val)=>formatter.format(val);
            }),
            currency: cf((lng, opt)=>{
                const formatter = new Intl.NumberFormat(lng, {
                    ...opt,
                    style: 'currency'
                });
                return (val)=>formatter.format(val);
            }),
            datetime: cf((lng, opt)=>{
                const formatter = new Intl.DateTimeFormat(lng, {
                    ...opt
                });
                return (val)=>formatter.format(val);
            }),
            relativetime: cf((lng, opt)=>{
                const formatter = new Intl.RelativeTimeFormat(lng, {
                    ...opt
                });
                return (val)=>formatter.format(val, opt.range || 'day');
            }),
            list: cf((lng, opt)=>{
                const formatter = new Intl.ListFormat(lng, {
                    ...opt
                });
                return (val)=>formatter.format(val);
            })
        };
    }
    add(name, fc) {
        this.formats[name.toLowerCase().trim()] = fc;
    }
    addCached(name, fc) {
        this.formats[name.toLowerCase().trim()] = $1bac384020b50752$var$createCachedFormatter(fc);
    }
    format(value, format, lng, options = {}) {
        const formats = format.split(this.formatSeparator);
        if (formats.length > 1 && formats[0].indexOf('(') > 1 && formats[0].indexOf(')') < 0 && formats.find((f)=>f.indexOf(')') > -1)) {
            const lastIndex = formats.findIndex((f)=>f.indexOf(')') > -1);
            formats[0] = [
                formats[0],
                ...formats.splice(1, lastIndex)
            ].join(this.formatSeparator);
        }
        const result = formats.reduce((mem, f)=>{
            const { formatName: formatName, formatOptions: formatOptions } = $1bac384020b50752$var$parseFormatStr(f);
            if (this.formats[formatName]) {
                let formatted = mem;
                try {
                    const valOptions = options?.formatParams?.[options.interpolationkey] || {};
                    const l = valOptions.locale || valOptions.lng || options.locale || options.lng || lng;
                    formatted = this.formats[formatName](mem, l, {
                        ...formatOptions,
                        ...options,
                        ...valOptions
                    });
                } catch (error) {
                    this.logger.warn(error);
                }
                return formatted;
            } else this.logger.warn(`there was no format function for ${formatName}`);
            return mem;
        }, value);
        return result;
    }
}
const $1bac384020b50752$var$removePending = (q, name)=>{
    if (q.pending[name] !== undefined) {
        delete q.pending[name];
        q.pendingCount--;
    }
};
class $1bac384020b50752$var$Connector extends $1bac384020b50752$var$EventEmitter {
    constructor(backend, store, services, options = {}){
        super();
        this.backend = backend;
        this.store = store;
        this.services = services;
        this.languageUtils = services.languageUtils;
        this.options = options;
        this.logger = $1bac384020b50752$var$baseLogger.create('backendConnector');
        this.waitingReads = [];
        this.maxParallelReads = options.maxParallelReads || 10;
        this.readingCalls = 0;
        this.maxRetries = options.maxRetries >= 0 ? options.maxRetries : 5;
        this.retryTimeout = options.retryTimeout >= 1 ? options.retryTimeout : 350;
        this.state = {};
        this.queue = [];
        this.backend?.init?.(services, options.backend, options);
    }
    queueLoad(languages, namespaces, options, callback) {
        const toLoad = {};
        const pending = {};
        const toLoadLanguages = {};
        const toLoadNamespaces = {};
        languages.forEach((lng)=>{
            let hasAllNamespaces = true;
            namespaces.forEach((ns)=>{
                const name = `${lng}|${ns}`;
                if (!options.reload && this.store.hasResourceBundle(lng, ns)) this.state[name] = 2;
                else if (this.state[name] < 0) ;
                else if (this.state[name] === 1) {
                    if (pending[name] === undefined) pending[name] = true;
                } else {
                    this.state[name] = 1;
                    hasAllNamespaces = false;
                    if (pending[name] === undefined) pending[name] = true;
                    if (toLoad[name] === undefined) toLoad[name] = true;
                    if (toLoadNamespaces[ns] === undefined) toLoadNamespaces[ns] = true;
                }
            });
            if (!hasAllNamespaces) toLoadLanguages[lng] = true;
        });
        if (Object.keys(toLoad).length || Object.keys(pending).length) this.queue.push({
            pending: pending,
            pendingCount: Object.keys(pending).length,
            loaded: {},
            errors: [],
            callback: callback
        });
        return {
            toLoad: Object.keys(toLoad),
            pending: Object.keys(pending),
            toLoadLanguages: Object.keys(toLoadLanguages),
            toLoadNamespaces: Object.keys(toLoadNamespaces)
        };
    }
    loaded(name, err, data) {
        const s = name.split('|');
        const lng = s[0];
        const ns = s[1];
        if (err) this.emit('failedLoading', lng, ns, err);
        if (!err && data) this.store.addResourceBundle(lng, ns, data, undefined, undefined, {
            skipCopy: true
        });
        this.state[name] = err ? -1 : 2;
        if (err && data) this.state[name] = 0;
        const loaded = {};
        this.queue.forEach((q)=>{
            $1bac384020b50752$var$pushPath(q.loaded, [
                lng
            ], ns);
            $1bac384020b50752$var$removePending(q, name);
            if (err) q.errors.push(err);
            if (q.pendingCount === 0 && !q.done) {
                Object.keys(q.loaded).forEach((l)=>{
                    if (!loaded[l]) loaded[l] = {};
                    const loadedKeys = q.loaded[l];
                    if (loadedKeys.length) loadedKeys.forEach((n)=>{
                        if (loaded[l][n] === undefined) loaded[l][n] = true;
                    });
                });
                q.done = true;
                if (q.errors.length) q.callback(q.errors);
                else q.callback();
            }
        });
        this.emit('loaded', loaded);
        this.queue = this.queue.filter((q)=>!q.done);
    }
    read(lng, ns, fcName, tried = 0, wait = this.retryTimeout, callback) {
        if (!lng.length) return callback(null, {});
        if (this.readingCalls >= this.maxParallelReads) {
            this.waitingReads.push({
                lng: lng,
                ns: ns,
                fcName: fcName,
                tried: tried,
                wait: wait,
                callback: callback
            });
            return;
        }
        this.readingCalls++;
        const resolver = (err, data)=>{
            this.readingCalls--;
            if (this.waitingReads.length > 0) {
                const next = this.waitingReads.shift();
                this.read(next.lng, next.ns, next.fcName, next.tried, next.wait, next.callback);
            }
            if (err && data && tried < this.maxRetries) {
                setTimeout(()=>{
                    this.read.call(this, lng, ns, fcName, tried + 1, wait * 2, callback);
                }, wait);
                return;
            }
            callback(err, data);
        };
        const fc = this.backend[fcName].bind(this.backend);
        if (fc.length === 2) {
            try {
                const r = fc(lng, ns);
                if (r && typeof r.then === 'function') r.then((data)=>resolver(null, data)).catch(resolver);
                else resolver(null, r);
            } catch (err) {
                resolver(err);
            }
            return;
        }
        return fc(lng, ns, resolver);
    }
    prepareLoading(languages, namespaces, options = {}, callback) {
        if (!this.backend) {
            this.logger.warn('No backend was added via i18next.use. Will not load resources.');
            return callback && callback();
        }
        if ($1bac384020b50752$var$isString(languages)) languages = this.languageUtils.toResolveHierarchy(languages);
        if ($1bac384020b50752$var$isString(namespaces)) namespaces = [
            namespaces
        ];
        const toLoad = this.queueLoad(languages, namespaces, options, callback);
        if (!toLoad.toLoad.length) {
            if (!toLoad.pending.length) callback();
            return null;
        }
        toLoad.toLoad.forEach((name)=>{
            this.loadOne(name);
        });
    }
    load(languages, namespaces, callback) {
        this.prepareLoading(languages, namespaces, {}, callback);
    }
    reload(languages, namespaces, callback) {
        this.prepareLoading(languages, namespaces, {
            reload: true
        }, callback);
    }
    loadOne(name, prefix = '') {
        const s = name.split('|');
        const lng = s[0];
        const ns = s[1];
        this.read(lng, ns, 'read', undefined, undefined, (err, data)=>{
            if (err) this.logger.warn(`${prefix}loading namespace ${ns} for language ${lng} failed`, err);
            if (!err && data) this.logger.log(`${prefix}loaded namespace ${ns} for language ${lng}`, data);
            this.loaded(name, err, data);
        });
    }
    saveMissing(languages, namespace, key, fallbackValue, isUpdate, options = {}, clb = ()=>{}) {
        if (this.services?.utils?.hasLoadedNamespace && !this.services?.utils?.hasLoadedNamespace(namespace)) {
            this.logger.warn(`did not save key "${key}" as the namespace "${namespace}" was not yet loaded`, 'This means something IS WRONG in your setup. You access the t function before i18next.init / i18next.loadNamespace / i18next.changeLanguage was done. Wait for the callback or Promise to resolve before accessing it!!!');
            return;
        }
        if (key === undefined || key === null || key === '') return;
        if (this.backend?.create) {
            const opts = {
                ...options,
                isUpdate: isUpdate
            };
            const fc = this.backend.create.bind(this.backend);
            if (fc.length < 6) try {
                let r;
                if (fc.length === 5) r = fc(languages, namespace, key, fallbackValue, opts);
                else r = fc(languages, namespace, key, fallbackValue);
                if (r && typeof r.then === 'function') r.then((data)=>clb(null, data)).catch(clb);
                else clb(null, r);
            } catch (err) {
                clb(err);
            }
            else fc(languages, namespace, key, fallbackValue, clb, opts);
        }
        if (!languages || !languages[0]) return;
        this.store.addResource(languages[0], namespace, key, fallbackValue);
    }
}
const $1bac384020b50752$var$get = ()=>({
        debug: false,
        initAsync: true,
        ns: [
            'translation'
        ],
        defaultNS: [
            'translation'
        ],
        fallbackLng: [
            'dev'
        ],
        fallbackNS: false,
        supportedLngs: false,
        nonExplicitSupportedLngs: false,
        load: 'all',
        preload: false,
        simplifyPluralSuffix: true,
        keySeparator: '.',
        nsSeparator: ':',
        pluralSeparator: '_',
        contextSeparator: '_',
        partialBundledLanguages: false,
        saveMissing: false,
        updateMissing: false,
        saveMissingTo: 'fallback',
        saveMissingPlurals: true,
        missingKeyHandler: false,
        missingInterpolationHandler: false,
        postProcess: false,
        postProcessPassResolved: false,
        returnNull: false,
        returnEmptyString: true,
        returnObjects: false,
        joinArrays: false,
        returnedObjectHandler: false,
        parseMissingKeyHandler: false,
        appendNamespaceToMissingKey: false,
        appendNamespaceToCIMode: false,
        overloadTranslationOptionHandler: (args)=>{
            let ret = {};
            if (typeof args[1] === 'object') ret = args[1];
            if ($1bac384020b50752$var$isString(args[1])) ret.defaultValue = args[1];
            if ($1bac384020b50752$var$isString(args[2])) ret.tDescription = args[2];
            if (typeof args[2] === 'object' || typeof args[3] === 'object') {
                const options = args[3] || args[2];
                Object.keys(options).forEach((key)=>{
                    ret[key] = options[key];
                });
            }
            return ret;
        },
        interpolation: {
            escapeValue: true,
            format: (value)=>value,
            prefix: '{{',
            suffix: '}}',
            formatSeparator: ',',
            unescapePrefix: '-',
            nestingPrefix: '$t(',
            nestingSuffix: ')',
            nestingOptionsSeparator: ',',
            maxReplaces: 1000,
            skipOnVariables: true
        },
        cacheInBuiltFormats: true
    });
const $1bac384020b50752$var$transformOptions = (options)=>{
    if ($1bac384020b50752$var$isString(options.ns)) options.ns = [
        options.ns
    ];
    if ($1bac384020b50752$var$isString(options.fallbackLng)) options.fallbackLng = [
        options.fallbackLng
    ];
    if ($1bac384020b50752$var$isString(options.fallbackNS)) options.fallbackNS = [
        options.fallbackNS
    ];
    if (options.supportedLngs?.indexOf?.('cimode') < 0) options.supportedLngs = options.supportedLngs.concat([
        'cimode'
    ]);
    if (typeof options.initImmediate === 'boolean') options.initAsync = options.initImmediate;
    return options;
};
const $1bac384020b50752$var$noop = ()=>{};
const $1bac384020b50752$var$bindMemberFunctions = (inst)=>{
    const mems = Object.getOwnPropertyNames(Object.getPrototypeOf(inst));
    mems.forEach((mem)=>{
        if (typeof inst[mem] === 'function') inst[mem] = inst[mem].bind(inst);
    });
};
class $1bac384020b50752$var$I18n extends $1bac384020b50752$var$EventEmitter {
    constructor(options = {}, callback){
        super();
        this.options = $1bac384020b50752$var$transformOptions(options);
        this.services = {};
        this.logger = $1bac384020b50752$var$baseLogger;
        this.modules = {
            external: []
        };
        $1bac384020b50752$var$bindMemberFunctions(this);
        if (callback && !this.isInitialized && !options.isClone) {
            if (!this.options.initAsync) {
                this.init(options, callback);
                return this;
            }
            setTimeout(()=>{
                this.init(options, callback);
            }, 0);
        }
    }
    init(options = {}, callback) {
        this.isInitializing = true;
        if (typeof options === 'function') {
            callback = options;
            options = {};
        }
        if (options.defaultNS == null && options.ns) {
            if ($1bac384020b50752$var$isString(options.ns)) options.defaultNS = options.ns;
            else if (options.ns.indexOf('translation') < 0) options.defaultNS = options.ns[0];
        }
        const defOpts = $1bac384020b50752$var$get();
        this.options = {
            ...defOpts,
            ...this.options,
            ...$1bac384020b50752$var$transformOptions(options)
        };
        this.options.interpolation = {
            ...defOpts.interpolation,
            ...this.options.interpolation
        };
        if (options.keySeparator !== undefined) this.options.userDefinedKeySeparator = options.keySeparator;
        if (options.nsSeparator !== undefined) this.options.userDefinedNsSeparator = options.nsSeparator;
        if (typeof this.options.overloadTranslationOptionHandler !== 'function') this.options.overloadTranslationOptionHandler = defOpts.overloadTranslationOptionHandler;
        if (this.options.debug === true) {
            if (typeof console !== 'undefined') console.warn("i18next is maintained with support from locize.com \u2014 consider powering your project with managed localization (AI, CDN, integrations): https://locize.com");
        }
        const createClassOnDemand = (ClassOrObject)=>{
            if (!ClassOrObject) return null;
            if (typeof ClassOrObject === 'function') return new ClassOrObject();
            return ClassOrObject;
        };
        if (!this.options.isClone) {
            if (this.modules.logger) $1bac384020b50752$var$baseLogger.init(createClassOnDemand(this.modules.logger), this.options);
            else $1bac384020b50752$var$baseLogger.init(null, this.options);
            let formatter;
            if (this.modules.formatter) formatter = this.modules.formatter;
            else formatter = $1bac384020b50752$var$Formatter;
            const lu = new $1bac384020b50752$var$LanguageUtil(this.options);
            this.store = new $1bac384020b50752$var$ResourceStore(this.options.resources, this.options);
            const s = this.services;
            s.logger = $1bac384020b50752$var$baseLogger;
            s.resourceStore = this.store;
            s.languageUtils = lu;
            s.pluralResolver = new $1bac384020b50752$var$PluralResolver(lu, {
                prepend: this.options.pluralSeparator,
                simplifyPluralSuffix: this.options.simplifyPluralSuffix
            });
            const usingLegacyFormatFunction = this.options.interpolation.format && this.options.interpolation.format !== defOpts.interpolation.format;
            if (usingLegacyFormatFunction) this.logger.deprecate(`init: you are still using the legacy format function, please use the new approach: https://www.i18next.com/translation-function/formatting`);
            if (formatter && (!this.options.interpolation.format || this.options.interpolation.format === defOpts.interpolation.format)) {
                s.formatter = createClassOnDemand(formatter);
                if (s.formatter.init) s.formatter.init(s, this.options);
                this.options.interpolation.format = s.formatter.format.bind(s.formatter);
            }
            s.interpolator = new $1bac384020b50752$var$Interpolator(this.options);
            s.utils = {
                hasLoadedNamespace: this.hasLoadedNamespace.bind(this)
            };
            s.backendConnector = new $1bac384020b50752$var$Connector(createClassOnDemand(this.modules.backend), s.resourceStore, s, this.options);
            s.backendConnector.on('*', (event, ...args)=>{
                this.emit(event, ...args);
            });
            if (this.modules.languageDetector) {
                s.languageDetector = createClassOnDemand(this.modules.languageDetector);
                if (s.languageDetector.init) s.languageDetector.init(s, this.options.detection, this.options);
            }
            if (this.modules.i18nFormat) {
                s.i18nFormat = createClassOnDemand(this.modules.i18nFormat);
                if (s.i18nFormat.init) s.i18nFormat.init(this);
            }
            this.translator = new $1bac384020b50752$var$Translator(this.services, this.options);
            this.translator.on('*', (event, ...args)=>{
                this.emit(event, ...args);
            });
            this.modules.external.forEach((m)=>{
                if (m.init) m.init(this);
            });
        }
        this.format = this.options.interpolation.format;
        if (!callback) callback = $1bac384020b50752$var$noop;
        if (this.options.fallbackLng && !this.services.languageDetector && !this.options.lng) {
            const codes = this.services.languageUtils.getFallbackCodes(this.options.fallbackLng);
            if (codes.length > 0 && codes[0] !== 'dev') this.options.lng = codes[0];
        }
        if (!this.services.languageDetector && !this.options.lng) this.logger.warn('init: no languageDetector is used and no lng is defined');
        const storeApi = [
            'getResource',
            'hasResourceBundle',
            'getResourceBundle',
            'getDataByLanguage'
        ];
        storeApi.forEach((fcName)=>{
            this[fcName] = (...args)=>this.store[fcName](...args);
        });
        const storeApiChained = [
            'addResource',
            'addResources',
            'addResourceBundle',
            'removeResourceBundle'
        ];
        storeApiChained.forEach((fcName)=>{
            this[fcName] = (...args)=>{
                this.store[fcName](...args);
                return this;
            };
        });
        const deferred = $1bac384020b50752$var$defer();
        const load = ()=>{
            const finish = (err, t)=>{
                this.isInitializing = false;
                if (this.isInitialized && !this.initializedStoreOnce) this.logger.warn('init: i18next is already initialized. You should call init just once!');
                this.isInitialized = true;
                if (!this.options.isClone) this.logger.log('initialized', this.options);
                this.emit('initialized', this.options);
                deferred.resolve(t);
                callback(err, t);
            };
            if (this.languages && !this.isInitialized) return finish(null, this.t.bind(this));
            this.changeLanguage(this.options.lng, finish);
        };
        if (this.options.resources || !this.options.initAsync) load();
        else setTimeout(load, 0);
        return deferred;
    }
    loadResources(language, callback = $1bac384020b50752$var$noop) {
        let usedCallback = callback;
        const usedLng = $1bac384020b50752$var$isString(language) ? language : this.language;
        if (typeof language === 'function') usedCallback = language;
        if (!this.options.resources || this.options.partialBundledLanguages) {
            if (usedLng?.toLowerCase() === 'cimode' && (!this.options.preload || this.options.preload.length === 0)) return usedCallback();
            const toLoad = [];
            const append = (lng)=>{
                if (!lng) return;
                if (lng === 'cimode') return;
                const lngs = this.services.languageUtils.toResolveHierarchy(lng);
                lngs.forEach((l)=>{
                    if (l === 'cimode') return;
                    if (toLoad.indexOf(l) < 0) toLoad.push(l);
                });
            };
            if (!usedLng) {
                const fallbacks = this.services.languageUtils.getFallbackCodes(this.options.fallbackLng);
                fallbacks.forEach((l)=>append(l));
            } else append(usedLng);
            this.options.preload?.forEach?.((l)=>append(l));
            this.services.backendConnector.load(toLoad, this.options.ns, (e)=>{
                if (!e && !this.resolvedLanguage && this.language) this.setResolvedLanguage(this.language);
                usedCallback(e);
            });
        } else usedCallback(null);
    }
    reloadResources(lngs, ns, callback) {
        const deferred = $1bac384020b50752$var$defer();
        if (typeof lngs === 'function') {
            callback = lngs;
            lngs = undefined;
        }
        if (typeof ns === 'function') {
            callback = ns;
            ns = undefined;
        }
        if (!lngs) lngs = this.languages;
        if (!ns) ns = this.options.ns;
        if (!callback) callback = $1bac384020b50752$var$noop;
        this.services.backendConnector.reload(lngs, ns, (err)=>{
            deferred.resolve();
            callback(err);
        });
        return deferred;
    }
    use(module) {
        if (!module) throw new Error('You are passing an undefined module! Please check the object you are passing to i18next.use()');
        if (!module.type) throw new Error('You are passing a wrong module! Please check the object you are passing to i18next.use()');
        if (module.type === 'backend') this.modules.backend = module;
        if (module.type === 'logger' || module.log && module.warn && module.error) this.modules.logger = module;
        if (module.type === 'languageDetector') this.modules.languageDetector = module;
        if (module.type === 'i18nFormat') this.modules.i18nFormat = module;
        if (module.type === 'postProcessor') $1bac384020b50752$var$postProcessor.addPostProcessor(module);
        if (module.type === 'formatter') this.modules.formatter = module;
        if (module.type === '3rdParty') this.modules.external.push(module);
        return this;
    }
    setResolvedLanguage(l) {
        if (!l || !this.languages) return;
        if ([
            'cimode',
            'dev'
        ].indexOf(l) > -1) return;
        for(let li = 0; li < this.languages.length; li++){
            const lngInLngs = this.languages[li];
            if ([
                'cimode',
                'dev'
            ].indexOf(lngInLngs) > -1) continue;
            if (this.store.hasLanguageSomeTranslations(lngInLngs)) {
                this.resolvedLanguage = lngInLngs;
                break;
            }
        }
        if (!this.resolvedLanguage && this.languages.indexOf(l) < 0 && this.store.hasLanguageSomeTranslations(l)) {
            this.resolvedLanguage = l;
            this.languages.unshift(l);
        }
    }
    changeLanguage(lng, callback) {
        this.isLanguageChangingTo = lng;
        const deferred = $1bac384020b50752$var$defer();
        this.emit('languageChanging', lng);
        const setLngProps = (l)=>{
            this.language = l;
            this.languages = this.services.languageUtils.toResolveHierarchy(l);
            this.resolvedLanguage = undefined;
            this.setResolvedLanguage(l);
        };
        const done = (err, l)=>{
            if (l) {
                if (this.isLanguageChangingTo === lng) {
                    setLngProps(l);
                    this.translator.changeLanguage(l);
                    this.isLanguageChangingTo = undefined;
                    this.emit('languageChanged', l);
                    this.logger.log('languageChanged', l);
                }
            } else this.isLanguageChangingTo = undefined;
            deferred.resolve((...args)=>this.t(...args));
            if (callback) callback(err, (...args)=>this.t(...args));
        };
        const setLng = (lngs)=>{
            if (!lng && !lngs && this.services.languageDetector) lngs = [];
            const fl = $1bac384020b50752$var$isString(lngs) ? lngs : lngs && lngs[0];
            const l = this.store.hasLanguageSomeTranslations(fl) ? fl : this.services.languageUtils.getBestMatchFromCodes($1bac384020b50752$var$isString(lngs) ? [
                lngs
            ] : lngs);
            if (l) {
                if (!this.language) setLngProps(l);
                if (!this.translator.language) this.translator.changeLanguage(l);
                this.services.languageDetector?.cacheUserLanguage?.(l);
            }
            this.loadResources(l, (err)=>{
                done(err, l);
            });
        };
        if (!lng && this.services.languageDetector && !this.services.languageDetector.async) setLng(this.services.languageDetector.detect());
        else if (!lng && this.services.languageDetector && this.services.languageDetector.async) {
            if (this.services.languageDetector.detect.length === 0) this.services.languageDetector.detect().then(setLng);
            else this.services.languageDetector.detect(setLng);
        } else setLng(lng);
        return deferred;
    }
    getFixedT(lng, ns, keyPrefix) {
        const fixedT = (key, opts, ...rest)=>{
            let o;
            if (typeof opts !== 'object') o = this.options.overloadTranslationOptionHandler([
                key,
                opts
            ].concat(rest));
            else o = {
                ...opts
            };
            o.lng = o.lng || fixedT.lng;
            o.lngs = o.lngs || fixedT.lngs;
            o.ns = o.ns || fixedT.ns;
            if (o.keyPrefix !== '') o.keyPrefix = o.keyPrefix || keyPrefix || fixedT.keyPrefix;
            const keySeparator = this.options.keySeparator || '.';
            let resultKey;
            if (o.keyPrefix && Array.isArray(key)) resultKey = key.map((k)=>{
                if (typeof k === 'function') k = $1bac384020b50752$export$ec2cccc18ab4c2ae(k, {
                    ...this.options,
                    ...opts
                });
                return `${o.keyPrefix}${keySeparator}${k}`;
            });
            else {
                if (typeof key === 'function') key = $1bac384020b50752$export$ec2cccc18ab4c2ae(key, {
                    ...this.options,
                    ...opts
                });
                resultKey = o.keyPrefix ? `${o.keyPrefix}${keySeparator}${key}` : key;
            }
            return this.t(resultKey, o);
        };
        if ($1bac384020b50752$var$isString(lng)) fixedT.lng = lng;
        else fixedT.lngs = lng;
        fixedT.ns = ns;
        fixedT.keyPrefix = keyPrefix;
        return fixedT;
    }
    t(...args) {
        return this.translator?.translate(...args);
    }
    exists(...args) {
        return this.translator?.exists(...args);
    }
    setDefaultNamespace(ns) {
        this.options.defaultNS = ns;
    }
    hasLoadedNamespace(ns, options = {}) {
        if (!this.isInitialized) {
            this.logger.warn('hasLoadedNamespace: i18next was not initialized', this.languages);
            return false;
        }
        if (!this.languages || !this.languages.length) {
            this.logger.warn('hasLoadedNamespace: i18n.languages were undefined or empty', this.languages);
            return false;
        }
        const lng = options.lng || this.resolvedLanguage || this.languages[0];
        const fallbackLng = this.options ? this.options.fallbackLng : false;
        const lastLng = this.languages[this.languages.length - 1];
        if (lng.toLowerCase() === 'cimode') return true;
        const loadNotPending = (l, n)=>{
            const loadState = this.services.backendConnector.state[`${l}|${n}`];
            return loadState === -1 || loadState === 0 || loadState === 2;
        };
        if (options.precheck) {
            const preResult = options.precheck(this, loadNotPending);
            if (preResult !== undefined) return preResult;
        }
        if (this.hasResourceBundle(lng, ns)) return true;
        if (!this.services.backendConnector.backend || this.options.resources && !this.options.partialBundledLanguages) return true;
        if (loadNotPending(lng, ns) && (!fallbackLng || loadNotPending(lastLng, ns))) return true;
        return false;
    }
    loadNamespaces(ns, callback) {
        const deferred = $1bac384020b50752$var$defer();
        if (!this.options.ns) {
            if (callback) callback();
            return Promise.resolve();
        }
        if ($1bac384020b50752$var$isString(ns)) ns = [
            ns
        ];
        ns.forEach((n)=>{
            if (this.options.ns.indexOf(n) < 0) this.options.ns.push(n);
        });
        this.loadResources((err)=>{
            deferred.resolve();
            if (callback) callback(err);
        });
        return deferred;
    }
    loadLanguages(lngs, callback) {
        const deferred = $1bac384020b50752$var$defer();
        if ($1bac384020b50752$var$isString(lngs)) lngs = [
            lngs
        ];
        const preloaded = this.options.preload || [];
        const newLngs = lngs.filter((lng)=>preloaded.indexOf(lng) < 0 && this.services.languageUtils.isSupportedCode(lng));
        if (!newLngs.length) {
            if (callback) callback();
            return Promise.resolve();
        }
        this.options.preload = preloaded.concat(newLngs);
        this.loadResources((err)=>{
            deferred.resolve();
            if (callback) callback(err);
        });
        return deferred;
    }
    dir(lng) {
        if (!lng) lng = this.resolvedLanguage || (this.languages?.length > 0 ? this.languages[0] : this.language);
        if (!lng) return 'rtl';
        try {
            const l = new Intl.Locale(lng);
            if (l && l.getTextInfo) {
                const ti = l.getTextInfo();
                if (ti && ti.direction) return ti.direction;
            }
        } catch (e) {}
        const rtlLngs = [
            'ar',
            'shu',
            'sqr',
            'ssh',
            'xaa',
            'yhd',
            'yud',
            'aao',
            'abh',
            'abv',
            'acm',
            'acq',
            'acw',
            'acx',
            'acy',
            'adf',
            'ads',
            'aeb',
            'aec',
            'afb',
            'ajp',
            'apc',
            'apd',
            'arb',
            'arq',
            'ars',
            'ary',
            'arz',
            'auz',
            'avl',
            'ayh',
            'ayl',
            'ayn',
            'ayp',
            'bbz',
            'pga',
            'he',
            'iw',
            'ps',
            'pbt',
            'pbu',
            'pst',
            'prp',
            'prd',
            'ug',
            'ur',
            'ydd',
            'yds',
            'yih',
            'ji',
            'yi',
            'hbo',
            'men',
            'xmn',
            'fa',
            'jpr',
            'peo',
            'pes',
            'prs',
            'dv',
            'sam',
            'ckb'
        ];
        const languageUtils = this.services?.languageUtils || new $1bac384020b50752$var$LanguageUtil($1bac384020b50752$var$get());
        if (lng.toLowerCase().indexOf('-latn') > 1) return 'ltr';
        return rtlLngs.indexOf(languageUtils.getLanguagePartFromCode(lng)) > -1 || lng.toLowerCase().indexOf('-arab') > 1 ? 'rtl' : 'ltr';
    }
    static createInstance(options = {}, callback) {
        const instance = new $1bac384020b50752$var$I18n(options, callback);
        instance.createInstance = $1bac384020b50752$var$I18n.createInstance;
        return instance;
    }
    cloneInstance(options = {}, callback = $1bac384020b50752$var$noop) {
        const forkResourceStore = options.forkResourceStore;
        if (forkResourceStore) delete options.forkResourceStore;
        const mergedOptions = {
            ...this.options,
            ...options,
            isClone: true
        };
        const clone = new $1bac384020b50752$var$I18n(mergedOptions);
        if (options.debug !== undefined || options.prefix !== undefined) clone.logger = clone.logger.clone(options);
        const membersToCopy = [
            'store',
            'services',
            'language'
        ];
        membersToCopy.forEach((m)=>{
            clone[m] = this[m];
        });
        clone.services = {
            ...this.services
        };
        clone.services.utils = {
            hasLoadedNamespace: clone.hasLoadedNamespace.bind(clone)
        };
        if (forkResourceStore) {
            const clonedData = Object.keys(this.store.data).reduce((prev, l)=>{
                prev[l] = {
                    ...this.store.data[l]
                };
                prev[l] = Object.keys(prev[l]).reduce((acc, n)=>{
                    acc[n] = {
                        ...prev[l][n]
                    };
                    return acc;
                }, prev[l]);
                return prev;
            }, {});
            clone.store = new $1bac384020b50752$var$ResourceStore(clonedData, mergedOptions);
            clone.services.resourceStore = clone.store;
        }
        if (options.interpolation) {
            const defOpts = $1bac384020b50752$var$get();
            const mergedInterpolation = {
                ...defOpts.interpolation,
                ...this.options.interpolation,
                ...options.interpolation
            };
            const mergedForInterpolator = {
                ...mergedOptions,
                interpolation: mergedInterpolation
            };
            clone.services.interpolator = new $1bac384020b50752$var$Interpolator(mergedForInterpolator);
        }
        clone.translator = new $1bac384020b50752$var$Translator(clone.services, mergedOptions);
        clone.translator.on('*', (event, ...args)=>{
            clone.emit(event, ...args);
        });
        clone.init(mergedOptions, callback);
        clone.translator.options = mergedOptions;
        clone.translator.backendConnector.services.utils = {
            hasLoadedNamespace: clone.hasLoadedNamespace.bind(clone)
        };
        return clone;
    }
    toJSON() {
        return {
            options: this.options,
            store: this.store,
            language: this.language,
            languages: this.languages,
            resolvedLanguage: this.resolvedLanguage
        };
    }
}
const $1bac384020b50752$export$2e2bcd8739ae039 = $1bac384020b50752$var$I18n.createInstance();
const $1bac384020b50752$export$99152e8d49ca4e7d = $1bac384020b50752$export$2e2bcd8739ae039.createInstance;
const $1bac384020b50752$export$147ec2801e896265 = $1bac384020b50752$export$2e2bcd8739ae039.dir;
const $1bac384020b50752$export$2cd8252107eb640b = $1bac384020b50752$export$2e2bcd8739ae039.init;
const $1bac384020b50752$export$d3d08d944062d7e = $1bac384020b50752$export$2e2bcd8739ae039.loadResources;
const $1bac384020b50752$export$a5d9bf5d83fcab09 = $1bac384020b50752$export$2e2bcd8739ae039.reloadResources;
const $1bac384020b50752$export$1f96ae73734a86cc = $1bac384020b50752$export$2e2bcd8739ae039.use;
const $1bac384020b50752$export$61465194746e7fd2 = $1bac384020b50752$export$2e2bcd8739ae039.changeLanguage;
const $1bac384020b50752$export$f90d180fc7da3b3b = $1bac384020b50752$export$2e2bcd8739ae039.getFixedT;
const $1bac384020b50752$export$625550452a3fa3ec = $1bac384020b50752$export$2e2bcd8739ae039.t;
const $1bac384020b50752$export$f7e9f41ea797a17 = $1bac384020b50752$export$2e2bcd8739ae039.exists;
const $1bac384020b50752$export$2b4b218e406d2d00 = $1bac384020b50752$export$2e2bcd8739ae039.setDefaultNamespace;
const $1bac384020b50752$export$93d9ee97c1ad3f31 = $1bac384020b50752$export$2e2bcd8739ae039.hasLoadedNamespace;
const $1bac384020b50752$export$83be934b53fff43b = $1bac384020b50752$export$2e2bcd8739ae039.loadNamespaces;
const $1bac384020b50752$export$8cd7e7a54fa865bc = $1bac384020b50752$export$2e2bcd8739ae039.loadLanguages;


const { slice: $199830a05f92d3d0$var$slice, forEach: $199830a05f92d3d0$var$forEach } = [];
function $199830a05f92d3d0$var$defaults(obj) {
    $199830a05f92d3d0$var$forEach.call($199830a05f92d3d0$var$slice.call(arguments, 1), (source)=>{
        if (source) {
            for(const prop in source)if (obj[prop] === undefined) obj[prop] = source[prop];
        }
    });
    return obj;
}
function $199830a05f92d3d0$var$hasXSS(input) {
    if (typeof input !== 'string') return false;
    // Common XSS attack patterns
    const xssPatterns = [
        /<\s*script.*?>/i,
        /<\s*\/\s*script\s*>/i,
        /<\s*img.*?on\w+\s*=/i,
        /<\s*\w+\s*on\w+\s*=.*?>/i,
        /javascript\s*:/i,
        /vbscript\s*:/i,
        /expression\s*\(/i,
        /eval\s*\(/i,
        /alert\s*\(/i,
        /document\.cookie/i,
        /document\.write\s*\(/i,
        /window\.location/i,
        /innerHTML/i
    ];
    return xssPatterns.some((pattern)=>pattern.test(input));
}
// eslint-disable-next-line no-control-regex
const $199830a05f92d3d0$var$fieldContentRegExp = /^[\u0009\u0020-\u007e\u0080-\u00ff]+$/;
const $199830a05f92d3d0$var$serializeCookie = function(name, val) {
    let options = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : {
        path: '/'
    };
    const opt = options;
    const value = encodeURIComponent(val);
    let str = `${name}=${value}`;
    if (opt.maxAge > 0) {
        const maxAge = opt.maxAge - 0;
        if (Number.isNaN(maxAge)) throw new Error('maxAge should be a Number');
        str += `; Max-Age=${Math.floor(maxAge)}`;
    }
    if (opt.domain) {
        if (!$199830a05f92d3d0$var$fieldContentRegExp.test(opt.domain)) throw new TypeError('option domain is invalid');
        str += `; Domain=${opt.domain}`;
    }
    if (opt.path) {
        if (!$199830a05f92d3d0$var$fieldContentRegExp.test(opt.path)) throw new TypeError('option path is invalid');
        str += `; Path=${opt.path}`;
    }
    if (opt.expires) {
        if (typeof opt.expires.toUTCString !== 'function') throw new TypeError('option expires is invalid');
        str += `; Expires=${opt.expires.toUTCString()}`;
    }
    if (opt.httpOnly) str += '; HttpOnly';
    if (opt.secure) str += '; Secure';
    if (opt.sameSite) {
        const sameSite = typeof opt.sameSite === 'string' ? opt.sameSite.toLowerCase() : opt.sameSite;
        switch(sameSite){
            case true:
                str += '; SameSite=Strict';
                break;
            case 'lax':
                str += '; SameSite=Lax';
                break;
            case 'strict':
                str += '; SameSite=Strict';
                break;
            case 'none':
                str += '; SameSite=None';
                break;
            default:
                throw new TypeError('option sameSite is invalid');
        }
    }
    if (opt.partitioned) str += '; Partitioned';
    return str;
};
const $199830a05f92d3d0$var$cookie = {
    create (name, value, minutes, domain) {
        let cookieOptions = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : {
            path: '/',
            sameSite: 'strict'
        };
        if (minutes) {
            cookieOptions.expires = new Date();
            cookieOptions.expires.setTime(cookieOptions.expires.getTime() + minutes * 60000);
        }
        if (domain) cookieOptions.domain = domain;
        document.cookie = $199830a05f92d3d0$var$serializeCookie(name, value, cookieOptions);
    },
    read (name) {
        const nameEQ = `${name}=`;
        const ca = document.cookie.split(';');
        for(let i = 0; i < ca.length; i++){
            let c = ca[i];
            while(c.charAt(0) === ' ')c = c.substring(1, c.length);
            if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
        }
        return null;
    },
    remove (name, domain) {
        this.create(name, '', -1, domain);
    }
};
var $199830a05f92d3d0$var$cookie$1 = {
    name: 'cookie',
    // Deconstruct the options object and extract the lookupCookie property
    lookup (_ref) {
        let { lookupCookie: lookupCookie } = _ref;
        if (lookupCookie && typeof document !== 'undefined') return $199830a05f92d3d0$var$cookie.read(lookupCookie) || undefined;
        return undefined;
    },
    // Deconstruct the options object and extract the lookupCookie, cookieMinutes, cookieDomain, and cookieOptions properties
    cacheUserLanguage (lng, _ref2) {
        let { lookupCookie: lookupCookie, cookieMinutes: cookieMinutes, cookieDomain: cookieDomain, cookieOptions: cookieOptions } = _ref2;
        if (lookupCookie && typeof document !== 'undefined') $199830a05f92d3d0$var$cookie.create(lookupCookie, lng, cookieMinutes, cookieDomain, cookieOptions);
    }
};
var $199830a05f92d3d0$var$querystring = {
    name: 'querystring',
    // Deconstruct the options object and extract the lookupQuerystring property
    lookup (_ref) {
        let { lookupQuerystring: lookupQuerystring } = _ref;
        let found;
        if (typeof window !== 'undefined') {
            let { search: search } = window.location;
            if (!window.location.search && window.location.hash?.indexOf('?') > -1) search = window.location.hash.substring(window.location.hash.indexOf('?'));
            const query = search.substring(1);
            const params = query.split('&');
            for(let i = 0; i < params.length; i++){
                const pos = params[i].indexOf('=');
                if (pos > 0) {
                    const key = params[i].substring(0, pos);
                    if (key === lookupQuerystring) found = params[i].substring(pos + 1);
                }
            }
        }
        return found;
    }
};
var $199830a05f92d3d0$var$hash = {
    name: 'hash',
    // Deconstruct the options object and extract the lookupHash property and the lookupFromHashIndex property
    lookup (_ref) {
        let { lookupHash: lookupHash, lookupFromHashIndex: lookupFromHashIndex } = _ref;
        let found;
        if (typeof window !== 'undefined') {
            const { hash: hash } = window.location;
            if (hash && hash.length > 2) {
                const query = hash.substring(1);
                if (lookupHash) {
                    const params = query.split('&');
                    for(let i = 0; i < params.length; i++){
                        const pos = params[i].indexOf('=');
                        if (pos > 0) {
                            const key = params[i].substring(0, pos);
                            if (key === lookupHash) found = params[i].substring(pos + 1);
                        }
                    }
                }
                if (found) return found;
                if (!found && lookupFromHashIndex > -1) {
                    const language = hash.match(/\/([a-zA-Z-]*)/g);
                    if (!Array.isArray(language)) return undefined;
                    const index = typeof lookupFromHashIndex === 'number' ? lookupFromHashIndex : 0;
                    return language[index]?.replace('/', '');
                }
            }
        }
        return found;
    }
};
let $199830a05f92d3d0$var$hasLocalStorageSupport = null;
const $199830a05f92d3d0$var$localStorageAvailable = ()=>{
    if ($199830a05f92d3d0$var$hasLocalStorageSupport !== null) return $199830a05f92d3d0$var$hasLocalStorageSupport;
    try {
        $199830a05f92d3d0$var$hasLocalStorageSupport = typeof window !== 'undefined' && window.localStorage !== null;
        if (!$199830a05f92d3d0$var$hasLocalStorageSupport) return false;
        const testKey = 'i18next.translate.boo';
        window.localStorage.setItem(testKey, 'foo');
        window.localStorage.removeItem(testKey);
    } catch (e) {
        $199830a05f92d3d0$var$hasLocalStorageSupport = false;
    }
    return $199830a05f92d3d0$var$hasLocalStorageSupport;
};
var $199830a05f92d3d0$var$localStorage = {
    name: 'localStorage',
    // Deconstruct the options object and extract the lookupLocalStorage property
    lookup (_ref) {
        let { lookupLocalStorage: lookupLocalStorage } = _ref;
        if (lookupLocalStorage && $199830a05f92d3d0$var$localStorageAvailable()) return window.localStorage.getItem(lookupLocalStorage) || undefined; // Undefined ensures type consistency with the previous version of this function
        return undefined;
    },
    // Deconstruct the options object and extract the lookupLocalStorage property
    cacheUserLanguage (lng, _ref2) {
        let { lookupLocalStorage: lookupLocalStorage } = _ref2;
        if (lookupLocalStorage && $199830a05f92d3d0$var$localStorageAvailable()) window.localStorage.setItem(lookupLocalStorage, lng);
    }
};
let $199830a05f92d3d0$var$hasSessionStorageSupport = null;
const $199830a05f92d3d0$var$sessionStorageAvailable = ()=>{
    if ($199830a05f92d3d0$var$hasSessionStorageSupport !== null) return $199830a05f92d3d0$var$hasSessionStorageSupport;
    try {
        $199830a05f92d3d0$var$hasSessionStorageSupport = typeof window !== 'undefined' && window.sessionStorage !== null;
        if (!$199830a05f92d3d0$var$hasSessionStorageSupport) return false;
        const testKey = 'i18next.translate.boo';
        window.sessionStorage.setItem(testKey, 'foo');
        window.sessionStorage.removeItem(testKey);
    } catch (e) {
        $199830a05f92d3d0$var$hasSessionStorageSupport = false;
    }
    return $199830a05f92d3d0$var$hasSessionStorageSupport;
};
var $199830a05f92d3d0$var$sessionStorage = {
    name: 'sessionStorage',
    lookup (_ref) {
        let { lookupSessionStorage: lookupSessionStorage } = _ref;
        if (lookupSessionStorage && $199830a05f92d3d0$var$sessionStorageAvailable()) return window.sessionStorage.getItem(lookupSessionStorage) || undefined;
        return undefined;
    },
    cacheUserLanguage (lng, _ref2) {
        let { lookupSessionStorage: lookupSessionStorage } = _ref2;
        if (lookupSessionStorage && $199830a05f92d3d0$var$sessionStorageAvailable()) window.sessionStorage.setItem(lookupSessionStorage, lng);
    }
};
var $199830a05f92d3d0$var$navigator$1 = {
    name: 'navigator',
    lookup (options) {
        const found = [];
        if (typeof navigator !== 'undefined') {
            const { languages: languages, userLanguage: userLanguage, language: language } = navigator;
            if (languages) // chrome only; not an array, so can't use .push.apply instead of iterating
            for(let i = 0; i < languages.length; i++)found.push(languages[i]);
            if (userLanguage) found.push(userLanguage);
            if (language) found.push(language);
        }
        return found.length > 0 ? found : undefined;
    }
};
var $199830a05f92d3d0$var$htmlTag = {
    name: 'htmlTag',
    // Deconstruct the options object and extract the htmlTag property
    lookup (_ref) {
        let { htmlTag: htmlTag } = _ref;
        let found;
        const internalHtmlTag = htmlTag || (typeof document !== 'undefined' ? document.documentElement : null);
        if (internalHtmlTag && typeof internalHtmlTag.getAttribute === 'function') found = internalHtmlTag.getAttribute('lang');
        return found;
    }
};
var $199830a05f92d3d0$var$path = {
    name: 'path',
    // Deconstruct the options object and extract the lookupFromPathIndex property
    lookup (_ref) {
        let { lookupFromPathIndex: lookupFromPathIndex } = _ref;
        if (typeof window === 'undefined') return undefined;
        const language = window.location.pathname.match(/\/([a-zA-Z-]*)/g);
        if (!Array.isArray(language)) return undefined;
        const index = typeof lookupFromPathIndex === 'number' ? lookupFromPathIndex : 0;
        return language[index]?.replace('/', '');
    }
};
var $199830a05f92d3d0$var$subdomain = {
    name: 'subdomain',
    lookup (_ref) {
        let { lookupFromSubdomainIndex: lookupFromSubdomainIndex } = _ref;
        // If given get the subdomain index else 1
        const internalLookupFromSubdomainIndex = typeof lookupFromSubdomainIndex === 'number' ? lookupFromSubdomainIndex + 1 : 1;
        // get all matches if window.location. is existing
        // first item of match is the match itself and the second is the first group match which should be the first subdomain match
        // is the hostname no public domain get the or option of localhost
        const language = typeof window !== 'undefined' && window.location?.hostname?.match(/^(\w{2,5})\.(([a-z0-9-]{1,63}\.[a-z]{2,6})|localhost)/i);
        // if there is no match (null) return undefined
        if (!language) return undefined;
        // return the given group match
        return language[internalLookupFromSubdomainIndex];
    }
};
// some environments, throws when accessing document.cookie
let $199830a05f92d3d0$var$canCookies = false;
try {
    // eslint-disable-next-line no-unused-expressions
    document.cookie;
    $199830a05f92d3d0$var$canCookies = true;
// eslint-disable-next-line no-empty
} catch (e) {}
const $199830a05f92d3d0$var$order = [
    'querystring',
    'cookie',
    'localStorage',
    'sessionStorage',
    'navigator',
    'htmlTag'
];
if (!$199830a05f92d3d0$var$canCookies) $199830a05f92d3d0$var$order.splice(1, 1);
const $199830a05f92d3d0$var$getDefaults = ()=>({
        order: $199830a05f92d3d0$var$order,
        lookupQuerystring: 'lng',
        lookupCookie: 'i18next',
        lookupLocalStorage: 'i18nextLng',
        lookupSessionStorage: 'i18nextLng',
        // cache user language
        caches: [
            'localStorage'
        ],
        excludeCacheFor: [
            'cimode'
        ],
        // cookieMinutes: 10,
        // cookieDomain: 'myDomain'
        convertDetectedLanguage: (l)=>l
    });
class $199830a05f92d3d0$export$2e2bcd8739ae039 {
    constructor(services){
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
        this.type = 'languageDetector';
        this.detectors = {};
        this.init(services, options);
    }
    init() {
        let services = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {
            languageUtils: {}
        };
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
        let i18nOptions = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : {};
        this.services = services;
        this.options = $199830a05f92d3d0$var$defaults(options, this.options || {}, $199830a05f92d3d0$var$getDefaults());
        if (typeof this.options.convertDetectedLanguage === 'string' && this.options.convertDetectedLanguage.indexOf('15897') > -1) this.options.convertDetectedLanguage = (l)=>l.replace('-', '_');
        // backwards compatibility
        if (this.options.lookupFromUrlIndex) this.options.lookupFromPathIndex = this.options.lookupFromUrlIndex;
        this.i18nOptions = i18nOptions;
        this.addDetector($199830a05f92d3d0$var$cookie$1);
        this.addDetector($199830a05f92d3d0$var$querystring);
        this.addDetector($199830a05f92d3d0$var$localStorage);
        this.addDetector($199830a05f92d3d0$var$sessionStorage);
        this.addDetector($199830a05f92d3d0$var$navigator$1);
        this.addDetector($199830a05f92d3d0$var$htmlTag);
        this.addDetector($199830a05f92d3d0$var$path);
        this.addDetector($199830a05f92d3d0$var$subdomain);
        this.addDetector($199830a05f92d3d0$var$hash);
    }
    addDetector(detector) {
        this.detectors[detector.name] = detector;
        return this;
    }
    detect() {
        let detectionOrder = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : this.options.order;
        let detected = [];
        detectionOrder.forEach((detectorName)=>{
            if (this.detectors[detectorName]) {
                let lookup = this.detectors[detectorName].lookup(this.options);
                if (lookup && typeof lookup === 'string') lookup = [
                    lookup
                ];
                if (lookup) detected = detected.concat(lookup);
            }
        });
        detected = detected.filter((d)=>d !== undefined && d !== null && !$199830a05f92d3d0$var$hasXSS(d)).map((d)=>this.options.convertDetectedLanguage(d));
        if (this.services && this.services.languageUtils && this.services.languageUtils.getBestMatchFromCodes) return detected; // new i18next v19.5.0
        return detected.length > 0 ? detected[0] : null; // a little backward compatibility
    }
    cacheUserLanguage(lng) {
        let caches = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : this.options.caches;
        if (!caches) return;
        if (this.options.excludeCacheFor && this.options.excludeCacheFor.indexOf(lng) > -1) return;
        caches.forEach((cacheName)=>{
            if (this.detectors[cacheName]) this.detectors[cacheName].cacheUserLanguage(lng, this.options);
        });
    }
}
$199830a05f92d3d0$export$2e2bcd8739ae039.type = 'languageDetector';


var $387689dfeea53611$exports = {};
$387689dfeea53611$exports = JSON.parse('{"auth":"authentication","book_not_found":"No books found. The program will exit.","cache_file_clean":"Clear all cache files on exit","cache_file_enable":"Save web image cache files to speed up subsequent reading (consumes disk space)","cache_file_dir":"Cache folder, defaults to the system temp directory and is cleared when the program exits","can_not_init_book":"Unable to initialize the book.","cannot_listen":"Cannot start listening on the server","check_image_completed":"Image analysis completed.","check_image_error":"Resolution analysis error:","check_image_ing":"Analyzing image resolution:","check_image_start":"Starting image analysis","check_mac_error":"Error detecting MAC address:","check_port_error":"Error detecting port: %v","clear_temp_file_completed":"Temporary files cleaned up successfully:","clear_temp_file_error":"Failed to clean up temporary files:","clear_temp_file_start":"Starting to clean up temporary files.","comigo_example":"  comi book.zip\\n\\nSet the web service port (default is 1234):\\n  comi -p 2345 book.zip\\n\\nWithout opening a browser (Windows):\\n  comi -o=false book.zip\\n\\nMultiple parameters:\\n  comi -p 2345 --host example.com test.zip\\n","comigo_use":"comi","completed_extract":"Decompression completed","completed_ls":"Compressed file scan completed","config":"Specify configuration file","config_change":"Comigo configuration changes:","config_file_not_found":"Configuration file not found, using default configuration.","config_file_not_resolve":"Failed to parse the configuration.","config_save_to":"Default location to save configuration files. Options: RAM, HomeDirectory, WorkingDirectory, ProgramDirectory","ctrl_c_hint":"Press CTRL-C to quit","debug_mode":"Debug mode","disable_lan":"Disable LAN sharing","enable_database":"Enable a local database to store scanned book data","enable_file_upload":"Enable file upload","enable_frpc":"Enable frp reverse proxy","enable_webp":"Enable webp transfer, requires webp-server","epub_cannot_resort":"Cannot rearrange pages in EPUB files:","exceeds_maximum_depth":"Exceeds maximum depth, MaxDepth =","file_exit":"File already exists","file_not_found":"No suitable file was found.","flip_page_template":"Default template: flip page","format_customization_error":"This format does not support customization:","found_config_file":"Found the configuration file:","found_in_bookstore":"Found in the bookstore, skipping scan:","how_many_books_update":"Path %v updated %v books","frp_command":"frpc command or frpc executable file path (default: frpc)","frp_random_remote_port":"frp remote random port (40000~50000)","frp_remote_port":"frp remote port. If set to -1, it uses the same port as local","frp_server_addr":"frps-addr (requires frpc)","frp_server_error":"Unable to start the frpc service. Check the command format and confirm that the frpc executable is in the PATH.","frp_server_port":"frps server_port (requires frpc)","frp_setting_save_completed":"Successfully saved frpc settings.","frp_token":"token (requires frpc)","frpc_ini_error":"frpc ini initialization error","frpc_server_start":"The frpc service has started.","generate_metadata":"Generate metadata for books","get_ip_error":"Error obtaining IP:","html_title":"Comigo Reader","init_database":"Initialize database:","init_locale":"The default language is English.","local_host":"Custom domain name","local_reading":"Local reading:","log_to_file":"Enable logging to a file","long_description":"comigo, a simple comic book reader.","max_depth":"Maximum search depth","min_media_num":"Minimum number of media files required before the ZIP is considered a comic archive","no_pages_in_pdf":"No pages in PDF.","not_a_valid_zip_file":"Not a valid ZIP file:","open_browser":"Open the browser simultaneously (Windows=true)","open_browser_error":"Failed to open the browser.","open_image_error":"Error opening image:","password":"password","path_not_exist":"Path does not exist","port":"Service port","port_busy":" %v port is occupied, trying a random port","print_all_ip":"Print all available network interface IP addresses","print_config":"Print configuration","re_enter_password":"Re-enter password","reading_url_maybe":"The reading link may be: ","reg_file_hint":"On Windows, double-click this file to import and register a Comigo right-click menu.","requires_login":"Requires login","rescan":"Rescan","save":"SAVE","delete":"DELETE","save_config_file":"Saving configuration:","scan_archive_error":"File scan error","scan_error":"Scan error, directory:","scan_ing":"Scanning file","scan_pdf":"Scan PDF:","scan_start_hint":"Scan:","scroll_template":"Default template: scroll","short_description":"A simple comic book reader.","shutdown_hint":"Shutting down. Press Ctrl+C again to force.","sketch_count_seconds":"Countdown seconds in sketch mode","sketch_template":"Default template: sketch","skip_path":"Skip the path:","settings_page":"Settings Page","upload_page":"Upload Page","404notfound":"404 Not Found","sort":"Image sorting rules (none, name, time)","sort_by_name":"Sort images by file name","sort_by_time":"Sort images by modification time","start_clear_file":"- Interrupted operation, cleaning the temporary folder -","start_extract":"Starting decompression:","start_in_background":"Start Comigo in background","start_ls":"Starting to scan compressed files:","static_file_mode":"Enable static-file mode. In static mode, all images and scripts are bundled into the HTML file, so it can be saved directly as a single html file (in development,Only effective on reading pages that use scroll-mode with infinite scrolling).","stop_background":"Stop background Comigo","target_path":"Target directory","temp_folder_create_error":"Failed to create temporary folder.","temp_folder_error":"Failed to set temporary folder.","temp_folder_path":"The temporary folder is:","template":"Default page template (scroll, flip, sketch)","timeout":"timeout (minutes)","tls_crt":"TLS/SSL certificate file path","tls_enable":"Enable TLS/SSL","auto_https_cert":"Automatically Request and Issue HTTPS Certificates (Let\'s Encrypt)","port443_busy_disable_auto_tls":"Port 443 is busy, disabling Auto TLS.","auto_tls_disabled_invalid_domain":"Auto TLS requires a valid domain to function. Auto TLS has been disabled.","auto_tls_disabled_custom_cert_set":"A custom certificate is configured. Auto TLS has been disabled.","auto_tls_disabled_lan_access_off":"Auto TLS has been disabled because LAN access is turned off.","tls_key":"TLS/SSL key file path","un_archive_error":"File decompression error","unable_to_extract_images_from_pdf":"Unable to extract images from PDF.","unmarshal_config_file_error":"Failed to load the configuration file, please check the format.","unsupported_extract":"Decompression not supported:","unsupported_file_type":"Unsupported file type:","upload_disable_hint":"File upload feature has been disabled","upload_path":"Path to save uploaded files","username":"username","web_server_error":"Failed to start the web service, port:","webp_command":"webp-server command or executable file path","webp_quality":"webp compression quality","webp_server_error":"Unable to start the webp conversion service. Check the command format and confirm that the webp-server executable is in the PATH.","webp_server_start":"The webp conversion service has started.","webp_setting_error":"webp setting error","webp_setting_save_completed":"Successfully saved webp settings.","websocket_error":"websocket error:","websocket_messages":"websocket messages:","zip_encode":"Manually specify ZIP file encoding (e.g., gbk, shiftjis)","lang":"Interface language setting (auto, zh, en, ja), default is auto (auto-detect)","CacheDir":"local cache location","CacheDir_Description":"Local image cache location, default system temporary folder.","CertFile":"CertFile","CertFile_Description":"TLS/SSL certificate file path (default: , \\"~/.config/.comigo/cert.crt\\")","ClearCacheExit":"Clean up on exit","ClearCacheExit_Description":"When exiting the program, clear the web image cache.","ClearDatabaseWhenExit":"Clear database books","ClearDatabaseWhenExit_Description":"When the local database is enabled, non-existing books are purged after the scan is completed.","ConfigManager":"Profile management","ConfigManagerDeleteSuccess":"Configuration has been deleted.","ConfigManagerDescription":"Clicking Save will upload the current configuration to the server and overwrite the existing configuration file.","ConfigManagerSaveHint":"There is already a configuration file, please change the save location.","ConfigManagerSaveSuccess":"Configuration saved.","ConfigSaveTo":"Configure default save path","Debug":"Turn on Debug mode","Debug_Description":"Enable Debug to print more debugging information and check some settings related to unfinished hidden features.","DisableLAN":"Disable LAN sharing","DisableLAN_Description":"Reading services are only provided on this machine and are not shared externally. ","EnableDatabase":"Enable database","EnableDatabase_Description":"Enable local database to save scanned book data. ","EnableFrpcServer":"EnableFrpcServer","RequiresLogin":"Requires Login","RequiresLogin_Description":"Whether to enable login.","EnableTLS":"Enable TLS","EnableTLS_Description":"Whether to enable HTTPS protocol. \\nThe certificate needs to be set in the key file.","EnableUpload":"Enable upload functionality","EnableUpload_Description":"Enable upload functionality.","EnablePlugin":"Enable plugin system","EnablePlugin_Description":"Enable the plugin system to allow inserting custom HTML, CSS, and JavaScript code into pages.","plugins_config":"Plugin System","buildin_plugins":"Built-in Plugins","enabled_plugins":"Enabled Plugins","auto_flip_drag":"Drag","auto_flip_drag_panel":"Drag Panel","auto_flip_start_flip":"Start Auto Flip","auto_flip_start_scroll":"Start Auto Scroll","auto_flip_pause_flip":"Pause Auto Flip","auto_flip_pause_scroll":"Pause Auto Scroll","auto_flip_interval":"Interval:","auto_flip_seconds":"seconds","auto_scroll_distance":"Scroll Distance:","plugin_name_clock":"Clock Plugin","plugin_name_auto_flip":"Auto Flip Plugin","plugin_name_auto_scroll":"Auto Scroll Plugin","plugin_name_comigo_xyz":"Comigo.xyz Plugin","plugin_name_sample":"Sample Plugin","plugin_name_sketch_practice":"Sketch Practice Plugin","sketch_practice_start":"Start Sketch Practice","sketch_practice_pause":"Pause Sketch Practice","sketch_practice_countdown":"Countdown","ExcludePath":"Exclude path","ExcludePath_Description":"When scanning books, the names of files or folders that need to be excluded","FrpClientConfig":"FrpClient settings","GenerateBookMetadata":"Generate book metadata","GenerateMetaData":"Generate metadata","GenerateMetaData_Description":"Generate book metadata. \\nNot currently in effect.","HomeDirectory":"HomeDirectory","Host":"Host name","Host_Description":"Customize the host name displayed by the QR code. \\nThe default is the network card IP.","KeyFile":"KeyFile","KeyFile_Description":"TLS/SSL key file path (default: \\"~/.config/.comigo/key.key\\")","StoreUrls":"Library folder","StoreUrls_Description":"Library folder; supports absolute and relative paths.Relative paths are based on the current working directory.<br>Experimental support for SFTP and SMB remote paths. Format examples:<br>sftp://user:pass@192.168.1.1:22/some/path <br>smb://guest@192.168.1.1:445/some/path","LogFileName":"Log file name","LogFileName_Description":"Log file name","LogFilePath":"Log save location","LogFilePath_Description":"Log file save location","LogToFile":"Record Log to local","LogToFile_Description":"Whether to save the program log to a local file. \\nNot saved by default.","MaxScanDepth":"Maximum scan depth","MaxScanDepth_Description":"Maximum scan depth. \\nFiles exceeding the depth will not be scanned. \\nThe current execution directory is the base.","MinImageNum":"Minimum number of pictures","MinImageNum_Description":"A compressed package or folder must contain at least a few pictures to be considered a book.","OpenBrowser":"Open browser","OpenBrowser_Description":"After the scan is completed, whether to open the browser at the same time. \\nThe default is true for windows and false for other platforms.","Password":"Password","Password_Description":"The password used to log in.","ReEnterPassword":"Re-enter password","ReEnterPassword_Description":"Re-enter the password to confirm.","Port":"Port","Port_Description":"Web service port.","PrintAllPossibleQRCode":"More QR codes","ProgramDirectory":"The directory where the program is located","StartFrpClientInBackground":"Start FrpClient","SupportFileType":"Supported compressed packages","SupportFileType_Description":"When scanning a file, it is used to decide whether to skip or count it as a file suffix for book processing.","SupportMediaType":"Supported image files","SupportMediaType_Description":"Image file suffix used to count the number of images when scanning compressed packages","Timeout":"Expiration","TimeoutLimitForScan":"Scan timeout","TimeoutLimitForScan_Description":"When scanning a file, if it takes more than a few seconds, it will give up scanning the file to avoid getting stuck on an overly large file.","Timeout_Description":"Cookie expiration time after enabling login. \\nThe unit is minutes.","UseCache":"Local image cache","UseCache_Description":"Enable local image extraction cache, disabled by default.","Username":"Username","Username_Description":"The username required for the login interface.","WorkingDirectory":"current working directory","ZipFileTextEncoding":"Not UTF-8","ZipFileTextEncoding_Description":"Non-utf-8 encoded ZIP file, what encoding should be used to parse it. \\nDefault GBK.","labs":"Labs","all_page_num":"Total Pages: {0}","author":"Author: {0}","auto_crop":"Edge trimming","auto_crop_num":"Cropping threshold: ","auto_double_page":"Auto Double Page (beta)","auto_hide_toolbar":"Auto Hide Toolbar","back-to-top":"Back to Top","back_button":"Back Button","back_to_bookshelf":"Back to Bookshelf","book_shelf":"Bookshelf","child_book_hint":"{0} books in the folder","click_to_toggle":"(Click to toggle)","do_you_reset_local_settings":"Do you want to reset to the default settings?","double_page_mode":"Double Page Mode","double_page_width":"Double Page Width:","download_sample_config_file":"Download Sample Config File","download_windows_reg_file":"Download Windows Reg File","drop_to_upload":"Click or drag files to this area to upload","energy_threshold":"Energy Threshold:","epub_info":"ePub Information","exit_fullscreen":"Exit Fullscreen Mode","filesize":"File Size: {0}","flip_mode":"Flip Mode","flip_odd_even_page":"Flip Odd/Even Pages","flip_odd_even_page_hint":"Click here if double pages do not align properly.","found_read_history":"Local Reading History Found","from_interrupt":"From Last Interruption","full_screen_hint":"Fullscreen Button","fullscreen":"Fullscreen","good_job_and_bye":"Good job! Goodbye.","gray_image":"Grayscale Image","hint":"hint","hint_first_page":"You are on the first page and cannot turn forward.","hint_last_page":"You are on the last page and cannot turn backward.","hour":"hours","compress_image":"Compress Image ","infinite_dropdown":"Infinite Dropdown","interval":"Interval:","manga_mode":"Manga(Right to Left)","load_all_pages":"Load All Pages","skip_and_load_full":"Read pages not loaded (first %d pages), click the button below to load all","completed_and_load_full":"You have reached the last page. Click the button below to load the full content.","load_from_interrupt":"Load from last reading position (page XX)?","login_success_hint":"Login successful. Returning to the previous page.","logout":"Logout","margin_bottom_on_scroll_mode":"Margin Bottom:","margin_on_scroll_mode":"Page gap:","limit_width":"Limit Width:","minute":"minutes","network":"Network","no_book_found_hint":"No books found. Try uploading a file?","no_support_upload_file":"File upload functionality has been disabled by the administrator.","not_support_fullscreen":"This browser does not support fullscreen mode.","now_is":"Now:","number_of_online_books":"Number of Online Books:","original_image":"Original Image","original_pdf_link":"Original PDF Link","page":"page","page_turning_seconds":"Page Turn Interval:","scroll_fixed_pagination":"Fixed Pagination","scroll_infinite_scroll":"Infinite Scroll","pdf_hint_message":"Supports pure image PDFs. If loading is slow or errors occur, please try the following:","please_enable_upload":"Please enable server upload support.","please_enter_content":"Please enter content","qrcode_hint":"Scan to read. Click to display QR code.","raw_resolution":"Raw Resolution","re_sort_book":"Resort Books","re_sort_page":"Resort Pages","reader_settings":"Settings","reading_progress_bar":"Reading Progress Bar","refresh_page":"Refresh Page","reset_local_settings":"Reset Settings","resort_file":"Resort File","comic_mode":"Comic(Left to Right)","save_page_num":"Save Progress","scan_qrcode":"Scan QR Code:","scanned_hint":"Scanned XX books. Do you want to view them now?","scroll_mode":"Scroll Mode","second":"seconds","select_language":"Select Language","server_config":"Server Configuration","server_setting":"Comigo Server Settings","set_back_color":"Background Color:","set_interface_color":"Interface Color:","show_filename":"Show Filename","show_file_icon":"Show Icon","show_header":"Show Header","showPageNum":"Show Page Number","simplify_filename":"Simplify Filename","single_page_mode":"Single Page Mode","single_page_width":"Single Page Width:","sort_by_default":"Default Sort","sort_by_filename":"Sort by Filename (A-Z)","sort_by_filename_reverse":"Sort by Filename (Z-A)","sort_by_filesize":"Sort by Filesize (Large to Small)","sort_by_filesize_reverse":"Sort by Filesize (Small to Large)","sort_by_modify_time":"Sort by Modify Time (Newest to Oldest)","sort_by_modify_time_reverse":"Sort by Modify Time (Oldest to Newest)","sort_reverse":"(Reverse)","start_sketch_message":"The countdown has begun. Have a nice day!","start_sketch_mode":"Start Sketch","starting_from_beginning":"Start from Beginning","starting_from_beginning_hint":"Load from the first page","stop_sketch_mode":"Stop Sketch","submit":"submit","add":"Add","success_fullscreen":"Entered Fullscreen Mode","successfully_loaded_reading_progress":"Successfully Loaded Reading Progress","sync_page":"Remote Page Sync","temp_future_hint":"Temporarily put some features that are not yet finished, under development and adjustment.","test":"Test","total_is":"Total:","total_time":"Total Time:","type_or_paste_content":"Type or paste content","upload_file":"Upload File","width_use_fixed_value":"Landscape Mode Width: Fixed Value","width_use_percent":"Landscape Mode Width: Percentage","portrait_width_percent":"Portrait width (percent)","auto_align":"Auto align view","swipe_turn":"Swipe Turn","login_title":"Login to Comigo","login_subtitle":"Please enter your username and password","login_failed":"Login failed, please check your username and password","login_error_teapot":"The server does not require authentication; please access the <a class=\\"font-semibold text-blue-600\\" href=\\"/\\">homepage</a> directly.","logging_in":"Logging in...","login":"Login","other_information":"Other information","login_forgot_password_hint":"Forgot your password? Please contact the system administrator","no_pattern":"Solid Color","grid_line":"Grid Line","grid_point":"Grid Point","mosaic":"Mosaic","open_pdf_in_browser":"Open PDF in Browser","StaticFileMode":"Static File Mode","StaticFileMode_Description":"Enable static-file mode. In static mode, all images and scripts are bundled into the HTML file, so it can be saved directly as a single html file (in development).","confirm_logout":"Are you sure you want to log out?","confirm_reset_settings":"Are you sure you want to reset local settings?","current_dir_scope":"When run in the current directory (local scope)","current_user_scope":"Effective for logged-in user (global scope)","portable_binary_scope":"Effective for this binary only (portable mode)","saveSuccessHint":"Settings saved successfully. The page will refresh automatically in 2 seconds.","port_change_hint":"Port has been changed. Redirecting to the new port in 2 seconds.","no_books_library_path_notice":"No readable books were found. Please configure a library path. The page will automatically refresh once the configuration is complete.","download_raw_archive":"AS File","download_portable_web_file":"AS HTML","download_as_zip":"AS ZIP","download_as_epub":"AS EPUB","clear_cache":"Clear Cache","confirm_clear_cache":"Are you sure you want to clear this book\'s cache files? (metadata and cover cache)","clear_cache_success":"Cache cleared","clear_cache_failed":"Failed to clear cache","download":"Download","switch":"Switch","EnableTailscale":"Enable Tailscale","EnableTailscale_Description":"Enable the Tailscale intranet penetration feature. The first time you enable it, verification is required in the Tailscale admin console.","TailscaleHostname":"Tailscale Hostname","TailscaleHostname_Description":"Tailscale hostname part. The full domain looks like {hostname}.example.ts.net","TailscalePort":"Tailscale Listening Port","TailscalePort_Description":"Tailscale listening port. Default is 443, TLS is enabled automatically.","FunnelTunnel":"Enable Funnel Tunnel (Public Access)","FunnelTunnel_Description":"Funnel Tunnel (public access) . If you don\u2019t want to publish publicly, it is recommended to set password protection. Funnel Tunnel can only use ports 443, 8443, and 10000.","read_only_mode":"Read-Only Mode","read_only_mode_description":"You are currently in read-only mode.You cannot change settings or upload files via the web interface.","enable_single_instance":"Enable single instance mode to ensure only one program instance runs at a time","tailscale_auth_url_is":"To start Tailscale service, restart with TS_AUTHKEY set, or go to: ","tailscale_server_start":"Tailscale service has started.","tailscale_reading_url":"Tailscale reading link: ","tailscale_not_connected_hint":"Tailscale is not connected. Please check your network or Tailscale settings.","tailscale_not_enabled":"Tailscale is not enabled. Please enable Tailscale first.","ServerSettings":"Comigo Server Settings","settings_stores":"Library Settings","settings_network":"Network Settings","settings_extra":"Experimental Features","settings_log_title":"Live Server Logs","settings_log_broadcast_placeholder":"Enter broadcast message","settings_log_broadcast_send":"Send","settings_log_sse_connected":"log server connected","settings_log_sse_closed":"closed","settings_log_sse_retrying":"retrying...","settings_log_broadcast_default_message":"Hello ~","remote_access":"Remote Access","ErrPasswordMismatch":"The password entered twice does not match, please re-enter it","PromptSetPassword":"Please set password","MsgLoginSettingsUpdated":"Login settings updated successfully. Redirecting to the login page...","CurrentPassword":"Current Password","AdminAccountSetup":"Administrator account and password","AdminAccountSetupDescription":"Please set the administrator account and password for logging into Comigo. After setup, accessing the service will require login.","ConfigStorageLocationPrompt":"Select the location to store the config file:","set_account_password":"Set account password","connect_tailscale":"Connect to Tailscale network","disconnect_tailscale":"Disconnect from Tailscale","tailscale_status":"Tailscale Status","service_status":"Service Status","running":"Running","client_count":"Client Count","host_system":"Host System","connection_status":"Connection Status","connected":"Connected","not_connected":"Not Connected","enable_funnel":"Enable Funnel","enable_funnel_public_access":"Enable Funnel\uFF08Public Access\uFF09","disabled":"Disabled","ip_address":"IP Address","service_version":"Service Version","read_link":"Read Link","submiting":"Submitting...","enable":"Enable","disable":"Disable","TailscaleAuthKey":"Tailscale Auto Auth Key","TailscaleAuthKeyDescription":"Tailscale automatic authentication key (TS_AUTHKEY), used for authentication in environments without a browser.","verify_tailscale":"Please click the link to verify Tailscale:","funnel_status":"Funnel Status","funnel_tunnel":"Funnel Tunnel","funnel_setup_done":"Funnel setup completed","funnel_setup_not_done":"Funnel needs to be configured","funnel_not_set_hint":"To use Funnel permissions, you need to:","funnel_require_dns_1":"In DNS panel","funnel_require_dns_2":"enable MagicDNS and HTTPS.","funnel_require_acl_1":"In ACL panel","funnel_require_acl_2":"edit ACL rules to enable the Funnel tunnel","funnel_require_acl_3":"(download the sample JSON file).","funnel_require_password_1":"When enabling Funnel Login Check, you must set an Comigo account and password to use the Funnel tunnel.","verify_link":"Verify link","copy_link":"Copy link","copy_success":"Copied to clipboard","copy_failed":"Copy failed, please copy manually","FunnelLoginCheck":"Funnel Login Check","FunnelLoginCheckDescription":"When enabling the Funnel tunnel, check whether password  is currently set.","funnel_login_check_enabled_but_no_password":"Funnel Login Check is enabled, but no login password is set. The Funnel tunnel cannot be activated.","tailscale_settings_submitted_check_status":"Tailscale settings submitted. Please check Tailscale status.","value_already_exists_do_not_add_again":"The value already exists. Please do not add it again.","file_uploaded_successfully":"File uploaded successfully.","content_empty_please_enter_before_submit":"Content is empty. Please enter content before submitting.","default_prompt_message":"Default Prompt Message","confirm":"Confirm","ok":"OK","cancel":"Cancel","uploading":"Uploading...","upload_failed_network_error":"Upload failed: Network error","drag_or_click_to_upload":"Drag and drop files here or click to select files","selected_file":"Selected File","select_upload_target_store":"Select Upload Target Store","store_not_exists":"Store path does not exist","no_available_stores":"No available stores, please add store paths in settings first","store_validation_failed":"Invalid store path","passwords_not_match":"The two passwords entered do not match.","please_delete_other_config_first":"Please delete configuration files from other locations first.","save_config_success":"Configuration file saved successfully!","save_config_failed":"Failed to save configuration file.","no_config_file_to_delete_in_path":"There is no configuration file to delete in the selected path.","delete_config_success":"Configuration file deleted successfully.","delete_config_failed":"Failed to delete configuration file.","err_invalid_number":"Please enter a valid number","err_number_range":"Value must be between {0} and {1}","err_update_config_failed":"Failed to update configuration","err_network_error":"Network error, please try again","err_update_login_settings_failed":"Failed to update login settings","err_add_config_failed":"Failed to add configuration","err_delete_config_failed":"Failed to delete configuration","err_save_config_failed":"Failed to save configuration","err_config_locked":"Configuration is locked and cannot be modified","read":"Read","unread":"Unread","mark_all_as_read":"Mark All as Read","new_version_available":"A New Version is Available!","download_now":"Download Now","infinite_scroll_button":"Infinite","paged_scroll_button":"Paged","page_flip_button":"Flip","infinite_scroll":"Infinite Scroll (Load All)","paged_scroll":"Paged Scroll (32 Pages)","page_flip":"Flip Mode (Swipe)","reading_progress_percent":"Progress(%)","reading_progress_page":"Progress(num)","reading_history":"History","view_all_reading_history":"View All Reading History","loading":"Loading...","no_reading_history":"No reading history","unknown":"Unknown","systray_tooltip":"Comigo Comic Reader","systray_open_browser":"Open Browser","systray_open_browser_tooltip":"Open Comigo in browser","systray_copy_url":"Copy Reading URL","systray_copy_url_tooltip":"Copy reading URL to clipboard","systray_enable_tailscale":"Enable Tailscale","systray_disable_tailscale":"Disable Tailscale","systray_toggle_tailscale_tooltip":"Toggle Tailscale status","systray_language":"Language","systray_language_tooltip":"Switch interface language","systray_language_zh":"\u4E2D\u6587","systray_language_zh_tooltip":"Switch to Chinese","systray_language_en":"English","systray_language_en_tooltip":"Switch to English","systray_language_ja":"\u65E5\u672C\u8A9E","systray_language_ja_tooltip":"\u65E5\u672C\u8A9E\u306B\u5207\u308A\u66FF\u3048","systray_open_directory":"Open Directory","systray_open_directory_tooltip":"Open related directories","systray_config_directory":"Config Directory","systray_config_directory_tooltip":"Open config directory","systray_store":"Library","systray_store_folder":"Library Folder","systray_refresh_directories":"Refresh","systray_refresh_directories_tooltip":"Refresh directory list","systray_quit":"Quit","systray_quit_tooltip":"Quit Comigo","context_menu_open_with_comigo":"Open with Comigo","register_context_menu":"Register context menu: Open with Comigo","unregister_context_menu":"Unregister context menu: Open with Comigo","systray_extra":"Extra","systray_extra_tooltip":"More integration options","systray_project":"Comigo Project Page","systray_project_tooltip":"Open Comigo GitHub repository","register_folder_context_menu":"Register folder context menu: Open with Comigo","unregister_folder_context_menu":"Unregister folder context menu: Open with Comigo","register_file_association":"Register archive file association (as a candidate)","unregister_file_association":"Unregister archive file association","create_desktop_shortcut":"Create desktop shortcut","log_language_changed_to_chinese":"Language changed to Chinese","log_language_changed_to_english":"Language changed to English","log_language_changed_to_japanese":"Language changed to Japanese","log_requesting_quit_from_systray":"Requesting quit from system tray","log_failed_to_read_icon_file":"Failed to read icon file: %v, using default icon","log_opening_browser":"Opening browser: %s","log_failed_to_copy_url":"Failed to copy URL to clipboard: %v","log_copied_url_to_clipboard":"Copied URL to clipboard: %s","log_failed_to_toggle_tailscale":"Failed to toggle Tailscale: %v","log_failed_to_set_language":"Failed to set language: %v","log_failed_to_get_config_dir":"Failed to get config dir: %v","log_failed_to_clear_folder_context_menu":"Failed to clear Windows folder context menu: %v","log_failed_to_register_folder_context_menu":"Failed to register Windows folder context menu: %v","log_failed_to_create_desktop_shortcut":"Failed to create desktop shortcut: %v","log_failed_to_unregister_archive_handler":"Failed to unregister archive handler: %v","log_failed_to_register_archive_handler":"Failed to register archive handler: %v","log_opening_comigo_project_page":"Opening Comigo project page: https://github.com/yumenaka/comigo","log_failed_to_open_directory":"Failed to open directory: %v","log_upload_file_count":"Upload file count: %d","log_upload_path_not_set":"Upload path not set","log_mkdir_failed":"mkdir failed: %s","log_mkdir_upload_folder_success":"mkdir upload folder success: %s","log_file_upload_success":"File uploaded successfully: %s","log_no_changes_skipped_scan":"No changes in cfg, skipped scan store path","log_no_changes_skipped_rescan":"No changes in cfg, skipped rescan dir","log_tailscale_config_changed_restart":"Tailscale config changed, will restart Tailscale server","log_tailscale_enabled_start":"Tailscale enabled, will start Tailscale server","log_tailscale_disabled_stop":"Tailscale disabled, will stop Tailscale server","log_scan_subdirectory_error":"Error scanning subdirectory: %v","log_get_file_info_failed":"Failed to get file info: %v","log_tailscale_disabled_skip_qrcode":"Tailscale is disabled, skipping ShowQRCodeTailscale function.","log_tailscale_status_check_exceeded":"Tailscale status check exceeded, stopping further checks.","log_tailscale_not_yet_fqdn":"Tailscale FQDN not yet available","log_flip_mode_book_id":"Flip Mode Book ID: %s","log_invalid_upload_path":"Invalid upload path.","log_time_elapsed":"Time elapsed: %v","log_waiting_for_api_health":"Waiting for API health endpoint...","log_api_healthy_ready":"Comigo API is healthy and ready!","log_path_error":"path error","log_checking_cfg_sharename":"Checking cfg ShareName","log_warning_failed_to_get_homedir":"Warning: failed to get HomeDir: %v","log_warning_failed_to_get_executable_path":"Warning: failed to get Executable path: %v","log_child_daemon_started":"child daemon started?","log_daemon_started":"daemon started","log_update_config":"Update config: %s = %s","log_failed_to_update_local_config":"Failed to update local config: %v","log_update_user_info_username":"Update user info: Username=%s","log_update_user_info_current_password":"Update user info: CurrentPassword=%s","log_update_user_info_password":"Update user info: Password=%s","log_update_user_info_reenter_password":"Update user info: ReEnterPassword=%s","log_add_array_config_handler":"AddArrayConfigHandler: %s = %s","log_delete_array_config_handler":"DeleteArrayConfigHandler: %s = %s","log_received_json_data":"Received JSON data: %s","log_failed_to_parse_json":"Failed to parse JSON data","log_failed_to_scan_store_path":"Failed to scan store path: %v","log_failed_to_save_results_to_database":"Failed to save results to database: %v","log_error_failed_save_to_directory":"error: Failed save to %s directory","log_error_find_config_in":"error: Find config in %s %s","log_error_failed_to_delete_config":"error: Failed to delete config in %s directory","log_store_url_already_exists":"Store Url already exists: %s","log_store_url_already_exists_in_config":"Store Url already exists in config: %s","log_failed_to_add_store_url":"Failed to add store url from config: %s","log_string_already_exists":"AddStringArrayConfig: string \'%s\' already exists","log_failed_to_unmarshal_json":"Failed to unmarshal JSON: %v","log_unknown_config_key":"Unknown config key: %s","log_error_getting_executable_path":"Error getting executable path: %s","log_executable_name":"Executable Name: %s","log_executable_path":"Executable Path: %s","log_failed_to_unregister_windows_context_menu":"Failed to unregister Windows context menu: %v","log_failed_to_register_windows_context_menu":"Failed to register Windows context menu: %v","log_received_new_args_from_instance":"Received new args from another instance: %v","log_single_instance_check_failed":"Single instance check failed: %v","log_args_sent_to_existing_instance":"Args sent to existing instance, exiting...","log_error_listing_books":"Error listing books: %s","log_error_clearing_temp_files":"Error clearing temp files: %s","log_cleared_temp_files":"Cleared temp files: %s","log_failed_to_get_homedirectory":"Failed to get HomeDirectory: %s","log_program_directory":"ProgramDirectory: %s","log_failed_to_get_program_directory":"Failed to get ProgramDirectory: %v","log_failed_to_get_working_directory":"Failed to get WorkingDirectory: %s","log_failed_to_open_file":"Failed to open file: %s, error: %v","log_failed_to_get_file_info":"Failed to get file info: %s, error: %v","log_add_book_error":"AddBook_error bookID: %s %s","log_get_child_books_for_bookid":"Get child books for bookID %s","log_get_child_books_count":"bookID %v, Get %v child books","log_get_bookshelf_error":"GetBookShelf Error: %v","log_error_getting_absolute_path":"Error getting absolute path: %v","log_cannot_shorten_id":"Cannot shorten ID: %s","log_error_getting_book_group":"Error getting book group: %s","log_warning_failed_to_set_socket_permissions":"Warning: failed to set socket permissions: %v","log_single_instance_server_started":"Single instance server started on: %s","log_failed_to_accept_connection":"Failed to accept connection: %v","log_failed_to_decode_message":"Failed to decode message: %v","log_failed_to_handle_new_args":"Failed to handle new args: %v","log_received_and_processed_new_args":"Received and processed new args: %v","log_failed_to_read_response":"Failed to read response, but message may have been sent: %v","log_successfully_sent_args":"Successfully sent args to existing instance: %v","log_error_closing_listener":"Error closing listener: %v","log_another_instance_running":"Another instance is already running, sending args to it...","log_get_book_error":"GetBook: %v","log_saving_books_meta_data_to":"Saving books metadata to %s","log_error_saving_book":"Error saving book %s: %s","log_successfully_saved_books":"Successfully saved %d books to %s","log_cfg_save_to":"cfg Save To %s","log_try_delete_cfg_in":"Try delete cfg in %s","log_websocket_server_received":"websocket server received: %v","log_failed_to_create_temp_config_dir":"Failed to create temp config dir: %s","log_failed_to_create_config_dir":"Failed to create config dir: %s","log_username_or_password_empty":"Username or password is empty. Using default Jwt Signing key.","log_invalid_port_number":"Invalid port number. Using default port: %d","log_tailscale_status_not_available":"Tailscale status not available yet: %v","log_server_action":"Server action: %v","log_server_action_string":"Server action: %s","log_loading_books_from":"Loading books from: %s","log_configured_store_urls":"Configured store URLs: %v","log_book_data_directory_not_exist":"Book data directory does not exist yet: %s","log_error_accessing_book_data_directory":"Error accessing book_data directory: %s","log_error_reading_book_data_directory":"Error reading book_data directory: %s","log_skipping_directory":"Skipping directory %s","log_skipping_non_json_file":"Skipping non-JSON file %s","log_error_reading_file":"Error reading file %s: %s","log_warning_corrupted_json_file":"Warning: corrupted JSON file %s, skipping: %s","log_error_deleting_corrupted_file":"Error deleting corrupted file %s: %s","log_error_adding_book_to_store":"Error adding book %s to store: %s","log_loaded_books_so_far":"Loaded %d books so far from %s","log_successfully_loaded_books":"Successfully loaded %d books from %s","log_error_adding_subfolder":"Error adding subfolder: %s","log_error_saving_book_to_json":"Error saving book %s to JSON: %s","log_error_adding_book":"Error adding book %s: %s","log_error_deleting_book_json_file":"Error deleting book %s JSON file: %s","err_add_bookstore_key_exists":"add Bookstore Error: The key already exists [%s]","err_add_book_empty_bookid":"add book error: empty BookID","err_add_bookstore_key_not_found":"add Bookstore Error: The key not found [%s]","err_cannot_find_group":"cannot find group, id=%s","err_getbook_cannot_find":"GetBook: cannot find book, id=%s","err_storebookmark_cannot_find":"StoreBookMark: cannot find book, id=%s","err_storebookmark_unknown_type":"StoreBookMark: unknown bookmark type","err_getbookmark_cannot_find":"GetBookMark: cannot find book, id=%s","err_deletebook_cannot_find":"DeleteBook: cannot find book, id=%s","err_cannot_find_book_topofshelf":"error: cannot find book in TopOfShelfInfo","err_cannot_find_child_books":"cannot find child books info\uFF0CBookID\uFF1A%s","err_getparentbook_cannot_find":"GetParentBook: cannot find book by childID=%s","err_cannot_find_book_parentfolder":"cannot find book, parentFolder=%s","log_book_data_already_exists":"Book data already exists: %s  %s","err_slice_not_supported":"This slice setting is not supported (only []string is supported)","err_must_be_nonempty_config_pointer":"Must be a non-empty *Config pointer","err_non_htmx_request":"non-htmx request","err_no_form_data":"no form data","err_failed_to_find_executable_path":"error: failed to find executable path","err_failed_to_find_home_directory":"error: failed to find home directory","err_funnel_mode_ports_only":"funnel mode only supports ports 443, 8443, and 10000","err_number_not_found":"number not found","err_file_not_zip_archive":"file is not a ZIP archive","err_file_not_rar_archive":"file is not a RAR archive","err_scan_file_error":"scan file error","err_countpages_pdf_invalid":"CountPagesOfPDF: invalid PDF: %s %s","err_name_in_archive_empty":"nameInArchive is empty","err_unsupported_archive_format":"unsupported archive format or file not found in archive","err_file_not_found_in_archive":"file not found in archive","err_extract_path_not_found":"extractPath not found in context","err_needfile_empty":"needFile is empty","err_getdata_from_epub_error":"getDataFromEpub Error. epubPath:%s  needFile:%s","err_container_xml_empty":"container.xml content is empty","err_no_valid_opf_path":"no valid OPF path found in container.xml","err_content_type_not_found":"contentType not found in cache","err_file_does_not_exist":"File does not exist:%s","err_charset_not_found":"charset not found","err_imaging_decode_error":"imaging.Decode() Error","err_imaging_encode_error":"imaging.Encode() Error","err_jpeg_encode_error":"digestImage jpeg.Encode() Error","err_imageresize_maxwidth_error":"ImageResizeByMaxWidth Error maxWidth(%d) > sourceWidth(%d)","err_imageresize_maxheight_error":"ImageResizeByMaxHeight Error maxHeight(%d) > sourceHeight(%d)","log_starting_server_on_port":"Starting Server...on port %d ...","log_auto_tls_enabled_for_domain":"Auto TLS enabled for domain: %s","log_custom_tls_cert":"Custom TLS Cert CertFile: %s KeyFile: %s","log_server_shutdown_successfully":"Server Shutdown Successfully, Starting Server...on port %d ...","log_tailscale_server_initialized":"Tailscale server initialized successfully on %s:%d","log_working_directory":"Working directory: %s","log_failed_to_add_working_directory_to_store_urls":"Failed to add working directory to store urls: %s","log_args_index":"args[%d]: %s","log_failed_to_add_store_url_from_args":"Failed to add store url from args: %s","log_failed_to_add_upload_path_to_store_urls":"Failed to add upload path to store urls: %s","log_starting_tailscale_http_server":"Starting Tailscale HTTP server on %s:%d","log_tailscale_server_stopped_successfully":"Tailscale server stopped successfully.","log_checking_book_files_exist":"Checking book files exist...","log_checking_store_exist":"Checking store exist...","log_error_deleting_book":"Error deleting book %s: %s","log_error_initializing_main_folder":"Error initializing main folder: %s","log_failed_to_get_absolute_path_scan":"Failed to get absolute path: %s","log_skip_to_scan_root_directory":"Skip scanning root directory: %s, %v","log_failed_to_get_relative_path":"Failed to get relative path: %s","log_skip_to_scan_directory":"Skip directory: %s, %v","log_failed_to_open_file_get_single":"Failed to open file %s: %v","log_timeout_identify_archive_format":"Operation timeout: identifying archive format took more than 30 seconds","log_failed_to_identify_archive_format":"Failed to identify archive format: %v","log_timeout_create_filesystem":"Operation timeout: creating filesystem took more than 30 seconds","log_failed_to_create_filesystem":"Failed to create filesystem: %v","log_timeout_read_file_content":"Operation timeout: reading file content took more than 30 seconds","log_failed_to_read_file_content":"Failed to read file content: %v","log_timeout_open_file_in_archive":"Operation timeout: opening file in archive took more than 30 seconds","log_failed_to_open_file_in_archive":"Failed to open file in archive: %v","log_timeout_extract_file":"Operation timeout: extracting file took more than 30 seconds","log_failed_to_extract_file":"Failed to extract file: %v","err_tailscale_netlistener_nil":"Tailscale netListener is nil; server will not start","err_tailscale_http_server_error":"Tailscale HTTP server error: %v","err_error_closing_network_listener":"Error closing network listener: %v","err_error_closing_tailscale_server":"Error closing Tailscale server: %v","err_failed_to_create_tailscale_funnel_listener":"Failed to create Tailscale funnel listener on %s: %v","err_failed_to_create_tailscale_listener":"Failed to create Tailscale listener on %s: %v","err_failed_to_create_tailscale_local_client":"Failed to create Tailscale local client: %v","err_error_stopping_tailscale_server":"Error stopping Tailscale server: %v","err_failed_to_set_config_value":"Failed to set config value: %v","err_failed_to_parse_bool":"Failed to parse \'%s\' as bool: %v","err_failed_to_parse_int":"Failed to parse \'%s\' as int: %v","err_failed_to_set_username":"Failed to set Username: %v","err_failed_to_render_userinfo_config_template":"Failed to render UserInfoConfig template: %v","err_failed_to_add_config_value":"Failed to add config value: %v","err_getbookshelf_error":"GetBookShelf Error: %v","err_failed_to_run_tailscale":"Failed to run Tailscale: %v","err_store_url_already_exists_error":"store Url already exists: %s","err_store_path_is_subdir_of_existing":"New store path is a subdirectory of existing store: %s is a subdirectory of %s","err_store_path_is_parent_of_existing":"New store path is a parent directory of existing store: %s is a parent directory of %s","err_store_path_conflict":"Store path conflict","err_invalid_store_path":"Invalid store path: %s","err_field_not_exists":"field \'%s\' does not exist","err_field_cannot_set":"cannot set field \'%s\'","err_field_not_slice_type":"field \'%s\' is not a slice type","err_field_element_not_string":"field \'%s\' element type is not string","err_field_type_not_supported":"field \'%s\' type is not supported for setting: %s","log_failed_to_create_extract_path":"Failed to create extract path: %v","log_failed_to_open_file_unarchive":"Failed to open file: %v","log_failed_to_identify_file_format":"Failed to identify file format: %v","log_failed_to_extract_zip_file":"Failed to extract zip file: %v","log_zip_file_extracted":"ZIP file extracted: %s to %s","log_failed_to_extract_rar_file":"Failed to extract RAR file: %v","log_rar_file_extracted":"RAR file extracted: %s to %s","log_file_close_error":"file.Close() Error: %s","log_failed_to_create_directory":"Failed to create directory: %v","log_failed_to_create_parent_directory":"Failed to create parent directory: %v","log_failed_to_create_file":"Failed to create file: %v","log_failed_to_copy_file_content":"Failed to copy file content: %v","log_failed_to_get_container_xml":"Failed to get container.xml: %s","log_failed_to_parse_container_xml":"Failed to parse container.xml: %s","log_html_tokenizer_error":"HTML tokenizer error: %v","log_failed_to_get_opf_file_path":"Failed to get OPF file path: %s","log_failed_to_read_opf_file":"Failed to read OPF file: %s","log_failed_to_parse_opf_file":"Failed to parse OPF file: %s","log_failed_to_write_file_to_cache":"Failed to write file to cache: %v","log_content_type_not_found_in_cache":"ContentType not found in cache for key: %+v","log_failed_to_read_file_from_cache":"Failed to read file from cache: %v","log_cover_file_not_found":"Cover file not found","log_created_new_book":"Created new book: %s","log_updated_existing_book":"Updated existing book: %s %s","log_saved_media_files_for_book":"Saved %d media files for book %s","log_saved_bookmarks_for_book":"Saved %d bookmarks for book %s","log_get_media_files_for_book_error":"Get media files for book %s error: %s","log_get_bookmarks_for_book_error":"Get bookmarks for book %s error: %s","log_error_creating_new_book_group":"Error creating new book group: %s","log_failed_to_read_directory":"Failed to read directory: %s, error: %v","log_failed_to_get_file_info_scan":"Failed to get file info: %s, error: %v","log_nonutf8_zip_error":"NonUTF-8 ZIP: %s, Error: %s","log_failed_to_get_image_list_from_epub":"Failed to get image list from EPUB: %s, error: %v","log_failed_to_get_metadata_from_epub":"Failed to get metadata from EPUB: %s, error: %v","log_failed_to_access_path_in_archive":"Failed to access path %s in archive: %v","log_skip_scan_path":"Skip Scan: %s","log_skip_unsupported_file_type":"Skipping unsupported file type: %s","log_found_supported_file_type":"Found supported file type (non-media): %s","log_processing_file":"Processing file: %s (path: %s)","log_failed_to_get_file_info_in_archive":"Failed to get file info in archive: %v","log_books_saved_to_database_successfully":"SaveBooksToDatabase: Books saved to database successfully: %d","log_failed_to_store_bookmark":"Failed to store bookmark: %s","log_updated_bookmarks_for_book_id":"Updated bookmarks for book ID %s: %s","log_getbook_error_common":"GetBook error: %s","log_getpicturedata_error":"GetPictureData error: %s","log_getbook_error_scroll":"GetBook: %v","log_server_not_ready_within_timeout":"Server not ready within %v, continue anyway","log_open_database_error":"OpenDatabase Error: %s","log_loadbooks_error":"LoadBooks_error %s","log_savebooks_error":"SaveBooks_error %s","log_error_listing_books_from_database":"Error listing books from database: %s","log_failed_savebookstodatabase":"Failed SaveBooksToDatabase: %v","log_failed_to_get_free_port":"Failed to get a free port: %v","log_using_port":"Using port: %d","log_http_newrequest_error":"http.NewRequest Error: %s","log_client_do_error":"client.Do Error: %s","log_body_close_error":"Body.Close() Error: %s","log_io_readall_error":"io.ReadAll Error: %s","log_response_write_error":"Response.Write Error: %s","log_getimagefrompdf_imgdata_nil":"GetImageFromPDF: imgData is nil","log_countpages_pdf_invalid_error":"CountPagesOfPDF: invalid PDF: %v Error: %v","log_getimagefrompdf_time":"GetImageFromPDF: %v","log_syncpage_message_to_flipmode":"SyncPage message toFlipMode: %v %v","log_syncpage_message_to_scrollmode":"SyncPage message to ScrollMode: %v %v","log_failed_to_open_database":"Failed to open database: %v","log_failed_to_ping_database":"Failed to ping database: %v","log_failed_to_create_tables":"Failed to create tables: %v","log_database_initialized_successfully":"Database initialized successfully","err_failed_to_read_embedded_image":"Failed to read embedded image: %v","err_failed_to_read_embedded_data":"Failed to read embedded data: %v","scroll_wheel_flip":"Scroll Wheel Flip","open_in_new_tab":"Open in New Tab","auto_rescan_interval_minutes":"Automatic rescan interval","AutoRescanIntervalMinutes":"Automatic rescan interval","auto_rescan_interval_minutes_desc":"In minutes. Set to 0 to disable automatic scanning","log_scheduler_task_still_running_skip":"Previous scan task is still running, skipping this scan","log_scheduler_task_execution_failed":"Scheduled scan task execution failed: %v","log_scheduler_task_execution_completed":"Scheduled scan task execution completed","log_scheduler_interval_zero_no_scheduled_scan":"Scan interval is 0, no scheduled scanning","log_scheduler_stop_task_failed":"Failed to stop scheduled task: %v","log_scheduler_stop_old_task_failed":"Failed to stop old scheduled task: %v","log_scheduler_create_scheduler_failed":"Failed to create scheduler: %v","log_scheduler_create_task_failed":"Failed to create scheduled task: %v","log_scheduler_task_started":"Scheduled scan task started, interval: %d minutes","log_scheduler_task_stopped":"Scheduled scan task stopped","auto_rescan_started":"Automatic rescan started, interval: %d minutes","auto_rescan_stopped":"Automatic rescan stopped","auto_rescan_updated":"Automatic rescan interval updated to: %d minutes","auto_rescan_status_running":"Running","auto_rescan_status_stopped":"Stopped","auto_rescan_current_interval":"Current scan interval: %d minutes","auto_rescan_enabled_hint":"Automatic rescan has been enabled. The system will periodically scan the library","auto_rescan_disabled_hint":"Automatic rescan has been disabled","other_settings":"Other Settings","plugin_enable":"Enable Plugin","err_failed_to_get_config_dir":"Failed to get config directory: %v","comigo_xyz_title":"Comigo","comigo_xyz_subtitle":"Comigo - Simple and EasyUse Manga Reader","comigo_xyz_description":"Help you read manga on all devices - Whether computer or phone, Windows, Linux,\u3000MacOS","comigo_xyz_feature_format":"\uD83D\uDCDA Multi-format Support","comigo_xyz_feature_format_desc":"Supports ZIP, RAR, CBZ, EPUB, PDF and other comic formats","comigo_xyz_feature_cross_platform":"\uD83C\uDF10 Cross-platform","comigo_xyz_feature_cross_platform_desc":"Supports Linux, Windows, Mac OS operating systems","comigo_xyz_feature_responsive":"\uD83D\uDCF1 Responsive Design","comigo_xyz_feature_responsive_desc":"Adapts to desktop and mobile devices","comigo_xyz_feature_reading_modes":"\uD83D\uDD04 Multiple Reading Modes","comigo_xyz_feature_reading_modes_desc":"Supports flip mode and scroll mode to meet different reading habits","comigo_xyz_feature_security":"\uD83D\uDD12 Secure and Reliable","comigo_xyz_feature_security_desc":"Supports HTTPS and user authentication, built-in Tailscale remote access","comigo_xyz_feature_performance":"\u26A1 High Performance","comigo_xyz_feature_performance_desc":"Fast scanning and loading for smooth reading experience","comigo_xyz_feature_plugin":"\uD83D\uDD0C Plugin System","comigo_xyz_feature_plugin_desc":"Auto page-turn, clock and more plugins, with custom plugin support","comigo_xyz_feature_media":"\uD83C\uDFAC Media Playback","comigo_xyz_feature_media_desc":"Built-in audio and video player","comigo_xyz_feature_download":"\uD83D\uDCE5 Flexible Download","comigo_xyz_feature_download_desc":"Batch download image folders, convert to EPUB format","comigo_xyz_feature_history":"\uD83D\uDCDC Reading History","comigo_xyz_feature_history_desc":"Automatic reading history tracking for easy continuation","comigo_xyz_github_repo":"\uD83D\uDCE6 GitHub Repository","comigo_xyz_github_button":"Visit GitHub Project","comigo_xyz_download_latest":"Download Latest Version","comigo_xyz_quick_start":"Quick Start","comigo_xyz_quick_start_step1":"Download and run Comigo","comigo_xyz_quick_start_step2":"Configure your library path in settings","comigo_xyz_quick_start_step3":"Start enjoying your reading experience!","comigo_xyz_download_for":"Download for","comigo_xyz_select_platform":"Select Platform","comigo_xyz_download":"\u2B07\uFE0F Download","comigo_xyz_cli_install_title":"One-Click CLI Installation:","comigo_xyz_cli_install_desc":"Use the following commands to quickly install Comigo","comigo_xyz_cli_install_copy":"Copy","comigo_xyz_cli_install_copied":"Copied","comigo_xyz_cli_install_cn_desc":"Recommended for Mainland China users:","comigo_xyz_docker_deploy_title":"Deploy with Docker:","rescan_store":"Scan","delete_store":"Delete","select_store_to_operate":"Please select a store to operate","confirm_delete_store":"Confirm to delete the store? This will also delete all book data in this store","rescan_store_success":"Store scan completed, {0} new books added","delete_store_success":"Store deleted successfully","rescan_store_in_progress":"Scanning store, please wait...","err_rescan_store_failed":"Failed to rescan store","err_delete_store_failed":"Failed to delete store","log_book_file_not_exist_skip":"Book file does not exist, skipping load: %s","log_error_deleting_orphan_metadata":"Error deleting orphan metadata file %s: %s","log_book_version_mismatch_skip":"Book %s version mismatch (cached: %s, current: %s), skipping load","log_book_version_minor_mismatch":"Book %s minor version differs (cached: %s, current: %s), will migrate bookmarks and rescan","log_bookmark_saved_for_migration":"Saved %d bookmarks for book %s, pending migration","log_bookmark_migrated":"Successfully migrated %d bookmarks for book %s","log_error_deleting_version_mismatch_metadata":"Error deleting version mismatch metadata file %s: %s","player":"Player","play":"Play","pause":"Pause","previous":"Previous","next":"Next","playlist":"Playlist","audio":"Audio","video":"Video","browser_not_support_video":"Your browser does not support video playback","browser_not_support_audio":"Your browser does not support audio playback","play_failed":"Playback failed","first_media":"This is the first item","last_media":"This is the last item","auto_play_next":"Auto play next","loop_playlist":"Loop playlist","player_autoplay_help":"Auto play next and loop require media playback permission in your browser. On mobile devices, it may also fail due to battery saver/background restrictions.","add_bookmark":"Add Bookmark","manual_bookmark":"User Bookmark","auto_bookmark":"Auto Bookmark","bookmark_added":"Bookmark added","bookmark_exists":"Bookmark exists for this page","continue_reading":"Continue","delete_record":"Delete","confirm_delete_bookmark":"Delete this reading record?","bookmark_deleted":"Deleted","log_remote_store_connect_failed":"Failed to connect to remote store: %s, error: %v","log_remote_file_stat_failed":"Failed to get remote file info: %s, error: %v","log_remote_file_open_failed":"Failed to open remote file: %s, error: %v","log_remote_pdf_download_on_demand":"Downloading remote PDF on demand: %s","log_webdav_filesystem_connected":"WebDAV filesystem connected: %s, base path: %s","log_sftp_filesystem_connected":"SFTP filesystem connected: %s, base path: %s","log_smb_filesystem_connected":"SMB filesystem connected: %s, base path: %s","log_ftp_filesystem_connected":"FTP filesystem connected: %s, base path: %s","log_s3_filesystem_connected":"S3 filesystem connected: %s, base path: %s","log_epub_metadata_remote_not_supported":"EPUB metadata extraction is not supported for remote streaming","log_remote_file_download_to_cache":"Downloading remote file to cache: %s -> %s","log_scan_remote_store_start":"Starting to scan remote store: %s","log_scan_start_hint_remote":"Scan: %s (remote path: %s)","log_remote_store_check_book_existence_failed":"Failed to connect to remote store to check book existence: %s, error: %v","log_remote_book_existence_check_failed":"Failed to check remote book existence: %s, error: %v","log_add_remote_store":"Adding remote store: %s (protocol: %s, host: %s)","log_s3_connecting":"Connecting to S3 service %s (bucket: %s, prefix: %s)","log_ftp_connecting":"Connecting to FTP server %s (TLS: %v, timeout: %v)","log_webdav_download_range":"Download range: %s [%d-%d]","log_smb_connecting":"Connecting to SMB server %s (timeout: %ds, user: %s, share: %s)","log_smb_mount_share":"Mounting SMB share: %s","log_using_cached_file":"Using cached file: %s","log_download_file":"Downloading file: %s","log_remote_book_existence_check_failed_detail":"Remote book existence check failed - BookID: %s, RemoteURL: %s, BookPath: %s, error: %v","log_local_book_existence_check_failed":"Local book existence check failed: %s, error: %v","log_failed_to_set_field":"Failed to set field %s: %v","log_cache_mkdir_failed":"Failed to create cache directory: %v","log_cache_hit_memory":"Cache hit (memory): %s","log_cache_hit_disk":"Cache hit (disk): %s","log_cache_write_disk_failed":"Failed to write to disk cache: %v","log_cached_to_disk":"Cached to disk: %s -> %s","log_get_file_error":"Get file error: %s","log_save_file_to_cache_error":"SaveFileToCache error: %s","log_save_cover_to_local_error":"SaveCoverToLocal error: %s","log_failed_to_parse_cover_url":"Failed to parse cover URL: %s","log_child_book_id_missing_in_cover_url":"Child book ID is missing in cover URL","log_failed_to_get_child_book":"Failed to get child book: %s","log_failed_to_read_embedded_image":"Failed to read embedded image: %s","log_failed_to_create_epub_generator":"Failed to create EPUB generator: %s","log_failed_to_get_image_epub":"Failed to get image %s: %v","log_failed_to_generate_epub":"Failed to generate EPUB: %s","log_plugin_enabled":"Plugin enabled: %s","log_disable_mutex_plugin_auto_flip":"Disabling mutex plugin: auto_flip","log_disable_mutex_plugin_sketch_practice":"Disabling mutex plugin: sketch_practice","log_plugin_disabled":"Plugin disabled: %s","log_rescan_store":"Rescanning store: %s","log_rescan_store_completed_new_books":"Store scan completed, %d new books added","log_delete_store":"Deleting store: %s","log_deleted_books_count":"Deleted %d books","log_upload_no_store_selected":"No upload target store selected","log_upload_invalid_store_path":"Invalid store path: %s","log_upload_store_path_not_exist":"Store path does not exist: %s","log_error_closing_zip_writer":"Error closing zip writer: %s","log_file_not_found_skipping":"File not found, skipping: %s","log_error_opening_file":"Error opening file: %s, error: %s","log_error_creating_zip_entry":"Error creating zip entry: %s, error: %s","log_error_writing_file_to_zip":"Error writing file to zip: %s, error: %s","log_failed_to_delete_bookmark":"Failed to delete bookmark: %v","log_get_config_dir_error":"GetConfigDir error: %s","log_delete_book_json_error":"DeleteBookJson error: %s","log_delete_cover_cache_error":"DeleteCoverCache error: %s","log_delete_book_cache_error":"DeleteBookCache error: %s","log_plugin_system_disabled_skip_scan":"Plugin system disabled, skipping custom plugin scan","log_plugin_dir_not_exist_skip_load":"Plugin directory does not exist: %s, skipping custom plugin load","log_plugin_scope_load_error":"Error loading plugin for scope %s: %v","log_plugin_custom_loaded_count":"Successfully loaded %d custom plugins","log_plugin_loaded_item":"  - [%s] %s (%s)","log_plugin_dir_not_exist":"Plugin directory does not exist: %s","log_plugin_read_file_failed":"Failed to read plugin file %s: %v","log_plugin_read_book_file_failed":"Failed to read book plugin file %s: %v","log_plugin_loaded_for_book":"For book %s, loaded %d %s plugin(s)","log_load_custom_plugin_failed":"Failed to load custom plugin: %v","log_cfg_host_enabled_plugin_list":"cfg.Host: %v , cfg.EnabledPluginList: %v","log_config_changed_restart_web":"Config changed, restarting web server...","log_config_changed_start_tailscale":"Config changed, starting tailscale...","log_config_changed_stop_tailscale":"Config changed, stopping tailscale...","log_config_changed_restart_tailscale":"Config changed, restart tailscale...","log_received_rescan_message":"Received rescan message: %s","log_failed_to_read_image_epub":"Failed to read image %s: %v","log_failed_to_decode_image_config_epub":"Failed to decode image config: %v, using default dimensions","log_cannot_set_username_or_password":"Cannot set Username or Password","log_login_failed":"Login failed: %s-%s, %s-%s","log_get_generated_image_params":"GetGeneratedImage: height=%s, width=%s, text=%s, font_size=%s","log_toml_marshal_error":"toml.Marshal Error","log_config_file_changed_detected":"Comigo config file change detected: %s Op: %s","log_error_num_value":"Error num value: %s"}');


var $8c66835f5815577e$exports = {};
$8c66835f5815577e$exports = JSON.parse('{"auth":"\u8BA4\u8BC1","book_not_found":"\u672A\u627E\u5230\u53EF\u9605\u8BFB\u7684\u4E66\u7C4D\uFF0C\u7A0B\u5E8F\u9000\u51FA\u3002","cache_file_clean":"\u9000\u51FA\u65F6\u6E05\u9664\u5168\u90E8\u7F13\u5B58\u6587\u4EF6","cache_file_enable":"\u662F\u5426\u4FDD\u5B58web\u56FE\u7247\u7F13\u5B58\uFF0C\u53EF\u52A0\u5FEB\u4E8C\u6B21\u8BFB\u53D6\u4F46\u4F1A\u5360\u7528\u786C\u76D8\u7A7A\u95F4","cache_file_dir":"\u7F13\u5B58\u6587\u4EF6\u5939\uFF0C\u9ED8\u8BA4\u4E3A\u7CFB\u7EDF\u4E34\u65F6\u76EE\u5F55\uFF0C\u7A0B\u5E8F\u9000\u51FA\u65F6\u53EF\u80FD\u88AB\u6E05\u7A7A","can_not_init_book":"\u65E0\u6CD5\u521D\u59CB\u5316\u4E66\u7C4D\u3002","cannot_listen":"\u65E0\u6CD5\u76D1\u542C\u7AEF\u53E3\uFF1A","check_image_completed":"\u56FE\u7247\u89E3\u6790\u5B8C\u6210\u3002","check_image_error":"\u5206\u8FA8\u7387\u5206\u6790\u51FA\u9519\uFF1A","check_image_ing":"\u6B63\u5728\u5206\u6790\u56FE\u7247\u5206\u8FA8\u7387\uFF1A","check_image_start":"\u5F00\u59CB\u89E3\u6790\u56FE\u7247\u2026\u2026","check_mac_error":"\u68C0\u6D4BMac\u5730\u5740\u51FA\u9519\uFF1A","check_port_error":"\u68C0\u6D4B\u7AEF\u53E3\u51FA\u9519\uFF1A%v","clear_temp_file_completed":"\u4E34\u65F6\u6587\u4EF6\u6E05\u7406\u6210\u529F\uFF1A","clear_temp_file_error":"\u4E34\u65F6\u6587\u4EF6\u6E05\u7406\u5931\u8D25\uFF1A","clear_temp_file_start":"\u5F00\u59CB\u6E05\u7406\u4E34\u65F6\u6587\u4EF6\u3002","comigo_example":"  comi book.zip\\n\\n\u8BBE\u5B9A\u7F51\u9875\u670D\u52A1\u7AEF\u53E3\uFF08\u9ED8\u8BA41234\uFF09\uFF1A\\n  comi -p 2345 book.zip\\n\\n\u4E0D\u6253\u5F00\u6D4F\u89C8\u5668\uFF08windows\uFF09\uFF1A\\n  comi -o=false book.zip\\n\\n\u6307\u5B9A\u591A\u4E2A\u53C2\u6570\uFF1A\\n  comi -p 2345 --host example.com test.zip\\n","comigo_use":"comi","completed_extract":"\u89E3\u538B\u5B8C\u6210\uFF1A","completed_ls":"\u538B\u7F29\u6587\u4EF6\u626B\u63CF\u5B8C\u6210","config":"\u6307\u5B9A\u914D\u7F6E\u6587\u4EF6","config_change":"comigo\u914D\u7F6E\u53D1\u751F\u53D8\u66F4\uFF1A","config_file_not_found":"\u672A\u627E\u5230\u914D\u7F6E\u6587\u4EF6\uFF0C\u4F7F\u7528\u9ED8\u8BA4\u8BBE\u7F6E\u3002","config_file_not_resolve":"\u89E3\u6790\u914D\u7F6E\u5931\u8D25\u3002","config_save_to":"\u914D\u7F6E\u6587\u4EF6\u9ED8\u8BA4\u4FDD\u5B58\u4F4D\u7F6E\uFF0C\u53EF\u9009\u503C\uFF1ARAM\u3001HomeDirectory\u3001WorkingDirectory\u3001ProgramDirectory","ctrl_c_hint":"\u6309 CTRL-C \u9000\u51FA","debug_mode":"\u8C03\u8BD5\u6A21\u5F0F","disable_lan":"\u7981\u7528LAN\u5171\u4EAB","enable_database":"\u542F\u7528\u672C\u5730\u6570\u636E\u5E93\u4FDD\u5B58\u5DF2\u626B\u63CF\u7684\u4E66\u7C4D\u6570\u636E","enable_file_upload":"\u542F\u7528\u6587\u4EF6\u4E0A\u4F20\u529F\u80FD","enable_frpc":"\u542F\u7528frp\u53CD\u5411\u4EE3\u7406","enable_webp":"\u542F\u7528webp\u4F20\u8F93\uFF0C\u9700\u8981webp-server","epub_cannot_resort":"\u65E0\u6CD5\u5BF9epub\u6587\u4EF6\u91CD\u65B0\u6392\u5E8F\uFF1A","exceeds_maximum_depth":"\u8D85\u8FC7\u6700\u5927\u641C\u7D22\u6DF1\u5EA6\uFF0CMaxDepth=","file_exit":"\u6587\u4EF6\u5DF2\u5B58\u5728\uFF1A","file_not_found":"\u672A\u627E\u5230\u5408\u9002\u7684\u6587\u4EF6\u3002","flip_page_template":"\u9ED8\u8BA4\u6A21\u677F\uFF1A\u7FFB\u9875\u9605\u8BFB","format_customization_error":"\u4E0D\u652F\u6301\u7684\u89E3\u538B\u683C\u5F0F\uFF1A","found_config_file":"\u53D1\u73B0\u914D\u7F6E\u6587\u4EF6\uFF1A","found_in_bookstore":"\u4E66\u7C4D\u5DF2\u5B58\u5728\u4E66\u5E93\u4E2D\uFF0C\u8DF3\u8FC7\u626B\u63CF\uFF1A","how_many_books_update":"\u8DEF\u5F84 %v \u66F4\u65B0 %v \u672C\u4E66","frp_command":"frpc\u547D\u4EE4\u6216frpc\u53EF\u6267\u884C\u6587\u4EF6\u8DEF\u5F84","frp_random_remote_port":"frp\u8FDC\u7A0B\u968F\u673A\u7AEF\u53E3\uFF0840000~50000\uFF09","frp_remote_port":"frp\u8FDC\u7A0B\u7AEF\u53E3\uFF0C\u5982\u8BBE\u4E3A-1\u5219\u4E0E\u672C\u5730\u7AEF\u53E3\u76F8\u540C","frp_server_addr":"frps-addr\uFF08\u9700\u8981frpc\uFF09","frp_server_error":"\u65E0\u6CD5\u542F\u52A8frpc\u670D\u52A1\uFF0C\u8BF7\u68C0\u67E5\u547D\u4EE4\u683C\u5F0F\u5E76\u786E\u8BA4PATH\u4E2D\u6709frpc\u53EF\u6267\u884C\u6587\u4EF6\u3002","frp_server_port":"frps server_port\uFF08\u9700\u8981frpc\uFF09","frp_setting_save_completed":"\u6210\u529F\u4FDD\u5B58frpc\u8BBE\u7F6E\u3002","frp_token":"token\uFF08\u9700\u8981frpc\uFF09","frpc_ini_error":"frpc ini\u521D\u59CB\u5316\u9519\u8BEF\u3002","frpc_server_start":"frpc\u5DF2\u542F\u52A8\u3002","generate_metadata":"\u751F\u6210\u4E66\u7C4D\u5143\u6570\u636E","get_ip_error":"\u83B7\u53D6IP\u51FA\u9519\uFF1A","html_title":"Comigo \u6F2B\u753B\u9605\u8BFB\u5668","init_database":"\u521D\u59CB\u5316\u6570\u636E\u5E93\uFF1A","init_locale":"\u9ED8\u8BA4\u8BED\u8A00\u4E3A\u4E2D\u6587\u3002","local_host":"\u81EA\u5B9A\u4E49\u57DF\u540D","local_reading":"\u672C\u673A\u9605\u8BFB\uFF1A","log_to_file":"\u8BB0\u5F55\u65E5\u5FD7\u5230\u6587\u4EF6","long_description":"comigo\uFF0C\u7B80\u5355\u7684\u6F2B\u753B\u9605\u8BFB\u5668\u3002","max_depth":"\u6700\u5927\u641C\u7D22\u6DF1\u5EA6","min_media_num":"\u81F3\u5C11\u5305\u542B\u591A\u5C11\u5A92\u4F53\u6587\u4EF6\u624D\u8BA4\u5B9A\u4E3A\u6F2B\u753B\u538B\u7F29\u5305","no_pages_in_pdf":"PDF\u4E2D\u65E0\u9875\u9762","not_a_valid_zip_file":"\u4E0D\u662F\u5408\u6CD5\u7684zip\u6587\u4EF6\uFF1A","open_browser":"\u542F\u52A8\u65F6\u540C\u65F6\u6253\u5F00\u6D4F\u89C8\u5668\uFF08windows=true\uFF09","open_browser_error":"\u6253\u5F00\u6D4F\u89C8\u5668\u5931\u8D25\u3002","open_image_error":"\u6253\u5F00\u56FE\u7247\u65F6\u51FA\u9519\uFF1A","password":"\u5BC6\u7801","path_not_exist":"\u8DEF\u5F84\u4E0D\u5B58\u5728","port":"\u670D\u52A1\u7AEF\u53E3","port_busy":"%v \u7AEF\u53E3\u88AB\u5360\u7528\uFF0C\u5C1D\u8BD5\u4F7F\u7528\u968F\u673A\u7AEF\u53E3","print_all_ip":"\u6253\u5370\u6240\u6709\u53EF\u7528\u7F51\u5361ip","print_config":"\u6253\u5370\u914D\u7F6E\uFF1A","re_enter_password":"\u91CD\u65B0\u8F93\u5165\u5BC6\u7801","reading_url_maybe":"\u9605\u8BFB\u94FE\u63A5\u53EF\u80FD\u4E3A\uFF1A","reg_file_hint":"\u53CC\u51FB\u5BFC\u5165\u6B64\u6587\u4EF6\u53EF\u5728Windows\u6CE8\u518C\u53F3\u952E\u83DC\u5355","requires_login":"\u542F\u7528\u767B\u5F55\u4FDD\u62A4","rescan":"\u91CD\u65B0\u626B\u63CF","save":"\u4FDD\u5B58","delete":"\u5220\u9664","save_config_file":"\u4FDD\u5B58\u914D\u7F6E\uFF1A","scan_archive_error":"\u6587\u4EF6\u626B\u63CF\u51FA\u9519","scan_error":"\u626B\u63CF\u51FA\u9519\uFF1A","scan_ing":"\u6B63\u5728\u626B\u63CF\u6587\u4EF6\u2026\u2026","scan_pdf":"\u626B\u63CFPDF\uFF1A","scan_start_hint":"\u5F00\u59CB\u626B\u63CF\uFF1A","scroll_template":"\u9ED8\u8BA4\u6A21\u677F\uFF1A\u5377\u8F74\u9605\u8BFB","short_description":"\u4E00\u4E2A\u7B80\u5355\u7684\u6F2B\u753B\u9605\u8BFB\u5668\u3002","shutdown_hint":"\u7A0B\u5E8F\u6B63\u5728\u9000\u51FA\uFF0C\u6309 Ctrl+C \u518D\u6B21\u5F3A\u5236\u9000\u51FA","sketch_count_seconds":"\u901F\u5199\u6A21\u5F0F\u5012\u8BA1\u65F6\uFF08\u79D2\uFF09","sketch_template":"\u9ED8\u8BA4\u6A21\u677F\uFF1A\u8349\u7A3F\u53C2\u8003","skip_path":"\u5FFD\u7565\u8DEF\u5F84\uFF1A","settings_page":"\u8BBE\u7F6E\u9875\u9762","upload_page":"\u4E0A\u4F20\u9875\u9762","404notfound":"404 \u9875\u9762\u672A\u627E\u5230","sort":"\u56FE\u7247\u6392\u5E8F\u89C4\u5219\uFF08none\u3001name\u3001time\uFF09","sort_by_name":"\u6309\u6587\u4EF6\u540D\u6392\u5E8F","sort_by_time":"\u6309\u4FEE\u6539\u65F6\u95F4\u6392\u5E8F","start_clear_file":"\u8FD0\u884C\u4E2D\u65AD\uFF0C\u6E05\u7406\u4E34\u65F6\u6587\u4EF6\u5939\u4E2D\u2026\u2026","start_extract":"\u5F00\u59CB\u89E3\u538B\uFF1A","start_in_background":"\u540E\u53F0\u8FD0\u884CComigo","start_ls":"\u5F00\u59CB\u626B\u63CF\u538B\u7F29\u6587\u4EF6\uFF1A","static_file_mode":"\u662F\u5426\u5F00\u542F\u9759\u6001\u6587\u4EF6\u6A21\u5F0F\u3002\u9759\u6001\u6A21\u5F0F\u4E0B\uFF0C\u6240\u6709\u56FE\u7247\u4E0E\u811A\u672C\u90FD\u6253\u5305\u5230html\u6587\u4EF6\u91CC\uFF0C\u53EF\u4EE5\u76F4\u63A5\u53E6\u5B58\u4E3A\u5355\u4E2A\u7F51\u9875\u6587\u4EF6\uFF08\u5F00\u53D1\u4E2D\uFF0C\u4EC5\u5377\u8F74\u6A21\u5F0F+\u65E0\u9650\u4E0B\u62C9\u7684\u9605\u8BFB\u9875\u9762\u6709\u6548\uFF09\u3002","stop_background":"\u505C\u6B62\u540E\u53F0\u8FD0\u884C\u7684\u8FDB\u7A0B","target_path":"\u76EE\u6807\u76EE\u5F55","temp_folder_create_error":"\u521B\u5EFA\u4E34\u65F6\u6587\u4EF6\u5939\u5931\u8D25\u3002","temp_folder_error":"\u4E34\u65F6\u6587\u4EF6\u5939\u8BBE\u7F6E\u5931\u8D25\u3002","temp_folder_path":"\u4E34\u65F6\u6587\u4EF6\u5939\u8DEF\u5F84\uFF1A","template":"\u9ED8\u8BA4\u9875\u9762\u6A21\u677F\uFF08scroll\u3001flip\u3001sketch\uFF09","timeout":"\u8D85\u65F6\u65F6\u95F4\uFF08\u5206\u949F\uFF09","tls_crt":"TLS/SSL\u8BC1\u4E66\u6587\u4EF6\u8DEF\u5F84","tls_enable":"\u542F\u7528TLS/SSL","auto_https_cert":"\u81EA\u52A8\u7533\u8BF7\u3001\u7B7E\u53D1 HTTPS \u8BC1\u4E66\uFF08Let\'s Encrypt\uFF09","port443_busy_disable_auto_tls":"443 \u7AEF\u53E3\u5DF2\u88AB\u5360\u7528\uFF0C\u5DF2\u7981\u7528\u81EA\u52A8 TLS\u3002","auto_tls_disabled_invalid_domain":"\u81EA\u52A8 TLS \u9700\u8981\u6709\u6548\u7684\u57DF\u540D\u624D\u80FD\u5DE5\u4F5C\uFF0C\u5DF2\u7981\u7528\u81EA\u52A8 TLS\u3002","auto_tls_disabled_custom_cert_set":"\u5DF2\u8BBE\u7F6E\u81EA\u5B9A\u4E49\u8BC1\u4E66\uFF0C\u5DF2\u7981\u7528\u81EA\u52A8 TLS\u3002","auto_tls_disabled_lan_access_off":"\u7981\u7528\u5C40\u57DF\u7F51\u8BBF\u95EE\u65F6\uFF0C\u5DF2\u7981\u7528\u81EA\u52A8 TLS\u3002","tls_key":"TLS/SSL\u5BC6\u94A5\u6587\u4EF6\u8DEF\u5F84","un_archive_error":"\u6587\u4EF6\u89E3\u538B\u51FA\u9519","unable_to_extract_images_from_pdf":"\u65E0\u6CD5\u4ECEPDF\u4E2D\u63D0\u53D6\u56FE\u7247","unmarshal_config_file_error":"\u5E94\u7528\u914D\u7F6E\u6587\u4EF6\u5931\u8D25\uFF0C\u8BF7\u68C0\u67E5\u683C\u5F0F\u3002","unsupported_extract":"\u4E0D\u652F\u6301\u89E3\u538B\uFF1A","unsupported_file_type":"\u4E0D\u652F\u6301\u7684\u6587\u4EF6\u7C7B\u578B\uFF1A","upload_disable_hint":"\u4E0A\u4F20\u529F\u80FD\u5DF2\u7981\u7528","upload_path":"\u4E0A\u4F20\u6587\u4EF6\u7684\u4FDD\u5B58\u8DEF\u5F84","username":"\u7528\u6237\u540D","web_server_error":"web\u670D\u52A1\u542F\u52A8\u5931\u8D25\uFF0C\u7AEF\u53E3\uFF1A","webp_command":"webp-server\u547D\u4EE4\u6216webp-server\u53EF\u6267\u884C\u6587\u4EF6\u8DEF\u5F84","webp_quality":"webp\u538B\u7F29\u8D28\u91CF","webp_server_error":"\u65E0\u6CD5\u542F\u52A8webp\u8F6C\u6362\u670D\u52A1\uFF0C\u8BF7\u68C0\u67E5\u547D\u4EE4\u683C\u5F0F\u5E76\u786E\u8BA4PATH\u4E2D\u6709webp-server\u53EF\u6267\u884C\u6587\u4EF6\u3002","webp_server_start":"webp\u8F6C\u6362\u670D\u52A1\u5DF2\u542F\u52A8","webp_setting_error":"webp\u8BBE\u7F6E\u9519\u8BEF\u3002","webp_setting_save_completed":"webp\u8BBE\u7F6E\u4FDD\u5B58\u9519\u8BEF\u3002","websocket_error":"websocket\u9519\u8BEF\uFF1A","websocket_messages":"websocket\u4FE1\u606F\uFF1A","zip_encode":"\u6307\u5B9Azip\u6587\u4EF6\u7F16\u7801\uFF08gbk\u3001shiftjis\u7B49\uFF09","lang":"\u754C\u9762\u8BED\u8A00\u8BBE\u7F6E\uFF08auto\u3001zh\u3001en\u3001ja\uFF09\uFF0C\u9ED8\u8BA4\u4E3Aauto\uFF08\u81EA\u52A8\u68C0\u6D4B\uFF09","CacheDir":"\u672C\u5730\u7F13\u5B58\u6587\u4EF6\u5939","CacheDir_Description":"\u672C\u5730\u6587\u4EF6\u7684\u7F13\u5B58\u4F4D\u7F6E\uFF0C\u9ED8\u8BA4\u7CFB\u7EDF\u4E34\u65F6\u6587\u4EF6\u5939\u3002","CertFile":"CertFile","CertFile_Description":"TLS/SSL \u8BC1\u4E66\u6587\u4EF6\u8DEF\u5F84 (default: \u3001\\"~/.config/.comigo/cert.crt\\")","ClearCacheExit":"\u9000\u51FA\u65F6\u6E05\u7406","ClearCacheExit_Description":"\u9000\u51FA\u7A0B\u5E8F\u7684\u65F6\u5019\uFF0C\u6E05\u7406web\u56FE\u7247\u7F13\u5B58\u3002","ClearDatabaseWhenExit":"\u6E05\u9664\u6570\u636E\u5E93\u4E66\u7C4D","ClearDatabaseWhenExit_Description":"\u542F\u7528\u672C\u5730\u6570\u636E\u5E93\u65F6\uFF0C\u626B\u63CF\u5B8C\u6210\u540E\uFF0C\u6E05\u9664\u4E0D\u5B58\u5728\u7684\u4E66\u7C4D\u3002","ConfigManager":"\u914D\u7F6E\u6587\u4EF6\u7BA1\u7406","ConfigManagerDeleteSuccess":"\u914D\u7F6E\u5DF2\u5220\u9664\u3002","ConfigManagerDescription":"\u70B9\u51FBSave\uFF0C\u4F1A\u5C06\u5F53\u524D\u914D\u7F6E\u4E0A\u4F20\u5230\u670D\u52A1\u5668\uFF0C\u5E76\u8986\u76D6\u5DF2\u7ECF\u5B58\u5728\u7684\u8BBE\u7F6E\u6587\u4EF6\u3002","ConfigManagerSaveHint":"\u5DF2\u6709\u914D\u7F6E\u6587\u4EF6,\u8BF7\u5207\u6362\u4FDD\u5B58\u4F4D\u7F6E\u3002","ConfigManagerSaveSuccess":"\u914D\u7F6E\u5DF2\u4FDD\u5B58\u3002","ConfigSaveTo":"\u914D\u7F6E\u4FDD\u5B58\u8DEF\u5F84","Debug":"\u5F00\u542FDebug\u6A21\u5F0F","Debug_Description":"\u542F\u7528Debug,\u6253\u5370\u66F4\u591A\u8C03\u8BD5\u4FE1\u606F\u3002\u5E76\u67E5\u770B\u4E00\u4E9B\u672A\u5B8C\u6210\u7684\u9690\u85CF\u529F\u80FD\u76F8\u5173\u8BBE\u7F6E\u3002","DisableLAN":"\u7981\u6B62\u5C40\u57DF\u7F51\u5171\u4EAB","DisableLAN_Description":"\u53EA\u5728\u672C\u673A\u63D0\u4F9B\u9605\u8BFB\u670D\u52A1\uFF0C\u4E0D\u5BF9\u5916\u5171\u4EAB","EnableDatabase":"\u542F\u7528\u6570\u636E\u5E93","EnableDatabase_Description":"\u542F\u7528\u672C\u5730\u6570\u636E\u5E93\uFF0C\u4FDD\u5B58\u626B\u63CF\u5230\u7684\u4E66\u7C4D\u6570\u636E\u3002","EnableFrpcServer":"EnableFrpcServer","RequiresLogin":"\u542F\u7528\u767B\u9646\u4FDD\u62A4","RequiresLogin_Description":"\u662F\u5426\u542F\u7528\u767B\u5F55\u3002","EnableTLS":"Enable TLS","EnableTLS_Description":"\u662F\u5426\u542F\u7528HTTPS\u534F\u8BAE\u3002\u9700\u8981\u8BBE\u7F6E\u8BC1\u4E66\u4E8Ekey\u6587\u4EF6\u3002","EnableUpload":"\u542F\u7528\u4E0A\u4F20\u529F\u80FD","EnableUpload_Description":"\u542F\u7528\u4E0A\u4F20\u529F\u80FD\u3002","EnablePlugin":"\u542F\u7528\u63D2\u4EF6\u7CFB\u7EDF","EnablePlugin_Description":"\u542F\u7528\u63D2\u4EF6\u7CFB\u7EDF\uFF0C\u5141\u8BB8\u5728\u9875\u9762\u4E2D\u63D2\u5165\u81EA\u5B9A\u4E49\u7684HTML\u3001CSS\u548CJavaScript\u4EE3\u7801(\u8BD5\u9A8C\u6027\u529F\u80FD)\u3002","plugins_config":"\u63D2\u4EF6\u7CFB\u7EDF","buildin_plugins":"\u5185\u7F6E\u63D2\u4EF6","enabled_plugins":"\u5DF2\u542F\u7528\u63D2\u4EF6","auto_flip_drag":"\u62D6\u52A8","auto_flip_drag_panel":"\u62D6\u52A8\u9762\u677F","auto_flip_start_flip":"\u5F00\u59CB\u81EA\u52A8\u7FFB\u9875","auto_flip_start_scroll":"\u5F00\u59CB\u81EA\u52A8\u6EDA\u52A8","auto_flip_pause_flip":"\u6682\u505C\u81EA\u52A8\u7FFB\u9875","auto_flip_pause_scroll":"\u6682\u505C\u81EA\u52A8\u6EDA\u52A8","auto_flip_interval":"\u95F4\u9694:","auto_flip_seconds":"\u79D2","auto_scroll_distance":"\u6EDA\u52A8\u8DDD\u79BB:","plugin_name_clock":"\u65F6\u949F\u63D2\u4EF6","plugin_name_auto_flip":"\u81EA\u52A8\u7FFB\u9875\u63D2\u4EF6","plugin_name_auto_scroll":"\u81EA\u52A8\u6EDA\u52A8\u63D2\u4EF6","plugin_name_comigo_xyz":"Comigo.xyz\u63D2\u4EF6","plugin_name_sample":"\u793A\u4F8B\u63D2\u4EF6","plugin_name_sketch_practice":"\u901F\u5199\u7EC3\u4E60\u63D2\u4EF6","sketch_practice_start":"\u5F00\u59CB\u901F\u5199\u7EC3\u4E60","sketch_practice_pause":"\u6682\u505C\u901F\u5199\u7EC3\u4E60","sketch_practice_countdown":"\u5012\u8BA1\u65F6","ExcludePath":"\u6392\u9664\u8DEF\u5F84","ExcludePath_Description":"\u626B\u63CF\u4E66\u7C4D\u7684\u65F6\u5019\uFF0C\u9700\u8981\u6392\u9664\u7684\u6587\u4EF6\u6216\u6587\u4EF6\u5939\u7684\u540D\u5B57","FrpClientConfig":"FrpClient\u8BBE\u7F6E","GenerateBookMetadata":"\u751F\u6210\u4E66\u7C4D\u5143\u6570\u636E","GenerateMetaData":"\u751F\u6210\u5143\u6570\u636E","GenerateMetaData_Description":"\u751F\u6210\u4E66\u7C4D\u5143\u6570\u636E\u3002\u5F53\u524D\u672A\u751F\u6548\u3002","HomeDirectory":"\u7528\u6237\u4E3B\u76EE\u5F55","Host":"\u57DF\u540D","Host_Description":"\u81EA\u5B9A\u4E49\u4E8C\u7EF4\u7801\u663E\u793A\u7684\u4E3B\u673A\u540D\u3002\u9ED8\u8BA4\u4E3A\u7F51\u5361IP\u3002","KeyFile":"KeyFile","KeyFile_Description":"TLS/SSL key\u6587\u4EF6\u8DEF\u5F84 (default: \\"~/.config/.comigo/key.key\\")","StoreUrls":"\u4E66\u5E93\u6587\u4EF6\u5939","StoreUrls_Description":"\u4E66\u5E93\u6587\u4EF6\u5939\uFF0C\u652F\u6301\u7EDD\u5BF9\u76EE\u5F55\u4E0E\u76F8\u5BF9\u76EE\u5F55\u3002\u76F8\u5BF9\u76EE\u5F55\u4EE5\u5F53\u524D\u6267\u884C\u76EE\u5F55\u4E3A\u57FA\u51C6\u3002<br>\u8BD5\u9A8C\u6027\u652F\u6301sftp\u4E0Esmb\u8FDC\u7A0B\u76EE\u5F55\u3002\u683C\u5F0F\u7C7B\u4F3C:<br>sftp://user:pass@192.168.1.1:22/some/path <br>smb://guest@192.168.1.1:445/some/path","LogFileName":"Log\u6587\u4EF6\u540D","LogFileName_Description":"Log\u6587\u4EF6\u540D","LogFilePath":"Log\u4FDD\u5B58\u4F4D\u7F6E","LogFilePath_Description":"Log\u6587\u4EF6\u7684\u4FDD\u5B58\u4F4D\u7F6E","LogToFile":"\u8BB0\u5F55Log\u5230\u672C\u5730","LogToFile_Description":"\u662F\u5426\u4FDD\u5B58\u7A0B\u5E8FLog\u5230\u672C\u5730\u6587\u4EF6\u3002\u9ED8\u8BA4\u4E0D\u4FDD\u5B58\u3002","MaxScanDepth":"\u6700\u5927\u626B\u63CF\u6DF1\u5EA6","MaxScanDepth_Description":"\u6700\u5927\u626B\u63CF\u6DF1\u5EA6\u3002\u8D85\u8FC7\u6DF1\u5EA6\u7684\u6587\u4EF6\u4E0D\u4F1A\u88AB\u626B\u63CF\u3002\u5F53\u524D\u6267\u884C\u76EE\u5F55\u4E3A\u57FA\u51C6\u3002","MinImageNum":"\u6700\u5C0F\u56FE\u7247\u6570","MinImageNum_Description":"\u538B\u7F29\u5305\u6216\u6587\u4EF6\u5939\u5185\u81F3\u5C11\u6709\u51E0\u5F20\u56FE\u7247\uFF0C\u624D\u7B97\u4F5C\u4E66\u7C4D\u3002","OpenBrowser":"\u6253\u5F00\u6D4F\u89C8\u5668","OpenBrowser_Description":"\u626B\u63CF\u5B8C\u6210\u540E\uFF0C\u662F\u5426\u540C\u65F6\u6253\u5F00\u6D4F\u89C8\u5668\u3002windows\u9ED8\u8BA4true\uFF0C\u5176\u4ED6\u5E73\u53F0\u9ED8\u8BA4false\u3002","Password":"\u5BC6\u7801","Password_Description":"\u767B\u5F55\u7528\u5BC6\u7801\u3002","ReEnterPassword":"\u91CD\u65B0\u8F93\u5165\u5BC6\u7801","ReEnterPassword_Description":"\u91CD\u65B0\u8F93\u5165\u5BC6\u7801\u3002","Port":"\u7AEF\u53E3","Port_Description":"\u7F51\u9875\u670D\u52A1\u7AEF\u53E3","PrintAllPossibleQRCode":"\u66F4\u591A\u4E8C\u7EF4\u7801","ProgramDirectory":"\u7A0B\u5E8F\u6240\u5728\u76EE\u5F55","StartFrpClientInBackground":"\u542F\u52A8FrpClient","SupportFileType":"\u652F\u6301\u7684\u538B\u7F29\u5305","SupportFileType_Description":"\u626B\u63CF\u6587\u4EF6\u65F6\uFF0C\u7528\u4E8E\u51B3\u5B9A\u8DF3\u8FC7\uFF0C\u8FD8\u662F\u7B97\u4F5C\u4E66\u7C4D\u5904\u7406\u7684\u6587\u4EF6\u540E\u7F00","SupportMediaType":"\u652F\u6301\u7684\u56FE\u7247\u6587\u4EF6","SupportMediaType_Description":"\u626B\u63CF\u538B\u7F29\u5305\u65F6\uFF0C\u7528\u4E8E\u7EDF\u8BA1\u56FE\u7247\u6570\u91CF\u7684\u56FE\u7247\u6587\u4EF6\u540E\u7F00","Timeout":"\u8FC7\u671F\u65F6\u95F4","TimeoutLimitForScan":"\u626B\u63CF\u8D85\u65F6","TimeoutLimitForScan_Description":"\u626B\u63CF\u6587\u4EF6\u65F6\uFF0C\u8D85\u8FC7\u51E0\u79D2\u949F\uFF0C\u5C31\u653E\u5F03\u626B\u63CF\u8FD9\u4E2A\u6587\u4EF6\uFF0C\u907F\u514D\u5361\u5728\u8FC7\u5927\u6587\u4EF6\u4E0A\u3002","Timeout_Description":"cookie\u8FC7\u671F\u65F6\u95F4\u3002\u5355\u4F4D\u4E3A\u5206\u949F\u3002","UseCache":"\u672C\u5730\u56FE\u7247\u7F13\u5B58","UseCache_Description":"\u542F\u7528\u672C\u5730\u56FE\u7247\u89E3\u538B\u7F13\u5B58\uFF0C\u9ED8\u8BA4\u7981\u7528\u3002","Username":"\u7528\u6237\u540D","Username_Description":"\u767B\u5F55\u7528\u7684\u7528\u6237\u540D\u3002","WorkingDirectory":"\u5F53\u524D\u5DE5\u4F5C\u76EE\u5F55","ZipFileTextEncoding":"\u975EUTF-8","ZipFileTextEncoding_Description":"\u975Eutf-8\u7F16\u7801ZIP\u6587\u4EF6\uFF0C\u5C1D\u8BD5\u7528\u4EC0\u4E48\u7F16\u7801\u89E3\u6790\u3002\u9ED8\u8BA4GBK\u3002","labs":"\u5B9E\u9A8C","all_page_num":"\u603B\u9875\u6570\uFF1A{0}","author":"\u4F5C\u8005\uFF1A{0}","auto_crop":"\u81EA\u52A8\u5207\u8FB9","auto_crop_num":"\u5207\u8FB9\u9608\u503C: ","auto_double_page":"\u81EA\u52A8\u5408\u5E76\u53CC\u9875\uFF08beta\uFF09","auto_hide_toolbar":"\u81EA\u52A8\u9690\u85CF\u5DE5\u5177\u680F","back-to-top":"\u8FD4\u56DE\u9876\u90E8","back_button":"\u8FD4\u56DE\u6309\u94AE","back_to_bookshelf":"\u8FD4\u56DE\u4E66\u67B6","book_shelf":"\u4E66\u5E93","child_book_hint":"\u6587\u4EF6\u5939\u5185\u6709{0}\u672C\u4E66","click_to_toggle":"(\u70B9\u51FB\u5207\u6362)","do_you_reset_local_settings":"\u662F\u5426\u8981\u91CD\u7F6E\u4E3A\u9ED8\u8BA4\u8BBE\u7F6E\uFF1F","double_page_mode":"\u53CC\u9875\u6A21\u5F0F","double_page_width":"\u6A2A\u5C4F\u53CC\u9875\u5BBD\u5EA6:","download_sample_config_file":"\u4E0B\u8F7D\u793A\u4F8B\u914D\u7F6E\u6587\u4EF6","download_windows_reg_file":"\u4E0B\u8F7D\u53F3\u952E\u6CE8\u518C\u6587\u4EF6","drop_to_upload":"\u70B9\u51FB\u6216\u5C06\u6587\u4EF6\u62D6\u52A8\u5230\u6B64\u533A\u57DF\u4EE5\u4E0A\u4F20","energy_threshold":"\u5207\u8FB9\u5F3A\u5EA6\uFF1A","epub_info":"ePub \u4FE1\u606F","exit_fullscreen":"\u9000\u51FA\u5168\u5C4F\u6A21\u5F0F","filesize":"\u5927\u5C0F\uFF1A{0}","flip_mode":"\u7FFB\u9875\u6A21\u5F0F","flip_odd_even_page":"\u66F4\u6539\u8DE8\u9875\u5339\u914D","flip_odd_even_page_hint":"\u5982\u679C\u8DE8\u9875\u5185\u5BB9\u4E0D\u5339\u914D\uFF0C\u53EF\u4EE5\u5C1D\u8BD5\u70B9\u51FB\u4FEE\u6B63","found_read_history":"\u53D1\u73B0\u672C\u5730\u9605\u8BFB\u8BB0\u5F55","from_interrupt":"\u4ECE\u4E2D\u65AD\u5904\u7EE7\u7EED","full_screen_hint":"\u5168\u5C4F\u6309\u94AE","fullscreen":"\u5207\u6362\u5168\u5C4F","good_job_and_bye":"\u505A\u5F97\u4E0D\u9519\uFF0C\u518D\u89C1~","gray_image":"\u9ED1\u767D\u5316","hint":"\u63D0\u793A","hint_first_page":"\u5F53\u524D\u662F\u7B2C\u4E00\u9875\uFF0C\u65E0\u6CD5\u5411\u524D\u7FFB\u9875","hint_last_page":"\u5F53\u524D\u662F\u6700\u540E\u4E00\u9875\uFF0C\u65E0\u6CD5\u5411\u540E\u7FFB\u9875","hour":"\u5C0F\u65F6","compress_image":"\u538B\u7F29\u56FE\u7247","infinite_dropdown":"\u65E0\u9650\u4E0B\u62C9\u6A21\u5F0F","interval":"\u95F4\u9694:","manga_mode":"\u65E5\u6F2B\uFF08\u53F3\u5F00\u672C\uFF09","load_all_pages":"\u52A0\u8F7D\u5168\u90E8\u9875\u9762","skip_and_load_full":"\u5DF2\u9605\u8BFB\u9875\u672A\u52A0\u8F7D\uFF08\u524D %d \u9875\uFF09\uFF0C\u70B9\u51FB\u4E0B\u9762\u6309\u94AE\u52A0\u8F7D\u5168\u90E8","completed_and_load_full":"\u5DF2\u7ECF\u9605\u8BFB\u5230\u6700\u540E\u4E00\u9875\uFF0C\u70B9\u51FB\u4E0B\u9762\u6309\u94AE\u52A0\u8F7D\u5168\u672C","load_from_interrupt":"\u662F\u5426\u4ECE\u7B2C %d \u9875\u5F00\u59CB\u52A0\u8F7D\uFF1F","login_success_hint":"\u767B\u5F55\u6210\u529F\uFF0C\u8FD4\u56DE\u4E0A\u4E00\u9875\u9762","logout":"\u9000\u51FA\u767B\u5F55","margin_bottom_on_scroll_mode":"\u9875\u9762\u95F4\u8DDD:","margin_on_scroll_mode":"\u9875\u9762\u95F4\u9699:","limit_width":"\u538B\u7F29\uFF08\u56FE\u7247\u5BBD\uFF09\uFF1A","minute":"\u5206","network":"\u7F51\u7EDC","no_book_found_hint":"\u672A\u627E\u5230\u4E66\u7C4D\uFF0C\u8BD5\u8BD5\u4E0A\u4F20\u6587\u4EF6\uFF1F","no_support_upload_file":"\u6587\u4EF6\u4E0A\u4F20\u529F\u80FD\u5DF2\u88AB\u7BA1\u7406\u5458\u5173\u95ED","not_support_fullscreen":"\u5F53\u524D\u6D4F\u89C8\u5668\u4E0D\u652F\u6301\u5168\u5C4F\u663E\u793A","now_is":"\u5F53\u524D:","number_of_online_books":"\u5728\u7EBF\u4E66\u7C4D\u6570\u91CF\uFF1A","original_image":"\u663E\u793A\u539F\u56FE","original_pdf_link":"\u67E5\u770B\u539F\u59CBPDF","page":"\u9875","page_turning_seconds":"\u7FFB\u9875\u95F4\u9694:","scroll_fixed_pagination":"\u56FA\u5B9A\u5206\u9875","scroll_infinite_scroll":"\u65E0\u9650\u4E0B\u62C9","pdf_hint_message":"\u652F\u6301\u7EAF\u56FE\u7247PDF\uFF0C\u5982\u679C\u52A0\u8F7D\u7F13\u6162\u6216\u51FA\u9519\uFF0C\u8BF7\u5C1D\u8BD5\uFF1A","please_enable_upload":"\u8BF7\u542F\u7528\u670D\u52A1\u5668\u7684\u4E0A\u4F20\u529F\u80FD","please_enter_content":"\u8BF7\u8F93\u5165\u5185\u5BB9","qrcode_hint":"\u626B\u7801\u9605\u8BFB\uFF0C\u70B9\u51FB\u663E\u793A\u4E8C\u7EF4\u7801","raw_resolution":"\u539F\u59CB\u5206\u8FA8\u7387","re_sort_book":"\u91CD\u65B0\u6392\u5217\u4E66\u7C4D","re_sort_page":"\u91CD\u65B0\u6392\u5E8F\u9875\u9762","reader_settings":"\u9605\u8BFB\u8BBE\u7F6E","reading_progress_bar":"\u9605\u8BFB\u8FDB\u5EA6\u6761","refresh_page":"\u5237\u65B0\u9875\u9762","reset_local_settings":"\u91CD\u7F6E\u672C\u5730\u8BBE\u7F6E","resort_file":"\u91CD\u65B0\u6392\u5E8F\u6587\u4EF6","comic_mode":"\u7F8E\u6F2B\uFF08\u5DE6\u5F00\u672C\uFF09","save_page_num":"\u4FDD\u5B58\u9605\u8BFB\u8FDB\u5EA6","scan_qrcode":"\u626B\u7801\u9605\u8BFB\uFF1A","scanned_hint":"\u626B\u63CF\u5230XX\u672C\u4E66\uFF0C\u7ACB\u5373\u67E5\u770B\uFF1F","scroll_mode":"\u5377\u8F74\u6A21\u5F0F","second":"\u79D2","select_language":"\u9009\u62E9\u8BED\u8A00","server_config":"\u670D\u52A1\u5668\u8BBE\u7F6E","server_setting":"Comigo \u670D\u52A1\u5668\u8BBE\u7F6E","set_back_color":"\u80CC\u666F\u989C\u8272:","set_interface_color":"\u754C\u9762\u989C\u8272:","show_filename":"\u663E\u793A\u6587\u4EF6\u540D","show_file_icon":"\u663E\u793A\u6587\u4EF6\u56FE\u6807","show_header":"\u663E\u793A\u6807\u9898","showPageNum":"\u663E\u793A\u9875\u7801","simplify_filename":"\u7B80\u5316\u6587\u4EF6\u540D","single_page_mode":"\u5355\u9875\u6A21\u5F0F","single_page_width":"\u6A2A\u5C4F\u5355\u9875\u5BBD\u5EA6:","sort_by_default":"\u4FDD\u6301\u9ED8\u8BA4\u987A\u5E8F","sort_by_filename":"\u6309\u6587\u4EF6\u540D\u6392\u5E8F (A-Z)","sort_by_filename_reverse":"\u6309\u6587\u4EF6\u540D\u6392\u5E8F (Z-A)","sort_by_filesize":"\u6309\u6587\u4EF6\u5927\u5C0F\u6392\u5E8F (\u4ECE\u5927\u5230\u5C0F)","sort_by_filesize_reverse":"\u6309\u6587\u4EF6\u5927\u5C0F\u6392\u5E8F (\u4ECE\u5C0F\u5230\u5927)","sort_by_modify_time":"\u6309\u4FEE\u6539\u65F6\u95F4\u6392\u5E8F (\u4ECE\u65B0\u5230\u65E7)","sort_by_modify_time_reverse":"\u6309\u4FEE\u6539\u65F6\u95F4\u6392\u5E8F (\u4ECE\u65E7\u5230\u65B0)","sort_reverse":"\uFF08\u53CD\u5411\uFF09","start_sketch_message":"\u5012\u8BA1\u65F6\u901F\u5199\u5DF2\u5F00\u59CB\uFF0C\u795D\u4F60\u5FC3\u60C5\u6109\u5FEB\u3002","start_sketch_mode":"\u5F00\u59CB\u901F\u5199","starting_from_beginning":"\u4ECE\u5934\u5F00\u59CB","starting_from_beginning_hint":"\u4ECE\u5934\u5F00\u59CB\u52A0\u8F7D","stop_sketch_mode":"\u505C\u6B62\u901F\u5199","submit":"\u63D0\u4EA4","add":"\u6DFB\u52A0","success_fullscreen":"\u5DF2\u8FDB\u5165\u5168\u5C4F\u6A21\u5F0F","successfully_loaded_reading_progress":"\u6210\u529F\u52A0\u8F7D\u9605\u8BFB\u8FDB\u5EA6","sync_page":"\u8FDC\u7A0B\u540C\u6B65\u7FFB\u9875","temp_future_hint":"\u653E\u4E00\u4E9B\u8FD8\u672A\u5B8C\u6210\u7684\u529F\u80FD\uFF0C\u5F00\u53D1\u4E0E\u8C03\u6574\u4E2D\u3002","test":"\u6D4B\u8BD5","total_is":"\u5B8C\u6210:","total_time":"\u603B\u65F6\u95F4:","type_or_paste_content":"\u952E\u5165\u6216\u7C98\u8D34\u5185\u5BB9","upload_file":"\u4E0A\u4F20\u6587\u4EF6","width_use_fixed_value":"\u6A2A\u5C4F\u5BBD\u5EA6: \u56FA\u5B9A\u503Cpx","width_use_percent":"\u6A2A\u5C4F\u5BBD\u5EA6: \u767E\u5206\u6BD4%","portrait_width_percent":"\u7AD6\u5C4F\u5BBD\u5EA6(\u767E\u5206\u6BD4)","auto_align":"\u81EA\u52A8\u5BF9\u9F50\u753B\u9762","swipe_turn":"\u89E6\u6478\u6ED1\u52A8\u7FFB\u9875","login_title":"\u767B\u5F55Comigo","login_subtitle":"\u8BF7\u8F93\u5165\u60A8\u7684\u8D26\u53F7\u548C\u5BC6\u7801","login_failed":"\u767B\u5F55\u5931\u8D25\uFF0C\u8BF7\u68C0\u67E5\u7528\u6237\u540D\u548C\u5BC6\u7801","login_error_teapot":"\u670D\u52A1\u5668\u5F53\u524D\u4E0D\u9700\u8981\u8BA4\u8BC1\uFF0C\u8BF7\u76F4\u63A5\u8BBF\u95EE<a class=\\"font-semibold text-blue-600\\" href=\\"/\\">\u9996\u9875</a>","logging_in":"\u767B\u5F55\u4E2D...","login":"\u767B\u5F55","other_information":"\u5176\u4ED6\u4FE1\u606F","login_forgot_password_hint":"\u5FD8\u8BB0\u5BC6\u7801\uFF1F\u8BF7\u8054\u7CFB\u7CFB\u7EDF\u7BA1\u7406\u5458","no_pattern":"\u7EAF\u8272","grid_line":"\u7F51\u683C\u7EBF","grid_point":"\u7F51\u683C\u70B9","mosaic":"\u9A6C\u8D5B\u514B","open_pdf_in_browser":"\u7528\u6D4F\u89C8\u5668\u6253\u5F00PDF","StaticFileMode":"\u9759\u6001\u6587\u4EF6\u6A21\u5F0F","StaticFileMode_Description":"\u662F\u5426\u5F00\u542F\u9759\u6001\u6587\u4EF6\u6A21\u5F0F\u3002\u9759\u6001\u6A21\u5F0F\u4E0B\uFF0C\u6240\u6709\u56FE\u7247\u4E0E\u811A\u672C\u90FD\u6253\u5305\u5230html\u6587\u4EF6\u91CC\uFF0C\u53EF\u4EE5\u76F4\u63A5\u53E6\u5B58\u4E3A\u5355\u4E2A\u7F51\u9875\u6587\u4EF6\uFF08\u5F00\u53D1\u4E2D\uFF09\u3002","confirm_logout":"\u786E\u8BA4\u8981\u9000\u51FA\u767B\u5F55\u5417\uFF1F","confirm_reset_settings":"\u786E\u8BA4\u8981\u91CD\u7F6E\u672C\u5730\u8BBE\u7F6E\u5417\uFF1F","current_dir_scope":"\u5728\u5F53\u524D\u76EE\u5F55\u8FD0\u884C\u65F6\uFF08\u5C40\u90E8\u6709\u6548\uFF09","current_user_scope":"\u5F53\u524D\u767B\u5F55\u7528\u6237\u6709\u6548\uFF08\u5168\u5C40\u6709\u6548\uFF09","portable_binary_scope":"\u6B64\u4E8C\u8FDB\u5236\u6587\u4EF6\u6709\u6548\uFF08\u4FBF\u643A\u6A21\u5F0F\uFF09","saveSuccessHint":"\u8BBE\u7F6E\u4FEE\u6539\u6210\u529F\uFF0C2\u79D2\u540E\u81EA\u52A8\u5237\u65B0\u9875\u9762\u3002","port_change_hint":"\u7AEF\u53E3\u5DF2\u66F4\u6539\uFF0C2\u79D2\u540E\u81EA\u52A8\u8DF3\u8F6C\u5230\u65B0\u7AEF\u53E3\u3002","no_books_library_path_notice":"\u672A\u68C0\u6D4B\u5230\u53EF\u9605\u8BFB\u7684\u4E66\u7C4D\uFF0C\u8BF7\u5148\u914D\u7F6E\u4E66\u5E93\u8DEF\u5F84\u3002\u914D\u7F6E\u5B8C\u6210\u540E\uFF0C\u9875\u9762\u5C06\u81EA\u52A8\u5237\u65B0\u3002","download_raw_archive":"\u539F\u59CB\u6587\u4EF6","download_portable_web_file":"\u4FBF\u643A\u7F51\u9875","download_as_zip":"ZIP\u6587\u4EF6","download_as_epub":"EPUB\u6587\u4EF6","clear_cache":"\u6E05\u9664\u7F13\u5B58","confirm_clear_cache":"\u786E\u5B9A\u8981\u6E05\u9664\u6B64\u4E66\u7C4D\u7684\u7F13\u5B58\u6587\u4EF6\u5417\uFF1F\uFF08\u5305\u62EC\u5143\u6570\u636E\u548C\u5C01\u9762\u7F13\u5B58\uFF09","clear_cache_success":"\u7F13\u5B58\u5DF2\u6E05\u9664","clear_cache_failed":"\u6E05\u9664\u7F13\u5B58\u5931\u8D25","download":"\u4E0B\u8F7D","switch":"\u5207\u6362","EnableTailscale":"Tailscale\u5185\u7F51\u7A7F\u900F","EnableTailscale_Description":"Tailscale\u5185\u7F51\u7A7F\u900F\u8BBE\u7F6E\uFF0C\u9996\u6B21\u8FDE\u63A5\uFF0C\u9700\u8981\u5728Tailscale\u63A7\u5236\u53F0\u9A8C\u8BC1\u3002","TailscaleHostname":"Tailscale\u4E3B\u673A\u540D","TailscaleHostname_Description":"Tailscale\u4E3B\u673A\u540D\u90E8\u5206\uFF0C\u5B8C\u6574\u57DF\u540D\u7C7B\u4F3C {hostname}.example.ts.net","TailscalePort":"Tailscale\u76D1\u542C\u7AEF\u53E3","TailscalePort_Description":"Tailscale\u76D1\u542C\u7AEF\u53E3\u3002\u9ED8\u8BA4443\uFF0C\u81EA\u52A8\u542F\u7528TLS\u3002","FunnelTunnel":"Funnel\u96A7\u9053","FunnelTunnel_Description":"Funnel\u96A7\u9053\uFF08\u516C\u7F51\u8BBF\u95EE\uFF09\u3002\u5982\u679C\u4F60\u4E0D\u60F3\u8981\u5BF9\u5916\u516C\u5F00\uFF0C\u5EFA\u8BAE\u8BBE\u7F6E\u5BC6\u7801\u4FDD\u62A4\u3002Funnel\u96A7\u9053\u53EA\u652F\u6301443, 8443, 10000\u7AEF\u53E3\u3002","read_only_mode":"\u53EA\u8BFB\u6A21\u5F0F","read_only_mode_description":"\u5F53\u524D\u5904\u4E8E\u53EA\u8BFB\u6A21\u5F0F\uFF0C\u65E0\u6CD5\u5728\u7F51\u9875\u7AEF\u66F4\u6539\u8BBE\u5B9A\u6216\u4E0A\u4F20\u6587\u4EF6\u3002","enable_single_instance":"\u542F\u7528\u5355\u5B9E\u4F8B\u6A21\u5F0F\uFF0C\u786E\u4FDD\u540C\u4E00\u65F6\u95F4\u53EA\u6709\u4E00\u4E2A\u7A0B\u5E8F\u5B9E\u4F8B\u8FD0\u884C","tailscale_auth_url_is":"\u8981\u542F\u52A8 Tailscale \u670D\u52A1\u5668\uFF0C\u8BF7\u8BBE\u7F6E TS_AUTHKEY \u540E\u91CD\u542F\uFF0C\u6216\u8BBF\u95EE\u8BA4\u8BC1\u94FE\u63A5\uFF1A","tailscale_server_start":"Tailscale \u5DF2\u542F\u52A8\uFF0C\u6B63\u5728\u83B7\u53D6 IP \u5730\u5740\u2026\u2026","tailscale_reading_url":"\u901A\u8FC7 Tailscale \u8BBF\u95EE\u7684\u9605\u8BFB\u94FE\u63A5\u4E3A\uFF1A","tailscale_not_connected_hint":"Tailscale \u672A\u8FDE\u63A5\uFF0C\u8BF7\u68C0\u67E5\u7F51\u7EDC\u6216\u8BA4\u8BC1\u72B6\u6001\u3002","tailscale_not_enabled":"Tailscale \u672A\u542F\u7528\uFF0C\u8BF7\u5148\u542F\u7528\u3002","ServerSettings":"Comigo\u670D\u52A1\u5668\u8BBE\u7F6E","settings_stores":"\u4E66\u5E93\u8BBE\u7F6E","settings_network":"\u7F51\u7EDC\u8BBE\u7F6E","settings_extra":"\u5B9E\u9A8C\u529F\u80FD","settings_log_title":"\u670D\u52A1\u5668\u5B9E\u65F6\u65E5\u5FD7","settings_log_broadcast_placeholder":"\u8F93\u5165\u5E7F\u64AD\u6D88\u606F","settings_log_broadcast_send":"\u53D1\u9001","settings_log_sse_connected":"\u65E5\u5FD7\u670D\u52A1\u5668\u5DF2\u8FDE\u63A5","settings_log_sse_closed":"\u8FDE\u63A5\u5DF2\u5173\u95ED","settings_log_sse_retrying":"\u91CD\u8BD5\u4E2D...","settings_log_broadcast_default_message":"Hello ~","remote_access":"\u8FDC\u7A0B\u8BBF\u95EE","ErrPasswordMismatch":"\u4E24\u6B21\u8F93\u5165\u7684\u5BC6\u7801\u4E0D\u4E00\u81F4\uFF0C\u8BF7\u91CD\u65B0\u8F93\u5165","PromptSetPassword":"\u8BF7\u8F93\u5165\u5BC6\u7801","MsgLoginSettingsUpdated":"\u767B\u5F55\u8BBE\u5B9A\u4FEE\u6539\u6210\u529F\uFF0C\u81EA\u52A8\u8DF3\u8F6C\u767B\u5F55\u9875\u9762","CurrentPassword":"\u5F53\u524D\u5BC6\u7801","AdminAccountSetup":"\u7BA1\u7406\u8D26\u53F7\u5BC6\u7801","AdminAccountSetupDescription":"\u8BF7\u8BBE\u7F6E\u7BA1\u7406\u5458\u8D26\u53F7\u5BC6\u7801\uFF0C\u7528\u4E8E\u767B\u5F55Comigo\u3002\u8BBE\u7F6E\u540E\uFF0C\u8BBF\u95EE\u670D\u52A1\u5C06\u9700\u8981\u767B\u5F55\u3002","ConfigStorageLocationPrompt":"\u8BF7\u9009\u62E9\u914D\u7F6E\u6587\u4EF6\u5B58\u50A8\u7684\u4F4D\u7F6E:","set_account_password":"\u8BBE\u7F6E\u8D26\u53F7\u5BC6\u7801","connect_tailscale":"\u8FDE\u63A5Tailscale\u7F51\u7EDC","disconnect_tailscale":"\u65AD\u5F00Tailscale\u8FDE\u63A5","tailscale_status":"Tailscale\u72B6\u6001","service_status":"\u670D\u52A1\u72B6\u6001","running":"\u8FD0\u884C\u4E2D","client_count":"\u5BA2\u6237\u7AEF\u6570","host_system":"\u5BBF\u4E3B\u7CFB\u7EDF","connection_status":"\u8FDE\u63A5\u72B6\u51B5","connected":"\u5DF2\u8FDE\u63A5","not_connected":"\u672A\u8FDE\u63A5","enable_funnel":"Funnel\u6A21\u5F0F","enable_funnel_public_access":"Funnel\u6A21\u5F0F\uFF08\u516C\u7F51\u8BBF\u95EE\uFF09","disabled":"\u672A\u542F\u7528","ip_address":"IP\u5730\u5740","service_version":"\u670D\u52A1\u7248\u672C","read_link":"\u9605\u8BFB\u94FE\u63A5","submiting":"\u63D0\u4EA4\u4E2D...","enable":"\u542F\u7528","disable":"\u7981\u7528","TailscaleAuthKey":"Tailscale\u9884\u6388\u6743\u5BC6\u94A5","TailscaleAuthKeyDescription":"Tailscale\u9884\u6388\u6743\u5BC6\u94A5\uFF08TS_AUTHKEY\uFF09\uFF0C\u7528\u4E8E\u5728\u65E0\u6D4F\u89C8\u5668\u73AF\u5883\u81EA\u52A8\u8BA4\u8BC1\u3002","verify_tailscale":"\u8BF7\u70B9\u51FB\u94FE\u63A5\uFF0C\u9A8C\u8BC1Tailscale:","funnel_status":"Funnel\u72B6\u6001","funnel_tunnel":"Funnel\u516C\u7F51\u96A7\u9053","funnel_setup_done":"Funnel\u8BBE\u7F6E\u5B8C\u6BD5\uFF0C\u53EF\u5F00\u542F\u516C\u7F51\u96A7\u9053","funnel_setup_not_done":"Funnel\u516C\u7F51\u96A7\u9053\u9700\u8981\u8BBE\u7F6E","funnel_not_set_hint":"\u542F\u7528Funnel\u516C\u7F51\u96A7\u9053\u3002\u9700\u8981\uFF1A","funnel_require_dns_1":"\u5728Tailscale\u63A7\u5236\u53F0DNS\u9762\u677F","funnel_require_dns_2":"\u5F00\u542FMagicDNS\u4E0EHTTPS\u529F\u80FD\u3002","funnel_require_acl_1":"\u5728Tailscale\u63A7\u5236\u53F0ACL\u9762\u677F","funnel_require_acl_2":"\u7F16\u8F91ACL\u89C4\u5219\uFF0C\u542F\u7528Funnel\u6743\u9650","funnel_require_acl_3":"\uFF08\u70B9\u6B64\u4E0B\u8F7D\u793A\u4F8BJSON\uFF09\u3002","funnel_require_password_1":"\u542F\u7528\u3010Funnel\u5BC6\u7801\u68C0\u67E5\u3011\u65F6\uFF0C\u9700\u8981\u8BBE\u7F6EComigo\u7BA1\u7406\u5458\u8D26\u6237\u4E0E\u5BC6\u7801\u624D\u80FD\u4F7F\u7528Funnel\u96A7\u9053","verify_link":"\u9A8C\u8BC1\u94FE\u63A5","copy_link":"\u590D\u5236\u94FE\u63A5","copy_success":"\u5DF2\u590D\u5236\u5230\u526A\u8D34\u677F","copy_failed":"\u590D\u5236\u5931\u8D25\uFF0C\u8BF7\u624B\u52A8\u590D\u5236","FunnelLoginCheck":"Funnel\u5BC6\u7801\u68C0\u67E5","FunnelLoginCheckDescription":"\u542F\u7528Funnel\u516C\u7F51\u96A7\u9053\u524D\uFF0C\u68C0\u67E5\u5F53\u524D\u662F\u5426\u5DF2\u8BBE\u7F6E\u767B\u5F55\u4FDD\u62A4\u3002","funnel_login_check_enabled_but_no_password":"\u3010Funnel\u5BC6\u7801\u68C0\u67E5\u3011\u5DF2\u5F00\u542F\uFF0C\u4E14\u672A\u8BBE\u7F6E\u767B\u5F55\u8D26\u53F7\u5BC6\u7801\uFF0C\u65E0\u6CD5\u542F\u7528Funnel\u96A7\u9053\u3002","tailscale_settings_submitted_check_status":"Tailscale\u8BBE\u7F6E\u5DF2\u63D0\u4EA4\uFF0C\u8BF7\u67E5\u770BTailscale\u72B6\u6001","value_already_exists_do_not_add_again":"\u8BE5\u503C\u5DF2\u5B58\u5728\uFF0C\u8BF7\u52FF\u91CD\u590D\u6DFB\u52A0","file_uploaded_successfully":"\u6587\u4EF6\u6210\u529F\u4E0A\u4F20","content_empty_please_enter_before_submit":"\u5185\u5BB9\u4E3A\u7A7A\uFF0C\u8BF7\u8F93\u5165\u5185\u5BB9\u540E\u63D0\u4EA4","default_prompt_message":"\u9ED8\u8BA4\u63D0\u793A\u4FE1\u606F","confirm":"\u786E\u5B9A","ok":"\u786E\u5B9A","cancel":"\u53D6\u6D88","uploading":"\u4E0A\u4F20\u4E2D...","upload_failed_network_error":"\u4E0A\u4F20\u5931\u8D25: \u7F51\u7EDC\u9519\u8BEF","drag_or_click_to_upload":"\u5C06\u6587\u4EF6\u62D6\u62FD\u5230\u6B64\u5904\u6216\u70B9\u51FB\u9009\u62E9\u6587\u4EF6","selected_file":"\u9009\u62E9\u7684\u6587\u4EF6","select_upload_target_store":"\u9009\u62E9\u4E0A\u4F20\u76EE\u6807\u4E66\u5E93","store_not_exists":"\u4E66\u5E93\u8DEF\u5F84\u4E0D\u5B58\u5728","no_available_stores":"\u6CA1\u6709\u53EF\u7528\u7684\u4E66\u5E93\uFF0C\u8BF7\u5148\u5728\u8BBE\u7F6E\u4E2D\u6DFB\u52A0\u4E66\u5E93\u8DEF\u5F84","store_validation_failed":"\u65E0\u6548\u7684\u4E66\u5E93\u8DEF\u5F84","passwords_not_match":"\u4E24\u6B21\u8F93\u5165\u7684\u5BC6\u7801\u4E0D\u4E00\u81F4","please_delete_other_config_first":"\u8BF7\u5148\u5220\u9664\u5176\u4ED6\u4F4D\u7F6E\u7684\u914D\u7F6E\u6587\u4EF6","save_config_success":"\u4FDD\u5B58\u8BBE\u7F6E\u6587\u4EF6\u6210\u529F\uFF01","save_config_failed":"\u4FDD\u5B58\u8BBE\u7F6E\u6587\u4EF6\u5931\u8D25","no_config_file_to_delete_in_path":"\u5F53\u524D\u9009\u62E9\u7684\u8DEF\u5F84\u4E0B\uFF0C\u6CA1\u6709\u53EF\u5220\u9664\u7684\u914D\u7F6E\u6587\u4EF6","delete_config_success":"\u5220\u9664\u8BBE\u7F6E\u6587\u4EF6\u6210\u529F","delete_config_failed":"\u5220\u9664\u8BBE\u7F6E\u6587\u4EF6\u5931\u8D25","err_invalid_number":"\u8BF7\u8F93\u5165\u6709\u6548\u7684\u6570\u5B57","err_number_range":"\u6570\u503C\u5FC5\u987B\u5728 {0} \u5230 {1} \u4E4B\u95F4","err_update_config_failed":"\u66F4\u65B0\u914D\u7F6E\u5931\u8D25","err_network_error":"\u7F51\u7EDC\u9519\u8BEF\uFF0C\u8BF7\u91CD\u8BD5","err_update_login_settings_failed":"\u66F4\u65B0\u767B\u5F55\u8BBE\u7F6E\u5931\u8D25","err_add_config_failed":"\u6DFB\u52A0\u914D\u7F6E\u5931\u8D25","err_delete_config_failed":"\u5220\u9664\u914D\u7F6E\u5931\u8D25","err_save_config_failed":"\u4FDD\u5B58\u914D\u7F6E\u5931\u8D25","err_config_locked":"\u914D\u7F6E\u5DF2\u9501\u5B9A\uFF0C\u65E0\u6CD5\u4FEE\u6539","read":"\u5DF2\u8BFB","unread":"\u672A\u8BFB","mark_all_as_read":"\u5168\u90E8\u6807\u8BB0\u4E3A\u5DF2\u8BFB","new_version_available":"\u6709\u65B0\u7248\u672C\u53EF\u7528\uFF01","download_now":"\u7ACB\u5373\u4E0B\u8F7D","infinite_scroll_button":"\u65E0\u9650\u5377\u8F74","paged_scroll_button":"\u5206\u9875\u5377\u8F74","page_flip_button":"\u7FFB\u9875\u6A21\u5F0F","infinite_scroll":"\u65E0\u9650\u5377\u8F74(\u52A0\u8F7D\u5168\u90E8)","paged_scroll":"\u5206\u9875\u5377\u8F74(\u52A0\u8F7D32\u9875)","page_flip":"\u7FFB\u9875\u6A21\u5F0F(\u5DE6\u53F3\u6ED1\u52A8)","reading_progress_percent":"\u9605\u8BFB\u8FDB\u5EA6\uFF08\u767E\u5206\u6BD4\uFF09","reading_progress_page":"\u9605\u8BFB\u8FDB\u5EA6\uFF08\u9875\u6570\uFF09","reading_history":"\u9605\u8BFB\u5386\u53F2","view_all_reading_history":"\u67E5\u770B\u5168\u90E8\u9605\u8BFB\u5386\u53F2","loading":"\u52A0\u8F7D\u4E2D...","no_reading_history":"\u6682\u65E0\u9605\u8BFB\u5386\u53F2","unknown":"\u672A\u77E5","systray_tooltip":"Comigo \u6F2B\u753B\u9605\u8BFB\u5668","systray_open_browser":"\u6253\u5F00\u6D4F\u89C8\u5668","systray_open_browser_tooltip":"\u5728\u6D4F\u89C8\u5668\u4E2D\u6253\u5F00 Comigo","systray_copy_url":"\u590D\u5236\u9605\u8BFB\u5730\u5740","systray_copy_url_tooltip":"\u590D\u5236\u9605\u8BFB\u5730\u5740\u5230\u526A\u8D34\u677F","systray_enable_tailscale":"\u542F\u7528Tailscale","systray_disable_tailscale":"\u7981\u7528Tailscale","systray_toggle_tailscale_tooltip":"\u5207\u6362Tailscale\u72B6\u6001","systray_language":"\u8BED\u8A00\u5207\u6362","systray_language_tooltip":"\u5207\u6362\u754C\u9762\u8BED\u8A00","systray_language_zh":"\u4E2D\u6587","systray_language_zh_tooltip":"\u5207\u6362\u5230\u4E2D\u6587","systray_language_en":"English","systray_language_en_tooltip":"Switch to English","systray_language_ja":"\u65E5\u672C\u8A9E","systray_language_ja_tooltip":"\u65E5\u672C\u8A9E\u306B\u5207\u308A\u66FF\u3048","systray_open_directory":"\u6253\u5F00\u76EE\u5F55","systray_open_directory_tooltip":"\u6253\u5F00\u76F8\u5173\u76EE\u5F55","systray_config_directory":"\u914D\u7F6E\u6587\u4EF6\u76EE\u5F55","systray_config_directory_tooltip":"\u6253\u5F00\u914D\u7F6E\u6587\u4EF6\u76EE\u5F55","systray_store":"\u4E66\u5E93","systray_store_folder":"\u4E66\u5E93\u6587\u4EF6\u5939","systray_refresh_directories":"\u5237\u65B0","systray_refresh_directories_tooltip":"\u5237\u65B0\u76EE\u5F55\u5217\u8868","systray_quit":"\u9000\u51FA","systray_quit_tooltip":"\u9000\u51FA Comigo","context_menu_open_with_comigo":"\u4F7F\u7528Comigo\u6253\u5F00","register_context_menu":"\u6CE8\u518C\u53F3\u952E\u83DC\u5355\uFF08\u201C\u4F7F\u7528Comigo\u6253\u5F00\u201D\uFF09","unregister_context_menu":"\u6E05\u7406\u53F3\u952E\u83DC\u5355\uFF08\u201C\u4F7F\u7528Comigo\u6253\u5F00\u201D\uFF09","systray_extra":"\u5176\u4ED6","systray_extra_tooltip":"\u66F4\u591A\u7CFB\u7EDF\u96C6\u6210\u529F\u80FD","systray_project":"Comigo \u9879\u76EE\u5730\u5740","systray_project_tooltip":"\u6253\u5F00 Comigo \u5728 GitHub \u4E0A\u7684\u9879\u76EE\u4E3B\u9875","register_folder_context_menu":"\u6CE8\u518C\u6587\u4EF6\u5939\u53F3\u952E\u83DC\u5355\uFF08\u201C\u4F7F\u7528Comigo\u6253\u5F00\u201D\uFF09","unregister_folder_context_menu":"\u6E05\u7406\u6587\u4EF6\u5939\u53F3\u952E\u83DC\u5355\uFF08\u201C\u4F7F\u7528Comigo\u6253\u5F00\u201D\uFF09","register_file_association":"\u6CE8\u518C\u538B\u7F29\u6587\u4EF6\u7C7B\u578B\u5173\u8054\uFF08\u4F5C\u4E3A\u5019\u9009\u6253\u5F00\u65B9\u5F0F\uFF09","unregister_file_association":"\u6E05\u7406\u538B\u7F29\u6587\u4EF6\u7C7B\u578B\u5173\u8054\uFF08\u79FB\u9664Comigo\u5019\u9009\u9879\uFF09","create_desktop_shortcut":"\u5728\u684C\u9762\u521B\u5EFA\u5FEB\u6377\u65B9\u5F0F","log_language_changed_to_chinese":"\u8BED\u8A00\u5DF2\u5207\u6362\u4E3A\u4E2D\u6587","log_language_changed_to_english":"\u8BED\u8A00\u5DF2\u5207\u6362\u4E3A\u82F1\u6587","log_language_changed_to_japanese":"\u8BED\u8A00\u5DF2\u5207\u6362\u4E3A\u65E5\u6587","log_requesting_quit_from_systray":"\u4ECE\u7CFB\u7EDF\u6258\u76D8\u8BF7\u6C42\u9000\u51FA","log_failed_to_read_icon_file":"\u8BFB\u53D6\u56FE\u6807\u6587\u4EF6\u5931\u8D25: %v\uFF0C\u4F7F\u7528\u9ED8\u8BA4\u56FE\u6807","log_opening_browser":"\u6B63\u5728\u6253\u5F00\u6D4F\u89C8\u5668: %s","log_failed_to_copy_url":"\u590D\u5236URL\u5230\u526A\u8D34\u677F\u5931\u8D25: %v","log_copied_url_to_clipboard":"\u5DF2\u590D\u5236URL\u5230\u526A\u8D34\u677F: %s","log_failed_to_toggle_tailscale":"\u5207\u6362Tailscale\u72B6\u6001\u5931\u8D25: %v","log_failed_to_set_language":"\u8BBE\u7F6E\u8BED\u8A00\u5931\u8D25: %v","log_failed_to_get_config_dir":"\u83B7\u53D6\u914D\u7F6E\u76EE\u5F55\u5931\u8D25: %v","log_failed_to_clear_folder_context_menu":"\u6E05\u9664Windows\u6587\u4EF6\u5939\u53F3\u952E\u83DC\u5355\u5931\u8D25: %v","log_failed_to_register_folder_context_menu":"\u6CE8\u518CWindows\u6587\u4EF6\u5939\u53F3\u952E\u83DC\u5355\u5931\u8D25: %v","log_failed_to_create_desktop_shortcut":"\u521B\u5EFA\u684C\u9762\u5FEB\u6377\u65B9\u5F0F\u5931\u8D25: %v","log_failed_to_unregister_archive_handler":"\u53D6\u6D88\u6CE8\u518C\u538B\u7F29\u6587\u4EF6\u5904\u7406\u5668\u5931\u8D25: %v","log_failed_to_register_archive_handler":"\u6CE8\u518C\u538B\u7F29\u6587\u4EF6\u5904\u7406\u5668\u5931\u8D25: %v","log_opening_comigo_project_page":"\u6B63\u5728\u6253\u5F00Comigo\u9879\u76EE\u9875\u9762: https://github.com/yumenaka/comigo","log_failed_to_open_directory":"\u6253\u5F00\u76EE\u5F55\u5931\u8D25: %v","log_upload_file_count":"\u4E0A\u4F20\u6587\u4EF6\u6570\u91CF: %d","log_upload_path_not_set":"\u4E0A\u4F20\u8DEF\u5F84\u672A\u8BBE\u7F6E","log_mkdir_failed":"\u521B\u5EFA\u76EE\u5F55\u5931\u8D25: %s","log_mkdir_upload_folder_success":"\u521B\u5EFA\u4E0A\u4F20\u6587\u4EF6\u5939\u6210\u529F: %s","log_file_upload_success":"\u6587\u4EF6\u4E0A\u4F20\u6210\u529F: %s","log_no_changes_skipped_scan":"\u914D\u7F6E\u65E0\u53D8\u5316\uFF0C\u8DF3\u8FC7\u626B\u63CF\u4E66\u5E93\u8DEF\u5F84","log_no_changes_skipped_rescan":"\u914D\u7F6E\u65E0\u53D8\u5316\uFF0C\u8DF3\u8FC7\u91CD\u65B0\u626B\u63CF\u76EE\u5F55","log_tailscale_config_changed_restart":"Tailscale\u914D\u7F6E\u5DF2\u66F4\u6539\uFF0C\u5C06\u91CD\u542FTailscale\u670D\u52A1\u5668","log_tailscale_enabled_start":"Tailscale\u5DF2\u542F\u7528\uFF0C\u5C06\u542F\u52A8Tailscale\u670D\u52A1\u5668","log_tailscale_disabled_stop":"Tailscale\u5DF2\u7981\u7528\uFF0C\u5C06\u505C\u6B62Tailscale\u670D\u52A1\u5668","log_scan_subdirectory_error":"\u626B\u63CF\u5B50\u76EE\u5F55\u51FA\u9519: %v","log_get_file_info_failed":"\u83B7\u53D6\u6587\u4EF6\u4FE1\u606F\u5931\u8D25: %v","log_tailscale_disabled_skip_qrcode":"Tailscale\u5DF2\u7981\u7528\uFF0C\u8DF3\u8FC7\u663E\u793A\u4E8C\u7EF4\u7801\u529F\u80FD","log_tailscale_status_check_exceeded":"Tailscale\u72B6\u6001\u68C0\u67E5\u6B21\u6570\u8D85\u9650\uFF0C\u505C\u6B62\u8FDB\u4E00\u6B65\u68C0\u67E5","log_tailscale_not_yet_fqdn":"Tailscale FQDN\u5C1A\u672A\u5C31\u7EEA","log_flip_mode_book_id":"\u7FFB\u9875\u6A21\u5F0F\u4E66\u7C4DID: %s","log_invalid_upload_path":"\u65E0\u6548\u7684\u4E0A\u4F20\u8DEF\u5F84","log_time_elapsed":"\u8017\u65F6: %v","log_waiting_for_api_health":"\u7B49\u5F85API\u5065\u5EB7\u68C0\u67E5\u7AEF\u70B9...","log_api_healthy_ready":"Comigo API\u5DF2\u5C31\u7EEA\u5E76\u6B63\u5E38\u8FD0\u884C","log_path_error":"\u8DEF\u5F84\u9519\u8BEF","log_checking_cfg_sharename":"\u6B63\u5728\u68C0\u67E5\u914D\u7F6EShareName","log_warning_failed_to_get_homedir":"\u8B66\u544A: \u83B7\u53D6\u4E3B\u76EE\u5F55\u5931\u8D25: %v","log_warning_failed_to_get_executable_path":"\u8B66\u544A: \u83B7\u53D6\u53EF\u6267\u884C\u6587\u4EF6\u8DEF\u5F84\u5931\u8D25: %v","log_child_daemon_started":"\u5B50\u5B88\u62A4\u8FDB\u7A0B\u5DF2\u542F\u52A8\uFF1F","log_daemon_started":"\u5B88\u62A4\u8FDB\u7A0B\u5DF2\u542F\u52A8","log_update_config":"\u66F4\u65B0\u914D\u7F6E: %s = %s","log_failed_to_update_local_config":"\u66F4\u65B0\u672C\u5730\u914D\u7F6E\u5931\u8D25: %v","log_update_user_info_username":"\u66F4\u65B0\u7528\u6237\u4FE1\u606F: \u7528\u6237\u540D=%s","log_update_user_info_current_password":"\u66F4\u65B0\u7528\u6237\u4FE1\u606F: \u5F53\u524D\u5BC6\u7801=%s","log_update_user_info_password":"\u66F4\u65B0\u7528\u6237\u4FE1\u606F: \u5BC6\u7801=%s","log_update_user_info_reenter_password":"\u66F4\u65B0\u7528\u6237\u4FE1\u606F: \u91CD\u65B0\u8F93\u5165\u5BC6\u7801=%s","log_add_array_config_handler":"\u6DFB\u52A0\u6570\u7EC4\u914D\u7F6E\u5904\u7406: %s = %s","log_delete_array_config_handler":"\u5220\u9664\u6570\u7EC4\u914D\u7F6E\u5904\u7406: %s = %s","log_received_json_data":"\u6536\u5230JSON\u6570\u636E: %s","log_failed_to_parse_json":"\u89E3\u6790JSON\u5931\u8D25","log_failed_to_scan_store_path":"\u626B\u63CF\u4E66\u5E93\u8DEF\u5F84\u5931\u8D25: %v","log_failed_to_save_results_to_database":"\u4FDD\u5B58\u7ED3\u679C\u5230\u6570\u636E\u5E93\u5931\u8D25: %v","log_error_failed_save_to_directory":"\u9519\u8BEF: \u4FDD\u5B58\u5230 %s \u76EE\u5F55\u5931\u8D25","log_error_find_config_in":"\u9519\u8BEF: \u5728 %s %s \u627E\u5230\u914D\u7F6E\u6587\u4EF6","log_error_failed_to_delete_config":"\u9519\u8BEF: \u5220\u9664 %s \u76EE\u5F55\u4E2D\u7684\u914D\u7F6E\u5931\u8D25","log_store_url_already_exists":"\u4E66\u5E93URL\u5DF2\u5B58\u5728: %s","log_store_url_already_exists_in_config":"\u914D\u7F6E\u4E2D\u4E66\u5E93URL\u5DF2\u5B58\u5728: %s","log_failed_to_add_store_url":"\u4ECE\u914D\u7F6E\u6DFB\u52A0\u4E66\u5E93URL\u5931\u8D25: %s","log_string_already_exists":"\u5B57\u7B26\u4E32 \'%s\' \u5DF2\u5B58\u5728","log_failed_to_unmarshal_json":"\u53CD\u5E8F\u5217\u5316JSON\u5931\u8D25: %v","log_unknown_config_key":"\u672A\u77E5\u914D\u7F6E\u952E: %s","log_error_getting_executable_path":"\u83B7\u53D6\u53EF\u6267\u884C\u6587\u4EF6\u8DEF\u5F84\u9519\u8BEF: %s","log_executable_name":"\u53EF\u6267\u884C\u6587\u4EF6\u540D: %s","log_executable_path":"\u53EF\u6267\u884C\u6587\u4EF6\u8DEF\u5F84: %s","log_failed_to_unregister_windows_context_menu":"\u53D6\u6D88\u6CE8\u518CWindows\u53F3\u952E\u83DC\u5355\u5931\u8D25: %v","log_failed_to_register_windows_context_menu":"\u6CE8\u518CWindows\u53F3\u952E\u83DC\u5355\u5931\u8D25: %v","log_received_new_args_from_instance":"\u4ECE\u53E6\u4E00\u4E2A\u5B9E\u4F8B\u6536\u5230\u65B0\u53C2\u6570: %v","log_single_instance_check_failed":"\u5355\u5B9E\u4F8B\u68C0\u67E5\u5931\u8D25: %v","log_args_sent_to_existing_instance":"\u53C2\u6570\u5DF2\u53D1\u9001\u5230\u73B0\u6709\u5B9E\u4F8B\uFF0C\u6B63\u5728\u9000\u51FA...","log_error_listing_books":"\u5217\u51FA\u4E66\u7C4D\u9519\u8BEF: %s","log_error_clearing_temp_files":"\u6E05\u7406\u4E34\u65F6\u6587\u4EF6\u9519\u8BEF: %s","log_cleared_temp_files":"\u5DF2\u6E05\u7406\u4E34\u65F6\u6587\u4EF6: %s","log_failed_to_get_homedirectory":"\u83B7\u53D6\u4E3B\u76EE\u5F55\u5931\u8D25: %s","log_program_directory":"\u7A0B\u5E8F\u76EE\u5F55: %s","log_failed_to_get_program_directory":"\u83B7\u53D6\u7A0B\u5E8F\u76EE\u5F55\u5931\u8D25: %v","log_failed_to_get_working_directory":"\u83B7\u53D6\u5DE5\u4F5C\u76EE\u5F55\u5931\u8D25: %s","log_failed_to_open_file":"\u6253\u5F00\u6587\u4EF6\u5931\u8D25: %s, \u9519\u8BEF: %v","log_failed_to_get_file_info":"\u83B7\u53D6\u6587\u4EF6\u4FE1\u606F\u5931\u8D25: %s, \u9519\u8BEF: %v","log_add_book_error":"\u6DFB\u52A0\u4E66\u7C4D\u9519\u8BEF \u4E66\u7C4DID:%s %s","log_get_child_books_for_bookid":"\u83B7\u53D6\u4E66\u7C4DID %s \u7684\u5B50\u4E66\u7C4D","log_get_child_books_count":"\u83B7\u53D6\u4E66\u7C4DID %v \u7684 %v \u672C\u5B50\u4E66\u7C4D","log_get_bookshelf_error":"\u83B7\u53D6\u4E66\u67B6\u9519\u8BEF: %v","log_error_getting_absolute_path":"\u83B7\u53D6\u7EDD\u5BF9\u8DEF\u5F84\u9519\u8BEF: %v","log_cannot_shorten_id":"\u65E0\u6CD5\u7F29\u77EDID: %s","log_error_getting_book_group":"\u83B7\u53D6\u4E66\u7C4D\u7EC4\u9519\u8BEF: %s","log_warning_failed_to_set_socket_permissions":"\u8B66\u544A: \u8BBE\u7F6E\u5957\u63A5\u5B57\u6743\u9650\u5931\u8D25: %v","log_single_instance_server_started":"\u5355\u5B9E\u4F8B\u670D\u52A1\u5668\u5DF2\u542F\u52A8: %s","log_failed_to_accept_connection":"\u63A5\u53D7\u8FDE\u63A5\u5931\u8D25: %v","log_failed_to_decode_message":"\u89E3\u7801\u6D88\u606F\u5931\u8D25: %v","log_failed_to_handle_new_args":"\u5904\u7406\u65B0\u53C2\u6570\u5931\u8D25: %v","log_received_and_processed_new_args":"\u5DF2\u63A5\u6536\u5E76\u5904\u7406\u65B0\u53C2\u6570: %v","log_failed_to_read_response":"\u8BFB\u53D6\u54CD\u5E94\u5931\u8D25\uFF0C\u4F46\u6D88\u606F\u53EF\u80FD\u5DF2\u53D1\u9001: %v","log_successfully_sent_args":"\u6210\u529F\u53D1\u9001\u53C2\u6570\u5230\u73B0\u6709\u5B9E\u4F8B: %v","log_error_closing_listener":"\u5173\u95ED\u76D1\u542C\u5668\u9519\u8BEF: %v","log_another_instance_running":"\u53E6\u4E00\u4E2A\u5B9E\u4F8B\u6B63\u5728\u8FD0\u884C\uFF0C\u6B63\u5728\u5411\u5176\u53D1\u9001\u53C2\u6570...","log_get_book_error":"\u83B7\u53D6\u4E66\u7C4D\u9519\u8BEF: %v","log_saving_books_meta_data_to":"\u6B63\u5728\u4FDD\u5B58\u4E66\u7C4DMetadata\u5230 %s","log_error_saving_book":"\u4FDD\u5B58\u4E66\u7C4D %s \u9519\u8BEF: %s","log_successfully_saved_books":"\u6210\u529F\u4FDD\u5B58 %d \u672C\u4E66\u7C4D\u5230 %s","log_cfg_save_to":"\u914D\u7F6E\u4FDD\u5B58\u5230 %s","log_try_delete_cfg_in":"\u5C1D\u8BD5\u5220\u9664 %s \u4E2D\u7684\u914D\u7F6E","log_websocket_server_received":"websocket\u670D\u52A1\u5668\u6536\u5230: %v","log_failed_to_create_temp_config_dir":"\u521B\u5EFA\u4E34\u65F6\u914D\u7F6E\u76EE\u5F55\u5931\u8D25: %s","log_failed_to_create_config_dir":"\u521B\u5EFA\u914D\u7F6E\u76EE\u5F55\u5931\u8D25: %s","log_username_or_password_empty":"\u7528\u6237\u540D\u6216\u5BC6\u7801\u4E3A\u7A7A\u3002\u4F7F\u7528\u9ED8\u8BA4JWT\u7B7E\u540D\u5BC6\u94A5\u3002","log_invalid_port_number":"\u65E0\u6548\u7684\u7AEF\u53E3\u53F7\u3002\u4F7F\u7528\u9ED8\u8BA4\u7AEF\u53E3: %d","log_tailscale_status_not_available":"Tailscale\u72B6\u6001\u5C1A\u4E0D\u53EF\u7528: %v","log_server_action":"\u670D\u52A1\u5668\u64CD\u4F5C: %v","log_server_action_string":"\u670D\u52A1\u5668\u64CD\u4F5C: %s","log_loading_books_from":"\u6B63\u5728\u4ECE %s \u52A0\u8F7D\u4E66\u7C4D","log_configured_store_urls":"\u5DF2\u914D\u7F6E\u7684\u4E66\u5E93URL: %v","log_book_data_directory_not_exist":"\u4E66\u7C4D\u6570\u636E\u76EE\u5F55\u5C1A\u4E0D\u5B58\u5728: %s","log_error_accessing_book_data_directory":"\u8BBF\u95EE\u4E66\u7C4D\u6570\u636E\u76EE\u5F55\u9519\u8BEF: %s","log_error_reading_book_data_directory":"\u8BFB\u53D6\u4E66\u7C4D\u6570\u636E\u76EE\u5F55\u9519\u8BEF: %s","log_skipping_directory":"\u8DF3\u8FC7\u76EE\u5F55 %s","log_skipping_non_json_file":"\u8DF3\u8FC7\u975EJSON\u6587\u4EF6 %s","log_error_reading_file":"\u8BFB\u53D6\u6587\u4EF6 %s \u9519\u8BEF: %s","log_warning_corrupted_json_file":"\u8B66\u544A: JSON\u6587\u4EF6 %s \u5DF2\u635F\u574F\uFF0C\u8DF3\u8FC7: %s","log_error_deleting_corrupted_file":"\u5220\u9664\u635F\u574F\u7684\u6587\u4EF6 %s \u9519\u8BEF: %s","log_error_adding_book_to_store":"\u6DFB\u52A0\u4E66\u7C4D %s \u5230\u4E66\u5E93\u9519\u8BEF: %s","log_loaded_books_so_far":"\u5DF2\u52A0\u8F7D %d \u672C\u4E66\u7C4D\uFF0C\u8DEF\u5F84: %s","log_successfully_loaded_books":"\u6210\u529F\u52A0\u8F7D %d \u672C\u4E66\u7C4D(%s)","log_error_adding_subfolder":"\u6DFB\u52A0\u5B50\u6587\u4EF6\u5939\u9519\u8BEF: %s","log_error_saving_book_to_json":"\u4FDD\u5B58\u4E66\u7C4D %s \u5230JSON\u9519\u8BEF: %s","log_error_adding_book":"\u6DFB\u52A0\u4E66\u7C4D %s \u9519\u8BEF: %s","log_error_deleting_book_json_file":"\u5220\u9664\u4E66\u7C4D %s JSON\u6587\u4EF6\u9519\u8BEF: %s","err_add_bookstore_key_exists":"\u6DFB\u52A0\u4E66\u5E93\u9519\u8BEF\uFF1A\u952E\u5DF2\u5B58\u5728 [%s]","err_add_book_empty_bookid":"\u6DFB\u52A0\u4E66\u7C4D\u9519\u8BEF\uFF1ABookID\u4E3A\u7A7A","err_add_bookstore_key_not_found":"\u6DFB\u52A0\u4E66\u5E93\u9519\u8BEF\uFF1A\u672A\u627E\u5230\u952E [%s]","err_cannot_find_group":"\u65E0\u6CD5\u627E\u5230\u7EC4\uFF0Cid=%s","err_getbook_cannot_find":"GetBook\uFF1A\u65E0\u6CD5\u627E\u5230\u4E66\u7C4D\uFF0Cid=%s","err_storebookmark_cannot_find":"StoreBookMark\uFF1A\u65E0\u6CD5\u627E\u5230\u4E66\u7C4D\uFF0Cid=%s","err_storebookmark_unknown_type":"StoreBookMark\uFF1A\u672A\u77E5\u4E66\u7B7E\u7C7B\u578B","err_getbookmark_cannot_find":"GetBookMark\uFF1A\u65E0\u6CD5\u627E\u5230\u4E66\u7C4D\uFF0Cid=%s","err_deletebook_cannot_find":"DeleteBook\uFF1A\u65E0\u6CD5\u627E\u5230\u4E66\u7C4D\uFF0Cid=%s","err_cannot_find_book_topofshelf":"\u9519\u8BEF\uFF1A\u5728TopOfShelfInfo\u4E2D\u65E0\u6CD5\u627E\u5230\u4E66\u7C4D","err_cannot_find_child_books":"\u65E0\u6CD5\u627E\u5230\u5B50\u4E66\u7C4D\u4FE1\u606F\uFF0CBookID\uFF1A%s","err_getparentbook_cannot_find":"GetParentBook: \u65E0\u6CD5\u901A\u8FC7childID\u627E\u5230\u4E66\u7C4D=%s","err_cannot_find_book_parentfolder":"\u65E0\u6CD5\u627E\u5230\u4E66\u7C4D\uFF0CparentFolder=%s","log_book_data_already_exists":"\u4E66\u7C4D\u6570\u636E\u5DF2\u5B58\u5728: %s  %s","err_slice_not_supported":"\u6682\u4E0D\u652F\u6301\u6B64 slice \u7684\u8BBE\u7F6E(\u4EC5\u652F\u6301 []string)","err_must_be_nonempty_config_pointer":"\u5FC5\u987B\u662F\u4E00\u4E2A\u975E\u7A7A\u7684 *Config \u6307\u9488","err_non_htmx_request":"non-htmx request","err_no_form_data":"no form data","err_failed_to_find_executable_path":"\u9519\u8BEF: \u65E0\u6CD5\u627E\u5230\u53EF\u6267\u884C\u6587\u4EF6\u8DEF\u5F84","err_failed_to_find_home_directory":"\u9519\u8BEF: \u65E0\u6CD5\u627E\u5230\u4E3B\u76EE\u5F55","err_funnel_mode_ports_only":"funnel\u6A21\u5F0F\u4EC5\u652F\u6301443\u30018443\u548C10000\u7AEF\u53E3","err_number_not_found":"\u672A\u627E\u5230\u6570\u5B57","err_file_not_zip_archive":"\u6587\u4EF6\u4E0D\u662FZIP\u538B\u7F29\u5305","err_file_not_rar_archive":"\u6587\u4EF6\u4E0D\u662FRAR\u538B\u7F29\u5305","err_scan_file_error":"\u626B\u63CF\u6587\u4EF6\u9519\u8BEF","err_countpages_pdf_invalid":"CountPagesOfPDF: \u65E0\u6548\u7684PDF: %s %s","err_name_in_archive_empty":"nameInArchive\u4E3A\u7A7A","err_unsupported_archive_format":"\u4E0D\u652F\u6301\u7684\u538B\u7F29\u683C\u5F0F\u6216\u5728\u538B\u7F29\u5305\u4E2D\u672A\u627E\u5230\u6587\u4EF6","err_file_not_found_in_archive":"\u5728\u538B\u7F29\u5305\u4E2D\u672A\u627E\u5230\u6587\u4EF6","err_extract_path_not_found":"extractPath\u5728\u4E0A\u4E0B\u6587\u4E2D\u672A\u627E\u5230","err_needfile_empty":"needFile\u4E3A\u7A7A","err_getdata_from_epub_error":"getDataFromEpub\u9519\u8BEF\u3002epubPath:%s  needFile:%s","err_container_xml_empty":"container.xml\u5185\u5BB9\u4E3A\u7A7A","err_no_valid_opf_path":"container.xml\u4E2D\u672A\u627E\u5230\u6709\u6548\u7684OPF\u8DEF\u5F84","err_content_type_not_found":"\u5728\u7F13\u5B58\u4E2D\u672A\u627E\u5230contentType","err_file_does_not_exist":"\u6587\u4EF6\u4E0D\u5B58\u5728:%s","err_charset_not_found":"\u672A\u627E\u5230\u5B57\u7B26\u96C6","err_imaging_decode_error":"imaging.Decode() \u9519\u8BEF","err_imaging_encode_error":"imaging.Encode() \u9519\u8BEF","err_jpeg_encode_error":"digestImage jpeg.Encode() \u9519\u8BEF","err_imageresize_maxwidth_error":"ImageResizeByMaxWidth \u9519\u8BEF maxWidth(%d) > sourceWidth(%d)","err_imageresize_maxheight_error":"ImageResizeByMaxHeight \u9519\u8BEF maxHeight(%d) > sourceHeight(%d)","log_starting_server_on_port":"\u6B63\u5728\u542F\u52A8\u670D\u52A1\u5668...\u7AEF\u53E3 %d ...","log_auto_tls_enabled_for_domain":"\u81EA\u52A8TLS\u5DF2\u542F\u7528\uFF0C\u57DF\u540D: %s","log_custom_tls_cert":"\u81EA\u5B9A\u4E49TLS\u8BC1\u4E66 CertFile: %s KeyFile: %s","log_server_shutdown_successfully":"\u670D\u52A1\u5668\u5DF2\u6210\u529F\u5173\u95ED\uFF0C\u6B63\u5728\u542F\u52A8\u670D\u52A1\u5668...\u7AEF\u53E3 %d ...","log_tailscale_server_initialized":"Tailscale\u670D\u52A1\u5668\u5DF2\u6210\u529F\u521D\u59CB\u5316 %s:%d","log_working_directory":"\u5DE5\u4F5C\u76EE\u5F55: %s","log_failed_to_add_working_directory_to_store_urls":"\u5C06\u5DE5\u4F5C\u76EE\u5F55\u6DFB\u52A0\u5230\u4E66\u5E93URL\u5931\u8D25: %s","log_args_index":"args[%d]: %s","log_failed_to_add_store_url_from_args":"\u4ECE\u53C2\u6570\u6DFB\u52A0\u4E66\u5E93URL\u5931\u8D25: %s","log_failed_to_add_upload_path_to_store_urls":"\u5C06\u4E0A\u4F20\u8DEF\u5F84\u6DFB\u52A0\u5230\u4E66\u5E93URL\u5931\u8D25: %s","log_starting_tailscale_http_server":"\u6B63\u5728\u542F\u52A8Tailscale HTTP\u670D\u52A1\u5668 %s:%d","log_tailscale_server_stopped_successfully":"Tailscale\u670D\u52A1\u5668\u5DF2\u6210\u529F\u505C\u6B62","log_checking_book_files_exist":"\u6B63\u5728\u68C0\u67E5\u4E66\u7C4D\u6587\u4EF6\u662F\u5426\u5B58\u5728...","log_checking_store_exist":"\u6B63\u5728\u68C0\u67E5\u4E66\u5E93\u662F\u5426\u5B58\u5728...","log_error_deleting_book":"\u5220\u9664\u4E66\u7C4D %s \u9519\u8BEF: %s","log_error_initializing_main_folder":"\u521D\u59CB\u5316\u4E3B\u6587\u4EF6\u5939\u9519\u8BEF: %s","log_failed_to_get_absolute_path_scan":"\u83B7\u53D6\u7EDD\u5BF9\u8DEF\u5F84\u5931\u8D25: %s","log_skip_to_scan_root_directory":"\u8DF3\u8FC7\u626B\u63CF\u6839\u76EE\u5F55: %s, %v","log_failed_to_get_relative_path":"\u83B7\u53D6\u76F8\u5BF9\u8DEF\u5F84\u5931\u8D25: %s","log_skip_to_scan_directory":"\u8DF3\u8FC7\u76EE\u5F55: %s, %v","log_failed_to_open_file_get_single":"\u65E0\u6CD5\u6253\u5F00\u6587\u4EF6 %s: %v","log_timeout_identify_archive_format":"\u64CD\u4F5C\u8D85\u65F6\uFF1A\u8BC6\u522B\u538B\u7F29\u683C\u5F0F\u82B1\u8D39\u4E86\u8D85\u8FC730\u79D2","log_failed_to_identify_archive_format":"\u8BC6\u522B\u538B\u7F29\u683C\u5F0F\u5931\u8D25: %v","log_timeout_create_filesystem":"\u64CD\u4F5C\u8D85\u65F6\uFF1A\u521B\u5EFA\u6587\u4EF6\u7CFB\u7EDF\u82B1\u8D39\u4E86\u8D85\u8FC730\u79D2","log_failed_to_create_filesystem":"\u521B\u5EFA\u6587\u4EF6\u7CFB\u7EDF\u5931\u8D25: %v","log_timeout_read_file_content":"\u64CD\u4F5C\u8D85\u65F6\uFF1A\u8BFB\u53D6\u6587\u4EF6\u5185\u5BB9\u82B1\u8D39\u4E86\u8D85\u8FC730\u79D2","log_failed_to_read_file_content":"\u8BFB\u53D6\u6587\u4EF6\u5185\u5BB9\u5931\u8D25: %v","log_timeout_open_file_in_archive":"\u64CD\u4F5C\u8D85\u65F6\uFF1A\u6253\u5F00\u538B\u7F29\u5305\u5185\u6587\u4EF6\u82B1\u8D39\u4E86\u8D85\u8FC730\u79D2","log_failed_to_open_file_in_archive":"\u6253\u5F00\u538B\u7F29\u5305\u5185\u6587\u4EF6\u5931\u8D25: %v","log_timeout_extract_file":"\u64CD\u4F5C\u8D85\u65F6\uFF1A\u63D0\u53D6\u6587\u4EF6\u82B1\u8D39\u4E86\u8D85\u8FC730\u79D2","log_failed_to_extract_file":"\u63D0\u53D6\u6587\u4EF6\u5931\u8D25: %v","err_tailscale_netlistener_nil":"Tailscale netListener\u4E3A\u7A7A\uFF1B\u670D\u52A1\u5668\u5C06\u65E0\u6CD5\u542F\u52A8","err_tailscale_http_server_error":"Tailscale HTTP\u670D\u52A1\u5668\u9519\u8BEF: %v","err_error_closing_network_listener":"\u5173\u95ED\u7F51\u7EDC\u76D1\u542C\u5668\u9519\u8BEF: %v","err_error_closing_tailscale_server":"\u5173\u95EDTailscale\u670D\u52A1\u5668\u9519\u8BEF: %v","err_failed_to_create_tailscale_funnel_listener":"\u521B\u5EFATailscale funnel\u76D1\u542C\u5668\u5931\u8D25 %s: %v","err_failed_to_create_tailscale_listener":"\u521B\u5EFATailscale\u76D1\u542C\u5668\u5931\u8D25 %s: %v","err_failed_to_create_tailscale_local_client":"\u521B\u5EFATailscale\u672C\u5730\u5BA2\u6237\u7AEF\u5931\u8D25: %v","err_error_stopping_tailscale_server":"\u505C\u6B62Tailscale\u670D\u52A1\u5668\u9519\u8BEF: %v","err_failed_to_set_config_value":"\u8BBE\u7F6E\u914D\u7F6E\u503C\u5931\u8D25: %v","err_failed_to_parse_bool":"\u65E0\u6CD5\u5C06 \'%s\' \u89E3\u6790\u4E3A bool: %v","err_failed_to_parse_int":"\u65E0\u6CD5\u5C06 \'%s\' \u89E3\u6790\u4E3A int: %v","err_failed_to_set_username":"\u8BBE\u7F6E\u7528\u6237\u540D\u5931\u8D25: %v","err_failed_to_render_userinfo_config_template":"\u6E32\u67D3UserInfoConfig\u6A21\u677F\u5931\u8D25: %v","err_failed_to_add_config_value":"\u6DFB\u52A0\u914D\u7F6E\u503C\u5931\u8D25: %v","err_getbookshelf_error":"GetBookShelf\u9519\u8BEF: %v","err_failed_to_run_tailscale":"\u8FD0\u884CTailscale\u5931\u8D25: %v","err_store_url_already_exists_error":"\u4E66\u5E93URL\u5DF2\u5B58\u5728: %s","err_store_path_is_subdir_of_existing":"\u65B0\u4E66\u5E93\u8DEF\u5F84\u662F\u5DF2\u6709\u4E66\u5E93\u7684\u5B50\u76EE\u5F55: %s \u662F %s \u7684\u5B50\u76EE\u5F55","err_store_path_is_parent_of_existing":"\u65B0\u4E66\u5E93\u8DEF\u5F84\u662F\u5DF2\u6709\u4E66\u5E93\u7684\u7236\u76EE\u5F55: %s \u662F %s \u7684\u7236\u76EE\u5F55","err_store_path_conflict":"\u4E66\u5E93\u8DEF\u5F84\u51B2\u7A81","err_invalid_store_path":"\u65E0\u6548\u7684\u4E66\u5E93\u8DEF\u5F84: %s","err_field_not_exists":"\u4E0D\u5B58\u5728\u540D\u4E3A \'%s\' \u7684\u5B57\u6BB5","err_field_cannot_set":"\u65E0\u6CD5\u5BF9\u5B57\u6BB5 \'%s\' \u8FDB\u884C\u8BBE\u7F6E","err_field_not_slice_type":"\u5B57\u6BB5 \'%s\' \u4E0D\u662F\u5207\u7247\u7C7B\u578B","err_field_element_not_string":"\u5B57\u6BB5 \'%s\' \u7684\u5143\u7D20\u7C7B\u578B\u4E0D\u662F string","err_field_type_not_supported":"\u6682\u4E0D\u652F\u6301\u8BBE\u7F6E\u5B57\u6BB5 \'%s\' \u7684\u7C7B\u578B: %s","log_failed_to_create_extract_path":"\u521B\u5EFA\u89E3\u538B\u8DEF\u5F84\u5931\u8D25: %v","log_failed_to_open_file_unarchive":"\u6253\u5F00\u6587\u4EF6\u5931\u8D25: %v","log_failed_to_identify_file_format":"\u8BC6\u522B\u6587\u4EF6\u683C\u5F0F\u5931\u8D25: %v","log_failed_to_extract_zip_file":"\u89E3\u538BZIP\u6587\u4EF6\u5931\u8D25: %v","log_zip_file_extracted":"ZIP\u6587\u4EF6\u89E3\u538B\u5B8C\u6210\uFF1A%s \u89E3\u538B\u5230\uFF1A%s","log_failed_to_extract_rar_file":"\u89E3\u538BRAR\u6587\u4EF6\u5931\u8D25: %v","log_rar_file_extracted":"RAR\u6587\u4EF6\u89E3\u538B\u5B8C\u6210\uFF1A%s \u89E3\u538B\u5230\uFF1A%s","log_file_close_error":"file.Close() \u9519\u8BEF: %s","log_failed_to_create_directory":"\u521B\u5EFA\u76EE\u5F55\u5931\u8D25: %v","log_failed_to_create_parent_directory":"\u521B\u5EFA\u7236\u76EE\u5F55\u5931\u8D25: %v","log_failed_to_create_file":"\u521B\u5EFA\u6587\u4EF6\u5931\u8D25: %v","log_failed_to_copy_file_content":"\u590D\u5236\u6587\u4EF6\u5185\u5BB9\u5931\u8D25: %v","log_failed_to_get_container_xml":"\u83B7\u53D6container.xml\u5931\u8D25: %s","log_failed_to_parse_container_xml":"\u89E3\u6790container.xml\u5931\u8D25: %s","log_html_tokenizer_error":"HTML\u5206\u8BCD\u5668\u9519\u8BEF: %v","log_failed_to_get_opf_file_path":"\u83B7\u53D6OPF\u6587\u4EF6\u8DEF\u5F84\u5931\u8D25: %s","log_failed_to_read_opf_file":"\u8BFB\u53D6OPF\u6587\u4EF6\u5931\u8D25: %s","log_failed_to_parse_opf_file":"\u89E3\u6790OPF\u6587\u4EF6\u5931\u8D25: %s","log_failed_to_write_file_to_cache":"\u5199\u5165\u6587\u4EF6\u5230\u7F13\u5B58\u5931\u8D25: %v","log_content_type_not_found_in_cache":"\u7F13\u5B58\u4E2D\u672A\u627E\u5230ContentType\uFF0C\u952E: %+v","log_failed_to_read_file_from_cache":"\u4ECE\u7F13\u5B58\u8BFB\u53D6\u6587\u4EF6\u5931\u8D25: %v","log_cover_file_not_found":"\u5C01\u9762\u6587\u4EF6\u672A\u627E\u5230","log_created_new_book":"\u521B\u5EFA\u65B0\u4E66\u7C4D: %s","log_updated_existing_book":"\u66F4\u65B0\u73B0\u6709\u4E66\u7C4D: %s %s","log_saved_media_files_for_book":"\u4E3A\u4E66\u7C4D %s \u4FDD\u5B58\u4E86 %d \u4E2A\u5A92\u4F53\u6587\u4EF6","log_saved_bookmarks_for_book":"\u4E3A\u4E66\u7C4D %s \u4FDD\u5B58\u4E86 %d \u4E2A\u4E66\u7B7E","log_get_media_files_for_book_error":"\u83B7\u53D6\u4E66\u7C4D %s \u7684\u5A92\u4F53\u6587\u4EF6\u9519\u8BEF: %s","log_get_bookmarks_for_book_error":"\u83B7\u53D6\u4E66\u7C4D %s \u7684\u4E66\u7B7E\u9519\u8BEF: %s","log_error_creating_new_book_group":"\u521B\u5EFA\u65B0\u4E66\u7EC4\u9519\u8BEF: %s","log_failed_to_read_directory":"\u8BFB\u53D6\u76EE\u5F55\u5931\u8D25: %s, \u9519\u8BEF: %v","log_failed_to_get_file_info_scan":"\u83B7\u53D6\u6587\u4EF6\u4FE1\u606F\u5931\u8D25: %s, \u9519\u8BEF: %v","log_nonutf8_zip_error":"NonUTF-8 ZIP: %s, \u9519\u8BEF: %s","log_failed_to_get_image_list_from_epub":"\u4ECEEPUB\u83B7\u53D6\u56FE\u7247\u5217\u8868\u5931\u8D25: %s, \u9519\u8BEF: %v","log_failed_to_get_metadata_from_epub":"\u4ECEEPUB\u83B7\u53D6\u5143\u6570\u636E\u5931\u8D25: %s, \u9519\u8BEF: %v","log_failed_to_access_path_in_archive":"\u8BBF\u95EE\u538B\u7F29\u5305\u4E2D\u7684\u8DEF\u5F84 %s \u5931\u8D25: %v","log_skip_scan_path":"\u8DF3\u8FC7\u626B\u63CF: %s","log_skip_unsupported_file_type":"\u8DF3\u8FC7\u4E0D\u652F\u6301\u7684\u6587\u4EF6\u7C7B\u578B: %s","log_found_supported_file_type":"\u627E\u5230\u652F\u6301\u7684\u6587\u4EF6\u7C7B\u578B\uFF08\u975E\u5A92\u4F53\u6587\u4EF6\uFF09: %s","log_processing_file":"\u5904\u7406\u6587\u4EF6: %s (\u8DEF\u5F84: %s)","log_failed_to_get_file_info_in_archive":"\u83B7\u53D6\u538B\u7F29\u5305\u4E2D\u7684\u6587\u4EF6\u4FE1\u606F\u5931\u8D25: %v","log_books_saved_to_database_successfully":"SaveBooksToDatabase: \u6210\u529F\u4FDD\u5B58 %d \u672C\u4E66\u7C4D\u5230\u6570\u636E\u5E93","log_failed_to_store_bookmark":"\u5B58\u50A8\u4E66\u7B7E\u5931\u8D25: %s","log_updated_bookmarks_for_book_id":"\u66F4\u65B0\u4E66\u7C4DID %s \u7684\u4E66\u7B7E: %s","log_getbook_error_common":"GetBook\u9519\u8BEF: %s","log_getpicturedata_error":"GetPictureData\u9519\u8BEF: %s","log_getbook_error_scroll":"GetBook: %v","log_server_not_ready_within_timeout":"\u670D\u52A1\u5668\u5728 %v \u5185\u672A\u5C31\u7EEA\uFF0C\u7EE7\u7EED\u6267\u884C","log_open_database_error":"OpenDatabase\u9519\u8BEF: %s","log_loadbooks_error":"LoadBooks\u9519\u8BEF %s","log_savebooks_error":"SaveBooks\u9519\u8BEF %s","log_error_listing_books_from_database":"\u4ECE\u6570\u636E\u5E93\u5217\u51FA\u4E66\u7C4D\u9519\u8BEF: %s","log_failed_savebookstodatabase":"\u4FDD\u5B58\u4E66\u7C4D\u5230\u6570\u636E\u5E93\u5931\u8D25: %v","log_failed_to_get_free_port":"\u83B7\u53D6\u7A7A\u95F2\u7AEF\u53E3\u5931\u8D25: %v","log_using_port":"\u4F7F\u7528\u7AEF\u53E3: %d","log_http_newrequest_error":"http.NewRequest\u9519\u8BEF: %s","log_client_do_error":"client.Do\u9519\u8BEF: %s","log_body_close_error":"Body.Close()\u9519\u8BEF: %s","log_io_readall_error":"io.ReadAll\u9519\u8BEF: %s","log_response_write_error":"Response.Write\u9519\u8BEF: %s","log_getimagefrompdf_imgdata_nil":"GetImageFromPDF: imgData\u4E3A\u7A7A","log_countpages_pdf_invalid_error":"CountPagesOfPDF: \u65E0\u6548\u7684PDF: %v \u9519\u8BEF: %v","log_getimagefrompdf_time":"GetImageFromPDF: %v","log_syncpage_message_to_flipmode":"SyncPage\u6D88\u606F\u53D1\u9001\u5230FlipMode: %v %v","log_syncpage_message_to_scrollmode":"SyncPage\u6D88\u606F\u53D1\u9001\u5230ScrollMode: %v %v","log_failed_to_open_database":"\u6253\u5F00\u6570\u636E\u5E93\u5931\u8D25: %v","log_failed_to_ping_database":"Ping\u6570\u636E\u5E93\u5931\u8D25: %v","log_failed_to_create_tables":"\u521B\u5EFA\u8868\u5931\u8D25: %v","log_database_initialized_successfully":"\u6570\u636E\u5E93\u521D\u59CB\u5316\u6210\u529F","err_failed_to_read_embedded_image":"\u8BFB\u53D6\u5D4C\u5165\u56FE\u7247\u5931\u8D25: %v","err_failed_to_read_embedded_data":"\u8BFB\u53D6\u5D4C\u5165\u6570\u636E\u5931\u8D25: %v","scroll_wheel_flip":"\u9F20\u6807\u6EDA\u8F6E\u7FFB\u9875","open_in_new_tab":"\u65B0\u6807\u7B7E\u9875\u6253\u5F00","auto_rescan_interval_minutes":"\u5B9A\u671F\u626B\u63CF\u95F4\u9694","AutoRescanIntervalMinutes":"\u5B9A\u671F\u626B\u63CF\u95F4\u9694","auto_rescan_interval_minutes_desc":"\u5B9A\u671F\u626B\u63CF\u4E66\u5E93\u95F4\u9694\u65F6\u95F4\u3002\u5355\u4F4D\u4E3A\u5206\u949F\u3002\u9ED8\u8BA4\u4E3A 0\uFF0C\u8868\u793A\u7981\u7528\u5B9A\u671F\u626B\u63CF","log_scheduler_task_still_running_skip":"\u4E0A\u4E00\u4E2A\u626B\u63CF\u4EFB\u52A1\u4ECD\u5728\u6267\u884C\u4E2D\uFF0C\u8DF3\u8FC7\u672C\u6B21\u626B\u63CF","log_scheduler_task_execution_failed":"\u5B9A\u65F6\u626B\u63CF\u4EFB\u52A1\u6267\u884C\u5931\u8D25: %v","log_scheduler_task_execution_completed":"\u5B9A\u65F6\u626B\u63CF\u4EFB\u52A1\u6267\u884C\u5B8C\u6210","log_scheduler_interval_zero_no_scheduled_scan":"\u626B\u63CF\u95F4\u9694\u4E3A 0\uFF0C\u4E0D\u81EA\u52A8\u626B\u63CF","log_scheduler_stop_task_failed":"\u505C\u6B62\u5B9A\u65F6\u4EFB\u52A1\u5931\u8D25: %v","log_scheduler_stop_old_task_failed":"\u505C\u6B62\u65E7\u5B9A\u65F6\u4EFB\u52A1\u5931\u8D25: %v","log_scheduler_create_scheduler_failed":"\u521B\u5EFA\u8C03\u5EA6\u5668\u5931\u8D25: %v","log_scheduler_create_task_failed":"\u521B\u5EFA\u5B9A\u65F6\u4EFB\u52A1\u5931\u8D25: %v","log_scheduler_task_started":"\u5B9A\u65F6\u626B\u63CF\u4EFB\u52A1\u5DF2\u542F\u52A8\uFF0C\u95F4\u9694: %d \u5206\u949F","log_scheduler_task_stopped":"\u5B9A\u65F6\u626B\u63CF\u4EFB\u52A1\u5DF2\u505C\u6B62","auto_rescan_started":"\u81EA\u52A8\u626B\u63CF\u5DF2\u542F\u52A8\uFF0C\u95F4\u9694: %d \u5206\u949F","auto_rescan_stopped":"\u81EA\u52A8\u626B\u63CF\u5DF2\u505C\u6B62","auto_rescan_updated":"\u81EA\u52A8\u626B\u63CF\u95F4\u9694\u5DF2\u66F4\u65B0\u4E3A: %d \u5206\u949F","auto_rescan_status_running":"\u8FD0\u884C\u4E2D","auto_rescan_status_stopped":"\u5DF2\u505C\u6B62","auto_rescan_current_interval":"\u5F53\u524D\u626B\u63CF\u95F4\u9694: %d \u5206\u949F","auto_rescan_enabled_hint":"\u81EA\u52A8\u91CD\u626B\u5DF2\u542F\u7528\uFF0C\u7CFB\u7EDF\u5C06\u5B9A\u671F\u626B\u63CF\u4E66\u5E93","auto_rescan_disabled_hint":"\u81EA\u52A8\u91CD\u626B\u5DF2\u7981\u7528","other_settings":"\u5176\u4ED6\u8BBE\u7F6E","plugin_enable":"\u542F\u7528\u63D2\u4EF6","err_failed_to_get_config_dir":"\u83B7\u53D6\u914D\u7F6E\u76EE\u5F55\u5931\u8D25: %v","comigo_xyz_title":"Comigo","comigo_xyz_subtitle":"Comigo - \u7B80\u5355\u65B9\u4FBF\u7684\u6F2B\u753B\u9605\u8BFB\u5668","comigo_xyz_description":"\u5728\u6240\u6709\u8BBE\u5907\u4E0A\u770B\u6F2B\u753B\u3002\u65E0\u8BBA\u7535\u8111\u8FD8\u662F\u624B\u673A\uFF0CWindows\u8FD8\u662FLinux\u3001MacOS","comigo_xyz_feature_format":"\uD83D\uDCDA \u591A\u683C\u5F0F\u652F\u6301","comigo_xyz_feature_format_desc":"\u652F\u6301 ZIP\u3001RAR\u3001CBZ\u3001EPUB\u3001PDF \u7B49\u591A\u79CD\u6F2B\u753B\u683C\u5F0F","comigo_xyz_feature_cross_platform":"\uD83C\uDF10 \u8DE8\u5E73\u53F0","comigo_xyz_feature_cross_platform_desc":"\u5355\u4E2A\u6587\u4EF6\uFF0C\u65E0\u9700\u5B89\u88C5\u3001\u652F\u6301 Linux\u3001Windows\u3001Mac OS \u4E09\u5927\u64CD\u4F5C\u7CFB\u7EDF","comigo_xyz_feature_responsive":"\uD83D\uDCF1 \u54CD\u5E94\u5F0F\u8BBE\u8BA1","comigo_xyz_feature_responsive_desc":"\u9002\u914D\u684C\u9762\u7AEF\u548C\u79FB\u52A8\u7AEF\uFF0C\u591A\u8BBE\u5907\u5B9E\u65F6\u540C\u6B65\u7FFB\u9875","comigo_xyz_feature_reading_modes":"\uD83D\uDD04 \u591A\u79CD\u9605\u8BFB\u6A21\u5F0F","comigo_xyz_feature_reading_modes_desc":"\u652F\u6301\u7FFB\u9875\u6A21\u5F0F\u548C\u5377\u8F74\u6A21\u5F0F\uFF0C\u6EE1\u8DB3\u4E0D\u540C\u9605\u8BFB\u4E60\u60EF","comigo_xyz_feature_security":"\uD83D\uDD12 \u5B89\u5168\u53EF\u9760","comigo_xyz_feature_security_desc":"\u652F\u6301 HTTPS \u548C\u7528\u6237\u8BA4\u8BC1\uFF0C\u5185\u7F6ETailscale\u8FDC\u7A0B\u8BBF\u95EE","comigo_xyz_feature_performance":"\u26A1 \u9AD8\u6027\u80FD","comigo_xyz_feature_performance_desc":"\u5FEB\u901F\u626B\u63CF\u548C\u52A0\u8F7D\uFF0C\u6D41\u7545\u7684\u9605\u8BFB\u4F53\u9A8C","comigo_xyz_feature_plugin":"\uD83D\uDD0C \u63D2\u4EF6\u7CFB\u7EDF","comigo_xyz_feature_plugin_desc":"\u652F\u6301\u81EA\u52A8\u7FFB\u9875\u3001\u65F6\u949F\u7B49\u63D2\u4EF6\uFF0C\u53EF\u6269\u5C55\u81EA\u5B9A\u4E49\u63D2\u4EF6","comigo_xyz_feature_media":"\uD83C\uDFAC \u5A92\u4F53\u64AD\u653E","comigo_xyz_feature_media_desc":"\u5185\u7F6E\u97F3\u9891\u3001\u89C6\u9891\u64AD\u653E\u5668","comigo_xyz_feature_download":"\uD83D\uDCE5 \u7075\u6D3B\u4E0B\u8F7D","comigo_xyz_feature_download_desc":"\u652F\u6301\u6253\u5305\u4E0B\u8F7D\u56FE\u7247\u6587\u4EF6\u5939\uFF0C\u652F\u6301\u8F6C\u6362\u4E3A EPUB \u683C\u5F0F","comigo_xyz_feature_history":"\uD83D\uDCDC \u9605\u8BFB\u5386\u53F2","comigo_xyz_feature_history_desc":"\u81EA\u52A8\u8BB0\u5F55\u9605\u8BFB\u5386\u53F2\uFF0C\u65B9\u4FBF\u7EED\u8BFB","comigo_xyz_github_repo":"\uD83D\uDCE6 GitHub \u4ED3\u5E93","comigo_xyz_github_button":"\u8BBF\u95EE GitHub \u9879\u76EE","comigo_xyz_download_latest":"\u4E0B\u8F7D\u6700\u65B0\u7248\u672C","comigo_xyz_quick_start":"\u5FEB\u901F\u5F00\u59CB","comigo_xyz_quick_start_step1":"\u4E0B\u8F7D\u5E76\u8FD0\u884C Comigo \u7A0B\u5E8F","comigo_xyz_quick_start_step2":"\u914D\u7F6E\u60A8\u7684\u4E66\u5E93\u8DEF\u5F84","comigo_xyz_quick_start_step3":"\u5F00\u59CB\u4EAB\u53D7\u9605\u8BFB\uFF01","comigo_xyz_download_for":"\u4E0B\u8F7D\u9002\u7528\u4E8E","comigo_xyz_select_platform":"\u9009\u62E9\u5E73\u53F0","comigo_xyz_download":"\u2B07\uFE0F \u4E0B\u8F7D","comigo_xyz_cli_install_title":"\u4E00\u952E\u5B89\u88C5\u547D\u4EE4\u884C\u7248\uFF1A","comigo_xyz_cli_install_desc":"\u4F7F\u7528\u4EE5\u4E0B\u547D\u4EE4\u5FEB\u901F\u5B89\u88C5 Comigo","comigo_xyz_cli_install_copy":"\u590D\u5236","comigo_xyz_cli_install_copied":"\u5DF2\u590D\u5236","comigo_xyz_cli_install_cn_desc":"\u4E2D\u56FD\u5927\u9646\u7528\u6237\u63A8\u8350\uFF1A","comigo_xyz_docker_deploy_title":"\u4F7F\u7528Docker\u90E8\u7F72\uFF1A","rescan_store":"\u626B\u63CF\u4E66\u5E93","delete_store":"\u5220\u9664\u4E66\u5E93","select_store_to_operate":"\u8BF7\u5148\u9009\u62E9\u8981\u64CD\u4F5C\u7684\u4E66\u5E93","confirm_delete_store":"\u786E\u8BA4\u5220\u9664\u4E66\u5E93\uFF1F\u8FD9\u5C06\u540C\u65F6\u5220\u9664\u8BE5\u4E66\u5E93\u7684\u6240\u6709\u4E66\u7C4D\u6570\u636E","rescan_store_success":"\u4E66\u5E93\u626B\u63CF\u5B8C\u6210\uFF0C\u65B0\u589E {0} \u672C\u4E66","delete_store_success":"\u4E66\u5E93\u5220\u9664\u6210\u529F","rescan_store_in_progress":"\u6B63\u5728\u626B\u63CF\u4E66\u5E93\uFF0C\u8BF7\u7A0D\u5019...","err_rescan_store_failed":"\u626B\u63CF\u4E66\u5E93\u5931\u8D25","err_delete_store_failed":"\u5220\u9664\u4E66\u5E93\u5931\u8D25","log_book_file_not_exist_skip":"\u4E66\u7C4D\u6587\u4EF6\u4E0D\u5B58\u5728\uFF0C\u8DF3\u8FC7\u52A0\u8F7D: %s","log_error_deleting_orphan_metadata":"\u5220\u9664\u5B64\u7ACB\u5143\u6570\u636E\u6587\u4EF6 %s \u9519\u8BEF: %s","log_book_version_mismatch_skip":"\u4E66\u7C4D %s \u7248\u672C\u4E0D\u5339\u914D\uFF08\u7F13\u5B58: %s, \u5F53\u524D: %s\uFF09\uFF0C\u8DF3\u8FC7\u52A0\u8F7D","log_book_version_minor_mismatch":"\u4E66\u7C4D %s \u6B21\u7248\u672C\u4E0D\u540C\uFF08\u7F13\u5B58: %s, \u5F53\u524D: %s\uFF09\uFF0C\u5C06\u8FC1\u79FB\u4E66\u7B7E\u540E\u91CD\u65B0\u626B\u63CF","log_bookmark_saved_for_migration":"\u4E66\u7C4D %s \u7684 %d \u4E2A\u4E66\u7B7E\u5DF2\u4FDD\u5B58\uFF0C\u7B49\u5F85\u8FC1\u79FB\u5230\u65B0\u6570\u636E","log_bookmark_migrated":"\u4E66\u7C4D %s \u6210\u529F\u8FC1\u79FB\u4E86 %d \u4E2A\u4E66\u7B7E","log_error_deleting_version_mismatch_metadata":"\u5220\u9664\u7248\u672C\u4E0D\u5339\u914D\u7684\u5143\u6570\u636E\u6587\u4EF6 %s \u9519\u8BEF: %s","player":"\u64AD\u653E\u5668","play":"\u64AD\u653E","pause":"\u6682\u505C","previous":"\u4E0A\u4E00\u66F2","next":"\u4E0B\u4E00\u66F2","playlist":"\u64AD\u653E\u5217\u8868","audio":"\u97F3\u9891","video":"\u89C6\u9891","browser_not_support_video":"\u60A8\u7684\u6D4F\u89C8\u5668\u4E0D\u652F\u6301\u89C6\u9891\u64AD\u653E","browser_not_support_audio":"\u60A8\u7684\u6D4F\u89C8\u5668\u4E0D\u652F\u6301\u97F3\u9891\u64AD\u653E","play_failed":"\u64AD\u653E\u5931\u8D25","first_media":"\u5DF2\u7ECF\u662F\u7B2C\u4E00\u4E2A\u4E86","last_media":"\u5DF2\u7ECF\u662F\u6700\u540E\u4E00\u4E2A\u4E86","auto_play_next":"\u81EA\u52A8\u4E0B\u4E00\u66F2","loop_playlist":"\u5FAA\u73AF\u64AD\u653E\u5217\u8868","player_autoplay_help":"\u81EA\u52A8\u4E0B\u4E00\u66F2\u4E0E\u5FAA\u73AF\u64AD\u653E\u529F\u80FD\uFF0C\u9700\u8981\u6D4F\u89C8\u5668\u5F00\u542F\u5A92\u4F53\u64AD\u653E\u6743\u9650\u3002\u5728\u79FB\u52A8\u8BBE\u5907\u4E0A\u4E5F\u53EF\u80FD\u56E0\u7701\u7535/\u540E\u53F0\u9650\u5236\u800C\u5931\u6548\u3002","add_bookmark":"\u6DFB\u52A0\u4E66\u7B7E","manual_bookmark":"\u624B\u52A8\u4E66\u7B7E","auto_bookmark":"\u81EA\u52A8\u4E66\u7B7E","bookmark_added":"\u4E66\u7B7E\u5DF2\u6DFB\u52A0","bookmark_exists":"\u8BE5\u9875\u5DF2\u6709\u4E66\u7B7E","continue_reading":"\u7EE7\u7EED\u9605\u8BFB","delete_record":"\u5220\u9664","confirm_delete_bookmark":"\u786E\u5B9A\u5220\u9664\u6B64\u9605\u8BFB\u8BB0\u5F55\uFF1F","bookmark_deleted":"\u5DF2\u5220\u9664","log_remote_store_connect_failed":"\u65E0\u6CD5\u8FDE\u63A5\u8FDC\u7A0B\u4E66\u5E93: %s, \u9519\u8BEF: %v","log_remote_file_stat_failed":"\u65E0\u6CD5\u83B7\u53D6\u8FDC\u7A0B\u6587\u4EF6\u4FE1\u606F: %s, \u9519\u8BEF: %v","log_remote_file_open_failed":"\u65E0\u6CD5\u6253\u5F00\u8FDC\u7A0B\u6587\u4EF6: %s, \u9519\u8BEF: %v","log_remote_pdf_download_on_demand":"\u6309\u9700\u4E0B\u8F7D\u8FDC\u7A0B PDF: %s","log_webdav_filesystem_connected":"WebDAV \u6587\u4EF6\u7CFB\u7EDF\u5DF2\u8FDE\u63A5: %s, \u57FA\u7840\u8DEF\u5F84: %s","log_sftp_filesystem_connected":"SFTP \u6587\u4EF6\u7CFB\u7EDF\u5DF2\u8FDE\u63A5: %s, \u57FA\u7840\u8DEF\u5F84: %s","log_smb_filesystem_connected":"SMB \u6587\u4EF6\u7CFB\u7EDF\u5DF2\u8FDE\u63A5: %s, \u57FA\u7840\u8DEF\u5F84: %s","log_ftp_filesystem_connected":"FTP \u6587\u4EF6\u7CFB\u7EDF\u5DF2\u8FDE\u63A5: %s, \u57FA\u7840\u8DEF\u5F84: %s","log_s3_filesystem_connected":"S3 \u6587\u4EF6\u7CFB\u7EDF\u5DF2\u8FDE\u63A5: %s, \u57FA\u7840\u8DEF\u5F84: %s","log_epub_metadata_remote_not_supported":"EPUB \u5143\u6570\u636E\u63D0\u53D6\u6682\u4E0D\u652F\u6301\u8FDC\u7A0B\u6D41\u5F0F\u8BFB\u53D6","log_remote_file_download_to_cache":"\u4E0B\u8F7D\u8FDC\u7A0B\u6587\u4EF6\u5230\u7F13\u5B58: %s -> %s","log_scan_remote_store_start":"\u5F00\u59CB\u626B\u63CF\u8FDC\u7A0B\u4E66\u5E93: %s","log_scan_start_hint_remote":"\u5F00\u59CB\u626B\u63CF\uFF1A%s (\u8FDC\u7A0B\u8DEF\u5F84: %s)","log_remote_store_check_book_existence_failed":"\u65E0\u6CD5\u8FDE\u63A5\u8FDC\u7A0B\u4E66\u5E93\u68C0\u67E5\u4E66\u7C4D\u5B58\u5728\u6027: %s, \u9519\u8BEF: %v","log_remote_book_existence_check_failed":"\u68C0\u67E5\u8FDC\u7A0B\u4E66\u7C4D\u5B58\u5728\u6027\u5931\u8D25: %s, \u9519\u8BEF: %v","log_add_remote_store":"\u6DFB\u52A0\u8FDC\u7A0B\u4E66\u5E93: %s (\u534F\u8BAE: %s, \u4E3B\u673A: %s)","log_s3_connecting":"\u6B63\u5728\u8FDE\u63A5 S3 \u670D\u52A1 %s (\u5B58\u50A8\u6876: %s, \u524D\u7F00: %s)","log_ftp_connecting":"\u6B63\u5728\u8FDE\u63A5 FTP \u670D\u52A1\u5668 %s (TLS: %v, \u8D85\u65F6: %v)","log_webdav_download_range":"\u4E0B\u8F7D\u7247\u6BB5: %s [%d-%d]","log_smb_connecting":"\u6B63\u5728\u8FDE\u63A5 SMB \u670D\u52A1\u5668 %s (\u8D85\u65F6: %d\u79D2, \u7528\u6237: %s, \u5171\u4EAB: %s)","log_smb_mount_share":"\u6B63\u5728\u6302\u8F7D SMB \u5171\u4EAB: %s","log_using_cached_file":"\u4F7F\u7528\u7F13\u5B58\u6587\u4EF6: %s","log_download_file":"\u4E0B\u8F7D\u6587\u4EF6\uFF1A%s","log_remote_book_existence_check_failed_detail":"\u8FDC\u7A0B\u4E66\u7C4D\u5B58\u5728\u6027\u68C0\u67E5\u5931\u8D25 - BookID: %s, RemoteURL: %s, BookPath: %s, \u9519\u8BEF: %v","log_local_book_existence_check_failed":"\u68C0\u67E5\u672C\u5730\u4E66\u7C4D\u5B58\u5728\u6027\u5931\u8D25: %s, \u9519\u8BEF: %v","log_failed_to_set_field":"\u8BBE\u7F6E\u5B57\u6BB5 %s \u5931\u8D25: %v","log_cache_mkdir_failed":"\u521B\u5EFA\u7F13\u5B58\u76EE\u5F55\u5931\u8D25: %v","log_cache_hit_memory":"\u4ECE\u5185\u5B58\u7F13\u5B58\u547D\u4E2D: %s","log_cache_hit_disk":"\u4ECE\u78C1\u76D8\u7F13\u5B58\u547D\u4E2D: %s","log_cache_write_disk_failed":"\u5199\u5165\u78C1\u76D8\u7F13\u5B58\u5931\u8D25: %v","log_cached_to_disk":"\u5DF2\u7F13\u5B58\u5230\u78C1\u76D8: %s -> %s","log_get_file_error":"\u83B7\u53D6\u6587\u4EF6\u9519\u8BEF: %s","log_save_file_to_cache_error":"SaveFileToCache \u9519\u8BEF: %s","log_save_cover_to_local_error":"SaveCoverToLocal \u9519\u8BEF: %s","log_failed_to_parse_cover_url":"\u89E3\u6790\u5C01\u9762 URL \u5931\u8D25: %s","log_child_book_id_missing_in_cover_url":"\u5C01\u9762 URL \u4E2D\u7F3A\u5C11\u5B50\u4E66\u7C4D ID","log_failed_to_get_child_book":"\u83B7\u53D6\u5B50\u4E66\u7C4D\u5931\u8D25: %s","log_failed_to_read_embedded_image":"\u8BFB\u53D6\u5185\u5D4C\u56FE\u7247\u5931\u8D25: %s","log_failed_to_create_epub_generator":"\u521B\u5EFA EPUB \u751F\u6210\u5668\u5931\u8D25: %s","log_failed_to_get_image_epub":"\u83B7\u53D6\u56FE\u7247\u5931\u8D25 %s: %v","log_failed_to_generate_epub":"\u751F\u6210 EPUB \u5931\u8D25: %s","log_plugin_enabled":"\u542F\u7528\u63D2\u4EF6: %s","log_disable_mutex_plugin_auto_flip":"\u7981\u7528\u4E92\u65A5\u63D2\u4EF6: auto_flip","log_disable_mutex_plugin_sketch_practice":"\u7981\u7528\u4E92\u65A5\u63D2\u4EF6: sketch_practice","log_plugin_disabled":"\u7981\u7528\u63D2\u4EF6: %s","log_rescan_store":"\u91CD\u65B0\u626B\u63CF\u4E66\u5E93: %s","log_rescan_store_completed_new_books":"\u4E66\u5E93\u626B\u63CF\u5B8C\u6210\uFF0C\u65B0\u589E %d \u672C\u4E66","log_delete_store":"\u5220\u9664\u4E66\u5E93: %s","log_deleted_books_count":"\u5220\u9664\u4E86 %d \u672C\u4E66\u7C4D","log_upload_no_store_selected":"\u672A\u9009\u62E9\u4E0A\u4F20\u76EE\u6807\u4E66\u5E93","log_upload_invalid_store_path":"\u65E0\u6548\u7684\u4E66\u5E93\u8DEF\u5F84: %s","log_upload_store_path_not_exist":"\u4E66\u5E93\u8DEF\u5F84\u4E0D\u5B58\u5728: %s","log_error_closing_zip_writer":"\u5173\u95ED zip \u5199\u5165\u5668\u9519\u8BEF: %s","log_file_not_found_skipping":"\u672A\u627E\u5230\u6587\u4EF6\uFF0C\u8DF3\u8FC7: %s","log_error_opening_file":"\u6253\u5F00\u6587\u4EF6\u9519\u8BEF: %s, \u9519\u8BEF: %s","log_error_creating_zip_entry":"\u521B\u5EFA zip \u6761\u76EE\u9519\u8BEF: %s, \u9519\u8BEF: %s","log_error_writing_file_to_zip":"\u5199\u5165\u6587\u4EF6\u5230 zip \u9519\u8BEF: %s, \u9519\u8BEF: %s","log_failed_to_delete_bookmark":"\u5220\u9664\u4E66\u7B7E\u5931\u8D25: %v","log_get_config_dir_error":"GetConfigDir \u9519\u8BEF: %s","log_delete_book_json_error":"DeleteBookJson \u9519\u8BEF: %s","log_delete_cover_cache_error":"DeleteCoverCache \u9519\u8BEF: %s","log_delete_book_cache_error":"DeleteBookCache \u9519\u8BEF: %s","log_plugin_system_disabled_skip_scan":"\u63D2\u4EF6\u7CFB\u7EDF\u672A\u542F\u7528\uFF0C\u8DF3\u8FC7\u81EA\u5B9A\u4E49\u63D2\u4EF6\u626B\u63CF","log_plugin_dir_not_exist_skip_load":"\u63D2\u4EF6\u76EE\u5F55\u4E0D\u5B58\u5728: %s\uFF0C\u8DF3\u8FC7\u81EA\u5B9A\u4E49\u63D2\u4EF6\u52A0\u8F7D","log_plugin_scope_load_error":"\u52A0\u8F7D %s \u8303\u56F4\u63D2\u4EF6\u65F6\u51FA\u9519: %v","log_plugin_custom_loaded_count":"\u6210\u529F\u52A0\u8F7D %d \u4E2A\u81EA\u5B9A\u4E49\u63D2\u4EF6","log_plugin_loaded_item":"  - [%s] %s (%s)","log_plugin_dir_not_exist":"\u63D2\u4EF6\u76EE\u5F55\u4E0D\u5B58\u5728: %s","log_plugin_read_file_failed":"\u8BFB\u53D6\u63D2\u4EF6\u6587\u4EF6\u5931\u8D25 %s: %v","log_plugin_read_book_file_failed":"\u8BFB\u53D6\u4E66\u7C4D\u63D2\u4EF6\u6587\u4EF6\u5931\u8D25 %s: %v","log_plugin_loaded_for_book":"\u52A0\u8F7D\u4E66\u7C4D %s \u7684 %s \u63D2\u4EF6: %d \u4E2A","log_load_custom_plugin_failed":"\u52A0\u8F7D\u81EA\u5B9A\u4E49\u63D2\u4EF6\u5931\u8D25: %v","log_cfg_host_enabled_plugin_list":"cfg.Host: %v , cfg.EnabledPluginList: %v","log_config_changed_restart_web":"\u914D\u7F6E\u5DF2\u53D8\u66F4\uFF0C\u6B63\u5728\u91CD\u542F Web \u670D\u52A1...","log_config_changed_start_tailscale":"\u914D\u7F6E\u5DF2\u53D8\u66F4\uFF0C\u6B63\u5728\u542F\u52A8 Tailscale...","log_config_changed_stop_tailscale":"\u914D\u7F6E\u5DF2\u53D8\u66F4\uFF0C\u6B63\u5728\u505C\u6B62 Tailscale...","log_config_changed_restart_tailscale":"\u914D\u7F6E\u5DF2\u53D8\u66F4\uFF0C\u6B63\u5728\u91CD\u542F Tailscale...","log_received_rescan_message":"\u6536\u5230\u91CD\u65B0\u626B\u63CF\u6D88\u606F\uFF1A%s","log_failed_to_read_image_epub":"\u8BFB\u53D6\u56FE\u7247\u5931\u8D25 %s: %v","log_failed_to_decode_image_config_epub":"\u89E3\u7801\u56FE\u7247\u914D\u7F6E\u5931\u8D25: %v\uFF0C\u4F7F\u7528\u9ED8\u8BA4\u5C3A\u5BF8","log_cannot_set_username_or_password":"\u65E0\u6CD5\u8BBE\u7F6E\u7528\u6237\u540D\u6216\u5BC6\u7801","log_login_failed":"\u767B\u5F55\u5931\u8D25: %s-%s, %s-%s","log_get_generated_image_params":"GetGeneratedImage: height=%s, width=%s, text=%s, font_size=%s","log_toml_marshal_error":"toml.Marshal \u9519\u8BEF","log_config_file_changed_detected":"\u68C0\u6D4B\u5230 Comigo \u914D\u7F6E\u6587\u4EF6\u6539\u53D8: %s Op: %s","log_error_num_value":"\u9519\u8BEF\u6570\u503C: %s"}');


var $654fabe7e18f36ba$exports = {};
$654fabe7e18f36ba$exports = JSON.parse('{"auth":"\u8A8D\u8A3C","book_not_found":"\u8AAD\u3081\u308B\u672C\u304C\u898B\u3064\u304B\u3089\u305A\u3001\u30D7\u30ED\u30B0\u30E9\u30E0\u306F\u7D42\u4E86\u3057\u307E\u3057\u305F\u3002","cache_file_clean":"\u7D42\u4E86\u6642\u306B\u30AD\u30E3\u30C3\u30B7\u30E5\u30D5\u30A1\u30A4\u30EB\u3092\u3059\u3079\u3066\u524A\u9664\u3059\u308B","cache_file_enable":"\u30A6\u30A7\u30D6\u753B\u50CF\u306E\u30AD\u30E3\u30C3\u30B7\u30E5\u3092\u6709\u52B9\u306B\u3059\u308B\u304B\uFF08\u518D\u8AAD\u8FBC\u307F\u3092\u9AD8\u901F\u5316\u3057\u307E\u3059\u304CHDD\u5BB9\u91CF\u3092\u6D88\u8CBB\u3057\u307E\u3059\uFF09","cache_file_dir":"\u30AD\u30E3\u30C3\u30B7\u30E5\u30D5\u30A9\u30EB\u30C0\uFF08\u30C7\u30D5\u30A9\u30EB\u30C8\u306F\u30B7\u30B9\u30C6\u30E0\u306E\u30C6\u30F3\u30DD\u30E9\u30EA\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\uFF09","can_not_init_book":"\u66F8\u7C4D\u306E\u521D\u671F\u5316\u304C\u3067\u304D\u307E\u305B\u3093\u3002","cannot_listen":"Web\u30B5\u30FC\u30D0\u30FC\u3092\u958B\u59CB\u3067\u304D\u307E\u305B\u3093","check_image_completed":"\u89E3\u50CF\u5EA6\u89E3\u6790\u304C\u5B8C\u4E86\u3057\u307E\u3057\u305F","check_image_error":"\u89E3\u50CF\u5EA6\u89E3\u6790\u30A8\u30E9\u30FC\uFF1A","check_image_ing":"\u89E3\u50CF\u5EA6\u89E3\u6790\u4E2D\uFF1A","check_image_start":"\u753B\u50CF\u306E\u89E3\u6790\u3092\u958B\u59CB\u3057\u307E\u3059","check_mac_error":"Mac\u30A2\u30C9\u30EC\u30B9\u691C\u51FA\u30A8\u30E9\u30FC\uFF1A","check_port_error":"\u30DD\u30FC\u30C8\u691C\u77E5\u30A8\u30E9\u30FC\uFF1A%v","clear_temp_file_completed":"\u30AD\u30E3\u30C3\u30B7\u30E5\u30D5\u30A9\u30EB\u30C0\u30FC\u306E\u30AF\u30EA\u30FC\u30F3\u30A2\u30C3\u30D7\u306B\u6210\u529F\u3057\u307E\u3057\u305F\uFF1A","clear_temp_file_error":"\u30AD\u30E3\u30C3\u30B7\u30E5\u30D5\u30A9\u30EB\u30C0\u30FC\u306E\u30AF\u30EA\u30FC\u30F3\u30A2\u30C3\u30D7\u306B\u5931\u6557\u3057\u307E\u3057\u305F\uFF1A","clear_temp_file_start":"---- \u30AD\u30E3\u30C3\u30B7\u30E5\u30D5\u30A9\u30EB\u30C0\u30FC\u306E\u30AF\u30EA\u30FC\u30F3\u30A2\u30C3\u30D7\u3092\u958B\u59CB ----","comigo_example":"  comi book.zip\\n\\n\u30DD\u30FC\u30C8\u3092\u8A2D\u5B9A\u3057\u307E\u3059\uFF08\u30C7\u30D5\u30A9\u30EB\u30C8\u306F1234\uFF09\uFF1A\\n  comi -p 2345 book.zip\\n\\n\u30D6\u30E9\u30A6\u30B6\u3092\u958B\u304B\u305A\u306B\uFF08Windows\uFF09\uFF1A\\n  comi -o=false book.zip\\n\\n\u8907\u6570\u306E\u30D1\u30E9\u30E1\u30FC\u30BF\uFF1A\\n  comi -p 2345  --host example.com test.zip\\n","comigo_use":"comi","completed_extract":"\u89E3\u51CD\u7D42\u4E86","completed_ls":"\u5727\u7E2E\u30D5\u30A1\u30A4\u30EB\u306E\u30B9\u30AD\u30E3\u30F3\u5B8C\u4E86","config":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u306E\u6307\u5B9A","config_change":"comigo\u8A2D\u5B9A\u306E\u5909\u66F4\u304C\u3042\u308A\u307E\u3057\u305F\uFF1A","config_file_not_found":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u304C\u898B\u3064\u304B\u3089\u305A\u3001\u30C7\u30D5\u30A9\u30EB\u30C8\u306E\u8A2D\u5B9A\u3092\u4F7F\u7528\u3057\u307E\u3059","config_file_not_resolve":"\u8A2D\u5B9A\u306E\u89E3\u6790\u306B\u5931\u6557\u3057\u307E\u3057\u305F\u3002","config_save_to":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u306E\u4FDD\u5B58\u5834\u6240\u30C7\u30D5\u30A9\u30EB\u30C8\u5024\uFF1ARAM\u3001HomeDirectory\u3001WorkingDirectory\u3001ProgramDirectory","ctrl_c_hint":"CTRL-C\u3092\u62BC\u3057\u3066\u7D42\u4E86\u3057\u307E\u3059","debug_mode":"\u30C7\u30D0\u30C3\u30B0\u30E2\u30FC\u30C9","disable_lan":"LAN\u5171\u6709\u3092\u7121\u52B9\u5316\u3059\u308B","enable_database":"\u30B9\u30AD\u30E3\u30F3\u3057\u305F\u672C\u3092\u30ED\u30FC\u30AB\u30EB\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u306B\u4FDD\u5B58\u3059\u308B","enable_file_upload":"\u30D5\u30A1\u30A4\u30EB\u306E\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u3092\u6709\u52B9\u306B\u3059\u308B","enable_frpc":"frpc\u3092\u6709\u52B9\u306B\u3059\u308B","enable_webp":"webp\u8EE2\u9001\u3092\u6709\u52B9\u306B\u3059\u308B\uFF08webp-server\u304C\u5FC5\u8981\uFF09","epub_cannot_resort":"epub\u30D5\u30A1\u30A4\u30EB\u306E\u4E26\u3079\u66FF\u3048\u304C\u3067\u304D\u307E\u305B\u3093\uFF1A","exceeds_maximum_depth":"\u6700\u5927\u691C\u7D22\u6DF1\u5EA6\u3092\u8D85\u3048\u3066\u3044\u307E\u3059\u3002MaxDepth =","file_exit":"\u30D5\u30A1\u30A4\u30EB\u306F\u65E2\u306B\u5B58\u5728\u3057\u3066\u3044\u307E\u3059","file_not_found":"\u6B63\u3057\u3044\u30D5\u30A1\u30A4\u30EB\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3067\u3057\u305F\u3002","flip_page_template":"\u30C7\u30D5\u30A9\u30EB\u30C8\uFF1A\u30DA\u30FC\u30B8\u3092\u3081\u304F\u3063\u3066\u8AAD\u3080","format_customization_error":"\u3053\u306E\u62BD\u51FA\u30D5\u30A9\u30FC\u30DE\u30C3\u30C8\u306F\u30B5\u30DD\u30FC\u30C8\u3057\u3066\u3044\u307E\u305B\u3093\uFF1A","found_config_file":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u304C\u898B\u3064\u304B\u308A\u307E\u3057\u305F\uFF1A","found_in_bookstore":"\u672C\u68DA\u3067\u898B\u3064\u304B\u308A\u307E\u3057\u305F\uFF1A","how_many_books_update":"\u30D1\u30B9 %v \u66F4\u65B0 %v \u518A\u306E\u66F8\u7C4D","frp_command":"frpc\u30B3\u30DE\u30F3\u30C9\u3001\u307E\u305F\u306Ffrpc\u5B9F\u884C\u30D5\u30A1\u30A4\u30EB\u306E\u30D1\u30B9","frp_random_remote_port":"frp \u30EA\u30E2\u30FC\u30C8\u30E9\u30F3\u30C0\u30E0\u30DD\u30FC\u30C8(40000~50000)","frp_remote_port":"frp \u30EA\u30E2\u30FC\u30C8\u30DD\u30FC\u30C8\u3001-1\u3067\u30ED\u30FC\u30AB\u30EB\u3068\u540C\u3058\u30DD\u30FC\u30C8\u3092\u4F7F\u7528","frp_server_addr":"frps-addr\uFF08frpc\u304C\u5FC5\u8981\uFF09","frp_server_error":"frpc\u30B5\u30FC\u30D3\u30B9\u3092\u8D77\u52D5\u3067\u304D\u307E\u305B\u3093\u3002\u30B3\u30DE\u30F3\u30C9\u5F62\u5F0F\u3092\u78BA\u8A8D\u3057\u3001PATH\u306Bfrpc\u5B9F\u884C\u30D5\u30A1\u30A4\u30EB\u304C\u5B58\u5728\u3059\u308B\u3053\u3068\u3092\u78BA\u8A8D\u3057\u3066\u304F\u3060\u3055\u3044\u3002","frp_server_port":"frps server_port\uFF08frpc\u304C\u5FC5\u8981\uFF09","frp_setting_save_completed":"frpc\u8A2D\u5B9A\u306E\u4FDD\u5B58\u306B\u6210\u529F\u3057\u307E\u3057\u305F\u3002","frp_token":"token\uFF08frpc\u304C\u5FC5\u8981\uFF09","frpc_ini_error":"frpc ini\u521D\u671F\u5316\u30A8\u30E9\u30FC","frpc_server_start":"FRPC\u30B5\u30FC\u30D3\u30B9\u306F\u958B\u59CB\u3057\u307E\u3057\u305F","generate_metadata":"\u66F8\u7C4D\u306E\u30E1\u30BF\u30C7\u30FC\u30BF\u3092\u751F\u6210\u3059\u308B","get_ip_error":"IP\u53D6\u5F97\u30A8\u30E9\u30FC\uFF1A","html_title":"Comigo Reader ","init_database":"\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u3092\u521D\u671F\u5316\uFF1A","init_locale":"\u30C7\u30D5\u30A9\u30EB\u30C8\u306E\u8A00\u8A9E\u306F\u65E5\u672C\u8A9E\u3002","local_host":"\u30C9\u30E1\u30A4\u30F3\u540D\u306E\u8A2D\u5B9A","local_reading":"\u30ED\u30FC\u30AB\u30EB\u30EA\u30FC\u30C7\u30A3\u30F3\u30B0\uFF1A","log_to_file":"\u30ED\u30B0\u30D5\u30A1\u30A4\u30EB\u3078\u30ED\u30B0\u51FA\u529B","long_description":"comigo \u30B7\u30F3\u30D7\u30EB\u306A\u30B3\u30DF\u30C3\u30AF\u30EA\u30FC\u30C0\u30FC","max_depth":"\u6700\u5927\u691C\u7D22\u6DF1\u5EA6","min_media_num":"zip\u5185\u306E\u753B\u50CF\u6570\u306E\u6700\u5C0F\u8AAD\u53D6\u57FA\u6E96","no_pages_in_pdf":"PDF\u5185\u306B\u30DA\u30FC\u30B8\u304C\u3042\u308A\u307E\u305B\u3093","not_a_valid_zip_file":"\u6709\u52B9\u306AZIP\u30D5\u30A1\u30A4\u30EB\u3067\u306F\u3042\u308A\u307E\u305B\u3093\uFF1A","open_browser":"\u30D6\u30E9\u30A6\u30B6\u30FC\u3092\u958B\u304F\uFF08Windows=true\uFF09","open_browser_error":"\u30D6\u30E9\u30A6\u30B6\u3092\u958B\u304F\u3053\u3068\u304C\u3067\u304D\u307E\u305B\u3093\u3067\u3057\u305F\u3002","open_image_error":"\u30D5\u30A1\u30A4\u30EB\u3092\u958B\u304F\u969B\u306B\u30A8\u30E9\u30FC\u304C\u767A\u751F\u3057\u307E\u3057\u305F\uFF1A","password":"\u30D1\u30B9\u30EF\u30FC\u30C9","path_not_exist":"\u6307\u5B9A\u3055\u308C\u305F\u30D1\u30B9\u306F\u5B58\u5728\u3057\u307E\u305B\u3093","port":"\u30B5\u30FC\u30D3\u30B9\u30DD\u30FC\u30C8","port_busy":"%v \u30DD\u30FC\u30C8\u304C\u5360\u6709\u3055\u308C\u3066\u3044\u307E\u3059\u3002\u30E9\u30F3\u30C0\u30E0\u30DD\u30FC\u30C8\u3092\u8A66\u3057\u307E\u3059","print_all_ip":"\u3059\u3079\u3066\u306E\u30ED\u30FC\u30AB\u30EBIP\u3092\u8868\u793A","print_config":"\u30D7\u30ED\u30D5\u30A1\u30A4\u30EB\u51FA\u529B","re_enter_password":"\u30D1\u30B9\u30EF\u30FC\u30C9\u518D\u5165\u529B","reading_url_maybe":"\u8AAD\u66F8\u30EA\u30F3\u30AF\uFF1A","reg_file_hint":"Windows\u30B7\u30B9\u30C6\u30E0\u306F\u3053\u306E\u30D5\u30A1\u30A4\u30EB\u3092\u30A4\u30F3\u30DD\u30FC\u30C8\u3057\u3066\u53F3\u30AF\u30EA\u30C3\u30AF\u30E1\u30CB\u30E5\u30FC\u3092\u767B\u9332\u3057\u307E\u3059\u3002","requires_login":"\u30ED\u30B0\u30A4\u30F3\u8A8D\u8A3C\u3092\u6709\u52B9\u306B\u3059\u308B","rescan":"\u518D\u30B9\u30AD\u30E3\u30F3","save":"\u4FDD\u5B58","delete":"\u524A\u9664","save_config_file":"\u8A2D\u5B9A\u306E\u4FDD\u5B58\uFF1A","scan_archive_error":"\u30D5\u30A1\u30A4\u30EB\u30B9\u30AD\u30E3\u30F3\u30A8\u30E9\u30FC","scan_error":"\u30B9\u30AD\u30E3\u30F3\u30A8\u30E9\u30FC\uFF1A","scan_ing":"\u30D5\u30A1\u30A4\u30EB\u3092\u30B9\u30AD\u30E3\u30F3\u4E2D","scan_pdf":"PDF\u3092\u30B9\u30AD\u30E3\u30F3\u3059\u308B\uFF1A","scan_start_hint":"\u30B9\u30AD\u30E3\u30F3\u958B\u59CB\uFF1A","scroll_template":"\u30C7\u30D5\u30A9\u30EB\u30C8\uFF1A\u30B9\u30AF\u30ED\u30FC\u30EB\u30E2\u30FC\u30C9","short_description":"\u30B7\u30F3\u30D7\u30EB\u306A\u30B3\u30DF\u30C3\u30AF\u30D6\u30C3\u30AF\u30EA\u30FC\u30C0\u30FC\u3067\u3059\u3002","shutdown_hint":"\u7D42\u4E86\u51E6\u7406\u4E2D\u3002\u518D\u5EA6Ctrl+C\u3067\u5F37\u5236\u7D42\u4E86\u3057\u307E\u3059\u3002","sketch_count_seconds":"\u30B9\u30B1\u30C3\u30C1\u30E2\u30FC\u30C9\u306E\u30AB\u30A6\u30F3\u30C8\u30C0\u30A6\u30F3\u79D2\u6570","sketch_template":"\u30C7\u30D5\u30A9\u30EB\u30C8\uFF1A\u30B9\u30B1\u30C3\u30C1\u30E2\u30FC\u30C9","skip_path":"\u30B9\u30AD\u30C3\u30D7\u30D1\u30B9\uFF1A","settings_page":"\u8A2D\u5B9A\u30DA\u30FC\u30B8","upload_page":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u30DA\u30FC\u30B8","404notfound":"404 \u30DA\u30FC\u30B8\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093","sort":"\u753B\u50CF\u306E\u4E26\u3079\u66FF\u3048\u30EB\u30FC\u30EB\uFF08none\u3001name\u3001time\uFF09","sort_by_name":"\u753B\u50CF\u3092\u30D5\u30A1\u30A4\u30EB\u540D\u3067\u4E26\u3079\u66FF\u3048","sort_by_time":"\u753B\u50CF\u3092\u6700\u7D42\u66F4\u65B0\u6642\u9593\u3067\u4E26\u3079\u66FF\u3048","start_clear_file":"- \u4E00\u6642\u30D5\u30A9\u30EB\u30C0\u306E\u6383\u9664\u3092\u958B\u59CB\u3057\u307E\u3059 -","start_extract":"\u89E3\u51CD\u3092\u958B\u59CB\u3057\u307E\u3059\uFF1A","start_in_background":"\u30D0\u30C3\u30AF\u30B0\u30E9\u30A6\u30F3\u30C9\u3067\u5B9F\u884C","start_ls":"\u5727\u7E2E\u30D5\u30A1\u30A4\u30EB\u306E\u30B9\u30AD\u30E3\u30F3\u3092\u958B\u59CB\uFF1A","static_file_mode":"\u9759\u7684\u30D5\u30A1\u30A4\u30EB\u30E2\u30FC\u30C9\u3092\u6709\u52B9\u306B\u3059\u308B\u304B\u3069\u3046\u304B\u3002\u9759\u7684\u30E2\u30FC\u30C9\u3067\u306F\u3001\u3059\u3079\u3066\u306E\u753B\u50CF\u3068\u30B9\u30AF\u30EA\u30D7\u30C8\u304C \u30D1\u30C3\u30B1\u30FC\u30B8\u3055\u308C\u3001\u5358\u4E00\u306EHTML\u30D5\u30A1\u30A4\u30EB\u3068\u3057\u3066\u76F4\u63A5\u4FDD\u5B58\u3067\u304D\u307E\u3059\uFF08\u958B\u767A\u4E2D\u3001\u30B9\u30AF\u30ED\u30FC\u30EB\u30E2\u30FC\u30C9\uFF0B\u7121\u9650\u30B9\u30AF\u30ED\u30FC\u30EB\u306E\u95B2\u89A7\u30DA\u30FC\u30B8\u3067\u306E\u307F\u6709\u52B9\uFF09\u3002","stop_background":"\u30D0\u30C3\u30AF\u30B0\u30E9\u30A6\u30F3\u30C9\u3092\u505C\u6B62","target_path":"\u5BFE\u8C61\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA","temp_folder_create_error":"\u4E00\u6642\u30D5\u30A9\u30EB\u30C0\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F\u3002","temp_folder_error":"\u4E00\u6642\u30D5\u30A9\u30EB\u30C0\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F\u3002","temp_folder_path":"\u30AD\u30E3\u30C3\u30B7\u30E5\u30D5\u30A9\u30EB\u30C0\u30FC\uFF1A","template":"\u30C7\u30D5\u30A9\u30EB\u30C8\u30C6\u30F3\u30D7\u30EC\u30FC\u30C8(scroll,flip,sketch)","timeout":"\u30BF\u30A4\u30E0\u30A2\u30A6\u30C8(\u5206)","tls_crt":"TLS/SSL \u8A3C\u660E\u66F8\u30D5\u30A1\u30A4\u30EB\u30D1\u30B9","tls_enable":"TLS/SSL\u3092\u6709\u52B9\u306B\u3059\u308B","auto_https_cert":"HTTPS\u8A3C\u660E\u66F8\u3092\u81EA\u52D5\u3067\u7533\u8ACB\u30FB\u767A\u884C\uFF08Let\'s Encrypt\uFF09","port443_busy_disable_auto_tls":"\u30DD\u30FC\u30C8 443 \u304C\u4F7F\u7528\u4E2D\u306E\u305F\u3081\u3001Auto TLS \u3092\u7121\u52B9\u5316\u3057\u307E\u3059\u3002","auto_tls_disabled_invalid_domain":"Auto TLS \u3092\u52D5\u4F5C\u3055\u305B\u308B\u306B\u306F\u6709\u52B9\u306A\u30C9\u30E1\u30A4\u30F3\u304C\u5FC5\u8981\u3067\u3059\u3002Auto TLS \u3092\u7121\u52B9\u5316\u3057\u307E\u3057\u305F\u3002","auto_tls_disabled_custom_cert_set":"\u30AB\u30B9\u30BF\u30E0\u8A3C\u660E\u66F8\u304C\u8A2D\u5B9A\u3055\u308C\u3066\u3044\u308B\u305F\u3081\u3001Auto TLS \u3092\u7121\u52B9\u5316\u3057\u307E\u3057\u305F\u3002","auto_tls_disabled_lan_access_off":"LAN \u30A2\u30AF\u30BB\u30B9\u304C\u7121\u52B9\u306E\u305F\u3081\u3001Auto TLS \u3092\u7121\u52B9\u5316\u3057\u307E\u3057\u305F\u3002","tls_key":"TLS/SSL \u30AD\u30FC\u30D5\u30A1\u30A4\u30EB\u30D1\u30B9","un_archive_error":"\u30D5\u30A1\u30A4\u30EB\u89E3\u51CD\u30A8\u30E9\u30FC","unable_to_extract_images_from_pdf":"PDF\u304B\u3089\u753B\u50CF\u3092\u62BD\u51FA\u3067\u304D\u307E\u305B\u3093\u3067\u3057\u305F\u3002","unmarshal_config_file_error":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u3067\u30A8\u30E9\u30FC\u304C\u767A\u751F\u3057\u307E\u3057\u305F\u3002\u30D5\u30A9\u30FC\u30DE\u30C3\u30C8\u3092\u78BA\u8A8D\u3057\u3066\u304F\u3060\u3055\u3044\u3002","unsupported_extract":"\u89E3\u51CD\u306B\u306F\u5BFE\u5FDC\u3057\u3066\u3044\u307E\u305B\u3093\uFF1A","unsupported_file_type":"\u30B5\u30DD\u30FC\u30C8\u3055\u308C\u3066\u3044\u306A\u3044\u30D5\u30A1\u30A4\u30EB\u30BF\u30A4\u30D7\uFF1A","upload_disable_hint":"\u30D5\u30A1\u30A4\u30EB\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u6A5F\u80FD\u306F\u7121\u52B9\u3067\u3059","upload_path":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u3055\u308C\u305F\u30D5\u30A1\u30A4\u30EB\u306E\u4FDD\u5B58\u5148","username":"\u30E6\u30FC\u30B6\u30FC\u540D","web_server_error":"\u30A6\u30A7\u30D6\u30B5\u30FC\u30D3\u30B9\u306E\u958B\u59CB\u306B\u5931\u6557\u3057\u307E\u3057\u305F","webp_command":"webp-server\u30B3\u30DE\u30F3\u30C9\u3001\u307E\u305F\u306Fwebp-server\u5B9F\u884C\u53EF\u80FD\u30D1\u30B9","webp_quality":"webp\u5727\u7E2E\u54C1\u8CEA","webp_server_error":"Webp\u5909\u63DB\u30B5\u30FC\u30D3\u30B9\u3092\u8D77\u52D5\u3067\u304D\u307E\u305B\u3093\u3002\u30B3\u30DE\u30F3\u30C9\u5F62\u5F0F\u3092\u78BA\u8A8D\u3057\u3001PATH\u306Bwebp-server\u5B9F\u884C\u30D5\u30A1\u30A4\u30EB\u304C\u5B58\u5728\u3059\u308B\u3053\u3068\u3092\u78BA\u8A8D\u3057\u3066\u304F\u3060\u3055\u3044\u3002","webp_server_start":"WEBP\u5909\u63DB\u30B5\u30FC\u30D3\u30B9\u3092\u958B\u59CB\u3057\u307E\u3057\u305F","webp_setting_error":"webp\u4FDD\u5B58\u30A8\u30E9\u30FC","webp_setting_save_completed":"webp\u8A2D\u5B9A\u306E\u4FDD\u5B58\u306B\u6210\u529F\u3057\u307E\u3057\u305F\u3002","websocket_error":"websocket\u30A8\u30E9\u30FC\uFF1A","websocket_messages":"WebSocket\u30E1\u30C3\u30BB\u30FC\u30B8\uFF1A","zip_encode":"ZIP\u30D5\u30A1\u30A4\u30EB\u306E\u30A8\u30F3\u30B3\u30FC\u30C9\u3092\u624B\u52D5\u6307\u5B9A\uFF08gbk, shiftjis\u306A\u3069\uFF09","lang":"\u30A4\u30F3\u30BF\u30FC\u30D5\u30A7\u30FC\u30B9\u8A00\u8A9E\u8A2D\u5B9A\uFF08auto\u3001zh\u3001en\u3001ja\uFF09\u3001\u30C7\u30D5\u30A9\u30EB\u30C8\u306Fauto\uFF08\u81EA\u52D5\u691C\u51FA\uFF09","CacheDir":"\u30ED\u30FC\u30AB\u30EB\u30AD\u30E3\u30C3\u30B7\u30E5\u306E\u5834\u6240","CacheDir_Description":"\u30ED\u30FC\u30AB\u30EB\u30AD\u30E3\u30C3\u30B7\u30E5\u306E\u5834\u6240\u3001\u30C7\u30D5\u30A9\u30EB\u30C8\u306E\u30B7\u30B9\u30C6\u30E0\u4E00\u6642\u30D5\u30A9\u30EB\u30C0\u30FC\u3002","CertFile":"\u8A3C\u660E\u66F8\u30D5\u30A1\u30A4\u30EB","CertFile_Description":"TLS/SSL \u8A3C\u660E\u66F8\u30D5\u30A1\u30A4\u30EB\u306E\u30D1\u30B9 (\u30C7\u30D5\u30A9\u30EB\u30C8: \\"~/.config/.comigo/cert.crt\\")","ClearCacheExit":"\u7D42\u4E86\u6642\u306B\u30AF\u30EA\u30FC\u30F3\u30A2\u30C3\u30D7","ClearCacheExit_Description":"\u30D7\u30ED\u30B0\u30E9\u30E0\u3092\u7D42\u4E86\u3059\u308B\u3068\u304D\u306F\u3001Web \u753B\u50CF\u30AD\u30E3\u30C3\u30B7\u30E5\u3092\u30AF\u30EA\u30A2\u3057\u307E\u3059\u3002","ClearDatabaseWhenExit":"\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u30D6\u30C3\u30AF\u3092\u30AF\u30EA\u30A2\u3059\u308B","ClearDatabaseWhenExit_Description":"\u30ED\u30FC\u30AB\u30EB \u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u304C\u6709\u52B9\u306B\u306A\u3063\u3066\u3044\u308B\u5834\u5408\u3001\u30B9\u30AD\u30E3\u30F3\u306E\u5B8C\u4E86\u5F8C\u306B\u5B58\u5728\u3057\u306A\u3044\u66F8\u7C4D\u306F\u6D88\u53BB\u3055\u308C\u307E\u3059\u3002","ConfigManager":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u7BA1\u7406","ConfigManagerDeleteSuccess":"\u8A2D\u5B9A\u304C\u524A\u9664\u3055\u308C\u307E\u3057\u305F\u3002","ConfigManagerDescription":"Save\u3092\u30AF\u30EA\u30C3\u30AF\u3059\u308B\u3068\u3001\u73FE\u5728\u306E\u8A2D\u5B9A\u304C\u30B5\u30FC\u30D0\u30FC\u306B\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u3055\u308C\u3001\u65E2\u5B58\u306E\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u304C\u4E0A\u66F8\u304D\u3055\u308C\u307E\u3059\u3002","ConfigManagerSaveHint":"\u3059\u3067\u306B\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u304C\u3042\u308B\u306E\u3067\u3001\u4FDD\u5B58\u5834\u6240\u3092\u5909\u66F4\u3057\u3066\u304F\u3060\u3055\u3044\u3002","ConfigManagerSaveSuccess":"\u8A2D\u5B9A\u304C\u4FDD\u5B58\u3055\u308C\u307E\u3057\u305F\u3002","ConfigSaveTo":"Config File Location","Debug":"\u30C7\u30D0\u30C3\u30B0\u30E2\u30FC\u30C9\u3092\u30AA\u30F3\u306B\u3059\u308B","Debug_Description":"Debug\u3092\u6709\u52B9\u306B\u3059\u308B\u3068\u3001\u3088\u308A\u591A\u304F\u306E\u30C7\u30D0\u30C3\u30B0\u60C5\u5831\u304C\u51FA\u529B\u3055\u308C\u3001\u672A\u5B8C\u6210\u306E\u96A0\u3057\u6A5F\u80FD\u306B\u95A2\u3059\u308B\u8A2D\u5B9A\u3092\u78BA\u8A8D\u3067\u304D\u307E\u3059\u3002","DisableLAN":"LAN\u5171\u6709\u3092\u7121\u52B9\u306B\u3059\u308B","DisableLAN_Description":"\u3053\u306E\u30DE\u30B7\u30F3\u3060\u3051\u3067\u8AAD\u3080\u3001\u5916\u90E8\u306B\u306F\u5171\u6709\u3057\u307E\u305B\u3093\u3002","EnableDatabase":"\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u3092\u6709\u52B9\u306B\u3059\u308B","EnableDatabase_Description":"\u30ED\u30FC\u30AB\u30EB\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u306B\u30B9\u30AD\u30E3\u30F3\u3057\u305F\u66F8\u7C4D\u30C7\u30FC\u30BF\u3092\u4FDD\u5B58\u3067\u304D\u308B\u3088\u3046\u306B\u3057\u307E\u3059\u3002","EnableFrpcServer":"Frpc\u30B5\u30FC\u30D0\u30FC\u3092\u6709\u52B9\u306B\u3059\u308B","RequiresLogin":"\u30ED\u30B0\u30A4\u30F3\u4FDD\u8B77","RequiresLogin_Description":"\u30ED\u30B0\u30A4\u30F3\u3092\u6709\u52B9\u306B\u3059\u308B\u304B\u3069\u3046\u304B\u3002","EnableTLS":"Enable TLS","EnableTLS_Description":"HTTPS \u30D7\u30ED\u30C8\u30B3\u30EB\u3092\u6709\u52B9\u306B\u3059\u308B\u304B\u3069\u3046\u304B\u3002\u8A3C\u660E\u66F8\u306F\u30AD\u30FC\u30D5\u30A1\u30A4\u30EB\u306B\u8A2D\u5B9A\u3059\u308B\u5FC5\u8981\u304C\u3042\u308A\u307E\u3059\u3002","EnableUpload":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u6A5F\u80FD\u3092\u6709\u52B9\u306B\u3059\u308B","EnableUpload_Description":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u6A5F\u80FD\u3092\u6709\u52B9\u306B\u3057\u307E\u3059\u3002","EnablePlugin":"\u30D7\u30E9\u30B0\u30A4\u30F3\u30B7\u30B9\u30C6\u30E0\u3092\u6709\u52B9\u306B\u3059\u308B","EnablePlugin_Description":"\u30D7\u30E9\u30B0\u30A4\u30F3\u30B7\u30B9\u30C6\u30E0\u3092\u6709\u52B9\u306B\u3057\u3001\u30DA\u30FC\u30B8\u306B\u30AB\u30B9\u30BF\u30E0\u306EHTML\u3001CSS\u3001JavaScript\u30B3\u30FC\u30C9\u3092\u633F\u5165\u3067\u304D\u308B\u3088\u3046\u306B\u3057\u307E\u3059\u3002","plugins_config":"\u30D7\u30E9\u30B0\u30A4\u30F3\u30B7\u30B9\u30C6\u30E0","buildin_plugins":"\u7D44\u307F\u8FBC\u307F\u30D7\u30E9\u30B0\u30A4\u30F3","enabled_plugins":"\u6709\u52B9\u306A\u30D7\u30E9\u30B0\u30A4\u30F3","auto_flip_drag":"\u30C9\u30E9\u30C3\u30B0","auto_flip_drag_panel":"\u30D1\u30CD\u30EB\u3092\u30C9\u30E9\u30C3\u30B0","auto_flip_start_flip":"\u81EA\u52D5\u3081\u304F\u308A\u3092\u958B\u59CB","auto_flip_start_scroll":"\u81EA\u52D5\u30B9\u30AF\u30ED\u30FC\u30EB\u3092\u958B\u59CB","auto_flip_pause_flip":"\u81EA\u52D5\u3081\u304F\u308A\u3092\u4E00\u6642\u505C\u6B62","auto_flip_pause_scroll":"\u81EA\u52D5\u30B9\u30AF\u30ED\u30FC\u30EB\u3092\u4E00\u6642\u505C\u6B62","auto_flip_interval":"\u9593\u9694:","auto_flip_seconds":"\u79D2","auto_scroll_distance":"\u30B9\u30AF\u30ED\u30FC\u30EB\u8DDD\u96E2:","plugin_name_clock":"\u6642\u8A08\u30D7\u30E9\u30B0\u30A4\u30F3","plugin_name_auto_flip":"\u81EA\u52D5\u3081\u304F\u308A\u30D7\u30E9\u30B0\u30A4\u30F3","plugin_name_auto_scroll":"\u81EA\u52D5\u30B9\u30AF\u30ED\u30FC\u30EB\u30D7\u30E9\u30B0\u30A4\u30F3","plugin_name_comigo_xyz":"Comigo.xyz\u30D7\u30E9\u30B0\u30A4\u30F3","plugin_name_sample":"\u30B5\u30F3\u30D7\u30EB\u30D7\u30E9\u30B0\u30A4\u30F3","plugin_name_sketch_practice":"\u30B9\u30B1\u30C3\u30C1\u7DF4\u7FD2\u30D7\u30E9\u30B0\u30A4\u30F3","sketch_practice_start":"\u30B9\u30B1\u30C3\u30C1\u7DF4\u7FD2\u3092\u958B\u59CB","sketch_practice_pause":"\u30B9\u30B1\u30C3\u30C1\u7DF4\u7FD2\u3092\u4E00\u6642\u505C\u6B62","sketch_practice_countdown":"\u30AB\u30A6\u30F3\u30C8\u30C0\u30A6\u30F3","ExcludePath":"\u30D1\u30B9\u3092\u9664\u5916\u3059\u308B","ExcludePath_Description":"\u672C\u3092\u30B9\u30AD\u30E3\u30F3\u3059\u308B\u3068\u304D\u306B\u3001\u9664\u5916\u3059\u308B\u5FC5\u8981\u304C\u3042\u308B\u30D5\u30A1\u30A4\u30EB\u307E\u305F\u306F\u30D5\u30A9\u30EB\u30C0\u30FC\u306E\u540D\u524D","FrpClientConfig":"FrpClient\u8A2D\u5B9A","GenerateBookMetadata":"\u66F8\u7C4D\u306E\u30E1\u30BF\u30C7\u30FC\u30BF\u3092\u751F\u6210\u3059\u308B","GenerateMetaData":"\u30E1\u30BF\u30C7\u30FC\u30BF\u306E\u751F\u6210","GenerateMetaData_Description":"\u66F8\u7C4D\u306E\u30E1\u30BF\u30C7\u30FC\u30BF\u3092\u751F\u6210\u3057\u307E\u3059\u3002\u73FE\u5728\u306F\u6709\u52B9\u3067\u306F\u3042\u308A\u307E\u305B\u3093\u3002","HomeDirectory":"\u30DB\u30FC\u30E0\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA","Host":"\u30C9\u30E1\u30A4\u30F3\u540D","Host_Description":"QR\u30B3\u30FC\u30C9\u3067\u8868\u793A\u3055\u308C\u308B\u30DB\u30B9\u30C8\u540D\u3092\u30AB\u30B9\u30BF\u30DE\u30A4\u30BA\u3057\u307E\u3059\u3002\\n\u30C7\u30D5\u30A9\u30EB\u30C8\u306F\u30CD\u30C3\u30C8\u30EF\u30FC\u30AF \u30AB\u30FC\u30C9\u306E IP \u3067\u3059\u3002","KeyFile":"\u30AD\u30FC\u30D5\u30A1\u30A4\u30EB","KeyFile_Description":"TLS/SSL \u30AD\u30FC \u30D5\u30A1\u30A4\u30EB \u30D1\u30B9 (\u30C7\u30D5\u30A9\u30EB\u30C8: \\"~/.config/.comigo/key.key\\")","StoreUrls":"\u30E9\u30A4\u30D6\u30E9\u30EA\u30D5\u30A9\u30EB\u30C0\u30FC","StoreUrls_Description":"\u66F8\u5EAB\u30D5\u30A9\u30EB\u30C0\u3002\u7D76\u5BFE\u30D1\u30B9\u3068\u76F8\u5BFE\u30D1\u30B9\u306B\u5BFE\u5FDC\u3002\u76F8\u5BFE\u30D1\u30B9\u306F\u73FE\u5728\u306E\u4F5C\u696D\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3092\u57FA\u6E96\u3068\u3057\u307E\u3059\u3002<br>SFTP\u30FBSMB \u30EA\u30E2\u30FC\u30C8\u30D1\u30B9\u306F\u8A66\u9A13\u7684\u5BFE\u5FDC\u3002\u5F62\u5F0F\u4F8B:<br>sftp://user:pass@192.168.1.1:22/some/path <br>smb://guest@192.168.1.1:445/some/path","LogFileName":"\u30ED\u30B0\u30D5\u30A1\u30A4\u30EB\u540D","LogFileName_Description":"\u30ED\u30B0\u30D5\u30A1\u30A4\u30EB\u540D","LogFilePath":"\u30ED\u30B0\u306E\u4FDD\u5B58\u5834\u6240","LogFilePath_Description":"\u30ED\u30B0\u30D5\u30A1\u30A4\u30EB\u306E\u4FDD\u5B58\u5834\u6240","LogToFile":"\u30ED\u30B0\u3092\u30ED\u30FC\u30AB\u30EB\u306B\u8A18\u9332\u3059\u308B","LogToFile_Description":"\u30D7\u30ED\u30B0\u30E9\u30E0\u30ED\u30B0\u3092\u30ED\u30FC\u30AB\u30EB\u30D5\u30A1\u30A4\u30EB\u306B\u4FDD\u5B58\u3059\u308B\u304B\u3069\u3046\u304B\u3002\u30C7\u30D5\u30A9\u30EB\u30C8\u3067\u306F\u4FDD\u5B58\u3055\u308C\u307E\u305B\u3093\u3002","MaxScanDepth":"\u6700\u5927\u30B9\u30AD\u30E3\u30F3\u6DF1\u5EA6","MaxScanDepth_Description":"\u6700\u5927\u30B9\u30AD\u30E3\u30F3\u6DF1\u5EA6\u3002\\n\u6DF1\u3055\u3092\u8D85\u3048\u308B\u30D5\u30A1\u30A4\u30EB\u306F\u30B9\u30AD\u30E3\u30F3\u3055\u308C\u307E\u305B\u3093\u3002\\n\u73FE\u5728\u306E\u5B9F\u884C\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u304C\u30D9\u30FC\u30B9\u3068\u306A\u308A\u307E\u3059\u3002","MinImageNum":"\u6700\u4F4E\u679A\u6570","MinImageNum_Description":"\u672C\u3068\u307F\u306A\u3055\u308C\u308B\u306B\u306F\u3001\u5C11\u306A\u304F\u3068\u3082\u6570\u679A\u306E\u753B\u50CF\u304C\u542B\u307E\u308C\u3066\u3044\u308B\u5FC5\u8981\u304C\u3042\u308A\u307E\u3059\u3002","OpenBrowser":"\u30D6\u30E9\u30A6\u30B6\u3092\u958B\u304F","OpenBrowser_Description":"windows\u306E\u30C7\u30D5\u30A9\u30EB\u30C8\u306Ftrue\u3001\u305D\u306E\u4ED6\u306E\u30D7\u30E9\u30C3\u30C8\u30D5\u30A9\u30FC\u30E0\u306E\u30C7\u30D5\u30A9\u30EB\u30C8\u306Ffalse\u3002","Password":"\u30D1\u30B9\u30EF\u30FC\u30C9","Password_Description":"\u30ED\u30B0\u30A4\u30F3\u306B\u4F7F\u7528\u3055\u308C\u308B\u30D1\u30B9\u30EF\u30FC\u30C9\u3002","ReEnterPassword":"\u30D1\u30B9\u30EF\u30FC\u30C9\u3092\u518D\u5165\u529B","ReEnterPassword_Description":"\u30D1\u30B9\u30EF\u30FC\u30C9\u3092\u518D\u5165\u529B\u3057\u307E\u3059\u3002","Port":"\u30DD\u30FC\u30C8","Port_Description":"Web \u30B5\u30FC\u30D3\u30B9 \u30DD\u30FC\u30C8\u3002","PrintAllPossibleQRCode":"\u305D\u306E\u4ED6\u306E QR \u30B3\u30FC\u30C9","ProgramDirectory":"\u30D7\u30ED\u30B0\u30E9\u30E0\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA","StartFrpClientInBackground":"Frp\u30AF\u30E9\u30A4\u30A2\u30F3\u30C8\u3092\u958B\u59CB\u3059\u308B","SupportFileType":"\u30B5\u30DD\u30FC\u30C8\u3055\u308C\u3066\u3044\u308B\u5727\u7E2E\u30D1\u30C3\u30B1\u30FC\u30B8","SupportFileType_Description":"\u30D5\u30A1\u30A4\u30EB\u3092\u30B9\u30AD\u30E3\u30F3\u3059\u308B\u3068\u304D\u306B\u3001\u305D\u308C\u3092\u30B9\u30AD\u30C3\u30D7\u3059\u308B\u304B\u3001\u30D6\u30C3\u30AF\u51E6\u7406\u306E\u30D5\u30A1\u30A4\u30EB\u306E\u63A5\u5C3E\u8F9E\u3068\u3057\u3066\u30AB\u30A6\u30F3\u30C8\u3059\u308B\u304B\u3092\u6C7A\u5B9A\u3059\u308B\u305F\u3081\u306B\u4F7F\u7528\u3055\u308C\u307E\u3059\u3002","SupportMediaType":"\u30B5\u30DD\u30FC\u30C8\u3055\u308C\u3066\u3044\u308B\u753B\u50CF\u30D5\u30A1\u30A4\u30EB","SupportMediaType_Description":"\u5727\u7E2E\u30D1\u30C3\u30B1\u30FC\u30B8\u3092\u30B9\u30AD\u30E3\u30F3\u3059\u308B\u3068\u304D\u306B\u753B\u50CF\u306E\u6570\u3092\u30AB\u30A6\u30F3\u30C8\u3059\u308B\u305F\u3081\u306B\u4F7F\u7528\u3055\u308C\u308B\u753B\u50CF\u30D5\u30A1\u30A4\u30EB\u306E\u63A5\u5C3E\u8F9E","Timeout":"\u6709\u52B9\u671F\u9650","TimeoutLimitForScan":"\u30B9\u30AD\u30E3\u30F3\u30BF\u30A4\u30E0\u30A2\u30A6\u30C8","TimeoutLimitForScan_Description":"\u30D5\u30A1\u30A4\u30EB\u306E\u30B9\u30AD\u30E3\u30F3\u6642\u306B\u6570\u79D2\u4EE5\u4E0A\u304B\u304B\u308B\u5834\u5408\u3001\u5927\u304D\u3059\u304E\u308B\u30D5\u30A1\u30A4\u30EB\u3067\u30B9\u30BF\u30C3\u30AF\u3059\u308B\u3053\u3068\u3092\u907F\u3051\u308B\u305F\u3081\u306B\u30D5\u30A1\u30A4\u30EB\u306E\u30B9\u30AD\u30E3\u30F3\u3092\u4E2D\u6B62\u3057\u307E\u3059\u3002","Timeout_Description":"\u30ED\u30B0\u30A4\u30F3\u3092\u6709\u52B9\u306B\u3057\u305F\u5F8C\u306E Cookie \u306E\u6709\u52B9\u671F\u9650\u3002\u5358\u4F4D\u306F\u5206\u3067\u3059\u3002","UseCache":"\u30ED\u30FC\u30AB\u30EB\u753B\u50CF\u30AD\u30E3\u30C3\u30B7\u30E5","UseCache_Description":"\u30ED\u30FC\u30AB\u30EB\u753B\u50CF\u89E3\u51CD\u30AD\u30E3\u30C3\u30B7\u30E5\u3092\u6709\u52B9\u306B\u3057\u307E\u3059\u3002\u30C7\u30D5\u30A9\u30EB\u30C8\u3067\u306F\u7121\u52B9\u3067\u3059\u3002","Username":"\u30E6\u30FC\u30B6\u30FC\u540D","Username_Description":"\u30ED\u30B0\u30A4\u30F3 \u30A4\u30F3\u30BF\u30FC\u30D5\u30A7\u30A4\u30B9\u306B\u5FC5\u8981\u306A\u30E6\u30FC\u30B6\u30FC\u540D\u3002","WorkingDirectory":"\u4F5C\u696D\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA","ZipFileTextEncoding":"UTF-8\u3067\u306F\u3042\u308A\u307E\u305B\u3093","ZipFileTextEncoding_Description":"utf-8 \u4EE5\u5916\u3067\u30A8\u30F3\u30B3\u30FC\u30C9\u3055\u308C\u305F ZIP \u30D5\u30A1\u30A4\u30EB\u3002\u89E3\u6790\u3059\u308B\u306B\u306F\u3069\u306E\u30A8\u30F3\u30B3\u30FC\u30C9\u3092\u4F7F\u7528\u3059\u308B\u5FC5\u8981\u304C\u3042\u308A\u307E\u3059\u3002\u30C7\u30D5\u30A9\u30EB\u30C8\u306EGBK\u3002","labs":"\u30E9\u30DC\u30BA","all_page_num":"\u7DCF\u30DA\u30FC\u30B8\u6570: {0}","author":"\u4F5C\u8005: {0}","auto_crop":"\u30A8\u30C3\u30B8\u30C8\u30EA\u30DF\u30F3\u30B0","auto_crop_num":"\u30C8\u30EA\u30DF\u30F3\u30B0\u95BE\u5024: ","auto_double_page":"\u81EA\u52D5\u898B\u958B\u304D\u30DA\u30FC\u30B8 (\u03B2)","auto_hide_toolbar":"\u30C4\u30FC\u30EB\u30D0\u30FC\u975E\u8868\u793A","back-to-top":"\u30DA\u30FC\u30B8\u30C8\u30C3\u30D7\u3078\u623B\u308B","back_button":"\u623B\u308B\u30DC\u30BF\u30F3","back_to_bookshelf":"\u672C\u68DA\u306B\u623B\u308B","book_shelf":"\u672C\u68DA","child_book_hint":"\u30D5\u30A9\u30EB\u30C0\u5185\u306B{0}\u518A\u306E\u66F8\u7C4D\u304C\u542B\u307E\u308C\u3066\u3044\u307E\u3059","click_to_toggle":"(\u30AF\u30EA\u30C3\u30AF\u5207\u308A\u66FF\u3048)","do_you_reset_local_settings":"\u3059\u3079\u3066\u306E\u8A2D\u5B9A\u3092\u521D\u671F\u5316\u3057\u307E\u3059\u304B\uFF1F","double_page_mode":"\u30C0\u30D6\u30EB\u30DA\u30FC\u30B8","double_page_width":"\u898B\u958B\u304D\u30DA\u30FC\u30B8\u5E45:","download_sample_config_file":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u3092\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","download_windows_reg_file":"REG\u30D5\u30A1\u30A4\u30EB\u3092\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","drop_to_upload":"\u30D5\u30A1\u30A4\u30EB\u3092\u30AF\u30EA\u30C3\u30AF\u3059\u308B\u304B\u3001\u3053\u306E\u30A8\u30EA\u30A2\u306B\u30C9\u30E9\u30C3\u30B0\uFF06\u30C9\u30ED\u30C3\u30D7\u3057\u3066\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9","energy_threshold":"\u30A8\u30CD\u30EB\u30AE\u30FC\u95BE\u5024:","epub_info":"ePub\u60C5\u5831","exit_fullscreen":"\u30D5\u30EB\u30B9\u30AF\u30EA\u30FC\u30F3\u3092\u7D42\u4E86","filesize":"\u30D5\u30A1\u30A4\u30EB\u30B5\u30A4\u30BA: {0}","flip_mode":"\u30E8\u30B3\u8AAD\u307F\u30E2\u30FC\u30C9","flip_odd_even_page":"\u30DA\u30FC\u30B8\u9593\u306E\u30DE\u30C3\u30C1\u30F3\u30B0\u3092\u5909\u66F4","flip_odd_even_page_hint":"\u30DA\u30FC\u30B8\u5185\u5BB9\u304C\u4E00\u81F4\u3057\u306A\u3044\u5834\u5408\u306F\u3001\u30AF\u30EA\u30C3\u30AF\u3057\u3066\u4FEE\u6B63\u3092\u8A66\u307F\u3066\u304F\u3060\u3055\u3044","found_read_history":"\u8AAD\u66F8\u5C65\u6B74\u304C\u898B\u3064\u304B\u308A\u307E\u3057\u305F","from_interrupt":"\u9014\u4E2D\u304B\u3089\u518D\u958B","full_screen_hint":"\u30D5\u30EB\u30B9\u30AF\u30EA\u30FC\u30F3\u30DC\u30BF\u30F3","fullscreen":"\u30D5\u30EB\u30B9\u30AF\u30EA\u30FC\u30F3","good_job_and_bye":"\u304A\u75B2\u308C\u69D8\u3067\u3057\u305F\u3002\u3055\u3088\u3046\u306A\u3089\u3002","gray_image":"\u30B0\u30EC\u30FC\u30B9\u30B1\u30FC\u30EB\u5316","hint":"Hint","hint_first_page":"\u6700\u521D\u306E\u30DA\u30FC\u30B8\u3067\u306F\u9032\u3081\u307E\u305B\u3093","hint_last_page":"\u3053\u308C\u304C\u6700\u5F8C\u306E\u30DA\u30FC\u30B8\u3067\u3059","hour":"\u6642\u9593","compress_image":"\u753B\u50CF\u5727\u7E2E","infinite_dropdown":"\u7121\u9650\u30C9\u30ED\u30C3\u30D7\u30C0\u30A6\u30F3\u30E2\u30FC\u30C9","interval":"\u9593\u9694:","manga_mode":"\u30DE\u30F3\u30AC(\u53F3\u5411\u304D)","load_all_pages":"\u5168\u30DA\u30FC\u30B8\u3092\u8AAD\u307F\u8FBC\u3080","skip_and_load_full":"\u65E2\u8AAD\u30DA\u30FC\u30B8\u304C\u8AAD\u307F\u8FBC\u307E\u308C\u3066\u3044\u307E\u305B\u3093\uFF08\u6700\u521D\u306E %d \u30DA\u30FC\u30B8\uFF09\u3001\u4E0B\u306E\u30DC\u30BF\u30F3\u3092\u30AF\u30EA\u30C3\u30AF\u3057\u3066\u3059\u3079\u3066\u3092\u8AAD\u307F\u8FBC\u307F\u307E\u3059","completed_and_load_full":"\u6700\u5F8C\u306E\u30DA\u30FC\u30B8\u307E\u3067\u8AAD\u307E\u308C\u307E\u3057\u305F\u3002\u4E0B\u306E\u30DC\u30BF\u30F3\u3092\u30AF\u30EA\u30C3\u30AF\u3057\u3066\u5168\u30DA\u30FC\u30B8\u3092\u8AAD\u307F\u8FBC\u307F\u307E\u3059\u3002","load_from_interrupt":"%d \u30DA\u30FC\u30B8\u304B\u3089\u8AAD\u307F\u8FBC\u307F\u3092\u958B\u59CB\u3057\u307E\u3059\u304B\uFF1F","login_success_hint":"\u30ED\u30B0\u30A4\u30F3\u306B\u6210\u529F\u3057\u307E\u3057\u305F\u3002\u524D\u306E\u30DA\u30FC\u30B8\u306B\u623B\u308A\u307E\u3059","logout":"\u30ED\u30B0\u30A2\u30A6\u30C8","margin_bottom_on_scroll_mode":"\u30DA\u30FC\u30B8\u9593\u306E\u4F59\u767D:","margin_on_scroll_mode":"\u30DA\u30FC\u30B8\u30AE\u30E3\u30C3\u30D7:","limit_width":"\u5E45\u5236\u9650:","minute":"\u5206","network":"\u30CD\u30C3\u30C8","no_book_found_hint":"\u672C\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3002\u30D5\u30A1\u30A4\u30EB\u3092\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u3057\u3066\u307F\u3066\u304F\u3060\u3055\u3044\u3002","no_support_upload_file":"\u30D5\u30A1\u30A4\u30EB\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u6A5F\u80FD\u306F\u7BA1\u7406\u8005\u306B\u3088\u308A\u7121\u52B9\u5316\u3055\u308C\u3066\u3044\u307E\u3059","not_support_fullscreen":"\u3053\u306E\u30D6\u30E9\u30A6\u30B6\u306F\u30D5\u30EB\u30B9\u30AF\u30EA\u30FC\u30F3\u3092\u30B5\u30DD\u30FC\u30C8\u3057\u3066\u3044\u307E\u305B\u3093","now_is":"\u73FE\u5728:","number_of_online_books":"\u30AA\u30F3\u30E9\u30A4\u30F3\u66F8\u7C4D\u6570\uFF1A","original_image":"\u30AA\u30EA\u30B8\u30CA\u30EB\u89E3\u50CF\u5EA6","original_pdf_link":"PDF\u30D5\u30A1\u30A4\u30EB\u306E\u30EA\u30F3\u30AF","page":"\u30DA\u30FC\u30B8","page_turning_seconds":"\u30DA\u30FC\u30B8\u3081\u304F\u308A\u6642\u9593:","scroll_fixed_pagination":"\u30DA\u30FC\u30B8\u30F3\u30B0","scroll_infinite_scroll":"\u7121\u9650\u30B9\u30AF\u30ED\u30FC\u30EB","pdf_hint_message":"\u7D14\u7C8B\u306A\u753B\u50CFPDF\u306B\u5BFE\u5FDC\u3057\u3066\u3044\u307E\u3059\u3002\u8AAD\u307F\u8FBC\u307F\u304C\u9045\u3044\u5834\u5408\u3084\u30A8\u30E9\u30FC\u304C\u767A\u751F\u3057\u305F\u5834\u5408\u306F\u3001\u4EE5\u4E0B\u3092\u8A66\u3057\u3066\u304F\u3060\u3055\u3044\uFF1A","please_enable_upload":"\u30B5\u30FC\u30D0\u30FC\u306E\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u6A5F\u80FD\u3092\u6709\u52B9\u306B\u3057\u3066\u304F\u3060\u3055\u3044","please_enter_content":"\u5185\u5BB9\u3092\u5165\u529B\u3057\u3066\u304F\u3060\u3055\u3044","qrcode_hint":"\u30AF\u30EA\u30C3\u30AF\u3057\u3066QR\u30B3\u30FC\u30C9\u3092\u8868\u793A","raw_resolution":"\u5143\u306E\u89E3\u50CF\u5EA6","re_sort_book":"\u66F8\u7C4D\u3092\u4E26\u3079\u66FF\u3048","re_sort_page":"\u30DA\u30FC\u30B8\u3092\u4E26\u3079\u66FF\u3048","reader_settings":"\u8A2D\u5B9A","reading_progress_bar":"\u8AAD\u66F8\u9032\u6357\u30D0\u30FC","refresh_page":"\u30DA\u30FC\u30B8\u3092\u66F4\u65B0","reset_local_settings":"\u30EA\u30FC\u30C0\u30FC\u8A2D\u5B9A\u3092\u30EA\u30BB\u30C3\u30C8","resort_file":"\u30D5\u30A1\u30A4\u30EB\u3092\u518D\u4E26\u3073\u66FF\u3048","comic_mode":"\u30B3\u30DF\u30C3\u30AF(\u5DE6\u5411\u304D)","save_page_num":"\u8AAD\u66F8\u9032\u6357\u3092\u4FDD\u5B58","scan_qrcode":"QR\u30B3\u30FC\u30C9\u3092\u30B9\u30AD\u30E3\u30F3:","scanned_hint":"XX\u518A\u898B\u3064\u304B\u308A\u307E\u3057\u305F\u3002\u8868\u793A\u3057\u307E\u3059\u304B\uFF1F","scroll_mode":"\u30BF\u30C6\u8AAD\u307F\u30E2\u30FC\u30C9","second":"\u79D2","select_language":"\u8A00\u8A9E\u3092\u9078\u629E","server_config":"\u30B5\u30FC\u30D0\u30FC\u8A2D\u5B9A","server_setting":"Comigo\u30B5\u30FC\u30D0\u30FC\u8A2D\u5B9A","set_back_color":"\u80CC\u666F\u8272\u3092\u8A2D\u5B9A:","set_interface_color":"UI\u30AB\u30E9\u30FC\u8A2D\u5B9A:","show_filename":"\u30D5\u30A1\u30A4\u30EB\u540D\u3092\u8868\u793A","show_file_icon":"\u30D5\u30A1\u30A4\u30EB\u30A2\u30A4\u30B3\u30F3","show_header":"\u30BF\u30A4\u30C8\u30EB\u3092\u8868\u793A","showPageNum":"\u30DA\u30FC\u30B8\u756A\u53F7\u3092\u8868\u793A","simplify_filename":"\u30D5\u30A1\u30A4\u30EB\u540D\u3092\u7C21\u7565\u5316","single_page_mode":"\u30B7\u30F3\u30B0\u30EB\u30DA\u30FC\u30B8","single_page_width":"\u30B7\u30F3\u30B0\u30EB\u30DA\u30FC\u30B8\u5E45:","sort_by_default":"\u30C7\u30D5\u30A9\u30EB\u30C8\u9806","sort_by_filename":"\u30D5\u30A1\u30A4\u30EB\u540D\u9806 (A-Z)","sort_by_filename_reverse":"\u30D5\u30A1\u30A4\u30EB\u540D\u9806 (Z-A)","sort_by_filesize":"\u30D5\u30A1\u30A4\u30EB\u30B5\u30A4\u30BA\u9806 (\u5927\u2192\u5C0F)","sort_by_filesize_reverse":"\u30D5\u30A1\u30A4\u30EB\u30B5\u30A4\u30BA\u9806 (\u5C0F\u2192\u5927)","sort_by_modify_time":"\u66F4\u65B0\u65E5\u6642\u9806 (\u65B0\u2192\u65E7)","sort_by_modify_time_reverse":"\u66F4\u65B0\u65E5\u6642\u9806 (\u65E7\u2192\u65B0)","sort_reverse":"\uFF08\u9006\u9806\uFF09","start_sketch_message":"\u30AF\u30ED\u30C3\u30AD\u30FC\u30E2\u30FC\u30C9\u304C\u958B\u59CB\u3055\u308C\u307E\u3057\u305F\u3002\u826F\u3044\u4E00\u65E5\u3092\uFF01","start_sketch_mode":"\u30AF\u30ED\u30C3\u30AD\u30FC\u30E2\u30FC\u30C9\u958B\u59CB","starting_from_beginning":"\u6700\u521D\u304B\u3089\u958B\u59CB","starting_from_beginning_hint":"\u6700\u521D\u306E\u30DA\u30FC\u30B8\u304B\u3089\u8AAD\u307F\u8FBC\u307F\u307E\u3059","stop_sketch_mode":"\u30AF\u30ED\u30C3\u30AD\u30FC\u30E2\u30FC\u30C9\u7D42\u4E86","submit":"\u9001\u4FE1","add":"\u8FFD\u52A0","success_fullscreen":"\u30D5\u30EB\u30B9\u30AF\u30EA\u30FC\u30F3\u306B\u6210\u529F\u3057\u307E\u3057\u305F","successfully_loaded_reading_progress":"\u8AAD\u66F8\u9032\u6357\u3092\u8AAD\u307F\u8FBC\u307F\u307E\u3057\u305F","sync_page":"\u30EA\u30E2\u30FC\u30C8\u30DA\u30FC\u30B8\u540C\u671F","temp_future_hint":"\u307E\u3060\u5B8C\u6210\u3057\u3066\u3044\u306A\u3044\u3044\u304F\u3064\u304B\u306E\u6A5F\u80FD\u3092\u3001\u4E00\u6642\u7684\u306B\u7F6E\u304F\u3002","test":"\u30C6\u30B9\u30C8","total_is":"\u5408\u8A08:","total_time":"\u5408\u8A08\u6642\u9593:","type_or_paste_content":"\u5165\u529B\u307E\u305F\u306F\u8CBC\u308A\u4ED8\u3051","upload_file":"\u30D5\u30A1\u30A4\u30EB\u3092\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9","width_use_fixed_value":"\u6A2A\u5E45: \u56FA\u5B9A\u5024","width_use_percent":"\u6A2A\u5E45: \u30D1\u30FC\u30BB\u30F3\u30C6\u30FC\u30B8","portrait_width_percent":"\u7E26\u753B\u9762\u5E45\uFF08\u30D1\u30FC\u30BB\u30F3\u30C8\uFF09","auto_align":"\u81EA\u52D5\u753B\u9762\u6574\u5217","swipe_turn":"\u30B9\u30EF\u30A4\u30D7\u30DA\u30FC\u30B8\u9001\u308A","login_title":"Comigo\u306B\u30ED\u30B0\u30A4\u30F3","login_subtitle":"\u30E6\u30FC\u30B6\u30FC\u540D\u3068\u30D1\u30B9\u30EF\u30FC\u30C9\u3092\u5165\u529B\u3057\u3066\u304F\u3060\u3055\u3044","login_failed":"\u30ED\u30B0\u30A4\u30F3\u306B\u5931\u6557\u3057\u307E\u3057\u305F\u3002\u30E6\u30FC\u30B6\u30FC\u540D\u3068\u30D1\u30B9\u30EF\u30FC\u30C9\u3092\u78BA\u8A8D\u3057\u3066\u304F\u3060\u3055\u3044","login_error_teapot":"\u30B5\u30FC\u30D0\u30FC\u306F\u8A8D\u8A3C\u3092\u5FC5\u8981\u3068\u3057\u307E\u305B\u3093\u3002<a class=\\"font-semibold text-blue-600\\" href=\\"/\\">\u30DB\u30FC\u30E0\u30DA\u30FC\u30B8</a> \u306B\u76F4\u63A5\u30A2\u30AF\u30BB\u30B9\u3057\u3066\u304F\u3060\u3055\u3044","logging_in":"\u30ED\u30B0\u30A4\u30F3\u4E2D...","login":"\u30ED\u30B0\u30A4\u30F3","other_information":"\u305D\u306E\u4ED6\u60C5\u5831","login_forgot_password_hint":"\u30D1\u30B9\u30EF\u30FC\u30C9\u3092\u304A\u5FD8\u308C\u3067\u3059\u304B\uFF1F\u30B7\u30B9\u30C6\u30E0\u7BA1\u7406\u8005\u306B\u304A\u554F\u3044\u5408\u308F\u305B\u304F\u3060\u3055\u3044","no_pattern":"\u7121\u5730","grid_line":"\u30B0\u30EA\u30C3\u30C9\u7DDA","grid_point":"\u30B0\u30EA\u30C3\u30C9\u70B9","mosaic":"\u30E2\u30B6\u30A4\u30AF","open_pdf_in_browser":"\u30D6\u30E9\u30A6\u30B6\u3067PDF\u3092\u958B\u304F","StaticFileMode":"\u9759\u7684\u30D5\u30A1\u30A4\u30EB\u30E2\u30FC\u30C9","StaticFileMode_Description":"\u9759\u7684\u30D5\u30A1\u30A4\u30EB\u30E2\u30FC\u30C9\u3092\u6709\u52B9\u306B\u3059\u308B\u304B\u3069\u3046\u304B\u3002\u9759\u7684\u30E2\u30FC\u30C9\u3067\u306F\u3001\u3059\u3079\u3066\u306E\u753B\u50CF\u3068\u30B9\u30AF\u30EA\u30D7\u30C8\u304C \u30D1\u30C3\u30B1\u30FC\u30B8\u3055\u308C\u3001\u5358\u4E00\u306EHTML\u30D5\u30A1\u30A4\u30EB\u3068\u3057\u3066\u76F4\u63A5\u4FDD\u5B58\u3067\u304D\u307E\u3059\uFF08\u958B\u767A\u4E2D\uFF09\u3002","confirm_logout":"\u30ED\u30B0\u30A2\u30A6\u30C8\u3057\u3066\u3082\u3088\u308D\u3057\u3044\u3067\u3059\u304B\uFF1F","confirm_reset_settings":"\u30ED\u30FC\u30AB\u30EB\u8A2D\u5B9A\u3092\u30EA\u30BB\u30C3\u30C8\u3057\u3066\u3082\u3088\u308D\u3057\u3044\u3067\u3059\u304B\uFF1F","current_dir_scope":"\u4ECA\u306E\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3067\u5B9F\u884C\u3057\u305F\u5834\u5408\uFF08\u30ED\u30FC\u30AB\u30EB\u9069\u7528\uFF09","current_user_scope":"\u30ED\u30B0\u30A4\u30F3\u30E6\u30FC\u30B6\u30FC\u306B\u5BFE\u3057\u3066\u6709\u52B9\uFF08\u30B0\u30ED\u30FC\u30D0\u30EB\u9069\u7528\uFF09","portable_binary_scope":"\u5F53\u8A72\u30D0\u30A4\u30CA\u30EA\u306B\u5BFE\u3057\u3066\u6709\u52B9\uFF08\u30DD\u30FC\u30BF\u30D6\u30EB\u30E2\u30FC\u30C9\uFF09","saveSuccessHint":"\u8A2D\u5B9A\u304C\u4FDD\u5B58\u3055\u308C\u307E\u3057\u305F\u30022\u79D2\u5F8C\u306B\u30DA\u30FC\u30B8\u304C\u81EA\u52D5\u7684\u306B\u30EA\u30ED\u30FC\u30C9\u3055\u308C\u307E\u3059\u3002","port_change_hint":"\u30DD\u30FC\u30C8\u304C\u5909\u66F4\u3055\u308C\u307E\u3057\u305F\u30022\u79D2\u5F8C\u306B\u65B0\u3057\u3044\u30DD\u30FC\u30C8\u306B\u30EA\u30C0\u30A4\u30EC\u30AF\u30C8\u3057\u307E\u3059\u3002","no_books_library_path_notice":"\u95B2\u89A7\u53EF\u80FD\u306A\u66F8\u7C4D\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3002\u30E9\u30A4\u30D6\u30E9\u30EA\u30D1\u30B9\u3092\u8A2D\u5B9A\u3057\u3066\u304F\u3060\u3055\u3044\u3002\u8A2D\u5B9A\u304C\u5B8C\u4E86\u3059\u308B\u3068\u3001\u30DA\u30FC\u30B8\u306F\u81EA\u52D5\u7684\u306B\u30EA\u30ED\u30FC\u30C9\u3055\u308C\u307E\u3059\u3002","download_raw_archive":"\u30D5\u30A1\u30A4\u30EB\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","download_portable_web_file":"HTML\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","download_as_zip":"ZIP\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","download_as_epub":"EPUB\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","clear_cache":"\u30AD\u30E3\u30C3\u30B7\u30E5\u524A\u9664","confirm_clear_cache":"\u3053\u306E\u66F8\u7C4D\u306E\u30AD\u30E3\u30C3\u30B7\u30E5\u30D5\u30A1\u30A4\u30EB\u3092\u524A\u9664\u3057\u307E\u3059\u304B\uFF1F\uFF08\u30E1\u30BF\u30C7\u30FC\u30BF\u3068\u8868\u7D19\u30AD\u30E3\u30C3\u30B7\u30E5\uFF09","clear_cache_success":"\u30AD\u30E3\u30C3\u30B7\u30E5\u3092\u524A\u9664\u3057\u307E\u3057\u305F","clear_cache_failed":"\u30AD\u30E3\u30C3\u30B7\u30E5\u306E\u524A\u9664\u306B\u5931\u6557\u3057\u307E\u3057\u305F","download":"\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","switch":"\u5207\u308A\u66FF\u3048","EnableTailscale":"Tailscale\u3092\u6709\u52B9\u5316","EnableTailscale_Description":"Tailscale\u306E\u5185\u7DB2\u900F\u904E\u6A5F\u80FD\u3092\u6709\u52B9\u306B\u3057\u307E\u3059\u3002\u521D\u56DE\u306E\u6709\u52B9\u5316\u6642\u306B\u306F\u3001Tailscale\u7BA1\u7406\u30B3\u30F3\u30BD\u30FC\u30EB\u3067\u306E\u8A8D\u8A3C\u304C\u5FC5\u8981\u3067\u3059\u3002","TailscaleHostname":"Tailscale\u30DB\u30B9\u30C8\u540D","TailscaleHostname_Description":"Tailscale\u306E\u30DB\u30B9\u30C8\u540D\u90E8\u5206\u3067\u3059\u3002\u5B8C\u5168\u306A\u30C9\u30E1\u30A4\u30F3\u306F {hostname}.example.ts.net \u306E\u3088\u3046\u306B\u306A\u308A\u307E\u3059\u3002","TailscalePort":"Tailscale\u5F85\u3061\u53D7\u3051\u30DD\u30FC\u30C8","TailscalePort_Description":"Tailscale\u306E\u5F85\u3061\u53D7\u3051\u30DD\u30FC\u30C8\u3067\u3059\u3002\u30C7\u30D5\u30A9\u30EB\u30C8\u306F443\u3067\u3001\u81EA\u52D5\u7684\u306BTLS\u304C\u6709\u52B9\u306B\u306A\u308A\u307E\u3059\u3002","FunnelTunnel":"Funnel\u30E2\u30FC\u30C9\uFF08\u30D1\u30D6\u30EA\u30C3\u30AF\u30A2\u30AF\u30BB\u30B9\uFF09","FunnelTunnel_Description":"Funnel\u30E2\u30FC\u30C9\uFF08\u30D1\u30D6\u30EA\u30C3\u30AF\u30A2\u30AF\u30BB\u30B9\uFF09\u3002\u516C\u958B\u3057\u305F\u304F\u306A\u3044\u5834\u5408\u306F\u3001\u30D1\u30B9\u30EF\u30FC\u30C9\u4FDD\u8B77\u3092\u8A2D\u5B9A\u3059\u308B\u3053\u3068\u3092\u63A8\u5968\u3057\u307E\u3059\u3002Funnel\u30E2\u30FC\u30C9\u3067\u306F443\u30018443\u300110000\u30DD\u30FC\u30C8\u306E\u307F\u4F7F\u7528\u3067\u304D\u307E\u3059\u3002","read_only_mode":"\u8AAD\u307F\u53D6\u308A\u5C02\u7528\u30E2\u30FC\u30C9","read_only_mode_description":"\u73FE\u5728\u306F\u8AAD\u307F\u53D6\u308A\u5C02\u7528\u30E2\u30FC\u30C9\u306E\u305F\u3081\u3001Web\u4E0A\u3067\u306E\u8A2D\u5B9A\u5909\u66F4\u3084\u30D5\u30A1\u30A4\u30EB\u306E\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u306F\u3067\u304D\u307E\u305B\u3093\u3002","enable_single_instance":"\u5358\u4E00\u30A4\u30F3\u30B9\u30BF\u30F3\u30B9\u30E2\u30FC\u30C9\u3092\u6709\u52B9\u306B\u3057\u3001\u540C\u6642\u306B1\u3064\u306E\u30D7\u30ED\u30B0\u30E9\u30E0\u30A4\u30F3\u30B9\u30BF\u30F3\u30B9\u306E\u307F\u304C\u5B9F\u884C\u3055\u308C\u308B\u3088\u3046\u306B\u3057\u307E\u3059","tailscale_auth_url_is":"Tailscale\u30B5\u30FC\u30D0\u30FC\u3092\u8D77\u52D5\u3059\u308B\u306B\u306F\u3001TS_AUTHKEY\u3092\u8A2D\u5B9A\u3057\u305F\u72B6\u614B\u3067\u518D\u8D77\u52D5\u3059\u308B\u304B\u3001 \u8A8D\u8A3CURL\u306B\u30A2\u30AF\u30BB\u30B9\u3057\u3066\u304F\u3060\u3055\u3044\uFF1A","tailscale_server_start":"Tailscale\u30B5\u30FC\u30D0\u30FC\u304C\u8D77\u52D5\u3057\u307E\u3057\u305F","tailscale_reading_url":"Tailscale\u8AAD\u66F8\u30EA\u30F3\u30AF\uFF1A","tailscale_not_connected_hint":"Tailscale\u304C\u63A5\u7D9A\u3055\u308C\u3066\u3044\u307E\u305B\u3093\u3002Tailscale\u30AF\u30E9\u30A4\u30A2\u30F3\u30C8\u3092\u30A4\u30F3\u30B9\u30C8\u30FC\u30EB\u3057\u3066\u30ED\u30B0\u30A4\u30F3\u3057\u3066\u304F\u3060\u3055\u3044\u3002","tailscale_not_enabled":"Tailscale\u304C\u6709\u52B9\u306B\u306A\u3063\u3066\u3044\u307E\u305B\u3093\u3002","ServerSettings":"Comigo\u30B5\u30FC\u30D0\u30FC\u8A2D\u5B9A","settings_stores":"\u66F8\u5EAB\u8A2D\u5B9A","settings_network":"\u30CD\u30C3\u30C8\u30EF\u30FC\u30AF\u8A2D\u5B9A","settings_extra":"\u5B9F\u9A13\u7684\u6A5F\u80FD","settings_log_title":"\u30EA\u30A2\u30EB\u30BF\u30A4\u30E0\u30B5\u30FC\u30D0\u30FC\u30ED\u30B0","settings_log_broadcast_placeholder":"\u30D6\u30ED\u30FC\u30C9\u30AD\u30E3\u30B9\u30C8\u30E1\u30C3\u30BB\u30FC\u30B8\u3092\u5165\u529B","settings_log_broadcast_send":"\u9001\u4FE1","settings_log_sse_connected":"\u30ED\u30B0\u30B5\u30FC\u30D0\u30FC\u306B\u63A5\u7D9A\u3057\u307E\u3057\u305F","settings_log_sse_closed":"\u63A5\u7D9A\u304C\u9589\u3058\u3089\u308C\u307E\u3057\u305F","settings_log_sse_retrying":"\u518D\u63A5\u7D9A\u4E2D...","settings_log_broadcast_default_message":"Hello ~","remote_access":"\u30EA\u30E2\u30FC\u30C8\u30A2\u30AF\u30BB\u30B9","ErrPasswordMismatch":"\u5165\u529B\u3057\u305F\u30D1\u30B9\u30EF\u30FC\u30C9\u304C\u4E00\u81F4\u3057\u307E\u305B\u3093\u3002\u518D\u5EA6\u5165\u529B\u3057\u3066\u304F\u3060\u3055\u3044","PromptSetPassword":"\u30D1\u30B9\u30EF\u30FC\u30C9\u3092\u8A2D\u5B9A\u3057\u3066\u304F\u3060\u3055\u3044","MsgLoginSettingsUpdated":"\u30ED\u30B0\u30A4\u30F3\u8A2D\u5B9A\u3092\u5909\u66F4\u3057\u307E\u3057\u305F\u3001\u30ED\u30B0\u30A4\u30F3\u30DA\u30FC\u30B8\u3078\u79FB\u52D5\u3057\u307E\u3059","CurrentPassword":"\u73FE\u5728\u306E\u30D1\u30B9\u30EF\u30FC\u30C9","AdminAccountSetup":"\u7BA1\u7406\u8005\u30A2\u30AB\u30A6\u30F3\u30C8\u3068\u30D1\u30B9\u30EF\u30FC\u30C9","AdminAccountSetupDescription":"Comigo\u306B\u30ED\u30B0\u30A4\u30F3\u3059\u308B\u305F\u3081\u306E\u7BA1\u7406\u8005\u30A2\u30AB\u30A6\u30F3\u30C8\u3068\u30D1\u30B9\u30EF\u30FC\u30C9\u3092\u8A2D\u5B9A\u3057\u3066\u304F\u3060\u3055\u3044\u3002\u8A2D\u5B9A\u5F8C\u3001\u30B5\u30FC\u30D3\u30B9\u3078\u30A2\u30AF\u30BB\u30B9\u3059\u308B\u969B\u306B\u306F\u30ED\u30B0\u30A4\u30F3\u304C\u5FC5\u8981\u3068\u306A\u308A\u307E\u3059\u3002","ConfigStorageLocationPrompt":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u306E\u4FDD\u5B58\u5834\u6240\u3092\u9078\u629E\u3057\u3066\u304F\u3060\u3055\u3044:","set_account_password":"\u30A2\u30AB\u30A6\u30F3\u30C8\u306E\u30D1\u30B9\u30EF\u30FC\u30C9\u3092\u8A2D\u5B9A","connect_tailscale":"Tailscale\u30CD\u30C3\u30C8\u30EF\u30FC\u30AF\u306B\u63A5\u7D9A","disconnect_tailscale":"Tailscale\u3092\u5207\u65AD","tailscale_status":"Tailscale\u30B9\u30C6\u30FC\u30BF\u30B9","service_status":"\u30B5\u30FC\u30D3\u30B9\u72B6\u614B","running":"\u5B9F\u884C\u4E2D","client_count":"\u30AF\u30E9\u30A4\u30A2\u30F3\u30C8\u6570","host_system":"\u30DB\u30B9\u30C8\u30B7\u30B9\u30C6\u30E0","connection_status":"\u63A5\u7D9A\u72B6\u6CC1","connected":"\u63A5\u7D9A\u6E08\u307F","not_connected":"\u672A\u63A5\u7D9A","enable_funnel":"Funnel\u6709\u52B9","enable_funnel_public_access":"Funnel\u6709\u52B9\uFF08\u30D1\u30D6\u30EA\u30C3\u30AF\u30A2\u30AF\u30BB\u30B9\uFF09","disabled":"\u672A\u6709\u52B9\u5316","ip_address":"IP\u30A2\u30C9\u30EC\u30B9","service_version":"\u30B5\u30FC\u30D3\u30B9\u30D0\u30FC\u30B8\u30E7\u30F3","read_link":"\u30EA\u30F3\u30AF\u3092\u8AAD\u3080","submiting":"\u9001\u4FE1\u4E2D...","enable":"\u6709\u52B9","disable":"\u7121\u52B9","TailscaleAuthKey":"Tailscale\u81EA\u52D5\u8A8D\u8A3C\u30AD\u30FC","TailscaleAuthKeyDescription":"Tailscale\u306E\u81EA\u52D5\u8A8D\u8A3C\u30AD\u30FC\uFF08TS_AUTHKEY\uFF09\u3002\u30D6\u30E9\u30A6\u30B6\u306E\u306A\u3044\u74B0\u5883\u3067\u8A8D\u8A3C\u3092\u884C\u3046\u305F\u3081\u306B\u4F7F\u7528\u3057\u307E\u3059\u3002","verify_tailscale":"\u30EA\u30F3\u30AF\u3092\u30AF\u30EA\u30C3\u30AF\u3057\u3066\u3001Tailscale\u3092\u8A8D\u8A3C\u3057\u3066\u304F\u3060\u3055\u3044:","funnel_status":"Funnel\u72B6\u614B","funnel_tunnel":"Funnel\u96A7\u9053","funnel_setup_done":"Funnel\u306E\u8A2D\u5B9A\u304C\u5B8C\u4E86\u3057\u307E\u3057\u305F","funnel_setup_not_done":"Funnel\u306E\u8A2D\u5B9A\u304C\u5FC5\u8981\u3067\u3059","funnel_not_set_hint":"Funnel\u6A29\u9650\u3092\u4F7F\u7528\u3059\u308B\u306B\u306F\u3001\u6B21\u306E\u8A2D\u5B9A\u304C\u5FC5\u8981\u3067\u3059\uFF1A","funnel_require_dns_1":"Tailscale\u30B3\u30F3\u30BD\u30FC\u30EB\u306EDNS\u30D1\u30CD\u30EB\u3067","funnel_require_dns_2":"MagicDNS\u3068HTTPS\u6A5F\u80FD\u3092\u6709\u52B9\u5316\u3057\u307E\u3059\u3002","funnel_require_acl_1":"Tailscale\u30B3\u30F3\u30BD\u30FC\u30EB\u306EACL\u30D1\u30CD\u30EB\u3067","funnel_require_acl_2":"ACL\u30EB\u30FC\u30EB\u3092\u7DE8\u96C6\u3057\u3001Funnel\u96A7\u9053\u3092\u6709\u52B9\u5316\u3057\u307E\u3059\u3002","funnel_require_acl_3":"\uFF08\u30B5\u30F3\u30D7\u30EBJSON\u30D5\u30A1\u30A4\u30EB\u3092\u3053\u3061\u3089\u304B\u3089\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9\uFF09\u3002","funnel_require_password_1":"\u300C\u30D5\u30A1\u30CD\u30EB\u30ED\u30B0\u30A4\u30F3\u306E\u78BA\u8A8D\u300D\u3092\u6709\u52B9\u306B\u3059\u308B\u5834\u5408\u3001\u30D7\u30E9\u30A4\u30D9\u30FC\u30C8\u30A2\u30AF\u30BB\u30B9\u3092\u78BA\u4FDD\u3059\u308B\u305F\u3081\u306B\u8A8D\u8A3C\uFF08\u30ED\u30B0\u30A4\u30F3\u30D1\u30B9\u30EF\u30FC\u30C9\uFF09\u3092\u6709\u52B9\u306B\u3059\u308B\u5FC5\u8981\u304C\u3042\u308A\u307E\u3059\u3002","verify_link":"\u8A8D\u8A3C\u30EA\u30F3\u30AF","copy_link":"\u30EA\u30F3\u30AF\u3092\u30B3\u30D4\u30FC","copy_success":"\u30AF\u30EA\u30C3\u30D7\u30DC\u30FC\u30C9\u306B\u30B3\u30D4\u30FC\u3057\u307E\u3057\u305F","copy_failed":"\u30B3\u30D4\u30FC\u306B\u5931\u6557\u3057\u307E\u3057\u305F\u3002\u624B\u52D5\u3067\u30B3\u30D4\u30FC\u3057\u3066\u304F\u3060\u3055\u3044","FunnelLoginCheck":"\u30D5\u30A1\u30CD\u30EB\u30ED\u30B0\u30A4\u30F3\u306E\u78BA\u8A8D","FunnelLoginCheckDescription":"Funnel\u30C8\u30F3\u30CD\u30EB\u3092\u6709\u52B9\u306B\u3059\u308B\u969B\u3001\u73FE\u5728\u30ED\u30B0\u30A4\u30F3\u30D1\u30B9\u30EF\u30FC\u30C9\u304C\u8A2D\u5B9A\u3055\u308C\u3066\u3044\u308B\u304B\u78BA\u8A8D\u3059\u308B\u3002","funnel_login_check_enabled_but_no_password":"\u300C\u30D5\u30A1\u30CD\u30EB\u30ED\u30B0\u30A4\u30F3\u306E\u78BA\u8A8D\u300D\u306F\u6709\u52B9\u3067\u3059\u304C\u3001\u30ED\u30B0\u30A4\u30F3\u30D1\u30B9\u30EF\u30FC\u30C9\u304C\u8A2D\u5B9A\u3055\u308C\u3066\u3044\u306A\u3044\u305F\u3081\u3001\u30D5\u30A1\u30CD\u30EB\u30C8\u30F3\u30CD\u30EB\u3092\u6709\u52B9\u306B\u3067\u304D\u307E\u305B\u3093\u3002","tailscale_settings_submitted_check_status":"Tailscale\u306E\u8A2D\u5B9A\u3092\u9001\u4FE1\u3057\u307E\u3057\u305F\u3002Tailscale\u306E\u30B9\u30C6\u30FC\u30BF\u30B9\u3092\u78BA\u8A8D\u3057\u3066\u304F\u3060\u3055\u3044\u3002","value_already_exists_do_not_add_again":"\u3053\u306E\u5024\u306F\u3059\u3067\u306B\u5B58\u5728\u3057\u307E\u3059\u3002\u91CD\u8907\u3057\u3066\u8FFD\u52A0\u3057\u306A\u3044\u3067\u304F\u3060\u3055\u3044\u3002","file_uploaded_successfully":"\u30D5\u30A1\u30A4\u30EB\u306E\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u306B\u6210\u529F\u3057\u307E\u3057\u305F\u3002","content_empty_please_enter_before_submit":"\u5185\u5BB9\u304C\u7A7A\u3067\u3059\u3002\u5165\u529B\u3057\u3066\u304B\u3089\u9001\u4FE1\u3057\u3066\u304F\u3060\u3055\u3044\u3002","default_prompt_message":"\u30C7\u30D5\u30A9\u30EB\u30C8\u306E\u30E1\u30C3\u30BB\u30FC\u30B8","confirm":"\u78BA\u8A8D","ok":"OK","cancel":"\u30AD\u30E3\u30F3\u30BB\u30EB","uploading":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u4E2D...","upload_failed_network_error":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u306B\u5931\u6557\u3057\u307E\u3057\u305F\uFF1A\u30CD\u30C3\u30C8\u30EF\u30FC\u30AF\u30A8\u30E9\u30FC","drag_or_click_to_upload":"\u3053\u3053\u306B\u30D5\u30A1\u30A4\u30EB\u3092\u30C9\u30E9\u30C3\u30B0\uFF06\u30C9\u30ED\u30C3\u30D7\u3059\u308B\u304B\u3001\u30AF\u30EA\u30C3\u30AF\u3057\u3066\u30D5\u30A1\u30A4\u30EB\u3092\u9078\u629E\u3057\u3066\u304F\u3060\u3055\u3044","selected_file":"\u9078\u629E\u3055\u308C\u305F\u30D5\u30A1\u30A4\u30EB","select_upload_target_store":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u5148\u306E\u66F8\u5EAB\u3092\u9078\u629E","store_not_exists":"\u66F8\u5EAB\u30D1\u30B9\u304C\u5B58\u5728\u3057\u307E\u305B\u3093","no_available_stores":"\u5229\u7528\u53EF\u80FD\u306A\u66F8\u5EAB\u304C\u3042\u308A\u307E\u305B\u3093\u3002\u8A2D\u5B9A\u3067\u66F8\u5EAB\u30D1\u30B9\u3092\u8FFD\u52A0\u3057\u3066\u304F\u3060\u3055\u3044","store_validation_failed":"\u7121\u52B9\u306A\u66F8\u5EAB\u30D1\u30B9","passwords_not_match":"\u5165\u529B\u3057\u305F2\u3064\u306E\u30D1\u30B9\u30EF\u30FC\u30C9\u304C\u4E00\u81F4\u3057\u307E\u305B\u3093\u3002","please_delete_other_config_first":"\u4ED6\u306E\u5834\u6240\u306E\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u3092\u5148\u306B\u524A\u9664\u3057\u3066\u304F\u3060\u3055\u3044\u3002","save_config_success":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u306E\u4FDD\u5B58\u306B\u6210\u529F\u3057\u307E\u3057\u305F\uFF01","save_config_failed":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u306E\u4FDD\u5B58\u306B\u5931\u6557\u3057\u307E\u3057\u305F\u3002","no_config_file_to_delete_in_path":"\u9078\u629E\u3057\u305F\u30D1\u30B9\u306B\u306F\u524A\u9664\u3067\u304D\u308B\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u304C\u3042\u308A\u307E\u305B\u3093\u3002","delete_config_success":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u3092\u524A\u9664\u3057\u307E\u3057\u305F\u3002","delete_config_failed":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u306E\u524A\u9664\u306B\u5931\u6557\u3057\u307E\u3057\u305F\u3002","err_invalid_number":"\u6709\u52B9\u306A\u6570\u5B57\u3092\u5165\u529B\u3057\u3066\u304F\u3060\u3055\u3044","err_number_range":"\u5024\u306F {0} \u304B\u3089 {1} \u306E\u9593\u3067\u3042\u308B\u5FC5\u8981\u304C\u3042\u308A\u307E\u3059","err_update_config_failed":"\u8A2D\u5B9A\u306E\u66F4\u65B0\u306B\u5931\u6557\u3057\u307E\u3057\u305F","err_network_error":"\u30CD\u30C3\u30C8\u30EF\u30FC\u30AF\u30A8\u30E9\u30FC\u3001\u518D\u8A66\u884C\u3057\u3066\u304F\u3060\u3055\u3044","err_update_login_settings_failed":"\u30ED\u30B0\u30A4\u30F3\u8A2D\u5B9A\u306E\u66F4\u65B0\u306B\u5931\u6557\u3057\u307E\u3057\u305F","err_add_config_failed":"\u8A2D\u5B9A\u306E\u8FFD\u52A0\u306B\u5931\u6557\u3057\u307E\u3057\u305F","err_delete_config_failed":"\u8A2D\u5B9A\u306E\u524A\u9664\u306B\u5931\u6557\u3057\u307E\u3057\u305F","err_save_config_failed":"\u8A2D\u5B9A\u306E\u4FDD\u5B58\u306B\u5931\u6557\u3057\u307E\u3057\u305F","err_config_locked":"\u8A2D\u5B9A\u304C\u30ED\u30C3\u30AF\u3055\u308C\u3066\u3044\u308B\u305F\u3081\u3001\u5909\u66F4\u3067\u304D\u307E\u305B\u3093","read":"\u65E2\u8AAD","unread":"\u672A\u8AAD","mark_all_as_read":"\u3059\u3079\u3066\u3092\u65E2\u8AAD\u306B\u3059\u308B","new_version_available":"\u65B0\u3057\u3044\u30D0\u30FC\u30B8\u30E7\u30F3\u304C\u5229\u7528\u53EF\u80FD\u3067\u3059\uFF01","download_now":"\u4ECA\u3059\u3050\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","infinite_scroll_button":"\u7121\u9650\u30B9\u30AF\u30ED\u30FC\u30EB","paged_scroll_button":"\u30B9\u30AF\u30ED\u30FC\u30EB","page_flip_button":"\u5DE6\u53F3\u3081\u304F\u308A","infinite_scroll":"\u7121\u9650\u30B9\u30AF\u30ED\u30FC\u30EB\uFF08\u3059\u3079\u3066\uFF09","paged_scroll":"\u30B9\u30AF\u30ED\u30FC\u30EB\uFF0832\u30DA\u30FC\u30B8\u307E\u3067\uFF09","page_flip":"\u30DA\u30FC\u30B8\u3081\u304F\u308A","reading_progress_percent":"\u8AAD\u66F8\u9032\u6357\uFF08\uFF05\uFF09","reading_progress_page":"\u8AAD\u66F8\u9032\u6357\uFF08\u6570\u5B57\uFF09","reading_history":"\u8AAD\u66F8\u5C65\u6B74","view_all_reading_history":"\u3059\u3079\u3066\u306E\u8AAD\u66F8\u5C65\u6B74\u3092\u8868\u793A","loading":"\u8AAD\u307F\u8FBC\u307F\u4E2D...","no_reading_history":"\u8AAD\u66F8\u5C65\u6B74\u304C\u3042\u308A\u307E\u305B\u3093","unknown":"\u4E0D\u660E","systray_tooltip":"Comigo \u30B3\u30DF\u30C3\u30AF\u30EA\u30FC\u30C0\u30FC","systray_open_browser":"\u30D6\u30E9\u30A6\u30B6\u3092\u958B\u304F","systray_open_browser_tooltip":"\u30D6\u30E9\u30A6\u30B6\u3067Comigo\u3092\u958B\u304F","systray_copy_url":"\u8AAD\u66F8URL\u3092\u30B3\u30D4\u30FC","systray_copy_url_tooltip":"\u8AAD\u66F8URL\u3092\u30AF\u30EA\u30C3\u30D7\u30DC\u30FC\u30C9\u306B\u30B3\u30D4\u30FC","systray_enable_tailscale":"Tailscale\u3092\u6709\u52B9\u5316","systray_disable_tailscale":"Tailscale\u3092\u7121\u52B9\u5316","systray_toggle_tailscale_tooltip":"Tailscale\u72B6\u614B\u3092\u5207\u308A\u66FF\u3048","systray_language":"\u8A00\u8A9E\u5207\u308A\u66FF\u3048","systray_language_tooltip":"\u30A4\u30F3\u30BF\u30FC\u30D5\u30A7\u30FC\u30B9\u8A00\u8A9E\u3092\u5207\u308A\u66FF\u3048","systray_language_zh":"\u4E2D\u6587","systray_language_zh_tooltip":"\u4E2D\u56FD\u8A9E\u306B\u5207\u308A\u66FF\u3048","systray_language_en":"English","systray_language_en_tooltip":"\u82F1\u8A9E\u306B\u5207\u308A\u66FF\u3048","systray_language_ja":"\u65E5\u672C\u8A9E","systray_language_ja_tooltip":"\u65E5\u672C\u8A9E\u306B\u5207\u308A\u66FF\u3048","systray_open_directory":"\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3092\u958B\u304F","systray_open_directory_tooltip":"\u95A2\u9023\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3092\u958B\u304F","systray_config_directory":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA","systray_config_directory_tooltip":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3092\u958B\u304F","systray_store":"\u66F8\u5EAB","systray_store_folder":"\u66F8\u5EAB\u30D5\u30A9\u30EB\u30C0\u30FC","systray_refresh_directories":"\u66F4\u65B0","systray_refresh_directories_tooltip":"\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u30EA\u30B9\u30C8\u3092\u66F4\u65B0","systray_quit":"\u7D42\u4E86","systray_quit_tooltip":"Comigo\u3092\u7D42\u4E86","context_menu_open_with_comigo":"Comigo\u3067\u958B\u304F","register_context_menu":"\u53F3\u30AF\u30EA\u30C3\u30AF\u30E1\u30CB\u30E5\u30FC\u3092\u767B\u9332\uFF08Comigo\u3067\u958B\u304F\uFF09","unregister_context_menu":"\u53F3\u30AF\u30EA\u30C3\u30AF\u30E1\u30CB\u30E5\u30FC\u3092\u524A\u9664\uFF08Comigo\u3067\u958B\u304F\uFF09","systray_extra":"\u305D\u306E\u4ED6","systray_extra_tooltip":"\u305D\u306E\u4ED6\u306E\u9023\u643A\u6A5F\u80FD","systray_project":"Comigo \u30D7\u30ED\u30B8\u30A7\u30AF\u30C8\u30DA\u30FC\u30B8","systray_project_tooltip":"Comigo \u306E GitHub \u30EA\u30DD\u30B8\u30C8\u30EA\u3092\u958B\u304F","register_folder_context_menu":"\u30D5\u30A9\u30EB\u30C0\u30FC\u53F3\u30AF\u30EA\u30C3\u30AF\u30E1\u30CB\u30E5\u30FC\u3092\u767B\u9332\uFF08Comigo\u3067\u958B\u304F\uFF09","unregister_folder_context_menu":"\u30D5\u30A9\u30EB\u30C0\u30FC\u53F3\u30AF\u30EA\u30C3\u30AF\u30E1\u30CB\u30E5\u30FC\u3092\u524A\u9664\uFF08Comigo\u3067\u958B\u304F\uFF09","register_file_association":"\u5727\u7E2E\u30D5\u30A1\u30A4\u30EB\u306E\u95A2\u9023\u4ED8\u3051\u3092\u767B\u9332\uFF08\u5019\u88DC\u3068\u3057\u3066\u8FFD\u52A0\uFF09","unregister_file_association":"\u5727\u7E2E\u30D5\u30A1\u30A4\u30EB\u306E\u95A2\u9023\u4ED8\u3051\u3092\u89E3\u9664","create_desktop_shortcut":"\u30C7\u30B9\u30AF\u30C8\u30C3\u30D7\u306B\u30B7\u30E7\u30FC\u30C8\u30AB\u30C3\u30C8\u3092\u4F5C\u6210","log_language_changed_to_chinese":"\u8A00\u8A9E\u304C\u4E2D\u56FD\u8A9E\u306B\u5909\u66F4\u3055\u308C\u307E\u3057\u305F","log_language_changed_to_english":"\u8A00\u8A9E\u304C\u82F1\u8A9E\u306B\u5909\u66F4\u3055\u308C\u307E\u3057\u305F","log_language_changed_to_japanese":"\u8A00\u8A9E\u304C\u65E5\u672C\u8A9E\u306B\u5909\u66F4\u3055\u308C\u307E\u3057\u305F","log_requesting_quit_from_systray":"\u30B7\u30B9\u30C6\u30E0\u30C8\u30EC\u30A4\u304B\u3089\u7D42\u4E86\u3092\u8981\u6C42\u3057\u3066\u3044\u307E\u3059","log_failed_to_read_icon_file":"\u30A2\u30A4\u30B3\u30F3\u30D5\u30A1\u30A4\u30EB\u306E\u8AAD\u307F\u8FBC\u307F\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v\u3001\u30C7\u30D5\u30A9\u30EB\u30C8\u30A2\u30A4\u30B3\u30F3\u3092\u4F7F\u7528\u3057\u307E\u3059","log_opening_browser":"\u30D6\u30E9\u30A6\u30B6\u3092\u958B\u3044\u3066\u3044\u307E\u3059: %s","log_failed_to_copy_url":"URL\u3092\u30AF\u30EA\u30C3\u30D7\u30DC\u30FC\u30C9\u306B\u30B3\u30D4\u30FC\u3067\u304D\u307E\u305B\u3093\u3067\u3057\u305F: %v","log_copied_url_to_clipboard":"URL\u3092\u30AF\u30EA\u30C3\u30D7\u30DC\u30FC\u30C9\u306B\u30B3\u30D4\u30FC\u3057\u307E\u3057\u305F: %s","log_failed_to_toggle_tailscale":"Tailscale\u306E\u72B6\u614B\u3092\u5207\u308A\u66FF\u3048\u3089\u308C\u307E\u305B\u3093\u3067\u3057\u305F: %v","log_failed_to_set_language":"\u8A00\u8A9E\u306E\u8A2D\u5B9A\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_get_config_dir":"\u8A2D\u5B9A\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_clear_folder_context_menu":"Windows\u30D5\u30A9\u30EB\u30C0\u30FC\u53F3\u30AF\u30EA\u30C3\u30AF\u30E1\u30CB\u30E5\u30FC\u306E\u524A\u9664\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_register_folder_context_menu":"Windows\u30D5\u30A9\u30EB\u30C0\u30FC\u53F3\u30AF\u30EA\u30C3\u30AF\u30E1\u30CB\u30E5\u30FC\u306E\u767B\u9332\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_create_desktop_shortcut":"\u30C7\u30B9\u30AF\u30C8\u30C3\u30D7\u30B7\u30E7\u30FC\u30C8\u30AB\u30C3\u30C8\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_unregister_archive_handler":"\u30A2\u30FC\u30AB\u30A4\u30D6\u30CF\u30F3\u30C9\u30E9\u30FC\u306E\u767B\u9332\u89E3\u9664\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_register_archive_handler":"\u30A2\u30FC\u30AB\u30A4\u30D6\u30CF\u30F3\u30C9\u30E9\u30FC\u306E\u767B\u9332\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_opening_comigo_project_page":"Comigo\u30D7\u30ED\u30B8\u30A7\u30AF\u30C8\u30DA\u30FC\u30B8\u3092\u958B\u3044\u3066\u3044\u307E\u3059: https://github.com/yumenaka/comigo","log_failed_to_open_directory":"\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3092\u958B\u3051\u307E\u305B\u3093\u3067\u3057\u305F: %v","log_upload_file_count":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u30D5\u30A1\u30A4\u30EB\u6570: %d","log_upload_path_not_set":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u30D1\u30B9\u304C\u8A2D\u5B9A\u3055\u308C\u3066\u3044\u307E\u305B\u3093","log_mkdir_failed":"\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_mkdir_upload_folder_success":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u30D5\u30A9\u30EB\u30C0\u30FC\u306E\u4F5C\u6210\u306B\u6210\u529F\u3057\u307E\u3057\u305F: %s","log_file_upload_success":"\u30D5\u30A1\u30A4\u30EB\u306E\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u306B\u6210\u529F\u3057\u307E\u3057\u305F: %s","log_no_changes_skipped_scan":"\u8A2D\u5B9A\u306B\u5909\u66F4\u304C\u306A\u3044\u305F\u3081\u3001\u66F8\u5EAB\u30D1\u30B9\u306E\u30B9\u30AD\u30E3\u30F3\u3092\u30B9\u30AD\u30C3\u30D7\u3057\u307E\u3057\u305F","log_no_changes_skipped_rescan":"\u8A2D\u5B9A\u306B\u5909\u66F4\u304C\u306A\u3044\u305F\u3081\u3001\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u518D\u30B9\u30AD\u30E3\u30F3\u3092\u30B9\u30AD\u30C3\u30D7\u3057\u307E\u3057\u305F","log_tailscale_config_changed_restart":"Tailscale\u8A2D\u5B9A\u304C\u5909\u66F4\u3055\u308C\u305F\u305F\u3081\u3001Tailscale\u30B5\u30FC\u30D0\u30FC\u3092\u518D\u8D77\u52D5\u3057\u307E\u3059","log_tailscale_enabled_start":"Tailscale\u304C\u6709\u52B9\u306B\u306A\u3063\u305F\u305F\u3081\u3001Tailscale\u30B5\u30FC\u30D0\u30FC\u3092\u8D77\u52D5\u3057\u307E\u3059","log_tailscale_disabled_stop":"Tailscale\u304C\u7121\u52B9\u306B\u306A\u3063\u305F\u305F\u3081\u3001Tailscale\u30B5\u30FC\u30D0\u30FC\u3092\u505C\u6B62\u3057\u307E\u3059","log_scan_subdirectory_error":"\u30B5\u30D6\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u30B9\u30AD\u30E3\u30F3\u30A8\u30E9\u30FC: %v","log_get_file_info_failed":"\u30D5\u30A1\u30A4\u30EB\u60C5\u5831\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_tailscale_disabled_skip_qrcode":"Tailscale\u304C\u7121\u52B9\u306E\u305F\u3081\u3001QR\u30B3\u30FC\u30C9\u8868\u793A\u6A5F\u80FD\u3092\u30B9\u30AD\u30C3\u30D7\u3057\u307E\u3059","log_tailscale_status_check_exceeded":"Tailscale\u72B6\u614B\u30C1\u30A7\u30C3\u30AF\u306E\u4E0A\u9650\u306B\u9054\u3057\u305F\u305F\u3081\u3001\u3053\u308C\u4EE5\u4E0A\u306E\u30C1\u30A7\u30C3\u30AF\u3092\u505C\u6B62\u3057\u307E\u3059","log_tailscale_not_yet_fqdn":"Tailscale FQDN\u304C\u307E\u3060\u5229\u7528\u3067\u304D\u307E\u305B\u3093","log_flip_mode_book_id":"\u30D5\u30EA\u30C3\u30D7\u30E2\u30FC\u30C9\u66F8\u7C4DID: %s","log_invalid_upload_path":"\u7121\u52B9\u306A\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u30D1\u30B9","log_time_elapsed":"\u7D4C\u904E\u6642\u9593: %v","log_waiting_for_api_health":"API\u30D8\u30EB\u30B9\u30C1\u30A7\u30C3\u30AF\u30A8\u30F3\u30C9\u30DD\u30A4\u30F3\u30C8\u3092\u5F85\u6A5F\u4E2D...","log_api_healthy_ready":"Comigo API\u304C\u6B63\u5E38\u306B\u52D5\u4F5C\u3057\u3001\u6E96\u5099\u304C\u6574\u3044\u307E\u3057\u305F","log_path_error":"\u30D1\u30B9\u30A8\u30E9\u30FC","log_checking_cfg_sharename":"\u8A2D\u5B9AShareName\u3092\u78BA\u8A8D\u4E2D","log_warning_failed_to_get_homedir":"\u8B66\u544A: \u30DB\u30FC\u30E0\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_warning_failed_to_get_executable_path":"\u8B66\u544A: \u5B9F\u884C\u30D5\u30A1\u30A4\u30EB\u30D1\u30B9\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_child_daemon_started":"\u5B50\u30C7\u30FC\u30E2\u30F3\u304C\u8D77\u52D5\u3057\u307E\u3057\u305F\u304B\uFF1F","log_daemon_started":"\u30C7\u30FC\u30E2\u30F3\u304C\u8D77\u52D5\u3057\u307E\u3057\u305F","log_update_config":"\u8A2D\u5B9A\u3092\u66F4\u65B0: %s = %s","log_failed_to_update_local_config":"\u30ED\u30FC\u30AB\u30EB\u8A2D\u5B9A\u306E\u66F4\u65B0\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_update_user_info_username":"\u30E6\u30FC\u30B6\u30FC\u60C5\u5831\u3092\u66F4\u65B0: \u30E6\u30FC\u30B6\u30FC\u540D=%s","log_update_user_info_current_password":"\u30E6\u30FC\u30B6\u30FC\u60C5\u5831\u3092\u66F4\u65B0: \u73FE\u5728\u306E\u30D1\u30B9\u30EF\u30FC\u30C9=%s","log_update_user_info_password":"\u30E6\u30FC\u30B6\u30FC\u60C5\u5831\u3092\u66F4\u65B0: \u30D1\u30B9\u30EF\u30FC\u30C9=%s","log_update_user_info_reenter_password":"\u30E6\u30FC\u30B6\u30FC\u60C5\u5831\u3092\u66F4\u65B0: \u30D1\u30B9\u30EF\u30FC\u30C9\u518D\u5165\u529B=%s","log_add_array_config_handler":"\u914D\u5217\u8A2D\u5B9A\u30CF\u30F3\u30C9\u30E9\u30FC\u3092\u8FFD\u52A0: %s = %s","log_delete_array_config_handler":"\u914D\u5217\u8A2D\u5B9A\u30CF\u30F3\u30C9\u30E9\u30FC\u3092\u524A\u9664: %s = %s","log_received_json_data":"JSON\u30C7\u30FC\u30BF\u3092\u53D7\u4FE1: %s","log_failed_to_parse_json":"JSON\u30C7\u30FC\u30BF\u306E\u89E3\u6790\u306B\u5931\u6557\u3057\u307E\u3057\u305F","log_failed_to_scan_store_path":"\u66F8\u5EAB\u30D1\u30B9\u306E\u30B9\u30AD\u30E3\u30F3\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_save_results_to_database":"\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u3078\u306E\u7D50\u679C\u4FDD\u5B58\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_error_failed_save_to_directory":"\u30A8\u30E9\u30FC: %s \u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3078\u306E\u4FDD\u5B58\u306B\u5931\u6557\u3057\u307E\u3057\u305F","log_error_find_config_in":"\u30A8\u30E9\u30FC: %s %s \u3067\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u304C\u898B\u3064\u304B\u308A\u307E\u3057\u305F","log_error_failed_to_delete_config":"\u30A8\u30E9\u30FC: %s \u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u5185\u306E\u8A2D\u5B9A\u306E\u524A\u9664\u306B\u5931\u6557\u3057\u307E\u3057\u305F","log_store_url_already_exists":"\u66F8\u5EABURL\u304C\u65E2\u306B\u5B58\u5728\u3057\u307E\u3059: %s","log_store_url_already_exists_in_config":"\u8A2D\u5B9A\u5185\u306B\u66F8\u5EABURL\u304C\u65E2\u306B\u5B58\u5728\u3057\u307E\u3059: %s","log_failed_to_add_store_url":"\u8A2D\u5B9A\u304B\u3089\u66F8\u5EABURL\u306E\u8FFD\u52A0\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_string_already_exists":"\u6587\u5B57\u5217 \'%s\' \u304C\u65E2\u306B\u5B58\u5728\u3057\u307E\u3059","log_failed_to_unmarshal_json":"JSON\u306E\u9006\u30B7\u30EA\u30A2\u30EB\u5316\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_unknown_config_key":"\u4E0D\u660E\u306A\u8A2D\u5B9A\u30AD\u30FC: %s","log_error_getting_executable_path":"\u5B9F\u884C\u30D5\u30A1\u30A4\u30EB\u30D1\u30B9\u306E\u53D6\u5F97\u30A8\u30E9\u30FC: %s","log_executable_name":"\u5B9F\u884C\u30D5\u30A1\u30A4\u30EB\u540D: %s","log_executable_path":"\u5B9F\u884C\u30D5\u30A1\u30A4\u30EB\u30D1\u30B9: %s","log_failed_to_unregister_windows_context_menu":"Windows\u53F3\u30AF\u30EA\u30C3\u30AF\u30E1\u30CB\u30E5\u30FC\u306E\u767B\u9332\u89E3\u9664\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_register_windows_context_menu":"Windows\u53F3\u30AF\u30EA\u30C3\u30AF\u30E1\u30CB\u30E5\u30FC\u306E\u767B\u9332\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_received_new_args_from_instance":"\u5225\u306E\u30A4\u30F3\u30B9\u30BF\u30F3\u30B9\u304B\u3089\u65B0\u3057\u3044\u5F15\u6570\u3092\u53D7\u4FE1: %v","log_single_instance_check_failed":"\u5358\u4E00\u30A4\u30F3\u30B9\u30BF\u30F3\u30B9\u30C1\u30A7\u30C3\u30AF\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_args_sent_to_existing_instance":"\u5F15\u6570\u304C\u65E2\u5B58\u306E\u30A4\u30F3\u30B9\u30BF\u30F3\u30B9\u306B\u9001\u4FE1\u3055\u308C\u307E\u3057\u305F\u3002\u7D42\u4E86\u3057\u3066\u3044\u307E\u3059...","log_error_listing_books":"\u66F8\u7C4D\u306E\u4E00\u89A7\u53D6\u5F97\u30A8\u30E9\u30FC: %s","log_error_clearing_temp_files":"\u4E00\u6642\u30D5\u30A1\u30A4\u30EB\u306E\u30AF\u30EA\u30A2\u30A8\u30E9\u30FC: %s","log_cleared_temp_files":"\u4E00\u6642\u30D5\u30A1\u30A4\u30EB\u3092\u30AF\u30EA\u30A2\u3057\u307E\u3057\u305F: %s","log_failed_to_get_homedirectory":"\u30DB\u30FC\u30E0\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_program_directory":"\u30D7\u30ED\u30B0\u30E9\u30E0\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA: %s","log_failed_to_get_program_directory":"\u30D7\u30ED\u30B0\u30E9\u30E0\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_get_working_directory":"\u4F5C\u696D\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_failed_to_open_file":"\u30D5\u30A1\u30A4\u30EB\u3092\u958B\u3051\u307E\u305B\u3093\u3067\u3057\u305F: %s, \u30A8\u30E9\u30FC: %v","log_failed_to_get_file_info":"\u30D5\u30A1\u30A4\u30EB\u60C5\u5831\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s, \u30A8\u30E9\u30FC: %v","log_add_book_error":"AddBook_error \u66F8\u7C4DID: %s %s","log_get_child_books_for_bookid":"\u66F8\u7C4DID %s \u306E\u5B50\u66F8\u7C4D\u3092\u53D6\u5F97","log_get_child_books_count":"\u66F8\u7C4DID %v \u306E %v \u518A\u306E\u5B50\u66F8\u7C4D\u3092\u53D6\u5F97","log_get_bookshelf_error":"\u672C\u68DA\u306E\u53D6\u5F97\u30A8\u30E9\u30FC: %v","log_error_getting_absolute_path":"\u7D76\u5BFE\u30D1\u30B9\u306E\u53D6\u5F97\u30A8\u30E9\u30FC: %v","log_cannot_shorten_id":"ID\u3092\u77ED\u7E2E\u3067\u304D\u307E\u305B\u3093: %s","log_error_getting_book_group":"\u66F8\u7C4D\u30B0\u30EB\u30FC\u30D7\u306E\u53D6\u5F97\u30A8\u30E9\u30FC: %s","log_warning_failed_to_set_socket_permissions":"\u8B66\u544A: \u30BD\u30B1\u30C3\u30C8\u6A29\u9650\u306E\u8A2D\u5B9A\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_single_instance_server_started":"\u5358\u4E00\u30A4\u30F3\u30B9\u30BF\u30F3\u30B9\u30B5\u30FC\u30D0\u30FC\u304C\u8D77\u52D5\u3057\u307E\u3057\u305F: %s","log_failed_to_accept_connection":"\u63A5\u7D9A\u306E\u53D7\u3051\u5165\u308C\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_decode_message":"\u30E1\u30C3\u30BB\u30FC\u30B8\u306E\u30C7\u30B3\u30FC\u30C9\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_handle_new_args":"\u65B0\u3057\u3044\u5F15\u6570\u306E\u51E6\u7406\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_received_and_processed_new_args":"\u65B0\u3057\u3044\u5F15\u6570\u3092\u53D7\u4FE1\u3057\u3001\u51E6\u7406\u3057\u307E\u3057\u305F: %v","log_failed_to_read_response":"\u5FDC\u7B54\u306E\u8AAD\u307F\u53D6\u308A\u306B\u5931\u6557\u3057\u307E\u3057\u305F\u304C\u3001\u30E1\u30C3\u30BB\u30FC\u30B8\u306F\u9001\u4FE1\u3055\u308C\u305F\u53EF\u80FD\u6027\u304C\u3042\u308A\u307E\u3059: %v","log_successfully_sent_args":"\u65E2\u5B58\u306E\u30A4\u30F3\u30B9\u30BF\u30F3\u30B9\u306B\u5F15\u6570\u3092\u6B63\u5E38\u306B\u9001\u4FE1\u3057\u307E\u3057\u305F: %v","log_error_closing_listener":"\u30EA\u30B9\u30CA\u30FC\u306E\u30AF\u30ED\u30FC\u30BA\u30A8\u30E9\u30FC: %v","log_another_instance_running":"\u5225\u306E\u30A4\u30F3\u30B9\u30BF\u30F3\u30B9\u304C\u65E2\u306B\u5B9F\u884C\u4E2D\u3067\u3059\u3002\u5F15\u6570\u3092\u9001\u4FE1\u3057\u3066\u3044\u307E\u3059...","log_get_book_error":"\u66F8\u7C4D\u306E\u53D6\u5F97: %v","log_saving_books_meta_data_to":"\u66F8\u7C4DMetadata\u3092 %s \u306B\u4FDD\u5B58\u4E2D","log_error_saving_book":"\u66F8\u7C4D %s \u306E\u4FDD\u5B58\u30A8\u30E9\u30FC: %s","log_successfully_saved_books":"%d \u518A\u306E\u66F8\u7C4D\u3092 %s \u306B\u6B63\u5E38\u306B\u4FDD\u5B58\u3057\u307E\u3057\u305F","log_cfg_save_to":"\u8A2D\u5B9A\u3092 %s \u306B\u4FDD\u5B58","log_try_delete_cfg_in":"%s \u5185\u306E\u8A2D\u5B9A\u3092\u524A\u9664\u3057\u3088\u3046\u3068\u3057\u3066\u3044\u307E\u3059","log_websocket_server_received":"websocket\u30B5\u30FC\u30D0\u30FC\u304C\u53D7\u4FE1: %v","log_failed_to_create_temp_config_dir":"\u4E00\u6642\u8A2D\u5B9A\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_failed_to_create_config_dir":"\u8A2D\u5B9A\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_username_or_password_empty":"\u30E6\u30FC\u30B6\u30FC\u540D\u307E\u305F\u306F\u30D1\u30B9\u30EF\u30FC\u30C9\u304C\u7A7A\u3067\u3059\u3002\u30C7\u30D5\u30A9\u30EB\u30C8\u306EJWT\u7F72\u540D\u30AD\u30FC\u3092\u4F7F\u7528\u3057\u307E\u3059\u3002","log_invalid_port_number":"\u7121\u52B9\u306A\u30DD\u30FC\u30C8\u756A\u53F7\u3067\u3059\u3002\u30C7\u30D5\u30A9\u30EB\u30C8\u30DD\u30FC\u30C8\u3092\u4F7F\u7528\u3057\u307E\u3059: %d","log_tailscale_status_not_available":"Tailscale\u72B6\u614B\u306F\u307E\u3060\u5229\u7528\u3067\u304D\u307E\u305B\u3093: %v","log_server_action":"\u30B5\u30FC\u30D0\u30FC\u64CD\u4F5C: %v","log_server_action_string":"\u30B5\u30FC\u30D0\u30FC\u64CD\u4F5C: %s","log_loading_books_from":"\u66F8\u7C4D\u3092 %s \u304B\u3089\u8AAD\u307F\u8FBC\u307F\u4E2D","log_configured_store_urls":"\u8A2D\u5B9A\u3055\u308C\u305F\u66F8\u5EABURL: %v","log_book_data_directory_not_exist":"\u66F8\u7C4D\u30C7\u30FC\u30BF\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u304C\u307E\u3060\u5B58\u5728\u3057\u307E\u305B\u3093: %s","log_error_accessing_book_data_directory":"\u66F8\u7C4D\u30C7\u30FC\u30BF\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3078\u306E\u30A2\u30AF\u30BB\u30B9\u30A8\u30E9\u30FC: %s","log_error_reading_book_data_directory":"\u66F8\u7C4D\u30C7\u30FC\u30BF\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u8AAD\u307F\u8FBC\u307F\u30A8\u30E9\u30FC: %s","log_skipping_directory":"\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA %s \u3092\u30B9\u30AD\u30C3\u30D7","log_skipping_non_json_file":"\u975EJSON\u30D5\u30A1\u30A4\u30EB %s \u3092\u30B9\u30AD\u30C3\u30D7","log_error_reading_file":"\u30D5\u30A1\u30A4\u30EB %s \u306E\u8AAD\u307F\u8FBC\u307F\u30A8\u30E9\u30FC: %s","log_warning_corrupted_json_file":"\u8B66\u544A: \u7834\u640D\u3057\u305FJSON\u30D5\u30A1\u30A4\u30EB %s\u3001\u30B9\u30AD\u30C3\u30D7: %s","log_error_deleting_corrupted_file":"\u7834\u640D\u3057\u305F\u30D5\u30A1\u30A4\u30EB %s \u306E\u524A\u9664\u30A8\u30E9\u30FC: %s","log_error_adding_book_to_store":"\u66F8\u7C4D %s \u3092\u66F8\u5EAB\u306B\u8FFD\u52A0\u3059\u308B\u30A8\u30E9\u30FC: %s","log_loaded_books_so_far":"\u3053\u308C\u307E\u3067\u306B %d \u518A\u306E\u66F8\u7C4D\u3092\u8AAD\u307F\u8FBC\u307F\u307E\u3057\u305F\uFF08%s\uFF09","log_successfully_loaded_books":"%d \u518A\u306E\u66F8\u7C4D\u3092\u6B63\u5E38\u306B\u8AAD\u307F\u8FBC\u307F\u307E\u3057\u305F(%s)","log_error_adding_subfolder":"\u30B5\u30D6\u30D5\u30A9\u30EB\u30C0\u30FC\u306E\u8FFD\u52A0\u30A8\u30E9\u30FC: %s","log_error_saving_book_to_json":"\u66F8\u7C4D %s \u3092JSON\u306B\u4FDD\u5B58\u3059\u308B\u30A8\u30E9\u30FC: %s","log_error_adding_book":"\u66F8\u7C4D %s \u306E\u8FFD\u52A0\u30A8\u30E9\u30FC: %s","log_error_deleting_book_json_file":"\u66F8\u7C4D %s \u306EJSON\u30D5\u30A1\u30A4\u30EB\u524A\u9664\u30A8\u30E9\u30FC: %s","err_add_bookstore_key_exists":"\u66F8\u5EAB\u8FFD\u52A0\u30A8\u30E9\u30FC\uFF1A\u30AD\u30FC\u304C\u65E2\u306B\u5B58\u5728\u3057\u307E\u3059 [%s]","err_add_book_empty_bookid":"\u66F8\u7C4D\u8FFD\u52A0\u30A8\u30E9\u30FC\uFF1ABookID\u304C\u7A7A\u3067\u3059","err_add_bookstore_key_not_found":"\u66F8\u5EAB\u8FFD\u52A0\u30A8\u30E9\u30FC\uFF1A\u30AD\u30FC\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093 [%s]","err_cannot_find_group":"\u30B0\u30EB\u30FC\u30D7\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3001id=%s","err_getbook_cannot_find":"GetBook\uFF1A\u66F8\u7C4D\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3001id=%s","err_storebookmark_cannot_find":"StoreBookMark\uFF1A\u66F8\u7C4D\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3001id=%s","err_storebookmark_unknown_type":"StoreBookMark\uFF1A\u4E0D\u660E\u306A\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF\u30BF\u30A4\u30D7","err_getbookmark_cannot_find":"GetBookMark\uFF1A\u66F8\u7C4D\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3001id=%s","err_deletebook_cannot_find":"DeleteBook\uFF1A\u66F8\u7C4D\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3001id=%s","err_cannot_find_book_topofshelf":"\u30A8\u30E9\u30FC\uFF1ATopOfShelfInfo\u3067\u66F8\u7C4D\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093","err_cannot_find_child_books":"\u5B50\u66F8\u7C4D\u30C7\u30FC\u30BF\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3001BookID\uFF1A%s","err_getparentbook_cannot_find":"GetParentBook: childID\u3067\u66F8\u7C4D\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093=%s","err_cannot_find_book_parentfolder":"\u66F8\u7C4D\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3001parentFolder=%s","log_book_data_already_exists":"\u66F8\u7C4D\u30C7\u30FC\u30BF\u304C\u65E2\u306B\u5B58\u5728\u3057\u307E\u3059: %s  %s","err_slice_not_supported":"\u3053\u306E\u30B9\u30E9\u30A4\u30B9\u8A2D\u5B9A\u306F\u30B5\u30DD\u30FC\u30C8\u3055\u308C\u3066\u3044\u307E\u305B\u3093\uFF08[]string\u306E\u307F\u30B5\u30DD\u30FC\u30C8\uFF09","err_must_be_nonempty_config_pointer":"\u7A7A\u3067\u306A\u3044 *Config \u30DD\u30A4\u30F3\u30BF\u3067\u3042\u308B\u5FC5\u8981\u304C\u3042\u308A\u307E\u3059","err_non_htmx_request":"non-htmx request","err_no_form_data":"no form data","err_failed_to_find_executable_path":"\u30A8\u30E9\u30FC: \u5B9F\u884C\u30D5\u30A1\u30A4\u30EB\u30D1\u30B9\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093","err_failed_to_find_home_directory":"\u30A8\u30E9\u30FC: \u30DB\u30FC\u30E0\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093","err_funnel_mode_ports_only":"funnel\u30E2\u30FC\u30C9\u306F443\u30018443\u300110000\u30DD\u30FC\u30C8\u306E\u307F\u30B5\u30DD\u30FC\u30C8\u3057\u307E\u3059","err_number_not_found":"\u6570\u5B57\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093","err_file_not_zip_archive":"\u30D5\u30A1\u30A4\u30EB\u306FZIP\u30A2\u30FC\u30AB\u30A4\u30D6\u3067\u306F\u3042\u308A\u307E\u305B\u3093","err_file_not_rar_archive":"\u30D5\u30A1\u30A4\u30EB\u306FRAR\u30A2\u30FC\u30AB\u30A4\u30D6\u3067\u306F\u3042\u308A\u307E\u305B\u3093","err_scan_file_error":"\u30D5\u30A1\u30A4\u30EB\u30B9\u30AD\u30E3\u30F3\u30A8\u30E9\u30FC","err_countpages_pdf_invalid":"CountPagesOfPDF: \u7121\u52B9\u306APDF: %s %s","err_name_in_archive_empty":"nameInArchive\u304C\u7A7A\u3067\u3059","err_unsupported_archive_format":"\u30B5\u30DD\u30FC\u30C8\u3055\u308C\u3066\u3044\u306A\u3044\u30A2\u30FC\u30AB\u30A4\u30D6\u5F62\u5F0F\u3001\u307E\u305F\u306F\u30A2\u30FC\u30AB\u30A4\u30D6\u5185\u3067\u30D5\u30A1\u30A4\u30EB\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093","err_file_not_found_in_archive":"\u30A2\u30FC\u30AB\u30A4\u30D6\u5185\u3067\u30D5\u30A1\u30A4\u30EB\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093","err_extract_path_not_found":"\u30B3\u30F3\u30C6\u30AD\u30B9\u30C8\u5185\u3067extractPath\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093","err_needfile_empty":"needFile\u304C\u7A7A\u3067\u3059","err_getdata_from_epub_error":"getDataFromEpub\u30A8\u30E9\u30FC\u3002epubPath:%s  needFile:%s","err_container_xml_empty":"container.xml\u306E\u5185\u5BB9\u304C\u7A7A\u3067\u3059","err_no_valid_opf_path":"container.xml\u3067\u6709\u52B9\u306AOPF\u30D1\u30B9\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093","err_content_type_not_found":"\u30AD\u30E3\u30C3\u30B7\u30E5\u5185\u3067contentType\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093","err_file_does_not_exist":"\u30D5\u30A1\u30A4\u30EB\u304C\u5B58\u5728\u3057\u307E\u305B\u3093:%s","err_charset_not_found":"\u6587\u5B57\u30BB\u30C3\u30C8\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093","err_imaging_decode_error":"imaging.Decode() \u30A8\u30E9\u30FC","err_imaging_encode_error":"imaging.Encode() \u30A8\u30E9\u30FC","err_jpeg_encode_error":"digestImage jpeg.Encode() \u30A8\u30E9\u30FC","err_imageresize_maxwidth_error":"ImageResizeByMaxWidth \u30A8\u30E9\u30FC maxWidth(%d) > sourceWidth(%d)","err_imageresize_maxheight_error":"ImageResizeByMaxHeight \u30A8\u30E9\u30FC maxHeight(%d) > sourceHeight(%d)","log_starting_server_on_port":"\u30B5\u30FC\u30D0\u30FC\u3092\u8D77\u52D5\u4E2D...\u30DD\u30FC\u30C8 %d ...","log_auto_tls_enabled_for_domain":"\u81EA\u52D5TLS\u304C\u6709\u52B9\u306B\u306A\u308A\u307E\u3057\u305F\u3001\u30C9\u30E1\u30A4\u30F3: %s","log_custom_tls_cert":"\u30AB\u30B9\u30BF\u30E0TLS\u8A3C\u660E\u66F8 CertFile: %s KeyFile: %s","log_server_shutdown_successfully":"\u30B5\u30FC\u30D0\u30FC\u304C\u6B63\u5E38\u306B\u30B7\u30E3\u30C3\u30C8\u30C0\u30A6\u30F3\u3057\u307E\u3057\u305F\u3002\u30B5\u30FC\u30D0\u30FC\u3092\u8D77\u52D5\u4E2D...\u30DD\u30FC\u30C8 %d ...","log_tailscale_server_initialized":"Tailscale\u30B5\u30FC\u30D0\u30FC\u304C\u6B63\u5E38\u306B\u521D\u671F\u5316\u3055\u308C\u307E\u3057\u305F %s:%d","log_working_directory":"\u4F5C\u696D\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA: %s","log_failed_to_add_working_directory_to_store_urls":"\u4F5C\u696D\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3092\u66F8\u5EABURL\u306B\u8FFD\u52A0\u3067\u304D\u307E\u305B\u3093\u3067\u3057\u305F: %s","log_args_index":"args[%d]: %s","log_failed_to_add_store_url_from_args":"\u5F15\u6570\u304B\u3089\u66F8\u5EABURL\u3092\u8FFD\u52A0\u3067\u304D\u307E\u305B\u3093\u3067\u3057\u305F: %s","log_failed_to_add_upload_path_to_store_urls":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u30D1\u30B9\u3092\u66F8\u5EABURL\u306B\u8FFD\u52A0\u3067\u304D\u307E\u305B\u3093\u3067\u3057\u305F: %s","log_starting_tailscale_http_server":"Tailscale HTTP\u30B5\u30FC\u30D0\u30FC\u3092\u8D77\u52D5\u4E2D %s:%d","log_tailscale_server_stopped_successfully":"Tailscale\u30B5\u30FC\u30D0\u30FC\u304C\u6B63\u5E38\u306B\u505C\u6B62\u3057\u307E\u3057\u305F","log_checking_book_files_exist":"\u66F8\u7C4D\u30D5\u30A1\u30A4\u30EB\u306E\u5B58\u5728\u3092\u78BA\u8A8D\u4E2D...","log_checking_store_exist":"\u66F8\u5EAB\u306E\u5B58\u5728\u3092\u78BA\u8A8D\u4E2D...","log_error_deleting_book":"\u66F8\u7C4D %s \u306E\u524A\u9664\u30A8\u30E9\u30FC: %s","log_error_initializing_main_folder":"\u30E1\u30A4\u30F3\u30D5\u30A9\u30EB\u30C0\u30FC\u306E\u521D\u671F\u5316\u30A8\u30E9\u30FC: %s","log_failed_to_get_absolute_path_scan":"\u7D76\u5BFE\u30D1\u30B9\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_skip_to_scan_root_directory":"\u30EB\u30FC\u30C8\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u30B9\u30AD\u30E3\u30F3\u3092\u30B9\u30AD\u30C3\u30D7: %s\u3001%v","log_failed_to_get_relative_path":"\u76F8\u5BFE\u30D1\u30B9\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_skip_to_scan_directory":"\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3092\u30B9\u30AD\u30C3\u30D7: %s\u3001%v","log_failed_to_open_file_get_single":"\u30D5\u30A1\u30A4\u30EB\u3092\u958B\u3051\u307E\u305B\u3093\u3067\u3057\u305F %s: %v","log_timeout_identify_archive_format":"\u64CD\u4F5C\u30BF\u30A4\u30E0\u30A2\u30A6\u30C8\uFF1A\u30A2\u30FC\u30AB\u30A4\u30D6\u5F62\u5F0F\u306E\u8B58\u5225\u306B30\u79D2\u4EE5\u4E0A\u304B\u304B\u308A\u307E\u3057\u305F","log_failed_to_identify_archive_format":"\u30A2\u30FC\u30AB\u30A4\u30D6\u5F62\u5F0F\u306E\u8B58\u5225\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_timeout_create_filesystem":"\u64CD\u4F5C\u30BF\u30A4\u30E0\u30A2\u30A6\u30C8\uFF1A\u30D5\u30A1\u30A4\u30EB\u30B7\u30B9\u30C6\u30E0\u306E\u4F5C\u6210\u306B30\u79D2\u4EE5\u4E0A\u304B\u304B\u308A\u307E\u3057\u305F","log_failed_to_create_filesystem":"\u30D5\u30A1\u30A4\u30EB\u30B7\u30B9\u30C6\u30E0\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_timeout_read_file_content":"\u64CD\u4F5C\u30BF\u30A4\u30E0\u30A2\u30A6\u30C8\uFF1A\u30D5\u30A1\u30A4\u30EB\u5185\u5BB9\u306E\u8AAD\u307F\u8FBC\u307F\u306B30\u79D2\u4EE5\u4E0A\u304B\u304B\u308A\u307E\u3057\u305F","log_failed_to_read_file_content":"\u30D5\u30A1\u30A4\u30EB\u5185\u5BB9\u306E\u8AAD\u307F\u8FBC\u307F\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_timeout_open_file_in_archive":"\u64CD\u4F5C\u30BF\u30A4\u30E0\u30A2\u30A6\u30C8\uFF1A\u30A2\u30FC\u30AB\u30A4\u30D6\u5185\u30D5\u30A1\u30A4\u30EB\u306E\u30AA\u30FC\u30D7\u30F3\u306B30\u79D2\u4EE5\u4E0A\u304B\u304B\u308A\u307E\u3057\u305F","log_failed_to_open_file_in_archive":"\u30A2\u30FC\u30AB\u30A4\u30D6\u5185\u30D5\u30A1\u30A4\u30EB\u306E\u30AA\u30FC\u30D7\u30F3\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_timeout_extract_file":"\u64CD\u4F5C\u30BF\u30A4\u30E0\u30A2\u30A6\u30C8\uFF1A\u30D5\u30A1\u30A4\u30EB\u306E\u62BD\u51FA\u306B30\u79D2\u4EE5\u4E0A\u304B\u304B\u308A\u307E\u3057\u305F","log_failed_to_extract_file":"\u30D5\u30A1\u30A4\u30EB\u306E\u62BD\u51FA\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","err_tailscale_netlistener_nil":"Tailscale netListener\u304Cnil\u3067\u3059\uFF1B\u30B5\u30FC\u30D0\u30FC\u306F\u8D77\u52D5\u3057\u307E\u305B\u3093","err_tailscale_http_server_error":"Tailscale HTTP\u30B5\u30FC\u30D0\u30FC\u30A8\u30E9\u30FC: %v","err_error_closing_network_listener":"\u30CD\u30C3\u30C8\u30EF\u30FC\u30AF\u30EA\u30B9\u30CA\u30FC\u306E\u30AF\u30ED\u30FC\u30BA\u30A8\u30E9\u30FC: %v","err_error_closing_tailscale_server":"Tailscale\u30B5\u30FC\u30D0\u30FC\u306E\u30AF\u30ED\u30FC\u30BA\u30A8\u30E9\u30FC: %v","err_failed_to_create_tailscale_funnel_listener":"Tailscale funnel\u30EA\u30B9\u30CA\u30FC\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F %s: %v","err_failed_to_create_tailscale_listener":"Tailscale\u30EA\u30B9\u30CA\u30FC\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F %s: %v","err_failed_to_create_tailscale_local_client":"Tailscale\u30ED\u30FC\u30AB\u30EB\u30AF\u30E9\u30A4\u30A2\u30F3\u30C8\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","err_error_stopping_tailscale_server":"Tailscale\u30B5\u30FC\u30D0\u30FC\u306E\u505C\u6B62\u30A8\u30E9\u30FC: %v","err_failed_to_set_config_value":"\u8A2D\u5B9A\u5024\u306E\u8A2D\u5B9A\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","err_failed_to_parse_bool":"\'%s\'\u3092bool\u3068\u3057\u3066\u89E3\u6790\u3067\u304D\u307E\u305B\u3093\u3067\u3057\u305F: %v","err_failed_to_parse_int":"\'%s\'\u3092int\u3068\u3057\u3066\u89E3\u6790\u3067\u304D\u307E\u305B\u3093\u3067\u3057\u305F: %v","err_failed_to_set_username":"\u30E6\u30FC\u30B6\u30FC\u540D\u306E\u8A2D\u5B9A\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","err_failed_to_render_userinfo_config_template":"UserInfoConfig\u30C6\u30F3\u30D7\u30EC\u30FC\u30C8\u306E\u30EC\u30F3\u30C0\u30EA\u30F3\u30B0\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","err_failed_to_add_config_value":"\u8A2D\u5B9A\u5024\u306E\u8FFD\u52A0\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","err_getbookshelf_error":"GetBookShelf\u30A8\u30E9\u30FC: %v","err_failed_to_run_tailscale":"Tailscale\u306E\u5B9F\u884C\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","err_store_url_already_exists_error":"\u66F8\u5EABURL\u304C\u65E2\u306B\u5B58\u5728\u3057\u307E\u3059: %s","err_store_path_is_subdir_of_existing":"\u65B0\u3057\u3044\u66F8\u5EAB\u30D1\u30B9\u306F\u65E2\u5B58\u306E\u66F8\u5EAB\u306E\u30B5\u30D6\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3067\u3059: %s \u306F %s \u306E\u30B5\u30D6\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3067\u3059","err_store_path_is_parent_of_existing":"\u65B0\u3057\u3044\u66F8\u5EAB\u30D1\u30B9\u306F\u65E2\u5B58\u306E\u66F8\u5EAB\u306E\u89AA\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3067\u3059: %s \u306F %s \u306E\u89AA\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3067\u3059","err_store_path_conflict":"\u66F8\u5EAB\u30D1\u30B9\u306E\u7AF6\u5408","err_invalid_store_path":"\u7121\u52B9\u306A\u66F8\u5EAB\u30D1\u30B9: %s","err_field_not_exists":"\u30D5\u30A3\u30FC\u30EB\u30C9 \'%s\' \u304C\u5B58\u5728\u3057\u307E\u305B\u3093","err_field_cannot_set":"\u30D5\u30A3\u30FC\u30EB\u30C9 \'%s\' \u3092\u8A2D\u5B9A\u3067\u304D\u307E\u305B\u3093","err_field_not_slice_type":"\u30D5\u30A3\u30FC\u30EB\u30C9 \'%s\' \u306F\u30B9\u30E9\u30A4\u30B9\u578B\u3067\u306F\u3042\u308A\u307E\u305B\u3093","err_field_element_not_string":"\u30D5\u30A3\u30FC\u30EB\u30C9 \'%s\' \u306E\u8981\u7D20\u578B\u304Cstring\u3067\u306F\u3042\u308A\u307E\u305B\u3093","err_field_type_not_supported":"\u30D5\u30A3\u30FC\u30EB\u30C9 \'%s\' \u306E\u578B\u306F\u8A2D\u5B9A\u306B\u5BFE\u5FDC\u3057\u3066\u3044\u307E\u305B\u3093: %s","log_failed_to_create_extract_path":"\u89E3\u51CD\u30D1\u30B9\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_open_file_unarchive":"\u30D5\u30A1\u30A4\u30EB\u3092\u958B\u3051\u307E\u305B\u3093\u3067\u3057\u305F: %v","log_failed_to_identify_file_format":"\u30D5\u30A1\u30A4\u30EB\u5F62\u5F0F\u306E\u8B58\u5225\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_extract_zip_file":"ZIP\u30D5\u30A1\u30A4\u30EB\u306E\u89E3\u51CD\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_zip_file_extracted":"ZIP\u30D5\u30A1\u30A4\u30EB\u89E3\u51CD\u5B8C\u4E86\uFF1A%s \u3092 %s \u306B\u89E3\u51CD","log_failed_to_extract_rar_file":"RAR\u30D5\u30A1\u30A4\u30EB\u306E\u89E3\u51CD\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_rar_file_extracted":"RAR\u30D5\u30A1\u30A4\u30EB\u89E3\u51CD\u5B8C\u4E86\uFF1A%s \u3092 %s \u306B\u89E3\u51CD","log_file_close_error":"file.Close() \u30A8\u30E9\u30FC: %s","log_failed_to_create_directory":"\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_create_parent_directory":"\u89AA\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_create_file":"\u30D5\u30A1\u30A4\u30EB\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_copy_file_content":"\u30D5\u30A1\u30A4\u30EB\u5185\u5BB9\u306E\u30B3\u30D4\u30FC\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_get_container_xml":"container.xml\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_failed_to_parse_container_xml":"container.xml\u306E\u89E3\u6790\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_html_tokenizer_error":"HTML\u30C8\u30FC\u30AF\u30CA\u30A4\u30B6\u30FC\u30A8\u30E9\u30FC: %v","log_failed_to_get_opf_file_path":"OPF\u30D5\u30A1\u30A4\u30EB\u30D1\u30B9\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_failed_to_read_opf_file":"OPF\u30D5\u30A1\u30A4\u30EB\u306E\u8AAD\u307F\u8FBC\u307F\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_failed_to_parse_opf_file":"OPF\u30D5\u30A1\u30A4\u30EB\u306E\u89E3\u6790\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_failed_to_write_file_to_cache":"\u30D5\u30A1\u30A4\u30EB\u3092\u30AD\u30E3\u30C3\u30B7\u30E5\u306B\u66F8\u304D\u8FBC\u307F\u5931\u6557: %v","log_content_type_not_found_in_cache":"\u30AD\u30E3\u30C3\u30B7\u30E5\u5185\u3067ContentType\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3001\u30AD\u30FC: %+v","log_failed_to_read_file_from_cache":"\u30AD\u30E3\u30C3\u30B7\u30E5\u304B\u3089\u30D5\u30A1\u30A4\u30EB\u3092\u8AAD\u307F\u8FBC\u307F\u5931\u6557: %v","log_cover_file_not_found":"\u30AB\u30D0\u30FC\u30D5\u30A1\u30A4\u30EB\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093","log_created_new_book":"\u65B0\u3057\u3044\u66F8\u7C4D\u3092\u4F5C\u6210\u3057\u307E\u3057\u305F: %s","log_updated_existing_book":"\u65E2\u5B58\u306E\u66F8\u7C4D\u3092\u66F4\u65B0\u3057\u307E\u3057\u305F: %s %s","log_saved_media_files_for_book":"\u66F8\u7C4D %s \u306B %d \u500B\u306E\u30E1\u30C7\u30A3\u30A2\u30D5\u30A1\u30A4\u30EB\u3092\u4FDD\u5B58\u3057\u307E\u3057\u305F","log_saved_bookmarks_for_book":"\u66F8\u7C4D %s \u306B %d \u500B\u306E\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF\u3092\u4FDD\u5B58\u3057\u307E\u3057\u305F","log_get_media_files_for_book_error":"\u66F8\u7C4D %s \u306E\u30E1\u30C7\u30A3\u30A2\u30D5\u30A1\u30A4\u30EB\u53D6\u5F97\u30A8\u30E9\u30FC: %s","log_get_bookmarks_for_book_error":"\u66F8\u7C4D %s \u306E\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF\u53D6\u5F97\u30A8\u30E9\u30FC: %s","log_error_creating_new_book_group":"\u65B0\u3057\u3044\u66F8\u7C4D\u30B0\u30EB\u30FC\u30D7\u306E\u4F5C\u6210\u30A8\u30E9\u30FC: %s","log_failed_to_read_directory":"\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u8AAD\u307F\u8FBC\u307F\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s, \u30A8\u30E9\u30FC: %v","log_failed_to_get_file_info_scan":"\u30D5\u30A1\u30A4\u30EB\u60C5\u5831\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s, \u30A8\u30E9\u30FC: %v","log_nonutf8_zip_error":"NonUTF-8 ZIP: %s, \u30A8\u30E9\u30FC: %s","log_failed_to_get_image_list_from_epub":"EPUB\u304B\u3089\u753B\u50CF\u30EA\u30B9\u30C8\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s, \u30A8\u30E9\u30FC: %v","log_failed_to_get_metadata_from_epub":"EPUB\u304B\u3089\u30E1\u30BF\u30C7\u30FC\u30BF\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s, \u30A8\u30E9\u30FC: %v","log_failed_to_access_path_in_archive":"\u30A2\u30FC\u30AB\u30A4\u30D6\u5185\u306E\u30D1\u30B9 %s \u3078\u306E\u30A2\u30AF\u30BB\u30B9\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_skip_scan_path":"\u30B9\u30AD\u30E3\u30F3\u3092\u30B9\u30AD\u30C3\u30D7: %s","log_skip_unsupported_file_type":"\u30B5\u30DD\u30FC\u30C8\u3055\u308C\u3066\u3044\u306A\u3044\u30D5\u30A1\u30A4\u30EB\u30BF\u30A4\u30D7\u3092\u30B9\u30AD\u30C3\u30D7: %s","log_found_supported_file_type":"\u30B5\u30DD\u30FC\u30C8\u3055\u308C\u3066\u3044\u308B\u30D5\u30A1\u30A4\u30EB\u30BF\u30A4\u30D7\uFF08\u975E\u30E1\u30C7\u30A3\u30A2\uFF09\u3092\u767A\u898B: %s","log_processing_file":"\u30D5\u30A1\u30A4\u30EB\u3092\u51E6\u7406\u4E2D: %s (\u30D1\u30B9: %s)","log_failed_to_get_file_info_in_archive":"\u30A2\u30FC\u30AB\u30A4\u30D6\u5185\u306E\u30D5\u30A1\u30A4\u30EB\u60C5\u5831\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_books_saved_to_database_successfully":"SaveBooksToDatabase: %d \u518A\u306E\u66F8\u7C4D\u3092\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u306B\u6B63\u5E38\u306B\u4FDD\u5B58\u3057\u307E\u3057\u305F","log_failed_to_store_bookmark":"\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF\u306E\u4FDD\u5B58\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s","log_updated_bookmarks_for_book_id":"\u66F8\u7C4DID %s \u306E\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF\u3092\u66F4\u65B0\u3057\u307E\u3057\u305F: %s","log_getbook_error_common":"GetBook\u30A8\u30E9\u30FC: %s","log_getpicturedata_error":"GetPictureData\u30A8\u30E9\u30FC: %s","log_getbook_error_scroll":"GetBook: %v","log_server_not_ready_within_timeout":"\u30B5\u30FC\u30D0\u30FC\u304C %v \u4EE5\u5185\u306B\u6E96\u5099\u3067\u304D\u307E\u305B\u3093\u3067\u3057\u305F\u304C\u3001\u7D9A\u884C\u3057\u307E\u3059","log_open_database_error":"OpenDatabase\u30A8\u30E9\u30FC: %s","log_loadbooks_error":"LoadBooks\u30A8\u30E9\u30FC %s","log_savebooks_error":"SaveBooks\u30A8\u30E9\u30FC %s","log_error_listing_books_from_database":"\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u304B\u3089\u66F8\u7C4D\u306E\u4E00\u89A7\u53D6\u5F97\u30A8\u30E9\u30FC: %s","log_failed_savebookstodatabase":"SaveBooksToDatabase\u5931\u6557: %v","log_failed_to_get_free_port":"\u7A7A\u304D\u30DD\u30FC\u30C8\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_using_port":"\u30DD\u30FC\u30C8\u3092\u4F7F\u7528: %d","log_http_newrequest_error":"http.NewRequest\u30A8\u30E9\u30FC: %s","log_client_do_error":"client.Do\u30A8\u30E9\u30FC: %s","log_body_close_error":"Body.Close()\u30A8\u30E9\u30FC: %s","log_io_readall_error":"io.ReadAll\u30A8\u30E9\u30FC: %s","log_response_write_error":"Response.Write\u30A8\u30E9\u30FC: %s","log_getimagefrompdf_imgdata_nil":"GetImageFromPDF: imgData\u304Cnil\u3067\u3059","log_countpages_pdf_invalid_error":"CountPagesOfPDF: \u7121\u52B9\u306APDF: %v \u30A8\u30E9\u30FC: %v","log_getimagefrompdf_time":"GetImageFromPDF: %v","log_syncpage_message_to_flipmode":"SyncPage\u30E1\u30C3\u30BB\u30FC\u30B8\u3092FlipMode\u306B\u9001\u4FE1: %v %v","log_syncpage_message_to_scrollmode":"SyncPage\u30E1\u30C3\u30BB\u30FC\u30B8\u3092ScrollMode\u306B\u9001\u4FE1: %v %v","log_failed_to_open_database":"\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u3092\u958B\u3051\u307E\u305B\u3093\u3067\u3057\u305F: %v","log_failed_to_ping_database":"\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u306EPing\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_failed_to_create_tables":"\u30C6\u30FC\u30D6\u30EB\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_database_initialized_successfully":"\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u304C\u6B63\u5E38\u306B\u521D\u671F\u5316\u3055\u308C\u307E\u3057\u305F","err_failed_to_read_embedded_image":"\u57CB\u3081\u8FBC\u307F\u753B\u50CF\u306E\u8AAD\u307F\u8FBC\u307F\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","err_failed_to_read_embedded_data":"\u57CB\u3081\u8FBC\u307F\u30C7\u30FC\u30BF\u306E\u8AAD\u307F\u8FBC\u307F\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","scroll_wheel_flip":"\u30DB\u30A4\u30FC\u30EB\u30DA\u30FC\u30B8\u3081\u304F\u308A","open_in_new_tab":"\u65B0\u3057\u3044\u30BF\u30D6\u3067\u958B\u304F","auto_rescan_interval_minutes":"\u81EA\u52D5\u518D\u30B9\u30AD\u30E3\u30F3\u9593\u9694","AutoRescanIntervalMinutes":"\u81EA\u52D5\u518D\u30B9\u30AD\u30E3\u30F3\u9593\u9694","auto_rescan_interval_minutes_desc":"\u5206\u5358\u4F4D\u30020\u3092\u6307\u5B9A\u3059\u308B\u3068\u81EA\u52D5\u30B9\u30AD\u30E3\u30F3\u3092\u7121\u52B9\u306B\u3057\u307E\u3059","log_scheduler_task_still_running_skip":"\u524D\u306E\u30B9\u30AD\u30E3\u30F3\u30BF\u30B9\u30AF\u304C\u307E\u3060\u5B9F\u884C\u4E2D\u3067\u3059\u3002\u3053\u306E\u30B9\u30AD\u30E3\u30F3\u3092\u30B9\u30AD\u30C3\u30D7\u3057\u307E\u3059","log_scheduler_task_execution_failed":"\u30B9\u30B1\u30B8\u30E5\u30FC\u30EB\u3055\u308C\u305F\u30B9\u30AD\u30E3\u30F3\u30BF\u30B9\u30AF\u306E\u5B9F\u884C\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_scheduler_task_execution_completed":"\u30B9\u30B1\u30B8\u30E5\u30FC\u30EB\u3055\u308C\u305F\u30B9\u30AD\u30E3\u30F3\u30BF\u30B9\u30AF\u306E\u5B9F\u884C\u304C\u5B8C\u4E86\u3057\u307E\u3057\u305F","log_scheduler_interval_zero_no_scheduled_scan":"\u30B9\u30AD\u30E3\u30F3\u9593\u9694\u304C0\u306E\u305F\u3081\u3001\u30B9\u30B1\u30B8\u30E5\u30FC\u30EB\u3055\u308C\u305F\u30B9\u30AD\u30E3\u30F3\u306F\u5B9F\u884C\u3057\u307E\u305B\u3093","log_scheduler_stop_task_failed":"\u30B9\u30B1\u30B8\u30E5\u30FC\u30EB\u3055\u308C\u305F\u30BF\u30B9\u30AF\u306E\u505C\u6B62\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_scheduler_stop_old_task_failed":"\u53E4\u3044\u30B9\u30B1\u30B8\u30E5\u30FC\u30EB\u3055\u308C\u305F\u30BF\u30B9\u30AF\u306E\u505C\u6B62\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_scheduler_create_scheduler_failed":"\u30B9\u30B1\u30B8\u30E5\u30FC\u30E9\u30FC\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_scheduler_create_task_failed":"\u30B9\u30B1\u30B8\u30E5\u30FC\u30EB\u3055\u308C\u305F\u30BF\u30B9\u30AF\u306E\u4F5C\u6210\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","log_scheduler_task_started":"\u30B9\u30B1\u30B8\u30E5\u30FC\u30EB\u3055\u308C\u305F\u30B9\u30AD\u30E3\u30F3\u30BF\u30B9\u30AF\u304C\u958B\u59CB\u3055\u308C\u307E\u3057\u305F\u3002\u9593\u9694: %d \u5206","log_scheduler_task_stopped":"\u30B9\u30B1\u30B8\u30E5\u30FC\u30EB\u3055\u308C\u305F\u30B9\u30AD\u30E3\u30F3\u30BF\u30B9\u30AF\u304C\u505C\u6B62\u3057\u307E\u3057\u305F","auto_rescan_started":"\u81EA\u52D5\u518D\u30B9\u30AD\u30E3\u30F3\u304C\u958B\u59CB\u3055\u308C\u307E\u3057\u305F\u3002\u9593\u9694: %d \u5206","auto_rescan_stopped":"\u81EA\u52D5\u518D\u30B9\u30AD\u30E3\u30F3\u304C\u505C\u6B62\u3057\u307E\u3057\u305F","auto_rescan_updated":"\u81EA\u52D5\u518D\u30B9\u30AD\u30E3\u30F3\u9593\u9694\u304C\u66F4\u65B0\u3055\u308C\u307E\u3057\u305F: %d \u5206","auto_rescan_status_running":"\u5B9F\u884C\u4E2D","auto_rescan_status_stopped":"\u505C\u6B62","auto_rescan_current_interval":"\u73FE\u5728\u306E\u30B9\u30AD\u30E3\u30F3\u9593\u9694: %d \u5206","auto_rescan_enabled_hint":"\u81EA\u52D5\u518D\u30B9\u30AD\u30E3\u30F3\u304C\u6709\u52B9\u306B\u306A\u308A\u307E\u3057\u305F\u3002\u30B7\u30B9\u30C6\u30E0\u306F\u5B9A\u671F\u7684\u306B\u66F8\u5EAB\u3092\u30B9\u30AD\u30E3\u30F3\u3057\u307E\u3059","auto_rescan_disabled_hint":"\u81EA\u52D5\u518D\u30B9\u30AD\u30E3\u30F3\u304C\u7121\u52B9\u306B\u306A\u308A\u307E\u3057\u305F","other_settings":"\u305D\u306E\u4ED6\u306E\u8A2D\u5B9A","plugin_enable":"\u30D7\u30E9\u30B0\u30A4\u30F3\u3092\u6709\u52B9\u306B\u3059\u308B","err_failed_to_get_config_dir":"\u8A2D\u5B9A\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %v","comigo_xyz_title":"Comigo","comigo_xyz_subtitle":"Comigo - \u30B7\u30F3\u30D7\u30EB\u306A\u30B3\u30DF\u30C3\u30AF\u30EA\u30FC\u30C0\u30FC","comigo_xyz_description":"\u3059\u3079\u3066\u306E\u30C7\u30D0\u30A4\u30B9\u3067\u6F2B\u753B\u3092\u8AAD\u3080 - \uFF30\uFF23\u3067\u3082\u30B9\u30DE\u30FC\u30C8\u30D5\u30A9\u30F3\u3067\u3082\u3002","comigo_xyz_feature_format":"\uD83D\uDCDA \u8907\u6570\u30D5\u30A9\u30FC\u30DE\u30C3\u30C8\u5BFE\u5FDC","comigo_xyz_feature_format_desc":"ZIP\u3001RAR\u3001CBZ\u3001EPUB\u3001PDF \u306A\u3069\u3001\u3055\u307E\u3056\u307E\u306A\u30B3\u30DF\u30C3\u30AF\u30D5\u30A9\u30FC\u30DE\u30C3\u30C8\u3092\u30B5\u30DD\u30FC\u30C8","comigo_xyz_feature_cross_platform":"\uD83C\uDF10 \u30AF\u30ED\u30B9\u30D7\u30E9\u30C3\u30C8\u30D5\u30A9\u30FC\u30E0","comigo_xyz_feature_cross_platform_desc":"Linux\u3001Windows\u3001Mac OS \u306E\u30AA\u30DA\u30EC\u30FC\u30C6\u30A3\u30F3\u30B0\u30B7\u30B9\u30C6\u30E0\u3092\u30B5\u30DD\u30FC\u30C8","comigo_xyz_feature_responsive":"\uD83D\uDCF1 \u30EC\u30B9\u30DD\u30F3\u30B7\u30D6\u30C7\u30B6\u30A4\u30F3","comigo_xyz_feature_responsive_desc":"\u30C7\u30B9\u30AF\u30C8\u30C3\u30D7\u3068\u30E2\u30D0\u30A4\u30EB\u30C7\u30D0\u30A4\u30B9\u306B\u5BFE\u5FDC","comigo_xyz_feature_reading_modes":"\uD83D\uDD04 \u8907\u6570\u306E\u8AAD\u66F8\u30E2\u30FC\u30C9","comigo_xyz_feature_reading_modes_desc":"\u30DA\u30FC\u30B8\u3081\u304F\u308A\u30E2\u30FC\u30C9\u3068\u30B9\u30AF\u30ED\u30FC\u30EB\u30E2\u30FC\u30C9\u3092\u30B5\u30DD\u30FC\u30C8\u3057\u3001\u3055\u307E\u3056\u307E\u306A\u8AAD\u66F8\u7FD2\u6163\u306B\u5BFE\u5FDC","comigo_xyz_feature_security":"\uD83D\uDD12 \u5B89\u5168\u3067\u4FE1\u983C\u6027\u304C\u9AD8\u3044","comigo_xyz_feature_security_desc":"HTTPS \u3068\u30E6\u30FC\u30B6\u30FC\u8A8D\u8A3C\u3092\u30B5\u30DD\u30FC\u30C8\u3057\u3001Tailscale \u30EA\u30E2\u30FC\u30C8\u30A2\u30AF\u30BB\u30B9\u3092\u5185\u8535","comigo_xyz_feature_performance":"\u26A1 \u9AD8\u6027\u80FD","comigo_xyz_feature_performance_desc":"\u9AD8\u901F\u30B9\u30AD\u30E3\u30F3\u3068\u8AAD\u307F\u8FBC\u307F\u3067\u3001\u30B9\u30E0\u30FC\u30BA\u306A\u8AAD\u66F8\u4F53\u9A13\u3092\u5B9F\u73FE","comigo_xyz_feature_plugin":"\uD83D\uDD0C \u30D7\u30E9\u30B0\u30A4\u30F3\u30B7\u30B9\u30C6\u30E0","comigo_xyz_feature_plugin_desc":"\u81EA\u52D5\u30DA\u30FC\u30B8\u3081\u304F\u308A\u3001\u6642\u8A08\u306A\u3069\u306E\u30D7\u30E9\u30B0\u30A4\u30F3\u306B\u5BFE\u5FDC\u3001\u30AB\u30B9\u30BF\u30E0\u30D7\u30E9\u30B0\u30A4\u30F3\u306E\u62E1\u5F35\u3082\u53EF\u80FD","comigo_xyz_feature_media":"\uD83C\uDFAC \u30E1\u30C7\u30A3\u30A2\u518D\u751F","comigo_xyz_feature_media_desc":"\u5185\u8535\u30AA\u30FC\u30C7\u30A3\u30AA\u30FB\u30D3\u30C7\u30AA\u30D7\u30EC\u30FC\u30E4\u30FC","comigo_xyz_feature_download":"\uD83D\uDCE5 \u67D4\u8EDF\u306A\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","comigo_xyz_feature_download_desc":"\u753B\u50CF\u30D5\u30A9\u30EB\u30C0\u306E\u4E00\u62EC\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9\u3001EPUB\u30D5\u30A9\u30FC\u30DE\u30C3\u30C8\u3078\u306E\u5909\u63DB\u306B\u5BFE\u5FDC","comigo_xyz_feature_history":"\uD83D\uDCDC \u95B2\u89A7\u5C65\u6B74","comigo_xyz_feature_history_desc":"\u95B2\u89A7\u5C65\u6B74\u3092\u81EA\u52D5\u8A18\u9332\u3001\u7D9A\u304D\u304B\u3089\u8AAD\u3081\u308B","comigo_xyz_github_repo":"\uD83D\uDCE6 GitHub \u30EA\u30DD\u30B8\u30C8\u30EA","comigo_xyz_github_button":"GitHub \u30D7\u30ED\u30B8\u30A7\u30AF\u30C8\u306B\u30A2\u30AF\u30BB\u30B9","comigo_xyz_download_latest":"\u6700\u65B0\u30D0\u30FC\u30B8\u30E7\u30F3\u3092\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","comigo_xyz_quick_start":"\u30AF\u30A4\u30C3\u30AF\u30B9\u30BF\u30FC\u30C8","comigo_xyz_quick_start_step1":"Comigo \u3092\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9\u3057\u3066\u5B9F\u884C","comigo_xyz_quick_start_step2":"\u8A2D\u5B9A\u3067\u30E9\u30A4\u30D6\u30E9\u30EA\u30D1\u30B9\u3092\u69CB\u6210","comigo_xyz_quick_start_step3":"\u8AAD\u66F8\u4F53\u9A13\u3092\u304A\u697D\u3057\u307F\u304F\u3060\u3055\u3044\uFF01","comigo_xyz_download_for":"\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","comigo_xyz_select_platform":"\u30D7\u30E9\u30C3\u30C8\u30D5\u30A9\u30FC\u30E0\u3092\u9078\u629E","comigo_xyz_download":"\u2B07\uFE0F \u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","comigo_xyz_cli_install_title":"\u30EF\u30F3\u30AF\u30EA\u30C3\u30AF CLI \u30A4\u30F3\u30B9\u30C8\u30FC\u30EB\uFF1A","comigo_xyz_cli_install_desc":"\u4EE5\u4E0B\u306E\u30B3\u30DE\u30F3\u30C9\u3092\u4F7F\u7528\u3057\u3066 Comigo \u3092\u3059\u3070\u3084\u304F\u30A4\u30F3\u30B9\u30C8\u30FC\u30EB","comigo_xyz_cli_install_copy":"\u30B3\u30D4\u30FC","comigo_xyz_cli_install_copied":"\u30B3\u30D4\u30FC\u3057\u307E\u3057\u305F","comigo_xyz_cli_install_cn_desc":"\u4E2D\u56FD\u5927\u9678\u30E6\u30FC\u30B6\u30FC\u306B\u304A\u3059\u3059\u3081\uFF1A","comigo_xyz_docker_deploy_title":"Docker\u3067\u30C7\u30D7\u30ED\u30A4\uFF1A","rescan_store":"\u66F8\u5EAB\u30B9\u30AD\u30E3\u30F3","delete_store":"\u66F8\u5EAB\u524A\u9664","select_store_to_operate":"\u64CD\u4F5C\u3059\u308B\u66F8\u5EAB\u3092\u9078\u629E\u3057\u3066\u304F\u3060\u3055\u3044","confirm_delete_store":"\u66F8\u5EAB\u3092\u524A\u9664\u3057\u307E\u3059\u304B\uFF1F\u3053\u308C\u306B\u3088\u308A\u3001\u3053\u306E\u66F8\u5EAB\u306E\u3059\u3079\u3066\u306E\u66F8\u7C4D\u30C7\u30FC\u30BF\u3082\u524A\u9664\u3055\u308C\u307E\u3059","rescan_store_success":"\u66F8\u5EAB\u306E\u30B9\u30AD\u30E3\u30F3\u304C\u5B8C\u4E86\u3057\u307E\u3057\u305F\u3002{0} \u518A\u306E\u65B0\u3057\u3044\u66F8\u7C4D\u304C\u8FFD\u52A0\u3055\u308C\u307E\u3057\u305F","delete_store_success":"\u66F8\u5EAB\u304C\u6B63\u5E38\u306B\u524A\u9664\u3055\u308C\u307E\u3057\u305F","rescan_store_in_progress":"\u66F8\u5EAB\u3092\u30B9\u30AD\u30E3\u30F3\u3057\u3066\u3044\u307E\u3059\u3002\u304A\u5F85\u3061\u304F\u3060\u3055\u3044...","err_rescan_store_failed":"\u66F8\u5EAB\u306E\u518D\u30B9\u30AD\u30E3\u30F3\u306B\u5931\u6557\u3057\u307E\u3057\u305F","err_delete_store_failed":"\u66F8\u5EAB\u306E\u524A\u9664\u306B\u5931\u6557\u3057\u307E\u3057\u305F","log_book_file_not_exist_skip":"\u66F8\u7C4D\u30D5\u30A1\u30A4\u30EB\u304C\u5B58\u5728\u3057\u306A\u3044\u305F\u3081\u3001\u8AAD\u307F\u8FBC\u307F\u3092\u30B9\u30AD\u30C3\u30D7\u3057\u307E\u3057\u305F: %s","log_error_deleting_orphan_metadata":"\u5B64\u7ACB\u3057\u305F\u30E1\u30BF\u30C7\u30FC\u30BF\u30D5\u30A1\u30A4\u30EB %s \u306E\u524A\u9664\u30A8\u30E9\u30FC: %s","log_book_version_mismatch_skip":"\u66F8\u7C4D %s \u306E\u30D0\u30FC\u30B8\u30E7\u30F3\u304C\u4E00\u81F4\u3057\u307E\u305B\u3093\uFF08\u30AD\u30E3\u30C3\u30B7\u30E5: %s\u3001\u73FE\u5728: %s\uFF09\u3001\u8AAD\u307F\u8FBC\u307F\u3092\u30B9\u30AD\u30C3\u30D7","log_book_version_minor_mismatch":"\u66F8\u7C4D %s \u306E\u30DE\u30A4\u30CA\u30FC\u30D0\u30FC\u30B8\u30E7\u30F3\u304C\u7570\u306A\u308A\u307E\u3059\uFF08\u30AD\u30E3\u30C3\u30B7\u30E5: %s\u3001\u73FE\u5728: %s\uFF09\u3001\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF\u3092\u79FB\u884C\u3057\u3066\u518D\u30B9\u30AD\u30E3\u30F3\u3057\u307E\u3059","log_bookmark_saved_for_migration":"\u66F8\u7C4D %s \u306E %d \u500B\u306E\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF\u3092\u4FDD\u5B58\u3057\u307E\u3057\u305F\u3001\u79FB\u884C\u5F85\u3061","log_bookmark_migrated":"\u66F8\u7C4D %s \u306E %d \u500B\u306E\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF\u3092\u6B63\u5E38\u306B\u79FB\u884C\u3057\u307E\u3057\u305F","log_error_deleting_version_mismatch_metadata":"\u30D0\u30FC\u30B8\u30E7\u30F3\u4E0D\u4E00\u81F4\u306E\u30E1\u30BF\u30C7\u30FC\u30BF\u30D5\u30A1\u30A4\u30EB %s \u306E\u524A\u9664\u30A8\u30E9\u30FC: %s","player":"\u30D7\u30EC\u30FC\u30E4\u30FC","play":"\u518D\u751F","pause":"\u4E00\u6642\u505C\u6B62","previous":"\u524D\u3078","next":"\u6B21\u3078","playlist":"\u30D7\u30EC\u30A4\u30EA\u30B9\u30C8","audio":"\u30AA\u30FC\u30C7\u30A3\u30AA","video":"\u30D3\u30C7\u30AA","browser_not_support_video":"\u304A\u4F7F\u3044\u306E\u30D6\u30E9\u30A6\u30B6\u306F\u30D3\u30C7\u30AA\u518D\u751F\u306B\u5BFE\u5FDC\u3057\u3066\u3044\u307E\u305B\u3093","browser_not_support_audio":"\u304A\u4F7F\u3044\u306E\u30D6\u30E9\u30A6\u30B6\u306F\u30AA\u30FC\u30C7\u30A3\u30AA\u518D\u751F\u306B\u5BFE\u5FDC\u3057\u3066\u3044\u307E\u305B\u3093","play_failed":"\u518D\u751F\u306B\u5931\u6557\u3057\u307E\u3057\u305F","first_media":"\u6700\u521D\u306E\u30A2\u30A4\u30C6\u30E0\u3067\u3059","last_media":"\u6700\u5F8C\u306E\u30A2\u30A4\u30C6\u30E0\u3067\u3059","auto_play_next":"\u6B21\u3092\u81EA\u52D5\u518D\u751F","loop_playlist":"\u30D7\u30EC\u30A4\u30EA\u30B9\u30C8\u3092\u30EB\u30FC\u30D7","player_autoplay_help":"\u6B21\u3092\u81EA\u52D5\u518D\u751F\uFF0F\u30EB\u30FC\u30D7\u518D\u751F\u306E\u6A5F\u80FD\u306B\u306F\u3001\u30D6\u30E9\u30A6\u30B6\u3067\u30E1\u30C7\u30A3\u30A2\u518D\u751F\u306E\u8A31\u53EF\u304C\u5FC5\u8981\u3067\u3059\u3002\u30E2\u30D0\u30A4\u30EB\u7AEF\u672B\u3067\u306F\u7701\u96FB\u529B\u3084\u30D0\u30C3\u30AF\u30B0\u30E9\u30A6\u30F3\u30C9\u5236\u9650\u306B\u3088\u308A\u52D5\u4F5C\u3057\u306A\u3044\u3053\u3068\u3082\u3042\u308A\u307E\u3059\u3002","add_bookmark":"\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF\u3092\u8FFD\u52A0","manual_bookmark":"\u624B\u52D5\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF","auto_bookmark":"\u81EA\u52D5\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF","bookmark_added":"\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF\u3092\u8FFD\u52A0\u3057\u307E\u3057\u305F","bookmark_exists":"\u3053\u306E\u30DA\u30FC\u30B8\u306B\u306F\u3059\u3067\u306B\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF\u304C\u3042\u308A\u307E\u3059","continue_reading":"\u7D9A\u304D\u3092\u8AAD\u3080","delete_record":"\u524A\u9664","confirm_delete_bookmark":"\u3053\u306E\u95B2\u89A7\u8A18\u9332\u3092\u524A\u9664\u3057\u307E\u3059\u304B\uFF1F","bookmark_deleted":"\u524A\u9664\u3057\u307E\u3057\u305F","log_remote_store_connect_failed":"\u30EA\u30E2\u30FC\u30C8\u66F8\u5EAB\u3078\u306E\u63A5\u7D9A\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s, \u30A8\u30E9\u30FC: %v","log_remote_file_stat_failed":"\u30EA\u30E2\u30FC\u30C8\u30D5\u30A1\u30A4\u30EB\u60C5\u5831\u306E\u53D6\u5F97\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s, \u30A8\u30E9\u30FC: %v","log_remote_file_open_failed":"\u30EA\u30E2\u30FC\u30C8\u30D5\u30A1\u30A4\u30EB\u306E\u30AA\u30FC\u30D7\u30F3\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s, \u30A8\u30E9\u30FC: %v","log_remote_pdf_download_on_demand":"\u30EA\u30E2\u30FC\u30C8PDF\u3092\u30AA\u30F3\u30C7\u30DE\u30F3\u30C9\u3067\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9\u4E2D: %s","log_webdav_filesystem_connected":"WebDAV\u30D5\u30A1\u30A4\u30EB\u30B7\u30B9\u30C6\u30E0\u306B\u63A5\u7D9A\u3057\u307E\u3057\u305F: %s, \u30D9\u30FC\u30B9\u30D1\u30B9: %s","log_sftp_filesystem_connected":"SFTP\u30D5\u30A1\u30A4\u30EB\u30B7\u30B9\u30C6\u30E0\u306B\u63A5\u7D9A\u3057\u307E\u3057\u305F: %s, \u30D9\u30FC\u30B9\u30D1\u30B9: %s","log_smb_filesystem_connected":"SMB\u30D5\u30A1\u30A4\u30EB\u30B7\u30B9\u30C6\u30E0\u306B\u63A5\u7D9A\u3057\u307E\u3057\u305F: %s, \u30D9\u30FC\u30B9\u30D1\u30B9: %s","log_ftp_filesystem_connected":"FTP\u30D5\u30A1\u30A4\u30EB\u30B7\u30B9\u30C6\u30E0\u306B\u63A5\u7D9A\u3057\u307E\u3057\u305F: %s, \u30D9\u30FC\u30B9\u30D1\u30B9: %s","log_s3_filesystem_connected":"S3\u30D5\u30A1\u30A4\u30EB\u30B7\u30B9\u30C6\u30E0\u306B\u63A5\u7D9A\u3057\u307E\u3057\u305F: %s, \u30D9\u30FC\u30B9\u30D1\u30B9: %s","log_epub_metadata_remote_not_supported":"EPUB\u30E1\u30BF\u30C7\u30FC\u30BF\u62BD\u51FA\u306F\u30EA\u30E2\u30FC\u30C8\u30B9\u30C8\u30EA\u30FC\u30DF\u30F3\u30B0\u3067\u306F\u30B5\u30DD\u30FC\u30C8\u3055\u308C\u3066\u3044\u307E\u305B\u3093","log_remote_file_download_to_cache":"\u30EA\u30E2\u30FC\u30C8\u30D5\u30A1\u30A4\u30EB\u3092\u30AD\u30E3\u30C3\u30B7\u30E5\u306B\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9\u4E2D: %s -> %s","log_scan_remote_store_start":"\u30EA\u30E2\u30FC\u30C8\u66F8\u5EAB\u306E\u30B9\u30AD\u30E3\u30F3\u3092\u958B\u59CB: %s","log_scan_start_hint_remote":"\u30B9\u30AD\u30E3\u30F3\u958B\u59CB\uFF1A%s (\u30EA\u30E2\u30FC\u30C8\u30D1\u30B9: %s)","log_remote_store_check_book_existence_failed":"\u66F8\u7C4D\u306E\u5B58\u5728\u78BA\u8A8D\u306E\u305F\u3081\u30EA\u30E2\u30FC\u30C8\u66F8\u5EAB\u306B\u63A5\u7D9A\u3067\u304D\u307E\u305B\u3093\u3067\u3057\u305F: %s, \u30A8\u30E9\u30FC: %v","log_remote_book_existence_check_failed":"\u30EA\u30E2\u30FC\u30C8\u66F8\u7C4D\u306E\u5B58\u5728\u78BA\u8A8D\u306B\u5931\u6557\u3057\u307E\u3057\u305F: %s, \u30A8\u30E9\u30FC: %v","log_add_remote_store":"\u30EA\u30E2\u30FC\u30C8\u66F8\u5EAB\u3092\u8FFD\u52A0: %s (\u30D7\u30ED\u30C8\u30B3\u30EB: %s, \u30DB\u30B9\u30C8: %s)","log_s3_connecting":"S3 \u30B5\u30FC\u30D3\u30B9\u306B\u63A5\u7D9A\u4E2D %s (\u30D0\u30B1\u30C3\u30C8: %s, \u30D7\u30EC\u30D5\u30A3\u30C3\u30AF\u30B9: %s)","log_ftp_connecting":"FTP \u30B5\u30FC\u30D0\u30FC\u306B\u63A5\u7D9A\u4E2D %s (TLS: %v, \u30BF\u30A4\u30E0\u30A2\u30A6\u30C8: %v)","log_webdav_download_range":"\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9\u7BC4\u56F2: %s [%d-%d]","log_smb_connecting":"SMB \u30B5\u30FC\u30D0\u30FC\u306B\u63A5\u7D9A\u4E2D %s (\u30BF\u30A4\u30E0\u30A2\u30A6\u30C8: %d\u79D2, \u30E6\u30FC\u30B6\u30FC: %s, \u5171\u6709: %s)","log_smb_mount_share":"SMB \u5171\u6709\u3092\u30DE\u30A6\u30F3\u30C8\u4E2D: %s","log_using_cached_file":"\u30AD\u30E3\u30C3\u30B7\u30E5\u30D5\u30A1\u30A4\u30EB\u3092\u4F7F\u7528: %s","log_download_file":"\u30D5\u30A1\u30A4\u30EB\u3092\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9\u4E2D: %s","log_remote_book_existence_check_failed_detail":"\u30EA\u30E2\u30FC\u30C8\u66F8\u7C4D\u306E\u5B58\u5728\u78BA\u8A8D\u306B\u5931\u6557 - BookID: %s, RemoteURL: %s, BookPath: %s, \u30A8\u30E9\u30FC: %v","log_local_book_existence_check_failed":"\u30ED\u30FC\u30AB\u30EB\u66F8\u7C4D\u306E\u5B58\u5728\u78BA\u8A8D\u306B\u5931\u6557: %s, \u30A8\u30E9\u30FC: %v","log_failed_to_set_field":"\u30D5\u30A3\u30FC\u30EB\u30C9 %s \u306E\u8A2D\u5B9A\u306B\u5931\u6557: %v","log_cache_mkdir_failed":"\u30AD\u30E3\u30C3\u30B7\u30E5\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u4F5C\u6210\u306B\u5931\u6557: %v","log_cache_hit_memory":"\u30E1\u30E2\u30EA\u30AD\u30E3\u30C3\u30B7\u30E5\u306B\u30D2\u30C3\u30C8: %s","log_cache_hit_disk":"\u30C7\u30A3\u30B9\u30AF\u30AD\u30E3\u30C3\u30B7\u30E5\u306B\u30D2\u30C3\u30C8: %s","log_cache_write_disk_failed":"\u30C7\u30A3\u30B9\u30AF\u30AD\u30E3\u30C3\u30B7\u30E5\u3078\u306E\u66F8\u304D\u8FBC\u307F\u306B\u5931\u6557: %v","log_cached_to_disk":"\u30C7\u30A3\u30B9\u30AF\u306B\u30AD\u30E3\u30C3\u30B7\u30E5\u3057\u307E\u3057\u305F: %s -> %s","log_get_file_error":"\u30D5\u30A1\u30A4\u30EB\u53D6\u5F97\u30A8\u30E9\u30FC: %s","log_save_file_to_cache_error":"SaveFileToCache \u30A8\u30E9\u30FC: %s","log_save_cover_to_local_error":"SaveCoverToLocal \u30A8\u30E9\u30FC: %s","log_failed_to_parse_cover_url":"\u30AB\u30D0\u30FC URL \u306E\u89E3\u6790\u306B\u5931\u6557: %s","log_child_book_id_missing_in_cover_url":"\u30AB\u30D0\u30FC URL \u306B\u5B50\u66F8\u7C4D ID \u304C\u3042\u308A\u307E\u305B\u3093","log_failed_to_get_child_book":"\u5B50\u66F8\u7C4D\u306E\u53D6\u5F97\u306B\u5931\u6557: %s","log_failed_to_read_embedded_image":"\u57CB\u3081\u8FBC\u307F\u753B\u50CF\u306E\u8AAD\u307F\u8FBC\u307F\u306B\u5931\u6557: %s","log_failed_to_create_epub_generator":"EPUB \u30B8\u30A7\u30CD\u30EC\u30FC\u30BF\u306E\u4F5C\u6210\u306B\u5931\u6557: %s","log_failed_to_get_image_epub":"\u753B\u50CF\u306E\u53D6\u5F97\u306B\u5931\u6557 %s: %v","log_failed_to_generate_epub":"EPUB \u306E\u751F\u6210\u306B\u5931\u6557: %s","log_plugin_enabled":"\u30D7\u30E9\u30B0\u30A4\u30F3\u3092\u6709\u52B9\u306B\u3057\u307E\u3057\u305F: %s","log_disable_mutex_plugin_auto_flip":"\u76F8\u4E92\u6392\u4ED6\u30D7\u30E9\u30B0\u30A4\u30F3\u3092\u7121\u52B9\u306B: auto_flip","log_disable_mutex_plugin_sketch_practice":"\u76F8\u4E92\u6392\u4ED6\u30D7\u30E9\u30B0\u30A4\u30F3\u3092\u7121\u52B9\u306B: sketch_practice","log_plugin_disabled":"\u30D7\u30E9\u30B0\u30A4\u30F3\u3092\u7121\u52B9\u306B\u3057\u307E\u3057\u305F: %s","log_rescan_store":"\u66F8\u5EAB\u3092\u518D\u30B9\u30AD\u30E3\u30F3\u4E2D: %s","log_rescan_store_completed_new_books":"\u66F8\u5EAB\u30B9\u30AD\u30E3\u30F3\u5B8C\u4E86\u3001%d \u518A\u306E\u65B0\u898F\u66F8\u7C4D","log_delete_store":"\u66F8\u5EAB\u3092\u524A\u9664\u4E2D: %s","log_deleted_books_count":"%d \u518A\u306E\u66F8\u7C4D\u3092\u524A\u9664\u3057\u307E\u3057\u305F","log_upload_no_store_selected":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u5148\u306E\u66F8\u5EAB\u304C\u9078\u629E\u3055\u308C\u3066\u3044\u307E\u305B\u3093","log_upload_invalid_store_path":"\u7121\u52B9\u306A\u66F8\u5EAB\u30D1\u30B9: %s","log_upload_store_path_not_exist":"\u66F8\u5EAB\u30D1\u30B9\u304C\u5B58\u5728\u3057\u307E\u305B\u3093: %s","log_error_closing_zip_writer":"zip \u30E9\u30A4\u30BF\u30FC\u306E\u30AF\u30ED\u30FC\u30BA\u30A8\u30E9\u30FC: %s","log_file_not_found_skipping":"\u30D5\u30A1\u30A4\u30EB\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3001\u30B9\u30AD\u30C3\u30D7: %s","log_error_opening_file":"\u30D5\u30A1\u30A4\u30EB\u3092\u958B\u304F\u30A8\u30E9\u30FC: %s, \u30A8\u30E9\u30FC: %s","log_error_creating_zip_entry":"zip \u30A8\u30F3\u30C8\u30EA\u4F5C\u6210\u30A8\u30E9\u30FC: %s, \u30A8\u30E9\u30FC: %s","log_error_writing_file_to_zip":"zip \u3078\u306E\u66F8\u304D\u8FBC\u307F\u30A8\u30E9\u30FC: %s, \u30A8\u30E9\u30FC: %s","log_failed_to_delete_bookmark":"\u30D6\u30C3\u30AF\u30DE\u30FC\u30AF\u306E\u524A\u9664\u306B\u5931\u6557: %v","log_get_config_dir_error":"GetConfigDir \u30A8\u30E9\u30FC: %s","log_delete_book_json_error":"DeleteBookJson \u30A8\u30E9\u30FC: %s","log_delete_cover_cache_error":"DeleteCoverCache \u30A8\u30E9\u30FC: %s","log_delete_book_cache_error":"DeleteBookCache \u30A8\u30E9\u30FC: %s","log_plugin_system_disabled_skip_scan":"\u30D7\u30E9\u30B0\u30A4\u30F3\u30B7\u30B9\u30C6\u30E0\u304C\u7121\u52B9\u306E\u305F\u3081\u3001\u30AB\u30B9\u30BF\u30E0\u30D7\u30E9\u30B0\u30A4\u30F3\u306E\u30B9\u30AD\u30E3\u30F3\u3092\u30B9\u30AD\u30C3\u30D7\u3057\u307E\u3059","log_plugin_dir_not_exist_skip_load":"\u30D7\u30E9\u30B0\u30A4\u30F3\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u304C\u5B58\u5728\u3057\u307E\u305B\u3093: %s\u3001\u30AB\u30B9\u30BF\u30E0\u30D7\u30E9\u30B0\u30A4\u30F3\u306E\u8AAD\u307F\u8FBC\u307F\u3092\u30B9\u30AD\u30C3\u30D7\u3057\u307E\u3059","log_plugin_scope_load_error":"\u30B9\u30B3\u30FC\u30D7 %s \u306E\u30D7\u30E9\u30B0\u30A4\u30F3\u8AAD\u307F\u8FBC\u307F\u3067\u30A8\u30E9\u30FC: %v","log_plugin_custom_loaded_count":"\u30AB\u30B9\u30BF\u30E0\u30D7\u30E9\u30B0\u30A4\u30F3 %d \u500B\u306E\u8AAD\u307F\u8FBC\u307F\u306B\u6210\u529F\u3057\u307E\u3057\u305F","log_plugin_loaded_item":"  - [%s] %s (%s)","log_plugin_dir_not_exist":"\u30D7\u30E9\u30B0\u30A4\u30F3\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u304C\u5B58\u5728\u3057\u307E\u305B\u3093: %s","log_plugin_read_file_failed":"\u30D7\u30E9\u30B0\u30A4\u30F3\u30D5\u30A1\u30A4\u30EB\u306E\u8AAD\u307F\u8FBC\u307F\u306B\u5931\u6557 %s: %v","log_plugin_read_book_file_failed":"\u66F8\u7C4D\u30D7\u30E9\u30B0\u30A4\u30F3\u30D5\u30A1\u30A4\u30EB\u306E\u8AAD\u307F\u8FBC\u307F\u306B\u5931\u6557 %s: %v","log_plugin_loaded_for_book":"\u66F8\u7C4D %s \u306E %s \u30D7\u30E9\u30B0\u30A4\u30F3 %d \u500B\u3092\u8AAD\u307F\u8FBC\u307F\u307E\u3057\u305F","log_load_custom_plugin_failed":"\u30AB\u30B9\u30BF\u30E0\u30D7\u30E9\u30B0\u30A4\u30F3\u306E\u8AAD\u307F\u8FBC\u307F\u306B\u5931\u6557: %v","log_cfg_host_enabled_plugin_list":"cfg.Host: %v , cfg.EnabledPluginList: %v","log_config_changed_restart_web":"\u8A2D\u5B9A\u304C\u5909\u66F4\u3055\u308C\u307E\u3057\u305F\u3001Web \u30B5\u30FC\u30D0\u30FC\u3092\u518D\u8D77\u52D5\u3057\u3066\u3044\u307E\u3059...","log_config_changed_start_tailscale":"\u8A2D\u5B9A\u304C\u5909\u66F4\u3055\u308C\u307E\u3057\u305F\u3001Tailscale \u3092\u8D77\u52D5\u3057\u3066\u3044\u307E\u3059...","log_config_changed_stop_tailscale":"\u8A2D\u5B9A\u304C\u5909\u66F4\u3055\u308C\u307E\u3057\u305F\u3001Tailscale \u3092\u505C\u6B62\u3057\u3066\u3044\u307E\u3059...","log_config_changed_restart_tailscale":"\u8A2D\u5B9A\u304C\u5909\u66F4\u3055\u308C\u307E\u3057\u305F\u3001Tailscale \u3092\u518D\u8D77\u52D5\u3057\u3066\u3044\u307E\u3059...","log_received_rescan_message":"\u518D\u30B9\u30AD\u30E3\u30F3\u30E1\u30C3\u30BB\u30FC\u30B8\u3092\u53D7\u4FE1: %s","log_failed_to_read_image_epub":"\u753B\u50CF\u306E\u8AAD\u307F\u8FBC\u307F\u306B\u5931\u6557 %s: %v","log_failed_to_decode_image_config_epub":"\u753B\u50CF\u8A2D\u5B9A\u306E\u30C7\u30B3\u30FC\u30C9\u306B\u5931\u6557: %v\u3001\u30C7\u30D5\u30A9\u30EB\u30C8\u5C3A\u5BF8\u3092\u4F7F\u7528\u3057\u307E\u3059","log_cannot_set_username_or_password":"\u30E6\u30FC\u30B6\u30FC\u540D\u307E\u305F\u306F\u30D1\u30B9\u30EF\u30FC\u30C9\u3092\u8A2D\u5B9A\u3067\u304D\u307E\u305B\u3093","log_login_failed":"\u30ED\u30B0\u30A4\u30F3\u306B\u5931\u6557: %s-%s, %s-%s","log_get_generated_image_params":"GetGeneratedImage: height=%s, width=%s, text=%s, font_size=%s","log_toml_marshal_error":"toml.Marshal \u30A8\u30E9\u30FC","log_config_file_changed_detected":"Comigo \u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u306E\u5909\u66F4\u3092\u691C\u51FA: %s Op: %s","log_error_num_value":"\u30A8\u30E9\u30FC\u6570\u5024: %s"}');


(0, $1bac384020b50752$export$2e2bcd8739ae039).use((0, $199830a05f92d3d0$export$2e2bcd8739ae039)).init({
    debug: false,
    initImmediate: true,
    supportedLngs: [
        'en-US',
        'ja-JP',
        'zh-CN',
        'en',
        'zh',
        'ja'
    ],
    fallbackLng: [
        'en',
        'zh',
        'ja'
    ],
    resources: {
        'en-US': {
            translation: (0, (/*@__PURE__*/$parcel$interopDefault($387689dfeea53611$exports)))
        },
        en: {
            translation: (0, (/*@__PURE__*/$parcel$interopDefault($387689dfeea53611$exports)))
        },
        'zh-CN': {
            translation: (0, (/*@__PURE__*/$parcel$interopDefault($8c66835f5815577e$exports)))
        },
        zh: {
            translation: (0, (/*@__PURE__*/$parcel$interopDefault($8c66835f5815577e$exports)))
        },
        'ja-JP': {
            translation: (0, (/*@__PURE__*/$parcel$interopDefault($654fabe7e18f36ba$exports)))
        },
        ja: {
            translation: (0, (/*@__PURE__*/$parcel$interopDefault($654fabe7e18f36ba$exports)))
        }
    }
});
window.i18next = (0, $1bac384020b50752$export$2e2bcd8739ae039 // 使i18next在全局作用域中可用 
);


// packages/alpinejs/src/scheduler.js
var $8c83eaf28779ff46$var$flushPending = false;
var $8c83eaf28779ff46$var$flushing = false;
var $8c83eaf28779ff46$var$queue = [];
var $8c83eaf28779ff46$var$lastFlushedIndex = -1;
function $8c83eaf28779ff46$var$scheduler(callback) {
    $8c83eaf28779ff46$var$queueJob(callback);
}
function $8c83eaf28779ff46$var$queueJob(job) {
    if (!$8c83eaf28779ff46$var$queue.includes(job)) $8c83eaf28779ff46$var$queue.push(job);
    $8c83eaf28779ff46$var$queueFlush();
}
function $8c83eaf28779ff46$var$dequeueJob(job) {
    let index = $8c83eaf28779ff46$var$queue.indexOf(job);
    if (index !== -1 && index > $8c83eaf28779ff46$var$lastFlushedIndex) $8c83eaf28779ff46$var$queue.splice(index, 1);
}
function $8c83eaf28779ff46$var$queueFlush() {
    if (!$8c83eaf28779ff46$var$flushing && !$8c83eaf28779ff46$var$flushPending) {
        $8c83eaf28779ff46$var$flushPending = true;
        queueMicrotask($8c83eaf28779ff46$var$flushJobs);
    }
}
function $8c83eaf28779ff46$var$flushJobs() {
    $8c83eaf28779ff46$var$flushPending = false;
    $8c83eaf28779ff46$var$flushing = true;
    for(let i = 0; i < $8c83eaf28779ff46$var$queue.length; i++){
        $8c83eaf28779ff46$var$queue[i]();
        $8c83eaf28779ff46$var$lastFlushedIndex = i;
    }
    $8c83eaf28779ff46$var$queue.length = 0;
    $8c83eaf28779ff46$var$lastFlushedIndex = -1;
    $8c83eaf28779ff46$var$flushing = false;
}
// packages/alpinejs/src/reactivity.js
var $8c83eaf28779ff46$var$reactive;
var $8c83eaf28779ff46$var$effect;
var $8c83eaf28779ff46$var$release;
var $8c83eaf28779ff46$var$raw;
var $8c83eaf28779ff46$var$shouldSchedule = true;
function $8c83eaf28779ff46$var$disableEffectScheduling(callback) {
    $8c83eaf28779ff46$var$shouldSchedule = false;
    callback();
    $8c83eaf28779ff46$var$shouldSchedule = true;
}
function $8c83eaf28779ff46$var$setReactivityEngine(engine) {
    $8c83eaf28779ff46$var$reactive = engine.reactive;
    $8c83eaf28779ff46$var$release = engine.release;
    $8c83eaf28779ff46$var$effect = (callback)=>engine.effect(callback, {
            scheduler: (task)=>{
                if ($8c83eaf28779ff46$var$shouldSchedule) $8c83eaf28779ff46$var$scheduler(task);
                else task();
            }
        });
    $8c83eaf28779ff46$var$raw = engine.raw;
}
function $8c83eaf28779ff46$var$overrideEffect(override) {
    $8c83eaf28779ff46$var$effect = override;
}
function $8c83eaf28779ff46$var$elementBoundEffect(el) {
    let cleanup2 = ()=>{};
    let wrappedEffect = (callback)=>{
        let effectReference = $8c83eaf28779ff46$var$effect(callback);
        if (!el._x_effects) {
            el._x_effects = /* @__PURE__ */ new Set();
            el._x_runEffects = ()=>{
                el._x_effects.forEach((i)=>i());
            };
        }
        el._x_effects.add(effectReference);
        cleanup2 = ()=>{
            if (effectReference === void 0) return;
            el._x_effects.delete(effectReference);
            $8c83eaf28779ff46$var$release(effectReference);
        };
        return effectReference;
    };
    return [
        wrappedEffect,
        ()=>{
            cleanup2();
        }
    ];
}
function $8c83eaf28779ff46$var$watch(getter, callback) {
    let firstTime = true;
    let oldValue;
    let effectReference = $8c83eaf28779ff46$var$effect(()=>{
        let value = getter();
        JSON.stringify(value);
        if (!firstTime) queueMicrotask(()=>{
            callback(value, oldValue);
            oldValue = value;
        });
        else oldValue = value;
        firstTime = false;
    });
    return ()=>$8c83eaf28779ff46$var$release(effectReference);
}
// packages/alpinejs/src/mutation.js
var $8c83eaf28779ff46$var$onAttributeAddeds = [];
var $8c83eaf28779ff46$var$onElRemoveds = [];
var $8c83eaf28779ff46$var$onElAddeds = [];
function $8c83eaf28779ff46$var$onElAdded(callback) {
    $8c83eaf28779ff46$var$onElAddeds.push(callback);
}
function $8c83eaf28779ff46$var$onElRemoved(el, callback) {
    if (typeof callback === "function") {
        if (!el._x_cleanups) el._x_cleanups = [];
        el._x_cleanups.push(callback);
    } else {
        callback = el;
        $8c83eaf28779ff46$var$onElRemoveds.push(callback);
    }
}
function $8c83eaf28779ff46$var$onAttributesAdded(callback) {
    $8c83eaf28779ff46$var$onAttributeAddeds.push(callback);
}
function $8c83eaf28779ff46$var$onAttributeRemoved(el, name, callback) {
    if (!el._x_attributeCleanups) el._x_attributeCleanups = {};
    if (!el._x_attributeCleanups[name]) el._x_attributeCleanups[name] = [];
    el._x_attributeCleanups[name].push(callback);
}
function $8c83eaf28779ff46$var$cleanupAttributes(el, names) {
    if (!el._x_attributeCleanups) return;
    Object.entries(el._x_attributeCleanups).forEach(([name, value])=>{
        if (names === void 0 || names.includes(name)) {
            value.forEach((i)=>i());
            delete el._x_attributeCleanups[name];
        }
    });
}
function $8c83eaf28779ff46$var$cleanupElement(el) {
    el._x_effects?.forEach($8c83eaf28779ff46$var$dequeueJob);
    while(el._x_cleanups?.length)el._x_cleanups.pop()();
}
var $8c83eaf28779ff46$var$observer = new MutationObserver($8c83eaf28779ff46$var$onMutate);
var $8c83eaf28779ff46$var$currentlyObserving = false;
function $8c83eaf28779ff46$var$startObservingMutations() {
    $8c83eaf28779ff46$var$observer.observe(document, {
        subtree: true,
        childList: true,
        attributes: true,
        attributeOldValue: true
    });
    $8c83eaf28779ff46$var$currentlyObserving = true;
}
function $8c83eaf28779ff46$var$stopObservingMutations() {
    $8c83eaf28779ff46$var$flushObserver();
    $8c83eaf28779ff46$var$observer.disconnect();
    $8c83eaf28779ff46$var$currentlyObserving = false;
}
var $8c83eaf28779ff46$var$queuedMutations = [];
function $8c83eaf28779ff46$var$flushObserver() {
    let records = $8c83eaf28779ff46$var$observer.takeRecords();
    $8c83eaf28779ff46$var$queuedMutations.push(()=>records.length > 0 && $8c83eaf28779ff46$var$onMutate(records));
    let queueLengthWhenTriggered = $8c83eaf28779ff46$var$queuedMutations.length;
    queueMicrotask(()=>{
        if ($8c83eaf28779ff46$var$queuedMutations.length === queueLengthWhenTriggered) while($8c83eaf28779ff46$var$queuedMutations.length > 0)$8c83eaf28779ff46$var$queuedMutations.shift()();
    });
}
function $8c83eaf28779ff46$var$mutateDom(callback) {
    if (!$8c83eaf28779ff46$var$currentlyObserving) return callback();
    $8c83eaf28779ff46$var$stopObservingMutations();
    let result = callback();
    $8c83eaf28779ff46$var$startObservingMutations();
    return result;
}
var $8c83eaf28779ff46$var$isCollecting = false;
var $8c83eaf28779ff46$var$deferredMutations = [];
function $8c83eaf28779ff46$var$deferMutations() {
    $8c83eaf28779ff46$var$isCollecting = true;
}
function $8c83eaf28779ff46$var$flushAndStopDeferringMutations() {
    $8c83eaf28779ff46$var$isCollecting = false;
    $8c83eaf28779ff46$var$onMutate($8c83eaf28779ff46$var$deferredMutations);
    $8c83eaf28779ff46$var$deferredMutations = [];
}
function $8c83eaf28779ff46$var$onMutate(mutations) {
    if ($8c83eaf28779ff46$var$isCollecting) {
        $8c83eaf28779ff46$var$deferredMutations = $8c83eaf28779ff46$var$deferredMutations.concat(mutations);
        return;
    }
    let addedNodes = [];
    let removedNodes = /* @__PURE__ */ new Set();
    let addedAttributes = /* @__PURE__ */ new Map();
    let removedAttributes = /* @__PURE__ */ new Map();
    for(let i = 0; i < mutations.length; i++){
        if (mutations[i].target._x_ignoreMutationObserver) continue;
        if (mutations[i].type === "childList") {
            mutations[i].removedNodes.forEach((node)=>{
                if (node.nodeType !== 1) return;
                if (!node._x_marker) return;
                removedNodes.add(node);
            });
            mutations[i].addedNodes.forEach((node)=>{
                if (node.nodeType !== 1) return;
                if (removedNodes.has(node)) {
                    removedNodes.delete(node);
                    return;
                }
                if (node._x_marker) return;
                addedNodes.push(node);
            });
        }
        if (mutations[i].type === "attributes") {
            let el = mutations[i].target;
            let name = mutations[i].attributeName;
            let oldValue = mutations[i].oldValue;
            let add2 = ()=>{
                if (!addedAttributes.has(el)) addedAttributes.set(el, []);
                addedAttributes.get(el).push({
                    name: name,
                    value: el.getAttribute(name)
                });
            };
            let remove = ()=>{
                if (!removedAttributes.has(el)) removedAttributes.set(el, []);
                removedAttributes.get(el).push(name);
            };
            if (el.hasAttribute(name) && oldValue === null) add2();
            else if (el.hasAttribute(name)) {
                remove();
                add2();
            } else remove();
        }
    }
    removedAttributes.forEach((attrs, el)=>{
        $8c83eaf28779ff46$var$cleanupAttributes(el, attrs);
    });
    addedAttributes.forEach((attrs, el)=>{
        $8c83eaf28779ff46$var$onAttributeAddeds.forEach((i)=>i(el, attrs));
    });
    for (let node of removedNodes){
        if (addedNodes.some((i)=>i.contains(node))) continue;
        $8c83eaf28779ff46$var$onElRemoveds.forEach((i)=>i(node));
    }
    for (let node of addedNodes){
        if (!node.isConnected) continue;
        $8c83eaf28779ff46$var$onElAddeds.forEach((i)=>i(node));
    }
    addedNodes = null;
    removedNodes = null;
    addedAttributes = null;
    removedAttributes = null;
}
// packages/alpinejs/src/scope.js
function $8c83eaf28779ff46$var$scope(node) {
    return $8c83eaf28779ff46$var$mergeProxies($8c83eaf28779ff46$var$closestDataStack(node));
}
function $8c83eaf28779ff46$var$addScopeToNode(node, data2, referenceNode) {
    node._x_dataStack = [
        data2,
        ...$8c83eaf28779ff46$var$closestDataStack(referenceNode || node)
    ];
    return ()=>{
        node._x_dataStack = node._x_dataStack.filter((i)=>i !== data2);
    };
}
function $8c83eaf28779ff46$var$closestDataStack(node) {
    if (node._x_dataStack) return node._x_dataStack;
    if (typeof ShadowRoot === "function" && node instanceof ShadowRoot) return $8c83eaf28779ff46$var$closestDataStack(node.host);
    if (!node.parentNode) return [];
    return $8c83eaf28779ff46$var$closestDataStack(node.parentNode);
}
function $8c83eaf28779ff46$var$mergeProxies(objects) {
    return new Proxy({
        objects: objects
    }, $8c83eaf28779ff46$var$mergeProxyTrap);
}
var $8c83eaf28779ff46$var$mergeProxyTrap = {
    ownKeys ({ objects: objects }) {
        return Array.from(new Set(objects.flatMap((i)=>Object.keys(i))));
    },
    has ({ objects: objects }, name) {
        if (name == Symbol.unscopables) return false;
        return objects.some((obj)=>Object.prototype.hasOwnProperty.call(obj, name) || Reflect.has(obj, name));
    },
    get ({ objects: objects }, name, thisProxy) {
        if (name == "toJSON") return $8c83eaf28779ff46$var$collapseProxies;
        return Reflect.get(objects.find((obj)=>Reflect.has(obj, name)) || {}, name, thisProxy);
    },
    set ({ objects: objects }, name, value, thisProxy) {
        const target = objects.find((obj)=>Object.prototype.hasOwnProperty.call(obj, name)) || objects[objects.length - 1];
        const descriptor = Object.getOwnPropertyDescriptor(target, name);
        if (descriptor?.set && descriptor?.get) return descriptor.set.call(thisProxy, value) || true;
        return Reflect.set(target, name, value);
    }
};
function $8c83eaf28779ff46$var$collapseProxies() {
    let keys = Reflect.ownKeys(this);
    return keys.reduce((acc, key)=>{
        acc[key] = Reflect.get(this, key);
        return acc;
    }, {});
}
// packages/alpinejs/src/interceptor.js
function $8c83eaf28779ff46$var$initInterceptors(data2) {
    let isObject2 = (val)=>typeof val === "object" && !Array.isArray(val) && val !== null;
    let recurse = (obj, basePath = "")=>{
        Object.entries(Object.getOwnPropertyDescriptors(obj)).forEach(([key, { value: value, enumerable: enumerable }])=>{
            if (enumerable === false || value === void 0) return;
            if (typeof value === "object" && value !== null && value.__v_skip) return;
            let path = basePath === "" ? key : `${basePath}.${key}`;
            if (typeof value === "object" && value !== null && value._x_interceptor) obj[key] = value.initialize(data2, path, key);
            else if (isObject2(value) && value !== obj && !(value instanceof Element)) recurse(value, path);
        });
    };
    return recurse(data2);
}
function $8c83eaf28779ff46$var$interceptor(callback, mutateObj = ()=>{}) {
    let obj = {
        initialValue: void 0,
        _x_interceptor: true,
        initialize (data2, path, key) {
            return callback(this.initialValue, ()=>$8c83eaf28779ff46$var$get(data2, path), (value)=>$8c83eaf28779ff46$var$set(data2, path, value), path, key);
        }
    };
    mutateObj(obj);
    return (initialValue)=>{
        if (typeof initialValue === "object" && initialValue !== null && initialValue._x_interceptor) {
            let initialize = obj.initialize.bind(obj);
            obj.initialize = (data2, path, key)=>{
                let innerValue = initialValue.initialize(data2, path, key);
                obj.initialValue = innerValue;
                return initialize(data2, path, key);
            };
        } else obj.initialValue = initialValue;
        return obj;
    };
}
function $8c83eaf28779ff46$var$get(obj, path) {
    return path.split(".").reduce((carry, segment)=>carry[segment], obj);
}
function $8c83eaf28779ff46$var$set(obj, path, value) {
    if (typeof path === "string") path = path.split(".");
    if (path.length === 1) obj[path[0]] = value;
    else if (path.length === 0) throw error;
    else {
        if (obj[path[0]]) return $8c83eaf28779ff46$var$set(obj[path[0]], path.slice(1), value);
        else {
            obj[path[0]] = {};
            return $8c83eaf28779ff46$var$set(obj[path[0]], path.slice(1), value);
        }
    }
}
// packages/alpinejs/src/magics.js
var $8c83eaf28779ff46$var$magics = {};
function $8c83eaf28779ff46$var$magic(name, callback) {
    $8c83eaf28779ff46$var$magics[name] = callback;
}
function $8c83eaf28779ff46$var$injectMagics(obj, el) {
    let memoizedUtilities = $8c83eaf28779ff46$var$getUtilities(el);
    Object.entries($8c83eaf28779ff46$var$magics).forEach(([name, callback])=>{
        Object.defineProperty(obj, `$${name}`, {
            get () {
                return callback(el, memoizedUtilities);
            },
            enumerable: false
        });
    });
    return obj;
}
function $8c83eaf28779ff46$var$getUtilities(el) {
    let [utilities, cleanup2] = $8c83eaf28779ff46$var$getElementBoundUtilities(el);
    let utils = {
        interceptor: $8c83eaf28779ff46$var$interceptor,
        ...utilities
    };
    $8c83eaf28779ff46$var$onElRemoved(el, cleanup2);
    return utils;
}
// packages/alpinejs/src/utils/error.js
function $8c83eaf28779ff46$var$tryCatch(el, expression, callback, ...args) {
    try {
        return callback(...args);
    } catch (e) {
        $8c83eaf28779ff46$var$handleError(e, el, expression);
    }
}
function $8c83eaf28779ff46$var$handleError(...args) {
    return $8c83eaf28779ff46$var$errorHandler(...args);
}
var $8c83eaf28779ff46$var$errorHandler = $8c83eaf28779ff46$var$normalErrorHandler;
function $8c83eaf28779ff46$var$setErrorHandler(handler4) {
    $8c83eaf28779ff46$var$errorHandler = handler4;
}
function $8c83eaf28779ff46$var$normalErrorHandler(error2, el, expression) {
    error2 = Object.assign(error2 ?? {
        message: "No error message given."
    }, {
        el: el,
        expression: expression
    });
    console.warn(`Alpine Expression Error: ${error2.message}

${expression ? 'Expression: "' + expression + '"\n\n' : ""}`, el);
    setTimeout(()=>{
        throw error2;
    }, 0);
}
// packages/alpinejs/src/evaluator.js
var $8c83eaf28779ff46$var$shouldAutoEvaluateFunctions = true;
function $8c83eaf28779ff46$var$dontAutoEvaluateFunctions(callback) {
    let cache = $8c83eaf28779ff46$var$shouldAutoEvaluateFunctions;
    $8c83eaf28779ff46$var$shouldAutoEvaluateFunctions = false;
    let result = callback();
    $8c83eaf28779ff46$var$shouldAutoEvaluateFunctions = cache;
    return result;
}
function $8c83eaf28779ff46$var$evaluate(el, expression, extras = {}) {
    let result;
    $8c83eaf28779ff46$var$evaluateLater(el, expression)((value)=>result = value, extras);
    return result;
}
function $8c83eaf28779ff46$var$evaluateLater(...args) {
    return $8c83eaf28779ff46$var$theEvaluatorFunction(...args);
}
var $8c83eaf28779ff46$var$theEvaluatorFunction = $8c83eaf28779ff46$var$normalEvaluator;
function $8c83eaf28779ff46$var$setEvaluator(newEvaluator) {
    $8c83eaf28779ff46$var$theEvaluatorFunction = newEvaluator;
}
var $8c83eaf28779ff46$var$theRawEvaluatorFunction;
function $8c83eaf28779ff46$var$setRawEvaluator(newEvaluator) {
    $8c83eaf28779ff46$var$theRawEvaluatorFunction = newEvaluator;
}
function $8c83eaf28779ff46$var$normalEvaluator(el, expression) {
    let overriddenMagics = {};
    $8c83eaf28779ff46$var$injectMagics(overriddenMagics, el);
    let dataStack = [
        overriddenMagics,
        ...$8c83eaf28779ff46$var$closestDataStack(el)
    ];
    let evaluator = typeof expression === "function" ? $8c83eaf28779ff46$var$generateEvaluatorFromFunction(dataStack, expression) : $8c83eaf28779ff46$var$generateEvaluatorFromString(dataStack, expression, el);
    return $8c83eaf28779ff46$var$tryCatch.bind(null, el, expression, evaluator);
}
function $8c83eaf28779ff46$var$generateEvaluatorFromFunction(dataStack, func) {
    return (receiver = ()=>{}, { scope: scope2 = {}, params: params = [], context: context } = {})=>{
        if (!$8c83eaf28779ff46$var$shouldAutoEvaluateFunctions) {
            $8c83eaf28779ff46$var$runIfTypeOfFunction(receiver, func, $8c83eaf28779ff46$var$mergeProxies([
                scope2,
                ...dataStack
            ]), params);
            return;
        }
        let result = func.apply($8c83eaf28779ff46$var$mergeProxies([
            scope2,
            ...dataStack
        ]), params);
        $8c83eaf28779ff46$var$runIfTypeOfFunction(receiver, result);
    };
}
var $8c83eaf28779ff46$var$evaluatorMemo = {};
function $8c83eaf28779ff46$var$generateFunctionFromString(expression, el) {
    if ($8c83eaf28779ff46$var$evaluatorMemo[expression]) return $8c83eaf28779ff46$var$evaluatorMemo[expression];
    let AsyncFunction = Object.getPrototypeOf(async function() {}).constructor;
    let rightSideSafeExpression = /^[\n\s]*if.*\(.*\)/.test(expression.trim()) || /^(let|const)\s/.test(expression.trim()) ? `(async()=>{ ${expression} })()` : expression;
    const safeAsyncFunction = ()=>{
        try {
            let func2 = new AsyncFunction([
                "__self",
                "scope"
            ], `with (scope) { __self.result = ${rightSideSafeExpression} }; __self.finished = true; return __self.result;`);
            Object.defineProperty(func2, "name", {
                value: `[Alpine] ${expression}`
            });
            return func2;
        } catch (error2) {
            $8c83eaf28779ff46$var$handleError(error2, el, expression);
            return Promise.resolve();
        }
    };
    let func = safeAsyncFunction();
    $8c83eaf28779ff46$var$evaluatorMemo[expression] = func;
    return func;
}
function $8c83eaf28779ff46$var$generateEvaluatorFromString(dataStack, expression, el) {
    let func = $8c83eaf28779ff46$var$generateFunctionFromString(expression, el);
    return (receiver = ()=>{}, { scope: scope2 = {}, params: params = [], context: context } = {})=>{
        func.result = void 0;
        func.finished = false;
        let completeScope = $8c83eaf28779ff46$var$mergeProxies([
            scope2,
            ...dataStack
        ]);
        if (typeof func === "function") {
            let promise = func.call(context, func, completeScope).catch((error2)=>$8c83eaf28779ff46$var$handleError(error2, el, expression));
            if (func.finished) {
                $8c83eaf28779ff46$var$runIfTypeOfFunction(receiver, func.result, completeScope, params, el);
                func.result = void 0;
            } else promise.then((result)=>{
                $8c83eaf28779ff46$var$runIfTypeOfFunction(receiver, result, completeScope, params, el);
            }).catch((error2)=>$8c83eaf28779ff46$var$handleError(error2, el, expression)).finally(()=>func.result = void 0);
        }
    };
}
function $8c83eaf28779ff46$var$runIfTypeOfFunction(receiver, value, scope2, params, el) {
    if ($8c83eaf28779ff46$var$shouldAutoEvaluateFunctions && typeof value === "function") {
        let result = value.apply(scope2, params);
        if (result instanceof Promise) result.then((i)=>$8c83eaf28779ff46$var$runIfTypeOfFunction(receiver, i, scope2, params)).catch((error2)=>$8c83eaf28779ff46$var$handleError(error2, el, value));
        else receiver(result);
    } else if (typeof value === "object" && value instanceof Promise) value.then((i)=>receiver(i));
    else receiver(value);
}
function $8c83eaf28779ff46$var$evaluateRaw(...args) {
    return $8c83eaf28779ff46$var$theRawEvaluatorFunction(...args);
}
function $8c83eaf28779ff46$var$normalRawEvaluator(el, expression, extras = {}) {
    let overriddenMagics = {};
    $8c83eaf28779ff46$var$injectMagics(overriddenMagics, el);
    let dataStack = [
        overriddenMagics,
        ...$8c83eaf28779ff46$var$closestDataStack(el)
    ];
    let scope2 = $8c83eaf28779ff46$var$mergeProxies([
        extras.scope ?? {},
        ...dataStack
    ]);
    let params = extras.params ?? [];
    if (expression.includes("await")) {
        let AsyncFunction = Object.getPrototypeOf(async function() {}).constructor;
        let rightSideSafeExpression = /^[\n\s]*if.*\(.*\)/.test(expression.trim()) || /^(let|const)\s/.test(expression.trim()) ? `(async()=>{ ${expression} })()` : expression;
        let func = new AsyncFunction([
            "scope"
        ], `with (scope) { let __result = ${rightSideSafeExpression}; return __result }`);
        let result = func.call(extras.context, scope2);
        return result;
    } else {
        let rightSideSafeExpression = /^[\n\s]*if.*\(.*\)/.test(expression.trim()) || /^(let|const)\s/.test(expression.trim()) ? `(()=>{ ${expression} })()` : expression;
        let func = new Function([
            "scope"
        ], `with (scope) { let __result = ${rightSideSafeExpression}; return __result }`);
        let result = func.call(extras.context, scope2);
        if (typeof result === "function" && $8c83eaf28779ff46$var$shouldAutoEvaluateFunctions) return result.apply(scope2, params);
        return result;
    }
}
// packages/alpinejs/src/directives.js
var $8c83eaf28779ff46$var$prefixAsString = "x-";
function $8c83eaf28779ff46$var$prefix(subject = "") {
    return $8c83eaf28779ff46$var$prefixAsString + subject;
}
function $8c83eaf28779ff46$var$setPrefix(newPrefix) {
    $8c83eaf28779ff46$var$prefixAsString = newPrefix;
}
var $8c83eaf28779ff46$var$directiveHandlers = {};
function $8c83eaf28779ff46$var$directive(name, callback) {
    $8c83eaf28779ff46$var$directiveHandlers[name] = callback;
    return {
        before (directive2) {
            if (!$8c83eaf28779ff46$var$directiveHandlers[directive2]) {
                console.warn(String.raw`Cannot find directive \`${directive2}\`. \`${name}\` will use the default order of execution`);
                return;
            }
            const pos = $8c83eaf28779ff46$var$directiveOrder.indexOf(directive2);
            $8c83eaf28779ff46$var$directiveOrder.splice(pos >= 0 ? pos : $8c83eaf28779ff46$var$directiveOrder.indexOf("DEFAULT"), 0, name);
        }
    };
}
function $8c83eaf28779ff46$var$directiveExists(name) {
    return Object.keys($8c83eaf28779ff46$var$directiveHandlers).includes(name);
}
function $8c83eaf28779ff46$var$directives(el, attributes, originalAttributeOverride) {
    attributes = Array.from(attributes);
    if (el._x_virtualDirectives) {
        let vAttributes = Object.entries(el._x_virtualDirectives).map(([name, value])=>({
                name: name,
                value: value
            }));
        let staticAttributes = $8c83eaf28779ff46$var$attributesOnly(vAttributes);
        vAttributes = vAttributes.map((attribute)=>{
            if (staticAttributes.find((attr)=>attr.name === attribute.name)) return {
                name: `x-bind:${attribute.name}`,
                value: `"${attribute.value}"`
            };
            return attribute;
        });
        attributes = attributes.concat(vAttributes);
    }
    let transformedAttributeMap = {};
    let directives2 = attributes.map($8c83eaf28779ff46$var$toTransformedAttributes((newName, oldName)=>transformedAttributeMap[newName] = oldName)).filter($8c83eaf28779ff46$var$outNonAlpineAttributes).map($8c83eaf28779ff46$var$toParsedDirectives(transformedAttributeMap, originalAttributeOverride)).sort($8c83eaf28779ff46$var$byPriority);
    return directives2.map((directive2)=>{
        return $8c83eaf28779ff46$var$getDirectiveHandler(el, directive2);
    });
}
function $8c83eaf28779ff46$var$attributesOnly(attributes) {
    return Array.from(attributes).map($8c83eaf28779ff46$var$toTransformedAttributes()).filter((attr)=>!$8c83eaf28779ff46$var$outNonAlpineAttributes(attr));
}
var $8c83eaf28779ff46$var$isDeferringHandlers = false;
var $8c83eaf28779ff46$var$directiveHandlerStacks = /* @__PURE__ */ new Map();
var $8c83eaf28779ff46$var$currentHandlerStackKey = Symbol();
function $8c83eaf28779ff46$var$deferHandlingDirectives(callback) {
    $8c83eaf28779ff46$var$isDeferringHandlers = true;
    let key = Symbol();
    $8c83eaf28779ff46$var$currentHandlerStackKey = key;
    $8c83eaf28779ff46$var$directiveHandlerStacks.set(key, []);
    let flushHandlers = ()=>{
        while($8c83eaf28779ff46$var$directiveHandlerStacks.get(key).length)$8c83eaf28779ff46$var$directiveHandlerStacks.get(key).shift()();
        $8c83eaf28779ff46$var$directiveHandlerStacks.delete(key);
    };
    let stopDeferring = ()=>{
        $8c83eaf28779ff46$var$isDeferringHandlers = false;
        flushHandlers();
    };
    callback(flushHandlers);
    stopDeferring();
}
function $8c83eaf28779ff46$var$getElementBoundUtilities(el) {
    let cleanups = [];
    let cleanup2 = (callback)=>cleanups.push(callback);
    let [effect3, cleanupEffect] = $8c83eaf28779ff46$var$elementBoundEffect(el);
    cleanups.push(cleanupEffect);
    let utilities = {
        Alpine: $8c83eaf28779ff46$var$alpine_default,
        effect: effect3,
        cleanup: cleanup2,
        evaluateLater: $8c83eaf28779ff46$var$evaluateLater.bind($8c83eaf28779ff46$var$evaluateLater, el),
        evaluate: $8c83eaf28779ff46$var$evaluate.bind($8c83eaf28779ff46$var$evaluate, el)
    };
    let doCleanup = ()=>cleanups.forEach((i)=>i());
    return [
        utilities,
        doCleanup
    ];
}
function $8c83eaf28779ff46$var$getDirectiveHandler(el, directive2) {
    let noop = ()=>{};
    let handler4 = $8c83eaf28779ff46$var$directiveHandlers[directive2.type] || noop;
    let [utilities, cleanup2] = $8c83eaf28779ff46$var$getElementBoundUtilities(el);
    $8c83eaf28779ff46$var$onAttributeRemoved(el, directive2.original, cleanup2);
    let fullHandler = ()=>{
        if (el._x_ignore || el._x_ignoreSelf) return;
        handler4.inline && handler4.inline(el, directive2, utilities);
        handler4 = handler4.bind(handler4, el, directive2, utilities);
        $8c83eaf28779ff46$var$isDeferringHandlers ? $8c83eaf28779ff46$var$directiveHandlerStacks.get($8c83eaf28779ff46$var$currentHandlerStackKey).push(handler4) : handler4();
    };
    fullHandler.runCleanups = cleanup2;
    return fullHandler;
}
var $8c83eaf28779ff46$var$startingWith = (subject, replacement)=>({ name: name, value: value })=>{
        if (name.startsWith(subject)) name = name.replace(subject, replacement);
        return {
            name: name,
            value: value
        };
    };
var $8c83eaf28779ff46$var$into = (i)=>i;
function $8c83eaf28779ff46$var$toTransformedAttributes(callback = ()=>{}) {
    return ({ name: name, value: value })=>{
        let { name: newName, value: newValue } = $8c83eaf28779ff46$var$attributeTransformers.reduce((carry, transform)=>{
            return transform(carry);
        }, {
            name: name,
            value: value
        });
        if (newName !== name) callback(newName, name);
        return {
            name: newName,
            value: newValue
        };
    };
}
var $8c83eaf28779ff46$var$attributeTransformers = [];
function $8c83eaf28779ff46$var$mapAttributes(callback) {
    $8c83eaf28779ff46$var$attributeTransformers.push(callback);
}
function $8c83eaf28779ff46$var$outNonAlpineAttributes({ name: name }) {
    return $8c83eaf28779ff46$var$alpineAttributeRegex().test(name);
}
var $8c83eaf28779ff46$var$alpineAttributeRegex = ()=>new RegExp(`^${$8c83eaf28779ff46$var$prefixAsString}([^:^.]+)\\b`);
function $8c83eaf28779ff46$var$toParsedDirectives(transformedAttributeMap, originalAttributeOverride) {
    return ({ name: name, value: value })=>{
        if (name === value) value = "";
        let typeMatch = name.match($8c83eaf28779ff46$var$alpineAttributeRegex());
        let valueMatch = name.match(/:([a-zA-Z0-9\-_:]+)/);
        let modifiers = name.match(/\.[^.\]]+(?=[^\]]*$)/g) || [];
        let original = originalAttributeOverride || transformedAttributeMap[name] || name;
        return {
            type: typeMatch ? typeMatch[1] : null,
            value: valueMatch ? valueMatch[1] : null,
            modifiers: modifiers.map((i)=>i.replace(".", "")),
            expression: value,
            original: original
        };
    };
}
var $8c83eaf28779ff46$var$DEFAULT = "DEFAULT";
var $8c83eaf28779ff46$var$directiveOrder = [
    "ignore",
    "ref",
    "data",
    "id",
    "anchor",
    "bind",
    "init",
    "for",
    "model",
    "modelable",
    "transition",
    "show",
    "if",
    $8c83eaf28779ff46$var$DEFAULT,
    "teleport"
];
function $8c83eaf28779ff46$var$byPriority(a, b) {
    let typeA = $8c83eaf28779ff46$var$directiveOrder.indexOf(a.type) === -1 ? $8c83eaf28779ff46$var$DEFAULT : a.type;
    let typeB = $8c83eaf28779ff46$var$directiveOrder.indexOf(b.type) === -1 ? $8c83eaf28779ff46$var$DEFAULT : b.type;
    return $8c83eaf28779ff46$var$directiveOrder.indexOf(typeA) - $8c83eaf28779ff46$var$directiveOrder.indexOf(typeB);
}
// packages/alpinejs/src/utils/dispatch.js
function $8c83eaf28779ff46$var$dispatch(el, name, detail = {}) {
    el.dispatchEvent(new CustomEvent(name, {
        detail: detail,
        bubbles: true,
        // Allows events to pass the shadow DOM barrier.
        composed: true,
        cancelable: true
    }));
}
// packages/alpinejs/src/utils/walk.js
function $8c83eaf28779ff46$var$walk(el, callback) {
    if (typeof ShadowRoot === "function" && el instanceof ShadowRoot) {
        Array.from(el.children).forEach((el2)=>$8c83eaf28779ff46$var$walk(el2, callback));
        return;
    }
    let skip = false;
    callback(el, ()=>skip = true);
    if (skip) return;
    let node = el.firstElementChild;
    while(node){
        $8c83eaf28779ff46$var$walk(node, callback, false);
        node = node.nextElementSibling;
    }
}
// packages/alpinejs/src/utils/warn.js
function $8c83eaf28779ff46$var$warn(message, ...args) {
    console.warn(`Alpine Warning: ${message}`, ...args);
}
// packages/alpinejs/src/lifecycle.js
var $8c83eaf28779ff46$var$started = false;
function $8c83eaf28779ff46$var$start() {
    if ($8c83eaf28779ff46$var$started) $8c83eaf28779ff46$var$warn("Alpine has already been initialized on this page. Calling Alpine.start() more than once can cause problems.");
    $8c83eaf28779ff46$var$started = true;
    if (!document.body) $8c83eaf28779ff46$var$warn("Unable to initialize. Trying to load Alpine before `<body>` is available. Did you forget to add `defer` in Alpine's `<script>` tag?");
    $8c83eaf28779ff46$var$dispatch(document, "alpine:init");
    $8c83eaf28779ff46$var$dispatch(document, "alpine:initializing");
    $8c83eaf28779ff46$var$startObservingMutations();
    $8c83eaf28779ff46$var$onElAdded((el)=>$8c83eaf28779ff46$var$initTree(el, $8c83eaf28779ff46$var$walk));
    $8c83eaf28779ff46$var$onElRemoved((el)=>$8c83eaf28779ff46$var$destroyTree(el));
    $8c83eaf28779ff46$var$onAttributesAdded((el, attrs)=>{
        $8c83eaf28779ff46$var$directives(el, attrs).forEach((handle)=>handle());
    });
    let outNestedComponents = (el)=>!$8c83eaf28779ff46$var$closestRoot(el.parentElement, true);
    Array.from(document.querySelectorAll($8c83eaf28779ff46$var$allSelectors().join(","))).filter(outNestedComponents).forEach((el)=>{
        $8c83eaf28779ff46$var$initTree(el);
    });
    $8c83eaf28779ff46$var$dispatch(document, "alpine:initialized");
    setTimeout(()=>{
        $8c83eaf28779ff46$var$warnAboutMissingPlugins();
    });
}
var $8c83eaf28779ff46$var$rootSelectorCallbacks = [];
var $8c83eaf28779ff46$var$initSelectorCallbacks = [];
function $8c83eaf28779ff46$var$rootSelectors() {
    return $8c83eaf28779ff46$var$rootSelectorCallbacks.map((fn)=>fn());
}
function $8c83eaf28779ff46$var$allSelectors() {
    return $8c83eaf28779ff46$var$rootSelectorCallbacks.concat($8c83eaf28779ff46$var$initSelectorCallbacks).map((fn)=>fn());
}
function $8c83eaf28779ff46$var$addRootSelector(selectorCallback) {
    $8c83eaf28779ff46$var$rootSelectorCallbacks.push(selectorCallback);
}
function $8c83eaf28779ff46$var$addInitSelector(selectorCallback) {
    $8c83eaf28779ff46$var$initSelectorCallbacks.push(selectorCallback);
}
function $8c83eaf28779ff46$var$closestRoot(el, includeInitSelectors = false) {
    return $8c83eaf28779ff46$var$findClosest(el, (element)=>{
        const selectors = includeInitSelectors ? $8c83eaf28779ff46$var$allSelectors() : $8c83eaf28779ff46$var$rootSelectors();
        if (selectors.some((selector)=>element.matches(selector))) return true;
    });
}
function $8c83eaf28779ff46$var$findClosest(el, callback) {
    if (!el) return;
    if (callback(el)) return el;
    if (el._x_teleportBack) el = el._x_teleportBack;
    if (el.parentNode instanceof ShadowRoot) return $8c83eaf28779ff46$var$findClosest(el.parentNode.host, callback);
    if (!el.parentElement) return;
    return $8c83eaf28779ff46$var$findClosest(el.parentElement, callback);
}
function $8c83eaf28779ff46$var$isRoot(el) {
    return $8c83eaf28779ff46$var$rootSelectors().some((selector)=>el.matches(selector));
}
var $8c83eaf28779ff46$var$initInterceptors2 = [];
function $8c83eaf28779ff46$var$interceptInit(callback) {
    $8c83eaf28779ff46$var$initInterceptors2.push(callback);
}
var $8c83eaf28779ff46$var$markerDispenser = 1;
function $8c83eaf28779ff46$var$initTree(el, walker = $8c83eaf28779ff46$var$walk, intercept = ()=>{}) {
    if ($8c83eaf28779ff46$var$findClosest(el, (i)=>i._x_ignore)) return;
    $8c83eaf28779ff46$var$deferHandlingDirectives(()=>{
        walker(el, (el2, skip)=>{
            if (el2._x_marker) return;
            intercept(el2, skip);
            $8c83eaf28779ff46$var$initInterceptors2.forEach((i)=>i(el2, skip));
            $8c83eaf28779ff46$var$directives(el2, el2.attributes).forEach((handle)=>handle());
            if (!el2._x_ignore) el2._x_marker = $8c83eaf28779ff46$var$markerDispenser++;
            el2._x_ignore && skip();
        });
    });
}
function $8c83eaf28779ff46$var$destroyTree(root, walker = $8c83eaf28779ff46$var$walk) {
    walker(root, (el)=>{
        $8c83eaf28779ff46$var$cleanupElement(el);
        $8c83eaf28779ff46$var$cleanupAttributes(el);
        delete el._x_marker;
    });
}
function $8c83eaf28779ff46$var$warnAboutMissingPlugins() {
    let pluginDirectives = [
        [
            "ui",
            "dialog",
            [
                "[x-dialog], [x-popover]"
            ]
        ],
        [
            "anchor",
            "anchor",
            [
                "[x-anchor]"
            ]
        ],
        [
            "sort",
            "sort",
            [
                "[x-sort]"
            ]
        ]
    ];
    pluginDirectives.forEach(([plugin2, directive2, selectors])=>{
        if ($8c83eaf28779ff46$var$directiveExists(directive2)) return;
        selectors.some((selector)=>{
            if (document.querySelector(selector)) {
                $8c83eaf28779ff46$var$warn(`found "${selector}", but missing ${plugin2} plugin`);
                return true;
            }
        });
    });
}
// packages/alpinejs/src/nextTick.js
var $8c83eaf28779ff46$var$tickStack = [];
var $8c83eaf28779ff46$var$isHolding = false;
function $8c83eaf28779ff46$var$nextTick(callback = ()=>{}) {
    queueMicrotask(()=>{
        $8c83eaf28779ff46$var$isHolding || setTimeout(()=>{
            $8c83eaf28779ff46$var$releaseNextTicks();
        });
    });
    return new Promise((res)=>{
        $8c83eaf28779ff46$var$tickStack.push(()=>{
            callback();
            res();
        });
    });
}
function $8c83eaf28779ff46$var$releaseNextTicks() {
    $8c83eaf28779ff46$var$isHolding = false;
    while($8c83eaf28779ff46$var$tickStack.length)$8c83eaf28779ff46$var$tickStack.shift()();
}
function $8c83eaf28779ff46$var$holdNextTicks() {
    $8c83eaf28779ff46$var$isHolding = true;
}
// packages/alpinejs/src/utils/classes.js
function $8c83eaf28779ff46$var$setClasses(el, value) {
    if (Array.isArray(value)) return $8c83eaf28779ff46$var$setClassesFromString(el, value.join(" "));
    else if (typeof value === "object" && value !== null) return $8c83eaf28779ff46$var$setClassesFromObject(el, value);
    else if (typeof value === "function") return $8c83eaf28779ff46$var$setClasses(el, value());
    return $8c83eaf28779ff46$var$setClassesFromString(el, value);
}
function $8c83eaf28779ff46$var$setClassesFromString(el, classString) {
    let split = (classString2)=>classString2.split(" ").filter(Boolean);
    let missingClasses = (classString2)=>classString2.split(" ").filter((i)=>!el.classList.contains(i)).filter(Boolean);
    let addClassesAndReturnUndo = (classes)=>{
        el.classList.add(...classes);
        return ()=>{
            el.classList.remove(...classes);
        };
    };
    classString = classString === true ? classString = "" : classString || "";
    return addClassesAndReturnUndo(missingClasses(classString));
}
function $8c83eaf28779ff46$var$setClassesFromObject(el, classObject) {
    let split = (classString)=>classString.split(" ").filter(Boolean);
    let forAdd = Object.entries(classObject).flatMap(([classString, bool])=>bool ? split(classString) : false).filter(Boolean);
    let forRemove = Object.entries(classObject).flatMap(([classString, bool])=>!bool ? split(classString) : false).filter(Boolean);
    let added = [];
    let removed = [];
    forRemove.forEach((i)=>{
        if (el.classList.contains(i)) {
            el.classList.remove(i);
            removed.push(i);
        }
    });
    forAdd.forEach((i)=>{
        if (!el.classList.contains(i)) {
            el.classList.add(i);
            added.push(i);
        }
    });
    return ()=>{
        removed.forEach((i)=>el.classList.add(i));
        added.forEach((i)=>el.classList.remove(i));
    };
}
// packages/alpinejs/src/utils/styles.js
function $8c83eaf28779ff46$var$setStyles(el, value) {
    if (typeof value === "object" && value !== null) return $8c83eaf28779ff46$var$setStylesFromObject(el, value);
    return $8c83eaf28779ff46$var$setStylesFromString(el, value);
}
function $8c83eaf28779ff46$var$setStylesFromObject(el, value) {
    let previousStyles = {};
    Object.entries(value).forEach(([key, value2])=>{
        previousStyles[key] = el.style[key];
        if (!key.startsWith("--")) key = $8c83eaf28779ff46$var$kebabCase(key);
        el.style.setProperty(key, value2);
    });
    setTimeout(()=>{
        if (el.style.length === 0) el.removeAttribute("style");
    });
    return ()=>{
        $8c83eaf28779ff46$var$setStyles(el, previousStyles);
    };
}
function $8c83eaf28779ff46$var$setStylesFromString(el, value) {
    let cache = el.getAttribute("style", value);
    el.setAttribute("style", value);
    return ()=>{
        el.setAttribute("style", cache || "");
    };
}
function $8c83eaf28779ff46$var$kebabCase(subject) {
    return subject.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
}
// packages/alpinejs/src/utils/once.js
function $8c83eaf28779ff46$var$once(callback, fallback = ()=>{}) {
    let called = false;
    return function() {
        if (!called) {
            called = true;
            callback.apply(this, arguments);
        } else fallback.apply(this, arguments);
    };
}
// packages/alpinejs/src/directives/x-transition.js
$8c83eaf28779ff46$var$directive("transition", (el, { value: value, modifiers: modifiers, expression: expression }, { evaluate: evaluate2 })=>{
    if (typeof expression === "function") expression = evaluate2(expression);
    if (expression === false) return;
    if (!expression || typeof expression === "boolean") $8c83eaf28779ff46$var$registerTransitionsFromHelper(el, modifiers, value);
    else $8c83eaf28779ff46$var$registerTransitionsFromClassString(el, expression, value);
});
function $8c83eaf28779ff46$var$registerTransitionsFromClassString(el, classString, stage) {
    $8c83eaf28779ff46$var$registerTransitionObject(el, $8c83eaf28779ff46$var$setClasses, "");
    let directiveStorageMap = {
        "enter": (classes)=>{
            el._x_transition.enter.during = classes;
        },
        "enter-start": (classes)=>{
            el._x_transition.enter.start = classes;
        },
        "enter-end": (classes)=>{
            el._x_transition.enter.end = classes;
        },
        "leave": (classes)=>{
            el._x_transition.leave.during = classes;
        },
        "leave-start": (classes)=>{
            el._x_transition.leave.start = classes;
        },
        "leave-end": (classes)=>{
            el._x_transition.leave.end = classes;
        }
    };
    directiveStorageMap[stage](classString);
}
function $8c83eaf28779ff46$var$registerTransitionsFromHelper(el, modifiers, stage) {
    $8c83eaf28779ff46$var$registerTransitionObject(el, $8c83eaf28779ff46$var$setStyles);
    let doesntSpecify = !modifiers.includes("in") && !modifiers.includes("out") && !stage;
    let transitioningIn = doesntSpecify || modifiers.includes("in") || [
        "enter"
    ].includes(stage);
    let transitioningOut = doesntSpecify || modifiers.includes("out") || [
        "leave"
    ].includes(stage);
    if (modifiers.includes("in") && !doesntSpecify) modifiers = modifiers.filter((i, index)=>index < modifiers.indexOf("out"));
    if (modifiers.includes("out") && !doesntSpecify) modifiers = modifiers.filter((i, index)=>index > modifiers.indexOf("out"));
    let wantsAll = !modifiers.includes("opacity") && !modifiers.includes("scale");
    let wantsOpacity = wantsAll || modifiers.includes("opacity");
    let wantsScale = wantsAll || modifiers.includes("scale");
    let opacityValue = wantsOpacity ? 0 : 1;
    let scaleValue = wantsScale ? $8c83eaf28779ff46$var$modifierValue(modifiers, "scale", 95) / 100 : 1;
    let delay = $8c83eaf28779ff46$var$modifierValue(modifiers, "delay", 0) / 1e3;
    let origin = $8c83eaf28779ff46$var$modifierValue(modifiers, "origin", "center");
    let property = "opacity, transform";
    let durationIn = $8c83eaf28779ff46$var$modifierValue(modifiers, "duration", 150) / 1e3;
    let durationOut = $8c83eaf28779ff46$var$modifierValue(modifiers, "duration", 75) / 1e3;
    let easing = `cubic-bezier(0.4, 0.0, 0.2, 1)`;
    if (transitioningIn) {
        el._x_transition.enter.during = {
            transformOrigin: origin,
            transitionDelay: `${delay}s`,
            transitionProperty: property,
            transitionDuration: `${durationIn}s`,
            transitionTimingFunction: easing
        };
        el._x_transition.enter.start = {
            opacity: opacityValue,
            transform: `scale(${scaleValue})`
        };
        el._x_transition.enter.end = {
            opacity: 1,
            transform: `scale(1)`
        };
    }
    if (transitioningOut) {
        el._x_transition.leave.during = {
            transformOrigin: origin,
            transitionDelay: `${delay}s`,
            transitionProperty: property,
            transitionDuration: `${durationOut}s`,
            transitionTimingFunction: easing
        };
        el._x_transition.leave.start = {
            opacity: 1,
            transform: `scale(1)`
        };
        el._x_transition.leave.end = {
            opacity: opacityValue,
            transform: `scale(${scaleValue})`
        };
    }
}
function $8c83eaf28779ff46$var$registerTransitionObject(el, setFunction, defaultValue = {}) {
    if (!el._x_transition) el._x_transition = {
        enter: {
            during: defaultValue,
            start: defaultValue,
            end: defaultValue
        },
        leave: {
            during: defaultValue,
            start: defaultValue,
            end: defaultValue
        },
        in (before = ()=>{}, after = ()=>{}) {
            $8c83eaf28779ff46$var$transition(el, setFunction, {
                during: this.enter.during,
                start: this.enter.start,
                end: this.enter.end
            }, before, after);
        },
        out (before = ()=>{}, after = ()=>{}) {
            $8c83eaf28779ff46$var$transition(el, setFunction, {
                during: this.leave.during,
                start: this.leave.start,
                end: this.leave.end
            }, before, after);
        }
    };
}
window.Element.prototype._x_toggleAndCascadeWithTransitions = function(el, value, show, hide) {
    const nextTick2 = document.visibilityState === "visible" ? requestAnimationFrame : setTimeout;
    let clickAwayCompatibleShow = ()=>nextTick2(show);
    if (value) {
        if (el._x_transition && (el._x_transition.enter || el._x_transition.leave)) el._x_transition.enter && (Object.entries(el._x_transition.enter.during).length || Object.entries(el._x_transition.enter.start).length || Object.entries(el._x_transition.enter.end).length) ? el._x_transition.in(show) : clickAwayCompatibleShow();
        else el._x_transition ? el._x_transition.in(show) : clickAwayCompatibleShow();
        return;
    }
    el._x_hidePromise = el._x_transition ? new Promise((resolve, reject)=>{
        el._x_transition.out(()=>{}, ()=>resolve(hide));
        el._x_transitioning && el._x_transitioning.beforeCancel(()=>reject({
                isFromCancelledTransition: true
            }));
    }) : Promise.resolve(hide);
    queueMicrotask(()=>{
        let closest = $8c83eaf28779ff46$var$closestHide(el);
        if (closest) {
            if (!closest._x_hideChildren) closest._x_hideChildren = [];
            closest._x_hideChildren.push(el);
        } else nextTick2(()=>{
            let hideAfterChildren = (el2)=>{
                let carry = Promise.all([
                    el2._x_hidePromise,
                    ...(el2._x_hideChildren || []).map(hideAfterChildren)
                ]).then(([i])=>i?.());
                delete el2._x_hidePromise;
                delete el2._x_hideChildren;
                return carry;
            };
            hideAfterChildren(el).catch((e)=>{
                if (!e.isFromCancelledTransition) throw e;
            });
        });
    });
};
function $8c83eaf28779ff46$var$closestHide(el) {
    let parent = el.parentNode;
    if (!parent) return;
    return parent._x_hidePromise ? parent : $8c83eaf28779ff46$var$closestHide(parent);
}
function $8c83eaf28779ff46$var$transition(el, setFunction, { during: during, start: start2, end: end } = {}, before = ()=>{}, after = ()=>{}) {
    if (el._x_transitioning) el._x_transitioning.cancel();
    if (Object.keys(during).length === 0 && Object.keys(start2).length === 0 && Object.keys(end).length === 0) {
        before();
        after();
        return;
    }
    let undoStart, undoDuring, undoEnd;
    $8c83eaf28779ff46$var$performTransition(el, {
        start () {
            undoStart = setFunction(el, start2);
        },
        during () {
            undoDuring = setFunction(el, during);
        },
        before: before,
        end () {
            undoStart();
            undoEnd = setFunction(el, end);
        },
        after: after,
        cleanup () {
            undoDuring();
            undoEnd();
        }
    });
}
function $8c83eaf28779ff46$var$performTransition(el, stages) {
    let interrupted, reachedBefore, reachedEnd;
    let finish = $8c83eaf28779ff46$var$once(()=>{
        $8c83eaf28779ff46$var$mutateDom(()=>{
            interrupted = true;
            if (!reachedBefore) stages.before();
            if (!reachedEnd) {
                stages.end();
                $8c83eaf28779ff46$var$releaseNextTicks();
            }
            stages.after();
            if (el.isConnected) stages.cleanup();
            delete el._x_transitioning;
        });
    });
    el._x_transitioning = {
        beforeCancels: [],
        beforeCancel (callback) {
            this.beforeCancels.push(callback);
        },
        cancel: $8c83eaf28779ff46$var$once(function() {
            while(this.beforeCancels.length)this.beforeCancels.shift()();
            finish();
        }),
        finish: finish
    };
    $8c83eaf28779ff46$var$mutateDom(()=>{
        stages.start();
        stages.during();
    });
    $8c83eaf28779ff46$var$holdNextTicks();
    requestAnimationFrame(()=>{
        if (interrupted) return;
        let duration = Number(getComputedStyle(el).transitionDuration.replace(/,.*/, "").replace("s", "")) * 1e3;
        let delay = Number(getComputedStyle(el).transitionDelay.replace(/,.*/, "").replace("s", "")) * 1e3;
        if (duration === 0) duration = Number(getComputedStyle(el).animationDuration.replace("s", "")) * 1e3;
        $8c83eaf28779ff46$var$mutateDom(()=>{
            stages.before();
        });
        reachedBefore = true;
        requestAnimationFrame(()=>{
            if (interrupted) return;
            $8c83eaf28779ff46$var$mutateDom(()=>{
                stages.end();
            });
            $8c83eaf28779ff46$var$releaseNextTicks();
            setTimeout(el._x_transitioning.finish, duration + delay);
            reachedEnd = true;
        });
    });
}
function $8c83eaf28779ff46$var$modifierValue(modifiers, key, fallback) {
    if (modifiers.indexOf(key) === -1) return fallback;
    const rawValue = modifiers[modifiers.indexOf(key) + 1];
    if (!rawValue) return fallback;
    if (key === "scale") {
        if (isNaN(rawValue)) return fallback;
    }
    if (key === "duration" || key === "delay") {
        let match = rawValue.match(/([0-9]+)ms/);
        if (match) return match[1];
    }
    if (key === "origin") {
        if ([
            "top",
            "right",
            "left",
            "center",
            "bottom"
        ].includes(modifiers[modifiers.indexOf(key) + 2])) return [
            rawValue,
            modifiers[modifiers.indexOf(key) + 2]
        ].join(" ");
    }
    return rawValue;
}
// packages/alpinejs/src/clone.js
var $8c83eaf28779ff46$var$isCloning = false;
function $8c83eaf28779ff46$var$skipDuringClone(callback, fallback = ()=>{}) {
    return (...args)=>$8c83eaf28779ff46$var$isCloning ? fallback(...args) : callback(...args);
}
function $8c83eaf28779ff46$var$onlyDuringClone(callback) {
    return (...args)=>$8c83eaf28779ff46$var$isCloning && callback(...args);
}
var $8c83eaf28779ff46$var$interceptors = [];
function $8c83eaf28779ff46$var$interceptClone(callback) {
    $8c83eaf28779ff46$var$interceptors.push(callback);
}
function $8c83eaf28779ff46$var$cloneNode(from, to) {
    $8c83eaf28779ff46$var$interceptors.forEach((i)=>i(from, to));
    $8c83eaf28779ff46$var$isCloning = true;
    $8c83eaf28779ff46$var$dontRegisterReactiveSideEffects(()=>{
        $8c83eaf28779ff46$var$initTree(to, (el, callback)=>{
            callback(el, ()=>{});
        });
    });
    $8c83eaf28779ff46$var$isCloning = false;
}
var $8c83eaf28779ff46$var$isCloningLegacy = false;
function $8c83eaf28779ff46$var$clone(oldEl, newEl) {
    if (!newEl._x_dataStack) newEl._x_dataStack = oldEl._x_dataStack;
    $8c83eaf28779ff46$var$isCloning = true;
    $8c83eaf28779ff46$var$isCloningLegacy = true;
    $8c83eaf28779ff46$var$dontRegisterReactiveSideEffects(()=>{
        $8c83eaf28779ff46$var$cloneTree(newEl);
    });
    $8c83eaf28779ff46$var$isCloning = false;
    $8c83eaf28779ff46$var$isCloningLegacy = false;
}
function $8c83eaf28779ff46$var$cloneTree(el) {
    let hasRunThroughFirstEl = false;
    let shallowWalker = (el2, callback)=>{
        $8c83eaf28779ff46$var$walk(el2, (el3, skip)=>{
            if (hasRunThroughFirstEl && $8c83eaf28779ff46$var$isRoot(el3)) return skip();
            hasRunThroughFirstEl = true;
            callback(el3, skip);
        });
    };
    $8c83eaf28779ff46$var$initTree(el, shallowWalker);
}
function $8c83eaf28779ff46$var$dontRegisterReactiveSideEffects(callback) {
    let cache = $8c83eaf28779ff46$var$effect;
    $8c83eaf28779ff46$var$overrideEffect((callback2, el)=>{
        let storedEffect = cache(callback2);
        $8c83eaf28779ff46$var$release(storedEffect);
        return ()=>{};
    });
    callback();
    $8c83eaf28779ff46$var$overrideEffect(cache);
}
// packages/alpinejs/src/utils/bind.js
function $8c83eaf28779ff46$var$bind(el, name, value, modifiers = []) {
    if (!el._x_bindings) el._x_bindings = $8c83eaf28779ff46$var$reactive({});
    el._x_bindings[name] = value;
    name = modifiers.includes("camel") ? $8c83eaf28779ff46$var$camelCase(name) : name;
    switch(name){
        case "value":
            $8c83eaf28779ff46$var$bindInputValue(el, value);
            break;
        case "style":
            $8c83eaf28779ff46$var$bindStyles(el, value);
            break;
        case "class":
            $8c83eaf28779ff46$var$bindClasses(el, value);
            break;
        case "selected":
        case "checked":
            $8c83eaf28779ff46$var$bindAttributeAndProperty(el, name, value);
            break;
        default:
            $8c83eaf28779ff46$var$bindAttribute(el, name, value);
            break;
    }
}
function $8c83eaf28779ff46$var$bindInputValue(el, value) {
    if ($8c83eaf28779ff46$var$isRadio(el)) {
        if (el.attributes.value === void 0) el.value = value;
        if (window.fromModel) {
            if (typeof value === "boolean") el.checked = $8c83eaf28779ff46$var$safeParseBoolean(el.value) === value;
            else el.checked = $8c83eaf28779ff46$var$checkedAttrLooseCompare(el.value, value);
        }
    } else if ($8c83eaf28779ff46$var$isCheckbox(el)) {
        if (Number.isInteger(value)) el.value = value;
        else if (!Array.isArray(value) && typeof value !== "boolean" && ![
            null,
            void 0
        ].includes(value)) el.value = String(value);
        else if (Array.isArray(value)) el.checked = value.some((val)=>$8c83eaf28779ff46$var$checkedAttrLooseCompare(val, el.value));
        else el.checked = !!value;
    } else if (el.tagName === "SELECT") $8c83eaf28779ff46$var$updateSelect(el, value);
    else {
        if (el.value === value) return;
        el.value = value === void 0 ? "" : value;
    }
}
function $8c83eaf28779ff46$var$bindClasses(el, value) {
    if (el._x_undoAddedClasses) el._x_undoAddedClasses();
    el._x_undoAddedClasses = $8c83eaf28779ff46$var$setClasses(el, value);
}
function $8c83eaf28779ff46$var$bindStyles(el, value) {
    if (el._x_undoAddedStyles) el._x_undoAddedStyles();
    el._x_undoAddedStyles = $8c83eaf28779ff46$var$setStyles(el, value);
}
function $8c83eaf28779ff46$var$bindAttributeAndProperty(el, name, value) {
    $8c83eaf28779ff46$var$bindAttribute(el, name, value);
    $8c83eaf28779ff46$var$setPropertyIfChanged(el, name, value);
}
function $8c83eaf28779ff46$var$bindAttribute(el, name, value) {
    if ([
        null,
        void 0,
        false
    ].includes(value) && $8c83eaf28779ff46$var$attributeShouldntBePreservedIfFalsy(name)) el.removeAttribute(name);
    else {
        if ($8c83eaf28779ff46$var$isBooleanAttr(name)) value = name;
        $8c83eaf28779ff46$var$setIfChanged(el, name, value);
    }
}
function $8c83eaf28779ff46$var$setIfChanged(el, attrName, value) {
    if (el.getAttribute(attrName) != value) el.setAttribute(attrName, value);
}
function $8c83eaf28779ff46$var$setPropertyIfChanged(el, propName, value) {
    if (el[propName] !== value) el[propName] = value;
}
function $8c83eaf28779ff46$var$updateSelect(el, value) {
    const arrayWrappedValue = [].concat(value).map((value2)=>{
        return value2 + "";
    });
    Array.from(el.options).forEach((option)=>{
        option.selected = arrayWrappedValue.includes(option.value);
    });
}
function $8c83eaf28779ff46$var$camelCase(subject) {
    return subject.toLowerCase().replace(/-(\w)/g, (match, char)=>char.toUpperCase());
}
function $8c83eaf28779ff46$var$checkedAttrLooseCompare(valueA, valueB) {
    return valueA == valueB;
}
function $8c83eaf28779ff46$var$safeParseBoolean(rawValue) {
    if ([
        1,
        "1",
        "true",
        "on",
        "yes",
        true
    ].includes(rawValue)) return true;
    if ([
        0,
        "0",
        "false",
        "off",
        "no",
        false
    ].includes(rawValue)) return false;
    return rawValue ? Boolean(rawValue) : null;
}
var $8c83eaf28779ff46$var$booleanAttributes = /* @__PURE__ */ new Set([
    "allowfullscreen",
    "async",
    "autofocus",
    "autoplay",
    "checked",
    "controls",
    "default",
    "defer",
    "disabled",
    "formnovalidate",
    "inert",
    "ismap",
    "itemscope",
    "loop",
    "multiple",
    "muted",
    "nomodule",
    "novalidate",
    "open",
    "playsinline",
    "readonly",
    "required",
    "reversed",
    "selected",
    "shadowrootclonable",
    "shadowrootdelegatesfocus",
    "shadowrootserializable"
]);
function $8c83eaf28779ff46$var$isBooleanAttr(attrName) {
    return $8c83eaf28779ff46$var$booleanAttributes.has(attrName);
}
function $8c83eaf28779ff46$var$attributeShouldntBePreservedIfFalsy(name) {
    return ![
        "aria-pressed",
        "aria-checked",
        "aria-expanded",
        "aria-selected"
    ].includes(name);
}
function $8c83eaf28779ff46$var$getBinding(el, name, fallback) {
    if (el._x_bindings && el._x_bindings[name] !== void 0) return el._x_bindings[name];
    return $8c83eaf28779ff46$var$getAttributeBinding(el, name, fallback);
}
function $8c83eaf28779ff46$var$extractProp(el, name, fallback, extract = true) {
    if (el._x_bindings && el._x_bindings[name] !== void 0) return el._x_bindings[name];
    if (el._x_inlineBindings && el._x_inlineBindings[name] !== void 0) {
        let binding = el._x_inlineBindings[name];
        binding.extract = extract;
        return $8c83eaf28779ff46$var$dontAutoEvaluateFunctions(()=>{
            return $8c83eaf28779ff46$var$evaluate(el, binding.expression);
        });
    }
    return $8c83eaf28779ff46$var$getAttributeBinding(el, name, fallback);
}
function $8c83eaf28779ff46$var$getAttributeBinding(el, name, fallback) {
    let attr = el.getAttribute(name);
    if (attr === null) return typeof fallback === "function" ? fallback() : fallback;
    if (attr === "") return true;
    if ($8c83eaf28779ff46$var$isBooleanAttr(name)) return !![
        name,
        "true"
    ].includes(attr);
    return attr;
}
function $8c83eaf28779ff46$var$isCheckbox(el) {
    return el.type === "checkbox" || el.localName === "ui-checkbox" || el.localName === "ui-switch";
}
function $8c83eaf28779ff46$var$isRadio(el) {
    return el.type === "radio" || el.localName === "ui-radio";
}
// packages/alpinejs/src/utils/debounce.js
function $8c83eaf28779ff46$var$debounce(func, wait) {
    let timeout;
    return function() {
        const context = this, args = arguments;
        const later = function() {
            timeout = null;
            func.apply(context, args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}
// packages/alpinejs/src/utils/throttle.js
function $8c83eaf28779ff46$var$throttle(func, limit) {
    let inThrottle;
    return function() {
        let context = this, args = arguments;
        if (!inThrottle) {
            func.apply(context, args);
            inThrottle = true;
            setTimeout(()=>inThrottle = false, limit);
        }
    };
}
// packages/alpinejs/src/entangle.js
function $8c83eaf28779ff46$var$entangle({ get: outerGet, set: outerSet }, { get: innerGet, set: innerSet }) {
    let firstRun = true;
    let outerHash;
    let innerHash;
    let reference = $8c83eaf28779ff46$var$effect(()=>{
        let outer = outerGet();
        let inner = innerGet();
        if (firstRun) {
            innerSet($8c83eaf28779ff46$var$cloneIfObject(outer));
            firstRun = false;
        } else {
            let outerHashLatest = JSON.stringify(outer);
            let innerHashLatest = JSON.stringify(inner);
            if (outerHashLatest !== outerHash) innerSet($8c83eaf28779ff46$var$cloneIfObject(outer));
            else if (outerHashLatest !== innerHashLatest) outerSet($8c83eaf28779ff46$var$cloneIfObject(inner));
        }
        outerHash = JSON.stringify(outerGet());
        innerHash = JSON.stringify(innerGet());
    });
    return ()=>{
        $8c83eaf28779ff46$var$release(reference);
    };
}
function $8c83eaf28779ff46$var$cloneIfObject(value) {
    return typeof value === "object" ? JSON.parse(JSON.stringify(value)) : value;
}
// packages/alpinejs/src/plugin.js
function $8c83eaf28779ff46$var$plugin(callback) {
    let callbacks = Array.isArray(callback) ? callback : [
        callback
    ];
    callbacks.forEach((i)=>i($8c83eaf28779ff46$var$alpine_default));
}
// packages/alpinejs/src/store.js
var $8c83eaf28779ff46$var$stores = {};
var $8c83eaf28779ff46$var$isReactive = false;
function $8c83eaf28779ff46$var$store(name, value) {
    if (!$8c83eaf28779ff46$var$isReactive) {
        $8c83eaf28779ff46$var$stores = $8c83eaf28779ff46$var$reactive($8c83eaf28779ff46$var$stores);
        $8c83eaf28779ff46$var$isReactive = true;
    }
    if (value === void 0) return $8c83eaf28779ff46$var$stores[name];
    $8c83eaf28779ff46$var$stores[name] = value;
    $8c83eaf28779ff46$var$initInterceptors($8c83eaf28779ff46$var$stores[name]);
    if (typeof value === "object" && value !== null && value.hasOwnProperty("init") && typeof value.init === "function") $8c83eaf28779ff46$var$stores[name].init();
}
function $8c83eaf28779ff46$var$getStores() {
    return $8c83eaf28779ff46$var$stores;
}
// packages/alpinejs/src/binds.js
var $8c83eaf28779ff46$var$binds = {};
function $8c83eaf28779ff46$var$bind2(name, bindings) {
    let getBindings = typeof bindings !== "function" ? ()=>bindings : bindings;
    if (name instanceof Element) return $8c83eaf28779ff46$var$applyBindingsObject(name, getBindings());
    else $8c83eaf28779ff46$var$binds[name] = getBindings;
    return ()=>{};
}
function $8c83eaf28779ff46$var$injectBindingProviders(obj) {
    Object.entries($8c83eaf28779ff46$var$binds).forEach(([name, callback])=>{
        Object.defineProperty(obj, name, {
            get () {
                return (...args)=>{
                    return callback(...args);
                };
            }
        });
    });
    return obj;
}
function $8c83eaf28779ff46$var$applyBindingsObject(el, obj, original) {
    let cleanupRunners = [];
    while(cleanupRunners.length)cleanupRunners.pop()();
    let attributes = Object.entries(obj).map(([name, value])=>({
            name: name,
            value: value
        }));
    let staticAttributes = $8c83eaf28779ff46$var$attributesOnly(attributes);
    attributes = attributes.map((attribute)=>{
        if (staticAttributes.find((attr)=>attr.name === attribute.name)) return {
            name: `x-bind:${attribute.name}`,
            value: `"${attribute.value}"`
        };
        return attribute;
    });
    $8c83eaf28779ff46$var$directives(el, attributes, original).map((handle)=>{
        cleanupRunners.push(handle.runCleanups);
        handle();
    });
    return ()=>{
        while(cleanupRunners.length)cleanupRunners.pop()();
    };
}
// packages/alpinejs/src/datas.js
var $8c83eaf28779ff46$var$datas = {};
function $8c83eaf28779ff46$var$data(name, callback) {
    $8c83eaf28779ff46$var$datas[name] = callback;
}
function $8c83eaf28779ff46$var$injectDataProviders(obj, context) {
    Object.entries($8c83eaf28779ff46$var$datas).forEach(([name, callback])=>{
        Object.defineProperty(obj, name, {
            get () {
                return (...args)=>{
                    return callback.bind(context)(...args);
                };
            },
            enumerable: false
        });
    });
    return obj;
}
// packages/alpinejs/src/alpine.js
var $8c83eaf28779ff46$var$Alpine = {
    get reactive () {
        return $8c83eaf28779ff46$var$reactive;
    },
    get release () {
        return $8c83eaf28779ff46$var$release;
    },
    get effect () {
        return $8c83eaf28779ff46$var$effect;
    },
    get raw () {
        return $8c83eaf28779ff46$var$raw;
    },
    version: "3.15.4",
    flushAndStopDeferringMutations: $8c83eaf28779ff46$var$flushAndStopDeferringMutations,
    dontAutoEvaluateFunctions: $8c83eaf28779ff46$var$dontAutoEvaluateFunctions,
    disableEffectScheduling: $8c83eaf28779ff46$var$disableEffectScheduling,
    startObservingMutations: $8c83eaf28779ff46$var$startObservingMutations,
    stopObservingMutations: $8c83eaf28779ff46$var$stopObservingMutations,
    setReactivityEngine: $8c83eaf28779ff46$var$setReactivityEngine,
    onAttributeRemoved: $8c83eaf28779ff46$var$onAttributeRemoved,
    onAttributesAdded: $8c83eaf28779ff46$var$onAttributesAdded,
    closestDataStack: $8c83eaf28779ff46$var$closestDataStack,
    skipDuringClone: $8c83eaf28779ff46$var$skipDuringClone,
    onlyDuringClone: $8c83eaf28779ff46$var$onlyDuringClone,
    addRootSelector: $8c83eaf28779ff46$var$addRootSelector,
    addInitSelector: $8c83eaf28779ff46$var$addInitSelector,
    setErrorHandler: $8c83eaf28779ff46$var$setErrorHandler,
    interceptClone: $8c83eaf28779ff46$var$interceptClone,
    addScopeToNode: $8c83eaf28779ff46$var$addScopeToNode,
    deferMutations: $8c83eaf28779ff46$var$deferMutations,
    mapAttributes: $8c83eaf28779ff46$var$mapAttributes,
    evaluateLater: $8c83eaf28779ff46$var$evaluateLater,
    interceptInit: $8c83eaf28779ff46$var$interceptInit,
    initInterceptors: $8c83eaf28779ff46$var$initInterceptors,
    injectMagics: $8c83eaf28779ff46$var$injectMagics,
    setEvaluator: $8c83eaf28779ff46$var$setEvaluator,
    setRawEvaluator: $8c83eaf28779ff46$var$setRawEvaluator,
    mergeProxies: $8c83eaf28779ff46$var$mergeProxies,
    extractProp: $8c83eaf28779ff46$var$extractProp,
    findClosest: $8c83eaf28779ff46$var$findClosest,
    onElRemoved: $8c83eaf28779ff46$var$onElRemoved,
    closestRoot: $8c83eaf28779ff46$var$closestRoot,
    destroyTree: $8c83eaf28779ff46$var$destroyTree,
    interceptor: $8c83eaf28779ff46$var$interceptor,
    transition: // INTERNAL: not public API and is subject to change without major release.
    $8c83eaf28779ff46$var$transition,
    setStyles: // INTERNAL
    $8c83eaf28779ff46$var$setStyles,
    mutateDom: // INTERNAL
    $8c83eaf28779ff46$var$mutateDom,
    directive: $8c83eaf28779ff46$var$directive,
    entangle: $8c83eaf28779ff46$var$entangle,
    throttle: $8c83eaf28779ff46$var$throttle,
    debounce: $8c83eaf28779ff46$var$debounce,
    evaluate: $8c83eaf28779ff46$var$evaluate,
    evaluateRaw: $8c83eaf28779ff46$var$evaluateRaw,
    initTree: $8c83eaf28779ff46$var$initTree,
    nextTick: $8c83eaf28779ff46$var$nextTick,
    prefixed: $8c83eaf28779ff46$var$prefix,
    prefix: $8c83eaf28779ff46$var$setPrefix,
    plugin: $8c83eaf28779ff46$var$plugin,
    magic: $8c83eaf28779ff46$var$magic,
    store: $8c83eaf28779ff46$var$store,
    start: $8c83eaf28779ff46$var$start,
    clone: $8c83eaf28779ff46$var$clone,
    cloneNode: // INTERNAL
    $8c83eaf28779ff46$var$cloneNode,
    // INTERNAL
    bound: $8c83eaf28779ff46$var$getBinding,
    $data: $8c83eaf28779ff46$var$scope,
    watch: $8c83eaf28779ff46$var$watch,
    walk: $8c83eaf28779ff46$var$walk,
    data: $8c83eaf28779ff46$var$data,
    bind: $8c83eaf28779ff46$var$bind2
};
var $8c83eaf28779ff46$var$alpine_default = $8c83eaf28779ff46$var$Alpine;
// node_modules/@vue/shared/dist/shared.esm-bundler.js
function $8c83eaf28779ff46$var$makeMap(str, expectsLowerCase) {
    const map = /* @__PURE__ */ Object.create(null);
    const list = str.split(",");
    for(let i = 0; i < list.length; i++)map[list[i]] = true;
    return expectsLowerCase ? (val)=>!!map[val.toLowerCase()] : (val)=>!!map[val];
}
var $8c83eaf28779ff46$var$specialBooleanAttrs = `itemscope,allowfullscreen,formnovalidate,ismap,nomodule,novalidate,readonly`;
var $8c83eaf28779ff46$var$isBooleanAttr2 = /* @__PURE__ */ $8c83eaf28779ff46$var$makeMap($8c83eaf28779ff46$var$specialBooleanAttrs + `,async,autofocus,autoplay,controls,default,defer,disabled,hidden,loop,open,required,reversed,scoped,seamless,checked,muted,multiple,selected`);
var $8c83eaf28779ff46$var$EMPTY_OBJ = Object.freeze({});
var $8c83eaf28779ff46$var$EMPTY_ARR = Object.freeze([]);
var $8c83eaf28779ff46$var$hasOwnProperty = Object.prototype.hasOwnProperty;
var $8c83eaf28779ff46$var$hasOwn = (val, key)=>$8c83eaf28779ff46$var$hasOwnProperty.call(val, key);
var $8c83eaf28779ff46$var$isArray = Array.isArray;
var $8c83eaf28779ff46$var$isMap = (val)=>$8c83eaf28779ff46$var$toTypeString(val) === "[object Map]";
var $8c83eaf28779ff46$var$isString = (val)=>typeof val === "string";
var $8c83eaf28779ff46$var$isSymbol = (val)=>typeof val === "symbol";
var $8c83eaf28779ff46$var$isObject = (val)=>val !== null && typeof val === "object";
var $8c83eaf28779ff46$var$objectToString = Object.prototype.toString;
var $8c83eaf28779ff46$var$toTypeString = (value)=>$8c83eaf28779ff46$var$objectToString.call(value);
var $8c83eaf28779ff46$var$toRawType = (value)=>{
    return $8c83eaf28779ff46$var$toTypeString(value).slice(8, -1);
};
var $8c83eaf28779ff46$var$isIntegerKey = (key)=>$8c83eaf28779ff46$var$isString(key) && key !== "NaN" && key[0] !== "-" && "" + parseInt(key, 10) === key;
var $8c83eaf28779ff46$var$cacheStringFunction = (fn)=>{
    const cache = /* @__PURE__ */ Object.create(null);
    return (str)=>{
        const hit = cache[str];
        return hit || (cache[str] = fn(str));
    };
};
var $8c83eaf28779ff46$var$camelizeRE = /-(\w)/g;
var $8c83eaf28779ff46$var$camelize = $8c83eaf28779ff46$var$cacheStringFunction((str)=>{
    return str.replace($8c83eaf28779ff46$var$camelizeRE, (_, c)=>c ? c.toUpperCase() : "");
});
var $8c83eaf28779ff46$var$hyphenateRE = /\B([A-Z])/g;
var $8c83eaf28779ff46$var$hyphenate = $8c83eaf28779ff46$var$cacheStringFunction((str)=>str.replace($8c83eaf28779ff46$var$hyphenateRE, "-$1").toLowerCase());
var $8c83eaf28779ff46$var$capitalize = $8c83eaf28779ff46$var$cacheStringFunction((str)=>str.charAt(0).toUpperCase() + str.slice(1));
var $8c83eaf28779ff46$var$toHandlerKey = $8c83eaf28779ff46$var$cacheStringFunction((str)=>str ? `on${$8c83eaf28779ff46$var$capitalize(str)}` : ``);
var $8c83eaf28779ff46$var$hasChanged = (value, oldValue)=>value !== oldValue && (value === value || oldValue === oldValue);
// node_modules/@vue/reactivity/dist/reactivity.esm-bundler.js
var $8c83eaf28779ff46$var$targetMap = /* @__PURE__ */ new WeakMap();
var $8c83eaf28779ff46$var$effectStack = [];
var $8c83eaf28779ff46$var$activeEffect;
var $8c83eaf28779ff46$var$ITERATE_KEY = Symbol("iterate");
var $8c83eaf28779ff46$var$MAP_KEY_ITERATE_KEY = Symbol("Map key iterate");
function $8c83eaf28779ff46$var$isEffect(fn) {
    return fn && fn._isEffect === true;
}
function $8c83eaf28779ff46$var$effect2(fn, options = $8c83eaf28779ff46$var$EMPTY_OBJ) {
    if ($8c83eaf28779ff46$var$isEffect(fn)) fn = fn.raw;
    const effect3 = $8c83eaf28779ff46$var$createReactiveEffect(fn, options);
    if (!options.lazy) effect3();
    return effect3;
}
function $8c83eaf28779ff46$var$stop(effect3) {
    if (effect3.active) {
        $8c83eaf28779ff46$var$cleanup(effect3);
        if (effect3.options.onStop) effect3.options.onStop();
        effect3.active = false;
    }
}
var $8c83eaf28779ff46$var$uid = 0;
function $8c83eaf28779ff46$var$createReactiveEffect(fn, options) {
    const effect3 = function reactiveEffect() {
        if (!effect3.active) return fn();
        if (!$8c83eaf28779ff46$var$effectStack.includes(effect3)) {
            $8c83eaf28779ff46$var$cleanup(effect3);
            try {
                $8c83eaf28779ff46$var$enableTracking();
                $8c83eaf28779ff46$var$effectStack.push(effect3);
                $8c83eaf28779ff46$var$activeEffect = effect3;
                return fn();
            } finally{
                $8c83eaf28779ff46$var$effectStack.pop();
                $8c83eaf28779ff46$var$resetTracking();
                $8c83eaf28779ff46$var$activeEffect = $8c83eaf28779ff46$var$effectStack[$8c83eaf28779ff46$var$effectStack.length - 1];
            }
        }
    };
    effect3.id = $8c83eaf28779ff46$var$uid++;
    effect3.allowRecurse = !!options.allowRecurse;
    effect3._isEffect = true;
    effect3.active = true;
    effect3.raw = fn;
    effect3.deps = [];
    effect3.options = options;
    return effect3;
}
function $8c83eaf28779ff46$var$cleanup(effect3) {
    const { deps: deps } = effect3;
    if (deps.length) {
        for(let i = 0; i < deps.length; i++)deps[i].delete(effect3);
        deps.length = 0;
    }
}
var $8c83eaf28779ff46$var$shouldTrack = true;
var $8c83eaf28779ff46$var$trackStack = [];
function $8c83eaf28779ff46$var$pauseTracking() {
    $8c83eaf28779ff46$var$trackStack.push($8c83eaf28779ff46$var$shouldTrack);
    $8c83eaf28779ff46$var$shouldTrack = false;
}
function $8c83eaf28779ff46$var$enableTracking() {
    $8c83eaf28779ff46$var$trackStack.push($8c83eaf28779ff46$var$shouldTrack);
    $8c83eaf28779ff46$var$shouldTrack = true;
}
function $8c83eaf28779ff46$var$resetTracking() {
    const last = $8c83eaf28779ff46$var$trackStack.pop();
    $8c83eaf28779ff46$var$shouldTrack = last === void 0 ? true : last;
}
function $8c83eaf28779ff46$var$track(target, type, key) {
    if (!$8c83eaf28779ff46$var$shouldTrack || $8c83eaf28779ff46$var$activeEffect === void 0) return;
    let depsMap = $8c83eaf28779ff46$var$targetMap.get(target);
    if (!depsMap) $8c83eaf28779ff46$var$targetMap.set(target, depsMap = /* @__PURE__ */ new Map());
    let dep = depsMap.get(key);
    if (!dep) depsMap.set(key, dep = /* @__PURE__ */ new Set());
    if (!dep.has($8c83eaf28779ff46$var$activeEffect)) {
        dep.add($8c83eaf28779ff46$var$activeEffect);
        $8c83eaf28779ff46$var$activeEffect.deps.push(dep);
        if ($8c83eaf28779ff46$var$activeEffect.options.onTrack) $8c83eaf28779ff46$var$activeEffect.options.onTrack({
            effect: $8c83eaf28779ff46$var$activeEffect,
            target: target,
            type: type,
            key: key
        });
    }
}
function $8c83eaf28779ff46$var$trigger(target, type, key, newValue, oldValue, oldTarget) {
    const depsMap = $8c83eaf28779ff46$var$targetMap.get(target);
    if (!depsMap) return;
    const effects = /* @__PURE__ */ new Set();
    const add2 = (effectsToAdd)=>{
        if (effectsToAdd) effectsToAdd.forEach((effect3)=>{
            if (effect3 !== $8c83eaf28779ff46$var$activeEffect || effect3.allowRecurse) effects.add(effect3);
        });
    };
    if (type === "clear") depsMap.forEach(add2);
    else if (key === "length" && $8c83eaf28779ff46$var$isArray(target)) depsMap.forEach((dep, key2)=>{
        if (key2 === "length" || key2 >= newValue) add2(dep);
    });
    else {
        if (key !== void 0) add2(depsMap.get(key));
        switch(type){
            case "add":
                if (!$8c83eaf28779ff46$var$isArray(target)) {
                    add2(depsMap.get($8c83eaf28779ff46$var$ITERATE_KEY));
                    if ($8c83eaf28779ff46$var$isMap(target)) add2(depsMap.get($8c83eaf28779ff46$var$MAP_KEY_ITERATE_KEY));
                } else if ($8c83eaf28779ff46$var$isIntegerKey(key)) add2(depsMap.get("length"));
                break;
            case "delete":
                if (!$8c83eaf28779ff46$var$isArray(target)) {
                    add2(depsMap.get($8c83eaf28779ff46$var$ITERATE_KEY));
                    if ($8c83eaf28779ff46$var$isMap(target)) add2(depsMap.get($8c83eaf28779ff46$var$MAP_KEY_ITERATE_KEY));
                }
                break;
            case "set":
                if ($8c83eaf28779ff46$var$isMap(target)) add2(depsMap.get($8c83eaf28779ff46$var$ITERATE_KEY));
                break;
        }
    }
    const run = (effect3)=>{
        if (effect3.options.onTrigger) effect3.options.onTrigger({
            effect: effect3,
            target: target,
            key: key,
            type: type,
            newValue: newValue,
            oldValue: oldValue,
            oldTarget: oldTarget
        });
        if (effect3.options.scheduler) effect3.options.scheduler(effect3);
        else effect3();
    };
    effects.forEach(run);
}
var $8c83eaf28779ff46$var$isNonTrackableKeys = /* @__PURE__ */ $8c83eaf28779ff46$var$makeMap(`__proto__,__v_isRef,__isVue`);
var $8c83eaf28779ff46$var$builtInSymbols = new Set(Object.getOwnPropertyNames(Symbol).map((key)=>Symbol[key]).filter($8c83eaf28779ff46$var$isSymbol));
var $8c83eaf28779ff46$var$get2 = /* @__PURE__ */ $8c83eaf28779ff46$var$createGetter();
var $8c83eaf28779ff46$var$readonlyGet = /* @__PURE__ */ $8c83eaf28779ff46$var$createGetter(true);
var $8c83eaf28779ff46$var$arrayInstrumentations = /* @__PURE__ */ $8c83eaf28779ff46$var$createArrayInstrumentations();
function $8c83eaf28779ff46$var$createArrayInstrumentations() {
    const instrumentations = {};
    [
        "includes",
        "indexOf",
        "lastIndexOf"
    ].forEach((key)=>{
        instrumentations[key] = function(...args) {
            const arr = $8c83eaf28779ff46$var$toRaw(this);
            for(let i = 0, l = this.length; i < l; i++)$8c83eaf28779ff46$var$track(arr, "get", i + "");
            const res = arr[key](...args);
            if (res === -1 || res === false) return arr[key](...args.map($8c83eaf28779ff46$var$toRaw));
            else return res;
        };
    });
    [
        "push",
        "pop",
        "shift",
        "unshift",
        "splice"
    ].forEach((key)=>{
        instrumentations[key] = function(...args) {
            $8c83eaf28779ff46$var$pauseTracking();
            const res = $8c83eaf28779ff46$var$toRaw(this)[key].apply(this, args);
            $8c83eaf28779ff46$var$resetTracking();
            return res;
        };
    });
    return instrumentations;
}
function $8c83eaf28779ff46$var$createGetter(isReadonly = false, shallow = false) {
    return function get3(target, key, receiver) {
        if (key === "__v_isReactive") return !isReadonly;
        else if (key === "__v_isReadonly") return isReadonly;
        else if (key === "__v_raw" && receiver === (isReadonly ? shallow ? $8c83eaf28779ff46$var$shallowReadonlyMap : $8c83eaf28779ff46$var$readonlyMap : shallow ? $8c83eaf28779ff46$var$shallowReactiveMap : $8c83eaf28779ff46$var$reactiveMap).get(target)) return target;
        const targetIsArray = $8c83eaf28779ff46$var$isArray(target);
        if (!isReadonly && targetIsArray && $8c83eaf28779ff46$var$hasOwn($8c83eaf28779ff46$var$arrayInstrumentations, key)) return Reflect.get($8c83eaf28779ff46$var$arrayInstrumentations, key, receiver);
        const res = Reflect.get(target, key, receiver);
        if ($8c83eaf28779ff46$var$isSymbol(key) ? $8c83eaf28779ff46$var$builtInSymbols.has(key) : $8c83eaf28779ff46$var$isNonTrackableKeys(key)) return res;
        if (!isReadonly) $8c83eaf28779ff46$var$track(target, "get", key);
        if (shallow) return res;
        if ($8c83eaf28779ff46$var$isRef(res)) {
            const shouldUnwrap = !targetIsArray || !$8c83eaf28779ff46$var$isIntegerKey(key);
            return shouldUnwrap ? res.value : res;
        }
        if ($8c83eaf28779ff46$var$isObject(res)) return isReadonly ? $8c83eaf28779ff46$var$readonly(res) : $8c83eaf28779ff46$var$reactive2(res);
        return res;
    };
}
var $8c83eaf28779ff46$var$set2 = /* @__PURE__ */ $8c83eaf28779ff46$var$createSetter();
function $8c83eaf28779ff46$var$createSetter(shallow = false) {
    return function set3(target, key, value, receiver) {
        let oldValue = target[key];
        if (!shallow) {
            value = $8c83eaf28779ff46$var$toRaw(value);
            oldValue = $8c83eaf28779ff46$var$toRaw(oldValue);
            if (!$8c83eaf28779ff46$var$isArray(target) && $8c83eaf28779ff46$var$isRef(oldValue) && !$8c83eaf28779ff46$var$isRef(value)) {
                oldValue.value = value;
                return true;
            }
        }
        const hadKey = $8c83eaf28779ff46$var$isArray(target) && $8c83eaf28779ff46$var$isIntegerKey(key) ? Number(key) < target.length : $8c83eaf28779ff46$var$hasOwn(target, key);
        const result = Reflect.set(target, key, value, receiver);
        if (target === $8c83eaf28779ff46$var$toRaw(receiver)) {
            if (!hadKey) $8c83eaf28779ff46$var$trigger(target, "add", key, value);
            else if ($8c83eaf28779ff46$var$hasChanged(value, oldValue)) $8c83eaf28779ff46$var$trigger(target, "set", key, value, oldValue);
        }
        return result;
    };
}
function $8c83eaf28779ff46$var$deleteProperty(target, key) {
    const hadKey = $8c83eaf28779ff46$var$hasOwn(target, key);
    const oldValue = target[key];
    const result = Reflect.deleteProperty(target, key);
    if (result && hadKey) $8c83eaf28779ff46$var$trigger(target, "delete", key, void 0, oldValue);
    return result;
}
function $8c83eaf28779ff46$var$has(target, key) {
    const result = Reflect.has(target, key);
    if (!$8c83eaf28779ff46$var$isSymbol(key) || !$8c83eaf28779ff46$var$builtInSymbols.has(key)) $8c83eaf28779ff46$var$track(target, "has", key);
    return result;
}
function $8c83eaf28779ff46$var$ownKeys(target) {
    $8c83eaf28779ff46$var$track(target, "iterate", $8c83eaf28779ff46$var$isArray(target) ? "length" : $8c83eaf28779ff46$var$ITERATE_KEY);
    return Reflect.ownKeys(target);
}
var $8c83eaf28779ff46$var$mutableHandlers = {
    get: $8c83eaf28779ff46$var$get2,
    set: $8c83eaf28779ff46$var$set2,
    deleteProperty: $8c83eaf28779ff46$var$deleteProperty,
    has: $8c83eaf28779ff46$var$has,
    ownKeys: $8c83eaf28779ff46$var$ownKeys
};
var $8c83eaf28779ff46$var$readonlyHandlers = {
    get: $8c83eaf28779ff46$var$readonlyGet,
    set (target, key) {
        console.warn(`Set operation on key "${String(key)}" failed: target is readonly.`, target);
        return true;
    },
    deleteProperty (target, key) {
        console.warn(`Delete operation on key "${String(key)}" failed: target is readonly.`, target);
        return true;
    }
};
var $8c83eaf28779ff46$var$toReactive = (value)=>$8c83eaf28779ff46$var$isObject(value) ? $8c83eaf28779ff46$var$reactive2(value) : value;
var $8c83eaf28779ff46$var$toReadonly = (value)=>$8c83eaf28779ff46$var$isObject(value) ? $8c83eaf28779ff46$var$readonly(value) : value;
var $8c83eaf28779ff46$var$toShallow = (value)=>value;
var $8c83eaf28779ff46$var$getProto = (v)=>Reflect.getPrototypeOf(v);
function $8c83eaf28779ff46$var$get$1(target, key, isReadonly = false, isShallow = false) {
    target = target["__v_raw"];
    const rawTarget = $8c83eaf28779ff46$var$toRaw(target);
    const rawKey = $8c83eaf28779ff46$var$toRaw(key);
    if (key !== rawKey) !isReadonly && $8c83eaf28779ff46$var$track(rawTarget, "get", key);
    !isReadonly && $8c83eaf28779ff46$var$track(rawTarget, "get", rawKey);
    const { has: has2 } = $8c83eaf28779ff46$var$getProto(rawTarget);
    const wrap = isShallow ? $8c83eaf28779ff46$var$toShallow : isReadonly ? $8c83eaf28779ff46$var$toReadonly : $8c83eaf28779ff46$var$toReactive;
    if (has2.call(rawTarget, key)) return wrap(target.get(key));
    else if (has2.call(rawTarget, rawKey)) return wrap(target.get(rawKey));
    else if (target !== rawTarget) target.get(key);
}
function $8c83eaf28779ff46$var$has$1(key, isReadonly = false) {
    const target = this["__v_raw"];
    const rawTarget = $8c83eaf28779ff46$var$toRaw(target);
    const rawKey = $8c83eaf28779ff46$var$toRaw(key);
    if (key !== rawKey) !isReadonly && $8c83eaf28779ff46$var$track(rawTarget, "has", key);
    !isReadonly && $8c83eaf28779ff46$var$track(rawTarget, "has", rawKey);
    return key === rawKey ? target.has(key) : target.has(key) || target.has(rawKey);
}
function $8c83eaf28779ff46$var$size(target, isReadonly = false) {
    target = target["__v_raw"];
    !isReadonly && $8c83eaf28779ff46$var$track($8c83eaf28779ff46$var$toRaw(target), "iterate", $8c83eaf28779ff46$var$ITERATE_KEY);
    return Reflect.get(target, "size", target);
}
function $8c83eaf28779ff46$var$add(value) {
    value = $8c83eaf28779ff46$var$toRaw(value);
    const target = $8c83eaf28779ff46$var$toRaw(this);
    const proto = $8c83eaf28779ff46$var$getProto(target);
    const hadKey = proto.has.call(target, value);
    if (!hadKey) {
        target.add(value);
        $8c83eaf28779ff46$var$trigger(target, "add", value, value);
    }
    return this;
}
function $8c83eaf28779ff46$var$set$1(key, value) {
    value = $8c83eaf28779ff46$var$toRaw(value);
    const target = $8c83eaf28779ff46$var$toRaw(this);
    const { has: has2, get: get3 } = $8c83eaf28779ff46$var$getProto(target);
    let hadKey = has2.call(target, key);
    if (!hadKey) {
        key = $8c83eaf28779ff46$var$toRaw(key);
        hadKey = has2.call(target, key);
    } else $8c83eaf28779ff46$var$checkIdentityKeys(target, has2, key);
    const oldValue = get3.call(target, key);
    target.set(key, value);
    if (!hadKey) $8c83eaf28779ff46$var$trigger(target, "add", key, value);
    else if ($8c83eaf28779ff46$var$hasChanged(value, oldValue)) $8c83eaf28779ff46$var$trigger(target, "set", key, value, oldValue);
    return this;
}
function $8c83eaf28779ff46$var$deleteEntry(key) {
    const target = $8c83eaf28779ff46$var$toRaw(this);
    const { has: has2, get: get3 } = $8c83eaf28779ff46$var$getProto(target);
    let hadKey = has2.call(target, key);
    if (!hadKey) {
        key = $8c83eaf28779ff46$var$toRaw(key);
        hadKey = has2.call(target, key);
    } else $8c83eaf28779ff46$var$checkIdentityKeys(target, has2, key);
    const oldValue = get3 ? get3.call(target, key) : void 0;
    const result = target.delete(key);
    if (hadKey) $8c83eaf28779ff46$var$trigger(target, "delete", key, void 0, oldValue);
    return result;
}
function $8c83eaf28779ff46$var$clear() {
    const target = $8c83eaf28779ff46$var$toRaw(this);
    const hadItems = target.size !== 0;
    const oldTarget = $8c83eaf28779ff46$var$isMap(target) ? new Map(target) : new Set(target);
    const result = target.clear();
    if (hadItems) $8c83eaf28779ff46$var$trigger(target, "clear", void 0, void 0, oldTarget);
    return result;
}
function $8c83eaf28779ff46$var$createForEach(isReadonly, isShallow) {
    return function forEach(callback, thisArg) {
        const observed = this;
        const target = observed["__v_raw"];
        const rawTarget = $8c83eaf28779ff46$var$toRaw(target);
        const wrap = isShallow ? $8c83eaf28779ff46$var$toShallow : isReadonly ? $8c83eaf28779ff46$var$toReadonly : $8c83eaf28779ff46$var$toReactive;
        !isReadonly && $8c83eaf28779ff46$var$track(rawTarget, "iterate", $8c83eaf28779ff46$var$ITERATE_KEY);
        return target.forEach((value, key)=>{
            return callback.call(thisArg, wrap(value), wrap(key), observed);
        });
    };
}
function $8c83eaf28779ff46$var$createIterableMethod(method, isReadonly, isShallow) {
    return function(...args) {
        const target = this["__v_raw"];
        const rawTarget = $8c83eaf28779ff46$var$toRaw(target);
        const targetIsMap = $8c83eaf28779ff46$var$isMap(rawTarget);
        const isPair = method === "entries" || method === Symbol.iterator && targetIsMap;
        const isKeyOnly = method === "keys" && targetIsMap;
        const innerIterator = target[method](...args);
        const wrap = isShallow ? $8c83eaf28779ff46$var$toShallow : isReadonly ? $8c83eaf28779ff46$var$toReadonly : $8c83eaf28779ff46$var$toReactive;
        !isReadonly && $8c83eaf28779ff46$var$track(rawTarget, "iterate", isKeyOnly ? $8c83eaf28779ff46$var$MAP_KEY_ITERATE_KEY : $8c83eaf28779ff46$var$ITERATE_KEY);
        return {
            // iterator protocol
            next () {
                const { value: value, done: done } = innerIterator.next();
                return done ? {
                    value: value,
                    done: done
                } : {
                    value: isPair ? [
                        wrap(value[0]),
                        wrap(value[1])
                    ] : wrap(value),
                    done: done
                };
            },
            // iterable protocol
            [Symbol.iterator] () {
                return this;
            }
        };
    };
}
function $8c83eaf28779ff46$var$createReadonlyMethod(type) {
    return function(...args) {
        {
            const key = args[0] ? `on key "${args[0]}" ` : ``;
            console.warn(`${$8c83eaf28779ff46$var$capitalize(type)} operation ${key}failed: target is readonly.`, $8c83eaf28779ff46$var$toRaw(this));
        }
        return type === "delete" ? false : this;
    };
}
function $8c83eaf28779ff46$var$createInstrumentations() {
    const mutableInstrumentations2 = {
        get (key) {
            return $8c83eaf28779ff46$var$get$1(this, key);
        },
        get size () {
            return $8c83eaf28779ff46$var$size(this);
        },
        has: $8c83eaf28779ff46$var$has$1,
        add: $8c83eaf28779ff46$var$add,
        set: $8c83eaf28779ff46$var$set$1,
        delete: $8c83eaf28779ff46$var$deleteEntry,
        clear: $8c83eaf28779ff46$var$clear,
        forEach: $8c83eaf28779ff46$var$createForEach(false, false)
    };
    const shallowInstrumentations2 = {
        get (key) {
            return $8c83eaf28779ff46$var$get$1(this, key, false, true);
        },
        get size () {
            return $8c83eaf28779ff46$var$size(this);
        },
        has: $8c83eaf28779ff46$var$has$1,
        add: $8c83eaf28779ff46$var$add,
        set: $8c83eaf28779ff46$var$set$1,
        delete: $8c83eaf28779ff46$var$deleteEntry,
        clear: $8c83eaf28779ff46$var$clear,
        forEach: $8c83eaf28779ff46$var$createForEach(false, true)
    };
    const readonlyInstrumentations2 = {
        get (key) {
            return $8c83eaf28779ff46$var$get$1(this, key, true);
        },
        get size () {
            return $8c83eaf28779ff46$var$size(this, true);
        },
        has (key) {
            return $8c83eaf28779ff46$var$has$1.call(this, key, true);
        },
        add: $8c83eaf28779ff46$var$createReadonlyMethod("add"),
        set: $8c83eaf28779ff46$var$createReadonlyMethod("set"),
        delete: $8c83eaf28779ff46$var$createReadonlyMethod("delete"),
        clear: $8c83eaf28779ff46$var$createReadonlyMethod("clear"),
        forEach: $8c83eaf28779ff46$var$createForEach(true, false)
    };
    const shallowReadonlyInstrumentations2 = {
        get (key) {
            return $8c83eaf28779ff46$var$get$1(this, key, true, true);
        },
        get size () {
            return $8c83eaf28779ff46$var$size(this, true);
        },
        has (key) {
            return $8c83eaf28779ff46$var$has$1.call(this, key, true);
        },
        add: $8c83eaf28779ff46$var$createReadonlyMethod("add"),
        set: $8c83eaf28779ff46$var$createReadonlyMethod("set"),
        delete: $8c83eaf28779ff46$var$createReadonlyMethod("delete"),
        clear: $8c83eaf28779ff46$var$createReadonlyMethod("clear"),
        forEach: $8c83eaf28779ff46$var$createForEach(true, true)
    };
    const iteratorMethods = [
        "keys",
        "values",
        "entries",
        Symbol.iterator
    ];
    iteratorMethods.forEach((method)=>{
        mutableInstrumentations2[method] = $8c83eaf28779ff46$var$createIterableMethod(method, false, false);
        readonlyInstrumentations2[method] = $8c83eaf28779ff46$var$createIterableMethod(method, true, false);
        shallowInstrumentations2[method] = $8c83eaf28779ff46$var$createIterableMethod(method, false, true);
        shallowReadonlyInstrumentations2[method] = $8c83eaf28779ff46$var$createIterableMethod(method, true, true);
    });
    return [
        mutableInstrumentations2,
        readonlyInstrumentations2,
        shallowInstrumentations2,
        shallowReadonlyInstrumentations2
    ];
}
var [$8c83eaf28779ff46$var$mutableInstrumentations, $8c83eaf28779ff46$var$readonlyInstrumentations, $8c83eaf28779ff46$var$shallowInstrumentations, $8c83eaf28779ff46$var$shallowReadonlyInstrumentations] = /* @__PURE__ */ $8c83eaf28779ff46$var$createInstrumentations();
function $8c83eaf28779ff46$var$createInstrumentationGetter(isReadonly, shallow) {
    const instrumentations = shallow ? isReadonly ? $8c83eaf28779ff46$var$shallowReadonlyInstrumentations : $8c83eaf28779ff46$var$shallowInstrumentations : isReadonly ? $8c83eaf28779ff46$var$readonlyInstrumentations : $8c83eaf28779ff46$var$mutableInstrumentations;
    return (target, key, receiver)=>{
        if (key === "__v_isReactive") return !isReadonly;
        else if (key === "__v_isReadonly") return isReadonly;
        else if (key === "__v_raw") return target;
        return Reflect.get($8c83eaf28779ff46$var$hasOwn(instrumentations, key) && key in target ? instrumentations : target, key, receiver);
    };
}
var $8c83eaf28779ff46$var$mutableCollectionHandlers = {
    get: /* @__PURE__ */ $8c83eaf28779ff46$var$createInstrumentationGetter(false, false)
};
var $8c83eaf28779ff46$var$readonlyCollectionHandlers = {
    get: /* @__PURE__ */ $8c83eaf28779ff46$var$createInstrumentationGetter(true, false)
};
function $8c83eaf28779ff46$var$checkIdentityKeys(target, has2, key) {
    const rawKey = $8c83eaf28779ff46$var$toRaw(key);
    if (rawKey !== key && has2.call(target, rawKey)) {
        const type = $8c83eaf28779ff46$var$toRawType(target);
        console.warn(`Reactive ${type} contains both the raw and reactive versions of the same object${type === `Map` ? ` as keys` : ``}, which can lead to inconsistencies. Avoid differentiating between the raw and reactive versions of an object and only use the reactive version if possible.`);
    }
}
var $8c83eaf28779ff46$var$reactiveMap = /* @__PURE__ */ new WeakMap();
var $8c83eaf28779ff46$var$shallowReactiveMap = /* @__PURE__ */ new WeakMap();
var $8c83eaf28779ff46$var$readonlyMap = /* @__PURE__ */ new WeakMap();
var $8c83eaf28779ff46$var$shallowReadonlyMap = /* @__PURE__ */ new WeakMap();
function $8c83eaf28779ff46$var$targetTypeMap(rawType) {
    switch(rawType){
        case "Object":
        case "Array":
            return 1;
        case "Map":
        case "Set":
        case "WeakMap":
        case "WeakSet":
            return 2;
        default:
            return 0;
    }
}
function $8c83eaf28779ff46$var$getTargetType(value) {
    return value["__v_skip"] || !Object.isExtensible(value) ? 0 : $8c83eaf28779ff46$var$targetTypeMap($8c83eaf28779ff46$var$toRawType(value));
}
function $8c83eaf28779ff46$var$reactive2(target) {
    if (target && target["__v_isReadonly"]) return target;
    return $8c83eaf28779ff46$var$createReactiveObject(target, false, $8c83eaf28779ff46$var$mutableHandlers, $8c83eaf28779ff46$var$mutableCollectionHandlers, $8c83eaf28779ff46$var$reactiveMap);
}
function $8c83eaf28779ff46$var$readonly(target) {
    return $8c83eaf28779ff46$var$createReactiveObject(target, true, $8c83eaf28779ff46$var$readonlyHandlers, $8c83eaf28779ff46$var$readonlyCollectionHandlers, $8c83eaf28779ff46$var$readonlyMap);
}
function $8c83eaf28779ff46$var$createReactiveObject(target, isReadonly, baseHandlers, collectionHandlers, proxyMap) {
    if (!$8c83eaf28779ff46$var$isObject(target)) {
        console.warn(`value cannot be made reactive: ${String(target)}`);
        return target;
    }
    if (target["__v_raw"] && !(isReadonly && target["__v_isReactive"])) return target;
    const existingProxy = proxyMap.get(target);
    if (existingProxy) return existingProxy;
    const targetType = $8c83eaf28779ff46$var$getTargetType(target);
    if (targetType === 0) return target;
    const proxy = new Proxy(target, targetType === 2 ? collectionHandlers : baseHandlers);
    proxyMap.set(target, proxy);
    return proxy;
}
function $8c83eaf28779ff46$var$toRaw(observed) {
    return observed && $8c83eaf28779ff46$var$toRaw(observed["__v_raw"]) || observed;
}
function $8c83eaf28779ff46$var$isRef(r) {
    return Boolean(r && r.__v_isRef === true);
}
// packages/alpinejs/src/magics/$nextTick.js
$8c83eaf28779ff46$var$magic("nextTick", ()=>$8c83eaf28779ff46$var$nextTick);
// packages/alpinejs/src/magics/$dispatch.js
$8c83eaf28779ff46$var$magic("dispatch", (el)=>$8c83eaf28779ff46$var$dispatch.bind($8c83eaf28779ff46$var$dispatch, el));
// packages/alpinejs/src/magics/$watch.js
$8c83eaf28779ff46$var$magic("watch", (el, { evaluateLater: evaluateLater2, cleanup: cleanup2 })=>(key, callback)=>{
        let evaluate2 = evaluateLater2(key);
        let getter = ()=>{
            let value;
            evaluate2((i)=>value = i);
            return value;
        };
        let unwatch = $8c83eaf28779ff46$var$watch(getter, callback);
        cleanup2(unwatch);
    });
// packages/alpinejs/src/magics/$store.js
$8c83eaf28779ff46$var$magic("store", $8c83eaf28779ff46$var$getStores);
// packages/alpinejs/src/magics/$data.js
$8c83eaf28779ff46$var$magic("data", (el)=>$8c83eaf28779ff46$var$scope(el));
// packages/alpinejs/src/magics/$root.js
$8c83eaf28779ff46$var$magic("root", (el)=>$8c83eaf28779ff46$var$closestRoot(el));
// packages/alpinejs/src/magics/$refs.js
$8c83eaf28779ff46$var$magic("refs", (el)=>{
    if (el._x_refs_proxy) return el._x_refs_proxy;
    el._x_refs_proxy = $8c83eaf28779ff46$var$mergeProxies($8c83eaf28779ff46$var$getArrayOfRefObject(el));
    return el._x_refs_proxy;
});
function $8c83eaf28779ff46$var$getArrayOfRefObject(el) {
    let refObjects = [];
    $8c83eaf28779ff46$var$findClosest(el, (i)=>{
        if (i._x_refs) refObjects.push(i._x_refs);
    });
    return refObjects;
}
// packages/alpinejs/src/ids.js
var $8c83eaf28779ff46$var$globalIdMemo = {};
function $8c83eaf28779ff46$var$findAndIncrementId(name) {
    if (!$8c83eaf28779ff46$var$globalIdMemo[name]) $8c83eaf28779ff46$var$globalIdMemo[name] = 0;
    return ++$8c83eaf28779ff46$var$globalIdMemo[name];
}
function $8c83eaf28779ff46$var$closestIdRoot(el, name) {
    return $8c83eaf28779ff46$var$findClosest(el, (element)=>{
        if (element._x_ids && element._x_ids[name]) return true;
    });
}
function $8c83eaf28779ff46$var$setIdRoot(el, name) {
    if (!el._x_ids) el._x_ids = {};
    if (!el._x_ids[name]) el._x_ids[name] = $8c83eaf28779ff46$var$findAndIncrementId(name);
}
// packages/alpinejs/src/magics/$id.js
$8c83eaf28779ff46$var$magic("id", (el, { cleanup: cleanup2 })=>(name, key = null)=>{
        let cacheKey = `${name}${key ? `-${key}` : ""}`;
        return $8c83eaf28779ff46$var$cacheIdByNameOnElement(el, cacheKey, cleanup2, ()=>{
            let root = $8c83eaf28779ff46$var$closestIdRoot(el, name);
            let id = root ? root._x_ids[name] : $8c83eaf28779ff46$var$findAndIncrementId(name);
            return key ? `${name}-${id}-${key}` : `${name}-${id}`;
        });
    });
$8c83eaf28779ff46$var$interceptClone((from, to)=>{
    if (from._x_id) to._x_id = from._x_id;
});
function $8c83eaf28779ff46$var$cacheIdByNameOnElement(el, cacheKey, cleanup2, callback) {
    if (!el._x_id) el._x_id = {};
    if (el._x_id[cacheKey]) return el._x_id[cacheKey];
    let output = callback();
    el._x_id[cacheKey] = output;
    cleanup2(()=>{
        delete el._x_id[cacheKey];
    });
    return output;
}
// packages/alpinejs/src/magics/$el.js
$8c83eaf28779ff46$var$magic("el", (el)=>el);
// packages/alpinejs/src/magics/index.js
$8c83eaf28779ff46$var$warnMissingPluginMagic("Focus", "focus", "focus");
$8c83eaf28779ff46$var$warnMissingPluginMagic("Persist", "persist", "persist");
function $8c83eaf28779ff46$var$warnMissingPluginMagic(name, magicName, slug) {
    $8c83eaf28779ff46$var$magic(magicName, (el)=>$8c83eaf28779ff46$var$warn(`You can't use [$${magicName}] without first installing the "${name}" plugin here: https://alpinejs.dev/plugins/${slug}`, el));
}
// packages/alpinejs/src/directives/x-modelable.js
$8c83eaf28779ff46$var$directive("modelable", (el, { expression: expression }, { effect: effect3, evaluateLater: evaluateLater2, cleanup: cleanup2 })=>{
    let func = evaluateLater2(expression);
    let innerGet = ()=>{
        let result;
        func((i)=>result = i);
        return result;
    };
    let evaluateInnerSet = evaluateLater2(`${expression} = __placeholder`);
    let innerSet = (val)=>evaluateInnerSet(()=>{}, {
            scope: {
                "__placeholder": val
            }
        });
    let initialValue = innerGet();
    innerSet(initialValue);
    queueMicrotask(()=>{
        if (!el._x_model) return;
        el._x_removeModelListeners["default"]();
        let outerGet = el._x_model.get;
        let outerSet = el._x_model.set;
        let releaseEntanglement = $8c83eaf28779ff46$var$entangle({
            get () {
                return outerGet();
            },
            set (value) {
                outerSet(value);
            }
        }, {
            get () {
                return innerGet();
            },
            set (value) {
                innerSet(value);
            }
        });
        cleanup2(releaseEntanglement);
    });
});
// packages/alpinejs/src/directives/x-teleport.js
$8c83eaf28779ff46$var$directive("teleport", (el, { modifiers: modifiers, expression: expression }, { cleanup: cleanup2 })=>{
    if (el.tagName.toLowerCase() !== "template") $8c83eaf28779ff46$var$warn("x-teleport can only be used on a <template> tag", el);
    let target = $8c83eaf28779ff46$var$getTarget(expression);
    let clone2 = el.content.cloneNode(true).firstElementChild;
    el._x_teleport = clone2;
    clone2._x_teleportBack = el;
    el.setAttribute("data-teleport-template", true);
    clone2.setAttribute("data-teleport-target", true);
    if (el._x_forwardEvents) el._x_forwardEvents.forEach((eventName)=>{
        clone2.addEventListener(eventName, (e)=>{
            e.stopPropagation();
            el.dispatchEvent(new e.constructor(e.type, e));
        });
    });
    $8c83eaf28779ff46$var$addScopeToNode(clone2, {}, el);
    let placeInDom = (clone3, target2, modifiers2)=>{
        if (modifiers2.includes("prepend")) target2.parentNode.insertBefore(clone3, target2);
        else if (modifiers2.includes("append")) target2.parentNode.insertBefore(clone3, target2.nextSibling);
        else target2.appendChild(clone3);
    };
    $8c83eaf28779ff46$var$mutateDom(()=>{
        placeInDom(clone2, target, modifiers);
        $8c83eaf28779ff46$var$skipDuringClone(()=>{
            $8c83eaf28779ff46$var$initTree(clone2);
        })();
    });
    el._x_teleportPutBack = ()=>{
        let target2 = $8c83eaf28779ff46$var$getTarget(expression);
        $8c83eaf28779ff46$var$mutateDom(()=>{
            placeInDom(el._x_teleport, target2, modifiers);
        });
    };
    cleanup2(()=>$8c83eaf28779ff46$var$mutateDom(()=>{
            clone2.remove();
            $8c83eaf28779ff46$var$destroyTree(clone2);
        }));
});
var $8c83eaf28779ff46$var$teleportContainerDuringClone = document.createElement("div");
function $8c83eaf28779ff46$var$getTarget(expression) {
    let target = $8c83eaf28779ff46$var$skipDuringClone(()=>{
        return document.querySelector(expression);
    }, ()=>{
        return $8c83eaf28779ff46$var$teleportContainerDuringClone;
    })();
    if (!target) $8c83eaf28779ff46$var$warn(`Cannot find x-teleport element for selector: "${expression}"`);
    return target;
}
// packages/alpinejs/src/directives/x-ignore.js
var $8c83eaf28779ff46$var$handler = ()=>{};
$8c83eaf28779ff46$var$handler.inline = (el, { modifiers: modifiers }, { cleanup: cleanup2 })=>{
    modifiers.includes("self") ? el._x_ignoreSelf = true : el._x_ignore = true;
    cleanup2(()=>{
        modifiers.includes("self") ? delete el._x_ignoreSelf : delete el._x_ignore;
    });
};
$8c83eaf28779ff46$var$directive("ignore", $8c83eaf28779ff46$var$handler);
// packages/alpinejs/src/directives/x-effect.js
$8c83eaf28779ff46$var$directive("effect", $8c83eaf28779ff46$var$skipDuringClone((el, { expression: expression }, { effect: effect3 })=>{
    effect3($8c83eaf28779ff46$var$evaluateLater(el, expression));
}));
// packages/alpinejs/src/utils/on.js
function $8c83eaf28779ff46$var$on(el, event, modifiers, callback) {
    let listenerTarget = el;
    let handler4 = (e)=>callback(e);
    let options = {};
    let wrapHandler = (callback2, wrapper)=>(e)=>wrapper(callback2, e);
    if (modifiers.includes("dot")) event = $8c83eaf28779ff46$var$dotSyntax(event);
    if (modifiers.includes("camel")) event = $8c83eaf28779ff46$var$camelCase2(event);
    if (modifiers.includes("passive")) options.passive = true;
    if (modifiers.includes("capture")) options.capture = true;
    if (modifiers.includes("window")) listenerTarget = window;
    if (modifiers.includes("document")) listenerTarget = document;
    if (modifiers.includes("debounce")) {
        let nextModifier = modifiers[modifiers.indexOf("debounce") + 1] || "invalid-wait";
        let wait = $8c83eaf28779ff46$var$isNumeric(nextModifier.split("ms")[0]) ? Number(nextModifier.split("ms")[0]) : 250;
        handler4 = $8c83eaf28779ff46$var$debounce(handler4, wait);
    }
    if (modifiers.includes("throttle")) {
        let nextModifier = modifiers[modifiers.indexOf("throttle") + 1] || "invalid-wait";
        let wait = $8c83eaf28779ff46$var$isNumeric(nextModifier.split("ms")[0]) ? Number(nextModifier.split("ms")[0]) : 250;
        handler4 = $8c83eaf28779ff46$var$throttle(handler4, wait);
    }
    if (modifiers.includes("prevent")) handler4 = wrapHandler(handler4, (next, e)=>{
        e.preventDefault();
        next(e);
    });
    if (modifiers.includes("stop")) handler4 = wrapHandler(handler4, (next, e)=>{
        e.stopPropagation();
        next(e);
    });
    if (modifiers.includes("once")) handler4 = wrapHandler(handler4, (next, e)=>{
        next(e);
        listenerTarget.removeEventListener(event, handler4, options);
    });
    if (modifiers.includes("away") || modifiers.includes("outside")) {
        listenerTarget = document;
        handler4 = wrapHandler(handler4, (next, e)=>{
            if (el.contains(e.target)) return;
            if (e.target.isConnected === false) return;
            if (el.offsetWidth < 1 && el.offsetHeight < 1) return;
            if (el._x_isShown === false) return;
            next(e);
        });
    }
    if (modifiers.includes("self")) handler4 = wrapHandler(handler4, (next, e)=>{
        e.target === el && next(e);
    });
    if ($8c83eaf28779ff46$var$isKeyEvent(event) || $8c83eaf28779ff46$var$isClickEvent(event)) handler4 = wrapHandler(handler4, (next, e)=>{
        if ($8c83eaf28779ff46$var$isListeningForASpecificKeyThatHasntBeenPressed(e, modifiers)) return;
        next(e);
    });
    listenerTarget.addEventListener(event, handler4, options);
    return ()=>{
        listenerTarget.removeEventListener(event, handler4, options);
    };
}
function $8c83eaf28779ff46$var$dotSyntax(subject) {
    return subject.replace(/-/g, ".");
}
function $8c83eaf28779ff46$var$camelCase2(subject) {
    return subject.toLowerCase().replace(/-(\w)/g, (match, char)=>char.toUpperCase());
}
function $8c83eaf28779ff46$var$isNumeric(subject) {
    return !Array.isArray(subject) && !isNaN(subject);
}
function $8c83eaf28779ff46$var$kebabCase2(subject) {
    if ([
        " ",
        "_"
    ].includes(subject)) return subject;
    return subject.replace(/([a-z])([A-Z])/g, "$1-$2").replace(/[_\s]/, "-").toLowerCase();
}
function $8c83eaf28779ff46$var$isKeyEvent(event) {
    return [
        "keydown",
        "keyup"
    ].includes(event);
}
function $8c83eaf28779ff46$var$isClickEvent(event) {
    return [
        "contextmenu",
        "click",
        "mouse"
    ].some((i)=>event.includes(i));
}
function $8c83eaf28779ff46$var$isListeningForASpecificKeyThatHasntBeenPressed(e, modifiers) {
    let keyModifiers = modifiers.filter((i)=>{
        return ![
            "window",
            "document",
            "prevent",
            "stop",
            "once",
            "capture",
            "self",
            "away",
            "outside",
            "passive",
            "preserve-scroll"
        ].includes(i);
    });
    if (keyModifiers.includes("debounce")) {
        let debounceIndex = keyModifiers.indexOf("debounce");
        keyModifiers.splice(debounceIndex, $8c83eaf28779ff46$var$isNumeric((keyModifiers[debounceIndex + 1] || "invalid-wait").split("ms")[0]) ? 2 : 1);
    }
    if (keyModifiers.includes("throttle")) {
        let debounceIndex = keyModifiers.indexOf("throttle");
        keyModifiers.splice(debounceIndex, $8c83eaf28779ff46$var$isNumeric((keyModifiers[debounceIndex + 1] || "invalid-wait").split("ms")[0]) ? 2 : 1);
    }
    if (keyModifiers.length === 0) return false;
    if (keyModifiers.length === 1 && $8c83eaf28779ff46$var$keyToModifiers(e.key).includes(keyModifiers[0])) return false;
    const systemKeyModifiers = [
        "ctrl",
        "shift",
        "alt",
        "meta",
        "cmd",
        "super"
    ];
    const selectedSystemKeyModifiers = systemKeyModifiers.filter((modifier)=>keyModifiers.includes(modifier));
    keyModifiers = keyModifiers.filter((i)=>!selectedSystemKeyModifiers.includes(i));
    if (selectedSystemKeyModifiers.length > 0) {
        const activelyPressedKeyModifiers = selectedSystemKeyModifiers.filter((modifier)=>{
            if (modifier === "cmd" || modifier === "super") modifier = "meta";
            return e[`${modifier}Key`];
        });
        if (activelyPressedKeyModifiers.length === selectedSystemKeyModifiers.length) {
            if ($8c83eaf28779ff46$var$isClickEvent(e.type)) return false;
            if ($8c83eaf28779ff46$var$keyToModifiers(e.key).includes(keyModifiers[0])) return false;
        }
    }
    return true;
}
function $8c83eaf28779ff46$var$keyToModifiers(key) {
    if (!key) return [];
    key = $8c83eaf28779ff46$var$kebabCase2(key);
    let modifierToKeyMap = {
        "ctrl": "control",
        "slash": "/",
        "space": " ",
        "spacebar": " ",
        "cmd": "meta",
        "esc": "escape",
        "up": "arrow-up",
        "down": "arrow-down",
        "left": "arrow-left",
        "right": "arrow-right",
        "period": ".",
        "comma": ",",
        "equal": "=",
        "minus": "-",
        "underscore": "_"
    };
    modifierToKeyMap[key] = key;
    return Object.keys(modifierToKeyMap).map((modifier)=>{
        if (modifierToKeyMap[modifier] === key) return modifier;
    }).filter((modifier)=>modifier);
}
// packages/alpinejs/src/directives/x-model.js
$8c83eaf28779ff46$var$directive("model", (el, { modifiers: modifiers, expression: expression }, { effect: effect3, cleanup: cleanup2 })=>{
    let scopeTarget = el;
    if (modifiers.includes("parent")) scopeTarget = el.parentNode;
    let evaluateGet = $8c83eaf28779ff46$var$evaluateLater(scopeTarget, expression);
    let evaluateSet;
    if (typeof expression === "string") evaluateSet = $8c83eaf28779ff46$var$evaluateLater(scopeTarget, `${expression} = __placeholder`);
    else if (typeof expression === "function" && typeof expression() === "string") evaluateSet = $8c83eaf28779ff46$var$evaluateLater(scopeTarget, `${expression()} = __placeholder`);
    else evaluateSet = ()=>{};
    let getValue = ()=>{
        let result;
        evaluateGet((value)=>result = value);
        return $8c83eaf28779ff46$var$isGetterSetter(result) ? result.get() : result;
    };
    let setValue = (value)=>{
        let result;
        evaluateGet((value2)=>result = value2);
        if ($8c83eaf28779ff46$var$isGetterSetter(result)) result.set(value);
        else evaluateSet(()=>{}, {
            scope: {
                "__placeholder": value
            }
        });
    };
    if (typeof expression === "string" && el.type === "radio") $8c83eaf28779ff46$var$mutateDom(()=>{
        if (!el.hasAttribute("name")) el.setAttribute("name", expression);
    });
    let event = el.tagName.toLowerCase() === "select" || [
        "checkbox",
        "radio"
    ].includes(el.type) || modifiers.includes("lazy") ? "change" : "input";
    let removeListener = $8c83eaf28779ff46$var$isCloning ? ()=>{} : $8c83eaf28779ff46$var$on(el, event, modifiers, (e)=>{
        setValue($8c83eaf28779ff46$var$getInputValue(el, modifiers, e, getValue()));
    });
    if (modifiers.includes("fill")) {
        if ([
            void 0,
            null,
            ""
        ].includes(getValue()) || $8c83eaf28779ff46$var$isCheckbox(el) && Array.isArray(getValue()) || el.tagName.toLowerCase() === "select" && el.multiple) setValue($8c83eaf28779ff46$var$getInputValue(el, modifiers, {
            target: el
        }, getValue()));
    }
    if (!el._x_removeModelListeners) el._x_removeModelListeners = {};
    el._x_removeModelListeners["default"] = removeListener;
    cleanup2(()=>el._x_removeModelListeners["default"]());
    if (el.form) {
        let removeResetListener = $8c83eaf28779ff46$var$on(el.form, "reset", [], (e)=>{
            $8c83eaf28779ff46$var$nextTick(()=>el._x_model && el._x_model.set($8c83eaf28779ff46$var$getInputValue(el, modifiers, {
                    target: el
                }, getValue())));
        });
        cleanup2(()=>removeResetListener());
    }
    el._x_model = {
        get () {
            return getValue();
        },
        set (value) {
            setValue(value);
        }
    };
    el._x_forceModelUpdate = (value)=>{
        if (value === void 0 && typeof expression === "string" && expression.match(/\./)) value = "";
        window.fromModel = true;
        $8c83eaf28779ff46$var$mutateDom(()=>$8c83eaf28779ff46$var$bind(el, "value", value));
        delete window.fromModel;
    };
    effect3(()=>{
        let value = getValue();
        if (modifiers.includes("unintrusive") && document.activeElement.isSameNode(el)) return;
        el._x_forceModelUpdate(value);
    });
});
function $8c83eaf28779ff46$var$getInputValue(el, modifiers, event, currentValue) {
    return $8c83eaf28779ff46$var$mutateDom(()=>{
        if (event instanceof CustomEvent && event.detail !== void 0) return event.detail !== null && event.detail !== void 0 ? event.detail : event.target.value;
        else if ($8c83eaf28779ff46$var$isCheckbox(el)) {
            if (Array.isArray(currentValue)) {
                let newValue = null;
                if (modifiers.includes("number")) newValue = $8c83eaf28779ff46$var$safeParseNumber(event.target.value);
                else if (modifiers.includes("boolean")) newValue = $8c83eaf28779ff46$var$safeParseBoolean(event.target.value);
                else newValue = event.target.value;
                return event.target.checked ? currentValue.includes(newValue) ? currentValue : currentValue.concat([
                    newValue
                ]) : currentValue.filter((el2)=>!$8c83eaf28779ff46$var$checkedAttrLooseCompare2(el2, newValue));
            } else return event.target.checked;
        } else if (el.tagName.toLowerCase() === "select" && el.multiple) {
            if (modifiers.includes("number")) return Array.from(event.target.selectedOptions).map((option)=>{
                let rawValue = option.value || option.text;
                return $8c83eaf28779ff46$var$safeParseNumber(rawValue);
            });
            else if (modifiers.includes("boolean")) return Array.from(event.target.selectedOptions).map((option)=>{
                let rawValue = option.value || option.text;
                return $8c83eaf28779ff46$var$safeParseBoolean(rawValue);
            });
            return Array.from(event.target.selectedOptions).map((option)=>{
                return option.value || option.text;
            });
        } else {
            let newValue;
            if ($8c83eaf28779ff46$var$isRadio(el)) {
                if (event.target.checked) newValue = event.target.value;
                else newValue = currentValue;
            } else newValue = event.target.value;
            if (modifiers.includes("number")) return $8c83eaf28779ff46$var$safeParseNumber(newValue);
            else if (modifiers.includes("boolean")) return $8c83eaf28779ff46$var$safeParseBoolean(newValue);
            else if (modifiers.includes("trim")) return newValue.trim();
            else return newValue;
        }
    });
}
function $8c83eaf28779ff46$var$safeParseNumber(rawValue) {
    let number = rawValue ? parseFloat(rawValue) : null;
    return $8c83eaf28779ff46$var$isNumeric2(number) ? number : rawValue;
}
function $8c83eaf28779ff46$var$checkedAttrLooseCompare2(valueA, valueB) {
    return valueA == valueB;
}
function $8c83eaf28779ff46$var$isNumeric2(subject) {
    return !Array.isArray(subject) && !isNaN(subject);
}
function $8c83eaf28779ff46$var$isGetterSetter(value) {
    return value !== null && typeof value === "object" && typeof value.get === "function" && typeof value.set === "function";
}
// packages/alpinejs/src/directives/x-cloak.js
$8c83eaf28779ff46$var$directive("cloak", (el)=>queueMicrotask(()=>$8c83eaf28779ff46$var$mutateDom(()=>el.removeAttribute($8c83eaf28779ff46$var$prefix("cloak")))));
// packages/alpinejs/src/directives/x-init.js
$8c83eaf28779ff46$var$addInitSelector(()=>`[${$8c83eaf28779ff46$var$prefix("init")}]`);
$8c83eaf28779ff46$var$directive("init", $8c83eaf28779ff46$var$skipDuringClone((el, { expression: expression }, { evaluate: evaluate2 })=>{
    if (typeof expression === "string") return !!expression.trim() && evaluate2(expression, {}, false);
    return evaluate2(expression, {}, false);
}));
// packages/alpinejs/src/directives/x-text.js
$8c83eaf28779ff46$var$directive("text", (el, { expression: expression }, { effect: effect3, evaluateLater: evaluateLater2 })=>{
    let evaluate2 = evaluateLater2(expression);
    effect3(()=>{
        evaluate2((value)=>{
            $8c83eaf28779ff46$var$mutateDom(()=>{
                el.textContent = value;
            });
        });
    });
});
// packages/alpinejs/src/directives/x-html.js
$8c83eaf28779ff46$var$directive("html", (el, { expression: expression }, { effect: effect3, evaluateLater: evaluateLater2 })=>{
    let evaluate2 = evaluateLater2(expression);
    effect3(()=>{
        evaluate2((value)=>{
            $8c83eaf28779ff46$var$mutateDom(()=>{
                el.innerHTML = value;
                el._x_ignoreSelf = true;
                $8c83eaf28779ff46$var$initTree(el);
                delete el._x_ignoreSelf;
            });
        });
    });
});
// packages/alpinejs/src/directives/x-bind.js
$8c83eaf28779ff46$var$mapAttributes($8c83eaf28779ff46$var$startingWith(":", $8c83eaf28779ff46$var$into($8c83eaf28779ff46$var$prefix("bind:"))));
var $8c83eaf28779ff46$var$handler2 = (el, { value: value, modifiers: modifiers, expression: expression, original: original }, { effect: effect3, cleanup: cleanup2 })=>{
    if (!value) {
        let bindingProviders = {};
        $8c83eaf28779ff46$var$injectBindingProviders(bindingProviders);
        let getBindings = $8c83eaf28779ff46$var$evaluateLater(el, expression);
        getBindings((bindings)=>{
            $8c83eaf28779ff46$var$applyBindingsObject(el, bindings, original);
        }, {
            scope: bindingProviders
        });
        return;
    }
    if (value === "key") return $8c83eaf28779ff46$var$storeKeyForXFor(el, expression);
    if (el._x_inlineBindings && el._x_inlineBindings[value] && el._x_inlineBindings[value].extract) return;
    let evaluate2 = $8c83eaf28779ff46$var$evaluateLater(el, expression);
    effect3(()=>evaluate2((result)=>{
            if (result === void 0 && typeof expression === "string" && expression.match(/\./)) result = "";
            $8c83eaf28779ff46$var$mutateDom(()=>$8c83eaf28779ff46$var$bind(el, value, result, modifiers));
        }));
    cleanup2(()=>{
        el._x_undoAddedClasses && el._x_undoAddedClasses();
        el._x_undoAddedStyles && el._x_undoAddedStyles();
    });
};
$8c83eaf28779ff46$var$handler2.inline = (el, { value: value, modifiers: modifiers, expression: expression })=>{
    if (!value) return;
    if (!el._x_inlineBindings) el._x_inlineBindings = {};
    el._x_inlineBindings[value] = {
        expression: expression,
        extract: false
    };
};
$8c83eaf28779ff46$var$directive("bind", $8c83eaf28779ff46$var$handler2);
function $8c83eaf28779ff46$var$storeKeyForXFor(el, expression) {
    el._x_keyExpression = expression;
}
// packages/alpinejs/src/directives/x-data.js
$8c83eaf28779ff46$var$addRootSelector(()=>`[${$8c83eaf28779ff46$var$prefix("data")}]`);
$8c83eaf28779ff46$var$directive("data", (el, { expression: expression }, { cleanup: cleanup2 })=>{
    if ($8c83eaf28779ff46$var$shouldSkipRegisteringDataDuringClone(el)) return;
    expression = expression === "" ? "{}" : expression;
    let magicContext = {};
    $8c83eaf28779ff46$var$injectMagics(magicContext, el);
    let dataProviderContext = {};
    $8c83eaf28779ff46$var$injectDataProviders(dataProviderContext, magicContext);
    let data2 = $8c83eaf28779ff46$var$evaluate(el, expression, {
        scope: dataProviderContext
    });
    if (data2 === void 0 || data2 === true) data2 = {};
    $8c83eaf28779ff46$var$injectMagics(data2, el);
    let reactiveData = $8c83eaf28779ff46$var$reactive(data2);
    $8c83eaf28779ff46$var$initInterceptors(reactiveData);
    let undo = $8c83eaf28779ff46$var$addScopeToNode(el, reactiveData);
    reactiveData["init"] && $8c83eaf28779ff46$var$evaluate(el, reactiveData["init"]);
    cleanup2(()=>{
        reactiveData["destroy"] && $8c83eaf28779ff46$var$evaluate(el, reactiveData["destroy"]);
        undo();
    });
});
$8c83eaf28779ff46$var$interceptClone((from, to)=>{
    if (from._x_dataStack) {
        to._x_dataStack = from._x_dataStack;
        to.setAttribute("data-has-alpine-state", true);
    }
});
function $8c83eaf28779ff46$var$shouldSkipRegisteringDataDuringClone(el) {
    if (!$8c83eaf28779ff46$var$isCloning) return false;
    if ($8c83eaf28779ff46$var$isCloningLegacy) return true;
    return el.hasAttribute("data-has-alpine-state");
}
// packages/alpinejs/src/directives/x-show.js
$8c83eaf28779ff46$var$directive("show", (el, { modifiers: modifiers, expression: expression }, { effect: effect3 })=>{
    let evaluate2 = $8c83eaf28779ff46$var$evaluateLater(el, expression);
    if (!el._x_doHide) el._x_doHide = ()=>{
        $8c83eaf28779ff46$var$mutateDom(()=>{
            el.style.setProperty("display", "none", modifiers.includes("important") ? "important" : void 0);
        });
    };
    if (!el._x_doShow) el._x_doShow = ()=>{
        $8c83eaf28779ff46$var$mutateDom(()=>{
            if (el.style.length === 1 && el.style.display === "none") el.removeAttribute("style");
            else el.style.removeProperty("display");
        });
    };
    let hide = ()=>{
        el._x_doHide();
        el._x_isShown = false;
    };
    let show = ()=>{
        el._x_doShow();
        el._x_isShown = true;
    };
    let clickAwayCompatibleShow = ()=>setTimeout(show);
    let toggle = $8c83eaf28779ff46$var$once((value)=>value ? show() : hide(), (value)=>{
        if (typeof el._x_toggleAndCascadeWithTransitions === "function") el._x_toggleAndCascadeWithTransitions(el, value, show, hide);
        else value ? clickAwayCompatibleShow() : hide();
    });
    let oldValue;
    let firstTime = true;
    effect3(()=>evaluate2((value)=>{
            if (!firstTime && value === oldValue) return;
            if (modifiers.includes("immediate")) value ? clickAwayCompatibleShow() : hide();
            toggle(value);
            oldValue = value;
            firstTime = false;
        }));
});
// packages/alpinejs/src/directives/x-for.js
$8c83eaf28779ff46$var$directive("for", (el, { expression: expression }, { effect: effect3, cleanup: cleanup2 })=>{
    let iteratorNames = $8c83eaf28779ff46$var$parseForExpression(expression);
    let evaluateItems = $8c83eaf28779ff46$var$evaluateLater(el, iteratorNames.items);
    let evaluateKey = $8c83eaf28779ff46$var$evaluateLater(el, // the x-bind:key expression is stored for our use instead of evaluated.
    el._x_keyExpression || "index");
    el._x_prevKeys = [];
    el._x_lookup = {};
    effect3(()=>$8c83eaf28779ff46$var$loop(el, iteratorNames, evaluateItems, evaluateKey));
    cleanup2(()=>{
        Object.values(el._x_lookup).forEach((el2)=>$8c83eaf28779ff46$var$mutateDom(()=>{
                $8c83eaf28779ff46$var$destroyTree(el2);
                el2.remove();
            }));
        delete el._x_prevKeys;
        delete el._x_lookup;
    });
});
function $8c83eaf28779ff46$var$loop(el, iteratorNames, evaluateItems, evaluateKey) {
    let isObject2 = (i)=>typeof i === "object" && !Array.isArray(i);
    let templateEl = el;
    evaluateItems((items)=>{
        if ($8c83eaf28779ff46$var$isNumeric3(items) && items >= 0) items = Array.from(Array(items).keys(), (i)=>i + 1);
        if (items === void 0) items = [];
        let lookup = el._x_lookup;
        let prevKeys = el._x_prevKeys;
        let scopes = [];
        let keys = [];
        if (isObject2(items)) items = Object.entries(items).map(([key, value])=>{
            let scope2 = $8c83eaf28779ff46$var$getIterationScopeVariables(iteratorNames, value, key, items);
            evaluateKey((value2)=>{
                if (keys.includes(value2)) $8c83eaf28779ff46$var$warn("Duplicate key on x-for", el);
                keys.push(value2);
            }, {
                scope: {
                    index: key,
                    ...scope2
                }
            });
            scopes.push(scope2);
        });
        else for(let i = 0; i < items.length; i++){
            let scope2 = $8c83eaf28779ff46$var$getIterationScopeVariables(iteratorNames, items[i], i, items);
            evaluateKey((value)=>{
                if (keys.includes(value)) $8c83eaf28779ff46$var$warn("Duplicate key on x-for", el);
                keys.push(value);
            }, {
                scope: {
                    index: i,
                    ...scope2
                }
            });
            scopes.push(scope2);
        }
        let adds = [];
        let moves = [];
        let removes = [];
        let sames = [];
        for(let i = 0; i < prevKeys.length; i++){
            let key = prevKeys[i];
            if (keys.indexOf(key) === -1) removes.push(key);
        }
        prevKeys = prevKeys.filter((key)=>!removes.includes(key));
        let lastKey = "template";
        for(let i = 0; i < keys.length; i++){
            let key = keys[i];
            let prevIndex = prevKeys.indexOf(key);
            if (prevIndex === -1) {
                prevKeys.splice(i, 0, key);
                adds.push([
                    lastKey,
                    i
                ]);
            } else if (prevIndex !== i) {
                let keyInSpot = prevKeys.splice(i, 1)[0];
                let keyForSpot = prevKeys.splice(prevIndex - 1, 1)[0];
                prevKeys.splice(i, 0, keyForSpot);
                prevKeys.splice(prevIndex, 0, keyInSpot);
                moves.push([
                    keyInSpot,
                    keyForSpot
                ]);
            } else sames.push(key);
            lastKey = key;
        }
        for(let i = 0; i < removes.length; i++){
            let key = removes[i];
            if (!(key in lookup)) continue;
            $8c83eaf28779ff46$var$mutateDom(()=>{
                $8c83eaf28779ff46$var$destroyTree(lookup[key]);
                lookup[key].remove();
            });
            delete lookup[key];
        }
        for(let i = 0; i < moves.length; i++){
            let [keyInSpot, keyForSpot] = moves[i];
            let elInSpot = lookup[keyInSpot];
            let elForSpot = lookup[keyForSpot];
            let marker = document.createElement("div");
            $8c83eaf28779ff46$var$mutateDom(()=>{
                if (!elForSpot) $8c83eaf28779ff46$var$warn(`x-for ":key" is undefined or invalid`, templateEl, keyForSpot, lookup);
                elForSpot.after(marker);
                elInSpot.after(elForSpot);
                elForSpot._x_currentIfEl && elForSpot.after(elForSpot._x_currentIfEl);
                marker.before(elInSpot);
                elInSpot._x_currentIfEl && elInSpot.after(elInSpot._x_currentIfEl);
                marker.remove();
            });
            elForSpot._x_refreshXForScope(scopes[keys.indexOf(keyForSpot)]);
        }
        for(let i = 0; i < adds.length; i++){
            let [lastKey2, index] = adds[i];
            let lastEl = lastKey2 === "template" ? templateEl : lookup[lastKey2];
            if (lastEl._x_currentIfEl) lastEl = lastEl._x_currentIfEl;
            let scope2 = scopes[index];
            let key = keys[index];
            let clone2 = document.importNode(templateEl.content, true).firstElementChild;
            let reactiveScope = $8c83eaf28779ff46$var$reactive(scope2);
            $8c83eaf28779ff46$var$addScopeToNode(clone2, reactiveScope, templateEl);
            clone2._x_refreshXForScope = (newScope)=>{
                Object.entries(newScope).forEach(([key2, value])=>{
                    reactiveScope[key2] = value;
                });
            };
            $8c83eaf28779ff46$var$mutateDom(()=>{
                lastEl.after(clone2);
                $8c83eaf28779ff46$var$skipDuringClone(()=>$8c83eaf28779ff46$var$initTree(clone2))();
            });
            if (typeof key === "object") $8c83eaf28779ff46$var$warn("x-for key cannot be an object, it must be a string or an integer", templateEl);
            lookup[key] = clone2;
        }
        for(let i = 0; i < sames.length; i++)lookup[sames[i]]._x_refreshXForScope(scopes[keys.indexOf(sames[i])]);
        templateEl._x_prevKeys = keys;
    });
}
function $8c83eaf28779ff46$var$parseForExpression(expression) {
    let forIteratorRE = /,([^,\}\]]*)(?:,([^,\}\]]*))?$/;
    let stripParensRE = /^\s*\(|\)\s*$/g;
    let forAliasRE = /([\s\S]*?)\s+(?:in|of)\s+([\s\S]*)/;
    let inMatch = expression.match(forAliasRE);
    if (!inMatch) return;
    let res = {};
    res.items = inMatch[2].trim();
    let item = inMatch[1].replace(stripParensRE, "").trim();
    let iteratorMatch = item.match(forIteratorRE);
    if (iteratorMatch) {
        res.item = item.replace(forIteratorRE, "").trim();
        res.index = iteratorMatch[1].trim();
        if (iteratorMatch[2]) res.collection = iteratorMatch[2].trim();
    } else res.item = item;
    return res;
}
function $8c83eaf28779ff46$var$getIterationScopeVariables(iteratorNames, item, index, items) {
    let scopeVariables = {};
    if (/^\[.*\]$/.test(iteratorNames.item) && Array.isArray(item)) {
        let names = iteratorNames.item.replace("[", "").replace("]", "").split(",").map((i)=>i.trim());
        names.forEach((name, i)=>{
            scopeVariables[name] = item[i];
        });
    } else if (/^\{.*\}$/.test(iteratorNames.item) && !Array.isArray(item) && typeof item === "object") {
        let names = iteratorNames.item.replace("{", "").replace("}", "").split(",").map((i)=>i.trim());
        names.forEach((name)=>{
            scopeVariables[name] = item[name];
        });
    } else scopeVariables[iteratorNames.item] = item;
    if (iteratorNames.index) scopeVariables[iteratorNames.index] = index;
    if (iteratorNames.collection) scopeVariables[iteratorNames.collection] = items;
    return scopeVariables;
}
function $8c83eaf28779ff46$var$isNumeric3(subject) {
    return !Array.isArray(subject) && !isNaN(subject);
}
// packages/alpinejs/src/directives/x-ref.js
function $8c83eaf28779ff46$var$handler3() {}
$8c83eaf28779ff46$var$handler3.inline = (el, { expression: expression }, { cleanup: cleanup2 })=>{
    let root = $8c83eaf28779ff46$var$closestRoot(el);
    if (!root._x_refs) root._x_refs = {};
    root._x_refs[expression] = el;
    cleanup2(()=>delete root._x_refs[expression]);
};
$8c83eaf28779ff46$var$directive("ref", $8c83eaf28779ff46$var$handler3);
// packages/alpinejs/src/directives/x-if.js
$8c83eaf28779ff46$var$directive("if", (el, { expression: expression }, { effect: effect3, cleanup: cleanup2 })=>{
    if (el.tagName.toLowerCase() !== "template") $8c83eaf28779ff46$var$warn("x-if can only be used on a <template> tag", el);
    let evaluate2 = $8c83eaf28779ff46$var$evaluateLater(el, expression);
    let show = ()=>{
        if (el._x_currentIfEl) return el._x_currentIfEl;
        let clone2 = el.content.cloneNode(true).firstElementChild;
        $8c83eaf28779ff46$var$addScopeToNode(clone2, {}, el);
        $8c83eaf28779ff46$var$mutateDom(()=>{
            el.after(clone2);
            $8c83eaf28779ff46$var$skipDuringClone(()=>$8c83eaf28779ff46$var$initTree(clone2))();
        });
        el._x_currentIfEl = clone2;
        el._x_undoIf = ()=>{
            $8c83eaf28779ff46$var$mutateDom(()=>{
                $8c83eaf28779ff46$var$destroyTree(clone2);
                clone2.remove();
            });
            delete el._x_currentIfEl;
        };
        return clone2;
    };
    let hide = ()=>{
        if (!el._x_undoIf) return;
        el._x_undoIf();
        delete el._x_undoIf;
    };
    effect3(()=>evaluate2((value)=>{
            value ? show() : hide();
        }));
    cleanup2(()=>el._x_undoIf && el._x_undoIf());
});
// packages/alpinejs/src/directives/x-id.js
$8c83eaf28779ff46$var$directive("id", (el, { expression: expression }, { evaluate: evaluate2 })=>{
    let names = evaluate2(expression);
    names.forEach((name)=>$8c83eaf28779ff46$var$setIdRoot(el, name));
});
$8c83eaf28779ff46$var$interceptClone((from, to)=>{
    if (from._x_ids) to._x_ids = from._x_ids;
});
// packages/alpinejs/src/directives/x-on.js
$8c83eaf28779ff46$var$mapAttributes($8c83eaf28779ff46$var$startingWith("@", $8c83eaf28779ff46$var$into($8c83eaf28779ff46$var$prefix("on:"))));
$8c83eaf28779ff46$var$directive("on", $8c83eaf28779ff46$var$skipDuringClone((el, { value: value, modifiers: modifiers, expression: expression }, { cleanup: cleanup2 })=>{
    let evaluate2 = expression ? $8c83eaf28779ff46$var$evaluateLater(el, expression) : ()=>{};
    if (el.tagName.toLowerCase() === "template") {
        if (!el._x_forwardEvents) el._x_forwardEvents = [];
        if (!el._x_forwardEvents.includes(value)) el._x_forwardEvents.push(value);
    }
    let removeListener = $8c83eaf28779ff46$var$on(el, value, modifiers, (e)=>{
        evaluate2(()=>{}, {
            scope: {
                "$event": e
            },
            params: [
                e
            ]
        });
    });
    cleanup2(()=>removeListener());
}));
// packages/alpinejs/src/directives/index.js
$8c83eaf28779ff46$var$warnMissingPluginDirective("Collapse", "collapse", "collapse");
$8c83eaf28779ff46$var$warnMissingPluginDirective("Intersect", "intersect", "intersect");
$8c83eaf28779ff46$var$warnMissingPluginDirective("Focus", "trap", "focus");
$8c83eaf28779ff46$var$warnMissingPluginDirective("Mask", "mask", "mask");
function $8c83eaf28779ff46$var$warnMissingPluginDirective(name, directiveName, slug) {
    $8c83eaf28779ff46$var$directive(directiveName, (el)=>$8c83eaf28779ff46$var$warn(`You can't use [x-${directiveName}] without first installing the "${name}" plugin here: https://alpinejs.dev/plugins/${slug}`, el));
}
// packages/alpinejs/src/index.js
$8c83eaf28779ff46$var$alpine_default.setEvaluator($8c83eaf28779ff46$var$normalEvaluator);
$8c83eaf28779ff46$var$alpine_default.setRawEvaluator($8c83eaf28779ff46$var$normalRawEvaluator);
$8c83eaf28779ff46$var$alpine_default.setReactivityEngine({
    reactive: $8c83eaf28779ff46$var$reactive2,
    effect: $8c83eaf28779ff46$var$effect2,
    release: $8c83eaf28779ff46$var$stop,
    raw: $8c83eaf28779ff46$var$toRaw
});
var $8c83eaf28779ff46$export$b7ee041e4ad2afec = $8c83eaf28779ff46$var$alpine_default;
// packages/alpinejs/builds/module.js
var $8c83eaf28779ff46$export$2e2bcd8739ae039 = $8c83eaf28779ff46$export$b7ee041e4ad2afec;


// packages/persist/src/index.js
function $9b2f94dab0f686ea$export$9a6132153fba2e0(Alpine) {
    let persist = ()=>{
        let alias;
        let storage;
        try {
            storage = localStorage;
        } catch (e) {
            console.error(e);
            console.warn("Alpine: $persist is using temporary storage since localStorage is unavailable.");
            let dummy = /* @__PURE__ */ new Map();
            storage = {
                getItem: dummy.get.bind(dummy),
                setItem: dummy.set.bind(dummy)
            };
        }
        return Alpine.interceptor((initialValue, getter, setter, path, key)=>{
            let lookup = alias || `_x_${path}`;
            let initial = $9b2f94dab0f686ea$var$storageHas(lookup, storage) ? $9b2f94dab0f686ea$var$storageGet(lookup, storage) : initialValue;
            setter(initial);
            Alpine.effect(()=>{
                let value = getter();
                $9b2f94dab0f686ea$var$storageSet(lookup, value, storage);
                setter(value);
            });
            return initial;
        }, (func)=>{
            func.as = (key)=>{
                alias = key;
                return func;
            }, func.using = (target)=>{
                storage = target;
                return func;
            };
        });
    };
    Object.defineProperty(Alpine, "$persist", {
        get: ()=>persist()
    });
    Alpine.magic("persist", persist);
    Alpine.persist = (key, { get: get, set: set }, storage = localStorage)=>{
        let initial = $9b2f94dab0f686ea$var$storageHas(key, storage) ? $9b2f94dab0f686ea$var$storageGet(key, storage) : get();
        set(initial);
        Alpine.effect(()=>{
            let value = get();
            $9b2f94dab0f686ea$var$storageSet(key, value, storage);
            set(value);
        });
    };
}
function $9b2f94dab0f686ea$var$storageHas(key, storage) {
    return storage.getItem(key) !== null;
}
function $9b2f94dab0f686ea$var$storageGet(key, storage) {
    let value = storage.getItem(key);
    if (value === void 0) return;
    return JSON.parse(value);
}
function $9b2f94dab0f686ea$var$storageSet(key, value, storage) {
    storage.setItem(key, JSON.stringify(value));
}
// packages/persist/builds/module.js
var $9b2f94dab0f686ea$export$2e2bcd8739ae039 = $9b2f94dab0f686ea$export$9a6132153fba2e0;


// packages/morph/src/morph.js
function $fd2717825660a9ff$var$morph(from, toHtml, options) {
    $fd2717825660a9ff$var$monkeyPatchDomSetAttributeToAllowAtSymbols();
    let context = $fd2717825660a9ff$var$createMorphContext(options);
    let toEl = typeof toHtml === "string" ? $fd2717825660a9ff$var$createElement(toHtml) : toHtml;
    if (window.Alpine && window.Alpine.closestDataStack && !from._x_dataStack) {
        toEl._x_dataStack = window.Alpine.closestDataStack(from);
        toEl._x_dataStack && window.Alpine.cloneNode(from, toEl);
    }
    context.patch(from, toEl);
    return from;
}
function $fd2717825660a9ff$var$morphBetween(startMarker, endMarker, toHtml, options = {}) {
    $fd2717825660a9ff$var$monkeyPatchDomSetAttributeToAllowAtSymbols();
    let context = $fd2717825660a9ff$var$createMorphContext(options);
    let fromContainer = startMarker.parentNode;
    let fromBlock = new $fd2717825660a9ff$var$Block(startMarker, endMarker);
    let toContainer = typeof toHtml === "string" ? (()=>{
        let container = document.createElement("div");
        container.insertAdjacentHTML("beforeend", toHtml);
        return container;
    })() : toHtml;
    let toStartMarker = document.createComment("[morph-start]");
    let toEndMarker = document.createComment("[morph-end]");
    toContainer.insertBefore(toStartMarker, toContainer.firstChild);
    toContainer.appendChild(toEndMarker);
    let toBlock = new $fd2717825660a9ff$var$Block(toStartMarker, toEndMarker);
    if (window.Alpine && window.Alpine.closestDataStack) {
        toContainer._x_dataStack = window.Alpine.closestDataStack(fromContainer);
        toContainer._x_dataStack && window.Alpine.cloneNode(fromContainer, toContainer);
    }
    context.patchChildren(fromBlock, toBlock);
}
function $fd2717825660a9ff$var$createMorphContext(options = {}) {
    let defaultGetKey = (el)=>el.getAttribute("key");
    let noop = ()=>{};
    let context = {
        key: options.key || defaultGetKey,
        lookahead: options.lookahead || false,
        updating: options.updating || noop,
        updated: options.updated || noop,
        removing: options.removing || noop,
        removed: options.removed || noop,
        adding: options.adding || noop,
        added: options.added || noop
    };
    context.patch = function(from, to) {
        if (context.differentElementNamesTypesOrKeys(from, to)) return context.swapElements(from, to);
        let updateChildrenOnly = false;
        let skipChildren = false;
        let skipUntil = (predicate)=>context.skipUntilCondition = predicate;
        if ($fd2717825660a9ff$var$shouldSkipChildren(context.updating, ()=>skipChildren = true, skipUntil, from, to, ()=>updateChildrenOnly = true)) return;
        if (from.nodeType === 1 && window.Alpine) {
            window.Alpine.cloneNode(from, to);
            if (from._x_teleport && to._x_teleport) context.patch(from._x_teleport, to._x_teleport);
        }
        if ($fd2717825660a9ff$var$textOrComment(to)) {
            context.patchNodeValue(from, to);
            context.updated(from, to);
            return;
        }
        if (!updateChildrenOnly) context.patchAttributes(from, to);
        context.updated(from, to);
        if (!skipChildren) context.patchChildren(from, to);
    };
    context.differentElementNamesTypesOrKeys = function(from, to) {
        return from.nodeType != to.nodeType || from.nodeName != to.nodeName || context.getKey(from) != context.getKey(to);
    };
    context.swapElements = function(from, to) {
        if ($fd2717825660a9ff$var$shouldSkip(context.removing, from)) return;
        let toCloned = to.cloneNode(true);
        if ($fd2717825660a9ff$var$shouldSkip(context.adding, toCloned)) return;
        from.replaceWith(toCloned);
        context.removed(from);
        context.added(toCloned);
    };
    context.patchNodeValue = function(from, to) {
        let value = to.nodeValue;
        if (from.nodeValue !== value) from.nodeValue = value;
    };
    context.patchAttributes = function(from, to) {
        if (from._x_transitioning) return;
        if (from._x_isShown && !to._x_isShown) return;
        if (!from._x_isShown && to._x_isShown) return;
        let domAttributes = Array.from(from.attributes);
        let toAttributes = Array.from(to.attributes);
        for(let i = domAttributes.length - 1; i >= 0; i--){
            let name = domAttributes[i].name;
            if (!to.hasAttribute(name)) from.removeAttribute(name);
        }
        for(let i = toAttributes.length - 1; i >= 0; i--){
            let name = toAttributes[i].name;
            let value = toAttributes[i].value;
            if (from.getAttribute(name) !== value) from.setAttribute(name, value);
        }
    };
    context.patchChildren = function(from, to) {
        let fromKeys = context.keyToMap(from.children);
        let fromKeyHoldovers = {};
        let currentTo = $fd2717825660a9ff$var$getFirstNode(to);
        let currentFrom = $fd2717825660a9ff$var$getFirstNode(from);
        while(currentTo){
            $fd2717825660a9ff$var$seedingMatchingId(currentTo, currentFrom);
            let toKey = context.getKey(currentTo);
            let fromKey = context.getKey(currentFrom);
            if (context.skipUntilCondition) {
                let fromDone = !currentFrom || context.skipUntilCondition(currentFrom);
                let toDone = !currentTo || context.skipUntilCondition(currentTo);
                if (fromDone && toDone) context.skipUntilCondition = null;
                else {
                    if (!fromDone) currentFrom = currentFrom && $fd2717825660a9ff$var$getNextSibling(from, currentFrom);
                    if (!toDone) currentTo = currentTo && $fd2717825660a9ff$var$getNextSibling(to, currentTo);
                    continue;
                }
            }
            if (!currentFrom) {
                if (toKey && fromKeyHoldovers[toKey]) {
                    let holdover = fromKeyHoldovers[toKey];
                    from.appendChild(holdover);
                    currentFrom = holdover;
                    fromKey = context.getKey(currentFrom);
                } else {
                    if (!$fd2717825660a9ff$var$shouldSkip(context.adding, currentTo)) {
                        let clone = currentTo.cloneNode(true);
                        from.appendChild(clone);
                        context.added(clone);
                    }
                    currentTo = $fd2717825660a9ff$var$getNextSibling(to, currentTo);
                    continue;
                }
            }
            let isIf = (node)=>node && node.nodeType === 8 && node.textContent === "[if BLOCK]><![endif]";
            let isEnd = (node)=>node && node.nodeType === 8 && node.textContent === "[if ENDBLOCK]><![endif]";
            if (isIf(currentTo) && isIf(currentFrom)) {
                let nestedIfCount = 0;
                let fromBlockStart = currentFrom;
                while(currentFrom){
                    let next = $fd2717825660a9ff$var$getNextSibling(from, currentFrom);
                    if (isIf(next)) nestedIfCount++;
                    else if (isEnd(next) && nestedIfCount > 0) nestedIfCount--;
                    else if (isEnd(next) && nestedIfCount === 0) {
                        currentFrom = next;
                        break;
                    }
                    currentFrom = next;
                }
                let fromBlockEnd = currentFrom;
                nestedIfCount = 0;
                let toBlockStart = currentTo;
                while(currentTo){
                    let next = $fd2717825660a9ff$var$getNextSibling(to, currentTo);
                    if (isIf(next)) nestedIfCount++;
                    else if (isEnd(next) && nestedIfCount > 0) nestedIfCount--;
                    else if (isEnd(next) && nestedIfCount === 0) {
                        currentTo = next;
                        break;
                    }
                    currentTo = next;
                }
                let toBlockEnd = currentTo;
                let fromBlock = new $fd2717825660a9ff$var$Block(fromBlockStart, fromBlockEnd);
                let toBlock = new $fd2717825660a9ff$var$Block(toBlockStart, toBlockEnd);
                context.patchChildren(fromBlock, toBlock);
                continue;
            }
            if (currentFrom.nodeType === 1 && context.lookahead && !currentFrom.isEqualNode(currentTo)) {
                let nextToElementSibling = $fd2717825660a9ff$var$getNextSibling(to, currentTo);
                let found = false;
                while(!found && nextToElementSibling){
                    if (nextToElementSibling.nodeType === 1 && currentFrom.isEqualNode(nextToElementSibling)) {
                        found = true;
                        currentFrom = context.addNodeBefore(from, currentTo, currentFrom);
                        fromKey = context.getKey(currentFrom);
                    }
                    nextToElementSibling = $fd2717825660a9ff$var$getNextSibling(to, nextToElementSibling);
                }
            }
            if (toKey !== fromKey) {
                if (!toKey && fromKey) {
                    fromKeyHoldovers[fromKey] = currentFrom;
                    currentFrom = context.addNodeBefore(from, currentTo, currentFrom);
                    fromKeyHoldovers[fromKey].remove();
                    currentFrom = $fd2717825660a9ff$var$getNextSibling(from, currentFrom);
                    currentTo = $fd2717825660a9ff$var$getNextSibling(to, currentTo);
                    continue;
                }
                if (toKey && !fromKey) {
                    if (fromKeys[toKey]) {
                        currentFrom.replaceWith(fromKeys[toKey]);
                        currentFrom = fromKeys[toKey];
                        fromKey = context.getKey(currentFrom);
                    }
                }
                if (toKey && fromKey) {
                    let fromKeyNode = fromKeys[toKey];
                    if (fromKeyNode) {
                        fromKeyHoldovers[fromKey] = currentFrom;
                        currentFrom.replaceWith(fromKeyNode);
                        currentFrom = fromKeyNode;
                        fromKey = context.getKey(currentFrom);
                    } else {
                        fromKeyHoldovers[fromKey] = currentFrom;
                        currentFrom = context.addNodeBefore(from, currentTo, currentFrom);
                        fromKeyHoldovers[fromKey].remove();
                        currentFrom = $fd2717825660a9ff$var$getNextSibling(from, currentFrom);
                        currentTo = $fd2717825660a9ff$var$getNextSibling(to, currentTo);
                        continue;
                    }
                }
            }
            let currentFromNext = currentFrom && $fd2717825660a9ff$var$getNextSibling(from, currentFrom);
            context.patch(currentFrom, currentTo);
            currentTo = currentTo && $fd2717825660a9ff$var$getNextSibling(to, currentTo);
            currentFrom = currentFromNext;
        }
        let removals = [];
        while(currentFrom){
            if (!$fd2717825660a9ff$var$shouldSkip(context.removing, currentFrom)) removals.push(currentFrom);
            currentFrom = $fd2717825660a9ff$var$getNextSibling(from, currentFrom);
        }
        while(removals.length){
            let domForRemoval = removals.shift();
            domForRemoval.remove();
            context.removed(domForRemoval);
        }
    };
    context.getKey = function(el) {
        return el && el.nodeType === 1 && context.key(el);
    };
    context.keyToMap = function(els) {
        let map = {};
        for (let el of els){
            let theKey = context.getKey(el);
            if (theKey) map[theKey] = el;
        }
        return map;
    };
    context.addNodeBefore = function(parent, node, beforeMe) {
        if (!$fd2717825660a9ff$var$shouldSkip(context.adding, node)) {
            let clone = node.cloneNode(true);
            parent.insertBefore(clone, beforeMe);
            context.added(clone);
            return clone;
        }
        return node;
    };
    return context;
}
$fd2717825660a9ff$var$morph.step = ()=>{};
$fd2717825660a9ff$var$morph.log = ()=>{};
function $fd2717825660a9ff$var$shouldSkip(hook, ...args) {
    let skip = false;
    hook(...args, ()=>skip = true);
    return skip;
}
function $fd2717825660a9ff$var$shouldSkipChildren(hook, skipChildren, skipUntil, ...args) {
    let skip = false;
    hook(...args, ()=>skip = true, skipChildren, skipUntil);
    return skip;
}
var $fd2717825660a9ff$var$patched = false;
function $fd2717825660a9ff$var$createElement(html) {
    const template = document.createElement("template");
    template.innerHTML = html;
    return template.content.firstElementChild;
}
function $fd2717825660a9ff$var$textOrComment(el) {
    return el.nodeType === 3 || el.nodeType === 8;
}
var $fd2717825660a9ff$var$Block = class {
    constructor(start, end){
        this.startComment = start;
        this.endComment = end;
    }
    get children() {
        let children = [];
        let currentNode = this.startComment.nextSibling;
        while(currentNode && currentNode !== this.endComment){
            children.push(currentNode);
            currentNode = currentNode.nextSibling;
        }
        return children;
    }
    appendChild(child) {
        this.endComment.before(child);
    }
    get firstChild() {
        let first = this.startComment.nextSibling;
        if (first === this.endComment) return;
        return first;
    }
    nextNode(reference) {
        let next = reference.nextSibling;
        if (next === this.endComment) return;
        return next;
    }
    insertBefore(newNode, reference) {
        reference.before(newNode);
        return newNode;
    }
};
function $fd2717825660a9ff$var$getFirstNode(parent) {
    return parent.firstChild;
}
function $fd2717825660a9ff$var$getNextSibling(parent, reference) {
    let next;
    if (parent instanceof $fd2717825660a9ff$var$Block) next = parent.nextNode(reference);
    else next = reference.nextSibling;
    return next;
}
function $fd2717825660a9ff$var$monkeyPatchDomSetAttributeToAllowAtSymbols() {
    if ($fd2717825660a9ff$var$patched) return;
    $fd2717825660a9ff$var$patched = true;
    let original = Element.prototype.setAttribute;
    let hostDiv = document.createElement("div");
    Element.prototype.setAttribute = function newSetAttribute(name, value) {
        if (!name.includes("@")) return original.call(this, name, value);
        hostDiv.innerHTML = `<span ${name}="${value}"></span>`;
        let attr = hostDiv.firstElementChild.getAttributeNode(name);
        hostDiv.firstElementChild.removeAttributeNode(attr);
        this.setAttributeNode(attr);
    };
}
function $fd2717825660a9ff$var$seedingMatchingId(to, from) {
    let fromId = from && from._x_bindings && from._x_bindings.id;
    if (!fromId) return;
    if (!to.setAttribute) return;
    to.setAttribute("id", fromId);
    to.id = fromId;
}
// packages/morph/src/index.js
function $fd2717825660a9ff$export$2e5e8c41f5d4e7c7(Alpine) {
    Alpine.morph = $fd2717825660a9ff$var$morph;
    Alpine.morphBetween = $fd2717825660a9ff$var$morphBetween;
}
// packages/morph/builds/module.js
var $fd2717825660a9ff$export$2e2bcd8739ae039 = $fd2717825660a9ff$export$2e5e8c41f5d4e7c7;


// packages/intersect/src/index.js
function $b3c3dd74fa47ca5d$export$1f4807a235930d45(Alpine) {
    Alpine.directive("intersect", Alpine.skipDuringClone((el, { value: value, expression: expression, modifiers: modifiers }, { evaluateLater: evaluateLater, cleanup: cleanup })=>{
        let evaluate = evaluateLater(expression);
        let options = {
            rootMargin: $b3c3dd74fa47ca5d$var$getRootMargin(modifiers),
            threshold: $b3c3dd74fa47ca5d$var$getThreshold(modifiers)
        };
        let observer = new IntersectionObserver((entries)=>{
            entries.forEach((entry)=>{
                if (entry.isIntersecting === (value === "leave")) return;
                evaluate();
                modifiers.includes("once") && observer.disconnect();
            });
        }, options);
        observer.observe(el);
        cleanup(()=>{
            observer.disconnect();
        });
    }));
}
function $b3c3dd74fa47ca5d$var$getThreshold(modifiers) {
    if (modifiers.includes("full")) return 0.99;
    if (modifiers.includes("half")) return 0.5;
    if (!modifiers.includes("threshold")) return 0;
    let threshold = modifiers[modifiers.indexOf("threshold") + 1];
    if (threshold === "100") return 1;
    if (threshold === "0") return 0;
    return Number(`.${threshold}`);
}
function $b3c3dd74fa47ca5d$var$getLengthValue(rawValue) {
    let match = rawValue.match(/^(-?[0-9]+)(px|%)?$/);
    return match ? match[1] + (match[2] || "px") : void 0;
}
function $b3c3dd74fa47ca5d$var$getRootMargin(modifiers) {
    const key = "margin";
    const fallback = "0px 0px 0px 0px";
    const index = modifiers.indexOf(key);
    if (index === -1) return fallback;
    let values = [];
    for(let i = 1; i < 5; i++)values.push($b3c3dd74fa47ca5d$var$getLengthValue(modifiers[index + i] || ""));
    values = values.filter((v)=>v !== void 0);
    return values.length ? values.join(" ").trim() : fallback;
}
// packages/intersect/builds/module.js
var $b3c3dd74fa47ca5d$export$2e2bcd8739ae039 = $b3c3dd74fa47ca5d$export$1f4807a235930d45;


window.Alpine = (0, $8c83eaf28779ff46$export$2e2bcd8739ae039 // 将 Alpine 实例添加到窗口对象中。
);
(0, $8c83eaf28779ff46$export$2e2bcd8739ae039).plugin((0, $9b2f94dab0f686ea$export$2e2bcd8739ae039)) // 用于在本地存储中持久化数据的插件
;
(0, $8c83eaf28779ff46$export$2e2bcd8739ae039).plugin((0, $fd2717825660a9ff$export$2e2bcd8739ae039)) // 不丢失 Alpine 页面状态的情况下，根据服务器请求更新 HTML
;
(0, $8c83eaf28779ff46$export$2e2bcd8739ae039).plugin((0, $b3c3dd74fa47ca5d$export$2e2bcd8739ae039)) //  Intersection Observer 的一个便捷封装，在元素进入视口时做出反应。
;


/* eslint-disable promise/prefer-await-to-then */ const $fe9db8f8165fa827$var$methodMap = [
    [
        'requestFullscreen',
        'exitFullscreen',
        'fullscreenElement',
        'fullscreenEnabled',
        'fullscreenchange',
        'fullscreenerror'
    ],
    // New WebKit
    [
        'webkitRequestFullscreen',
        'webkitExitFullscreen',
        'webkitFullscreenElement',
        'webkitFullscreenEnabled',
        'webkitfullscreenchange',
        'webkitfullscreenerror'
    ],
    // Old WebKit
    [
        'webkitRequestFullScreen',
        'webkitCancelFullScreen',
        'webkitCurrentFullScreenElement',
        'webkitCancelFullScreen',
        'webkitfullscreenchange',
        'webkitfullscreenerror'
    ],
    [
        'mozRequestFullScreen',
        'mozCancelFullScreen',
        'mozFullScreenElement',
        'mozFullScreenEnabled',
        'mozfullscreenchange',
        'mozfullscreenerror'
    ],
    [
        'msRequestFullscreen',
        'msExitFullscreen',
        'msFullscreenElement',
        'msFullscreenEnabled',
        'MSFullscreenChange',
        'MSFullscreenError'
    ]
];
const $fe9db8f8165fa827$var$nativeAPI = (()=>{
    if (typeof document === 'undefined') return false;
    const unprefixedMethods = $fe9db8f8165fa827$var$methodMap[0];
    const returnValue = {};
    for (const methodList of $fe9db8f8165fa827$var$methodMap){
        const exitFullscreenMethod = methodList?.[1];
        if (exitFullscreenMethod in document) {
            for (const [index, method] of methodList.entries())returnValue[unprefixedMethods[index]] = method;
            return returnValue;
        }
    }
    return false;
})();
const $fe9db8f8165fa827$var$eventNameMap = {
    change: $fe9db8f8165fa827$var$nativeAPI.fullscreenchange,
    error: $fe9db8f8165fa827$var$nativeAPI.fullscreenerror
};
// eslint-disable-next-line import/no-mutable-exports
let $fe9db8f8165fa827$var$screenfull = {
    // eslint-disable-next-line default-param-last
    request (element = document.documentElement, options) {
        return new Promise((resolve, reject)=>{
            const onFullScreenEntered = ()=>{
                $fe9db8f8165fa827$var$screenfull.off('change', onFullScreenEntered);
                resolve();
            };
            $fe9db8f8165fa827$var$screenfull.on('change', onFullScreenEntered);
            const returnPromise = element[$fe9db8f8165fa827$var$nativeAPI.requestFullscreen](options);
            if (returnPromise instanceof Promise) returnPromise.then(onFullScreenEntered).catch(reject);
        });
    },
    exit () {
        return new Promise((resolve, reject)=>{
            if (!$fe9db8f8165fa827$var$screenfull.isFullscreen) {
                resolve();
                return;
            }
            const onFullScreenExit = ()=>{
                $fe9db8f8165fa827$var$screenfull.off('change', onFullScreenExit);
                resolve();
            };
            $fe9db8f8165fa827$var$screenfull.on('change', onFullScreenExit);
            const returnPromise = document[$fe9db8f8165fa827$var$nativeAPI.exitFullscreen]();
            if (returnPromise instanceof Promise) returnPromise.then(onFullScreenExit).catch(reject);
        });
    },
    toggle (element, options) {
        return $fe9db8f8165fa827$var$screenfull.isFullscreen ? $fe9db8f8165fa827$var$screenfull.exit() : $fe9db8f8165fa827$var$screenfull.request(element, options);
    },
    onchange (callback) {
        $fe9db8f8165fa827$var$screenfull.on('change', callback);
    },
    onerror (callback) {
        $fe9db8f8165fa827$var$screenfull.on('error', callback);
    },
    on (event, callback) {
        const eventName = $fe9db8f8165fa827$var$eventNameMap[event];
        if (eventName) document.addEventListener(eventName, callback, false);
    },
    off (event, callback) {
        const eventName = $fe9db8f8165fa827$var$eventNameMap[event];
        if (eventName) document.removeEventListener(eventName, callback, false);
    },
    raw: $fe9db8f8165fa827$var$nativeAPI
};
Object.defineProperties($fe9db8f8165fa827$var$screenfull, {
    isFullscreen: {
        get: ()=>Boolean(document[$fe9db8f8165fa827$var$nativeAPI.fullscreenElement])
    },
    element: {
        enumerable: true,
        get: ()=>document[$fe9db8f8165fa827$var$nativeAPI.fullscreenElement] ?? undefined
    },
    isEnabled: {
        enumerable: true,
        // Coerce to boolean in case of old WebKit.
        get: ()=>Boolean(document[$fe9db8f8165fa827$var$nativeAPI.fullscreenEnabled])
    }
});
if (!$fe9db8f8165fa827$var$nativeAPI) $fe9db8f8165fa827$var$screenfull = {
    isEnabled: false
};
var $fe9db8f8165fa827$export$2e2bcd8739ae039 = $fe9db8f8165fa827$var$screenfull;


window.Screenfull = (0, $fe9db8f8165fa827$export$2e2bcd8739ae039 // 将 screenfull 实例添加到窗口对象中。
);


// Alpine 使用 Persist 插件，用会话 cookie 作为存储
// https://alpinejs.dev/plugins/persist#custom-storage
// 定义自定义存储对象，公开 getItem 函数和 setItem 函数
window.cookieStorage = {
    getItem (key) {
        let cookies = document.cookie.split(";");
        for(let i = 0; i < cookies.length; i++){
            let cookie = cookies[i].split("=");
            if (key === cookie[0].trim()) return decodeURIComponent(cookie[1]);
        }
        return null;
    },
    setItem (key, value) {
        document.cookie = `${key}=${encodeURIComponent(value)}; SameSite=Lax`; //SameSite设置默认值（Lax），防止控制台报错。加载图像或框架（frame）的请求将不会包含用户的 Cookie。
    }
} // // 然后就可以这样使用使用 cookieStorage 作为 Persist 插件的存储了
 // Alpine.store('cookie', {
 //     someCookieKey: Alpine.$persist(false).using(cookieStorage).as('someCookieKey'),
 // })
;


// 用Alpine Persist 注册全局变量
// https://alpinejs.dev/plugins/persist#using-alpine-persist-global
/**
 * 解析UserAgent获取浏览器信息
 * @returns {string} 浏览器名称
 */ function $3d5cedefbd41a457$var$getBrowserInfo() {
    const ua = navigator.userAgent;
    let browser = 'Unknown';
    if (ua.indexOf('Firefox') > -1) browser = 'Firefox';
    else if (ua.indexOf('Edg') > -1) browser = 'Edge';
    else if (ua.indexOf('Chrome') > -1) browser = 'Chrome';
    else if (ua.indexOf('Safari') > -1) browser = 'Safari';
    else if (ua.indexOf('Opera') > -1 || ua.indexOf('OPR') > -1) browser = 'Opera';
    else if (ua.indexOf('Trident') > -1 || ua.indexOf('MSIE') > -1) browser = 'IE';
    return browser;
}
/**
 * 解析UserAgent获取系统信息
 * @returns {string} 系统名称
 */ function $3d5cedefbd41a457$var$getSystemInfo() {
    const ua = navigator.userAgent;
    let os = 'Unknown';
    if (ua.indexOf('Win') > -1) os = 'Windows';
    else if (ua.indexOf('Mac') > -1) os = 'MacOS';
    else if (ua.indexOf('Linux') > -1) os = 'Linux';
    else if (ua.indexOf('Android') > -1) os = 'Android';
    else if (ua.indexOf('iOS') > -1 || ua.indexOf('iPhone') > -1 || ua.indexOf('iPad') > -1) os = 'iOS';
    return os;
}
/**
 * 生成随机字符串
 * @returns {string} 随机字符串
 */ function $3d5cedefbd41a457$var$generateRandomString() {
    return (Date.now() % 10000000).toString(36) + Math.random().toString(36).substring(2, 5);
}
// 浏览器 系统信息 随机字符串
const $3d5cedefbd41a457$var$browser = $3d5cedefbd41a457$var$getBrowserInfo();
const $3d5cedefbd41a457$var$system = $3d5cedefbd41a457$var$getSystemInfo();
const $3d5cedefbd41a457$var$randomString = $3d5cedefbd41a457$var$generateRandomString();
// 生成userID: 使用UserAgent的哈希值 + 随机字符串，确保唯一性且长度适中
const $3d5cedefbd41a457$var$initClientID = `Client_${$3d5cedefbd41a457$var$randomString}_${$3d5cedefbd41a457$var$system}_${$3d5cedefbd41a457$var$browser}`;
Alpine.store('global', {
    nowPageNum: 1,
    allPageNum: 1,
    onlineBook: true,
    // 播放器：音量（0~100）
    playerVolume: Alpine.$persist(100).as('global.playerVolume'),
    // 播放器：是否静音
    playerMuted: Alpine.$persist(false).as('global.playerMuted'),
    // 播放器：是否自动播放下一曲
    autoPlayNext: Alpine.$persist(true).as('global.autoPlayNext'),
    // 播放器：是否循环播放播放列表
    loopPlaylist: Alpine.$persist(true).as('global.loopPlaylist'),
    // 自动切边
    autoCrop: Alpine.$persist(false).as('global.autoCrop'),
    // 自动切边阈值,范围是0~100。多数情况下 1 就够了。
    autoCropNum: Alpine.$persist(1).as('global.autoCropNum'),
    // 是否压缩图片
    autoResize: Alpine.$persist(false).as('global.autoResize'),
    // 压缩图片限宽
    autoResizeWidth: Alpine.$persist(800).as('global.autoResizeWidth'),
    // bgPattern 背景花纹
    bgPattern: Alpine.$persist('grid-line').as('global.bgPattern'),
    // 是否禁止缓存（TODO：缓存功能优化与测试）
    noCache: Alpine.$persist(false).as('global.noCache'),
    // clientID 用于识别匿名用户与设备
    clientID: Alpine.$persist($3d5cedefbd41a457$var$initClientID).as('global.clientID'),
    // debugMode 是否开启调试模式
    debugMode: Alpine.$persist(true).as('global.debugMode'),
    //是否通过websocket同步翻页
    syncPageByWS: Alpine.$persist(true).as('global.syncPageByWS'),
    // bookSortBy 书籍排序方式 以按照文件名、修改时间、文件大小排序（或反向排序）
    bookSortBy: Alpine.$persist('name').as('global.bookSortBy'),
    // pageSortBy 书页排序顺序 以按照文件名、修改时间、文件大小排序（或反向排序）
    pageSortBy: Alpine.$persist('name').as('global.pageSortBy'),
    language: Alpine.$persist('en').as('global.language'),
    //是否保存阅读进度（页数）到本地存储
    saveReadingProgress: Alpine.$persist(true).as('global.saveReadingProgress'),
    // 从本地存储加载页码并跳转
    loadPageNumFromLocalStorage (book_id, callbackFunction) {
        if (!this.saveReadingProgress) return;
        try {
            const key = `pageNum_${book_id}`;
            const savedPageNum = localStorage.getItem(key);
            if (savedPageNum !== null && !isNaN(parseInt(savedPageNum))) {
                const pageNum = parseInt(savedPageNum);
                // 确保页码在有效范围内
                if (pageNum > 0 && pageNum <= Alpine.store('global').allPageNum) {
                    console.log(`\u{52A0}\u{8F7D}\u{5230}\u{672C}\u{5730}\u{5B58}\u{50A8}\u{7684}\u{9875}\u{7801}: ${pageNum}`);
                    callbackFunction(); // 跳转函数,或发送书签更新信息
                }
            }
        } catch (e) {
            console.error("Error loading page number from localStorage:", e);
        }
    },
    // 保存当前页码到本地存储
    savePageNumToLocalStorage (book_id) {
        if (!this.saveReadingProgress) return;
        if (!book_id) {
            console.warn("savePageNumToLocalStorage: book_id is required");
            return;
        }
        try {
            const key = `pageNum_${book_id}`;
            const nowPageNum = Alpine.store('global').nowPageNum;
            localStorage.setItem(key, nowPageNum);
        } catch (e) {
            console.error("Error saving page number to localStorage:", e);
        }
    },
    // readerMode 当前阅读模式: infinite_scroll  paged_scroll  flip_page
    readMode: Alpine.$persist('infinite_scroll').as('global.readMode'),
    // 切换阅读模式
    infiniteScrollLoadAllPage (mode) {
        this.readMode = "infinite_scroll";
        const url = new URL(window.location.href);
        url.searchParams.delete("start");
        window.location.href = url.href;
    },
    onChangeReadMode () {
        // 切换阅读模式时，如果在阅读，就修改URL路径 参考文献：https://developer.mozilla.org/zh-CN/docs/Web/API/URL
        const url = new URL(window.location.href);
        const pathname = url.pathname;
        // 使用 URLSearchParams 提取键值对
        const params = new URLSearchParams(url.search);
        // 分割路径为各层级关键词, filter(Boolean) 的作用是去除空字符串 如//aa/bb/ 会产生空字符串(虽然这里不会这么做)
        const pathSegments = url.pathname.split('/').filter(Boolean); // like ["nginx_test", "scroll", "id3DcA1v9"]
        const book_id = pathSegments[pathSegments.length - 1];
        console.log(`\u{5207}\u{6362}\u{9605}\u{8BFB}\u{6A21}\u{5F0F}\u{5230}: ${this.readMode}, \u{5F53}\u{524D}\u{8DEF}\u{5F84}: ${pathname},${pathSegments}, \u{67E5}\u{8BE2}\u{53C2}\u{6570}: ${params.toString()}`);
        // 卷轴(无限)模式
        if (this.readMode === 'infinite_scroll') // 如果已经是无限卷轴模式
        {
            if (pathSegments.includes("scroll") && params.get("page") === null) {
                console.log(`${pathSegments} , ${params.get("page")}`);
                console.log("\u5DF2\u7ECF\u662F\u65E0\u9650\u5377\u8F74\u6A21\u5F0F\uFF0C\u65E0\u9700\u5207\u6362");
                return;
            }
        }
        // 卷轴(分页)模式
        if (this.readMode === 'paged_scroll') // 如果已经是分页卷轴模式
        {
            if (pathSegments.includes("scroll") && params.get("page") !== null) {
                console.log(`${pathSegments} , ${params.get("page")}`);
                console.log("\u5DF2\u7ECF\u662F\u5206\u9875\u5377\u8F74\u6A21\u5F0F\uFF0C\u65E0\u9700\u5207\u6362");
                return;
            }
        }
        // 翻页模式
        if (this.readMode === 'page_flip') // 如果已经是翻页模式
        {
            if (pathSegments.includes("flip")) {
                console.log("\u5DF2\u7ECF\u662F\u7FFB\u9875\u6A21\u5F0F\uFF0C\u65E0\u9700\u5207\u6362");
                console.log(`${pathSegments} , ${params.get("start")}`);
                return;
            }
        }
        // 跳转到新的阅读模式URL
        if (pathSegments.includes("scroll") || pathSegments.includes("flip")) window.location.href = this.getReadURL(book_id, this.nowPageNum);
    },
    getReadURL (book_id, start_index) {
        // TODO: 处理旧版本数据干扰的问题。若干个版本后大概就不需要了，到时候删除这段代码。
        if (this.readMode !== 'page_flip' && this.readMode !== 'paged_scroll' && this.readMode !== 'infinite_scroll') {
            console.error(`\u{672A}\u{77E5}\u{7684}\u{9605}\u{8BFB}\u{6A21}\u{5F0F}: ${this.readMode}, \u{53EF}\u{80FD}\u{662F}\u{65E7}\u{7248}\u{672C}\u{6570}\u{636E}\u{5E72}\u{6270}, \u{91CD}\u{7F6E}\u{4E3A} infinite_scroll`);
            this.readMode = 'infinite_scroll';
        }
        let PAGED_SIZE = 32;
        // console.log(`生成阅读模式URL: ${this.readMode}`);
        // console.log(`当前页码: ${start_index}`);
        const url = new URL(window.location.href);
        // 卷轴(无限)
        if (this.readMode === 'infinite_scroll') {
            let new_url = new URL(`/scroll/${book_id}`, url.origin);
            if (start_index > 1) new_url.searchParams.set("start", start_index.toString());
            return new_url.href;
        }
        // 卷轴(分页)
        if (this.readMode === 'paged_scroll') {
            let new_url = new URL(`/scroll/${book_id}`, url.origin);
            let page = Math.floor(start_index / PAGED_SIZE) + 1;
            new_url.searchParams.set("page", page.toString());
            return new_url.href;
        }
        // 翻页(左右)
        if (this.readMode === 'page_flip') {
            let new_url = new URL(`/flip/${book_id}`, url.origin);
            if (start_index > 1) new_url.searchParams.set("start", start_index.toString());
            return new_url.href;
        }
        return "";
    },
    // 竖屏模式
    isPortrait: false,
    // 横屏模式
    isLandscape: true,
    // 获取cookie里面存储的值
    getCookieValue (bookID, valueName) {
        let pgCookie = "";
        const paramName = bookID === "" ? `$${valueName}` : `${bookID}_${valueName}`;
        const cookies = document.cookie.split(";");
        for(let i = 0; i < cookies.length; i++){
            const cookie = cookies[i].trim();
            if (cookie.startsWith(paramName)) pgCookie = decodeURIComponent(cookie.substring(paramName.length + 1));
        }
        return pgCookie;
    },
    setPaginationIndex (bookID, valueName, value) {
        const paramName = bookID === "" ? `$${valueName}` : `${bookID}_${valueName}`;
        // 设置cookie，过期时间为365天
        const expirationDate = new Date();
        expirationDate.setDate(expirationDate.getDate() + 365);
        document.cookie = `${paramName}${encodeURIComponent(value)}; expires=${expirationDate.toUTCString()}; path=/; SameSite=Lax`;
        window.location.reload();
    },
    /**
     * 调用后端 /api/store-bookmark 接口，更新书签信息
     * @param {Object} params
     * @param {string} params.type - 书签类型，例如 'auto'
     * @param {string} params.bookId - 书籍ID
     * @param {number} params.pageIndex - 页码（1 起始）
     * @param {string} [params.label='自动书签'] - 书签名称，当前后端固定为自动书签，仅用于日志
     * @returns {Promise<Object|string>} 后端返回的响应体
     */ async UpdateBookmark ({ type: type = 'auto', bookId: bookId, pageIndex: pageIndex, description: description = '' } = {}) {
        if (!bookId) {
            const error = new Error('UpdateBookmark: bookId is required');
            if (this.debugMode) console.error(error);
            throw error;
        }
        if (!Number.isInteger(pageIndex) || pageIndex <= 0) {
            const error = new Error('UpdateBookmark: pageIndex must be a positive integer');
            if (this.debugMode) console.error(error);
            throw error;
        }
        if (description === '') description = `${$3d5cedefbd41a457$var$browser} in ${$3d5cedefbd41a457$var$system}`;
        const payload = {
            type: type,
            book_id: bookId,
            page_index: pageIndex,
            description: description
        };
        const response = await fetch('/api/store-bookmark', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify(payload)
        });
        const contentType = response.headers.get('content-type') || '';
        const isJSON = contentType.includes('application/json');
        const responseBody = isJSON ? await response.json() : await response.text();
        if (!response.ok) {
            const error = new Error(`UpdateBookmark failed: ${response.status} ${response.statusText}`);
            if (this.debugMode) console.error('[UpdateBookmark] error', error, responseBody);
            throw error;
        }
        return responseBody;
    },
    // 检测并设置视口方向
    checkOrientation () {
        const isPortrait = window.innerHeight > window.innerWidth;
        this.isPortrait = isPortrait;
        this.isLandscape = !isPortrait;
    //console.log(`当前视口方向: ${isPortrait ? '竖屏' : '横屏'}`);
    },
    // 初始化方法
    init () {
        // 设置初始方向
        this.checkOrientation();
        // 添加视口变化监听
        window.addEventListener('resize', ()=>{
            this.checkOrientation();
        });
    }
});
// 初始化全局存储
document.addEventListener('alpine:initialized', ()=>{
    Alpine.store('global').init();
});
const $3d5cedefbd41a457$var$url = new URL(window.location.href);
if (($3d5cedefbd41a457$var$url.protocol === 'http:' || $3d5cedefbd41a457$var$url.protocol === 'https:') && !window.location.toString().endsWith('.html')) Alpine.store('global').onlineBook = true;
else Alpine.store('global').onlineBook = false;


// BookShelf 书架设置
Alpine.store('shelf', {
    bookCardMode: Alpine.$persist('gird').as('shelf.bookCardMode'),
    showFilename: Alpine.$persist(true).as('shelf.showFilename'),
    showFileIcon: Alpine.$persist(true).as('shelf.showFileIcon'),
    simplifyTitle: Alpine.$persist(true).as('shelf.simplifyTitle'),
    openInNewTab: Alpine.$persist(false).as('shelf.openInNewTab'),
    InfiniteDropdown: Alpine.$persist(false).as('shelf.InfiniteDropdown'),
    bookCardShowTitleFlag: Alpine.$persist(true).as('shelf.bookCardShowTitleFlag'),
    readingProgressPercent: Alpine.$persist(false).as('shelf.readingProgressPercent'),
    syncScrollFlag: false,
    // 屏幕宽横比,inLandscapeMode的判断依据
    aspectRatio: 1.2,
    // 可见范围宽高的具体值
    clientWidth: 0,
    clientHeight: 0
});


// Scroll 卷轴模式
Alpine.store("scroll", {
    simplifyTitle: Alpine.$persist(true).as("scroll.simplifyTitle"),
    //下拉模式下，漫画页面的底部间距。单位px。
    marginBottomOnScrollMode: Alpine.$persist(0).as("scroll.marginBottomOnScrollMode"),
    //卷轴模式下，是否分页加载（反之则无限下拉）
    fixedPagination: Alpine.$persist(false).as("scroll.fixedPagination"),
    // 卷轴模式的同步滚动,目前还没做
    syncScrollFlag: Alpine.$persist(false).as("scroll.syncScrollFlag"),
    imageMaxWidth: 400,
    // 屏幕宽横比,inLandscapeMode的判断依据
    aspectRatio: 1.2,
    // 可见范围宽高的具体值
    clientWidth: 0,
    clientHeight: 0,
    //漫画页的单位,是否使用固定值
    widthUseFixedValue: Alpine.$persist(true).as("scroll.widthUseFixedValue"),
    //竖屏(Portrait)状态的漫画页宽度,百分比
    portraitWidthPercent: Alpine.$persist(100).as("scroll.portraitWidthPercent"),
    //横屏(Landscape)状态的漫画页宽度,百分比
    singlePageWidth_Percent: Alpine.$persist(60).as("scroll.singlePageWidth_Percent"),
    doublePageWidth_Percent: Alpine.$persist(60).as("scroll.doublePageWidth_Percent"),
    //横屏(Landscape)状态的漫画页宽度。px。
    singlePageWidth_PX: Alpine.$persist(800).as("scroll.singlePageWidth_PX"),
    doublePageWidth_PX: Alpine.$persist(800).as("scroll.doublePageWidth_PX"),
    //书籍数据,需要从远程拉取
    //是否显示顶部页头
    showHeaderFlag: true,
    //是否显示页数
    showPageNum: Alpine.$persist(false).as("scroll.showPageNum"),
    //ws翻页相关
    syncPageByWS: Alpine.$persist(false).as("scroll.syncPageByWS")
});


// Flip 翻页模式
Alpine.store('flip', {
    imageMaxWidth: 400,
    //自动隐藏工具条
    autoHideToolbar: Alpine.$persist(false).as('flip.autoHideToolbar'),
    // autoHideToolbar: Alpine.$persist((() => {
    //     const ua = navigator.userAgent || '';
    //     const isMobileUA = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini|Mobi/i.test(ua);
    //     const isTouch = (('ontouchstart' in window) || (navigator.maxTouchPoints > 1));
    //     //console.log('isMobileUA', isMobileUA);
    //     //console.log('isTouch', isTouch);
    //     // return isMobileUA || isTouch;
    //     return isMobileUA;
    // })()).as('flip.autoHideToolbar'),
    //自动对齐
    autoAlign: Alpine.$persist(true).as('flip.autoAlignTop'),
    //是否显示页头
    show_header: Alpine.$persist(true).as('flip.show_header'),
    //是否显示页脚
    showFooter: Alpine.$persist(true).as('flip.showFooter'),
    //是否显示页数
    showPageNum: Alpine.$persist(true).as('flip.showPageNum'),
    //是否是日本漫画【右半屏翻页,从左到右(true)】【右半屏翻页,从右到左(false)】
    mangaMode: Alpine.$persist(true).as('flip.mangaMode'),
    //触摸滑动翻页
    swipeTurn: Alpine.$persist(true).as('flip.swipeTurn'),
    //鼠标滚轮翻页
    wheelFlip: Alpine.$persist(false).as('flip.wheelFlip'),
    //双页模式
    doublePageMode: Alpine.$persist(false).as('flip.doublePageMode'),
    //自动拼合双页(TODO)
    autoDoublePageMode: Alpine.$persist(false).as('flip.autoDoublePageModeFlag'),
    //素描模式标记
    sketchModeFlag: false,
    //是否显示素描提示
    showPageHint: Alpine.$persist(false).as('flip.showPageHint'),
    //翻页间隔时间
    sketchFlipSecond: 30,
    //计时用,从0开始
    sketchSecondCount: 0,
    // ============ 滑动动画配置（新增）============
    // 滑动动画持续时间（毫秒）
    swipeAnimationDuration: Alpine.$persist(300).as('flip.swipeAnimationDuration'),
    // 回弹动画持续时间（毫秒）
    resetAnimationDuration: Alpine.$persist(400).as('flip.resetAnimationDuration'),
    // 滑动阈值（像素）- 超过这个值才会触发翻页
    swipeThreshold: Alpine.$persist(100).as('flip.swipeThreshold'),
    // 快速滑动超时时间（毫秒）
    swipeTimeout: Alpine.$persist(300).as('flip.swipeTimeout'),
    // ============ 其他可配置参数（新增）============
    // 预加载图片范围
    preloadRange: Alpine.$persist(10).as('flip.preloadRange'),
    // 滚轮节流延迟（毫秒）
    wheelThrottleDelay: Alpine.$persist(250).as('flip.wheelThrottleDelay'),
    // WebSocket 最大重连次数
    websocketMaxReconnect: Alpine.$persist(200).as('flip.websocketMaxReconnect'),
    // WebSocket 重连间隔（毫秒）
    websocketReconnectInterval: Alpine.$persist(3000).as('flip.websocketReconnectInterval')
});


// 自定义主题
Alpine.store('theme', {
    theme: Alpine.$persist('light').as('theme'),
    interfaceColor: '#F5F5E4',
    backgroundColor: '#E0D9CD',
    textColor: '#000000',
    toggleTheme () {
        this.theme = this.theme === 'light' ? 'dark' : 'light';
    }
});


//请求图片文件时，可添加的额外参数
const $bbe0293ef6bc5e52$export$444dbc9dc04ca304 = {
    resize_width: -1,
    resize_height: -1,
    do_compress_image: false,
    resize_max_width: 800,
    resize_max_height: -1,
    do_auto_crop: false,
    auto_crop_num: 1,
    gray: false
};
//添加各种字符串参数,不需要的话为空
const $bbe0293ef6bc5e52$var$resize_width_str = $bbe0293ef6bc5e52$export$444dbc9dc04ca304.resize_width > 0 ? "&resize_width=" + $bbe0293ef6bc5e52$export$444dbc9dc04ca304.resize_width : "";
const $bbe0293ef6bc5e52$var$resize_height_str = $bbe0293ef6bc5e52$export$444dbc9dc04ca304.resize_height > 0 ? "&resize_height=" + $bbe0293ef6bc5e52$export$444dbc9dc04ca304.resize_height : "";
const $bbe0293ef6bc5e52$var$gray_str = $bbe0293ef6bc5e52$export$444dbc9dc04ca304.gray ? "&gray=true" : "";
const $bbe0293ef6bc5e52$var$do_compress_image_str = $bbe0293ef6bc5e52$export$444dbc9dc04ca304.do_compress_image ? "&resize_max_width=" + $bbe0293ef6bc5e52$export$444dbc9dc04ca304.resize_max_width : "";
const $bbe0293ef6bc5e52$var$resize_max_height_str = $bbe0293ef6bc5e52$export$444dbc9dc04ca304.resize_max_height > 0 ? "&resize_max_height=" + $bbe0293ef6bc5e52$export$444dbc9dc04ca304.resize_max_height : "";
const $bbe0293ef6bc5e52$var$auto_crop_str = $bbe0293ef6bc5e52$export$444dbc9dc04ca304.do_auto_crop ? "&auto_crop=" + $bbe0293ef6bc5e52$export$444dbc9dc04ca304.auto_crop_num : "";
//所有附加的转换参数
let $bbe0293ef6bc5e52$export$52550a8b6a1f2afe = $bbe0293ef6bc5e52$var$resize_width_str + $bbe0293ef6bc5e52$var$resize_height_str + $bbe0293ef6bc5e52$var$do_compress_image_str + $bbe0293ef6bc5e52$var$resize_max_height_str + $bbe0293ef6bc5e52$var$auto_crop_str + $bbe0293ef6bc5e52$var$gray_str;
if ($bbe0293ef6bc5e52$export$52550a8b6a1f2afe !== "") {
    $bbe0293ef6bc5e52$export$52550a8b6a1f2afe = "?" + $bbe0293ef6bc5e52$export$52550a8b6a1f2afe.substring(1);
    console.log("addStr:", $bbe0293ef6bc5e52$export$52550a8b6a1f2afe);
}


// Start Alpine.
Alpine.start();
// Document ready function to ensure the DOM is fully loaded.
document.addEventListener('DOMContentLoaded', function() {
    initFlowbite() // initialize Flowbite
    ;
});
// Add event listeners for all HTMX events.
document.body.addEventListener('htmx:afterSwap htmx:afterRequest htmx:afterSettle', function() {
    initFlowbite() // initialize Flowbite
    ;
});

})();
//# sourceMappingURL=main.js.map
