<template>
	<div class="MainBody" v-if="this.book">
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
			class="main"
			@click="getMouseXY($event)"
			@mouseover="mouseOver"
			:style="randomColor"
			@mouseleave="mouseLeave"
		>
			<img
				v-if="now_page <= this.book.all_page_num && now_page >= 1"
				lazy-src="/resources/favicon.ico"
				v-bind:src="this.book.pages[now_page - 1].url"
			/>
			<img />
		</div>

		<n-pagination
			class="footer"
			v-if="this.showPaginationFlag"
			v-model:page="now_page"
			:page-count="this.book.all_page_num"
		/>

		<n-drawer
			v-model:show="drawerActive"
			@update:show="saveConfigToCookie"
			:height="275"
			:width="251"
			:placement="drawerPlacement"
		>
			<n-drawer-content title="页面设置" closable>
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
							:checked="selectedTemplate === 'single'"
							@change="onChangeTemplate"
							value="single"
							name="basic-demo"
						>单页模式</n-radio-button>
					</n-radio-group>
				</n-space>

				<!-- 开关：开启自动双页模式 -->
				<n-space>
					<n-switch
						size="large"
						v-model:value="this.showHeaderFlag_FlipMode"
						@update:value="setShowHeaderChange"
					>
						<template #checked>显示页头</template>
						<template #unchecked>显示页头</template>
					</n-switch>
				</n-space>

				<!-- 开关：是否显示页数导航条） -->
				<n-space>
					<n-switch
						size="large"
						v-model:value="this.showPaginationFlag"
						@update:value="setShowPaginationFlagChange"
					>
						<template #checked>显示底部导航条</template>
						<template #unchecked>显示底部导航条</template>
					</n-switch>
				</n-space>

				<n-space vertical>
					<!-- 记忆阅读进度？ -->
					<p v-if="this.width_usePercentFlag">单页漫画宽度（%）</p>
					<n-slider
						v-if="this.width_usePercentFlag"
						v-model:value="this.singlePageWidth_Percent"
						:step="1"
						:max="100"
						:min="10"
						:format-tooltip="value => `${value}%`"
						:marks="marks"
					/>
				</n-space>
			</n-drawer-content>
		</n-drawer>
	</div>
</template>

