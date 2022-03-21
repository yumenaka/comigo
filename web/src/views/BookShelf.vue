<template>
    <div id="BookShelf" class="manga">
        <Header
            class="footer"
            v-if="this.showHeaderFlag"
            :bookIsFolder="book.IsFolder"
            :bookName="book.name"
        >
            <!-- 右边的设置图标，点击屏幕中央也可以打开 -->
            <n-icon size="40" @click="drawerActivate('right')">
                <settings-outline />
            </n-icon>
        </Header>

        <!-- 渲染漫画部分 -->
        <div
            class="main_manga"
            v-for="(page, key) in book.pages"
            :key="page.url"
            @click="onMouseClick($event)"
            @mousemove="onMouseMove"
            @mouseleave="onMouseLeave"
        >
            <img v-lazy="page.url" v-bind:alt="key + 1" v-bind:key="key" />
            <div
                class="page_hint"
                v-if="showPageNumFlag_BookShelf"
            >{{ key + 1 }}/{{ book.all_page_num }}</div>
        </div>

        <Drawer
            :initDrawerActive="this.drawerActive"
            :initDrawerPlacement="this.drawerPlacement"
            @saveConfig="this.saveConfigToCookie"
            @startSketch="this.startSketchMode"
            @closeDrawer="this.drawerDeactivate"
            @setT="this.OnSetTemplate"
            :nowTemplate="this.nowTemplate"
        >
            <span>{{ $t("setBackColor") }}</span>
            <n-color-picker v-model:value="model.color" :modes="['rgb']" :show-alpha="false" />

            <!-- 分割线 -->
            <n-divider />
        </Drawer>

        <n-back-top :show="showBackTopFlag" type="info" color="#8a2be2" :right="20" :bottom="20" />
        <n-button @click="scrollToTop(90);" size="large" secondary strong>Back To Top</n-button>
    </div>
</template>

<script>
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
import { NButton, NBackTop, NIcon, NDivider, NColorPicker, } from 'naive-ui'
import Header from "@/components/Header.vue";
import Drawer from "@/components/Drawer.vue";
import { defineComponent, reactive } from 'vue'
import { useCookies } from "vue3-cookies";// https://github.com/KanHarI/vue3-cookies
import { SettingsOutline } from '@vicons/ionicons5'
import axios from "axios";

