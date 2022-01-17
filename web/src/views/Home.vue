<template>
	<div class="home">
		<Header v-if="this.showHeader">
			<h2>
				<a
					v-if="!this.$store.state.book.IsFolder"
					v-bind:href="'raw/' + this.$store.state.book.name"
				>{{ this.$store.state.book.name }}【Download】</a>
				<a
					v-if="this.$store.state.book.IsFolder"
					v-bind:href="'raw/' + this.$store.state.book.name"
				>{{ this.$store.state.book.name }}</a>
			</h2>
		</Header>
	</div>
</template>

<script>
import {
	onBeforeMount,
	onMounted,
	onBeforeUpdate,
	onUpdated,
	onBeforeUnmount,
	onUnmounted,
	onActivated,
	onDeactivated,
	onErrorCaptured,
} from "vue";

// @ is an alias to /src
import Header from "@/components/Header.vue";
// import ScrollMode from "@/components/ScrollMode.vue";

export default {
	name: "Home", //默认为 default。如果 <router-view>设置了名称，则会渲染对应的路由配置中 components 下的相应组件。
	components: {
		Header,
		// ScrollMode,
	},
	setup() {
		onBeforeMount(() => {
			console.log("Before Mount!");
		});
		onMounted(() => {
			console.log("Mounted!");
		});
		onBeforeUpdate(() => {
			console.log("Before Update!");
		});
		onUpdated(() => {
			console.log("Updated!");
		});
		onBeforeUnmount(() => {
			console.log("Before Unmount!");
		});
		onUnmounted(() => {
			console.log("Unmounted!");
		});
		onActivated(() => {
			console.log("Activated!");
		});
		onDeactivated(() => {
			console.log("Deactivated!");
		});
		onErrorCaptured(() => {
			console.log("Error Captured!");
		});
	},
	//组件的 data 选项必须是一个函数
	//每个实例可以维护一份被返回对象的独立的拷贝
	data() {
		return {
			book: null,
			showHeader: true,
			//如果你知道你会在晚些时候需要一个 property，但是一开始它为空或不存在，那么你仅需要设置一些初始值。
			bookshelf: {},
			setting: {
				template: "scroll",
				sketch_count_seconds: 90,
			},
			now_page: 1,
			duration: 300,
			offset: 0,
			easing: "easeInOutCubic",
			message: {
				user_uuid: "",
				server_status: "",
				now_book_uuid: "",
				read_percent: 0.0,
				msg: "",
			},
		};
	},
	computed: {
		// 计算属性的 getter
		nowTemplate: function () {
			var localValue = this.$cookies.get("nowTemplate");
			console.log("computed 1:" + localValue);
			if (localValue != null) {
				return localValue;
			} else {
				return this.$store.state.setting.template;
			}
		},
	},
	methods: {
		initPage() {
			this.$store.dispatch("syncBookDataAction");
			this.$store.dispatch("syncSettingDataAction");
			this.$store.dispatch("syncBookShelfDataAction");
		},
		getNumber: function (number) {
			this.page = number;
			console.log(number);
		},
		getNowTemplate: function () {
			var localValue = this.$cookies.get("nowTemplate");
			console.log("computed 1:" + localValue);
			if (localValue != null) {
				return localValue;
			} else {
				this.$cookies.set("nowTemplate", this.$store.state.setting.template);
				console.log("computed 2:" + this.$store.state.setting.template);
				return this.$store.state.setting.template;
			}
		},
	},
};
</script>
