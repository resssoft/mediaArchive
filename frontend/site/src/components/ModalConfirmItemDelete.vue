<template>
    <modal
            name="confirmItemDelete"
            height="auto"
            width="400px"
    >
        <div class="window-header">
            Confirm item {{ itemId }} delete
            <a href="javascript:void(0)" class="close-modal pull-right" @click="$modal.hide('confirmItemDelete')"><i
                    class="fa fa-times"></i></a>
        </div>
        <div class="modal-body">
            Are you sure you want to delete item?
        </div>
        <div class="modal-footer">
            <button
                    class="btn btn-primary"
                    href="javascript:void(0);"
                    v-on:click="deleteCurrentItem()"
                    type="button"
            >
                <span>Yes, delete</span>
            </button>
            <button
                    class="btn btn-default"
                    href="javascript:void(0);"
                    v-on:click="$modal.hide('confirmItemDelete')"
                    type="button"
            >
                <span>Cancel</span>
            </button>
        </div>
    </modal>
</template>

<script>
    import axios from "axios";
    import {TokenMixin} from '../assets/js/mixins/token-mixin';

    export default {
        name: 'ModalConfirmItemDelete',
        mixins: [
            TokenMixin,
        ],
        props: {
            url: {
                type: String,
                required: true
            },
            itemId: {
                type: Number,
                required: true
            },
        },
        methods: {
            deleteCurrentItem() {
                let url = this.url + this.itemId;
                if (!url) {
                    this.$emit('onError', 'Invalid params');
                    return;
                }
                axios.delete(url, {headers: {'Authorization': 'Bearer ' + this.actualToken()}})
                    .then(() => {
                        this.$modal.hide('confirmItemDelete');
                        this.$emit('onComplete');
                    })
                    .catch((err) => {
                        this.$modal.hide('confirmItemDelete');
                        this.$emit('onError', err);
                    });
            }
        }
    }
</script>
