<template>
	<div id="BookScroll" v-if="this.book" class="manga">
		<Header v-if="this.showHeader">
			<h2>
				<p v-if="book.IsFolder" :href="'raw/' + book.name">{{ book.name }}</p>
				<a v-if="!book.IsFolder" :href="'raw/' + book.name">{{ book.name }}【Download】</a>
			</h2>
		</Header>
		<div v-for="(page, key) in book.pages" :key="page.url">
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
		<button @click="scrollToTop(90);">Back To Top</button>
	</div>
</template>

<script>
import Header from "@/components/Header.vue";
import { onMounted } from 'vue'
export default {
	setup() {
		onMounted(() => {
			console.log('mounted in the composition api!')
		})
	},
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
	beforeMount() {
		console.log("BookScroll beforeMount");
		this.axios
			.get("/book.json")
			.then((response) => {
				if (response.status == 200) {
					this.book = response.data;
					// console.log(this.book);
				}
			})
			.catch((error) => alert(error));
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

	// computed: {
	// 	// 计算属性的 getter
	// 	getBook() {
	// 		// `this` 指向 vm 实例
	// 		return this.author.books.length > 0 ? 'Yes' : 'No'
	// 	}
	// },
};


</script>

<style></style>

<style scoped>
.manga {
	max-width: 100%;
}

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
	width: 100px;
	height: 50px;
	background-color: #bbcbff;
	font-family: Avenir, Helvetica, Arial, sans-serif;
	-webkit-font-smoothing: antialiased;
	-moz-osx-font-smoothing: grayscale;
	font-size: 13px;
}

.LoadingImage {
	width: 80vw;
	max-width: 80vw;
}
.ErrorImage {
	width: 80vw;
	max-width: 80vw;
}

/* 大致上是横向较宽（显示区域）时的CSS样式，IE无效 */
@media screen and (min-aspect-ratio: 19/19) {
	.SinglePageImage {
		width: 50vw;
	}
	.DoublePageImage {
		width: 90vw;
	}
}

/* 竖屏(显示区域)CSS样式，IE无效 */
@media screen and (max-aspect-ratio: 19/19) {
	.manga img {
		max-width: 100%;
	}
	.SinglePageImage {
		width: 100vw;
	}
	.DoublePage {
		width: 100vw;
	}
}

/* 高分横屏（显示区域）时的CSS样式，IE无效 */
/* min-width 输出设备中页面最小可视区域宽度 大于这个width时，其中的css起作用 超宽屏 */
@media screen and (min-aspect-ratio: 19/19) and (min-width: 1922px) {
	.SinglePageImage {
		width: 1000px;
	}
	.DoublePageImage {
		width: 1900px;
	}
}
</style>
