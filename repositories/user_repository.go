package repositories

import (
	"errors"
	"sync"

	"github.com/c-my/lottery_iris/datamodels"
)

//Query defines query actions
type Query func(datamodels.User) bool

//UserRepository is a
type UserRepository interface {
	Exec(query Query, action Query, limit int, mode int) (ok bool)
	Select(query Query) (user datamodels.User, found bool)
	SelectMany(query Query, limit int) (results []datamodels.User)

	InsertOrUpdate(user datamodels.User) (updatedUser datamodels.User, err error)
	Delete(query Query, limit int) (deleted bool)
}

//NewUserRepository is
func NewUserRepository(source map[uint]datamodels.User) UserRepository {
	return &userMemoryRepository{source: source}
}

type userMemoryRepository struct {
	source map[uint]datamodels.User
	mu     sync.RWMutex
}

const (
	ReadOnlyMode = iota
	ReadWriteMode
)

func (r *userMemoryRepository) Exec(query Query, action Query, actionLimit int, mode int) (ok bool) {
	loops := 0
	if mode == ReadOnlyMode {
		r.mu.RLock()
		defer r.mu.RUnlock()
	} else {
		r.mu.Lock()
		defer r.mu.Unlock()
	}

	for _, user := range r.source {
		ok = query(user)
		if ok {
			if action(user) {
				loops++
				if actionLimit >= loops {
					break
				}
			}
		}
	}
	return
}

func (r *userMemoryRepository) Select(query Query) (user datamodels.User, found bool) {
	found = r.Exec(query, func(u datamodels.User) bool {
		user = u
		return true
	}, 1, ReadOnlyMode)

	// set an empty datamodels.Movie if not found at all.
	if !found {
		user = datamodels.User{}
	}

	return
}

func (r *userMemoryRepository) SelectMany(query Query, limit int) (results []datamodels.User) {
	r.Exec(query, func(u datamodels.User) bool {
		results = append(results, u)
		return true
	}, limit, ReadOnlyMode)

	return
}

func (r *userMemoryRepository) InsertOrUpdate(user datamodels.User) (datamodels.User, error) {
	id := user.ID

	if id == 0 { // 创建一个新的操作
		var lastID uint
		// 找到最大的ID，避免重复。
		// 在实际使用时您可以使用第三方库去生成
		// 一个string类型的UUID
		r.mu.RLock()
		for _, item := range r.source {
			if item.ID > lastID {
				lastID = item.ID
			}
		}
		r.mu.RUnlock()

		id = lastID + 1
		user.ID = id

		// map-specific thing
		r.mu.Lock()
		r.source[id] = user
		r.mu.Unlock()

		return user, nil
	}

	current, exists := r.Select(func(m datamodels.User) bool {
		return m.ID == id
	})

	if !exists { // 当ID不存在时抛出一个error
		return datamodels.User{}, errors.New("failed to update a nonexistent user")
	}

	// 或者注释下面这段然后用 r.source[id] = m 做单纯替换
	if user.Nickname != "" {
		current.Nickname = user.Nickname
	}

	// map-specific thing
	r.mu.Lock()
	r.source[id] = current
	r.mu.Unlock()

	return user, nil
}

func (r *userMemoryRepository) Delete(query Query, limit int) bool {
	return r.Exec(query, func(u datamodels.User) bool {
		delete(r.source, u.ID)
		return true
	}, limit, ReadWriteMode)
}
