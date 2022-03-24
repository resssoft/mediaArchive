
let TokenMixin = {
    data: function () {
        return {}
    },
    created: function () {
    },
    methods: {
        /**
         * @return {string}
         */
        actualToken() {
            const token = JSON.parse(localStorage.getItem('token'));
            if (!token) {
                return '';
            }
            if (!token.accessToken) {
                return '';
            }
            return token.accessToken;
        }
    }
};

export {TokenMixin};