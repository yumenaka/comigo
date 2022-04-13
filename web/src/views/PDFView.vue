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
    </Drawer>

    <PDF :url="this.pdfUrl"></PDF>

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
import { NIcon, } from "naive-ui";

// https://juejin.cn/post/6995856687106261000
export default {
    name: 'PDFView',
    components: {
        Header,//自定义页头
        Drawer,//自定义抽屉
        Bottom,//自定义页尾
        PDF,
        SettingsOutline, //图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）
        NIcon, //图标  https://www.naiveui.com/zh-CN/os-theme/components/icon
    },
    setup() {
    },
    data() {
        return {
            pdfUrl: "",
            drawerActive: false,
            readerMode: "pdfview",
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
            // localStorage.setItem("debugModeFlag", this.debugModeFlag);
        },
    },

}
</script>