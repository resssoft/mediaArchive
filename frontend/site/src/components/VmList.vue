<template>
  <div class="page">
    <Header/>
    <div v-show="isDataVisible" class="container dashboard card">
      <div class="card-header">
        Virtual Machines

        <span class="special-buttons">
          <button @click="$modal.show('confirmVmDeploy')" :class="'btn btn-default run-btn ' + DeployVM.status ">Deploy VM</button>
          <button @click="$modal.show('confirmVmDestroy')" :class="'btn btn-default run-btn run-btn-destroy ' + DestroyVM.status ">Destroy VM</button>
          <button @click="runDiscoverVM()" :class="'btn btn-default run-btn ' + DiscoverVM.status ">Discover VM</button>
        </span>
        <span class="other-buttons">
          <button @click="showCreatedNodes()" class="btn btn-primary">{{buttonFilterCreatedText}}</button>
          <button @click="showChangedNodes()" class="btn btn-primary">{{buttonFilterChangedText}}</button>
          <button @click="showDestroyedNodes()" class="btn btn-primary">{{buttonFilterDestroyedText}}</button>
          <button @click="$modal.show('VmNewNode')" class="btn btn-primary">New node</button>

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
            <th
                v-for="column in columns"
                v-bind:key="column.value"
                scope="col"
            >
              <div class="flex">
                <span class="col-text" @click="SortInvert(column.value)">{{ column.name }}</span>
                <span class="icon-up" v-show="sortField==column.value"><b-icon icon="arrow-up"></b-icon></span>
                <span class="icon-down" v-show="sortField=='-'+column.value"><b-icon icon="arrow-down"></b-icon></span>
              </div>
              <div class="resizer" @mousedown="resizerMousedown"></div>
            </th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="item in items" v-bind:key="item.id">
            <td v-for="column in columns"
                v-bind:key="column.value"
                :class="(column.name == 'State' ? 'text-center' : 'text-left')"
            >
              <span v-if="column.name == 'CRUD'" class="cel-buttons">
                <button class="btn btn-default icon-only" @click="editItem(item)">
                  <i class="fa fa-pencil"></i>
                </button>
                <button class="btn btn-default delete-btn small icon-only" @click="openDeleteItemConfirmation(item)">
                  <i class="fa fa-trash"></i>
                </button>
              </span>
              <span v-else-if="column.name == 'State'" class="cel-text" :title="item[column.value]">
                <span v-if="item[column.value] == 'poweredOn'" class="vm-state powered-on"></span>
                <span v-if="item[column.value] == 'poweredOff'" class="vm-state powered-off"></span>
                <span v-if="item[column.value] == 'CREATED'" class="vm-state created">🚀</span>
                <span v-if="item[column.value] == 'CHANGED'" class="vm-state created"><i class="fa fa-refresh"></i></span>
                <span v-if="item[column.value] == 'MOVE_PENDING'" class="vm-state created"><i class="fa fa-arrow-right"></i></span>
                <span v-if="item[column.value] == 'DESTROYED'" class="vm-state created"><i class="fa fa-trash"></i></span>
                <span v-if="item[column.value] == 'NONE'" class="vm-state created"><i class="fa fa-circle-o"></i></span>
              </span>
              <span v-else class="cel-text">
                {{ item[column.value] }}
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

    <VmNewNode successCreateCallback="loadData()"/>

    <ModalConfirmVmDeploy v-on:onConfirm="runDeployCreatedNodes" />

    <ModalConfirmVmDestroy v-on:onConfirm="runDestroyNodes" />

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
import {UserStoredDataMixin} from '../assets/js/mixins/user-stored-data-mixin';
import {ListLogicMixin} from '../assets/js/mixins/list-logic-mixin';
import VmNewNode from "./VmNewNode";
import ModalConfirmVmDeploy from "./ModalConfirmVmDeploy";
import ModalConfirmVmDestroy from "./ModalConfirmVmDestroy";

