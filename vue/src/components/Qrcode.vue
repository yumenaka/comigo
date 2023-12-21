<template>
    <!-- <qrcode-vue :value="getValue()" :size="size" level="H" /> -->
    <n-image width="200" :src="this.getValue()" />
</template>
<script>
// import QrcodeVue from 'qrcode.vue'  //https://github.com/fengyuanchen/vue-qrcode  //yarn add qrcode.vue  //yarn remove qrcode.vue
import { NImage, } from 'naive-ui'
export default {
    data() {
        return {
            // value:"test",
            size: 200,
        }
    },
    components: {
        // QrcodeVue,
        NImage,
    },
    methods: {
        // 回首页
        getValue() {
            // 当前URL
            const url = document.location.toString();
            // 按照“/”分割字符串
            const arrUrl = url.split("/");
            let host = this.$store.state.server_status.ServerHost
            //添加本地连接里的端口部分
            if (arrUrl[2].indexOf(":") != -1) {
                let pose = arrUrl[2].indexOf(":")
                host = this.$store.state.server_status.ServerHost + arrUrl[2].slice(pose)
            }
            //如果 URI 组件中含有分隔符，比如 ? 和 #，则应当使用 encodeURIComponent() 方法分别对各组件进行编码
            const qrcode_src = "api/qrcode.png?qrcode_str=" + encodeURIComponent(url.replace(arrUrl[2], host))
            // console.log(qrcode_src)
            return qrcode_src
        },
    },
}
</script>