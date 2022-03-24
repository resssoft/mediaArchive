<template>
  <div class="page">
    <Header/>
    <div v-show="isDataVisible" class="container dashboard card">
      <div class="card-header">
        Servers
        <span class="other-buttons">
          <button @click="$modal.show('createItem')" class="btn btn-primary">New</button>
        </span>
        <span class="special-buttons">
          <button @click="runDiscoverHV()" :class="'btn btn-default run-btn ' + DiscoverHV.status ">Discover HV</button>
        </span>
      </div>
      <div class="card-body">
        <!--<div class="alert alert-danger alert-dismissable" v-show="error != ''">
          <strong>Error! </strong>
          <div class="alert-text" style="display: inline-block;"> {{ error }}</div>
        </div>-->
        <div class="alert alert-info alert-dismissable" v-show="info != ''">
          <strong>
            <b-icon icon="info-circle-fill"></b-icon>
          </strong>
          <div class="alert-text" style="display: inline-block;"> {{ info }}</div>
        </div>
        <div v-bind:class="loading"></div>
        <Filters
            :onChanged="filterChanged"
            v-on:onChanged="filterChanged"
            :filter="filter"
            :filterValues="filterValues"
        />
        <table class="table table-responsive-sm" id="resizable-table">
          <thead>
          <tr>
            <th scope="col" style="max-width: 40px">
              <div class="flex">
                <span class="col-text" @click="SortInvert('Id')">Id</span>
                <span class="icon-up" v-show="sortField=='Id'"><b-icon icon="arrow-up"></b-icon></span>
                <span class="icon-down" v-show="sortField==`-Id`"><b-icon icon="arrow-down"></b-icon></span>
              </div>
              <div class="resizer" @mousedown="resizerMousedown"></div>
            </th>
            <th scope="col" style="max-width: 150px">
              <div class="flex">
                <span class="col-text" @click="SortInvert('Name')">Name</span>
                <span class="icon-up" v-show="sortField=='Name'"><b-icon icon="arrow-up"></b-icon></span>
                <span class="icon-down" v-show="sortField==`-Name`"><b-icon icon="arrow-down"></b-icon></span>
              </div>
              <div class="resizer" @mousedown="resizerMousedown"></div>
            </th>
            <th scope="col" style="max-width: 100px">
              <div class="flex">
                <span class="col-text" @click="SortInvert('Ip')">Ip</span>
                <span class="icon-up" v-show="sortField=='Ip'"><b-icon icon="arrow-up"></b-icon></span>
                <span class="icon-down" v-show="sortField==`-Ip`"><b-icon icon="arrow-down"></b-icon></span>
              </div>
              <div class="resizer" @mousedown="resizerMousedown"></div>
            </th>
            <th scope="col">
              <div class="flex">
                <span class="col-text" @click="SortInvert('Cpu')">Cpu/Ram/VmRam/Ssd/Hdd</span>
                <span class="icon-up" v-show="sortField=='Cpu'"><b-icon icon="arrow-up"></b-icon></span>
                <span class="icon-down" v-show="sortField==`-Cpu`"><b-icon icon="arrow-down"></b-icon></span>
              </div>
              <div class="resizer" @mousedown="resizerMousedown"></div>
            </th>
            <th scope="col">
              <div class="flex">
                <span class="col-text" >CRUD</span>
              </div>
              <div class="resizer" @mousedown="resizerMousedown"></div>
            </th>
          </tr>
          </thead>
          <tbody>
            <tr v-for="item in items" v-bind:key="item.id">
              <td class="text-left"><span class="cel-text"> {{ item.Id }} </span></td>
              <td class="text-left">
                <span class="cel-text">
                  <span :class="'status status-'+item.State"></span>
                  <span class='server-name'>{{ item.Name }}</span>
                  <br><span class='server-model'>{{ item.Model }}</span>
                </span></td>
              <td class="text-left"><span class="cel-text"> {{ item.Ip }} </span></td>
              <td class="text-left">
                <span class="cel-text">
                  <div class="bar-wrapper">
                    <ul :id="'item-'+item.Id+'-cpu-progress'" class='hardware-progress-bar cpu-progress'></ul>
                    <span class=''>Cpu: {{ item.Cpu }}/{{ item.MaxCpu }}</span><br>
                  </div>
                  <div class="bar-wrapper">
                    <ul :id="'item-'+item.Id+'-ram-progress'" class='hardware-progress-bar ram-progress'></ul>
                    <span class=''>Ram: {{ item.Ram }}/{{ item.MaxRam }}</span><br>
                  </div>
                  <div class="bar-wrapper">
                    <ul :id="'item-'+item.Id+'-free-ram-progress'" class='hardware-progress-bar free-ram-progress'></ul>
                    <span class=''>VmRam: {{ item.VmSum }}/{{ item.Ram }}</span><br>
                  </div>
                  <div class="bar-wrapper">
                    <ul :id="'item-'+item.Id+'-disk-progress'" class='hardware-progress-bar disk-progress'></ul>
                    <span class=''>Ssd/Hdd: {{ item.Ssd }}/{{ item.Hdd }}</span>
                  </div>
                </span>
              </td>
              <td class="text-left">
                <span class="cel-buttons">
                  <button class="btn btn-default icon-only" @click="editItem(item)">
                    <i class="fa fa-pencil"></i>
                  </button>
                  <button class="btn btn-default delete-btn small icon-only" @click="openDeleteItemConfirmation(item)">
                    <i class="fa fa-trash"></i>
                  </button>
                </span>
              </td>
            </tr>
          </tbody>
        </table>
        <ListControls
            v-on:onLimitChanged="ChangeLimit"
            :total="total"
            :options="options"
        />
      </div>
    </div>

    <ModalItem
        type="edit"
        :itemId="currentItemId"
        :modalName="'editItem'"
        v-on:onSuccess="itemSuccess"/>

    <ModalItem
        type="create"
        :itemId="0"
        :modalName="'createItem'"
        v-on:onSuccess="itemSuccess"/>

    <ModalConfirmItemDelete
        v-on:onComplete="deletedItem"
        v-on:onError="deletedItemErr"
        :url="apiUrl"
        :itemId="currentItemId"
    />

    <ModalLog :data="textareaData" :modalName="'modalTextarea'"/>

    <ModalError name="listError" :message="error"/>

    <ModalSuccess name="listSuccess" :message="success"/>

  </div>
