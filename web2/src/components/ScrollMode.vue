<template>
	<div id="ScrollMode" v-if="this.book" class="manga">
		<!-- @click.self 只当在 event.target 是当前元素自身时才触发的处理函数 内部元素不会触发-->
		<Header v-if="this.showHeader">
			<h3 @click.self="drawerActivate('right')">
				<p v-if="book.IsFolder" :href="'raw/' + book.name">{{ book.name }}</p>
				<a v-if="!book.IsFolder" :href="'raw/' + book.name">{{ book.name }}【Download】</a>
			</h3>
		</Header>
		<div v-for="(page, key) in book.pages" :key="page.url" @click="getMouseXY($event)">
			<!-- v-lazy="page.url"  :src="page.url" -->
			<img
				v-lazy="page.url"
				v-bind:H="page.height"
				v-bind:W="page.width"
				v-bind:key="key"
				v-demo="{ singleWidth: this.getSinglePageWidth(), doubleWidth: this.getDoublePageWidth() ,isLandscape:this.isLandscapeMode()}"
			/>
			<p v-if="showPageNum">{{ key + 1 }}/{{ book.all_page_num }}</p>
		</div>
		<n-button @click="scrollToTop(90);">Back To Top</n-button>
	</div>
	<n-back-top :show="showBackTop" :right="20" :bottom="20" />
	<n-drawer v-model:show="drawerActive" :height="275" :width="251" :placement="drawerPlacement">
		<div>
			<n-drawer-content title="阅读设置" closable>
				<n-space>
					<n-radio-group v-model:value="selectedTemplate">
						<n-radio-button
							:checked="selectedTemplate === 'scroll'"
							@change="onChangeTemplate"
							value="scroll"
							name="basic-demo"
						>卷轴模式</n-radio-button>
						<n-radio-button
							:checked="selectedTemplate === 'single'"
							@change="onChangeTemplate"
							value="single"
							name="basic-demo"
						>单页模式</n-radio-button>
						<n-radio-button
							:checked="selectedTemplate === 'sketch'"
							@change="onChangeTemplate"
							value="sketch"
							name="basic-demo"
						>Sketch模式</n-radio-button>
					</n-radio-group>
				</n-space>

				<!-- 横屏模式 -->
				<n-space vertical v-if="this.isLandscapeMode()">
					<p>双开页宽度（横屏状态）</p>
					<n-slider
						v-model:value="doublePageWidth_LandscapeMode"
						:step="1"
						:max="100"
						:min="10"
						:format-tooltip="value => `${value}%`"
						:marks="marks"
					/>

					<p>单页宽度（横屏状态）</p>
					<n-slider
						v-model:value="singlePageWidth_LandscapeMode"
						:step="1"
						:max="100"
						:min="10"
						:format-tooltip="value => `${value}%`"
						:marks="marks"
					/>
				</n-space>

				<!-- 竖屏模式 -->
				<n-space vertical v-if="!this.isLandscapeMode()">
					双开页宽度（竖屏状态）:
					<n-slider
						v-model:value="doublePageWidth_PortraitMode"
						:step="1"
						:max="100"
						:min="10"
						:format-tooltip="value => `${value}%`"
						:marks="marks"
					/>
					<p>单页宽度（竖屏状态）</p>
					<n-slider
						v-model:value="singlePageWidth_PortraitMode"
						:step="1"
						:max="100"
						:min="10"
						:format-tooltip="value => `${value}%`"
						:marks="marks"
					/>
				</n-space>

				<!-- <p>同步滚动</p> -->
				<!-- <n-switch v-model:value="syncScroll">syncScroll</n-switch> -->
			</n-drawer-content>
		</div>
	</n-drawer>
</template>

