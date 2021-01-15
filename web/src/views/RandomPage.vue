<template>
  <div id="SinglePage">
      <Header>
        <h2>
          <a v-bind:href="'raw/' + book.name">{{ book.name }}</a>
        </h2>
      </Header>

      <div
        class="randombox"
        style="
          width: 600px;
          height: 800px;
          border: 3px #cccccc dashed;
          margin: auto;
        "
      >
        <img
          v-lazy="book.pages[1].url"
          v-bind:H="book.pages[1].height"
          v-bind:W="book.pages[1].width"
          height="95%"
        />
        <p>{{ book.page_num }}</p>
      </div>
    <slot></slot>

  </div>
</template>

<script>
export default {
    
}
</script>

<style>
</style>