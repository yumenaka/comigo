<template>
	<div id="ScrollMode" class="manga" @mouseenter="onMouseEnter" @mouseleave="onMouseLeave">
		<Header :inShelf="false" :readMode="'scroll'" :setDownLoadLink="needDownloadLink()" :headerTitle="book.title"
			:bookID="book.id" :depth="book.depth" :showReturnIcon="true" :showSettingsIcon="true"
			:InfiniteDropdown="ScrollModeConfig.InfiniteDropdown" v-bind:style="{ background: model.interfaceColor }"
			@drawerActivate="drawerActivate">
		</Header>

		<!-- 分页条 -->
		<n-pagination class="flex flex-row justify-center m-2 text-lg font-semibold" v-if="!ScrollModeConfig.InfiniteDropdown"
			v-model:page="paginationPage" :page-count="Math.ceil(book.page_count / LOAD_LIMIT)"
			:on-update:page="onPaginationPageChange" />

		<!-- PDF格式的额外提示信息。-->
		<div v-if="book.type === '.pdf' && ((!ScrollModeConfig.InfiniteDropdown && paginationPage === 1) || ScrollModeConfig.InfiniteDropdown)"
			class="flex flex-row self-center justify-center w-full p-1 m-1 font-semibold text-center rounded shadow-2xl">
			<div class="pb-1">{{ $t('pdf_hint_message') }}</div>
			<a class="text-blue-700 underline " :href="'api/raw/' + book.id + '/' + encodeURIComponent(book.title)"
				target="_blank">{{ $t('original_pdf_link') }} </a>
		</div>

		<!-- 渲染漫画的主体部分 -->
		<div class="main_manga" v-for="(single_image, n) in localImages" :key="single_image.url"
			@click="onMouseClick($event)" @mousemove="onMouseMove">
			<ImageScroll  v-if="(!single_image.url.endsWith('.html'))&&(!single_image.url.includes('.hidden.'))"  :image_url="imageParametersString(single_image.url)" :sPWL="sPWL" :dPWL="dPWL" :sPWP="sPWP"
				:dPWP="dPWP" :nowPageNum="nowPageNum" :page_count="book.page_count" :book_id="book.id"
				:showPageNumFlag_ScrollMode="ScrollModeConfig.showPageNumFlag_ScrollMode" :syncPageByWS="ScrollModeConfig.syncPageByWS"
				:autoScrolling="autoScrolling" :userControlling="userControlling" :margin="ScrollModeConfig.marginOnScrollMode"
				@refreshNowPageNum="refreshNowPageNum">
			</ImageScroll>
		</div>

		<Observer @loadNextBlock="loadNextBlock" />
		<!-- 返回顶部的圆形按钮，向上滑动的时候出现 -->
		<n-back-top class="bg-blue-200" :show="showBackTopFlag" type="info" :right="20" :bottom="20" />

		<!-- 分页条 -->
		<n-pagination class="flex flex-row justify-center m-2 text-lg font-semibold" v-if="!ScrollModeConfig.InfiniteDropdown"
			v-model:page="paginationPage" :page-count="Math.ceil(book.page_count / LOAD_LIMIT)"
			:on-update:page="onPaginationPageChange" />

		<!-- 底部最下面的返回顶部按钮 -->
		<button v-if="ScrollModeConfig.InfiniteDropdown" class="w-24 h-12 m-2 text-gray-900 bg-blue-300 rounded hover:bg-blue-500"
			@click="scrollToTop(90);" size="large">{{ $t('back-to-top') }}</button>

		<QuickJumpBar class="self-center mt-2 mb-2" :nowBookID="book.id" :readMode="'scroll'"
			:InfiniteDropdown="ScrollModeConfig.InfiniteDropdown"></QuickJumpBar>

		<!-- 底部页脚 -->
		<Bottom v-bind:style="{ background: model.interfaceColor }" class="flex-none h-12" :ServerName="$store.state.server_status.ServerName
		? $store.state.server_status.ServerName
		: 'Comigo'
		"></Bottom>

		<Drawer :initDrawerActive="drawerActive" :initDrawerPlacement="drawerPlacement"
			@saveConfig="saveConfigToLocalStorage" @startSketch="startSketchMode" @closeDrawer="drawerDeactivate"
			@setT="OnSetTemplate" :readerMode="ScrollModeConfig.readerMode" :inBookShelf="false" :sketching="false">

			<!-- 选择：切换页面模式 -->
			<n-button @click="changeReaderModeToFlipMode">{{ $t('switch_to_flip_mode') }}</n-button>
			<!-- 开关：卷轴模式下，是无限下拉，还是分页加载 -->
			<n-switch size="large" v-model:value="ScrollModeConfig.InfiniteDropdown" :rail-style="railStyle"
				@update:value="setInfiniteDropdown">
				<template #checked>{{ $t("infinite_dropdown") }}</template>
				<template #unchecked>{{ $t("pagination_mode") }}</template>
			</n-switch>

			<!-- 页面加载模式 -->
			<n-select :placeholder='$t("re_sort_page")' @update:value="onResort" :options="options" />

			<!-- 无限加载下拉，还是分页模式的切换按钮 -->
			<!-- <n-button v-if="!InfiniteDropdown" @click="setInfiniteDropdown(true)">{{ $t('to_infinite_dropdown_mode')
			}}</n-button>
			<n-button v-if="InfiniteDropdown" @click="setInfiniteDropdown(false)">{{ $t('to_pagination_mode')
			}}</n-button> -->

			<!-- 同步翻页 -->
			<n-switch size="large" v-model:value="ScrollModeConfig.syncPageByWS" @update:value="setSyncPageByWS">
				<template #checked>{{ $t("sync_page") }}</template>
				<template #unchecked>{{ $t("sync_page") }}</template>
			</n-switch>

			<!-- 显示页数 -->
			<n-switch size="large" v-model:value="ScrollModeConfig.showPageNumFlag_ScrollMode" @update:value="setShowPageNumChange">
				<template #checked>{{ $t('showPageNum') }}</template>
				<template #unchecked>{{ $t('showPageNum') }}</template>
			</n-switch>

			<!-- 开关：横屏状态下,宽度单位是百分比还是固定值 -->
			<n-switch size="large" v-model:value="ScrollModeConfig.imageWidth_usePercentFlag" :rail-style="railStyle"
				@update:value="setImageWidthUsePercentFlag">
				<template #checked>{{ $t('width_usePercent') }}</template>
				<template #unchecked>{{ $t('width_useFixedValue') }}</template>
			</n-switch>

			<!-- 单页-漫画宽度-使用百分比 -->
			<!-- 数字输入% -->
			<n-input-number v-if="ScrollModeConfig.imageWidth_usePercentFlag" size="small" :show-button="false"
				v-model:value="ScrollModeConfig.singlePageWidth_Percent" :max="100" :min="10" :update-value-on-input="false">
				<template #prefix>{{ $t('singlePageWidth') }}</template>
				<template #suffix>%</template>
			</n-input-number>

			<!-- 滑动选择% -->
			<n-slider v-if="ScrollModeConfig.imageWidth_usePercentFlag" v-model:value="ScrollModeConfig.singlePageWidth_Percent" :step="1" :max="100"
				:min="10" :format-tooltip="(value: any) => `${value}%`" />

			<!-- 开页-漫画宽度-使用百分比  -->
			<!-- 数字输入% -->
			<n-input-number v-if="ScrollModeConfig.imageWidth_usePercentFlag" size="small" :show-button="false"
				v-model:value="ScrollModeConfig.doublePageWidth_Percent" :max="100" :min="10" :update-value-on-input="false">
				<template #prefix>{{ $t('doublePageWidth') }}</template>
				<template #suffix>%</template>
			</n-input-number>
			<!-- 滑动选择% -->
			<n-slider v-if="ScrollModeConfig.imageWidth_usePercentFlag" v-model:value="ScrollModeConfig.doublePageWidth_Percent" :step="1" :max="100"
				:min="10" :format-tooltip="(value: number) => `${value}%`" />

			<!-- 单页-漫画宽度-使用固定值PX -->
			<!-- 数字输入框 -->
			<n-input-number v-if="!ScrollModeConfig.imageWidth_usePercentFlag" size="small" :show-button="false"
				v-model:value="ScrollModeConfig.singlePageWidth_PX" :min="50" :update-value-on-input="false">
				<template #prefix>{{ $t('singlePageWidth') }}</template>
				<template #suffix>px</template>
			</n-input-number>

			<!-- 滑动选择PX -->
			<n-slider v-if="!ScrollModeConfig.imageWidth_usePercentFlag" v-model:value="ScrollModeConfig.singlePageWidth_PX" :step="10" :max="1600"
				:min="50" :format-tooltip="(value: any) => `${value}px`" />

			<!-- 数字输入框 -->
			<n-input-number v-if="!ScrollModeConfig.imageWidth_usePercentFlag" size="small" :show-button="false"
				v-model:value="ScrollModeConfig.doublePageWidth_PX" :min="50" :update-value-on-input="false">
				<template #prefix>{{ $t('doublePageWidth') }}</template>
				<template #suffix>px</template>
			</n-input-number>

			<!-- 滑动选择PX -->
			<n-slider v-if="!ScrollModeConfig.imageWidth_usePercentFlag" v-model:value="ScrollModeConfig.doublePageWidth_PX" :step="10" :max="1600"
				:min="50" :format-tooltip="(value: any) => `${value}px`" />

			<!-- 页面间隙px -->
			<n-input-number v-if="!ScrollModeConfig.imageWidth_usePercentFlag" size="small" :show-button="false"
				v-model:value="ScrollModeConfig.marginOnScrollMode" :min="0" :update-value-on-input="false">
				<template #prefix>{{ $t('marginOnScrollMode') }}</template>
				<template #suffix>px</template>
			</n-input-number>
			<n-slider v-model:value="ScrollModeConfig.marginOnScrollMode" :step="1" :max="100" :min="0"
				:format-tooltip="(value: any) => `${value}px`" />

			<!-- 开关：自动切边 -->
			<n-switch size="large" v-model:value="imageParameters.do_auto_crop"
				@update:value="setImageParameters_DoAutoCrop">
				<template #checked>{{ $t('auto_crop') }}</template>
				<template #unchecked>{{ $t('auto_crop') }}</template>
			</n-switch>
			<!-- 切白边阈值 -->
			<n-input-number :show-button="false" v-if="imageParameters.do_auto_crop"
				v-model:value="imageParameters.auto_crop_num" :max="10" :min="0">
				<template #prefix>{{ $t('energy_threshold') }}</template>
			</n-input-number>

			<!-- 开关：压缩图片 -->
			<n-switch size="large" :rail-style="railStyle" v-model:value="imageParameters.do_auto_resize"
				@update:value="setImageParameters_DoAutoResize">
				<template #checked>{{ $t('image_width_limit') }}</template>
				<template #unchecked>{{ $t('raw_resolution') }}</template>
			</n-switch>
			<!-- 压缩图片参数：数字输入框 -->
			<n-input-number v-if="imageParameters.do_auto_resize" size="small" :show-button="false"
				v-model:value="imageParameters.resize_max_width" :min="100">
				<template #prefix>{{ $t('max_width') }}</template>
				<template #suffix>px</template>
			</n-input-number>
		</Drawer>
	</div>
