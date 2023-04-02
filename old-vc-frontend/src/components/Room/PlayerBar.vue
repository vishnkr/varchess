<template>
<div>
<v-card flat>
    <v-container fluid>
      <v-row class="child-flex">
        <div>
          <v-toolbar height="50%">
            <v-btn
              icon
              class="hidden-xs-only"
            >
              <v-icon>mdi-arrow-left</v-icon>
            </v-btn>

            <v-toolbar-title>Title</v-toolbar-title>

            <v-spacer></v-spacer>

            <v-btn
              icon
              class="hidden-xs-only"
            >
              <v-icon>mdi-magnify</v-icon>
            </v-btn>
          </v-toolbar>
        </div>

        <div style="flex-basis: 20%">
          <v-toolbar dark height="50%">
            <v-spacer></v-spacer>
          </v-toolbar>
        </div>
      </v-row>
    </v-container>
  </v-card>

  </div>
</template>

<script>
/*
<div>

        <v-tool-bar floating>
            
            <v-avatar color="gray">
                <img
                    :src="getRandomPieceURL"
                >
            ></v-avatar>
            <v-toolbar-title>{{username}}</v-toolbar-title>
            <v-spacer></v-spacer>
        </v-tool-bar>
    </div>
 <div>
              <v-toolbar>
                <v-btn
                  icon
                  class="hidden-xs-only"
                >
                  <v-avatar>
                    <img
                    :src="getRandomPieceURL"
                    alt="color"
                  >
                  </v-avatar>
                </v-btn>

                <v-toolbar-title>username</v-toolbar-title>
              </v-toolbar>
            </div>*/
export default {/*
    props:['color','username'],
    methods:{
        getRandomPieceURL(){ 
            var pieces = ['p','q','r','b','k','n'];
            var pieceType = pieces[Math.floor(Math.random() * pieces.length)];
            console.log(this.color,pieceType)
            if (this.color==null){
                return
            }
            return require(`../../assets/images/pieces/${this.color}/${pieceType}.svg`)
        },
    }*/
}
</script>

<style>
img {
        -khtml-user-select: none;
        -o-user-select: none;
        -moz-user-select: none;
        -webkit-user-select: none;
        user-select: none;
    }
</style>