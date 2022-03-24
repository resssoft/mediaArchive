<template>
  <modal
      :name="name"
      height="auto"
      width="400px"
      class="error-modal"
  >
    <div class="window-header">
      Error!
      <a href="javascript:void(0)" class="close-modal pull-right" @click="$modal.hide(name)"><i
          class="fa fa-times"></i></a>
    </div>
    <div class="modal-body">
      {{ message }}
    </div>
    <div class="modal-footer">
      <button
          class="btn btn-default"
          href="javascript:void(0);"
          v-on:click="$modal.hide(name)"
          type="button"
      >
        <span>OK</span>
      </button>
    </div>
  </modal>
</template>

<script>

export default {
  name: 'ModalError',
  components: {

  },
  props: [
    'name',
    'message',
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
