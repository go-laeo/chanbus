package chanbus

type ChanList []chan interface{}

func (cl ChanList) IndexOf(el chan interface{}) int {
	for i, v := range cl {
		if v == el {
			return i
		}
	}

	return -1
}
