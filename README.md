# 选课系统后端

## API列表

下面API没有说明权限，只是给了功能性列表。

- 系统初始化
  - [x] 检查是否初始化
  - [x] 初始化
- 公开API
  - [x] 登录
  - [x] 注册
  - [x] 登出
  - [x] 检测可否注册
  - [x] 开放或关闭注册
- 用户相关
  - [x] 获取当前用户信息
  - [x] 获取其他用户信息
  - [x] 添加用户（与注册不一样）
  - [x] 获取用户列表
  - [x] 删除用户
  - [x] 更新用户信息
  - [x] 更新个人信息（非管理员仅限修改用户名）
  - [x] 修改密码
  - [x] 修改他人密码（管理员）
- 课程相关
  - [x] 获取可选课程列表
  - [x] 添加课程共信息
  - [x] 开设课头
  - [x] 更新课程信息
  - [x] 更新课头信息
  - [x] 获取用户课程安排表
  - [ ] 查询成绩
  - [x] 查询课头内学生及相关信息
  - [x] 选课
  - [x] 撤课
  - [x] 强制选课
  - [x] 强制撤课
  - [ ] 登记成绩
- 学院相关
  - [x] 获取学院列表
  - [x] 开设学院
- 学期相关
  - [x] 获取学期列表
  - [x] 开设学期
  - [x] 获取当前学期
  - [x] 设置当前学期