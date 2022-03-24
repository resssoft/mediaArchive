<template>
  <modal
      name="confirmVmDeploy"
      height="auto"
      width="400px"
  >
    <div class="window-header">
      Confirm VM deploy
      <a href="javascript:void(0)" class="close-modal pull-right" @click="$modal.hide('confirmVmDeploy')"><i
          class="fa fa-times"></i></a>
    </div>
    <div class="modal-body">
      Are you sure you want to deploy VM?
      <br>
      <br>
      <b-form-checkbox
          id="checkbox-unavailable"
          v-model="unavailable"
          name="checkbox-unavailable"
          value="y"
          unchecked-value=""
      >
        Ignore unavailable servers
      </b-form-checkbox>
      <b-form-checkbox
          id="checkbox-recreate"
          v-model="recreate"
          name="checkbox-recreate"
          value="y"
          unchecked-value=""
      >
        Recreate existing VMs
      </b-form-checkbox>
      <b-form-checkbox
          id="checkbox-changed"
          v-model="changed"
          name="checkbox-changed"
          value="y"
          unchecked-value=""
      >
        VM was changed
      </b-form-checkbox>
      <b-form-checkbox
          id="checkbox-moved"
          v-model="moved"
          name="checkbox-moved"
          value="y"
          unchecked-value=""
      >
        VM was moved
      </b-form-checkbox>
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
          v-on:click="$modal.hide('confirmVmDeploy')"
          type="button"
      >
        <span>Cancel</span>
      </button>
    </div>
  </modal>
</template>

<script>
export default {
  name: 'ModalConfirmVmDeploy',
  props: {},
  data() {
    return {
      unavailable: false,
      recreate: false,
      changed: false,
      moved: false,
    }
  },
  methods: {
    confirm() {
      this.$emit(
          'onConfirm',
          this.unavailable === "y",
          this.recreate === "y",
          this.changed === "y",
          this.moved === "y"
      );
    }
  }
}
</script>
