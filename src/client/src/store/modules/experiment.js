import api from '../../api'

const CREATING = 'experiment/CREATING'
const LOADING = 'experiment/LOADING'
const LOADED = 'experiment/LOADED'
const ERROR = 'experiment/ERROR'

const emptyExperiment = {
}

const initialState = {
  creating: false,
  loaded: false,
  loading: false,
  data: { ...emptyExperiment },
  error: null
}

const mutations = {
  [CREATING] (state) {
    state.creating = true
    state.loaded = false
    state.loading = false
    state.error = null
  },

  [LOADING] (state) {
    state.creating = false
    state.loading = true
    state.error = null
  },

  [LOADED] (state, experiment) {
    state.creating = false
    state.loading = false
    state.loaded = true
    state.data = experiment
    state.error = null
  },

  [ERROR] (state, error) {
    state.creating = false
    state.loading = false
    state.loaded = false
    state.error = error
  }
}

const actions = {
  createExperiment ({ commit }, { modelID, experimentParams }) {
    commit(CREATING)
    return api
      .createExperiment(
        modelID, experimentParams,
        (data) => commit(LOADED, data),
        (error) => commit(ERROR, { error })
      )
  },
  loadExperiment ({ commit }, { modelID, experimentID }) {
    commit(LOADING)
    return api
      .loadExperiment(
        modelID, experimentID,
        (data) => commit(LOADED, data),
        (error) => commit(ERROR, { error })
      )
  }
}

export default {
  state: initialState,
  mutations,
  actions
}
