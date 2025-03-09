## My Approach

### The Algorithm

For this project, I implemented Minimax with Alpha-Beta Pruning. The only difference is that before running the algorithm to generate the best move, it first checks for an immediate win or loss scenario. If such a scenario exists, it generates either a blocking move or a winning move accordingly.

Overall, the algorithm follows a standard approach. It starts by placing pieces at the center of the board and then evaluates the best possible moves around the existing pieces, assessing them based on their potential outcomes.

### How to Run

To execute the program, run the following command. You need to pass the URL through the command line and execute both files simultaneously:

```sh
go run ..\Intro_to_AI_hw\midterm\algorithm.go ..\Intro_to_AI_hw\midterm\aiMain.go -url="target_url"
```