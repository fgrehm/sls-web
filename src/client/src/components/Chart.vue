<template lang="pug">
#chart
</template>

<script>
const Highcharts = require('highcharts')

export default {
  name: 'chart',
  props: {
    options: { type: Object, required: true }
  },
  mounted () {
    require('highcharts/modules/exporting')(Highcharts)
    this.chart = Highcharts.chart('chart', this.options)
  },
  watch: {
    options (newOptions) {
      this.chart = Highcharts.chart('chart', newOptions)
    }
  },
  beforeDestroy () {
    if (this.chart) {
      this.chart.destroy()
    }
  }
}
</script>
