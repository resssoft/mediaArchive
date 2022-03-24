<template>
  <div class="page">
    <Header/>
    <div v-show="isDataVisible" class="container dashboard card">
      <div class="card-header">
        {{ title }}
        <span class="other-buttons">
              <button @click="$modal.show('createItem')" class="btn btn-primary">New</button>
            </span>
        <span class="special-buttons" v-for="(filterField, filterIndex) in buttons" v-bind:key="filterIndex">
              <button @click="button(filterIndex)" class="btn btn-primary run-btn"
                      :class="filterField.class"
                      v-if="filterField.name !== '' || filterField.name !== undefined"
              >
                  {{ filterField.name }}
              </button>
            </span>
      </div>
      <div class="card-body">
        <div class="alert alert-info alert-dismissable" v-show="info != ''">
          <strong>
            <b-icon icon="info-circle-fill"></b-icon>
          </strong>
          <div class="alert-text" style="display: inline-block;"> {{ info }}</div>
        </div>
        <div v-bind:class="loading"></div>
        <Filters
            v-on:onChanged="filterChanged"
            :filter="filter"
            :filterValues="filterValues"
        />
        <table class="table table-responsive-sm" id="resizable-table">
          <thead>
          <tr>
            <th v-for="column in columns" v-bind:key="column.value" scope="col">
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
          <tr v-for="item in items" v-bind:key="item.ID">
            <td v-for="column in columns"
                v-bind:key="column.value"
                @click="rowCellClick(column.value, item[column.value])">
              <span v-if="column.buttons !== undefined" class="cel-buttons">
                <button v-if="column.buttons.includes('edit')" class="btn btn-default icon-only"
                        @click="editItem(item)">
                  <i class="fa fa-pencil"></i>
                </button>
                <button v-if="column.buttons.includes('delete')" class="btn btn-default delete-btn small icon-only"
                        @click="openDeleteItemConfirmation(item)">
                  <i class="fa fa-trash"></i>
                </button>
              </span>
              <span v-else-if="column.template !== undefined" class="cel-text"
                    v-html="fillTemplate(column.template.value, item)"></span>
              <span v-else class="cel-text">{{ item[column.value] }}</span>
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
        :itemId=currentItemId
        :modalName="'editItem'"
        v-on:onSuccess="itemSuccess"/>
    <ModalItem
        type="create"
        :itemId="0"
        :modalName="'createItem'"
        v-on:onSuccess="itemSuccess"/>
    <ModalLog :data="textareaData" :modalName="'modalTextarea'"/>
    <ModalError name="listError" :message="error"/>
    <ModalSuccess name="listSuccess" :message="success"/>
    <ModalConfirmItemDelete
        v-on:onComplete="deletedItem"
        v-on:onError="deletedItemErr"
        :url="apiUrl"
        :itemId="currentItemId"
    />
  </div>
</template>
<script>

import axios from "axios";
import {ListLogicMixin} from '../assets/js/mixins/list-logic-mixin';

export default {
  name: 'listPage',
  mixins: [
    ListLogicMixin
  ],
  components: {

  },
  data() {
    return {
      title: '...',
      items: null,
      currentItemId: 0,
      apiList: '',
      apiItem: '',
      apiUrl: '',
      apiConfig: '/api/config/list/',
      columns: [],
      options: [],
      textareaData: '',
      dataKeysForStore:['limit', 'sortField'],
    };
  },
  methods: {

    rowCellClick(cellName, val) {
      if (cellName.toLowerCase() === 'id') {
        this.$router.push({path: this.apiItem + val})
      }
    },
    editItem(item) {
      //this.$router.push({path: this.apiItem + item.Id});
      this.currentItemId = item.ID; //parseInt(item.Id);
      this.$modal.show('editItem');
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
      let pageName = this.$route.params.page;
      console.log(this.$route.name, this.$route.params.page);
      if (pageName == '') {
        pageName = this.$route.name;
      }
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
            this.errorHandle(err);
          })
          .then(res => {
            loader.hide();
            this.loading = 'loaded';
            console.log('loaded data');
            this.items = res.data.data;
            this.total = res.data.total;
            setTimeout(this.checkButtonsStatus, 2000);
          })
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
          [{"Condition":"ilike","Data":{"Entity":"DiscoverHV"}}],
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
            }
            else{
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
};
</script>
