<template>
  <modal
      :name=modalName
      :scrollable="true"
      height="80%"
      width="60%"
      class="logs-modal"
  >
    <div class="window-header">
      Logs data
      <a href="javascript:void(0)" class="close-modal pull-right" @click="$modal.hide(modalName)"><i
          class="fa fa-times"></i></a>
    </div>
    <div class="modal-body" style="overflow: scroll; height: calc(100% - 46px)">
      <pre>{{ data }}</pre>
    </div>
  </modal>
</template>

<script>

export default {
  name: 'ModalLog',
  components: {

  },
  props: [
    'data',
    'modalName',
  ],
  data() {
    return {
      success: '',
    }
  },
  methods: {

  },
  mounted() {

  },
  watch: {

  },
  computed: {

  },
  created: function () {

  }
}
</script>
