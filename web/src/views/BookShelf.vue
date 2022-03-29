<template>
    <div id="BookShelf">
        <Header
            class="footer"
            v-if="this.showHeaderFlag"
            :bookIsFolder="false"
            :bookName="this.bookShelfTitle"
            :showReturnIcon="this.headerShowReturnIcon"
            :setDownLoadLink="false"
        >
            <!-- 右边的设置图标，点击屏幕中央也可以打开 -->
            <n-icon size="40" @click="drawerActivate('right')">
                <settings-outline />
            </n-icon>
        </Header>

        <!-- 渲染书架部分 -->
        <div class="shelf">
            <!-- 使用tailwindcss提供的flex布局： -->
            <!-- flex-row：https://www.tailwindcss.cn/docs/flex-direction -->
            <!-- 使用 flex-wrap 允许 flex 项目换行 https://www.tailwindcss.cn/docs/flex-wrap -->
            <!-- 在组件中使用v-for时，key是必须的 -->
            <div class="flex flex-row flex-wrap">
                <BookCard
                    v-for="(book_info, key) in this.bookshelf"
                    :key="key"
                    :title="book_info.name"
                    :id="book_info.id"
                    :image_src="book_info.cover.url"
                    :ReaderMode="this.ReaderMode"
                    :showTitle="this.bookCardShowTitleFlag"
                    :childBookNum="book_info.child_book_num ? 'x' + book_info.child_book_num : ''"
                    @click="onOpenBook(book_info.id, book_info.book_type)"
                ></BookCard>
            </div>
        </div>

        <Drawer
            :initDrawerActive="this.drawerActive"
            :initDrawerPlacement="this.drawerPlacement"
            @saveConfig="this.saveConfigToCookie"
            @startSketch="this.startSketchMode"
            @closeDrawer="this.drawerDeactivate"
            @setRM="this.OnSetReaderMode"
            :ReaderMode="this.ReaderMode"
            :sketching="false"
        >
            <span>{{ $t("setBackColor") }}</span>
            <n-color-picker v-model:value="model.color" :modes="['rgb']" :show-alpha="false" />
            <!-- 分割线 -->
            <n-divider />
            <!-- 开关：显示原图 黑白 -->
            <n-space>
                <n-switch
                    size="large"
                    v-model:value="this.bookCardShowTitleFlag"
                    @update:value="setBookCardShowTitleFlag"
                >
                    <template #checked>显示小标题</template>
                    <template #unchecked>显示小标题</template>
                </n-switch>
            </n-space>
        </Drawer>
        <Bottom></Bottom>
    </div>
</template>

<script>
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
//Firefox插件Textarea Cache 报错：源映射错误：Error: NetworkError when attempting to fetch resource.
//Firefox插件Video DownloadHelper报错:已不赞成使用 CanvasRenderingContext2D 中的 drawWindow 方法
import { NIcon, NDivider, NColorPicker, NSwitch, NSpace } from 'naive-ui'
import Header from "@/components/Header.vue";
import Drawer from "@/components/Drawer.vue";
import BookCard from "@/components/BookCard.vue";
import Bottom from "@/components/Bottom.vue";
import { defineComponent, reactive } from 'vue'
import { useCookies } from "vue3-cookies";// https://github.com/KanHarI/vue3-cookies
import { SettingsOutline } from '@vicons/ionicons5'
import axios from "axios";

