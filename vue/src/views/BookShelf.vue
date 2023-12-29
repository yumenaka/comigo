<template>
    <div class="BookShelf w-full h-screen flex flex-col">
        <Header class="flex-none h-12" in-shelf="true" :bookIsFolder="false" :headerTitle="bookShelfTitle"
            :showReSortIcon="true" :showReturnIcon="headerShowReturnIcon" :showSettingsIcon="true"
            :bookID="$route.params.group_id" :depth="bookshelf[0].depth > 1 ? bookshelf[0].depth - 1 : 0" :setDownLoadLink="false"
            @drawerActivate="drawerActivate" @onResort="onResort">
        </Header>

        <!-- Flex Grow 控制 flex 项目放大的功能类 https://www.tailwindcss.cn/docs/flex-grow -->
        <div v-if="bookCardMode == 'gird'" class="bookshelf flex-grow">
            <!-- 使用tailwindcss提供的flex布局： -->
            <!-- flex-row：https://www.tailwindcss.cn/docs/flex-direction -->
            <div class="flex flex-row flex-wrap justify-center min-h-48">
                <BookCard v-for="(book_info, key) in bookshelf" :key="key" :book_info="book_info"
                    :bookCardMode="bookCardMode" :simplifyTitle="simplifyTitle" :readerMode="readerMode"
                    :showTitle="bookCardShowTitleFlag">
                </BookCard>
            </div>
        </div>

        <div v-if="bookCardMode == 'list'" class="bookshelf flex-grow flex flex-col justify-center items-center">
            <BookList v-for="(book_info, key) in bookshelf" :key="key" :book_info="book_info" :simplifyTitle="simplifyTitle"
                :showTitle="bookCardShowTitleFlag" :readerMode="readerMode">
            </BookList>
        </div>

        <div v-if="bookCardMode == 'text'" class="bookshelf flex-grow">
            <div class="flex flex-row flex-wrap mt-4 mb-4 px-4 justify-left items-center min-w-4">
                <BookText v-for="(book_info, key) in bookshelf" :key="key" :book_info="book_info" :readerMode="readerMode">
                </BookText>
            </div>
        </div>

        <!-- 返回顶部的圆形按钮，向上滑动的时候出现 -->
        <n-back-top class="bg-blue-200" :show="showBackTopFlag" type="info" :right="20" :bottom="20" />

        <Bottom class="flex-none h-12" :ServerName="$store.state.server_status.ServerName
            ? $store.state.server_status.ServerName
            : 'Comigo'
            "></Bottom>
        <Drawer :initDrawerActive="drawerActive" :initDrawerPlacement="drawerPlacement" @saveConfig="saveConfigToLocal"
            @startSketch="startSketchMode" @closeDrawer="drawerDeactivate" @setRM="OnSetReaderMode" :readerMode="readerMode"
            :sketching="false" :inBookShelf="true">
            <!-- 设置颜色 -->
            <span>{{ $t("setInterfaceColor") }}</span>
            <n-color-picker v-model:value="model.interfaceColor" :modes="['hex']" :show-alpha="false"
                @update:value="onInterfaceColorChange" />
            <span>{{ $t("setBackColor") }}</span>
            <n-color-picker v-model:value="model.backgroundColor" :modes="['hex']" :show-alpha="false"
                @update:value="onBackgroundColorChange" />

            <!-- 开关：下拉阅读 -->
            <n-switch size="large" :rail-style="railStyle" v-model:value="readerModeIsScroll"
                @update:value="setReaderModeIsScroll">
                <template #checked>{{ $t("scroll_mode") }}</template>
                <template #unchecked>{{ $t("flip_mode") }}</template>
            </n-switch>

            <!-- 开关：显示书名 -->
            <n-switch size="large" v-model:value="bookCardShowTitleFlag" @update:value="setBookCardShowTitleFlag">
                <template #checked>{{ $t("show_book_titles") }}</template>
                <template #unchecked>{{ $t("show_book_titles") }}</template>
            </n-switch>

            <!-- 开关：简化书名 -->
            <n-switch size="large" v-model:value="simplifyTitle" @update:value="setSimplifyTitle">
                <template #checked>{{ $t("simplify_book_titles") }}</template>
                <template #unchecked>{{ $t("simplify_book_titles") }}</template>
            </n-switch>

            <!-- 页面重新排序
            <n-select :placeholder="$t('re_sort_book')" @update:value="onResort" :options="options" />
            <p>&nbsp;</p> -->
            <!-- 下载示例配置文件的按钮 -->
            <a href="api/config.toml" target="_blank">
                <n-button>{{ $t("DownloadSampleConfigFile") }}</n-button>
            </a>
            <!-- 下载windows reg文件的按钮 -->
            <a v-if="remoteIsWindows()" href="api/comigo.reg" target="_blank">
                <n-button>{{ $t("DownloadWindowsRegFile") }}</n-button>
            </a>
        </Drawer>
    </div>
