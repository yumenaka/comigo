<template>
  <header class="header">
    <slot></slot>
    <!-- <div slot="start">手动调整图片宽度（{{rangeValue}}）： </div>
    <div slot="end">2400</div>-->

    <!-- <div v-for="b in bookshelf" :key="b.uuid" class="bookshelf">
      <v-btn @click="onChangeBook(b.uuid)" class="book_button">{{ b.name }} ({{ b.page_num }})</v-btn>
    </div>-->
  </header>
</template>

<script>
export default {
  name: "Header",

  data() {
    return {
      mybook:this.book,
      rangeValue: 500,
    };
  },

  methods: {
    changeEvent() {
      console.log("change to:" + this.rangeValue);
    },
  },
};
</script>

<style scoped>
.header {
  background: rgb(69, 224, 177);
  color: #111;
  text-align: center;
  padding: 10px;
}
.header a {
  padding-right: 5px;
  text-decoration: none;
}
</style>
