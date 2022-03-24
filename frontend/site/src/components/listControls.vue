<template>
    <div class="table-controls-list">
        <div class="table-controls">
            <div class="select-total float-right">
                Total: {{ total }}
            </div>
            <div class="select-limit">
                <b-form-select @change="limitChanged(limit)" v-model="limit" :options="options"></b-form-select>
            </div>
        </div>
        <div class="pagination-container" v-show="total > limit">
            <button @click="prevPage" :disabled="pageNumber <= 1" class="float-left btn btn-default">
                Previous
            </button>
            <button @click="nextPage" :disabled="pageNumber >= pageCount" class="btn btn-default">
                Next
            </button>
            <ul class="float-left">
                <li @click="pageChanged(p)" v-for="p in paginatedData" v-bind:key="p.first" class="pagination-item"
                    v-bind:class="{ 'pagination-current' : (p == pageNumber)}">
                    {{ p }}
                </li>
            </ul>
        </div>
    </div>
</template>
<script>
    import {UserStoredDataMixin} from '../assets/js/mixins/user-stored-data-mixin';
    export default {
        name: 'listControls',
        mixins: [
            UserStoredDataMixin,
        ],
        props: {
            total: Number,
            options: {
                type: Array,
                default: function () {
                    return [
                        {
                            value: '10',
                            text: '10',
                        },
                        {
                            value: '100',
                            text: '100',
                        },
                        {
                            value: '0',
                            text: 'All',
                        },
                    ]
                }
            }
        },
        data() {
            return {
                pageNumber: 1,
                limit: 20,
                dataKeysForStore:['limit'],
            };
        },
        methods: {
            onRouteChange() {
                if (this.$route.query.page != undefined) {
                    this.pageNumber = this.$route.query.page;
                } else {
                    this.pageNumber = 1;
                }
                this.limit = 20;
                this.ApplyUserConfig();
                this.UpdateUserConfig();
            },
            pageChanged(page) {
                if (page === '...') {
                    return
                }
                this.pageNumber = page;
                this.$router.push({query: {page: this.pageNumber}});
                this.$emit('onPageChanged', page);
            },
            limitChanged(limit) {
                this.$emit('onLimitChanged', limit);
            },
            prevPage() {
                this.pageNumber--;
                this.pageChanged(this.pageNumber);
            },
            nextPage() {
                this.pageNumber++;
                this.pageChanged(this.pageNumber);
            },
        },
        computed: {
            pageCount() {
                let l = this.total,
                    s = parseInt(this.limit) || 1;
                if (s === 0) {
                    return 1
                }
                return Math.ceil(l / s) || 1;
            },
            pagesViewed() {
                return this.pageCount > this.limit;
            },
            paginatedData() {
                let pageNumbersFull = Array.from(Array(parseInt(this.pageCount) + 1).keys());
                pageNumbersFull = pageNumbersFull.filter((val) => val !== 0);
                let pageItems = pageNumbersFull.slice(0, pageNumbersFull.length);
                if (this.pageCount > 10) {
                    if (this.pageNumber > 6) {
                        pageItems = [];
                        pageItems.push(1);
                        pageItems.push("...");
                        if ((pageNumbersFull.length - this.pageNumber) > 10) {
                            pageItems.push(...pageNumbersFull.slice(this.pageNumber, this.pageNumber + 10));
                            pageItems.push("...");
                            pageItems.push(this.pageCount);
                        } else {
                            pageItems.push(...pageNumbersFull.slice(this.pageNumber, this.pageCount));
                        }
                    } else {
                        pageItems = pageNumbersFull.slice(0, 10);
                        pageItems.push("...");
                        pageItems.push(this.pageCount);
                    }
                }
                return pageItems;
            },
        },
        watch: {
            '$route': 'onRouteChange'
        },
        mounted: function () {
            this.pageNumber = this.$route.query.page || 1;
        },
        created: function () {
          this.pageNumber = this.$route.query.page || 1;
          this.ApplyUserConfig();
        },
    }
</script>
