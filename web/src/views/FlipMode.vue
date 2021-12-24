<template>
	<div class="body" v-if="this.book">
		<!-- 顶部，标题页头 -->
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

		<!-- 主题，漫画div -->
		<div
			class="manga_area"
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
				v-if="now_page_FlipMode <= this.book.all_page_num && now_page_FlipMode >= 1"
				lazy-src="/resources/favicon.ico"
				v-bind:src="this.book.pages[now_page_FlipMode - 1].url"
			/>
			<img />
			<div v-if="this.showPageNumFlag_FlipMode" class="sketch_hint">{{ pageNumOrSketchHint }}</div>
		</div>

		<!-- 页脚 拖动条 -->
		<div class="footer" v-if="this.showFooterFlag_FlipMode">
			<!-- 右手模式用 ，底部滑动条 -->
			<div v-if="this.useRightHalfScreen_FlipMode">
				<span>{{ this.now_page_FlipMode }}</span>
				<n-slider
					v-model:value="now_page_FlipMode"
					:max="this.book.all_page_num"
					:min="1"
					:step="1"
					:format-tooltip="value => `${value}`"
				/>
				<span>{{ this.book.all_page_num }}</span>
			</div>
			<!-- 左手模式用 底部滑动条，设置reverse翻转计数方向 -->
			<div v-if="!this.useRightHalfScreen_FlipMode">
				<span>{{ this.book.all_page_num }}</span>
				<n-slider
					reverse
					v-model:value="now_page_FlipMode"
					:max="this.book.all_page_num"
					:min="1"
					:step="1"
					:format-tooltip="value => `${value}`"
				/>
				<span>{{ this.now_page_FlipMode }}</span>
			</div>
		</div>
	</div>

	<!-- 设置抽屉，一开始隐藏 -->
	<Drawer
		:initDrawerActive="this.drawerActive"
		:initDrawerPlacement="this.drawerPlacement"
		@saveConfig="this.saveConfigToCookie"
		@startSketch="this.startSketchMode"
		@stopSketch="this.stopSketchMode"
		@closeDrawer="this.drawerDeactivate"
		@setT="OnSetTemplate"
		:nowTemplateDrawer="this.nowTemplate"
	>

		<span>{{ this.$t('message.setBackColor') }}</span>
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
				<template #checked>{{ this.$t('message.showHeader') }}</template>
				<template #unchecked>{{ this.$t('message.showHeader') }}</template>
			</n-switch>
		</n-space>

		<!-- 开关：显示阅读进度条） -->
		<n-space>
			<n-switch
				size="large"
				v-model:value="this.showFooterFlag_FlipMode"
				@update:value="setShowFooterFlagChange"
			>
				<template #checked>{{ this.$t('message.readingProgressBar') }}</template>
				<template #unchecked>{{ this.$t('message.readingProgressBar') }}</template>
			</n-switch>
		</n-space>

		<!-- 开关：显示当前页数 -->
		<n-space>
			<n-switch
				size="large"
				v-model:value="this.showPageNumFlag_FlipMode"
				@update:value="setShowPageNumChange"
			>
				<template #checked>{{ this.$t('message.showPageNum') }}</template>
				<template #unchecked>{{ this.$t('message.showPageNum') }}</template>
			</n-switch>
		</n-space>

		<!-- 分割线 -->
		<n-divider />

		<!-- 开关：Debug，现在只会随机背景色 -->
		<n-space>
			<n-switch
				size="large"
				v-model:value="this.debugModeFlag"
				@update:value="this.setDebugModeFlagFlag"
			>
				<template #checked>{{ this.$t('message.debugMode') }}</template>
				<template #unchecked>{{ this.$t('message.debugMode') }}</template>
			</n-switch>
		</n-space>

		<!-- 保存阅读进度 -->
		<n-space>
			<n-switch
				size="large"
				v-model:value="this.savePageNumFlag_FlipMode"
				@update:value="this.setSavePageNumFlag"
			>
				<template #checked>{{ this.$t('message.savePageNum') }}</template>
				<template #unchecked>{{ this.$t('message.savePageNum') }}</template>
			</n-switch>
		</n-space>

		<!-- 开关：翻页模式，默认右半屏-->
		<n-space>
			<n-switch
				size="large"
				v-model:value="this.useRightHalfScreen_FlipMode"
				:rail-style="railStyle"
				@update:value="this.setFlipScreenFlag"
			>
				<template #checked>{{ this.$t('message.rightScreenToNext') }}</template>
				<template #unchecked>{{ this.$t('message.leftScreenToNext') }}</template>
			</n-switch>
		</n-space>
	</Drawer>
