<template>
  <div id="app">
    <router-view />
  </div>
</template>
<script>
export default {
  name: 'app'
}
</script>

<style lang="scss">

$fa-font-path: "~font-awesome/fonts";
@import '~font-awesome/scss/font-awesome';

//@import "~simple-line-icons/scss/simple-line-icons";

@import "~@coreui/coreui/scss/coreui";

//@import "~@coreui/icons/scss/coreui-icons";

//@import "~pace-progress/themes/blue/pace-theme-minimal.css";

@import "~toastr/toastr";

@import "~codemirror/lib/codemirror.css";

@import "~bootstrap-select/sass/bootstrap-select";

@import "./assets/css/bootstrap-custom.scss";

@import "./assets/css/main.scss";

@import "./assets/css/dashboard.scss";

@import "./assets/css/dashboard-table.scss";

@import "./assets/css/login.scss";

</style>