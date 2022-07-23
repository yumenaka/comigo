<template>
    <div class="w-full h-screen flex flex-col">
        <Header class="header flex-none h-12" :bookIsFolder="false" :headerTitle="this.getUploadTitile()"
            :showReturnIcon="true" :showSettingsIcon="false"  :bookID='null' :setDownLoadLink="false">
        </Header>
        <!-- 渲染书架部分 有书的时候显示书  没有的时候显示上传控件-->
        <!-- Flex Grow 控制 flex 项目放大的功能类 https://www.tailwindcss.cn/docs/flex-grow -->
        <!-- 上传控件 -->
        <div class="mian_area flex-grow">
        </div>
        <Bottom class="bottom flex-none h-12" :softVersion="
            this.$store.state.server_status.ServerName
                ? this.$store.state.server_status.ServerName
                : 'Comigo'
        "></Bottom>
    </div>
</template>

<script>
import Header from "@/components/Header.vue";
import Bottom from "@/components/Bottom.vue";

import { defineComponent, reactive } from "vue";
export default defineComponent({
    name: "LoginPage",
    props: ['readMode'],
    emits: ["setSome"],
    components: {
        Header, // 自定义页头
        Bottom, // 自定义页尾
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
            PageTitle: "",
        };
    },
    created() {
        // 当前颜色
        if (localStorage.getItem("BackgroundColor") != null) {
            this.model.backgroundColor = localStorage.getItem("BackgroundColor");
        }
        if (localStorage.getItem("InterfaceColor") != null) {
            this.model.interfaceColor = localStorage.getItem("InterfaceColor");
        }
    },
    methods: {
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
