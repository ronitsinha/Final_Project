# Final Project: A Doom Clone

This is the final project for Game Programming at Concord Academy. My final project is a clone of Doom/Wolfenstein 3D using raycasting for psuedo-3D. 
The objective of this game is, however, much simpler than Doom or Wolfenstein 3D: To shoot all the targets in a given amount of time.

Controls:
W- forward
S- backward
A- rotate counter-clockwise
D- rotate clockwise
TAB- Open top-down view
1,2,3- switch between weapons (chainsaw, pistol, shotgun, respectively)
Left click to use weapon (hold down for chainsaw)

Here's a screenshot of my progress so far:

![Alt text](./assets/images/screenshot.png?raw=true)

The DOOM weapon sprites should be creditted to Falcon Delta,
and the raycasting uses a Go binding of SFML found
here: https://github.com/zyedidia/sfml

The raycasting is based on Lode Vandevenne's tutorial on raycasting using
QuickCG in C++, which can be found here: http://lodev.org/cgtutor/raycasting.html