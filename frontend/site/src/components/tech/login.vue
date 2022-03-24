<template>
  <!--<div class="container">
    <div class="alert alert-danger alert-dismissable" v-show="error != ''">
      <strong>Error!</strong>
      <div class="errorMsg" style="display: inline-block;">{{ error }}</div>
    </div>
    <div class="panel-body">
      <form class="form-horizontal" role="form" @submit="submitForm">
        <div class="form-group">
          <div class="col-md-12">
            <input type="text" class="form-control" placeholder="Username or email" v-model="from.Username">
            <i class="fa fa-user"></i>
          </div>
        </div>
        <div class="form-group">
          <div class="col-md-12">
            <input type="password" class="form-control" placeholder="Password" v-model="from.Password">
            <i class="fa fa-lock"></i>
          </div>
        </div>
        <div class="form-group">
          <div class="col-md-12">
            <button type="submit" class="btn btn-primary btn-block">Sign in</button>
          </div>
        </div>
      </form>
    </div>
  </div>-->

  <div class="container login">
    <div class="row justify-content-center">
      <div class="col-12 col-lg-4" style="height: 200px">
        <div class="top-logo">VM Manager</div>
      </div>
    </div>
    <div class="row justify-content-center">
      <div class="col-12 col-lg-4">
        <div class="card-group">
          <div class="card p-4 login">
            <div class="card-body">
              <div class="alert alert-danger alert-dismissable" v-show="error != ''">
                <strong>Error!</strong>
                <div class="errorMsg" style="display: inline-block;">{{ error }}</div>
              </div>
              <form @submit.prevent="submitForm" role="form">
                <h2 class="text-center">VM Manager login</h2>
                <p class="text-muted text-center">Sign in to Dashboard</p>
                <div class="input-group mb-3">
                  <div class="input-group-prepend">
                    <span class="input-group-text">
                      <i class="fa fa-user"></i>
                    </span>
                  </div>
                  <input
                      class="form-control"
                      type="text"
                      name="email"
                      placeholder='Username or email'
                      value=""
                      v-model="from.Username"
                  />
                </div>
                <div class="input-group mb-3">
                  <div class="input-group-prepend">
                    <span class="input-group-text">
                      <i class="fa fa-lock"></i>
                    </span>
                  </div>
                  <input
                      class="form-control"
                      type="password"
                      name="password"
                      placeholder="Password"
                      v-model="from.Password"
                  />
                </div>
                <div class="input-group mb-3">
                  <div class="form-check">
                    <input class="form-check-input" checked="checked" name="remember" type="checkbox" id="remember">
                    <label class="form-check-label" for="remember">
                      Remember me
                    </label>
                  </div>

                </div>
                <div class="row">
                  <div class="col-6">
                    <button style="width:100%" class="btn btn-primary px-4" type="submit">Sign in</button>
                  </div>
                  <div class="col-6">
                    <a style="white-space:nowrap" class="forgot-link" href="">Forgot password?</a>
                  </div>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

</template>

<script>
import axios from 'axios';

export default {
  name: 'dcnew',
  data() {
    return {
      error: '',
      from: {
        Username: '',
        Password: ''
      },
      apiUrl: '/api/auth/sign-in/',
      redirectToUrl: '/'
    }
  },
  mounted () {
    this.setAdditionalStyles();
  },
  methods: {
    submitForm() {
      let loader = this.$loading.show();

      axios.post(this.apiUrl, this.from)
          .then(res => {
            loader.hide()

            if (res.data.error == undefined) {
              localStorage.setItem('token', JSON.stringify(res.data));
              this.$router.push({path: this.redirectToUrl})
            }
          })
          .catch((err) => {
            loader.hide()

            if (err.response) {
              if (err.response.status === 401) {
                this.$router.push({path: this.login});
              }
              this.error = '[' + err.response.status + ' ' + err.response.statusText + '] ';
              if (err.response.data != '') {
                this.error += err.response.data;
              }
            } else {
              this.error = 'Server error:  no connection'
            }
          });
    },
    setAdditionalStyles() {
      let body = document.querySelector('body');
      body.classList.add('login');
    },
  }
}
</script>