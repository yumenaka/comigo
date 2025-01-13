(function () {

function $parcel$interopDefault(a) {
  return a && a.__esModule ? a.default : a;
}

      var $parcel$global =
        typeof globalThis !== 'undefined'
          ? globalThis
          : typeof self !== 'undefined'
          ? self
          : typeof window !== 'undefined'
          ? window
          : typeof global !== 'undefined'
          ? global
          : {};
  
var $parcel$modules = {};
var $parcel$inits = {};

var parcelRequire = $parcel$global["parcelRequire94c2"];

if (parcelRequire == null) {
  parcelRequire = function(id) {
    if (id in $parcel$modules) {
      return $parcel$modules[id].exports;
    }
    if (id in $parcel$inits) {
      var init = $parcel$inits[id];
      delete $parcel$inits[id];
      var module = {id: id, exports: {}};
      $parcel$modules[id] = module;
      init.call(module.exports, module, module.exports);
      return module.exports;
    }
    var err = new Error("Cannot find module '" + id + "'");
    err.code = 'MODULE_NOT_FOUND';
    throw err;
  };

  parcelRequire.register = function register(id, init) {
    $parcel$inits[id] = init;
  };

  $parcel$global["parcelRequire94c2"] = parcelRequire;
}

var parcelRegister = parcelRequire.register;
parcelRegister("eZIU4", function(module, exports) {
(function(e, t) {
    if (typeof define === "function" && define.amd) define([], t);
    else if (0, module.exports) module.exports = t();
    else e.htmx = e.htmx || t();
})(typeof self !== "undefined" ? self : this, function() {
    return function() {
        "use strict";
        var Q = {
            onLoad: F,
            process: zt,
            on: de,
            off: ge,
            trigger: ce,
            ajax: Nr,
            find: C,
            findAll: f,
            closest: v,
            values: function(e, t) {
                var r = dr(e, t || "post");
                return r.values;
            },
            remove: _,
            addClass: z,
            removeClass: n,
            toggleClass: $,
            takeClass: W,
            defineExtension: Ur,
            removeExtension: Br,
            logAll: V,
            logNone: j,
            logger: null,
            config: {
                historyEnabled: true,
                historyCacheSize: 10,
                refreshOnHistoryMiss: false,
                defaultSwapStyle: "innerHTML",
                defaultSwapDelay: 0,
                defaultSettleDelay: 20,
                includeIndicatorStyles: true,
                indicatorClass: "htmx-indicator",
                requestClass: "htmx-request",
                addedClass: "htmx-added",
                settlingClass: "htmx-settling",
                swappingClass: "htmx-swapping",
                allowEval: true,
                allowScriptTags: true,
                inlineScriptNonce: "",
                attributesToSettle: [
                    "class",
                    "style",
                    "width",
                    "height"
                ],
                withCredentials: false,
                timeout: 0,
                wsReconnectDelay: "full-jitter",
                wsBinaryType: "blob",
                disableSelector: "[hx-disable], [data-hx-disable]",
                useTemplateFragments: false,
                scrollBehavior: "smooth",
                defaultFocusScroll: false,
                getCacheBusterParam: false,
                globalViewTransitions: false,
                methodsThatUseUrlParams: [
                    "get"
                ],
                selfRequestsOnly: false,
                ignoreTitle: false,
                scrollIntoViewOnBoost: true,
                triggerSpecsCache: null
            },
            parseInterval: d,
            _: t,
            createEventSource: function(e) {
                return new EventSource(e, {
                    withCredentials: true
                });
            },
            createWebSocket: function(e) {
                var t = new WebSocket(e, []);
                t.binaryType = Q.config.wsBinaryType;
                return t;
            },
            version: "1.9.12"
        };
        var r = {
            addTriggerHandler: Lt,
            bodyContains: se,
            canAccessLocalStorage: U,
            findThisElement: xe,
            filterValues: yr,
            hasAttribute: o,
            getAttributeValue: te,
            getClosestAttributeValue: ne,
            getClosestMatch: c,
            getExpressionVars: Hr,
            getHeaders: xr,
            getInputValues: dr,
            getInternalData: ae,
            getSwapSpecification: wr,
            getTriggerSpecs: it,
            getTarget: ye,
            makeFragment: l,
            mergeObjects: le,
            makeSettleInfo: T,
            oobSwap: Ee,
            querySelectorExt: ue,
            selectAndSwap: je,
            settleImmediately: nr,
            shouldCancel: ut,
            triggerEvent: ce,
            triggerErrorEvent: fe,
            withExtensions: R
        };
        var w = [
            "get",
            "post",
            "put",
            "delete",
            "patch"
        ];
        var i = w.map(function(e) {
            return "[hx-" + e + "], [data-hx-" + e + "]";
        }).join(", ");
        var S = e("head"), q = e("title"), H = e("svg", true);
        function e(e1, t) {
            return new RegExp("<" + e1 + "(\\s[^>]*>|>)([\\s\\S]*?)<\\/" + e1 + ">", !!t ? "gim" : "im");
        }
        function d(e) {
            if (e == undefined) return undefined;
            let t = NaN;
            if (e.slice(-2) == "ms") t = parseFloat(e.slice(0, -2));
            else if (e.slice(-1) == "s") t = parseFloat(e.slice(0, -1)) * 1e3;
            else if (e.slice(-1) == "m") t = parseFloat(e.slice(0, -1)) * 60000;
            else t = parseFloat(e);
            return isNaN(t) ? undefined : t;
        }
        function ee(e, t) {
            return e.getAttribute && e.getAttribute(t);
        }
        function o(e, t) {
            return e.hasAttribute && (e.hasAttribute(t) || e.hasAttribute("data-" + t));
        }
        function te(e, t) {
            return ee(e, t) || ee(e, "data-" + t);
        }
        function u(e) {
            return e.parentElement;
        }
        function re() {
            return document;
        }
        function c(e, t) {
            while(e && !t(e))e = u(e);
            return e ? e : null;
        }
        function L(e, t, r) {
            var n = te(t, r);
            var i = te(t, "hx-disinherit");
            if (e !== t && i && (i === "*" || i.split(" ").indexOf(r) >= 0)) return "unset";
            else return n;
        }
        function ne(t, r) {
            var n = null;
            c(t, function(e) {
                return n = L(t, e, r);
            });
            if (n !== "unset") return n;
        }
        function h(e, t) {
            var r = e.matches || e.matchesSelector || e.msMatchesSelector || e.mozMatchesSelector || e.webkitMatchesSelector || e.oMatchesSelector;
            return r && r.call(e, t);
        }
        function A(e) {
            var t = /<([a-z][^\/\0>\x20\t\r\n\f]*)/i;
            var r = t.exec(e);
            if (r) return r[1].toLowerCase();
            else return "";
        }
        function s(e, t) {
            var r = new DOMParser;
            var n = r.parseFromString(e, "text/html");
            var i = n.body;
            while(t > 0){
                t--;
                i = i.firstChild;
            }
            if (i == null) i = re().createDocumentFragment();
            return i;
        }
        function N(e) {
            return /<body/.test(e);
        }
        function l(e) {
            var t = !N(e);
            var r = A(e);
            var n = e;
            if (r === "head") n = n.replace(S, "");
            if (Q.config.useTemplateFragments && t) {
                var i = s("<body><template>" + n + "</template></body>", 0);
                var a = i.querySelector("template").content;
                if (Q.config.allowScriptTags) oe(a.querySelectorAll("script"), function(e) {
                    if (Q.config.inlineScriptNonce) e.nonce = Q.config.inlineScriptNonce;
                    e.htmxExecuted = navigator.userAgent.indexOf("Firefox") === -1;
                });
                else oe(a.querySelectorAll("script"), function(e) {
                    _(e);
                });
                return a;
            }
            switch(r){
                case "thead":
                case "tbody":
                case "tfoot":
                case "colgroup":
                case "caption":
                    return s("<table>" + n + "</table>", 1);
                case "col":
                    return s("<table><colgroup>" + n + "</colgroup></table>", 2);
                case "tr":
                    return s("<table><tbody>" + n + "</tbody></table>", 2);
                case "td":
                case "th":
                    return s("<table><tbody><tr>" + n + "</tr></tbody></table>", 3);
                case "script":
                case "style":
                    return s("<div>" + n + "</div>", 1);
                default:
                    return s(n, 0);
            }
        }
        function ie(e) {
            if (e) e();
        }
        function I(e, t) {
            return Object.prototype.toString.call(e) === "[object " + t + "]";
        }
        function k(e) {
            return I(e, "Function");
        }
        function P(e) {
            return I(e, "Object");
        }
        function ae(e) {
            var t = "htmx-internal-data";
            var r = e[t];
            if (!r) r = e[t] = {};
            return r;
        }
        function M(e) {
            var t = [];
            if (e) for(var r = 0; r < e.length; r++)t.push(e[r]);
            return t;
        }
        function oe(e, t) {
            if (e) for(var r = 0; r < e.length; r++)t(e[r]);
        }
        function X(e) {
            var t = e.getBoundingClientRect();
            var r = t.top;
            var n = t.bottom;
            return r < window.innerHeight && n >= 0;
        }
        function se(e) {
            if (e.getRootNode && e.getRootNode() instanceof window.ShadowRoot) return re().body.contains(e.getRootNode().host);
            else return re().body.contains(e);
        }
        function D(e) {
            return e.trim().split(/\s+/);
        }
        function le(e, t) {
            for(var r in t)if (t.hasOwnProperty(r)) e[r] = t[r];
            return e;
        }
        function E(e) {
            try {
                return JSON.parse(e);
            } catch (e) {
                b(e);
                return null;
            }
        }
        function U() {
            var e = "htmx:localStorageTest";
            try {
                localStorage.setItem(e, e);
                localStorage.removeItem(e);
                return true;
            } catch (e) {
                return false;
            }
        }
        function B(t) {
            try {
                var e = new URL(t);
                if (e) t = e.pathname + e.search;
                if (!/^\/$/.test(t)) t = t.replace(/\/+$/, "");
                return t;
            } catch (e) {
                return t;
            }
        }
        function t(e) {
            return Tr(re().body, function() {
                return eval(e);
            });
        }
        function F(t) {
            var e = Q.on("htmx:load", function(e) {
                t(e.detail.elt);
            });
            return e;
        }
        function V() {
            Q.logger = function(e, t, r) {
                if (console) console.log(t, e, r);
            };
        }
        function j() {
            Q.logger = null;
        }
        function C(e, t) {
            if (t) return e.querySelector(t);
            else return C(re(), e);
        }
        function f(e, t) {
            if (t) return e.querySelectorAll(t);
            else return f(re(), e);
        }
        function _(e, t) {
            e = p(e);
            if (t) setTimeout(function() {
                _(e);
                e = null;
            }, t);
            else e.parentElement.removeChild(e);
        }
        function z(e, t, r) {
            e = p(e);
            if (r) setTimeout(function() {
                z(e, t);
                e = null;
            }, r);
            else e.classList && e.classList.add(t);
        }
        function n(e, t, r) {
            e = p(e);
            if (r) setTimeout(function() {
                n(e, t);
                e = null;
            }, r);
            else if (e.classList) {
                e.classList.remove(t);
                if (e.classList.length === 0) e.removeAttribute("class");
            }
        }
        function $(e, t) {
            e = p(e);
            e.classList.toggle(t);
        }
        function W(e, t) {
            e = p(e);
            oe(e.parentElement.children, function(e) {
                n(e, t);
            });
            z(e, t);
        }
        function v(e, t) {
            e = p(e);
            if (e.closest) return e.closest(t);
            else {
                do {
                    if (e == null || h(e, t)) return e;
                }while (e = e && u(e));
                return null;
            }
        }
        function g(e, t) {
            return e.substring(0, t.length) === t;
        }
        function G(e, t) {
            return e.substring(e.length - t.length) === t;
        }
        function J(e) {
            var t = e.trim();
            if (g(t, "<") && G(t, "/>")) return t.substring(1, t.length - 2);
            else return t;
        }
        function Z(e, t) {
            if (t.indexOf("closest ") === 0) return [
                v(e, J(t.substr(8)))
            ];
            else if (t.indexOf("find ") === 0) return [
                C(e, J(t.substr(5)))
            ];
            else if (t === "next") return [
                e.nextElementSibling
            ];
            else if (t.indexOf("next ") === 0) return [
                K(e, J(t.substr(5)))
            ];
            else if (t === "previous") return [
                e.previousElementSibling
            ];
            else if (t.indexOf("previous ") === 0) return [
                Y(e, J(t.substr(9)))
            ];
            else if (t === "document") return [
                document
            ];
            else if (t === "window") return [
                window
            ];
            else if (t === "body") return [
                document.body
            ];
            else return re().querySelectorAll(J(t));
        }
        var K = function(e, t) {
            var r = re().querySelectorAll(t);
            for(var n = 0; n < r.length; n++){
                var i = r[n];
                if (i.compareDocumentPosition(e) === Node.DOCUMENT_POSITION_PRECEDING) return i;
            }
        };
        var Y = function(e, t) {
            var r = re().querySelectorAll(t);
            for(var n = r.length - 1; n >= 0; n--){
                var i = r[n];
                if (i.compareDocumentPosition(e) === Node.DOCUMENT_POSITION_FOLLOWING) return i;
            }
        };
        function ue(e, t) {
            if (t) return Z(e, t)[0];
            else return Z(re().body, e)[0];
        }
        function p(e) {
            if (I(e, "String")) return C(e);
            else return e;
        }
        function ve(e, t, r) {
            if (k(t)) return {
                target: re().body,
                event: e,
                listener: t
            };
            else return {
                target: p(e),
                event: t,
                listener: r
            };
        }
        function de(t, r, n) {
            jr(function() {
                var e = ve(t, r, n);
                e.target.addEventListener(e.event, e.listener);
            });
            var e = k(r);
            return e ? r : n;
        }
        function ge(t, r, n) {
            jr(function() {
                var e = ve(t, r, n);
                e.target.removeEventListener(e.event, e.listener);
            });
            return k(r) ? r : n;
        }
        var pe = re().createElement("output");
        function me(e, t) {
            var r = ne(e, t);
            if (r) {
                if (r === "this") return [
                    xe(e, t)
                ];
                else {
                    var n = Z(e, r);
                    if (n.length === 0) {
                        b('The selector "' + r + '" on ' + t + " returned no matches!");
                        return [
                            pe
                        ];
                    } else return n;
                }
            }
        }
        function xe(e, t) {
            return c(e, function(e) {
                return te(e, t) != null;
            });
        }
        function ye(e) {
            var t = ne(e, "hx-target");
            if (t) {
                if (t === "this") return xe(e, "hx-target");
                else return ue(e, t);
            } else {
                var r = ae(e);
                if (r.boosted) return re().body;
                else return e;
            }
        }
        function be(e) {
            var t = Q.config.attributesToSettle;
            for(var r = 0; r < t.length; r++){
                if (e === t[r]) return true;
            }
            return false;
        }
        function we(t, r) {
            oe(t.attributes, function(e) {
                if (!r.hasAttribute(e.name) && be(e.name)) t.removeAttribute(e.name);
            });
            oe(r.attributes, function(e) {
                if (be(e.name)) t.setAttribute(e.name, e.value);
            });
        }
        function Se(e, t) {
            var r = Fr(t);
            for(var n = 0; n < r.length; n++){
                var i = r[n];
                try {
                    if (i.isInlineSwap(e)) return true;
                } catch (e) {
                    b(e);
                }
            }
            return e === "outerHTML";
        }
        function Ee(e, i, a) {
            var t = "#" + ee(i, "id");
            var o = "outerHTML";
            if (e === "true") ;
            else if (e.indexOf(":") > 0) {
                o = e.substr(0, e.indexOf(":"));
                t = e.substr(e.indexOf(":") + 1, e.length);
            } else o = e;
            var r = re().querySelectorAll(t);
            if (r) {
                oe(r, function(e) {
                    var t;
                    var r = i.cloneNode(true);
                    t = re().createDocumentFragment();
                    t.appendChild(r);
                    if (!Se(o, e)) t = r;
                    var n = {
                        shouldSwap: true,
                        target: e,
                        fragment: t
                    };
                    if (!ce(e, "htmx:oobBeforeSwap", n)) return;
                    e = n.target;
                    if (n["shouldSwap"]) Fe(o, e, e, t, a);
                    oe(a.elts, function(e) {
                        ce(e, "htmx:oobAfterSwap", n);
                    });
                });
                i.parentNode.removeChild(i);
            } else {
                i.parentNode.removeChild(i);
                fe(re().body, "htmx:oobErrorNoTarget", {
                    content: i
                });
            }
            return e;
        }
        function Ce(e, t, r) {
            var n = ne(e, "hx-select-oob");
            if (n) {
                var i = n.split(",");
                for(var a = 0; a < i.length; a++){
                    var o = i[a].split(":", 2);
                    var s = o[0].trim();
                    if (s.indexOf("#") === 0) s = s.substring(1);
                    var l = o[1] || "true";
                    var u = t.querySelector("#" + s);
                    if (u) Ee(l, u, r);
                }
            }
            oe(f(t, "[hx-swap-oob], [data-hx-swap-oob]"), function(e) {
                var t = te(e, "hx-swap-oob");
                if (t != null) Ee(t, e, r);
            });
        }
        function Re(e) {
            oe(f(e, "[hx-preserve], [data-hx-preserve]"), function(e) {
                var t = te(e, "id");
                var r = re().getElementById(t);
                if (r != null) e.parentNode.replaceChild(r, e);
            });
        }
        function Te(o, e, s) {
            oe(e.querySelectorAll("[id]"), function(e) {
                var t = ee(e, "id");
                if (t && t.length > 0) {
                    var r = t.replace("'", "\\'");
                    var n = e.tagName.replace(":", "\\:");
                    var i = o.querySelector(n + "[id='" + r + "']");
                    if (i && i !== o) {
                        var a = e.cloneNode();
                        we(e, i);
                        s.tasks.push(function() {
                            we(e, a);
                        });
                    }
                }
            });
        }
        function Oe(e) {
            return function() {
                n(e, Q.config.addedClass);
                zt(e);
                Nt(e);
                qe(e);
                ce(e, "htmx:load");
            };
        }
        function qe(e) {
            var t = "[autofocus]";
            var r = h(e, t) ? e : e.querySelector(t);
            if (r != null) r.focus();
        }
        function a(e, t, r, n) {
            Te(e, r, n);
            while(r.childNodes.length > 0){
                var i = r.firstChild;
                z(i, Q.config.addedClass);
                e.insertBefore(i, t);
                if (i.nodeType !== Node.TEXT_NODE && i.nodeType !== Node.COMMENT_NODE) n.tasks.push(Oe(i));
            }
        }
        function He(e, t) {
            var r = 0;
            while(r < e.length)t = (t << 5) - t + e.charCodeAt(r++) | 0;
            return t;
        }
        function Le(e) {
            var t = 0;
            if (e.attributes) for(var r = 0; r < e.attributes.length; r++){
                var n = e.attributes[r];
                if (n.value) {
                    t = He(n.name, t);
                    t = He(n.value, t);
                }
            }
            return t;
        }
        function Ae(e) {
            var t = ae(e);
            if (t.onHandlers) {
                for(var r = 0; r < t.onHandlers.length; r++){
                    const n = t.onHandlers[r];
                    e.removeEventListener(n.event, n.listener);
                }
                delete t.onHandlers;
            }
        }
        function Ne(e) {
            var t = ae(e);
            if (t.timeout) clearTimeout(t.timeout);
            if (t.webSocket) t.webSocket.close();
            if (t.sseEventSource) t.sseEventSource.close();
            if (t.listenerInfos) oe(t.listenerInfos, function(e) {
                if (e.on) e.on.removeEventListener(e.trigger, e.listener);
            });
            Ae(e);
            oe(Object.keys(t), function(e) {
                delete t[e];
            });
        }
        function m(e) {
            ce(e, "htmx:beforeCleanupElement");
            Ne(e);
            if (e.children) oe(e.children, function(e) {
                m(e);
            });
        }
        function Ie(t, e, r) {
            if (t.tagName === "BODY") return Ue(t, e, r);
            else {
                var n;
                var i = t.previousSibling;
                a(u(t), t, e, r);
                if (i == null) n = u(t).firstChild;
                else n = i.nextSibling;
                r.elts = r.elts.filter(function(e) {
                    return e != t;
                });
                while(n && n !== t){
                    if (n.nodeType === Node.ELEMENT_NODE) r.elts.push(n);
                    n = n.nextElementSibling;
                }
                m(t);
                u(t).removeChild(t);
            }
        }
        function ke(e, t, r) {
            return a(e, e.firstChild, t, r);
        }
        function Pe(e, t, r) {
            return a(u(e), e, t, r);
        }
        function Me(e, t, r) {
            return a(e, null, t, r);
        }
        function Xe(e, t, r) {
            return a(u(e), e.nextSibling, t, r);
        }
        function De(e, t, r) {
            m(e);
            return u(e).removeChild(e);
        }
        function Ue(e, t, r) {
            var n = e.firstChild;
            a(e, n, t, r);
            if (n) {
                while(n.nextSibling){
                    m(n.nextSibling);
                    e.removeChild(n.nextSibling);
                }
                m(n);
                e.removeChild(n);
            }
        }
        function Be(e, t, r) {
            var n = r || ne(e, "hx-select");
            if (n) {
                var i = re().createDocumentFragment();
                oe(t.querySelectorAll(n), function(e) {
                    i.appendChild(e);
                });
                t = i;
            }
            return t;
        }
        function Fe(e, t, r, n, i) {
            switch(e){
                case "none":
                    return;
                case "outerHTML":
                    Ie(r, n, i);
                    return;
                case "afterbegin":
                    ke(r, n, i);
                    return;
                case "beforebegin":
                    Pe(r, n, i);
                    return;
                case "beforeend":
                    Me(r, n, i);
                    return;
                case "afterend":
                    Xe(r, n, i);
                    return;
                case "delete":
                    De(r, n, i);
                    return;
                default:
                    var a = Fr(t);
                    for(var o = 0; o < a.length; o++){
                        var s = a[o];
                        try {
                            var l = s.handleSwap(e, r, n, i);
                            if (l) {
                                if (typeof l.length !== "undefined") for(var u = 0; u < l.length; u++){
                                    var f = l[u];
                                    if (f.nodeType !== Node.TEXT_NODE && f.nodeType !== Node.COMMENT_NODE) i.tasks.push(Oe(f));
                                }
                                return;
                            }
                        } catch (e) {
                            b(e);
                        }
                    }
                    if (e === "innerHTML") Ue(r, n, i);
                    else Fe(Q.config.defaultSwapStyle, t, r, n, i);
            }
        }
        function Ve(e) {
            if (e.indexOf("<title") > -1) {
                var t = e.replace(H, "");
                var r = t.match(q);
                if (r) return r[2];
            }
        }
        function je(e, t, r, n, i, a) {
            i.title = Ve(n);
            var o = l(n);
            if (o) {
                Ce(r, o, i);
                o = Be(r, o, a);
                Re(o);
                return Fe(e, r, t, o, i);
            }
        }
        function _e(e, t, r) {
            var n = e.getResponseHeader(t);
            if (n.indexOf("{") === 0) {
                var i = E(n);
                for(var a in i)if (i.hasOwnProperty(a)) {
                    var o = i[a];
                    if (!P(o)) o = {
                        value: o
                    };
                    ce(r, a, o);
                }
            } else {
                var s = n.split(",");
                for(var l = 0; l < s.length; l++)ce(r, s[l].trim(), []);
            }
        }
        var ze = /\s/;
        var x = /[\s,]/;
        var $e = /[_$a-zA-Z]/;
        var We = /[_$a-zA-Z0-9]/;
        var Ge = [
            '"',
            "'",
            "/"
        ];
        var Je = /[^\s]/;
        var Ze = /[{(]/;
        var Ke = /[})]/;
        function Ye(e) {
            var t = [];
            var r = 0;
            while(r < e.length){
                if ($e.exec(e.charAt(r))) {
                    var n = r;
                    while(We.exec(e.charAt(r + 1)))r++;
                    t.push(e.substr(n, r - n + 1));
                } else if (Ge.indexOf(e.charAt(r)) !== -1) {
                    var i = e.charAt(r);
                    var n = r;
                    r++;
                    while(r < e.length && e.charAt(r) !== i){
                        if (e.charAt(r) === "\\") r++;
                        r++;
                    }
                    t.push(e.substr(n, r - n + 1));
                } else {
                    var a = e.charAt(r);
                    t.push(a);
                }
                r++;
            }
            return t;
        }
        function Qe(e, t, r) {
            return $e.exec(e.charAt(0)) && e !== "true" && e !== "false" && e !== "this" && e !== r && t !== ".";
        }
        function et(e, t, r) {
            if (t[0] === "[") {
                t.shift();
                var n = 1;
                var i = " return (function(" + r + "){ return (";
                var a = null;
                while(t.length > 0){
                    var o = t[0];
                    if (o === "]") {
                        n--;
                        if (n === 0) {
                            if (a === null) i = i + "true";
                            t.shift();
                            i += ")})";
                            try {
                                var s = Tr(e, function() {
                                    return Function(i)();
                                }, function() {
                                    return true;
                                });
                                s.source = i;
                                return s;
                            } catch (e) {
                                fe(re().body, "htmx:syntax:error", {
                                    error: e,
                                    source: i
                                });
                                return null;
                            }
                        }
                    } else if (o === "[") n++;
                    if (Qe(o, a, r)) i += "((" + r + "." + o + ") ? (" + r + "." + o + ") : (window." + o + "))";
                    else i = i + o;
                    a = t.shift();
                }
            }
        }
        function y(e, t) {
            var r = "";
            while(e.length > 0 && !t.test(e[0]))r += e.shift();
            return r;
        }
        function tt(e) {
            var t;
            if (e.length > 0 && Ze.test(e[0])) {
                e.shift();
                t = y(e, Ke).trim();
                e.shift();
            } else t = y(e, x);
            return t;
        }
        var rt = "input, textarea, select";
        function nt(e, t, r) {
            var n = [];
            var i = Ye(t);
            do {
                y(i, Je);
                var a = i.length;
                var o = y(i, /[,\[\s]/);
                if (o !== "") {
                    if (o === "every") {
                        var s = {
                            trigger: "every"
                        };
                        y(i, Je);
                        s.pollInterval = d(y(i, /[,\[\s]/));
                        y(i, Je);
                        var l = et(e, i, "event");
                        if (l) s.eventFilter = l;
                        n.push(s);
                    } else if (o.indexOf("sse:") === 0) n.push({
                        trigger: "sse",
                        sseEvent: o.substr(4)
                    });
                    else {
                        var u = {
                            trigger: o
                        };
                        var l = et(e, i, "event");
                        if (l) u.eventFilter = l;
                        while(i.length > 0 && i[0] !== ","){
                            y(i, Je);
                            var f = i.shift();
                            if (f === "changed") u.changed = true;
                            else if (f === "once") u.once = true;
                            else if (f === "consume") u.consume = true;
                            else if (f === "delay" && i[0] === ":") {
                                i.shift();
                                u.delay = d(y(i, x));
                            } else if (f === "from" && i[0] === ":") {
                                i.shift();
                                if (Ze.test(i[0])) var c = tt(i);
                                else {
                                    var c = y(i, x);
                                    if (c === "closest" || c === "find" || c === "next" || c === "previous") {
                                        i.shift();
                                        var h = tt(i);
                                        if (h.length > 0) c += " " + h;
                                    }
                                }
                                u.from = c;
                            } else if (f === "target" && i[0] === ":") {
                                i.shift();
                                u.target = tt(i);
                            } else if (f === "throttle" && i[0] === ":") {
                                i.shift();
                                u.throttle = d(y(i, x));
                            } else if (f === "queue" && i[0] === ":") {
                                i.shift();
                                u.queue = y(i, x);
                            } else if (f === "root" && i[0] === ":") {
                                i.shift();
                                u[f] = tt(i);
                            } else if (f === "threshold" && i[0] === ":") {
                                i.shift();
                                u[f] = y(i, x);
                            } else fe(e, "htmx:syntax:error", {
                                token: i.shift()
                            });
                        }
                        n.push(u);
                    }
                }
                if (i.length === a) fe(e, "htmx:syntax:error", {
                    token: i.shift()
                });
                y(i, Je);
            }while (i[0] === "," && i.shift());
            if (r) r[t] = n;
            return n;
        }
        function it(e) {
            var t = te(e, "hx-trigger");
            var r = [];
            if (t) {
                var n = Q.config.triggerSpecsCache;
                r = n && n[t] || nt(e, t, n);
            }
            if (r.length > 0) return r;
            else if (h(e, "form")) return [
                {
                    trigger: "submit"
                }
            ];
            else if (h(e, 'input[type="button"], input[type="submit"]')) return [
                {
                    trigger: "click"
                }
            ];
            else if (h(e, rt)) return [
                {
                    trigger: "change"
                }
            ];
            else return [
                {
                    trigger: "click"
                }
            ];
        }
        function at(e) {
            ae(e).cancelled = true;
        }
        function ot(e, t, r) {
            var n = ae(e);
            n.timeout = setTimeout(function() {
                if (se(e) && n.cancelled !== true) {
                    if (!ct(r, e, Wt("hx:poll:trigger", {
                        triggerSpec: r,
                        target: e
                    }))) t(e);
                    ot(e, t, r);
                }
            }, r.pollInterval);
        }
        function st(e) {
            return location.hostname === e.hostname && ee(e, "href") && ee(e, "href").indexOf("#") !== 0;
        }
        function lt(t, r, e) {
            if (t.tagName === "A" && st(t) && (t.target === "" || t.target === "_self") || t.tagName === "FORM") {
                r.boosted = true;
                var n, i;
                if (t.tagName === "A") {
                    n = "get";
                    i = ee(t, "href");
                } else {
                    var a = ee(t, "method");
                    n = a ? a.toLowerCase() : "get";
                    n;
                    i = ee(t, "action");
                }
                e.forEach(function(e) {
                    ht(t, function(e, t) {
                        if (v(e, Q.config.disableSelector)) {
                            m(e);
                            return;
                        }
                        he(n, i, e, t);
                    }, r, e, true);
                });
            }
        }
        function ut(e, t) {
            if (e.type === "submit" || e.type === "click") {
                if (t.tagName === "FORM") return true;
                if (h(t, 'input[type="submit"], button') && v(t, "form") !== null) return true;
                if (t.tagName === "A" && t.href && (t.getAttribute("href") === "#" || t.getAttribute("href").indexOf("#") !== 0)) return true;
            }
            return false;
        }
        function ft(e, t) {
            return ae(e).boosted && e.tagName === "A" && t.type === "click" && (t.ctrlKey || t.metaKey);
        }
        function ct(e, t, r) {
            var n = e.eventFilter;
            if (n) try {
                return n.call(t, r) !== true;
            } catch (e) {
                fe(re().body, "htmx:eventFilter:error", {
                    error: e,
                    source: n.source
                });
                return true;
            }
            return false;
        }
        function ht(a, o, e, s, l) {
            var u = ae(a);
            var t;
            if (s.from) t = Z(a, s.from);
            else t = [
                a
            ];
            if (s.changed) t.forEach(function(e) {
                var t = ae(e);
                t.lastValue = e.value;
            });
            oe(t, function(n) {
                var i = function(e) {
                    if (!se(a)) {
                        n.removeEventListener(s.trigger, i);
                        return;
                    }
                    if (ft(a, e)) return;
                    if (l || ut(e, a)) e.preventDefault();
                    if (ct(s, a, e)) return;
                    var t = ae(e);
                    t.triggerSpec = s;
                    if (t.handledFor == null) t.handledFor = [];
                    if (t.handledFor.indexOf(a) < 0) {
                        t.handledFor.push(a);
                        if (s.consume) e.stopPropagation();
                        if (s.target && e.target) {
                            if (!h(e.target, s.target)) return;
                        }
                        if (s.once) {
                            if (u.triggeredOnce) return;
                            else u.triggeredOnce = true;
                        }
                        if (s.changed) {
                            var r = ae(n);
                            if (r.lastValue === n.value) return;
                            r.lastValue = n.value;
                        }
                        if (u.delayed) clearTimeout(u.delayed);
                        if (u.throttle) return;
                        if (s.throttle > 0) {
                            if (!u.throttle) {
                                o(a, e);
                                u.throttle = setTimeout(function() {
                                    u.throttle = null;
                                }, s.throttle);
                            }
                        } else if (s.delay > 0) u.delayed = setTimeout(function() {
                            o(a, e);
                        }, s.delay);
                        else {
                            ce(a, "htmx:trigger");
                            o(a, e);
                        }
                    }
                };
                if (e.listenerInfos == null) e.listenerInfos = [];
                e.listenerInfos.push({
                    trigger: s.trigger,
                    listener: i,
                    on: n
                });
                n.addEventListener(s.trigger, i);
            });
        }
        var vt = false;
        var dt = null;
        function gt() {
            if (!dt) {
                dt = function() {
                    vt = true;
                };
                window.addEventListener("scroll", dt);
                setInterval(function() {
                    if (vt) {
                        vt = false;
                        oe(re().querySelectorAll("[hx-trigger='revealed'],[data-hx-trigger='revealed']"), function(e) {
                            pt(e);
                        });
                    }
                }, 200);
            }
        }
        function pt(t) {
            if (!o(t, "data-hx-revealed") && X(t)) {
                t.setAttribute("data-hx-revealed", "true");
                var e = ae(t);
                if (e.initHash) ce(t, "revealed");
                else t.addEventListener("htmx:afterProcessNode", function(e) {
                    ce(t, "revealed");
                }, {
                    once: true
                });
            }
        }
        function mt(e, t, r) {
            var n = D(r);
            for(var i = 0; i < n.length; i++){
                var a = n[i].split(/:(.+)/);
                if (a[0] === "connect") xt(e, a[1], 0);
                if (a[0] === "send") bt(e);
            }
        }
        function xt(s, r, n) {
            if (!se(s)) return;
            if (r.indexOf("/") == 0) {
                var e = location.hostname + (location.port ? ":" + location.port : "");
                if (location.protocol == "https:") r = "wss://" + e + r;
                else if (location.protocol == "http:") r = "ws://" + e + r;
            }
            var t = Q.createWebSocket(r);
            t.onerror = function(e) {
                fe(s, "htmx:wsError", {
                    error: e,
                    socket: t
                });
                yt(s);
            };
            t.onclose = function(e) {
                if ([
                    1006,
                    1012,
                    1013
                ].indexOf(e.code) >= 0) {
                    var t = wt(n);
                    setTimeout(function() {
                        xt(s, r, n + 1);
                    }, t);
                }
            };
            t.onopen = function(e) {
                n = 0;
            };
            ae(s).webSocket = t;
            t.addEventListener("message", function(e) {
                if (yt(s)) return;
                var t = e.data;
                R(s, function(e) {
                    t = e.transformResponse(t, null, s);
                });
                var r = T(s);
                var n = l(t);
                var i = M(n.children);
                for(var a = 0; a < i.length; a++){
                    var o = i[a];
                    Ee(te(o, "hx-swap-oob") || "true", o, r);
                }
                nr(r.tasks);
            });
        }
        function yt(e) {
            if (!se(e)) {
                ae(e).webSocket.close();
                return true;
            }
        }
        function bt(u) {
            var f = c(u, function(e) {
                return ae(e).webSocket != null;
            });
            if (f) u.addEventListener(it(u)[0].trigger, function(e) {
                var t = ae(f).webSocket;
                var r = xr(u, f);
                var n = dr(u, "post");
                var i = n.errors;
                var a = n.values;
                var o = Hr(u);
                var s = le(a, o);
                var l = yr(s, u);
                l["HEADERS"] = r;
                if (i && i.length > 0) {
                    ce(u, "htmx:validation:halted", i);
                    return;
                }
                t.send(JSON.stringify(l));
                if (ut(e, u)) e.preventDefault();
            });
            else fe(u, "htmx:noWebSocketSourceError");
        }
        function wt(e) {
            var t = Q.config.wsReconnectDelay;
            if (typeof t === "function") return t(e);
            if (t === "full-jitter") {
                var r = Math.min(e, 6);
                var n = 1e3 * Math.pow(2, r);
                return n * Math.random();
            }
            b('htmx.config.wsReconnectDelay must either be a function or the string "full-jitter"');
        }
        function St(e, t, r) {
            var n = D(r);
            for(var i = 0; i < n.length; i++){
                var a = n[i].split(/:(.+)/);
                if (a[0] === "connect") Et(e, a[1]);
                if (a[0] === "swap") Ct(e, a[1]);
            }
        }
        function Et(t, e) {
            var r = Q.createEventSource(e);
            r.onerror = function(e) {
                fe(t, "htmx:sseError", {
                    error: e,
                    source: r
                });
                Tt(t);
            };
            ae(t).sseEventSource = r;
        }
        function Ct(a, o) {
            var s = c(a, Ot);
            if (s) {
                var l = ae(s).sseEventSource;
                var u = function(e) {
                    if (Tt(s)) return;
                    if (!se(a)) {
                        l.removeEventListener(o, u);
                        return;
                    }
                    var t = e.data;
                    R(a, function(e) {
                        t = e.transformResponse(t, null, a);
                    });
                    var r = wr(a);
                    var n = ye(a);
                    var i = T(a);
                    je(r.swapStyle, n, a, t, i);
                    nr(i.tasks);
                    ce(a, "htmx:sseMessage", e);
                };
                ae(a).sseListener = u;
                l.addEventListener(o, u);
            } else fe(a, "htmx:noSSESourceError");
        }
        function Rt(e, t, r) {
            var n = c(e, Ot);
            if (n) {
                var i = ae(n).sseEventSource;
                var a = function() {
                    if (!Tt(n)) {
                        if (se(e)) t(e);
                        else i.removeEventListener(r, a);
                    }
                };
                ae(e).sseListener = a;
                i.addEventListener(r, a);
            } else fe(e, "htmx:noSSESourceError");
        }
        function Tt(e) {
            if (!se(e)) {
                ae(e).sseEventSource.close();
                return true;
            }
        }
        function Ot(e) {
            return ae(e).sseEventSource != null;
        }
        function qt(e, t, r, n) {
            var i = function() {
                if (!r.loaded) {
                    r.loaded = true;
                    t(e);
                }
            };
            if (n > 0) setTimeout(i, n);
            else i();
        }
        function Ht(t, i, e) {
            var a = false;
            oe(w, function(r) {
                if (o(t, "hx-" + r)) {
                    var n = te(t, "hx-" + r);
                    a = true;
                    i.path = n;
                    i.verb = r;
                    e.forEach(function(e) {
                        Lt(t, e, i, function(e, t) {
                            if (v(e, Q.config.disableSelector)) {
                                m(e);
                                return;
                            }
                            he(r, n, e, t);
                        });
                    });
                }
            });
            return a;
        }
        function Lt(n, e, t, r) {
            if (e.sseEvent) Rt(n, r, e.sseEvent);
            else if (e.trigger === "revealed") {
                gt();
                ht(n, r, t, e);
                pt(n);
            } else if (e.trigger === "intersect") {
                var i = {};
                if (e.root) i.root = ue(n, e.root);
                if (e.threshold) i.threshold = parseFloat(e.threshold);
                var a = new IntersectionObserver(function(e) {
                    for(var t = 0; t < e.length; t++){
                        var r = e[t];
                        if (r.isIntersecting) {
                            ce(n, "intersect");
                            break;
                        }
                    }
                }, i);
                a.observe(n);
                ht(n, r, t, e);
            } else if (e.trigger === "load") {
                if (!ct(e, n, Wt("load", {
                    elt: n
                }))) qt(n, r, t, e.delay);
            } else if (e.pollInterval > 0) {
                t.polling = true;
                ot(n, r, e);
            } else ht(n, r, t, e);
        }
        function At(e) {
            if (!e.htmxExecuted && Q.config.allowScriptTags && (e.type === "text/javascript" || e.type === "module" || e.type === "")) {
                var t = re().createElement("script");
                oe(e.attributes, function(e) {
                    t.setAttribute(e.name, e.value);
                });
                t.textContent = e.textContent;
                t.async = false;
                if (Q.config.inlineScriptNonce) t.nonce = Q.config.inlineScriptNonce;
                var r = e.parentElement;
                try {
                    r.insertBefore(t, e);
                } catch (e) {
                    b(e);
                } finally{
                    if (e.parentElement) e.parentElement.removeChild(e);
                }
            }
        }
        function Nt(e) {
            if (h(e, "script")) At(e);
            oe(f(e, "script"), function(e) {
                At(e);
            });
        }
        function It(e) {
            var t = e.attributes;
            if (!t) return false;
            for(var r = 0; r < t.length; r++){
                var n = t[r].name;
                if (g(n, "hx-on:") || g(n, "data-hx-on:") || g(n, "hx-on-") || g(n, "data-hx-on-")) return true;
            }
            return false;
        }
        function kt(e) {
            var t = null;
            var r = [];
            if (It(e)) r.push(e);
            if (document.evaluate) {
                var n = document.evaluate('.//*[@*[ starts-with(name(), "hx-on:") or starts-with(name(), "data-hx-on:") or starts-with(name(), "hx-on-") or starts-with(name(), "data-hx-on-") ]]', e);
                while(t = n.iterateNext())r.push(t);
            } else if (typeof e.getElementsByTagName === "function") {
                var i = e.getElementsByTagName("*");
                for(var a = 0; a < i.length; a++)if (It(i[a])) r.push(i[a]);
            }
            return r;
        }
        function Pt(e) {
            if (e.querySelectorAll) {
                var t = ", [hx-boost] a, [data-hx-boost] a, a[hx-boost], a[data-hx-boost]";
                var r = e.querySelectorAll(i + t + ", form, [type='submit'], [hx-sse], [data-hx-sse], [hx-ws]," + " [data-hx-ws], [hx-ext], [data-hx-ext], [hx-trigger], [data-hx-trigger], [hx-on], [data-hx-on]");
                return r;
            } else return [];
        }
        function Mt(e) {
            var t = v(e.target, "button, input[type='submit']");
            var r = Dt(e);
            if (r) r.lastButtonClicked = t;
        }
        function Xt(e) {
            var t = Dt(e);
            if (t) t.lastButtonClicked = null;
        }
        function Dt(e) {
            var t = v(e.target, "button, input[type='submit']");
            if (!t) return;
            var r = p("#" + ee(t, "form")) || v(t, "form");
            if (!r) return;
            return ae(r);
        }
        function Ut(e) {
            e.addEventListener("click", Mt);
            e.addEventListener("focusin", Mt);
            e.addEventListener("focusout", Xt);
        }
        function Bt(e) {
            var t = Ye(e);
            var r = 0;
            for(var n = 0; n < t.length; n++){
                const i = t[n];
                if (i === "{") r++;
                else if (i === "}") r--;
            }
            return r;
        }
        function Ft(t, e, r) {
            var n = ae(t);
            if (!Array.isArray(n.onHandlers)) n.onHandlers = [];
            var i;
            var a = function(e) {
                return Tr(t, function() {
                    if (!i) i = new Function("event", r);
                    i.call(t, e);
                });
            };
            t.addEventListener(e, a);
            n.onHandlers.push({
                event: e,
                listener: a
            });
        }
        function Vt(e) {
            var t = te(e, "hx-on");
            if (t) {
                var r = {};
                var n = t.split("\n");
                var i = null;
                var a = 0;
                while(n.length > 0){
                    var o = n.shift();
                    var s = o.match(/^\s*([a-zA-Z:\-\.]+:)(.*)/);
                    if (a === 0 && s) {
                        o.split(":");
                        i = s[1].slice(0, -1);
                        r[i] = s[2];
                    } else r[i] += o;
                    a += Bt(o);
                }
                for(var l in r)Ft(e, l, r[l]);
            }
        }
        function jt(e) {
            Ae(e);
            for(var t = 0; t < e.attributes.length; t++){
                var r = e.attributes[t].name;
                var n = e.attributes[t].value;
                if (g(r, "hx-on") || g(r, "data-hx-on")) {
                    var i = r.indexOf("-on") + 3;
                    var a = r.slice(i, i + 1);
                    if (a === "-" || a === ":") {
                        var o = r.slice(i + 1);
                        if (g(o, ":")) o = "htmx" + o;
                        else if (g(o, "-")) o = "htmx:" + o.slice(1);
                        else if (g(o, "htmx-")) o = "htmx:" + o.slice(5);
                        Ft(e, o, n);
                    }
                }
            }
        }
        function _t(t) {
            if (v(t, Q.config.disableSelector)) {
                m(t);
                return;
            }
            var r = ae(t);
            if (r.initHash !== Le(t)) {
                Ne(t);
                r.initHash = Le(t);
                Vt(t);
                ce(t, "htmx:beforeProcessNode");
                if (t.value) r.lastValue = t.value;
                var e = it(t);
                var n = Ht(t, r, e);
                if (!n) {
                    if (ne(t, "hx-boost") === "true") lt(t, r, e);
                    else if (o(t, "hx-trigger")) e.forEach(function(e) {
                        Lt(t, e, r, function() {});
                    });
                }
                if (t.tagName === "FORM" || ee(t, "type") === "submit" && o(t, "form")) Ut(t);
                var i = te(t, "hx-sse");
                if (i) St(t, r, i);
                var a = te(t, "hx-ws");
                if (a) mt(t, r, a);
                ce(t, "htmx:afterProcessNode");
            }
        }
        function zt(e) {
            e = p(e);
            if (v(e, Q.config.disableSelector)) {
                m(e);
                return;
            }
            _t(e);
            oe(Pt(e), function(e) {
                _t(e);
            });
            oe(kt(e), jt);
        }
        function $t(e) {
            return e.replace(/([a-z0-9])([A-Z])/g, "$1-$2").toLowerCase();
        }
        function Wt(e, t) {
            var r;
            if (window.CustomEvent && typeof window.CustomEvent === "function") r = new CustomEvent(e, {
                bubbles: true,
                cancelable: true,
                detail: t
            });
            else {
                r = re().createEvent("CustomEvent");
                r.initCustomEvent(e, true, true, t);
            }
            return r;
        }
        function fe(e, t, r) {
            ce(e, t, le({
                error: t
            }, r));
        }
        function Gt(e) {
            return e === "htmx:afterProcessNode";
        }
        function R(e, t) {
            oe(Fr(e), function(e) {
                try {
                    t(e);
                } catch (e) {
                    b(e);
                }
            });
        }
        function b(e) {
            if (console.error) console.error(e);
            else if (console.log) console.log("ERROR: ", e);
        }
        function ce(e, t, r) {
            e = p(e);
            if (r == null) r = {};
            r["elt"] = e;
            var n = Wt(t, r);
            if (Q.logger && !Gt(t)) Q.logger(e, t, r);
            if (r.error) {
                b(r.error);
                ce(e, "htmx:error", {
                    errorInfo: r
                });
            }
            var i = e.dispatchEvent(n);
            var a = $t(t);
            if (i && a !== t) {
                var o = Wt(a, n.detail);
                i = i && e.dispatchEvent(o);
            }
            R(e, function(e) {
                i = i && e.onEvent(t, n) !== false && !n.defaultPrevented;
            });
            return i;
        }
        var Jt = location.pathname + location.search;
        function Zt() {
            var e = re().querySelector("[hx-history-elt],[data-hx-history-elt]");
            return e || re().body;
        }
        function Kt(e, t, r, n) {
            if (!U()) return;
            if (Q.config.historyCacheSize <= 0) {
                localStorage.removeItem("htmx-history-cache");
                return;
            }
            e = B(e);
            var i = E(localStorage.getItem("htmx-history-cache")) || [];
            for(var a = 0; a < i.length; a++)if (i[a].url === e) {
                i.splice(a, 1);
                break;
            }
            var o = {
                url: e,
                content: t,
                title: r,
                scroll: n
            };
            ce(re().body, "htmx:historyItemCreated", {
                item: o,
                cache: i
            });
            i.push(o);
            while(i.length > Q.config.historyCacheSize)i.shift();
            while(i.length > 0)try {
                localStorage.setItem("htmx-history-cache", JSON.stringify(i));
                break;
            } catch (e) {
                fe(re().body, "htmx:historyCacheError", {
                    cause: e,
                    cache: i
                });
                i.shift();
            }
        }
        function Yt(e) {
            if (!U()) return null;
            e = B(e);
            var t = E(localStorage.getItem("htmx-history-cache")) || [];
            for(var r = 0; r < t.length; r++){
                if (t[r].url === e) return t[r];
            }
            return null;
        }
        function Qt(e) {
            var t = Q.config.requestClass;
            var r = e.cloneNode(true);
            oe(f(r, "." + t), function(e) {
                n(e, t);
            });
            return r.innerHTML;
        }
        function er() {
            var e = Zt();
            var t = Jt || location.pathname + location.search;
            var r;
            try {
                r = re().querySelector('[hx-history="false" i],[data-hx-history="false" i]');
            } catch (e) {
                r = re().querySelector('[hx-history="false"],[data-hx-history="false"]');
            }
            if (!r) {
                ce(re().body, "htmx:beforeHistorySave", {
                    path: t,
                    historyElt: e
                });
                Kt(t, Qt(e), re().title, window.scrollY);
            }
            if (Q.config.historyEnabled) history.replaceState({
                htmx: true
            }, re().title, window.location.href);
        }
        function tr(e) {
            if (Q.config.getCacheBusterParam) {
                e = e.replace(/org\.htmx\.cache-buster=[^&]*&?/, "");
                if (G(e, "&") || G(e, "?")) e = e.slice(0, -1);
            }
            if (Q.config.historyEnabled) history.pushState({
                htmx: true
            }, "", e);
            Jt = e;
        }
        function rr(e) {
            if (Q.config.historyEnabled) history.replaceState({
                htmx: true
            }, "", e);
            Jt = e;
        }
        function nr(e) {
            oe(e, function(e) {
                e.call();
            });
        }
        function ir(a) {
            var e = new XMLHttpRequest;
            var o = {
                path: a,
                xhr: e
            };
            ce(re().body, "htmx:historyCacheMiss", o);
            e.open("GET", a, true);
            e.setRequestHeader("HX-Request", "true");
            e.setRequestHeader("HX-History-Restore-Request", "true");
            e.setRequestHeader("HX-Current-URL", re().location.href);
            e.onload = function() {
                if (this.status >= 200 && this.status < 400) {
                    ce(re().body, "htmx:historyCacheMissLoad", o);
                    var e = l(this.response);
                    e = e.querySelector("[hx-history-elt],[data-hx-history-elt]") || e;
                    var t = Zt();
                    var r = T(t);
                    var n = Ve(this.response);
                    if (n) {
                        var i = C("title");
                        if (i) i.innerHTML = n;
                        else window.document.title = n;
                    }
                    Ue(t, e, r);
                    nr(r.tasks);
                    Jt = a;
                    ce(re().body, "htmx:historyRestore", {
                        path: a,
                        cacheMiss: true,
                        serverResponse: this.response
                    });
                } else fe(re().body, "htmx:historyCacheMissLoadError", o);
            };
            e.send();
        }
        function ar(e) {
            er();
            e = e || location.pathname + location.search;
            var t = Yt(e);
            if (t) {
                var r = l(t.content);
                var n = Zt();
                var i = T(n);
                Ue(n, r, i);
                nr(i.tasks);
                document.title = t.title;
                setTimeout(function() {
                    window.scrollTo(0, t.scroll);
                }, 0);
                Jt = e;
                ce(re().body, "htmx:historyRestore", {
                    path: e,
                    item: t
                });
            } else if (Q.config.refreshOnHistoryMiss) window.location.reload(true);
            else ir(e);
        }
        function or(e) {
            var t = me(e, "hx-indicator");
            if (t == null) t = [
                e
            ];
            oe(t, function(e) {
                var t = ae(e);
                t.requestCount = (t.requestCount || 0) + 1;
                e.classList["add"].call(e.classList, Q.config.requestClass);
            });
            return t;
        }
        function sr(e) {
            var t = me(e, "hx-disabled-elt");
            if (t == null) t = [];
            oe(t, function(e) {
                var t = ae(e);
                t.requestCount = (t.requestCount || 0) + 1;
                e.setAttribute("disabled", "");
            });
            return t;
        }
        function lr(e, t) {
            oe(e, function(e) {
                var t = ae(e);
                t.requestCount = (t.requestCount || 0) - 1;
                if (t.requestCount === 0) e.classList["remove"].call(e.classList, Q.config.requestClass);
            });
            oe(t, function(e) {
                var t = ae(e);
                t.requestCount = (t.requestCount || 0) - 1;
                if (t.requestCount === 0) e.removeAttribute("disabled");
            });
        }
        function ur(e, t) {
            for(var r = 0; r < e.length; r++){
                var n = e[r];
                if (n.isSameNode(t)) return true;
            }
            return false;
        }
        function fr(e) {
            if (e.name === "" || e.name == null || e.disabled || v(e, "fieldset[disabled]")) return false;
            if (e.type === "button" || e.type === "submit" || e.tagName === "image" || e.tagName === "reset" || e.tagName === "file") return false;
            if (e.type === "checkbox" || e.type === "radio") return e.checked;
            return true;
        }
        function cr(e, t, r) {
            if (e != null && t != null) {
                var n = r[e];
                if (n === undefined) r[e] = t;
                else if (Array.isArray(n)) {
                    if (Array.isArray(t)) r[e] = n.concat(t);
                    else n.push(t);
                } else if (Array.isArray(t)) r[e] = [
                    n
                ].concat(t);
                else r[e] = [
                    n,
                    t
                ];
            }
        }
        function hr(t, r, n, e, i) {
            if (e == null || ur(t, e)) return;
            else t.push(e);
            if (fr(e)) {
                var a = ee(e, "name");
                var o = e.value;
                if (e.multiple && e.tagName === "SELECT") o = M(e.querySelectorAll("option:checked")).map(function(e) {
                    return e.value;
                });
                if (e.files) o = M(e.files);
                cr(a, o, r);
                if (i) vr(e, n);
            }
            if (h(e, "form")) {
                var s = e.elements;
                oe(s, function(e) {
                    hr(t, r, n, e, i);
                });
            }
        }
        function vr(e, t) {
            if (e.willValidate) {
                ce(e, "htmx:validation:validate");
                if (!e.checkValidity()) {
                    t.push({
                        elt: e,
                        message: e.validationMessage,
                        validity: e.validity
                    });
                    ce(e, "htmx:validation:failed", {
                        message: e.validationMessage,
                        validity: e.validity
                    });
                }
            }
        }
        function dr(e, t) {
            var r = [];
            var n = {};
            var i = {};
            var a = [];
            var o = ae(e);
            if (o.lastButtonClicked && !se(o.lastButtonClicked)) o.lastButtonClicked = null;
            var s = h(e, "form") && e.noValidate !== true || te(e, "hx-validate") === "true";
            if (o.lastButtonClicked) s = s && o.lastButtonClicked.formNoValidate !== true;
            if (t !== "get") hr(r, i, a, v(e, "form"), s);
            hr(r, n, a, e, s);
            if (o.lastButtonClicked || e.tagName === "BUTTON" || e.tagName === "INPUT" && ee(e, "type") === "submit") {
                var l = o.lastButtonClicked || e;
                var u = ee(l, "name");
                cr(u, l.value, i);
            }
            var f = me(e, "hx-include");
            oe(f, function(e) {
                hr(r, n, a, e, s);
                if (!h(e, "form")) oe(e.querySelectorAll(rt), function(e) {
                    hr(r, n, a, e, s);
                });
            });
            n = le(n, i);
            return {
                errors: a,
                values: n
            };
        }
        function gr(e, t, r) {
            if (e !== "") e += "&";
            if (String(r) === "[object Object]") r = JSON.stringify(r);
            var n = encodeURIComponent(r);
            e += encodeURIComponent(t) + "=" + n;
            return e;
        }
        function pr(e) {
            var t = "";
            for(var r in e)if (e.hasOwnProperty(r)) {
                var n = e[r];
                if (Array.isArray(n)) oe(n, function(e) {
                    t = gr(t, r, e);
                });
                else t = gr(t, r, n);
            }
            return t;
        }
        function mr(e) {
            var t = new FormData;
            for(var r in e)if (e.hasOwnProperty(r)) {
                var n = e[r];
                if (Array.isArray(n)) oe(n, function(e) {
                    t.append(r, e);
                });
                else t.append(r, n);
            }
            return t;
        }
        function xr(e, t, r) {
            var n = {
                "HX-Request": "true",
                "HX-Trigger": ee(e, "id"),
                "HX-Trigger-Name": ee(e, "name"),
                "HX-Target": te(t, "id"),
                "HX-Current-URL": re().location.href
            };
            Rr(e, "hx-headers", false, n);
            if (r !== undefined) n["HX-Prompt"] = r;
            if (ae(e).boosted) n["HX-Boosted"] = "true";
            return n;
        }
        function yr(t, e) {
            var r = ne(e, "hx-params");
            if (r) {
                if (r === "none") return {};
                else if (r === "*") return t;
                else if (r.indexOf("not ") === 0) {
                    oe(r.substr(4).split(","), function(e) {
                        e = e.trim();
                        delete t[e];
                    });
                    return t;
                } else {
                    var n = {};
                    oe(r.split(","), function(e) {
                        e = e.trim();
                        n[e] = t[e];
                    });
                    return n;
                }
            } else return t;
        }
        function br(e) {
            return ee(e, "href") && ee(e, "href").indexOf("#") >= 0;
        }
        function wr(e, t) {
            var r = t ? t : ne(e, "hx-swap");
            var n = {
                swapStyle: ae(e).boosted ? "innerHTML" : Q.config.defaultSwapStyle,
                swapDelay: Q.config.defaultSwapDelay,
                settleDelay: Q.config.defaultSettleDelay
            };
            if (Q.config.scrollIntoViewOnBoost && ae(e).boosted && !br(e)) n["show"] = "top";
            if (r) {
                var i = D(r);
                if (i.length > 0) for(var a = 0; a < i.length; a++){
                    var o = i[a];
                    if (o.indexOf("swap:") === 0) n["swapDelay"] = d(o.substr(5));
                    else if (o.indexOf("settle:") === 0) n["settleDelay"] = d(o.substr(7));
                    else if (o.indexOf("transition:") === 0) n["transition"] = o.substr(11) === "true";
                    else if (o.indexOf("ignoreTitle:") === 0) n["ignoreTitle"] = o.substr(12) === "true";
                    else if (o.indexOf("scroll:") === 0) {
                        var s = o.substr(7);
                        var l = s.split(":");
                        var u = l.pop();
                        var f = l.length > 0 ? l.join(":") : null;
                        n["scroll"] = u;
                        n["scrollTarget"] = f;
                    } else if (o.indexOf("show:") === 0) {
                        var c = o.substr(5);
                        var l = c.split(":");
                        var h = l.pop();
                        var f = l.length > 0 ? l.join(":") : null;
                        n["show"] = h;
                        n["showTarget"] = f;
                    } else if (o.indexOf("focus-scroll:") === 0) {
                        var v = o.substr(13);
                        n["focusScroll"] = v == "true";
                    } else if (a == 0) n["swapStyle"] = o;
                    else b("Unknown modifier in hx-swap: " + o);
                }
            }
            return n;
        }
        function Sr(e) {
            return ne(e, "hx-encoding") === "multipart/form-data" || h(e, "form") && ee(e, "enctype") === "multipart/form-data";
        }
        function Er(t, r, n) {
            var i = null;
            R(r, function(e) {
                if (i == null) i = e.encodeParameters(t, n, r);
            });
            if (i != null) return i;
            else {
                if (Sr(r)) return mr(n);
                else return pr(n);
            }
        }
        function T(e) {
            return {
                tasks: [],
                elts: [
                    e
                ]
            };
        }
        function Cr(e, t) {
            var r = e[0];
            var n = e[e.length - 1];
            if (t.scroll) {
                var i = null;
                if (t.scrollTarget) i = ue(r, t.scrollTarget);
                if (t.scroll === "top" && (r || i)) {
                    i = i || r;
                    i.scrollTop = 0;
                }
                if (t.scroll === "bottom" && (n || i)) {
                    i = i || n;
                    i.scrollTop = i.scrollHeight;
                }
            }
            if (t.show) {
                var i = null;
                if (t.showTarget) {
                    var a = t.showTarget;
                    if (t.showTarget === "window") a = "body";
                    i = ue(r, a);
                }
                if (t.show === "top" && (r || i)) {
                    i = i || r;
                    i.scrollIntoView({
                        block: "start",
                        behavior: Q.config.scrollBehavior
                    });
                }
                if (t.show === "bottom" && (n || i)) {
                    i = i || n;
                    i.scrollIntoView({
                        block: "end",
                        behavior: Q.config.scrollBehavior
                    });
                }
            }
        }
        function Rr(e, t, r, n) {
            if (n == null) n = {};
            if (e == null) return n;
            var i = te(e, t);
            if (i) {
                var a = i.trim();
                var o = r;
                if (a === "unset") return null;
                if (a.indexOf("javascript:") === 0) {
                    a = a.substr(11);
                    o = true;
                } else if (a.indexOf("js:") === 0) {
                    a = a.substr(3);
                    o = true;
                }
                if (a.indexOf("{") !== 0) a = "{" + a + "}";
                var s;
                if (o) s = Tr(e, function() {
                    return Function("return (" + a + ")")();
                }, {});
                else s = E(a);
                for(var l in s){
                    if (s.hasOwnProperty(l)) {
                        if (n[l] == null) n[l] = s[l];
                    }
                }
            }
            return Rr(u(e), t, r, n);
        }
        function Tr(e, t, r) {
            if (Q.config.allowEval) return t();
            else {
                fe(e, "htmx:evalDisallowedError");
                return r;
            }
        }
        function Or(e, t) {
            return Rr(e, "hx-vars", true, t);
        }
        function qr(e, t) {
            return Rr(e, "hx-vals", false, t);
        }
        function Hr(e) {
            return le(Or(e), qr(e));
        }
        function Lr(t, r, n) {
            if (n !== null) try {
                t.setRequestHeader(r, n);
            } catch (e) {
                t.setRequestHeader(r, encodeURIComponent(n));
                t.setRequestHeader(r + "-URI-AutoEncoded", "true");
            }
        }
        function Ar(t) {
            if (t.responseURL && typeof URL !== "undefined") try {
                var e = new URL(t.responseURL);
                return e.pathname + e.search;
            } catch (e) {
                fe(re().body, "htmx:badResponseUrl", {
                    url: t.responseURL
                });
            }
        }
        function O(e, t) {
            return t.test(e.getAllResponseHeaders());
        }
        function Nr(e, t, r) {
            e = e.toLowerCase();
            if (r) {
                if (r instanceof Element || I(r, "String")) return he(e, t, null, null, {
                    targetOverride: p(r),
                    returnPromise: true
                });
                else return he(e, t, p(r.source), r.event, {
                    handler: r.handler,
                    headers: r.headers,
                    values: r.values,
                    targetOverride: p(r.target),
                    swapOverride: r.swap,
                    select: r.select,
                    returnPromise: true
                });
            } else return he(e, t, null, null, {
                returnPromise: true
            });
        }
        function Ir(e) {
            var t = [];
            while(e){
                t.push(e);
                e = e.parentElement;
            }
            return t;
        }
        function kr(e, t, r) {
            var n;
            var i;
            if (typeof URL === "function") {
                i = new URL(t, document.location.href);
                var a = document.location.origin;
                n = a === i.origin;
            } else {
                i = t;
                n = g(t, document.location.origin);
            }
            if (Q.config.selfRequestsOnly) {
                if (!n) return false;
            }
            return ce(e, "htmx:validateUrl", le({
                url: i,
                sameHost: n
            }, r));
        }
        function he(t, r, n, i, a, e) {
            var o = null;
            var s = null;
            a = a != null ? a : {};
            if (a.returnPromise && typeof Promise !== "undefined") var l = new Promise(function(e, t) {
                o = e;
                s = t;
            });
            if (n == null) n = re().body;
            var M = a.handler || Mr;
            var X = a.select || null;
            if (!se(n)) {
                ie(o);
                return l;
            }
            var u = a.targetOverride || ye(n);
            if (u == null || u == pe) {
                fe(n, "htmx:targetError", {
                    target: te(n, "hx-target")
                });
                ie(s);
                return l;
            }
            var f = ae(n);
            var c = f.lastButtonClicked;
            if (c) {
                var h = ee(c, "formaction");
                if (h != null) r = h;
                var v = ee(c, "formmethod");
                if (v != null) {
                    if (v.toLowerCase() !== "dialog") t = v;
                }
            }
            var d = ne(n, "hx-confirm");
            if (e === undefined) {
                var D = function(e) {
                    return he(t, r, n, i, a, !!e);
                };
                var U = {
                    target: u,
                    elt: n,
                    path: r,
                    verb: t,
                    triggeringEvent: i,
                    etc: a,
                    issueRequest: D,
                    question: d
                };
                if (ce(n, "htmx:confirm", U) === false) {
                    ie(o);
                    return l;
                }
            }
            var g = n;
            var p = ne(n, "hx-sync");
            var m = null;
            var x = false;
            if (p) {
                var B = p.split(":");
                var F = B[0].trim();
                if (F === "this") g = xe(n, "hx-sync");
                else g = ue(n, F);
                p = (B[1] || "drop").trim();
                f = ae(g);
                if (p === "drop" && f.xhr && f.abortable !== true) {
                    ie(o);
                    return l;
                } else if (p === "abort") {
                    if (f.xhr) {
                        ie(o);
                        return l;
                    } else x = true;
                } else if (p === "replace") ce(g, "htmx:abort");
                else if (p.indexOf("queue") === 0) {
                    var V = p.split(" ");
                    m = (V[1] || "last").trim();
                }
            }
            if (f.xhr) {
                if (f.abortable) ce(g, "htmx:abort");
                else {
                    if (m == null) {
                        if (i) {
                            var y = ae(i);
                            if (y && y.triggerSpec && y.triggerSpec.queue) m = y.triggerSpec.queue;
                        }
                        if (m == null) m = "last";
                    }
                    if (f.queuedRequests == null) f.queuedRequests = [];
                    if (m === "first" && f.queuedRequests.length === 0) f.queuedRequests.push(function() {
                        he(t, r, n, i, a);
                    });
                    else if (m === "all") f.queuedRequests.push(function() {
                        he(t, r, n, i, a);
                    });
                    else if (m === "last") {
                        f.queuedRequests = [];
                        f.queuedRequests.push(function() {
                            he(t, r, n, i, a);
                        });
                    }
                    ie(o);
                    return l;
                }
            }
            var b = new XMLHttpRequest;
            f.xhr = b;
            f.abortable = x;
            var w = function() {
                f.xhr = null;
                f.abortable = false;
                if (f.queuedRequests != null && f.queuedRequests.length > 0) {
                    var e = f.queuedRequests.shift();
                    e();
                }
            };
            var j = ne(n, "hx-prompt");
            if (j) {
                var S = prompt(j);
                if (S === null || !ce(n, "htmx:prompt", {
                    prompt: S,
                    target: u
                })) {
                    ie(o);
                    w();
                    return l;
                }
            }
            if (d && !e) {
                if (!confirm(d)) {
                    ie(o);
                    w();
                    return l;
                }
            }
            var E = xr(n, u, S);
            if (t !== "get" && !Sr(n)) E["Content-Type"] = "application/x-www-form-urlencoded";
            if (a.headers) E = le(E, a.headers);
            var _ = dr(n, t);
            var C = _.errors;
            var R = _.values;
            if (a.values) R = le(R, a.values);
            var z = Hr(n);
            var $ = le(R, z);
            var T = yr($, n);
            if (Q.config.getCacheBusterParam && t === "get") T["org.htmx.cache-buster"] = ee(u, "id") || "true";
            if (r == null || r === "") r = re().location.href;
            var O = Rr(n, "hx-request");
            var W = ae(n).boosted;
            var q = Q.config.methodsThatUseUrlParams.indexOf(t) >= 0;
            var H = {
                boosted: W,
                useUrlParams: q,
                parameters: T,
                unfilteredParameters: $,
                headers: E,
                target: u,
                verb: t,
                errors: C,
                withCredentials: a.credentials || O.credentials || Q.config.withCredentials,
                timeout: a.timeout || O.timeout || Q.config.timeout,
                path: r,
                triggeringEvent: i
            };
            if (!ce(n, "htmx:configRequest", H)) {
                ie(o);
                w();
                return l;
            }
            r = H.path;
            t = H.verb;
            E = H.headers;
            T = H.parameters;
            C = H.errors;
            q = H.useUrlParams;
            if (C && C.length > 0) {
                ce(n, "htmx:validation:halted", H);
                ie(o);
                w();
                return l;
            }
            var G = r.split("#");
            var J = G[0];
            var L = G[1];
            var A = r;
            if (q) {
                A = J;
                var Z = Object.keys(T).length !== 0;
                if (Z) {
                    if (A.indexOf("?") < 0) A += "?";
                    else A += "&";
                    A += pr(T);
                    if (L) A += "#" + L;
                }
            }
            if (!kr(n, A, H)) {
                fe(n, "htmx:invalidPath", H);
                ie(s);
                return l;
            }
            b.open(t.toUpperCase(), A, true);
            b.overrideMimeType("text/html");
            b.withCredentials = H.withCredentials;
            b.timeout = H.timeout;
            if (O.noHeaders) ;
            else {
                for(var N in E)if (E.hasOwnProperty(N)) {
                    var K = E[N];
                    Lr(b, N, K);
                }
            }
            var I = {
                xhr: b,
                target: u,
                requestConfig: H,
                etc: a,
                boosted: W,
                select: X,
                pathInfo: {
                    requestPath: r,
                    finalRequestPath: A,
                    anchor: L
                }
            };
            b.onload = function() {
                try {
                    var e = Ir(n);
                    I.pathInfo.responsePath = Ar(b);
                    M(n, I);
                    lr(k, P);
                    ce(n, "htmx:afterRequest", I);
                    ce(n, "htmx:afterOnLoad", I);
                    if (!se(n)) {
                        var t = null;
                        while(e.length > 0 && t == null){
                            var r = e.shift();
                            if (se(r)) t = r;
                        }
                        if (t) {
                            ce(t, "htmx:afterRequest", I);
                            ce(t, "htmx:afterOnLoad", I);
                        }
                    }
                    ie(o);
                    w();
                } catch (e) {
                    fe(n, "htmx:onLoadError", le({
                        error: e
                    }, I));
                    throw e;
                }
            };
            b.onerror = function() {
                lr(k, P);
                fe(n, "htmx:afterRequest", I);
                fe(n, "htmx:sendError", I);
                ie(s);
                w();
            };
            b.onabort = function() {
                lr(k, P);
                fe(n, "htmx:afterRequest", I);
                fe(n, "htmx:sendAbort", I);
                ie(s);
                w();
            };
            b.ontimeout = function() {
                lr(k, P);
                fe(n, "htmx:afterRequest", I);
                fe(n, "htmx:timeout", I);
                ie(s);
                w();
            };
            if (!ce(n, "htmx:beforeRequest", I)) {
                ie(o);
                w();
                return l;
            }
            var k = or(n);
            var P = sr(n);
            oe([
                "loadstart",
                "loadend",
                "progress",
                "abort"
            ], function(t) {
                oe([
                    b,
                    b.upload
                ], function(e) {
                    e.addEventListener(t, function(e) {
                        ce(n, "htmx:xhr:" + t, {
                            lengthComputable: e.lengthComputable,
                            loaded: e.loaded,
                            total: e.total
                        });
                    });
                });
            });
            ce(n, "htmx:beforeSend", I);
            var Y = q ? null : Er(b, n, T);
            b.send(Y);
            return l;
        }
        function Pr(e, t) {
            var r = t.xhr;
            var n = null;
            var i = null;
            if (O(r, /HX-Push:/i)) {
                n = r.getResponseHeader("HX-Push");
                i = "push";
            } else if (O(r, /HX-Push-Url:/i)) {
                n = r.getResponseHeader("HX-Push-Url");
                i = "push";
            } else if (O(r, /HX-Replace-Url:/i)) {
                n = r.getResponseHeader("HX-Replace-Url");
                i = "replace";
            }
            if (n) {
                if (n === "false") return {};
                else return {
                    type: i,
                    path: n
                };
            }
            var a = t.pathInfo.finalRequestPath;
            var o = t.pathInfo.responsePath;
            var s = ne(e, "hx-push-url");
            var l = ne(e, "hx-replace-url");
            var u = ae(e).boosted;
            var f = null;
            var c = null;
            if (s) {
                f = "push";
                c = s;
            } else if (l) {
                f = "replace";
                c = l;
            } else if (u) {
                f = "push";
                c = o || a;
            }
            if (c) {
                if (c === "false") return {};
                if (c === "true") c = o || a;
                if (t.pathInfo.anchor && c.indexOf("#") === -1) c = c + "#" + t.pathInfo.anchor;
                return {
                    type: f,
                    path: c
                };
            } else return {};
        }
        function Mr(l, u) {
            var f = u.xhr;
            var c = u.target;
            var e = u.etc;
            var t = u.requestConfig;
            var h = u.select;
            if (!ce(l, "htmx:beforeOnLoad", u)) return;
            if (O(f, /HX-Trigger:/i)) _e(f, "HX-Trigger", l);
            if (O(f, /HX-Location:/i)) {
                er();
                var r = f.getResponseHeader("HX-Location");
                var v;
                if (r.indexOf("{") === 0) {
                    v = E(r);
                    r = v["path"];
                    delete v["path"];
                }
                Nr("GET", r, v).then(function() {
                    tr(r);
                });
                return;
            }
            var n = O(f, /HX-Refresh:/i) && "true" === f.getResponseHeader("HX-Refresh");
            if (O(f, /HX-Redirect:/i)) {
                location.href = f.getResponseHeader("HX-Redirect");
                n && location.reload();
                return;
            }
            if (n) {
                location.reload();
                return;
            }
            if (O(f, /HX-Retarget:/i)) {
                if (f.getResponseHeader("HX-Retarget") === "this") u.target = l;
                else u.target = ue(l, f.getResponseHeader("HX-Retarget"));
            }
            var d = Pr(l, u);
            var i = f.status >= 200 && f.status < 400 && f.status !== 204;
            var g = f.response;
            var a = f.status >= 400;
            var p = Q.config.ignoreTitle;
            var o = le({
                shouldSwap: i,
                serverResponse: g,
                isError: a,
                ignoreTitle: p
            }, u);
            if (!ce(c, "htmx:beforeSwap", o)) return;
            c = o.target;
            g = o.serverResponse;
            a = o.isError;
            p = o.ignoreTitle;
            u.target = c;
            u.failed = a;
            u.successful = !a;
            if (o.shouldSwap) {
                if (f.status === 286) at(l);
                R(l, function(e) {
                    g = e.transformResponse(g, f, l);
                });
                if (d.type) er();
                var s = e.swapOverride;
                if (O(f, /HX-Reswap:/i)) s = f.getResponseHeader("HX-Reswap");
                var v = wr(l, s);
                if (v.hasOwnProperty("ignoreTitle")) p = v.ignoreTitle;
                c.classList.add(Q.config.swappingClass);
                var m = null;
                var x = null;
                var y = function() {
                    try {
                        var e = document.activeElement;
                        var t = {};
                        try {
                            t = {
                                elt: e,
                                start: e ? e.selectionStart : null,
                                end: e ? e.selectionEnd : null
                            };
                        } catch (e) {}
                        var r;
                        if (h) r = h;
                        if (O(f, /HX-Reselect:/i)) r = f.getResponseHeader("HX-Reselect");
                        if (d.type) {
                            ce(re().body, "htmx:beforeHistoryUpdate", le({
                                history: d
                            }, u));
                            if (d.type === "push") {
                                tr(d.path);
                                ce(re().body, "htmx:pushedIntoHistory", {
                                    path: d.path
                                });
                            } else {
                                rr(d.path);
                                ce(re().body, "htmx:replacedInHistory", {
                                    path: d.path
                                });
                            }
                        }
                        var n = T(c);
                        je(v.swapStyle, c, l, g, n, r);
                        if (t.elt && !se(t.elt) && ee(t.elt, "id")) {
                            var i = document.getElementById(ee(t.elt, "id"));
                            var a = {
                                preventScroll: v.focusScroll !== undefined ? !v.focusScroll : !Q.config.defaultFocusScroll
                            };
                            if (i) {
                                if (t.start && i.setSelectionRange) try {
                                    i.setSelectionRange(t.start, t.end);
                                } catch (e) {}
                                i.focus(a);
                            }
                        }
                        c.classList.remove(Q.config.swappingClass);
                        oe(n.elts, function(e) {
                            if (e.classList) e.classList.add(Q.config.settlingClass);
                            ce(e, "htmx:afterSwap", u);
                        });
                        if (O(f, /HX-Trigger-After-Swap:/i)) {
                            var o = l;
                            if (!se(l)) o = re().body;
                            _e(f, "HX-Trigger-After-Swap", o);
                        }
                        var s = function() {
                            oe(n.tasks, function(e) {
                                e.call();
                            });
                            oe(n.elts, function(e) {
                                if (e.classList) e.classList.remove(Q.config.settlingClass);
                                ce(e, "htmx:afterSettle", u);
                            });
                            if (u.pathInfo.anchor) {
                                var e = re().getElementById(u.pathInfo.anchor);
                                if (e) e.scrollIntoView({
                                    block: "start",
                                    behavior: "auto"
                                });
                            }
                            if (n.title && !p) {
                                var t = C("title");
                                if (t) t.innerHTML = n.title;
                                else window.document.title = n.title;
                            }
                            Cr(n.elts, v);
                            if (O(f, /HX-Trigger-After-Settle:/i)) {
                                var r = l;
                                if (!se(l)) r = re().body;
                                _e(f, "HX-Trigger-After-Settle", r);
                            }
                            ie(m);
                        };
                        if (v.settleDelay > 0) setTimeout(s, v.settleDelay);
                        else s();
                    } catch (e) {
                        fe(l, "htmx:swapError", u);
                        ie(x);
                        throw e;
                    }
                };
                var b = Q.config.globalViewTransitions;
                if (v.hasOwnProperty("transition")) b = v.transition;
                if (b && ce(l, "htmx:beforeTransition", u) && typeof Promise !== "undefined" && document.startViewTransition) {
                    var w = new Promise(function(e, t) {
                        m = e;
                        x = t;
                    });
                    var S = y;
                    y = function() {
                        document.startViewTransition(function() {
                            S();
                            return w;
                        });
                    };
                }
                if (v.swapDelay > 0) setTimeout(y, v.swapDelay);
                else y();
            }
            if (a) fe(l, "htmx:responseError", le({
                error: "Response Status Error Code " + f.status + " from " + u.pathInfo.requestPath
            }, u));
        }
        var Xr = {};
        function Dr() {
            return {
                init: function(e) {
                    return null;
                },
                onEvent: function(e, t) {
                    return true;
                },
                transformResponse: function(e, t, r) {
                    return e;
                },
                isInlineSwap: function(e) {
                    return false;
                },
                handleSwap: function(e, t, r, n) {
                    return false;
                },
                encodeParameters: function(e, t, r) {
                    return null;
                }
            };
        }
        function Ur(e, t) {
            if (t.init) t.init(r);
            Xr[e] = le(Dr(), t);
        }
        function Br(e) {
            delete Xr[e];
        }
        function Fr(e, r, n) {
            if (e == undefined) return r;
            if (r == undefined) r = [];
            if (n == undefined) n = [];
            var t = te(e, "hx-ext");
            if (t) oe(t.split(","), function(e) {
                e = e.replace(/ /g, "");
                if (e.slice(0, 7) == "ignore:") {
                    n.push(e.slice(7));
                    return;
                }
                if (n.indexOf(e) < 0) {
                    var t = Xr[e];
                    if (t && r.indexOf(t) < 0) r.push(t);
                }
            });
            return Fr(u(e), r, n);
        }
        var Vr = false;
        re().addEventListener("DOMContentLoaded", function() {
            Vr = true;
        });
        function jr(e) {
            if (Vr || re().readyState === "complete") e();
            else re().addEventListener("DOMContentLoaded", e);
        }
        function _r() {
            if (Q.config.includeIndicatorStyles !== false) re().head.insertAdjacentHTML("beforeend", "<style>                      ." + Q.config.indicatorClass + "{opacity:0}                      ." + Q.config.requestClass + " ." + Q.config.indicatorClass + "{opacity:1; transition: opacity 200ms ease-in;}                      ." + Q.config.requestClass + "." + Q.config.indicatorClass + "{opacity:1; transition: opacity 200ms ease-in;}                    </style>");
        }
        function zr() {
            var e = re().querySelector('meta[name="htmx-config"]');
            if (e) return E(e.content);
            else return null;
        }
        function $r() {
            var e = zr();
            if (e) Q.config = le(Q.config, e);
        }
        jr(function() {
            $r();
            _r();
            var e = re().body;
            zt(e);
            var t = re().querySelectorAll("[hx-trigger='restored'],[data-hx-trigger='restored']");
            e.addEventListener("htmx:abort", function(e) {
                var t = e.target;
                var r = ae(t);
                if (r && r.xhr) r.xhr.abort();
            });
            const r = window.onpopstate ? window.onpopstate.bind(window) : null;
            window.onpopstate = function(e) {
                if (e.state && e.state.htmx) {
                    ar();
                    oe(t, function(e) {
                        ce(e, "htmx:restored", {
                            document: re(),
                            triggerEvent: ce
                        });
                    });
                } else if (r) r(e);
            };
            setTimeout(function() {
                ce(e, "htmx:load", {});
                e = null;
            }, 0);
        });
        return Q;
    }();
});

});

parcelRegister("2Gj6O", function(module, exports) {
"use strict";
Object.defineProperty(module.exports, "__esModule", {
    value: true
});
Object.defineProperty(module.exports, "default", {
    enumerable: true,
    get: function() {
        return $1f3e53a042807708$var$_default;
    }
});
function $1f3e53a042807708$var$createPlugin(plugin, config) {
    return {
        handler: plugin,
        config: config
    };
}
$1f3e53a042807708$var$createPlugin.withOptions = function(pluginFunction, configFunction = ()=>({})) {
    const optionsFunction = function(options) {
        return {
            __options: options,
            handler: pluginFunction(options),
            config: configFunction(options)
        };
    };
    optionsFunction.__isOptionsFunction = true;
    // Expose plugin dependencies so that `object-hash` returns a different
    // value if anything here changes, to ensure a rebuild is triggered.
    optionsFunction.__pluginFunction = pluginFunction;
    optionsFunction.__configFunction = configFunction;
    return optionsFunction;
};
const $1f3e53a042807708$var$_default = $1f3e53a042807708$var$createPlugin;

});

parcelRequire("eZIU4");
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
function $8c83eaf28779ff46$var$handleError(error2, el, expression) {
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
    return (receiver = ()=>{}, { scope: scope2 = {}, params: params = [] } = {})=>{
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
    return (receiver = ()=>{}, { scope: scope2 = {}, params: params = [] } = {})=>{
        func.result = void 0;
        func.finished = false;
        let completeScope = $8c83eaf28779ff46$var$mergeProxies([
            scope2,
            ...dataStack
        ]);
        if (typeof func === "function") {
            let promise = func(func, completeScope).catch((error2)=>$8c83eaf28779ff46$var$handleError(error2, el, expression));
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
    var timeout;
    return function() {
        var context = this, args = arguments;
        var later = function() {
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
    version: "3.14.7",
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
    interceptClone: $8c83eaf28779ff46$var$interceptClone,
    addScopeToNode: $8c83eaf28779ff46$var$addScopeToNode,
    deferMutations: $8c83eaf28779ff46$var$deferMutations,
    mapAttributes: $8c83eaf28779ff46$var$mapAttributes,
    evaluateLater: $8c83eaf28779ff46$var$evaluateLater,
    interceptInit: $8c83eaf28779ff46$var$interceptInit,
    setEvaluator: $8c83eaf28779ff46$var$setEvaluator,
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
            "passive"
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
    var event = el.tagName.toLowerCase() === "select" || [
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
    let value = storage.getItem(key, storage);
    if (value === void 0) return;
    return JSON.parse(value);
}
function $9b2f94dab0f686ea$var$storageSet(key, value, storage) {
    storage.setItem(key, JSON.stringify(value));
}
// packages/persist/builds/module.js
var $9b2f94dab0f686ea$export$2e2bcd8739ae039 = $9b2f94dab0f686ea$export$9a6132153fba2e0;


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
        if (last && last.obj && typeof last.obj[`${last.k}.${e}`] !== 'undefined') last.obj = undefined;
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
const $1bac384020b50752$var$deepFind = function(obj, path) {
    let keySeparator = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : '.';
    if (!obj) return undefined;
    if (obj[path]) return obj[path];
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
const $1bac384020b50752$var$getCleanedCode = (code)=>code && code.replace('_', '-');
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
        if (console && console[type]) console[type].apply(console, args);
    }
};
class $1bac384020b50752$var$Logger {
    constructor(concreteLogger){
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
        this.init(concreteLogger, options);
    }
    init(concreteLogger) {
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
        this.prefix = options.prefix || 'i18next:';
        this.logger = concreteLogger || $1bac384020b50752$var$consoleLogger;
        this.options = options;
        this.debug = options.debug;
    }
    log() {
        for(var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++)args[_key] = arguments[_key];
        return this.forward(args, 'log', '', true);
    }
    warn() {
        for(var _len2 = arguments.length, args = new Array(_len2), _key2 = 0; _key2 < _len2; _key2++)args[_key2] = arguments[_key2];
        return this.forward(args, 'warn', '', true);
    }
    error() {
        for(var _len3 = arguments.length, args = new Array(_len3), _key3 = 0; _key3 < _len3; _key3++)args[_key3] = arguments[_key3];
        return this.forward(args, 'error', '');
    }
    deprecate() {
        for(var _len4 = arguments.length, args = new Array(_len4), _key4 = 0; _key4 < _len4; _key4++)args[_key4] = arguments[_key4];
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
    emit(event) {
        for(var _len = arguments.length, args = new Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++)args[_key - 1] = arguments[_key];
        if (this.observers[event]) {
            const cloned = Array.from(this.observers[event].entries());
            cloned.forEach((_ref)=>{
                let [observer, numTimesAdded] = _ref;
                for(let i = 0; i < numTimesAdded; i++)observer(...args);
            });
        }
        if (this.observers['*']) {
            const cloned = Array.from(this.observers['*'].entries());
            cloned.forEach((_ref2)=>{
                let [observer, numTimesAdded] = _ref2;
                for(let i = 0; i < numTimesAdded; i++)observer.apply(observer, [
                    event,
                    ...args
                ]);
            });
        }
    }
}
class $1bac384020b50752$var$ResourceStore extends $1bac384020b50752$var$EventEmitter {
    constructor(data){
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {
            ns: [
                'translation'
            ],
            defaultNS: 'translation'
        };
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
    getResource(lng, ns, key) {
        let options = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : {};
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
        return $1bac384020b50752$var$deepFind(this.data && this.data[lng] && this.data[lng][ns], key, keySeparator);
    }
    addResource(lng, ns, key, value) {
        let options = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : {
            silent: false
        };
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
    addResources(lng, ns, resources) {
        let options = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : {
            silent: false
        };
        for(const m in resources)if ($1bac384020b50752$var$isString(resources[m]) || Array.isArray(resources[m])) this.addResource(lng, ns, m, resources[m], {
            silent: true
        });
        if (!options.silent) this.emit('added', lng, ns, resources);
    }
    addResourceBundle(lng, ns, resources, deep, overwrite) {
        let options = arguments.length > 5 && arguments[5] !== undefined ? arguments[5] : {
            silent: false,
            skipCopy: false
        };
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
        if (this.options.compatibilityAPI === 'v1') return {
            ...this.getResource(lng, ns)
        };
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
            if (this.processors[processor]) value = this.processors[processor].process(value, key, options, translator);
        });
        return value;
    }
};
const $1bac384020b50752$var$checkedLoadedFor = {};
class $1bac384020b50752$var$Translator extends $1bac384020b50752$var$EventEmitter {
    constructor(services){
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
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
    exists(key) {
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {
            interpolation: {}
        };
        if (key === undefined || key === null) return false;
        const resolved = this.resolve(key, options);
        return resolved && resolved.res !== undefined;
    }
    extractFromKey(key, options) {
        let nsSeparator = options.nsSeparator !== undefined ? options.nsSeparator : this.options.nsSeparator;
        if (nsSeparator === undefined) nsSeparator = ':';
        const keySeparator = options.keySeparator !== undefined ? options.keySeparator : this.options.keySeparator;
        let namespaces = options.ns || this.options.defaultNS || [];
        const wouldCheckForNsInKey = nsSeparator && key.indexOf(nsSeparator) > -1;
        const seemsNaturalLanguage = !this.options.userDefinedKeySeparator && !options.keySeparator && !this.options.userDefinedNsSeparator && !options.nsSeparator && !$1bac384020b50752$var$looksLikeObjectPath(key, nsSeparator, keySeparator);
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
    translate(keys, options, lastKey) {
        if (typeof options !== 'object' && this.options.overloadTranslationOptionHandler) options = this.options.overloadTranslationOptionHandler(arguments);
        if (typeof options === 'object') options = {
            ...options
        };
        if (!options) options = {};
        if (keys === undefined || keys === null) return '';
        if (!Array.isArray(keys)) keys = [
            String(keys)
        ];
        const returnDetails = options.returnDetails !== undefined ? options.returnDetails : this.options.returnDetails;
        const keySeparator = options.keySeparator !== undefined ? options.keySeparator : this.options.keySeparator;
        const { key: key, namespaces: namespaces } = this.extractFromKey(keys[keys.length - 1], options);
        const namespace = namespaces[namespaces.length - 1];
        const lng = options.lng || this.language;
        const appendNamespaceToCIMode = options.appendNamespaceToCIMode || this.options.appendNamespaceToCIMode;
        if (lng && lng.toLowerCase() === 'cimode') {
            if (appendNamespaceToCIMode) {
                const nsSeparator = options.nsSeparator || this.options.nsSeparator;
                if (returnDetails) return {
                    res: `${namespace}${nsSeparator}${key}`,
                    usedKey: key,
                    exactUsedKey: key,
                    usedLng: lng,
                    usedNS: namespace,
                    usedParams: this.getUsedParamsDetails(options)
                };
                return `${namespace}${nsSeparator}${key}`;
            }
            if (returnDetails) return {
                res: key,
                usedKey: key,
                exactUsedKey: key,
                usedLng: lng,
                usedNS: namespace,
                usedParams: this.getUsedParamsDetails(options)
            };
            return key;
        }
        const resolved = this.resolve(keys, options);
        let res = resolved && resolved.res;
        const resUsedKey = resolved && resolved.usedKey || key;
        const resExactUsedKey = resolved && resolved.exactUsedKey || key;
        const resType = Object.prototype.toString.apply(res);
        const noObject = [
            '[object Number]',
            '[object Function]',
            '[object RegExp]'
        ];
        const joinArrays = options.joinArrays !== undefined ? options.joinArrays : this.options.joinArrays;
        const handleAsObjectInI18nFormat = !this.i18nFormat || this.i18nFormat.handleAsObject;
        const handleAsObject = !$1bac384020b50752$var$isString(res) && typeof res !== 'boolean' && typeof res !== 'number';
        if (handleAsObjectInI18nFormat && res && handleAsObject && noObject.indexOf(resType) < 0 && !($1bac384020b50752$var$isString(joinArrays) && Array.isArray(res))) {
            if (!options.returnObjects && !this.options.returnObjects) {
                if (!this.options.returnedObjectHandler) this.logger.warn('accessing an object - but returnObjects options is not enabled!');
                const r = this.options.returnedObjectHandler ? this.options.returnedObjectHandler(resUsedKey, res, {
                    ...options,
                    ns: namespaces
                }) : `key '${key} (${this.language})' returned an object instead of string.`;
                if (returnDetails) {
                    resolved.res = r;
                    resolved.usedParams = this.getUsedParamsDetails(options);
                    return resolved;
                }
                return r;
            }
            if (keySeparator) {
                const resTypeIsArray = Array.isArray(res);
                const copy = resTypeIsArray ? [] : {};
                const newKeyToUse = resTypeIsArray ? resExactUsedKey : resUsedKey;
                for(const m in res)if (Object.prototype.hasOwnProperty.call(res, m)) {
                    const deepKey = `${newKeyToUse}${keySeparator}${m}`;
                    copy[m] = this.translate(deepKey, {
                        ...options,
                        joinArrays: false,
                        ns: namespaces
                    });
                    if (copy[m] === deepKey) copy[m] = res[m];
                }
                res = copy;
            }
        } else if (handleAsObjectInI18nFormat && $1bac384020b50752$var$isString(joinArrays) && Array.isArray(res)) {
            res = res.join(joinArrays);
            if (res) res = this.extendTranslation(res, keys, options, lastKey);
        } else {
            let usedDefault = false;
            let usedKey = false;
            const needsPluralHandling = options.count !== undefined && !$1bac384020b50752$var$isString(options.count);
            const hasDefaultValue = $1bac384020b50752$var$Translator.hasDefaultValue(options);
            const defaultValueSuffix = needsPluralHandling ? this.pluralResolver.getSuffix(lng, options.count, options) : '';
            const defaultValueSuffixOrdinalFallback = options.ordinal && needsPluralHandling ? this.pluralResolver.getSuffix(lng, options.count, {
                ordinal: false
            }) : '';
            const needsZeroSuffixLookup = needsPluralHandling && !options.ordinal && options.count === 0 && this.pluralResolver.shouldUseIntlApi();
            const defaultValue = needsZeroSuffixLookup && options[`defaultValue${this.options.pluralSeparator}zero`] || options[`defaultValue${defaultValueSuffix}`] || options[`defaultValue${defaultValueSuffixOrdinalFallback}`] || options.defaultValue;
            if (!this.isValidLookup(res) && hasDefaultValue) {
                usedDefault = true;
                res = defaultValue;
            }
            if (!this.isValidLookup(res)) {
                usedKey = true;
                res = key;
            }
            const missingKeyNoValueFallbackToKey = options.missingKeyNoValueFallbackToKey || this.options.missingKeyNoValueFallbackToKey;
            const resForMissing = missingKeyNoValueFallbackToKey && usedKey ? undefined : res;
            const updateMissing = hasDefaultValue && defaultValue !== res && this.options.updateMissing;
            if (usedKey || usedDefault || updateMissing) {
                this.logger.log(updateMissing ? 'updateKey' : 'missingKey', lng, namespace, key, updateMissing ? defaultValue : res);
                if (keySeparator) {
                    const fk = this.resolve(key, {
                        ...options,
                        keySeparator: false
                    });
                    if (fk && fk.res) this.logger.warn('Seems the loaded translations were in flat JSON format instead of nested. Either set keySeparator: false on init or make sure your translations are published in nested format.');
                }
                let lngs = [];
                const fallbackLngs = this.languageUtils.getFallbackCodes(this.options.fallbackLng, options.lng || this.language);
                if (this.options.saveMissingTo === 'fallback' && fallbackLngs && fallbackLngs[0]) for(let i = 0; i < fallbackLngs.length; i++)lngs.push(fallbackLngs[i]);
                else if (this.options.saveMissingTo === 'all') lngs = this.languageUtils.toResolveHierarchy(options.lng || this.language);
                else lngs.push(options.lng || this.language);
                const send = (l, k, specificDefaultValue)=>{
                    const defaultForMissing = hasDefaultValue && specificDefaultValue !== res ? specificDefaultValue : resForMissing;
                    if (this.options.missingKeyHandler) this.options.missingKeyHandler(l, namespace, k, defaultForMissing, updateMissing, options);
                    else if (this.backendConnector && this.backendConnector.saveMissing) this.backendConnector.saveMissing(l, namespace, k, defaultForMissing, updateMissing, options);
                    this.emit('missingKey', l, namespace, k, res);
                };
                if (this.options.saveMissing) {
                    if (this.options.saveMissingPlurals && needsPluralHandling) lngs.forEach((language)=>{
                        const suffixes = this.pluralResolver.getSuffixes(language, options);
                        if (needsZeroSuffixLookup && options[`defaultValue${this.options.pluralSeparator}zero`] && suffixes.indexOf(`${this.options.pluralSeparator}zero`) < 0) suffixes.push(`${this.options.pluralSeparator}zero`);
                        suffixes.forEach((suffix)=>{
                            send([
                                language
                            ], key + suffix, options[`defaultValue${suffix}`] || defaultValue);
                        });
                    });
                    else send(lngs, key, defaultValue);
                }
            }
            res = this.extendTranslation(res, keys, options, resolved, lastKey);
            if (usedKey && res === key && this.options.appendNamespaceToMissingKey) res = `${namespace}:${key}`;
            if ((usedKey || usedDefault) && this.options.parseMissingKeyHandler) {
                if (this.options.compatibilityAPI !== 'v1') res = this.options.parseMissingKeyHandler(this.options.appendNamespaceToMissingKey ? `${namespace}:${key}` : key, usedDefault ? res : undefined);
                else res = this.options.parseMissingKeyHandler(res);
            }
        }
        if (returnDetails) {
            resolved.res = res;
            resolved.usedParams = this.getUsedParamsDetails(options);
            return resolved;
        }
        return res;
    }
    extendTranslation(res, key, options, resolved, lastKey) {
        var _this = this;
        if (this.i18nFormat && this.i18nFormat.parse) res = this.i18nFormat.parse(res, {
            ...this.options.interpolation.defaultVariables,
            ...options
        }, options.lng || this.language || resolved.usedLng, resolved.usedNS, resolved.usedKey, {
            resolved: resolved
        });
        else if (!options.skipInterpolation) {
            if (options.interpolation) this.interpolator.init({
                ...options,
                interpolation: {
                    ...this.options.interpolation,
                    ...options.interpolation
                }
            });
            const skipOnVariables = $1bac384020b50752$var$isString(res) && (options && options.interpolation && options.interpolation.skipOnVariables !== undefined ? options.interpolation.skipOnVariables : this.options.interpolation.skipOnVariables);
            let nestBef;
            if (skipOnVariables) {
                const nb = res.match(this.interpolator.nestingRegexp);
                nestBef = nb && nb.length;
            }
            let data = options.replace && !$1bac384020b50752$var$isString(options.replace) ? options.replace : options;
            if (this.options.interpolation.defaultVariables) data = {
                ...this.options.interpolation.defaultVariables,
                ...data
            };
            res = this.interpolator.interpolate(res, data, options.lng || this.language || resolved.usedLng, options);
            if (skipOnVariables) {
                const na = res.match(this.interpolator.nestingRegexp);
                const nestAft = na && na.length;
                if (nestBef < nestAft) options.nest = false;
            }
            if (!options.lng && this.options.compatibilityAPI !== 'v1' && resolved && resolved.res) options.lng = this.language || resolved.usedLng;
            if (options.nest !== false) res = this.interpolator.nest(res, function() {
                for(var _len = arguments.length, args = new Array(_len), _key = 0; _key < _len; _key++)args[_key] = arguments[_key];
                if (lastKey && lastKey[0] === args[0] && !options.context) {
                    _this.logger.warn(`It seems you are nesting recursively key: ${args[0]} in key: ${key[0]}`);
                    return null;
                }
                return _this.translate(...args, key);
            }, options);
            if (options.interpolation) this.interpolator.reset();
        }
        const postProcess = options.postProcess || this.options.postProcess;
        const postProcessorNames = $1bac384020b50752$var$isString(postProcess) ? [
            postProcess
        ] : postProcess;
        if (res !== undefined && res !== null && postProcessorNames && postProcessorNames.length && options.applyPostProcessor !== false) res = $1bac384020b50752$var$postProcessor.handle(postProcessorNames, res, key, this.options && this.options.postProcessPassResolved ? {
            i18nResolved: {
                ...resolved,
                usedParams: this.getUsedParamsDetails(options)
            },
            ...options
        } : options, this);
        return res;
    }
    resolve(keys) {
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
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
            const extracted = this.extractFromKey(k, options);
            const key = extracted.key;
            usedKey = key;
            let namespaces = extracted.namespaces;
            if (this.options.fallbackNS) namespaces = namespaces.concat(this.options.fallbackNS);
            const needsPluralHandling = options.count !== undefined && !$1bac384020b50752$var$isString(options.count);
            const needsZeroSuffixLookup = needsPluralHandling && !options.ordinal && options.count === 0 && this.pluralResolver.shouldUseIntlApi();
            const needsContextHandling = options.context !== undefined && ($1bac384020b50752$var$isString(options.context) || typeof options.context === 'number') && options.context !== '';
            const codes = options.lngs ? options.lngs : this.languageUtils.toResolveHierarchy(options.lng || this.language, options.fallbackLng);
            namespaces.forEach((ns)=>{
                if (this.isValidLookup(found)) return;
                usedNS = ns;
                if (!$1bac384020b50752$var$checkedLoadedFor[`${codes[0]}-${ns}`] && this.utils && this.utils.hasLoadedNamespace && !this.utils.hasLoadedNamespace(usedNS)) {
                    $1bac384020b50752$var$checkedLoadedFor[`${codes[0]}-${ns}`] = true;
                    this.logger.warn(`key "${usedKey}" for languages "${codes.join(', ')}" won't get resolved as namespace "${usedNS}" was not yet loaded`, 'This means something IS WRONG in your setup. You access the t function before i18next.init / i18next.loadNamespace / i18next.changeLanguage was done. Wait for the callback or Promise to resolve before accessing it!!!');
                }
                codes.forEach((code)=>{
                    if (this.isValidLookup(found)) return;
                    usedLng = code;
                    const finalKeys = [
                        key
                    ];
                    if (this.i18nFormat && this.i18nFormat.addLookupKeys) this.i18nFormat.addLookupKeys(finalKeys, key, code, ns, options);
                    else {
                        let pluralSuffix;
                        if (needsPluralHandling) pluralSuffix = this.pluralResolver.getSuffix(code, options.count, options);
                        const zeroSuffix = `${this.options.pluralSeparator}zero`;
                        const ordinalPrefix = `${this.options.pluralSeparator}ordinal${this.options.pluralSeparator}`;
                        if (needsPluralHandling) {
                            finalKeys.push(key + pluralSuffix);
                            if (options.ordinal && pluralSuffix.indexOf(ordinalPrefix) === 0) finalKeys.push(key + pluralSuffix.replace(ordinalPrefix, this.options.pluralSeparator));
                            if (needsZeroSuffixLookup) finalKeys.push(key + zeroSuffix);
                        }
                        if (needsContextHandling) {
                            const contextKey = `${key}${this.options.contextSeparator}${options.context}`;
                            finalKeys.push(contextKey);
                            if (needsPluralHandling) {
                                finalKeys.push(contextKey + pluralSuffix);
                                if (options.ordinal && pluralSuffix.indexOf(ordinalPrefix) === 0) finalKeys.push(contextKey + pluralSuffix.replace(ordinalPrefix, this.options.pluralSeparator));
                                if (needsZeroSuffixLookup) finalKeys.push(contextKey + zeroSuffix);
                            }
                        }
                    }
                    let possibleKey;
                    while(possibleKey = finalKeys.pop())if (!this.isValidLookup(found)) {
                        exactUsedKey = possibleKey;
                        found = this.getResource(code, ns, possibleKey, options);
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
    getResource(code, ns, key) {
        let options = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : {};
        if (this.i18nFormat && this.i18nFormat.getResource) return this.i18nFormat.getResource(code, ns, key, options);
        return this.resourceStore.getResource(code, ns, key, options);
    }
    getUsedParamsDetails() {
        let options = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};
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
const $1bac384020b50752$var$capitalize = (string)=>string.charAt(0).toUpperCase() + string.slice(1);
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
            if (typeof Intl !== 'undefined' && typeof Intl.getCanonicalLocales !== 'undefined') try {
                let formattedCode = Intl.getCanonicalLocales(code)[0];
                if (formattedCode && this.options.lowerCaseLng) formattedCode = formattedCode.toLowerCase();
                if (formattedCode) return formattedCode;
            } catch (e) {}
            const specialCases = [
                'hans',
                'hant',
                'latn',
                'cyrl',
                'cans',
                'mong',
                'arab'
            ];
            let p = code.split('-');
            if (this.options.lowerCaseLng) p = p.map((part)=>part.toLowerCase());
            else if (p.length === 2) {
                p[0] = p[0].toLowerCase();
                p[1] = p[1].toUpperCase();
                if (specialCases.indexOf(p[1].toLowerCase()) > -1) p[1] = $1bac384020b50752$var$capitalize(p[1].toLowerCase());
            } else if (p.length === 3) {
                p[0] = p[0].toLowerCase();
                if (p[1].length === 2) p[1] = p[1].toUpperCase();
                if (p[0] !== 'sgn' && p[2].length === 2) p[2] = p[2].toUpperCase();
                if (specialCases.indexOf(p[1].toLowerCase()) > -1) p[1] = $1bac384020b50752$var$capitalize(p[1].toLowerCase());
                if (specialCases.indexOf(p[2].toLowerCase()) > -1) p[2] = $1bac384020b50752$var$capitalize(p[2].toLowerCase());
            }
            return p.join('-');
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
        const fallbackCodes = this.getFallbackCodes(fallbackCode || this.options.fallbackLng || [], code);
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
let $1bac384020b50752$var$sets = [
    {
        lngs: [
            'ach',
            'ak',
            'am',
            'arn',
            'br',
            'fil',
            'gun',
            'ln',
            'mfe',
            'mg',
            'mi',
            'oc',
            'pt',
            'pt-BR',
            'tg',
            'tl',
            'ti',
            'tr',
            'uz',
            'wa'
        ],
        nr: [
            1,
            2
        ],
        fc: 1
    },
    {
        lngs: [
            'af',
            'an',
            'ast',
            'az',
            'bg',
            'bn',
            'ca',
            'da',
            'de',
            'dev',
            'el',
            'en',
            'eo',
            'es',
            'et',
            'eu',
            'fi',
            'fo',
            'fur',
            'fy',
            'gl',
            'gu',
            'ha',
            'hi',
            'hu',
            'hy',
            'ia',
            'it',
            'kk',
            'kn',
            'ku',
            'lb',
            'mai',
            'ml',
            'mn',
            'mr',
            'nah',
            'nap',
            'nb',
            'ne',
            'nl',
            'nn',
            'no',
            'nso',
            'pa',
            'pap',
            'pms',
            'ps',
            'pt-PT',
            'rm',
            'sco',
            'se',
            'si',
            'so',
            'son',
            'sq',
            'sv',
            'sw',
            'ta',
            'te',
            'tk',
            'ur',
            'yo'
        ],
        nr: [
            1,
            2
        ],
        fc: 2
    },
    {
        lngs: [
            'ay',
            'bo',
            'cgg',
            'fa',
            'ht',
            'id',
            'ja',
            'jbo',
            'ka',
            'km',
            'ko',
            'ky',
            'lo',
            'ms',
            'sah',
            'su',
            'th',
            'tt',
            'ug',
            'vi',
            'wo',
            'zh'
        ],
        nr: [
            1
        ],
        fc: 3
    },
    {
        lngs: [
            'be',
            'bs',
            'cnr',
            'dz',
            'hr',
            'ru',
            'sr',
            'uk'
        ],
        nr: [
            1,
            2,
            5
        ],
        fc: 4
    },
    {
        lngs: [
            'ar'
        ],
        nr: [
            0,
            1,
            2,
            3,
            11,
            100
        ],
        fc: 5
    },
    {
        lngs: [
            'cs',
            'sk'
        ],
        nr: [
            1,
            2,
            5
        ],
        fc: 6
    },
    {
        lngs: [
            'csb',
            'pl'
        ],
        nr: [
            1,
            2,
            5
        ],
        fc: 7
    },
    {
        lngs: [
            'cy'
        ],
        nr: [
            1,
            2,
            3,
            8
        ],
        fc: 8
    },
    {
        lngs: [
            'fr'
        ],
        nr: [
            1,
            2
        ],
        fc: 9
    },
    {
        lngs: [
            'ga'
        ],
        nr: [
            1,
            2,
            3,
            7,
            11
        ],
        fc: 10
    },
    {
        lngs: [
            'gd'
        ],
        nr: [
            1,
            2,
            3,
            20
        ],
        fc: 11
    },
    {
        lngs: [
            'is'
        ],
        nr: [
            1,
            2
        ],
        fc: 12
    },
    {
        lngs: [
            'jv'
        ],
        nr: [
            0,
            1
        ],
        fc: 13
    },
    {
        lngs: [
            'kw'
        ],
        nr: [
            1,
            2,
            3,
            4
        ],
        fc: 14
    },
    {
        lngs: [
            'lt'
        ],
        nr: [
            1,
            2,
            10
        ],
        fc: 15
    },
    {
        lngs: [
            'lv'
        ],
        nr: [
            1,
            2,
            0
        ],
        fc: 16
    },
    {
        lngs: [
            'mk'
        ],
        nr: [
            1,
            2
        ],
        fc: 17
    },
    {
        lngs: [
            'mnk'
        ],
        nr: [
            0,
            1,
            2
        ],
        fc: 18
    },
    {
        lngs: [
            'mt'
        ],
        nr: [
            1,
            2,
            11,
            20
        ],
        fc: 19
    },
    {
        lngs: [
            'or'
        ],
        nr: [
            2,
            1
        ],
        fc: 2
    },
    {
        lngs: [
            'ro'
        ],
        nr: [
            1,
            2,
            20
        ],
        fc: 20
    },
    {
        lngs: [
            'sl'
        ],
        nr: [
            5,
            1,
            2,
            3
        ],
        fc: 21
    },
    {
        lngs: [
            'he',
            'iw'
        ],
        nr: [
            1,
            2,
            20,
            21
        ],
        fc: 22
    }
];
let $1bac384020b50752$var$_rulesPluralsTypes = {
    1: (n)=>Number(n > 1),
    2: (n)=>Number(n != 1),
    3: (n)=>0,
    4: (n)=>Number(n % 10 == 1 && n % 100 != 11 ? 0 : n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20) ? 1 : 2),
    5: (n)=>Number(n == 0 ? 0 : n == 1 ? 1 : n == 2 ? 2 : n % 100 >= 3 && n % 100 <= 10 ? 3 : n % 100 >= 11 ? 4 : 5),
    6: (n)=>Number(n == 1 ? 0 : n >= 2 && n <= 4 ? 1 : 2),
    7: (n)=>Number(n == 1 ? 0 : n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20) ? 1 : 2),
    8: (n)=>Number(n == 1 ? 0 : n == 2 ? 1 : n != 8 && n != 11 ? 2 : 3),
    9: (n)=>Number(n >= 2),
    10: (n)=>Number(n == 1 ? 0 : n == 2 ? 1 : n < 7 ? 2 : n < 11 ? 3 : 4),
    11: (n)=>Number(n == 1 || n == 11 ? 0 : n == 2 || n == 12 ? 1 : n > 2 && n < 20 ? 2 : 3),
    12: (n)=>Number(n % 10 != 1 || n % 100 == 11),
    13: (n)=>Number(n !== 0),
    14: (n)=>Number(n == 1 ? 0 : n == 2 ? 1 : n == 3 ? 2 : 3),
    15: (n)=>Number(n % 10 == 1 && n % 100 != 11 ? 0 : n % 10 >= 2 && (n % 100 < 10 || n % 100 >= 20) ? 1 : 2),
    16: (n)=>Number(n % 10 == 1 && n % 100 != 11 ? 0 : n !== 0 ? 1 : 2),
    17: (n)=>Number(n == 1 || n % 10 == 1 && n % 100 != 11 ? 0 : 1),
    18: (n)=>Number(n == 0 ? 0 : n == 1 ? 1 : 2),
    19: (n)=>Number(n == 1 ? 0 : n == 0 || n % 100 > 1 && n % 100 < 11 ? 1 : n % 100 > 10 && n % 100 < 20 ? 2 : 3),
    20: (n)=>Number(n == 1 ? 0 : n == 0 || n % 100 > 0 && n % 100 < 20 ? 1 : 2),
    21: (n)=>Number(n % 100 == 1 ? 1 : n % 100 == 2 ? 2 : n % 100 == 3 || n % 100 == 4 ? 3 : 0),
    22: (n)=>Number(n == 1 ? 0 : n == 2 ? 1 : (n < 0 || n > 10) && n % 10 == 0 ? 2 : 3)
};
const $1bac384020b50752$var$nonIntlVersions = [
    'v1',
    'v2',
    'v3'
];
const $1bac384020b50752$var$intlVersions = [
    'v4'
];
const $1bac384020b50752$var$suffixesOrder = {
    zero: 0,
    one: 1,
    two: 2,
    few: 3,
    many: 4,
    other: 5
};
const $1bac384020b50752$var$createRules = ()=>{
    const rules = {};
    $1bac384020b50752$var$sets.forEach((set)=>{
        set.lngs.forEach((l)=>{
            rules[l] = {
                numbers: set.nr,
                plurals: $1bac384020b50752$var$_rulesPluralsTypes[set.fc]
            };
        });
    });
    return rules;
};
class $1bac384020b50752$var$PluralResolver {
    constructor(languageUtils){
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
        this.languageUtils = languageUtils;
        this.options = options;
        this.logger = $1bac384020b50752$var$baseLogger.create('pluralResolver');
        if ((!this.options.compatibilityJSON || $1bac384020b50752$var$intlVersions.includes(this.options.compatibilityJSON)) && (typeof Intl === 'undefined' || !Intl.PluralRules)) {
            this.options.compatibilityJSON = 'v3';
            this.logger.error('Your environment seems not to be Intl API compatible, use an Intl.PluralRules polyfill. Will fallback to the compatibilityJSON v3 format handling.');
        }
        this.rules = $1bac384020b50752$var$createRules();
        this.pluralRulesCache = {};
    }
    addRule(lng, obj) {
        this.rules[lng] = obj;
    }
    clearCache() {
        this.pluralRulesCache = {};
    }
    getRule(code) {
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
        if (this.shouldUseIntlApi()) {
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
                if (!code.match(/-|_/)) return;
                const lngPart = this.languageUtils.getLanguagePartFromCode(code);
                rule = this.getRule(lngPart, options);
            }
            this.pluralRulesCache[cacheKey] = rule;
            return rule;
        }
        return this.rules[code] || this.rules[this.languageUtils.getLanguagePartFromCode(code)];
    }
    needsPlural(code) {
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
        const rule = this.getRule(code, options);
        if (this.shouldUseIntlApi()) return rule && rule.resolvedOptions().pluralCategories.length > 1;
        return rule && rule.numbers.length > 1;
    }
    getPluralFormsOfKey(code, key) {
        let options = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : {};
        return this.getSuffixes(code, options).map((suffix)=>`${key}${suffix}`);
    }
    getSuffixes(code) {
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
        const rule = this.getRule(code, options);
        if (!rule) return [];
        if (this.shouldUseIntlApi()) return rule.resolvedOptions().pluralCategories.sort((pluralCategory1, pluralCategory2)=>$1bac384020b50752$var$suffixesOrder[pluralCategory1] - $1bac384020b50752$var$suffixesOrder[pluralCategory2]).map((pluralCategory)=>`${this.options.prepend}${options.ordinal ? `ordinal${this.options.prepend}` : ''}${pluralCategory}`);
        return rule.numbers.map((number)=>this.getSuffix(code, number, options));
    }
    getSuffix(code, count) {
        let options = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : {};
        const rule = this.getRule(code, options);
        if (rule) {
            if (this.shouldUseIntlApi()) return `${this.options.prepend}${options.ordinal ? `ordinal${this.options.prepend}` : ''}${rule.select(count)}`;
            return this.getSuffixRetroCompatible(rule, count);
        }
        this.logger.warn(`no plural rule found for: ${code}`);
        return '';
    }
    getSuffixRetroCompatible(rule, count) {
        const idx = rule.noAbs ? rule.plurals(count) : rule.plurals(Math.abs(count));
        let suffix = rule.numbers[idx];
        if (this.options.simplifyPluralSuffix && rule.numbers.length === 2 && rule.numbers[0] === 1) {
            if (suffix === 2) suffix = 'plural';
            else if (suffix === 1) suffix = '';
        }
        const returnSuffix = ()=>this.options.prepend && suffix.toString() ? this.options.prepend + suffix.toString() : suffix.toString();
        if (this.options.compatibilityJSON === 'v1') {
            if (suffix === 1) return '';
            if (typeof suffix === 'number') return `_plural_${suffix.toString()}`;
            return returnSuffix();
        } else if (this.options.compatibilityJSON === 'v2') return returnSuffix();
        else if (this.options.simplifyPluralSuffix && rule.numbers.length === 2 && rule.numbers[0] === 1) return returnSuffix();
        return this.options.prepend && idx.toString() ? this.options.prepend + idx.toString() : idx.toString();
    }
    shouldUseIntlApi() {
        return !$1bac384020b50752$var$nonIntlVersions.includes(this.options.compatibilityJSON);
    }
}
const $1bac384020b50752$var$deepFindWithDefaults = function(data, defaultData, key) {
    let keySeparator = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : '.';
    let ignoreJSONStructure = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : true;
    let path = $1bac384020b50752$var$getPathWithDefaults(data, defaultData, key);
    if (!path && ignoreJSONStructure && $1bac384020b50752$var$isString(key)) {
        path = $1bac384020b50752$var$deepFind(data, key, keySeparator);
        if (path === undefined) path = $1bac384020b50752$var$deepFind(defaultData, key, keySeparator);
    }
    return path;
};
const $1bac384020b50752$var$regexSafe = (val)=>val.replace(/\$/g, '$$$$');
class $1bac384020b50752$var$Interpolator {
    constructor(){
        let options = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};
        this.logger = $1bac384020b50752$var$baseLogger.create('interpolator');
        this.options = options;
        this.format = options.interpolation && options.interpolation.format || ((value)=>value);
        this.init(options);
    }
    init() {
        let options = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};
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
            if (existingRegExp && existingRegExp.source === pattern) {
                existingRegExp.lastIndex = 0;
                return existingRegExp;
            }
            return new RegExp(pattern, 'g');
        };
        this.regexp = getOrResetRegExp(this.regexp, `${this.prefix}(.+?)${this.suffix}`);
        this.regexpUnescape = getOrResetRegExp(this.regexpUnescape, `${this.prefix}${this.unescapePrefix}(.+?)${this.unescapeSuffix}${this.suffix}`);
        this.nestingRegexp = getOrResetRegExp(this.nestingRegexp, `${this.nestingPrefix}(.+?)${this.nestingSuffix}`);
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
        const missingInterpolationHandler = options && options.missingInterpolationHandler || this.options.missingInterpolationHandler;
        const skipOnVariables = options && options.interpolation && options.interpolation.skipOnVariables !== undefined ? options.interpolation.skipOnVariables : this.options.interpolation.skipOnVariables;
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
    nest(str, fc) {
        let options = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : {};
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
            if (matchedSingleQuotes && matchedSingleQuotes.length % 2 === 0 && !matchedDoubleQuotes || matchedDoubleQuotes.length % 2 !== 0) optionsString = optionsString.replace(/'/g, '"');
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
            let doReduce = false;
            if (match[0].indexOf(this.formatSeparator) !== -1 && !/{.*}/.test(match[1])) {
                const r = match[1].split(this.formatSeparator).map((elem)=>elem.trim());
                match[1] = r.shift();
                formatters = r;
                doReduce = true;
            }
            value = fc(handleHasOptions.call(this, match[1].trim(), clonedOptions), clonedOptions);
            if (value && match[0] === str && !$1bac384020b50752$var$isString(value)) return value;
            if (!$1bac384020b50752$var$isString(value)) value = $1bac384020b50752$var$makeString(value);
            if (!value) {
                this.logger.warn(`missed to resolve ${match[1]} for nesting ${str}`);
                value = '';
            }
            if (doReduce) value = formatters.reduce((v, f)=>this.format(v, f, options.lng, {
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
    return (val, lng, options)=>{
        let optForCache = options;
        if (options && options.interpolationkey && options.formatParams && options.formatParams[options.interpolationkey] && options[options.interpolationkey]) optForCache = {
            ...optForCache,
            [options.interpolationkey]: undefined
        };
        const key = lng + JSON.stringify(optForCache);
        let formatter = cache[key];
        if (!formatter) {
            formatter = fn($1bac384020b50752$var$getCleanedCode(lng), options);
            cache[key] = formatter;
        }
        return formatter(val);
    };
};
class $1bac384020b50752$var$Formatter {
    constructor(){
        let options = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};
        this.logger = $1bac384020b50752$var$baseLogger.create('formatter');
        this.options = options;
        this.formats = {
            number: $1bac384020b50752$var$createCachedFormatter((lng, opt)=>{
                const formatter = new Intl.NumberFormat(lng, {
                    ...opt
                });
                return (val)=>formatter.format(val);
            }),
            currency: $1bac384020b50752$var$createCachedFormatter((lng, opt)=>{
                const formatter = new Intl.NumberFormat(lng, {
                    ...opt,
                    style: 'currency'
                });
                return (val)=>formatter.format(val);
            }),
            datetime: $1bac384020b50752$var$createCachedFormatter((lng, opt)=>{
                const formatter = new Intl.DateTimeFormat(lng, {
                    ...opt
                });
                return (val)=>formatter.format(val);
            }),
            relativetime: $1bac384020b50752$var$createCachedFormatter((lng, opt)=>{
                const formatter = new Intl.RelativeTimeFormat(lng, {
                    ...opt
                });
                return (val)=>formatter.format(val, opt.range || 'day');
            }),
            list: $1bac384020b50752$var$createCachedFormatter((lng, opt)=>{
                const formatter = new Intl.ListFormat(lng, {
                    ...opt
                });
                return (val)=>formatter.format(val);
            })
        };
        this.init(options);
    }
    init(services) {
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {
            interpolation: {}
        };
        this.formatSeparator = options.interpolation.formatSeparator || ',';
    }
    add(name, fc) {
        this.formats[name.toLowerCase().trim()] = fc;
    }
    addCached(name, fc) {
        this.formats[name.toLowerCase().trim()] = $1bac384020b50752$var$createCachedFormatter(fc);
    }
    format(value, format, lng) {
        let options = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : {};
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
                    const valOptions = options && options.formatParams && options.formatParams[options.interpolationkey] || {};
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
    constructor(backend, store, services){
        let options = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : {};
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
        if (this.backend && this.backend.init) this.backend.init(services, options.backend, options);
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
    read(lng, ns, fcName) {
        let tried = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : 0;
        let wait = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : this.retryTimeout;
        let callback = arguments.length > 5 ? arguments[5] : undefined;
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
    prepareLoading(languages, namespaces) {
        let options = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : {};
        let callback = arguments.length > 3 ? arguments[3] : undefined;
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
    loadOne(name) {
        let prefix = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : '';
        const s = name.split('|');
        const lng = s[0];
        const ns = s[1];
        this.read(lng, ns, 'read', undefined, undefined, (err, data)=>{
            if (err) this.logger.warn(`${prefix}loading namespace ${ns} for language ${lng} failed`, err);
            if (!err && data) this.logger.log(`${prefix}loaded namespace ${ns} for language ${lng}`, data);
            this.loaded(name, err, data);
        });
    }
    saveMissing(languages, namespace, key, fallbackValue, isUpdate) {
        let options = arguments.length > 5 && arguments[5] !== undefined ? arguments[5] : {};
        let clb = arguments.length > 6 && arguments[6] !== undefined ? arguments[6] : ()=>{};
        if (this.services.utils && this.services.utils.hasLoadedNamespace && !this.services.utils.hasLoadedNamespace(namespace)) {
            this.logger.warn(`did not save key "${key}" as the namespace "${namespace}" was not yet loaded`, 'This means something IS WRONG in your setup. You access the t function before i18next.init / i18next.loadNamespace / i18next.changeLanguage was done. Wait for the callback or Promise to resolve before accessing it!!!');
            return;
        }
        if (key === undefined || key === null || key === '') return;
        if (this.backend && this.backend.create) {
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
        initImmediate: true,
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
        }
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
    if (options.supportedLngs && options.supportedLngs.indexOf('cimode') < 0) options.supportedLngs = options.supportedLngs.concat([
        'cimode'
    ]);
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
    constructor(){
        let options = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};
        let callback = arguments.length > 1 ? arguments[1] : undefined;
        super();
        this.options = $1bac384020b50752$var$transformOptions(options);
        this.services = {};
        this.logger = $1bac384020b50752$var$baseLogger;
        this.modules = {
            external: []
        };
        $1bac384020b50752$var$bindMemberFunctions(this);
        if (callback && !this.isInitialized && !options.isClone) {
            if (!this.options.initImmediate) {
                this.init(options, callback);
                return this;
            }
            setTimeout(()=>{
                this.init(options, callback);
            }, 0);
        }
    }
    init() {
        var _this = this;
        let options = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};
        let callback = arguments.length > 1 ? arguments[1] : undefined;
        this.isInitializing = true;
        if (typeof options === 'function') {
            callback = options;
            options = {};
        }
        if (!options.defaultNS && options.defaultNS !== false && options.ns) {
            if ($1bac384020b50752$var$isString(options.ns)) options.defaultNS = options.ns;
            else if (options.ns.indexOf('translation') < 0) options.defaultNS = options.ns[0];
        }
        const defOpts = $1bac384020b50752$var$get();
        this.options = {
            ...defOpts,
            ...this.options,
            ...$1bac384020b50752$var$transformOptions(options)
        };
        if (this.options.compatibilityAPI !== 'v1') this.options.interpolation = {
            ...defOpts.interpolation,
            ...this.options.interpolation
        };
        if (options.keySeparator !== undefined) this.options.userDefinedKeySeparator = options.keySeparator;
        if (options.nsSeparator !== undefined) this.options.userDefinedNsSeparator = options.nsSeparator;
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
            else if (typeof Intl !== 'undefined') formatter = $1bac384020b50752$var$Formatter;
            const lu = new $1bac384020b50752$var$LanguageUtil(this.options);
            this.store = new $1bac384020b50752$var$ResourceStore(this.options.resources, this.options);
            const s = this.services;
            s.logger = $1bac384020b50752$var$baseLogger;
            s.resourceStore = this.store;
            s.languageUtils = lu;
            s.pluralResolver = new $1bac384020b50752$var$PluralResolver(lu, {
                prepend: this.options.pluralSeparator,
                compatibilityJSON: this.options.compatibilityJSON,
                simplifyPluralSuffix: this.options.simplifyPluralSuffix
            });
            if (formatter && (!this.options.interpolation.format || this.options.interpolation.format === defOpts.interpolation.format)) {
                s.formatter = createClassOnDemand(formatter);
                s.formatter.init(s, this.options);
                this.options.interpolation.format = s.formatter.format.bind(s.formatter);
            }
            s.interpolator = new $1bac384020b50752$var$Interpolator(this.options);
            s.utils = {
                hasLoadedNamespace: this.hasLoadedNamespace.bind(this)
            };
            s.backendConnector = new $1bac384020b50752$var$Connector(createClassOnDemand(this.modules.backend), s.resourceStore, s, this.options);
            s.backendConnector.on('*', function(event) {
                for(var _len = arguments.length, args = new Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++)args[_key - 1] = arguments[_key];
                _this.emit(event, ...args);
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
            this.translator.on('*', function(event) {
                for(var _len2 = arguments.length, args = new Array(_len2 > 1 ? _len2 - 1 : 0), _key2 = 1; _key2 < _len2; _key2++)args[_key2 - 1] = arguments[_key2];
                _this.emit(event, ...args);
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
            this[fcName] = function() {
                return _this.store[fcName](...arguments);
            };
        });
        const storeApiChained = [
            'addResource',
            'addResources',
            'addResourceBundle',
            'removeResourceBundle'
        ];
        storeApiChained.forEach((fcName)=>{
            this[fcName] = function() {
                _this.store[fcName](...arguments);
                return _this;
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
            if (this.languages && this.options.compatibilityAPI !== 'v1' && !this.isInitialized) return finish(null, this.t.bind(this));
            this.changeLanguage(this.options.lng, finish);
        };
        if (this.options.resources || !this.options.initImmediate) load();
        else setTimeout(load, 0);
        return deferred;
    }
    loadResources(language) {
        let callback = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : $1bac384020b50752$var$noop;
        let usedCallback = callback;
        const usedLng = $1bac384020b50752$var$isString(language) ? language : this.language;
        if (typeof language === 'function') usedCallback = language;
        if (!this.options.resources || this.options.partialBundledLanguages) {
            if (usedLng && usedLng.toLowerCase() === 'cimode' && (!this.options.preload || this.options.preload.length === 0)) return usedCallback();
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
            if (this.options.preload) this.options.preload.forEach((l)=>append(l));
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
    }
    changeLanguage(lng, callback) {
        var _this2 = this;
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
                setLngProps(l);
                this.translator.changeLanguage(l);
                this.isLanguageChangingTo = undefined;
                this.emit('languageChanged', l);
                this.logger.log('languageChanged', l);
            } else this.isLanguageChangingTo = undefined;
            deferred.resolve(function() {
                return _this2.t(...arguments);
            });
            if (callback) callback(err, function() {
                return _this2.t(...arguments);
            });
        };
        const setLng = (lngs)=>{
            if (!lng && !lngs && this.services.languageDetector) lngs = [];
            const l = $1bac384020b50752$var$isString(lngs) ? lngs : this.services.languageUtils.getBestMatchFromCodes(lngs);
            if (l) {
                if (!this.language) setLngProps(l);
                if (!this.translator.language) this.translator.changeLanguage(l);
                if (this.services.languageDetector && this.services.languageDetector.cacheUserLanguage) this.services.languageDetector.cacheUserLanguage(l);
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
        var _this3 = this;
        const fixedT = function(key, opts) {
            let options;
            if (typeof opts !== 'object') {
                for(var _len3 = arguments.length, rest = new Array(_len3 > 2 ? _len3 - 2 : 0), _key3 = 2; _key3 < _len3; _key3++)rest[_key3 - 2] = arguments[_key3];
                options = _this3.options.overloadTranslationOptionHandler([
                    key,
                    opts
                ].concat(rest));
            } else options = {
                ...opts
            };
            options.lng = options.lng || fixedT.lng;
            options.lngs = options.lngs || fixedT.lngs;
            options.ns = options.ns || fixedT.ns;
            if (options.keyPrefix !== '') options.keyPrefix = options.keyPrefix || keyPrefix || fixedT.keyPrefix;
            const keySeparator = _this3.options.keySeparator || '.';
            let resultKey;
            if (options.keyPrefix && Array.isArray(key)) resultKey = key.map((k)=>`${options.keyPrefix}${keySeparator}${k}`);
            else resultKey = options.keyPrefix ? `${options.keyPrefix}${keySeparator}${key}` : key;
            return _this3.t(resultKey, options);
        };
        if ($1bac384020b50752$var$isString(lng)) fixedT.lng = lng;
        else fixedT.lngs = lng;
        fixedT.ns = ns;
        fixedT.keyPrefix = keyPrefix;
        return fixedT;
    }
    t() {
        return this.translator && this.translator.translate(...arguments);
    }
    exists() {
        return this.translator && this.translator.exists(...arguments);
    }
    setDefaultNamespace(ns) {
        this.options.defaultNS = ns;
    }
    hasLoadedNamespace(ns) {
        let options = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : {};
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
        if (!lng) lng = this.resolvedLanguage || (this.languages && this.languages.length > 0 ? this.languages[0] : this.language);
        if (!lng) return 'rtl';
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
        const languageUtils = this.services && this.services.languageUtils || new $1bac384020b50752$var$LanguageUtil($1bac384020b50752$var$get());
        return rtlLngs.indexOf(languageUtils.getLanguagePartFromCode(lng)) > -1 || lng.toLowerCase().indexOf('-arab') > 1 ? 'rtl' : 'ltr';
    }
    static createInstance() {
        let options = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};
        let callback = arguments.length > 1 ? arguments[1] : undefined;
        return new $1bac384020b50752$var$I18n(options, callback);
    }
    cloneInstance() {
        let options = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};
        let callback = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : $1bac384020b50752$var$noop;
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
            clone.store = new $1bac384020b50752$var$ResourceStore(this.store.data, mergedOptions);
            clone.services.resourceStore = clone.store;
        }
        clone.translator = new $1bac384020b50752$var$Translator(clone.services, mergedOptions);
        clone.translator.on('*', function(event) {
            for(var _len4 = arguments.length, args = new Array(_len4 > 1 ? _len4 - 1 : 0), _key4 = 1; _key4 < _len4; _key4++)args[_key4 - 1] = arguments[_key4];
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
$1bac384020b50752$export$2e2bcd8739ae039.createInstance = $1bac384020b50752$var$I18n.createInstance;
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
    activeClasses: 'bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-white',
    inactiveClasses: 'text-gray-500 dark:text-gray-400',
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
    backdropClasses: 'bg-gray-900/50 dark:bg-gray-900/80 fixed inset-0 z-40',
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
    backdropClasses: 'bg-gray-900/50 dark:bg-gray-900/80 fixed inset-0 z-30',
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
    activeClasses: 'text-blue-600 hover:text-blue-600 dark:text-blue-500 dark:hover:text-blue-500 border-blue-600 dark:border-blue-500',
    inactiveClasses: 'dark:border-transparent text-gray-500 hover:text-gray-600 dark:text-gray-400 border-gray-100 hover:border-gray-300 dark:border-gray-700 dark:hover:text-gray-300',
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
    nextArrow: '<svg class="w-4 h-4 rtl:rotate-180 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 10"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M1 5h12m0 0L9 1m4 4L9 9"/></svg>',
    orientation: 'auto',
    pickLevel: 0,
    prevArrow: '<svg class="w-4 h-4 rtl:rotate-180 text-gray-800 dark:text-white" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 10"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5H1m0 0 4 4M1 5l4-4"/></svg>',
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
var $a18c32cb29f9cbaf$var$range = document.createRange();
function $a18c32cb29f9cbaf$var$parseHTML(html) {
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
var $a18c32cb29f9cbaf$var$pickerTemplate = $a18c32cb29f9cbaf$var$optimizeTemplateHTML("<div class=\"datepicker hidden\">\n  <div class=\"datepicker-picker inline-block rounded-lg bg-white dark:bg-gray-700 shadow-lg p-4\">\n    <div class=\"datepicker-header\">\n      <div class=\"datepicker-title bg-white dark:bg-gray-700 dark:text-white px-2 py-3 text-center font-semibold\"></div>\n      <div class=\"datepicker-controls flex justify-between mb-2\">\n        <button type=\"button\" class=\"bg-white dark:bg-gray-700 rounded-lg text-gray-500 dark:text-white hover:bg-gray-100 dark:hover:bg-gray-600 hover:text-gray-900 dark:hover:text-white text-lg p-2.5 focus:outline-none focus:ring-2 focus:ring-gray-200 prev-btn\"></button>\n        <button type=\"button\" class=\"text-sm rounded-lg text-gray-900 dark:text-white bg-white dark:bg-gray-700 font-semibold py-2.5 px-5 hover:bg-gray-100 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-200 view-switch\"></button>\n        <button type=\"button\" class=\"bg-white dark:bg-gray-700 rounded-lg text-gray-500 dark:text-white hover:bg-gray-100 dark:hover:bg-gray-600 hover:text-gray-900 dark:hover:text-white text-lg p-2.5 focus:outline-none focus:ring-2 focus:ring-gray-200 next-btn\"></button>\n      </div>\n    </div>\n    <div class=\"datepicker-main p-1\"></div>\n    <div class=\"datepicker-footer\">\n      <div class=\"datepicker-controls flex space-x-2 rtl:space-x-reverse mt-2\">\n        <button type=\"button\" class=\"%buttonClass% today-btn text-white bg-blue-700 !bg-primary-700 dark:bg-blue-600 dark:!bg-primary-600 hover:bg-blue-800 hover:!bg-primary-800 dark:hover:bg-blue-700 dark:hover:!bg-primary-700 focus:ring-4 focus:ring-blue-300 focus:!ring-primary-300 font-medium rounded-lg text-sm px-5 py-2 text-center w-1/2\"></button>\n        <button type=\"button\" class=\"%buttonClass% clear-btn text-gray-900 dark:text-white bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-600 focus:ring-4 focus:ring-blue-300 focus:!ring-primary-300 font-medium rounded-lg text-sm px-5 py-2 text-center w-1/2\"></button>\n      </div>\n    </div>\n  </div>\n</div>");
var $a18c32cb29f9cbaf$var$daysTemplate = $a18c32cb29f9cbaf$var$optimizeTemplateHTML("<div class=\"days\">\n  <div class=\"days-of-week grid grid-cols-7 mb-1\">".concat($a18c32cb29f9cbaf$var$createTagRepeat('span', 7, {
    "class": 'dow block flex-1 leading-9 border-0 rounded-lg cursor-default text-center text-gray-900 font-semibold text-sm'
}), "</div>\n  <div class=\"datepicker-grid w-64 grid grid-cols-7\">").concat($a18c32cb29f9cbaf$var$createTagRepeat('span', 42, {
    "class": 'block flex-1 leading-9 border-0 rounded-lg cursor-default text-center text-gray-900 font-semibold text-sm h-6 leading-6 text-sm font-medium text-gray-500 dark:text-gray-400'
}), "</div>\n</div>"));
var $a18c32cb29f9cbaf$var$calendarWeeksTemplate = $a18c32cb29f9cbaf$var$optimizeTemplateHTML("<div class=\"calendar-weeks\">\n  <div class=\"days-of-week flex\"><span class=\"dow h-6 leading-6 text-sm font-medium text-gray-500 dark:text-gray-400\"></span></div>\n  <div class=\"weeks\">".concat($a18c32cb29f9cbaf$var$createTagRepeat('span', 6, {
    "class": 'week block flex-1 leading-9 border-0 rounded-lg cursor-default text-center text-gray-900 font-semibold text-sm'
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
                    el.className = _this.daysOfWeekDisabled.includes(dow) ? 'dow disabled text-center h-6 leading-6 text-sm font-medium text-gray-500 dark:text-gray-400 cursor-not-allowed' : 'dow text-center h-6 leading-6 text-sm font-medium text-gray-500 dark:text-gray-400';
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
                    el.className = "datepicker-cell hover:bg-gray-100 dark:hover:bg-gray-600 block flex-1 leading-9 border-0 rounded-lg cursor-pointer text-center text-gray-900 dark:text-white font-semibold text-sm ".concat(_this2.cellClass);
                    el.dataset.date = current;
                    el.textContent = date.getDate();
                    if (current < _this2.first) classList.add('prev', 'text-gray-500', 'dark:text-white');
                    else if (current > _this2.last) classList.add('next', 'text-gray-500', 'dark:text-white');
                    if (_this2.today === current) classList.add('today', 'bg-gray-100', 'dark:bg-gray-600');
                    if (current < _this2.minDate || current > _this2.maxDate || _this2.disabled.includes(current)) {
                        classList.add('disabled', 'cursor-not-allowed', 'text-gray-400', 'dark:text-gray-500');
                        classList.remove('hover:bg-gray-100', 'dark:hover:bg-gray-600', 'text-gray-900', 'dark:text-white', 'cursor-pointer');
                    }
                    if (_this2.daysOfWeekDisabled.includes(day)) {
                        classList.add('disabled', 'cursor-not-allowed', 'text-gray-400', 'dark:text-gray-500');
                        classList.remove('hover:bg-gray-100', 'dark:hover:bg-gray-600', 'text-gray-900', 'dark:text-white', 'cursor-pointer');
                        $a18c32cb29f9cbaf$var$pushUnique(_this2.disabled, current);
                    }
                    if (_this2.daysOfWeekHighlighted.includes(day)) classList.add('highlighted');
                    if (_this2.range) {
                        var _this2$range = $a18c32cb29f9cbaf$var$_slicedToArray(_this2.range, 2), rangeStart = _this2$range[0], rangeEnd = _this2$range[1];
                        if (current > rangeStart && current < rangeEnd) {
                            classList.add('range', 'bg-gray-200', 'dark:bg-gray-600');
                            classList.remove('rounded-lg', 'rounded-l-lg', 'rounded-r-lg');
                        }
                        if (current === rangeStart) {
                            classList.add('range-start', 'bg-gray-100', 'dark:bg-gray-600', 'rounded-l-lg');
                            classList.remove('rounded-lg', 'rounded-r-lg');
                        }
                        if (current === rangeEnd) {
                            classList.add('range-end', 'bg-gray-100', 'dark:bg-gray-600', 'rounded-r-lg');
                            classList.remove('rounded-lg', 'rounded-l-lg');
                        }
                    }
                    if (_this2.selected.includes(current)) {
                        classList.add('selected', 'bg-blue-700', '!bg-primary-700', 'text-white', 'dark:bg-blue-600', 'dark:!bg-primary-600', 'dark:text-white');
                        classList.remove('text-gray-900', 'text-gray-500', 'hover:bg-gray-100', 'dark:text-white', 'dark:hover:bg-gray-600', 'dark:bg-gray-600', 'bg-gray-100', 'bg-gray-200');
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
                    el.classList.remove('range', 'range-start', 'range-end', 'selected', 'bg-blue-700', '!bg-primary-700', 'text-white', 'dark:bg-blue-600', 'dark:!bg-primary-600', 'dark:text-white', 'focused');
                    el.classList.add('text-gray-900', 'rounded-lg', 'dark:text-white');
                });
                Array.from(this.grid.children).forEach(function(el) {
                    var current = Number(el.dataset.date);
                    var classList = el.classList;
                    classList.remove('bg-gray-200', 'dark:bg-gray-600', 'rounded-l-lg', 'rounded-r-lg');
                    if (current > rangeStart && current < rangeEnd) {
                        classList.add('range', 'bg-gray-200', 'dark:bg-gray-600');
                        classList.remove('rounded-lg');
                    }
                    if (current === rangeStart) {
                        classList.add('range-start', 'bg-gray-200', 'dark:bg-gray-600', 'rounded-l-lg');
                        classList.remove('rounded-lg');
                    }
                    if (current === rangeEnd) {
                        classList.add('range-end', 'bg-gray-200', 'dark:bg-gray-600', 'rounded-r-lg');
                        classList.remove('rounded-lg');
                    }
                    if (_this3.selected.includes(current)) {
                        classList.add('selected', 'bg-blue-700', '!bg-primary-700', 'text-white', 'dark:bg-blue-600', 'dark:!bg-primary-600', 'dark:text-white');
                        classList.remove('text-gray-900', 'hover:bg-gray-100', 'dark:text-white', 'dark:hover:bg-gray-600', 'bg-gray-100', 'bg-gray-200', 'dark:bg-gray-600');
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
                    el.className = "datepicker-cell hover:bg-gray-100 dark:hover:bg-gray-600 block flex-1 leading-9 border-0 rounded-lg cursor-pointer text-center text-gray-900 dark:text-white font-semibold text-sm ".concat(_this.cellClass);
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
                        classList.add('selected', 'bg-blue-700', '!bg-primary-700', 'text-white', 'dark:bg-blue-600', 'dark:!bg-primary-600', 'dark:text-white');
                        classList.remove('text-gray-900', 'hover:bg-gray-100', 'dark:text-white', 'dark:hover:bg-gray-600');
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
                    el.classList.remove('range', 'range-start', 'range-end', 'selected', 'bg-blue-700', '!bg-primary-700', 'dark:bg-blue-600', 'dark:!bg-primary-700', 'dark:text-white', 'text-white', 'focused');
                    el.classList.add('text-gray-900', 'hover:bg-gray-100', 'dark:text-white', 'dark:hover:bg-gray-600');
                });
                Array.from(this.grid.children).forEach(function(el, index) {
                    var classList = el.classList;
                    if (index > rangeStart && index < rangeEnd) classList.add('range');
                    if (index === rangeStart) classList.add('range-start');
                    if (index === rangeEnd) classList.add('range-end');
                    if (selected.includes(index)) {
                        classList.add('selected', 'bg-blue-700', '!bg-primary-700', 'text-white', 'dark:bg-blue-600', 'dark:!bg-primary-600', 'dark:text-white');
                        classList.remove('text-gray-900', 'hover:bg-gray-100', 'dark:text-white', 'dark:hover:bg-gray-600');
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
                    el.className = "datepicker-cell hover:bg-gray-100 dark:hover:bg-gray-600 block flex-1 leading-9 border-0 rounded-lg cursor-pointer text-center text-gray-900 dark:text-white font-semibold text-sm ".concat(_this2.cellClass);
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
                        classList.add('selected', 'bg-blue-700', '!bg-primary-700', 'text-white', 'dark:bg-blue-600', 'dark:!bg-primary-600', 'dark:text-white');
                        classList.remove('text-gray-900', 'hover:bg-gray-100', 'dark:text-white', 'dark:hover:bg-gray-600');
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
                    el.classList.remove('range', 'range-start', 'range-end', 'selected', 'bg-blue-700', '!bg-primary-700', 'text-white', 'dark:bg-blue-600', 'dark!bg-primary-600', 'dark:text-white', 'focused');
                });
                Array.from(this.grid.children).forEach(function(el) {
                    var current = Number(el.textContent);
                    var classList = el.classList;
                    if (current > rangeStart && current < rangeEnd) classList.add('range');
                    if (current === rangeStart) classList.add('range-start');
                    if (current === rangeEnd) classList.add('range-end');
                    if (_this3.selected.includes(current)) {
                        classList.add('selected', 'bg-blue-700', '!bg-primary-700', 'text-white', 'dark:bg-blue-600', 'dark:!bg-primary-600', 'dark:text-white');
                        classList.remove('text-gray-900', 'hover:bg-gray-100', 'dark:text-white', 'dark:hover:bg-gray-600');
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
                this.inputField.classList.add('in-edit', 'border-blue-700', '!border-primary-700');
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
                this.inputField.classList.remove('in-edit', 'border-blue-700', '!border-primary-700');
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


const { slice: $199830a05f92d3d0$var$slice, forEach: $199830a05f92d3d0$var$forEach } = [];
function $199830a05f92d3d0$var$defaults(obj) {
    $199830a05f92d3d0$var$forEach.call($199830a05f92d3d0$var$slice.call(arguments, 1), (source)=>{
        if (source) {
            for(const prop in source)if (obj[prop] === undefined) obj[prop] = source[prop];
        }
    });
    return obj;
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
        document.cookie = $199830a05f92d3d0$var$serializeCookie(name, encodeURIComponent(value), cookieOptions);
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
    remove (name) {
        this.create(name, '', -1);
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
let $199830a05f92d3d0$var$hasLocalStorageSupport = null;
const $199830a05f92d3d0$var$localStorageAvailable = ()=>{
    if ($199830a05f92d3d0$var$hasLocalStorageSupport !== null) return $199830a05f92d3d0$var$hasLocalStorageSupport;
    try {
        $199830a05f92d3d0$var$hasLocalStorageSupport = window !== 'undefined' && window.localStorage !== null;
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
        $199830a05f92d3d0$var$hasSessionStorageSupport = window !== 'undefined' && window.sessionStorage !== null;
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
        detected = detected.map((d)=>this.options.convertDetectedLanguage(d));
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


// lib/index.ts
var $040ba8e83b63f6c5$exports = {};
var $41148a221f03ee13$exports = {};
/* MIT license */ var $72c9eba4dd26aa1d$exports = {};
'use strict';
$72c9eba4dd26aa1d$exports = {
    "aliceblue": [
        240,
        248,
        255
    ],
    "antiquewhite": [
        250,
        235,
        215
    ],
    "aqua": [
        0,
        255,
        255
    ],
    "aquamarine": [
        127,
        255,
        212
    ],
    "azure": [
        240,
        255,
        255
    ],
    "beige": [
        245,
        245,
        220
    ],
    "bisque": [
        255,
        228,
        196
    ],
    "black": [
        0,
        0,
        0
    ],
    "blanchedalmond": [
        255,
        235,
        205
    ],
    "blue": [
        0,
        0,
        255
    ],
    "blueviolet": [
        138,
        43,
        226
    ],
    "brown": [
        165,
        42,
        42
    ],
    "burlywood": [
        222,
        184,
        135
    ],
    "cadetblue": [
        95,
        158,
        160
    ],
    "chartreuse": [
        127,
        255,
        0
    ],
    "chocolate": [
        210,
        105,
        30
    ],
    "coral": [
        255,
        127,
        80
    ],
    "cornflowerblue": [
        100,
        149,
        237
    ],
    "cornsilk": [
        255,
        248,
        220
    ],
    "crimson": [
        220,
        20,
        60
    ],
    "cyan": [
        0,
        255,
        255
    ],
    "darkblue": [
        0,
        0,
        139
    ],
    "darkcyan": [
        0,
        139,
        139
    ],
    "darkgoldenrod": [
        184,
        134,
        11
    ],
    "darkgray": [
        169,
        169,
        169
    ],
    "darkgreen": [
        0,
        100,
        0
    ],
    "darkgrey": [
        169,
        169,
        169
    ],
    "darkkhaki": [
        189,
        183,
        107
    ],
    "darkmagenta": [
        139,
        0,
        139
    ],
    "darkolivegreen": [
        85,
        107,
        47
    ],
    "darkorange": [
        255,
        140,
        0
    ],
    "darkorchid": [
        153,
        50,
        204
    ],
    "darkred": [
        139,
        0,
        0
    ],
    "darksalmon": [
        233,
        150,
        122
    ],
    "darkseagreen": [
        143,
        188,
        143
    ],
    "darkslateblue": [
        72,
        61,
        139
    ],
    "darkslategray": [
        47,
        79,
        79
    ],
    "darkslategrey": [
        47,
        79,
        79
    ],
    "darkturquoise": [
        0,
        206,
        209
    ],
    "darkviolet": [
        148,
        0,
        211
    ],
    "deeppink": [
        255,
        20,
        147
    ],
    "deepskyblue": [
        0,
        191,
        255
    ],
    "dimgray": [
        105,
        105,
        105
    ],
    "dimgrey": [
        105,
        105,
        105
    ],
    "dodgerblue": [
        30,
        144,
        255
    ],
    "firebrick": [
        178,
        34,
        34
    ],
    "floralwhite": [
        255,
        250,
        240
    ],
    "forestgreen": [
        34,
        139,
        34
    ],
    "fuchsia": [
        255,
        0,
        255
    ],
    "gainsboro": [
        220,
        220,
        220
    ],
    "ghostwhite": [
        248,
        248,
        255
    ],
    "gold": [
        255,
        215,
        0
    ],
    "goldenrod": [
        218,
        165,
        32
    ],
    "gray": [
        128,
        128,
        128
    ],
    "green": [
        0,
        128,
        0
    ],
    "greenyellow": [
        173,
        255,
        47
    ],
    "grey": [
        128,
        128,
        128
    ],
    "honeydew": [
        240,
        255,
        240
    ],
    "hotpink": [
        255,
        105,
        180
    ],
    "indianred": [
        205,
        92,
        92
    ],
    "indigo": [
        75,
        0,
        130
    ],
    "ivory": [
        255,
        255,
        240
    ],
    "khaki": [
        240,
        230,
        140
    ],
    "lavender": [
        230,
        230,
        250
    ],
    "lavenderblush": [
        255,
        240,
        245
    ],
    "lawngreen": [
        124,
        252,
        0
    ],
    "lemonchiffon": [
        255,
        250,
        205
    ],
    "lightblue": [
        173,
        216,
        230
    ],
    "lightcoral": [
        240,
        128,
        128
    ],
    "lightcyan": [
        224,
        255,
        255
    ],
    "lightgoldenrodyellow": [
        250,
        250,
        210
    ],
    "lightgray": [
        211,
        211,
        211
    ],
    "lightgreen": [
        144,
        238,
        144
    ],
    "lightgrey": [
        211,
        211,
        211
    ],
    "lightpink": [
        255,
        182,
        193
    ],
    "lightsalmon": [
        255,
        160,
        122
    ],
    "lightseagreen": [
        32,
        178,
        170
    ],
    "lightskyblue": [
        135,
        206,
        250
    ],
    "lightslategray": [
        119,
        136,
        153
    ],
    "lightslategrey": [
        119,
        136,
        153
    ],
    "lightsteelblue": [
        176,
        196,
        222
    ],
    "lightyellow": [
        255,
        255,
        224
    ],
    "lime": [
        0,
        255,
        0
    ],
    "limegreen": [
        50,
        205,
        50
    ],
    "linen": [
        250,
        240,
        230
    ],
    "magenta": [
        255,
        0,
        255
    ],
    "maroon": [
        128,
        0,
        0
    ],
    "mediumaquamarine": [
        102,
        205,
        170
    ],
    "mediumblue": [
        0,
        0,
        205
    ],
    "mediumorchid": [
        186,
        85,
        211
    ],
    "mediumpurple": [
        147,
        112,
        219
    ],
    "mediumseagreen": [
        60,
        179,
        113
    ],
    "mediumslateblue": [
        123,
        104,
        238
    ],
    "mediumspringgreen": [
        0,
        250,
        154
    ],
    "mediumturquoise": [
        72,
        209,
        204
    ],
    "mediumvioletred": [
        199,
        21,
        133
    ],
    "midnightblue": [
        25,
        25,
        112
    ],
    "mintcream": [
        245,
        255,
        250
    ],
    "mistyrose": [
        255,
        228,
        225
    ],
    "moccasin": [
        255,
        228,
        181
    ],
    "navajowhite": [
        255,
        222,
        173
    ],
    "navy": [
        0,
        0,
        128
    ],
    "oldlace": [
        253,
        245,
        230
    ],
    "olive": [
        128,
        128,
        0
    ],
    "olivedrab": [
        107,
        142,
        35
    ],
    "orange": [
        255,
        165,
        0
    ],
    "orangered": [
        255,
        69,
        0
    ],
    "orchid": [
        218,
        112,
        214
    ],
    "palegoldenrod": [
        238,
        232,
        170
    ],
    "palegreen": [
        152,
        251,
        152
    ],
    "paleturquoise": [
        175,
        238,
        238
    ],
    "palevioletred": [
        219,
        112,
        147
    ],
    "papayawhip": [
        255,
        239,
        213
    ],
    "peachpuff": [
        255,
        218,
        185
    ],
    "peru": [
        205,
        133,
        63
    ],
    "pink": [
        255,
        192,
        203
    ],
    "plum": [
        221,
        160,
        221
    ],
    "powderblue": [
        176,
        224,
        230
    ],
    "purple": [
        128,
        0,
        128
    ],
    "rebeccapurple": [
        102,
        51,
        153
    ],
    "red": [
        255,
        0,
        0
    ],
    "rosybrown": [
        188,
        143,
        143
    ],
    "royalblue": [
        65,
        105,
        225
    ],
    "saddlebrown": [
        139,
        69,
        19
    ],
    "salmon": [
        250,
        128,
        114
    ],
    "sandybrown": [
        244,
        164,
        96
    ],
    "seagreen": [
        46,
        139,
        87
    ],
    "seashell": [
        255,
        245,
        238
    ],
    "sienna": [
        160,
        82,
        45
    ],
    "silver": [
        192,
        192,
        192
    ],
    "skyblue": [
        135,
        206,
        235
    ],
    "slateblue": [
        106,
        90,
        205
    ],
    "slategray": [
        112,
        128,
        144
    ],
    "slategrey": [
        112,
        128,
        144
    ],
    "snow": [
        255,
        250,
        250
    ],
    "springgreen": [
        0,
        255,
        127
    ],
    "steelblue": [
        70,
        130,
        180
    ],
    "tan": [
        210,
        180,
        140
    ],
    "teal": [
        0,
        128,
        128
    ],
    "thistle": [
        216,
        191,
        216
    ],
    "tomato": [
        255,
        99,
        71
    ],
    "turquoise": [
        64,
        224,
        208
    ],
    "violet": [
        238,
        130,
        238
    ],
    "wheat": [
        245,
        222,
        179
    ],
    "white": [
        255,
        255,
        255
    ],
    "whitesmoke": [
        245,
        245,
        245
    ],
    "yellow": [
        255,
        255,
        0
    ],
    "yellowgreen": [
        154,
        205,
        50
    ]
};


var $e97dd6e5b87d46f7$exports = {};
'use strict';
var $a8cbc72403bdcef6$exports = {};
$a8cbc72403bdcef6$exports = function isArrayish(obj) {
    if (!obj || typeof obj === 'string') return false;
    return obj instanceof Array || Array.isArray(obj) || obj.length >= 0 && (obj.splice instanceof Function || Object.getOwnPropertyDescriptor(obj, obj.length - 1) && obj.constructor.name !== 'String');
};


var $e97dd6e5b87d46f7$var$concat = Array.prototype.concat;
var $e97dd6e5b87d46f7$var$slice = Array.prototype.slice;
var $e97dd6e5b87d46f7$var$swizzle = $e97dd6e5b87d46f7$exports = function swizzle(args) {
    var results = [];
    for(var i = 0, len = args.length; i < len; i++){
        var arg = args[i];
        if ($a8cbc72403bdcef6$exports(arg)) // http://jsperf.com/javascript-array-concat-vs-push/98
        results = $e97dd6e5b87d46f7$var$concat.call(results, $e97dd6e5b87d46f7$var$slice.call(arg));
        else results.push(arg);
    }
    return results;
};
$e97dd6e5b87d46f7$var$swizzle.wrap = function(fn) {
    return function() {
        return fn($e97dd6e5b87d46f7$var$swizzle(arguments));
    };
};


var $41148a221f03ee13$var$hasOwnProperty = Object.hasOwnProperty;
var $41148a221f03ee13$var$reverseNames = Object.create(null);
// create a list of reverse color names
for(var $41148a221f03ee13$var$name in $72c9eba4dd26aa1d$exports)if ($41148a221f03ee13$var$hasOwnProperty.call($72c9eba4dd26aa1d$exports, $41148a221f03ee13$var$name)) $41148a221f03ee13$var$reverseNames[$72c9eba4dd26aa1d$exports[$41148a221f03ee13$var$name]] = $41148a221f03ee13$var$name;
var $41148a221f03ee13$var$cs = $41148a221f03ee13$exports = {
    to: {},
    get: {}
};
$41148a221f03ee13$var$cs.get = function(string) {
    var prefix = string.substring(0, 3).toLowerCase();
    var val;
    var model;
    switch(prefix){
        case 'hsl':
            val = $41148a221f03ee13$var$cs.get.hsl(string);
            model = 'hsl';
            break;
        case 'hwb':
            val = $41148a221f03ee13$var$cs.get.hwb(string);
            model = 'hwb';
            break;
        default:
            val = $41148a221f03ee13$var$cs.get.rgb(string);
            model = 'rgb';
            break;
    }
    if (!val) return null;
    return {
        model: model,
        value: val
    };
};
$41148a221f03ee13$var$cs.get.rgb = function(string) {
    if (!string) return null;
    var abbr = /^#([a-f0-9]{3,4})$/i;
    var hex = /^#([a-f0-9]{6})([a-f0-9]{2})?$/i;
    var rgba = /^rgba?\(\s*([+-]?\d+)(?=[\s,])\s*(?:,\s*)?([+-]?\d+)(?=[\s,])\s*(?:,\s*)?([+-]?\d+)\s*(?:[,|\/]\s*([+-]?[\d\.]+)(%?)\s*)?\)$/;
    var per = /^rgba?\(\s*([+-]?[\d\.]+)\%\s*,?\s*([+-]?[\d\.]+)\%\s*,?\s*([+-]?[\d\.]+)\%\s*(?:[,|\/]\s*([+-]?[\d\.]+)(%?)\s*)?\)$/;
    var keyword = /^(\w+)$/;
    var rgb = [
        0,
        0,
        0,
        1
    ];
    var match;
    var i;
    var hexAlpha;
    if (match = string.match(hex)) {
        hexAlpha = match[2];
        match = match[1];
        for(i = 0; i < 3; i++){
            // https://jsperf.com/slice-vs-substr-vs-substring-methods-long-string/19
            var i2 = i * 2;
            rgb[i] = parseInt(match.slice(i2, i2 + 2), 16);
        }
        if (hexAlpha) rgb[3] = parseInt(hexAlpha, 16) / 255;
    } else if (match = string.match(abbr)) {
        match = match[1];
        hexAlpha = match[3];
        for(i = 0; i < 3; i++)rgb[i] = parseInt(match[i] + match[i], 16);
        if (hexAlpha) rgb[3] = parseInt(hexAlpha + hexAlpha, 16) / 255;
    } else if (match = string.match(rgba)) {
        for(i = 0; i < 3; i++)rgb[i] = parseInt(match[i + 1], 0);
        if (match[4]) {
            if (match[5]) rgb[3] = parseFloat(match[4]) * 0.01;
            else rgb[3] = parseFloat(match[4]);
        }
    } else if (match = string.match(per)) {
        for(i = 0; i < 3; i++)rgb[i] = Math.round(parseFloat(match[i + 1]) * 2.55);
        if (match[4]) {
            if (match[5]) rgb[3] = parseFloat(match[4]) * 0.01;
            else rgb[3] = parseFloat(match[4]);
        }
    } else if (match = string.match(keyword)) {
        if (match[1] === 'transparent') return [
            0,
            0,
            0,
            0
        ];
        if (!$41148a221f03ee13$var$hasOwnProperty.call($72c9eba4dd26aa1d$exports, match[1])) return null;
        rgb = $72c9eba4dd26aa1d$exports[match[1]];
        rgb[3] = 1;
        return rgb;
    } else return null;
    for(i = 0; i < 3; i++)rgb[i] = $41148a221f03ee13$var$clamp(rgb[i], 0, 255);
    rgb[3] = $41148a221f03ee13$var$clamp(rgb[3], 0, 1);
    return rgb;
};
$41148a221f03ee13$var$cs.get.hsl = function(string) {
    if (!string) return null;
    var hsl = /^hsla?\(\s*([+-]?(?:\d{0,3}\.)?\d+)(?:deg)?\s*,?\s*([+-]?[\d\.]+)%\s*,?\s*([+-]?[\d\.]+)%\s*(?:[,|\/]\s*([+-]?(?=\.\d|\d)(?:0|[1-9]\d*)?(?:\.\d*)?(?:[eE][+-]?\d+)?)\s*)?\)$/;
    var match = string.match(hsl);
    if (match) {
        var alpha = parseFloat(match[4]);
        var h = (parseFloat(match[1]) % 360 + 360) % 360;
        var s = $41148a221f03ee13$var$clamp(parseFloat(match[2]), 0, 100);
        var l = $41148a221f03ee13$var$clamp(parseFloat(match[3]), 0, 100);
        var a = $41148a221f03ee13$var$clamp(isNaN(alpha) ? 1 : alpha, 0, 1);
        return [
            h,
            s,
            l,
            a
        ];
    }
    return null;
};
$41148a221f03ee13$var$cs.get.hwb = function(string) {
    if (!string) return null;
    var hwb = /^hwb\(\s*([+-]?\d{0,3}(?:\.\d+)?)(?:deg)?\s*,\s*([+-]?[\d\.]+)%\s*,\s*([+-]?[\d\.]+)%\s*(?:,\s*([+-]?(?=\.\d|\d)(?:0|[1-9]\d*)?(?:\.\d*)?(?:[eE][+-]?\d+)?)\s*)?\)$/;
    var match = string.match(hwb);
    if (match) {
        var alpha = parseFloat(match[4]);
        var h = (parseFloat(match[1]) % 360 + 360) % 360;
        var w = $41148a221f03ee13$var$clamp(parseFloat(match[2]), 0, 100);
        var b = $41148a221f03ee13$var$clamp(parseFloat(match[3]), 0, 100);
        var a = $41148a221f03ee13$var$clamp(isNaN(alpha) ? 1 : alpha, 0, 1);
        return [
            h,
            w,
            b,
            a
        ];
    }
    return null;
};
$41148a221f03ee13$var$cs.to.hex = function() {
    var rgba = $e97dd6e5b87d46f7$exports(arguments);
    return '#' + $41148a221f03ee13$var$hexDouble(rgba[0]) + $41148a221f03ee13$var$hexDouble(rgba[1]) + $41148a221f03ee13$var$hexDouble(rgba[2]) + (rgba[3] < 1 ? $41148a221f03ee13$var$hexDouble(Math.round(rgba[3] * 255)) : '');
};
$41148a221f03ee13$var$cs.to.rgb = function() {
    var rgba = $e97dd6e5b87d46f7$exports(arguments);
    return rgba.length < 4 || rgba[3] === 1 ? 'rgb(' + Math.round(rgba[0]) + ', ' + Math.round(rgba[1]) + ', ' + Math.round(rgba[2]) + ')' : 'rgba(' + Math.round(rgba[0]) + ', ' + Math.round(rgba[1]) + ', ' + Math.round(rgba[2]) + ', ' + rgba[3] + ')';
};
$41148a221f03ee13$var$cs.to.rgb.percent = function() {
    var rgba = $e97dd6e5b87d46f7$exports(arguments);
    var r = Math.round(rgba[0] / 255 * 100);
    var g = Math.round(rgba[1] / 255 * 100);
    var b = Math.round(rgba[2] / 255 * 100);
    return rgba.length < 4 || rgba[3] === 1 ? 'rgb(' + r + '%, ' + g + '%, ' + b + '%)' : 'rgba(' + r + '%, ' + g + '%, ' + b + '%, ' + rgba[3] + ')';
};
$41148a221f03ee13$var$cs.to.hsl = function() {
    var hsla = $e97dd6e5b87d46f7$exports(arguments);
    return hsla.length < 4 || hsla[3] === 1 ? 'hsl(' + hsla[0] + ', ' + hsla[1] + '%, ' + hsla[2] + '%)' : 'hsla(' + hsla[0] + ', ' + hsla[1] + '%, ' + hsla[2] + '%, ' + hsla[3] + ')';
};
// hwb is a bit different than rgb(a) & hsl(a) since there is no alpha specific syntax
// (hwb have alpha optional & 1 is default value)
$41148a221f03ee13$var$cs.to.hwb = function() {
    var hwba = $e97dd6e5b87d46f7$exports(arguments);
    var a = '';
    if (hwba.length >= 4 && hwba[3] !== 1) a = ', ' + hwba[3];
    return 'hwb(' + hwba[0] + ', ' + hwba[1] + '%, ' + hwba[2] + '%' + a + ')';
};
$41148a221f03ee13$var$cs.to.keyword = function(rgb) {
    return $41148a221f03ee13$var$reverseNames[rgb.slice(0, 3)];
};
// helpers
function $41148a221f03ee13$var$clamp(num, min, max) {
    return Math.min(Math.max(min, num), max);
}
function $41148a221f03ee13$var$hexDouble(num) {
    var str = Math.round(num).toString(16).toUpperCase();
    return str.length < 2 ? '0' + str : str;
}


var $b0ce43f4286ff9ac$exports = {};
var $f48c6abccc789ef5$exports = {};
/* MIT license */ /* eslint-disable no-mixed-operators */ 
// NOTE: conversions should only return primitive values (i.e. arrays, or
//       values that give correct `typeof` results).
//       do not use box values types (i.e. Number(), String(), etc.)
const $f48c6abccc789ef5$var$reverseKeywords = {};
for (const key of Object.keys($72c9eba4dd26aa1d$exports))$f48c6abccc789ef5$var$reverseKeywords[$72c9eba4dd26aa1d$exports[key]] = key;
const $f48c6abccc789ef5$var$convert = {
    rgb: {
        channels: 3,
        labels: 'rgb'
    },
    hsl: {
        channels: 3,
        labels: 'hsl'
    },
    hsv: {
        channels: 3,
        labels: 'hsv'
    },
    hwb: {
        channels: 3,
        labels: 'hwb'
    },
    cmyk: {
        channels: 4,
        labels: 'cmyk'
    },
    xyz: {
        channels: 3,
        labels: 'xyz'
    },
    lab: {
        channels: 3,
        labels: 'lab'
    },
    lch: {
        channels: 3,
        labels: 'lch'
    },
    hex: {
        channels: 1,
        labels: [
            'hex'
        ]
    },
    keyword: {
        channels: 1,
        labels: [
            'keyword'
        ]
    },
    ansi16: {
        channels: 1,
        labels: [
            'ansi16'
        ]
    },
    ansi256: {
        channels: 1,
        labels: [
            'ansi256'
        ]
    },
    hcg: {
        channels: 3,
        labels: [
            'h',
            'c',
            'g'
        ]
    },
    apple: {
        channels: 3,
        labels: [
            'r16',
            'g16',
            'b16'
        ]
    },
    gray: {
        channels: 1,
        labels: [
            'gray'
        ]
    }
};
$f48c6abccc789ef5$exports = $f48c6abccc789ef5$var$convert;
// Hide .channels and .labels properties
for (const model of Object.keys($f48c6abccc789ef5$var$convert)){
    if (!('channels' in $f48c6abccc789ef5$var$convert[model])) throw new Error('missing channels property: ' + model);
    if (!('labels' in $f48c6abccc789ef5$var$convert[model])) throw new Error('missing channel labels property: ' + model);
    if ($f48c6abccc789ef5$var$convert[model].labels.length !== $f48c6abccc789ef5$var$convert[model].channels) throw new Error('channel and label counts mismatch: ' + model);
    const { channels: channels, labels: labels } = $f48c6abccc789ef5$var$convert[model];
    delete $f48c6abccc789ef5$var$convert[model].channels;
    delete $f48c6abccc789ef5$var$convert[model].labels;
    Object.defineProperty($f48c6abccc789ef5$var$convert[model], 'channels', {
        value: channels
    });
    Object.defineProperty($f48c6abccc789ef5$var$convert[model], 'labels', {
        value: labels
    });
}
$f48c6abccc789ef5$var$convert.rgb.hsl = function(rgb) {
    const r = rgb[0] / 255;
    const g = rgb[1] / 255;
    const b = rgb[2] / 255;
    const min = Math.min(r, g, b);
    const max = Math.max(r, g, b);
    const delta = max - min;
    let h;
    let s;
    if (max === min) h = 0;
    else if (r === max) h = (g - b) / delta;
    else if (g === max) h = 2 + (b - r) / delta;
    else if (b === max) h = 4 + (r - g) / delta;
    h = Math.min(h * 60, 360);
    if (h < 0) h += 360;
    const l = (min + max) / 2;
    if (max === min) s = 0;
    else if (l <= 0.5) s = delta / (max + min);
    else s = delta / (2 - max - min);
    return [
        h,
        s * 100,
        l * 100
    ];
};
$f48c6abccc789ef5$var$convert.rgb.hsv = function(rgb) {
    let rdif;
    let gdif;
    let bdif;
    let h;
    let s;
    const r = rgb[0] / 255;
    const g = rgb[1] / 255;
    const b = rgb[2] / 255;
    const v = Math.max(r, g, b);
    const diff = v - Math.min(r, g, b);
    const diffc = function(c) {
        return (v - c) / 6 / diff + 0.5;
    };
    if (diff === 0) {
        h = 0;
        s = 0;
    } else {
        s = diff / v;
        rdif = diffc(r);
        gdif = diffc(g);
        bdif = diffc(b);
        if (r === v) h = bdif - gdif;
        else if (g === v) h = 1 / 3 + rdif - bdif;
        else if (b === v) h = 2 / 3 + gdif - rdif;
        if (h < 0) h += 1;
        else if (h > 1) h -= 1;
    }
    return [
        h * 360,
        s * 100,
        v * 100
    ];
};
$f48c6abccc789ef5$var$convert.rgb.hwb = function(rgb) {
    const r = rgb[0];
    const g = rgb[1];
    let b = rgb[2];
    const h = $f48c6abccc789ef5$var$convert.rgb.hsl(rgb)[0];
    const w = 1 / 255 * Math.min(r, Math.min(g, b));
    b = 1 - 1 / 255 * Math.max(r, Math.max(g, b));
    return [
        h,
        w * 100,
        b * 100
    ];
};
$f48c6abccc789ef5$var$convert.rgb.cmyk = function(rgb) {
    const r = rgb[0] / 255;
    const g = rgb[1] / 255;
    const b = rgb[2] / 255;
    const k = Math.min(1 - r, 1 - g, 1 - b);
    const c = (1 - r - k) / (1 - k) || 0;
    const m = (1 - g - k) / (1 - k) || 0;
    const y = (1 - b - k) / (1 - k) || 0;
    return [
        c * 100,
        m * 100,
        y * 100,
        k * 100
    ];
};
function $f48c6abccc789ef5$var$comparativeDistance(x, y) {
    /*
		See https://en.m.wikipedia.org/wiki/Euclidean_distance#Squared_Euclidean_distance
	*/ return (x[0] - y[0]) ** 2 + (x[1] - y[1]) ** 2 + (x[2] - y[2]) ** 2;
}
$f48c6abccc789ef5$var$convert.rgb.keyword = function(rgb) {
    const reversed = $f48c6abccc789ef5$var$reverseKeywords[rgb];
    if (reversed) return reversed;
    let currentClosestDistance = Infinity;
    let currentClosestKeyword;
    for (const keyword of Object.keys($72c9eba4dd26aa1d$exports)){
        const value = $72c9eba4dd26aa1d$exports[keyword];
        // Compute comparative distance
        const distance = $f48c6abccc789ef5$var$comparativeDistance(rgb, value);
        // Check if its less, if so set as closest
        if (distance < currentClosestDistance) {
            currentClosestDistance = distance;
            currentClosestKeyword = keyword;
        }
    }
    return currentClosestKeyword;
};
$f48c6abccc789ef5$var$convert.keyword.rgb = function(keyword) {
    return $72c9eba4dd26aa1d$exports[keyword];
};
$f48c6abccc789ef5$var$convert.rgb.xyz = function(rgb) {
    let r = rgb[0] / 255;
    let g = rgb[1] / 255;
    let b = rgb[2] / 255;
    // Assume sRGB
    r = r > 0.04045 ? ((r + 0.055) / 1.055) ** 2.4 : r / 12.92;
    g = g > 0.04045 ? ((g + 0.055) / 1.055) ** 2.4 : g / 12.92;
    b = b > 0.04045 ? ((b + 0.055) / 1.055) ** 2.4 : b / 12.92;
    const x = r * 0.4124 + g * 0.3576 + b * 0.1805;
    const y = r * 0.2126 + g * 0.7152 + b * 0.0722;
    const z = r * 0.0193 + g * 0.1192 + b * 0.9505;
    return [
        x * 100,
        y * 100,
        z * 100
    ];
};
$f48c6abccc789ef5$var$convert.rgb.lab = function(rgb) {
    const xyz = $f48c6abccc789ef5$var$convert.rgb.xyz(rgb);
    let x = xyz[0];
    let y = xyz[1];
    let z = xyz[2];
    x /= 95.047;
    y /= 100;
    z /= 108.883;
    x = x > 0.008856 ? x ** (1 / 3) : 7.787 * x + 16 / 116;
    y = y > 0.008856 ? y ** (1 / 3) : 7.787 * y + 16 / 116;
    z = z > 0.008856 ? z ** (1 / 3) : 7.787 * z + 16 / 116;
    const l = 116 * y - 16;
    const a = 500 * (x - y);
    const b = 200 * (y - z);
    return [
        l,
        a,
        b
    ];
};
$f48c6abccc789ef5$var$convert.hsl.rgb = function(hsl) {
    const h = hsl[0] / 360;
    const s = hsl[1] / 100;
    const l = hsl[2] / 100;
    let t2;
    let t3;
    let val;
    if (s === 0) {
        val = l * 255;
        return [
            val,
            val,
            val
        ];
    }
    if (l < 0.5) t2 = l * (1 + s);
    else t2 = l + s - l * s;
    const t1 = 2 * l - t2;
    const rgb = [
        0,
        0,
        0
    ];
    for(let i = 0; i < 3; i++){
        t3 = h + 1 / 3 * -(i - 1);
        if (t3 < 0) t3++;
        if (t3 > 1) t3--;
        if (6 * t3 < 1) val = t1 + (t2 - t1) * 6 * t3;
        else if (2 * t3 < 1) val = t2;
        else if (3 * t3 < 2) val = t1 + (t2 - t1) * (2 / 3 - t3) * 6;
        else val = t1;
        rgb[i] = val * 255;
    }
    return rgb;
};
$f48c6abccc789ef5$var$convert.hsl.hsv = function(hsl) {
    const h = hsl[0];
    let s = hsl[1] / 100;
    let l = hsl[2] / 100;
    let smin = s;
    const lmin = Math.max(l, 0.01);
    l *= 2;
    s *= l <= 1 ? l : 2 - l;
    smin *= lmin <= 1 ? lmin : 2 - lmin;
    const v = (l + s) / 2;
    const sv = l === 0 ? 2 * smin / (lmin + smin) : 2 * s / (l + s);
    return [
        h,
        sv * 100,
        v * 100
    ];
};
$f48c6abccc789ef5$var$convert.hsv.rgb = function(hsv) {
    const h = hsv[0] / 60;
    const s = hsv[1] / 100;
    let v = hsv[2] / 100;
    const hi = Math.floor(h) % 6;
    const f = h - Math.floor(h);
    const p = 255 * v * (1 - s);
    const q = 255 * v * (1 - s * f);
    const t = 255 * v * (1 - s * (1 - f));
    v *= 255;
    switch(hi){
        case 0:
            return [
                v,
                t,
                p
            ];
        case 1:
            return [
                q,
                v,
                p
            ];
        case 2:
            return [
                p,
                v,
                t
            ];
        case 3:
            return [
                p,
                q,
                v
            ];
        case 4:
            return [
                t,
                p,
                v
            ];
        case 5:
            return [
                v,
                p,
                q
            ];
    }
};
$f48c6abccc789ef5$var$convert.hsv.hsl = function(hsv) {
    const h = hsv[0];
    const s = hsv[1] / 100;
    const v = hsv[2] / 100;
    const vmin = Math.max(v, 0.01);
    let sl;
    let l;
    l = (2 - s) * v;
    const lmin = (2 - s) * vmin;
    sl = s * vmin;
    sl /= lmin <= 1 ? lmin : 2 - lmin;
    sl = sl || 0;
    l /= 2;
    return [
        h,
        sl * 100,
        l * 100
    ];
};
// http://dev.w3.org/csswg/css-color/#hwb-to-rgb
$f48c6abccc789ef5$var$convert.hwb.rgb = function(hwb) {
    const h = hwb[0] / 360;
    let wh = hwb[1] / 100;
    let bl = hwb[2] / 100;
    const ratio = wh + bl;
    let f;
    // Wh + bl cant be > 1
    if (ratio > 1) {
        wh /= ratio;
        bl /= ratio;
    }
    const i = Math.floor(6 * h);
    const v = 1 - bl;
    f = 6 * h - i;
    if ((i & 0x01) !== 0) f = 1 - f;
    const n = wh + f * (v - wh); // Linear interpolation
    let r;
    let g;
    let b;
    /* eslint-disable max-statements-per-line,no-multi-spaces */ switch(i){
        default:
        case 6:
        case 0:
            r = v;
            g = n;
            b = wh;
            break;
        case 1:
            r = n;
            g = v;
            b = wh;
            break;
        case 2:
            r = wh;
            g = v;
            b = n;
            break;
        case 3:
            r = wh;
            g = n;
            b = v;
            break;
        case 4:
            r = n;
            g = wh;
            b = v;
            break;
        case 5:
            r = v;
            g = wh;
            b = n;
            break;
    }
    /* eslint-enable max-statements-per-line,no-multi-spaces */ return [
        r * 255,
        g * 255,
        b * 255
    ];
};
$f48c6abccc789ef5$var$convert.cmyk.rgb = function(cmyk) {
    const c = cmyk[0] / 100;
    const m = cmyk[1] / 100;
    const y = cmyk[2] / 100;
    const k = cmyk[3] / 100;
    const r = 1 - Math.min(1, c * (1 - k) + k);
    const g = 1 - Math.min(1, m * (1 - k) + k);
    const b = 1 - Math.min(1, y * (1 - k) + k);
    return [
        r * 255,
        g * 255,
        b * 255
    ];
};
$f48c6abccc789ef5$var$convert.xyz.rgb = function(xyz) {
    const x = xyz[0] / 100;
    const y = xyz[1] / 100;
    const z = xyz[2] / 100;
    let r;
    let g;
    let b;
    r = x * 3.2406 + y * -1.5372 + z * -0.4986;
    g = x * -0.9689 + y * 1.8758 + z * 0.0415;
    b = x * 0.0557 + y * -0.204 + z * 1.0570;
    // Assume sRGB
    r = r > 0.0031308 ? 1.055 * r ** (1.0 / 2.4) - 0.055 : r * 12.92;
    g = g > 0.0031308 ? 1.055 * g ** (1.0 / 2.4) - 0.055 : g * 12.92;
    b = b > 0.0031308 ? 1.055 * b ** (1.0 / 2.4) - 0.055 : b * 12.92;
    r = Math.min(Math.max(0, r), 1);
    g = Math.min(Math.max(0, g), 1);
    b = Math.min(Math.max(0, b), 1);
    return [
        r * 255,
        g * 255,
        b * 255
    ];
};
$f48c6abccc789ef5$var$convert.xyz.lab = function(xyz) {
    let x = xyz[0];
    let y = xyz[1];
    let z = xyz[2];
    x /= 95.047;
    y /= 100;
    z /= 108.883;
    x = x > 0.008856 ? x ** (1 / 3) : 7.787 * x + 16 / 116;
    y = y > 0.008856 ? y ** (1 / 3) : 7.787 * y + 16 / 116;
    z = z > 0.008856 ? z ** (1 / 3) : 7.787 * z + 16 / 116;
    const l = 116 * y - 16;
    const a = 500 * (x - y);
    const b = 200 * (y - z);
    return [
        l,
        a,
        b
    ];
};
$f48c6abccc789ef5$var$convert.lab.xyz = function(lab) {
    const l = lab[0];
    const a = lab[1];
    const b = lab[2];
    let x;
    let y;
    let z;
    y = (l + 16) / 116;
    x = a / 500 + y;
    z = y - b / 200;
    const y2 = y ** 3;
    const x2 = x ** 3;
    const z2 = z ** 3;
    y = y2 > 0.008856 ? y2 : (y - 16 / 116) / 7.787;
    x = x2 > 0.008856 ? x2 : (x - 16 / 116) / 7.787;
    z = z2 > 0.008856 ? z2 : (z - 16 / 116) / 7.787;
    x *= 95.047;
    y *= 100;
    z *= 108.883;
    return [
        x,
        y,
        z
    ];
};
$f48c6abccc789ef5$var$convert.lab.lch = function(lab) {
    const l = lab[0];
    const a = lab[1];
    const b = lab[2];
    let h;
    const hr = Math.atan2(b, a);
    h = hr * 360 / 2 / Math.PI;
    if (h < 0) h += 360;
    const c = Math.sqrt(a * a + b * b);
    return [
        l,
        c,
        h
    ];
};
$f48c6abccc789ef5$var$convert.lch.lab = function(lch) {
    const l = lch[0];
    const c = lch[1];
    const h = lch[2];
    const hr = h / 360 * 2 * Math.PI;
    const a = c * Math.cos(hr);
    const b = c * Math.sin(hr);
    return [
        l,
        a,
        b
    ];
};
$f48c6abccc789ef5$var$convert.rgb.ansi16 = function(args, saturation = null) {
    const [r, g, b] = args;
    let value = saturation === null ? $f48c6abccc789ef5$var$convert.rgb.hsv(args)[2] : saturation; // Hsv -> ansi16 optimization
    value = Math.round(value / 50);
    if (value === 0) return 30;
    let ansi = 30 + (Math.round(b / 255) << 2 | Math.round(g / 255) << 1 | Math.round(r / 255));
    if (value === 2) ansi += 60;
    return ansi;
};
$f48c6abccc789ef5$var$convert.hsv.ansi16 = function(args) {
    // Optimization here; we already know the value and don't need to get
    // it converted for us.
    return $f48c6abccc789ef5$var$convert.rgb.ansi16($f48c6abccc789ef5$var$convert.hsv.rgb(args), args[2]);
};
$f48c6abccc789ef5$var$convert.rgb.ansi256 = function(args) {
    const r = args[0];
    const g = args[1];
    const b = args[2];
    // We use the extended greyscale palette here, with the exception of
    // black and white. normal palette only has 4 greyscale shades.
    if (r === g && g === b) {
        if (r < 8) return 16;
        if (r > 248) return 231;
        return Math.round((r - 8) / 247 * 24) + 232;
    }
    const ansi = 16 + 36 * Math.round(r / 255 * 5) + 6 * Math.round(g / 255 * 5) + Math.round(b / 255 * 5);
    return ansi;
};
$f48c6abccc789ef5$var$convert.ansi16.rgb = function(args) {
    let color = args % 10;
    // Handle greyscale
    if (color === 0 || color === 7) {
        if (args > 50) color += 3.5;
        color = color / 10.5 * 255;
        return [
            color,
            color,
            color
        ];
    }
    const mult = (~~(args > 50) + 1) * 0.5;
    const r = (color & 1) * mult * 255;
    const g = (color >> 1 & 1) * mult * 255;
    const b = (color >> 2 & 1) * mult * 255;
    return [
        r,
        g,
        b
    ];
};
$f48c6abccc789ef5$var$convert.ansi256.rgb = function(args) {
    // Handle greyscale
    if (args >= 232) {
        const c = (args - 232) * 10 + 8;
        return [
            c,
            c,
            c
        ];
    }
    args -= 16;
    let rem;
    const r = Math.floor(args / 36) / 5 * 255;
    const g = Math.floor((rem = args % 36) / 6) / 5 * 255;
    const b = rem % 6 / 5 * 255;
    return [
        r,
        g,
        b
    ];
};
$f48c6abccc789ef5$var$convert.rgb.hex = function(args) {
    const integer = ((Math.round(args[0]) & 0xFF) << 16) + ((Math.round(args[1]) & 0xFF) << 8) + (Math.round(args[2]) & 0xFF);
    const string = integer.toString(16).toUpperCase();
    return '000000'.substring(string.length) + string;
};
$f48c6abccc789ef5$var$convert.hex.rgb = function(args) {
    const match = args.toString(16).match(/[a-f0-9]{6}|[a-f0-9]{3}/i);
    if (!match) return [
        0,
        0,
        0
    ];
    let colorString = match[0];
    if (match[0].length === 3) colorString = colorString.split('').map((char)=>{
        return char + char;
    }).join('');
    const integer = parseInt(colorString, 16);
    const r = integer >> 16 & 0xFF;
    const g = integer >> 8 & 0xFF;
    const b = integer & 0xFF;
    return [
        r,
        g,
        b
    ];
};
$f48c6abccc789ef5$var$convert.rgb.hcg = function(rgb) {
    const r = rgb[0] / 255;
    const g = rgb[1] / 255;
    const b = rgb[2] / 255;
    const max = Math.max(Math.max(r, g), b);
    const min = Math.min(Math.min(r, g), b);
    const chroma = max - min;
    let grayscale;
    let hue;
    if (chroma < 1) grayscale = min / (1 - chroma);
    else grayscale = 0;
    if (chroma <= 0) hue = 0;
    else if (max === r) hue = (g - b) / chroma % 6;
    else if (max === g) hue = 2 + (b - r) / chroma;
    else hue = 4 + (r - g) / chroma;
    hue /= 6;
    hue %= 1;
    return [
        hue * 360,
        chroma * 100,
        grayscale * 100
    ];
};
$f48c6abccc789ef5$var$convert.hsl.hcg = function(hsl) {
    const s = hsl[1] / 100;
    const l = hsl[2] / 100;
    const c = l < 0.5 ? 2.0 * s * l : 2.0 * s * (1.0 - l);
    let f = 0;
    if (c < 1.0) f = (l - 0.5 * c) / (1.0 - c);
    return [
        hsl[0],
        c * 100,
        f * 100
    ];
};
$f48c6abccc789ef5$var$convert.hsv.hcg = function(hsv) {
    const s = hsv[1] / 100;
    const v = hsv[2] / 100;
    const c = s * v;
    let f = 0;
    if (c < 1.0) f = (v - c) / (1 - c);
    return [
        hsv[0],
        c * 100,
        f * 100
    ];
};
$f48c6abccc789ef5$var$convert.hcg.rgb = function(hcg) {
    const h = hcg[0] / 360;
    const c = hcg[1] / 100;
    const g = hcg[2] / 100;
    if (c === 0.0) return [
        g * 255,
        g * 255,
        g * 255
    ];
    const pure = [
        0,
        0,
        0
    ];
    const hi = h % 1 * 6;
    const v = hi % 1;
    const w = 1 - v;
    let mg = 0;
    /* eslint-disable max-statements-per-line */ switch(Math.floor(hi)){
        case 0:
            pure[0] = 1;
            pure[1] = v;
            pure[2] = 0;
            break;
        case 1:
            pure[0] = w;
            pure[1] = 1;
            pure[2] = 0;
            break;
        case 2:
            pure[0] = 0;
            pure[1] = 1;
            pure[2] = v;
            break;
        case 3:
            pure[0] = 0;
            pure[1] = w;
            pure[2] = 1;
            break;
        case 4:
            pure[0] = v;
            pure[1] = 0;
            pure[2] = 1;
            break;
        default:
            pure[0] = 1;
            pure[1] = 0;
            pure[2] = w;
    }
    /* eslint-enable max-statements-per-line */ mg = (1.0 - c) * g;
    return [
        (c * pure[0] + mg) * 255,
        (c * pure[1] + mg) * 255,
        (c * pure[2] + mg) * 255
    ];
};
$f48c6abccc789ef5$var$convert.hcg.hsv = function(hcg) {
    const c = hcg[1] / 100;
    const g = hcg[2] / 100;
    const v = c + g * (1.0 - c);
    let f = 0;
    if (v > 0.0) f = c / v;
    return [
        hcg[0],
        f * 100,
        v * 100
    ];
};
$f48c6abccc789ef5$var$convert.hcg.hsl = function(hcg) {
    const c = hcg[1] / 100;
    const g = hcg[2] / 100;
    const l = g * (1.0 - c) + 0.5 * c;
    let s = 0;
    if (l > 0.0 && l < 0.5) s = c / (2 * l);
    else if (l >= 0.5 && l < 1.0) s = c / (2 * (1 - l));
    return [
        hcg[0],
        s * 100,
        l * 100
    ];
};
$f48c6abccc789ef5$var$convert.hcg.hwb = function(hcg) {
    const c = hcg[1] / 100;
    const g = hcg[2] / 100;
    const v = c + g * (1.0 - c);
    return [
        hcg[0],
        (v - c) * 100,
        (1 - v) * 100
    ];
};
$f48c6abccc789ef5$var$convert.hwb.hcg = function(hwb) {
    const w = hwb[1] / 100;
    const b = hwb[2] / 100;
    const v = 1 - b;
    const c = v - w;
    let g = 0;
    if (c < 1) g = (v - c) / (1 - c);
    return [
        hwb[0],
        c * 100,
        g * 100
    ];
};
$f48c6abccc789ef5$var$convert.apple.rgb = function(apple) {
    return [
        apple[0] / 65535 * 255,
        apple[1] / 65535 * 255,
        apple[2] / 65535 * 255
    ];
};
$f48c6abccc789ef5$var$convert.rgb.apple = function(rgb) {
    return [
        rgb[0] / 255 * 65535,
        rgb[1] / 255 * 65535,
        rgb[2] / 255 * 65535
    ];
};
$f48c6abccc789ef5$var$convert.gray.rgb = function(args) {
    return [
        args[0] / 100 * 255,
        args[0] / 100 * 255,
        args[0] / 100 * 255
    ];
};
$f48c6abccc789ef5$var$convert.gray.hsl = function(args) {
    return [
        0,
        0,
        args[0]
    ];
};
$f48c6abccc789ef5$var$convert.gray.hsv = $f48c6abccc789ef5$var$convert.gray.hsl;
$f48c6abccc789ef5$var$convert.gray.hwb = function(gray) {
    return [
        0,
        100,
        gray[0]
    ];
};
$f48c6abccc789ef5$var$convert.gray.cmyk = function(gray) {
    return [
        0,
        0,
        0,
        gray[0]
    ];
};
$f48c6abccc789ef5$var$convert.gray.lab = function(gray) {
    return [
        gray[0],
        0,
        0
    ];
};
$f48c6abccc789ef5$var$convert.gray.hex = function(gray) {
    const val = Math.round(gray[0] / 100 * 255) & 0xFF;
    const integer = (val << 16) + (val << 8) + val;
    const string = integer.toString(16).toUpperCase();
    return '000000'.substring(string.length) + string;
};
$f48c6abccc789ef5$var$convert.rgb.gray = function(rgb) {
    const val = (rgb[0] + rgb[1] + rgb[2]) / 3;
    return [
        val / 255 * 100
    ];
};


var $145e82ba91440c6e$exports = {};

/*
	This function routes a model to all other models.

	all functions that are routed have a property `.conversion` attached
	to the returned synthetic function. This property is an array
	of strings, each with the steps in between the 'from' and 'to'
	color models (inclusive).

	conversions that are not possible simply are not included.
*/ function $145e82ba91440c6e$var$buildGraph() {
    const graph = {};
    // https://jsperf.com/object-keys-vs-for-in-with-closure/3
    const models = Object.keys($f48c6abccc789ef5$exports);
    for(let len = models.length, i = 0; i < len; i++)graph[models[i]] = {
        // http://jsperf.com/1-vs-infinity
        // micro-opt, but this is simple.
        distance: -1,
        parent: null
    };
    return graph;
}
// https://en.wikipedia.org/wiki/Breadth-first_search
function $145e82ba91440c6e$var$deriveBFS(fromModel) {
    const graph = $145e82ba91440c6e$var$buildGraph();
    const queue = [
        fromModel
    ]; // Unshift -> queue -> pop
    graph[fromModel].distance = 0;
    while(queue.length){
        const current = queue.pop();
        const adjacents = Object.keys($f48c6abccc789ef5$exports[current]);
        for(let len = adjacents.length, i = 0; i < len; i++){
            const adjacent = adjacents[i];
            const node = graph[adjacent];
            if (node.distance === -1) {
                node.distance = graph[current].distance + 1;
                node.parent = current;
                queue.unshift(adjacent);
            }
        }
    }
    return graph;
}
function $145e82ba91440c6e$var$link(from, to) {
    return function(args) {
        return to(from(args));
    };
}
function $145e82ba91440c6e$var$wrapConversion(toModel, graph) {
    const path = [
        graph[toModel].parent,
        toModel
    ];
    let fn = $f48c6abccc789ef5$exports[graph[toModel].parent][toModel];
    let cur = graph[toModel].parent;
    while(graph[cur].parent){
        path.unshift(graph[cur].parent);
        fn = $145e82ba91440c6e$var$link($f48c6abccc789ef5$exports[graph[cur].parent][cur], fn);
        cur = graph[cur].parent;
    }
    fn.conversion = path;
    return fn;
}
$145e82ba91440c6e$exports = function(fromModel) {
    const graph = $145e82ba91440c6e$var$deriveBFS(fromModel);
    const conversion = {};
    const models = Object.keys(graph);
    for(let len = models.length, i = 0; i < len; i++){
        const toModel = models[i];
        const node = graph[toModel];
        if (node.parent === null) continue;
        conversion[toModel] = $145e82ba91440c6e$var$wrapConversion(toModel, graph);
    }
    return conversion;
};


const $b0ce43f4286ff9ac$var$convert = {};
const $b0ce43f4286ff9ac$var$models = Object.keys($f48c6abccc789ef5$exports);
function $b0ce43f4286ff9ac$var$wrapRaw(fn) {
    const wrappedFn = function(...args) {
        const arg0 = args[0];
        if (arg0 === undefined || arg0 === null) return arg0;
        if (arg0.length > 1) args = arg0;
        return fn(args);
    };
    // Preserve .conversion property if there is one
    if ('conversion' in fn) wrappedFn.conversion = fn.conversion;
    return wrappedFn;
}
function $b0ce43f4286ff9ac$var$wrapRounded(fn) {
    const wrappedFn = function(...args) {
        const arg0 = args[0];
        if (arg0 === undefined || arg0 === null) return arg0;
        if (arg0.length > 1) args = arg0;
        const result = fn(args);
        // We're assuming the result is an array here.
        // see notice in conversions.js; don't use box types
        // in conversion functions.
        if (typeof result === 'object') for(let len = result.length, i = 0; i < len; i++)result[i] = Math.round(result[i]);
        return result;
    };
    // Preserve .conversion property if there is one
    if ('conversion' in fn) wrappedFn.conversion = fn.conversion;
    return wrappedFn;
}
$b0ce43f4286ff9ac$var$models.forEach((fromModel)=>{
    $b0ce43f4286ff9ac$var$convert[fromModel] = {};
    Object.defineProperty($b0ce43f4286ff9ac$var$convert[fromModel], 'channels', {
        value: $f48c6abccc789ef5$exports[fromModel].channels
    });
    Object.defineProperty($b0ce43f4286ff9ac$var$convert[fromModel], 'labels', {
        value: $f48c6abccc789ef5$exports[fromModel].labels
    });
    const routes = $145e82ba91440c6e$exports(fromModel);
    const routeModels = Object.keys(routes);
    routeModels.forEach((toModel)=>{
        const fn = routes[toModel];
        $b0ce43f4286ff9ac$var$convert[fromModel][toModel] = $b0ce43f4286ff9ac$var$wrapRounded(fn);
        $b0ce43f4286ff9ac$var$convert[fromModel][toModel].raw = $b0ce43f4286ff9ac$var$wrapRaw(fn);
    });
});
$b0ce43f4286ff9ac$exports = $b0ce43f4286ff9ac$var$convert;


const $040ba8e83b63f6c5$var$skippedModels = [
    // To be honest, I don't really feel like keyword belongs in color convert, but eh.
    'keyword',
    // Gray conflicts with some method names, and has its own method defined.
    'gray',
    // Shouldn't really be in color-convert either...
    'hex'
];
const $040ba8e83b63f6c5$var$hashedModelKeys = {};
for (const model of Object.keys($b0ce43f4286ff9ac$exports))$040ba8e83b63f6c5$var$hashedModelKeys[[
    ...$b0ce43f4286ff9ac$exports[model].labels
].sort().join('')] = model;
const $040ba8e83b63f6c5$var$limiters = {};
function $040ba8e83b63f6c5$var$Color(object, model) {
    if (!(this instanceof $040ba8e83b63f6c5$var$Color)) return new $040ba8e83b63f6c5$var$Color(object, model);
    if (model && model in $040ba8e83b63f6c5$var$skippedModels) model = null;
    if (model && !(model in $b0ce43f4286ff9ac$exports)) throw new Error('Unknown model: ' + model);
    let i;
    let channels;
    if (object == null) {
        this.model = 'rgb';
        this.color = [
            0,
            0,
            0
        ];
        this.valpha = 1;
    } else if (object instanceof $040ba8e83b63f6c5$var$Color) {
        this.model = object.model;
        this.color = [
            ...object.color
        ];
        this.valpha = object.valpha;
    } else if (typeof object === 'string') {
        const result = $41148a221f03ee13$exports.get(object);
        if (result === null) throw new Error('Unable to parse color from string: ' + object);
        this.model = result.model;
        channels = $b0ce43f4286ff9ac$exports[this.model].channels;
        this.color = result.value.slice(0, channels);
        this.valpha = typeof result.value[channels] === 'number' ? result.value[channels] : 1;
    } else if (object.length > 0) {
        this.model = model || 'rgb';
        channels = $b0ce43f4286ff9ac$exports[this.model].channels;
        const newArray = Array.prototype.slice.call(object, 0, channels);
        this.color = $040ba8e83b63f6c5$var$zeroArray(newArray, channels);
        this.valpha = typeof object[channels] === 'number' ? object[channels] : 1;
    } else if (typeof object === 'number') {
        // This is always RGB - can be converted later on.
        this.model = 'rgb';
        this.color = [
            object >> 16 & 0xFF,
            object >> 8 & 0xFF,
            object & 0xFF
        ];
        this.valpha = 1;
    } else {
        this.valpha = 1;
        const keys = Object.keys(object);
        if ('alpha' in object) {
            keys.splice(keys.indexOf('alpha'), 1);
            this.valpha = typeof object.alpha === 'number' ? object.alpha : 0;
        }
        const hashedKeys = keys.sort().join('');
        if (!(hashedKeys in $040ba8e83b63f6c5$var$hashedModelKeys)) throw new Error('Unable to parse color from object: ' + JSON.stringify(object));
        this.model = $040ba8e83b63f6c5$var$hashedModelKeys[hashedKeys];
        const { labels: labels } = $b0ce43f4286ff9ac$exports[this.model];
        const color = [];
        for(i = 0; i < labels.length; i++)color.push(object[labels[i]]);
        this.color = $040ba8e83b63f6c5$var$zeroArray(color);
    }
    // Perform limitations (clamping, etc.)
    if ($040ba8e83b63f6c5$var$limiters[this.model]) {
        channels = $b0ce43f4286ff9ac$exports[this.model].channels;
        for(i = 0; i < channels; i++){
            const limit = $040ba8e83b63f6c5$var$limiters[this.model][i];
            if (limit) this.color[i] = limit(this.color[i]);
        }
    }
    this.valpha = Math.max(0, Math.min(1, this.valpha));
    if (Object.freeze) Object.freeze(this);
}
$040ba8e83b63f6c5$var$Color.prototype = {
    toString () {
        return this.string();
    },
    toJSON () {
        return this[this.model]();
    },
    string (places) {
        let self = this.model in $41148a221f03ee13$exports.to ? this : this.rgb();
        self = self.round(typeof places === 'number' ? places : 1);
        const args = self.valpha === 1 ? self.color : [
            ...self.color,
            this.valpha
        ];
        return $41148a221f03ee13$exports.to[self.model](args);
    },
    percentString (places) {
        const self = this.rgb().round(typeof places === 'number' ? places : 1);
        const args = self.valpha === 1 ? self.color : [
            ...self.color,
            this.valpha
        ];
        return $41148a221f03ee13$exports.to.rgb.percent(args);
    },
    array () {
        return this.valpha === 1 ? [
            ...this.color
        ] : [
            ...this.color,
            this.valpha
        ];
    },
    object () {
        const result = {};
        const { channels: channels } = $b0ce43f4286ff9ac$exports[this.model];
        const { labels: labels } = $b0ce43f4286ff9ac$exports[this.model];
        for(let i = 0; i < channels; i++)result[labels[i]] = this.color[i];
        if (this.valpha !== 1) result.alpha = this.valpha;
        return result;
    },
    unitArray () {
        const rgb = this.rgb().color;
        rgb[0] /= 255;
        rgb[1] /= 255;
        rgb[2] /= 255;
        if (this.valpha !== 1) rgb.push(this.valpha);
        return rgb;
    },
    unitObject () {
        const rgb = this.rgb().object();
        rgb.r /= 255;
        rgb.g /= 255;
        rgb.b /= 255;
        if (this.valpha !== 1) rgb.alpha = this.valpha;
        return rgb;
    },
    round (places) {
        places = Math.max(places || 0, 0);
        return new $040ba8e83b63f6c5$var$Color([
            ...this.color.map($040ba8e83b63f6c5$var$roundToPlace(places)),
            this.valpha
        ], this.model);
    },
    alpha (value) {
        if (value !== undefined) return new $040ba8e83b63f6c5$var$Color([
            ...this.color,
            Math.max(0, Math.min(1, value))
        ], this.model);
        return this.valpha;
    },
    // Rgb
    red: $040ba8e83b63f6c5$var$getset('rgb', 0, $040ba8e83b63f6c5$var$maxfn(255)),
    green: $040ba8e83b63f6c5$var$getset('rgb', 1, $040ba8e83b63f6c5$var$maxfn(255)),
    blue: $040ba8e83b63f6c5$var$getset('rgb', 2, $040ba8e83b63f6c5$var$maxfn(255)),
    hue: $040ba8e83b63f6c5$var$getset([
        'hsl',
        'hsv',
        'hsl',
        'hwb',
        'hcg'
    ], 0, (value)=>(value % 360 + 360) % 360),
    saturationl: $040ba8e83b63f6c5$var$getset('hsl', 1, $040ba8e83b63f6c5$var$maxfn(100)),
    lightness: $040ba8e83b63f6c5$var$getset('hsl', 2, $040ba8e83b63f6c5$var$maxfn(100)),
    saturationv: $040ba8e83b63f6c5$var$getset('hsv', 1, $040ba8e83b63f6c5$var$maxfn(100)),
    value: $040ba8e83b63f6c5$var$getset('hsv', 2, $040ba8e83b63f6c5$var$maxfn(100)),
    chroma: $040ba8e83b63f6c5$var$getset('hcg', 1, $040ba8e83b63f6c5$var$maxfn(100)),
    gray: $040ba8e83b63f6c5$var$getset('hcg', 2, $040ba8e83b63f6c5$var$maxfn(100)),
    white: $040ba8e83b63f6c5$var$getset('hwb', 1, $040ba8e83b63f6c5$var$maxfn(100)),
    wblack: $040ba8e83b63f6c5$var$getset('hwb', 2, $040ba8e83b63f6c5$var$maxfn(100)),
    cyan: $040ba8e83b63f6c5$var$getset('cmyk', 0, $040ba8e83b63f6c5$var$maxfn(100)),
    magenta: $040ba8e83b63f6c5$var$getset('cmyk', 1, $040ba8e83b63f6c5$var$maxfn(100)),
    yellow: $040ba8e83b63f6c5$var$getset('cmyk', 2, $040ba8e83b63f6c5$var$maxfn(100)),
    black: $040ba8e83b63f6c5$var$getset('cmyk', 3, $040ba8e83b63f6c5$var$maxfn(100)),
    x: $040ba8e83b63f6c5$var$getset('xyz', 0, $040ba8e83b63f6c5$var$maxfn(95.047)),
    y: $040ba8e83b63f6c5$var$getset('xyz', 1, $040ba8e83b63f6c5$var$maxfn(100)),
    z: $040ba8e83b63f6c5$var$getset('xyz', 2, $040ba8e83b63f6c5$var$maxfn(108.833)),
    l: $040ba8e83b63f6c5$var$getset('lab', 0, $040ba8e83b63f6c5$var$maxfn(100)),
    a: $040ba8e83b63f6c5$var$getset('lab', 1),
    b: $040ba8e83b63f6c5$var$getset('lab', 2),
    keyword (value) {
        if (value !== undefined) return new $040ba8e83b63f6c5$var$Color(value);
        return $b0ce43f4286ff9ac$exports[this.model].keyword(this.color);
    },
    hex (value) {
        if (value !== undefined) return new $040ba8e83b63f6c5$var$Color(value);
        return $41148a221f03ee13$exports.to.hex(this.rgb().round().color);
    },
    hexa (value) {
        if (value !== undefined) return new $040ba8e83b63f6c5$var$Color(value);
        const rgbArray = this.rgb().round().color;
        let alphaHex = Math.round(this.valpha * 255).toString(16).toUpperCase();
        if (alphaHex.length === 1) alphaHex = '0' + alphaHex;
        return $41148a221f03ee13$exports.to.hex(rgbArray) + alphaHex;
    },
    rgbNumber () {
        const rgb = this.rgb().color;
        return (rgb[0] & 0xFF) << 16 | (rgb[1] & 0xFF) << 8 | rgb[2] & 0xFF;
    },
    luminosity () {
        // http://www.w3.org/TR/WCAG20/#relativeluminancedef
        const rgb = this.rgb().color;
        const lum = [];
        for (const [i, element] of rgb.entries()){
            const chan = element / 255;
            lum[i] = chan <= 0.04045 ? chan / 12.92 : ((chan + 0.055) / 1.055) ** 2.4;
        }
        return 0.2126 * lum[0] + 0.7152 * lum[1] + 0.0722 * lum[2];
    },
    contrast (color2) {
        // http://www.w3.org/TR/WCAG20/#contrast-ratiodef
        const lum1 = this.luminosity();
        const lum2 = color2.luminosity();
        if (lum1 > lum2) return (lum1 + 0.05) / (lum2 + 0.05);
        return (lum2 + 0.05) / (lum1 + 0.05);
    },
    level (color2) {
        // https://www.w3.org/TR/WCAG/#contrast-enhanced
        const contrastRatio = this.contrast(color2);
        if (contrastRatio >= 7) return 'AAA';
        return contrastRatio >= 4.5 ? 'AA' : '';
    },
    isDark () {
        // YIQ equation from http://24ways.org/2010/calculating-color-contrast
        const rgb = this.rgb().color;
        const yiq = (rgb[0] * 2126 + rgb[1] * 7152 + rgb[2] * 722) / 10000;
        return yiq < 128;
    },
    isLight () {
        return !this.isDark();
    },
    negate () {
        const rgb = this.rgb();
        for(let i = 0; i < 3; i++)rgb.color[i] = 255 - rgb.color[i];
        return rgb;
    },
    lighten (ratio) {
        const hsl = this.hsl();
        hsl.color[2] += hsl.color[2] * ratio;
        return hsl;
    },
    darken (ratio) {
        const hsl = this.hsl();
        hsl.color[2] -= hsl.color[2] * ratio;
        return hsl;
    },
    saturate (ratio) {
        const hsl = this.hsl();
        hsl.color[1] += hsl.color[1] * ratio;
        return hsl;
    },
    desaturate (ratio) {
        const hsl = this.hsl();
        hsl.color[1] -= hsl.color[1] * ratio;
        return hsl;
    },
    whiten (ratio) {
        const hwb = this.hwb();
        hwb.color[1] += hwb.color[1] * ratio;
        return hwb;
    },
    blacken (ratio) {
        const hwb = this.hwb();
        hwb.color[2] += hwb.color[2] * ratio;
        return hwb;
    },
    grayscale () {
        // http://en.wikipedia.org/wiki/Grayscale#Converting_color_to_grayscale
        const rgb = this.rgb().color;
        const value = rgb[0] * 0.3 + rgb[1] * 0.59 + rgb[2] * 0.11;
        return $040ba8e83b63f6c5$var$Color.rgb(value, value, value);
    },
    fade (ratio) {
        return this.alpha(this.valpha - this.valpha * ratio);
    },
    opaquer (ratio) {
        return this.alpha(this.valpha + this.valpha * ratio);
    },
    rotate (degrees) {
        const hsl = this.hsl();
        let hue = hsl.color[0];
        hue = (hue + degrees) % 360;
        hue = hue < 0 ? 360 + hue : hue;
        hsl.color[0] = hue;
        return hsl;
    },
    mix (mixinColor, weight) {
        // Ported from sass implementation in C
        // https://github.com/sass/libsass/blob/0e6b4a2850092356aa3ece07c6b249f0221caced/functions.cpp#L209
        if (!mixinColor || !mixinColor.rgb) throw new Error('Argument to "mix" was not a Color instance, but rather an instance of ' + typeof mixinColor);
        const color1 = mixinColor.rgb();
        const color2 = this.rgb();
        const p = weight === undefined ? 0.5 : weight;
        const w = 2 * p - 1;
        const a = color1.alpha() - color2.alpha();
        const w1 = ((w * a === -1 ? w : (w + a) / (1 + w * a)) + 1) / 2;
        const w2 = 1 - w1;
        return $040ba8e83b63f6c5$var$Color.rgb(w1 * color1.red() + w2 * color2.red(), w1 * color1.green() + w2 * color2.green(), w1 * color1.blue() + w2 * color2.blue(), color1.alpha() * p + color2.alpha() * (1 - p));
    }
};
// Model conversion methods and static constructors
for (const model of Object.keys($b0ce43f4286ff9ac$exports)){
    if ($040ba8e83b63f6c5$var$skippedModels.includes(model)) continue;
    const { channels: channels } = $b0ce43f4286ff9ac$exports[model];
    // Conversion methods
    $040ba8e83b63f6c5$var$Color.prototype[model] = function(...args) {
        if (this.model === model) return new $040ba8e83b63f6c5$var$Color(this);
        if (args.length > 0) return new $040ba8e83b63f6c5$var$Color(args, model);
        return new $040ba8e83b63f6c5$var$Color([
            ...$040ba8e83b63f6c5$var$assertArray($b0ce43f4286ff9ac$exports[this.model][model].raw(this.color)),
            this.valpha
        ], model);
    };
    // 'static' construction methods
    $040ba8e83b63f6c5$var$Color[model] = function(...args) {
        let color = args[0];
        if (typeof color === 'number') color = $040ba8e83b63f6c5$var$zeroArray(args, channels);
        return new $040ba8e83b63f6c5$var$Color(color, model);
    };
}
function $040ba8e83b63f6c5$var$roundTo(number, places) {
    return Number(number.toFixed(places));
}
function $040ba8e83b63f6c5$var$roundToPlace(places) {
    return function(number) {
        return $040ba8e83b63f6c5$var$roundTo(number, places);
    };
}
function $040ba8e83b63f6c5$var$getset(model, channel, modifier) {
    model = Array.isArray(model) ? model : [
        model
    ];
    for (const m of model)($040ba8e83b63f6c5$var$limiters[m] || ($040ba8e83b63f6c5$var$limiters[m] = []))[channel] = modifier;
    model = model[0];
    return function(value) {
        let result;
        if (value !== undefined) {
            if (modifier) value = modifier(value);
            result = this[model]();
            result.color[channel] = value;
            return result;
        }
        result = this[model]().color[channel];
        if (modifier) result = modifier(result);
        return result;
    };
}
function $040ba8e83b63f6c5$var$maxfn(max) {
    return function(v) {
        return Math.max(0, Math.min(max, v));
    };
}
function $040ba8e83b63f6c5$var$assertArray(value) {
    return Array.isArray(value) ? value : [
        value
    ];
}
function $040ba8e83b63f6c5$var$zeroArray(array, length) {
    for(let i = 0; i < length; i++)if (typeof array[i] !== 'number') array[i] = 0;
    return array;
}
$040ba8e83b63f6c5$exports = $040ba8e83b63f6c5$var$Color;


var $84df4a417d8748cd$exports = {};
var $a041c22fe0ffdbd8$exports = {};
"use strict";
Object.defineProperty($a041c22fe0ffdbd8$exports, "__esModule", {
    value: true
});
Object.defineProperty($a041c22fe0ffdbd8$exports, "default", {
    enumerable: true,
    get: function() {
        return $a041c22fe0ffdbd8$var$_default;
    }
});

const $a041c22fe0ffdbd8$var$_createPlugin = /*#__PURE__*/ $a041c22fe0ffdbd8$var$_interop_require_default((parcelRequire("2Gj6O")));
function $a041c22fe0ffdbd8$var$_interop_require_default(obj) {
    return obj && obj.__esModule ? obj : {
        default: obj
    };
}
const $a041c22fe0ffdbd8$var$_default = $a041c22fe0ffdbd8$var$_createPlugin.default;


$84df4a417d8748cd$exports = ($a041c22fe0ffdbd8$exports.__esModule ? $a041c22fe0ffdbd8$exports : {
    default: $a041c22fe0ffdbd8$exports
}).default;


var $ff1d430d0589e4d1$exports = {};
/**
 * lodash (Custom Build) <https://lodash.com/>
 * Build: `lodash modularize exports="npm" -o ./`
 * Copyright jQuery Foundation and other contributors <https://jquery.org/>
 * Released under MIT license <https://lodash.com/license>
 * Based on Underscore.js 1.8.3 <http://underscorejs.org/LICENSE>
 * Copyright Jeremy Ashkenas, DocumentCloud and Investigative Reporters & Editors
 */ /** Used as references for various `Number` constants. */ var $ff1d430d0589e4d1$var$MAX_SAFE_INTEGER = 9007199254740991;
/** `Object#toString` result references. */ var $ff1d430d0589e4d1$var$argsTag = '[object Arguments]', $ff1d430d0589e4d1$var$funcTag = '[object Function]', $ff1d430d0589e4d1$var$genTag = '[object GeneratorFunction]';
/** Used to detect unsigned integer values. */ var $ff1d430d0589e4d1$var$reIsUint = /^(?:0|[1-9]\d*)$/;
/**
 * A specialized version of `_.forEach` for arrays without support for
 * iteratee shorthands.
 *
 * @private
 * @param {Array} [array] The array to iterate over.
 * @param {Function} iteratee The function invoked per iteration.
 * @returns {Array} Returns `array`.
 */ function $ff1d430d0589e4d1$var$arrayEach(array, iteratee) {
    var index = -1, length = array ? array.length : 0;
    while(++index < length){
        if (iteratee(array[index], index, array) === false) break;
    }
    return array;
}
/**
 * The base implementation of `_.times` without support for iteratee shorthands
 * or max array length checks.
 *
 * @private
 * @param {number} n The number of times to invoke `iteratee`.
 * @param {Function} iteratee The function invoked per iteration.
 * @returns {Array} Returns the array of results.
 */ function $ff1d430d0589e4d1$var$baseTimes(n, iteratee) {
    var index = -1, result = Array(n);
    while(++index < n)result[index] = iteratee(index);
    return result;
}
/**
 * Creates a unary function that invokes `func` with its argument transformed.
 *
 * @private
 * @param {Function} func The function to wrap.
 * @param {Function} transform The argument transform.
 * @returns {Function} Returns the new function.
 */ function $ff1d430d0589e4d1$var$overArg(func, transform) {
    return function(arg) {
        return func(transform(arg));
    };
}
/** Used for built-in method references. */ var $ff1d430d0589e4d1$var$objectProto = Object.prototype;
/** Used to check objects for own properties. */ var $ff1d430d0589e4d1$var$hasOwnProperty = $ff1d430d0589e4d1$var$objectProto.hasOwnProperty;
/**
 * Used to resolve the
 * [`toStringTag`](http://ecma-international.org/ecma-262/7.0/#sec-object.prototype.tostring)
 * of values.
 */ var $ff1d430d0589e4d1$var$objectToString = $ff1d430d0589e4d1$var$objectProto.toString;
/** Built-in value references. */ var $ff1d430d0589e4d1$var$propertyIsEnumerable = $ff1d430d0589e4d1$var$objectProto.propertyIsEnumerable;
/* Built-in method references for those with the same name as other `lodash` methods. */ var $ff1d430d0589e4d1$var$nativeKeys = $ff1d430d0589e4d1$var$overArg(Object.keys, Object);
/**
 * Creates an array of the enumerable property names of the array-like `value`.
 *
 * @private
 * @param {*} value The value to query.
 * @param {boolean} inherited Specify returning inherited property names.
 * @returns {Array} Returns the array of property names.
 */ function $ff1d430d0589e4d1$var$arrayLikeKeys(value, inherited) {
    // Safari 8.1 makes `arguments.callee` enumerable in strict mode.
    // Safari 9 makes `arguments.length` enumerable in strict mode.
    var result = $ff1d430d0589e4d1$var$isArray(value) || $ff1d430d0589e4d1$var$isArguments(value) ? $ff1d430d0589e4d1$var$baseTimes(value.length, String) : [];
    var length = result.length, skipIndexes = !!length;
    for(var key in value)if ((inherited || $ff1d430d0589e4d1$var$hasOwnProperty.call(value, key)) && !(skipIndexes && (key == 'length' || $ff1d430d0589e4d1$var$isIndex(key, length)))) result.push(key);
    return result;
}
/**
 * The base implementation of `_.forEach` without support for iteratee shorthands.
 *
 * @private
 * @param {Array|Object} collection The collection to iterate over.
 * @param {Function} iteratee The function invoked per iteration.
 * @returns {Array|Object} Returns `collection`.
 */ var $ff1d430d0589e4d1$var$baseEach = $ff1d430d0589e4d1$var$createBaseEach($ff1d430d0589e4d1$var$baseForOwn);
/**
 * The base implementation of `baseForOwn` which iterates over `object`
 * properties returned by `keysFunc` and invokes `iteratee` for each property.
 * Iteratee functions may exit iteration early by explicitly returning `false`.
 *
 * @private
 * @param {Object} object The object to iterate over.
 * @param {Function} iteratee The function invoked per iteration.
 * @param {Function} keysFunc The function to get the keys of `object`.
 * @returns {Object} Returns `object`.
 */ var $ff1d430d0589e4d1$var$baseFor = $ff1d430d0589e4d1$var$createBaseFor();
/**
 * The base implementation of `_.forOwn` without support for iteratee shorthands.
 *
 * @private
 * @param {Object} object The object to iterate over.
 * @param {Function} iteratee The function invoked per iteration.
 * @returns {Object} Returns `object`.
 */ function $ff1d430d0589e4d1$var$baseForOwn(object, iteratee) {
    return object && $ff1d430d0589e4d1$var$baseFor(object, iteratee, $ff1d430d0589e4d1$var$keys);
}
/**
 * The base implementation of `_.keys` which doesn't treat sparse arrays as dense.
 *
 * @private
 * @param {Object} object The object to query.
 * @returns {Array} Returns the array of property names.
 */ function $ff1d430d0589e4d1$var$baseKeys(object) {
    if (!$ff1d430d0589e4d1$var$isPrototype(object)) return $ff1d430d0589e4d1$var$nativeKeys(object);
    var result = [];
    for(var key in Object(object))if ($ff1d430d0589e4d1$var$hasOwnProperty.call(object, key) && key != 'constructor') result.push(key);
    return result;
}
/**
 * Creates a `baseEach` or `baseEachRight` function.
 *
 * @private
 * @param {Function} eachFunc The function to iterate over a collection.
 * @param {boolean} [fromRight] Specify iterating from right to left.
 * @returns {Function} Returns the new base function.
 */ function $ff1d430d0589e4d1$var$createBaseEach(eachFunc, fromRight) {
    return function(collection, iteratee) {
        if (collection == null) return collection;
        if (!$ff1d430d0589e4d1$var$isArrayLike(collection)) return eachFunc(collection, iteratee);
        var length = collection.length, index = fromRight ? length : -1, iterable = Object(collection);
        while(fromRight ? index-- : ++index < length){
            if (iteratee(iterable[index], index, iterable) === false) break;
        }
        return collection;
    };
}
/**
 * Creates a base function for methods like `_.forIn` and `_.forOwn`.
 *
 * @private
 * @param {boolean} [fromRight] Specify iterating from right to left.
 * @returns {Function} Returns the new base function.
 */ function $ff1d430d0589e4d1$var$createBaseFor(fromRight) {
    return function(object, iteratee, keysFunc) {
        var index = -1, iterable = Object(object), props = keysFunc(object), length = props.length;
        while(length--){
            var key = props[fromRight ? length : ++index];
            if (iteratee(iterable[key], key, iterable) === false) break;
        }
        return object;
    };
}
/**
 * Checks if `value` is a valid array-like index.
 *
 * @private
 * @param {*} value The value to check.
 * @param {number} [length=MAX_SAFE_INTEGER] The upper bounds of a valid index.
 * @returns {boolean} Returns `true` if `value` is a valid index, else `false`.
 */ function $ff1d430d0589e4d1$var$isIndex(value, length) {
    length = length == null ? $ff1d430d0589e4d1$var$MAX_SAFE_INTEGER : length;
    return !!length && (typeof value == 'number' || $ff1d430d0589e4d1$var$reIsUint.test(value)) && value > -1 && value % 1 == 0 && value < length;
}
/**
 * Checks if `value` is likely a prototype object.
 *
 * @private
 * @param {*} value The value to check.
 * @returns {boolean} Returns `true` if `value` is a prototype, else `false`.
 */ function $ff1d430d0589e4d1$var$isPrototype(value) {
    var Ctor = value && value.constructor, proto = typeof Ctor == 'function' && Ctor.prototype || $ff1d430d0589e4d1$var$objectProto;
    return value === proto;
}
/**
 * Iterates over elements of `collection` and invokes `iteratee` for each element.
 * The iteratee is invoked with three arguments: (value, index|key, collection).
 * Iteratee functions may exit iteration early by explicitly returning `false`.
 *
 * **Note:** As with other "Collections" methods, objects with a "length"
 * property are iterated like arrays. To avoid this behavior use `_.forIn`
 * or `_.forOwn` for object iteration.
 *
 * @static
 * @memberOf _
 * @since 0.1.0
 * @alias each
 * @category Collection
 * @param {Array|Object} collection The collection to iterate over.
 * @param {Function} [iteratee=_.identity] The function invoked per iteration.
 * @returns {Array|Object} Returns `collection`.
 * @see _.forEachRight
 * @example
 *
 * _([1, 2]).forEach(function(value) {
 *   console.log(value);
 * });
 * // => Logs `1` then `2`.
 *
 * _.forEach({ 'a': 1, 'b': 2 }, function(value, key) {
 *   console.log(key);
 * });
 * // => Logs 'a' then 'b' (iteration order is not guaranteed).
 */ function $ff1d430d0589e4d1$var$forEach(collection, iteratee) {
    var func = $ff1d430d0589e4d1$var$isArray(collection) ? $ff1d430d0589e4d1$var$arrayEach : $ff1d430d0589e4d1$var$baseEach;
    return func(collection, typeof iteratee == 'function' ? iteratee : $ff1d430d0589e4d1$var$identity);
}
/**
 * Checks if `value` is likely an `arguments` object.
 *
 * @static
 * @memberOf _
 * @since 0.1.0
 * @category Lang
 * @param {*} value The value to check.
 * @returns {boolean} Returns `true` if `value` is an `arguments` object,
 *  else `false`.
 * @example
 *
 * _.isArguments(function() { return arguments; }());
 * // => true
 *
 * _.isArguments([1, 2, 3]);
 * // => false
 */ function $ff1d430d0589e4d1$var$isArguments(value) {
    // Safari 8.1 makes `arguments.callee` enumerable in strict mode.
    return $ff1d430d0589e4d1$var$isArrayLikeObject(value) && $ff1d430d0589e4d1$var$hasOwnProperty.call(value, 'callee') && (!$ff1d430d0589e4d1$var$propertyIsEnumerable.call(value, 'callee') || $ff1d430d0589e4d1$var$objectToString.call(value) == $ff1d430d0589e4d1$var$argsTag);
}
/**
 * Checks if `value` is classified as an `Array` object.
 *
 * @static
 * @memberOf _
 * @since 0.1.0
 * @category Lang
 * @param {*} value The value to check.
 * @returns {boolean} Returns `true` if `value` is an array, else `false`.
 * @example
 *
 * _.isArray([1, 2, 3]);
 * // => true
 *
 * _.isArray(document.body.children);
 * // => false
 *
 * _.isArray('abc');
 * // => false
 *
 * _.isArray(_.noop);
 * // => false
 */ var $ff1d430d0589e4d1$var$isArray = Array.isArray;
/**
 * Checks if `value` is array-like. A value is considered array-like if it's
 * not a function and has a `value.length` that's an integer greater than or
 * equal to `0` and less than or equal to `Number.MAX_SAFE_INTEGER`.
 *
 * @static
 * @memberOf _
 * @since 4.0.0
 * @category Lang
 * @param {*} value The value to check.
 * @returns {boolean} Returns `true` if `value` is array-like, else `false`.
 * @example
 *
 * _.isArrayLike([1, 2, 3]);
 * // => true
 *
 * _.isArrayLike(document.body.children);
 * // => true
 *
 * _.isArrayLike('abc');
 * // => true
 *
 * _.isArrayLike(_.noop);
 * // => false
 */ function $ff1d430d0589e4d1$var$isArrayLike(value) {
    return value != null && $ff1d430d0589e4d1$var$isLength(value.length) && !$ff1d430d0589e4d1$var$isFunction(value);
}
/**
 * This method is like `_.isArrayLike` except that it also checks if `value`
 * is an object.
 *
 * @static
 * @memberOf _
 * @since 4.0.0
 * @category Lang
 * @param {*} value The value to check.
 * @returns {boolean} Returns `true` if `value` is an array-like object,
 *  else `false`.
 * @example
 *
 * _.isArrayLikeObject([1, 2, 3]);
 * // => true
 *
 * _.isArrayLikeObject(document.body.children);
 * // => true
 *
 * _.isArrayLikeObject('abc');
 * // => false
 *
 * _.isArrayLikeObject(_.noop);
 * // => false
 */ function $ff1d430d0589e4d1$var$isArrayLikeObject(value) {
    return $ff1d430d0589e4d1$var$isObjectLike(value) && $ff1d430d0589e4d1$var$isArrayLike(value);
}
/**
 * Checks if `value` is classified as a `Function` object.
 *
 * @static
 * @memberOf _
 * @since 0.1.0
 * @category Lang
 * @param {*} value The value to check.
 * @returns {boolean} Returns `true` if `value` is a function, else `false`.
 * @example
 *
 * _.isFunction(_);
 * // => true
 *
 * _.isFunction(/abc/);
 * // => false
 */ function $ff1d430d0589e4d1$var$isFunction(value) {
    // The use of `Object#toString` avoids issues with the `typeof` operator
    // in Safari 8-9 which returns 'object' for typed array and other constructors.
    var tag = $ff1d430d0589e4d1$var$isObject(value) ? $ff1d430d0589e4d1$var$objectToString.call(value) : '';
    return tag == $ff1d430d0589e4d1$var$funcTag || tag == $ff1d430d0589e4d1$var$genTag;
}
/**
 * Checks if `value` is a valid array-like length.
 *
 * **Note:** This method is loosely based on
 * [`ToLength`](http://ecma-international.org/ecma-262/7.0/#sec-tolength).
 *
 * @static
 * @memberOf _
 * @since 4.0.0
 * @category Lang
 * @param {*} value The value to check.
 * @returns {boolean} Returns `true` if `value` is a valid length, else `false`.
 * @example
 *
 * _.isLength(3);
 * // => true
 *
 * _.isLength(Number.MIN_VALUE);
 * // => false
 *
 * _.isLength(Infinity);
 * // => false
 *
 * _.isLength('3');
 * // => false
 */ function $ff1d430d0589e4d1$var$isLength(value) {
    return typeof value == 'number' && value > -1 && value % 1 == 0 && value <= $ff1d430d0589e4d1$var$MAX_SAFE_INTEGER;
}
/**
 * Checks if `value` is the
 * [language type](http://www.ecma-international.org/ecma-262/7.0/#sec-ecmascript-language-types)
 * of `Object`. (e.g. arrays, functions, objects, regexes, `new Number(0)`, and `new String('')`)
 *
 * @static
 * @memberOf _
 * @since 0.1.0
 * @category Lang
 * @param {*} value The value to check.
 * @returns {boolean} Returns `true` if `value` is an object, else `false`.
 * @example
 *
 * _.isObject({});
 * // => true
 *
 * _.isObject([1, 2, 3]);
 * // => true
 *
 * _.isObject(_.noop);
 * // => true
 *
 * _.isObject(null);
 * // => false
 */ function $ff1d430d0589e4d1$var$isObject(value) {
    var type = typeof value;
    return !!value && (type == 'object' || type == 'function');
}
/**
 * Checks if `value` is object-like. A value is object-like if it's not `null`
 * and has a `typeof` result of "object".
 *
 * @static
 * @memberOf _
 * @since 4.0.0
 * @category Lang
 * @param {*} value The value to check.
 * @returns {boolean} Returns `true` if `value` is object-like, else `false`.
 * @example
 *
 * _.isObjectLike({});
 * // => true
 *
 * _.isObjectLike([1, 2, 3]);
 * // => true
 *
 * _.isObjectLike(_.noop);
 * // => false
 *
 * _.isObjectLike(null);
 * // => false
 */ function $ff1d430d0589e4d1$var$isObjectLike(value) {
    return !!value && typeof value == 'object';
}
/**
 * Creates an array of the own enumerable property names of `object`.
 *
 * **Note:** Non-object values are coerced to objects. See the
 * [ES spec](http://ecma-international.org/ecma-262/7.0/#sec-object.keys)
 * for more details.
 *
 * @static
 * @since 0.1.0
 * @memberOf _
 * @category Object
 * @param {Object} object The object to query.
 * @returns {Array} Returns the array of property names.
 * @example
 *
 * function Foo() {
 *   this.a = 1;
 *   this.b = 2;
 * }
 *
 * Foo.prototype.c = 3;
 *
 * _.keys(new Foo);
 * // => ['a', 'b'] (iteration order is not guaranteed)
 *
 * _.keys('hi');
 * // => ['0', '1']
 */ function $ff1d430d0589e4d1$var$keys(object) {
    return $ff1d430d0589e4d1$var$isArrayLike(object) ? $ff1d430d0589e4d1$var$arrayLikeKeys(object) : $ff1d430d0589e4d1$var$baseKeys(object);
}
/**
 * This method returns the first argument it receives.
 *
 * @static
 * @since 0.1.0
 * @memberOf _
 * @category Util
 * @param {*} value Any value.
 * @returns {*} Returns `value`.
 * @example
 *
 * var object = { 'a': 1 };
 *
 * console.log(_.identity(object) === object);
 * // => true
 */ function $ff1d430d0589e4d1$var$identity(value) {
    return value;
}
$ff1d430d0589e4d1$exports = $ff1d430d0589e4d1$var$forEach;


var $4c07cb87e79a40c4$exports = {};
$4c07cb87e79a40c4$exports = $4c07cb87e79a40c4$var$flatten;
$4c07cb87e79a40c4$var$flatten.flatten = $4c07cb87e79a40c4$var$flatten;
$4c07cb87e79a40c4$var$flatten.unflatten = $4c07cb87e79a40c4$var$unflatten;
function $4c07cb87e79a40c4$var$isBuffer(obj) {
    return obj && obj.constructor && typeof obj.constructor.isBuffer === 'function' && obj.constructor.isBuffer(obj);
}
function $4c07cb87e79a40c4$var$keyIdentity(key) {
    return key;
}
function $4c07cb87e79a40c4$var$flatten(target, opts) {
    opts = opts || {};
    const delimiter = opts.delimiter || '.';
    const maxDepth = opts.maxDepth;
    const transformKey = opts.transformKey || $4c07cb87e79a40c4$var$keyIdentity;
    const output = {};
    function step(object, prev, currentDepth) {
        currentDepth = currentDepth || 1;
        Object.keys(object).forEach(function(key) {
            const value = object[key];
            const isarray = opts.safe && Array.isArray(value);
            const type = Object.prototype.toString.call(value);
            const isbuffer = $4c07cb87e79a40c4$var$isBuffer(value);
            const isobject = type === '[object Object]' || type === '[object Array]';
            const newKey = prev ? prev + delimiter + transformKey(key) : transformKey(key);
            if (!isarray && !isbuffer && isobject && Object.keys(value).length && (!opts.maxDepth || currentDepth < maxDepth)) return step(value, newKey, currentDepth + 1);
            output[newKey] = value;
        });
    }
    step(target);
    return output;
}
function $4c07cb87e79a40c4$var$unflatten(target, opts) {
    opts = opts || {};
    const delimiter = opts.delimiter || '.';
    const overwrite = opts.overwrite || false;
    const transformKey = opts.transformKey || $4c07cb87e79a40c4$var$keyIdentity;
    const result = {};
    const isbuffer = $4c07cb87e79a40c4$var$isBuffer(target);
    if (isbuffer || Object.prototype.toString.call(target) !== '[object Object]') return target;
    // safely ensure that the key is
    // an integer.
    function getkey(key) {
        const parsedKey = Number(key);
        return isNaN(parsedKey) || key.indexOf('.') !== -1 || opts.object ? key : parsedKey;
    }
    function addKeys(keyPrefix, recipient, target) {
        return Object.keys(target).reduce(function(result, key) {
            result[keyPrefix + delimiter + key] = target[key];
            return result;
        }, recipient);
    }
    function isEmpty(val) {
        const type = Object.prototype.toString.call(val);
        const isArray = type === '[object Array]';
        const isObject = type === '[object Object]';
        if (!val) return true;
        else if (isArray) return !val.length;
        else if (isObject) return !Object.keys(val).length;
    }
    target = Object.keys(target).reduce(function(result, key) {
        const type = Object.prototype.toString.call(target[key]);
        const isObject = type === '[object Object]' || type === '[object Array]';
        if (!isObject || isEmpty(target[key])) {
            result[key] = target[key];
            return result;
        } else return addKeys(key, result, $4c07cb87e79a40c4$var$flatten(target[key], opts));
    }, {});
    Object.keys(target).forEach(function(key) {
        const split = key.split(delimiter).map(transformKey);
        let key1 = getkey(split.shift());
        let key2 = getkey(split[0]);
        let recipient = result;
        while(key2 !== undefined){
            if (key1 === '__proto__') return;
            const type = Object.prototype.toString.call(recipient[key1]);
            const isobject = type === '[object Object]' || type === '[object Array]';
            // do not write over falsey, non-undefined values if overwrite is false
            if (!overwrite && !isobject && typeof recipient[key1] !== 'undefined') return;
            if (overwrite && !isobject || !overwrite && recipient[key1] == null) recipient[key1] = typeof key2 === 'number' && !opts.object ? [] : {};
            recipient = recipient[key1];
            if (split.length > 0) {
                key1 = getkey(split.shift());
                key2 = getkey(split[0]);
            }
        }
        // unflatten again for 'messy objects'
        recipient[key1] = $4c07cb87e79a40c4$var$unflatten(target[key], opts);
    });
    return result;
}


var $1cab5c01d65a66cd$var$SCHEME = Symbol("color-scheme");
var $1cab5c01d65a66cd$var$emptyConfig = {};
var $1cab5c01d65a66cd$export$231aedbc8f7c4964 = (config = $1cab5c01d65a66cd$var$emptyConfig, { produceCssVariable: produceCssVariable = $1cab5c01d65a66cd$var$defaultProduceCssVariable, produceThemeClass: produceThemeClass = $1cab5c01d65a66cd$var$defaultProduceThemeClass, produceThemeVariant: produceThemeVariant = produceThemeClass, defaultTheme: defaultTheme, strict: strict = false } = {})=>{
    const resolved = {
        variants: [],
        utilities: {},
        colors: {}
    };
    const configObject = typeof config === "function" ? config({
        dark: $1cab5c01d65a66cd$var$dark,
        light: $1cab5c01d65a66cd$var$light
    }) : config;
    (0, (/*@__PURE__*/$parcel$interopDefault($ff1d430d0589e4d1$exports)))(configObject, (colors, themeName)=>{
        const themeClassName = produceThemeClass(themeName);
        const themeVariant = produceThemeVariant(themeName);
        const flatColors = $1cab5c01d65a66cd$var$flattenColors(colors);
        resolved.variants.push({
            name: themeVariant,
            // tailwind will generate only the first matched definition
            definition: [
                $1cab5c01d65a66cd$var$generateVariantDefinitions(`.${themeClassName}`),
                $1cab5c01d65a66cd$var$generateVariantDefinitions(`[data-theme='${themeName}']`),
                $1cab5c01d65a66cd$var$generateRootVariantDefinitions(themeName, defaultTheme)
            ].flat()
        });
        const cssSelector = `.${themeClassName},[data-theme="${themeName}"]`;
        resolved.utilities[cssSelector] = colors[$1cab5c01d65a66cd$var$SCHEME] ? {
            "color-scheme": colors[$1cab5c01d65a66cd$var$SCHEME]
        } : {};
        (0, (/*@__PURE__*/$parcel$interopDefault($ff1d430d0589e4d1$exports)))(flatColors, (colorValue, colorName)=>{
            if (colorName === $1cab5c01d65a66cd$var$SCHEME) return;
            const safeColorName = $1cab5c01d65a66cd$var$escapeChars(colorName, "/");
            let [h, s, l, defaultAlphaValue] = [
                0,
                0,
                0,
                1
            ];
            try {
                [h, s, l, defaultAlphaValue] = $1cab5c01d65a66cd$var$toHslaArray(colorValue);
            } catch (error) {
                const message = `\r
Warning - In theme "${themeName}" color "${colorName}". ${error.message}`;
                if (strict) throw new Error(message);
                return console.error(message);
            }
            const twcColorVariable = produceCssVariable(safeColorName);
            const twcOpacityVariable = `${produceCssVariable(safeColorName)}-opacity`;
            const hslValues = `${h} ${s}% ${l}%`;
            resolved.utilities[cssSelector][twcColorVariable] = hslValues;
            $1cab5c01d65a66cd$var$addRootUtilities(resolved.utilities, {
                key: twcColorVariable,
                value: hslValues,
                defaultTheme: defaultTheme,
                themeName: themeName
            });
            if (typeof defaultAlphaValue === "number") {
                const alphaValue = defaultAlphaValue.toFixed(2);
                resolved.utilities[cssSelector][twcOpacityVariable] = alphaValue;
                $1cab5c01d65a66cd$var$addRootUtilities(resolved.utilities, {
                    key: twcOpacityVariable,
                    value: alphaValue,
                    defaultTheme: defaultTheme,
                    themeName: themeName
                });
            }
            resolved.colors[colorName] = ({ opacityVariable: opacityVariable, opacityValue: opacityValue })=>{
                if (!isNaN(+opacityValue)) return `hsl(var(${twcColorVariable}) / ${opacityValue})`;
                if (opacityVariable) return `hsl(var(${twcColorVariable}) / var(${twcOpacityVariable}, var(${opacityVariable})))`;
                return `hsl(var(${twcColorVariable}) / var(${twcOpacityVariable}, 1))`;
            };
        });
    });
    return resolved;
};
var $1cab5c01d65a66cd$export$d712115a11635717 = (config = $1cab5c01d65a66cd$var$emptyConfig, options = {})=>{
    const resolved = $1cab5c01d65a66cd$export$231aedbc8f7c4964(config, options);
    return (0, (/*@__PURE__*/$parcel$interopDefault($84df4a417d8748cd$exports)))(({ addUtilities: addUtilities, addVariant: addVariant })=>{
        addUtilities(resolved.utilities);
        resolved.variants.forEach(({ name: name, definition: definition })=>addVariant(name, definition));
    }, // extend the colors config
    {
        theme: {
            extend: {
                // @ts-ignore tailwind types are broken
                colors: resolved.colors
            }
        }
    });
};
function $1cab5c01d65a66cd$var$escapeChars(str, ...chars) {
    let result = str;
    for (let char of chars){
        const regexp = new RegExp(char, "g");
        result = str.replace(regexp, "\\" + char);
    }
    return result;
}
function $1cab5c01d65a66cd$var$flattenColors(colors) {
    const flatColorsWithDEFAULT = (0, (/*@__PURE__*/$parcel$interopDefault($4c07cb87e79a40c4$exports)))(colors, {
        safe: true,
        delimiter: "-"
    });
    return Object.entries(flatColorsWithDEFAULT).reduce((acc, [key, value])=>{
        acc[key.replace(/\-DEFAULT$/, "")] = value;
        return acc;
    }, {});
}
function $1cab5c01d65a66cd$var$toHslaArray(colorValue) {
    return (0, (/*@__PURE__*/$parcel$interopDefault($040ba8e83b63f6c5$exports)))(colorValue).hsl().round(1).array();
}
function $1cab5c01d65a66cd$var$defaultProduceCssVariable(themeName) {
    return `--twc-${themeName}`;
}
function $1cab5c01d65a66cd$var$defaultProduceThemeClass(themeName) {
    return themeName;
}
function $1cab5c01d65a66cd$var$dark(colors) {
    return {
        ...colors,
        [$1cab5c01d65a66cd$var$SCHEME]: "dark"
    };
}
function $1cab5c01d65a66cd$var$light(colors) {
    return {
        ...colors,
        [$1cab5c01d65a66cd$var$SCHEME]: "light"
    };
}
function $1cab5c01d65a66cd$var$generateVariantDefinitions(selector) {
    return [
        `${selector}&`,
        `:is(${selector} > &:not([data-theme]))`,
        `:is(${selector} &:not(${selector} [data-theme]:not(${selector}) * ))`,
        `:is(${selector}:not(:has([data-theme])) &:not([data-theme]))`
    ];
}
function $1cab5c01d65a66cd$var$generateRootVariantDefinitions(themeName, defaultTheme) {
    const baseDefinitions = [
        `:root&`,
        `:is(:root > &:not([data-theme]))`,
        `:is(:root &:not([data-theme] *):not([data-theme]))`
    ];
    if (typeof defaultTheme === "string" && themeName === defaultTheme) return baseDefinitions;
    if (typeof defaultTheme === "object" && themeName === defaultTheme.light) return baseDefinitions.map((definition)=>`@media (prefers-color-scheme: light){${definition}}`);
    if (typeof defaultTheme === "object" && themeName === defaultTheme.dark) return baseDefinitions.map((definition)=>`@media (prefers-color-scheme: dark){${definition}}`);
    return [];
}
function $1cab5c01d65a66cd$var$addRootUtilities(utilities, { key: key, value: value, defaultTheme: defaultTheme, themeName: themeName }) {
    if (!defaultTheme) return;
    if (typeof defaultTheme === "string") {
        if (themeName === defaultTheme) {
            if (!utilities[":root"]) utilities[":root"] = {};
            utilities[":root"][key] = value;
        }
    } else if (themeName === defaultTheme.light) {
        if (!utilities["@media (prefers-color-scheme: light)"]) utilities["@media (prefers-color-scheme: light)"] = {
            ":root": {}
        };
        utilities["@media (prefers-color-scheme: light)"][":root"][key] = value;
    } else if (themeName === defaultTheme.dark) {
        if (!utilities["@media (prefers-color-scheme: dark)"]) utilities["@media (prefers-color-scheme: dark)"] = {
            ":root": {}
        };
        utilities["@media (prefers-color-scheme: dark)"][":root"][key] = value;
    }
}


// packages/morph/src/morph.js
function $fd2717825660a9ff$var$morph(from, toHtml, options) {
    $fd2717825660a9ff$var$monkeyPatchDomSetAttributeToAllowAtSymbols();
    let fromEl;
    let toEl;
    let key, lookahead, updating, updated, removing, removed, adding, added;
    function assignOptions(options2 = {}) {
        let defaultGetKey = (el)=>el.getAttribute("key");
        let noop = ()=>{};
        updating = options2.updating || noop;
        updated = options2.updated || noop;
        removing = options2.removing || noop;
        removed = options2.removed || noop;
        adding = options2.adding || noop;
        added = options2.added || noop;
        key = options2.key || defaultGetKey;
        lookahead = options2.lookahead || false;
    }
    function patch(from2, to) {
        if (differentElementNamesTypesOrKeys(from2, to)) return swapElements(from2, to);
        let updateChildrenOnly = false;
        if ($fd2717825660a9ff$var$shouldSkip(updating, from2, to, ()=>updateChildrenOnly = true)) return;
        if (from2.nodeType === 1 && window.Alpine) {
            window.Alpine.cloneNode(from2, to);
            if (from2._x_teleport && to._x_teleport) patch(from2._x_teleport, to._x_teleport);
        }
        if ($fd2717825660a9ff$var$textOrComment(to)) {
            patchNodeValue(from2, to);
            updated(from2, to);
            return;
        }
        if (!updateChildrenOnly) patchAttributes(from2, to);
        updated(from2, to);
        patchChildren(from2, to);
    }
    function differentElementNamesTypesOrKeys(from2, to) {
        return from2.nodeType != to.nodeType || from2.nodeName != to.nodeName || getKey(from2) != getKey(to);
    }
    function swapElements(from2, to) {
        if ($fd2717825660a9ff$var$shouldSkip(removing, from2)) return;
        let toCloned = to.cloneNode(true);
        if ($fd2717825660a9ff$var$shouldSkip(adding, toCloned)) return;
        from2.replaceWith(toCloned);
        removed(from2);
        added(toCloned);
    }
    function patchNodeValue(from2, to) {
        let value = to.nodeValue;
        if (from2.nodeValue !== value) from2.nodeValue = value;
    }
    function patchAttributes(from2, to) {
        if (from2._x_transitioning) return;
        if (from2._x_isShown && !to._x_isShown) return;
        if (!from2._x_isShown && to._x_isShown) return;
        let domAttributes = Array.from(from2.attributes);
        let toAttributes = Array.from(to.attributes);
        for(let i = domAttributes.length - 1; i >= 0; i--){
            let name = domAttributes[i].name;
            if (!to.hasAttribute(name)) from2.removeAttribute(name);
        }
        for(let i = toAttributes.length - 1; i >= 0; i--){
            let name = toAttributes[i].name;
            let value = toAttributes[i].value;
            if (from2.getAttribute(name) !== value) from2.setAttribute(name, value);
        }
    }
    function patchChildren(from2, to) {
        let fromKeys = keyToMap(from2.children);
        let fromKeyHoldovers = {};
        let currentTo = $fd2717825660a9ff$var$getFirstNode(to);
        let currentFrom = $fd2717825660a9ff$var$getFirstNode(from2);
        while(currentTo){
            $fd2717825660a9ff$var$seedingMatchingId(currentTo, currentFrom);
            let toKey = getKey(currentTo);
            let fromKey = getKey(currentFrom);
            if (!currentFrom) {
                if (toKey && fromKeyHoldovers[toKey]) {
                    let holdover = fromKeyHoldovers[toKey];
                    from2.appendChild(holdover);
                    currentFrom = holdover;
                    fromKey = getKey(currentFrom);
                } else {
                    if (!$fd2717825660a9ff$var$shouldSkip(adding, currentTo)) {
                        let clone = currentTo.cloneNode(true);
                        from2.appendChild(clone);
                        added(clone);
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
                    let next = $fd2717825660a9ff$var$getNextSibling(from2, currentFrom);
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
                patchChildren(fromBlock, toBlock);
                continue;
            }
            if (currentFrom.nodeType === 1 && lookahead && !currentFrom.isEqualNode(currentTo)) {
                let nextToElementSibling = $fd2717825660a9ff$var$getNextSibling(to, currentTo);
                let found = false;
                while(!found && nextToElementSibling){
                    if (nextToElementSibling.nodeType === 1 && currentFrom.isEqualNode(nextToElementSibling)) {
                        found = true;
                        currentFrom = addNodeBefore(from2, currentTo, currentFrom);
                        fromKey = getKey(currentFrom);
                    }
                    nextToElementSibling = $fd2717825660a9ff$var$getNextSibling(to, nextToElementSibling);
                }
            }
            if (toKey !== fromKey) {
                if (!toKey && fromKey) {
                    fromKeyHoldovers[fromKey] = currentFrom;
                    currentFrom = addNodeBefore(from2, currentTo, currentFrom);
                    fromKeyHoldovers[fromKey].remove();
                    currentFrom = $fd2717825660a9ff$var$getNextSibling(from2, currentFrom);
                    currentTo = $fd2717825660a9ff$var$getNextSibling(to, currentTo);
                    continue;
                }
                if (toKey && !fromKey) {
                    if (fromKeys[toKey]) {
                        currentFrom.replaceWith(fromKeys[toKey]);
                        currentFrom = fromKeys[toKey];
                        fromKey = getKey(currentFrom);
                    }
                }
                if (toKey && fromKey) {
                    let fromKeyNode = fromKeys[toKey];
                    if (fromKeyNode) {
                        fromKeyHoldovers[fromKey] = currentFrom;
                        currentFrom.replaceWith(fromKeyNode);
                        currentFrom = fromKeyNode;
                        fromKey = getKey(currentFrom);
                    } else {
                        fromKeyHoldovers[fromKey] = currentFrom;
                        currentFrom = addNodeBefore(from2, currentTo, currentFrom);
                        fromKeyHoldovers[fromKey].remove();
                        currentFrom = $fd2717825660a9ff$var$getNextSibling(from2, currentFrom);
                        currentTo = $fd2717825660a9ff$var$getNextSibling(to, currentTo);
                        continue;
                    }
                }
            }
            let currentFromNext = currentFrom && $fd2717825660a9ff$var$getNextSibling(from2, currentFrom);
            patch(currentFrom, currentTo);
            currentTo = currentTo && $fd2717825660a9ff$var$getNextSibling(to, currentTo);
            currentFrom = currentFromNext;
        }
        let removals = [];
        while(currentFrom){
            if (!$fd2717825660a9ff$var$shouldSkip(removing, currentFrom)) removals.push(currentFrom);
            currentFrom = $fd2717825660a9ff$var$getNextSibling(from2, currentFrom);
        }
        while(removals.length){
            let domForRemoval = removals.shift();
            domForRemoval.remove();
            removed(domForRemoval);
        }
    }
    function getKey(el) {
        return el && el.nodeType === 1 && key(el);
    }
    function keyToMap(els) {
        let map = {};
        for (let el of els){
            let theKey = getKey(el);
            if (theKey) map[theKey] = el;
        }
        return map;
    }
    function addNodeBefore(parent, node, beforeMe) {
        if (!$fd2717825660a9ff$var$shouldSkip(adding, node)) {
            let clone = node.cloneNode(true);
            parent.insertBefore(clone, beforeMe);
            added(clone);
            return clone;
        }
        return node;
    }
    assignOptions(options);
    fromEl = from;
    toEl = typeof toHtml === "string" ? $fd2717825660a9ff$var$createElement(toHtml) : toHtml;
    if (window.Alpine && window.Alpine.closestDataStack && !from._x_dataStack) {
        toEl._x_dataStack = window.Alpine.closestDataStack(from);
        toEl._x_dataStack && window.Alpine.cloneNode(from, toEl);
    }
    patch(from, toEl);
    fromEl = void 0;
    toEl = void 0;
    return from;
}
$fd2717825660a9ff$var$morph.step = ()=>{};
$fd2717825660a9ff$var$morph.log = ()=>{};
function $fd2717825660a9ff$var$shouldSkip(hook, ...args) {
    let skip = false;
    hook(...args, ()=>skip = true);
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
}
// packages/morph/builds/module.js
var $fd2717825660a9ff$export$2e2bcd8739ae039 = $fd2717825660a9ff$export$2e5e8c41f5d4e7c7;


var $ae6560896f9c4f0c$exports = {};
$ae6560896f9c4f0c$exports = JSON.parse("{\"CachePath\":\"local cache location\",\"CachePath_Description\":\"Local image cache location, default system temporary folder.\",\"CertFile\":\"CertFile\",\"CertFile_Description\":\"TLS/SSL certificate file path (default: , \\\"~/.config/.comigo/cert.crt\\\")\",\"ClearCacheExit\":\"Clean up on exit\",\"ClearCacheExit_Description\":\"When exiting the program, clear the web image cache.\",\"ClearDatabaseWhenExit\":\"Clear database books\",\"ClearDatabaseWhenExit_Description\":\"When the local database is enabled, non-existing books are purged after the scan is completed.\",\"ConfigManager\":\"Profile management\",\"ConfigManagerDeleteSuccess\":\"Configuration has been deleted.\",\"ConfigManagerDescription\":\"Clicking Save will upload the current configuration to the server and overwrite the existing configuration file.\",\"ConfigManagerSaveHint\":\"There is already a configuration file, please change the save location.\",\"ConfigManagerSaveSuccess\":\"Configuration saved.\",\"ConfigSaveTo\":\"Configure default save path\",\"Debug\":\"Turn on Debug mode\",\"Debug_Description\":\"Enable Debug function\",\"DisableLAN\":\"Disable LAN sharing\",\"DisableLAN_Description\":\"Reading services are only provided on this machine and are not shared externally. This configuration does not support hot reloading.\",\"EnableDatabase\":\"Enable database\",\"EnableDatabase_Description\":\"Enable local database to save scanned book data. \\nThis configuration does not support hot reload.\",\"EnableFrpcServer\":\"EnableFrpcServer\",\"EnableLogin\":\"Enable login\",\"EnableLogin_Description\":\"Whether to enable login. \\nNo login is required by default. \\nThis configuration does not support hot reload.\",\"EnableTLS\":\"Enable TLS\",\"EnableTLS_Description\":\"Whether to enable HTTPS protocol. \\nThe certificate needs to be set in the key file.\",\"EnableUpload\":\"Enable upload functionality\",\"EnableUpload_Description\":\"Enable upload functionality.\",\"ExcludePath\":\"exclude path\",\"ExcludePath_Description\":\"When scanning books, the names of files or folders that need to be excluded\",\"FrpClientConfig\":\"FrpClient settings\",\"GenerateBookMetadata\":\"Generate book metadata\",\"GenerateMetaData\":\"Generate metadata\",\"GenerateMetaData_Description\":\"Generate book metadata. \\nNot currently in effect.\",\"HomeDirectory\":\"HomeDirectory\",\"Host\":\"domain name\",\"Host_Description\":\"Customize the host name displayed by the QR code. \\nThe default is the network card IP.\",\"KeyFile\":\"KeyFile\",\"KeyFile_Description\":\"TLS/SSL key file path (default: \\\"~/.config/.comigo/key.key\\\")\",\"LocalStores\":\"library folder\",\"LocalStores_Description\":\"The library folder supports absolute and relative directories. \\nRelative directories are based on the current execution directory.\",\"LogFileName\":\"Log file name\",\"LogFileName_Description\":\"Log file name\",\"LogFilePath\":\"Log save location\",\"LogFilePath_Description\":\"Log file save location\",\"LogToFile\":\"Record Log to local\",\"LogToFile_Description\":\"Whether to save the program log to a local file. \\nNot saved by default.\",\"MaxScanDepth\":\"Maximum scan depth\",\"MaxScanDepth_Description\":\"Maximum scan depth. \\nFiles exceeding the depth will not be scanned. \\nThe current execution directory is the base.\",\"MinImageNum\":\"Minimum number of pictures\",\"MinImageNum_Description\":\"A compressed package or folder must contain at least a few pictures to be considered a book.\",\"OpenBrowser\":\"Open browser\",\"OpenBrowser_Description\":\"After the scan is completed, whether to open the browser at the same time. \\nThe default is true for windows and false for other platforms.\",\"Password\":\"password\",\"Password_Description\":\"After enabling login, the password used to log in.\",\"Port\":\"port\",\"Port_Description\":\"Web service port. This configuration does not support hot reloading.\",\"PrintAllPossibleQRCode\":\"More QR codes\",\"ProgramDirectory\":\"The directory where the program is located\",\"StartFrpClientInBackground\":\"Start FrpClient\",\"SupportFileType\":\"Supported compressed packages\",\"SupportFileType_Description\":\"When scanning a file, it is used to decide whether to skip or count it as a file suffix for book processing.\",\"SupportMediaType\":\"Supported image files\",\"SupportMediaType_Description\":\"Image file suffix used to count the number of images when scanning compressed packages\",\"Timeout\":\"Expiration\",\"TimeoutLimitForScan\":\"Scan timeout\",\"TimeoutLimitForScan_Description\":\"When scanning a file, if it takes more than a few seconds, it will give up scanning the file to avoid getting stuck on an overly large file.\",\"Timeout_Description\":\"Cookie expiration time after enabling login. \\nThe unit is minutes. \\nIt expires in 180 minutes by default.\",\"UploadPath\":\"Upload location\",\"UploadPath_Description\":\"Customize the storage location of the uploaded files, by default, the upload folder is created under the current execution directory or the first bookstore directory.\",\"UseCache\":\"Local image cache\",\"Username\":\"username\",\"Username_Description\":\"After enabling login, the username required for the login interface.\",\"WorkingDirectory\":\"current working directory\",\"ZipFileTextEncoding\":\"Not UTF-8\",\"ZipFileTextEncoding_Description\":\"Non-utf-8 encoded ZIP file, what encoding should be used to parse it. \\nDefault GBK.\",\"abs\":\"labs\",\"all_page_num\":\"Total Pages: {0}\",\"author\":\"Author: {0}\",\"auto_crop\":\"Automatic edge trimming\",\"auto_double_page\":\"Auto Double Page (beta)\",\"auto_hide_toolbar\":\"Auto Hide Toolbar\",\"back-to-top\":\"Back to Top\",\"back_button\":\"Back Button\",\"back_to_bookshelf\":\"Back to Bookshelf\",\"book_shelf\":\"Bookshelf\",\"child_book_hint\":\"{0} books in the folder\",\"debug_mode\":\"Debug Mode\",\"do_you_reset_all_settings\":\"Do you want to reset to the default settings?\",\"double_page_mode\":\"Double Page Mode\",\"double_page_width\":\"Double Page Width:\",\"download_sample_config_file\":\"Download Sample Config File\",\"download_windows_reg_file\":\"Download Windows Reg File\",\"drop_to_upload\":\"Click or drag files to this area to upload\",\"energy_threshold\":\"Energy Threshold:\",\"epub_info\":\"ePub Information\",\"exit_fullscreen\":\"Exit Fullscreen Mode\",\"filesize\":\"File Size: {0}\",\"flip_mode\":\"Flip Mode\",\"flip_odd_even_page\":\"Flip Odd/Even Pages\",\"flip_odd_even_page_hint\":\"Click here if double pages do not align properly.\",\"found_read_history\":\"Local Reading History Found\",\"from_interrupt\":\"From Last Interruption\",\"full_screen_hint\":\"Fullscreen Button\",\"fullscreen\":\"Fullscreen\",\"good_job_and_byebye\":\"Good job! Goodbye.\",\"gray_image\":\"Grayscale Image\",\"hint\":\"hint\",\"hint_first_page\":\"You are on the first page and cannot turn forward.\",\"hint_last_page\":\"You are on the last page and cannot turn backward.\",\"hour\":\"hours\",\"image_width_limit\":\"Limit Width\",\"infinite_dropdown\":\"Infinite Dropdown\",\"interval\":\"Interval:\",\"left_screen_to_next\":\"Left to Right - Comic\",\"load_all_pages\":\"Load All Pages\",\"load_from_interrupt\":\"Load from last reading position (page XX)?\",\"login_success_hint\":\"Login successful. Returning to the previous page.\",\"logout\":\"Logout\",\"margin_bottom_on_scroll_mode\":\"Margin Bottom:\",\"margin_on_scroll_mode\":\"Page gap:\",\"max_width\":\"Max Width:\",\"minute\":\"minutes\",\"network\":\"Network\",\"no_book_found_hint\":\"No books found. Try uploading a file?\",\"no_support_upload_file\":\"File upload functionality has been disabled by the administrator.\",\"not_support_fullscreen\":\"This browser does not support fullscreen mode.\",\"now_is\":\"Now:\",\"number_of_online_books\":\"Number of Online Books:\",\"original_image\":\"Original Image\",\"original_pdf_link\":\"Original PDF Link\",\"page\":\"page\",\"page_turning_seconds\":\"Page Turn Interval:\",\"pagination_mode\":\"Pagination Mode\",\"pdf_hint_message\":\"Supports pure image PDFs. If loading is slow or errors occur, please try the following:\",\"please_enable_upload\":\"Please enable server upload support.\",\"please_enter_content\":\"Please enter content\",\"qrcode_hint\":\"Scan to read. Click to display QR code.\",\"raw_resolution\":\"Raw Resolution\",\"re_sort_book\":\"Resort Books\",\"re_sort_page\":\"Resort Pages\",\"reader_settings\":\"Reader Settings\",\"reading_progress_bar\":\"Reading Progress Bar\",\"refresh_page\":\"Refresh Page\",\"reset_all_settings\":\"Reset Settings\",\"resort_file\":\"Resort File\",\"right_screen_to_next\":\"Right to Left - Manga\",\"save_page_num\":\"Save Progress\",\"scan_qrcode\":\"Scan QR Code:\",\"scanned_hint\":\"Scanned XX books. Do you want to view them now?\",\"scroll_mode\":\"Scroll Mode\",\"second\":\"seconds\",\"select_language\":\"Select Language\",\"server_config\":\"Server Configuration\",\"server_setting\":\"Comigo Server Settings\",\"set_back_color\":\"Background Color:\",\"set_interface_color\":\"Interface Color:\",\"show_book_titles\":\"Show Book Titles\",\"show_file_icon\":\"Show File Icons\",\"show_header\":\"Show Header\",\"show_page_num\":\"Show Page Number\",\"simplify_book_titles\":\"Simplify Book Titles\",\"single_page_mode\":\"Single Page Mode\",\"single_page_width\":\"Single Page Width:\",\"sort_by_default\":\"Default Sort\",\"sort_by_filename\":\"Sort by Filename (A-Z)\",\"sort_by_filename_reverse\":\"Sort by Filename (Z-A)\",\"sort_by_filesize\":\"Sort by Filesize (Large to Small)\",\"sort_by_filesize_reverse\":\"Sort by Filesize (Small to Large)\",\"sort_by_modify_time\":\"Sort by Modify Time (Newest to Oldest)\",\"sort_by_modify_time_reverse\":\"Sort by Modify Time (Oldest to Newest)\",\"sort_reverse\":\"(Reverse)\",\"start_sketch_message\":\"The countdown has begun. Have a nice day!\",\"start_sketch_mode\":\"Start Sketch\",\"starting_from_beginning\":\"Start from Beginning\",\"starting_from_beginning_hint\":\"Load from the first page\",\"stop_sketch_mode\":\"Stop Sketch\",\"submit\":\"submit\",\"success_fullscreen\":\"Entered Fullscreen Mode\",\"successfully_loaded_reading_progress\":\"Successfully Loaded Reading Progress\",\"switch_to_flip_mode\":\"Switch to Flip Mode\",\"switch_to_scrolling_mode\":\"Switch to Scroll Mode\",\"sync_page\":\"Remote Page Sync\",\"temp_future_hint\":\"Temporarily put some features that are not yet finished, under development and adjustment.\",\"test\":\"Test\",\"to_flip_mode\":\"Switch to Flip Mode\",\"to_infinite_dropdown_mode\":\"Switch to Infinite Dropdown Mode\",\"to_pagination_mode\":\"Switch to Pagination Mode\",\"to_scroll_mode\":\"Switch to Scroll Mode\",\"total_is\":\"Total:\",\"total_time\":\"Total Time:\",\"type_or_paste_content\":\"Type or paste content\",\"upload_file\":\"Upload File\",\"uploaded_folder_hint\":\"Files will be uploaded to the upload folder under the default library directory.\",\"width_use_fixed_value\":\"Landscape Mode Width: Fixed Value\",\"width_use_percent\":\"Landscape Mode Width: Percentage\"}");


var $c8762dc563519366$exports = {};
$c8762dc563519366$exports = JSON.parse('{"CachePath":"\u672C\u5730\u7F13\u5B58\u4F4D\u7F6E","CachePath_Description":"\u672C\u5730\u56FE\u7247\u7F13\u5B58\u4F4D\u7F6E\uFF0C\u9ED8\u8BA4\u7CFB\u7EDF\u4E34\u65F6\u6587\u4EF6\u5939\u3002","CertFile":"CertFile","CertFile_Description":"TLS/SSL \u8BC1\u4E66\u6587\u4EF6\u8DEF\u5F84 (default: \u3001\\"~/.config/.comigo/cert.crt\\")","ClearCacheExit":"\u9000\u51FA\u65F6\u6E05\u7406","ClearCacheExit_Description":"\u9000\u51FA\u7A0B\u5E8F\u7684\u65F6\u5019\uFF0C\u6E05\u7406web\u56FE\u7247\u7F13\u5B58\u3002","ClearDatabaseWhenExit":"\u6E05\u9664\u6570\u636E\u5E93\u4E66\u7C4D","ClearDatabaseWhenExit_Description":"\u542F\u7528\u672C\u5730\u6570\u636E\u5E93\u65F6\uFF0C\u626B\u63CF\u5B8C\u6210\u540E\uFF0C\u6E05\u9664\u4E0D\u5B58\u5728\u7684\u4E66\u7C4D\u3002","ConfigManager":"\u914D\u7F6E\u6587\u4EF6\u7BA1\u7406","ConfigManagerDeleteSuccess":"\u914D\u7F6E\u5DF2\u5220\u9664\u3002","ConfigManagerDescription":"\u70B9\u51FBSave\uFF0C\u4F1A\u5C06\u5F53\u524D\u914D\u7F6E\u4E0A\u4F20\u5230\u670D\u52A1\u5668\uFF0C\u5E76\u8986\u76D6\u5DF2\u7ECF\u5B58\u5728\u7684\u8BBE\u5B9A\u6587\u4EF6\u3002","ConfigManagerSaveHint":"\u5DF2\u6709\u914D\u7F6E\u6587\u4EF6,\u8BF7\u5207\u6362\u4FDD\u5B58\u4F4D\u7F6E\u3002","ConfigManagerSaveSuccess":"\u914D\u7F6E\u5DF2\u4FDD\u5B58\u3002","ConfigSaveTo":"\u914D\u7F6E\u9ED8\u8BA4\u4FDD\u5B58\u8DEF\u5F84","Debug":"\u5F00\u542FDebug\u6A21\u5F0F","Debug_Description":"\u542F\u7528Debug\u529F\u80FD","DisableLAN":"\u7981\u6B62\u5C40\u57DF\u7F51\u5171\u4EAB","DisableLAN_Description":"\u53EA\u5728\u672C\u673A\u63D0\u4F9B\u9605\u8BFB\u670D\u52A1\uFF0C\u4E0D\u5BF9\u5916\u5171\u4EAB\uFF0C\u6B64\u9879\u914D\u7F6E\u4E0D\u652F\u6301\u70ED\u91CD\u8F7D","EnableDatabase":"\u542F\u7528\u6570\u636E\u5E93","EnableDatabase_Description":"\u542F\u7528\u672C\u5730\u6570\u636E\u5E93\uFF0C\u4FDD\u5B58\u626B\u63CF\u5230\u7684\u4E66\u7C4D\u6570\u636E\u3002\u6B64\u9879\u914D\u7F6E\u4E0D\u652F\u6301\u70ED\u91CD\u8F7D\u3002","EnableFrpcServer":"EnableFrpcServer","EnableLogin":"\u542F\u7528\u767B\u9646","EnableLogin_Description":"\u662F\u5426\u542F\u7528\u767B\u5F55\u3002\u9ED8\u8BA4\u4E0D\u9700\u8981\u767B\u9646\u3002\u6B64\u9879\u914D\u7F6E\u4E0D\u652F\u6301\u70ED\u91CD\u8F7D\u3002","EnableTLS":"Enable TLS","EnableTLS_Description":"\u662F\u5426\u542F\u7528HTTPS\u534F\u8BAE\u3002\u9700\u8981\u8BBE\u7F6E\u8BC1\u4E66\u4E8Ekey\u6587\u4EF6\u3002","EnableUpload":"\u542F\u7528\u4E0A\u4F20\u529F\u80FD","EnableUpload_Description":"\u542F\u7528\u4E0A\u4F20\u529F\u80FD\u3002","ExcludePath":"\u6392\u9664\u8DEF\u5F84","ExcludePath_Description":"\u626B\u63CF\u4E66\u7C4D\u7684\u65F6\u5019\uFF0C\u9700\u8981\u6392\u9664\u7684\u6587\u4EF6\u6216\u6587\u4EF6\u5939\u7684\u540D\u5B57","FrpClientConfig":"FrpClient\u8BBE\u7F6E","GenerateBookMetadata":"\u751F\u6210\u4E66\u7C4D\u5143\u6570\u636E","GenerateMetaData":"\u751F\u6210\u5143\u6570\u636E","GenerateMetaData_Description":"\u751F\u6210\u4E66\u7C4D\u5143\u6570\u636E\u3002\u5F53\u524D\u672A\u751F\u6548\u3002","HomeDirectory":"\u7528\u6237\u4E3B\u76EE\u5F55","Host":"\u57DF\u540D","Host_Description":"\u81EA\u5B9A\u4E49\u4E8C\u7EF4\u7801\u663E\u793A\u7684\u4E3B\u673A\u540D\u3002\u9ED8\u8BA4\u4E3A\u7F51\u5361IP\u3002","KeyFile":"KeyFile","KeyFile_Description":"TLS/SSL key\u6587\u4EF6\u8DEF\u5F84 (default: \\"~/.config/.comigo/key.key\\")","LocalStores":"\u4E66\u5E93\u6587\u4EF6\u5939","LocalStores_Description":"\u4E66\u5E93\u6587\u4EF6\u5939\uFF0C\u652F\u6301\u7EDD\u5BF9\u76EE\u5F55\u4E0E\u76F8\u5BF9\u76EE\u5F55\u3002\u76F8\u5BF9\u76EE\u5F55\u4EE5\u5F53\u524D\u6267\u884C\u76EE\u5F55\u4E3A\u57FA\u51C6","LogFileName":"Log\u6587\u4EF6\u540D","LogFileName_Description":"Log\u6587\u4EF6\u540D","LogFilePath":"Log\u4FDD\u5B58\u4F4D\u7F6E","LogFilePath_Description":"Log\u6587\u4EF6\u7684\u4FDD\u5B58\u4F4D\u7F6E","LogToFile":"\u8BB0\u5F55Log\u5230\u672C\u5730","LogToFile_Description":"\u662F\u5426\u4FDD\u5B58\u7A0B\u5E8FLog\u5230\u672C\u5730\u6587\u4EF6\u3002\u9ED8\u8BA4\u4E0D\u4FDD\u5B58\u3002","MaxScanDepth":"\u6700\u5927\u626B\u63CF\u6DF1\u5EA6","MaxScanDepth_Description":"\u6700\u5927\u626B\u63CF\u6DF1\u5EA6\u3002\u8D85\u8FC7\u6DF1\u5EA6\u7684\u6587\u4EF6\u4E0D\u4F1A\u88AB\u626B\u63CF\u3002\u5F53\u524D\u6267\u884C\u76EE\u5F55\u4E3A\u57FA\u51C6\u3002","MinImageNum":"\u6700\u5C0F\u56FE\u7247\u6570","MinImageNum_Description":"\u538B\u7F29\u5305\u6216\u6587\u4EF6\u5939\u5185\u81F3\u5C11\u6709\u51E0\u5F20\u56FE\u7247\uFF0C\u624D\u7B97\u4F5C\u4E66\u7C4D\u3002","OpenBrowser":"\u6253\u5F00\u6D4F\u89C8\u5668","OpenBrowser_Description":"\u626B\u63CF\u5B8C\u6210\u540E\uFF0C\u662F\u5426\u540C\u65F6\u6253\u5F00\u6D4F\u89C8\u5668\u3002windows\u9ED8\u8BA4true\uFF0C\u5176\u4ED6\u5E73\u53F0\u9ED8\u8BA4false\u3002","Password":"\u5BC6\u7801","Password_Description":"\u542F\u7528\u767B\u9646\u540E\uFF0C\u767B\u5F55\u7528\u7684\u5BC6\u7801\u3002","Port":"\u7AEF\u53E3","Port_Description":"\u7F51\u9875\u670D\u52A1\u7AEF\u53E3\uFF0C\u6B64\u9879\u914D\u7F6E\u4E0D\u652F\u6301\u70ED\u91CD\u8F7D","PrintAllPossibleQRCode":"\u66F4\u591A\u4E8C\u7EF4\u7801","ProgramDirectory":"\u7A0B\u5E8F\u6240\u5728\u76EE\u5F55","StartFrpClientInBackground":"\u542F\u52A8FrpClient","SupportFileType":"\u652F\u6301\u7684\u538B\u7F29\u5305","SupportFileType_Description":"\u626B\u63CF\u6587\u4EF6\u65F6\uFF0C\u7528\u4E8E\u51B3\u5B9A\u8DF3\u8FC7\uFF0C\u8FD8\u662F\u7B97\u4F5C\u4E66\u7C4D\u5904\u7406\u7684\u6587\u4EF6\u540E\u7F00","SupportMediaType":"\u652F\u6301\u7684\u56FE\u7247\u6587\u4EF6","SupportMediaType_Description":"\u626B\u63CF\u538B\u7F29\u5305\u65F6\uFF0C\u7528\u4E8E\u7EDF\u8BA1\u56FE\u7247\u6570\u91CF\u7684\u56FE\u7247\u6587\u4EF6\u540E\u7F00","Timeout":"\u8FC7\u671F\u65F6\u95F4","TimeoutLimitForScan":"\u626B\u63CF\u8D85\u65F6","TimeoutLimitForScan_Description":"\u626B\u63CF\u6587\u4EF6\u65F6\uFF0C\u8D85\u8FC7\u51E0\u79D2\u949F\uFF0C\u5C31\u653E\u5F03\u626B\u63CF\u8FD9\u4E2A\u6587\u4EF6\uFF0C\u907F\u514D\u5361\u5728\u8FC7\u5927\u6587\u4EF6\u4E0A\u3002","Timeout_Description":"\u542F\u7528\u767B\u9646\u540E\uFF0Ccookie\u8FC7\u671F\u65F6\u95F4\u3002\u5355\u4F4D\u4E3A\u5206\u949F\u3002\u9ED8\u8BA4180\u5206\u8FC7\u671F\u3002","UploadPath":"\u4E0A\u4F20\u4F4D\u7F6E","UploadPath_Description":"\u81EA\u5B9A\u4E49\u4E0A\u4F20\u6587\u4EF6\u5B58\u50A8\u4F4D\u7F6E\uFF0C\u9ED8\u8BA4\u5728\u5F53\u524D\u6267\u884C\u76EE\u5F55\u6216\u7B2C\u4E00\u4E2A\u4E66\u5E93\u76EE\u5F55\u4E0B\u9762\u521B\u5EFA upload \u6587\u4EF6\u5939\u3002","UseCache":"\u672C\u5730\u56FE\u7247\u7F13\u5B58","Username":"\u7528\u6237\u540D","Username_Description":"\u542F\u7528\u767B\u9646\u540E\uFF0C\u767B\u5F55\u754C\u9762\u9700\u8981\u7684\u7528\u6237\u540D\u3002","WorkingDirectory":"\u5F53\u524D\u5DE5\u4F5C\u76EE\u5F55","ZipFileTextEncoding":"\u975EUTF-8","ZipFileTextEncoding_Description":"\u975Eutf-8\u7F16\u7801ZIP\u6587\u4EF6\uFF0C\u5C1D\u8BD5\u7528\u4EC0\u4E48\u7F16\u7801\u89E3\u6790\u3002\u9ED8\u8BA4GBK\u3002","abs":"\u5B9E\u9A8C","all_page_num":"\u603B\u9875\u6570\uFF1A{0}","author":"\u4F5C\u8005\uFF1A{0}","auto_crop":"\u81EA\u52A8\u5207\u8FB9","auto_double_page":"\u81EA\u52A8\u5408\u5E76\u53CC\u9875\uFF08beta\uFF09","auto_hide_toolbar":"\u81EA\u52A8\u9690\u85CF\u5DE5\u5177\u680F","back-to-top":"\u8FD4\u56DE\u9876\u90E8","back_button":"\u8FD4\u56DE\u6309\u94AE","back_to_bookshelf":"\u8FD4\u56DE\u4E66\u67B6","book_shelf":"\u4E66\u5E93","child_book_hint":"\u6587\u4EF6\u5939\u5185\u6709{0}\u672C\u4E66","debug_mode":"\u8C03\u8BD5\u6A21\u5F0F","do_you_reset_all_settings":"\u662F\u5426\u8981\u91CD\u7F6E\u4E3A\u9ED8\u8BA4\u8BBE\u7F6E\uFF1F","double_page_mode":"\u53CC\u9875\u6A21\u5F0F","double_page_width":"\u6A2A\u5C4F\u53CC\u9875\u5BBD\u5EA6:","download_sample_config_file":"\u4E0B\u8F7D\u793A\u4F8B\u914D\u7F6E\u6587\u4EF6","download_windows_reg_file":"\u4E0B\u8F7D\u53F3\u952E\u6CE8\u518C\u6587\u4EF6","drop_to_upload":"\u70B9\u51FB\u6216\u5C06\u6587\u4EF6\u62D6\u52A8\u5230\u6B64\u533A\u57DF\u4EE5\u4E0A\u4F20","energy_threshold":"\u5207\u8FB9\u5F3A\u5EA6\uFF1A","epub_info":"ePub \u4FE1\u606F","exit_fullscreen":"\u9000\u51FA\u5168\u5C4F\u6A21\u5F0F","filesize":"\u5927\u5C0F\uFF1A{0}","flip_mode":"\u7FFB\u9875\u6A21\u5F0F","flip_odd_even_page":"\u66F4\u6539\u8DE8\u9875\u5339\u914D","flip_odd_even_page_hint":"\u5982\u679C\u8DE8\u9875\u5185\u5BB9\u4E0D\u5339\u914D\uFF0C\u53EF\u4EE5\u5C1D\u8BD5\u70B9\u51FB\u4FEE\u6B63","found_read_history":"\u53D1\u73B0\u672C\u5730\u9605\u8BFB\u8BB0\u5F55","from_interrupt":"\u4ECE\u4E2D\u65AD\u5904\u7EE7\u7EED","full_screen_hint":"\u5168\u5C4F\u6309\u94AE","fullscreen":"\u5207\u6362\u5168\u5C4F","good_job_and_byebye":"\u505A\u5F97\u4E0D\u9519\uFF0C\u518D\u89C1~","gray_image":"\u9ED1\u767D\u5316","hint":"\u63D0\u793A","hint_first_page":"\u5F53\u524D\u662F\u7B2C\u4E00\u9875\uFF0C\u65E0\u6CD5\u5411\u524D\u7FFB\u9875","hint_last_page":"\u5F53\u524D\u662F\u6700\u540E\u4E00\u9875\uFF0C\u65E0\u6CD5\u5411\u540E\u7FFB\u9875","hour":"\u5C0F\u65F6","image_width_limit":"\u9650\u5236\u5BBD\u5EA6","infinite_dropdown":"\u65E0\u9650\u4E0B\u62C9\u6A21\u5F0F","interval":"\u95F4\u9694:","left_screen_to_next":"\u5DE6\u5F00\u672C\uFF08\u7F8E\u6F2B\uFF09","load_all_pages":"\u52A0\u8F7D\u6240\u6709\u9875\u9762","load_from_interrupt":"\u662F\u5426\u4ECE\u7B2CXX\u9875\u5F00\u59CB\u52A0\u8F7D\uFF1F","login_success_hint":"\u767B\u5F55\u6210\u529F\uFF0C\u8FD4\u56DE\u4E0A\u4E00\u9875\u9762","logout":"\u9000\u51FA\u767B\u5F55","margin_bottom_on_scroll_mode":"\u9875\u9762\u95F4\u8DDD:","margin_on_scroll_mode":"\u9875\u9762\u95F4\u9699:","max_width":"\u6700\u5927\u5BBD\u5EA6\uFF1A","minute":"\u5206","network":"\u7F51\u7EDC","no_book_found_hint":"\u672A\u627E\u5230\u4E66\u7C4D\uFF0C\u8BD5\u8BD5\u4E0A\u4F20\u6587\u4EF6\uFF1F","no_support_upload_file":"\u6587\u4EF6\u4E0A\u4F20\u529F\u80FD\u5DF2\u88AB\u7BA1\u7406\u5458\u5173\u95ED","not_support_fullscreen":"\u6B64\u6D4F\u89C8\u5668\u4E0D\u652F\u6301\u5168\u5C4F\u6A21\u5F0F","now_is":"\u5F53\u524D:","number_of_online_books":"\u5728\u7EBF\u4E66\u7C4D\u6570\u91CF\uFF1A","original_image":"\u663E\u793A\u539F\u56FE","original_pdf_link":"\u67E5\u770B\u539F\u59CBPDF","page":"\u9875","page_turning_seconds":"\u7FFB\u9875\u95F4\u9694:","pagination_mode":"\u5206\u9875\u52A0\u8F7D\u6A21\u5F0F","pdf_hint_message":"\u652F\u6301\u7EAF\u56FE\u7247PDF\uFF0C\u5982\u679C\u52A0\u8F7D\u7F13\u6162\u6216\u51FA\u9519\uFF0C\u8BF7\u5C1D\u8BD5\uFF1A","please_enable_upload":"\u8BF7\u542F\u7528\u670D\u52A1\u5668\u7684\u4E0A\u4F20\u529F\u80FD","please_enter_content":"\u8BF7\u8F93\u5165\u5185\u5BB9","qrcode_hint":"\u626B\u7801\u9605\u8BFB\uFF0C\u70B9\u51FB\u663E\u793A\u4E8C\u7EF4\u7801","raw_resolution":"\u539F\u59CB\u5206\u8FA8\u7387","re_sort_book":"\u91CD\u65B0\u6392\u5217\u4E66\u7C4D","re_sort_page":"\u91CD\u65B0\u6392\u5E8F\u9875\u9762","reader_settings":"\u9605\u8BFB\u5668\u8BBE\u7F6E","reading_progress_bar":"\u9605\u8BFB\u8FDB\u5EA6\u6761","refresh_page":"\u5237\u65B0\u9875\u9762","reset_all_settings":"\u91CD\u7F6E\u8BBE\u7F6E","resort_file":"\u91CD\u65B0\u6392\u5E8F\u6587\u4EF6","right_screen_to_next":"\u53F3\u5F00\u672C\uFF08\u65E5\u6F2B\uFF09","save_page_num":"\u4FDD\u5B58\u9605\u8BFB\u8FDB\u5EA6","scan_qrcode":"\u626B\u7801\u9605\u8BFB\uFF1A","scanned_hint":"\u626B\u63CF\u5230XX\u672C\u4E66\uFF0C\u7ACB\u5373\u67E5\u770B\uFF1F","scroll_mode":"\u5377\u8F74\u6A21\u5F0F","second":"\u79D2","select_language":"\u9009\u62E9\u8BED\u8A00","server_config":"\u670D\u52A1\u5668\u8BBE\u7F6E","server_setting":"Comigo \u670D\u52A1\u5668\u8BBE\u7F6E","set_back_color":"\u80CC\u666F\u989C\u8272:","set_interface_color":"\u754C\u9762\u989C\u8272:","show_book_titles":"\u663E\u793A\u4E66\u540D","show_file_icon":"\u663E\u793A\u6587\u4EF6\u56FE\u6807","show_header":"\u663E\u793A\u6807\u9898","show_page_num":"\u663E\u793A\u9875\u7801","simplify_book_titles":"\u7B80\u5316\u4E66\u540D","single_page_mode":"\u5355\u9875\u6A21\u5F0F","single_page_width":"\u6A2A\u5C4F\u5355\u9875\u5BBD\u5EA6:","sort_by_default":"\u4FDD\u6301\u9ED8\u8BA4\u987A\u5E8F","sort_by_filename":"\u6309\u6587\u4EF6\u540D\u6392\u5E8F (A-Z)","sort_by_filename_reverse":"\u6309\u6587\u4EF6\u540D\u6392\u5E8F (Z-A)","sort_by_filesize":"\u6309\u6587\u4EF6\u5927\u5C0F\u6392\u5E8F (\u4ECE\u5927\u5230\u5C0F)","sort_by_filesize_reverse":"\u6309\u6587\u4EF6\u5927\u5C0F\u6392\u5E8F (\u4ECE\u5C0F\u5230\u5927)","sort_by_modify_time":"\u6309\u4FEE\u6539\u65F6\u95F4\u6392\u5E8F (\u4ECE\u65B0\u5230\u65E7)","sort_by_modify_time_reverse":"\u6309\u4FEE\u6539\u65F6\u95F4\u6392\u5E8F (\u4ECE\u65E7\u5230\u65B0)","sort_reverse":"\uFF08\u53CD\u5411\uFF09","start_sketch_message":"\u5012\u8BA1\u65F6\u901F\u5199\u5DF2\u5F00\u59CB\uFF0C\u795D\u4F60\u5FC3\u60C5\u6109\u5FEB\u3002","start_sketch_mode":"\u5F00\u59CB\u901F\u5199","starting_from_beginning":"\u4ECE\u5934\u5F00\u59CB","starting_from_beginning_hint":"\u4ECE\u5934\u5F00\u59CB\u52A0\u8F7D","stop_sketch_mode":"\u505C\u6B62\u901F\u5199","submit":"\u63D0\u4EA4","success_fullscreen":"\u5DF2\u8FDB\u5165\u5168\u5C4F\u6A21\u5F0F","successfully_loaded_reading_progress":"\u6210\u529F\u52A0\u8F7D\u9605\u8BFB\u8FDB\u5EA6","switch_to_flip_mode":"\u5207\u6362\u5230\u7FFB\u9875\u6A21\u5F0F","switch_to_scrolling_mode":"\u5207\u6362\u5230\u5377\u8F74\u6A21\u5F0F","sync_page":"\u8FDC\u7A0B\u540C\u6B65\u7FFB\u9875","temp_future_hint":"\u4E34\u65F6\u653E\u4E00\u4E9B\u8FD8\u672A\u5B8C\u6210\u7684\u529F\u80FD\uFF0C\u5F00\u53D1\u4E0E\u8C03\u6574\u4E2D\u3002","test":"\u6D4B\u8BD5","to_flip_mode":"\u5207\u6362\u5230\u7FFB\u9875\u6A21\u5F0F","to_infinite_dropdown_mode":"\u5207\u6362\u5230\u65E0\u9650\u4E0B\u62C9\u6A21\u5F0F","to_pagination_mode":"\u5207\u6362\u5230\u5206\u9875\u6A21\u5F0F","to_scroll_mode":"\u5207\u6362\u5230\u5377\u8F74\u6A21\u5F0F","total_is":"\u5B8C\u6210:","total_time":"\u603B\u65F6\u95F4:","type_or_paste_content":"\u952E\u5165\u6216\u7C98\u8D34\u5185\u5BB9","upload_file":"\u4E0A\u4F20\u6587\u4EF6","uploaded_folder_hint":"\u6587\u4EF6\u5C06\u4E0A\u4F20\u5230\u9ED8\u8BA4\u4E66\u5E93\u76EE\u5F55\u4E0B\u7684 upload \u6587\u4EF6\u5939\u3002","width_use_fixed_value":"\u6A2A\u5C4F\u5BBD\u5EA6: \u56FA\u5B9A\u503Cpx","width_use_percent":"\u6A2A\u5C4F\u5BBD\u5EA6: \u767E\u5206\u6BD4%"}');


var $3849b1934e918b8a$exports = {};
$3849b1934e918b8a$exports = JSON.parse('{"CachePath":"\u30ED\u30FC\u30AB\u30EB\u30AD\u30E3\u30C3\u30B7\u30E5\u306E\u5834\u6240","CachePath_Description":"\u30ED\u30FC\u30AB\u30EB\u306E\u753B\u50CF\u30AD\u30E3\u30C3\u30B7\u30E5\u306E\u5834\u6240\u3001\u30C7\u30D5\u30A9\u30EB\u30C8\u306E\u30B7\u30B9\u30C6\u30E0\u4E00\u6642\u30D5\u30A9\u30EB\u30C0\u30FC\u3002","CertFile":"\u8A3C\u660E\u66F8\u30D5\u30A1\u30A4\u30EB","CertFile_Description":"TLS/SSL \u8A3C\u660E\u66F8\u30D5\u30A1\u30A4\u30EB\u306E\u30D1\u30B9 (\u30C7\u30D5\u30A9\u30EB\u30C8: \\"~/.config/.comigo/cert.crt\\")","ClearCacheExit":"\u7D42\u4E86\u6642\u306B\u30AF\u30EA\u30FC\u30F3\u30A2\u30C3\u30D7","ClearCacheExit_Description":"\u30D7\u30ED\u30B0\u30E9\u30E0\u3092\u7D42\u4E86\u3059\u308B\u3068\u304D\u306F\u3001Web \u753B\u50CF\u30AD\u30E3\u30C3\u30B7\u30E5\u3092\u30AF\u30EA\u30A2\u3057\u307E\u3059\u3002","ClearDatabaseWhenExit":"\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u30D6\u30C3\u30AF\u3092\u30AF\u30EA\u30A2\u3059\u308B","ClearDatabaseWhenExit_Description":"\u30ED\u30FC\u30AB\u30EB \u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u304C\u6709\u52B9\u306B\u306A\u3063\u3066\u3044\u308B\u5834\u5408\u3001\u30B9\u30AD\u30E3\u30F3\u306E\u5B8C\u4E86\u5F8C\u306B\u5B58\u5728\u3057\u306A\u3044\u66F8\u7C4D\u306F\u6D88\u53BB\u3055\u308C\u307E\u3059\u3002","ConfigManager":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u7BA1\u7406","ConfigManagerDeleteSuccess":"\u8A2D\u5B9A\u304C\u524A\u9664\u3055\u308C\u307E\u3057\u305F\u3002","ConfigManagerDescription":"Save\u3092\u30AF\u30EA\u30C3\u30AF\u3059\u308B\u3068\u3001\u73FE\u5728\u306E\u8A2D\u5B9A\u304C\u30B5\u30FC\u30D0\u30FC\u306B\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u3055\u308C\u3001\u65E2\u5B58\u306E\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u304C\u4E0A\u66F8\u304D\u3055\u308C\u307E\u3059\u3002","ConfigManagerSaveHint":"\u3059\u3067\u306B\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u304C\u3042\u308B\u306E\u3067\u3001\u4FDD\u5B58\u5834\u6240\u3092\u5909\u66F4\u3057\u3066\u304F\u3060\u3055\u3044\u3002","ConfigManagerSaveSuccess":"\u8A2D\u5B9A\u304C\u4FDD\u5B58\u3055\u308C\u307E\u3057\u305F\u3002","ConfigSaveTo":"Config File Location","Debug":"\u30C7\u30D0\u30C3\u30B0\u30E2\u30FC\u30C9\u3092\u30AA\u30F3\u306B\u3059\u308B","Debug_Description":"\u30C7\u30D0\u30C3\u30B0\u6A5F\u80FD\u3092\u6709\u52B9\u306B\u3059\u308B","DisableLAN":"LAN\u5171\u6709\u3092\u7121\u52B9\u306B\u3059\u308B","DisableLAN_Description":"\u8AAD\u307F\u53D6\u308A\u30B5\u30FC\u30D3\u30B9\u306F\u3053\u306E\u30DE\u30B7\u30F3\u3067\u306E\u307F\u63D0\u4F9B\u3055\u308C\u3001\u5916\u90E8\u306B\u306F\u5171\u6709\u3055\u308C\u307E\u305B\u3093\u3002\u30DB\u30C3\u30C8\u30EA\u30ED\u30FC\u30C9\u3092\u30B5\u30DD\u30FC\u30C8\u3057\u3066\u3044\u307E\u305B\u3093\u3002","EnableDatabase":"\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u3092\u6709\u52B9\u306B\u3059\u308B","EnableDatabase_Description":"\u30ED\u30FC\u30AB\u30EB\u30C7\u30FC\u30BF\u30D9\u30FC\u30B9\u306B\u30B9\u30AD\u30E3\u30F3\u3057\u305F\u66F8\u7C4D\u30C7\u30FC\u30BF\u3092\u4FDD\u5B58\u3067\u304D\u308B\u3088\u3046\u306B\u3057\u307E\u3059\u3002\u3053\u306E\u69CB\u6210\u306F\u30DB\u30C3\u30C8 \u30EA\u30ED\u30FC\u30C9\u3092\u30B5\u30DD\u30FC\u30C8\u3057\u307E\u305B\u3093\u3002","EnableFrpcServer":"Frpc\u30B5\u30FC\u30D0\u30FC\u3092\u6709\u52B9\u306B\u3059\u308B","EnableLogin":"Enable Login","EnableLogin_Description":"\u30ED\u30B0\u30A4\u30F3\u3092\u6709\u52B9\u306B\u3059\u308B\u304B\u3069\u3046\u304B\u3002\u30C7\u30D5\u30A9\u30EB\u30C8\u3067\u306F\u30ED\u30B0\u30A4\u30F3\u306F\u5FC5\u8981\u3042\u308A\u307E\u305B\u3093\u3002\u30DB\u30C3\u30C8\u30EA\u30ED\u30FC\u30C9\u3092\u30B5\u30DD\u30FC\u30C8\u3057\u307E\u305B\u3093\u3002","EnableTLS":"Enable TLS","EnableTLS_Description":"HTTPS \u30D7\u30ED\u30C8\u30B3\u30EB\u3092\u6709\u52B9\u306B\u3059\u308B\u304B\u3069\u3046\u304B\u3002\u8A3C\u660E\u66F8\u306F\u30AD\u30FC\u30D5\u30A1\u30A4\u30EB\u306B\u8A2D\u5B9A\u3059\u308B\u5FC5\u8981\u304C\u3042\u308A\u307E\u3059\u3002","EnableUpload":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u6A5F\u80FD\u3092\u6709\u52B9\u306B\u3059\u308B","EnableUpload_Description":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u6A5F\u80FD\u3092\u6709\u52B9\u306B\u3057\u307E\u3059\u3002","ExcludePath":"\u30D1\u30B9\u3092\u9664\u5916\u3059\u308B","ExcludePath_Description":"\u672C\u3092\u30B9\u30AD\u30E3\u30F3\u3059\u308B\u3068\u304D\u306B\u3001\u9664\u5916\u3059\u308B\u5FC5\u8981\u304C\u3042\u308B\u30D5\u30A1\u30A4\u30EB\u307E\u305F\u306F\u30D5\u30A9\u30EB\u30C0\u30FC\u306E\u540D\u524D","FrpClientConfig":"FrpClient\u8A2D\u5B9A","GenerateBookMetadata":"\u66F8\u7C4D\u306E\u30E1\u30BF\u30C7\u30FC\u30BF\u3092\u751F\u6210\u3059\u308B","GenerateMetaData":"\u30E1\u30BF\u30C7\u30FC\u30BF\u306E\u751F\u6210","GenerateMetaData_Description":"\u66F8\u7C4D\u306E\u30E1\u30BF\u30C7\u30FC\u30BF\u3092\u751F\u6210\u3057\u307E\u3059\u3002\u73FE\u5728\u306F\u6709\u52B9\u3067\u306F\u3042\u308A\u307E\u305B\u3093\u3002","HomeDirectory":"\u30DB\u30FC\u30E0\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA","Host":"\u30C9\u30E1\u30A4\u30F3\u540D","Host_Description":"QR\u30B3\u30FC\u30C9\u3067\u8868\u793A\u3055\u308C\u308B\u30DB\u30B9\u30C8\u540D\u3092\u30AB\u30B9\u30BF\u30DE\u30A4\u30BA\u3057\u307E\u3059\u3002\\n\u30C7\u30D5\u30A9\u30EB\u30C8\u306F\u30CD\u30C3\u30C8\u30EF\u30FC\u30AF \u30AB\u30FC\u30C9\u306E IP \u3067\u3059\u3002","KeyFile":"\u30AD\u30FC\u30D5\u30A1\u30A4\u30EB","KeyFile_Description":"TLS/SSL \u30AD\u30FC \u30D5\u30A1\u30A4\u30EB \u30D1\u30B9 (\u30C7\u30D5\u30A9\u30EB\u30C8: \\"~/.config/.comigo/key.key\\")","LocalStores":"\u30E9\u30A4\u30D6\u30E9\u30EA\u30D5\u30A9\u30EB\u30C0\u30FC","LocalStores_Description":"\u7D76\u5BFE\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3068\u76F8\u5BFE\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u304C\u30B5\u30DD\u30FC\u30C8\u3055\u308C\u3066\u3044\u308B\u3002 \u76F8\u5BFE\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306F\u5B9F\u884C\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u3092\u57FA\u6E96\u3068\u3057\u307E\u3059\u3002","LogFileName":"\u30ED\u30B0\u30D5\u30A1\u30A4\u30EB\u540D","LogFileName_Description":"\u30ED\u30B0\u30D5\u30A1\u30A4\u30EB\u540D","LogFilePath":"\u30ED\u30B0\u306E\u4FDD\u5B58\u5834\u6240","LogFilePath_Description":"\u30ED\u30B0\u30D5\u30A1\u30A4\u30EB\u306E\u4FDD\u5B58\u5834\u6240","LogToFile":"\u30ED\u30B0\u3092\u30ED\u30FC\u30AB\u30EB\u306B\u8A18\u9332\u3059\u308B","LogToFile_Description":"\u30D7\u30ED\u30B0\u30E9\u30E0\u30ED\u30B0\u3092\u30ED\u30FC\u30AB\u30EB\u30D5\u30A1\u30A4\u30EB\u306B\u4FDD\u5B58\u3059\u308B\u304B\u3069\u3046\u304B\u3002\u30C7\u30D5\u30A9\u30EB\u30C8\u3067\u306F\u4FDD\u5B58\u3055\u308C\u307E\u305B\u3093\u3002","MaxScanDepth":"\u6700\u5927\u30B9\u30AD\u30E3\u30F3\u6DF1\u5EA6","MaxScanDepth_Description":"\u6700\u5927\u30B9\u30AD\u30E3\u30F3\u6DF1\u5EA6\u3002\\n\u6DF1\u3055\u3092\u8D85\u3048\u308B\u30D5\u30A1\u30A4\u30EB\u306F\u30B9\u30AD\u30E3\u30F3\u3055\u308C\u307E\u305B\u3093\u3002\\n\u73FE\u5728\u306E\u5B9F\u884C\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u304C\u30D9\u30FC\u30B9\u3068\u306A\u308A\u307E\u3059\u3002","MinImageNum":"\u6700\u4F4E\u679A\u6570","MinImageNum_Description":"\u672C\u3068\u307F\u306A\u3055\u308C\u308B\u306B\u306F\u3001\u5C11\u306A\u304F\u3068\u3082\u6570\u679A\u306E\u753B\u50CF\u304C\u542B\u307E\u308C\u3066\u3044\u308B\u5FC5\u8981\u304C\u3042\u308A\u307E\u3059\u3002","OpenBrowser":"\u30D6\u30E9\u30A6\u30B6\u3092\u958B\u304F","OpenBrowser_Description":"windows\u306E\u30C7\u30D5\u30A9\u30EB\u30C8\u306Ftrue\u3001\u305D\u306E\u4ED6\u306E\u30D7\u30E9\u30C3\u30C8\u30D5\u30A9\u30FC\u30E0\u306E\u30C7\u30D5\u30A9\u30EB\u30C8\u306Ffalse\u3002","Password":"\u30D1\u30B9\u30EF\u30FC\u30C9","Password_Description":"\u30ED\u30B0\u30A4\u30F3\u3092\u6709\u52B9\u306B\u3057\u305F\u5F8C\u3001\u30ED\u30B0\u30A4\u30F3\u306B\u4F7F\u7528\u3055\u308C\u308B\u30D1\u30B9\u30EF\u30FC\u30C9\u3002","Port":"\u30DD\u30FC\u30C8","Port_Description":"Web \u30B5\u30FC\u30D3\u30B9 \u30DD\u30FC\u30C8\u3002\u3053\u306E\u69CB\u6210\u306F\u30DB\u30C3\u30C8 \u30EA\u30ED\u30FC\u30C9\u3092\u30B5\u30DD\u30FC\u30C8\u3057\u307E\u305B\u3093\u3002","PrintAllPossibleQRCode":"\u305D\u306E\u4ED6\u306E QR \u30B3\u30FC\u30C9","ProgramDirectory":"\u30D7\u30ED\u30B0\u30E9\u30E0\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA","StartFrpClientInBackground":"Frp\u30AF\u30E9\u30A4\u30A2\u30F3\u30C8\u3092\u958B\u59CB\u3059\u308B","SupportFileType":"\u30B5\u30DD\u30FC\u30C8\u3055\u308C\u3066\u3044\u308B\u5727\u7E2E\u30D1\u30C3\u30B1\u30FC\u30B8","SupportFileType_Description":"\u30D5\u30A1\u30A4\u30EB\u3092\u30B9\u30AD\u30E3\u30F3\u3059\u308B\u3068\u304D\u306B\u3001\u305D\u308C\u3092\u30B9\u30AD\u30C3\u30D7\u3059\u308B\u304B\u3001\u30D6\u30C3\u30AF\u51E6\u7406\u306E\u30D5\u30A1\u30A4\u30EB\u306E\u63A5\u5C3E\u8F9E\u3068\u3057\u3066\u30AB\u30A6\u30F3\u30C8\u3059\u308B\u304B\u3092\u6C7A\u5B9A\u3059\u308B\u305F\u3081\u306B\u4F7F\u7528\u3055\u308C\u307E\u3059\u3002","SupportMediaType":"\u30B5\u30DD\u30FC\u30C8\u3055\u308C\u3066\u3044\u308B\u753B\u50CF\u30D5\u30A1\u30A4\u30EB","SupportMediaType_Description":"\u5727\u7E2E\u30D1\u30C3\u30B1\u30FC\u30B8\u3092\u30B9\u30AD\u30E3\u30F3\u3059\u308B\u3068\u304D\u306B\u753B\u50CF\u306E\u6570\u3092\u30AB\u30A6\u30F3\u30C8\u3059\u308B\u305F\u3081\u306B\u4F7F\u7528\u3055\u308C\u308B\u753B\u50CF\u30D5\u30A1\u30A4\u30EB\u306E\u63A5\u5C3E\u8F9E","Timeout":"\u6709\u52B9\u671F\u9650","TimeoutLimitForScan":"\u30B9\u30AD\u30E3\u30F3\u30BF\u30A4\u30E0\u30A2\u30A6\u30C8","TimeoutLimitForScan_Description":"\u30D5\u30A1\u30A4\u30EB\u306E\u30B9\u30AD\u30E3\u30F3\u6642\u306B\u6570\u79D2\u4EE5\u4E0A\u304B\u304B\u308B\u5834\u5408\u3001\u5927\u304D\u3059\u304E\u308B\u30D5\u30A1\u30A4\u30EB\u3067\u30B9\u30BF\u30C3\u30AF\u3059\u308B\u3053\u3068\u3092\u907F\u3051\u308B\u305F\u3081\u306B\u30D5\u30A1\u30A4\u30EB\u306E\u30B9\u30AD\u30E3\u30F3\u3092\u4E2D\u6B62\u3057\u307E\u3059\u3002","Timeout_Description":"\u30ED\u30B0\u30A4\u30F3\u3092\u6709\u52B9\u306B\u3057\u305F\u5F8C\u306E Cookie \u306E\u6709\u52B9\u671F\u9650\u3002\u5358\u4F4D\u306F\u5206\u3067\u3059\u3002\u30C7\u30D5\u30A9\u30EB\u30C8\u3067\u306F 180 \u5206\u3067\u671F\u9650\u5207\u308C\u306B\u306A\u308A\u307E\u3059\u3002","UploadPath":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u5834\u6240","UploadPath_Description":"\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u3055\u308C\u305F\u30D5\u30A1\u30A4\u30EB\u306E\u4FDD\u5B58\u5834\u6240\u3092\u30AB\u30B9\u30BF\u30DE\u30A4\u30BA\u3057\u307E\u3059\u3002\u30C7\u30D5\u30A9\u30EB\u30C8\u3067\u306F\u3001\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u30D5\u30A9\u30EB\u30C0\u306F\u73FE\u5728\u306E\u5B9F\u884C\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u307E\u305F\u306F\u6700\u521D\u306E\u30D6\u30C3\u30AF\u30B9\u30C8\u30A2\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u306E\u4E0B\u306B\u4F5C\u6210\u3055\u308C\u307E\u3059\u3002","UseCache":"\u30ED\u30FC\u30AB\u30EB\u753B\u50CF\u30AD\u30E3\u30C3\u30B7\u30E5","Username":"\u30E6\u30FC\u30B6\u30FC\u540D","Username_Description":"\u30ED\u30B0\u30A4\u30F3\u3092\u6709\u52B9\u306B\u3057\u305F\u5F8C\u3001\u30ED\u30B0\u30A4\u30F3 \u30A4\u30F3\u30BF\u30FC\u30D5\u30A7\u30A4\u30B9\u306B\u5FC5\u8981\u306A\u30E6\u30FC\u30B6\u30FC\u540D\u3002","WorkingDirectory":"\u4F5C\u696D\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA","ZipFileTextEncoding":"UTF-8\u3067\u306F\u3042\u308A\u307E\u305B\u3093","ZipFileTextEncoding_Description":"utf-8 \u4EE5\u5916\u3067\u30A8\u30F3\u30B3\u30FC\u30C9\u3055\u308C\u305F ZIP \u30D5\u30A1\u30A4\u30EB\u3002\u89E3\u6790\u3059\u308B\u306B\u306F\u3069\u306E\u30A8\u30F3\u30B3\u30FC\u30C9\u3092\u4F7F\u7528\u3059\u308B\u5FC5\u8981\u304C\u3042\u308A\u307E\u3059\u3002\u30C7\u30D5\u30A9\u30EB\u30C8\u306EGBK\u3002","abs":"\u30DC","all_page_num":"\u7DCF\u30DA\u30FC\u30B8\u6570: {0}","author":"\u4F5C\u8005: {0}","auto_crop":"\u81EA\u52D5\u30A8\u30C3\u30B8\u30C8\u30EA\u30DF\u30F3\u30B0","auto_double_page":"\u81EA\u52D5\u898B\u958B\u304D\u30DA\u30FC\u30B8 (\u03B2)","auto_hide_toolbar":"\u30C4\u30FC\u30EB\u30D0\u30FC\u3092\u81EA\u52D5\u975E\u8868\u793A","back-to-top":"\u30DA\u30FC\u30B8\u30C8\u30C3\u30D7\u3078\u623B\u308B","back_button":"\u623B\u308B\u30DC\u30BF\u30F3","back_to_bookshelf":"\u672C\u68DA\u306B\u623B\u308B","book_shelf":"\u672C\u68DA","child_book_hint":"\u30D5\u30A9\u30EB\u30C0\u5185\u306B{0}\u518A\u306E\u66F8\u7C4D\u304C\u542B\u307E\u308C\u3066\u3044\u307E\u3059","debug_mode":"\u30C7\u30D0\u30C3\u30B0\u30E2\u30FC\u30C9","do_you_reset_all_settings":"\u3059\u3079\u3066\u306E\u8A2D\u5B9A\u3092\u521D\u671F\u5316\u3057\u307E\u3059\u304B\uFF1F","double_page_mode":"\u30C0\u30D6\u30EB\u30DA\u30FC\u30B8\u30E2\u30FC\u30C9","double_page_width":"\u898B\u958B\u304D\u30DA\u30FC\u30B8\u5E45:","download_sample_config_file":"\u8A2D\u5B9A\u30D5\u30A1\u30A4\u30EB\u3092\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","download_windows_reg_file":"REG\u30D5\u30A1\u30A4\u30EB\u3092\u30C0\u30A6\u30F3\u30ED\u30FC\u30C9","drop_to_upload":"\u30D5\u30A1\u30A4\u30EB\u3092\u30AF\u30EA\u30C3\u30AF\u3059\u308B\u304B\u3001\u3053\u306E\u30A8\u30EA\u30A2\u306B\u30C9\u30E9\u30C3\u30B0\uFF06\u30C9\u30ED\u30C3\u30D7\u3057\u3066\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9","energy_threshold":"\u30A8\u30CD\u30EB\u30AE\u30FC\u95BE\u5024:","epub_info":"ePub\u60C5\u5831","exit_fullscreen":"\u30D5\u30EB\u30B9\u30AF\u30EA\u30FC\u30F3\u3092\u7D42\u4E86","filesize":"\u30D5\u30A1\u30A4\u30EB\u30B5\u30A4\u30BA: {0}","flip_mode":"\u30E8\u30B3\u8AAD\u307F\u30E2\u30FC\u30C9","flip_odd_even_page":"\u30DA\u30FC\u30B8\u9593\u306E\u30DE\u30C3\u30C1\u30F3\u30B0\u3092\u5909\u66F4","flip_odd_even_page_hint":"\u30DA\u30FC\u30B8\u5185\u5BB9\u304C\u4E00\u81F4\u3057\u306A\u3044\u5834\u5408\u306F\u3001\u30AF\u30EA\u30C3\u30AF\u3057\u3066\u4FEE\u6B63\u3092\u8A66\u307F\u3066\u304F\u3060\u3055\u3044","found_read_history":"\u8AAD\u66F8\u5C65\u6B74\u304C\u898B\u3064\u304B\u308A\u307E\u3057\u305F","from_interrupt":"\u9014\u4E2D\u304B\u3089\u518D\u958B","full_screen_hint":"\u30D5\u30EB\u30B9\u30AF\u30EA\u30FC\u30F3\u30DC\u30BF\u30F3","fullscreen":"\u30D5\u30EB\u30B9\u30AF\u30EA\u30FC\u30F3","good_job_and_byebye":"\u304A\u75B2\u308C\u69D8\u3067\u3057\u305F\u3002\u3055\u3088\u3046\u306A\u3089\u3002","gray_image":"\u30B0\u30EC\u30FC\u30B9\u30B1\u30FC\u30EB\u5316","hint":"Hint","hint_first_page":"\u6700\u521D\u306E\u30DA\u30FC\u30B8\u3067\u306F\u9032\u3081\u307E\u305B\u3093","hint_last_page":"\u3053\u308C\u304C\u6700\u5F8C\u306E\u30DA\u30FC\u30B8\u3067\u3059","hour":"\u6642\u9593","image_width_limit":"\u5E45\u306E\u30B5\u30A4\u30BA\u5236\u9650","infinite_dropdown":"\u7121\u9650\u30C9\u30ED\u30C3\u30D7\u30C0\u30A6\u30F3\u30E2\u30FC\u30C9","interval":"\u9593\u9694:","left_screen_to_next":"\u5DE6\u5411\u304D: \u30B3\u30DF\u30C3\u30AF","load_all_pages":"\u3059\u3079\u3066\u306E\u30DA\u30FC\u30B8\u3092\u8AAD\u307F\u8FBC\u3080","load_from_interrupt":"XX\u30DA\u30FC\u30B8\u304B\u3089\u8AAD\u307F\u8FBC\u307F\u3092\u958B\u59CB\u3057\u307E\u3059\u304B\uFF1F","login_success_hint":"\u30ED\u30B0\u30A4\u30F3\u306B\u6210\u529F\u3057\u307E\u3057\u305F\u3002\u524D\u306E\u30DA\u30FC\u30B8\u306B\u623B\u308A\u307E\u3059","logout":"\u30ED\u30B0\u30A2\u30A6\u30C8","margin_bottom_on_scroll_mode":"\u30DA\u30FC\u30B8\u9593\u306E\u4F59\u767D:","margin_on_scroll_mode":"\u30DA\u30FC\u30B8\u30AE\u30E3\u30C3\u30D7:","max_width":"\u6700\u5927\u5E45:","minute":"\u5206","network":"\u30C3\u30C8\u30EF\u30FC\u30AF","no_book_found_hint":"\u66F8\u7C4D\u304C\u898B\u3064\u304B\u308A\u307E\u305B\u3093\u3002\u30D5\u30A1\u30A4\u30EB\u3092\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u3057\u3066\u307F\u3066\u304F\u3060\u3055\u3044\u3002","no_support_upload_file":"\u30D5\u30A1\u30A4\u30EB\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u6A5F\u80FD\u306F\u7BA1\u7406\u8005\u306B\u3088\u308A\u7121\u52B9\u5316\u3055\u308C\u3066\u3044\u307E\u3059","not_support_fullscreen":"\u3053\u306E\u30D6\u30E9\u30A6\u30B6\u306F\u30D5\u30EB\u30B9\u30AF\u30EA\u30FC\u30F3\u3092\u30B5\u30DD\u30FC\u30C8\u3057\u3066\u3044\u307E\u305B\u3093","now_is":"\u73FE\u5728:","number_of_online_books":"\u30AA\u30F3\u30E9\u30A4\u30F3\u66F8\u7C4D\u6570\uFF1A","original_image":"\u30AA\u30EA\u30B8\u30CA\u30EB\u89E3\u50CF\u5EA6","original_pdf_link":"PDF\u30D5\u30A1\u30A4\u30EB\u306E\u30EA\u30F3\u30AF","page":"\u30DA\u30FC\u30B8","page_turning_seconds":"\u30DA\u30FC\u30B8\u3081\u304F\u308A\u6642\u9593:","pagination_mode":"\u30DA\u30FC\u30B8\u30F3\u30B0\u30E2\u30FC\u30C9","pdf_hint_message":"\u7D14\u7C8B\u306A\u753B\u50CFPDF\u306B\u5BFE\u5FDC\u3057\u3066\u3044\u307E\u3059\u3002\u8AAD\u307F\u8FBC\u307F\u304C\u9045\u3044\u5834\u5408\u3084\u30A8\u30E9\u30FC\u304C\u767A\u751F\u3057\u305F\u5834\u5408\u306F\u3001\u4EE5\u4E0B\u3092\u8A66\u3057\u3066\u304F\u3060\u3055\u3044\uFF1A","please_enable_upload":"\u30B5\u30FC\u30D0\u30FC\u306E\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u6A5F\u80FD\u3092\u6709\u52B9\u306B\u3057\u3066\u304F\u3060\u3055\u3044","please_enter_content":"\u5185\u5BB9\u3092\u5165\u529B\u3057\u3066\u304F\u3060\u3055\u3044","qrcode_hint":"\u30AF\u30EA\u30C3\u30AF\u3057\u3066QR\u30B3\u30FC\u30C9\u3092\u8868\u793A","raw_resolution":"\u5143\u306E\u89E3\u50CF\u5EA6","re_sort_book":"\u66F8\u7C4D\u3092\u4E26\u3079\u66FF\u3048","re_sort_page":"\u30DA\u30FC\u30B8\u3092\u4E26\u3079\u66FF\u3048","reader_settings":"\u30EA\u30FC\u30C0\u30FC\u8A2D\u5B9A","reading_progress_bar":"\u8AAD\u66F8\u9032\u6357\u30D0\u30FC","refresh_page":"\u30DA\u30FC\u30B8\u3092\u66F4\u65B0","reset_all_settings":"\u3059\u3079\u3066\u306E\u8A2D\u5B9A\u3092\u30EA\u30BB\u30C3\u30C8","resort_file":"\u30D5\u30A1\u30A4\u30EB\u3092\u518D\u4E26\u3073\u66FF\u3048","right_screen_to_next":"\u53F3\u5411\u304D: \u30DE\u30F3\u30AC","save_page_num":"\u8AAD\u66F8\u9032\u6357\u3092\u4FDD\u5B58","scan_qrcode":"QR\u30B3\u30FC\u30C9\u3092\u30B9\u30AD\u30E3\u30F3:","scanned_hint":"XX\u518A\u898B\u3064\u304B\u308A\u307E\u3057\u305F\u3002\u8868\u793A\u3057\u307E\u3059\u304B\uFF1F","scroll_mode":"\u30BF\u30C6\u8AAD\u307F\u30E2\u30FC\u30C9","second":"\u79D2","select_language":"\u8A00\u8A9E\u3092\u9078\u629E","server_config":"\u30B5\u30FC\u30D0\u30FC\u8A2D\u5B9A","server_setting":"Comigo\u30B5\u30FC\u30D0\u30FC\u8A2D\u5B9A","set_back_color":"\u80CC\u666F\u8272\u3092\u8A2D\u5B9A:","set_interface_color":"UI\u30AB\u30E9\u30FC\u8A2D\u5B9A:","show_book_titles":"\u66F8\u7C4D\u30BF\u30A4\u30C8\u30EB\u3092\u8868\u793A","show_file_icon":"\u30D5\u30A1\u30A4\u30EB\u30A2\u30A4\u30B3\u30F3\u3092\u8868\u793A","show_header":"\u30BF\u30A4\u30C8\u30EB\u3092\u8868\u793A","show_page_num":"\u30DA\u30FC\u30B8\u756A\u53F7\u3092\u8868\u793A","simplify_book_titles":"\u66F8\u7C4D\u30BF\u30A4\u30C8\u30EB\u3092\u7C21\u7565\u5316","single_page_mode":"\u30B7\u30F3\u30B0\u30EB\u30DA\u30FC\u30B8\u30E2\u30FC\u30C9","single_page_width":"\u30B7\u30F3\u30B0\u30EB\u30DA\u30FC\u30B8\u5E45:","sort_by_default":"\u30C7\u30D5\u30A9\u30EB\u30C8\u9806","sort_by_filename":"\u30D5\u30A1\u30A4\u30EB\u540D\u9806 (A-Z)","sort_by_filename_reverse":"\u30D5\u30A1\u30A4\u30EB\u540D\u9806 (Z-A)","sort_by_filesize":"\u30D5\u30A1\u30A4\u30EB\u30B5\u30A4\u30BA\u9806 (\u5927\u2192\u5C0F)","sort_by_filesize_reverse":"\u30D5\u30A1\u30A4\u30EB\u30B5\u30A4\u30BA\u9806 (\u5C0F\u2192\u5927)","sort_by_modify_time":"\u66F4\u65B0\u65E5\u6642\u9806 (\u65B0\u2192\u65E7)","sort_by_modify_time_reverse":"\u66F4\u65B0\u65E5\u6642\u9806 (\u65E7\u2192\u65B0)","sort_reverse":"\uFF08\u9006\u9806\uFF09","start_sketch_message":"\u30AF\u30ED\u30C3\u30AD\u30FC\u30E2\u30FC\u30C9\u304C\u958B\u59CB\u3055\u308C\u307E\u3057\u305F\u3002\u826F\u3044\u4E00\u65E5\u3092\uFF01","start_sketch_mode":"\u30AF\u30ED\u30C3\u30AD\u30FC\u30E2\u30FC\u30C9\u958B\u59CB","starting_from_beginning":"\u6700\u521D\u304B\u3089\u958B\u59CB","starting_from_beginning_hint":"\u6700\u521D\u306E\u30DA\u30FC\u30B8\u304B\u3089\u8AAD\u307F\u8FBC\u307F\u307E\u3059","stop_sketch_mode":"\u30AF\u30ED\u30C3\u30AD\u30FC\u30E2\u30FC\u30C9\u7D42\u4E86","submit":"\u9001\u4FE1","success_fullscreen":"\u30D5\u30EB\u30B9\u30AF\u30EA\u30FC\u30F3\u306B\u6210\u529F\u3057\u307E\u3057\u305F","successfully_loaded_reading_progress":"\u8AAD\u66F8\u9032\u6357\u3092\u6B63\u5E38\u306B\u8AAD\u307F\u8FBC\u307F\u307E\u3057\u305F","switch_to_flip_mode":"\u30DA\u30FC\u30B8\u3081\u304F\u308A\u30E2\u30FC\u30C9\u306B\u5207\u66FF","switch_to_scrolling_mode":"\u30B9\u30AF\u30ED\u30FC\u30EB\u30E2\u30FC\u30C9\u306B\u5207\u66FF","sync_page":"\u30EA\u30E2\u30FC\u30C8\u30DA\u30FC\u30B8\u540C\u671F","temp_future_hint":"\u307E\u3060\u5B8C\u6210\u3057\u3066\u3044\u306A\u3044\u3044\u304F\u3064\u304B\u306E\u6A5F\u80FD\u3092\u3001\u4E00\u6642\u7684\u306B\u306B\u7F6E\u304F\u3002","test":"\u30C6\u30B9\u30C8","to_flip_mode":"\u30E8\u30B3\u8AAD\u307F\u30E2\u30FC\u30C9\u3078","to_infinite_dropdown_mode":"\u7121\u9650\u30C9\u30ED\u30C3\u30D7\u30C0\u30A6\u30F3\u30E2\u30FC\u30C9\u3078","to_pagination_mode":"\u30DA\u30FC\u30B8\u30F3\u30B0\u30E2\u30FC\u30C9\u3078","to_scroll_mode":"\u30BF\u30C6\u8AAD\u30E2\u30FC\u30C9\u3078","total_is":"\u5408\u8A08:","total_time":"\u5408\u8A08\u6642\u9593:","type_or_paste_content":"\u5165\u529B\u307E\u305F\u306F\u8CBC\u308A\u4ED8\u3051","upload_file":"\u30D5\u30A1\u30A4\u30EB\u3092\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9","uploaded_folder_hint":"\u30D5\u30A1\u30A4\u30EB\u306F\u5B9F\u884C\u30C7\u30A3\u30EC\u30AF\u30C8\u30EA\u5185\u306EComigoUpload\u30D5\u30A9\u30EB\u30C0\u306B\u30A2\u30C3\u30D7\u30ED\u30FC\u30C9\u3055\u308C\u307E\u3059","width_use_fixed_value":"\u6A2A\u5E45: \u56FA\u5B9A\u5024","width_use_percent":"\u6A2A\u5E45: \u30D1\u30FC\u30BB\u30F3\u30C6\u30FC\u30B8"}');


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


// 将 Alpine 实例添加到窗口对象中。
window.Alpine = (0, $8c83eaf28779ff46$export$2e2bcd8739ae039);
// Alpine Persist 插件，用于持久化存储。默认存储到 localStorage。
// 详细用法参见： https://alpinejs.dev/plugins/persist
(0, $8c83eaf28779ff46$export$2e2bcd8739ae039).plugin((0, $9b2f94dab0f686ea$export$2e2bcd8739ae039));
(0, $8c83eaf28779ff46$export$2e2bcd8739ae039).plugin((0, $fd2717825660a9ff$export$2e2bcd8739ae039));
(0, $1bac384020b50752$export$2e2bcd8739ae039).use((0, $199830a05f92d3d0$export$2e2bcd8739ae039)).init({
    debug: false,
    // // 在 setTimeout（默认异步行为）内的 init（） 中触发资源加载。如果您的后端同步加载资源，请将其设置为 false - 这样，
    // // 可以在 init（） 之后调用 i18next.t（） 而无需依赖初始化回调。此选项仅适用于同步（阻塞）加载后端，例如 i18next-fs-backend 和 i18next-sync-fs-backend！
    initImmediate: true,
    //lng: 'en', // if you're using a language detector, do not define the lng option
    // supportedLngs: ['en', 'cn', 'ja'],
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
            translation: (0, (/*@__PURE__*/$parcel$interopDefault($ae6560896f9c4f0c$exports)))
        },
        en: {
            translation: (0, (/*@__PURE__*/$parcel$interopDefault($ae6560896f9c4f0c$exports)))
        },
        'zh-CN': {
            translation: (0, (/*@__PURE__*/$parcel$interopDefault($c8762dc563519366$exports)))
        },
        zh: {
            translation: (0, (/*@__PURE__*/$parcel$interopDefault($c8762dc563519366$exports)))
        },
        'ja-JP': {
            translation: (0, (/*@__PURE__*/$parcel$interopDefault($3849b1934e918b8a$exports)))
        },
        ja: {
            translation: (0, (/*@__PURE__*/$parcel$interopDefault($3849b1934e918b8a$exports)))
        }
    }
}).then(function(t) {
//console.log(t('test'))
// i18next.changeLanguage('en', (err, t) => {
//     if (err) return console.log('something went wrong loading', err);
//     console.log(t('test'));
// });
});
window.i18next = (0, $1bac384020b50752$export$2e2bcd8739ae039 // 使i18next在全局作用域中可用
);
if (document.getElementById('FullScreenIcon')) document.getElementById('FullScreenIcon').addEventListener('click', ()=>{
    if ((0, $fe9db8f8165fa827$export$2e2bcd8739ae039).isEnabled) (0, $fe9db8f8165fa827$export$2e2bcd8739ae039).toggle();
    else // Ignore or do something else
    (0, $1bac384020b50752$export$2e2bcd8739ae039).t('not_support_fullscreen');
});
// 用Alpine Persist 注册全局变量
// https://alpinejs.dev/plugins/persist#using-alpine-persist-global
// global 全局设置
(0, $8c83eaf28779ff46$export$2e2bcd8739ae039).store('global', {
    // bgPattern 背景花纹
    bgPattern: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist('grid-line').as('global.bgPattern'),
    autoCrop: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(false).as('global.autoCrop'),
    autoCropNum: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(1).as('global.autoCropNum'),
    // userID 当前用户ID  用于同步阅读进度 随机生成
    userID: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(Math.random().toString(36).substring(2)).as('global.userID'),
    // debugMode 是否开启调试模式
    debugMode: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(true).as('global.debugMode'),
    // readerMode 当前阅读模式
    readMode: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist('scroll').as('global.readMode'),
    //是否通过websocket同步翻页
    syncPageByWS: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(true).as('global.syncPageByWS'),
    // bookSortBy 书籍排序方式 以按照文件名、修改时间、文件大小排序（或反向排序）
    bookSortBy: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist('name').as('global.bookSortBy'),
    // pageSortBy 书页排序顺序 以按照文件名、修改时间、文件大小排序（或反向排序）
    pageSortBy: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist('name').as('global.pageSortBy'),
    language: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist('en').as('global.language'),
    toggleReadMode () {
        this.readMode = this.readMode === 'flip' ? 'scroll' : 'flip';
    }
});
// BookShelf 书架设置
(0, $8c83eaf28779ff46$export$2e2bcd8739ae039).store('shelf', {
    bookCardMode: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist('gird').as('shelf.bookCardMode'),
    showTitle: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(true).as('shelf.showTitle'),
    showFileIcon: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(true).as('shelf.showFileIcon'),
    simplifyTitle: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(true).as('shelf.simplifyTitle'),
    InfiniteDropdown: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(false).as('shelf.InfiniteDropdown'),
    bookCardShowTitleFlag: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(true).as('shelf.bookCardShowTitleFlag'),
    syncScrollFlag: false,
    // 屏幕宽横比,inLandscapeMode的判断依据
    aspectRatio: 1.2,
    // 可见范围宽高的具体值
    clientWidth: 0,
    clientHeight: 0
});
// Scroll 卷轴模式
(0, $8c83eaf28779ff46$export$2e2bcd8739ae039).store('scroll', {
    nowPageNum: 1,
    simplifyTitle: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(true).as('scroll.simplifyTitle'),
    //下拉模式下，漫画页面的底部间距。单位px。
    marginBottomOnScrollMode: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(10).as('scroll.marginBottomOnScrollMode'),
    //卷轴模式下，是否无限下拉
    InfiniteDropdown: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(true).as('scroll.InfiniteDropdown'),
    syncScrollFlag: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(false).as('scroll.syncScrollFlag'),
    imageMaxWidth: 400,
    // 屏幕宽横比,inLandscapeMode的判断依据
    aspectRatio: 1.2,
    // 可见范围宽高的具体值
    clientWidth: 0,
    clientHeight: 0,
    //漫画页的单位,是否使用固定值
    widthUseFixedValue: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(true).as('scroll.widthUseFixedValue'),
    //横屏(Landscape)状态的漫画页宽度,百分比
    singlePageWidth_Percent: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(60).as('scroll.singlePageWidth_Percent'),
    doublePageWidth_Percent: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(95).as('scroll.doublePageWidth_Percent'),
    //横屏(Landscape)状态的漫画页宽度。px。
    singlePageWidth_PX: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(720).as('scroll.singlePageWidth_PX'),
    doublePageWidth_PX: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(1200).as('scroll.doublePageWidth_PX'),
    //书籍数据,需要从远程拉取
    //是否显示顶部页头
    showHeaderFlag: true,
    //是否显示页数
    show_page_num: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(false).as('scroll.show_page_num'),
    //ws翻页相关
    syncPageByWS: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(false).as('scroll.syncPageByWS')
});
// Flip 翻页模式
(0, $8c83eaf28779ff46$export$2e2bcd8739ae039).store('flip', {
    nowPageNum: 1,
    allPageNum: 100,
    imageMaxWidth: 400,
    isLandscapeMode: true,
    isPortraitMode: false,
    //自动隐藏工具条
    autoHideToolbar: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(true).as('flip.autoHideToolbar'),
    //是否显示页头
    show_header: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(true).as('flip.show_header'),
    //是否显示页脚
    showFooter: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(true).as('flip.showFooter'),
    //是否显示页数
    show_page_num: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(false).as('flip.show_page_num'),
    //是否是右半屏翻页（从右到左）?日本漫画从左到右(false)
    rightToLeft: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(false).as('flip.rightToLeft'),
    //双页模式
    doublePageMode: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(false).as('flip.doublePageMode'),
    //自动拼合双页(TODO)
    autoDoublePageMode: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(false).as('flip.autoDoublePageModeFlag'),
    //是否保存阅读进度（页数）
    saveReadingProgress: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(true).as('flip.saveReadingProgress'),
    //素描模式标记
    sketchModeFlag: false,
    //是否显示素描提示
    showPageHint: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(false).as('flip.showPageHint'),
    //翻页间隔时间
    sketchFlipSecond: 30,
    //计时用,从0开始
    sketchSecondCount: 0
});
// 自定义主题
(0, $8c83eaf28779ff46$export$2e2bcd8739ae039).store('theme', {
    theme: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist('light').as('theme'),
    interfaceColor: '#F5F5E4',
    backgroundColor: '#E0D9CD',
    textColor: '#000000',
    toggleTheme () {
        this.theme = this.theme === 'light' ? 'dark' : 'light';
    }
});
// 由于 Cookie “cookie.someCookieKey”缺少正确的“sameSite”属性值，缺少“SameSite”或含有无效值的 Cookie
// 即将被视作指定为“Lax”，该 Cookie 将无法发送至第三方上下文中。若您的应用程序依赖这组 Cookie 以在不同上下文中工作，
// 请添加“SameSite=None”属性。若要了解“SameSite”属性的更多信息，请参阅：https://developer.mozilla.org/docs/Web/HTTP/Headers/Set-Cookie/SameSite
// https://alpinejs.dev/plugins/persist#custom-storage
// 定义自定义存储对象，公开 getItem 函数和 setItem 函数
// 使用会话 cookie 作为存储
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
};
// 使用 cookieStorage 作为存储
(0, $8c83eaf28779ff46$export$2e2bcd8739ae039).store('cookie', {
    someCookieKey: (0, $8c83eaf28779ff46$export$2e2bcd8739ae039).$persist(false).using(cookieStorage).as('cookie.someCookieKey')
});
//请求图片文件时，可添加的额外参数
const $0c7aa0f20b85e54e$var$imageParameters = {
    resize_width: -1,
    resize_height: -1,
    do_auto_resize: false,
    resize_max_width: 800,
    resize_max_height: -1,
    do_auto_crop: false,
    auto_crop_num: 1,
    gray: false
};
//添加各种字符串参数,不需要的话为空
const $0c7aa0f20b85e54e$var$resize_width_str = $0c7aa0f20b85e54e$var$imageParameters.resize_width > 0 ? "&resize_width=" + $0c7aa0f20b85e54e$var$imageParameters.resize_width : "";
const $0c7aa0f20b85e54e$var$resize_height_str = $0c7aa0f20b85e54e$var$imageParameters.resize_height > 0 ? "&resize_height=" + $0c7aa0f20b85e54e$var$imageParameters.resize_height : "";
const $0c7aa0f20b85e54e$var$gray_str = $0c7aa0f20b85e54e$var$imageParameters.gray ? "&gray=true" : "";
const $0c7aa0f20b85e54e$var$do_auto_resize_str = $0c7aa0f20b85e54e$var$imageParameters.do_auto_resize ? "&resize_max_width=" + $0c7aa0f20b85e54e$var$imageParameters.resize_max_width : "";
const $0c7aa0f20b85e54e$var$resize_max_height_str = $0c7aa0f20b85e54e$var$imageParameters.resize_max_height > 0 ? "&resize_max_height=" + $0c7aa0f20b85e54e$var$imageParameters.resize_max_height : "";
const $0c7aa0f20b85e54e$var$auto_crop_str = $0c7aa0f20b85e54e$var$imageParameters.do_auto_crop ? "&auto_crop=" + $0c7aa0f20b85e54e$var$imageParameters.auto_crop_num : "";
//所有附加的转换参数
let $0c7aa0f20b85e54e$var$addStr = $0c7aa0f20b85e54e$var$resize_width_str + $0c7aa0f20b85e54e$var$resize_height_str + $0c7aa0f20b85e54e$var$do_auto_resize_str + $0c7aa0f20b85e54e$var$resize_max_height_str + $0c7aa0f20b85e54e$var$auto_crop_str + $0c7aa0f20b85e54e$var$gray_str;
if ($0c7aa0f20b85e54e$var$addStr !== "") {
    $0c7aa0f20b85e54e$var$addStr = "?" + $0c7aa0f20b85e54e$var$addStr.substring(1);
    console.log("addStr:", $0c7aa0f20b85e54e$var$addStr);
}
// Start Alpine.
(0, $8c83eaf28779ff46$export$2e2bcd8739ae039).start();
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
//# sourceMappingURL=scripts.js.map
