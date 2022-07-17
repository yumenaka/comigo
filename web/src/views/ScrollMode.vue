<template>
	<div id="ScrollMode" class="manga">
		<Header :setDownLoadLink="this.needDownloadLink()" :headerTitle="book.name" :bookID="this.book.id"
			:showReturnIcon="true" v-bind:style="{ background: model.interfaceColor }"
			@drawerActivate="this.drawerActivate">
		</Header>
		<!-- 顶部的加载全部页面顶部按钮 -->
		<button v-if="((this.startLoadPageNum > 1) && (!this.inLoading))"
			class="w-24 h-12 m-2 bg-blue-300 text-gray-900 hover:bg-blue-500 rounded" @click="this.loadAllPage"
			size="large">{{ $t('load_all_pages')}}</button>
		<!-- 渲染漫画部分 -->
		<div class="main_manga" v-for="(image, n) in this.localImages" :key="image.url" @click="onMouseClick($event)"
			@mousemove="onMouseMove" @mouseleave="onMouseLeave">
			<ImageScroll :image_url="this.imageParametersString(image.url)" :sPWL="this.sPWL" :dPWL="this.dPWL"
				:sPWP="this.sPWP" :dPWP="this.dPWP" :pageNum="n + this.startLoadPageNum"
				:all_page_num="this.book.all_page_num" :book_id="this.book.id"
				:showPageNumFlag_ScrollMode="this.showPageNumFlag_ScrollMode" :sendWSMessage="this.sendWSMessage"
				:syncPageFlag="this.syncPageFlag" :id="'JUMP_ID:' + (n + 1)" @refreshNowPageNum="refreshNowPageNum">
			</ImageScroll>
		</div>

		<Observer @loadNextBlock="loadNextBlock" />
		<!-- 返回顶部的圆形按钮，向上滑动的时候出现 -->
		<n-back-top :show="showBackTopFlag" type="info" color="#8a2be2" :right="20" :bottom="20" />
		<!-- 底部最下面的返回顶部按钮 -->
		<button class="w-24 h-12 m-2 bg-blue-300 text-gray-900 hover:bg-blue-500 rounded" @click="scrollToTop(90);"
			size="large">{{ $t('back-to-top') }}</button>

		<Bottom v-bind:style="{ background: model.interfaceColor }"
			:softVersion="this.$store.state.server_status.ServerName ? this.$store.state.server_status.ServerName : 'Comigo'">
		</Bottom>


		<Drawer :initDrawerActive="this.drawerActive" :initDrawerPlacement="this.drawerPlacement"
			@saveConfig="this.saveConfigToLocalStorage" @startSketch="this.startSketchMode"
			@closeDrawer="this.drawerDeactivate" @setT="this.OnSetTemplate" :readerMode="this.readerMode"
			:inBookShelf="false" :sketching="false">
			<!-- 选择：切换页面模式 -->
			<n-space>
				<n-button @click="changeReaderModeToFlipMode">{{ $t('switch_to_flip_mode') }}</n-button>
			</n-space>
			<p> &nbsp;</p>
			<!-- 同步页数 -->
			<n-space>
				<n-switch size="large" v-model:value="this.syncPageFlag" @update:value="this.setSyncPageFlag">
					<template #checked>{{ $t("sync_page") }}</template>
					<template #unchecked>{{ $t("sync_page") }}</template>
				</n-switch>
			</n-space>
			<!-- 保存页数 -->
			<n-space>
				<n-switch size="large" v-model:value="this.saveNowPageNumFlag" @update:value="this.setSavePageNumFlag">
					<template #checked>{{ $t("savePageNum") }}</template>
					<template #unchecked>{{ $t("savePageNum") }}</template>
				</n-switch>
			</n-space>
			<!-- 显示页数 -->
			<n-space>
				<n-switch size="large" v-model:value="this.showPageNumFlag_ScrollMode"
					@update:value="setShowPageNumChange">
					<template #checked>{{ $t('showPageNum') }}</template>
					<template #unchecked>{{ $t('showPageNum') }}</template>
				</n-switch>
			</n-space>

			<p> &nbsp;</p>
			<n-space vertical>
				<!-- 单页-漫画宽度-使用百分比 -->
				<!-- 数字输入% -->
				<n-input-number v-if="this.imageWidth_usePercentFlag" size="small" :show-button="false"
					v-model:value="this.singlePageWidth_Percent" :max="100" :min="10" :update-value-on-input="false">
					<template #prefix>{{ $t('singlePageWidth') }}</template>
					<template #suffix>%</template>
				</n-input-number>
				<!-- 滑动选择% -->
				<n-slider v-if="this.imageWidth_usePercentFlag" v-model:value="this.singlePageWidth_Percent" :step="1"
					:max="100" :min="10" :format-tooltip="value => `${value}%`" />

				<!-- 开页-漫画宽度-使用百分比  -->
				<!-- 数字输入% -->
				<n-input-number v-if="this.imageWidth_usePercentFlag" size="small" :show-button="false"
					v-model:value="this.doublePageWidth_Percent" :max="100" :min="10" :update-value-on-input="false">
					<template #prefix>{{ $t('doublePageWidth') }}</template>
					<template #suffix>%</template>
				</n-input-number>
				<!-- 滑动选择% -->
				<n-slider v-if="this.imageWidth_usePercentFlag" v-model:value="this.doublePageWidth_Percent" :step="1"
					:max="100" :min="10" :format-tooltip="value => `${value}%`" />

				<!-- 单页-漫画宽度-使用固定值PX -->
				<!-- 数字输入框 -->
				<n-input-number v-if="!this.imageWidth_usePercentFlag" size="small" :show-button="false"
					v-model:value="this.singlePageWidth_PX" :min="50" :update-value-on-input="false">
					<template #prefix>{{ $t('singlePageWidth') }}</template>
					<template #suffix>px</template>
				</n-input-number>
				<!-- 滑动选择PX -->
				<n-slider v-if="!this.imageWidth_usePercentFlag" v-model:value="this.singlePageWidth_PX" :step="10"
					:max="1600" :min="50" :format-tooltip="value => `${value}px`" />

				<!-- 数字输入框 -->
				<n-input-number v-if="!this.imageWidth_usePercentFlag" size="small" :show-button="false"
					v-model:value="this.doublePageWidth_PX" :min="50" :update-value-on-input="false">
					<template #prefix>{{ $t('doublePageWidth') }}</template>
					<template #suffix>px</template>
				</n-input-number>

				<!-- 滑动选择PX -->
				<n-slider v-if="!this.imageWidth_usePercentFlag" v-model:value="this.doublePageWidth_PX" :step="10"
					:max="1600" :min="50" :format-tooltip="value => `${value}px`" />
			</n-space>

			<p></p>
			<!-- 开关：横屏状态下,宽度单位是百分比还是固定值 -->
			<n-space>
				<n-switch size="large" v-model:value="this.imageWidth_usePercentFlag" :rail-style="railStyle"
					@update:value="this.setImageWidthUsePercentFlag">
					<template #checked>{{ $t('width_usePercent') }}</template>
					<template #unchecked>{{ $t('width_useFixedValue') }}</template>
				</n-switch>
			</n-space>

			<p> &nbsp;</p>

			<!-- 开关：是否显示页头
			<n-space>
				<n-switch size="large" v-model:value="this.showHeaderFlag" @update:value="setShowHeaderChange">
					<template #checked>{{ $t('showHeader') }}</template>
					<template #unchecked>{{ $t('showHeader') }}</template>
				</n-switch>
				<p></p>
			</n-space> -->

			<!-- 开关：自动切边 -->
			<n-space>
				<n-switch size="large" v-model:value="this.imageParameters.do_auto_crop"
					@update:value="setImageParameters_DoAutoCrop">
					<template #checked>{{ $t('auto_crop') }}</template>
					<template #unchecked>{{ $t('auto_crop') }}</template>
				</n-switch>
				<!-- 切白边阈值 -->
				<n-input-number :show-button="false" v-if="this.imageParameters.do_auto_crop"
					v-model:value="this.imageParameters.auto_crop_num" @update:value="setImageParameters_AutoCropNum"
					:max="10" :min="0">
					<template #prefix>{{ $t('energy_threshold') }}</template>
				</n-input-number>
			</n-space>
			<!-- 开关：压缩图片 -->
			<n-space>
				<n-switch size="large" :rail-style="railStyle" v-model:value="this.imageParameters.do_auto_resize"
					@update:value="setImageParameters_DoAutoResize">
					<template #checked>{{ $t('image_width_limit') }}</template>
					<template #unchecked>{{ $t('raw_resolution') }}</template>
				</n-switch>
				<!-- 压缩图片参数：数字输入框 -->
				<n-input-number v-if="this.imageParameters.do_auto_resize" size="small" :show-button="false"
					v-model:value="this.imageParameters.resize_max_width" :min="100">
					<template #prefix>{{ $t('max_width') }}</template>
					<template #suffix>px</template>
				</n-input-number>
			</n-space>

			<!-- 开关：显示原图 黑白 -->
			<!-- <n-space>
				<n-switch size="large" v-model:value="this.imageParameters.gray"
					@update:value="setImageParameters_Gray">
					<template #checked>{{ $t('gray_image') }}</template>
					<template #unchecked>{{ $t('gray_image') }}</template>
				</n-switch>
			</n-space> -->

			<!-- 分割线 -->
			<!-- <n-divider /> -->
			<p> &nbsp;</p>
			<n-dropdown trigger="hover" :options="options" @select="onResort">
				<n-button>{{ this.getSortHintText(this.resort_hint_key) }}</n-button>
			</n-dropdown>
		</Drawer>
	</div>
