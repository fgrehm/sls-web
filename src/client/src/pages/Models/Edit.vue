<template lang="pug">
model-form
</template>

<script>
import { mapState } from 'vuex'
import Form from './Form'

export default {
  created () {
    const { modelID } = this.$route.params
    if (!this.modelLoaded || this.modelId !== modelID) {
      this.$store.dispatch('loadModel', modelID)
    }
  },
  components: {
    'model-form': Form
  },
  computed: {
    ...mapState({
      modelId: ({ selectedModel }) => selectedModel.data.id,
      modelLoaded: ({ selectedModel }) => selectedModel.loaded
    })
  }
}
</script>
