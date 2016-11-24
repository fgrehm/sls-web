<template lang="pug">
div
  editable-model-hero(:model="model",:loading="loading",@save="saveModel")

  .contents
    template(v-if="loading")
      p.center
        .load.small
        em  Loading...
    template(v-else)
      h2 New Experiment
      form.full-width-form(@submit.prevent="submitExperiment")
        fieldset
          legend Identifier Parameters
          .row
            .col-1-4
              label
                | Name
                select.full(v-model="experimentParams.identifier.name",@change="onIdentifierChanged")
                  option(value="") Select an identifier
                  option(v-for="id in modelConstantIdentifiers",:value="id.name") {{ id.name }}
            .col-1-4
              label
                | From
                input.full(v-model.number="experimentParams.identifier.from",:disabled="!experimentParams.identifier.name")
            .col-1-4
              label
                | To
                input.full(v-model.number="experimentParams.identifier.to",:disabled="!experimentParams.identifier.name")
            .col-1-4
              label
                | Increment
                input.full(v-model.number="experimentParams.identifier.increment",:disabled="!experimentParams.identifier.name")
        fieldset
          legend Solution Parameters
          .row
            .col-1-3
              label
                | Max Iterations
                input.full(v-model.number="experimentParams.maxIterations",placeholder="Default: 1000000")
            .col-1-3
              label
                | Tolerance
                input.full(v-model.number="experimentParams.tolerance",placeholder="Default: 1e-10")
            .col-1-3
              br
              button.btn.green.solid.float-right(:disabled="!canSubmitForm")
                span(v-if="creating")
                  span.load.small
                  |  Processing...
                span(v-else)
                  i.fa.fa-flask
                  |  Create
    .clear
</template>

<script>
import { mapState, mapActions } from 'vuex'
import EditableModelHero from 'components/EditableModelHero'

export default {
  name: 'experiments-new',
  components: {
    EditableModelHero
  },
  created () {
    const { modelID } = this.$route.params
    if (!this.loaded || this.model.id !== modelID) {
      this.loadModel(modelID)
    }
  },
  data () {
    return {
      experimentParams: {
        identifier: {
          name: '',
          from: null,
          to: null,
          increment: null
        },
        maxIterations: null,
        tolerance: null
      }
    }
  },
  computed: {
    ...mapState({
      loading: ({ selectedModel }) => selectedModel.loading,
      loaded: ({ selectedModel }) => selectedModel.loaded,
      model: ({ selectedModel }) => selectedModel.data,
      experiment: ({ experiment }) => experiment.data,
      creating: ({ experiment }) => experiment.creating
    }),
    modelConstantIdentifiers () {
      if (!this.loaded) return []
      return this.model.parsedModel.identifiers.filter((id) => (
        id.type === 'constant'
      ))
    },
    canSubmitForm () {
      const { identifier } = this.experimentParams
      return !this.creating && (identifier.name && identifier.from && identifier.to && identifier.increment)
    }
  },
  methods: {
    ...mapActions([
      'loadModel',
      'saveModel',
      'createExperiment'
    ]),
    onIdentifierChanged () {
      if (this.experimentParams.identifier.from) return
      const { identifier } = this.experimentParams
      for (let i = 0; i < this.modelConstantIdentifiers.length; i++) {
        const id = this.modelConstantIdentifiers[i]
        if (id.name !== identifier.name) {
          continue
        }
        identifier.from = id.value
        break
      }
    },
    submitExperiment () {
      this.createExperiment({
        modelID: this.model.id,
        experimentParams: this.experimentParams
      })
    }
  },
  // REFACTOR: Ideally we should not rely on watchers for this
  watch: {
    experiment (newExperiment) {
      const { modelID } = this.$route.params
      const experimentID = this.experiment.id
      this.$router.push({ name: 'showExperiment', params: { modelID, experimentID } })
    }
  }
}
</script>

<style lang="sass" scoped>
h2
  font-weight: bold
  margin: 0
.col-1-3, .col-1-4
  padding: 5px
</style>
