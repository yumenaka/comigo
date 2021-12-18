<template>
	<div id="SingleMode" v-if="this.book">
		<Header v-if="this.showHeader">
			<h2>
				<p v-if="book.IsFolder" :href="'raw/' + book.name">{{ book.name }}</p>
				<a v-if="!book.IsFolder" :href="'raw/' + book.name">{{ book.name }}【Download】</a>
			</h2>
		</Header>
		<div id="SinglePageTemplate">
			<div class="single_page_main">
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
	</div>
</template>

<script>
import { useCookies } from "vue3-cookies";
import Header from "@/components/Header.vue";
export default {
	components: {
		Header,
	},
	setup() {
		const { cookies } = useCookies();
		return { cookies };
	},
	data() {
		return {
			book: null,
			showHeader: true,
			showPageNum: false,
			now_page: 1,
			showPagination: true,
		};
	},
	//在选项API中使用 Vue 生命周期钩子
	mounted() {
		// window.addEventListener("scroll", this.handleScroll);
		console.log("SingleMode created!");
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
	deactivated() {
		// 销毁监听
		window.removeEventListener("keyup", this.handleKeyup);
		// window.removeEventListener("scroll", this.handleScroll);
	},

	methods: {
		scrollToTop(scrollDuration) {
			var scrollStep = -window.scrollY / (scrollDuration / 15),
				scrollInterval = setInterval(function () {
					if (window.scrollY != 0) {
						window.scrollBy(0, scrollStep);
					}
					else clearInterval(scrollInterval);
				}, 15);
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
};
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
