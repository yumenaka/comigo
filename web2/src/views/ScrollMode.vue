<template>
	<div id="ScrollMode" v-if="this.book" class="manga">
		<Header
			class="footer"
			v-if="this.showHeaderFlag"
			:bookIsFolder="book.IsFolder"
			:bookName="book.name"
		>
			<!-- 右边的设置图标，点击屏幕中央也可以打开 -->
			<n-icon size="40" @click="drawerActivate('right')">
				<settings-outline />
			</n-icon>
		</Header>

		<!-- 渲染漫画部分 -->
		<div class="main_manga"
			v-for="(page, key) in book.pages"
			:key="page.url"
			@click="onMouseClick($event)"
			@mousemove="onMouseMove"
			@mouseleave="onMouseLeave"
		>
			<!-- v-lazy="page.url"  :src="page.url" -->
			<img v-lazy="page.url" v-bind:H="page.height" v-bind:W="page.width" v-bind:key="key" />
			<div class="page_num" v-if="showPageNumFlag_ScrollMode">{{ key + 1 }}/{{ book.all_page_num }}</div>
		</div>

		<!-- 设置抽屉 -->
		<n-drawer
			v-model:show="drawerActive"
			:height="275"
			:width="251"
			:placement="drawerPlacement"
			@update:show="saveConfigToCookie"
		>
			<!-- 抽屉内容 -->
			<n-drawer-content closable>
				<!-- 抽屉：自定义头部 -->
				<template #header>
					<span>页面设置</span>
				</template>
				<!-- 切换页面模式 -->
				<n-space v-if="this.debugModeFlag">
					<n-radio-group v-model:value="selectedTemplate">
						<n-radio-button
							:checked="selectedTemplate === 'scroll'"
							@change="onChangeTemplate"
							value="scroll"
							name="basic-demo"
						>卷轴模式</n-radio-button>
						<n-radio-button
							:checked="selectedTemplate === 'flip'"
							@change="onChangeTemplate"
							value="flip"
							name="basic-demo"
						>翻页模式</n-radio-button>
					</n-radio-group>
				</n-space>
				<!-- 分割线 -->
				<n-divider />

				<n-space vertical>
					<!-- 单页-漫画宽度-使用百分比 -->
					<!-- 数字输入% -->
					<n-input-number
						v-if="this.imageWidth_usePercentFlag"
						size="small"
						:show-button="false"
						v-model:value="this.singlePageWidth_Percent"
						:max="100"
						:min="10"
					>
						<template #prefix>单页漫画宽度：</template>
						<template #suffix>%</template>
					</n-input-number>
					<!-- 滑动选择% -->
					<n-slider
						v-if="this.imageWidth_usePercentFlag"
						v-model:value="this.singlePageWidth_Percent"
						:step="1"
						:max="100"
						:min="10"
						:format-tooltip="value => `${value}%`"
					/>

					<!-- 开页-漫画宽度-使用百分比  -->

					<!-- 数字输入% -->
					<n-input-number
						v-if="this.imageWidth_usePercentFlag"
						size="small"
						:show-button="false"
						v-model:value="this.doublePageWidth_Percent"
						:max="100"
						:min="10"
					>
						<template #prefix>双开页漫画宽：</template>
						<template #suffix>%</template>
					</n-input-number>
					<!-- 滑动选择% -->
					<n-slider
						v-if="this.imageWidth_usePercentFlag"
						v-model:value="this.doublePageWidth_Percent"
						:step="1"
						:max="100"
						:min="10"
						:format-tooltip="value => `${value}%`"
					/>

					<!-- 单页-漫画宽度-使用固定值PX -->

					<!-- 数字输入PX -->
					<n-input-number
						v-if="!this.imageWidth_usePercentFlag"
						size="small"
						:show-button="false"
						v-model:value="this.singlePageWidth_PX"
						:max="this.imageMaxWidth"
						:min="50"
					>
						<template #prefix>单页漫画宽度：</template>
						<template #suffix>px</template>
					</n-input-number>
					<!-- 滑动选择PX -->
					<n-slider
						v-if="!this.imageWidth_usePercentFlag"
						v-model:value="this.singlePageWidth_PX"
						:step="10"
						:max="this.imageMaxWidth"
						:min="50"
						:format-tooltip="value => `${value}px`"
					/>

					<!-- 数字输入PX -->
					<n-input-number
						v-if="!this.imageWidth_usePercentFlag"
						size="small"
						:show-button="false"
						v-model:value="this.doublePageWidth_PX"
						:max="this.imageMaxWidth"
						:min="50"
					>
						<template #prefix>双开页漫画宽：</template>
						<template #suffix>px</template>
					</n-input-number>

					<!-- 滑动选择PX -->
					<n-slider
						v-if="!this.imageWidth_usePercentFlag"
						v-model:value="this.doublePageWidth_PX"
						:step="10"
						:max="this.imageMaxWidth"
						:min="50"
						:format-tooltip="value => `${value}px`"
					/>
				</n-space>

				<p></p>
				<!-- 开关：横屏状态下，宽度单位是百分比还是固定值 -->
				<n-space>
					<n-switch
						size="large"
						v-model:value="this.imageWidth_usePercentFlag"
						:rail-style="railStyle"
						@update:value="this.setImageWidthUsePercentFlag"
					>
						<template #checked>宽度:百分比%</template>
						<template #unchecked>宽度:固定值px</template>
					</n-switch>
				</n-space>

				<!-- 分割线 -->
				<n-divider />

				<!-- 开关：是否显示页头 -->
				<n-space>
					<n-switch size="large" v-model:value="this.showHeaderFlag" @update:value="setShowHeaderChange">
						<template #checked>显示页头</template>
						<template #unchecked>显示页头</template>
					</n-switch>
					<p></p>
				</n-space>

				<!-- 开关：是否显示当前页数 -->
				<n-space>
					<n-switch
						size="large"
						v-model:value="this.showPageNumFlag_ScrollMode"
						@update:value="setShowPageNumChange"
					>
						<template #checked>显示页数</template>
						<template #unchecked>显示页数</template>
					</n-switch>
				</n-space>

			</n-drawer-content>
		</n-drawer>
		<n-back-top :show="showBackTopFlag" type="info" color="#8a2be2" :right="20" :bottom="20" />
		<n-button @click="scrollToTop(90);" size="large" secondary strong>Back To Top</n-button>
	</div>
