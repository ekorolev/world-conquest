import Vue from '../node_modules/vue/dist/vue.js'
import VueRouter from 'vue-router'
import Vuex from 'vuex'
import App from './App.vue'
import routes from './routes'

import jQuery from 'jquery'
import Popper from 'popper.js'
import moment from 'moment'

window.jQuery = window.$ = jQuery
window.Popper = Popper
window.moment = moment

require('./custom.scss')
require('bootstrap')

Vue.use(VueRouter)
Vue.use(Vuex)

let router = new VueRouter({ routes })

moment.locale('ru')
Vue.filter('formateDate', function (value) {
  if (value) {
    return moment(value).calendar()
  }
})

new Vue({
  router,
  store: require('./store'),
  template: '<App></App>',
  components: {
    App
  }
}).$mount('#app')