</template>

<script lang="ts">
// 直接导入组件并使用它。这种情况下,只有导入的组件才会被打包。
import { NBackTop, NButton, NInputNumber, NSelect, NSlider, NSwitch, useDialog, useMessage, NPagination, } from 'naive-ui'
import Header from "@/components/Header.vue";
import Drawer from "@/components/Drawer.vue";
import Bottom from "@/components/Bottom.vue";
import Observer from "@/components/Observer_in_Scroll.vue";
import ImageScroll from "@/components/Image_in_Scroll.vue";
import QuickJumpBar from "@/components/QuickJumpBar.vue";
import { CSSProperties, defineComponent, reactive, ref } from 'vue'
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
		NSlider,//滑动选择  Slider https://www.naiveui.com/zh-CN/os-theme/components/slider
		NSwitch,//开关   https://www.naiveui.com/zh-CN/os-theme/components/switch
		NInputNumber,//数字输入 https://www.naiveui.com/zh-CN/os-theme/components/input-number
		NButton,//按钮，来自:https://www.naiveui.com/zh-CN/os-theme/components/button
		NSelect,
		QuickJumpBar,
		NPagination,
	},
	// setup在创建组件前执行，因此没有this
	setup() {
		//此处不能使用this,但可以用getCurrentInstance 这个vue函数取得Proxy，实现类似 proxy.$socket.onmessage 这样的调用(https://github.com/likaia/vue-native-websocket-vue3)。
		// const { cookies } = useCookies();
		//在setup中访问 vuex 需要通过useStore()来访问  https://juejin.cn/post/6917592199140458504#heading-22=
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
			paginationPage: ref(1),
			// cookies,
			//背景色
			model,
			message,
			dialog,
			imageParameters,//获取图片所用参数
			imageParametersString: (source_url: string) => {
				//console.log("source_url:" + source_url)
				if (source_url.startsWith("api/get_file?")) {
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
					return base_str + addStr + nocache_str;
				} else {
					return source_url
				}
			},
			//开关的颜色
			// 开关的颜色
			railStyle: ({
				focused,
				checked
			}: {
				focused: boolean
				checked: boolean
			}) => {
				const style: CSSProperties = {}
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
		}
	},
	data() {
		return {
			debug: false,
			//加载页面相关
			nowPageNum: 1,//当前页数,从1开始,,不是数组下标,在pages.images数组当中用的时候需要-1
			LOAD_LIMIT: 20,//一次载入的漫画张数，默认为20.
			lastPageNum: 20,//载入漫画的最后一页，默认为20.
			localImages: [
				{
					filename: "",
					url: ""
				},
				{
					filename: "",
					url: ""
				}
			],
			nowLoading: true,//是否正在加载，延迟执行操作、隐藏按钮用
			userControlling: false,//用户是否正在操控，操控的时候不接收、也不发送WS翻页消息
			autoScrolling: false,//是否正在自动翻页，为真的时候，不发送WS消息
			book: {
				title: "loading",
				id: "abcde",
				page_count: 2,
				type: ".zip",
				depth: 0,
				pages: {
					sort_by: "",
					images: [
						{
							filename: "001",
							url: "/images/loading.gif",
						},
						{
							filename: "002",
							url: "/images/loading.gif",
						},
					],
				}
			},
			resort_hint_key: "resort",
			options: [
				{
					label: this.$t('sort_by_filename'),
					value: "filename",
				},
				{
					label: this.$t('sort_by_modify_time'),
					value: "modify_time"
				},
				{
					label: this.$t('sort_by_filesize'),
					value: "filesize"
				},
				{
					label: this.$t('sort_by_filename') + this.$t('sort_reverse'),
					value: "filename_reverse",
				},
				{
					label: this.$t('sort_by_modify_time') + this.$t('sort_reverse'),
					value: "modify_time_reverse"
				},
				{
					label: this.$t('sort_by_filesize') + this.$t('sort_reverse'),
					value: "filesize_reverse"
				},
			],

			drawerActive: false,
			drawerPlacement: 'right',
			//开发模式 还没有做的功能与设置,设置Debug以后才能见到
			debugModeFlag: true,
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
			//屏幕宽横比,inLandscapeMode的判断依据
			aspectRatio: 1.2,
			//可见范围宽高的具体值
			clientWidth: 0,
			clientHeight: 0,
			ScrollModeConfig: {
				InfiniteDropdown: true,//卷轴模式下，是否无限下拉
				marginOnScrollMode: 10,//下拉模式下，漫画页面的底部间距。px。
				//状态驱动的动态 CSS https://v3.cn.vuejs.org/api/sfc-style.html#%E7%8A%B6%E6%80%81%E9%A9%B1%E5%8A%A8%E7%9A%84%E5%8A%A8%E6%80%81-css
				//图片宽度的单位,是否使用百分比
				imageWidth_usePercentFlag: false,
				//横屏(Landscape)状态的漫画页宽度,百分比
				singlePageWidth_Percent: 50,
				doublePageWidth_Percent: 95,
				//横屏(Landscape)状态的漫画页宽度。px。
				singlePageWidth_PX: 720,
				doublePageWidth_PX: 720,
				//可见范围是否是横向
				isLandscapeMode: true,
				isPortraitMode: false,
				//书籍数据,需要从远程拉取
				//是否显示顶部页头
				showHeaderFlag: true,
				//是否显示页数
				showPageNumFlag_ScrollMode: false,
				readerMode: "scroll",
				imageMaxWidth: 10,
				//ws翻页相关
				syncPageByWS: true,//是否通过websocket同步翻页
				interfaceColor: "#F5F5E4",
                backgroundColor: "#E0D9CD",
			},
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
		// 从本次存储读取设置
		let configString = localStorage.getItem('ScrollModeConfig');
		if (localStorage.getItem('ScrollModeConfig') !== null && typeof configString === "string") {
				this.ScrollModeConfig = JSON.parse(configString)
		}
		//UI配色
		const tempBackgroundColor = localStorage.getItem("BackgroundColor")
		if (typeof (tempBackgroundColor) === 'string') {
			this.ScrollModeConfig.backgroundColor = tempBackgroundColor;
		}
		const tempInterfaceColor = localStorage.getItem("InterfaceColor")
		if (typeof (tempInterfaceColor) === 'string') {
			this.ScrollModeConfig.interfaceColor = tempInterfaceColor
		}
		this.model.backgroundColor = this.ScrollModeConfig.backgroundColor;
		this.model.interfaceColor = this.ScrollModeConfig.interfaceColor;

		// 消息监听，即接收websocket服务端推送的消息. optionsAPI用法
		this.$options.sockets.onmessage = (data: any) => this.handlePacket(data);
		//根据文件名、修改时间、文件大小等要素排序的参数
		let sort_image_by_str = "";
		if (this.$route.query.sort_by) {
			sort_image_by_str = "&sort_by=" + this.$route.query.sort_by
		}

		this.ScrollModeConfig.InfiniteDropdown = true;
		this.LOAD_LIMIT = 20;
		// 使用Vue Router的this.$route.query获取查询参数
		// 判断是否是无限下拉模式
		if (this.$route.query.page !== null && this.$route.query.page !== undefined) {
			this.ScrollModeConfig.InfiniteDropdown = false;
			this.LOAD_LIMIT = 35;
			//如果没有指定页数,则默认为1
			this.paginationPage = 1;
			// 如果链接上有多个同名query，route.query会返回一个数组
			// 但是我不想这样玩，query预计为string，所以as一下防止报错
			const pageAsNumber = parseInt(this.$route.query.page as string, 10);//转换为数字,10进制
			if (!isNaN(pageAsNumber)) {
				this.paginationPage = pageAsNumber;
			}
		}

		//根据路由参数获取特定书籍
		this.nowLoading = true;
		axios
			.get("/get_book?id=" + this.$route.params.id + sort_image_by_str)
			.then((response) => {
				//请求接口成功的逻辑
				this.book = response.data;
				//确定一开始要加载多少页
				if (this.book.page_count >= this.LOAD_LIMIT) {
					this.lastPageNum = this.LOAD_LIMIT;
				} else {
					this.lastPageNum = this.book.page_count;
				}
				if (this.book.type == ".epub") {
					this.options.push(
						{
							label: this.$t('epub_info'),
							value: "epub_info",
						}
					)
				}
				this.loadPages();
			}).catch((error) => {
				console.log("请求接口失败" + error);
			})
			.finally(
				() => {
					document.title = this.book.title;
					// console.log("成功获取书籍数据,书籍ID:" + this.$route.params.id);
				}
			);
		//监听路由参数的变化,刷新本地的Book数据
		this.$watch(
			() => this.$route.params.id,
			(id: any) => {
				if (id) {
					axios
						.get("/get_book?id=" + this.$route.params.id + sort_image_by_str)
						.then((response) => (this.book = response.data))
						.finally(() => {
							console.log("路由参数改变,书籍ID:" + id);
							//refresh web page
							window.location.reload();
						}
						);
				}
			}
		);

		window.addEventListener("scroll", this.onScroll);
		//文档视图调整大小时会触发 resize 事件。 https://developer.mozilla.org/zh-CN/docs/Web/API/Window/resize_event
		window.addEventListener("resize", this.onResize);
		this.ScrollModeConfig.imageMaxWidth = window.innerWidth;
	},
	// 挂载前 : 指令第一次绑定到元素并且在挂载父组件之前调用。
	beforeMount() {
	},
	onMounted() {
		this.ScrollModeConfig.isLandscapeMode = this.inLandscapeModeCheck();
		this.ScrollModeConfig.isPortraitMode = !this.inLandscapeModeCheck();
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
	},
	methods: {
		//接收服务器发来的websocket消息，做各种反应（翻页、提示信息）
		handlePacket(data: { data: string; }) {
			if (!this.debug) {
				return;
			}
			if (!this.ScrollModeConfig.syncPageByWS) {
				return;
			}
			const msg = JSON.parse(data.data);
			//心跳信息,直接返回
			if (msg.type === "heartbeat") {
				return;
			}
			//用户正在操作，不对翻页消息作反应。
			if (this.userControlling) {
				console.log("handlePacket:Return,Because User Controlling");
				return;
			}
			var localUserID = this.$store.userID;
			//确保服务器发来翻页信息，来自于另一个用户
			if (msg.user_id === localUserID) {
				console.log(this.$store.userID + "接收到Message:msg.user_id=" + msg.user_id);
				return;
			}
			if (msg.type === "scroll_mode_sync_page") {
				// console.log("自动翻页 msg：", msg);
				const syncData = JSON.parse(msg.data_string);
				if (syncData.book_id !== this.book.id) {
					console.log("虽然有翻页信息，但" + syncData.book_id + "≠" + this.book.id);
					return;
				}
			}
		},

		//刷新到底部的时候,改变images数据
		loadPages() {
			if (this.ScrollModeConfig.InfiniteDropdown) {
				//slice() 方法返回一个新的数组对象，这一对象是一个由 begin 和 end 决定的原数组的浅拷贝（包括 begin，不包括end）
				this.localImages = this.book.pages.images.slice(0, this.lastPageNum);//javascript的接片不能直接用[a,b]，而是需要调用.slice()函数
			} else {
				let start = (this.paginationPage - 1) * this.LOAD_LIMIT;
				let end = this.paginationPage * this.LOAD_LIMIT;
				if (end > this.book.page_count) {
					end = this.book.page_count;
				}
				this.localImages = this.book.pages.images.slice(start, end);
			}
		},

		//无限加载用,底部刷新
		loadNextBlock() {
			if (!this.ScrollModeConfig.InfiniteDropdown) {
				return;
			}
			if (this.book.page_count >= this.lastPageNum + this.LOAD_LIMIT) {
				this.lastPageNum = this.lastPageNum + this.LOAD_LIMIT;
			} else {
				this.lastPageNum = this.book.page_count;
			}
			this.loadPages();
		},
		//监听子组件事件: https://v3.cn.vuejs.org/guide/component-basics.html#%E7%9B%91%E5%90%AC%E5%AD%90%E7%BB%84%E4%BB%B6%E4%BA%8B%E4%BB%B6
		//滚动页面的时候刷新页数
		refreshNowPageNum(n: number) {
			if (!this.ScrollModeConfig.InfiniteDropdown) {
				return;
			}
			if (this.nowLoading) {
				return
			}
			this.nowPageNum = n;
			this.loadPages();
		},

		//页面排序相关
		onResort(key: string) {
			axios
				.get("/get_book?id=" + this.$route.params.id + "&sort_by=" + key)
				.then((response) => (this.book = response.data))
				.finally(
					() => {
						document.title = this.book.title;
						this.resort_hint_key = key;
						this.loadPages();
						// 带查询参数，结果是 /#/scroll/abc123?sort_by="filesize"
						this.$router.push({ name: "ScrollMode", replace: true, query: { sort_by: key } });
						console.log("成功刷新书籍数据,书籍ID:" + this.$route.params.id + "  sort_by=" + key);
					}
				);
		},

		//“重新排序”选择菜单的文字提示
		getSortHintText(key: string) {
			switch (key) {
				case "filename": return this.$t('sort_by_filename');
				case "modify_time": return this.$t('sort_by_modify_time');
				case "filesize": return this.$t('sort_by_filesize');
				case "filename_reverse": return this.$t('sort_by_filename') + this.$t('sort_reverse');
				case "modify_time_reverse": return this.$t('sort_by_modify_time') + this.$t('sort_reverse');
				case "filesize_reverse": return this.$t('sort_by_filesize') + this.$t('sort_reverse');
				default:
					return this.$t('re_sort_page');
			}
		},

		//切换到翻页模式
		changeReaderModeToFlipMode() {
			localStorage.setItem("ReaderMode", "flip");
			//replace的作用类似于 router.push，唯一不同的是，它在导航时不会向 history 添加新记录，正如它的名字所暗示的那样——它取代了当前的条目。
			this.$router.push({ name: "FlipMode", replace: true, params: { id: this.$route.params.id } });
		},
		needDownloadLink() {
			return this.book.type !== "dir"
		},
		//打开抽屉
		drawerActivate(place: string) {
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
		OnSetTemplate(value: string) {
			this.$emit("setTemplate", value);
		},
		onBackgroundColorChange(value: string) {
			this.model.backgroundColor = value;
			this.ScrollModeConfig.backgroundColor = value;
      		localStorage.setItem("BackgroundColor", value);
			localStorage.setItem('ScrollModeConfig', JSON.stringify(this.ScrollModeConfig));
		},
		//关闭抽屉的时候保存配置 
		//如果在一个组件上使用了 v-model:xxx,应该使用 @update:xxx  https://www.naiveui.com/zh-CN/os-theme/docs/common-issues
		saveConfigToLocalStorage() {
			localStorage.setItem('ScrollModeConfig', JSON.stringify(this.ScrollModeConfig));
		},
		setSyncPageByWS(value: boolean) {
			this.ScrollModeConfig.syncPageByWS = value;
			localStorage.setItem("SyncPageFlag", value ? "true" : "false");
		},

		//InfiniteDropdown
		setInfiniteDropdown(value: boolean) {
			this.ScrollModeConfig.InfiniteDropdown = value;
			localStorage.setItem("InfiniteDropdown", value ? "true" : "false");
			if (!this.ScrollModeConfig.InfiniteDropdown) {
				this.paginationPage = 1;
				this.LOAD_LIMIT = 35;
				//  Vue Router, 重定向到新 URL
				this.$router.push({ name: "ScrollMode", query: { page: this.paginationPage } }).then(() => {
					//刷新页面
					window.location.reload();
				});
			} else {
				this.LOAD_LIMIT = 20;
				// 重定向到新的 URL
				this.$router.push({ name: "ScrollMode" }).then(() => {
					//刷新页面
					window.location.reload();
				});
			}
		},
		onPaginationPageChange(value: number) {
			this.paginationPage = value;
			this.$router.push({ name: "ScrollMode", query: { page: value } });
			this.loadPages();
		},
		//设置是否显示顶部页头
		setShowHeaderChange(value: boolean) {
			this.ScrollModeConfig.showHeaderFlag = value;
			localStorage.setItem("showHeaderFlag", value ? "true" : "false");
		},
		setShowPageNumChange(value: boolean) {
			this.ScrollModeConfig.showPageNumFlag_ScrollMode = value;
			localStorage.setItem("showPageNumFlag_ScrollMode", value ? "true" : "false");
		},
		//图片处理相关
		//黑白化参数
		setImageParameters_Gray(value: boolean) {
			this.imageParameters.gray = value;
			localStorage.setItem("ImageParameters_Gray", value ? "true" : "false");
		},
		//缩放图片大小的参数
		setImageParameters_DoAutoResize(value: boolean) {
			this.imageParameters.do_auto_resize = value;
			localStorage.setItem("ImageParameters_DoAutoResize", value ? "true" : "false");
		},
		//设置是否切白边
		setImageParameters_DoAutoCrop(value: boolean) {
			this.imageParameters.do_auto_crop = value;
			localStorage.setItem("ImageParameters_DoAutoCrop", this.imageParameters.do_auto_crop ? "true" : "false");
		},
		//切白边参数
		setImageParameters_AutoCropNum(value: number) {
			this.imageParameters.auto_crop_num = value;
			localStorage.setItem("ImageParameters_AutoCropNum", this.imageParameters.auto_crop_num.toString());
		},

		setImageWidthUsePercentFlag(value: boolean) {
			this.ScrollModeConfig.imageWidth_usePercentFlag = value;
			localStorage.setItem("imageWidth_usePercentFlag", value ? "true" : "false");
		},

		setDebugModeFlag(value: boolean) {
			this.debugModeFlag = value;
			localStorage.setItem("debugModeFlag", value ? "true" : "false");
		},

		//可见区域变化的时候改变页面状态
		onResize() {
			this.ScrollModeConfig.imageMaxWidth = window.innerWidth
			this.clientWidth = document.documentElement.clientWidth
			this.clientHeight = document.documentElement.clientHeight
			// var aspectRatio = window.innerWidth / window.innerHeight
			this.aspectRatio = this.clientWidth / this.clientHeight
			// 为了调试的时候方便,阈值是正方形
			if (this.aspectRatio > (19 / 19)) {
				this.ScrollModeConfig.isLandscapeMode = true
				this.ScrollModeConfig.isPortraitMode = false
			} else {
				this.ScrollModeConfig.isLandscapeMode = false
				this.ScrollModeConfig.isPortraitMode = true
			}
		},
		//页面滚动的时候,改变返回顶部按钮的显隐
		onScroll() {
			let scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
			this.scrollDownFlag = scrollTop > this.scrollTopSave;
			//防手抖,小于一定数值状态就不变 Math.abs()会导致报错
			let step = this.scrollTopSave - scrollTop;
			// console.log("this.scrollDownFlag:",this.scrollDownFlag,"scrollTop:",scrollTop,"step:", step);
			this.scrollTopSave = scrollTop
			if (step < -5 || step > 5) {
				this.showBackTopFlag = ((scrollTop > 400) && !this.scrollDownFlag);
			}
		},
		//Vue鼠标事件和键盘事件、以及触屏事件 http://www.zhuanbike.com/thread-288.htm
		// 只有进入绑定该事件的元素才会触发事件,也就是不会冒泡。其对应的离开事件mouseleave
		onMouseEnter() {
			//进入控制模式是立刻出发
			this.userControlling = true;
			//console.log("MouseEnter User Controlling =" + this.userControlling);

		},
		onMouseLeave(e: any) {
			//退出控制模式、延迟2秒执行
			let _this = this;
			setTimeout(function () {
				_this.userControlling = true;
				//console.log("User Controlling=" + _this.userControlling);
			}, 500);
			//离开区域的时候,清空鼠标样式
			e.currentTarget.style.cursor = '';

		},
		//获取鼠标位置,决定是否打开设置面板
		onMouseClick(e: any) {
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
		},

		onMouseMove(e: any) {
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
		},

		scrollToTop(scrollDuration: number) {
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
			this.aspectRatio = window.innerWidth / window.innerHeight
			// console.log("aspectRatio=" + this.aspectRatio);
			// 为了测试方便,阈值是正方形
			return this.aspectRatio > (19 / 19);
		},
	},

	computed: {
		sPWL() {
			if (this.ScrollModeConfig.imageWidth_usePercentFlag) {
				return this.ScrollModeConfig.singlePageWidth_Percent + '%';
			} else {
				return this.ScrollModeConfig.singlePageWidth_PX + 'px';
			}
		},
		dPWL() {
			if (this.ScrollModeConfig.imageWidth_usePercentFlag) {
				return this.ScrollModeConfig.doublePageWidth_Percent + '%';
			} else {
				return this.ScrollModeConfig.doublePageWidth_PX + 'px';
			}
		},
		sPWP() {
			if (this.ScrollModeConfig.imageWidth_usePercentFlag) {
				return this.ScrollModeConfig.singlePageWidth_Percent + '%';
			} else {
				return this.ScrollModeConfig.singlePageWidth_PX + 'px';
			}
		},
		dPWP() {
			if (this.ScrollModeConfig.imageWidth_usePercentFlag) {
				return this.ScrollModeConfig.doublePageWidth_Percent + '%';
			} else {
				return this.ScrollModeConfig.doublePageWidth_PX + 'px';
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
</style>