<script>
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
import { NButton, NBackTop, NDrawer, NDrawerContent, NSpace, NSlider, NRadioButton, NRadioGroup } from 'naive-ui'
import Header from "@/components/Header.vue";
import { defineComponent, ref } from 'vue'
import { useCookies } from "vue3-cookies";
export default defineComponent({
	components: {
		Header,//标题，有点丑
		NButton,//按钮，来自：https://www.naiveui.com/zh-CN/os-theme/components/button
		NBackTop,//回到顶部按钮，来自：https://www.naiveui.com/zh-CN/os-theme/components/back-top
		NDrawer,//抽屉，可以从上下左右4个方向冒出. https://www.naiveui.com/zh-CN/os-theme/components/drawer
		NDrawerContent,//抽屉内容
		NSpace,//间距 https://www.naiveui.com/zh-CN/os-theme/components/space
		// NRadio,//单选  https://www.naiveui.com/zh-CN/os-theme/components/radio
		NRadioButton,//单选  用按钮显得更优雅一点
		NRadioGroup,
		NSlider,//滑动选择  Slider https://www.naiveui.com/zh-CN/os-theme/components/slider
		// NSwitch,//开关   https://www.naiveui.com/zh-CN/os-theme/components/switch
		// NLayout,//布局 https://www.naiveui.com/zh-CN/os-theme/components/layout
		// NLayoutSider,
		// NLayoutContent,
	},
	setup() {
		//此处不能使用this
		const { cookies } = useCookies();
		//抽屉相关
		const drawerActive = ref(false)
		const drawerPlacement = ref('right')
		const drawerActivate = (place) => {
			drawerActive.value = true
			drawerPlacement.value = place
		}
		//单选按钮绑定的数值
		// const checkedValueRef = ref(null)
		return {
			cookies,
			//抽屉激活状态
			drawerActive,
			//抽屉从哪个方向出现
			drawerPlacement,
			//激活抽屉的函数
			drawerActivate,
			//单选相关
			// checkedValue: checkedValueRef,
			// handleChange(e) {
			// 	checkedValueRef.value = e.target.value
			// },
			//滑动选择用
			marks: {
				30: '25%',
				50: '50%',
				75: '75%',
				95: '95%',
			}
		}
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
			//是否正在向下滚动
			scrollDown: false,
			//存储现在滚动的位置
			scrollTopSave: 0,
			//同步滚动，目前还没做
			syncScroll: false,
			//鼠标点击或触摸的位置
			clickX: 0,
			clickY: 0,
			landscapeMode: true,
			//横屏(Landscape)模式的漫画页宽度
			singlePageWidth_LandscapeMode: 50,
			doublePageWidth_LandscapeMode: 95,
			//竖屏(Portrait)模式的漫画页宽度
			singlePageWidth_PortraitMode: 100,
			doublePageWidth_PortraitMode: 100,
			selectedTemplate: "",
			clientWidth: 0,
			clientHeight: 0,
		};
	},
	//挂载前
	beforeMount() {
		//console.log('mounted in the composition api!')
		window.addEventListener("scroll", this.onScroll);
		window.addEventListener("resize", this.onResize);
		this.landscapeMode = this.isLandscapeMode();
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
	//卸载前
	beforeUnmount() {
		// 组件销毁时，销毁监听事件
		window.removeEventListener("scroll", this.onScroll);
		window.removeEventListener('resize', this.onResize)
	},
	methods: {
		onChangeTemplate() {
			// this.selectedTemplate = e.target.value
			if (this.selectedTemplate === "scroll") {
				this.cookies.set("nowTemplate", "scroll");
			}
			if (this.selectedTemplate === "single") {
				this.cookies.set("nowTemplate", "single");
			}
			if (this.selectedTemplate === "sketch") {
				//this.cookies.set("nowTemplate", "sketch");
			}
			location.reload(); //暂时无法动态刷新，研究好了再去掉
		},
		onResize() {
			this.clientWidth = document.documentElement.clientWidth
			this.clientHeight = document.documentElement.clientHeight
			// var aspectRatio = window.innerWidth / window.innerHeight
			var aspectRatio = this.clientWidth / this.clientHeight
			console.log("OnReSize,aspectRatio="+aspectRatio);
			// 为了半屏的时候更方便，阈值并不是正方形
			if (aspectRatio > (17 / 19)) {
				this.landscapeMode = true
			} else {
				this.landscapeMode = false
			}
		},
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
			//防手抖，小于一定数值状态就不变
			var step = Math.abs(this.scrollTopSave - scrollTop)
			// console.log("step:", step);
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
			//浏览器的视口，不包括工具栏和滚动条：
			// document.documentElement.clientHeight document.documentElement.ClientWidth不兼容手机？ 
			// var availHeight = document.documentElement.clientHeight
			// var availWidth = document.documentElement.clientWidth
			// console.log("clientHeigh=", document.documentElement.clientHeight, "ClientWidth=", document.documentElement.clientWidth);

			var availHeight = window.innerWidth
			var availWidth = window.innerHeight
			var MinX = availHeight * 0.40
			var MaxX = availHeight * 0.60
			var MinY = availWidth * 0.40
			var MaxY = availWidth * 0.60
			if ((this.clickX > MinX && this.clickX < MaxX) && (this.clickY > MinY && this.clickY < MaxY)) {
				//alert("点中了设置区域！")
				//console.log("点中了设置区域！");
				this.drawerActivate('right')
			}
			console.log("window.innerWidth=", window.innerWidth, "window.innerHeight=", window.innerHeight);
			// console.log("MinX=", MinX, "MaxX=", MaxX);
			// console.log("MinY=", MinY, "MaxY=", MaxY);
			// console.log("x=", e.x, "y=", e.y);
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
		//根据可视区域(viewport)长宽比，确认是横屏还是竖屏
		// aspect-ratio https://developer.mozilla.org/zh-CN/docs/Web/CSS/@media/aspect-ratio
		// window.innerWidth  不是响应式依赖，所以不能用计算属性
		isLandscapeMode() {
			//避免除数等于0，虽然正常情况下不会触发
			// if (window.innerHeight == 0) {
			// 	return false
			// }
			var aspectRatio = window.innerWidth / window.innerHeight
			// var aspectRatio = document.documentElement.clientWidth / document.documentElement.clientHeight
			
			var landscape =true
			// 为了半屏的时候更方便，阈值并不是正方形
			if (aspectRatio > (17 / 19)) {
				landscape= true
			} else {
				landscape= false
			}
			console.log("LandscapeMode="+landscape+" aspectRatio="+aspectRatio);
			return landscape
		},
		//单张漫画的宽度，暂时没有做Cookie保存
		getSinglePageWidth() {
			if (this.isLandscapeMode()) {
				return this.singlePageWidth_LandscapeMode + "vw"
			} else {
				return this.singlePageWidth_PortraitMode + "vw"
			}
		},
		//单张漫画的宽度，暂时没有做Cookie保存
		getDoublePageWidth() {
			if (this.isLandscapeMode()) {
				return this.doublePageWidth_LandscapeMode + "vw"
			} else {
				return this.doublePageWidth_PortraitMode + "vw"
			}
		},
	},
	computed: {

	}
});
</script>

<style></style>

<style scoped>
.manga {
	max-width: 100%;
}

/* https://developer.mozilla.org/zh-CN/docs/Web/CSS/object-fit */
.manga img {
	margin: auto;
	object-fit: scale-down;
	padding-top: 3px;
	padding-bottom: 3px;
	padding-right: 0px;
	padding-left: 0px;
	border-radius: 7px;
	box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

.LoadingImage {
	width: 80vw;
	max-width: 80vw;
}
.ErrorImage {
	width: 80vw;
	max-width: 80vw;
}

/* 竖屏(显示区域)CSS样式，IE无效 */
@media screen and (max-aspect-ratio: 17/19) {
  .SinglePageImage {
    width: 100%;
  }
  .DoublePageImage {
    width: 100%;
  }
}

/* 横屏（显示区域）时的CSS样式，IE无效 */
@media screen and (min-aspect-ratio: 17/19) {
  .SinglePageImage {
    width: 900px;
  }
  .DoublePageImage {
    width: 95%;
  }
}




</style>
