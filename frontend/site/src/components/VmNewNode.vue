<template>
  <div class="wrapper">
    <modal
        name="VmNewNode"
        height="auto"
        width="600px"
        class=""
    >
      <div class="window-header">
        Create new node
        <a href="javascript:void(0)" class="close-modal pull-right" @click="$modal.hide('VmNewNode')"><i
            class="fa fa-times"></i></a>
      </div>
      <div class="modal-body">
        <div class="page" style="background: white;">
          <!--<div class="alert alert-danger alert-dismissable" v-show="error !== ''">
            <strong>Error </strong>
            <div class="errorMsg" style="display: inline-block;">{{ error }}</div>
          </div>-->
          <form class="form-horizontal">
            <div class="form-group">
              <label class="col-sm-2 control-label">Name</label>
              <div class="col-sm-10">
                <input
                    type="text"
                    v-model="currentName"
                    placeholder="Type node name"
                    class="form-control"
                />
              </div>
            </div>
            <div class="form-group">
              <label class="col-sm-2 control-label">DC</label>
              <div class="col-sm-10">
                <select class="custom-select" v-model="currentDCid" @change="onDCselect()">
                  <option value="">Select DC...</option>
                  <option v-for="(item, index) in DCsSelectItems" :value="item.id" :key="index">{{
                      item.value
                    }}
                  </option>
                </select>
              </div>
            </div>
            <div class="form-group">
              <label class="col-sm-2 control-label">Net</label>
              <div class="col-sm-10">
                <select class="custom-select" v-model="currentNetId" @change="loadHyperVisorByNet()">
                  <option value="" selected>Select Net...</option>
                  <option v-for="(item, index) in NetsSelectItems" :value="item.id" :key="index">{{
                      item.value
                    }}
                  </option>
                </select>
              </div>
            </div>
            <div class="form-group">
              <label class="col-sm-2 control-label">HyperVisor</label>
              <div class="col-sm-10">
                <select class="custom-select" v-model="currentHyperVisorId" @change="onHyperVisorSelect()">
                  <option value="" selected>Select HyperVisor...</option>
                  <option v-for="(item, index) in HyperVisorsSelectItems" :value="item.id" :key="index">{{
                      item.value
                    }}
                  </option>
                </select>
              </div>
            </div>
            <div class="form-group">
              <label class="col-sm-2 control-label">Node type</label>
              <div class="col-sm-10">
                <select class="custom-select" v-model="currentNodeType">
                  <option value="" selected>Select node type...</option>
                  <option v-for="(item, index) in nodeTypes" :value="item" :key="index">{{ item }}</option>
                </select>
              </div>
            </div>
            <div class="form-group">
              <!--
              Ram должно быть большим, чем 4 гигабайта и меньше чем фактическая свободная память на гипервизоре минус 17 гиг - резерв для старта
              min="4"
              :max="hv-free-ram - 17"
              !-->
              <label class="col-sm-2 control-label">Ram</label>
              <div class="col-sm-10">
                <input
                    type="number"
                    v-model="currentRam"
                    :placeholder="vm_ram_placeholder"
                    :min="vm_ram_min"
                    :max="vm_ram_max"
                    class="form-control"
                />
              </div>
            </div>
            <div class="form-group">
              <!--
              Cpu оно должно быть равно (желательно) или меньше (неоптимально, но допустимо) количества ядер на сервере. Больше быть не должно
              :max="hv-cpu"
              !-->
              <label class="col-sm-2 control-label">Cpu</label>
              <div class="col-sm-10">
                <input
                    type="number"
                    v-model="currentCpu"
                    :placeholder="vm_cpu_placeholder"
                    :max="vm_cpu_max"
                    class="form-control"
                />
              </div>
            </div>
            <div class="form-group">
              <label class="col-sm-2 control-label">IP</label>
              <div class="col-sm-10">
                <select class="custom-select" v-model="currentIPid">
                  <option value="">Select IP...</option>
                  <option v-for="(item, index) in IPsSelectItems" :value="item.id" :key="index">{{
                      item.value
                    }}
                  </option>
                </select>
              </div>
            </div>
            <div class="form-group">
              <label class="col-sm-2 control-label">ISO-image</label>
              <div class="col-sm-10 editable-select-container">
                <select class="custom-select" v-model="currentIsoImage">
                  <option v-for="(item, index) in ISOsSelectItems" :value="item" :key="index">{{
                    item
                    }}
                  </option>
                </select>
                <input
                    type="text"
                    v-model="currentIsoImage"
                    class="form-control"
                />
              </div>
            </div>
          </form>
        </div>
      </div>
      <div class="modal-footer">
        <button @click="submitNewNode()" class="btn btn-primary">Create node</button>
        <button
            class="btn btn-default"
            href="javascript:void(0);"
            v-on:click="$modal.hide('VmNewNode')"
            type="button"
        >
          <span>Cancel</span>
        </button>
      </div>
    </modal>

    <ModalSuccess name="VmNewNodeSuccess" :message="success"/>

    <ModalError name="VmNewNodeError" :message="error" />

  </div>
