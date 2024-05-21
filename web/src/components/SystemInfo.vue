

<template>
  <div id="SystemInfo" class="w-4/6">
    <n-space vertical>
      <p>CPU:</p>
      <n-progress type="circle" :percentage="parseFloat(cpu_percentage.toString())" />
      <p>{{ $store.state.server_status.OSInfo.cpu_num_physical }}Core
        {{ $store.state.server_status.OSInfo.cpu_num_logical_total }}Thread</p>
      <n-divider></n-divider>

      <p>RAM:</p>
      <n-progress type="circle" :percentage="parseFloat(ram_percentage.toString())" color="#ffaa66" />
      <p>{{ ($store.state.server_status.OSInfo.memory_total*($store.state.server_status.OSInfo.memory_used_percent/100)/
      (1024 * 1024 * 1024)).toFixed(2) }}GB/{{ ($store.state.server_status.OSInfo.memory_total / (1024 * 1024 *
        1024)).toFixed(2) }}GB</p>
      <n-divider></n-divider>

      <p>Books:{{ $store.state.server_status.NumberOfBooks }}</p>
      <p>Upload:{{ $store.state.server_status.SupportUploadFile }}</p>
      <!-- <p>Devices:{{ $store.state.server_status.NumberOfOnLineDevices }}</p> -->
      <!-- <p>cpu_num_physical:{{$store.state.server_status.OSInfo.cpu_num_physical}}</p>
      <p>cpu_num_logical_total:{{$store.state.server_status.OSInfo.cpu_num_logical_total}}</p> -->

    </n-space>
  </div>
</template>

<script lang="ts">
import { NProgress, useMessage, NSpace, NDivider, } from "naive-ui";
import { defineComponent, ref } from 'vue'
export default defineComponent({
  name: "AboutPage",
  props: ['showSystemInfo'],
  emits: ["setSome"],
  components: {
    NProgress,//进度条：https://www.naiveui.com/zh-CN/os-theme/components/progress
    NSpace,
    NDivider,//分割线
  },
  setup() {
    const cpu_percentageRef = ref(0);
    const ram_percentageRef = ref(0);
    const message = useMessage()
    return {
      message,
      cpu_percentage: cpu_percentageRef,
      ram_percentage: ram_percentageRef,
    }
  },
  data() {
    return {
      readerMode: "",
      upLoadHint: "",
      // cpu_percentage: 0.01,
      // ram_percentage: 0.01,
    };
  },
  created() {
    // this.beforeBookNum = this.$store.state.server_status.NumberOfBooks;
    this.$store.dispatch("syncSeverStatusDataAction");
    this.showStatus();
  },
  //挂载前
  beforeMount() {
  },
  onMounted() {
  },
  //卸载前
  beforeUnmount() {
  },
  methods: {
    //开始展示  api/server_info
    // ServerName	"Comigo v0.8.5"
    // ServerHost	"192.168.3.145"
    // ServerPort	1234
    // NumberOfBooks	79
    // NumberOfOnLineUser	1
    // NumberOfOnLineDevices	1
    // SupportUploadFile	true
    // ClientIP	"127.0.0.1"

    // OSInfo	
    // cpu_num_logical_total	16
    // cpu_num_physical	8
    // cpu_used_percent	63.60414336345319
    // memory_total	66889330688
    // memory_free	14006677504
    // memory_used_percent	36.62989823337489
    // description	"linux amd64"

    showStatus() {
      // this.message.success("");
      //每次上传完成后，触发轮询的次数
      const pollTimer = setInterval(() => {
        //服务器拉取最新状态，看是否新加了书籍
        this.$store.dispatch("syncSeverStatusDataAllAction");
        //console.log(this.$store.state.server_status.OSInfo.cpu_used_percent);
        this.cpu_percentage = this.$store.state.server_status.OSInfo.cpu_used_percent.toFixed(2);
        this.ram_percentage = this.$store.state.server_status.OSInfo.memory_used_percent.toFixed(2);
        if (!this.showSystemInfo) {
          clearInterval(pollTimer);
        }
      }, 1000);
    },
    onRefreshPage() {
      location.reload();
    },
  },
  computed: {
  },
});
</script>

<style scoped>
#UploadFile {
  padding: 20px;
}

.n-upload {
  --tw-bg-opacity: 0.5;
  background-color: rgba(249, 250, 251, var(--tw-bg-opacity));
}
</style>
