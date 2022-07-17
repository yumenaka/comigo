<!-- Observer.vue Observer组件,下拉刷新用 -->
<template>
    <div class="manga">
        <img v-lazy="this.image_url" />
        <div class="page_hint" v-if="this.showPageNumFlag_ScrollMode">{{ this.pageNum }}/{{
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
    props: ['options', "syncPageFlag", 'sendWSMessage', 'book_id', 'image_url', 'pageNum', 'all_page_num', "showPageNumFlag_ScrollMode", "sPWL", "dPWL", "sPWP", "dPWP"],
    emits: ['refreshNowPageNum'],
    data: () => ({
        observer: null,
        tempThreshold: 0,
        sendMessageNum: 0,
    }),
    mounted() {
        const options = {
            root: null,//指定根 (root) 元素，用于检查目标的可见性。必须是目标元素的父级元素。如果未指定或者为null，则默认为浏览器视窗。。document.querySelector('#scrollArea')
            rootMargin: '0px',//根 (root) 元素的外边距。类似于 CSS 中的 margin 属性
            threshold: [0, 0.01, 0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09, 0.1, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16, 0.17, 0.18, 0.19, 0.2, 0.21, 0.22, 0.23, 0.24, 0.25, 0.26, 0.27, 0.28, 0.29, 0.3, 0.31, 0.32, 0.33, 0.34, 0.35, 0.36, 0.37, 0.38, 0.39, 0.4, 0.41, 0.42, 0.43, 0.44, 0.45, 0.46, 0.47, 0.48, 0.49, 0.5, 0.51, 0.52, 0.53, 0.54, 0.55, 0.56, 0.57, 0.58, 0.59, 0.6, 0.61, 0.62, 0.63, 0.64, 0.65, 0.66, 0.67, 0.68, 0.69, 0.7, 0.71, 0.72, 0.73, 0.74, 0.75, 0.76, 0.77, 0.78, 0.79, 0.8, 0.81, 0.82, 0.83, 0.84, 0.85, 0.86, 0.87, 0.88, 0.89, 0.9, 0.91, 0.92, 0.93, 0.94, 0.95, 0.96, 0.97, 0.98, 0.99, 1],//可以是单一的 number 也可以是 number 数组，target 元素和 root 元素相交程度达到该值的时候 IntersectionObserver 注册的回调函数将会被执行。
        }
        //let observer = new IntersectionObserver(callback, options);
        // https://developer.mozilla.org/en-US/docs/Web/API/IntersectionObserverEntry/intersectionRatio
        this.observer = new IntersectionObserver(([entry]) => {
            //isIntersecting：被监视对象与root是否交叉。如果没交叉，不做任何事情
            if (!entry.isIntersecting) {
                this.sendMessageNum = 0;
                return
            }
            // DOMRect 对象，提供元素的大小及其相对于视口的位置。
            // https://developer.mozilla.org/zh-CN/docs/Web/API/Element/getBoundingClientRect
            // console.dir(entry.boundingClientRect);
            //viewport
            // console.dir(entry.target);
            // const viewportHeight = entry.target.clientHeight;
            // console.log("viewport_clientHeight:", viewport_clientHeight);
            // //以底部为准：从下往上，进入viewport的时候才增加页数
            // if (entry.boundingClientRect.bottom >= viewportHeight) {
            //以顶部为准：从下往上，接触顶部的时候才改变页数。因为要留一些提前量，所以不是0而是50.
            if (entry.boundingClientRect.top <= 50) {
                this.$emit("refreshNowPageNum", this.pageNum);
                // console.log("pageNum", this.pageNum);
                //如果发送过一次了，就不再重复发送，避免大量发包
                if (this.syncPageFlag && this.sendMessageNum === 0) {
                    this.sendNowPage(entry.intersectionRatio);
                }
            }
        }, options);
        this.observer.observe(this.$el);//使用this.$el作为root元素以便观察DOM元素。$el指向当前组件的DOM元素。this.$el在mounted中才会出现的，在created的时候没有。
    },
    beforeUnmount() {
        this.observer.disconnect();
    },
    methods: {
        //Websocket 发送消息
        sendNowPage(now_page_num_percent) {
            this.sendMessageNum = this.sendMessageNum + 1
            const readPercent =
                parseFloat(this.pageNum) / parseFloat(this.all_page_num);
            // console.debug("ReadPercent: " + readPercent)
            const data = {
                book_id: this.book_id,
                now_page_num: this.pageNum,
                now_page_num_percent: now_page_num_percent,
                read_percent: readPercent,
            };
            // console.log("this.$store.userID: " + this.$store.state.userID)
            const newMsg = {
                type: "sync_page",
                status_code: 200,
                user_id: this.$store.state.userID,
                token: this.$store.state.token,
                detail: "翻页模式，发送数据", // 消息描述
                data_string: JSON.stringify(data),
            };
            // 配置为了json，可调用sendObj方法来发送数据
            if (this.$socket.readyState === 1) {
                this.$socket.sendObj(newMsg);
            }
            // console.log("send:", newMsg); 
        },
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