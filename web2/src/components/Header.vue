<template>
  <header class="header">
    <select v-model="selectedTemplate" @change="onChange()">
      <option disabled value>切换阅读模式</option>
      <option>scroll</option>
      <option>single</option>
      <!-- <option>sketch</option> -->
    </select>
    <slot></slot>
  </header>
</template>

<script>
import { useCookies } from "vue3-cookies";
export default {
  name: "Header",
  setup() {
    const { cookies } = useCookies();
    return { cookies };
  },
  data() {
    return {
      mybook: this.book,
      selectedTemplate: "",
    };
  },

  methods: {
    onChange() {
      if (this.selectedTemplate === "scroll") {
        this.cookies.set("nowTemplate", "scroll");
      }
      if (this.selectedTemplate === "single") {
        this.cookies.set("nowTemplate", "single");
      }
      if (this.selectedTemplate === "sketch") {
        this.cookies.set("nowTemplate", "sketch");
      }
      location.reload(); //暂时无法动态刷新，研究好了再去掉
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.header {
  background: rgb(229, 229, 229);
  color: #111;
  text-align: center;
  padding: 2px;
  font-size: 12px;
}
.header select {
  background: rgb(228, 248, 174);
  color: #111;
  float: left;
  padding: 5px;
      margin-left: auto;
    margin-right: auto;
}
</style>




