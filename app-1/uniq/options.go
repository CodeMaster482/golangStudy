package uniq

import "flag"

type Options struct {
	count  bool
	repeat bool
	unique bool
	ignore bool
	fields int
	chars  int
}

func Init(options *Options) {
	flag.BoolVar(&options.count, "c", false, "num of repetition of line")
	flag.BoolVar(&options.repeat, "d", false, "lines dublicated")
	flag.BoolVar(&options.unique, "u", false, "unique lines")
	flag.BoolVar(&options.ignore, "i", false, "symbols case")
	flag.IntVar(&options.fields, "f", 0, "dont count first n fields")
	flag.IntVar(&options.chars, "s", 0, "dont count first n chars")
}