export default defineComponent({
    name: "BookShelf",
    props: ['nowTemplate'],
    emits: ["setTemplate"],
    components: {
        Header,//自定义页头，有点丑
        Drawer,//自定义抽屉，还行
        NButton,//按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/button
        NBackTop,//回到顶部按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/back-top
        // NDrawer,//抽屉，可以从上下左右4个方向冒出. https://www.naiveui.com/zh-CN/os-theme/components/drawer
        // NDrawerContent,//抽屉内容
        // NSpace,//间距 https://www.naiveui.com/zh-CN/os-theme/components/space
        // NSlider,//滑动选择  Slider https://www.naiveui.com/zh-CN/os-theme/components/slider
        // NSwitch,//开关   https://www.naiveui.com/zh-CN/os-theme/components/switch
        // NRadio,//单选  https://www.naiveui.com/zh-CN/os-theme/components/radio
        // NRadioButton,//单选  用按钮显得更优雅一点
        // NInputNumber,//数字输入 https://www.naiveui.com/zh-CN/os-theme/components/input-number
        // NRadioGroup,
        // NLayout,//布局 https://www.naiveui.com/zh-CN/os-theme/components/layout
        // NLayoutSider,
        // NLayoutContent,
        NIcon,//图标  https://www.naiveui.com/zh-CN/os-theme/components/icon
        // NPageHeader,//页头 https://www.naiveui.com/zh-CN/os-theme/components/page-header
        // NAvatar, //头像 https://www.naiveui.com/zh-CN/os-theme/components/avatar

        SettingsOutline,//图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）
        NDivider,//分割线  https://www.naiveui.com/zh-CN/os-theme/components/divider
        NColorPicker,
    },
    setup() {
        //此处不能使用this
        const { cookies } = useCookies();
        //背景颜色，颜色选择器用
        const model = reactive({
            color: "#E0D9CD",
        });

        //单选按钮绑定的数值
        // const checkedValueRef = ref(null)
        return {
            cookies,
            //背景色
            model,
            //开关的颜色
            railStyle: ({ focused, checked }) => {
                const style = {}
                if (checked) {
                    style.background = '#d03050'
                    if (focused) {
                        style.boxShadow = '0 0 0 2px #d0305040'
                    }
                } else {
                    style.background = '#2080f0'
                    if (focused) {
                        style.boxShadow = '0 0 0 2px #2080f040'
                    }
                }
                return style
            },
            //滑动选择用建议值
            marks: {
                30: '25%',
                50: '50%',
                75: '75%',
                95: '95%',
            },
        }
    },
    data() {
        return {
            book: {
                name: "loading",
                all_page_num: 2,
                pages: [
                    {
                        height: 500,
                        width: 449,
                        url: "/images/loading.jpg",
                    },
                    {
                        height: 500,
                        width: 449,
                        url: "/images/loading.jpg",
                    },
                ],
            },
            drawerActive: false,
            drawerPlacement: 'right',
            //开发模式 还没有做的功能与设置，设置Debug以后才能见到
            debugModeFlag: true,
            //书籍数据，需要从远程拉取
            //是否显示顶部页头
            showHeaderFlag: true,
            //是否显示页数
            showPageNumFlag_BookShelf: false,
            //是否显示回到顶部按钮
            showBackTopFlag: false,
            //是否正在向下滚动
            scrollDownFlag: false,
            //存储现在滚动的位置
            scrollTopSave: 0,
            //同步滚动，目前还没做
            syncScrollFlag: false,
            //鼠标点击或触摸的位置
            clickX: 0,
            clickY: 0,
            //可见范围是否是横向
            isLandscapeMode: true,
            isPortraitMode: false,
            imageMaxWidth: 10,
            //屏幕宽横比，inLandscapeMode的判断依据
            aspectRatio: 1.2,

            //状态驱动的动态 CSS!!!!!
            // https://v3.cn.vuejs.org/api/sfc-style.html#%E7%8A%B6%E6%80%81%E9%A9%B1%E5%8A%A8%E7%9A%84%E5%8A%A8%E6%80%81-css
            //图片宽度的单位，是否使用百分比
            imageWidth_usePercentFlag: false,

            //横屏(Landscape)状态的漫画页宽度，百分比
            singlePageWidth_Percent: 50,
            doublePageWidth_Percent: 95,

            //横屏(Landscape)状态的漫画页宽度，PX
            singlePageWidth_PX: 720,
            doublePageWidth_PX: 1080,

            //可见范围宽高的具体值
            clientWidth: 0,
            clientHeight: 0,
        };
    },
    //Vue3生命周期:  https://v3.cn.vuejs.org/api/options-lifecycle-hooks.html#beforecreate
    // created : 在绑定元素的属性或事件监听器被应用之前调用。
    // beforeMount : 指令第一次绑定到元素并且在挂载父组件之前调用。
    // mounted : 在绑定元素的父组件被挂载后调用。
    // beforeUpdate: 在更新包含组件的 VNode 之前调用。。
    // updated: 在包含组件的 VNode 及其子组件的 VNode 更新后调用。
    // beforeUnmount: 当指令与在绑定元素父组件卸载之前时，只调用一次。
    // unmounted: 当指令与元素解除绑定且父组件已卸载时，只调用一次。

    created() {

        axios
            .get("/getbook?id=" + this.$route.params.id)
            .then((response) => (this.book = response.data))
            .finally(console.log("成功获取书籍数据,书籍ID：" + this.$route.params.id));
        //监听路由参数的变化，刷新本地的Book数据
        this.$watch(
            () => this.$route.params,
            (toParams) => {
                console.log(toParams)
                axios
                    .get("/getbook?id=" + this.$route.params.id)
                    .then((response) => (this.book = response.data))
                    .finally(console.log("成功获取书籍数据,书籍ID：" + this.$route.params.id));
            }
        )




        window.addEventListener("scroll", this.onScroll);
        //文档视图调整大小时会触发 resize 事件。 https://developer.mozilla.org/zh-CN/docs/Web/API/Window/resize_event
        window.addEventListener("resize", this.onResize);
        this.imageMaxWidth = window.innerWidth;
        //根据cookie初始化默认值,或初始化cookie值,cookie读取出来的都是字符串，不要直接用
        //是否显示顶部页头
        if (localStorage.getItem("showHeaderFlag") === "true") {
            this.showHeaderFlag = true;
        } else if (localStorage.getItem("showHeaderFlag") === "false") {
            this.showHeaderFlag = false;
        }
        //console.log("读取cookie并初始化: showHeaderFlag=" + this.showHeaderFlag);

        //是否显示页数
        if (localStorage.getItem("showPageNumFlag_BookShelf") === "true") {
            this.showPageNumFlag_BookShelf = true;
        } else if (localStorage.getItem("showPageNumFlag_BookShelf") === "false") {
            this.showPageNumFlag_BookShelf = false;
        }
        //console.log("读取cookie并初始化: showPageNumFlag_BookShelf=" + this.showPageNumFlag_BookShelf);

        //宽度是否使用百分比
        if (localStorage.getItem("imageWidth_usePercentFlag") === "true") {
            this.imageWidth_usePercentFlag = true;
        } else if (localStorage.getItem("imageWidth_usePercentFlag") === "false") {
            this.imageWidth_usePercentFlag = false;
        }

        //javascript 数字类型转换：https://www.runoob.com/js/js-type-conversion.html
        // NaN不能通过相等操作符（== 和 ===）来判断

        //漫画页宽度，Percent
        if (localStorage.getItem("singlePageWidth_Percent") != null) {
            let saveNum = Number(localStorage.getItem("singlePageWidth_Percent"));
            if (!isNaN(saveNum)) {
                this.singlePageWidth_Percent = saveNum;
            }
        }

        if (localStorage.getItem("doublePageWidth_Percent") != null) {
            let saveNum = Number(localStorage.getItem("doublePageWidth_Percent"));
            if (!isNaN(saveNum)) {
                this.doublePageWidth_Percent = saveNum;
            }
        }

        //漫画页宽度，PX
        if (localStorage.getItem("singlePageWidth_PX") != null) {
            let saveNum = Number(localStorage.getItem("singlePageWidth_PX"));
            if (!isNaN(saveNum)) {
                this.singlePageWidth_PX = saveNum;
            }
        }
        if (localStorage.getItem("doublePageWidth_PX") != null) {
            let saveNum = Number(localStorage.getItem("doublePageWidth_PX"));
            if (!isNaN(saveNum)) {
                this.doublePageWidth_PX = saveNum;
            }
        }

        //当前颜色
        if (localStorage.getItem("BookShelfDefaultColor") != null) {
            this.model.color = localStorage.getItem("BookShelfDefaultColor");
        }

    },

    // //挂载前
    beforeMount() {

    },
    onMounted() {
        //console.log('mounted in the composition api!')
        this.isLandscapeMode = this.inLandscapeModeCheck();
        this.isPortraitMode = !this.inLandscapeModeCheck();
        // https://v3.cn.vuejs.org/api/options-lifecycle-hooks.html#beforemount
        this.$nextTick(function () {
            //视图渲染之后运行的代码
        })
    },
    //卸载前
    beforeUnmount() {
        //组件销毁前，销毁监听事件
        window.removeEventListener("scroll", this.onScroll);
        window.removeEventListener('resize', this.onResize)
    },
    methods: {
        //打开抽屉
        drawerActivate(place) {
            this.drawerActive = true
            this.drawerPlacement = place
        },
        //关闭抽屉
        drawerDeactivate() {
            this.drawerActive = false
        },
        //开始素描模式
        startSketchMode() {
            this.$emit("setTemplate", "sketch");
        },
        //接收Draw的参数，继续往父组件传
        OnSetTemplate(value) {
            this.$emit("setTemplate", value);
        },
        //如果在一个组件上使用了 v-model:xxx，应该使用 @update:xxx  https://www.naiveui.com/zh-CN/os-theme/docs/common-issues
        saveConfigToCookie() {
            // 储存cookie
            localStorage.setItem("showHeaderFlag", this.showHeaderFlag);
            localStorage.setItem("showPageNumFlag_BookShelf", this.showPageNumFlag_BookShelf);
            localStorage.setItem("imageWidth_usePercentFlag", this.imageWidth_usePercentFlag);
            localStorage.setItem("singlePageWidth_Percent", this.singlePageWidth_Percent);
            localStorage.setItem("doublePageWidth_Percent", this.doublePageWidth_Percent);
            localStorage.setItem("singlePageWidth_PX", this.singlePageWidth_PX);
            localStorage.setItem("doublePageWidth_PX", this.doublePageWidth_PX);
            localStorage.setItem("BookShelfDefaultColor", this.model.color);
        },
        setShowHeaderChange(value) {
            console.log("value:" + value);
            this.showHeaderFlag = value;
            localStorage.setItem("showHeaderFlag", value);
            console.log("cookie设置完毕: showHeaderFlag=" + localStorage.getItem("showHeaderFlag"));
        },
        setShowPageNumChange(value) {
            console.log("value:" + value);
            this.showPageNumFlag_BookShelf = value;
            localStorage.setItem("showPageNumFlag_BookShelf", value);
            console.log("cookie设置完毕: showPageNumFlag_BookShelf=" + localStorage.getItem("showPageNumFlag_BookShelf"));
        },

        setImageWidthUsePercentFlag(value) {
            console.log("value:" + value);
            this.imageWidth_usePercentFlag = value;
            localStorage.setItem("imageWidth_usePercentFlag", value);
            console.log("cookie设置完毕: imageWidth_usePercentFlag=" + this.imageWidth_usePercentFlag);
        },

        //可见区域变化的时候改变页面状态
        onResize() {
            this.imageMaxWidth = window.innerWidth
            // document.querySelectorAll(".name");
            this.clientWidth = document.documentElement.clientWidth
            this.clientHeight = document.documentElement.clientHeight
            // var aspectRatio = window.innerWidth / window.innerHeight
            this.aspectRatio = this.clientWidth / this.clientHeight
            //console.log("OnReSize,aspectRatio=" + this.aspectRatio);
            // 为了调试的时候方便，阈值是正方形
            if (this.aspectRatio > (19 / 19)) {
                this.isLandscapeMode = true
                this.isPortraitMode = false
            } else {
                this.isLandscapeMode = false
                this.isPortraitMode = true
            }
        },
        //页面滚动的时候改变各种值
        onScroll() {
            let scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
            this.scrollDownFlag = scrollTop > this.scrollTopSave;
            //防手抖，小于一定数值状态就不变
            let step = Math.abs(this.scrollTopSave - scrollTop);
            // console.log("step:", step);
            this.scrollTopSave = scrollTop
            if (step > 5) {
                this.showBackTopFlag = scrollTop > 400 && !this.scrollDownFlag;
            }
        },
        //获取鼠标位置，决定是否打开设置面板
        onMouseClick(e) {
            this.clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离，单位为px）
            this.clickY = e.y //获取鼠标的Y坐标（鼠标与屏幕顶部的距离，单位为px）
            //浏览器的视口，不包括工具栏和滚动条:
            let innerWidth = window.innerWidth
            let innerHeight = window.innerHeight
            let MinX = innerWidth * 0.4
            let MaxX = innerWidth * 0.6
            let MinY = innerHeight * 0.4
            let MaxY = innerHeight * 0.6
            if ((this.clickX > MinX && this.clickX < MaxX) && (this.clickY > MinY && this.clickY < MaxY)) {
                //alert("点中了设置区域！")
                //console.log("点中了设置区域！");
                this.drawerActivate('right')
            }
        },

        onMouseMove(e) {
            this.clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离，单位为px）
            this.clickY = e.y //获取鼠标的Y坐标（鼠标与屏幕顶部的距离，单位为px）
            //浏览器的视口，不包括工具栏和滚动条:
            let innerWidth = window.innerWidth
            let innerHeight = window.innerHeight
            let MinX = innerWidth * 0.4
            let MaxX = innerWidth * 0.6
            let MinY = innerHeight * 0.4
            let MaxY = innerHeight * 0.6
            if ((this.clickX > MinX && this.clickX < MaxX) && (this.clickY > MinY && this.clickY < MaxY)) {
                //console.log("在设置区域！");
                e.currentTarget.style.cursor = 'url(/images/SettingsOutline.png), pointer';
            } else {
                e.currentTarget.style.cursor = '';
            }
        },
        onMouseLeave(e) {
            //离开区域的时候，清空鼠标样式
            e.currentTarget.style.cursor = '';
        },

        scrollToTop(scrollDuration) {
            let scrollStep = -window.scrollY / (scrollDuration / 15),
                scrollInterval = setInterval(function () {
                    if (window.scrollY !== 0) {
                        window.scrollBy(0, scrollStep);
                    }
                    else clearInterval(scrollInterval);
                }, 15);
        },
        //根据可视区域(viewport)长宽比，确认是横屏还是竖屏
        // aspect-ratio https://developer.mozilla.org/zh-CN/docs/Web/CSS/@media/aspect-ratio
        // window.innerWidth  不是响应式依赖，所以不能用计算属性
        inLandscapeModeCheck() {
            this.aspectRatio = window.innerWidth / window.innerHeight
            // 为了测试方便，阈值是正方形
            return this.aspectRatio > (19 / 19);
        },
    },

    computed: {
        sPWL() {
            if (this.imageWidth_usePercentFlag) {
                return this.singlePageWidth_Percent + '%';
            } else {
                return this.singlePageWidth_PX + 'px';
            }
        },
        dPWL() {
            if (this.imageWidth_usePercentFlag) {
                return this.doublePageWidth_Percent + '%';
            } else {
                return this.doublePageWidth_PX + 'px';
            }
        },
        sPWP() {
            if (this.imageWidth_usePercentFlag) {
                return this.singlePageWidth_Percent + '%';
            } else {
                return this.singlePageWidth_PX + 'px';
            }
        },
        dPWP() {
            if (this.imageWidth_usePercentFlag) {
                return this.doublePageWidth_Percent + '%';
            } else {
                return this.doublePageWidth_PX + 'px';
            }
        },
    }
});
</script>

