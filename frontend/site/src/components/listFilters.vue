<template>
    <div class="table-filters-container">
        <div class="table-filters">
            <div class="table-filter-wrapper">
                <label class="hidden-xs control-label inline list-filter">Filter:</label>
                <span class="icon-clear" @click="clear()" title="Clear filters">
                  <i class="fa fa-solid fa-trash" style="font-size: 17px"></i>
                </span>
                <span class="input-wrapper" v-for="(filterField, filterIndex) in filter" v-bind:key="filterField.name">
                          <label>{{ filterField.name }}:</label>
                          <b-form-input
                                  v-if="filterField.type === 'input' || filterField.type === '' || filterField.type === undefined"
                                  v-model="filterValues[filter[filterIndex].name]"
                                  :placeholder="filterField.name"
                                  :style="{width: filterField.w + 'px'}"
                                  @update="filterChanged()"></b-form-input>
                          <b-form-select
                                  v-if="filterField.type === 'select'"
                                  v-model="filterValues[filter[filterIndex].name]"
                                  :options="filterField.options"
                                  :style="{width: filterField.w + 'px'}"
                                  @change="filterChanged()"></b-form-select>
                      </span>
            </div>
        </div>
    </div>
</template>
<script>
    function debounce(that, fn, delay) {
        let timeoutID = undefined;
        return () => {
            clearTimeout(timeoutID);
            const args = arguments;
            timeoutID = setTimeout(() => {
                fn.apply(that, args);
            }, delay);
        };
    }

    export default {
        name: 'listFilters',
        props: { 
            result: Array,
            filter: {
                type: Array,
                default: function () {
                    return {
                        name: 'Name',
                        value: '',
                        w: 100
                    }
                }
            },
            filterValues: {
                type: Object,
                default: function () {
                    return {
                        Name: '',
                    }
                }
            }
        },
        data() {
            return {
                filters: []
            };
        },
        methods: {
            filterChanged() {
                this.debounced();
            },
            clear() {
                for (let FieldName in this.filterValues) {
                    this.filterValues[FieldName] = '';
                }
                this.filters = [];
                this.filterChanged();
            },
            filterSend() {
                console.log('component filter changed');
                this.filters = [];
                for (let FieldName in this.filterValues) {
                    if (this.filterValues[FieldName] != '') {
                        let data = {};
                        data[FieldName] = this.filterValues[FieldName];
                        this.filters.push({
                            Condition: "ilike",
                            Data: data,
                        });
                    }
                }
                this.$emit('onChanged', this.filters);
            },
        },
        created: function () {
            this.debounced = debounce(this, this.filterSend, 700);
        }
    }
</script>
