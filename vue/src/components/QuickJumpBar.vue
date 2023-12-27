<template>
  <div class="quick_jump_bar w-full my-0 flex flex-row justify-center items-center">
    <a :href="prevLink" class="text-blue-700  text-2xl font-semibold">⬅️</a>
    <select
      class="mx-4 px-2 py-4 w-3/4 border-gray-200 rounded-lg text-xl font-semibold text-center disabled:opacity-50 disabled:pointer-events-none"
      :onchange="handleChange">
      <option v-for="book in group_info_filter.BookInfos" :value="book.id" :key="book.id" :selected="book.id == nowBookID">
        {{ book.title }}
      </option>
    </select>
    <a :href="nextLink" class="text-blue-700  text-2xl font-semibold ">➡️</a>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import axios from 'axios';

export default defineComponent({
  name: 'QuickJumpBar',
  props: ['nowBookID','readMode'],
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

<style scoped></style>




