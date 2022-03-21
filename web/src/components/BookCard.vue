<template>
    <n-card :title="this.title" class="book_card">
        <template #cover>
            <img :src="this.image_src + this.resize_str" />
        </template>
    </n-card>
</template>

<script>
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
import { NCard } from 'naive-ui'
import { useCookies } from "vue3-cookies";
import { defineComponent } from 'vue'
export default defineComponent({
    name: "BookCover",
    props: ['title', 'image_src', 'id'],
    components: {
        NCard,
    },
    setup() {
        const { cookies } = useCookies();
        return { cookies };
    },
    data() {
        return {
            // resize_str: "&resize_height=340",
            resize_str: "&resize_cut=true&resize_width=240&resize_height=340",
        };
    },
});
</script>
// 
<style scoped>
.n-card {
    /* max-width: 100px; */
    /* width: 40px;
     */
    width: 240px;
}
</style>
