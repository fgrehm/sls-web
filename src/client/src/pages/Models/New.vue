<template lang="pug">
model-form
</template>

<script>
import { mapState } from 'vuex'
import Form from './Form'

export default {
  created () {
    this.$store.dispatch('initializeNewModel')
  },
  components: {
    'model-form': Form
  },
  computed: {
    ...mapState({
      model: ({ selectedModel }) => selectedModel.data
    })
  },
  watch: {
    model () {
      if (this.model.id !== null) {
        this.$router.push({ name: 'editModel', params: { modelID: this.model.id } })
      }
    }
  }
}
</script>
