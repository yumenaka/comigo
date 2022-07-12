<!-- Observer.vue Observer组件,下拉刷新用 -->
<template>
    <div class="observer" />
</template>

<script>
// 使用Intersection Observer API构建无限滚动组件
// https://www.w3cplus.com/vue/build-an-infinite-scroll-component-using-intersection-observer-api.html 
// https://vueschool.io/articles/vuejs-tutorials/build-an-infinite-scroll-component-using-intersection-observer-api/
export default {
    props: ['options'],
    data: () => ({
        observer: null,
    }),
    mounted() {
        const options = this.options || {};
        this.observer = new IntersectionObserver(([entry]) => {
            if (entry && entry.isIntersecting) {//如果相交了
                this.$emit("intersect");
                // console.log("observer intersect");
            }
        }, options);
        this.observer.observe(this.$el);//使用this.$el作为root元素以便观察DOM元素。$el指向当前组件的DOM元素。this.$el在mounted中才会出现的，在created的时候没有。
    },
    beforeUnmount() {
        this.observer.disconnect();
    },
};
</script>