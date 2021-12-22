<template>
	<div class="body" v-if="this.book">
		<Header
			class="header"
			v-if="this.showHeaderFlag_FlipMode"
			:bookIsFolder="book.IsFolder"
			:bookName="book.name"
		>
			<!-- 右边的设置图标，点击屏幕中央也可以打开 -->
			<n-icon size="40" @click="drawerActivate('right')">
				<settings-outline />
			</n-icon>
		</Header>
		<div
			class="main_manga"
			id="MangaMain"
			@click.self="onMouseClick"
			@mousemove.self="onMouseMove"
			@mouseleave.self="onMouseLeave"
		>
			<img
				@click.self="onMouseClick"
				@mousemove.self="onMouseMove"
				@mouseleave.self="onMouseLeave"
				class="manga_img"
				v-if="now_page <= this.book.all_page_num && now_page >= 1"
				lazy-src="/resources/favicon.ico"
				v-bind:src="this.book.pages[now_page - 1].url"
			/>
			<img />
			<span v-if="this.showPageNumFlag_FlipMode">{{ now_page }}/{{ book.all_page_num }}</span>
		</div>
		<div class="footer" v-if="this.showFooterFlag_FlipMode">
			<!-- 右手模式用 ，底部滑动条 -->

			<div v-if="this.useRightHandFlag">
				<span>第{{ this.now_page }}页</span>
				<n-slider
					v-model:value="now_page"
					:max="this.book.all_page_num"
					:min="1"
					:step="1"
					:format-tooltip="value => `第${value}页`"
				/>
				<span>第{{ this.book.all_page_num }}页</span>
			</div>

			<!-- 左手模式用 底部滑动条，设置reverse翻转计数方向 -->

			<div v-if="!this.useRightHandFlag">
				<span>第{{ this.book.all_page_num }}页</span>
				<n-slider
					reverse
					v-model:value="now_page"
					:max="this.book.all_page_num"
					:min="1"
					:step="1"
					:format-tooltip="value => `第${value}页`"
				/>
				<span>第{{ this.now_page }}页</span>
			</div>
		</div>
	</div>

	<n-drawer
		v-model:show="drawerActive"
		@update:show="saveConfigToCookie"
		:height="275"
		:width="251"
		:placement="drawerPlacement"
	>
		<n-drawer-content title="页面设置" closable>
			<!-- 切换页面模式 -->
			<n-space>
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
					>单页模式</n-radio-button>
				</n-radio-group>
			</n-space>

			<!-- 分割线 -->
			<n-divider />

			<span>设定背景色：</span>
			<n-color-picker v-model:value="model.color" :modes="['rgb']" :show-alpha="false" />

			<!-- 分割线 -->
			<n-divider />

			<!-- 开关：页头与书名 -->
			<n-space>
				<n-switch
					size="large"
					v-model:value="this.showHeaderFlag_FlipMode"
					@update:value="setShowHeaderChange"
				>
					<template #checked>页头与书名</template>
					<template #unchecked>页头与书名</template>
				</n-switch>
			</n-space>

			<!-- 开关：显示阅读进度条） -->
			<n-space>
				<n-switch
					size="large"
					v-model:value="this.showFooterFlag_FlipMode"
					@update:value="setShowFooterFlagChange"
				>
					<template #checked>阅读进度条</template>
					<template #unchecked>阅读进度条</template>
				</n-switch>
			</n-space>

			<!-- 开关：显示当前页数 -->
			<n-space>
				<n-switch
					size="large"
					v-model:value="this.showPageNumFlag_FlipMode"
					@update:value="setShowPageNumChange"
				>
					<template #checked>显示页数</template>
					<template #unchecked>显示页数</template>
				</n-switch>
			</n-space>

			<!-- 保存阅读进度 -->
			<n-space>
				<n-switch
					size="large"
					v-model:value="this.savePageNumFlag"
					@update:value="this.setSavePageNumFlag"
				>
					<template #checked>保存进度</template>
					<template #unchecked>保存进度</template>
				</n-switch>
			</n-space>

			<!-- 开关：翻页模式，默认右手-->
			<n-space>
				<n-switch
					size="large"
					v-model:value="this.useRightHandFlag"
					:rail-style="railStyle"
					@update:value="this.setUseLeftHandFlag"
				>
					<template #checked>右手翻页</template>
					<template #unchecked>左手翻页</template>
				</n-switch>
			</n-space>
			<!-- 分割线 -->
			<n-divider />

			<!-- 随机切换背景色 -->
			<n-space>
				<n-switch
					size="large"
					v-model:value="this.debugModeFlag"
					@update:value="this.setDebugModeFlagFlag"
				>
					<template #checked>Debug</template>
					<template #unchecked>Debug</template>
				</n-switch>
			</n-space>

			<!-- 抽屉：自定义底部 -->
			<template #footer>
				<n-button @click="changeTemplateToSketch">倒计时素描-开始</n-button>
				<n-avatar size="small" src="/favicon.ico" />
			</template>
		</n-drawer-content>
	</n-drawer>
