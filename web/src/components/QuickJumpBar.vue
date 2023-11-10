<template>
  <div class="w-full my-3 flex flex-row justify-center content-center">
    <div class=" text-gray-900 h-12 py-4 self-center text-base">
      <a href="" class="text-blue-700 hover:underline text-3xl font-semibold">⬅️</a>
    </div>
    <select class="mx-2 p-2 w-1/2 border-gray-200 rounded-lg text-xl font-semibold
     text-center disabled:opacity-50 disabled:pointer-events-none"
      onchange="location = '/#/scroll/'+this.value;location.reload();">
      <option v-for="book in quick_jump_info.BaseBooks" :value="book.id" :key="book.id" :selected="book.id == nowBookID">
        {{ book.name}}
      </option>
    </select>
    <div class=" text-gray-900 h-12 py-4 space-x-2 text-base content-center">
      <a href="" class="text-blue-700 hover:underline text-3xl font-semibold">➡️</a>
    </div>
  </div>
</template>

<script lang ="ts">
import { defineComponent } from 'vue'
import axios from "axios";
export default defineComponent({
  name: 'QuickJumpBar',
  props: ['nowBookID'],
  data() {
    return {
      SomeFlag: "",
      quick_jump_info: {
        BaseBooks: [
          {
            id: 0,
            name: "",
          }
        ]
      },
    };
  },
  created() {
    axios
      .get("/quick_jump_info?id=" + this.$route.params.id)
      .then((response) => {
        //请求接口成功的逻辑
        this.quick_jump_info = response.data;
        console.log(this.quick_jump_info);
        console.log(this.quick_jump_info.BaseBooks);
        console.log(this.quick_jump_info.BaseBooks.length);
        //debugger
      }).catch((error: any) => {
        console.log(error);
      });
  },
  methods: {
    //选择书籍的时候，跳转到对应章节
    onQuickJump() {
      if (this.nowBookID == 0) {
        return;
      }
      console.log("onQuickJump");
      window.location.href = '/#/scroll/' + this.nowBookID;
      window.location.reload();
    },
  },
});
</script>

<style scoped></style>




