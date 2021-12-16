<template>
	<div id="BookScroll">
		<div v-for="(page, key) in this.$store.state.book.pages" :key="page.url" class="manga">
			<img v-lazy="page.url" v-bind:H="page.height" v-bind:W="page.width" v-bind:key="key"
				v-bind:class="page.image_type" />
			<p v-if="showPageNum">{{ key + 1 }}/{{ AllPageNum }}</p>
		</div>
		<p></p>
		<button v-scroll="onScroll" v-show="btnFlag" fab color="#bbcbff" bottom right @click="toTop">▲</button>
		<slot></slot>
	</div>
</template>

<script>
	export default {
		components: {
			// Header,
		},
		data() {
			return {
				page_mode: "multi",
				btnFlag: false,
				showPageNum: false,
				duration: 300,
				offset: 0,
				easing: "easeInOutCubic",
				AllPageNum: this.$store.state.book.all_page_num,
				message: {
					user_uuid: "",
					server_status: "",
					now_book_uuid: "",
					read_percent: 0.0,
					msg: "",
				},
			};
		},
		mounted() {
			this.initPage();
		},

		methods: {
			initPage() {
				this.$cookies.keys();
			},
			getBook: function() {
				return this.$store.state.book;
			},
			getNumber: function(number) {
				this.page = number;
				console.log(number);
			},
			onScroll(e) {
				if (typeof window === "undefined") return;
				const top = window.pageYOffset || e.target.scrollTop || 0;
				this.btnFlag = top > 20;
			},
			toTop() {
				this.$vuetify.goTo(0);
			},
			onChangeBook: function(e, uuid) {
				// 当前元素
				this.message.now_book_uuid = uuid;
				this.message.msg = "ChangeBook";
				this.$socket.send(JSON.stringify(this.message));
				this.getBook();
			},
		},
	};
</script>

<style>
</style>

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
