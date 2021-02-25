package storage

import (
	"fmt"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"time"
)
// https://github.com/etcd-io/bbolt
// https://github.com/filebrowser/filebrowser/blob/master/storage/bolt/bolt.go
// https://github.com/filebrowser/filebrowser/blob/master/storage/bolt/config.go
// 或许直接用 https://github.com/asdine/storm 更简单？

func initDatabase() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	//db, err := bolt.Open("comigo.db", 0600, nil)
	//创建或读取数据库,并添加超时控制
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	//读写事务，如果是db.View ，则为只读。
	db.Update(func(tx *bolt.Tx) error {
		//建立一个新桶，参数是桶的名字
		b, err := tx.CreateBucketIfNotExists([]byte("bookBucket"))
		if err != nil {
			return fmt.Errorf("create book bucket: %v", err)
		}

		//使用 Bucket.Put() 存储键值对，接收两个 []byte 类型的参数
		err = b.Put([]byte("book"), []byte("test"))
		//使用 Bucket.Get() 来查询键值。参数是一个 []byte.Get 不会返回 error ，Get() 一定能正常工作（除非系统错误）
		v := b.Get([]byte("book"))
		fmt.Printf("The answer is: %s\n", v)

		//根据名字来删除桶
		if err = tx.DeleteBucket([]byte("bookBucket")); err != nil {
			return err
		}

		return nil
	})





}

func DataBase() {
//使用 db.Update 开启一个读写事务
	err := db.Update(func(tx *bolt.Tx) error{

		return nil
	})


}