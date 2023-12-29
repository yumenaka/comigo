<template>
  <div class="quick_jump_bar w-full my-0 flex flex-row justify-center items-center">
    <!-- 左箭头 -->
    <a :href="prevLink" class="text-blue-700  text-2xl font-semibold">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-7 h-7">
        <path fill-rule="evenodd"
          d="M11.03 3.97a.75.75 0 0 1 0 1.06l-6.22 6.22H21a.75.75 0 0 1 0 1.5H4.81l6.22 6.22a.75.75 0 1 1-1.06 1.06l-7.5-7.5a.75.75 0 0 1 0-1.06l7.5-7.5a.75.75 0 0 1 1.06 0Z"
          clip-rule="evenodd" />
      </svg>
    </a>
    
    <!-- 选择框 -->
    <select
      class="rounded mx-4 px-3 py-1.5 w-3/4 border border-gray-400 text-lg font-semibold text-center  disabled:opacity-50 disabled:pointer-events-none"
      :onchange="handleChange">
      <option class=" border-2 border-red-500" v-for="book in group_info_filter.BookInfos" :value="book.id" :key="book.id"
        :selected="book.id == nowBookID">
        {{ book.title }}
      </option>
    </select>
    <!-- 右箭头 -->
    <a :href="nextLink" class="text-blue-700  text-2xl font-semibold ">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
        class="w-7 h-7">
        <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3" />
      </svg>
    </a>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import axios from 'axios';

export default defineComponent({
  name: 'QuickJumpBar',
  props: ['nowBookID', 'readMode'],
  data() {
    return {
      SomeFlag: 'filename',
      group_info_filter: {
        BookInfos: [
          {
            id: 0,
            title: '',
            type: '',
          },
        ],
      },
      selectedBook: "",
    };
  },
  computed: {
    prevLink() {
      for (let i = 0; i < this.group_info_filter.BookInfos.length; i++) {
        if (this.group_info_filter.BookInfos[i].id === this.nowBookID) {
          if (i === 0) {
            return `/#/${this.readMode}/${this.nowBookID}`;
          }
          return `/#/${this.readMode}/${this.group_info_filter.BookInfos[i - 1].id}`;
        }
      }
    },
    nextLink() {
      for (let i = 0; i < this.group_info_filter.BookInfos.length; i++) {
        if (this.group_info_filter.BookInfos[i].id === this.nowBookID) {
          if (i === this.group_info_filter.BookInfos.length - 1) {
            return `/#/${this.readMode}/${this.nowBookID}`;
          }
          return `/#/${this.readMode}/${this.group_info_filter.BookInfos[i + 1].id}`;
        }
      }
    },
  },
  created() {
    this.selectedBook = this.nowBookID;
    this.fetchQuickJumpInfo();
  },
  methods: {
    handleChange(event: Event) {
      const target = event.target as HTMLInputElement;
      location.href = `/\#/${this.readMode}/` + target.value;
      location.reload();
    },
    async fetchQuickJumpInfo() {
      try {
        const response = await axios.get(`/group_info_filter?id=${this.$route.params.id}`);
        this.group_info_filter = response.data;
      } catch (error) {
        console.log(error);
      }
    }
  },
});
</script>

<style scoped>

</style>