</template>

<script>
import { useCookies } from "vue3-cookies";
import Header from "@/components/Header.vue";
import { defineComponent, ref, reactive } from 'vue'
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
import { NDrawer, NDrawerContent, NSpace, NSlider, NRadioButton, NRadioGroup, NSwitch, NIcon, NColorPicker, NAvatar, NButton, NDivider, } from 'naive-ui'
import { SettingsOutline } from '@vicons/ionicons5'
export default defineComponent({
	name: "FlipMode",
	props: ['book'],
	components: {
		Header,
		NDrawer,//抽屉，可以从上下左右4个方向冒出. https://www.naiveui.com/zh-CN/os-theme/components/drawer
		NDrawerContent,//抽屉内容
		NButton,//按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/button
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
		NAvatar, //头像 https://www.naiveui.com/zh-CN/os-theme/components/avatar
		SettingsOutline,//图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）
		NColorPicker, //颜色选择器 Color Picker https://www.naiveui.com/zh-CN/os-theme/components/color-picker
		NDivider,//分割线  https://www.naiveui.com/zh-CN/os-theme/components/divider
	},
	setup() {
		const { cookies } = useCookies();
		//设置用的抽屉
		const drawerActive = ref(false)
		const drawerPlacement = ref('right')
		const drawerActivate = (place) => {
			drawerActive.value = true
			drawerPlacement.value = place
		}
		//颜色选择器
		const model = reactive({
			color: '#252525'
		})
		return {
			model,
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
		};
	},
	data() {
		return {
			//开发模式 没做的功能与设置，设置Debug以后才能见到
			debugModeFlag: true,
			//是否显示页头
			showHeaderFlag_FlipMode: true,
			//是否显示页脚
			showFooterFlag_FlipMode: true,
			//是否显示页数
			showPageNumFlag_FlipMode: false,
			//是否是右手翻页
			useRightHandFlag: true,
			//是否拼合双叶
			autoDoublepage_FlipMode: true,
			//是否保存当前页数
			savePageNumFlag: true,
			//当前页数
			now_page: 1,
			selectedTemplate: "flip",
			sketchModeFlag: false,
			sketchSecondCount: 0,
			sketchFlipSecond: 30,
			interval: null,
		};
	},
	//在选项API中使用 Vue 生命周期钩子

	created() {
		//从cookies初始化默认值
		if (this.cookies.get("debugModeFlag") === "true") {
			this.debugModeFlag = true;
		} else if (this.cookies.get("debugModeFlag") === "false") {
			this.debugModeFlag = false;
		}

		if (this.cookies.get("showHeaderFlag_FlipMode") === "true") {
			this.showHeaderFlag_FlipMode = true;
		} else if (this.cookies.get("showHeaderFlag_FlipMode") === "false") {
			this.showHeaderFlag_FlipMode = false;
		}

		if (this.cookies.get("showFooterFlag_FlipMode_FlipMode") === "true") {
			this.showFooterFlag_FlipMode = true;
		} else if (this.cookies.get("showFooterFlag_FlipMode_FlipMode") === "false") {
			this.showFooterFlag_FlipMode = false;
		}
		//是否显示页数
		if (this.cookies.get("showPageNumFlag_FlipMode") === "true") {
			this.showPageNumFlag_FlipMode = true;
		} else if (this.cookies.get("showPageNumFlag_FlipMode") === "false") {
			this.showPageNumFlag_FlipMode = false;
		}
		if (this.cookies.get("useRightHandFlag") === "true") {
			this.useRightHandFlag = true;
		} else if (this.cookies.get("useRightHandFlag") === "false") {
			this.useRightHandFlag = false;
		}

		if (this.cookies.get("autoDoublepage_FlipMode") === "true") {
			this.autoDoublepage_FlipMode = true;
		} else if (this.cookies.get("autoDoublepage_FlipMode") === "false") {
			this.autoDoublepage_FlipMode = false;
		}

		if (this.cookies.get("savePageNumFlag") === "true") {
			this.savePageNumFlag = true;
		} else if (this.cookies.get("savePageNumFlag") === "false") {
			this.savePageNumFlag = false;
		}

		//当前颜色
		if (this.cookies.get("FlipModeDefaultColor") != null) {
			this.model.color = this.cookies.get("FlipModeDefaultColor");
		}
	},
	//挂载前
	beforeMount() {
		// 注册监听
		window.addEventListener("keyup", this.handleKeyup);
	},
	//卸载前
	beforeUnmount() {
		// 销毁监听
		window.removeEventListener("keyup", this.handleKeyup);
	},
	//挂载后
	mounted() {
		// setInterval(this.changeRandomColor, 1000);

	},

	updated() {
		//当前页数
		if (this.cookies.get("now_page" + this.book.name) != null) {
			let saveNum = Number(this.cookies.get("now_page" + this.book.name));
			if (!isNaN(saveNum)) {
				this.now_page = saveNum;
			}
		}
	},

	methods: {

		changeTemplateToSketch(e) {
			this.sketchModeFlag = !this.sketchModeFlag;
			if (this.sketchModeFlag) {
				//开始倒计时素描
				e.currentTarget.value = "倒计时素描-停止"
				this.cookies.set("nowTemplate", "sketch");
				this.selectedTemplate = "sketch";
				this.sketchModeFlag = true;
				this.interval = setInterval(this.sketchCount, 1000);
			} else {
				//停止倒计时素描
				e.currentTarget.value = "倒计时素描-开始"
				this.sketchModeFlag = false;
				this.cookies.set("nowTemplate", "flip");
				this.selectedTemplate = "flip";
				this.sketchSecondCount = 0;
			}
		},
		sketchCount() {
			if (!this.sketchModeFlag) {
				clearInterval(this.interval); // 清除定时器
				this.sketchSecondCount = 0;
			}
			this.sketchSecondCount = this.sketchSecondCount + 1;
			let nowSeconnd =this.sketchSecondCount % this.sketchFlipSecond
			console.log("sketchSecondCount="+this.sketchSecondCount+" nowSeconnd:"+nowSeconnd)
			if (nowSeconnd == 0) {
				if (this.now_page < this.book.all_page_num) {
					this.flipPage(1);
				} else {
					this.toPage(1);
				}
			}

		},
		// 关闭抽屉时，保存设置到cookies
		saveConfigToCookie() {
			this.cookies.set("debugModeFlag", this.debugModeFlag);
			this.cookies.set("showHeaderFlag_FlipMode", this.showHeaderFlag_FlipMode);
			this.cookies.set("showFooterFlag_FlipMode", this.showFooterFlag_FlipMode);
			this.cookies.set("showPageNumFlag_FlipMode", this.showPageNumFlag_FlipMode);
			this.cookies.set("useRightHandFlag", this.useRightHandFlag);
			this.cookies.set("autoDoublepage_FlipMode", this.autoDoublepage_FlipMode);
			this.cookies.set("savePageNumFlag", this.savePageNumFlag);
			this.cookies.set("now_page" + this.book.name, this.now_page);
			this.cookies.set("FlipModeDefaultColor", this.model.color);
		},
		changeRandomColor() {
			let R = Math.ceil(Math.random() * 155) + 100;
			let G = Math.ceil(Math.random() * 155) + 100;
			let B = Math.ceil(Math.random() * 100) + 100;
			//rgb(185,175,145)
			let RGB = 'rgb(' + R + "," + G + "," + B + ")";
			// console.log(RGB);
			this.model.color = RGB;
		},
		//HTML DOM 事件 https://www.runoob.com/jsref/dom-obj-event.html
		// 进入绑定该事件的元素和其子元素均会触发该事件，所以有一个重复触发，冒泡过程。其对应的离开事件 mouseout
		onMouseOver() {
		},
		// 只有进入绑定该事件的元素才会触发事件，也就是不会冒泡。其对应的离开事件mouseleave
		onMouseEnter() {
			// this.randomColor = 'background-color: rgb(235,235,235)';
		},
		onMouseLeave(e) {
			//离开区域的时候，清空鼠标样式
			e.currentTarget.style.cursor = '';
		},
		//事件修饰符: https://v3.cn.vuejs.org/guide/events.html#%E4%BA%8B%E4%BB%B6%E4%BF%AE%E9%A5%B0%E7%AC%A6
		onMouseMove(e) {
			// offsetX/Y获取到是触发点相对被触发dom的左上角距离
			let offsetX = e.offsetX;
			let offsetY = e.offsetY;
			//根据ID获取元素
			// let mangaDiv =document.getElementById("MangaMain")
			//不用自己获取元素
			let offsetWidth = e.currentTarget.offsetWidth;
			let offsetHeight = e.currentTarget.offsetHeight;
			// console.log("e.offsetX=" + offsetX, "e.offsetY=" + offsetY);
			// console.log("e.offsetWidth=" + offsetWidth, "e.offsetHeight=" + offsetHeight);
			var MinX = offsetWidth * 0.40
			var MaxX = offsetWidth * 0.60
			var MinY = offsetHeight * 0.40
			var MaxY = offsetHeight * 0.60
			if ((offsetX > MinX && offsetX < MaxX) && (offsetY > MinY && offsetY < MaxY)) {
				//设置区域;
				e.currentTarget.style.cursor = 'url(/images/SettingsOutline.png), pointer';
			} else {
				if (offsetX < (offsetWidth * 0.50)) {
					//左边的翻页
					e.currentTarget.style.cursor = 'url(/images/ArrowLeft.png), pointer';

				} else {
					//右边的翻页
					e.currentTarget.style.cursor = 'url(/images/ArrowRight.png), pointer';
				}
			}
		},

		//根据鼠标点击事件的位置，决定是左右翻页还是打开设置
		onMouseClick(e) {
			let offsetX = e.offsetX;
			let offsetY = e.offsetY;
			let offsetWidth = e.currentTarget.offsetWidth;
			let offsetHeight = e.currentTarget.offsetHeight;
			var MinX = offsetWidth * 0.40
			var MaxX = offsetWidth * 0.60
			var MinY = offsetHeight * 0.40
			var MaxY = offsetHeight * 0.60
			// console.log("鼠标点击：e.offsetX=" + offsetX, "e.offsetY=" + offsetY);
			if ((offsetX > MinX && offsetX < MaxX) && (offsetY > MinY && offsetY < MaxY)) {
				//点中了设置区域
				this.drawerActivate('right')
			} else {
				//随机一下背景色，只是为了好玩
				if (this.debugModeFlag) {
					this.changeRandomColor();
				}
				//决定如何翻页
				if (offsetX <= (offsetWidth / 2.0)) {
					//左边的翻页
					if (this.useRightHandFlag) {
						this.flipPage(-1);
					} else {
						this.flipPage(1);
					}
				} else {
					//右边的翻页
					if (this.useRightHandFlag) {
						this.flipPage(1);
					} else {
						this.flipPage(-1);
					}
				}
			}
		},
		setShowHeaderChange(value) {
			console.log("value:" + value);
			this.showHeaderFlag_FlipMode = value;
			this.cookies.set("showHeaderFlag_FlipMode", value);
			console.log("cookie设置完毕: showHeaderFlag_FlipMode=" + this.cookies.get("showHeaderFlag_FlipMode"));
		},
		setShowFooterFlagChange(value) {
			console.log("value:" + value);
			this.showFooterFlag_FlipMode = value;
			this.cookies.set("showFooterFlag_FlipMode", value);
			console.log("cookie设置完毕: showFooterFlag_FlipMode=" + this.cookies.get("showFooterFlag_FlipMode"));
		},

		setShowPageNumChange(value) {
			console.log("value:" + value);
			this.showPageNumFlag_FlipMode = value;
			this.cookies.set("showPageNumFlag_FlipMode", value);
			console.log("cookie设置完毕: showPageNumFlag_FlipMode=" + this.cookies.get("showPageNumFlag_FlipMode"));
		},

		setUseLeftHandFlag(value) {
			console.log("value:" + value);
			this.useRightHandFlag = value;
			this.cookies.set("useRightHandFlag", value);
			console.log("cookie设置完毕: useRightHandFlag=" + this.cookies.get("useRightHandFlag"));
		},

		setSavePageNumFlag(value) {
			console.log("value:" + value);
			this.savePageNumFlag = value;
			this.cookies.set("savePageNumFlag", value);
			console.log("cookie设置完毕: savePageNumFlag=" + this.cookies.get("savePageNumFlag"));
		},

		setDebugModeFlagFlag(value) {
			console.log("value:" + value);
			this.debugModeFlag = value;
			this.cookies.set("debugModeFlag", value);
			console.log("cookie设置完毕: debugModeFlag=" + this.cookies.get("debugModeFlag"));
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
			if (this.selectedTemplate === "sketch") {
				this.cookies.set("nowTemplate", "sketch");
			}
			location.reload(); //暂时无法动态刷新，研究vue-router去掉
		},

		flipPage: function (num) {
			if (
				this.now_page + num <= this.book.all_page_num &&
				this.now_page + num >= 1
			) {
				this.now_page = this.now_page + num;
			}
			if (this.savePageNumFlag) {
				this.cookies.set("now_page" + this.book.name, this.now_page);
			}
		},
		toPage: function (num) {
			if (num <= this.book.all_page_num && num >= 1) {
				this.now_page = num;
			}
			if (this.savePageNumFlag) {
				this.cookies.set("now_page" + this.book.name, this.now_page);
			}
			// console.log(num);
		},
		// 键盘事件
		handleKeyup(event) {
			const e = event || window.event || arguments.callee.caller.arguments[0];
			if (!e) return;
			//https://developer.mozilla.org/zh-CN/docs/Web/API/KeyboardEvent/keyCode
			switch (e.key) {
				case "ArrowUp":
				case "PageUp":
				case "ArrowLeft":
					this.flipPage(-1); //上一页
					break;
				case "Space":
				case "ArrowDown":
				case "PageDown":
				case "ArrowRight":
					this.flipPage(1); //下一页
					break;
				case "Home":
					this.toPage(1); //跳转到第一页
					break;
				case "End":
					this.toPage(this.book.all_page_num); //跳转到最后一页
					break;
				case "Ctrl":
					// Ctrl key pressed //组合键？
					//openOverlay();
					break;
			}
			// console.log(e.keyCode);
			// console.log(e.key);
		},
		//  滑轮事件
		handleScroll() {
			var e = document.body.scrollTop || document.documentElement.scrollTop;
			if (!e) return;
			// console.log(e);
		},
	},

	computed: {
		mainHeight() {
			let Height = "100vh";
			if (this.showFooterFlag_FlipMode || this.showHeaderFlag_FlipMode) {
				if (this.showPageNumFlag_FlipMode) {
					Height = "95vh"
				} else {
					Height = "95vh"
				}

			}
			return Height;
		},
		imgHeight() {
			let Height = "100vh";
			if (this.showFooterFlag_FlipMode || this.showHeaderFlag_FlipMode) {
				if (this.showPageNumFlag_FlipMode) {
					Height = "92vh"
				} else {
					Height = "95vh"
				}

			}
			return Height;
		},
	},

});
</script>

