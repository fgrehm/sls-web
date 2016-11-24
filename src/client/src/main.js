import Vue from 'vue'

import App from './App'
import router from './router'
import store from './store'

/* eslint-disable no-new */
new Vue({
  el: '#app',
  template: '<App/>',
  components: { App },
  router,
  store
})
