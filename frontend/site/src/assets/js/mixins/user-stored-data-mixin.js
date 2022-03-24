
let UserStoredDataMixin = {
    data: function () {
        return {
            defaultSavedFields: [],
        }
    },
    created: function () {
    },
    methods: {
        ApplyUserConfig() {
            let userPageConfig = null;
            let pageName = this.$route.name;
            if (pageName !== '') {
                let upcJson = localStorage.getItem(pageName + '-page');
                if (upcJson !== null) {
                    userPageConfig = JSON.parse(upcJson);
                    for (let configItem in userPageConfig) {
                        if (Object.keys(this).includes(configItem)) {
                            this[configItem] = userPageConfig[configItem];
                        }
                    }
                }
            }
        },
        UpdateUserConfig() {
            let userPageConfig = {};
            let pageName = this.$route.name;
            if (pageName !== '') {
                let savedJson = localStorage.getItem(pageName + '-page');
                if (savedJson !== null) {
                    userPageConfig = JSON.parse(savedJson);
                    for (let configItem in userPageConfig) {
                        if (Object.keys(this).includes(configItem)) {
                            userPageConfig[configItem] = this[configItem];
                        }
                    }
                }
                let keysForSave = this.defaultSavedFields;
                if (Object.keys(this).includes('dataKeysForStore')) {
                    keysForSave = this['dataKeysForStore'];
                }
                if (keysForSave.length === 0 && Object.keys(userPageConfig).length !== 0) {
                    keysForSave = Object.keys(userPageConfig);
                }
                for (let key in keysForSave) {
                    if (Object.keys(this).includes(keysForSave[key])) {
                        userPageConfig[keysForSave[key]] = this[keysForSave[key]];
                    }
                }
                localStorage.setItem(pageName + '-page', JSON.stringify(userPageConfig));
            }
        },
    }
};

export {UserStoredDataMixin};