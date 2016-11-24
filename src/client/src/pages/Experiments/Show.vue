<template lang="pug">
div
  editable-model-hero(:model="model",:loading="loadingModel",@save="saveModel")
  .contents
    template(v-if="!experimentLoaded")
      p.center
        .load.small
        em  Loading...
    template(v-else)
      h3
        | Experimenting with
        strong.identifier  {{ experiment.identifier.name }}
        |  from
        strong.from  {{ experiment.identifier.from }}
        |  to
        strong.to  {{ experiment.identifier.to }}
        |  in increments of
        strong.to  {{ experiment.identifier.increment }}
      chart(:options="chartOptions")
      table
        thead
          tr
            th Identifier Value
            th(v-for="(exp, label) in integrationFunctions",:title="exp")
              | {{ label }}
        tbody
          tr(v-for="solution in experiment.solutions",:title="formatSolutionTooltip(solution)",:class="{ found: solution.found, 'not-found': (solution.found == false)}")
            template(v-if="solution.finishedAt")
              td
                strong {{ solution.identifierValue }}
              td(v-for="result in solution.results")
                | {{ result.value | formatResultValue }}
            template(v-else)
              td
                | {{ solution.identifierValue }}
              td(:colspan="integrationFunctionsLength")
                .notice
                  .load.small
                  em  Processing...

    .clear
</template>

<script>
import { every as _every, keys as _keys, map as _map, each as _each } from 'lodash'
import { mapState, mapActions } from 'vuex'
import EditableModelHero from 'components/EditableModelHero'
import Chart from 'components/Chart'

export default {
  name: 'experiments-show',
  components: {
    EditableModelHero,
    Chart
  },
  created () {
    const { modelID, experimentID } = this.$route.params
    if (!this.modelLoaded || this.model.id !== modelID) {
      this.loadModel(modelID)
    }
    if (!this.experimentLoaded || this.experiment.id !== experimentID) {
      this.loadExperiment({ modelID, experimentID })
          .then(() => this.waitForExperimentToComplete())
    } else {
      this.waitForExperimentToComplete()
    }
  },
  computed: {
    ...mapState({
      loadingModel: ({ selectedModel }) => selectedModel.loading,
      modelLoaded: ({ selectedModel }) => selectedModel.loaded,
      model: ({ selectedModel }) => selectedModel.data,
      loadingExperiment: ({ experiment }) => experiment.loading,
      experimentLoaded: ({ experiment }) => experiment.loaded,
      experiment: ({ experiment }) => experiment.data,
      integrationFunctions: ({ experiment }) => experiment.data.integrationFunctions,
      integrationFunctionsLength: ({ experiment }) => _keys(experiment.data.integrationFunctions).length
    }),
    experimentCompleted () {
      if (!this.experimentLoaded) return false

      const { solutions } = this.experiment
      return solutions && solutions.length > 0 && _every(solutions, (sol) => (
        sol.finishedAt !== null
      ))
    },
    chartOptions () {
      if (!this.experimentCompleted) return { }

      const { solutions } = this.experiment
      const categories = _map(solutions, (sol) => sol.identifierValue)

      let series = { }
      _each(solutions, (sol) => {
        // REFACTOR: This only works because things are ordered but ideally
        // it should be a bit more robust
        _each(sol.results, (res) => {
          if (!series[res.label]) series[res.label] = []
          series[res.label].push(res.value)
        })
      })
      series = _map(series, (v, k) => ({ name: k, data: v }))

      return {
        title: { text: null },
        xAxis: { categories },
        yAxis: { title: { text: null } },
        legend: {
          layout: 'vertical',
          align: 'right',
          verticalAlign: 'middle',
          borderWidth: 0
        },
        series
      }
    }
  },
  filters: {
    formatResultValue (value) {
      return value.toPrecision(6)
    }
  },
  methods: {
    ...mapActions([
      'loadModel',
      'loadExperiment',
      'saveModel'
    ]),
    waitForExperimentToComplete () {
      if (this.experimentCompleted) return

      this.experimentCompletionTimeout = setTimeout(() => {
        const { modelID, experimentID } = this.$route.params
        this.loadExperiment({ modelID, experimentID })
            .then(() => {
              if (!this.experimentCompleted) {
                this.waitForExperimentToComplete(this.experiment)
              }
            })
      }, 1000)
    },
    onIdentifierChanged () {
      if (this.experimentParams.identifier.from) return

      const { identifier } = this.experimentParams
      for (let i = 0; i < this.availableIdentifiers.length; i++) {
        const id = this.availableIdentifiers[i]
        if (id.name !== identifier.name) {
          continue
        }
        identifier.from = id.value
        break
      }
    },
    formatSolutionTooltip (solution) {
      const { found, steps, finishedAt, executionTime } = solution
      if (!finishedAt) return ''
      const foundStr = found && 'found' || 'NOT found'
      return `Solution ${foundStr} in ${executionTime.toFixed(3)}ms (${steps} steps)`
    }
  },
  beforeDestroy () {
    clearTimeout(this.experimentCompletionTimeout)
  }
}
</script>

<style lang="sass" scoped>
@import '~assets/variables';
h2
  font-weight: bold
  margin: 0
.col-1-3, .col-1-4
  padding: 5px
table
  border: solid 1px $gray-lighter
  .notice
    background: lighten($notice, 5%)
    margin: 0
h3
  font-size: 1.5em
  background: $gray-light
  text-align: center
td, th
  font-size: 0.83em
  text-align: center
.found
  background: lighten($green, 25%)
.not-found
  background: lighten($red, 25%)
</style>
