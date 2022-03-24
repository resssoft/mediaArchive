<template>
  <div class="page">
    <Header
        :title=actionText
    />
    <div class="container dashboard card">
      <div class="card-body">
        <div class="alert alert-danger alert-dismissable" v-show="error !== ''">
          <strong>Error </strong>
          <div class="errorMsg" style="display: inline-block;">{{ error }}</div>
        </div>
        <button type="button" class="btn btn-primary delete-btn" v-if="action === 'update'" @click="deleteItem()">
          Delete
        </button>
        <button type="button" @click="returnBackToList()" class="btn btn-default">Back</button>
        <form ref="form" class="form-horizontal" @submit="submitForm">
          <span v-for="(field, index) in formData" v-bind:key="index">
          <div class="form-group" v-if="formData[index].show &&
        ([].concat(formData[index].form).includes(action) || formData[index].form === undefined)">
            <label class="col-sm-2 control-label">{{ field.name }}</label>
            <div class="col-sm-10">
              <input
                      @blur="validate(index, formData[index].value, formData[index].rule)"
                      :type="formData[index].type"
                      class="form-control"
                      v-model="formData[index].value"
                      v-if="formData[index].show && ['text', 'number', 'email'].includes(formData[index].type)"
                      :disabled="field.name === 'Id'"
                      v-bind:class="formData[index].class">
              <b-form-select
                      v-if="formData[index].type === 'select'"
                      v-model="formData[index].value"
                      :options="formData[index].options"></b-form-select>
              <b-form-textarea
                      class="col-sm-11 log-edit"
                      v-if="formData[index].type === 'textarea'"
                      v-model="formData[index].value"></b-form-textarea>
              <v-jsoneditor
                      v-if="formData[index].type === 'json'"
                      v-model="jsonData[index]"
                      :plus="false"
                      :options="jsonOptions"></v-jsoneditor>
              <button
                      type="button"
                      v-if="formData[index].type === 'textarea' && formData[index].name !== 'Import'"
                      class="btn btn-default col-sm-1"
                      @click="$modal.show('showTextareaData')"
              >Open
              </button>
            </div>
          </div>
          </span>
          <div class="form-group">
            <div class="col-sm-offset-2 col-sm-10">
              <button type="submit" class="btn btn-primary">{{ actionText }}</button>
            </div>
          </div>
        </form>
      </div>
    </div>

    <modal
        name="showTextareaData"
        :scrollable="true"
        height="80%"
        width="60%"
        class="logs-modal"
    >
      <div class="window-header">
        Modal
        <a href="javascript:void(0)" class="close-modal pull-right" @click="$modal.hide('showTextareaData')"><i
            class="fa fa-times"></i></a>
      </div>
      <div class="modal-body" style="overflow: scroll; height: calc(100% - 46px)">
        <pre>{{ textareaData }}</pre>
      </div>
    </modal>

    <modal
            name="notifyData"
            :scrollable="true"
            height="80%"
            width="60%"
            class="logs-modal"
    >
      <div class="window-header">
        Modal
        <a href="javascript:void(0)" class="close-modal pull-right" @click="$modal.hide('notifyData')"><i
                class="fa fa-times"></i></a>
      </div>
      <div class="modal-body" style="overflow: scroll; height: calc(100% - 46px)">
        <pre>{{ notifyData }}</pre>
      </div>
    </modal>
  </div>
</template>

<script>
import Header from './Header.vue';
import axios from 'axios';
import VJsoneditor from 'v-jsoneditor/src/index';

