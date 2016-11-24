import Vue from 'vue'
import VueRouter from 'vue-router'
Vue.use(VueRouter)

import Home from './pages/Home'
import Models from './pages/Models'
import Solutions from './pages/Solutions'
import Experiments from './pages/Experiments'

// REFACTOR: Might want to use some nested routes here
const routes = [
  { path: '/', component: Home },
  { path: '/models/new',
    name: 'newModel',
    component: Models.New },
  { path: '/models/:modelID',
    name: 'editModel',
    component: Models.Edit },
  { path: '/models/:modelID/solutions/new',
    name: 'newSolution',
    component: Solutions.Form },
  { path: '/models/:modelID/experiments/new',
    name: 'newExperiment',
    component: Experiments.New },
  { path: '/models/:modelID/experiments/:experimentID',
    name: 'showExperiment',
    component: Experiments.Show }
]

export default new VueRouter({
  routes
})
