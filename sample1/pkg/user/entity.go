package user
/**
定义user的实体模型, 通常对应数据库的表结构
 */
import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Password    string `json:"password,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
	Address     string `json:"address,omitempty"`
	DisplayPic  string `json:"display_pic,omitempty"`
}