</template>

<script lang="ts">
// 直接导入组件并使用它。这种情况下,只有导入的组件才会被打包。
// Firefox插件Textarea Cache 报错：源映射错误：Error: NetworkError when attempting to fetch resource.
// Firefox插件Video DownloadHelper报错:已不赞成使用 CanvasRenderingContext2D 中的 drawWindow 方法
import { NColorPicker, NSwitch, NButton, NSelect, NBackTop, } from "naive-ui";
import Header from "@/components/Header.vue";
import Drawer from "@/components/Drawer.vue";
import BookCard from "@/components/BookCard.vue";
import BookList from "@/components/BookList.vue";
import BookText from "@/components/BookText.vue";
import Bottom from "@/components/Bottom.vue";
import { CSSProperties, defineComponent, reactive } from "vue";
import { useCookies } from "vue3-cookies"; // https://github.com/KanHarI/vue3-cookies
import axios from "axios";

export default defineComponent({
    name: "BookShelf",
    props: ["testProp"],
    emits: ["setTemplate"],
    components: {
        Header, // 自定义页头
        Drawer, // 自定义抽屉
        BookCard, //自定义书本
        BookList,
        BookText,
        Bottom, // 自定义页尾
        NButton, //按钮,来自:https://www.naiveui.com/zh-CN/os-theme/components/button
        NSwitch,
        // NDivider, // 分割线  https://www.naiveui.com/zh-CN/os-theme/components/divider
        NColorPicker, //颜色选择器 https://www.naiveui.com/zh-CN/os-theme/components/color-picker
        // NDropdown,//下拉菜单 https://www.naiveui.com/zh-CN/os-theme/components/dropdown
        NSelect,
        NBackTop, //回到顶部按钮,来自:https://www.naiveui.com/zh-CN/os-theme/components/back-top
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
            railStyle: ({
                focused,
                checked,
            }: {
                focused: boolean;
                checked: boolean;
            }) => {
                const style: CSSProperties = {};
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
            bookCardMode: "gird",//gird,list,text
            simplifyTitle: true, //简化显示标题
            //是否显示回到顶部按钮
            showBackTopFlag: false,
            //是否正在向下滚动
            scrollDownFlag: false,
            //存储现在滚动的位置
            scrollTopSave: 0,
            resort_hint_key: "filename", //书籍的排序方式。可以按照文件名、修改时间、文件大小排序（或反向排序）
            readerMode: "scroll",
            readerModeIsScroll: true,
            bookShelfTitle: "Loading",
            headerShowReturnIcon: false,
            bookCardShowTitleFlag: true, // 书库中的书籍是否显示文字版标题
            maxDepth: 1,
            bookshelf: [
                {
                    title: "loading",
                    page_count: 1,
                    id: "12345",
                    type: ".zip",
                    parent_folder: "",
                    depth: 1,
                    child_book_num: 0,
                    pages: [
                        {
                            height: 500,
                            width: 449,
                            url: "/images/loading.gif",
                        },
                    ],
                    cover: {
                        filename: "loading.jpg",
                        height: 500,
                        width: 449,
                        url: "/images/loading.gif",
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
        //监听滚动，返回顶部按钮用
        window.addEventListener("scroll", this.onScroll);
        // 初始化默认值,读取出来的都是字符串,不要直接用
        //书籍排序方式。可以按照文件名、修改时间、文件大小排序（或反向排序）
        const tempBookShelf_SortBy = localStorage.getItem("BookShelf_SortBy");
        if (typeof tempBookShelf_SortBy === "string") {
            this.resort_hint_key = tempBookShelf_SortBy;
        }
        // 从服务器上拉取书架信息
        this.setBookShelf();
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
                this.setBookShelf();
            }
        );
        // 继续初始化默认值
        // 书籍卡片是否显示文字版标题
        if (localStorage.getItem("BookCardShowTitleFlag") === "true") {
            this.bookCardShowTitleFlag = true;
        } else if (localStorage.getItem("BookCardShowTitleFlag") === "false") {
            this.bookCardShowTitleFlag = false;
        }

        // 简化标题
        if (localStorage.getItem("SimplifyTitle") === "true") {
            this.simplifyTitle = true;
        } else if (localStorage.getItem("SimplifyTitle") === "false") {
            this.simplifyTitle = false;
        }

        // 当前颜色
        let tempBackgroundColor = localStorage.getItem("BackgroundColor");
        if (typeof tempBackgroundColor === "string") {
            this.model.backgroundColor = tempBackgroundColor;
        }

        let tempInterfaceColor = localStorage.getItem("InterfaceColor");
        if (typeof tempInterfaceColor === "string") {
            this.model.interfaceColor = tempInterfaceColor;
        }

        let tempBookCardMode = localStorage.getItem("bookCardMode");
        if (typeof tempBookCardMode === "string") {
            this.bookCardMode = tempBookCardMode;
        }
    },
    // 挂载前
    beforeMount() {
        // this.$store.dispatch("syncSeverStatusDataAction");
    },
    // 在绑定元素的父组件被挂载后调用。
    mounted() {
        this.isLandscapeMode = this.inLandscapeModeCheck();
        this.isPortraitMode = !this.inLandscapeModeCheck();
    },

    // 卸载前
    beforeUnmount() {
        //组件销毁前,销毁监听事件
        window.removeEventListener("scroll", this.onScroll);
    },
    methods: {
        // 初始化或者路由变化时,更新本地BookShelf相关数据
        setBookShelf() {
            if (!this.$route.params.group_id) {
                this.initBookShelf();
                return;
            }
            const temp_id = this.$route.params.group_id;
            if (typeof temp_id === "string") {
                this.getBooksGroupByBookID(temp_id);
            }
            this.headerShowReturnIcon = true;
        },

        // 获取所有书籍信息
        initBookShelf() {
            //根据文件名、修改时间、文件大小等要素排序的参数
            let sort_by = "";
            //按照本地的存储值或默认值排序
            if (this.resort_hint_key !== "") {
                sort_by = "&sort_by=" + this.resort_hint_key;
            }
            //按照路由里的查询参数排序
            if (this.$route.query.sort_by) {
                sort_by = "&sort_by=" + this.$route.query.sort_by;
            }
            let _this = this;
            axios
                .get("book_infos?depth=0" + sort_by)
                .then((response) => {
                    if (response.data !== "") {
                        this.bookshelf = response.data;
                    } else {
                        this.bookShelfTitle = _this.$t("no_book_found_hint");
                        this.$router.push({
                            name: "UploadPage",
                        });
                    }
                }).catch((error) => {
                    console.log("请求接口失败" + error);
                })
                .finally(() => {
                    this.setBookShelfTitle();
                });
            this.headerShowReturnIcon = false;
        },
        // 根据路由参数获取特定书籍组
        getBooksGroupByBookID(group_id: string) {
            //排序参数（文件名、修改时间、文件大小等）
            var sort_image_by_str = "";
            //按照本地的存储值或默认值排序
            if (this.resort_hint_key !== "") {
                sort_image_by_str = "&sort_by=" + this.resort_hint_key;
            }
            //按照路由里的查询参数排序
            if (this.$route.query.sort_by) {
                sort_image_by_str = "&sort_by=" + this.$route.query.sort_by;
            }
            axios
                .get("book_infos?book_group_id=" + group_id + sort_image_by_str)
                .then((response) => {
                    this.bookshelf = response.data;
                }).catch((error) => {
                    console.log("请求接口失败" + error);
                })
                .finally(() => {
                    this.setBookShelfTitle();
                });
        },
        // 设置书架名
        setBookShelfTitle() {
            // console.log(this.$route.params.id);
            // 路由里面有id这个参数、也就是处于scroll或flip模式，正在阅读某本书的时候,不需要设置书架名
            if (this.$route.params.id !== undefined) {
                //不是null而是undefined
                return;
            }
            if (this.bookshelf === null) {
                return;
            }
            //如果是书籍组的时候，如何设置标题
            if (
                this.bookshelf[0].parent_folder !== null &&
                this.bookshelf[0].parent_folder !== ""
            ) {
                document.title = this.bookshelf[0].parent_folder;
                this.bookShelfTitle = this.bookshelf[0].parent_folder;
            }
        },

        //页面滚动的时候,改变返回顶部按钮的显隐
        onScroll() {
            let scrollTop =
                document.documentElement.scrollTop || document.body.scrollTop;
            this.scrollDownFlag = scrollTop > this.scrollTopSave;
            //防手抖,小于一定数值状态就不变 Math.abs()会导致报错
            let step = this.scrollTopSave - scrollTop;
            // console.log("this.scrollDownFlag:",this.scrollDownFlag,"scrollTop:",scrollTop,"step:", step);
            this.scrollTopSave = scrollTop;
            if (step < -5 || step > 5) {
                this.showBackTopFlag = scrollTop > 400 && !this.scrollDownFlag;
            }
        },
        //查看服务器是否windows，来决定显示不显示reg文件下载按钮
        remoteIsWindows() {
            // 非空判断
            if (!this.$store.state.server_status) {
                return false;
            }
            if (!this.$store.state.server_status.OSInfo) {
                return false;
            }
            if (!this.$store.state.server_status.OSInfo.description) {
                return false;
            }
            if (!this.$store.state.server_status.OSInfo) {
                return false;
            }
            return (
                this.$store.state.server_status.OSInfo.description.indexOf(
                    "windows"
                ) !== -1
            );
        },
        //根据文件名、修改时间、文件大小等参数重新排序
        onResort(key: string) {
            if (key === "gird" || key === "list" || key === "text") {
                this.bookCardMode = key;
                localStorage.setItem("bookCardMode", this.bookCardMode);
                return;
            }
            this.resort_hint_key = key;
            localStorage.setItem("BookShelf_SortBy", key);
            if (this.$route.params.group_id) {
                console.log(
                    "onResort  bookID：" + this.$route.params.group_id + ", key:" + key
                );
                this.$router.push({
                    name: "ChildBookShelf",
                    params: { group_id: this.$route.params.group_id },
                    replace: true,
                    query: { sort_by: key },
                });
            } else {
                console.log("onResort  key：" + key);
                this.$router.push({
                    name: "BookShelf",
                    replace: true,
                    query: { sort_by: key },
                });
            }
        },

        refreshReadMode() {
            // 初始化或者路由变化时,读取其他页面的更改,并存储到本地存储的阅读器模式（ReaderMode）这个值,
            if (localStorage.getItem("ReaderMode") === "scroll") {
                this.readerModeIsScroll = true;
                this.readerMode = "scroll";
                // localStorage.setItem("ReaderMode", "scroll");
            }
            if (localStorage.getItem("ReaderMode") === "flip") {
                this.readerModeIsScroll = false;
                this.readerMode = "flip";
                // localStorage.setItem("ReaderMode", "flip");
            }
        },

        // 切换下拉、翻页阅读模式
        setReaderModeIsScroll(value: boolean) {
            this.readerModeIsScroll = value;
            if (this.readerModeIsScroll) {
                this.readerMode = "scroll";
                localStorage.setItem("ReaderMode", "scroll");
            } else {
                this.readerMode = "flip";
                localStorage.setItem("ReaderMode", "flip");
            }
            localStorage.setItem("ReaderModeIsScroll", value ? "true" : "false");
        },
        // 设置背景色的时候
        onBackgroundColorChange(value: string) {
            this.model.backgroundColor = value;
            localStorage.setItem("BackgroundColor", value);
        },
        // 设置UI颜色的时候
        onInterfaceColorChange(value: string) {
            this.model.interfaceColor = value;
            localStorage.setItem("InterfaceColor", value);
        },
        // 书籍卡片是否显示文字标题
        setBookCardShowTitleFlag(value: boolean) {
            this.bookCardShowTitleFlag = value;
            localStorage.setItem("BookCardShowTitleFlag", value ? "true" : "false");
            // console.log("成功保存设置: BookCardShowTitleFlag=" + localStorage.getItem("BookCardShowTitleFlag"));
        },

        // 书籍卡片是否显示文字标题
        setSimplifyTitle(value: boolean) {
            this.simplifyTitle = value;
            localStorage.setItem("SimplifyTitle", value ? "true" : "false");
        },

        // 打开抽屉
        drawerActivate(place: string) {
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
        OnSetReaderMode(value: string) {
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

.bookshelf {
    background: v-bind("model.backgroundColor");
}
</style>
