# Tic Tac Toe Squared Server

A server for the Tic Tac Toe Squared (AKA Meta Tic Tac Toe)

This is written in Go

### Tech Stack:
 *  MySQL database to store game and user information
 *  Go server to interface with database and respond to requests
 *  NGINX server to reverse proxy (so separate apps are separate processes)
 *  RESTful API as well as Websockets (to update when game changes)

## API
    URL | Function
    --- | --------
    GET /games | lists all games
    GET /games/{ID}/info | lists info on game with specified ID
    GET /games/{ID}/board | gives board in JSON
    GET /games/{ID}/string | gives board in string format (use monospaced font)
    GET /games/{ID}/ws | gives websocket that broadcasts when a game changes
    POST /games?Player1={PID1}&Player2={PID2} | makes a new game with specified ID's and returns the ID of the game created
    POST /games/{ID}/move?Player={PID}&Box={BID}&Square={SID} | makes a move and responds with an error if unsucessful; broadcasts on ws if succesful 

## Class Organization
 Image will be created later

