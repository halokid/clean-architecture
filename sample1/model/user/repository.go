package user
/**
实体模型具体实现的方法，也就是定义接口, 这些接口方法都是直接操作数据库的
 */

import (
	"context"
)


type Repository interface {
	/** 定义接口 */
	FindByID(ctx context.Context, id uint) (*User, error)

	BuildProfile(ctx context.Context, user *User) (*User, error)

	CreateMinimal(ctx context.Context, email, password, phoneNumber string) (*User, error)

	FindByEmailAndPassword(ctx context.Context, email, password string) (*User, error)

	FindByEmail(ctx context.Context, email string) (*User, error)

	DoesEmailExist(ctx context.Context, email string) (bool, error)

	ChangePassword(ctx context.Context, email, password string) error
}
