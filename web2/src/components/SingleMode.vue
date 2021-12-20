<template>
	<div id="SingleMode" v-if="this.book">
		<Header v-if="this.showHeader">
			<n-space justify="space-between">
				<!-- 放本书占位，以后放返回箭头 -->
				<!-- SVG资源来自 https://www.xicons.org/#/ -->
				<n-icon size="40">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						xmlns:xlink="http://www.w3.org/1999/xlink"
						viewBox="0 0 512 512"
					>
						<path
							d="M256 160c16-63.16 76.43-95.41 208-96a15.94 15.94 0 0 1 16 16v288a16 16 0 0 1-16 16c-128 0-177.45 25.81-208 64c-30.37-38-80-64-208-64c-9.88 0-16-8.05-16-17.93V80a15.94 15.94 0 0 1 16-16c131.57.59 192 32.84 208 96z"
							fill="none"
							stroke="currentColor"
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="32"
						/>
						<path
							fill="none"
							stroke="currentColor"
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="32"
							d="M256 160v288"
						/>
					</svg>
				</n-icon>
				<!-- 标题，可下载压缩包 -->
				<n-space>
					<h2 v-if="book.IsFolder" :href="'raw/' + book.name">{{ book.name }}</h2>
					<h2>
						<a v-if="!book.IsFolder" :href="'raw/' + book.name">{{ book.name }}</a>
					</h2>
				</n-space>
				<!-- 右边的设置图标，点击屏幕中央也可以打开设置 -->
				<n-icon size="40" @click="drawerActivate('right')">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						xmlns:xlink="http://www.w3.org/1999/xlink"
						viewBox="0 0 512 512"
					>
						<path
							d="M262.29 192.31a64 64 0 1 0 57.4 57.4a64.13 64.13 0 0 0-57.4-57.4zM416.39 256a154.34 154.34 0 0 1-1.53 20.79l45.21 35.46a10.81 10.81 0 0 1 2.45 13.75l-42.77 74a10.81 10.81 0 0 1-13.14 4.59l-44.9-18.08a16.11 16.11 0 0 0-15.17 1.75A164.48 164.48 0 0 1 325 400.8a15.94 15.94 0 0 0-8.82 12.14l-6.73 47.89a11.08 11.08 0 0 1-10.68 9.17h-85.54a11.11 11.11 0 0 1-10.69-8.87l-6.72-47.82a16.07 16.07 0 0 0-9-12.22a155.3 155.3 0 0 1-21.46-12.57a16 16 0 0 0-15.11-1.71l-44.89 18.07a10.81 10.81 0 0 1-13.14-4.58l-42.77-74a10.8 10.8 0 0 1 2.45-13.75l38.21-30a16.05 16.05 0 0 0 6-14.08c-.36-4.17-.58-8.33-.58-12.5s.21-8.27.58-12.35a16 16 0 0 0-6.07-13.94l-38.19-30A10.81 10.81 0 0 1 49.48 186l42.77-74a10.81 10.81 0 0 1 13.14-4.59l44.9 18.08a16.11 16.11 0 0 0 15.17-1.75A164.48 164.48 0 0 1 187 111.2a15.94 15.94 0 0 0 8.82-12.14l6.73-47.89A11.08 11.08 0 0 1 213.23 42h85.54a11.11 11.11 0 0 1 10.69 8.87l6.72 47.82a16.07 16.07 0 0 0 9 12.22a155.3 155.3 0 0 1 21.46 12.57a16 16 0 0 0 15.11 1.71l44.89-18.07a10.81 10.81 0 0 1 13.14 4.58l42.77 74a10.8 10.8 0 0 1-2.45 13.75l-38.21 30a16.05 16.05 0 0 0-6.05 14.08c.33 4.14.55 8.3.55 12.47z"
							fill="none"
							stroke="currentColor"
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="32"
						/>
					</svg>
				</n-icon>
			</n-space>
		</Header>

		<div id="SinglePageTemplate">
			<div class="single_page_main" @click="getMouseXY($event)">
				<img
					v-on:click="addPage(1)"
					v-if="now_page <= this.book.all_page_num && now_page >= 1"
					lazy-src="/resources/favicon.ico"
					v-bind:src="this.book.pages[now_page - 1].url"
				/>
				<img />
			</div>
			<slot></slot>
		</div>

		<n-drawer v-model:show="drawerActive" :height="275" :width="251" :placement="drawerPlacement">
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
					<n-switch size="large" v-model:value="this.showHeaderFlag" @update:value="setShowHeaderChange">
						<template #checked>显示页头</template>
						<template #unchecked>显示页头</template>
					</n-switch>
				</n-space>

				<!-- 开关：是否显示页头 -->
				<n-space>
					<n-switch size="large" v-model:value="this.showHeaderFlag" @update:value="setShowHeaderChange">
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
						<template #checked>显示页数</template>
						<template #unchecked>显示页数</template>
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
import { NDrawer, NDrawerContent, NSpace, NSlider, NRadioButton, NRadioGroup, NSwitch, NIcon, } from 'naive-ui'

