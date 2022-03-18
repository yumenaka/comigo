<template>
    <n-result status="404" title="404 资源不存在" description="生活总归带点荒谬">
        <template #footer>
            <n-button>找点乐子吧</n-button>
        </template>
    </n-result>
</template>

<script>
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
import { NButton, NResult } from 'naive-ui'
import { useCookies } from "vue3-cookies";
import { defineComponent } from 'vue'
export default defineComponent({
    name: "404Page", //默认为 default。如果 <router-view>设置了名称，则会渲染对应的路由配置中 components 下的相应组件。
    components: {
        NButton,
        NResult,
    },
    setup() {
        const { cookies } = useCookies();
        return { cookies };
    },
});
</script>

<style scoped>
</style>
