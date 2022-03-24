<template>
  <modal
          :name=modalName
          height="auto"
          width="600px"
          :scrollable="true"
  >
  <div class="window-header">
    {{ itemTitle }}
    <a href="javascript:void(0)" class="close-modal pull-right" @click="$modal.hide(modalName)">
      <i class="fa fa-times"></i>
    </a>
    </div>
    <div class="modal-body">
    <div class="page" style="background:white">
      <div class="alert alert-danger alert-dismissable" v-show="error !== ''">
        <strong>Error </strong>
        <div class="errorMsg" style="display: inline-block;">{{ error }}</div>
      </div>
      <form ref="form" class="form-horizontal" @submit.prevent="submitForm">
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
                    :class="'col-sm-11 ' + formData[index].class"
                    :rows="formData[index].rows || 5"
                    v-if="formData[index].type === 'textarea'"
                    v-model="formData[index].value"></b-form-textarea>
            <v-jsoneditor
                    v-if="formData[index].type === 'json'"
                    v-model="jsonData[index]"
                    :plus="false"
                    :options="jsonOptions"
                    :height="'400px'"></v-jsoneditor>
            <button
                    type="button"
                    v-if="formData[index].type === 'textarea' && formData[index].name !== 'Import'"
                    class="btn btn-default col-sm-1"
                    @click="showTextareaData()"
            >Open
            </button>
          </div>
        </div>
        </span>
        <div class="form-group">
          <div class="col-sm-offset-2 col-sm-10" style="text-align:right">
            <button type="submit" class="btn btn-primary">{{ actionText }}</button>
          </div>
        </div>
      </form>
    </div>
  </div>
  <ModalLog :data="textareaData" :modalName="'itemTextarea'" />
  </modal>
</template>

<script>
import axios from 'axios';
import VJsoneditor from 'v-jsoneditor/src/index';
import ModalLog from "./ModalLog";

export default {
  name: 'newItem',
  components: {
    VJsoneditor,
    ModalLog,
  },
  props: {
    type: String,
    itemId: Number,
    modalName: String,
  },
  data() {
    return {
      error: '',
      actionText: '...',
      action: 'add',
      apiUrl: '',
      redirectToUrl: '',
      helperApiUrl: '',
      apiConfig: '/api/config/item-form/',
      login: '/login',
      loading: 'loaded',
      formData: {},
      invalidClass: 'invalidated',
      jsonOptions: {
        "modes": ["tree", "code", "text"]
      },
      jsonData: {},
      itemTitle: '',
    }
  },
  methods: {
    returnBackToList() {
      this.$router.go(-1);
    },
    LoadConfig() {
      let loader = this.$loading.show();

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
            loader.hide();
            this.loading = 'loaded';
            this.errorHandle(err);
          })
          .then(res => {
            loader.hide();
            this.loading = 'loaded';
            console.log('loaded config');
            this.apiUrl = res.data.apiUrl;
            this.redirectToUrl = res.data.redirectToUrl;
            if (res.data.helperApiUrl !== undefined) {
              this.helperApiUrl = res.data.helperApiUrl;
            }
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
      let re = new RegExp(rule, "i");
      if (!re.test(val)) {
        console.log('INVALID');
        this.formData[field].class = this.invalidClass;
      } else {
        this.formData[field].class = '';
      }
      console.log(field, val, rule, this.formData);
      if (this.helperApiUrl !== '') {
        console.log(this.helperApiUrl);
        this.getValidData();
      }
    },
    moveUp() {
      window.scrollTo(0, 0);
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
      let loader = this.$loading.show();
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
      let itemId = this.itemId;
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
              loader.hide();
              if (res !== null && res !== undefined) {
                if (res.data.message !== null && res.data.message !== undefined) {
                  this.notifyData = res.data.message;
                  this.showNotifyData();
                } else {
                  this.$emit('onSuccess');
                }
              } else {
                this.$emit('onSuccess');
              }
            })
            .catch((err) => {
              loader.hide();
              this.errorHandle(err);
            });
      } else {
        axios.patch(this.apiUrl, data, {headers: {'Authorization': 'Bearer ' + token.accessToken}})
            .then(() => {
              loader.hide();
              this.$emit('onSuccess');
            })
            .catch((err) => {
              loader.hide();
              this.errorHandle(err);
            });
      }
    },
    configureFields() {
      let loader = this.$loading.show();
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        this.moveUp();
        return;
      }
      this.loading = 'loading';
      for (let field in this.formData) {
        if (this.formData[field].type === 'select' && this.formData[field].url != '') {
          this.formData[field].options = [];
          axios.post(this.formData[field].url, [],
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
                  for (let dataItem in res.data.data) {
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
      if (this.type !== 'create') {
        if (this.itemId == 0) {
          loader.hide();
          return;
        }
        this.actionText = 'Update';
        let entityId = this.itemId;
        axios.get(this.apiUrl + entityId, {headers: {'Authorization': 'Bearer ' + token.accessToken}})
            .then(res => {
              loader.hide();
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
              loader.hide();
              this.errorHandle(err);
            });
      } else {
        loader.hide();
        this.actionText = 'Create';
        this.action = 'add';
      }
    },
    getValidData() {
      let loader = this.$loading.show();
      this.error = '';
      let data = {};
      for (let realFieldName in this.formData) {
        if (this.formData[realFieldName].class === this.invalidClass) {
          this.error = 'invalid format for field: ' + this.formData[realFieldName].name;
          this.moveUp();
          return;
        }
        data[realFieldName] = this.formData[realFieldName].value;
        if (this.formData[realFieldName].type === 'number' || this.formData[realFieldName].toType === 'number') {
          data[realFieldName] = parseInt(this.formData[realFieldName].value)
        }
        if (this.formData[realFieldName].type === 'json') {
          data[realFieldName] = JSON.stringify(this.jsonData[realFieldName]);
        }
      }
      let itemId = this.itemId;
      if (itemId > 0) {
        data['id'] = itemId;
      }
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        this.moveUp();
        return;
      }
      if (this.helperApiUrl !== '') {
        axios.post(this.helperApiUrl, data, {headers: {'Authorization': 'Bearer ' + token.accessToken}})
                .then(res => {
                  loader.hide();
                  if (res !== null && res !== undefined) {
                    if (res.data !== null && res.data !== undefined) {
                      for (let fieldName in res.data) {
                        if (this.formData[fieldName] !== undefined) {
                          if (res.data[fieldName] != '') {
                            this.formData[fieldName].value = res.data[fieldName];
                          }
                        }
                      }
                    }
                  }
                })
                .catch((err) => {
                  loader.hide();
                  this.errorHandle(err);
                });
      }
    },
    onChangeJson() {
      console.log('onChange json');
    },
    showTextareaData() {
      this.$modal.show('itemTextarea');
    },
    showNotifyData() {
      this.$modal.show('itemTextarea');
    },
  },
  mounted() {

  },
  watch: {
    '$route': 'LoadConfig',
    'itemId': 'LoadConfig'
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
    if (this.type === 'create') {
      this.itemTitle = 'New item';
    } else {
      this.itemTitle = 'Edit item ' + this.itemId;
    }
    this.LoadConfig();
  }
}
</script>
