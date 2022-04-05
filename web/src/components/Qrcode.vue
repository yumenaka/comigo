<template>
    <qrcode-vue :value="getValue()" :size="size" level="H" />
</template>
<script>
import QrcodeVue from 'qrcode.vue'

export default {
    data() {
        return {
            // value:"test",
            size: 200,
        }
    },
    components: {
        QrcodeVue,
    },
    methods: {
        // 回首页
        getValue() {
            // 当前URL
            const url = document.location.toString();
            // 按照“/”分割字符串
            const arrUrl = url.split("/");
            const host = this.$store.state.server_status.ServerHostAndPort
            return url.replace(arrUrl[2], host);
        },
    },
}
</script>