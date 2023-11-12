<template>
  <div class="w-full my-3 flex flex-row justify-center items-center">
    <a :href="prevLink" class="text-blue-700 hover:underline text-4xl font-semibold">⬅️</a>
    <select
      class="mx-2 p-2 w-1/2 border-gray-200 rounded-lg text-xl font-semibold text-center disabled:opacity-50 disabled:pointer-events-none"
      onchange="location = '/#/scroll/'+this.value;location.reload();">
      <option v-for="book in same_group_books.BaseBooks" :value="book.id" :key="book.id" :selected="book.id == nowBookID">
        {{ book.name }}
      </option>
    </select>
    <a :href="nextLink" class="text-blue-700 hover:underline text-4xl font-semibold ">➡️</a>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import axios from 'axios';

export default defineComponent({
  name: 'QuickJumpBar',
  props: ['nowBookID'],
  data() {
    return {
      SomeFlag: '',
      same_group_books: {
        BaseBooks: [
          {
            id: 0,
            name: '',
          },
        ],
      },
      selectedBook: "",
    };
  },
  computed: {
    prevLink() {
      for (let i = 0; i < this.same_group_books.BaseBooks.length; i++) {
        if (this.same_group_books.BaseBooks[i].id === this.nowBookID) {
          if (i === 0) {
            return `/#/scroll/${this.nowBookID}`;
          }
          return `/#/scroll/${this.same_group_books.BaseBooks[i - 1].id}`;
        }
      }
    },
    nextLink() {
      for (let i = 0; i < this.same_group_books.BaseBooks.length; i++) {
        if (this.same_group_books.BaseBooks[i].id === this.nowBookID) {
          if (i === this.same_group_books.BaseBooks.length - 1) {
            return `/#/scroll/${this.nowBookID}`;
          }
          return `/#/scroll/${this.same_group_books.BaseBooks[i + 1].id}`;
        }
      }
    },
  },
  created() {
    this.selectedBook = this.nowBookID;
    this.fetchQuickJumpInfo();
  },
  methods: {
    async fetchQuickJumpInfo() {
      try {
        const response = await axios.get(`/same_group_books?id=${this.$route.params.id}`);
        this.same_group_books = response.data;
      } catch (error) {
        console.log(error);
      }
    }
  },
});
</script>

<style scoped></style>




