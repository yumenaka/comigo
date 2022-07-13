<!-- Observer.vue Observer组件,下拉刷新用 -->
<template>
    <div class="manga">
        <img v-lazy="this.image_url" />
        <div class="page_hint" v-if="this.showPageNumFlag_ScrollMode">{{ this.pageNum + 1 }}/{{
                this.all_page_num
        }}</div>
    </div>
</template>

<script>
// 使用Intersection Observer API构建无限滚动组件
// https://developer.mozilla.org/zh-CN/docs/Web/API/Intersection_Observer_API
// https://www.w3cplus.com/vue/build-an-infinite-scroll-component-using-intersection-observer-api.html 
// https://vueschool.io/articles/vuejs-tutorials/build-an-infinite-scroll-component-using-intersection-observer-api/
export default {
    props: ['options', 'image_url', 'pageNum', 'all_page_num', "showPageNumFlag_ScrollMode", "sPWL", "dPWL", "sPWP", "dPWP"],
    data: () => ({
        observer: null,
        tempThreshold: 0,
        entering: false,
    }),
    mounted() {
        const options = {
            root: null,//指定根 (root) 元素，用于检查目标的可见性。必须是目标元素的父级元素。如果未指定或者为null，则默认为浏览器视窗。。document.querySelector('#scrollArea')
            rootMargin: '0px',//根 (root) 元素的外边距。类似于 CSS 中的 margin 属性
            threshold: [0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1],//可以是单一的 number 也可以是 number 数组，target 元素和 root 元素相交程度达到该值的时候 IntersectionObserver 注册的回调函数将会被执行。
        }
        //let observer = new IntersectionObserver(callback, options);
        // https://developer.mozilla.org/en-US/docs/Web/API/IntersectionObserverEntry/intersectionRatio
        this.observer = new IntersectionObserver(([entry]) => {
            if (entry && entry.isIntersecting) {//如果相交了
                this.$emit("intersect");
                // console.log("observer intersect");
            }
            if(this.tempThreshold<entry.intersectionRatio){
                this.entering=true;
            }
            this.tempThreshold=entry.intersectionRatio;
            //如果交叉为0，则目标不可视,我们无需做任何事情。
            if (entry.intersectionRatio <= 0) return;
            //如果正在离开，不做任何事情
            // if (!this.entering) return;
            console.log("this.pageNum:", this.pageNum, " entry.intersectionRatio:", entry.intersectionRatio);

        }, options);
        this.observer.observe(this.$el);//使用this.$el作为root元素以便观察DOM元素。$el指向当前组件的DOM元素。this.$el在mounted中才会出现的，在created的时候没有。
    },
    beforeUnmount() {
        this.observer.disconnect();
    },
};
</script>


<style scoped>
/* https://developer.mozilla.org/zh-CN/docs/Web/CSS/object-fit */
.manga img {
    margin: auto;
    /* object-fit: scale-down; */
    padding: 3px 0px;
    border-radius: 7px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

.manga canvas {
    margin: auto;
    /* object-fit: scale-down; */
    padding: 3px 0px;
    border-radius: 7px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

.page_hint {
    /* 文字颜色 */
    color: #413d3d;
    /* 文字阴影：https://www.w3school.com.cn/css/css3_shadows.asp*/
    text-shadow: -1px 0 rgb(240, 229, 229), 0 1px rgb(253, 242, 242),
        1px 0 rgb(206, 183, 183), 0 -1px rgb(196, 175, 175);
}


.LoadingImage {
    width: 90vw;
    max-width: 90vw;
}

.ErrorImage {
    width: 90vw;
    max-width: 90vw;
}

/* 横屏（显示区域）时的CSS样式,IE无效 */
@media screen and (min-aspect-ratio: 19/19) {
    .SinglePageImage {
        width: v-bind(sPWL);
        max-width: 100%;
    }

    .DoublePageImage {
        width: v-bind(dPWL);
        max-width: 100%;
    }
}

/* 竖屏(显示区域)CSS样式,IE无效 */
@media screen and (max-aspect-ratio: 19/19) {
    .SinglePageImage {
        width: v-bind(sPWP);
        max-width: 100%;
    }

    .DoublePageImage {
        /* width: 100%; */
        width: v-bind(dPWP);
        max-width: 100%;
    }
}
</style>