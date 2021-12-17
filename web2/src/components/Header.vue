<template>
  <header class="header">
    <slot></slot>
  </header>
</template>

<script>
export default {
  name: "Header",

  data() {
    return {
      mybook: this.book,
      selectedTemplate: "",
    };
  },

  methods: {
    onChange() {
      if (this.selectedTemplate === "scroll") {
        this.change_template_to_scroll();
      }
      if (this.selectedTemplate === "double") {
        this.change_template_to_double();
      }
      if (this.selectedTemplate === "single") {
        this.change_template_to_single();
      }
      if (this.selectedTemplate === "sketch") {
        this.change_template_to_sketch();
      }
      location.reload(); //暂时无法动态刷新，研究好了再去掉
      //console.log(this.$store.state.setting.template);
    },
    change_template_to_scroll() {
      this.$cookies.set("nowTemplate", "scroll");
      this.$store.commit("change_template_to_scroll");
    },
    change_template_to_double() {
      this.$cookies.set("nowTemplate", "double");
      this.$store.commit("change_template_to_double");
    },
    change_template_to_single() {
      this.$cookies.set("nowTemplate", "single");
      this.$store.commit("change_template_to_single");
    },
    change_template_to_sketch() {
      this.$cookies.set("nowTemplate", "sketch");
      this.$store.commit("change_template_to_sketch");
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
  padding: 10px;
}
.header a {
  padding-right: 5px;
  text-decoration: none;
}

.header select {
  background: rgb(228, 248, 174);
  color: #111;
  font-size: 23px;
  float: left;
}

.header button {
  background: rgb(192, 205, 238);
  color: #111;
  font-size: 23px;
  float: left;
}
</style>




