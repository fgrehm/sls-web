import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

import models from './modules/models'
import selectedModel from './modules/selectedModel'
import solution from './modules/solution'
import experiment from './modules/experiment'

export default new Vuex.Store({
  modules: {
    models,
    selectedModel,
    solution,
    experiment
  }
})
