<template>
    <div id="BookShelf">
        <Header
            class="footer"
            v-if="this.showHeaderFlag"
            :bookIsFolder="false"
            :showReturnIcon="false"
        >
            <!-- 右边的设置图标，点击屏幕中央也可以打开 -->
            <n-icon size="40" @click="drawerActivate('right')">
                <settings-outline />
            </n-icon>
        </Header>

        <!-- 渲染书架部分 -->
        <div class="shelf">
            <!-- cols: 显示的栅格数量 -->
            <!-- x-gap: 横向间隔槽 -->
            <!-- y-gap: 纵向间隔槽 -->
            <!-- responsive: 'self' 根据自身宽度进行响应式布局，'screen' 根据屏幕断点进行响应式布局 -->
            <n-grid cols="2 s:4 m:5 l:6 xl:8 2xl:10" x-gap="2" y-gap="23" responsive="screen">
                <!-- 在组件中使用v-for时，key是必须的 -->
                <n-grid-item v-for="(book_info, key) in this.bookshelf" :key="key">
                    <BookCard
                        :title="book_info.name"
                        :id="book_info.id"
                        :image_src="book_info.cover.url"
                        :nowMode="this.nowMode"
                    ></BookCard>
                </n-grid-item>
            </n-grid>
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
    </div>
</template>

<script>
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
import { NIcon, NDivider, NColorPicker, NGrid, NGridItem, } from 'naive-ui'
import Header from "@/components/Header.vue";
import Drawer from "@/components/Drawer.vue";
import BookCard from "@/components/BookCard.vue";
import { defineComponent, reactive } from 'vue'
import { useCookies } from "vue3-cookies";// https://github.com/KanHarI/vue3-cookies
import { SettingsOutline } from '@vicons/ionicons5'
import axios from "axios";

export default defineComponent({
    name: "BookShelf",
    props: ['nowMode'],
    emits: ["setTemplate"],
    components: {
        Header,//自定义页头，有点丑
        Drawer,//自定义抽屉，还行
        BookCard,//自定义抽屉，还行
        // NButton,//按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/button
        // NBackTop,//回到顶部按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/back-top
        NGrid,//https://www.naiveui.com/zh-CN/os-theme/components/grid
        NGridItem,
        NIcon,//图标  https://www.naiveui.com/zh-CN/os-theme/components/icon
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
            colorHeader: "#d1c9c1",
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
            bookshelf: {},
            drawerActive: false,
            drawerPlacement: 'right',
            //开发模式 还没有做的功能与设置，设置Debug以后才能见到
            debugModeFlag: true,
            //书籍数据，需要从远程拉取
            //是否显示顶部页头
            showHeaderFlag: true,
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
            .get("/bookshelf.json")
            .then((response) => (this.bookshelf = response.data))
            .finally(console.log(this.bookshelf));
        // window.addEventListener("scroll", this.onScroll);
        // window.addEventListener("resize", this.onResize);
        this.imageMaxWidth = window.innerWidth;
        //初始化默认值,读取出来的都是字符串，不要直接用
        //是否显示顶部页头
        if (localStorage.getItem("showHeaderFlag") === "true") {
            this.showHeaderFlag = true;
        } else if (localStorage.getItem("showHeaderFlag") === "false") {
            this.showHeaderFlag = false;
        }
        //当前颜色
        if (localStorage.getItem("BookShelfDefaultColor") != null) {
            this.model.color = localStorage.getItem("BookShelfDefaultColor");
        }
    },
    //挂载前
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
        // window.removeEventListener("scroll", this.onScroll);
        // window.removeEventListener('resize', this.onResize)
    },
    methods: {
        //打开书籍
        onOpenBook(bookID) {
            if (this.nowTemplate == "flip" || this.nowTemplate == "sketch") {
                // 命名路由，并加上参数，让路由建立 url
                this.$router.push({ name: 'FlipMode', params: { id: bookID } })
            } else if (this.nowTemplate == "scroll") {
                // 命名路由，并加上参数，让路由建立 url
                this.$router.push({ name: 'ScrollMode', params: { id: bookID } })
            } else {
                // 命名路由，并加上参数，让路由建立 url
                this.$router.push({ name: 'ScrollMode', params: { id: bookID } })
            }
        },
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
            localStorage.setItem("BookShelfDefaultColor", this.model.color);
        },
        setShowHeaderChange(value) {
            console.log("value:" + value);
            this.showHeaderFlag = value;
            localStorage.setItem("showHeaderFlag", value);
            console.log("cookie设置完毕: showHeaderFlag=" + localStorage.getItem("showHeaderFlag"));
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
    }
});
</script>

<style scoped>
.shelf {
    padding-bottom: 10px;
    padding-left: 20px;
    padding-right: 10px;
    padding-top: 30px;
    max-width: 100%;
    min-height: 93vh;
    background: v-bind("model.color");
}
</style>
