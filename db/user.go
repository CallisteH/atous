package db

import (
	"atous/model"
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/muyo/sno"
)

func (s *DB) CreateUser(u *model.User) error {
	return s.conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketUsers))

		u.ID = "us_" + sno.New(byte(1)).String()

		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}

		return b.Put([]byte(u.ID), buf)
	})
}

func (s *DB) UpdateUser(id string, u *model.User) error {
	return s.conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketUsers))

		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), buf)
	})
}

func (s *DB) GetUser(id string) (*model.User, error) {
	var u model.User

	err := s.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketUsers))
		log.Println("GetUser id:", id)
		v := b.Get([]byte(id))
		if v == nil {
			return fmt.Errorf("User not found")
		}
		return json.Unmarshal(v, &u)
	})

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *DB) DeleteUser(id string) error {
	return s.conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketUsers))
		return b.Delete([]byte(id))
	})
}

func (s *DB) GetListUsers() ([]*model.User, error) {
	var users []*model.User
	err := s.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketUsers))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var u model.User
			err := json.Unmarshal(v, &u)
			if err != nil {
				return err
			}
			users = append(users, &u)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *DB) GetUserByEmail(email string) (*model.User, error) {
	var u model.User

	err := s.conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketUsers))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &u)
			if err != nil {
				return err
			}
			if u.Email == email {
				return nil
			}
		}
		return fmt.Errorf("User not found")
	})

	if err != nil {
		return nil, err
	}

	return &u, nil
}