<script>
import { useCookies } from "vue3-cookies";
import Header from "@/components/Header.vue";
import { defineComponent, ref } from 'vue'
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
import { NDrawer, NDrawerContent, NSpace, NSlider, NRadioButton, NRadioGroup, NSwitch, NIcon, NPagination } from 'naive-ui'
import { SettingsOutline } from '@vicons/ionicons5'
export default defineComponent({
	name: "FlipMode",
	props: ['book'],
	components: {
		Header,
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
		NPagination, //分页 https://www.naiveui.com/zh-CN/os-theme/components/pagination
		SettingsOutline,//图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）
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
		};
	},
	data() {
		return {
			//开发模式 还没有做的功能与设置，设置Debug以后才能见到
			debugModeFlag: true,
			//书籍数据，需要从远程拉取
			//是否显示页头
			showHeaderFlag_FlipMode: true,
			//是否显示分页
			showPaginationFlag: true,//https://www.naiveui.com/zh-CN/os-theme/components/pagination
			//鼠标点击或触摸的位置
			clickX: 0,
			clickY: 0,
			//图片宽度的单位，是否使用百分比
			imageWidth_usePercentFlag: true,
			now_page: 1,
			randomColor: ""
		};
	},
	//在选项API中使用 Vue 生命周期钩子

	created() {
		//是否显示页头
		if (this.cookies.get("showHeaderFlag_FlipMode") === "true") {
			this.showHeaderFlag_FlipMode = true;
		} else if (this.cookies.get("showHeaderFlag_FlipMode") === "false") {
			this.showHeaderFlag_FlipMode = false;
		}

		//是否显示分页导航
		if (this.cookies.get("showPaginationFlag_FlipMode") === "true") {
			this.showPaginationFlag = true;
		} else if (this.cookies.get("showPaginationFlag_FlipMode") === "false") {
			this.showPaginationFlag = false;
		}
		setInterval(this.changeColor, 1000);

	},
	//挂载前
	beforeMount() {
		// window.addEventListener("scroll", this.handleScroll);
		// 注册监听
		window.addEventListener("keyup", this.handleKeyup);
	},
	//卸载前
	beforeUnmount() {
		// 销毁监听
		window.removeEventListener("keyup", this.handleKeyup);
		// window.removeEventListener("scroll", this.handleScroll);
	},

	methods: {
		changeColor() {
			// let R = Math.ceil(Math.random() * 255);
			// let G = Math.ceil(Math.random() * 255);
			// let B = Math.ceil(Math.random() * 255);
			let R = Math.ceil(Math.random() * 155)+100;
			let G = Math.ceil(Math.random() * 155)+100;
			let B = Math.ceil(Math.random() * 100)+100;
			// this.randomColor = 'background-color: rgb(235,235,235)';
			//rgb(235,235,235)
			let RGB = 'rgb(' + R + "," + G + "," + B + ")";
			// console.log(RGB);
			this.randomColor = RGB;
		},
		//HTML DOM 事件 https://www.runoob.com/jsref/dom-obj-event.html
		mouseOver() {
			//鼠标移入改变样式
			// this.randomColor = 'background-color: rgb(235,235,235)';
		},
		mouseLeave() {
			//清空样式
			// this.randomColor = '';
		},

		// 关闭抽屉的时候保存设置
		saveConfigToCookie() {
			// console.log("show:" + show)
			// 组件销毁前，储存cookie
			this.cookies.set("showHeaderFlag_FlipMode", this.showHeaderFlag_FlipMode);
			this.cookies.set("showPaginationFlag", this.showPaginationFlag);

		},
		setShowHeaderChange(value) {
			console.log("value:" + value);
			this.showHeaderFlag_FlipMode = value;
			this.cookies.set("showHeaderFlag_FlipMode", value);
			console.log("cookie设置完毕: showHeaderFlag_FlipMode=" + this.cookies.get("showHeaderFlag_FlipMode"));
		},
		setShowPaginationFlagChange(value) {
			console.log("value:" + value);
			this.showPaginationFlag = value;
			this.cookies.set("showPaginationFlag", value);
			console.log("cookie设置完毕: showPaginationFlag=" + this.cookies.get("showPaginationFlag"));
		},


		//切换模板的函数，需要配合vue-router
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
			location.reload(); //暂时无法动态刷新，研究vue-router去掉
		},

		//获取鼠标位置，然后做点什么
		getMouseXY(e) {
			this.clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离，单位为px）
			this.clickY = e.y //获取鼠标的Y坐标（鼠标与屏幕顶部的距离，单位为px）
			//浏览器的可视范围，不包括工具栏和滚动条:
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
			} else {
				this.flipPage(1);
			}
			// console.log("window.innerWidth=", window.innerWidth, "window.innerHeight=", window.innerHeight);
			// console.log("MinX=", MinX, "MaxX=", MaxX);
			// console.log("MinY=", MinY, "MaxY=", MaxY);
			// console.log("x=", e.x, "y=", e.y);
		},

		flipPage: function (num) {
			if (
				this.now_page + num <= this.book.all_page_num &&
				this.now_page + num >= 1
			) {
				this.now_page = this.now_page + num;
			}
			// console.log(this.now_page);
		},
		toPage: function (num) {
			if (num <= this.book.all_page_num && num >= 1) {
				this.now_page = num;
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

.MainBody {
	font-family: Arial;
	margin: 0;
	padding: 0px;
	display: flex;
	/* column 值设置垂直堆叠 flex 项目（从上到下）： */
	flex-direction: column;
	/* justify-content 属性用于对齐 flex 项目： */
	/* 将 justify-content 和 align-items 属性设置为居中，flex 项目将完美居中： */
	justify-content: center;
	/* center 值将 flex 项目在容器中间对齐： */
	align-items: center;
	background: v-bind(randomColor);
}

.header {
	padding: 0px;
	text-align: center;
	width: 100%;
}

/* 漫画本身 */
.main {
	height: 100vh;
	padding: 0px;
}

/* 漫画div中的图片*/
.main img {
	background-color: #aaa;
	height: 100vh;
	/* width: 100%; */
	padding: 0px;
	border-radius: 3px;
	box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

/* 页脚 */
.footer {
	padding: 10px;
	text-align: center;
	background: rgba(221, 221, 221, 0.842);
	width: 100vw;
	justify-content: center;
	/* center 值将 flex 项目在容器中间对齐： */
	align-items: center;
}
</style>
