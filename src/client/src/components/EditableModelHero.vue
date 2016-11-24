<template lang="pug">
.hero
  template(v-if="loading")
    h1
      .load.small
      em  Loading...
  template(v-else)
    template(v-if="editingTitle")
      input.full(type="text",ref="title",v-model="title",@keyup.enter="doneEditingTitle",@keyup.esc="cancelEditingTitle")
    template(v-else)
      div(@click="editTitle")
        h1(v-if="model.title")
          | {{ model.title }}
        h1(v-else)
          em No Title
    div
      template(v-if="editingDescription")
        input.full(type="text",ref="description",v-model="description",@keyup.enter="doneEditingDescription",@keyup.esc="cancelEditingDescription")
      template(v-else)
        div(@click="editDescription")
          p(v-if="model.description")
            | {{ model.description }}
          p(v-else)
            em No description
</template>

<script>
export default {
  name: 'editable-model-hero',
  props: {
    model: { type: Object, required: true },
    loading: { type: Boolean, required: true }
  },
  data () {
    return {
      editingTitle: false,
      editingDescription: false,
      title: '',
      description: ''
    }
  },
  methods: {
    saveModel (model) {
      this.$emit('save', model)
    },

    editTitle () {
      this.title = this.model.title
      this.editingTitle = true
      this.$nextTick(() => {
        this.$refs.title.focus()
      })
    },
    doneEditingTitle () {
      if (this.savingModel) return
      this.saveModel({ title: this.title })
      this.editingTitle = false
    },
    cancelEditingTitle () {
      this.editingTitle = false
    },

    editDescription () {
      this.description = this.model.description
      this.editingDescription = true
      this.$nextTick(() => {
        this.$refs.description.focus()
      })
    },
    doneEditingDescription () {
      if (this.savingModel) return
      this.saveModel({ description: this.description })
      this.editingDescription = false
    },
    cancelEditingDescription () {
      this.editingDescription = false
    }
  }
}
</script>

<style lang="sass" scoped>
@import '~assets/variables';
.hero
  background: $white
  margin-top: 0
  padding: 1em
  h1
    font-weight: 400
    color: $black
    text-align: left
    font-size: 25px
    display: inline-block
  h1:hover, p:hover
    background: lighten($yellow, 25%)
    cursor: pointer
    display: inline-block
  p
    font-size: 15px
    font-style: italic
    color: $gray
    margin: 0
    display: inline-block
</style>
