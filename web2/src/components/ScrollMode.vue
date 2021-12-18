<template>
	<div id="ScrollMode" v-if="this.book" class="manga">
		<Header v-if="this.showHeader">
			<h2>
				<p v-if="book.IsFolder" :href="'raw/' + book.name">{{ book.name }}</p>
				<a v-if="!book.IsFolder" :href="'raw/' + book.name">{{ book.name }}【Download】</a>
			</h2>
		</Header>
		<div v-for="(page, key) in book.pages" :key="page.url" @click="getMouseXY($event)">
			<!-- v-lazy="page.url"  :src="page.url" -->
			<img
				v-lazy="page.url"
				v-bind:H="page.height"
				v-bind:W="page.width"
				v-bind:key="key"
				v-bind:class="page.image_type"
			/>
			<p v-if="showPageNum">{{ key + 1 }}/{{ book.all_page_num }}</p>
		</div>
		<n-button @click="scrollToTop(90);">Back To Top</n-button>
	</div>
	<n-back-top :show="showBackTop" :right="20" :bottom="20" />
</template>

<script>
// //原生JS监听页面滚动
// window.onscroll = function () {
// 	var scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
// 	console.log("距离", scrollTop);
// }
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
import { NButton, NBackTop } from 'naive-ui'
import Header from "@/components/Header.vue";
import { onMounted } from 'vue'
export default {
	setup() {
		//不能使用this
		onMounted(() => {
		})
	},
	components: {
		Header,
		//按钮，来自：https://www.naiveui.com/zh-CN/os-theme/components/button
		NButton,
		//回到顶部按钮，来自：https://www.naiveui.com/zh-CN/os-theme/components/back-top
		NBackTop,
	},
	data() {
		return {
			//书籍数据，需要从远程拉取
			book: null,
			//是否显示Header
			showHeader: true,
			//是否显示页数
			showPageNum: false,
			//是否显示回到顶部按钮
			showBackTop: false,
			//回到顶部按钮的透明度,opacity="showBackTopOpacity" 无效？
			// showBackTopOpacity: 1.0,  
			//是否正在向下滚动
			scrollDown: false,
			//存储现在滚动的位置
			scrollTopSave: 0,
			clickX: 0,
			clickY: 0,
		};
	},
	beforeMount() {
		//console.log('mounted in the composition api!')
		window.addEventListener("scroll", this.onScroll);
		this.axios
			.get("/book.json")
			.then((response) => {
				if (response.status == 200) {
					this.book = response.data;
					// console.log(this.book);
				}
			})
			.catch((error) => alert(error));
	},
	methods: {
		//页面滚动的时候改变各种值
		onScroll() {
			var scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
			if (scrollTop > this.scrollTopSave) {
				this.scrollDown = true
				// console.log("下滚中，距离", scrollTop);
			} else {
				this.scrollDown = false
				// console.log("上滚中，距离", scrollTop);
			}
			var step = Math.abs(this.scrollTopSave - scrollTop)
			console.log("step:", step);
			this.scrollTopSave = scrollTop
			if (step > 5) {
				if (scrollTop > 400 && !this.scrollDown) {
					//页面高度大于400，且正在上滚的时候显示按钮
					this.showBackTop = true
				} else {
					//页面高度小于200执行操作
					this.showBackTop = false
				}
			}
		},
		//获取鼠标位置，然后做点什么
		getMouseXY(e) {
			this.clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离，单位为px）
			this.clickY = e.y //获取鼠标的Y坐标（鼠标与屏幕顶部的距离，单位为px）
			document.body.scrollWidth 
			//浏览器窗口的尺寸（浏览器的视口，不包括工具栏和滚动条）
			console.log("window.innerWidth=", window.innerWidth , "window.innerHeight=", window.innerHeight);
			var availHeight = window.innerWidth
			var availWidth = window.innerHeight
			var MinX = availHeight * 0.40
			var MaxX = availHeight * 0.60
			var MinY = availWidth * 0.40
			var MaxY = availWidth * 0.60
			console.log("MinX=", MinX, "MaxX=", MaxX);
			console.log("MinY=", MinY, "MaxY=", MaxY);

			if ((this.clickX > MinX && this.clickX < MaxX) && (this.clickY > MinY && this.clickY < MaxY)) {
				alert("点中了设置区域！")
				console.log("点中了设置区域！");
			}

			console.log("x=", e.x, "y=", e.y);
		},
		scrollToTop(scrollDuration) {
			var scrollStep = -window.scrollY / (scrollDuration / 15),
				scrollInterval = setInterval(function () {
					if (window.scrollY != 0) {
						window.scrollBy(0, scrollStep);
					}
					else clearInterval(scrollInterval);
				}, 15);
		},
	},
	computed: {
	}
};
</script>

<style></style>

<style scoped>
.manga {
	max-width: 100%;
}

/* https://developer.mozilla.org/zh-CN/docs/Web/CSS/object-fit */
.manga img {
	margin: auto;
	max-width: inherit;
	object-fit: scale-down;
	padding-top: 3px;
	padding-bottom: 3px;
	padding-right: 0px;
	padding-left: 0px;
	border-radius: 7px;
	box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

button {
	width: 100px;
	height: 40px;
	padding: 6px;
	background-color: #bbcbff;
	font-family: Avenir, Helvetica, Arial, sans-serif;
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;
	font-size: 13px;
}

.LoadingImage {
	width: 80vw;
	max-width: 80vw;
}
.ErrorImage {
	width: 80vw;
	max-width: 80vw;
}

/* 大致上是横向较宽（显示区域）时的CSS样式，IE无效 */
@media screen and (min-aspect-ratio: 19/19) {
	.SinglePageImage {
		width: 50vw;
	}
	.DoublePageImage {
		width: 90vw;
	}
}

/* 竖屏(显示区域)CSS样式，IE无效 */
@media screen and (max-aspect-ratio: 19/19) {
	.manga img {
		max-width: 100%;
	}
	.SinglePageImage {
		width: 100vw;
	}
	.DoublePage {
		width: 100vw;
	}
}

/* 高分横屏（显示区域）时的CSS样式，IE无效 */
/* min-width 输出设备中页面最小可视区域宽度 大于这个width时，其中的css起作用 超宽屏 */
@media screen and (min-aspect-ratio: 19/19) and (min-width: 1922px) {
	.SinglePageImage {
		width: 1000px;
	}
	.DoublePageImage {
		width: 1900px;
	}
}
</style>
