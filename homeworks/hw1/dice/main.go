package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	MaxN = 100
	MaxM = 100
)

var exprRe = regexp.MustCompile(`^\s*(\d+)\s*[dD]\s*(\d+)\s*$`)

func printHelp() {
	fmt.Print(`Dice Roller

Usage:
  dice [flags] NdM

Where:
  N - number of dice (N <= 100)
  M - number of faces per die (M <= 100)

Flags:
  -n, --count <number>   Run multiple iterations (default 1)
  -s, --sum              Print only sum per iteration (incompatible with -v)
  -v, --verbose          Verbose (incompatible with -s)
  -h, --help             Show this help
Examples:
  dice 1d6
  dice 4d20
  dice --sum 4d20
  dice -v 4d20
  dice -n 3 2d10
  dice --count 3 --verbose 2d10
  dice -n 5 -s 3d6
`)
}

func main() {
	var (
		iterations  int
		sumFlag     bool
		verboseFlag bool
		helpFlag    bool
	)

	flag.IntVar(&iterations, "n", 1, "iterations count")
	flag.BoolVar(&sumFlag, "s", false, "sum only (incompatible with -v)")
	flag.BoolVar(&verboseFlag, "v", false, "verbose (incompatible with -s)")
	flag.BoolVar(&helpFlag, "h", false, "help")

	flag.IntVar(&iterations, "count", 1, "iterations count")
	flag.BoolVar(&sumFlag, "sum", false, "sum only (incompatible with -v)")
	flag.BoolVar(&verboseFlag, "verbose", false, "verbose (incompatible with -s)")
	flag.BoolVar(&helpFlag, "help", false, "help")

	flag.Parse()

	if helpFlag {
		printHelp()
		return
	}
	if sumFlag && verboseFlag {
		fmt.Fprintln(os.Stderr, "flags -s/--sum and -v/--verbose are incompatible")
		os.Exit(1)
	}
	if iterations <= 0 {
		fmt.Fprintln(os.Stderr, "count must be > 0")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: dice [flags] NdM (use -h for help)")
		os.Exit(1)
	}

	spec, err := parseSpec(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if verboseFlag {
		if iterations == 1 {
			fmt.Printf("Rolling %dd%d:\n", spec.N, spec.M)
		} else {
			fmt.Printf("Rolling %dd%d (%d iterations):\n", spec.N, spec.M, iterations)
		}
	}

	for i := 1; i <= iterations; i++ {
		rolls, sum := rollOnce(r, spec.N, spec.M)

		if verboseFlag {
			fmt.Printf("Iteration %d:\n", i)
			for j, v := range rolls {
				fmt.Printf("  Dice %d: %d\n", j+1, v)
			}
			fmt.Printf("  Sum: %d\n", sum)
			fmt.Printf("  Average: %g\n", float64(sum)/float64(spec.N))
			continue
		}

		fmt.Printf("Iteration %d:\n", i)
		if sumFlag {
			fmt.Printf("  Sum: %d\n", sum)
		} else {
			fmt.Printf("  Rolls: %s\n", formatRolls(rolls, sum))
		}
	}
}

type Spec struct {
	N int
	M int
}

func parseSpec(expr string) (Spec, error) {
	m := exprRe.FindStringSubmatch(expr)
	if m == nil {
		return Spec{}, fmt.Errorf("invalid dice expression: %q", expr)
	}
	n, _ := strconv.Atoi(m[1])
	mm, _ := strconv.Atoi(m[2])
	if n <= 0 || mm <= 0 {
		return Spec{}, fmt.Errorf("n and m must be > 0")
	}
	if n > MaxN || mm > MaxM {
		return Spec{}, fmt.Errorf("limits exceeded: N<=%d, M<=%d", MaxN, MaxM)
	}
	return Spec{N: n, M: mm}, nil
}

func rollOnce(r *rand.Rand, n, m int) ([]int, int) {
	res := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		val := r.Intn(m) + 1
		res[i] = val
		sum += val
	}
	return res, sum
}

func formatRolls(vals []int, sum int) string {
	if len(vals) == 1 {
		return fmt.Sprintf("%d", vals[0])
	}
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = strconv.Itoa(v)
	}
	return fmt.Sprintf("%s=%d", strings.Join(parts, "+"), sum)
}
