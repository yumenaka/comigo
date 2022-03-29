<template>
    <!-- class="w-28 md:w-33 lg:w-48"   Width of 28 by default, 32 on medium screens, and 48 on large screens -->
    <!-- 响应式设计：https://www.tailwindcss.cn/docs/responsive-design -->
    <!-- sm<640px  md<768px lg<1024px  lg<1280px 2xl<1536px-->
    <!-- 宽，只有一些典型值：https://www.tailwindcss.cn/docs/width  -->
    <!-- 高，只有一些典型值：https://www.tailwindcss.cn/docs/height -->
    <!-- 什么是REM：https://www.runoob.com/w3cnote/px-em-rem-different.html -->
    <!-- 边框圆角 rounded-xl：https://www.tailwindcss.cn/docs/border-radius -->
    <!-- 盒阴影 shadow: https://www.tailwindcss.cn/docs/box-shadow -->
    <!-- 外边距 m-x m-y  https://www.tailwindcss.cn/docs/margin -->
    <div
        class="relative w-32 h-44 mx-3 my-4 text-gray-800 dark:text-gray-300 rounded shadow-xl hover:shadow-2xl ring-1 ring-gray-400 hover:ring hover:ring-blue-500 .bg-top bg-cover"
        :style="getBackgroundImageStyle()"
    >
        <div v-if="this.childBookNum != ''" class="absolute inset-x-0 top-0 text-right">
            <span
                class="text-2xl text-yellow-500 font-extrabold text-shadow"
            >{{ this.childBookNum }}</span>
        </div>
        <!-- 文本对齐:       https://www.tailwindcss.cn/docs/text-align -->
        <!-- 定位:          https://www.tailwindcss.cn/docs/top-right-bottom-left -->
        <!-- 背景色:        https://www.tailwindcss.cn/docs/background-color -->
        <!-- 背景色不透明度: https://www.tailwindcss.cn/docs/background-opacity -->
        <!-- 文本溢出：      https://www.tailwindcss.cn/docs/text-overflow -->
        <!-- 字体粗细：     https://www.tailwindcss.cn/docs/font-weight -->
        <div
            v-if="this.showTitle"
            class="absolute inset-x-0 bottom-0 h-1/4 bg-gray-200 bg-opacity-70 font-semibold border-blue-800 rounded-b"
        >
            <span class="absolute inset-x-0 bottom-0 align-middle">{{ this.shortTitle }}</span>
        </div>
    </div>
</template>

<script>
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
// import { NCard, } from 'naive-ui'
import { useCookies } from "vue3-cookies";
import { defineComponent } from 'vue'
export default defineComponent({
    name: "BookCover",
    props: ['title', 'image_src', 'id', 'readerMode', 'showTitle', 'childBookNum'],
    components: {
        // NCard,
        // NEllipsis,
    },
    setup() {
        const { cookies } = useCookies();
        return { cookies };
    },
    computed: {
        shortTitle() {
            if (this.title.length <= 17) {
                return this.title
            } else {
                return this.title.substr(0, 17) + "…"
            }
        },
    },
    data() {
        return {
            // resize_str: "&resize_height=340",   &auto_crop=1
            resize_str: "&resize_width=256&resize_height=360&resize_cut=true",
        };
    },
    methods: {
        //回首页
        onBackTop() {
            // 字符串路径
            this.$router.push('/')
        },
        getBackgroundImageStyle() {
            return "background-image: url(" + this.getThumbnailsImageUrl() + ");"
        },
        getThumbnailsImageUrl() {
            //按照“/”分割字符串
            var arrUrl = this.image_src.split("/");
            // console.log(arrUrl)
            if (arrUrl[0] == "api") {
                return this.image_src + "&resize_width=256&resize_height=360&resize_cut=true"
            } else {
                return this.image_src
            }
        },
        //自己构建一个<a>链接，后来发现不如可以直接用router-link与命名路由
        getHref() {
            //当前URL
            var url = document.location.toString();
            //按照“/”分割字符串
            var arrUrl = url.split("/");
            //拼一个完整的图片URL
            if (this.readerMode == "flip") {
                var new_url = arrUrl[0] + "//" + arrUrl[2] + "/#" + "f/" + this.id
            }
            if (this.readerMode == "scroll") {
                new_url = arrUrl[0] + "//" + arrUrl[2] + "/#" + "s/" + this.id
            }
            console.log(new_url)
            return new_url
        },
    },
});
</script>
// 
<style scoped>
</style>
