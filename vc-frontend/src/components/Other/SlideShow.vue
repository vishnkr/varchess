<template>
    <div ref="container" class="keen-slider slider">
       <div class="keen-slider__slide number-slide1">
                <div class="slide-desc">
                    <div class="slide-img">
                        <img src="../../assets/mac.png"/>
                    </div>
                    <div class="slide-text">
                        Play on Varying Board Sizes from 5x5 to 16x16
                    </div>
                </div>
       </div>
       <div class="keen-slider__slide number-slide2"> 
                <div class="slide-desc">
                    <div class="slide-text">
                        Disable squares to restrict movement on certain parts of the board
                    </div>
                    <div class="slide-img">
                        <img src="../../assets/ipad.png"/>
                    </div>
                </div>
        </div>
       <div class="keen-slider__slide number-slide3"> 
            <div class="slide-desc">
                <div class="slide-text">
                    Create Your Own Chess Pieces with Custom Move Patterns
                </div>
                <div class="slide-img">
                    <img src="../../assets/mp2.png"/>
                </div>
            </div>
       </div>
       <div class="keen-slider__slide number-slide4"> Create and share variant templates </div>
       <div class="keen-slider__slide number-slide5"> Play against Stonkfish!</div>
     </div>
   </template>
   
   <script lang="ts">
   import { useKeenSlider } from 'keen-slider/vue.es'
   import 'keen-slider/keen-slider.min.css'
   export default {
    setup() {
    const [container] = useKeenSlider(
      {
        loop: true,
      },
      [
        (slider) => {
          let timeout: NodeJS.Timeout;
          let mouseOver = false
          function clearNextTimeout() {
            clearTimeout(timeout)
          }
          function nextTimeout() {
            clearTimeout(timeout)
            if (mouseOver) return
            timeout = setTimeout(() => {
              slider.next()
            }, 8000)
          }
          slider.on("created",nextTimeout)
          slider.on("dragStarted", clearNextTimeout)
          slider.on("animationEnded", nextTimeout)
          slider.on("updated", nextTimeout)
        },
        ]
        )
        return { container }
    },
   }
   </script>
   
   <style >
   body {
     margin: 0;
     font-family: 'Inter', sans-serif;
     -webkit-font-smoothing: antialiased;
     -moz-osx-font-smoothing: grayscale;
   }
   .slider{
    height: 100;
   }
   .slide-title{
        display: flex;

   }
   .slide-desc{
        display: flex;
        justify-content: center;
        flex-direction: row;
    }
    .slide-img{
       display: flex;
       max-width: 60%;
    }

   .slide-img img{
        max-width: 100%;
        max-height: 100%;
   }

   .slide-text{
        max-width:40%;
        display: flex;
        text-align: center;
        align-items: center;
        justify-content: center;
        font-size: 35px;
        color: #fff;
        font-weight: 500;
   }

   [class^="number-slide"],
   [class*=" number-slide"] {
     background: grey;
     overflow-y: hidden;
     display: flex;
     flex-direction: column;
     align-items: center;
     justify-content: center;
     font-size: 50px;
     color: #fff;
     font-weight: 500;
     max-height: 50vh;
   }
   .number-slide1 {
     background: rgb(64, 175, 255);
     background: linear-gradient(
       128deg,
       rgba(64, 175, 255, 1) 0%,
       rgba(63, 97, 255, 1) 100%
     );
   }
   .number-slide2 {
     background: rgb(255, 75, 64);
     background: linear-gradient(
       128deg,
       rgba(255, 154, 63, 1) 0%,
       rgba(255, 75, 64, 1) 100%
     );
   }
   .number-slide6 {
     background: rgb(182, 255, 64);
     background: linear-gradient(
       128deg,
       rgb(157, 241, 23) 0%,
       rgba(63, 255, 71, 1) 100%
     );
     background: linear-gradient(
       128deg,
       rgba(189, 255, 83, 1) 0%,
       rgba(43, 250, 82, 1) 100%
     );
   }
   .number-slide4 {
     background: rgb(45, 226, 214);
     background: linear-gradient(
       128deg,
       rgba(64, 255, 242, 1) 0%,
       rgba(63, 188, 255, 1) 100%
     );
   }
   .number-slide5 {
     background: rgb(255, 64, 156);
     background: linear-gradient(
       128deg,
       rgba(255, 64, 156, 1) 0%,
       rgba(255, 63, 63, 1) 100%
     );
   }
   .number-slide3 {
     background: rgb(64, 76, 255);
     background: linear-gradient(
       128deg,
       rgba(64, 76, 255, 1) 0%,
       rgba(174, 63, 255, 1) 100%
     );
   }
</style>