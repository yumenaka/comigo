<template>
    <Header :bookIsFolder="false" :headerTitle="book.name" :showReturnIcon="true" :bookID="book.id"
        :setDownLoadLink="true">
        <!-- 右边的设置图标,点击屏幕中央也可以打开 -->
        <n-icon size="40" @click="drawerActivate('right')">
            <settings-outline />
        </n-icon>
    </Header>
    <Drawer :initDrawerActive="this.drawerActive" :initDrawerPlacement="this.drawerPlacement"
        @saveConfig="this.saveConfigToLocalStorage" @closeDrawer="this.drawerDeactivate" :readerMode="this.readerMode"
        :inBookShelf="false" :sketching="false">

        <!-- 单页-漫画宽度-使用固定值PX -->

        <!-- 数字输入框 -->
        <n-input-number size="small" :show-button="false" v-model:value="this.pdfPageWidth" :min="50"
            :update-value-on-input="false">
            <template #prefix>{{ $t('singlePageWidth') }}</template>
            <template #suffix>px</template>
        </n-input-number>
        <!-- 滑动选择PX -->
        <n-slider v-model:value="this.pdfPageWidth" :step="10" :max="1600" :min="50"
            :format-tooltip="value => `${value}px`" />

    </Drawer>

    <PDF :pdfPageWidth="this.pdfPageWidth" :url="this.pdfUrl"></PDF>

    <Bottom
        :softVersion="this.$store.state.server_status.ServerName ? this.$store.state.server_status.ServerName : 'Comigo'">
    </Bottom>
</template>
<script>
import Header from "@/components/Header.vue";
import Drawer from "@/components/Drawer.vue";
import Bottom from "@/components/Bottom.vue";
import PDF from "@/components/PDF.vue";
import axios from "axios";
import { SettingsOutline } from "@vicons/ionicons5";
import { defineComponent, reactive } from 'vue'
import { NSlider, NIcon, NInputNumber, } from 'naive-ui'
// https://juejin.cn/post/6995856687106261000
export default defineComponent({
    name: 'PDFView',
    components: {
        Header,//自定义页头
        Drawer,//自定义抽屉
        Bottom,//自定义页尾
        PDF,
        SettingsOutline, //图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）
        NIcon, //图标  https://www.naiveui.com/zh-CN/os-theme/components/icon
        NSlider,//滑动选择  Slider https://www.naiveui.com/zh-CN/os-theme/components/slider
        NInputNumber,//数字输入 https://www.naiveui.com/zh-CN/os-theme/components/input-number
    },
    setup() {
        //背景色
        //reactive({}) 创建并返回一个响应式对象: https://www.bilibili.com/video/av925511720/?p=4  也讲到了toRefs()
        const model = reactive({
            backgroundColor: "#E0D9CD",
            interfaceColor: "#f5f5e4",
        });
        return {
            //背景色
            model,
        }
    },
    data() {
        return {
            pdfUrl: "",
            drawerActive: false,
            readerMode: "pdfview",
            //横屏(Landscape)状态的漫画页宽度,PX
            pdfPageWidth: 720,
            book: {
                name: "loading",
                all_page_num: 2,
                book_type: "dir",
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
            drawerPlacement: 'right',
        };
    },
    created() {


        //漫画页宽度,PX
        if (localStorage.getItem("pdfPageWidth") != null) {
            let saveNum = Number(localStorage.getItem("pdfPageWidth"));
            if (!isNaN(saveNum)) {
                this.pdfPageWidth = saveNum;
            }
        }

        //根据路由参数获取特定书籍
        axios
            .get("/getbook?id=" + this.$route.params.id)
            .then((response) => (this.book = response.data))
            .finally(
                () => {
                    document.title = this.book.name;
                    console.log("成功获取书籍数据,书籍ID:" + this.$route.params.id);
                    if (this.book.book_type === ".pdf") {
                        this.pdfUrl = 'api/raw/' + this.$route.params.id + '/' + this.book.name
                        console.log("pdfUrl:" + this.pdfUrl);
                    }
                }
            );
    },
    methods: {
        //打开抽屉
        drawerActivate(place) {
            this.drawerActive = true;
            this.drawerPlacement = place;
        },
        //关闭抽屉
        drawerDeactivate() {
            this.drawerActive = false;
        },
        // 关闭抽屉时,保存设置到cookies
        saveConfigToLocalStorage() {
            localStorage.setItem("pdfPageWidth", this.pdfPageWidth);
        },
    },

});
</script>


<style scoped>
.header {
    background: v-bind("model.interfaceColor");
}

.bottom {
    background: v-bind("model.interfaceColor");
}

.manga {
    max-width: 100%;
    background: v-bind("model.backgroundColor");
}

/* https://developer.mozilla.org/zh-CN/docs/Web/CSS/object-fit */
.manga img {
    margin: auto;
    /* object-fit: scale-down; */
    padding: 3px 0px;
    border-radius: 7px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

.manga canvas {
    margin: auto;
    /* object-fit: scale-down; */
    padding: 3px 0px;
    border-radius: 7px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

.page_hint {
    /* 文字颜色 */
    color: #413d3d;
    /* 文字阴影：https://www.w3school.com.cn/css/css3_shadows.asp*/
    text-shadow: -1px 0 rgb(240, 229, 229), 0 1px rgb(253, 242, 242),
        1px 0 rgb(206, 183, 183), 0 -1px rgb(196, 175, 175);
}

.LoadingImage {
    width: 90vw;
    max-width: 90vw;
}

.ErrorImage {
    width: 90vw;
    max-width: 90vw;
}

/* 横屏（显示区域）时的CSS样式,IE无效 */
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

/* 竖屏(显示区域)CSS样式,IE无效 */
@media screen and (max-aspect-ratio: 19/19) {
    .SinglePageImage {
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