export default defineComponent({
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
			book: null,
			//是否显示页头
			showHeaderFlag: true,
			//是否显示分页
			showPaginationFlag: true,//https://www.naiveui.com/zh-CN/os-theme/components/pagination
			//鼠标点击或触摸的位置
			clickX: 0,
			clickY: 0,
			//图片宽度的单位，是否使用百分比
			imageWidth_usePercentFlag: true,
			now_page: 1,
		};
	},
	//在选项API中使用 Vue 生命周期钩子

	created() {
		//是否显示页头
		if (this.cookies.get("showHeaderFlag_SingleMode") === "true") {
			this.showHeaderFlag = true;
		} else if (this.cookies.get("showHeaderFlag_SingleMode") === "false") {
			this.showHeaderFlag = false;
		} else {
			this.cookies.set("showHeaderFlag_SingleMode", this.showHeaderFlag);
		}

		//是否显示分页导航
		if (this.cookies.get("showPaginationFlag_SingleMode") === "true") {
			this.showPaginationFlag = true;
		} else if (this.cookies.get("showPaginationFlag_SingleMode") === "false") {
			this.showPaginationFlag = false;
		} else {
			this.cookies.set("showPaginationFlag_SingleMode", this.showPaginationFlag);
		}
	},
	//挂载前
	beforeMount() {
		// window.addEventListener("scroll", this.handleScroll);
		// 注册监听
		window.addEventListener("keyup", this.handleKeyup);
		this.axios
			.get("/book.json")
			.then((response) => {
				if (response.status == 200) {
					this.book = response.data;
					//console.log(this.book);
				}
			}).catch((error) => console.log(error),);
	},
	//卸载前
	beforeUnmount() {
		// 销毁监听
		window.removeEventListener("keyup", this.handleKeyup);
		// window.removeEventListener("scroll", this.handleScroll);
	},

	methods: {

		setShowHeaderChange(value) {
			console.log("value:" + value);
			this.showHeaderFlag = value;
			this.cookies.set("showHeaderFlag", value);
			console.log("cookie设置完毕: showHeaderFlag=" + this.cookies.get("showHeaderFlag"));
		},
		setShowPaginationFlagChange(value) {
			console.log("value:" + value);
			this.showPageshowPaginationFlagNumFlag = value;
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
			//浏览器的视口，不包括工具栏和滚动条:
			// document.documentElement.clientHeight document.documentElement.ClientWidth不兼容手机？ 
			// var availHeight = document.documentElement.clientHeight
			// var availWidth = document.documentElement.clientWidth
			// console.log("clientHeigh=", document.documentElement.clientHeight, "ClientWidth=", document.documentElement.clientWidth);

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

		addPage: function (num) {
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
					this.addPage(-1); //上一页
					break;
				case "Space":
				case "ArrowDown":
				case "PageDown":
				case "ArrowRight":
					this.addPage(1); //下一页
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
.single_page_main {
	width: 100%;
	height: 95vh;
	display: flex;
	justify-content: center;
	align-items: center;
}

.single_page_main img {
	max-width: 100%;
	max-height: 100%;
	display: block;
	margin: center;
}
</style>
