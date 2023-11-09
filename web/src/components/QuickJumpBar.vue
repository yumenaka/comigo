<template>
  <div class=" text-gray-900 h-12 py-4 space-x-2 text-base content-center">
    <a href="" class="text-blue-700 hover:underline font-semibold">⬅️</a>
  </div>
  <select onchange="location = this.value;">
    <option value="ss">01話 </option>
    <option value="ss">02話 </option>
    <!-- <option v-for="item in quickJumpList" :value="item.url" :key="item.id" :selected="item.id == nowBookID">{{ item.name
    }}
    </option> -->
  </select>
  <div class=" text-gray-900 h-12 py-4 space-x-2 text-base content-center">
    <a href="" class="text-blue-700 hover:underline font-semibold">➡️</a>
  </div>
</template>

<script lang ="ts">
import { defineComponent } from 'vue'
import axios from "axios";
export default defineComponent({
  name: 'QuickJumpBar',
  props: [],
  data() {
    return {
      SomeFlag: "",
      quick_jump_info: {},
    };
  },
  created() {
    // console.log("QuickJumpBar created");
		axios
			.get("/quick_jump_info?id=" + this.$route.params.id)
			.then((response) => {
				//请求接口成功的逻辑
				this.quick_jump_info = response.data;
				console.log(this.quick_jump_info);
        console.log(this.quick_jump_info.BaseBooks);
        console.log(this.quick_jump_info.BaseBooks.length);
        debugger
			}).catch((error: any) => {
          console.log(error);
       });
	},
  methods: {
    //选择书籍的时候，跳转到对应章节
    onQuickJump() {

    },
  },
});
</script>

<style scoped></style>




