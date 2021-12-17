<template>
	<div id="BookSingle" v-if="this.book">
		<Header v-if="this.showHeader">
			<h2>
				<p v-if="book.IsFolder" :href="'raw/' + book.name">{{ book.name }}</p>
				<a v-if="!book.IsFolder" :href="'raw/' + book.name">{{ book.name }}【Download】</a>
			</h2>
		</Header>
		<div class="single_page_main">
			<img
				v-on:click="addPage(1)"
				v-if="now_page <= this.$store.state.book.all_page_num && now_page >= 1"
				lazy-src="/resources/favicon.ico"
				v-bind:src="this.$store.state.book.pages[now_page - 1].url"
			/>
			<img />
		</div>
	</div>
</template>

<script>
import Header from "@/components/Header.vue";
export default {
	components: {
		Header,
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
	created() {
		// 注册监听
		window.addEventListener("keyup", this.handleKeyup);
		// window.addEventListener("scroll", this.handleScroll);
		console.log("BookScroll created!");
		this.axios
			.get("/book.json")
			.then((response) => {
				if (response.status == 200) {
					this.book = response.data;
					console.log("response.status == 200");
					console.log(this.book);
				}
			})
			.catch(console.log("this.book"), (error) => alert(error));
	},


	destroyed() {
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
