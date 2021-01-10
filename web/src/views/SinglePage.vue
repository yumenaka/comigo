<template>
  <div id="SinglePage">
        <Header>
        <h2>
          <a v-bind:href="'raw/' + book.name">{{ book.name }}</a>
        </h2>
      </Header>

      <div
        class="singlebox"
        style="
          width: 600px;
          height: 800px;
          border: 3px #cccccc dashed;
          margin: auto;
        "
      >
        <v-img
          contain
          lazy-src=book.pages[page].url
          max-height="600"
          max-width="800"
          src=book.pages[page].url
        ></v-img>
        <img v-bind:url="book.pages[0].url" height="95%" />
        <v-pagination v-model="page" :length="book.page_num" @input = "getNumber"> </v-pagination>
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