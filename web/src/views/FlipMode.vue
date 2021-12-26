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
		<!-- 事件修饰符： https://v3.cn.vuejs.org/guide/events.html#%E4%BA%8B%E4%BB%B6%E4%BF%AE%E9%A5%B0%E7%AC%A6 -->
		<div
			class="manga_area"
			id="MangaMain"
			@click.stop="onMouseClick"
			@mousemove.stop="onMouseMove"
			@mouseleave.stop="onMouseLeave"
		>
			<div class="manga_area_img_div">
				<!-- 非自动拼合模式最简单，直接显示一张图就成了 -->
				<img
					v-if="!this.autoDoublepage_FlipMode"
					v-bind:src="this.book.pages[nowPageNum_FlipMode - 1].url"
				/>

				<!-- 右手模式当前页，如果指定了自动拼合，右边可能显示拼合页 -->
				<!-- v-if="nowPageNum_FlipMode <= this.book.all_page_num && nowPageNum_FlipMode >= 1" -->
				<img
					v-if="this.rightScreenToNextFlag_FlipMode && this.autoDoublepage_FlipMode"
					v-bind:src="this.book.pages[nowPageNum_FlipMode - 1].url"
				/>
				<!-- 右手模式下一页，根据情况显示，自动拼合用 -->
				<img
					v-if="this.nowPageNum_FlipMode < book.all_page_num && this.rightScreenToNextFlag_FlipMode && this.autoDoublepage_FlipMode && !this.checkImageIsDoublePage_byPageNum(nowPageNum_FlipMode)&& !this.checkImageIsDoublePage_byPageNum(nowPageNum_FlipMode+1)"
					v-bind:src="this.book.pages[nowPageNum_FlipMode].url"
				/>

				<!-- 左手模式下一页，根据情况显示，自动拼合用 -->
				<img
					v-if="this.nowPageNum_FlipMode < book.all_page_num && !this.rightScreenToNextFlag_FlipMode && this.autoDoublepage_FlipMode && !this.checkImageIsDoublePage_byPageNum(nowPageNum_FlipMode)&& !this.checkImageIsDoublePage_byPageNum(nowPageNum_FlipMode+1)"
					v-bind:src="this.book.pages[nowPageNum_FlipMode].url"
				/>
				<!-- 左手模式当前页，如果指定了自动拼合，右边可能显示拼合页 -->
				<!-- v-if="nowPageNum_FlipMode <= this.book.all_page_num && nowPageNum_FlipMode >= 1" -->
				<img
					v-if="!this.rightScreenToNextFlag_FlipMode && this.autoDoublepage_FlipMode"
					v-bind:src="this.book.pages[nowPageNum_FlipMode - 1].url"
				/>
			</div>

			<div v-if="this.showPageHintFlag_FlipMode" class="sketch_hint">{{ pageNumOrSketchHint }}</div>
		</div>

		<!-- 页脚 拖动条 -->
		<div class="footer" v-if="this.showFooterFlag_FlipMode">
			<!-- 右手模式用 ，底部滑动条 -->
			<div v-if="this.rightScreenToNextFlag_FlipMode">
				<span>{{ this.nowPageNum_FlipMode }}</span>
				<n-slider
					v-model:value="nowPageNum_FlipMode"
					:max="this.book.all_page_num"
					:min="1"
					:step="1"
					:format-tooltip="(value) => `${value}`"
					@update:value="this.saveNowPageNumOnUpdate"
				/>
				<span>{{ this.book.all_page_num }}</span>
			</div>
			<!-- 左手模式用 底部滑动条，设置reverse翻转计数方向 -->
			<div v-if="!this.rightScreenToNextFlag_FlipMode">
				<span>{{ this.book.all_page_num }}</span>
				<n-slider
					reverse
					v-model:value="nowPageNum_FlipMode"
					:max="this.book.all_page_num"
					:min="1"
					:step="1"
					:format-tooltip="(value) => `${value}`"
					@update:value="this.saveNowPageNumOnUpdate"
				/>
				<span>{{ this.nowPageNum_FlipMode }}</span>
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
		:nowTemplate="this.nowTemplate"
	>
		<span>{{ $t("setBackColor") }}</span>
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
				<template #checked>{{ $t("showHeader") }}</template>
				<template #unchecked>{{ $t("showHeader") }}</template>
			</n-switch>
		</n-space>

		<!-- 开关：显示阅读进度条） -->
		<n-space>
			<n-switch
				size="large"
				v-model:value="this.showFooterFlag_FlipMode"
				@update:value="setShowFooterFlagChange"
			>
				<template #checked>{{ $t("readingProgressBar") }}</template>
				<template #unchecked>{{ $t("readingProgressBar") }}</template>
			</n-switch>
		</n-space>

		<!-- 开关：显示当前页数 -->
		<n-space>
			<n-switch
				size="large"
				v-model:value="this.showPageHintFlag_FlipMode"
				@update:value="setShowPageNumChange"
			>
				<template #checked>{{ $t("showPageNum") }}</template>
				<template #unchecked>{{ $t("showPageNum") }}</template>
			</n-switch>
		</n-space>

		<!-- 分割线 -->
		<n-divider />

		<!-- 保存阅读进度 -->
		<n-space>
			<n-switch
				size="large"
				v-model:value="this.saveNowPageNumFlag_FlipMode"
				@update:value="this.setSavePageNumFlag"
			>
				<template #checked>{{ $t("savePageNum") }}</template>
				<template #unchecked>{{ $t("savePageNum") }}</template>
			</n-switch>
		</n-space>

		<!-- 开关：翻页模式，默认右半屏-->
		<n-space>
			<n-switch
				size="large"
				v-model:value="this.rightScreenToNextFlag_FlipMode"
				:rail-style="railStyle"
				@update:value="this.setFlipScreenFlag"
			>
				<template #checked>{{ $t("rightScreenToNext") }}</template>
				<template #unchecked>{{ $t("leftScreenToNext") }}</template>
			</n-switch>
		</n-space>

		<!-- 开关：Debug，现在只会随机背景色 -->
		<n-space>
			<n-switch
				size="large"
				v-model:value="this.debugModeFlag"
				@update:value="this.setDebugModeFlag"
			>
				<template #checked>{{ $t("debugMode") }}</template>
				<template #unchecked>{{ $t("debugMode") }}</template>
			</n-switch>
		</n-space>


		<n-space v-if="this.debugModeFlag">
			<n-switch
				size="large"
				v-model:value="this.autoDoublepage_FlipMode"
				@update:value="this.setAutoDoublepage_FlipMode"
			>
				<template #checked>自动合并单双页</template>
				<template #unchecked>自动合并单双页</template>
			</n-switch>
		</n-space>


		<!-- 分割线 -->
		<n-divider v-if="this.nowTemplate == 'sketch'" />
		<!-- 自动翻页秒数 -->
		<!-- 数字输入% -->
		<n-input-number
			v-if="this.nowTemplate == 'sketch'"
			size="small"
			:show-button="false"
			v-model:value="this.sketchFlipSecond"
			:max="65535"
			:min="1"
			:update-value-on-input="false"
			@update:value="this.resetSketchSecondCount"
		>
			<template #prefix>{{ $t('pageTurningSeconds') }}</template>
			<template #suffix>{{ $t("second") }}</template>
		</n-input-number>
		<!-- 滑动选择% -->
		<n-slider
			v-if="this.nowTemplate == 'sketch'"
			v-model:value="this.sketchFlipSecond"
			:step="1"
			:max="120"
			:min="1"
			:marks="marks"
			:format-tooltip="value => `${value}s`"
			@update:value="this.resetSketchSecondCount"
		/>
	</Drawer>
