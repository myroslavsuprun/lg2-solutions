package main

import "fmt"

func LLExec() {
	var ll LinkedList[int]

	ll.Add(1)
	ll.Add(2)
	ll.Add(3)
	ll.Add(4)

	fmt.Println(ll)

	ll.Insert(22, 2)
	fmt.Println(ll)

	fmt.Println(ll.Index(22))
}

type LinkedList[T comparable] struct {
	val  T
	next *LinkedList[T]
}

func (ll LinkedList[T]) String() string {
	var next string
	if ll.next != nil {
		next = ll.next.String()
	}

	return fmt.Sprintf("value: %v, next: (%v)", ll.val, next)
}

func (ll *LinkedList[T]) Add(v T) {
	if ll == nil {
		ll = &LinkedList[T]{
			val: v,
		}
		return
	}

	if ll.val == *new(T) {
		ll.val = v
		return
	}

	if ll.next == nil {
		ll.next = &LinkedList[T]{
			val: v,
		}
		return
	}

	ll.next.Add(v)
}

func (ll *LinkedList[T]) Insert(v T, idx int) {
	if idx == 0 {
		ll.val = v
	}

	cur := ll
	for i := 1; i <= idx; i++ {
		if cur.next == nil {
			*cur.next = LinkedList[T]{}
		}

		cur = cur.next

		if i == idx {
			cur.val = v
		}
	}
}

func (ll LinkedList[T]) Index(v T) int {
	cur := ll
	for i := 0; ; i++ {
		if cur.val == v {
			return i
		}

		if cur.next == nil {
			return -1
		}

		cur = *cur.next
	}
}
