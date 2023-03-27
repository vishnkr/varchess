<template>
  <v-overlay :value="loading" :opacity="0.8">
    <v-card>
      <v-card-text>
        Loading...
      </v-card-text>
      <v-card-actions>
        <v-btn color="primary" @click="loading = false">Cancel</v-btn>
      </v-card-actions>
    </v-card>
  </v-overlay>
</template>

<script>
export default {
  props: {
    loading: {
      type: Boolean,
      required: true
    }
  }
}
</script>