</template>
<script>

import axios from 'axios';

// mixins
import {HwProgressBarMixin} from '../assets/js/mixins/hw-progress-bar-mixin';
import {UserStoredDataMixin} from '../assets/js/mixins/user-stored-data-mixin';
import {ListLogicMixin} from '../assets/js/mixins/list-logic-mixin';

export default {
  name: 'listPage',
  mixins: [
    HwProgressBarMixin,
    UserStoredDataMixin,
    ListLogicMixin
  ],
  components: {

  },
  data() {
    return {
      title: '...',
      items: null,
      currentItemId: 0,
      apiList: '/api/manager/esxi/list',
      apiItem: '',
      apiUrl: '',
      apiConfig: '/api/config/list/',
      columns: [
        {name: "Id", value: "Id"},
        {name: "Name", value: "Name"},
        {name: "Ip", value: "Ip"},
        {name: "Cpu/Ram/VmRam/Ssd/Hdd", value: "Cpu"},
        {name: "CRUD", value: "CRUD"},
      ],
      options: [],
      DiscoverHV: {
        api: "/api/launcher/start",
        apiStatus: "/api/launcher/status",
        status: "",
        name: "DiscoverHV",
        timer: true,
      },
      textareaData: '',
    };
  },
  methods: {

    LoadConfig() {
      if (this.$route.query.page != undefined) {
        this.pageNumber = this.$route.query.page;
      } else {
        this.pageNumber = 1;
      }

      this.title = '...';
      this.filter = [];
      this.items = null;
      this.error = '';
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        return;
      }
      this.loading = 'loading';
      let pageName = this.$route.name;
      axios.get(this.apiConfig + pageName,
          {headers: {'Authorization': 'Bearer ' + token.accessToken}})
          .catch((err) => {
            this.loading = 'loaded';
            this.errorHandle(err);
          })
          .then(res => {
            this.loading = 'loaded';
            console.log('loaded config');
            if (res !== null && res !== undefined) {
              if (res.data !== null && res.data !== undefined) {
                this.title = res.data.title;
                this.filter = res.data.filter;
                this.filters = [];
                this.filterValues = res.data.filterValues;
                this.apiList = res.data.apiList;
                this.apiItem = res.data.apiItem;
                this.apiUrl = res.data.apiUrl;
                this.sortField = res.data.sortField;
                this.columns = res.data.columns;
                this.options = res.data.options;
                this.limit = res.data.limit;
                this.ApplyUserConfig();
                this.filterPreparing();
                this.LoadData();
              }
            }
          });
    },

    LoadData() {
      let loader = this.$loading.show();

      this.error = '';
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        return;
      }
      this.loading = 'loading'
      axios.post(this.apiList + '?limit=' + this.limit + '&sort=' + this.sortField + '&page=' + this.pageNumber,
          this.filters,
          {headers: {'Authorization': 'Bearer ' + token.accessToken}})
          .catch((err) => {
            loader.hide();
            this.loading = 'loaded';
            if (err.response) {
              if (err.response.status === 401) {
                this.$router.push({path: this.login});
              }
              this.error = '[' + err.response.status + ' ' + err.response.statusText + '] ';
              if (err.response.data != '') {
                this.error += err.response.data;
              }
            } else {
              this.error = 'Server error:  no connection'
            }
          })
          .then(res => {
            loader.hide();
            console.log(res)
            this.loading = 'loaded';
            console.log('loaded data');
            this.items = res.data.data;
            this.total = res.data.total;
            //setTimeout(this.checkDiscoverStatus, 2000);
            setTimeout(this.fillHwProgressBars, 10);
          })
    },

    runDiscoverHV() {
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        return;
      }
      this.error = '';
      this.DiscoverHV.status = 'running';
      axios.post(this.DiscoverHV.api,
          {
            name: 'DiscoverHV'
          },
          {headers: {'Authorization': 'Bearer ' + token.accessToken}})
          .catch((err) => {
            this.errorHandle(err);
            this.DiscoverHV.status = '';
          })
          .then(res => {
            if (res !== null && res !== undefined) {
              this.info = 'Running';
              setTimeout(this.clearInfo, 5000);
              setTimeout(this.checkDiscoverStatus, 2000);
            }
          })
    },

    checkDiscoverStatus() {
      if (this.DiscoverHV.status === "running") {
        const token = JSON.parse(localStorage.getItem('token'));
        if (token == null) {
          this.error = 'token not found';
          return;
        }
        axios.post(this.DiscoverHV.apiStatus,
            {
              name: this.DiscoverHV.name
            },
            {headers: {'Authorization': 'Bearer ' + token.accessToken}})
            .catch((err) => {
              if (err.response) {
                if (err.response.status === 401) {
                  this.$router.push({path: this.login});
                }
                console.log('check Buttons Status: ' + '[' + err.response.status + ' ' + err.response.statusText + '] ');
              } else {
                console.log('check Buttons Status: server error:  no connection');
              }

              this.DiscoverHV.status = '';
              this.errorHandle(err);
            })
            .then(res => {
              if (res !== null && res !== undefined) {
                if (res.data != "") {
                  if (res.data.data !== "processed") {
                    this.DiscoverHV.status = "";

                    this.getAndShowDiscoverLogs();
                    this.LoadData();
                    this.success = 'Script successfully finished!';
                    this.$modal.show('success');
                  } else {
                    setTimeout(this.checkDiscoverStatus, 2000);
                  }
                }
              }
            });
      }
    },

    editItem(item) {
      //this.$router.push({path: this.apiItem + item.Id});
      this.currentItemId = item.Id;
      this.$modal.show('editItem');
    },

    fillHwProgressBars() {
      this.items.forEach((item) => {
        this.fillCpuProgressBar(item.Id, item.MinCpu, item.Cpu, item.MaxCpu);
        this.fillRamProgressBar(item.Id, item.MinRam, item.Ram, item.MaxRam);
        this.fillFreeRamProgressBar(item.Id, 0, item.VmSum, item.Ram);
        this.fillDiskProgressBar(item.Id, 0, item.Ssd, item.Hdd);
      });

      let spans = document.querySelectorAll('.hardware-progress-bar + span');
      spans.forEach((span) => {
        span.style.display = 'inline';
      });
    },

    getAndShowDiscoverLogs() {
      let loader = this.$loading.show();
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        this.moveUp();
        return;
      }
      this.loading = 'loading';
      axios.post('/api/log/list?limit=20&sort=-id&page=1',
          [{"Condition": "ilike", "Data": {"Entity": "DiscoverHV"}}],
          {headers: {'Authorization': 'Bearer ' + token.accessToken}})
          .catch((err) => {
            loader.hide();
            this.loading = 'loaded';
            this.errorHandle(err);
          })
          .then(res => {
            loader.hide();
            this.loading = 'loaded';
            if (res.data.data !== undefined) {
              let last_log = res.data.data[0];
              this.textareaData = last_log.Data;
              this.$modal.show('modalTextarea');
            } else {
              this.errorHandle("Something wrong, no logs found.");
            }
          });
    },

  },
  computed: {

  },
  watch: {
    '$route': 'onRouteChange'
  },
  created: function () {
    this.LoadConfig();
  },
  mounted() {
    this.$nextTick(function () {
      //this.initResizable();
    });
  },
  updated: function () {
    //this.columnsReset();
    //this.updateResizable();
  }
}
</script>
