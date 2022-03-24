<template>
  <!--<div class="page">
    <Header />
    <div class="index center">
    <h1>Welcome</h1>
    </div>
  </div>-->
  <div class="page">
    <Header />
    <div class="app-body">
      <main class="main">
        <div class="container-fluid">
          <div class="animated fadeIn index center">

          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
  import Header from './Header.vue';
export default {
  name: 'Index',
  components: {
    Header
  },
  props: {
    msg: String
  },
  mounted () {
    this.setAdditionalStyles();
  },
  methods: {
    setAdditionalStyles() {
      let body = document.querySelector('body');
      body.classList.remove('login');
    },
  }
}
</script>

<style scoped>
.center {
  vertical-align: middle;
  text-align: center;
}
</style>
