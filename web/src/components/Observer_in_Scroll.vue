<!-- Observer.vue Observer组件,下拉刷新用 -->
<template>
    <div class="observer" />
</template>

<script>
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
            }
        }, options);
        this.observer.observe(this.$el);
    },
    beforeUnmount() {
        this.observer.disconnect();
    },
};
</script>