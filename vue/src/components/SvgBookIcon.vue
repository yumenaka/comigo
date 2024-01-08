<template>
    <!-- 书籍组：显示书籍数量 -->
    <div v-if="childBookNum != ''" class="absolute top-0 z-10 right-1">
        <span class="flex text-xl text-yellow-500 font-black text-shadow">
            <svg class="opacity-60 w-7 h-7 " xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="rgb(234 179 8)"
                stroke-width="1.0" stroke="black">
                <path
                    d="M11.25 4.533A9.707 9.707 0 0 0 6 3a9.735 9.735 0 0 0-3.25.555.75.75 0 0 0-.5.707v14.25a.75.75 0 0 0 1 .707A8.237 8.237 0 0 1 6 18.75c1.995 0 3.823.707 5.25 1.886V4.533ZM12.75 20.636A8.214 8.214 0 0 1 18 18.75c.966 0 1.89.166 2.75.47a.75.75 0 0 0 1-.708V4.262a.75.75 0 0 0-.5-.707A9.735 9.735 0 0 0 18 3a9.707 9.707 0 0 0-5.25 1.533v16.103Z" />
            </svg>{{ childBookNum }}</span>
    </div>

    <!-- 文件夹 -->
    <i class="icon-[material-symbols--arrow-back-ios-new]"></i>

    <svg v-if='book_info.type == "dir"'
        class="opacity-70 w-7 h-7 m-0 absolute top-1 z-10 right-1 shadow-2xl hover:shadow-2xl"
        xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 512 512">
        <path d="M432 112V96a48.14 48.14 0 0 0-48-48H64a48.14 48.14 0 0 0-48 48v256a48.14 48.14 0 0 0 48 48h16"
            fill="rgb(234 179 8)" stroke="currentColor" stroke-linejoin="round" stroke-width="32"></path>
        <rect x="96" y="128" width="400" height="336" rx="45.99" ry="45.99" fill="rgb(234 179 8)" stroke="currentColor"
            stroke-linejoin="round" stroke-width="32"></rect>
        <ellipse cx="372.92" cy="219.64" rx="30.77" ry="30.55" fill="none" stroke="currentColor" stroke-miterlimit="10"
            stroke-width="32"></ellipse>
        <path d="M342.15 372.17L255 285.78a30.93 30.93 0 0 0-42.18-1.21L96 387.64" fill="none" stroke="currentColor"
            stroke-linecap="round" stroke-linejoin="round" stroke-width="32"></path>
        <path d="M265.23 464l118.59-117.73a31 31 0 0 1 41.46-1.87L496 402.91" fill="none" stroke="currentColor"
            stroke-linecap="round" stroke-linejoin="round" stroke-width="32"></path>
    </svg>

    <!-- 视频  -->
    <svg v-if='book_info.type == "video"'
        class="opacity-70 w-7 h-7 m-0 absolute top-1 z-10 right-1 shadow-2xl hover:shadow-2xl"
        xmlns="http://www.w3.org/2000/svg" fill="rgb(234 179 8)" viewBox="0 0 24 24" stroke-width="1.0"
        stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round"
            d="m15.75 10.5 4.72-4.72a.75.75 0 0 1 1.28.53v11.38a.75.75 0 0 1-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 0 0 2.25-2.25v-9a2.25 2.25 0 0 0-2.25-2.25h-9A2.25 2.25 0 0 0 2.25 7.5v9a2.25 2.25 0 0 0 2.25 2.25Z" />
    </svg>

    <!-- 音乐  -->
    <svg v-if='book_info.type == "audio"'
        class="opacity-70 w-7 h-7 m-0 absolute top-1 z-10 right-1 shadow-2xl hover:shadow-2xl"
        xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="rgb(234 179 8)" stroke-width="1.0"
        stroke="currentColor">
        <path fill-rule="evenodd"
            d="M19.952 1.651a.75.75 0 0 1 .298.599V16.303a3 3 0 0 1-2.176 2.884l-1.32.377a2.553 2.553 0 1 1-1.403-4.909l2.311-.66a1.5 1.5 0 0 0 1.088-1.442V6.994l-9 2.572v9.737a3 3 0 0 1-2.176 2.884l-1.32.377a2.553 2.553 0 1 1-1.402-4.909l2.31-.66a1.5 1.5 0 0 0 1.088-1.442V5.25a.75.75 0 0 1 .544-.721l10.5-3a.75.75 0 0 1 .658.122Z"
            clip-rule="evenodd" />
    </svg>

    <!-- PDF  -->
    <svg v-if='book_info.type == ".pdf"'
        class="opacity-70 w-7 h-7 m-0 absolute top-1 z-10 right-1 shadow-2xl hover:shadow-2xl"
        xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 16 16">
        <g fill="currentColor">
            <path
                d="M4.5 9.003a.5.5 0 0 0-.5.5v2a.5.5 0 0 0 1 0v-.166h.333a1.167 1.167 0 0 0 0-2.334H4.5zm.833 1.334H5v-.334h.333a.167.167 0 0 1 0 .334zm4.668-.835a.5.5 0 0 1 .5-.499h.998a.5.5 0 0 1 0 1h-.5v.335h.5a.5.5 0 1 1 0 1H11v.164a.5.5 0 0 1-1 .002L10 10.837L10 9.502zm-2.503-.499a.5.5 0 0 0-.5.5v2a.5.5 0 0 0 .5.5H8a1.5 1.5 0 0 0 0-3h-.502zm.5 2v-1H8a.5.5 0 0 1 0 1h-.002zM9 2.002H4.5a1.5 1.5 0 0 0-1.5 1.5v3.582A1.5 1.5 0 0 0 2 8.5v4.003a1.5 1.5 0 0 0 1.5 1.5h9a1.5 1.5 0 0 0 1.5-1.5V8.499a1.5 1.5 0 0 0-1-1.415V6h-2.5A1.5 1.5 0 0 1 9 4.5V2.002zM3.5 7.999h9a.5.5 0 0 1 .5.5v4.003a.5.5 0 0 1-.5.5h-9a.5.5 0 0 1-.5-.5V8.499a.5.5 0 0 1 .5-.5zm9.206-3H10.5a.5.5 0 0 1-.5-.5V2.298L12.706 5z"
                fill="rgb(234 179 8)" stroke-width="0.5" stroke="currentColor"></path>
        </g>
    </svg>
</template>

<script lang="ts">
import { defineComponent } from "vue";
export default defineComponent({
    name: "SvgBookIcon",
    props: ["book_info", "childBookNum"],
    components: {
    },
    computed: {
        childBookNum(): string {
            if (this.book_info.type === 'dir') {
                return ''
            }
            if (this.book_info.child_book_num > 0) {
                return "x" + this.book_info.child_book_num.toString();
            }
            return "";
        },
    },
    data() {
        return {
            test: "",
        };
    },
    methods: {
    },
});
</script>