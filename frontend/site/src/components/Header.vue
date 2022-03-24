<template>
  <div>
    <header class="app-header navbar">
      <button onclick="navbar.toggleSidebar()" class="navbar-toggler d-lg-none mr-auto" type="button"
              data-toggle="sidebar-show">
        <span class="navbar-toggler-icon"></span>
      </button>
      <button onclick="navbar.toggleSidebar()" class="navbar-toggler d-md-down-none" type="button"
              data-toggle="sidebar-lg-show">
        <span class="navbar-toggler-icon"></span>
      </button>
      <button onclick="navbar.toggleSidebar()" class="navbar-toggler navbar-toggler-mobile" type="button"
              data-toggle="sidebar-lg-show">
        <span class="navbar-toggler-icon"></span>
      </button>
      <a class="navbar-brand d-none d-md-inline-flex d-lg-inline-flex d-xl-inline-flex" href="/">
        vm manager
      </a>

      <dropdown-menu class="nav navbar-nav px-3 ml-auto">
        <a slot="trigger" class="nav-link" data-toggle="dropdown" href="#" role="button" aria-haspopup="true" aria-expanded="false">
          <i class="fa fa-user"></i>&nbsp;&nbsp;<span>Admin@admin</span>
        </a>
        <!--<div slot="header">Dropdown Header</div>-->
        <ul slot="body">
          <li>
            <a class="dropdown-item" href="#" @click="SignOut()">
              <i class="fa fa-lock"></i> Sign out
            </a>
          </li>
        </ul>
        <!--<div slot="footer">Dropdown Footer</div>-->
      </dropdown-menu>
    </header>

    <div class="sidebar">
      <nav class="sidebar-nav">
        <ul class="nav">
          <li v-for="(link, linkIndex) in links"
              v-bind:key="linkIndex"
              :class="'nav-item ' +  link.icon.type"
          >
            <router-link :to="link.to" :class="($route.path == link.to ? ' active ' : ' ') + ' nav-link ' + link.icon.type">
              <span v-if=" link.icon.type == 'bitmap' ">
                <img :src="`/img/${link.icon.value}`" />
              </span>
              <b-icon  v-else :icon="link.icon.value"></b-icon>
              <span class="link-text">{{ link.name }}</span>
            </router-link>
          </li>
        </ul>
      </nav>
      <button class="sidebar-minimizer brand-minimizer" type="button"></button>
    </div>

  </div>
</template>

<script>
import axios from "axios";
import DropdownMenu from 'v-dropdown-menu'
import 'v-dropdown-menu/dist/v-dropdown-menu.css' // Base style, required.

