<template>
    <div class="UploadPage w-full h-screen flex flex-col">
        <Header class="header flex-none h-12"  in-shelf="true" :bookIsFolder="false" :headerTitle="getUploadTitile()"
            :showReturnIcon="true" :showSettingsIcon="true" :bookID='null' :setDownLoadLink="false"
            @drawerActivate="drawerActivate">
        </Header>
        <!-- 渲染书架部分 有书的时候显示书  没有的时候显示上传控件-->
        <!-- Flex Grow 控制 flex 项目放大的功能类 https://www.tailwindcss.cn/docs/flex-grow -->
        <!-- 上传控件 -->
        <div class="mian_area flex-grow">
            <UploadFile>
            </UploadFile>
        </div>

        <Bottom class="bottom flex-none h-12" :ServerName="
            $store.state.server_status.ServerName
                ? $store.state.server_status.ServerName
                : 'Comigo'
        "></Bottom>

        <Drawer :initDrawerActive="drawerActive" :initDrawerPlacement="drawerPlacement" @closeDrawer="drawerDeactivate"
            :sketching="false" :inBookShelf="true">
            <SystemInfo :showSystemInfo="drawerActive">
            </SystemInfo>
        </Drawer>
    </div>
</template>

<script lang="ts">
import Header from "@/components/Header.vue";
import Bottom from "@/components/Bottom.vue";
import UploadFile from "@/components/UploadArea.vue";
import SystemInfo from "@/components/SystemInfo.vue";
import Drawer from "@/components/Drawer.vue";
import { defineComponent, reactive } from "vue";
export default defineComponent({
    name: "UploadPage",
    components: {
        Header, // 自定义页头
        Bottom, // 自定义页尾
        Drawer, // 自定义抽屉
        SystemInfo,//自定义系统信息
        UploadFile,//自定义的文件上传领域，一本书也没有的时候用
    },
    setup() {
        // 背景颜色,颜色选择器用  // 此处不能使用this
        const model = reactive({
            interfaceColor: "#F5F5E4",
            backgroundColor: "#E0D9CD",
        });
        return {
            model,
        };
    },

    data() {
        return {
            book_num: 0,
            drawerActive: false,
            drawerPlacement: "right",
            UploadPageTitle: "",
        };
    },
    created() {
        // 当前颜色
        const tempBackgroundColor = localStorage.getItem("BackgroundColor")
        if (typeof (tempBackgroundColor) === 'string') {
            this.model.backgroundColor = tempBackgroundColor;
        }
        const tempInterfaceColor = localStorage.getItem("tempInterfaceColor")
        if (typeof (tempInterfaceColor) === 'string') {
            this.model.interfaceColor = tempInterfaceColor
        }
    },
    methods: {
        // 打开抽屉
        drawerActivate(place: string) {
            this.drawerActive = true;
            this.drawerPlacement = place;
        },
        // 关闭抽屉
        drawerDeactivate() {
            this.drawerActive = false;
        },
        getUploadTitile() {
            //如果没有一本书
            if (this.$store.state.server_status.SupportUploadFile === false) {
                return this.$t('no_support_upload_file');
            }
            //如果没有一本书
            if (this.$store.state.server_status.NumberOfBooks === 0) {
                return this.$t('no_book_found_hint');
            }
            return this.$t('number_of_online_books') + this.$store.state.server_status.NumberOfBooks;
        },
        remoteIsWindows() {
            if (!this.$store.state.server_status) {
                return false
            }
            console.dir(this.$store.state.server_status.Description);
            return this.$store.state.server_status.Description.indexOf("windows") !== -1
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

.mian_area {
    background: v-bind("model.backgroundColor");
}
</style>