<style></style>

<style scoped>
.manga {
    max-width: 100%;
    background: v-bind("model.color");
}

.header {
    padding: 0px;
    width: 100%;
    height: 7vh;
}

/* https://developer.mozilla.org/zh-CN/docs/Web/CSS/object-fit */
.manga img {
    margin: auto;
    /* object-fit: scale-down; */
    padding: 3px 0px;
    border-radius: 7px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

.page_hint {
    /* 文字颜色 */
    color: #7e6e6e;
    /* 文字阴影：https://www.w3school.com.cn/css/css3_shadows.asp*/
    text-shadow: -1px 0 black, 0 1px black, 1px 0 black, 0 -1px black;
}

.LoadingImage {
    width: 90vw;
    max-width: 90vw;
}
.ErrorImage {
    width: 90vw;
    max-width: 90vw;
}

/* 横屏（显示区域）时的CSS样式，IE无效 */
@media screen and (min-aspect-ratio: 19/19) {
    .SinglePageImage {
        width: v-bind(sPWL);
        max-width: 100%;
    }
    .DoublePageImage {
        width: v-bind(dPWL);
        max-width: 100%;
    }
}

/* 竖屏(显示区域)CSS样式，IE无效 */
@media screen and (max-aspect-ratio: 19/19) {
    .SinglePageImage {
        /* width: 100%; */
        width: v-bind(sPWP);
        max-width: 100%;
    }
    .DoublePageImage {
        /* width: 100%; */
        width: v-bind(dPWP);
        max-width: 100%;
    }
}
</style>
