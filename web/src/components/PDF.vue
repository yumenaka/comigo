<template>
    <div class="pdfCanvas" v-for="item in pageNum" :key="item">
        <canvas :id="`pdf-canvas-${item}`" class="pdf-page" />
    </div>
</template>
<script>

// canvas结合pdfjs-dist在vue项目中预览pdf文件 
// https://juejin.cn/post/7025960857427247111
// 1、pdfjs需要特定版本（链接中说是2.2.228，测试后,发现2.9.359正常展示不报错……）

const PDF = require('pdfjs-dist')

// https://juejin.cn/post/6995856687106261000
export default {
    name: 'DisplayPdf',
    props: {
        url: { type: String },
        index: {
            type: Number
        },
    },
    components: {
    },
    setup() {
    },
    data() {
        return {
            // url: "",
            pageNum: 0,
            pdfCtx: null,
            book: {
                name: "loading",
                all_page_num: 2,
                book_type: "dir",
                pages: [
                    {
                        height: 500,
                        width: 449,
                        url: "/images/loading.jpg",
                    },
                    {
                        height: 500,
                        width: 449,
                        url: "/images/loading.jpg",
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
            // NOTE: 这一步要特别注意，网上很多关于pdfjs的使用教程里漏了这一步，会出现workerSrc未定义的错误
            PDF.GlobalWorkerOptions.workerSrc = require('pdfjs-dist/build/pdf.worker.entry')
            // console.log("!!!!!!!!!!!!!" + this.url)
            const loadingTask = PDF.getDocument(this.url)
            // var _this = this
            loadingTask.promise.then((pdf) => {
                this.pdfCtx = pdf
                this.pageNum = pdf.numPages
                this.$nextTick(() => {
                    this.renderPdf()
                })
            })
        },
        renderPdf(num = 1) {
            this.pdfCtx.getPage(num).then(page => {
                const canvas = document.getElementById(`pdf-canvas-${num}`)
                const ctx = canvas.getContext('2d')
                const viewport = page.getViewport({ scale: 1 })
                // 画布大小,默认值是width:300px,height:150px
                canvas.height = viewport.height
                canvas.width = viewport.width
                // // 画布的dom大小, 设置移动端,宽度设置铺满整个屏幕
                // const clientWidth = document.body.clientWidth
                // canvas.style.width = clientWidth + 'px'
                // // 根据pdf每页的宽高比例设置canvas的高度
                // canvas.style.height = clientWidth * (viewport.height / viewport.width) + 'px'
                page.render({
                    canvasContext: ctx,
                    viewport
                })
                if (num < this.pageNum) {
                    this.renderPdf(num + 1)
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