export default {
  name: 'newItem',
  components: {
    Header,
    VJsoneditor
  },
  data() {
    return {
      error: '',
      actionText: '...',
      action: 'add',
      apiUrl: '',
      redirectToUrl: '',
      apiConfig: '/api/config/item-form/',
      login: '/login',
      loading: 'loaded',
      formData: {},
      invalidClass: 'invalidated',
      jsonOptions: {
        "modes": ["tree", "code", "text"]
      },
      jsonData: {},
      notifyData: '',
    }
  },
  methods: {
    returnBackToList() {
      this.$router.go(-1);
    },
    LoadConfig() {
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
            this.loading = 'loaded';
            console.log('loaded config');
            this.apiUrl = res.data.apiUrl;
            this.redirectToUrl = res.data.redirectToUrl;
            this.formData = res.data.formData;
            console.log('title', this.title);
            console.log('title', res.data.title);
            console.log('filterValues', res.data.filterValues);
            this.configureFields();
          });
      console.log('after load config func');
    },
    validate(field, val, rule) {
      if (rule === undefined ||
          rule === null ||
          rule === '' ||
          (this.formData[field].required !== true && val == '')
      ) {
        this.formData[field].class = '';
        return
      }
      if (!rule.test(val)) {
        console.log('INVALID');
        this.formData[field].class = this.invalidClass;
      } else {
        this.formData[field].class = '';
      }
      console.log(field, val, rule, this.formData);
    },
    moveUp() {
      window.scrollTo(0, 0);
    },
    deleteItem() {
      let entityId = parseInt(this.$route.params.entity);
      this.actionText = 'Update';
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        return;
      }
      axios.delete(this.apiUrl + entityId, {headers: {'Authorization': 'Bearer ' + token.accessToken}})
          .then(() => {
            this.$router.push({path: this.redirectToUrl});
          })
          .catch((err) => {
            this.errorHandle(err);
          });
    },
    errorHandle(err) {
      if (err.response) {
        if (err.response.status === 401) {
          this.$router.push({path: this.login});
        }
        this.error = '[' + err.response.status + ' ' + err.response.statusText + '] ';
        if (err.response.data != '') {
          if (err.response.data.error !== undefined) {
            this.error += err.response.data.error;
          } else {
            this.error += err.response.data;
          }
        }
      } else {
        this.error = 'Server error:  no connection'
      }
      this.moveUp();
    },
    submitForm() {
      this.error = '';
      let data = {};
      for (let realFieldName in this.formData) {
        if (this.formData[realFieldName].class === this.invalidClass) {
          this.error = 'invalid format for field: ' + this.formData[realFieldName].name;
          this.moveUp();
          return;
        }
        data[realFieldName] = this.formData[realFieldName].value;
        if (this.formData[realFieldName].type === 'number') {
          data[realFieldName] = parseInt(this.formData[realFieldName].value)
        }
        if (this.formData[realFieldName].type === 'json') {
          data[realFieldName] = JSON.stringify(this.jsonData[realFieldName]);
        }
      }
      let itemId = parseInt(this.$route.params.entity);
      if (itemId > 0) {
        data['id'] = itemId;
      }
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        this.moveUp();
        return;
      }
      if (this.action == 'add') {
        axios.post(this.apiUrl, data, {headers: {'Authorization': 'Bearer ' + token.accessToken}})
                .then(res => {
                  if (res !== null && res !== undefined) {
                    if (res.data.message !== null && res.data.message !== undefined) {
                      this.notifyData = res.data.message;
                      this.$modal.show('notifyData');
                    } else {
                      this.$router.push({path: this.redirectToUrl});
                    }
                  } else {
                    this.$router.push({path: this.redirectToUrl});
                  }
            })
            .catch((err) => {
              this.errorHandle(err);
            });
      } else {
        axios.patch(this.apiUrl, data, {headers: {'Authorization': 'Bearer ' + token.accessToken}})
            .then(() => {
              this.$router.push({path: this.redirectToUrl})
            })
            .catch((err) => {
              this.errorHandle(err);
            });
      }
    },
    configureFields() {
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        this.moveUp();
        return;
      }
      this.loading = 'loading';
      for (let field in this.formData) {
        if (this.formData[field].type === 'select' && this.formData[field].url != '') {
          //console.log('load for select', this.formData[field].url);
          this.formData[field].options = [];
          axios.post(this.formData[field].url, [],
              {headers: {'Authorization': 'Bearer ' + token.accessToken}})
              .catch((err) => {
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
                this.loading = 'loaded';
                if (res.data.data !== undefined) {
                  for (let dataItem in res.data.data) {
                    console.log(res.data.data[dataItem]);
                    this.formData[field].options.push(
                        {
                          value: res.data.data[dataItem]["Id"],
                          text: res.data.data[dataItem][this.formData[field].field]
                        }
                    )
                  }
                }
              })
        }
      }
      if (this.error !== '') {
        return;
      }
      let entity = this.$route.params.entity;
      if (entity !== 'new') {
        this.actionText = 'Update';
        let entityId = parseInt(entity);
        axios.get(this.apiUrl + entityId, {headers: {'Authorization': 'Bearer ' + token.accessToken}})
            .then(res => {
              if (res.data != '') {
                this.action = 'update';
                for (let field in res.data.data) {
                  if (this.formData[field] != undefined) {
                    if (field == 'Id') {
                      this.formData[field].type = 'number';
                      this.formData[field].show = true;
                    }
                    if (this.formData[field].type === 'json') {
                      this.jsonData[field] = JSON.parse(res.data.data[field]);
                    }
                    this.formData[field].value = res.data.data[field];
                  }
                }
              }
            })
            .catch((err) => {
              this.errorHandle(err);
            });
      } else {
        this.actionText = 'Create';
        this.action = 'add';
      }
    },
    onChangeJson() {
      console.log('onChange json');
    },
  },
  watch: {
    '$route': 'LoadConfig'
  },
  computed: {
    textareaData() {
      for (let field in this.formData) {
        if (this.formData[field].type === 'textarea') {
          return this.formData[field].value;
        }
      }

      return '';
    },
  },
  created: function () {
    this.LoadConfig();
  }
}
</script>
