<template>
  <header class="header">
  <!-- <button v-on:click="change_template_to_scroll">scroll</button>
  <button v-on:click="change_template_to_double">double</button>
  <button v-on:click="change_template_to_sketch">sketch</button>
  <button v-on:click="change_template_to_single">single</button> -->
  <select v-model="selectedTemplate"  @change='onChange()'  >
    <option disabled value="">点此切换阅读模式</option>
    <option>scroll</option>
    <option>double</option>
    <option>single</option>
    <option>sketch</option>
  </select>
  <slot></slot>
  </header>
</template>

<script>
export default {
  name: "Header",

  data() {
    return {
      mybook: this.book,
      selectedTemplate: ''
    };
  },

  methods: {
    onChange() {
      if (this.selectedTemplate ==="scroll") {
        this.change_template_to_scroll();
      }
      if (this.selectedTemplate ==="double") {
        this.change_template_to_double();
      }
      if (this.selectedTemplate ==="single") {
        this.change_template_to_single();
      }
      if (this.selectedTemplate ==="sketch") {
        this.change_template_to_sketch();
      }
      //console.log(this.$store.state.setting.template);
    },
    change_template_to_scroll() {
      this.$store.commit("change_template_to_scroll");
    },
    change_template_to_double() {
      this.$store.commit("change_template_to_double");
    },
    change_template_to_single() {
      this.$store.commit("change_template_to_single");
    },
    change_template_to_sketch() {
      this.$store.commit("change_template_to_sketch");
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

.header select {
  background: rgb(192, 205, 238);
  color: #111;
  font-size: 20px;
  float: left;
}

</style>
