package main

import (
	"argentum/db"
	"context"
	"fmt"
	"time"
)

func autoFetch(d *Data, interval, timeout time.Duration) {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), timeout*2)
			if err := d.fetchUWs(ctx); err != nil {
				fmt.Print("UW\t")
				fmt.Println(err)
			}
			cancel()
			time.Sleep(interval)
		}
	}()

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			if err := d.fetchTasks(ctx); err != nil {
				fmt.Print("tasks\t")
				fmt.Println(err)
			}
			cancel()
			time.Sleep(interval)
		}
	}()

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			if err := d.fetchCats(ctx); err != nil {
				fmt.Print("cats\t")
				fmt.Println(err)
			}
			cancel()
			time.Sleep(interval)
		}
	}()

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			if err := d.fetchAddrs(ctx); err != nil {
				fmt.Print("addrs\t")
				fmt.Println(err)
			}
			cancel()
			time.Sleep(interval)
		}
	}()

}

func (d *Data) fetchUWs(ctx context.Context) error {
	workers, err := db.GetWorkers(ctx)
	if err != nil {
		return err
	}

	if len(workers) == 0 {
		d.MWorker = 0
	} else {
		d.MWorker = workers[len(workers)-1].Id
	}

	users, err := db.GetUsers(ctx)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		d.MUser = 0
	} else {
		d.MUser = users[len(users)-1].Id
	}

	um := make(map[int]db.User)
	for _, u := range users {
		um[u.Worker] = u
	}

	temp := make(map[int]*UserWorker)

	for _, w := range workers {
		temp[w.Id] = &UserWorker{Id: w.Id, W: w, U: um[w.Id]}
	}

	if len(temp) != 0 {
		d.UWs = temp
	}

	return nil
}

func (d *Data) fetchTasks(ctx context.Context) error {
	tasks, err := db.GetTasks(ctx)
	if err != nil {
		return err
	}

	temp := make(map[int]*db.Task)

	for _, v := range tasks {
		temp[v.Id] = &v
	}

	if len(tasks) == 0 {
		d.MTask = 0
	} else {
		d.Tasks = temp
		d.MTask = tasks[len(tasks)-1].Id
	}

	return nil
}

func (d *Data) fetchCats(ctx context.Context) error {
	cats, err := db.GetCategories(ctx)
	if err != nil {
		return err
	}

	temp := make(map[int]*db.Category)

	for _, v := range cats {
		temp[v.Id] = &v
	}

	if len(cats) == 0 {
		d.MCat = 0
	} else {
		d.Cats = temp
		d.MCat = cats[len(cats)-1].Id
	}

	return nil
}

func (d *Data) fetchAddrs(ctx context.Context) error {
	addrs, err := db.GetAddresses(ctx)
	if err != nil {
		return err
	}

	temp := make(map[int]db.Address)

	for _, v := range addrs {
		temp[v.Id] = v
	}

	if len(addrs) == 0 {
		d.MAddr = 0
	} else {
		d.Addrs = temp
		d.MAddr = addrs[len(addrs)-1].Id
	}

	return nil
}
