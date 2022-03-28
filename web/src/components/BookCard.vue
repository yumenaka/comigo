<template>
    <n-card class="book_card">
        <template #cover>
            <img :src="getThumbnailsImageUrl()" />
        </template>
        {{ this.shortTitle }}
    </n-card>
</template>

<script>
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
import { NCard, } from 'naive-ui'
import { useCookies } from "vue3-cookies";
import { defineComponent } from 'vue'
export default defineComponent({
    name: "BookCover",
    props: ['title', 'image_src', 'id', 'readerMode', 'showTitle'],
    components: {
        NCard,
        // NEllipsis,
    },
    setup() {
        const { cookies } = useCookies();
        return { cookies };
    },
    computed: {
        shortTitle() {
            if (this.title.length <= 12) {
                return this.title
            } else {
                return this.title.substr(0, 12) + "…"
            }
        },
    },
    data() {
        return {
            // resize_str: "&resize_height=340",
            resize_str: "&resize_width=256&resize_height=360&resize_cut=true",
        };
    },
    methods: {
        //回首页
        onBackTop() {
            // 字符串路径
            this.$router.push('/')
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
.n-card {
    /* 内边距区域 */
    padding-bottom: 0px;
    padding-left: 5px;
    padding-right: 5px;
    padding-top: 7px;

    width: 140px;
    height: 195px;
    border-radius: 6px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

/* 样式化链接 */
/* https://developer.mozilla.org/zh-CN/docs/Learn/CSS/Styling_text/Styling_links */
a {
    outline: none;
    text-decoration: none;
    padding: 2px 1px 0;
}

a:link {
    color: #265301;
}

a:visited {
    color: #437a16;
}

a:focus {
    border-bottom: 1px solid;
    background: #bae498;
}

a:hover {
    border-bottom: 1px solid;
    background: #cdfeaa;
}

a:active {
    background: #265301;
    color: #cdfeaa;
}
</style>
