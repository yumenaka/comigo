<template>
  <header class="header">
    <mt-range
            v-model="rangeValue"
            :min="0"
            :max="2400"
            :step="20"
            :bar-height="5"
            @change="changeEvent()"
            >
      <div slot="start">手动调整图片宽度（{{rangeValue}}）： </div>
      <div slot="end">2400</div>
    </mt-range>
  </header>
</template>

<script>
export default {
  name: "Header",
  data() {
      return {
          rangeValue:500,
      };
  },
  methods:{
      changeEvent(){
          console.log("change to:"+this.rangeValue);
      }
  }
};
</script>

<style scoped>
.header {
  background: #333;
  color: #fff;
  text-align: center;
  padding: 10px;
}
.header a {
  color: #fff;
  padding-right: 5px;
  text-decoration: none;
}
</style>
