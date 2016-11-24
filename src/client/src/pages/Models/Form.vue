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
        p
          template(v-if="model.id")
            router-link.btn.solid.blue.small(:to="{ name: 'newSolution', params: { modelID: model.id } }")
              i.fa.fa-play
              |  Solve
            | &nbsp;
            router-link.btn.solid.yellow.small(:to="{ name: 'newExperiment', params: { modelID: model.id } }")
              i.fa.fa-flask
              |  Experiment
            | &nbsp;
          button.btn.solid.green.small(href="#",:disabled="!dirty || saving",@click="doneEditingSource")
            i.fa.fa-save
            |  Save
        template(v-if="loaded")
          ace-editor(:source="model.source",@changed="updateModelSource")
      .col-1-2
        p.center
          button.btn.solid.small(@click="renderModelGraph",:disabled="rendered")
            i.fa.fa-file-image-o
            |  Render
        div
          pre(v-show="error")
            | {{ error && error.message }}
            | {{ error && error.response && error.response.data && error.response.data.error }}
          p.center(v-show="rendering")
            span.load.small
            em  Rendering...
          p.center(v-show="rendered")
            img(:src="model.graphUrl")
    .clear
</template>

<script>
import { mapState, mapActions } from 'vuex'

import EditableModelHero from 'components/EditableModelHero'
import AceEditor from 'components/AceEditor'

export default {
  name: 'models-form',
  components: {
    AceEditor,
    EditableModelHero
  },
  computed: {
    ...mapState({
      loading: ({ selectedModel }) => selectedModel.loading,
      loaded: ({ selectedModel }) => selectedModel.loaded,
      model: ({ selectedModel }) => selectedModel.data,
      error: ({ selectedModel }) => selectedModel.error,
      dirty: ({ selectedModel }) => selectedModel.dirty,
      saving: ({ selectedModel }) => selectedModel.saving,
      source: ({ selectedModel }) => selectedModel.data.source,
      rendering: ({ selectedModel }) => selectedModel.rendering,
      rendered: ({ selectedModel }) => selectedModel.rendered
    })
  },
  methods: {
    ...mapActions([
      'updateModelSource',
      'saveModel',
      'renderModelGraph'
    ]),
    doneEditingSource () {
      if (!this.saving) {
        this.saveModel({ source: this.source })
      }
    }
  }
}
</script>

<style lang="sass" scoped>
.col-1-2:nth-child(1)
  padding-right: 10px

img
  max-width: 100%
</style>
