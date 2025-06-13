| Hành động        | Endpoint               | Method   | Ghi chú                             |
| ---------------- | ---------------------- | -------- | ----------------------------------- |
| Tạo task         | `/tasks`               | `POST`   | Gọi `AddTask()`                     |
| Vô hiệu hoá task | `/tasks/:name/disable` | `PATCH`  | Gọi `DisableTaskByName()`           |
| Bật lại task     | `/tasks/:name/enable`  | `PATCH`  | Gọi `EnableTaskByName()` *(nên có)* |
| Xoá task         | `/tasks/:name`         | `DELETE` | Gọi `DeleteTaskByName()`            |