export default {
  name: 'Header',
  components: {
    DropdownMenu
  },
  props: ['title'],
  data() {
     return {
         sidebarMini: true,
         error: '',
         apiConfig: '/api/config/menu/left',
         apiSignOut: '/api/auth/sign-out/',
         apiUserInfo: '/api/user/',
         login: '/login',
         links: [
           {name: 'Items',  to: '/item/list',    icon: {type: 'bitmap', value: 'dc.png'}},
           {name: 'Groups',       to: '/group/list',  icon: {type: 'bitmap', value: 'esxi.png'}},
           {name: 'Configs',       to: '/config/list',icon: {type: 'bitmap', value: 'config.png'}},
           {name: 'Logs',          to: '/log/list',   icon: {type: 'bitmap', value: 'log.png'}},
        ],
         userMenuShow: false,
         userName: 'User@user',
     }
  },
    methods: {
        SignOut() {
            const token = JSON.parse(localStorage.getItem('token'));
            if (token == null) {
                this.error = 'token not found';
                return;
            }
            localStorage.removeItem('token');
            axios.post(this.apiSignOut,
                {},
                {headers: { 'Authorization': 'Bearer ' + token.accessToken}});
            this.$router.push({ path: this.login});
        },
      LoadConfig() {
          this.error = '';
          const token = JSON.parse(localStorage.getItem('token'));
          if (token == null) {
              this.error = 'token not found';
              return;
          }
          axios.get(this.apiConfig,
              {headers: { 'Authorization': 'Bearer ' + token.accessToken}})
              .catch((err) => {
                  if (err.response) {
                      if (err.response.status === 401) {
                          this.$router.push({ path: this.login});
                      }
                      this.error = 'Load config data: [' + err.response.status + ' ' + err.response.statusText + '] ';
                      if (err.response.data != '') {
                          this.error += err.response.data;
                      }
                  } else {
                      this.error = 'Load config data: Server error:  no connection'
                  }
              })
              .then(res => {
                  if (res !== null && res !== undefined) {
                      if (res.data !== null && res.data !== undefined) {
                          //this.links = res.data.links;
                          this.sidebarMini = res.data.sidebarMini;
                      }
                  }
              });
          axios.get(this.apiUserInfo,
              {headers: { 'Authorization': 'Bearer ' + token.accessToken}})
              .catch((err) => {
                  //TODO: create errors stack
                  console.log(err.response.status);
              })
              .then(res => {
                  if (res !== null && res !== undefined) {
                      if (res.data !== null && res.data !== undefined) {
                          this.userName = res.data.data.Email;
                          localStorage.setItem('user', JSON.stringify(res.data.data));
                      }
                  }
              });
          console.log('after load config func');

          axios.get(this.apiConfig,
          {headers: {'Authorization': 'Bearer ' + token.accessToken}})
          .catch((err) => {
            if (err.response) {
              if (err.response.status === 401) {
                this.$router.push({path: this.login});
              }
              this.error = 'Load config data: [' + err.response.status + ' ' + err.response.statusText + '] ';
              if (err.response.data != '') {
                this.error += err.response.data;
              }
            } else {
              this.error = 'Load config data: Server error:  no connection'
            }
          })
          .then(res => {
            if (res !== null && res !== undefined) {
              if (res.data !== null && res.data !== undefined) {
                //this.links = res.data.links;
                this.sidebarMini = res.data.sidebarMini;
              }
            }
          });
      console.log('after load config func');
    }
  },
  created: function () {
    this.LoadConfig();
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
.header-button {
  position: absolute;
  top: 10px;
  left: 15px;
  cursor: pointer;
}

.header-menu {
  height: 100%;
  width: 0;
  position: fixed;
  z-index: 1;
  left: 0px;
  background-color: #fff;
  overflow-x: hidden;
  transition: 0.5s;
  border: 2px solid;
}

.header-user {
    float: right;
    margin: 10px 20px 0 0;
    padding: 2px;
    font-weight: 500;
}

.header-user:hover {
    color: #666;
    cursor: pointer;
}

.header-user {
  float: right;
  margin: 10px 20px 0 0;
}

.header-menu-close {
  float: right;
  cursor: pointer;
  padding: 0 4px 0 4px;
}

.header-menu-close:hover {
  background: #000;
  color: #fff;
  border-radius: 50%;
}

.header-menu-list {
  padding-top: 15px;
  padding: 0 10px 10px 10px;
}

.header-menu-list a {
  font-size: 23px;
  color: #202020;
}

.header-top {
  height: 60px;
  z-index: 999;
}

.header-top h1 {
    font-size: 1.75rem;
    margin-left: 80px;
    display: inline;
    line-height: 50px;
}
.header-user-menu-links a:hover {
    background: #212529;
    color: #fff;
    text-decoration: none;
}

.header-user-menu-links a {
    color: #000;
    display: block;
    padding: 2px 11px;
}

.sidebar-left {
  background-color: #1d212a;
}

aside.sidebar-mini {
  width: 50px;
  -webkit-transition: all 0.1s ease-in-out;
  -moz-transition: all 0.1s ease-in-out;
  -o-transition: all 0.1s ease-in-out;
  -ms-transition: all 0.1s ease-in-out;
  transition: all 0.1s ease-in-out;
}

.sidebar {
  display: block;
}

.sidebar {
  position: fixed;
  top: 0;
  min-height: 100%;
  width: 240px;
  -webkit-transition: all 0.1s ease-in-out;
  -moz-transition: all 0.1s ease-in-out;
  -o-transition: all 0.1s ease-in-out;
  -ms-transition: all 0.1s ease-in-out;
  transition: all 0.1s ease-in-out;
}

article, aside, details, figcaption, figure, footer, header, hgroup, main, menu, nav, section, summary {
  display: block;
}

aside.sidebar-mini .sidebar-header, aside.sidebar-mini .sidebar-profile, aside.sidebar-mini .sidebar-summary {
  display: none;
}

.sidebar-header {
  color: #6F737E;
  font-weight: 600;
  line-height: 20px;
  margin: 0;
  padding: 10px 10px 5px;
  text-transform: uppercase;
}

h5, .h5 {
  font-size: 0.813em;
}

h1, h2, h3, h4, h5, h6, .h1, .h2, .h3, .h4, .h5, .h6 {
  font-family: "Montserrat", sans-serif;
  color: #424856;
  margin: 0;
}

aside.sidebar-mini .nav-pills {
  margin-left: 5px;
  margin-right: 5px;
}

.sidebar .nav-pills {
  margin-left: 10px;
  margin-right: 10px;
}

.nav {
  padding-left: 0;
  margin-bottom: 0;
  list-style: none;
}

.nav-stacked > li {
  float: none;
}

.nav-pills > li {
  float: left;
}

.nav > li {
  position: relative;
  display: block;
}

aside.sidebar-mini .nav > li > a {
  padding: 10px 9px;
  white-space: nowrap;
  overflow: hidden;
}

.sidebar .nav-pills > li > a {
  padding: 9px 10px;
  font-size: 0.875em;
}

.sidebar li.nav-item.vector {
  font-size: 21px;
}

.sidebar .nav a {
  text-decoration: none;
  font-weight: 600;
}

.sidebar .nav a {
  padding: 0.75rem 0.7rem;
}

.sidebar .nav a span {
  margin-right: 10px;
}

.sidebar .nav a span img {
  width: 22px;
}

.sidebar .nav a svg {

}

.sidebar .nav a .link-text
{
  font-size: 12px;
  display: inline-block;
  margin-left: 5px;
}

.nav-pills > li > a {
  border-radius: 4px;
}

.nav > li > a {
  position: relative;
  display: block;
  padding: 10px 15px;
  width: 100%;
}

.sidebar-left a {
  /*color: #B3B8C3;*/
  color: #fff;
}

aside.sidebar-mini .nav-pills > li.nav-dropdown-open, aside.sidebar-mini .nav-pills > li:hover {
  width: 235px;
}

aside.sidebar-mini .nav > li.nav-dropdown-open, aside.sidebar-mini .nav > li:hover {
  width: 240px;
  z-index: 10;
}

.nav-stacked > li + li {
  margin-top: 2px;
  margin-left: 0;
}

.nav-pills > li + li {
  margin-left: 2px;
}

.nav-stacked > li {
  float: none;
}

.nav-pills > li {
  float: left;
}

.nav > li {
  position: relative;
  display: contents;
}

aside.sidebar-mini .nav > li {
  display: block;
}

aside.sidebar-mini .nav > li.nav-dropdown-open > a, aside.sidebar-mini .nav > li:hover > a {
  color: #1d2939;
  background-color: #fff;
}

.sidebar-left .nav > li.open > a, .sidebar-left .nav > li > a:hover {
  color: #1d2939;
  background-color: #ffffff;
}

aside.sidebar-mini .nav > li.nav-dropdown-open a, aside.sidebar-mini .nav > li:hover a {
  display: block;
  overflow: visible;
}

.sidebar .nav-pills > li > a {
  padding: 9px 10px;
  font-size: 0.875em;
}

aside.sidebar-mini .nav-pills > li > a > svg {
  margin-right: 15px;
  width: 23px;
}

aside.sidebar-mini .nav > li.nav-dropdown-open a, aside.sidebar-mini .nav > li:hover a {
  /*display: block;*/
  overflow: visible;
}

aside.sidebar-mini .nav > li > a {
  padding: 10px 9px;
  white-space: nowrap;
  overflow: hidden;
}

.sidebar .nav-pills > li > a {
  padding: 9px 10px;
  font-size: 0.875em;
}

.nav > li > a:hover, .nav > li > a:focus {
  text-decoration: none;
  background-color: #eee;
}

.sidebar-left a:focus, .sidebar-left a:hover {
  color: #fff;
  background-color: transparent;
}

.sidebar .nav a {
  text-decoration: none;
  font-weight: 600;
}

aside.sidebar-mini .nav > li > a {
  display: flex;
}

.nav > li:after {
  content: '';
  width: 100%;
  order: 0;
}

aside.sidebar:hover {
  overflow: visible;
}

aside.sidebar {
  overflow: hidden;
}

.header-top h1 {
  font-size: 1.75rem;
  margin-left: 80px;
  display: inline;
  line-height: 50px;
}

</style>
