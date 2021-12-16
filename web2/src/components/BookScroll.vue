<template>
	<div id="BookScroll" v-if="this.book">
		<Header v-if="this.showHeader">
			<h2>
				<p v-if="book.IsFolder" :href="'raw/' + book.name">{{ book.name }}</p>
				<a v-if="!book.IsFolder" :href="'raw/' + book.name">{{ book.name }}【Download】</a>
			</h2>
		</Header>
		<div v-for="(page, key) in book.pages" :key="page.url" class="manga">
			<!-- v-lazy="page.url"  :src="page.url" -->
			<img
				v-lazy="page.url"
				v-bind:H="page.height"
				v-bind:W="page.width"
				v-bind:key="key"
				v-bind:class="page.image_type"
			/>
			<p v-if="showPageNum">{{ key + 1 }}/{{ book.all_page_num }}</p>
		</div>
		<button @click="toTop">回到顶部</button>
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
		};
	},
	created() {
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
	methods: {
		toTop() { 
			window.scrollTo(0,0);
		},
	},
};
</script>

<style></style>

<style scoped>
.manga img {
	margin: auto;
	max-width: inherit;
	padding-top: 3px;
	padding-bottom: 3px;
	padding-right: 0px;
	padding-left: 0px;
	border-radius: 7px;
	box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
}

button {
	width: 70px;
	height: 40px;
	background-color: #bbcbff;
}

/* 竖屏(显示区域)CSS样式，IE无效 */
@media screen and (max-aspect-ratio: 19/19) {
	.SinglePage {
		width: 100%;
	}
	.DoublePage {
		width: 100%;
	}
}

/* 横屏（显示区域）时的CSS样式，IE无效 */
@media screen and (min-aspect-ratio: 19/19) {
	.SinglePage {
		width: 900px;
	}
	.DoublePage {
		width: 95%;
	}
}

/* 高分横屏（显示区域）时的CSS样式，IE无效 */
/* min-width 输出设备中页面最小可视区域宽度 大于这个width时，其中的css起作用 超宽屏 */
@media screen and (min-aspect-ratio: 19/19) and (min-width: 1922px) {
	.SinglePage {
		width: 1000px;
	}
	.DoublePage {
		width: 1900px;
	}
}
</style>
