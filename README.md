is a machine learning tool based on [genetic programming](https://www.wikiwand.com/en/Genetic_programming) which can be used for classification and regression problems.

## To do

- Create lookup tables for operators that time to evaluate at runtime
- http://www.genetic-programming.com/gpquadraticexample.html
- http://www.genetic-programming.com/gpflowchart.html
- Organize code after reading [this](https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project)
- https://github.com/gonum/blas
- [Classification Strategies for Image Classification in Genetic Programming](http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.475.3010&rep=rep1&type=pdf)
- [Multi-class overview](http://dynamics.org/~altenber/UH_ICS/EC_REFS/GP_REFS/IEEE/CEC2001/395.pdf)
- [Sampling data for fitness evaluation](http://eplex.cs.ucf.edu/papers/morse_gecco16.pdf)
- Integrate [golearn](https://github.com/sjwhitworth/golearn)

## Architecture

![architecture](https://docs.google.com/drawings/d/1en_XKo3L65RCiFtu2ftutXYpPE3DO7SBW3qLL36Rdg4/pub?w=960&h=720)

## Development

```sh
go run cmd/xgp/*.go fit examples/iris/train.csv
```

## Dependencies

- [MaxHalford/gago](https://github.com/MaxHalford/gago)
- [urfave/cli](https://github.com/urfave/cli)