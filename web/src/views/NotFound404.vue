<template>
    <n-result status="404" title="404 资源不存在" description="没找到你想要的东西">
        <template #footer>
            <n-button @click="onBackTop()">回首页╮(￣▽￣")╭</n-button>
        </template>
    </n-result>
</template>

<script lang="ts">
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
    methods: {
        //回到首页
        onBackTop() {
            // 字符串路径
            this.$router.push('/')

        },
    },
});
</script>

<style scoped>
.n-result {
    /* 内边距区域 */
    padding-bottom: 0px;
    padding-left: 3px;
    padding-right: 3px;
    padding-top: 200px;

    align-content: center;
    width: 100%;
    height: 100vh;
    border-radius: 7px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}
</style>
