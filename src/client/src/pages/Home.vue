<template lang="pug">
.contents
  h1 Models
  p.center(v-if="loading")
    span.load.small
    |  Loading...
  ul(v-if="loaded")
    li.center(v-for="model in models")
      router-link(:to="{ name: 'editModel', params: { modelID: model.id } }")
        | {{ model | modelTitle }}
  .clear
  .notice.error(v-if="error")
    p Error loading models
    p
      button.btn.green.solid(@click="reload") Click here to try again
</template>

<script>
import { mapState, mapActions } from 'vuex'

export default {
  name: 'home',
  filters: {
    modelTitle (model) {
      return model.title && model.title || model.id
    }
  },
  created () {
    this.reload()
  },
  computed: {
    ...mapState({
      loading: ({ models }) => models.loading,
      loaded: ({ models }) => models.loaded,
      models: ({ models }) => models.data,
      error: ({ models }) => models.error
    })
  },
  methods: {
    ...mapActions({
      reload: 'loadModels'
    })
  }
}
</script>

<style scoped lang="sass">
@import '~assets/variables';

h1
  font-weight: bold
  margin-top: 0
  text-align: center
li a
  color: $gray
  display: block
  padding: 0.4em
  border: solid 1px $gray-lighter
  margin-bottom: -1px
  &:hover
    background: $gray-lighter
</style>
