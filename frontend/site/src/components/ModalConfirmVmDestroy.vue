<template>
  <modal
      name="confirmVmDestroy"
      height="auto"
      width="400px"
  >
    <div class="window-header">
      Confirm VM destroy
      <a href="javascript:void(0)" class="close-modal pull-right" @click="$modal.hide('confirmVmDestroy')"><i
          class="fa fa-times"></i></a>
    </div>
    <div class="modal-body">
      Are you sure you want to destroy VM?
      <br>
    </div>
    <div class="modal-footer">
      <button
          class="btn btn-primary"
          href="javascript:void(0);"
          v-on:click="confirm()"
          type="button"
      >
        <span>Yes</span>
      </button>
      <button
          class="btn btn-default"
          href="javascript:void(0);"
          v-on:click="$modal.hide('confirmVmDestroy')"
          type="button"
      >
        <span>Cancel</span>
      </button>
    </div>
  </modal>
</template>

<script>
export default {
  name: 'ModalConfirmVmDestroy',
  props: {},
  data() {
    return {}
  },
  methods: {
    confirm() {
      this.$emit('onConfirm');
    }
  }
}
</script>
