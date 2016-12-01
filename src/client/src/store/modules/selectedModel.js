import api from '../../api'

const SOURCE_UPDATED = 'selectedModel/SOURCE_UPDATED'
const SAVING = 'selectedModel/SAVING'
const LOADING = 'selectedModel/LOADING'
const LOADED = 'selectedModel/LOADED'
const LOAD_ERROR = 'selectedModel/LOAD_ERROR'
const RENDERING = 'selectedModel/RENDERING'
const RENDERED = 'selectedModel/RENDERED'
const RENDER_ERROR = 'selectedModel/RENDER_ERROR'

const emptyModel = {
  url: null,
  id: null,
  source: '// New SAN Model',
  parsedModel: {},
  modelHash: null,
  transitionsHash: null,
  graphUrl: null,
  createdAt: null,
  updatedAt: null
}

const initialState = {
  loading: false,
  loaded: false,
  dirty: false,
  saving: false,
  rendering: false,
  rendered: false,
  data: { ...emptyModel },
  error: null
}

const mutations = {
  [LOADING] (state) {
    state.loading = true
  },

  [LOADED] (state, model) {
    state.dirty = false
    state.loading = false
    state.saving = false
    state.loaded = true
    state.rendering = false
    state.rendered = true
    state.data = model
    state.error = null
  },

  [LOAD_ERROR] (state, error) {
    state.loading = false
    state.error = error
  },

  [SOURCE_UPDATED] (state, newSource) {
    state.dirty = true
    state.data.source = newSource
  },

  [SAVING] (state) {
    state.saving = true
  },

  [RENDERING] (state) {
    state.rendering = true
    state.error = null
  },

  [RENDERED] (state, { url, transitionsHash }) {
    state.data.graphUrl = `data:image/png;base64,${url}`
    state.rendering = false
    state.rendered = true
  },

  [RENDER_ERROR] (state, error) {
    state.rendering = false
    state.rendered = false
    state.error = error
  }
}

const actions = {
  initializeNewModel ({ commit }) {
    commit(LOADED, { ...emptyModel })
  },
  parseModel ({ commit, dispatch, state }, cb) {
    // HACK: Should have mutations specific to parsing insted of using the render ones
    commit(RENDERING)
    api.parseModel(state.data.source,
      () => dispatch('renderModelGraph'),
      (error) => commit(RENDER_ERROR, error)
    )
  },
  updateModelSource ({ commit, dispatch }, newSource) {
    commit(SOURCE_UPDATED, newSource)
    dispatch('parseModel')
  },
  saveModel ({ commit, state }, params) {
    commit(SAVING)

    const { id } = state.data
    if (id) {
      api.updateModel(id, params, (data) => commit(LOADED, data))
    } else {
      api.createModel({...state.data, ...params},
        (data) => commit(LOADED, data)
      )
    }
  },
  loadModel ({ commit }, modelID) {
    commit(LOADING)
    api.findModel(
      modelID,
      (data) => commit(LOADED, data),
      (error) => commit(LOAD_ERROR, error)
    )
  },
  renderModelGraph ({ commit, state }) {
    commit(RENDERING)

    const { source } = state.data
    api.renderGraph(source,
      (data) => commit(RENDERED, data),
      (error) => commit(RENDER_ERROR, error)
    )
  }
}

export default {
  state: initialState,
  mutations,
  actions
}
