package worker

import (
	"github.com/dgraph-io/dgraph/x"
	"strconv"
)

const INRANGE = "ir"

func inRange(bytes []byte, ranges []string) (bool, error) {

	if len(ranges) != 2 {
		return false, x.Errorf("Two argument expected in inRange, but got %d.", len(ranges))
	}

	start, err := strconv.Atoi(ranges[0])
	if err != nil {
		return false, err
	}

	end, err := strconv.Atoi(ranges[1])
	if err != nil {
		return false, err
	}

	isLeading := false
	hasRange := false
	position := 0
	value := 0
	var shift uint32 = 0

	readInt := func() {
		shift = 0
		if isLeading { //leading数字首个byte有两个标志位
			shift = 6
			b := bytes[position]
			position++
			value = int(b & 0x3F)
			hasRange = (b & 0x80) != 0 //hasRange在首位bit
			if (b & 0x40) == 0 {
				return
			}
		} else {
			value = 0
		}

		for ; shift < 32; shift += 7 {
			b := bytes[position]
			position++
			value |= int(b&0x7F) << shift
			if (b & 0x80) == 0 {
				return
			}
		}
	}

	current := -2

	for position < len(bytes) {
		isLeading = true
		readInt()
		current += value + 2
		if current >= start && current <= end {
			return true, nil
		} else if current > end {
			return false, nil
		} else if hasRange {
			isLeading = false
			readInt()
			current += value + 1
			if current >= start {
				return true, nil
			}
		}
	}

	return false, nil

}
