import axios from 'axios'

const API_HOST = process.env.API_HOST

export default {
  allModels (cb, errorCb) {
    return axios.get(`${API_HOST}/models`)
                .then(({ data }) => cb(data))
                .catch(errorCb)
  },
  findModel (modelID, cb, errorCb) {
    return axios.get(`${API_HOST}/models/${modelID}`)
                .then(({ data }) => cb(data))
                .catch(errorCb)
  },
  createModel (params, cb, errorCb) {
    return axios.post(`${API_HOST}/models`, params)
                .then(({ data }) => cb(data))
                .catch(errorCb)
  },
  updateModel (id, params, cb, errorCb) {
    return axios.put(`${API_HOST}/models/${id}`, params)
                .then(({ data }) => cb(data))
                .catch(errorCb)
  },
  solveModel (modelID, solutionParams, cb, errorCb) {
    return axios.post(`${API_HOST}/models/${modelID}/solutions`, solutionParams)
                .then(({ data }) => cb(data))
                .catch(errorCb)
  },
  createExperiment (modelID, experimentParams, cb, errorCb) {
    return axios.post(`${API_HOST}/models/${modelID}/experiments`, experimentParams)
                .then(({ data }) => cb(data))
                .catch(errorCb)
  },
  parseModel (source, cb, errorCb) {
    return axios.post(`${API_HOST}/parse`, source, { headers: { 'Content-Type': 'text/plain' } })
                .then(({ data }) => cb(data))
                .catch(errorCb)
  },
  loadExperiment (modelID, experimentID, cb, errorCb) {
    return axios.get(`${API_HOST}/models/${modelID}/experiments/${experimentID}`)
                .then(({ data }) => cb(data))
                .catch(errorCb)
  },
  renderGraph (source, cb, errorCb) {
    return axios.post(`${API_HOST}/graph`, source, { headers: { 'Content-Type': 'text/plain' } })
                .then(({ data }) => cb(data))
                .catch(errorCb)
  }
}
