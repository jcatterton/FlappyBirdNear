#Flappy Bird NEAT
A basic implementation of a "Flappy Bird" game which learns to progress using NEAT. This is simply the program to run
the game, not the NEAT library. The intention of this program is to give the NEAT library a very simple task for initial
testing.

To run, simply use `go run main.go` in terminal.

Due to the nature of Flappy Bird, the birds never actually move forward. Instead, the bird stays centered on the screen
while the pipes move left. The bird is affected by gravity, but is able to "flap" in order to apply an impulse force
upward. The bird "dies" if it touches the floor or a pipe. Pipes always appear in pairs, one growing from the top of
the screen and one from the bottom, with a gap in between. The size of the gap is consistent, although the height is
variable. Fitness is simply an additive function which increases over time.

The birds see using sight vectors, which are lines radiating away from the bird which check for certain statuses in the
game. One sight vector projects from the center of the bird to the location of the right-most pipe-pair on screen. One 
vector projects from the center of the bird directly up, until it matches the height of the top pipe of the right-most 
pipe-pair on screen. One vector projects from the center of the bird directly down, until it matches the higher of the
bottom pipe of the right-most pipe-pair on screen. Additionally, the bird is always aware of its currently vertical
velocity. The length of each sight vector, as well as the bird's vertical velocity, is given to the NEAT library, and
the output of the bird's corresponding genome is used to determine when the bird should flap.

Display Options:
* Enable/Disable Sight Vectors - Click 1
* Enable/Disable genome drawing - Click 2

Expected Performance:

Due to the simple nature of Flappy Bird. It is expected that a bird will be able to continue playing indefinitely after
only a few generations. Additionally, it is common that the first bird to successfully pass the first pipe-pair, is able
to continue indefinitely. This is due to there being only a single output (to flap or not to flap) and very little
variation in desired behavior (The bird should always try to stay between the top and bottom pipe of a pipe-pair).
