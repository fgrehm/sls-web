import api from '../../api'

const LOADING = 'models/LOADING'
const LOADED = 'models/LOADED'
const ERROR = 'models/ERROR'

const initialState = {
  loading: false,
  loaded: false,
  data: [],
  error: null
}

const mutations = {
  [LOADING] (state) {
    state.loading = true
  },

  [LOADED] (state, models) {
    state.loading = false
    state.loaded = true
    state.data = models
  },

  [ERROR] (state, error) {
    state.loading = false
    state.error = error
  }
}

const actions = {
  loadModels ({ commit }) {
    commit(LOADING)
    api.allModels(
      (data) => commit(LOADED, data),
      (error) => commit(ERROR, error)
    )
  }
}

export default {
  state: initialState,
  mutations,
  actions
}
