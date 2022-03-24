import Vue from 'vue'
import App from './App.vue'
import VueRouter from "vue-router";
import Index from "./components/Index";
import Settings from "./components/Settings";
import page404 from "./components/tech/404";
import Login from "./components/tech/login";
import ListPage from "./components/list";
import VmList from "./components/VmList";
import EsxiList from "./components/EsxiList";
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'
import VueLoading from 'vue-loading-overlay';
import 'vue-loading-overlay/dist/vue-loading.css';
import VModal from 'vue-js-modal/dist/index.nocss.js'
import 'vue-js-modal/dist/styles.css'

Vue.use(BootstrapVue);
Vue.use(IconsPlugin);
Vue.use(VueRouter);
Vue.use(VueLoading);
Vue.use(VModal);

const routes = [
  {
    path: "/",
    name: "index",
    component: Index,
    meta: { title: 'Home' }
  },
  {
    path: "/vm/list",
    name: "vm",
    component: VmList,
    meta: { title: 'Virtual Machines' }
  },
  {
    path: "/esxi/list",
    name: "esxi",
    component: EsxiList,
    meta: { title: 'Servers' }
  },
  {
    path: "/dc/list",
    name: "dc",
    component: ListPage,
    meta: { title: 'Data Centers' }
  },
  {
    path: "/net/list",
    name: "net",
    component: ListPage,
    meta: { title: 'Networks' }
  },
  {
    path: "/ip/list",
    name: "ip",
    component: ListPage,
    meta: { title: 'IPs' }
  },
  {
    path: "/user/list",
    name: "user",
    component: ListPage,
    meta: { title: 'Users' }
  },
  {
    path: "/settings",
    name: "settings",
    component: Settings,
    meta: { title: 'Settings' }
  },
  {
    path: "/config/list",
    name: "config",
    component: ListPage,
    meta: { title: 'Configs' }
  },
  {
    path: "/log/list",
    name: "log",
    component: ListPage,
    meta: { title: 'Logs' }
  },
  {
    path: "/login",
    name: "login",
    component: Login,
    meta: { title: 'Login' }
  },
  {
    path: "/:page/list",
    name: "listPage",
    component: ListPage,
  },
  {
    path: "*",
    name: "404",
    component: page404,
    meta: { title: 'Not found' }
  },
];

const DEFAULT_TITLE = 'VM Manager';
const router = new VueRouter({
  routes,
  mode: "hash",
});

router.beforeEach((to, from, next) => {
  // redirect to login page if not logged in and trying to access a restricted page
  const publicPages = ['/login'];
  const authRequired = !publicPages.includes(to.path);
  const loggedIn = localStorage.getItem('token');

  if (authRequired && !loggedIn) {
    return next('/login');
  }

  next();
});
router.afterEach((to) => {
  Vue.nextTick(() => {
    let new_title = `${to.meta.title} - VMM`;
    document.title = new_title || DEFAULT_TITLE;
  });
});

Vue.config.productionTip = false;

new Vue({
  router,
  render: (h) => h(App),
}).$mount("#app");
