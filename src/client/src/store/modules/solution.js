import axios from 'axios'
import { reduce as _reduce } from 'lodash'

import api from '../../api'

const SOLVING = 'solution/SOLVING'
const SOLVED = 'solution/SOLVED'
const SOLVE_ERROR = 'solution/SOLVE_ERROR'

const emptySolution = {
  executionTime: 0,
  errored: false,
  results: []
}

const initialState = {
  solving: false,
  solved: false,
  data: { ...emptySolution },
  error: null
}

const mutations = {
  [SOLVING] (state) {
    state.solving = true
    state.solved = false
    state.error = null
    state.data = { ...emptySolution }
  },

  [SOLVED] (state, solution) {
    state.solving = false
    state.solved = true
    state.data = solution
    state.error = null
  },

  [SOLVE_ERROR] (state, error) {
    state.solving = false
    state.solved = false
    state.error = error
  }
}

let solutionTimeout
function waitForSolutionToComplete (data, commit) {
  if (solutionTimeout) clearInterval(solutionTimeout)

  if (data.finishedAt) {
    commit(SOLVED, data)
    return
  }

  const url = data.url
  solutionTimeout = setTimeout(() => {
    axios
      .get(url)
      .then(({ data }) => {
        if (data.finishedAt) {
          commit(SOLVED, data)
        } else {
          waitForSolutionToComplete(data, commit)
        }
      })
      .catch(error => commit(SOLVE_ERROR, { error }))
  }, 700)
}

const actions = {
  solveModel ({ commit }, { modelID, solutionParams }) {
    if (solutionTimeout) clearInterval(solutionTimeout)

    // Clean up parameters for the request
    solutionParams.maxIterations = solutionParams.maxIterations || null
    solutionParams.tolerance = solutionParams.tolerance || null
    solutionParams.customIdentifiers =
      _reduce(solutionParams.customIdentifiers, (acc, value, key) => {
        if (value) acc[key] = value
        return acc
      }, {})

    commit(SOLVING)
    api.solveModel(
      modelID, solutionParams,
      (data) => waitForSolutionToComplete(data, commit),
      (error) => commit(SOLVE_ERROR, { error })
    )
  }
}

export default {
  state: initialState,
  mutations,
  actions
}
