<template>
  <div class="home">
    <div class="info-card">
      <v-card class="mx-auto"
    max-width="344">
    <v-card-text>
      <p class="display-1 text--primary">
        VarChess
      </p>
      <v-container>
        <v-row>
        <v-col
          cols="12"
        >
      <div class="headline">Create your own chess variants and play with friends! The perfect place to experiment your playing style with varying board sizes, customized piece movements and placements. More features coming soon.</div>
      </v-col>
      <v-col
          cols="12"
        >
          <v-text-field
            v-model="username"
            label="Enter Username"
            solo-inverted
          ></v-text-field>
        </v-col>
      </v-row>
      <v-dialog
        transition="dialog-bottom-transition"
        max-width="600"
      >
        <template v-slot:activator="{ on, attrs }">
          <div class="text-center">
            <v-btn
              rounded
              color="primary"
              v-bind="attrs"
              v-on="on"
              dark>
              Create Room
            </v-btn>
          </div>
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-toolbar
              color="primary"
              dark
            ><h3>Choose Game Setup</h3></v-toolbar>
             
            <v-card-actions class="justify-end">
              <v-radio-group
                    v-model="mode"
                    row
                  >
                    <v-radio
                      label="Standard 8x8 board"
                      value="standard"
                    ></v-radio>
                    <v-radio
                      label="Custom Variant"
                      value="custom"
                    ></v-radio>
                  </v-radio-group>
              <v-btn
                color="error"
                @click="dialog.value=false"
              >Cancel</v-btn>
              <v-btn
                color="success"
                @click="enterRoom"
              >Enter Room</v-btn>
            </v-card-actions>
          </v-card>
        </template>
        </v-dialog>
      </v-container>
    </v-card-text>
    </v-card>
  </div>
  </div>
</template>

<script>
//import EditorBoard from '../components/Editor/EditorBoard.vue';

//import EditorDialog from '../components/EditorDialog'
export default {
  components:{},//EditorBoard},//EditorDialog},
  methods:{
    redirectToEditor() {
      this.$router.push({ path: '/editor' });
    },

    enterRoom(){
      console.log(this.mode)
      this.mode=='custom'? this.redirectToEditor() : this.$router.push({ path: '/game' }) ;
    }
  },
  data:()=>{
    return {
      createClicked: false,
      username: '',
      mode:'standard',
    }
  }
}
</script>

<style scoped>

</style>