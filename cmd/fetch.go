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

	users, err := db.GetUsers(ctx)
	if err != nil {
		return err
	}

	um := make(map[int]db.User)
	for _, v := range users {
		um[v.Worker] = v
	}

	d.UWs = nil
	d.UWs = make(map[int]UserWorker)

	for _, v := range workers {
		d.UWs[v.Id] = UserWorker{Id: v.Id, W: v, U: um[v.Id]}
	}

	return nil
}

func (d *Data) fetchTasks(ctx context.Context) error {
	tasks, err := db.GetTasks(ctx)
	if err != nil {
		return err
	}

	d.Tasks = nil
	d.Tasks = make(map[int]db.Task)

	for _, v := range tasks {
		d.Tasks[v.Id] = v
	}

	if len(tasks) == 0 {
		d.MTask = 0
	} else {
		d.MTask = tasks[len(tasks)-1].Id
	}

	return nil
}

func (d *Data) fetchCats(ctx context.Context) error {
	cats, err := db.GetCategories(ctx)
	if err != nil {
		return err
	}

	d.Cats = nil
	d.Cats = make(map[int]db.Category)

	for _, v := range cats {
		d.Cats[v.Id] = v
	}

	if len(cats) == 0 {
		d.MCat = 0
	} else {
		d.MCat = cats[len(cats)-1].Id
	}

	return nil
}

func (d *Data) fetchAddrs(ctx context.Context) error {
	addrs, err := db.GetAddresses(ctx)
	if err != nil {
		return err
	}

	d.Addrs = nil
	d.Addrs = make(map[int]db.Address)

	for _, v := range addrs {
		d.Addrs[v.Id] = v
	}

	if len(addrs) == 0 {
		d.MAddr = 0
	} else {
		d.MAddr = addrs[len(addrs)-1].Id
	}

	return nil
}
