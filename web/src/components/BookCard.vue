<template>
    <n-card class="book_card">
        <template #cover>
            <img :src="this.image_src + this.resize_str" />
        </template>
        <span tooltip="false">{{ this.ellipsisTitle }}</span>
    </n-card>
</template>

<script>
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
import { NCard, } from 'naive-ui'
import { useCookies } from "vue3-cookies";
import { defineComponent } from 'vue'
export default defineComponent({
    name: "BookCover",
    props: ['title', 'image_src', 'id'],
    components: {
        NCard,
        // NEllipsis,
    },
    setup() {
        const { cookies } = useCookies();
        return { cookies };
    },
    computed: {
        ellipsisTitle() {
            if (this.title.length <= 15) {
                return this.title
            } else {
                return this.title.substr(0, 15) + "……"
            }
        },
    },
    data() {
        return {
            // resize_str: "&resize_height=340",
            resize_str: "&resize_width=128&resize_height=180&resize_cut=true",
        };
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

    width: 134px;
    height: 200px;
    border-radius: 6px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}
</style>
