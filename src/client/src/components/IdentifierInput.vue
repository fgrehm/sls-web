<template lang="pug">
label
  | {{ identifier.name }}
  input(:placeholder="'Default: ' + identifier.value",:value="value",@input="onInput")
</template>

<script>
export default {
  props: ['identifier', 'value'],
  methods: {
    onInput (event) {
      this.$emit('input', event.target.value)
    }
  }
}
</script>