export default {
  name: 'listPage',
  mixins: [
    UserStoredDataMixin,
    ListLogicMixin
  ],
  components: {
    ModalConfirmVmDeploy,
    ModalConfirmVmDestroy,
    VmNewNode
  },
  data() {
    return {
      title: '...',
      items: null,
      currentItemId: 0,
      apiList: '/api/manager/vm/list',
      apiItem: '',
      apiUrl: '',
      apiConfig: '/api/config/list/',
      buttonFilterCreatedText: 'Show CREATED',
      buttonFilterCreatedUsed: false,
      buttonFilterChangedText: 'Show CHANGED',
      buttonFilterChangedUsed: false,
      buttonFilterDestroyedText: 'Show DESTROYED',
      buttonFilterDestroyedUsed: false,
      columns: [
        {name: 'Id', value: 'Id'},
        {name: 'Server', value: 'Server'},
        {name: 'Prev. server', value: 'Prev. server'},
        {name: 'Name', value: 'Name'},
        {name: 'Ip', value: 'Ip'},
        {name: 'State', value: 'State'},
        {name: 'Mac', value: 'Mac'},
        {name: 'Ram', value: 'Ram'},
        {name: 'Cpu', value: 'Cpu'},
        {name: 'Comment', value: 'Comment'},
        {name: 'ProjectId', value: 'ProjectId'},
        {name: 'DiscoveredAt', value: 'DiscoveredAt'},
        {name: 'CRUD', value: 'CRUD'},
      ],
      options: [],
      DiscoverVM: {
        api: "/api/launcher/start",
        apiStatus: "/api/launcher/status",
        status: "",
        name: "DiscoverVM",
        timer: true,
      },
      DeployVM: {
        api: "/api/manager/vm/deploy/?custom=y",
        apiStatus: "/api/launcher/status",
        status: "",
        name: "deployVm",
        timer: true,
      },
      DestroyVM: {
        api: "/api/manager/vm/destroy/?custom=y",
        apiStatus: "/api/launcher/status",
        status: "",
        name: "destroyVm",
        timer: true,
      },
      textareaData: '',
    };
  },
  methods: {
    deployErr(err) {
      this.errorHandle(err);
    },
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
                this.buttons = res.data.buttons;
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
          })
    },
    runDiscoverVM() {
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        return;
      }
      this.error = '';
      this.DiscoverVM.status = 'running';
      axios.post(this.DiscoverVM.api,
          {
            name: 'DiscoverVM'
          },
          {headers: {'Authorization': 'Bearer ' + token.accessToken}})
          .catch((err) => {
            this.errorHandle(err);
            this.DiscoverVM.status = '';
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
      if (this.DiscoverVM.status === "running") {
        const token = JSON.parse(localStorage.getItem('token'));
        if (token == null) {
          this.error = 'token not found';
          return;
        }
        axios.post(this.DiscoverVM.apiStatus,
            {
              name: this.DiscoverVM.name
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

              this.DiscoverVM.status = '';
              this.errorHandle(err);
            })
            .then(res => {
              if (res !== null && res !== undefined) {
                if (res.data != "") {
                  if (res.data.data !== "processed") {
                    this.DiscoverVM.status = "";

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
    showCreatedNodes() {
      this.buttonFilterCreatedUsed = !this.buttonFilterCreatedUsed;
      if (this.buttonFilterCreatedUsed === true) {
        this.buttonFilterCreatedText = "Show all"
      } else {
        this.buttonFilterCreatedText = "Show CREATED";
        this.filters = [];
        this.LoadData();
        return;
      }
      this.sortField = "State";
      let data = {};
      data["State"] = "CREATED";
      this.filters.push({
        Condition: "ilike",
        Data: data,
      });
      this.UpdateUserConfig();
      this.LoadData();
    },
    showChangedNodes() {
      this.buttonFilterChangedUsed = !this.buttonFilterChangedUsed;
      if (this.buttonFilterChangedUsed === true) {
        this.buttonFilterChangedText = "Show all"
      } else {
        this.buttonFilterChangedText = "Show CHANGED";
        this.filters = [];
        this.LoadData();
        return;
      }
      this.sortField = "State";
      let data = {};
      data["State"] = "CHANGED";
      this.filters.push({
        Condition: "ilike",
        Data: data,
      });
      this.UpdateUserConfig();
      this.LoadData();
    },
    showDestroyedNodes() {
      this.buttonFilterDestroyedUsed = !this.buttonFilterDestroyedUsed;
      if (this.buttonFilterDestroyedUsed === true) {
        this.buttonFilterDestroyedText = "Show all"
      } else {
        this.buttonFilterDestroyedText = "Show DESTROYED";
        this.filters = [];
        this.LoadData();
        return;
      }
      this.sortField = "State";
      let data = {};
      data["State"] = "DESTROYED";
      this.filters.push({
        Condition: "ilike",
        Data: data,
      });
      this.UpdateUserConfig();
      this.LoadData();
    },
    runDeployCreatedNodes(unavailable, recreate, changed, moved) {
      this.$modal.hide('confirmVmDeploy');
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        return;
      }
      this.error = '';
      this.DeployVM.status = 'running';
      let deployApiUrl = this.DeployVM.api;
      if (unavailable === true) {
        deployApiUrl += "&unavailable=y"
      }
      if (recreate === true) {
        deployApiUrl += "&recreate=y"
      }
      if (changed === true) {
        deployApiUrl += "&changed=y"
      }
      if (moved === true) {
        deployApiUrl += "&moved=y"
      }
      console.log(unavailable, recreate, changed, moved, deployApiUrl);
      axios.post(deployApiUrl,
          [],
          {headers: {'Authorization': 'Bearer ' + token.accessToken}})
          .catch((err) => {
            this.errorHandle(err);
            this.DeployVM.status = '';
          })
          .then(res => {
            if (res !== null && res !== undefined) {
              //this.info = 'Running';
              //setTimeout(this.clearInfo, 5000);
              setTimeout(this.checkDeployStatus, 2000);
            }
          })
    },
    runDestroyNodes() {
      this.$modal.hide('confirmVmDestroy');
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        return;
      }
      this.error = '';
      this.DestroyVM.status = 'running';
      let deployApiUrl = this.DestroyVM.api;
      axios.post(deployApiUrl,
              [],
              {headers: {'Authorization': 'Bearer ' + token.accessToken}})
              .catch((err) => {
                this.errorHandle(err);
                this.DestroyVM.status = '';
              })
              .then(res => {
                if (res !== null && res !== undefined) {
                  setTimeout(this.checkDestroyStatus, 2000);
                }
              })
    },
    checkDeployStatus() {
      if (this.DeployVM.status === "running") {
        const token = JSON.parse(localStorage.getItem('token'));
        if (token == null) {
          this.error = 'token not found';
          return;
        }
        axios.post(this.DeployVM.apiStatus,
            {
              name: this.DeployVM.name
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

              this.DeployVM.status = '';
              this.errorHandle(err);
            })
            .then(res => {
              if (res !== null && res !== undefined) {
                if (res.data != "") {
                  if (res.data.data !== "processed") {
                    this.DeployVM.status = "";
                    this.getAndShowLogsBy("vm-deploy");
                    this.LoadData();
                    this.success = 'Script successfully finished!';
                    this.$modal.show('VmListSuccess');
                  } else {
                    setTimeout(this.checkDeployStatus, 2000);
                  }
                }
              }
            });
      }
    },
    checkDestroyStatus() {
      if (this.DestroyVM.status === "running") {
        const token = JSON.parse(localStorage.getItem('token'));
        if (token == null) {
          this.error = 'token not found';
          return;
        }
        axios.post(this.DestroyVM.apiStatus,
                {
                  name: this.DestroyVM.name
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

                  this.DestroyVM.status = '';
                  this.errorHandle(err);
                })
                .then(res => {
                  if (res !== null && res !== undefined) {
                    if (res.data != "") {
                      if (res.data.data !== "processed") {
                        this.DestroyVM.status = "";
                        this.getAndShowLogsBy("vm-destroy");
                        this.LoadData();
                        this.success = 'Script successfully finished!';
                        this.$modal.show('VmListSuccess');
                      } else {
                        setTimeout(this.checkDestroyStatus, 2000);
                      }
                    }
                  }
                });
      }
    },
    editItem(item) {
      //this.$router.push({path: this.apiItem + item.Id});
      this.currentItemId = item.Id;
      this.currentItemName = item.Name;
      this.$modal.show('editItem');
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
          [{"Condition": "ilike", "Data": {"Entity": "DiscoverVM"}}],
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

    getAndShowLogsBy(entityName) {
      let loader = this.$loading.show();
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        this.moveUp();
        return;
      }
      this.loading = 'loading';
      axios.post('/api/log/list?limit=20&sort=-id&page=1',
              [{"Condition": "ilike", "Data": {"Entity": entityName}}],
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
