<template>
  <div class="home" id="ScrollPage">
    <Header v-if="this.showHeader">
      <h2>
        <a v-if="!book.IsFolder" :href="'raw/' + book.name"
          >{{ book.name }}【Download】</a
        >
        <a v-if="book.IsFolder" :href="'raw/' + book.name">{{
          book.name
        }}</a>
      </h2>
    </Header>
    <!-- <BookScroll></BookScroll> -->

    <div v-for="(page, key) in book.pages" :key="page.url" class="manga">
      <!-- <p v-if="showPageNum">{{ key + 1 }}/{{ AllPageNum }}</p> -->
      <!-- v-lazy="page.url"  :src="page.url" -->
      <img
        :src="page.url" 
        v-bind:H="page.height"
        v-bind:W="page.width"
        v-bind:key="key"
        v-bind:class="page.image_type"
      />
      <p v-if="showPageNum">{{ key + 1 }}/{{ book.all_page_num }}</p>
    </div>
    <button>▲</button>
  </div>
</template>

<script>
// @ is an alias to /src
import Header from "@/components/Header.vue";
export default {
  name: "Home", //默认为 default。如果 <router-view>设置了名称，则会渲染对应的路由配置中 components 下的相应组件。
  components: {
    Header,
    // BookScroll,
  },
  setup() {},
  created() {
    console.log("created!");
    this.axios
      .get("/book.json")
      .then((response) => {
        if (response.status == 200) {
          this.book = response.data;
          console.log(this.book);
        }
      })
      .catch((error) => alert(error));
  },
  //组件的 data 选项必须是一个函数
  //每个实例可以维护一份被返回对象的独立的拷贝
  data() {
    return {
      book: null,
      showHeader: true,
      showPageNum: true,
    };
  },
  computed: {
    // 计算属性的 getter
    nowTemplate: function () {
      var localValue = this.$cookies.get("nowTemplate");
      console.log("computed 1:" + localValue);
      if (localValue !== null) {
        return localValue;
      } else {
        return this.$store.state.setting.template;
      }
    },
  },
  methods: {
    getNumber: function (number) {
      this.page = number;
      console.log(number);
    },
    getNowTemplate: function () {
      var localValue = this.$cookies.get("nowTemplate");
      console.log("computed 1:" + localValue);
      if (localValue !== null) {
        return localValue;
      } else {
        this.$cookies.set("nowTemplate", this.$store.state.setting.template);
        console.log("computed 2:" + this.$store.state.setting.template);
        return this.$store.state.setting.template;
      }
    },
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  background-color: #f6f7eb;
  align-items: center;
}

.manga img {
  margin: auto;
  max-width: inherit%;
  padding-top: 3px;
  padding-bottom: 3px;
  padding-right: 0px;
  padding-left: 0px;
  border-radius: 7px;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

/* 竖屏(显示区域)CSS样式，IE无效 */
@media screen and (max-aspect-ratio: 19/19) {
  .SinglePage {
    width: 100%;
  }
  .DoublePage {
    width: 100%;
  }
}

/* 横屏（显示区域）时的CSS样式，IE无效 */
@media screen and (min-aspect-ratio: 19/19) {
  .SinglePage {
    width: 900px;
  }
  .DoublePage {
    width: 95%;
  }
}

/* 高分横屏（显示区域）时的CSS样式，IE无效 */
/* min-width 输出设备中页面最小可视区域宽度 大于这个width时，其中的css起作用 超宽屏 */
@media screen and (min-aspect-ratio: 19/19) and (min-width: 1922px) {
  .SinglePage {
    width: 1000px;
  }
  .DoublePage {
    width: 1900px;
  }
}
</style>