</template>

<script>
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
import { NButton, NBackTop, NDrawer, NDrawerContent, NSpace, NSlider, NRadioButton, NRadioGroup, NSwitch, NIcon, NInputNumber,NDivider  } from 'naive-ui'
import Header from "@/components/Header.vue";
import { defineComponent, ref } from 'vue'
import { useCookies } from "vue3-cookies";// https://github.com/KanHarI/vue3-cookies
import { SettingsOutline } from '@vicons/ionicons5'

export default defineComponent({
	name: "ScrollMode",
	props: ['book'],
	components: {
		Header,//页头，有点丑
		NButton,//按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/button
		NBackTop,//回到顶部按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/back-top
		NDrawer,//抽屉，可以从上下左右4个方向冒出. https://www.naiveui.com/zh-CN/os-theme/components/drawer
		NDrawerContent,//抽屉内容
		NSpace,//间距 https://www.naiveui.com/zh-CN/os-theme/components/space
		// NRadio,//单选  https://www.naiveui.com/zh-CN/os-theme/components/radio
		NRadioButton,//单选  用按钮显得更优雅一点
		NRadioGroup,
		NSlider,//滑动选择  Slider https://www.naiveui.com/zh-CN/os-theme/components/slider
		NSwitch,//开关   https://www.naiveui.com/zh-CN/os-theme/components/switch
		// NLayout,//布局 https://www.naiveui.com/zh-CN/os-theme/components/layout
		// NLayoutSider,
		// NLayoutContent,
		NIcon,//图标  https://www.naiveui.com/zh-CN/os-theme/components/icon
		// NPageHeader,//页头 https://www.naiveui.com/zh-CN/os-theme/components/page-header
		// NAvatar, //头像 https://www.naiveui.com/zh-CN/os-theme/components/avatar
		NInputNumber,//数字输入 https://www.naiveui.com/zh-CN/os-theme/components/input-number
		SettingsOutline,//图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）
		NDivider,//分割线  https://www.naiveui.com/zh-CN/os-theme/components/divider
	},
	setup() {
		//此处不能使用this
		const { cookies } = useCookies();
		//设置用的抽屉
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
			//开关用的颜色
			railStyle: ({ focused, checked }) => {
				const style = {}
				if (checked) {
					style.background = '#d03050'
					if (focused) {
						style.boxShadow = '0 0 0 2px #d0305040'
					}
				} else {
					style.background = '#2080f0'
					if (focused) {
						style.boxShadow = '0 0 0 2px #2080f040'
					}
				}
				return style
			},

			//滑动选择用建议值
			marks: {
				30: '25%',
				50: '50%',
				75: '75%',
				95: '95%',
			},
		}
	},
	data() {
		return {
			//开发模式 还没有做的功能与设置，设置Debug以后才能见到
			debugModeFlag: true,
			//书籍数据，需要从远程拉取
			//是否显示顶部页头
			showHeaderFlag: true,
			//是否显示页数
			showPageNumFlag_ScrollMode: false,
			//是否显示回到顶部按钮
			showBackTopFlag: false,
			//是否正在向下滚动
			scrollDownFlag: false,
			//存储现在滚动的位置
			scrollTopSave: 0,
			//同步滚动，目前还没做
			syncScrollFlag: false,
			//鼠标点击或触摸的位置
			clickX: 0,
			clickY: 0,
			//可见范围是否是横向
			isLandscapeMode: true,
			isPortraitMode: false,
			imageMaxWidth: 10,
			//屏幕宽横比，inLandscapeMode的判断依据
			aspectRatio: 1.2,

			//状态驱动的动态 CSS!!!!!
			// https://v3.cn.vuejs.org/api/sfc-style.html#%E7%8A%B6%E6%80%81%E9%A9%B1%E5%8A%A8%E7%9A%84%E5%8A%A8%E6%80%81-css
			//图片宽度的单位，是否使用百分比
			imageWidth_usePercentFlag: true,

			//横屏(Landscape)状态的漫画页宽度，百分比
			singlePageWidth_Percent: 50,
			doublePageWidth_Percent: 95,

			//横屏(Landscape)状态的漫画页宽度，PX
			singlePageWidth_PX: 720,
			doublePageWidth_PX: 1080,

			//选择了哪个阅读模板
			selectedTemplate: "scroll",
			//可见范围宽高的具体值
			clientWidth: 0,
			clientHeight: 0,
		};
	},
	//Vue3生命周期:  https://v3.cn.vuejs.org/api/options-lifecycle-hooks.html#beforecreate
	// created : 在绑定元素的属性或事件监听器被应用之前调用。
	// beforeMount : 指令第一次绑定到元素并且在挂载父组件之前调用。。
	// mounted : 在绑定元素的父组件被挂载后调用。。
	// beforeUpdate: 在更新包含组件的 VNode 之前调用。。
	// updated: 在包含组件的 VNode 及其子组件的 VNode 更新后调用。
	// beforeUnmount: 当指令与在绑定元素父组件卸载之前时，只调用一次。
	// unmounted: 当指令与元素解除绑定且父组件已卸载时，只调用一次。

	created() {
		window.addEventListener("scroll", this.onScroll);
		//文档视图调整大小时会触发 resize 事件。 https://developer.mozilla.org/zh-CN/docs/Web/API/Window/resize_event
		window.addEventListener("resize", this.onResize);
		this.imageMaxWidth = window.innerWidth;
		//根据cookie初始化默认值,或初始化cookie值,cookie读取出来的都是字符串，不要直接用
		//是否显示顶部页头
		if (this.cookies.get("showHeaderFlag") === "true") {
			this.showHeaderFlag = true;
		} else if (this.cookies.get("showHeaderFlag") === "false") {
			this.showHeaderFlag = false;
		}
		//console.log("读取cookie并初始化: showHeaderFlag=" + this.showHeaderFlag);

		//是否显示页数
		if (this.cookies.get("showPageNumFlag_ScrollMode") === "true") {
			this.showPageNumFlag_ScrollMode = true;
		} else if (this.cookies.get("showPageNumFlag_ScrollMode") === "false") {
			this.showPageNumFlag_ScrollMode = false;
		}
		//console.log("读取cookie并初始化: showPageNumFlag_ScrollMode=" + this.showPageNumFlag_ScrollMode);

		//宽度是否使用百分比
		if (this.cookies.get("imageWidth_usePercentFlag") === "true") {
			this.imageWidth_usePercentFlag = true;
		} else if (this.cookies.get("imageWidth_usePercentFlag") === "false") {
			this.imageWidth_usePercentFlag = false;
		}

		//javascript 数字类型转换：https://www.runoob.com/js/js-type-conversion.html
		// NaN不能通过相等操作符（== 和 ===）来判断

		//漫画页宽度，Percent
		if (this.cookies.get("singlePageWidth_Percent") != null) {
			let saveNum = Number(this.cookies.get("singlePageWidth_Percent"));
			if (!isNaN(saveNum)) {
				this.singlePageWidth_Percent = saveNum;
			}
		}

		if (this.cookies.get("doublePageWidth_Percent") != null) {
			let saveNum = Number(this.cookies.get("doublePageWidth_Percent"));
			if (!isNaN(saveNum)) {
				this.doublePageWidth_Percent = saveNum;
			}
		}

		//漫画页宽度，PX
		if (this.cookies.get("singlePageWidth_PX") != null) {
			let saveNum = Number(this.cookies.get("singlePageWidth_PX"));
			if (!isNaN(saveNum)) {
				this.singlePageWidth_PX = saveNum;
			}
		}
		if (this.cookies.get("doublePageWidth_PX") != null) {
			let saveNum = Number(this.cookies.get("doublePageWidth_PX"));
			if (!isNaN(saveNum)) {
				this.doublePageWidth_PX = saveNum;
			}
		}

	},

	// //挂载前
	beforeMount() {

	},
	onMounted() {
		//console.log('mounted in the composition api!')

		this.isLandscapeMode = this.inLandscapeModeCheck();
		this.isPortraitMode = !this.inLandscapeModeCheck();
		// https://v3.cn.vuejs.org/api/options-lifecycle-hooks.html#beforemount
		this.$nextTick(function () {
			// 仅在整个视图都被渲染之后才会运行的代码
		})
	},
	//卸载前
	beforeUnmount() {
		// 组件销毁前，销毁监听事件
		window.removeEventListener("scroll", this.onScroll);
		window.removeEventListener('resize', this.onResize)
	},
	methods: {
		//如果在一个组件上使用了 v-model:xxx，应该使用 @update:xxx  https://www.naiveui.com/zh-CN/os-theme/docs/common-issues
		saveConfigToCookie(show) {
			console.log("show:" + show)
			// 组件销毁前，储存cookie
			this.cookies.set("showHeaderFlag", this.showHeaderFlag);
			this.cookies.set("showPageNumFlag_ScrollMode", this.showPageNumFlag_ScrollMode);
			this.cookies.set("imageWidth_usePercentFlag", this.imageWidth_usePercentFlag);
			this.cookies.set("singlePageWidth_Percent", this.singlePageWidth_Percent);
			this.cookies.set("doublePageWidth_Percent", this.doublePageWidth_Percent);
			this.cookies.set("singlePageWidth_PX", this.singlePageWidth_PX);
			this.cookies.set("doublePageWidth_PX", this.doublePageWidth_PX);
		},
		setShowHeaderChange(value) {
			console.log("value:" + value);
			this.showHeaderFlag = value;
			this.cookies.set("showHeaderFlag", value);
			console.log("cookie设置完毕: showHeaderFlag=" + this.cookies.get("showHeaderFlag"));
		},
		setShowPageNumChange(value) {
			console.log("value:" + value);
			this.showPageNumFlag_ScrollMode = value;
			this.cookies.set("showPageNumFlag_ScrollMode", value);
			console.log("cookie设置完毕: showPageNumFlag_ScrollMode=" + this.cookies.get("showPageNumFlag_ScrollMode"));
		},

		setImageWidthUsePercentFlag(value) {
			console.log("value:" + value);
			this.imageWidth_usePercentFlag = value;
			this.cookies.set("imageWidth_usePercentFlag", value);
			console.log("cookie设置完毕: imageWidth_usePercentFlag=" + this.imageWidth_usePercentFlag);
		},

		//切换模板的函数，需要配合vue-router
		onChangeTemplate() {
			// this.selectedTemplate = e.target.value
			if (this.selectedTemplate === "scroll") {
				this.cookies.set("nowTemplate", "scroll");
			}
			if (this.selectedTemplate === "flip") {
				this.cookies.set("nowTemplate", "flip");
			}
			location.reload(); //暂时无法动态刷新，研究vue-router去掉
		},
		//可见区域变化的时候改变页面状态
		onResize() {
			this.imageMaxWidth = window.innerWidth
			// document.querySelectorAll(".name");
			this.clientWidth = document.documentElement.clientWidth
			this.clientHeight = document.documentElement.clientHeight
			// var aspectRatio = window.innerWidth / window.innerHeight
			this.aspectRatio = this.clientWidth / this.clientHeight
			//console.log("OnReSize,aspectRatio=" + this.aspectRatio);
			// 为了调试的时候方便，阈值是正方形
			if (this.aspectRatio > (19 / 19)) {
				this.isLandscapeMode = true
				this.isPortraitMode = false
			} else {
				this.isLandscapeMode = false
				this.isPortraitMode = true
			}
		},
		//页面滚动的时候改变各种值
		onScroll() {
			var scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
			if (scrollTop > this.scrollTopSave) {
				this.scrollDownFlag = true
				// console.log("下滚中，距离", scrollTop);
			} else {
				this.scrollDownFlag = false
				// console.log("上滚中，距离", scrollTop);
			}
			//防手抖，小于一定数值状态就不变
			var step = Math.abs(this.scrollTopSave - scrollTop)
			// console.log("step:", step);
			this.scrollTopSave = scrollTop
			if (step > 5) {
				if (scrollTop > 400 && !this.scrollDownFlag) {
					//页面高度大于400，且正在上滚的时候显示按钮
					this.showBackTopFlag = true
				} else {
					//页面高度小于200执行操作
					this.showBackTopFlag = false
				}
			}
		},
		//获取鼠标位置，决定是否打开设置面板
		onMouseClick(e) {
			this.clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离，单位为px）
			this.clickY = e.y //获取鼠标的Y坐标（鼠标与屏幕顶部的距离，单位为px）
			//浏览器的视口，不包括工具栏和滚动条:
			var availHeight = window.innerWidth
			var availWidth = window.innerHeight
			var MinX = availHeight * 0.37
			var MaxX = availHeight * 0.65
			var MinY = availWidth * 0.37
			var MaxY = availWidth * 0.65
			if ((this.clickX > MinX && this.clickX < MaxX) && (this.clickY > MinY && this.clickY < MaxY)) {
				//alert("点中了设置区域！")
				//console.log("点中了设置区域！");
				this.drawerActivate('right')
			}
			// console.log("window.innerWidth=", window.innerWidth, "window.innerHeight=", window.innerHeight);
			// console.log("MinX=", MinX, "MaxX=", MaxX);
			// console.log("MinY=", MinY, "MaxY=", MaxY);
			// console.log("x=", e.x, "y=", e.y);
		},

		onMouseMove(e) {
			this.clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离，单位为px）
			this.clickY = e.y //获取鼠标的Y坐标（鼠标与屏幕顶部的距离，单位为px）
			//浏览器的视口，不包括工具栏和滚动条:
			var availHeight = window.innerWidth
			var availWidth = window.innerHeight
			var MinX = availHeight * 0.37
			var MaxX = availHeight * 0.65
			var MinY = availWidth * 0.37
			var MaxY = availWidth * 0.65
			if ((this.clickX > MinX && this.clickX < MaxX) && (this.clickY > MinY && this.clickY < MaxY)) {
				//console.log("在设置区域！");
				e.currentTarget.style.cursor = 'url(/images/SettingsOutline.png), pointer';
			} else {
				e.currentTarget.style.cursor = '';
			}
		},
		onMouseLeave(e) {
			//离开区域的时候，清空鼠标样式
			e.currentTarget.style.cursor = '';
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
		inLandscapeModeCheck() {
			//避免除数等于0，虽然正常情况下不会触发
			// if (window.innerHeight == 0) {
			// 	return false
			// }
			// var aspectRatio = document.documentElement.clientWidth / document.documentElement.clientHeight
			this.aspectRatio = window.innerWidth / window.innerHeight
			// console.log("aspectRatio=" + this.aspectRatio);
			// 为了半屏的时候更方便，阈值是正方形
			if (this.aspectRatio > (19 / 19)) {
				return true
			} else {
				return false
			}
		},
	},

	computed: {
		sPWL() {
			if (this.imageWidth_usePercentFlag) {
				return this.singlePageWidth_Percent + '%';
			} else {
				return this.singlePageWidth_PX + 'px';
			}
		},
		dPWL() {
			if (this.imageWidth_usePercentFlag) {
				return this.doublePageWidth_Percent + '%';
			} else {
				return this.doublePageWidth_PX + 'px';
			}

		},
		sPWP() {
			if (this.imageWidth_usePercentFlag) {
				return this.singlePageWidth_Percent + '%';
			} else {
				return this.singlePageWidth_PX + 'px';
			}
		},
		dPWP() {
			if (this.imageWidth_usePercentFlag) {
				return this.doublePageWidth_Percent + '%';
			} else {
				return this.doublePageWidth_PX + 'px';
			}
		},
	}
});
</script>

<style></style>

<style scoped>
.manga {
	max-width: 100%;
}

.header {
	padding: 0px;
	width: 100%;
	height: 7vh;
}

/* https://developer.mozilla.org/zh-CN/docs/Web/CSS/object-fit */
.manga img {
	margin: auto;
	/* object-fit: scale-down; */
	padding-top: 3px;
	padding-bottom: 3px;
	padding-right: 0px;
	padding-left: 0px;
	border-radius: 7px;
	box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

.LoadingImage {
	width: 90vw;
	max-width: 90vw;
}
.ErrorImage {
	width: 90vw;
	max-width: 90vw;
}

/* 横屏（显示区域）时的CSS样式，IE无效 */
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

/* 竖屏(显示区域)CSS样式，IE无效 */
@media screen and (max-aspect-ratio: 19/19) {
	.SinglePageImage {
		/* width: 100%; */
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
