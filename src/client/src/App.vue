<template lang="pug">
#app
  .row.header
    nav
      ul
        li
          router-link(to="/") SAN Lite-Solver
        li
          router-link(to="/models/new") New Model
  .clear
  router-view
</template>

<script>
export default {
  name: 'app'
}
</script>

<style lang="sass" src="assets/styles.sass"></style>