export default defineComponent({
    name: "BookShelf",
    props: ['readMode'],
    emits: ["setTemplate"],
    components: {
        Header,//自定义页头
        Drawer,//自定义抽屉
        BookCard,//自定义抽屉
        Bottom,//自定义页尾
        // NButton,//按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/button
        NSpace,
        NSwitch,
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
            bookShelfTitle: "loading",
            headerShowReturnIcon: false,
            bookCardShowTitleFlag: true,//书库中的书籍是否显示文字版标题
            ReaderMode: "scroll",
            maxDepth: 1,
            bookshelf: [{
                name: "loading",
                all_page_num: 1,
                id: "12345",
                book_type: "dir",
                parent_folder: "",
                depth: 1,
                pages: [
                    {
                        height: 500,
                        width: 449,
                        url: "/images/loading.jpg",
                    },
                ],
                cover: {
                    filename: "loading.jpg",
                    height: 500,
                    width: 449,
                    url: "/images/loading.jpg",
                },
            }],
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
        this.getBookShelfData();
        this.setReaderMode();
        //监听路由参数的变化，刷新本地的ReaderMode
        this.$watch(
            () => this.$route.params,
            () => {//想知道参数的变化的话，可把参数设置为 toParams, previousParams
                // console.log(toParams);
                // console.log(previousParams);
                // console.log("BookShelf: route change");
                this.setReaderMode();
                this.setBookShelfTitle();
                this.getBookShelfData();//会导致500错误，emmm……
            }
        )
        //初始化默认值,读取出来的都是字符串，不要直接用  			
        //书籍卡片是否显示文字版标题
        if (localStorage.getItem("BookCardShowTitleFlag") === "true") {
            this.bookCardShowTitleFlag = true;
        } else if (localStorage.getItem("BookCardShowTitleFlag") === "false") {
            this.bookCardShowTitleFlag = false;
        }
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
    //在绑定元素的父组件被挂载后调用。
    mounted() {
        //this.bookshelf.lenth != 1 && 
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
    },
    methods: {
        //书籍卡片是否显示文字标题
        setBookCardShowTitleFlag(value) {
            this.bookCardShowTitleFlag = value;
            localStorage.setItem("BookCardShowTitleFlag", value);
            // console.log("成功保存设置: BookCardShowTitleFlag=" + localStorage.getItem("BookCardShowTitleFlag"));
        },

        //初始化或者路由变化时，更新本地BookShelf相关数据
        getBookShelfData() {
            if (this.$route.params.group_id) {
                // console.log("BookShelf getBookShelfData!  this.$route.params.id" + this.$route.params.id)
                this.getBooksGroupByBookID(this.$route.params.group_id)
                this.headerShowReturnIcon = true
            } else {
                this.initBookShelf()
                this.headerShowReturnIcon = false
            }
        },
        //初始化或者路由变化时，读取其他页面的更改，并存储到本地存储的阅读器模式（ReaderMode）这个值，
        setReaderMode() {
            //阅读器模式，scroll或flip
            if (localStorage.getItem("ReaderMode") != null) {
                this.ReaderMode = localStorage.getItem("ReaderMode");
            }
        },
        //获取所有书籍信息
        initBookShelf() {
            axios
                .get("getlist?max_depth=1")
                .then((response) => (this.bookshelf = response.data))
                .finally(() => {
                    this.setBookShelfTitle();
                })
        },
        getBooksGroupByBookID(group_id) {
            // console.log("getBooksGroupByBookID bookID:" + group_id);
            //根据路由参数获取特定书籍组
            axios
                .get("getlist?book_group_book_id=" + group_id)
                .then((response) => {
                    if (response.data[0].name != null) {
                        this.bookshelf = response.data
                    }
                })
                .finally(() => {
                    this.setBookShelfTitle();
                })
        },
        setBookShelfTitle() {
            if (this.bookshelf[0].depth != null) {
                this.max_depth = this.bookshelf[0].depth
            }
            //默认显示服务器版本
            if (this.$store.state.server_status.ServerName) {
                //设置浏览器标签标题
                document.title = this.$store.state.server_status.ServerName
                //设置Header标题
                this.bookShelfTitle = this.$store.state.server_status.ServerName
            }
            //如果路由参数里面有ID（正在展示某个书籍组）
            if (this.$route.params.id) {
                if (this.bookshelf[0].parent_folder != null && this.bookshelf[0].parent_folder != "") {
                    // console.log("this.bookshelf[0]:" + this.bookshelf[0].parent_folder)
                    //设置浏览器标签标题
                    document.title = this.bookshelf[0].parent_folder;
                    //设置Header标题
                    this.bookShelfTitle = this.bookshelf[0].parent_folder;
                }
            }
        },
        //打开书籍
        onOpenBook(bookID, bookType) {
            // console.log("onOpenBook  bookID：" + bookID + " bookType：" + bookType)
            if (bookType == "book_group") {
                this.$router.push({ name: 'ChildBookShelf', params: { group_id: bookID } })
                return
            }
            if (this.ReaderMode == "flip" || this.ReaderMode == "sketch") {
                // 命名路由，并加上参数，让路由建立 url
                this.$router.push({ name: 'FlipMode', params: { id: bookID } })
            }
            if (this.ReaderMode == "scroll") {
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
            // this.$emit("setTemplate", "sketch");
        },
        //接收Draw的参数，继续往父组件传
        OnSetReaderMode(value) {
            if (value == "scroll" || value == "scroll")
                this.ReaderMode = value
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
.header {
    background: v-bind("model.colorHeader");
}
.shelf {
    /* padding-bottom: 20px;
    padding-left: 20px;
    padding-right: 20px;
    padding-top: 20px; */
    max-width: 100%;
    min-height: 95vh;
    height: auto;
    background: v-bind("model.color");
}
</style>
