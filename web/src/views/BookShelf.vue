<template>
    <div id="BookShelf">
        <Header :bookIsFolder="false" :headerTitle="this.bookShelfTitle" :showReturnIcon="this.headerShowReturnIcon"
            :bookID="this.bookshelf[0].id" :setDownLoadLink="false">
            <!-- 右边的设置图标,点击屏幕中央也可以打开 -->
            <n-icon size="40" @click="drawerActivate('right')">
                <settings-outline />
            </n-icon>
        </Header>

        <!-- 渲染书架部分 -->
        <div class="shelf">
            <!-- 使用tailwindcss提供的flex布局： -->
            <!-- flex-row：https://www.tailwindcss.cn/docs/flex-direction -->
            <!-- 使用 flex-wrap 允许 flex 项目换行 https://www.tailwindcss.cn/docs/flex-wrap -->
            <!-- 在组件中使用v-for时,key是必须的 -->
            <!-- justify-center：让项目沿着容器主轴的中心点对齐 https://www.tailwindcss.cn/docs/justify-content -->
            <div class="flex flex-row flex-wrap justify-center">
                <BookCard v-for="(book_info, key) in this.bookshelf" :key="key" :title="book_info.name"
                    :id="book_info.id" :image_src="book_info.cover.url" :readerMode="this.readerMode"
                    :showTitle="this.bookCardShowTitleFlag"
                    :childBookNum="book_info.child_book_num ? 'x' + book_info.child_book_num : ''"
                    :openURL="getBookCardOpenURL(book_info.id, book_info.book_type, book_info.name)"
                    :a_target="getBookCardTarget(book_info.book_type)">
                </BookCard>
            </div>
        </div>

        <Drawer :initDrawerActive="this.drawerActive" :initDrawerPlacement="this.drawerPlacement"
            @saveConfig="this.saveConfigToLocal" @startSketch="this.startSketchMode"
            @closeDrawer="this.drawerDeactivate" @setRM="this.OnSetReaderMode" :readerMode="this.readerMode"
            :sketching="false" :inBookShelf="true">
            <span>{{ $t('setInterfaceColor') }}</span>
            <n-color-picker v-model:value="model.interfaceColor" :modes="['hex']" :show-alpha="false"
                @update:value="onInterfaceColorChange" />
            <span>{{ $t("setBackColor") }}</span>
            <n-color-picker v-model:value="model.backgroundColor" :modes="['hex']" :show-alpha="false"
                @update:value="onBackgroundColorChange" />
            <p> &nbsp;</p>
            <!-- 开关：下拉阅读 -->
            <n-space>
                <n-switch size="large" :rail-style="railStyle" v-model:value="this.readerModeIsScroll"
                    @update:value="setReaderModeIsScroll">
                    <template #checked>{{ $t('scroll_mode') }}</template>
                    <template #unchecked>{{ $t('flip_mode') }}</template>
                </n-switch>
            </n-space>

            <!-- 开关：显示书名 -->
            <n-space>
                <n-switch size="large" v-model:value="this.bookCardShowTitleFlag"
                    @update:value="setBookCardShowTitleFlag">
                    <template #checked>{{ $t('show_book_titles') }}</template>
                    <template #unchecked>{{ $t('show_book_titles') }}</template>
                </n-switch>
            </n-space>
            <!-- 分割线 -->
            <n-divider />
            <n-space>
                <!-- 切换排序方式的按钮 -->
                <n-dropdown trigger="hover" :options="options" @select="onResort">
                    <n-button>{{ this.getSortHintText(this.resort_hint_key) }}</n-button>
                </n-dropdown>
            </n-space>
            <p> &nbsp;</p>
            <!-- 下载示例配置文件的按钮 -->
            <a href="/api/config.toml" target="_blank">
                <n-button>{{ $t('DownloadSampleConfigFile') }}</n-button>
            </a>
        </Drawer>
        <Bottom :softVersion="
            this.$store.state.server_status.ServerName
                ? this.$store.state.server_status.ServerName
                : 'Comigo'
        "></Bottom>
    </div>
</template>