</template>

<script>
import { useCookies } from "vue3-cookies";
// 自定义组件
import Header from "@/components/Header.vue";
import Drawer from "@/components/Drawer.vue";
import { defineComponent, reactive } from 'vue'
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
import { NSpace, NSlider, NSwitch, NIcon, NColorPicker, NDivider, } from 'naive-ui'
import { SettingsOutline } from '@vicons/ionicons5'
export default defineComponent({
	name: "FlipMode",
	props: ['book', 'nowTemplate'],
	emits: ["setTemplate"],
	components: {
		Header,
		Drawer,
		NSpace,//间距 https://www.naiveui.com/zh-CN/os-theme/components/space
		NSlider,//滑动选择  Slider https://www.naiveui.com/zh-CN/os-theme/components/slider
		NSwitch,//开关   https://www.naiveui.com/zh-CN/os-theme/components/switch
		// NLayout,//布局 https://www.naiveui.com/zh-CN/os-theme/components/layout
		// NLayoutSider,
		// NLayoutContent,
		NIcon,//图标  https://www.naiveui.com/zh-CN/os-theme/components/icon
		SettingsOutline,//图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）
		NColorPicker, //颜色选择器 Color Picker https://www.naiveui.com/zh-CN/os-theme/components/color-picker
		NDivider,//分割线  https://www.naiveui.com/zh-CN/os-theme/components/divider
	},
	setup() {
		const { cookies } = useCookies();
		//设置抽屉
		//颜色选择器
		const model = reactive({
			color: '#252525'
		})
		return {
			model,
			cookies,
			//开关用的颜色
			railStyle: ({ focused, checked }) => {
				const style = {}
				if (checked) {
					style.background = '#18a058'
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
			drawerActive: false,
			drawerPlacement: 'right',
			//开发模式 没做的功能与设置，设置Debug以后才能见到
			debugModeFlag: false,
			//是否显示页头
			showHeaderFlag_FlipMode: true,
			//是否显示页脚
			showFooterFlag_FlipMode: true,
			//是否是右半屏翻页
			useRightHalfScreen_FlipMode: true,
			//是否拼合双叶
			autoDoublepage_FlipMode: true,
			//是否保存当前页数
			savePageNumFlag_FlipMode: true,
			//当前页数
			now_page_FlipMode: 1,
			//素描模式标记
			sketchModeFlag: false,
			//是否显示素描提示
			showPageNumFlag_FlipMode: false,
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
		if (this.cookies.get("useRightHalfScreen_FlipMode") === "true") {
			this.useRightHalfScreen_FlipMode = true;
		} else if (this.cookies.get("useRightHalfScreen_FlipMode") === "false") {
			this.useRightHalfScreen_FlipMode = false;
		}

		if (this.cookies.get("autoDoublepage_FlipMode") === "true") {
			this.autoDoublepage_FlipMode = true;
		} else if (this.cookies.get("autoDoublepage_FlipMode") === "false") {
			this.autoDoublepage_FlipMode = false;
		}

		if (this.cookies.get("savePageNumFlag_FlipMode") === "true") {
			this.savePageNumFlag_FlipMode = true;
		} else if (this.cookies.get("savePageNumFlag_FlipMode") === "false") {
			this.savePageNumFlag_FlipMode = false;
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
		//当前页数
		if (this.cookies.get("now_page_FlipMode" + this.book.name) != null) {
			let saveNum = Number(this.cookies.get("now_page_FlipMode" + this.book.name));
			if (!isNaN(saveNum)) {
				this.now_page_FlipMode = saveNum;
			}
		}
		//自動開始Sketch模式
		if (this.nowTemplate == "sketch") {
			this.startSketchMode();
		}
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
		//界面有更新就会调用，随便乱放会引起难以调试的BUG
	},
	methods: {
		//接收Drawer的参数，继续往父组件传
		OnSetTemplate(value) {
			this.$emit("setTemplate", value);
		},
		sayHello(name) {
			alert("hello " + name);
		},
		//打开抽屉
		drawerActivate(place) {
			this.drawerActive = true
			this.drawerPlacement = place
		},
		//关闭抽屉
		drawerDeactivate() {
			this.drawerActive = false
		},
		//开始速写倒计时
		startSketchMode() {
			this.sketchModeFlag = true;
			this.showPageNumFlag_FlipMode = true;
			//是否显示页头
			this.showHeaderFlag_FlipMode=false,
			//是否显示页脚
			this.showFooterFlag_FlipMode=false,
			this.$emit("setTemplate", "sketch");
			//setTimeout和setInterval函数，都返回一个表示计数器编号的整数值，将该整数传入clearTimeout和clearInterval函数，就可以取消对应的定时器。setInterval指定某个任务每隔一段时间就执行一次。setTimeout()用于在指定的毫秒数后调用函数或计算表达式  setTimeout('console.log(2)',1000);
			this.interval = setInterval(this.sketchCount, 1000);
		},
		//停止速写倒计时
		stopSketchMode() {
			this.sketchModeFlag = false;
			this.showPageNumFlag_FlipMode = false;
			this.sketchSecondCount = 0;
			//是否显示页头
			this.showHeaderFlag_FlipMode=true,
			//是否显示页脚
			this.showFooterFlag_FlipMode=true,
			this.$emit("setTemplate", "flip");
			clearInterval(this.interval); // 清除定时器
		},
		sketchCount() {
			this.sketchSecondCount = this.sketchSecondCount + 1;
			let nowSeconnd = this.sketchSecondCount % this.sketchFlipSecond
			console.log("sketchSecondCount=" + this.sketchSecondCount + " nowSeconnd:" + nowSeconnd)
			if (nowSeconnd == 0) {
				if (this.now_page_FlipMode < this.book.all_page_num) {
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
			this.cookies.set("useRightHalfScreen_FlipMode", this.useRightHalfScreen_FlipMode);
			this.cookies.set("autoDoublepage_FlipMode", this.autoDoublepage_FlipMode);
			this.cookies.set("savePageNumFlag_FlipMode", this.savePageNumFlag_FlipMode);
			this.cookies.set("now_page_FlipMode" + this.book.name, this.now_page_FlipMode);
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
					if (this.useRightHalfScreen_FlipMode) {
						this.flipPage(-1);
					} else {
						this.flipPage(1);
					}
				} else {
					//右边的翻页
					if (this.useRightHalfScreen_FlipMode) {
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

		setFlipScreenFlag(value) {
			console.log("value:" + value);
			this.useRightHalfScreen_FlipMode = value;
			this.cookies.set("useRightHalfScreen_FlipMode", value);
			console.log("cookie设置完毕: useRightHalfScreen_FlipMode=" + this.cookies.get("useRightHalfScreen_FlipMode"));
		},

		setSavePageNumFlag(value) {
			console.log("value:" + value);
			this.savePageNumFlag_FlipMode = value;
			this.cookies.set("savePageNumFlag_FlipMode", value);
			console.log("cookie设置完毕: savePageNumFlag_FlipMode=" + this.cookies.get("savePageNumFlag_FlipMode"));
		},

		setDebugModeFlagFlag(value) {
			console.log("value:" + value);
			this.debugModeFlag = value;
			this.cookies.set("debugModeFlag", value);
			console.log("cookie设置完毕: debugModeFlag=" + this.cookies.get("debugModeFlag"));
		},

		flipPage: function (num) {
			if (
				this.now_page_FlipMode + num <= this.book.all_page_num &&
				this.now_page_FlipMode + num >= 1
			) {
				this.now_page_FlipMode = this.now_page_FlipMode + num;
			} else {
				// console.log("无法继续翻，Num:" + num)
				if (num > 0) {
					alert(this.$t('message.hintLastPage'));
				} else {
					alert(this.$t('message.hintFirstPage'));
				}
			}
			if (this.savePageNumFlag_FlipMode) {
				this.cookies.set("now_page_FlipMode" + this.book.name, this.now_page_FlipMode);
			}
		},
		toPage: function (num) {
			if (num <= this.book.all_page_num && num >= 1) {
				this.now_page_FlipMode = num;
			}
			if (this.savePageNumFlag_FlipMode) {
				this.cookies.set("now_page_FlipMode" + this.book.name, this.now_page_FlipMode);
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
		//页数或素描模式的提示
		pageNumOrSketchHint() {
			if (this.sketchModeFlag) {
				let nowSeconnd = (this.sketchSecondCount % this.sketchFlipSecond)+1
				let AllString =parseInt((this.sketchSecondCount+1)/60)+"分 "+ (this.sketchSecondCount+1)%60+"秒"
				let hintString = "现在:" + nowSeconnd + "秒    总共:" + AllString + "   翻页间隔:" + this.sketchFlipSecond
				return hintString
			} else {
				return this.now_page_FlipMode + "/" + this.book.all_page_num
			}
		},
		mangaAreaHeight() {
			let Height = 95
			//页头和底部拖动条都显示,或有一个显示的时候，95%
			if (this.showFooterFlag_FlipMode && this.showHeaderFlag_FlipMode) {
				Height = 95
			}
			if (this.showFooterFlag_FlipMode && !this.showHeaderFlag_FlipMode) {
				Height = 95
			}
			if (!this.showFooterFlag_FlipMode && this.showHeaderFlag_FlipMode) {
				Height = 95
			}
			//页头和底部拖动条都不显示的时候，漫画占满屏幕
			if ((!this.showFooterFlag_FlipMode) && (!this.showHeaderFlag_FlipMode)) {
				Height = 100
			}
			return Height + "vh";

		},
		mangaImageHeight() {
			let Height = 95
			//页头和底部拖动条都显示,或有一个显示的时候，95%
			if (this.showFooterFlag_FlipMode && this.showHeaderFlag_FlipMode) {
				Height = 95
			}
			if (this.showFooterFlag_FlipMode && !this.showHeaderFlag_FlipMode) {
				Height = 95
			}
			if (!this.showFooterFlag_FlipMode && this.showHeaderFlag_FlipMode) {
				Height = 95
			}
			//页头和拖动条都不显示的时候，漫画占满屏幕
			if ((!this.showFooterFlag_FlipMode) && (!this.showHeaderFlag_FlipMode)) {
				Height = 100
			}
			//与上面唯一的不同，减去素描提示的空间
			if (this.showPageNumFlag_FlipMode) {
				Height = Height - 3
			}
			return Height + "vh";
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
.manga_area {
	width: 100vw;
	height: v-bind(mangaAreaHeight);
	max-height: 100vh;
	max-width: 100vw;
	padding: 0px;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: baseline;
	user-select: none; /* 不可以被选中 */
	-moz-user-select: none; /* 火狐 */
	-webkit-user-select: none; /* 谷歌 */
}

/* 漫画div中的图片*/
.manga_area img {
	/* max-height: inherit 继承 */
	/* max-height: inherit; */
	max-height: v-bind(mangaImageHeight);
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

/* 漫画div图片下面的页数*/
.sketch_hint {
	height: 3vh;
	padding: 0px;
	text-align: center;
	font-size: 16px;
	/* 文字颜色 */
	color: rgb(238, 238, 238);
	/* 文字阴影：https://www.w3school.com.cn/css/css3_shadows.asp*/
	text-shadow: -1px 0 black, 0 1px black, 1px 0 black, 0 -1px black;
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
	color: rgb(80, 80, 255);
	text-shadow: 2px 2px 5px yellowgreen;
}

.footer div > span {
	width: 10vw;
}
</style>
