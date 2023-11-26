package model

import (
	"fmt"
	"sort"

	"github.com/ankar71/todo/helper"
)

type Todo struct {
	Id     int
	Label  string
	IsDone bool
}

type TodoDict struct {
	nextId    int
	data      map[int]*Todo
	doneCount int
}

func DefaultDict() *TodoDict {
	return &TodoDict{
		nextId:    0,
		data:      make(map[int]*Todo, 0),
		doneCount: 0,
	}
}

// methods

func (todo *Todo) String() string {
	return fmt.Sprintf("{Id: %d, Label: %q, IsDone: %t}", todo.Id, todo.Label, todo.IsDone)
}

func (td *TodoDict) List() []*Todo {
	result := make([]*Todo, 0, len(td.data))
	for _, v := range td.data {
		result = append(result, v)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Id < result[j].Id
	})
	return result
}

func (td *TodoDict) FilteredList(isDone bool) []*Todo {
	return helper.Filter(td.List(), func(todo *Todo) bool {
		if todo.IsDone == isDone {
			return true
		} else {
			return false
		}
	})
}

func (td *TodoDict) SetStatus(id int, isDone bool) {
	entry, ok := td.data[id]
	if ok && isDone != entry.IsDone {
		entry.IsDone = isDone
		if isDone {
			td.doneCount += 1
		} else {
			td.doneCount -= 1
		}
	}
}

func (td *TodoDict) SetStatusAll(isDone bool) {
	for k := range td.data {
		td.SetStatus(k, isDone)
	}
}

func (td *TodoDict) RemoveTask(id int) {
	if td.data[id].IsDone {
		td.doneCount -= 1
	}
	delete(td.data, id)
}

func (td *TodoDict) RemoveCompleted() {
	completed := td.FilteredList(true)
	for _, todo := range completed {
		td.RemoveTask(todo.Id)
	}
}

func (td *TodoDict) AddTask(descr string) {
	entry := &Todo{
		Id:     td.nextId,
		Label:  descr,
		IsDone: false,
	}
	td.nextId += 1
	td.data[entry.Id] = entry
}

func (td *TodoDict) AllDone() bool {
	return td.doneCount == len(td.data)
}

func (td *TodoDict) ItemsLeft() int {
	return len(td.data) - td.doneCount
}