<script>
// 直接导入组件并使用它。这种情况下,只有导入的组件才会被打包。
// Firefox插件Textarea Cache 报错：源映射错误：Error: NetworkError when attempting to fetch resource.
// Firefox插件Video DownloadHelper报错:已不赞成使用 CanvasRenderingContext2D 中的 drawWindow 方法
import { NIcon, NDivider, NColorPicker, NSwitch, NButton, NSpace, NDropdown, } from "naive-ui";
import Header from "@/components/Header.vue";
import Drawer from "@/components/Drawer.vue";
import BookCard from "@/components/BookCard.vue";
import Bottom from "@/components/Bottom.vue";

import { defineComponent, reactive } from "vue";
import { useCookies } from "vue3-cookies"; // https://github.com/KanHarI/vue3-cookies
import { SettingsOutline } from "@vicons/ionicons5";
import axios from "axios";

export default defineComponent({
    name: "BookShelf",
    props: ["testProp"],
    emits: ["setTemplate"],
    components: {
        Header, // 自定义页头
        Drawer, // 自定义抽屉
        BookCard, // 自定义抽屉
        Bottom, // 自定义页尾
        NButton,//按钮,来自:https://www.naiveui.com/zh-CN/os-theme/components/button
        NSpace,
        NSwitch,
        NIcon, // 图标  https://www.naiveui.com/zh-CN/os-theme/components/icon
        SettingsOutline, // 图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）
        NDivider, // 分割线  https://www.naiveui.com/zh-CN/os-theme/components/divider
        NColorPicker,//颜色选择器 https://www.naiveui.com/zh-CN/os-theme/components/color-picker
        NDropdown,//下拉菜单 https://www.naiveui.com/zh-CN/os-theme/components/dropdown
    },
    setup() {
        // 此处不能使用this
        const { cookies } = useCookies();
        // 背景颜色,颜色选择器用
        const model = reactive({
            interfaceColor: "#F5F5E4",
            backgroundColor: "#E0D9CD",
        });
        // 单选按钮绑定的数值
        // const checkedValueRef = ref(null)
        return {
            cookies,
            // 背景色
            model,
            // 开关的颜色
            railStyle: ({ focused, checked }) => {
                const style = {};
                if (checked) {
                    style.background = "#d03050";
                    if (focused) {
                        style.boxShadow = "0 0 0 2px #d0305040";
                    }
                } else {
                    style.background = "#2080f0";
                    if (focused) {
                        style.boxShadow = "0 0 0 2px #2080f040";
                    }
                }
                return style;
            },
            // 滑动选择用建议值
            marks: {
                30: "25%",
                50: "50%",
                75: "75%",
                95: "95%",
            },
        };
    },
    data() {
        return {
            resort_hint_key: "filename",//书籍的排序方式。可以按照文件名、修改时间、文件大小排序（或反向排序）
            options: [
                {
                    label: this.$t('sort_by_filename'),
                    key: "filename",
                },
                {
                    label: this.$t('sort_by_modify_time'),
                    key: "modify_time"
                },
                {
                    label: this.$t('sort_by_filesize'),
                    key: "filesize"
                },
                {
                    label: this.$t('sort_by_filename') + this.$t('sort_reverse'),
                    key: "filename_reverse",
                },
                {
                    label: this.$t('sort_by_modify_time') + this.$t('sort_reverse'),
                    key: "modify_time_reverse"
                },
                {
                    label: this.$t('sort_by_filesize') + this.$t('sort_reverse'),
                    key: "filesize_reverse"
                },
            ],
            readerMode: "scroll",
            readerModeIsScroll: true,
            bookShelfTitle: "loading",
            headerShowReturnIcon: false,
            bookCardShowTitleFlag: true, // 书库中的书籍是否显示文字版标题
            maxDepth: 1,
            bookshelf: [
                {
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
                },
            ],
            drawerActive: false,
            drawerPlacement: "right",
            // 开发模式 还没有做的功能与设置,设置Debug以后才能见到
            debugModeFlag: true,
            // 书籍数据,需要从远程拉取
            // 同步滚动,目前还没做
            syncScrollFlag: false,
            // 鼠标点击或触摸的位置
            clickX: 0,
            clickY: 0,
            // 可见范围是否是横向
            isLandscapeMode: true,
            isPortraitMode: false,
            imageMaxWidth: 10,
            // 屏幕宽横比,inLandscapeMode的判断依据
            aspectRatio: 1.2,
            // 可见范围宽高的具体值
            clientWidth: 0,
            clientHeight: 0,
        };
    },
    // Vue3生命周期:  https://v3.cn.vuejs.org/api/options-lifecycle-hooks.html#beforecreate
    // created : 在绑定元素的属性或事件监听器被应用之前调用。
    // beforeMount : 指令第一次绑定到元素并且在挂载父组件之前调用。
    // mounted : 在绑定元素的父组件被挂载后调用。
    // beforeUpdate: 在更新包含组件的 VNode 之前调用。。
    // updated: 在包含组件的 VNode 及其子组件的 VNode 更新后调用。
    // beforeUnmount: 当指令与在绑定元素父组件卸载之前时,只调用一次。
    // unmounted: 当指令与元素解除绑定且父组件已卸载时,只调用一次。
    created() {

        // 初始化默认值,读取出来的都是字符串,不要直接用
        //书籍排序方式。可以按照文件名、修改时间、文件大小排序（或反向排序）
        if (localStorage.getItem("ResortKey_BookShelf") !== null) {
            this.resort_hint_key = localStorage.getItem("ResortKey_BookShelf")
        }


        //TODO:读取服务器书籍总数，避免重复加载（尤其是书很多的时候）
        // 从服务器上拉取书架信息
        this.getBookShelfData();
        // 刷新ReadMode
        this.refreshReadMode();
        // 监听路由参数的变化,刷新本地数据
        this.$watch(
            () => this.$route.params,
            () => {
                // 想知道参数的变化的话,可把参数设置为 toParams, previousParams
                // console.log(toParams);
                // console.log(previousParams);
                // console.log("BookShelf: route change");
                this.refreshReadMode();
                this.setBookShelfTitle();
                this.getBookShelfData();
            }
        );
        // 继续初始化默认值
        // 书籍卡片是否显示文字版标题
        if (localStorage.getItem("BookCardShowTitleFlag") === "true") {
            this.bookCardShowTitleFlag = true;
        } else if (localStorage.getItem("BookCardShowTitleFlag") === "false") {
            this.bookCardShowTitleFlag = false;
        }

        // 当前颜色
        if (localStorage.getItem("BackgroundColor") != null) {
            this.model.backgroundColor = localStorage.getItem("BackgroundColor");
        }
        if (localStorage.getItem("InterfaceColor") != null) {
            this.model.interfaceColor = localStorage.getItem("InterfaceColor");
        }
    },
    // 挂载前
    beforeMount() { },
    // 在绑定元素的父组件被挂载后调用。
    mounted() {
        // this.bookshelf.lenth != 1 &&
        // console.log('mounted in the composition api!')
        this.isLandscapeMode = this.inLandscapeModeCheck();
        this.isPortraitMode = !this.inLandscapeModeCheck();
        // https://v3.cn.vuejs.org/api/options-lifecycle-hooks.html#beforemount
        this.$nextTick(function () {
            // 视图渲染之后运行的代码
        });
    },

    // 卸载前
    beforeUnmount() { },
    methods: {
        //根据文件名、修改时间、文件大小等参数重新排序
        onResort(key) {
            this.resort_hint_key = key
            localStorage.setItem("ResortKey_BookShelf", key)
            if (this.$route.params.group_id) {
                console.log("onResort  bookID：" + this.$route.params.group_id + ", key:" + key)
                this.$router.push({
                    name: "ChildBookShelf",
                    params: { group_id: this.$route.params.group_id },
                    replace: true,
                    query: { sort_by: key }
                });
            } else {
                console.log("onResort  key：" + key)
                this.$router.push({ name: "BookShelf", replace: true, query: { sort_by: key } })
            }
        },
        //返回“重新排序”选择菜单的文字提示
        getSortHintText(key) {
            switch (key) {
                case "filename": return this.$t('sort_by_filename');
                case "modify_time": return this.$t('sort_by_modify_time');
                case "filesize": return this.$t('sort_by_filesize');
                case "filename_reverse": return this.$t('sort_by_filename') + this.$t('sort_reverse');
                case "modify_time_reverse": return this.$t('sort_by_modify_time') + this.$t('sort_reverse');
                case "filesize_reverse": return this.$t('sort_by_filesize') + this.$t('sort_reverse');
                default:
                    return this.$t('re_sort');
            }
        },

        getBookCardTarget(bookType) {
            if (bookType === ".pdf" || bookType === "video" || bookType === "audio" || bookType === "unknown") {
                return "_blank";
            }
            // if (this.readerMode == "flip" || this.readerMode == "sketch" || this.readerMode == "scroll") {
            //     return "_self";
            // }
            return "_self";
        },

        getBookCardOpenURL(bookID, bookType, bookName) {
            // console.log("getBookCardOpenURL  bookID：" + bookID + " bookType：" + bookType)
            if (bookType == "book_group") {
                return "/#/child_shelf/" + bookID + "/";
            }
            if (bookType === ".pdf" || bookType === "video" || bookType === "audio" || bookType === "unknown") {
                return "/api/raw/" + bookID + "/" + bookName;
            }
            if (this.readerMode == "flip" || this.readerMode == "sketch") {
                return "/#/flip/" + bookID
            }
            if (this.readerMode == "scroll") {
                // 命名路由,并加上参数,让路由建立 url
                return "/#/scroll/" + bookID
            }
        },

        //已经被 getBookCardOpenURL与getBookCardOpenURL代替
        // 打开书阅读，或继续在书架里展示一组书
        onOpenBook(bookID, bookType) {
            // console.log("onOpenBook  bookID：" + bookID + " bookType：" + bookType)
            if (bookType == "book_group") {
                this.$router.push({
                    name: "ChildBookShelf",
                    params: { group_id: bookID },
                });
                return;
            }
            if (bookType === ".pdf" || bookType === ".mp4" || bookType === "video" || bookType === "audio" || bookType === "unknown") {
                //打开新的标签页,跳转到浏览器自带的打开功能
                // Window open() 方法 https://www.runoob.com/jsref/met-win-open.html
                axios
                    .get("/getbook?id=" + bookID)
                    .then((response) => (window.open('api/raw/' + bookID + '/' + response.data.name, '_blank')))// _blank - URL加载到一个新的窗口。默认值
                    .finally(
                        () => {
                            console.log("成功刷新书籍数据,书籍ID:" + bookID);
                        }
                    );
                return;
            }
            if (this.readerMode == "flip" || this.readerMode == "sketch") {
                // 命名路由,并加上参数,让路由建立 url
                this.$router.push({ name: "FlipMode", params: { id: bookID } });
            }
            if (this.readerMode == "scroll") {
                // 命名路由,并加上参数,让路由建立 url
                this.$router.push({ name: "ScrollMode", params: { id: bookID } });
            }
        },

        refreshReadMode() {
            // 初始化或者路由变化时,读取其他页面的更改,并存储到本地存储的阅读器模式（ReaderMode）这个值,
            if (localStorage.getItem("ReaderModeIsScroll") == "true") {
                this.readerModeIsScroll = true;
                this.readerMode = "scroll";
                localStorage.setItem("ReaderMode", "scroll");
            }
            if (localStorage.getItem("ReaderModeIsScroll") == "false") {
                this.readerModeIsScroll = false;
                this.readerMode = "flip";
                localStorage.setItem("ReaderMode", "flip");
            }
        },
        // 切换下拉、翻页阅读模式
        setReaderModeIsScroll(value) {
            this.readerModeIsScroll = value;
            if (this.readerModeIsScroll == true) {
                this.readerMode = "scroll";
                localStorage.setItem("ReaderMode", "scroll");
            } else {
                this.readerMode = "flip";
                localStorage.setItem("ReaderMode", "flip");
            }
            localStorage.setItem("ReaderModeIsScroll", value);
        },
        // 设置背景色的时候
        onBackgroundColorChange(value) {
            this.model.backgroundColor = value;
            localStorage.setItem("BackgroundColor", value);
        },
        // 设置UI颜色的时候
        onInterfaceColorChange(value) {
            this.model.interfaceColor = value;
            localStorage.setItem("InterfaceColor", value);
        },
        // 书籍卡片是否显示文字标题
        setBookCardShowTitleFlag(value) {
            this.bookCardShowTitleFlag = value;
            localStorage.setItem("BookCardShowTitleFlag", value);
            // console.log("成功保存设置: BookCardShowTitleFlag=" + localStorage.getItem("BookCardShowTitleFlag"));
        },
        // 初始化或者路由变化时,更新本地BookShelf相关数据
        getBookShelfData() {
            if (this.$route.params.group_id) {
                // console.log("BookShelf getBookShelfData!  this.$route.params.id" + this.$route.params.id)
                this.getBooksGroupByBookID(this.$route.params.group_id);
                this.headerShowReturnIcon = true;
            } else {
                this.initBookShelf();
                this.headerShowReturnIcon = false;
            }
        },

        // 获取所有书籍信息
        initBookShelf() {
            //根据文件名、修改时间、文件大小等要素排序的参数
            var sort_image_by_str = ""
            if (this.$route.query.sort_by) {
                //有路由查询参数的时候，按照路由里的查询参数排序
                sort_image_by_str = "&sort_by=" + this.$route.query.sort_by
            } else if (this.resort_hint_key !== "") {
                //没有的时候，就按照本地的存储值或默认值排序
                sort_image_by_str = "&sort_by=" + this.resort_hint_key
            }
            axios
                .get("getlist?max_depth=1" + sort_image_by_str)
                .then((response) => (this.bookshelf = response.data))
                .finally(() => {
                    this.setBookShelfTitle();
                });
        },
        // 根据路由参数获取特定书籍组
        getBooksGroupByBookID(group_id) {
            // console.log("getBooksGroupByBookID bookID:" + group_id);
            //排序参数（文件名、修改时间、文件大小等）
            var sort_image_by_str = ""
            if (this.$route.query.sort_by) {
                //有路由查询参数的时候，按照路由里的查询参数排序
                sort_image_by_str = "&sort_by=" + this.$route.query.sort_by
            } else if (this.resort_hint_key !== "") {
                //没有的时候，就按照本地的存储值或默认值排序
                sort_image_by_str = "&sort_by=" + this.resort_hint_key
            }
            axios
                .get("getlist?book_group_book_id=" + group_id + sort_image_by_str)
                .then((response) => {
                    if (response.data[0].name != null) {
                        this.bookshelf = response.data;
                    }
                })
                .finally(() => {
                    this.setBookShelfTitle();
                });
        },
        // 设置书架名
        setBookShelfTitle() {
            // 阅读某本书的时候,当然不需要设置
            if (this.$route.params.id != null) {
                return;
            }
            // 设置当前深度,这个值目前没用到
            if (this.bookshelf[0].depth != null) {
                this.max_depth = this.bookshelf[0].depth;
            }
            if (
                this.bookshelf[0].parent_folder != null &&
                this.bookshelf[0].parent_folder != ""
            ) {
                document.title = this.bookshelf[0].parent_folder;
                this.bookShelfTitle = this.bookshelf[0].parent_folder;
            }
            // //默认显示服务器版本
            // if (this.$store.state.server_status.ServerName) {
            //     //设置浏览器标签标题
            //     document.title = this.$store.state.server_status.ServerName
            //     //设置Header标题
            //     this.bookShelfTitle = this.$store.state.server_status.ServerName
            // }
        },

        // 打开抽屉
        drawerActivate(place) {
            this.drawerActive = true;
            this.drawerPlacement = place;
        },
        // 关闭抽屉
        drawerDeactivate() {
            this.drawerActive = false;
        },
        // 开始素描模式
        startSketchMode() {
            // this.$emit("setTemplate", "sketch");
        },
        // 接收Draw的参数,继续往父组件传
        OnSetReaderMode(value) {
            if (value == "scroll" || value == "scroll") this.readerMode = value;
        },
        // 如果在一个组件上使用了 v-model:xxx,应该使用 @update:xxx  https://www.naiveui.com/zh-CN/os-theme/docs/common-issues
        saveConfigToLocal() {
            // 储存cookie
            localStorage.setItem("BackgroundColor", this.model.backgroundColor);
            localStorage.setItem("InterfaceColor", this.model.interfaceColor);
        },

        // 根据可视区域(viewport)长宽比,确认是横屏还是竖屏
        // aspect-ratio https://developer.mozilla.org/zh-CN/docs/Web/CSS/@media/aspect-ratio
        // window.innerWidth  不是响应式依赖,所以不能用计算属性
        inLandscapeModeCheck() {
            this.aspectRatio = window.innerWidth / window.innerHeight;
            // 为了测试方便,阈值是正方形
            return this.aspectRatio > 19 / 19;
        },
    },
    computed: {},
});
</script>

<style scoped>
.header {
    background: v-bind("model.interfaceColor");
}

.bottom {
    background: v-bind("model.interfaceColor");
}

.shelf {
    max-width: 100%;
    min-height: 90vh;
    height: auto;
    background: v-bind("model.backgroundColor");
}
</style>
