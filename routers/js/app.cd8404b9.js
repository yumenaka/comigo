(function(e){function t(t){for(var a,i,r=t[0],l=t[1],c=t[2],u=0,h=[];u<r.length;u++)i=r[u],Object.prototype.hasOwnProperty.call(s,i)&&s[i]&&h.push(s[i][0]),s[i]=0;for(a in l)Object.prototype.hasOwnProperty.call(l,a)&&(e[a]=l[a]);g&&g(t);while(h.length)h.shift()();return n.push.apply(n,c||[]),o()}function o(){for(var e,t=0;t<n.length;t++){for(var o=n[t],a=!0,r=1;r<o.length;r++){var l=o[r];0!==s[l]&&(a=!1)}a&&(n.splice(t--,1),e=i(i.s=o[0]))}return e}var a={},s={app:0},n=[];function i(t){if(a[t])return a[t].exports;var o=a[t]={i:t,l:!1,exports:{}};return e[t].call(o.exports,o,o.exports,i),o.l=!0,o.exports}i.m=e,i.c=a,i.d=function(e,t,o){i.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:o})},i.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},i.t=function(e,t){if(1&t&&(e=i(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var o=Object.create(null);if(i.r(o),Object.defineProperty(o,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var a in e)i.d(o,a,function(t){return e[t]}.bind(null,a));return o},i.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return i.d(t,"a",t),t},i.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},i.p="";var r=window["webpackJsonp"]=window["webpackJsonp"]||[],l=r.push.bind(r);r.push=t,r=r.slice();for(var c=0;c<r.length;c++)t(r[c]);var g=l;n.push([0,"chunk-vendors"]),o()})({0:function(e,t,o){e.exports=o("56d7")},"034f":function(e,t,o){"use strict";o("85ec")},"0a90":function(e,t,o){},"37eb":function(e,t,o){},"3bc4":function(e,t,o){"use strict";o("d1e9")},"56d7":function(e,t,o){"use strict";o.r(t);o("4de4"),o("96cf");var a=o("1da1"),s=(o("e260"),o("e6cf"),o("cca6"),o("a79d"),o("2b0e")),n=function(){var e=this,t=e.$createElement,o=e._self._c||t;return o("div",{attrs:{id:"app"}},[this.$store.state.setting?o("div",["scroll"===e.nowTemplate?o("ScrollTemplate"):e._e(),"sketch"===e.nowTemplate?o("SketchTemplate"):e._e(),"single"===e.nowTemplate?o("SinglePageTemplate"):e._e(),"double"===e.nowTemplate?o("DoublePageTemplate"):e._e()],1):e._e()])},i=[],r=(o("d3b7"),o("ddb0"),o("bc3a")),l=o.n(r),c=function(){var e=this,t=e.$createElement,o=e._self._c||t;return o("div",{attrs:{id:"ScrollPage"}},[o("Header",[o("h2",[this.$store.state.book.IsFolder?e._e():o("a",{attrs:{href:"raw/"+this.$store.state.book.name}},[e._v(e._s(this.$store.state.book.name)+"【Download】")]),this.$store.state.book.IsFolder?o("a",{attrs:{href:"raw/"+this.$store.state.book.name}},[e._v(e._s(this.$store.state.book.name))]):e._e()]),o("h4",[e._v("总页数："+e._s(this.$store.state.book.all_page_num))])]),e._l(this.$store.state.book.pages,(function(t,a){return o("div",{key:t.url,staticClass:"manga"},[o("img",{directives:[{name:"lazy",rawName:"v-lazy",value:t.url,expression:"page.url"}],key:a,class:e._f("check_image")(t.image_type,t.url),attrs:{H:t.height,W:t.width}}),e.showPageNum?o("p",[e._v(e._s(a+1)+"/"+e._s(e.AllPageNum))]):e._e()])})),o("p"),o("v-btn",{directives:[{name:"scroll",rawName:"v-scroll",value:e.onScroll,expression:"onScroll"},{name:"show",rawName:"v-show",value:e.btnFlag,expression:"btnFlag"}],attrs:{fab:"",color:"#bbcbff",bottom:"",right:""},on:{click:e.toTop}},[e._v("▲")]),e._t("default")],2)},g=[],u=(o("25f0"),o("93d3")),h={components:{Header:u["a"]},data:function(){return{page_mode:"multi",btnFlag:!1,showPageNum:!1,duration:300,offset:0,easing:"easeInOutCubic",AllPageNum:this.$store.state.book.all_page_num,message:{user_uuid:"",server_status:"",now_book_uuid:"",read_percent:0,msg:""}}},mounted:function(){this.initPage()},destroyed:function(){},methods:{initPage:function(){this.$cookies.keys()},getBook:function(){return this.$store.state.book},getNumber:function(e){this.page=e,console.log(e)},onScroll:function(e){if("undefined"!==typeof window){var t=window.pageYOffset||e.target.scrollTop||0;this.btnFlag=t>20}},toTop:function(){this.$vuetify.goTo(0)},initWebSocket:function(){this.$socket.onopen=this.websocketonopen,this.$socket.onerror=this.websocketonerror,this.$socket.onmessage=this.websocketonmessage,this.$socket.onclose=this.websocketclose,this.hint="连接建立"},websocketonopen:function(e){this.hint="连接成功",console.log("连接建立",e)},websocketonerror:function(e){this.hint="连接出错",this.initWebSocket(),console.log("Connection Error !!!",e)},websocketonmessage:function(e){console.log(e),this.msgList.push(JSON.parse(e.data)),this.hint="接收消息"},onChangeBook:function(e,t){this.message.now_book_uuid=t,this.message.msg="ChangeBook",this.$socket.send(JSON.stringify(this.message)),this.getBook()},websocketsend:function(e){this.$socket.send(JSON.stringify(this.message)),console.log(this.$socket.readyState,e)},websocketclose:function(e){var t=this;this.hint="连接断开",console.log("断开连接",e);var o=e.code,a=e.reason,s=e.wasClean;console.log(o,a,s);var n=setInterval((function(){t.$socket.onopen(),0==e.target.readyState&&clearInterval(n)}),3e3)}},filters:{check_image:function(e,t){if(e=e.toString(),"SinglePage"==e||"DoublePage"==e)return e;function o(e){var t=new Image;if(t.src=e,t.complete)return t.width<t.height?"SinglePage":"DoublePage";t.onload=function(){return t.onload=null,t.width<t.height?"SinglePage":"DoublePage"}}return""==e&&console.log("图片信息为空，开始本地JS分析"+t),e=o(t),e}}},m=h,p=(o("9803"),o("2877")),_=o("6544"),d=o.n(_),k=o("8336"),f=o("269a"),b=o.n(f),v=o("f977"),w=Object(p["a"])(m,c,g,!1,null,null,null),$=w.exports;d()(w,{VBtn:k["a"]}),b()(w,{Scroll:v["a"]});var P=o("d47c"),y=function(){var e=this,t=e.$createElement,o=e._self._c||t;return o("div",{attrs:{id:"SinglePageTemplate"}},[e.showHeader?o("Header",[o("h2",[this.$store.state.book.IsFolder?e._e():o("a",{attrs:{href:"raw/"+this.$store.state.book.name}},[e._v(e._s(this.$store.state.book.name)+"【Download】")]),this.$store.state.book.IsFolder?o("a",{attrs:{href:"raw/"+this.$store.state.book.name}},[e._v(e._s(this.$store.state.book.name))]):e._e()])]):e._e(),o("div",{staticClass:"single_page_main"},[e.now_page<=this.$store.state.book.all_page_num&&e.now_page>=1?o("img",{attrs:{"lazy-src":"/resources/favicon.ico",src:this.$store.state.book.pages[e.now_page-1].url},on:{click:function(t){return e.addPage(1)}}}):e._e(),o("img")]),o("v-pagination",{attrs:{length:this.$store.state.book.all_page_num,"total-visible":10},on:{input:e.toPage},model:{value:e.now_page,callback:function(t){e.now_page=t},expression:"now_page"}}),e._t("default")],2)},S=[],T={components:{Header:u["a"]},data:function(){return{now_page:1,showHeader:!0,showPagination:!0,alert:!1,easing:"easeInOutCubic",book:null,bookshelf:null,setting:null}},mounted:function(){this.book=this.$store.state.book,this.bookshelf=this.$store.state.bookshelf,this.setting=this.$store.state.setting,window.addEventListener("keyup",this.handleKeyup)},destroyed:function(){window.removeEventListener("keyup",this.handleKeyup)},methods:{initPage:function(){},addPage:function(e){this.now_page+e<this.book.all_page_num&&this.now_page+e>=1&&(this.now_page=this.now_page+e)},toPage:function(e){e<=this.book.all_page_num&&e>=1&&(this.now_page=e)},handleKeyup:function(e){var t=e||window.event||arguments.callee.caller.arguments[0];if(t)switch(t.key){case"ArrowUp":case"PageUp":case"ArrowLeft":this.addPage(-1);break;case"Space":case"ArrowDown":case"PageDown":case"ArrowRight":this.addPage(1);break;case"Home":this.toPage(1);break;case"End":this.toPage(this.book.all_page_num);break;case"Ctrl":break}},handleScroll:function(){document.body.scrollTop||document.documentElement.scrollTop}}},x=T,A=(o("9f5f"),o("891e")),C=Object(p["a"])(x,y,S,!1,null,null,null),O=C.exports;d()(C,{VPagination:A["a"]});var D=function(){var e=this,t=e.$createElement,o=e._self._c||t;return o("div",{attrs:{id:"DoublePageTemplate"}},[e.showHeader?o("Header",[o("h2",[this.$store.state.book.IsFolder?e._e():o("a",{attrs:{href:"raw/"+this.$store.state.book.name}},[e._v(e._s(this.$store.state.book.name)+"【Download】现在时刻："+e._s(e.currentTime))]),this.$store.state.book.IsFolder?o("a",{attrs:{href:"raw/"+this.$store.state.book.name}},[e._v(e._s(this.$store.state.book.name)+"现在时刻："+e._s(e.currentTime))]):e._e()])]):e._e(),o("div",{staticClass:"double_page_main"},[e.page_mark<this.$store.state.book.all_page_num?o("img",{attrs:{id:"image1","lazy-src":"/resources/favicon.ico",src:this.$store.state.book.pages[e.page_mark].url},on:{click:e.nextPageClick}}):e._e(),o("img"),e.page_mark-1>=0&&e.page_mark<this.$store.state.book.all_page_num&&"SinglePage"==this.$store.state.book.pages[e.page_mark].image_type&&"SinglePage"==this.$store.state.book.pages[e.page_mark-1].image_type?o("img",{attrs:{id:"image2","lazy-src":"/resources/favicon.ico",src:this.$store.state.book.pages[e.page_mark-1].url},on:{click:e.previousPageClick}}):e._e(),o("img")]),e.showPagination?o("v-pagination",{attrs:{length:this.$store.state.book.all_page_num-1,"total-visible":15},on:{input:e.toPage},model:{value:e.page_mark,callback:function(t){e.page_mark=t},expression:"page_mark"}}):e._e(),e._t("default")],2)},H=[],N=(o("b0c0"),{components:{Header:u["a"]},data:function(){return{showHeader:!0,localbook:{name:this.$store.state.book.name,all_image_num:this.$store.state.book.all_page_num,images:this.$store.state.book.pages,pages:null,all_page_num:0},bookshelf:null,setting:null,page_mark:0,showPagination:!0,AllPageNum:this.$store.state.book.all_page_num-1,time_cont:0,alert:!1,easing:"easeInOutCubic",timer:"",currentTime:new Date}},methods:{initLocalBook:function(){},initPageMark:function(){this.$store.state.book.all_page_num<2?this.page_mark=0:"SinglePage"==this.$store.state.book.pages[0].image_type&&"SinglePage"==this.$store.state.book.pages[1].image_type?this.page_mark=1:this.page_mark=0},toPage:function(e){(e>this.$store.state.book.all_page_num||e<0)&&console.log("page_mark error",e),this.page_mark=e,console.log(e)},nextPageClick:function(){this.page_mark>=this.AllPageNum||this.page_mark<0?console.log(this.page_mark):this.page_mark!=this.AllPageNum-1?this.page_mark!=this.AllPageNum-2?"SinglePage"==this.$store.state.book.pages[this.page_mark+1].image_type&&"SinglePage"==this.$store.state.book.pages[this.page_mark+2].image_type?(this.page_mark=this.page_mark+2,console.log(this.page_mark)):(this.page_mark=this.page_mark+1,console.log(this.page_mark)):"SinglePage"==this.$store.state.book.pages[this.AllPageNum-1].image_type&&"SinglePage"==this.$store.state.book.pages[this.AllPageNum-2].image_type?(this.page_mark=this.page_mark+2,console.log(this.page_mark)):(this.page_mark=this.page_mark+1,console.log(this.page_mark)):this.page_mark=this.page_mark+1},previousPageClick:function(){if(this.page_mark<=0)console.log(this.page_mark);else{if(1==this.page_mark)return this.page_mark=this.page_mark-1,void console.log(this.page_mark);if(2==this.page_mark)return this.page_mark=this.page_mark-1,void console.log(this.page_mark);if(this.page_mark>=this.AllPageNum)return this.page_mark=this.AllPageNum-1,void console.log(this.page_mark);if("SinglePage"==this.$store.state.book.pages[this.page_mark-2].image_type&&"SinglePage"==this.$store.state.book.pages[this.page_mark-2-1].image_type)return this.page_mark=this.page_mark-2,void console.log(this.page_mark);this.page_mark=this.page_mark-1,console.log(this.page_mark)}},nextPage:function(){this.page_mark>this.$store.state.book.all_page_num?console.log(this.page_mark):this.page_mark!=this.$store.state.book.all_page_num?this.nextPageClick():console.log(this.page_mark)},previousPage:function(){this.page_mark>this.$store.state.book.all_page_num?console.log(this.page_mark):this.page_mark!=this.$store.state.book.all_page_num?this.previousPageClick():this.page_mark=this.page_mark-1},handleKeyup:function(e){var t=e||window.event||arguments.callee.caller.arguments[0];if(t)switch(t.key){case"PageUp":case"ArrowUp":case"ArrowLeft":this.previousPage();break;case"Space":case"ArrowDown":case"PageDown":case"ArrowRight":this.nextPage();break;case"Home":this.toPage(1);break;case"End":this.toPage(this.$store.state.book.all_page_num-1);break;case"Ctrl":break}},handleScroll:function(){document.body.scrollTop||document.documentElement.scrollTop}},created:function(){var e=this;this.timer=setInterval((function(){var t=new Date,o=t.getFullYear(),a=t.getMonth()+1,s=t.getDate();a>=1&&a<=9&&(a="0"+a),s>=0&&s<=9&&(s="0"+s);var n=o+" 年 "+a+" 月 "+s+" 日 ",i=t.getHours();i>=0&&i<=9&&(i="0"+i);var r=t.getMinutes();r>=0&&r<=9&&(r="0"+r);var l=t.getSeconds();l>=0&&l<=9&&(l="0"+l),e.currentTime=n+" "+i+":"+r+":"+l}),1e3)},mounted:function(){this.time_cont=0,this.$cookies.keys(),this.initPageMark(),window.addEventListener("keyup",this.handleKeyup)},destroyed:function(){window.removeEventListener("keyup",this.handleKeyup),this.timer&&clearInterval(this.timer)}}),j=N,E=(o("810a"),Object(p["a"])(j,D,H,!1,null,null,null)),I=E.exports;d()(E,{VPagination:A["a"]});var L={name:"app",components:{ScrollTemplate:$,SketchTemplate:P["default"],SinglePageTemplate:O,DoublePageTemplate:I},data:function(){return{book:null,bookshelf:{},setting:{template:"scroll",sketch_count_seconds:90},now_page:1,duration:300,offset:0,easing:"easeInOutCubic",message:{user_uuid:"",server_status:"",now_book_uuid:"",read_percent:0,msg:""}}},mounted:function(){this.initPage(),this.$cookies.keys()},destroyed:function(){},computed:{nowTemplate:function(){return console.log("computed:"+this.$store.state.setting.template),this.$store.state.setting.template}},methods:{initPage:function(){var e=this;l.a.get("/book.json").then((function(t){return e.$store.state.book=t.data})).finally(this.book=this.$store.book),l.a.get("/setting.json").then((function(t){return e.$store.state.setting=t.data})).finally(this.setting=this.$store.setting),l.a.get("/bookshelf.json").then((function(t){return e.$store.state.bookshelf=t.data})).finally(this.bookshelf=this.$store.bookshelf)},getNumber:function(e){this.page=e,console.log(e)}}},F=L,M=(o("034f"),Object(p["a"])(F,n,i,!1,null,null,null)),B=M.exports,K=o("8c4f");s["a"].use(K["a"]);var z=[{path:"/",name:"Scrool",component:$},{path:"/sketch",name:"Sketch",component:function(){return Promise.resolve().then(o.bind(null,"d47c"))}}],J=new K["a"]({routes:z}),U=J,W=o("caf9"),R=o("f309");s["a"].use(R["a"]);var V=new R["a"]({}),Y=o("2b27"),q=o.n(Y),G=o("2f62");s["a"].config.productionTip=!1,s["a"].use(W["a"],{preLoad:4.5,attempt:10}),s["a"].use(q.a),s["a"].$cookies.config("30d"),s["a"].use(G["a"]);var Q=new G["a"].Store({state:{count:0,todos:[{id:1,text:"...",done:!0},{id:2,text:"...",done:!1}],now_page:1,book:{name:"loading",page_num:1,pages:[{height:500,width:449,url:"/resources/favicon.ico",class:"Vertical"}]},bookshelf:{},setting:{template:"scroll",sketch_count_seconds:90},message:{user_uuid:"",server_status:"",now_book_uuid:"",read_percent:0,msg:""}},mutations:{change_template_to_scroll:function(e){e.setting.template="scroll",console.log("template:"+e.setting.template)},change_template_to_double:function(e){e.setting.template="double",console.log("template:"+e.setting.template)},change_template_to_single:function(e){e.setting.template="single",console.log("template:"+e.setting.template)},change_template_to_sketch:function(e){e.setting.template="sketch",console.log("template:"+e.setting.template)},increment:function(e){e.count++},syncBookDate:function(e,t){e.book=t.msg,console.log(e.book),console.log("syncBookDate run")}},getters:{doneTodos:function(e){return e.todos.filter((function(e){return e.done}))},now_page:function(e){return e.now_page},book:function(e){return e.book},bookshelf:function(e){return e.bookshelf},setting:function(e){return e.setting},message:function(e){return e.message}},actions:{incrementAction:function(e){e.commit("increment")},getMessageAction:function(e){return Object(a["a"])(regeneratorRuntime.mark((function t(){var o,a;return regeneratorRuntime.wrap((function(t){while(1)switch(t.prev=t.next){case 0:return t.next=2,l.a.get("/book.json").then((function(e){return e.data}),(function(){return""}));case 2:o=t.sent,a={message:o},e.commit("syncBookDate",a);case 5:case"end":return t.stop()}}),t)})))()}}});new s["a"]({router:U,vuetify:V,store:Q,render:function(e){return e(B)}}).$mount("#app")},6528:function(e,t,o){},"810a":function(e,t,o){"use strict";o("0a90")},"85ec":function(e,t,o){},"93d3":function(e,t,o){"use strict";var a=function(){var e=this,t=e.$createElement,o=e._self._c||t;return o("header",{staticClass:"header"},[o("select",{directives:[{name:"model",rawName:"v-model",value:e.selectedTemplate,expression:"selectedTemplate"}],on:{change:[function(t){var o=Array.prototype.filter.call(t.target.options,(function(e){return e.selected})).map((function(e){var t="_value"in e?e._value:e.value;return t}));e.selectedTemplate=t.target.multiple?o:o[0]},function(t){return e.onChange()}]}},[o("option",{attrs:{disabled:"",value:""}},[e._v("点此切换阅读模式")]),o("option",[e._v("scroll")]),o("option",[e._v("double")]),o("option",[e._v("single")]),o("option",[e._v("sketch")])]),e._t("default")],2)},s=[],n={name:"Header",data:function(){return{mybook:this.book,selectedTemplate:""}},methods:{onChange:function(){"scroll"===this.selectedTemplate&&this.change_template_to_scroll(),"double"===this.selectedTemplate&&this.change_template_to_double(),"single"===this.selectedTemplate&&this.change_template_to_single(),"sketch"===this.selectedTemplate&&this.change_template_to_sketch()},change_template_to_scroll:function(){this.$store.commit("change_template_to_scroll")},change_template_to_double:function(){this.$store.commit("change_template_to_double")},change_template_to_single:function(){this.$store.commit("change_template_to_single")},change_template_to_sketch:function(){this.$store.commit("change_template_to_sketch")}}},i=n,r=(o("e5c7"),o("2877")),l=Object(r["a"])(i,a,s,!1,null,"5270917b",null);t["a"]=l.exports},9803:function(e,t,o){"use strict";o("6528")},"9f5f":function(e,t,o){"use strict";o("f2b5")},d1e9:function(e,t,o){},d47c:function(e,t,o){"use strict";o.r(t);var a=function(){var e=this,t=e.$createElement,o=e._self._c||t;return o("div",{attrs:{id:"SketchPage"}},[e.showHeader?o("Header",[o("h2",[o("a",{attrs:{href:"raw/"+this.$store.state.book.name}},[e._v(e._s(this.$store.state.book.name)+"现在时刻："+e._s(e.currentTime))])])]):e._e(),o("div",{staticClass:"sketch_main"},[o("div",{attrs:{id:"SketchHint"}},[o("p",[e._v(e._s(e.getNowCount())+"/"+e._s(e.getALLSeconds())+"⏳")])]),o("img",{attrs:{"lazy-src":"/resources/favicon.ico",src:this.$store.state.book.pages[e.now_page-1].url},on:{click:function(t){return e.addPage(1)}}}),o("img"),o("div",{attrs:{id:"SketchHint"}},[o("p",[e._v("🕒"+e._s(e.currentTime))])])]),e.showPagination?o("v-pagination",{attrs:{length:this.$store.state.book.all_page_num,"total-visible":10},on:{input:e.toPage},model:{value:e.now_page,callback:function(t){e.now_page=t},expression:"now_page"}}):e._e(),e._t("default")],2)},s=[],n=(o("d3b7"),o("ddb0"),o("93d3")),i={components:{Header:n["a"]},data:function(){return{showHeader:!0,time_cont:1,WaitSeconds:this.$store.state.setting.sketch_count_seconds,book:null,bookshelf:null,setting:null,showPagination:!0,now_page:1,alert:!1,easing:"easeInOutCubic",timer:"",currentTime:null}},created:function(){var e=this;this.timer=setInterval((function(){var t=new Date,o=t.getFullYear(),a=t.getMonth()+1,s=t.getDate();a>=1&&a<=9&&(a="0"+a),s>=0&&s<=9&&(s="0"+s);var n=o+"年"+a+"月"+s+"日",i=t.getHours();i>=0&&i<=9&&(i="0"+i);var r=t.getMinutes();r>=0&&r<=9&&(r="0"+r);var l=t.getSeconds();l>=0&&l<=9&&(l="0"+l),e.currentTime=i+":"+r+":"+l,console.log(n+"time_cont："+e.time_cont),e.time_cont<e.WaitSeconds?e.time_cont++:(e.time_cont=0,console.log("时间到，翻页："+e.currentTime+"秒"),e.now_page<e.$store.state.book.all_page_num?e.now_page+=1:e.now_page=1)}),1e3)},mounted:function(){this.time_cont=0,this.$cookies.keys(),window.addEventListener("keyup",this.handleKeyup)},destroyed:function(){window.removeEventListener("keyup",this.handleKeyup),this.timer&&clearInterval(this.timer)},methods:{initPage:function(){},getWaitSeconds:function(){return this.$store.state.setting.sketch_count_seconds},getNowCount:function(){var e=this.$store.state.setting.sketch_count_seconds-this.time_cont;return e>=0&&e<=9&&(e="0"+e),e},getALLSeconds:function(){var e=this.$store.state.setting.sketch_count_seconds;return e>=0&&e<=9&&(e="0"+e),e},addPage:function(e){this.now_page+e<this.$store.state.book.all_page_num&&this.now_page+e>=1&&(this.now_page=this.now_page+e)},toPage:function(e){this.now_page=e,console.log(e)},handleKeyup:function(e){var t=e||window.event||arguments.callee.caller.arguments[0];if(t)switch(t.key){case"PageUp":case"ArrowUp":case"ArrowLeft":this.addPage(-1);break;case"Space":case"ArrowDown":case"PageDown":case"ArrowRight":this.addPage(1);break;case"Home":this.toPage(1);break;case"End":this.toPage(this.$store.state.book.all_page_num-1);break;case"Ctrl":break}},handleScroll:function(){document.body.scrollTop||document.documentElement.scrollTop}}},r=i,l=(o("3bc4"),o("2877")),c=o("6544"),g=o.n(c),u=o("891e"),h=Object(l["a"])(r,a,s,!1,null,null,null);t["default"]=h.exports;g()(h,{VPagination:u["a"]})},e5c7:function(e,t,o){"use strict";o("37eb")},f2b5:function(e,t,o){}});
//# sourceMappingURL=app.cd8404b9.js.map