<style scoped>
/* 参考CSS盒子模型慢慢改 */
/* https://www.runoob.com/css/css-boxmodel.html */
/* CSS 高度和宽度 */
/* https://www.w3school.com.cn/css/css_dimension.asp */
/* CSS Flexbox 布局 */
/* https://www.w3school.com.cn/css/css3_flexbox.asp */

* {
	box-sizing: border-box;
}

.body {
	font-family: Arial;
	margin: 0;
	padding: 0px;
	display: flex;
	max-width: 100%;
	/* flex-direction: column垂直堆叠 flex 项目（从上到下），column-reverse从下到上 row从左到右 row-reverse从右到左 */
	flex-direction: column;
	/* justify-content 属性用于对齐 flex 项目。 将 justify-content 和 align-items 属性设置为居中，flex 项目将完美居中： */
	justify-content: center;
	/* center 值将 flex 项目在容器中间对齐： */
	align-items: center;
	background: v-bind("model.color");
}

.header {
	padding: 0px;
	width: 100%;
	height: 5vh;
}

/* 漫画div */
.main_manga {
	width: 100vw;
	height: v-bind(mainHeight);
	max-height: 100vh;
	max-width: 100vw;
	padding: 0px;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: baseline;
}

/* 漫画div中的图片*/
.main_manga img {
	/* max-height: inherit 继承 */
	/* max-height: inherit; */
	max-height: v-bind(imgHeight);
	max-width: 100vw;

	margin: 0 auto; /* center */
	padding: 0px;
	left: 0; /* center */
	right: 0; /* center */
	background-color: #aaa;
	padding: 0px;
	border-radius: 3px;
	box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
	/* flex-grow: 100; */
}

.main_manga span {
	height: 3vh;
	/* padding: 0px; */
	text-align: center;
	font-size: 20px;
	/* flex-grow: 10; */
	width: 100vw;
}

/* 页脚 */
.footer {
	height: 5vh;
	padding: 10px;
	text-align: center;
	background: #f6f7eb;
	width: 80vw;
	padding: 0px;
}

.footer div {
	height: 5vh;
	display: flex;
	justify-content: row;
	/* center 值将 flex 项目在容器中间对齐： */
	align-items: center;
	.span {
		width: 28px;
	}
}

.footer div > span {
	width: 8vw;
}
</style>