</template>

<script>
import { useCookies } from "vue3-cookies";
// 自定义组件
import Header from "@/components/Header.vue";
import Drawer from "@/components/Drawer.vue";
import { defineComponent, reactive } from "vue";
// 直接导入组件并使用它。这种情况下，只有导入的组件才会被打包。
import { NSpace, NSlider, NSwitch, NIcon, NColorPicker, NDivider, NInputNumber, useMessage, } from "naive-ui";
import { SettingsOutline } from "@vicons/ionicons5";
export default defineComponent({
	name: "FlipMode",
	props: ["book", "nowTemplate"],
	emits: ["setTemplate"],
	components: {
		Header,
		Drawer,
		NSpace, //间距 https://www.naiveui.com/zh-CN/os-theme/components/space
		NSlider, //滑动选择  Slider https://www.naiveui.com/zh-CN/os-theme/components/slider
		NSwitch, //开关   https://www.naiveui.com/zh-CN/os-theme/components/switch
		// NLayout,//布局 https://www.naiveui.com/zh-CN/os-theme/components/layout
		// NLayoutSider,
		// NLayoutContent,
		NIcon, //图标  https://www.naiveui.com/zh-CN/os-theme/components/icon
		SettingsOutline, //图标,来自 https://www.xicons.org/#/   需要安装（npm i -D @vicons/ionicons5）
		NColorPicker, //颜色选择器 Color Picker https://www.naiveui.com/zh-CN/os-theme/components/color-picker
		NDivider, //分割线  https://www.naiveui.com/zh-CN/os-theme/components/divider
		NInputNumber,///  https://www.naiveui.com/zh-CN/os-theme/components/input-number
		// useNotification, // https://www.naiveui.com/zh-CN/os-theme/components/notification
	},
	setup() {
		const { cookies } = useCookies();
		//背景颜色，颜色选择器用
		const model = reactive({
			color: "#E0D9CD",
		});
		//警告信息
		const message = useMessage()
		// const notification = useNotification()
		return {
			message,
			//背景色
			model,
			cookies,
			//开关的颜色
			railStyle: ({ focused, checked }) => {
				const style = {};
				if (checked) {
					style.background = "#18a058";
					if (focused) {
						style.boxShadow = "0 0 0 2px #d0305040";
					}
				} else {
					style.background = "#2080f0";
					if (focused) {
						style.boxShadow = "0 0 0 2px #2080f040";
					}
				}
				return style;
			},
			//滑动秒数建议值
			marks: {
				30: '30',
				60: '60',
				90: '90',
				120: '120',
			},
		};
	},
	data() {
		return {
			drawerActive: false,
			drawerPlacement: "right",
			//开发模式 没做的功能与设置，设置Debug以后才能见到
			debugModeFlag: false,
			//是否显示页头
			showHeaderFlag_FlipMode: true,
			//是否显示页脚
			showFooterFlag_FlipMode: true,
			//是否是右半屏翻页
			rightScreenToNextFlag_FlipMode: true,
			//尝试拼合双叶
			autoDoublepage_FlipMode: true,
			//是否保存当前页数
			saveNowPageNumFlag_FlipMode: true,
			//当前页数，注意语义，直接就是1开始的页数，不是数组下标，在pages数组当中用的时候需要-1
			nowPageNum_FlipMode: 1,
			//素描模式标记
			sketchModeFlag: false,
			//是否显示素描提示
			showPageHintFlag_FlipMode: false,
			//翻页间隔时间
			sketchFlipSecond: 30,
			//计时用，从0开始
			sketchSecondCount: 0,
			//计时器ID
			interval: null,
		};
	},
	//在选项API中使用 Vue 生命周期钩子：
	created() {
		//从cookies初始化默认值
		if (this.cookies.get("debugModeFlag") === "true") {
			this.debugModeFlag = true;
		} else if (this.cookies.get("debugModeFlag") === "false") {
			this.debugModeFlag = false;
		}
		//是否显示标题
		if (this.cookies.get("showHeaderFlag_FlipMode") === "true") {
			this.showHeaderFlag_FlipMode = true;
		} else if (this.cookies.get("showHeaderFlag_FlipMode") === "false") {
			this.showHeaderFlag_FlipMode = false;
		}
		//是否显示页脚
		if (this.cookies.get("showFooterFlag_FlipMode_FlipMode") === "true") {
			this.showFooterFlag_FlipMode = true;
		} else if (this.cookies.get("showFooterFlag_FlipMode_FlipMode") === "false") {
			this.showFooterFlag_FlipMode = false;
		}
		//是否显示页数
		if (this.cookies.get("showPageHintFlag_FlipMode") === "true") {
			this.showPageHintFlag_FlipMode = true;
		} else if (this.cookies.get("showPageHintFlag_FlipMode") === "false") {
			this.showPageHintFlag_FlipMode = false;
		}
		//是否用右半屏翻页
		if (this.cookies.get("rightScreenToNextFlag_FlipMode") === "true") {
			this.rightScreenToNextFlag_FlipMode = true;
		} else if (this.cookies.get("rightScreenToNextFlag_FlipMode") === "false") {
			this.rightScreenToNextFlag_FlipMode = false;
		}
		//自动合并单页
		if (this.cookies.get("autoDoublepage_FlipMode") === "true") {
			this.autoDoublepage_FlipMode = true;
		} else if (this.cookies.get("autoDoublepage_FlipMode") === "false") {
			this.autoDoublepage_FlipMode = false;
		}
		//是否保存阅读进度
		if (this.cookies.get("saveNowPageNumFlag_FlipMode") === "true") {
			this.saveNowPageNumFlag_FlipMode = true;
		} else if (this.cookies.get("saveNowPageNumFlag_FlipMode") === "false") {
			this.saveNowPageNumFlag_FlipMode = false;
		}
		//当前背景色
		if (this.cookies.get("FlipModeDefaultColor") != null) {
			this.model.color = this.cookies.get("FlipModeDefaultColor");
		}
		//倒计时秒数
		if (this.cookies.get("sketchFlipSecond") != null) {
			let saveNum = Number(this.cookies.get("sketchFlipSecond"));
			if (!isNaN(saveNum)) {
				this.sketchFlipSecond = saveNum;
			}
		}
	},
	// beforeMount : 指令第一次绑定到元素并且在挂载父组件之前调用。
	beforeMount() {
		// 注册监听
		window.addEventListener("keyup", this.handleKeyup);
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
	// mounted : 在绑定元素的父组件被挂载后调用。
	mounted() {
		//需要得书籍远程数据，避免初始化失败，所以延迟0.5秒执行
		setTimeout(this.setNowPageNumByCookie, 500);
	},
	updated() {
		//界面有更新就会调用，随便乱放会引起难以调试的BUG
	},
	methods: {
		//根据书籍UUID，设定当前页数，因为需要取得远程书籍数据（this.book），所以延迟执行
		setNowPageNumByCookie() {
			if (this.saveNowPageNumFlag_FlipMode) {
				let cookieValue = this.cookies.get("nowPageNum_FlipMode" + this.book.uuid);
				if (cookieValue != null) {
					let saveNum = Number(cookieValue);
					if (!isNaN(saveNum)) {
						this.nowPageNum_FlipMode = saveNum;
						// console.log("成功读取页数" + saveNum);
					} else {
						console.log("读取页数失败，this.nowPageNum_FlipMode=" + this.nowPageNum_FlipMode);
					}
				} else {
					console.log("cookie里面没找到:" + "nowPageNum_FlipMode" + this.book.uuid);
				}
			} else {
				console.log("不读取页数，this.saveNowPageNumFlag_FlipMode=" + this.saveNowPageNumFlag_FlipMode);
			}
		},

		//设置当前模板-接收Drawer的参数，继续往父组件传
		OnSetTemplate(value) {
			this.$emit("setTemplate", value);
		},
		//打开抽屉
		drawerActivate(place) {
			this.drawerActive = true;
			this.drawerPlacement = place;
		},
		//关闭抽屉
		drawerDeactivate() {
			this.drawerActive = false;
		},
		//开始速写倒计时
		startSketchMode() {
			this.message.success(this.$t('startSketchMessage'));
			this.drawerActive = false; //关闭设置抽屉
			this.sketchModeFlag = true;
			this.showPageHintFlag_FlipMode = true;
			//是否显示页头
			(this.showHeaderFlag_FlipMode = false),
				//是否显示页脚
				(this.showFooterFlag_FlipMode = false),
				this.$emit("setTemplate", "sketch");
			//setTimeout和setInterval函数，都返回一个表示计数器编号的整数值，将该整数传入clearTimeout和clearInterval函数，就可以取消对应的定时器。setInterval指定某个任务每隔一段时间就执行一次。setTimeout()用于在指定的毫秒数后调用函数或计算表达式  setTimeout('console.log(2)',1000);
			this.interval = setInterval(this.sketchCount, 1000);
		},
		//修改间隔的时候重新计秒数
		resetSketchSecondCount() {
			this.sketchSecondCount = 0;
		},
		//停止速写倒计时
		stopSketchMode() {
			this.message.success(this.$t('goodjob_and_byebye'));
			this.sketchModeFlag = false;
			this.showPageHintFlag_FlipMode = false;
			this.sketchSecondCount = 0;
			//是否显示页头
			(this.showHeaderFlag_FlipMode = true),
				//是否显示页脚
				(this.showFooterFlag_FlipMode = true),
				this.$emit("setTemplate", "flip");
			clearInterval(this.interval); // 清除定时器
		},
		//开始速写（quick sketch），每秒执行一次
		sketchCount() {
			if (this.sketchModeFlag == false || this.nowTemplate != "sketch") {
				this.stopSketchMode();
			}
			this.sketchSecondCount = this.sketchSecondCount + 1;
			let nowSeconnd = this.sketchSecondCount % this.sketchFlipSecond;
			// console.log("sketchSecondCount=" + this.sketchSecondCount + " nowSeconnd:" + nowSeconnd)
			if (nowSeconnd == 0) {
				if (this.nowPageNum_FlipMode < this.book.all_page_num) {
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
			this.cookies.set("showPageHintFlag_FlipMode", this.showPageHintFlag_FlipMode);
			this.cookies.set("rightScreenToNextFlag_FlipMode", this.rightScreenToNextFlag_FlipMode);
			this.cookies.set("autoDoublepage_FlipMode", this.autoDoublepage_FlipMode);
			this.cookies.set("saveNowPageNumFlag_FlipMode", this.saveNowPageNumFlag_FlipMode);
			this.cookies.set("nowPageNum_FlipMode" + this.book.uuid, this.nowPageNum_FlipMode);
			this.cookies.set("FlipModeDefaultColor", this.model.color);
			this.cookies.set("sketchFlipSecond", this.sketchFlipSecond);
		},
		randomBackgroundColor() {
			let R = Math.ceil(Math.random() * 155) + 100;
			let G = Math.ceil(Math.random() * 155) + 100;
			let B = Math.ceil(Math.random() * 100) + 100;
			//rgb(185,175,145)
			let RGB = "rgb(" + R + "," + G + "," + B + ")";
			// console.log(RGB);
			this.model.color = RGB;
		},
		//HTML DOM 事件 https://www.runoob.com/jsref/dom-obj-event.html
		// 进入绑定该事件的元素和其子元素均会触发该事件，所以有一个重复触发，冒泡过程。其对应的离开事件 mouseout
		onMouseOver() { },
		// 只有进入绑定该事件的元素才会触发事件，也就是不会冒泡。其对应的离开事件mouseleave
		onMouseEnter() {
			// this.randomColor = 'background-color: rgb(235,235,235)';
		},
		onMouseLeave(e) {
			//离开区域的时候，清空鼠标样式
			e.currentTarget.style.cursor = "";
		},
		//事件修饰符: https://v3.cn.vuejs.org/guide/events.html#%E4%BA%8B%E4%BB%B6%E4%BF%AE%E9%A5%B0%E7%AC%A6
		onMouseMove(e) {
			// // offsetX/Y获取到是触发点相对被触发dom的左上角距离
			// let offsetX = e.offsetX;
			// let offsetY = e.offsetY;
			//根据ID获取元素
			// let mangaDiv =document.getElementById("MangaMain")
			//不用自己获取元素
			// let offsetWidth = e.currentTarget.offsetWidth;
			// let offsetHeight = e.currentTarget.offsetHeight;

			let clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离，单位为px）
			let clickY = e.y //获取鼠标的Y坐标（鼠标与屏幕顶部的距离，单位为px）
			//浏览器的视口，不包括工具栏和滚动条:
			let innerWidth = window.innerWidth
			let innerHeight = window.innerHeight

			var MinX = innerWidth * 0.4;
			var MaxX = innerWidth * 0.6;
			var MinY = innerHeight * 0.4;
			var MaxY = innerHeight * 0.6;
			if (clickX > MinX && clickX < MaxX && clickY > MinY && clickY < MaxY) {
				//设置区域;
				e.currentTarget.style.cursor = "url(/images/SettingsOutline.png), pointer";
			} else {
				if (clickX < innerWidth * 0.5) {
					//设置左边的鼠标指针
					if (this.rightScreenToNextFlag_FlipMode && this.nowPageNum_FlipMode == 1) {
						//右边翻下一页，且目前是第一页的时候，左边的鼠标指针，设置为禁止翻页
						e.currentTarget.style.cursor = "url(/images/Prohibited28Filled.png), pointer";
					} else if (!this.rightScreenToNextFlag_FlipMode && this.nowPageNum_FlipMode == this.book.all_page_num) {
						//左边翻下一页，且目前是最后一页的时候，左边的鼠标指针，设置为禁止翻页
						e.currentTarget.style.cursor = "url(/images/Prohibited28Filled.png), pointer";
					} else {
						//正常情况下，左边是向左的箭头
						e.currentTarget.style.cursor = "url(/images/ArrowLeft.png), pointer";
					}
				} else {
					//设置右边的鼠标指针
					if (this.rightScreenToNextFlag_FlipMode && this.nowPageNum_FlipMode == this.book.all_page_num) {
						//右边翻下一页，且目前是最后页的时候，右边的鼠标指针，设置为禁止翻页
						e.currentTarget.style.cursor = "url(/images/Prohibited28Filled.png), pointer";
					} else if (!this.rightScreenToNextFlag_FlipMode && this.nowPageNum_FlipMode == 1) {
						//左边翻下一页，且目前是第一页的时候，右边的鼠标指针，设置为禁止翻页
						e.currentTarget.style.cursor = "url(/images/Prohibited28Filled.png), pointer";
					} else {
						//正常情况下，右边是向右的箭头
						e.currentTarget.style.cursor = "url(/images/ArrowRight.png), pointer";
					}
				}
			}
		},

		//根据鼠标点击事件的位置，决定是左右翻页还是打开设置
		onMouseClick(e) {
			//元素自身的宽高，会有子元素问题
			// let offsetX = e.offsetX;
			// let offsetY = e.offsetY;
			// let offsetWidth = e.currentTarget.offsetWidth;
			// let offsetHeight = e.currentTarget.offsetHeight;

			let clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离，单位为px）
			let clickY = e.y //获取鼠标的Y坐标（鼠标与屏幕顶部的距离，单位为px）
			//浏览器的可视区域宽高，不包括工具栏和滚动条:
			let availHeight = window.innerWidth
			let availWidth = window.innerHeight
			var MinX = availHeight * 0.4;
			var MaxX = availHeight * 0.6;
			var MinY = availWidth * 0.4;
			var MaxY = availWidth * 0.6;
			// console.log("鼠标点击：e.offsetX=" + offsetX, "e.offsetY=" + offsetY);
			if (clickX > MinX && clickX < MaxX && clickY > MinY && clickY < MaxY) {
				//点中了设置区域
				this.drawerActivate("right");
			} else {
				// //随机一下背景色，只是为了好玩
				// if (this.debugModeFlag) {
				// 	this.randomBackgroundColor();
				// }
				//决定如何翻页
				if (clickX <= availHeight / 2.0) {
					//左边的翻页
					if (this.rightScreenToNextFlag_FlipMode) {
						this.toPerviousPage();
					} else {
						this.toNextPage();
					}
				} else {
					//右边的翻页
					if (this.rightScreenToNextFlag_FlipMode) {
						this.toNextPage();
					} else {
						this.toPerviousPage();
					}
				}
			}
		},
		toNextPage() {
			//如果开启了自动合并模式，并且当前页应该被合并
			if (this.autoDoublepage_FlipMode && this.checkMergedStatus_ByPageNum(this.nowPageNum_FlipMode)) {
				if (this.nowPageNum_FlipMode < this.book.all_page_num - 1) {
					this.flipPage(2);
				} else {
					this.flipPage(1);
				}
			} else {
				this.flipPage(1);
			}
		},
		toPerviousPage() {
			//错误值
			if (this.book.nowPageNum_FlipMode <= 1) {
				console.log("Error toPerviousPage");
				return;
			}
			//如果没有开启自动合并模式，或现在是第二页
			if (this.book.nowPageNum_FlipMode == 2 || !this.autoDoublepage_FlipMode) {
				this.flipPage(-1);
				return
			}
			//如果前一页是双开叶
			let pervIsDouble=this.checkImageIsDoublePage_byPageNum(this.nowPageNum_FlipMode - 2)
			if (pervIsDouble) {
				this.flipPage(-1);
				return
			}
			//如果前前页是双开叶
			let PervPervIsDouble=this.checkImageIsDoublePage_byPageNum(this.nowPageNum_FlipMode - 3)
			if (PervPervIsDouble) {
				this.flipPage(-1);
				return
			}
			//都没return掉，那么前两张可以合并，减两页
			this.flipPage(-2);
		},

		//给一个页数，然后判断自动双页模式下，是否应该预读并合并显示下一页
		checkMergedStatus_ByPageNum(pageNum) {
			//如果没有开启自动双页模式，当然不需要
			if (!this.autoDoublepage_FlipMode) {
				return false;
			}
			// //预读下下张图片，不做任何判断，相当于手动的lazy load
			// if (pageNum + 2 < this.book.all_page_num) {
			// 	this.book.all_page_num
			// 	let temp_image = new Image();
			// 	temp_image.src = this.book.pages[pageNum + 2].url;
			// }
			//已经是最后一页了，显然不需要自动合并下一页
			if (pageNum == this.book.all_page_num) {
				return false;
			}
			//可能传入的错误值，打印到控制台
			if (pageNum <= 0 || pageNum >= this.book.all_page_num) {
				console.log("Error pageNum :" + pageNum);
				return false;
			}
			//分析下一张漫画的宽高比,看是不是双开页
			let next_page_is_double_page = this.checkImageIsDoublePage_byPageNum(pageNum+1);
			//分析现在这张图片的宽高比,看是不是双开页
			let now_page_is_double_page = this.checkImageIsDoublePage_byPageNum(pageNum);
			//如果现在这张就是开页漫画，直接不用比了
			if (now_page_is_double_page) {
				return false;
			} else if (next_page_is_double_page) {
				//如果下一张漫画是开页，显然也没法合并过来
				return false;
			} else {
				//前两个条件不满足，也就是说两边都是单页，可以合并显示
				return true;
			}
		},
		checkImageIsDoublePage_byPageNum(pageNum) {
			let image = new Image();
			image.src = this.book.pages[pageNum-1].url;
			// 图片是否完全加载完成。
			//https://developer.mozilla.org/zh-CN/docs/Web/API/HTMLImageElement/complete
			if (image.complete) {
				if (image.width < image.height) {
					return false;//宽小于高，是是竖着的，单页漫画
				} else {
					return true;//宽大于高，很可能是开页
				}
			} else {
				//否则加载图片
				let temp_flag = false;
				image.onload = function () {
					//是单页漫画
					if (image.width < image.height) {
						image.onload = null;	// 避免重复加载
						temp_flag = false;//宽小于高，是是竖着的，单页漫画
					} else {
						//是双页漫画
						image.onload = null;	// 避免重复加载
						temp_flag = true;//宽大于高，很可能是开页
					}
				};
				return temp_flag;
			}
		},

		//翻页，其实不限页数
		flipPage: function (num) {
			if (
				this.nowPageNum_FlipMode + num <= this.book.all_page_num &&
				this.nowPageNum_FlipMode + num >= 1
			) {
				this.nowPageNum_FlipMode = this.nowPageNum_FlipMode + num;
			} else {
				// console.log("无法继续翻，Num:" + num)
				if (num > 0) {
					this.message.info(this.$t("hintLastPage"));
				} else {
					this.message.info(this.$t("hintFirstPage"));
				}
			}
			//保存页数
			this.saveNowPageNumOnUpdate(this.nowPageNum_FlipMode);
		},
		//拖动进度条,或翻页的时候保存页数
		saveNowPageNumOnUpdate(value) {
			if (this.saveNowPageNumFlag_FlipMode) {
				this.cookies.set("nowPageNum_FlipMode" + this.book.uuid, value);
			}
		},
		//跳转到指定页数
		toPage: function (num) {
			if (num <= this.book.all_page_num && num >= 1) {
				this.nowPageNum_FlipMode = num;
			}
			if (this.saveNowPageNumFlag_FlipMode) {
				this.cookies.set("nowPageNum_FlipMode" + this.book.uuid, this.nowPageNum_FlipMode);
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

		setShowHeaderChange(value) {
			console.log("value:" + value);
			this.showHeaderFlag_FlipMode = value;
			this.cookies.set("showHeaderFlag_FlipMode", value);
			console.log(
				"cookie设置完毕: showHeaderFlag_FlipMode=" +
				this.cookies.get("showHeaderFlag_FlipMode")
			);
		},
		setShowFooterFlagChange(value) {
			console.log("value:" + value);
			this.showFooterFlag_FlipMode = value;
			this.cookies.set("showFooterFlag_FlipMode", value);
			console.log(
				"cookie设置完毕: showFooterFlag_FlipMode=" +
				this.cookies.get("showFooterFlag_FlipMode")
			);
		},

		setShowPageNumChange(value) {
			console.log("value:" + value);
			this.showPageHintFlag_FlipMode = value;
			this.cookies.set("showPageHintFlag_FlipMode", value);
			console.log(
				"cookie设置完毕: showPageHintFlag_FlipMode=" +
				this.cookies.get("showPageHintFlag_FlipMode")
			);
		},

		setFlipScreenFlag(value) {
			console.log("value:" + value);
			this.rightScreenToNextFlag_FlipMode = value;
			this.cookies.set("rightScreenToNextFlag_FlipMode", value);
			console.log(
				"cookie设置完毕: rightScreenToNextFlag_FlipMode=" +
				this.cookies.get("rightScreenToNextFlag_FlipMode")
			);
		},

		setSavePageNumFlag(value) {
			console.log("value:" + value);
			this.saveNowPageNumFlag_FlipMode = value;
			this.cookies.set("saveNowPageNumFlag_FlipMode", value);
			console.log(
				"cookie设置完毕: saveNowPageNumFlag_FlipMode=" +
				this.cookies.get("saveNowPageNumFlag_FlipMode")
			);
		},

		setDebugModeFlag(value) {
			console.log("value:" + value);
			this.debugModeFlag = value;
			this.cookies.set("debugModeFlag", value);
			console.log("cookie设置完毕: debugModeFlag=" + this.cookies.get("debugModeFlag"));
		},

		setAutoDoublepage_FlipMode(value) {
			console.log("value:" + value);
			this.autoDoublepage_FlipMode = value;
			this.cookies.set("autoDoublepage_FlipMode", value);
			console.log("cookie设置完毕: autoDoublepage_FlipMode=" + this.cookies.get("autoDoublepage_FlipMode"));
		},



	},

	computed: {
		//页数或素描模式的提示
		pageNumOrSketchHint() {
			if (this.sketchModeFlag) {
				let nowSecond = (this.sketchSecondCount % this.sketchFlipSecond) + 1;
				let donePage = parseInt(this.sketchSecondCount / this.sketchFlipSecond);
				let totalMinutes = parseInt((this.sketchSecondCount + 1) / 60);
				//计算几小时几分
				let MinutesAndHourString = "";
				//如果不满意1小时，就不显示小时
				if (parseInt(totalMinutes / 60) == 0) {
					MinutesAndHourString = totalMinutes + this.$t("minute");
				} else {
					MinutesAndHourString =
						parseInt(totalMinutes / 60) +
						this.$t("hour") +
						(totalMinutes % 60) +
						this.$t("minute");
				}
				let AllTimeString =
					MinutesAndHourString + ((this.sketchSecondCount + 1) % 60) + this.$t("second");
				let hintString =
					this.$t("now_is") +
					nowSecond +
					this.$t("second") +
					"  " +
					this.$t("total_is") +
					donePage +
					this.$t("page") +
					"  " +
					this.$t("totalTime") +
					AllTimeString +
					"  " +
					this.$t("interval") +
					this.sketchFlipSecond +
					this.$t("second");
				return hintString;
			} else {
				return this.nowPageNum_FlipMode + "/" + this.book.all_page_num;
			}
		},
		mangaAreaHeight() {
			let Height = 95;
			//页头和底部拖动条都显示,或有一个显示的时候，95%
			if (this.showFooterFlag_FlipMode && this.showHeaderFlag_FlipMode) {
				Height = 95;
			}
			if (this.showFooterFlag_FlipMode && !this.showHeaderFlag_FlipMode) {
				Height = 95;
			}
			if (!this.showFooterFlag_FlipMode && this.showHeaderFlag_FlipMode) {
				Height = 95;
			}
			//页头和底部拖动条都不显示的时候，漫画占满屏幕
			if (!this.showFooterFlag_FlipMode && !this.showHeaderFlag_FlipMode) {
				Height = 100;
			}
			return Height + "vh";
		},
		mangaImageHeight() {
			let Height = 95;
			//页头和底部拖动条都显示,或有一个显示的时候，95%
			if (this.showFooterFlag_FlipMode && this.showHeaderFlag_FlipMode) {
				Height = 95;
			}
			if (this.showFooterFlag_FlipMode && !this.showHeaderFlag_FlipMode) {
				Height = 95;
			}
			if (!this.showFooterFlag_FlipMode && this.showHeaderFlag_FlipMode) {
				Height = 95;
			}
			//页头和拖动条都不显示的时候，漫画占满屏幕
			if (!this.showFooterFlag_FlipMode && !this.showHeaderFlag_FlipMode) {
				Height = 100;
			}
			//与上面唯一的不同，减去素描提示的空间
			if (this.showPageHintFlag_FlipMode) {
				if (this.nowTemplate == "sketch") {
					Height = Height - 6;
				} else {
					Height = Height - 3;
				}
			}
			return Height + "vh";
		},
		//进入素描模式的时候，把字体与高度放大一倍
		sketchHintHeight() {
			if (this.nowTemplate == "sketch") {
				return "6vh";
			} else {
				return "3vh";
			}
		},
		sketchHintFontSize() {
			if (this.nowTemplate == "sketch") {
				return "32px";
			} else {
				return "16px";
			}
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
	flex-direction: row;
	justify-content: center;
	align-items: baseline;
	user-select: none; /* 不可以被选中 */
	-moz-user-select: none; /* 火狐 */
	-webkit-user-select: none; /* 谷歌 */
}

/* 漫画div中的图片div*/
.manga_area_img_div {
	width: 100vw;
	height: v-bind(mangaAreaHeight);
	display: flex;
	flex-direction: row;
	justify-content: center;
	align-items: center;
	margin: 0px;
	border-radius: 3px;
	box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

/* 最后的一或两张图片*/
.manga_area_img_div img {
	/* max-height: inherit 继承 */
	/* max-height: inherit; */
	max-height: v-bind(mangaImageHeight);
	max-width: 100vw;
	margin: 0px; /* 两张图片之间不要留空间*/
	padding: 0px;
	background-color: #aaa;

	/* flex-grow: 100; */
}

/* 漫画div图片下面的页数*/
.sketch_hint {
	height: v-bind(sketchHintHeight);
	padding: 0px;
	text-align: center;
	font-size: v-bind(sketchHintFontSize);
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
