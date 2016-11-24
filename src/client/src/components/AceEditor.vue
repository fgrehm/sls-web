<template lang="pug">
#ace-editor
</template>

<script>
import { debounce } from 'lodash'

export default {
  name: 'ace-editor',
  props: {
    source: {
      type: String,
      required: true
    }
  },
  mounted () {
    const ace = require('brace')
    require('brace/mode/text')
    require('brace/theme/mono_industrial')

    const editor = ace.edit('ace-editor')
    editor.$blockScrolling = Infinity
    editor.setTheme('ace/theme/mono_industrial')
    editor.getSession().setMode('ace/mode/text')
    editor.setShowPrintMargin(false)
    editor.setFontSize(14)
    editor.setValue(this.source)
    editor.clearSelection()
    const debouncedOnChange = debounce(() => {
      this.$emit('changed', this.editor.getValue())
    }, 1000)
    editor.on('change', debouncedOnChange)
    this.editor = editor
  },
  beforeDestroy () {
    this.editor.destroy()
  }
}
</script>

<style lang="sass" scoped>
#ace-editor
  width: 100%
  height: 500px
</style>
