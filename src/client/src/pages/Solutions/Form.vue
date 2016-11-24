<template lang="pug">
div
  editable-model-hero(:model="model",:loading="loading",@save="saveModel")

  .contents
    template(v-if="loading")
      p.center
        .load.small
        em  Loading...
    template(v-else)
      .col-1-2
        form.full-width-forms(@submit.prevent="submitSolution")
          fieldset
            legend Solution parameters
            .col-1-2
              label
                | Iterations
                input(placeholder="Default: 1000000",v-model.number="solutionParams.maxIterations")
            .col-1-2
              label
                | Tolerance
                input(placeholder="Default: 1e-10",v-model.number="solutionParams.tolerance")
          fieldset
            legend Custom identifiers
            template(v-for="idPair in modelIdentifiersInPairs")
              .col-1-2(v-for="id in idPair")
                identifier-input(:identifier="id",v-model.number="solutionParams.customIdentifiers[id.name]")
              .clear
          button.btn.blue.solid.full(:disabled="solving")
            .load.small(v-if="solving")
            i.fa.fa-play(v-else)
            | {{ solving && ' Solving...' || ' Solve' }}
      .col-1-2
        fieldset
          legend Results
          pre(v-if="error")
            | {{ error && error.message }}
            | {{ error && error.response && error.response.data && error.response.data.error }}
          template(v-if="!solutionAttempted")
            .notice
              p.center Waiting for execution.
          template(v-else)
            template(v-if="solving")
              p.center
                .load.small
                em  Processing
            template(v-if="solved")
              template(v-if="!solution.errored")
                h3.center
                  | Solution
                  strong.found(v-show="solution.found")  found
                  strong.not-found(v-show="!solution.found")  not found
                  |  in
                  em  {{solution.executionTime.toFixed(3)}}ms
                  |  ({{solution.steps}} steps)
                chart(:options="chartOptions")
                table
                  thead
                    tr
                      th Function
                      th Result
                  tbody
                    tr(v-for="result in solution.results")
                      td {{result.label}}
                      td {{result.value}}
              template(v-if="solution.errored")
                .notice.error
                  p.center
                    i.fa.fa-warning
                    | The solution attempt has errored!
      .clear
</template>

<script>
import { reduce as _reduce, map as _map } from 'lodash'
import { mapState, mapActions } from 'vuex'
import EditableModelHero from 'components/EditableModelHero'
import IdentifierInput from 'components/IdentifierInput'
import Chart from 'components/Chart'

export default {
  name: 'solutions-form',
  components: {
    EditableModelHero,
    IdentifierInput,
    Chart
  },
  created () {
    const { modelID } = this.$route.params
    if (!this.loaded || this.model.id !== modelID) {
      this.loadModel(modelID)
    }
  },
  data () {
    return {
      solutionParams: {
        maxIterations: null,
        tolerance: null,
        customIdentifiers: {}
      },
      solutionAttempted: false
    }
  },
  computed: {
    ...mapState({
      loading: ({ selectedModel }) => selectedModel.loading,
      loaded: ({ selectedModel }) => selectedModel.loaded,
      model: ({ selectedModel }) => selectedModel.data,
      solution: ({ solution }) => solution.data,
      solving: ({ solution }) => solution.solving,
      solved: ({ solution }) => solution.solved,
      error: ({ selectedModel, solution }) => (selectedModel.error || solution.error)
    }),
    modelIdentifiersInPairs () {
      if (!this.loaded) return []

      const modelConstantIdentifiers = this.model.parsedModel.identifiers.filter((id) => (
        id.type === 'constant'
      ))

      let [row, idx] = [0, 0]
      return _reduce(modelConstantIdentifiers, (groups, identifier) => {
        if (!groups[row]) groups[row] = []

        groups[row].push(identifier)
        if (idx === 1) {
          idx = 0; row++
        } else {
          idx++
        }

        return groups
      }, [])
    },
    chartOptions () {
      if (!this.solved) return { }

      const data = _map(this.solution.results, (res) => (
        [ res.label, res.value ]
      ))

      return {
        chart: { type: 'column' },
        title: { text: null },
        xAxis: {
          type: 'category',
          labels: {
            rotation: -45,
            style: { fontSize: '13px', fontFamily: 'Verdana, sans-serif' }
          }
        },
        yAxis: { min: 0, max: 1, title: { text: null } },
        legend: { enabled: false },
        series: [{
          name: 'Result',
          data,
          dataLabels: {
            enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.4f}',
            y: 10, // 10 pixels down from the top
            style: { fontSize: '13px', fontFamily: 'Verdana, sans-serif' }
          }
        }]
      }
    }
  },
  methods: {
    ...mapActions([
      'loadModel',
      'solveModel',
      'saveModel'
    ]),
    submitSolution () {
      this.solutionAttempted = true
      this.solveModel({
        modelID: this.model.id,
        solutionParams: this.solutionParams
      })
    }
  }
}
</script>

<style lang="sass" scoped>
@import '~assets/variables';

form .col-1-2
  padding: 5px

.col-1-2:nth-child(1)
  padding-right: 10px

.found
  color: darken($green, 5%)
  font-weight: 600

.not-found
  color: $red
  font-weight: bold
</style>