</template>

<script>

import axios from 'axios';
import ModalSuccess from './ModalSuccess.vue';
import ModalError from './ModalError.vue';

export default {
  name: 'VmNewNode',
  components: {
    ModalError,
    ModalSuccess,
  },
  props: [
    'type',
    'itemId',
    'successCreateCallback'
  ],
  data() {
    return {
      // urls
      submitNewNodeUrl: '/api/manager/vm/',
      loadDCsUrl: '/api/manager/dc/list?sort=-FreeIp',
      loadNetsUrl: '/api/manager/net/list?sort=-FreeIp',
      loadIPsUrl: '/api/manager/ip/list?sort=-DiscoveredAt',
      loadHyperVisorsUrl: 'api/manager/esxi/list?sort=id',
      loadISOsUrl: '/api/launcher/content?parsed=y',
      apiConfig: '/api/config/item-form/',
      login: '/login',

      // data
      currentName: '',
      DCs: {},
      DCsSelectItems: [],
      currentDCid: '',
      Nets: {},
      NetsSelectItems: [],
      currentNetId: '',
      IPs: {},
      IPsSelectItems: [],
      ISOsSelectItems: [],
      currentIPid: '',
      HyperVisors: {},
      HyperVisorsSelectItems: [],
      currentHyperVisorId: '',
      nodeTypes: [
        'p2p',
        'astro',
        'prx'
      ],
      currentNodeType: '',
      currentCpu: '',
      currentRam: '',
      currentIsoImage: '',

      vm_ram_min: 4,
      vm_ram_max: 0,
      vm_ram_placeholder: '',
      vm_cpu_min: 0,
      vm_cpu_max: 0,
      vm_cpu_placeholder: '',

      // service
      error: '',
      success: '',
      loading: 'loaded',
    }
  },
  methods: {
    init(){
      this.currentName = '';
      this.currentDCid = '';
      this.currentNetId = '';
      this.currentIPid = '';
      this.currentHyperVisorId = '';
      this.currentNodeType = '';
      this.currentCpu = '';
      this.currentRam = '';
      this.currentIsoImage = '';

      this.vm_ram_min = 4;
      this.vm_ram_max = 0;
      this.vm_ram_placeholder = '';
      this.vm_cpu_min = 0;
      this.vm_cpu_max = 0;
      this.vm_cpu_placeholder = '';
    },
    loadDCs() {
      let loader = this.$loading.show();
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        this.moveUp();
        return;
      }
      axios.post(this.loadDCsUrl, [],
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
              this.DCs = res.data.data;
              this.genDCsSelectItems();
            }
          })
    },

    genDCsSelectItems() {
      this.DCsSelectItems = [];
      this.DCs.forEach((dc) => {
        let select_option_text = `${dc.Name}, ${dc.ServerCount} servers, ${dc.VmCount} VMs, ${dc.FreeIp}/${dc.AllIp} free IPs`;
        let select_option = {id: dc.Id, value: select_option_text};
        this.DCsSelectItems.push(select_option);
      });
    },

    onDCselect() {
      this.loadNetsByDC();
      this.loadIPsByDC();
    },

    loadNetsByDC() {
      let loader = this.$loading.show();
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        this.moveUp();
        return;
      }
      axios.post(this.loadNetsUrl,
          [{"Condition": "=", "Data": {"DcId": this.currentDCid}}],
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
              console.log(res.data.data)
              this.Nets = res.data.data;
              this.genNetsSelectItems();
            }
          })
    },

    genNetsSelectItems() {
      this.NetsSelectItems = [];
      this.Nets.forEach((net) => {
        let select_option_text = `${net.Address}, ${net.FreeIp}/${net.AllIp} free IPs`;
        let select_option = {id: net.Id, value: select_option_text};
        this.NetsSelectItems.push(select_option);
      });
    },

    loadIPsByDC() {
      let loader = this.$loading.show();
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        this.moveUp();
        return;
      }
      axios.post(this.loadIPsUrl,
          [
            {
              "Condition": "=",
              "Data": {
                "DcId": this.currentDCid
              }
            },
            {
              "Condition": "ilike",
              "Data": {
                "UsedBy": "none"
              }
            },
            {
              "Condition": "!=",
              "Data": {
                "State": "OK"
              }
            },
          ],
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
              console.log(res.data.data)
              this.IPs = res.data.data;
              this.genIPsSelectItems();
            }
          })
    },

    genIPsSelectItems() {
      this.IPsSelectItems = [];
      this.IPs.forEach((ip) => {
        let select_option_text = `${ip.Address}, network: ${ip.Network}, state: ${ip.State}, DC: ${this.currentDCname}, used by nobody`;
        let select_option = {id: ip.Id, value: select_option_text};
        this.IPsSelectItems.push(select_option);
      });
    },
    genISOsSelectItems() {
      this.ISOsSelectItems = [];
      let loader = this.$loading.show();
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        this.moveUp();
        return;
      }
      axios.post(this.loadISOsUrl,
              {
                "name": "isolist"
              },
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
                  console.log(res.data.data);
                  this.ISOsSelectItems = res.data.data;
                }
              });
    },

    loadHyperVisorByIP() {
      let loader = this.$loading.show();
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        this.moveUp();
        return;
      }
      axios.post(this.loadHyperVisorsUrl,
          [
            {
              "Condition": "=",
              "Data": {
                "IpId": this.currentIPid
              }
            },
          ],
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
              console.log(res.data.data)
              this.HyperVisors = res.data.data;
              this.genHyperVisorsSelectItems();
            }
          })
    },

    loadHyperVisorByNet() {
      let loader = this.$loading.show();
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        this.moveUp();
        return;
      }
      axios.post(this.loadHyperVisorsUrl,
          [
            {
              "Condition": "=",
              "Data": {
                "NetId": this.currentNetId
              }
            },
          ],
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
              console.log(res.data.data)
              this.HyperVisors = res.data.data;
              this.genHyperVisorsSelectItems();
            }
          })
    },

    genHyperVisorsSelectItems() {
      this.HyperVisorsSelectItems = [];
      this.HyperVisors.forEach((hv) => {
        let select_option_text = `${hv.Name} ${hv.Model}, CPU: ${hv.Cpu}/${hv.MaxCpu}, RAM: ${hv.Ram}/${hv.MaxRam}, VmRAM: ${hv.VmSum}/${hv.Ram}, SSD: ${hv.Ssd}, HDD: ${hv.Hdd}`;
        let select_option = {id: hv.Id, value: select_option_text};
        this.HyperVisorsSelectItems.push(select_option);
      });
    },

    onHyperVisorSelect() {
      let current_hv = this.HyperVisors.find(hv => hv.Id == this.currentHyperVisorId);
      this.genVMresourcesLimitsByHyperVisor(current_hv);
    },

    genVMresourcesLimitsByHyperVisor(hypervisor) {
      this.vm_ram_max = hypervisor.Ram - 17;
      this.vm_ram_placeholder = `From ${this.vm_ram_min} to ${this.vm_ram_max}`;
      this.vm_cpu_max = hypervisor.Cpu;
      this.vm_cpu_placeholder = `From ${this.vm_cpu_min} to ${this.vm_cpu_max} (${this.vm_cpu_max} is optimal)`;
    },

    submitNewNode() {
      const token = JSON.parse(localStorage.getItem('token'));
      if (token == null) {
        this.error = 'token not found';
        this.moveUp();
        return;
      }
      let data = {
        Id: null,
        Name: this.currentName,
        State: 'CREATED',
        VmType: this.currentNodeType,
        EsxiId: this.currentHyperVisorId,
        IpId: this.currentIPid,
        Mac: this.currentIPmac,
        Cpu: parseInt(this.currentCpu),
        Ram: parseInt(this.currentRam),
        VmIso: this.currentIsoImage,
        ProjectId: null,
        Comment: '',
        DiscoveredAt: '2021-12-03 06:28:57', //TODO обнулить позже
      };
      axios.post(this.submitNewNodeUrl, data, {headers: {'Authorization': 'Bearer ' + token.accessToken}})
          .then(() => {
            this.$root.$emit('onCreateNodeSuccess');
            this.success = 'Node succesfully created';
            this.$modal.show('VmNewNodeSuccess');
            this.$modal.hide('VmNewNode');
            this.init();
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
      this.$modal.show('VmNewNodeError');
    },

    moveUp() {
      window.scrollTo(0, 0);
    },

  },
  mounted() {
    this.init();
  },
  watch: {},
  computed: {
    currentIPmac: function () {
      let currentIP = this.IPs.find(ip => ip.Id == this.currentIPid);
      return currentIP.Mac;
    },
    currentDCname() {
      let currentDCid = this.DCs.find(dc => dc.Id == this.currentDCid);
      return currentDCid.Name;
    },
  },
  created: function () {
    this.init();
    this.loadDCs();
    this.genISOsSelectItems();
  }
}
</script>
