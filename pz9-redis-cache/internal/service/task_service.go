package service

import (
	"context"
	"errors"
	"log"
	"time"

	"example.com/pz9-redis-cache/internal/cache"
	"example.com/pz9-redis-cache/internal/config"
	"example.com/pz9-redis-cache/internal/task"
	"github.com/redis/go-redis/v9"
)

type TaskService struct {
	repo       *task.Repo
	redis      *redis.Client
	cfg        config.Config
	keys       cache.KeyBuilder
	serializer cache.Serializer
}

func NewTaskService(repo *task.Repo, redisClient *redis.Client, cfg config.Config) *TaskService {
	return &TaskService{
		repo:       repo,
		redis:      redisClient,
		cfg:        cfg,
		keys:       cache.NewKeyBuilder(),
		serializer: cache.NewSerializer(),
	}
}

func (s *TaskService) GetTaskByID(ctx context.Context, id int64) (task.Task, error) {
	key := s.keys.TaskByID(id)

	if s.redis != nil {
		cached, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			if s.serializer.IsNotFound([]byte(cached)) {
				log.Println("cache hit:", key, "(negative)")
				return task.Task{}, task.ErrTaskNotFound
			}
			t, err := s.serializer.UnmarshalTask([]byte(cached))
			if err == nil {
				log.Println("cache hit:", key)
				return t, nil
			}
			log.Println("cache decode error:", err)
		} else if !errors.Is(err, redis.Nil) {
			log.Println("redis read error:", err)
		} else {
			log.Println("cache miss:", key)
		}
	}

	t, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotFound) {
			s.cacheNotFound(ctx, key)
		}
		return task.Task{}, err
	}

	s.cacheTask(ctx, key, t, s.cfg.CacheTTL, s.cfg.CacheTTLJitter)
	return t, nil
}

func (s *TaskService) ListTasks(ctx context.Context, page, limit int) ([]task.Task, error) {
	key := s.keys.TasksListPaged(page, limit)

	if s.redis != nil {
		cached, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			items, err := s.serializer.UnmarshalTaskList([]byte(cached))
			if err == nil {
				log.Println("cache hit:", key)
				return items, nil
			}
			log.Println("cache decode error:", err)
		} else if !errors.Is(err, redis.Nil) {
			log.Println("redis read error:", err)
		} else {
			log.Println("cache miss:", key)
		}
	}

	items := s.repo.List(page, limit)
	s.cacheTaskList(ctx, key, items, s.cfg.ListCacheTTL, s.cfg.CacheTTLJitter)
	return items, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, t task.Task) error {
	if err := s.repo.Update(t); err != nil {
		return err
	}
	s.invalidateTask(ctx, t.ID)
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.invalidateTask(ctx, id)
	return nil
}

func (s *TaskService) cacheTask(ctx context.Context, key string, t task.Task, base, jitter time.Duration) {
	if s.redis == nil {
		return
	}
	bytes, err := s.serializer.MarshalTask(t)
	if err != nil {
		log.Println("cache encode error:", err)
		return
	}
	ttl := cache.TTLWithJitter(base, jitter)
	if err := s.redis.Set(ctx, key, bytes, ttl).Err(); err != nil {
		log.Println("redis write error:", err)
		return
	}
	log.Println("cache set:", key, "ttl=", ttl)
}

func (s *TaskService) cacheTaskList(ctx context.Context, key string, items []task.Task, base, jitter time.Duration) {
	if s.redis == nil {
		return
	}
	bytes, err := s.serializer.MarshalTaskList(items)
	if err != nil {
		log.Println("cache encode error:", err)
		return
	}
	ttl := cache.TTLWithJitter(base, jitter)
	if err := s.redis.Set(ctx, key, bytes, ttl).Err(); err != nil {
		log.Println("redis write error:", err)
		return
	}
	log.Println("cache set:", key, "ttl=", ttl)
}

func (s *TaskService) cacheNotFound(ctx context.Context, key string) {
	if s.redis == nil {
		return
	}
	ttl := cache.TTLWithJitter(s.cfg.NegativeCacheTTL, s.cfg.CacheTTLJitter/2)
	if err := s.redis.Set(ctx, key, s.serializer.MarshalNotFound(), ttl).Err(); err != nil {
		log.Println("redis write error:", err)
		return
	}
	log.Println("cache set:", key, "(negative) ttl=", ttl)
}

func (s *TaskService) invalidateTask(ctx context.Context, id int64) {
	if s.redis == nil {
		return
	}
	key := s.keys.TaskByID(id)
	if err := s.redis.Del(ctx, key).Err(); err != nil {
		log.Println("redis delete error:", err)
		return
	}
	log.Println("cache invalidated:", key)

	listKey := s.keys.TasksList()
	if err := s.redis.Del(ctx, listKey).Err(); err != nil {
		log.Println("redis delete error:", err)
	}
	log.Println("cache invalidated:", listKey)

	// Remove paged list keys (small dataset for practice).
	for page := 1; page <= 5; page++ {
		for limit := 1; limit <= 20; limit++ {
			pagedKey := s.keys.TasksListPaged(page, limit)
			_ = s.redis.Del(ctx, pagedKey).Err()
		}
	}
}
