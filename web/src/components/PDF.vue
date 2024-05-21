<template>
    <!-- 定位: https://www.tailwindcss.cn/docs/top-right-bottom-left -->
    <div class="pdfCanvas flex flex-col" v-for="item in pageNum" :key="item">
        <canvas :id="`pdf-canvas-${item}`" class="pdf-page" />
    </div>
</template>
<script>

// canvas结合pdfjs-dist在vue项目中预览pdf文件 
// https://juejin.cn/post/7025960857427247111

// 在线查看PDF文件,pdf.js使用方法
// https://www.cnblogs.com/zhanggf/p/8504317.html

// Vue3实现各种附件预览 
// https://juejin.cn/post/6995856687106261000

// 1、pdfjs需要特定版本（链接中说是2.2.228，测试后,发现2.9.359正常展示不报错……）

const PDF = require('pdfjs-dist')
// import useLazyData from '@/compositions'
// PDF.GlobalWorkerOptions.workerSrc = require('/pdf.worker.min.js')

export default {
    name: 'DisplayPdf',
    props: {
        url: { type: String },
        index: {
            type: Number
        },
        pdfPageWidth: {
            type: Number
        },
    },
    components: {
    },
    setup() {
        // const fn = () => {
        //     findNew().then(data => {
        //         console.log('findNew', data)
        //         list.value = data.result
        //     })
        // }
        // const target = useLazyLoad(fn)

        // return { target }
    },
    data() {
        return {
            // url: "",
            pageNum: 0,
            pdfCtx: null,
            book: {
                name: "loading",
                page_count: 2,
                book_type: "dir",
                pages: [
                    {
                        height: 500,
                        width: 449,
                        url: "/images/loading.gif",
                    },
                    {
                        height: 500,
                        width: 449,
                        url: "/images/loading.gif",
                    },
                ],
            },
        };
    },
    created() {
        // PDFJS.GlobalWorkerOptions.workerSrc = require("pdfjs-dist/build/pdf.worker.min.js");
    },
    beforeUpdate() {
        this.loadFile()
    },
    methods: {
        loadFile() {
            // 设定pdfjs的 workerSrc 参数
            // NOTE: 这一步要特别注意，漏了这一步，会出现workerSrc未定义的错误
            //打包后会出错？
            // PDF.GlobalWorkerOptions.workerSrc = require('pdfjs-dist/build/pdf.worker.entry')

            // console.log("!!!!!!!!!!!!!" + this.url)
            const loadingTask = PDF.getDocument(this.url)
            // var _this = this
            loadingTask.promise.then((pdf) => {
                this.pdfCtx = pdf
                this.pageNum = pdf.numPages
                this.$nextTick(() => {
                    console.log("PDF loaded");
                    this.renderPdf()
                })
            })
        },
        renderPdf(num = 1) {
            this.pdfCtx.getPage(num).then(page => {
                const canvas = document.getElementById(`pdf-canvas-${num}`)
                const ctx = canvas.getContext('2d')
                var viewport = page.getViewport({ scale: 1.0 });

                //默認dpi爲96，然後網頁一般使用dpi爲72。 
                var CSS_UNITS = 96.0 / 72.0;

                var s = this.pdfPageWidth / viewport.width;
                var scaledViewport = page.getViewport({ scale: s });//缩放过的PDF
                canvas.width = Math.floor(scaledViewport.width * CSS_UNITS)
                canvas.height = Math.floor(scaledViewport.height * CSS_UNITS);
                console.log("s:" + s + " canvas.width:" + canvas.width + " canvas.height:" + canvas.height)
                // 画布大小,默认值是width:300px,height:150px
                // canvas.height = viewport.height
                // canvas.width = viewport.width
                // // 画布的dom大小, 设置移动端,宽度设置铺满整个屏幕
                // const clientWidth = document.body.clientWidth
                // canvas.style.width = clientWidth + 'px'
                // // 根据pdf每页的宽高比例设置canvas的高度
                // canvas.style.height = clientWidth * (viewport.height / viewport.width) + 'px'

                page.render({
                    transform: [CSS_UNITS, 0, 0, CSS_UNITS, 0, 0],
                    canvasContext: ctx,
                    viewport
                })
                if (num < this.pageNum) {
                    // this.renderPdf(num + 1)
                    if (num < 10) {
                        this.renderPdf(num + 1)
                    }
                } else {
                    console.log('onRendered')
                }
            })
        }
    },
}
</script>


<style scoped>
.canvas {
    margin: auto;
    /* object-fit: scale-down; */
    padding: 3px 0px;
    border-radius: 7px;
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}
</style>