</template>

<script>
// 直接导入组件并使用它。这种情况下,只有导入的组件才会被打包。
import { NBackTop, NSpace, NSlider, NSwitch, NInputNumber, NButton, NDropdown, useMessage, useDialog } from 'naive-ui'
import Header from "@/components/Header.vue";
import Drawer from "@/components/Drawer.vue";
import Bottom from "@/components/Bottom.vue";
import Observer from "@/components/Observer_in_Scroll.vue";
import ImageScroll from "@/components/Image_in_Scroll.vue";
import { defineComponent, reactive } from 'vue'
// import { useCookies } from "vue3-cookies";// https://github.com/KanHarI/vue3-cookies
import axios from "axios";

export default defineComponent({
	name: "ScrollMode",
	props: ['test_prop'],
	emits: ["setTemplate"],
	components: {
		Header,//自定义页头
		Drawer,//自定义抽屉
		Bottom,//自定义页尾
		Observer,//Observer组件,下拉刷新用
		ImageScroll,//漫画页，包含Observer组，获取当前页用
		NBackTop,//回到顶部按钮,来自:https://www.naiveui.com/zh-CN/os-theme/components/back-top
		// NDrawer,//抽屉,可以从上下左右4个方向冒出. https://www.naiveui.com/zh-CN/os-theme/components/drawer
		// NDrawerContent,//抽屉内容
		NSpace,//间距 https://www.naiveui.com/zh-CN/os-theme/components/space
		// NRadio,//单选  https://www.naiveui.com/zh-CN/os-theme/components/radio
		// NRadioButton,//单选  用按钮显得更优雅一点
		// NRadioGroup,
		NSlider,//滑动选择  Slider https://www.naiveui.com/zh-CN/os-theme/components/slider
		NSwitch,//开关   https://www.naiveui.com/zh-CN/os-theme/components/switch
		// NLayout,//布局 https://www.naiveui.com/zh-CN/os-theme/components/layout
		// NLayoutSider,
		// NLayoutContent,
		// NPageHeader,//页头 https://www.naiveui.com/zh-CN/os-theme/components/page-header
		// NAvatar, //头像 https://www.naiveui.com/zh-CN/os-theme/components/avatar
		NInputNumber,//数字输入 https://www.naiveui.com/zh-CN/os-theme/components/input-number
		// NDivider,//分割线  https://www.naiveui.com/zh-CN/os-theme/components/divider
		// NColorPicker,
		NButton,//按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/button
		NDropdown,//下拉菜单 https://www.naiveui.com/zh-CN/os-theme/components/dropdown
	},
	// setup在创建组件前执行，因此没有this
	setup() {
		//此处不能使用this,但可以用getCurrentInstance 这个vue函数取得Proxy，实现类似 proxy.$socket.onmessage 这样的调用(https://github.com/likaia/vue-native-websocket-vue3)。
		// const { cookies } = useCookies();
		//在setup中访问vuex需要通过useStore()来访问  https://juejin.cn/post/6917592199140458504#heading-22=
		//背景颜色,颜色选择器用
		//reactive({}) 创建并返回一个响应式对象: https://www.bilibili.com/video/av925511720/?p=4  也讲到了toRefs()
		const model = reactive({
			backgroundColor: "#E0D9CD",
			interfaceColor: "#f5f5e4",
		});

		//请求图片文件时，可添加的额外参数
		const imageParameters = reactive({
			resize_width: -1,// 缩放图片,指定宽度
			resize_height: -1,// 指定高度,缩放图片
			do_auto_resize: false,
			resize_max_width: 800,//图片宽度大于这个上限时缩小 
			resize_max_height: -1,//图片高度大于这个上限时缩小
			do_auto_crop: false,
			auto_crop_num: 1,// 自动切白边阈值,范围是0~100,其实为1就够了	
			gray: false,//黑白化
		});
		//单选按钮绑定的数值,ref函数：返回一个响应式的引用
		// const checkedValueRef = ref(null)
		const message = useMessage();
		const dialog = useDialog();
		return {
			pdfUrl: "",
			// cookies,
			//背景色
			model,
			message,
			dialog,
			imageParameters,//获取图片所用参数
			imageParametersString: (source_url) => {
				// console.log("source_url:" + source_url)
				if (source_url.substr(0, 12) == "api/getfile?") {
					//当前URL
					const url = document.location.toString();
					//按照“/”分割字符串
					const arrUrl = url.split("/");
					//拼一个完整的图片URL（因为路由路径会变化,所以不能用相对路径？）
					const base_str = arrUrl[0] + "//" + arrUrl[2] + "/" + source_url;
					//添加各种字符串参数,不需要的话为空
					const resize_width_str = (imageParameters.resize_width > 0 ? "&resize_width=" + imageParameters.resize_width : "");
					const resize_height_str = (imageParameters.resize_height > 0 ? "&resize_height=" + imageParameters.resize_height : "");
					const gray_str = (imageParameters.gray ? "&gray=true" : "");
					const do_auto_resize_str = (imageParameters.do_auto_resize ? ("&resize_max_width=" + imageParameters.resize_max_width) : "");
					const resize_max_height_str = (imageParameters.resize_max_height > 0 ? "&resize_max_height=" + imageParameters.resize_max_height : "");
					const auto_crop_str = (imageParameters.do_auto_crop ? "&auto_crop=" + imageParameters.auto_crop_num : "");
					//所有附加的转换参数
					const addStr = resize_width_str + resize_height_str + do_auto_resize_str + resize_max_height_str + auto_crop_str + gray_str;
					//如果有附加转换参数，则设置成不缓存
					const nocache_str = (addStr === "" ? "" : "&no-cache=true");
					const full_url = base_str + addStr + nocache_str;
					// console.log(full_url);
					return full_url;
				} else {
					return source_url
				}
			},
			//开关的颜色
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
			//当前页数,注意语义,直接就是1开始的页数,不是数组下标,在pages\images数组当中用的时候需要-1
			nowPageNum: 1,
			startLoadPageNum: 1,
			StartFromBreakpoint: false,
			loadPageLimit: 20,//一次最多载入的漫画张数，默认为20.
			endLoadPageNum: 20,//载入漫画的最后一页，默认为20.
			syncPageFlag: true,//是否通过websocket同步翻页
			saveNowPageNumFlag: true,//是否（在本地存储里面）保存与恢复页数
			sendWSMessage: false,
			firstloadComplete: true,
			localImages: null,
			inLoading: true,//是否正在加载，延迟执行操作、隐藏按钮用
			book: {
				name: "loading",
				id: "abcde",
				all_page_num: 2,
				book_type: "dir",
				pages: {
					sort_by: "",
					images: [
						{
							height: 500,
							width: 449,
							url: "/images/loading.gif",
						},
						{
							height: 500,
							width: 449,
							url: "/images/loading.gif",
						},
					],
				}
			},
			resort_hint_key: "resort",
			options: [
				{
					label: this.$t('sort_by_filename'),
					key: "filename",
				},
				{
					label: this.$t('sort_by_modify_time'),
					key: "modify_time"
				},
				{
					label: this.$t('sort_by_filesize'),
					key: "filesize"
				},
				{
					label: this.$t('sort_by_filename') + this.$t('sort_reverse'),
					key: "filename_reverse",
				},
				{
					label: this.$t('sort_by_modify_time') + this.$t('sort_reverse'),
					key: "modify_time_reverse"
				},
				{
					label: this.$t('sort_by_filesize') + this.$t('sort_reverse'),
					key: "filesize_reverse"
				},
			],
			readerMode: "scroll",
			drawerActive: false,
			drawerPlacement: 'right',
			//开发模式 还没有做的功能与设置,设置Debug以后才能见到
			debugModeFlag: true,
			//书籍数据,需要从远程拉取
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
			//同步滚动,目前还没做
			syncScrollFlag: false,
			//鼠标点击或触摸的位置
			clickX: 0,
			clickY: 0,
			//可见范围是否是横向
			isLandscapeMode: true,
			isPortraitMode: false,
			imageMaxWidth: 10,
			//屏幕宽横比,inLandscapeMode的判断依据
			aspectRatio: 1.2,
			//状态驱动的动态 CSS
			// https://v3.cn.vuejs.org/api/sfc-style.html#%E7%8A%B6%E6%80%81%E9%A9%B1%E5%8A%A8%E7%9A%84%E5%8A%A8%E6%80%81-css
			//图片宽度的单位,是否使用百分比
			imageWidth_usePercentFlag: false,
			//横屏(Landscape)状态的漫画页宽度,百分比
			singlePageWidth_Percent: 50,
			doublePageWidth_Percent: 95,
			//横屏(Landscape)状态的漫画页宽度,PX
			singlePageWidth_PX: 720,
			doublePageWidth_PX: 1080,
			//可见范围宽高的具体值
			clientWidth: 0,
			clientHeight: 0,
		};
	},
	//Vue3生命周期:  https://v3.cn.vuejs.org/api/options-lifecycle-hooks.html#beforecreate
	// created : 在绑定元素的属性或事件监听器被应用之前调用。
	// beforeMount : 指令第一次绑定到元素并且在挂载父组件之前调用。
	// mounted : 在绑定元素的父组件被挂载后调用。
	// beforeUpdate: 在更新包含组件的 VNode 之前调用。。
	// updated: 在包含组件的 VNode 及其子组件的 VNode 更新后调用。
	// beforeUnmount: 当指令与在绑定元素父组件卸载之前时,只调用一次。
	// unmounted: 当指令与元素解除绑定且父组件已卸载时,只调用一次。
	created() {
		// 消息监听，即接收websocket服务端推送的消息. optionsAPI用法
		this.$options.sockets.onmessage = (data) => this.handlePacket(data);
		//根据文件名、修改时间、文件大小等要素排序的参数
		let sort_image_by_str = "";
		if (this.$route.query.sort_by) {
			sort_image_by_str = "&sort_by=" + this.$route.query.sort_by
		}
		//是否保存页数
		if (localStorage.getItem("saveNowPageNumFlag") === "true") {
			this.saveNowPageNumFlag = true;
		} else if (localStorage.getItem("saveNowPageNumFlag") === "false") {
			this.saveNowPageNumFlag = false;
		}
		//是否通过websocket同步页数
		if (localStorage.getItem("SyncPageFlag") === "true") {
			this.syncPageFlag = true;
		} else if (localStorage.getItem("SyncPageFlag") === "false") {
			this.syncPageFlag = false;
		}

		//根据路由参数获取特定书籍
		axios
			.get("/getbook?id=" + this.$route.params.id + sort_image_by_str)
			.then((response) => {
				//请求接口成功的逻辑
				this.book = response.data;
				this.loadPages();
			})
			.finally(
				() => {
					document.title = this.book.name;
					// console.log("成功获取书籍数据,书籍ID:" + this.$route.params.id);
					// 询问用户是否从中间开始加载，延迟1.5秒执行
					var _this = this
					setTimeout(function () {
						_this.loadLocalBookMark();
						_this.inLoading = false
					}, 1500);
				}
			);
		//监听路由参数的变化,刷新本地的Book数据
		this.$watch(
			() => this.$route.params.id,
			(id) => {
				if (id) {
					axios
						.get("/getbook?id=" + this.$route.params.id + sort_image_by_str)
						.then((response) => (this.book = response.data))
						.finally(console.log("路由参数改变,书籍ID:" + id));
				}
			}
		);

		window.addEventListener("scroll", this.onScroll);
		//文档视图调整大小时会触发 resize 事件。 https://developer.mozilla.org/zh-CN/docs/Web/API/Window/resize_event
		window.addEventListener("resize", this.onResize);
		this.imageMaxWidth = window.innerWidth;
		//根据本地存储初始化默认值,读取出来的是字符串,不要直接用

		//是否通过websocket同步页数
		if (localStorage.getItem("SyncPageFlag") === "true") {
			this.syncPageFlag = true;
		} else if (localStorage.getItem("SyncPageFlag") === "false") {
			this.syncPageFlag = false;
		}

		//是否显示页头
		if (localStorage.getItem("showHeaderFlag") === "true") {
			this.showHeaderFlag = true;
		} else if (localStorage.getItem("showHeaderFlag") === "false") {
			this.showHeaderFlag = false;
		}
		//console.log("读取设置并初始化: showHeaderFlag=" + this.showHeaderFlag);

		//是否显示页数
		if (localStorage.getItem("showPageNumFlag_ScrollMode") === "true") {
			this.showPageNumFlag_ScrollMode = true;
		} else if (localStorage.getItem("showPageNumFlag_ScrollMode") === "false") {
			this.showPageNumFlag_ScrollMode = false;
		}
		//console.log("读取设置并初始化: showPageNumFlag_ScrollMode=" + this.showPageNumFlag_ScrollMode);
		//javascript 数字类型转换：https://www.runoob.com/js/js-type-conversion.html
		// NaN不能通过相等操作符（== 和 ===）来判断

		//漫画页宽度,Percent
		if (localStorage.getItem("singlePageWidth_Percent") != null) {
			let saveNum = Number(localStorage.getItem("singlePageWidth_Percent"));
			if (!isNaN(saveNum)) {
				this.singlePageWidth_Percent = saveNum;
			}
		}

		if (localStorage.getItem("doublePageWidth_Percent") != null) {
			let saveNum = Number(localStorage.getItem("doublePageWidth_Percent"));
			if (!isNaN(saveNum)) {
				this.doublePageWidth_Percent = saveNum;
			}
		}

		//漫画页宽度,PX
		if (localStorage.getItem("singlePageWidth_PX") != null) {
			let saveNum = Number(localStorage.getItem("singlePageWidth_PX"));
			if (!isNaN(saveNum)) {
				this.singlePageWidth_PX = saveNum;
			}
		}
		if (localStorage.getItem("doublePageWidth_PX") != null) {
			let saveNum = Number(localStorage.getItem("doublePageWidth_PX"));
			if (!isNaN(saveNum)) {
				this.doublePageWidth_PX = saveNum;
			}
		}
		//当前颜色
		if (localStorage.getItem("BackgroundColor") != null) {
			this.model.backgroundColor = localStorage.getItem("BackgroundColor");
		}
		if (localStorage.getItem("InterfaceColor") != null) {
			this.model.interfaceColor = localStorage.getItem("InterfaceColor");
		}

		//宽度是否使用百分比
		if (localStorage.getItem("imageWidth_usePercentFlag") === "true") {
			this.imageWidth_usePercentFlag = true;
		} else if (localStorage.getItem("imageWidth_usePercentFlag") === "false") {
			this.imageWidth_usePercentFlag = false;
		}

		// 图片处理相关
		//是否获取黑白图片
		if (localStorage.getItem("ImageParameters_Gray") === "true") {
			this.imageParameters.gray = true;
		} else if (localStorage.getItem("ImageParameters_Gray") === "false") {
			this.imageParameters.gray = false;
		}
		// console.log("读取设置并初始化: ImageParameters_Gray=" + this.imageParameters.gray);

		//是否压缩图片
		if (localStorage.getItem("ImageParameters_DoAutoResize") === "true") {
			this.imageParameters.do_auto_resize = true;
		} else if (localStorage.getItem("ImageParameters_DoAutoResize") === "false") {
			this.imageParameters.do_auto_resize = false;
		}

		//启用压缩的Width下限
		if (localStorage.getItem("ImageParametersResizeMaxWidth") != null) {
			let saveNum = Number(localStorage.getItem("ImageParametersResizeMaxWidth"));
			if (!isNaN(saveNum)) {
				this.imageParameters.resize_max_width = saveNum;
			}
		}

		//是否自动切白边
		if (localStorage.getItem("ImageParameters_DoAutoCrop") === "true") {
			this.imageParameters.do_auto_crop = true;
		} else if (localStorage.getItem("ImageParameters_DoAutoCrop") === "false") {
			this.imageParameters.do_auto_crop = false;
		}

		//切白边参数
		if (localStorage.getItem("ImageParameters_AutoCropNum") != null) {
			let saveNum = Number(localStorage.getItem("ImageParameters_AutoCropNum"));
			if (!isNaN(saveNum)) {
				this.imageParameters.auto_crop_num = saveNum;
			}
		}
	},
	// 挂载前 : 指令第一次绑定到元素并且在挂载父组件之前调用。
	beforeMount() {
	},
	onMounted() {
		//console.log('mounted in the composition api!')
		this.isLandscapeMode = this.inLandscapeModeCheck();
		this.isPortraitMode = !this.inLandscapeModeCheck();
		// https://v3.cn.vuejs.org/api/options-lifecycle-hooks.html#beforemount
		this.$nextTick(function () {
			//视图渲染之后运行的代码
		})

	},
	//卸载前
	beforeUnmount() {
		//组件销毁前,销毁监听事件
		window.removeEventListener("scroll", this.onScroll);
		window.removeEventListener('resize', this.onResize);
		// 组件被销毁时清理观察者,不做的话会造成内存泄漏，因为事件监听器不会被清除。
		// this.observer.disconnect();//停止观察所有元素
	},
	methods: {
		//接收服务器发来的websocket消息，做各种反应（翻页、提示信息）
		handlePacket(data) {
			if (!this.syncPageFlag) {
				return;
			}
			const msg = JSON.parse(data.data);
			//心跳信息,直接返回
			if (msg.type === "heartbeat") {
				return;
			}
			//服务器发来翻页信息，来自于另一个用户，
			if (msg.type === "sync_page" && msg.user_id !== this.$store.userID) {
				const syncData = JSON.parse(msg.data_string);
				//正在读的是同一本书、就翻页。
				if (syncData.book_id === this.book.id && syncData.now_page_num !== this.nowPageNum) {
					// console.log(syncData);
					// this.toPage(syncData.now_page_num, false);
				}
			}
		},

		loadBookMarkDialog() {
			this.dialog.warning({
				title: this.$t('found_read_history'),
				content: this.$t('load_from_interrupt').replace("XX", this.startLoadPageNum),
				positiveText: this.$t('from_interrupt'),
				negativeText: this.$t('starting_from_beginning'),
				onPositiveClick: () => {
					if (this.startLoadPageNum + this.loadPageLimit <= this.book.all_page_num) {
						this.endLoadPageNum = this.startLoadPageNum + this.loadPageLimit;
					} else {
						this.endLoadPageNum = this.book.all_page_num;
					}
					this.message.success(this.$t('successfully_loaded_reading_progress'));
					this.loadPages();
				},
				onNegativeClick: () => {
					this.startLoadPageNum = 1;
					this.nowPageNum = 1;
					this.message.success(this.$t('starting_from_beginning_hint'));
					this.loadPages();
				}
			});
		},

		//刷新到底部的时候,改变images数据 默认参数： NowBlockNumPlusOne = false
		loadPages() {
			// const MaxPageNum = this.book.all_page_num
			// const LoadPageLimit = this.loadPageLimit
			// var NowBlockNum = Math.ceil(this.nowPageNum / LoadPageLimit)//现在在哪个区块（向上取整，有小数，则整数部分加1）取整：parseInt()
			// const AllBlockNum = Math.ceil(MaxPageNum/LoadPageLimit)//总区块数（向上取整，有小数，则整数部分加1）

			// console.log("startLoadPageNum:", this.startLoadPageNum)
			// console.log("endLoadPageNum:", this.endLoadPageNum)
			// console.dir(this.localImages)
			//slice() 方法返回一个新的数组对象，这一对象是一个由 begin 和 end 决定的原数组的浅拷贝（包括 begin，不包括end）
			this.localImages = this.book.pages.images.slice(this.startLoadPageNum - 1, this.endLoadPageNum);//javascript的接片不能直接用[a,b]，而是需要调用.slice()函数
		},
		loadAllPage() {
			this.startLoadPageNum = 1;
			this.nowPageNum = 1;
			this.loadPages()
		},
		//无限加载用,底部刷新
		loadNextBlock() {
			if (this.endLoadPageNum + this.loadPageLimit <= this.book.all_page_num) {
				this.endLoadPageNum = this.endLoadPageNum + this.loadPageLimit;
			} else {
				this.endLoadPageNum = this.book.all_page_num;
			}
			// console.log("loadNextBlock");
			this.loadPages(true);
		},
		//监听子组件事件: https://v3.cn.vuejs.org/guide/component-basics.html#%E7%9B%91%E5%90%AC%E5%AD%90%E7%BB%84%E4%BB%B6%E4%BA%8B%E4%BB%B6
		//滚动页面的时候刷新页数
		refreshNowPageNum(n) {
			if (this.inLoading) {
				return
			}
			this.nowPageNum = n;
			//保存页数
			this.saveLocalBookMark(this.nowPageNum, false);
			this.loadPages();
		},
		//滚动跳转到指定页数，需要写一个滚动函数
		toPage: function (num, sendWSMessage = true) {
			// this.sendWSMessage = sendWSMessage;
			if (num <= this.book.all_page_num && num >= 1) {
				this.nowPageNum = num;
			}
			this.loadPages();
			// 前端页面滚动到某个位置的方式
			//获取目标元素
			let element = document.getElementById("JUMP_ID:" + this.nowPageNum);
			if (!element) {
				console.log("没找到：", element);
				this.loadNextBlock();
				element = document.getElementById("JUMP_ID:" + this.nowPageNum);
			}
			//元素方法调用,自动刷新的时候，需要避免死锁,将syncPageFlag暂时设置为false不管用，
			// https://blog.csdn.net/weixin_41804429/article/details/102954146
			// https://developer.mozilla.org/zh-CN/docs/Web/API/Element/scrollIntoView

			if (element) {
				element.scrollIntoView(true)
			}
			//保存页数
			this.saveLocalBookMark(this.nowPageNum, sendWSMessage);
			// let _this = this
			// setTimeout(() => { _this.sendWSMessage = true; }, 2000);
		},
		//滑动页面、停止滚动的时候保存页数
		saveLocalBookMark(value, sendWSMessage = true) {
			if (this.saveNowPageNumFlag) {
				localStorage.setItem("nowPageNum" + this.book.id, value);
			}
			//发送翻页消息到服务器
			if (sendWSMessage && this.syncPageFlag) {
				this.sendNowPage();
			}
		},

		//根据书籍UUID,设定当前页数,显示的时候需要远程书籍数据（this.book）,可能需要延迟执行
		loadLocalBookMark() {
			if (!this.saveNowPageNumFlag) {
				return
			}
			let cookieValue = localStorage.getItem("nowPageNum" + this.book.id);
			if (cookieValue === null) {
				console.log("本队存储里没找到:" + "nowPageNum = " + this.nowPageNum);
				return
			}
			let saveNum = Number(cookieValue);
			if (isNaN(saveNum)) {
				console.log("读取页数失败,this.nowPageNum = " + this.nowPageNum);
				return
			}
			this.nowPageNum = saveNum;
			this.startLoadPageNum = saveNum;
			//至少读到第三页才开始提醒中途加载
			if (this.nowPageNum >= 3) {
				this.loadBookMarkDialog();
			}
			console.log("成功读取页数" + saveNum);
		},
		// 各种翻页函数，或许不需要都写
		flipPage() {
		},
		toNextPage() {
		},
		toPerviousPage() {
		},
		// 键盘事件，需要改造成支持滚动的逻辑
		handleKeyup(event) {
			//错误:(815, 49) 不允许从实参中引用 'caller' 和 'callee'
			const e = event || window.event;
			if (!e) return;
			//https://developer.mozilla.org/zh-CN/docs/Web/API/KeyboardEvent/keyCode
			switch (e.key) {
				case "ArrowUp":
				case "PageUp":
					this.flipPage(-1); //上一页
					break;
				case "ArrowLeft":
					this.toNextPage();
					break;
				case "ArrowRight":
					this.toPerviousPage();
					break;
				case "Space":
				case "ArrowDown":
				case "PageDown":
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
					break;
			}
			// console.log(e.keyCode);
			// console.log(e.key);
		},

		//页面排序相关
		onResort(key) {
			axios
				.get("/getbook?id=" + this.$route.params.id + "&sort_by=" + key)
				.then((response) => (this.book = response.data))
				.finally(
					() => {
						document.title = this.book.name;
						this.resort_hint_key = key;
						this.loadPages();
						// 带查询参数，结果是 /#/scroll/abc123?sort_by="filesize"
						this.$router.push({ name: "ScrollMode", replace: true, query: { sort_by: key } });
						console.log("成功刷新书籍数据,书籍ID:" + this.$route.params.id + "  sort_by=" + key);
					}
				);
		},

		//返回“重新排序”选择菜单的文字提示
		getSortHintText(key) {
			switch (key) {
				case "filename": return this.$t('sort_by_filename');
				case "modify_time": return this.$t('sort_by_modify_time');
				case "filesize": return this.$t('sort_by_filesize');
				case "filename_reverse": return this.$t('sort_by_filename') + this.$t('sort_reverse');
				case "modify_time_reverse": return this.$t('sort_by_modify_time') + this.$t('sort_reverse');
				case "filesize_reverse": return this.$t('sort_by_filesize') + this.$t('sort_reverse');
				default:
					return this.$t('re_sort');
			}
		},

		//切换到翻页模式
		changeReaderModeToFlipMode() {
			localStorage.setItem("ReaderMode", "flip");
			//replace的作用类似于 router.push，唯一不同的是，它在导航时不会向 history 添加新记录，正如它的名字所暗示的那样——它取代了当前的条目。
			this.$router.push({ name: "FlipMode", replace: true, params: { id: this.$route.params.id } });
		},
		needDownloadLink() {
			return this.book.book_type !== "dir"
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
		//开始素描模式
		startSketchMode() {
			localStorage.setItem("ReaderMode", "sketch");
			this.$router.push({ name: "FlipMode", replace: true, params: { id: this.$route.params.id } });
		},
		//接收Draw的参数,继续往父组件传
		OnSetTemplate(value) {
			this.$emit("setTemplate", value);
		},
		onBackgroundColorChange(value) {
			this.model.backgroundColor = value
			localStorage.setItem("BackgroundColor", value);
		},
		//关闭抽屉的时候保存配置 
		//如果在一个组件上使用了 v-model:xxx,应该使用 @update:xxx  https://www.naiveui.com/zh-CN/os-theme/docs/common-issues
		saveConfigToLocalStorage() {
			// 储存配置
			localStorage.setItem("nowPageNum" + this.book.id, this.nowPageNum);
			localStorage.setItem("SyncPageFlag", this.syncPageFlag);
			localStorage.setItem("showHeaderFlag", this.showHeaderFlag);
			localStorage.setItem("showPageNumFlag_ScrollMode", this.showPageNumFlag_ScrollMode);
			localStorage.setItem("imageWidth_usePercentFlag", this.imageWidth_usePercentFlag);
			localStorage.setItem("singlePageWidth_Percent", this.singlePageWidth_Percent);
			localStorage.setItem("doublePageWidth_Percent", this.doublePageWidth_Percent);
			localStorage.setItem("singlePageWidth_PX", this.singlePageWidth_PX);
			localStorage.setItem("doublePageWidth_PX", this.doublePageWidth_PX);
			localStorage.setItem("BackgroundColor", this.model.backgroundColor);
			//set对有setXXXChange函数的来说有些多余,但没有set函数的话就有必要了
			localStorage.setItem("ImageParameters_DoAutoCrop", this.imageParameters.do_auto_crop);
			localStorage.setItem("ImageParametersResizeMaxWidth", this.imageParameters.resize_max_width);
		},
		setSyncPageFlag(value) {
			console.log("value:" + value);
			this.syncPageFlag = value;
			localStorage.setItem("SyncPageFlag", value);
			console.log(
				"cookie设置完毕: SyncPageFlag=" +
				localStorage.getItem("SyncPageFlag")
			);
		},
		setShowHeaderChange(value) {
			// console.log("value:" + value);
			this.showHeaderFlag = value;
			localStorage.setItem("showHeaderFlag", value);
			// console.log("成功保存设置: showHeaderFlag=" + localStorage.getItem("showHeaderFlag"));
		},
		setShowPageNumChange(value) {
			// console.log("value:" + value);
			this.showPageNumFlag_ScrollMode = value;
			localStorage.setItem("showPageNumFlag_ScrollMode", value);
			// console.log("成功保存设置: showPageNumFlag_ScrollMode=" + localStorage.getItem("showPageNumFlag_ScrollMode"));
		},
		//图片处理相关
		//黑白化参数
		setImageParameters_Gray(value) {
			// console.log("value:" + value);
			this.imageParameters.gray = value;
			localStorage.setItem("ImageParameters_Gray", value);
			// console.log("成功保存设置: ImageParameters_Gray=" + localStorage.getItem("ImageParameters_Gray"));
		},
		//缩放图片大小的参数
		setImageParameters_DoAutoResize(value) {
			this.imageParameters.do_auto_resize = value;
			localStorage.setItem("ImageParameters_DoAutoResize", value);
			// console.log("成功保存设置: ImageParameters_DoAutoResize=" + localStorage.getItem("ImageParameters_DoAutoResize"));
		},
		//设置是否切白边
		setImageParameters_DoAutoCrop(value) {
			this.imageParameters.do_auto_crop = value;
			localStorage.setItem("ImageParameters_DoAutoCrop", this.imageParameters.do_auto_crop);
			// console.log("成功保存设置: ImageParameters_DoAutoCrop=" + localStorage.getItem("ImageParameters_DoAutoCrop"));
		},
		//切白边参数
		setImageParameters_AutoCropNum(value) {
			this.imageParameters.auto_crop_num = value;
			localStorage.setItem("ImageParameters_AutoCropNum", this.imageParameters.auto_crop_num);
		},

		setImageWidthUsePercentFlag(value) {
			console.log("value:" + value);
			this.imageWidth_usePercentFlag = value;
			localStorage.setItem("imageWidth_usePercentFlag", value);
			// console.log("成功保存设置: imageWidth_usePercentFlag=" + this.imageWidth_usePercentFlag);
		},

		setSavePageNumFlag(value) {
			console.log("value:" + value);
			this.saveNowPageNumFlag = value;
			localStorage.setItem("saveNowPageNumFlag", value);
			console.log(
				"cookie设置完毕: saveNowPageNumFlag=" +
				localStorage.getItem("saveNowPageNumFlag")
			);
		},

		setDebugModeFlag(value) {
			console.log("value:" + value);
			this.debugModeFlag = value;
			//关闭Debug模式的时候顺便也关上“自动合并单双页”的功能（因为还有BUG）
			if (value === false) {
				this.autoDoublePageModeFlag = false;
			}
			localStorage.setItem("debugModeFlag", value);
			console.log(
				"cookie设置完毕: debugModeFlag=" + localStorage.getItem("debugModeFlag")
			);
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
			// 为了调试的时候方便,阈值是正方形
			if (this.aspectRatio > (19 / 19)) {
				this.isLandscapeMode = true
				this.isPortraitMode = false
			} else {
				this.isLandscapeMode = false
				this.isPortraitMode = true
			}
		},
		//页面滚动的时候,改变返回顶部按钮的显隐
		onScroll() {
			let scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
			this.scrollDownFlag = scrollTop > this.scrollTopSave;
			//防手抖,小于一定数值状态就不变 Math.abs()会导致报错
			// let step = Math.abs(this.scrollTopSave - scrollTop);
			let step = this.scrollTopSave - scrollTop;
			// console.log("this.scrollDownFlag:",this.scrollDownFlag,"scrollTop:",scrollTop,"step:", step);
			this.scrollTopSave = scrollTop
			// this.showBackTopFlag=true
			if (step < -5 || step > 5) {
				this.showBackTopFlag = ((scrollTop > 400) && !this.scrollDownFlag);
			}
		},
		//获取鼠标位置,决定是否打开设置面板
		onMouseClick(e) {
			this.clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
			this.clickY = e.y //获取鼠标的Y坐标（鼠标与屏幕顶部的距离,单位为px）
			//浏览器的视口,不包括工具栏和滚动条:
			let innerWidth = window.innerWidth
			let innerHeight = window.innerHeight
			//设置区域为正方形，边长按照宽或高里面，比较小的值决定
			const setArea = 0.15;
			// innerWidth >= innerHeight 的情况下
			let MinY = innerHeight * (0.5 - setArea);
			let MaxY = innerHeight * (0.5 + setArea);
			let MinX = innerWidth * 0.5 - (MaxY - MinY) * 0.5;
			let MaxX = innerWidth * 0.5 + (MaxY - MinY) * 0.5;
			if (innerWidth < innerHeight) {
				MinX = innerWidth * (0.5 - setArea);
				MaxX = innerWidth * (0.5 + setArea);
				MinY = innerHeight * 0.5 - (MaxX - MinX) * 0.5;
				MaxY = innerHeight * 0.5 + (MaxX - MinX) * 0.5;
			}
			//在设置区域
			let inSetArea = false
			if ((this.clickX > MinX && this.clickX < MaxX) && (this.clickY > MinY && this.clickY < MaxY)) {
				//console.log("点中了设置区域！");
				inSetArea = true
			}
			if (inSetArea) {
				this.drawerActivate('right')
			}

			// console.log("window.innerWidth=", window.innerWidth, "window.innerHeight=", window.innerHeight);
			// console.log("MinX=", MinX, "MaxX=", MaxX);
			// console.log("MinY=", MinY, "MaxY=", MaxY);
			// console.log("x=", e.x, "y=", e.y);
		},

		onMouseMove(e) {
			this.clickX = e.x //获取鼠标的X坐标（鼠标与屏幕左侧的距离,单位为px）
			this.clickY = e.y //获取鼠标的Y坐标（鼠标与屏幕顶部的距离,单位为px）
			//浏览器的视口,不包括工具栏和滚动条:
			let innerWidth = window.innerWidth
			let innerHeight = window.innerHeight
			//设置区域为正方形，边长按照宽或高里面，比较小的值决定
			const setArea = 0.15;
			// innerWidth >= innerHeight 的情况下
			let MinY = innerHeight * (0.5 - setArea);
			let MaxY = innerHeight * (0.5 + setArea);
			let MinX = innerWidth * 0.5 - (MaxY - MinY) * 0.5;
			let MaxX = innerWidth * 0.5 + (MaxY - MinY) * 0.5;
			if (innerWidth < innerHeight) {
				MinX = innerWidth * (0.5 - setArea);
				MaxX = innerWidth * (0.5 + setArea);
				MinY = innerHeight * 0.5 - (MaxX - MinX) * 0.5;
				MaxY = innerHeight * 0.5 + (MaxX - MinX) * 0.5;
			}
			//在设置区域
			let inSetArea = false
			if ((this.clickX > MinX && this.clickX < MaxX) && (this.clickY > MinY && this.clickY < MaxY)) {
				inSetArea = true
			}
			if (inSetArea) {
				//console.log("在设置区域！");
				e.currentTarget.style.cursor = 'url(/images/SettingsOutline.png), pointer';
			} else {
				e.currentTarget.style.cursor = '';
			}

			//获取元素,统计页数
			// let offsetWidth = e.currentTarget.offsetWidth;
			// let offsetHeight = e.currentTarget.offsetHeight;



		},
		onMouseLeave(e) {
			//离开区域的时候,清空鼠标样式
			e.currentTarget.style.cursor = '';
		},

		scrollToTop(scrollDuration) {
			let scrollStep = -window.scrollY / (scrollDuration / 15),
				scrollInterval = setInterval(function () {
					if (window.scrollY !== 0) {
						window.scrollBy(0, scrollStep);
					}
					else clearInterval(scrollInterval);
				}, 15);
		},
		//根据可视区域(viewport)长宽比,确认是横屏还是竖屏
		// aspect-ratio https://developer.mozilla.org/zh-CN/docs/Web/CSS/@media/aspect-ratio
		// window.innerWidth  不是响应式依赖,所以不能用计算属性
		inLandscapeModeCheck() {
			//避免除数等于0,虽然正常情况下不会触发
			// if (window.innerHeight == 0) {
			// 	return false
			// }
			// var aspectRatio = document.documentElement.clientWidth / document.documentElement.clientHeight
			this.aspectRatio = window.innerWidth / window.innerHeight
			// console.log("aspectRatio=" + this.aspectRatio);
			// 为了测试方便,阈值是正方形
			return this.aspectRatio > (19 / 19);
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


<style scoped>
.manga {
	max-width: 100%;
	background: v-bind("model.backgroundColor");
}

/* https://developer.mozilla.org/zh-CN/docs/Web/CSS/object-fit */
.manga img {
	margin: auto;
	/* object-fit: scale-down; */
	padding: 3px 0px;
	border-radius: 7px;
	box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

.manga canvas {
	margin: auto;
	/* object-fit: scale-down; */
	padding: 3px 0px;
	border-radius: 7px;
	box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}
</style>
