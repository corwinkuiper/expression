# Expression

Some remnants of a project to make something that can perform a least squares fit using some of the knowledge that I had.
Most of this is entirely useless, however I will go through what it can do.

Expression does indeed parse expressions and can produce an output number.
It can only accept functions of one variable and only deals in numbers.
It can also numerically differentiate expressions although this is pretty useless and I'm not entirely sure why I wrote it.

Fit can sometimes fit functions using least squares and steepest decent.
I tested it on some linear functions and it *seemed* to work okay.
However with some more complicated functions it gave quite a bit of nonsense.

The main package can be compiled and it takes an expression and returns some answer.
Just look in `main.go` for the actual details of how this works.

## Why

I'm quite happy with how this worked despite the ultimate failure.
I am happy because I had never written something that could parse a mathematical expression before and that project didn't totally fail.
I did implement known algorithms for manipulating the tokens but I'm still not too bothered by any of it.

As for the fitting, that didn't go quite as well as I had hoped.
I will have to do some more reading around least squares fitting, so I plan on doing that.