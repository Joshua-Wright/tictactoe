#Tic Tac Toe

This is a re-creation of https://xkcd.com/832/ using golang and svg.

[Link to finished product](http://i.imgur.com/QY9LEcr.jpg) (30MG PNG file rendered from svg output file)

I used a minimax tree algorithm to generate all the moves. 
No pruning of the tree was used, as it was not necessary for the 3x3 case. 
While this code should work for higher-order tic-tac-toe, my laptop (with 20GB ram) runs our of memory when trying to compute it. 
Pruning the tree could help with that case, but that would reduce the completeness of he resulting map.
