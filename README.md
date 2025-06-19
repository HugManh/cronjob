myapp/
├── cmd/
│   └── myapp/            # entrypoint chính
│       └── main.go       # khởi tạo server, router, và engine template
├── internal/             # logic/handler chỉ dùng nội bộ
│   ├── handler/
│   └── service/
├── pkg/                  # thư viện có thể public (nếu cần)
├── views/                # chứa file .jet
│   ├── layouts/
│   │    └── main.jet     # layout chung với {{yield}} hoặc {{ embed() }}
│   ├── partials/
│   │    ├── header.jet
│   │    └── footer.jet
│   └── pages/
│        ├── index.jet
│        └── about.jet
├── web/ (tuỳ chọn)       # assets static: css/js/img
├── go.mod
└── go.sum


| Hành động        | Endpoint               | Method   | Ghi chú                             |
| ---------------- | ---------------------- | -------- | ----------------------------------- |
| Tạo task         | `/tasks`               | `POST`   | Gọi `AddTask()`                     |
| Vô hiệu hoá task | `/tasks/:name/disable` | `PATCH`  | Gọi `DisableTaskByName()`           |
| Bật lại task     | `/tasks/:name/enable`  | `PATCH`  | Gọi `EnableTaskByName()` *(nên có)* |
| Xoá task         | `/tasks/:name`         | `DELETE` | Gọi `DeleteTaskByName()`            |
