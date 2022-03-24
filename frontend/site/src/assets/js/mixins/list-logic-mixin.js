
import axios from "axios";
import Header from "../../../components/Header";
import ModalItem from "../../../components/ModalItem";
import Filters from "../../../components/listFilters";
import ModalError from "../../../components/ModalError";
import ModalSuccess from "../../../components/ModalSuccess";
import ModalLog from "../../../components/ModalLog";
import ModalConfirmItemDelete from "../../../components/ModalConfirmItemDelete";
import ListControls from "../../../components/listControls";

// mixins
import {HwProgressBarMixin} from './hw-progress-bar-mixin';
import {TokenMixin} from './token-mixin';
import {UserStoredDataMixin} from './user-stored-data-mixin';

let ListLogicMixin = {
    mixins: [
        HwProgressBarMixin,
        TokenMixin,
        UserStoredDataMixin,
    ],
    components: {
        Header,
        ModalItem,
        Filters,
        ModalError,
        ModalSuccess,
        ModalLog,
        ModalConfirmItemDelete,
        ListControls,
    },
    data() {
        return {
            error: '',
            success: '',
            info: '',
            filter: [],
            filterValues: {},
            filters: [],
            total: 0,
            login: '/login',
            show: false,
            limit: 20,
            pageNumber: 1,
            sortField: 'id',
            loading: 'loaded',
            moving: false,
            currentColumn: null,
            nextColumn: null,
            buttons: null,
            dataKeysForStore:['limit', 'sortField'],
        };
    },
    methods: {
        onRouteChange() {
            this.sortField = 'id';
            this.limit = 20;
            this.ApplyUserConfig();
            this.UpdateUserConfig();
            this.LoadConfig();
        },
        button(buttonName) {
            console.log(buttonName);
            if (this.buttons === null) {
                return;
            }
            if (
                this.buttons[buttonName] === undefined ||
                this.buttons[buttonName]["api"] === undefined ||
                this.buttons[buttonName]["class"] === "running"
            ) {
                return;
            }

            const token = JSON.parse(localStorage.getItem('token'));
            if (token == null) {
                this.error = 'token not found';
                return;
            }
            this.error = '';
            axios.post(this.buttons[buttonName]["api"],
                {
                    name: buttonName
                },
                {headers: {'Authorization': 'Bearer ' + token.accessToken}})
                .catch((err) => {
                    this.errorHandle(err);
                })
                .then(res => {
                    this.buttons[buttonName]["class"] = "running";
                    if (res !== null && res !== undefined) {
                        this.info = 'Running';
                        setTimeout(this.clearInfo, 5000);
                        setTimeout(this.checkButtonsStatus, 2000);
                    }
                })
        },
        clearInfo() {
            this.info = '';
        },
        fillTemplate(tmp, itemFields) {
            for (let field in itemFields) {
                tmp = tmp.replaceAll("##" + field + "##", itemFields[field])
            }
            let matches = [...tmp.matchAll(/##to_low_deep_trim@([^#]*?)##/g)];
            let processedString = '';
            // replace something like this "##to_low_deep_trim@MY VALUE##" to "my-value"
            for (const match of matches) {
                if (match.length === 2) {
                    processedString = match[1].replace(' ', '-').toLocaleLowerCase();
                    tmp = tmp.replaceAll(match[0], processedString)
                }
            }
            return tmp
        },
        checkButtonsStatus() {
            for (let buttonName in this.buttons) {
                if (this.buttons[buttonName]["class"] === "running") {
                    const token = JSON.parse(localStorage.getItem('token'));
                    if (token == null) {
                        this.error = 'token not found';
                        return;
                    }
                    axios.post(this.buttons[buttonName]["apiStatus"],
                        {
                            name: buttonName
                        },
                        {headers: {'Authorization': 'Bearer ' + token.accessToken}})
                        .catch((err) => {
                            this.errorHandle(err);
                        })
                        .then(res => {
                            if (res !== null && res !== undefined) {
                                if (res.data != "") {
                                    if (res.data.data !== "processed") {
                                        this.buttons[buttonName]["class"] = "";

                                        this.getAndShowDiscoverLogs();
                                        this.LoadData();
                                        this.success = 'Script successfully finished!';
                                        this.$modal.show('listSuccess');
                                    } else {
                                        setTimeout(this.checkButtonsStatus, 2000);
                                    }
                                }
                            }
                        });
                }
            }
        },

        startLoading() {
            this.loading = 'loading';
        },
        endLoading() {
            this.loading = 'loaded';
        },
        openDeleteItemConfirmation(item) {
            this.currentItemId = parseInt(item.Id);
            this.$modal.show('confirmItemDelete');
        },
        deletedItem() {
            this.endLoading();
            this.success = 'Item successfully deleted!';
            this.$modal.show('listSuccess');
            this.endLoading();
            this.LoadData();
        },
        deletedItemErr(err) {
            this.endLoading();
            this.errorHandle(err);
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
            }
            else if (err != "") {
                this.error = err;
            }
            else {
                this.error = 'Server error:  no connection'
            }
            this.$modal.show('listError');
        },
        moveUp() {
            window.scrollTo(0, 0);
        },
        updateResizable() {
            let el = document.getElementById('resizable-table');
            if (el == null) {
                return;
            }
            const nodeName = el.nodeName;
            if (['TABLE', 'THEAD'].indexOf(nodeName) < 0) return;
            const table = nodeName === 'TABLE' ? el : el.parentElement;
            const thead = table.querySelector('thead');
            const ths = thead.querySelectorAll('th');
            ths.forEach((th) => {
                th.style.width = th.offsetWidth + 'px';
            });
        },
        columnsReset() {
            let el = document.getElementById('resizable-table');
            if (el == null) {
                return;
            }
            const nodeName = el.nodeName;
            if (['TABLE', 'THEAD'].indexOf(nodeName) < 0) return;
            const table = nodeName === 'TABLE' ? el : el.parentElement;
            const thead = table.querySelector('thead');
            const ths = thead.querySelectorAll('th');
            ths.forEach((th) => {
                th.style.width = null;
            });
        },
        resizerMousedown(event) {
            this.moving = true;
            this.currentColumn = event.target.parentElement;
            this.nextColumn = event.target.parentElement.nextSibling;
        },
        initResizable() {
            let el = document.getElementById('resizable-table');
            const nodeName = el.nodeName;
            if (['TABLE', 'THEAD'].indexOf(nodeName) < 0) return;
            const table = nodeName === 'TABLE' ? el : el.parentElement;
            table.style.position = 'relative';

            document.addEventListener('mouseup', () => {
                this.moving = false;
                this.currentColumn = null;
                this.nextColumn = null;
            });

            const cutPx = str => +str.replace('px', '');
            const handleResize = e => {
                if (this.moving) {
                    if (this.currentColumn !== null) {
                        this.currentColumn.style.width = cutPx(this.currentColumn.style.width) + e.movementX + 'px';
                    }
                    if (this.currentColumn !== null) {
                        this.nextColumn.style.width = cutPx(this.nextColumn.style.width) - e.movementX + 'px';
                    }
                }
            };
            table.addEventListener('mousemove', handleResize);
        },
        filterChanged(resultFilter) {
            this.filters = resultFilter;
            this.pageNumber = 1;
            if (this.$route.query.page != undefined) {
                history.pushState({}, null, `#${this.$route.path}`);
            }
            this.LoadData();
        },
        SortInvert(field) {
            if (field == this.sortField) {
                this.sortField = '-' + field
            } else {
                this.sortField = field
            }
            this.show = !this.show;
            this.UpdateUserConfig();
            this.LoadData();
        },
        filterPreparing() {
            const token = JSON.parse(localStorage.getItem('token'));
            if (token == null) {
                this.error = 'token not found';
                return;
            }
            for (let filter in this.filter) {
                if (this.filter[filter]["external"] !== undefined) {
                    if (this.filter[filter]["external"]["url"] === undefined || this.filter[filter]["external"]["url"] === "") {
                        continue;
                    }
                    axios.post(this.filter[filter]["external"]["url"],
                        [],
                        {headers: {'Authorization': 'Bearer ' + token.accessToken}})
                        .then(res => {
                            if (res === undefined || res === null) {
                                return;
                            }
                            if (res.data.data !== undefined) {
                                for (let dataIndex in res.data.data) {
                                    this.filter[filter]["options"].push({
                                        value: res.data.data[dataIndex][this.filter[filter]["external"]["fieldValue"]],
                                        text: res.data.data[dataIndex][this.filter[filter]["external"]["fieldText"]],
                                    });
                                }
                            }
                        });
                }
            }
        },
        ChangeLimit(newLimit) {
            this.limit = newLimit;
            this.UpdateUserConfig();
            this.LoadData();
        },

        itemSuccess() {
            this.$modal.hide('editItem');
            this.$modal.hide('createItem');
            this.LoadData();
        },

    },
    computed: {
        pageCount() {
            let l = this.total,
                s = parseInt(this.limit);
            if (s === 0) {
                return 1
            }
            return Math.ceil(l / s);
        },
        pagesViewed() {
            return this.pageCount > this.limit;
        },
        paginatedData() {
            let pageNumbersFull = Array.from(Array(parseInt(this.pageCount) + 1).keys())
            let pageItems = pageNumbersFull.slice(1, pageNumbersFull.length)
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
                    pageItems = pageNumbersFull.slice(0, 10)
                    pageItems.push("...");
                    pageItems.push(this.pageCount);
                }
            }
            return pageItems;
        },
        isDataVisible() {
            if (this.loading == 'loading') {
                return false;
            } else {
                return true;
            }
        },
        isLoading() {
            if (this.loading == 'loading') {
                return true;
            } else {
                return false;
            }
        },
    },
    watch: {

    },
    created: function () {

    },
    mounted() {

    },
    updated: function () {

    }
};

export {ListLogicMixin};