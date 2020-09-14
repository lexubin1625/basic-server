package meta

import (
	"time"
)

var (
	_ = time.Second
)

type Logs struct {
	//[ 0] id                                             uint                 null: false  primary: true   isArray: false  auto: true   col: uint            len: -1      default: []
	ID uint32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:uint;" json:"id"`
	//[ 1] func                                           varchar              null: false  primary: false  isArray: false  auto: false  col: varchar         len: -1      default: []
	Func string `gorm:"column:func;type:varchar;" json:"func"`
	//[ 2] action                                         smallint             null: false  primary: false  isArray: false  auto: false  col: smallint        len: -1      default: []
	Action int32 `gorm:"column:action;type:smallint;" json:"action"` // 行为 1 添加 2删除  3 修改
	//[ 3] ext                                            varchar              null: false  primary: false  isArray: false  auto: false  col: varchar         len: -1      default: []
	Ext string `gorm:"column:ext;type:varchar;" json:"ext"` // 扩展字段
	//[ 4] ext_value                                      varchar              null: false  primary: false  isArray: false  auto: false  col: varchar         len: -1      default: []
	ExtValue string `gorm:"column:ext_value;type:varchar;" json:"ext_value"` // 扩展字段值
	//[ 5] logs                                           varchar              null: false  primary: false  isArray: false  auto: false  col: varchar         len: -1      default: []
	Logs string `gorm:"column:logs;type:varchar;" json:"logs"` // 具体行为描述
	//[ 6] created_at                                     timestamp            null: true   primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: []
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;" json:"created_at"`
	//[ 7] created_user                                   varchar              null: false  primary: false  isArray: false  auto: false  col: varchar         len: -1      default: []
	CreatedUser string `gorm:"column:created_user;type:varchar;" json:"created_user"` // 操作人
	//[ 8] status                                         tinyint              null: false  primary: false  isArray: false  auto: false  col: tinyint         len: -1      default: []
	Status int32 `gorm:"column:status;type:tinyint;" json:"status"` // 操作状态 0 失败 1成功

}
