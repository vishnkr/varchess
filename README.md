# Varchess
![GitHub](https://img.shields.io/github/license/vishnkr/varchess)
## Overview
Varchess is a multiplayer chess website where you can create your own chess variants. Customizations include variable board dimensions, piece placements, custom pieces and movement patterns.

#### Built using Vue.js (frontend), Go (websocket server and game logic).

#### LIVE LINK - https://varchess.tech/


![Editor](https://i.imgur.com/aWI8KoW.png)


![Game Room](https://i.imgur.com/eoTxp7S.png)

Future Plans : 
- More pre-defined variants (teleportation portals, golem chess, long range attacks, poisoned pawn/hidden queen)
- Chess engine to mimic computer opponent (In the works - [Stonkfish](https://github.com/vishnkr/stonkfish))

## Running Locally
Clone this project and cd into the local repo.
```
git clone git@github.com:vishnkr/varchess.git
```
To run the application using Docker (docker-compose):

Ignore --build flag after first run
```
docker-compose up --build 
```
To stop the containers:
```
docker-compose down
```
The web app will be running locally at localhost:8080
